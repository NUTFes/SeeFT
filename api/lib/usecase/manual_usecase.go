package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

// マニュアル閲覧用セッションCookieの有効期限（issue #444: 30日間ログイン状態を維持する）
const manualSessionTTL = 30 * 24 * time.Hour

// OAuth state の有効期限（issue #444: ログイン開始から10分以内にコールバックを完了させる）
const manualStateTTL = 10 * time.Minute

// Googleのトークンエンドポイント。資格情報ではなく公開URL
const googleTokenEndpoint = "https://oauth2.googleapis.com/token" //nolint:gosec // 公開URLでありハードコードされた資格情報ではない

// IDトークンの発行者として受理する値（Google公式ドキュメントに両表記が明記されている）
var manualAllowedIssuers = []string{"https://accounts.google.com", "accounts.google.com"}

// 例外的に許可するメール。小文字で記載
var manualAllowExtra = []string{}

// nutfesドメインの許可判定用（大文字混じりを許容するため case-insensitive）
var manualAllowedDomainRe = regexp.MustCompile(`(?i)\.nutfes@gmail\.com$`)

// マニュアルIDのバリデーション（パストラバーサル対策）
var manualIDRe = regexp.MustCompile(`^[a-z0-9_-]{1,64}$`)

// ManualConfig は ManualUseCase の生成に必要な設定値。
// 環境変数の読み取りは呼び出し側（di.go）の責務とし、usecase では os.Getenv を呼ばない。
type ManualConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ManualDir    string
}

// ManualUseCase はGoogleアカウントによるマニュアル閲覧制限（issue #444）を扱う。
type ManualUseCase interface {
	// BuildAuthURL は指定マニュアルへ戻るstate付きのGoogle認可URLを組み立てる。
	BuildAuthURL(manualID string) string
	// MakeState はCSRF対策とリダイレクト先の保持を兼ねたOAuth stateを発行する。
	MakeState(manualID string) string
	// VerifyState はstateの署名と有効期限を検証し、埋め込まれたマニュアルIDを返す。
	VerifyState(state string) (manualID string, ok bool)
	// ExchangeCode は認可コードをGoogleのトークンエンドポイントに渡し、認証済みメールアドレスを取得する。
	ExchangeCode(ctx context.Context, code string) (email string, err error)
	// IsAllowed はメールアドレスが閲覧許可対象（nutfesドメイン or 例外許可）かどうかを判定する。
	IsAllowed(email string) bool
	// MakeSessionCookie は閲覧許可済みメールアドレスからセッションCookie値を発行する。
	MakeSessionCookie(email string) string
	// VerifySessionCookie はセッションCookie値の署名・期限・現在の許可リストを検証する。
	VerifySessionCookie(value string) (email string, ok bool)
	// ManualPath はマニュアルIDからHTMLファイルパスを解決する（存在しなければ ok=false）。
	ManualPath(manualID string) (string, bool)
}

type manualUseCase struct {
	cfg ManualConfig
	// テストからGoogleのトークンエンドポイントを差し替えられるようにフィールドにしている
	tokenEndpoint string
}

// NewManualUseCase は ManualUseCase を生成する。
func NewManualUseCase(cfg ManualConfig) ManualUseCase {
	return &manualUseCase{
		cfg:           cfg,
		tokenEndpoint: googleTokenEndpoint,
	}
}

// BuildAuthURL は指定マニュアルへ戻るstate付きのGoogle認可URLを組み立てる。
func (u *manualUseCase) BuildAuthURL(manualID string) string {
	v := url.Values{}
	v.Set("client_id", u.cfg.ClientID)
	v.Set("redirect_uri", u.cfg.RedirectURL)
	v.Set("response_type", "code")
	v.Set("scope", "openid email")
	v.Set("state", u.MakeState(manualID))
	return "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
}

type manualStatePayload struct {
	ID string `json:"id"`
	TS int64  `json:"ts"`
}

// MakeState はCSRF対策とリダイレクト先の保持を兼ねたOAuth stateを発行する。
// payload と署名を "." で連結した文字列を返す（Apps Script版と同じ構成）。
func (u *manualUseCase) MakeState(manualID string) string {
	payload := manualStatePayload{ID: manualID, TS: time.Now().UnixMilli()}
	return u.signPayload(payload)
}

// VerifyState はstateの署名と有効期限を検証し、埋め込まれたマニュアルIDを返す。
// 署名が正当なら期限切れの場合でもmanualIDを返す（呼び出し側でリトライ導線に使うため）。
// 署名が不正／JSONが壊れている場合はmanualIDを信用できないため空文字を返す。
func (u *manualUseCase) VerifyState(state string) (string, bool) {
	raw, ok := u.verifySignedPayload(state)
	if !ok {
		return "", false
	}
	var p manualStatePayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", false
	}
	if p.ID == "" {
		return "", false
	}
	if time.Since(time.UnixMilli(p.TS)) > manualStateTTL {
		return p.ID, false
	}
	return p.ID, true
}

type manualTokenResponse struct {
	IDToken string `json:"id_token"`
}

type manualIDTokenClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Aud           string `json:"aud"`
	Iss           string `json:"iss"`
}

// ExchangeCode は認可コードをGoogleのトークンエンドポイントに渡し、認証済みメールアドレスを取得する。
func (u *manualUseCase) ExchangeCode(ctx context.Context, code string) (string, error) {
	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", u.cfg.ClientID)
	form.Set("client_secret", u.cfg.ClientSecret)
	form.Set("redirect_uri", u.cfg.RedirectURL)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("トークンリクエスト作成失敗: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("トークンエンドポイントへの通信失敗: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("トークンエンドポイントが非OKステータスを返却: %d", resp.StatusCode)
	}

	var tokenResp manualTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("トークンレスポンスのJSON解析失敗: %w", err)
	}
	if tokenResp.IDToken == "" {
		return "", fmt.Errorf("id_tokenがレスポンスに含まれていない")
	}

	// 署名検証はGoogleのトークンエンドポイントからTLS直取得のため不要
	parts := strings.Split(tokenResp.IDToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("id_tokenの形式が不正")
	}
	payloadJSON, err := decodeBase64URLSegment(parts[1])
	if err != nil {
		return "", fmt.Errorf("id_tokenのpayloadデコード失敗: %w", err)
	}
	var claims manualIDTokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return "", fmt.Errorf("id_tokenのpayload解析失敗: %w", err)
	}
	// TLS直取得のため署名検証は省略するが、他アプリ向けトークンの誤用・混入への
	// 防御としてaud（宛先クライアント）とiss（発行者）は検証する
	if claims.Aud != u.cfg.ClientID {
		return "", fmt.Errorf("id_tokenのaudが本アプリのクライアントIDと一致しない")
	}
	if !slices.Contains(manualAllowedIssuers, claims.Iss) {
		return "", fmt.Errorf("id_tokenのissがGoogleの発行者と一致しない")
	}
	if !claims.EmailVerified {
		return "", fmt.Errorf("メールアドレスが未検証のためログイン拒否")
	}
	email := strings.ToLower(strings.TrimSpace(claims.Email))
	if email == "" {
		return "", fmt.Errorf("id_tokenにメールアドレスが含まれていない")
	}
	return email, nil
}

// IsAllowed はメールアドレスが閲覧許可対象（nutfesドメイン or 例外許可）かどうかを判定する。
func (u *manualUseCase) IsAllowed(email string) bool {
	if email == "" {
		return false
	}
	if manualAllowedDomainRe.MatchString(email) {
		return true
	}
	return slices.Contains(manualAllowExtra, strings.ToLower(email))
}

type manualSessionPayload struct {
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
}

// MakeSessionCookie は閲覧許可済みメールアドレスからセッションCookie値を発行する。
func (u *manualUseCase) MakeSessionCookie(email string) string {
	payload := manualSessionPayload{
		Email: email,
		Exp:   time.Now().Add(manualSessionTTL).UnixMilli(),
	}
	return u.signPayload(payload)
}

// VerifySessionCookie はセッションCookie値の署名・期限・現在の許可リストを検証する。
// 許可リストからの削除を即座に反映させるため、署名検証後にIsAllowedを再チェックする。
func (u *manualUseCase) VerifySessionCookie(value string) (string, bool) {
	raw, ok := u.verifySignedPayload(value)
	if !ok {
		return "", false
	}
	var p manualSessionPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", false
	}
	if time.Now().UnixMilli() > p.Exp {
		return "", false
	}
	if !u.IsAllowed(p.Email) {
		return "", false
	}
	return p.Email, true
}

// ManualPath はマニュアルIDからHTMLファイルパスを解決する（存在しなければ ok=false）。
func (u *manualUseCase) ManualPath(manualID string) (string, bool) {
	if !manualIDRe.MatchString(manualID) {
		return "", false
	}
	path := filepath.Join(u.cfg.ManualDir, manualID+".html")
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return "", false
	}
	return path, true
}

// signPayload はJSONペイロードをbase64url化した上でHMAC-SHA256署名を付与し、
// "payload.signature" 形式の文字列を返す（state・セッションCookie共通の署名フォーマット）。
func (u *manualUseCase) signPayload(payload any) string {
	// 空鍵で署名するとCookie偽造による認証バイパスにつながるため常に失敗させる
	// （設定の存在検証はdi.go側でも行うが、多層防御としてここでも拒否する）
	if u.cfg.ClientSecret == "" {
		return ""
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		// 呼び出し元の構造体は固定なのでMarshalが失敗することは想定しない
		return ""
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(raw)
	sig := u.sign(encodedPayload)
	return encodedPayload + "." + sig
}

// verifySignedPayload は "payload.signature" 形式の文字列を検証し、デコード済みJSONを返す。
func (u *manualUseCase) verifySignedPayload(value string) ([]byte, bool) {
	// 空鍵は署名の意味を持たないため検証を常に失敗させる（signPayloadと対の多層防御）
	if u.cfg.ClientSecret == "" {
		return nil, false
	}
	idx := strings.LastIndex(value, ".")
	if idx < 0 {
		return nil, false
	}
	encodedPayload := value[:idx]
	sig := value[idx+1:]

	expectedSig := u.sign(encodedPayload)
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil, false
	}

	raw, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, false
	}
	return raw, true
}

// sign はClientSecretを鍵としてbase64url payloadのHMAC-SHA256を計算する。
func (u *manualUseCase) sign(encodedPayload string) string {
	mac := hmac.New(sha256.New, []byte(u.cfg.ClientSecret))
	mac.Write([]byte(encodedPayload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// decodeBase64URLSegment はJWTの各セグメント（パディングなしbase64url）をデコードする。
// パディングが省略されているため、デコード前に不足分を補う。
func decodeBase64URLSegment(seg string) ([]byte, error) {
	if m := len(seg) % 4; m != 0 {
		seg += strings.Repeat("=", 4-m)
	}
	return base64.URLEncoding.DecodeString(seg)
}
