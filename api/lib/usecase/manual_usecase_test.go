package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestManualUseCase() *manualUseCase {
	return &manualUseCase{ //nolint:gosec // ClientSecretはテスト用のダミー値であり実資格情報ではない
		cfg: ManualConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RedirectURL:  "https://example.com/manuals/oauth/callback",
			ManualDir:    "",
		},
		tokenEndpoint: "https://oauth2.googleapis.com/token",
	}
}

// IsAllowed のテーブルテスト。nutfesドメイン（大文字混じり含む）とallowExtraのみ許可され、
// それ以外の大学アドレスや空文字は拒否されることを確認する。
func TestManualUseCase_IsAllowed(t *testing.T) {
	u := newTestManualUseCase()

	// テスト用にallowExtraへ例外アドレスを一時追加する
	originalAllowExtra := manualAllowExtra
	manualAllowExtra = []string{"vip@example.com"}
	defer func() { manualAllowExtra = originalAllowExtra }()

	cases := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "nutfesアドレスは許可",
			email: "taro.nutfes@gmail.com",
			want:  true,
		},
		{
			name:  "nutfesアドレスの大文字混じりも許可（大小無視）",
			email: "Taro.NUTFES@Gmail.com",
			want:  true,
		},
		{
			name:  "大学アドレスは拒否",
			email: "taro@stn.nagaokaut.ac.jp",
			want:  false,
		},
		{
			name:  "空文字は拒否",
			email: "",
			want:  false,
		},
		{
			name:  "allowExtraに含まれるアドレスは許可",
			email: "vip@example.com",
			want:  true,
		},
		{
			name:  "allowExtraに含まれないgmailアドレスは拒否",
			email: "someone@gmail.com",
			want:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := u.IsAllowed(c.email)
			if got != c.want {
				t.Errorf("IsAllowed(%q) = %v, want %v", c.email, got, c.want)
			}
		})
	}
}

// MakeState / VerifyState のラウンドトリップと、改ざん・期限切れ時の拒否を確認する。
func TestManualUseCase_State(t *testing.T) {
	u := newTestManualUseCase()

	t.Run("発行したstateはそのまま検証を通りマニュアルIDを復元できる", func(t *testing.T) {
		state := u.MakeState("summer-fes")
		id, ok := u.VerifyState(state)
		if !ok {
			t.Fatalf("VerifyState() ok = false, want true")
		}
		if id != "summer-fes" {
			t.Errorf("VerifyState() id = %q, want %q", id, "summer-fes")
		}
	})

	t.Run("署名を改ざんしたstateは拒否される", func(t *testing.T) {
		state := u.MakeState("summer-fes")
		tampered := state + "x"
		if _, ok := u.VerifyState(tampered); ok {
			t.Errorf("VerifyState(tampered) ok = true, want false")
		}
	})

	t.Run("11分前に発行されたstateは期限切れとして拒否される", func(t *testing.T) {
		expired := u.signPayload(manualStatePayload{
			ID: "summer-fes",
			TS: time.Now().Add(-11 * time.Minute).UnixMilli(),
		})
		id, ok := u.VerifyState(expired)
		if ok {
			t.Errorf("VerifyState(expired) ok = true, want false")
		}
		// 署名自体は正当なので、リトライ導線のためIDは復元できる想定
		if id != "summer-fes" {
			t.Errorf("VerifyState(expired) id = %q, want %q", id, "summer-fes")
		}
	})
}

// MakeSessionCookie / VerifySessionCookie のラウンドトリップと改ざん検知を確認する。
func TestManualUseCase_SessionCookie(t *testing.T) {
	u := newTestManualUseCase()

	t.Run("発行したセッションCookieはそのまま検証を通りメールアドレスを復元できる", func(t *testing.T) {
		cookie := u.MakeSessionCookie("taro.nutfes@gmail.com")
		email, ok := u.VerifySessionCookie(cookie)
		if !ok {
			t.Fatalf("VerifySessionCookie() ok = false, want true")
		}
		if email != "taro.nutfes@gmail.com" {
			t.Errorf("VerifySessionCookie() email = %q, want %q", email, "taro.nutfes@gmail.com")
		}
	})

	t.Run("改ざんしたセッションCookieは拒否される", func(t *testing.T) {
		cookie := u.MakeSessionCookie("taro.nutfes@gmail.com")
		tampered := cookie + "x"
		if _, ok := u.VerifySessionCookie(tampered); ok {
			t.Errorf("VerifySessionCookie(tampered) ok = true, want false")
		}
	})

	t.Run("許可リストから外れたメールアドレスのCookieは拒否される", func(t *testing.T) {
		// allowExtraにもnutfesドメインにも該当しないメールアドレスで発行したCookieは、
		// 署名が正当でもIsAllowedの再チェックで拒否される
		cookie := u.MakeSessionCookie("taro@stn.nagaokaut.ac.jp")
		if _, ok := u.VerifySessionCookie(cookie); ok {
			t.Errorf("VerifySessionCookie(disallowed email) ok = true, want false")
		}
	})
}

// ManualPath のバリデーションとファイル存在確認を確認する。
func TestManualUseCase_ManualPath(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "summer-fes.html"), []byte("<html></html>"), 0o600); err != nil {
		t.Fatalf("failed to create dummy manual file: %v", err)
	}

	u := &manualUseCase{cfg: ManualConfig{ManualDir: dir}, tokenEndpoint: googleTokenEndpoint}

	t.Run("正常なIDは存在するファイルパスを返す", func(t *testing.T) {
		path, ok := u.ManualPath("summer-fes")
		if !ok {
			t.Fatalf("ManualPath() ok = false, want true")
		}
		want := filepath.Join(dir, "summer-fes.html")
		if path != want {
			t.Errorf("ManualPath() path = %q, want %q", path, want)
		}
	})

	t.Run("パストラバーサルを含むIDは拒否される", func(t *testing.T) {
		if _, ok := u.ManualPath("../etc/passwd"); ok {
			t.Errorf("ManualPath(traversal) ok = true, want false")
		}
	})

	t.Run("大文字を含むIDは拒否される", func(t *testing.T) {
		if _, ok := u.ManualPath("Summer-Fes"); ok {
			t.Errorf("ManualPath(uppercase) ok = true, want false")
		}
	})

	t.Run("存在しないファイルのIDはfalseを返す", func(t *testing.T) {
		if _, ok := u.ManualPath("no-such-manual"); ok {
			t.Errorf("ManualPath(missing) ok = true, want false")
		}
	})
}

// fakeIDToken は署名検証をしない前提で、ヘッダ・ペイロード・ダミー署名からなる
// 未署名相当のJWT文字列を組み立てる（ExchangeCodeはTLS越しの直接取得を前提に署名検証を省略している）。
func fakeIDToken(t *testing.T, email string, emailVerified bool) string {
	t.Helper()
	return fakeIDTokenWithClaims(t, map[string]any{
		"email":          email,
		"email_verified": emailVerified,
		"aud":            "id",
		"iss":            "https://accounts.google.com",
	})
}

// fakeIDTokenWithClaims は任意のclaimsからJWT文字列を組み立てる（aud/iss不一致の検証用）。
func fakeIDTokenWithClaims(t *testing.T, claims map[string]any) string {
	t.Helper()

	header := map[string]string{"alg": "RS256", "typ": "JWT"}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("failed to marshal claims: %v", err)
	}

	encode := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
	return encode(headerJSON) + "." + encode(claimsJSON) + ".dummy-signature"
}

// ExchangeCode がトークンエンドポイントのレスポンスからメールアドレスを取り出せること、
// および email_verified=false のレスポンスは拒否されることを確認する。
func TestManualUseCase_ExchangeCode(t *testing.T) {
	t.Run("email_verified=trueのレスポンスからメールアドレスを取得できる", func(t *testing.T) {
		idToken := fakeIDToken(t, "Taro.Nutfes@Gmail.com", true)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"id_token": idToken})
		}))
		defer server.Close()

		u := &manualUseCase{cfg: ManualConfig{ClientID: "id", ClientSecret: "secret", RedirectURL: "https://example.com/cb"}, tokenEndpoint: server.URL}
		email, err := u.ExchangeCode(context.Background(), "dummy-code")
		if err != nil {
			t.Fatalf("ExchangeCode() error = %v, want nil", err)
		}
		if email != "taro.nutfes@gmail.com" {
			t.Errorf("ExchangeCode() email = %q, want %q（小文字化・trim想定）", email, "taro.nutfes@gmail.com")
		}
	})

	t.Run("email_verified=falseのレスポンスは拒否される", func(t *testing.T) {
		idToken := fakeIDToken(t, "taro.nutfes@gmail.com", false)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"id_token": idToken})
		}))
		defer server.Close()

		u := &manualUseCase{cfg: ManualConfig{ClientID: "id", ClientSecret: "secret", RedirectURL: "https://example.com/cb"}, tokenEndpoint: server.URL}
		email, err := u.ExchangeCode(context.Background(), "dummy-code")
		if err == nil {
			t.Fatalf("ExchangeCode() error = nil, want error（email_verified=falseは拒否想定）")
		}
		if email != "" {
			t.Errorf("ExchangeCode() email = %q, want empty on error", email)
		}
	})

	t.Run("audが自アプリのクライアントIDと異なるトークンは拒否される", func(t *testing.T) {
		idToken := fakeIDTokenWithClaims(t, map[string]any{
			"email": "taro.nutfes@gmail.com", "email_verified": true,
			"aud": "other-app", "iss": "https://accounts.google.com",
		})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"id_token": idToken})
		}))
		defer server.Close()

		u := &manualUseCase{cfg: ManualConfig{ClientID: "id", ClientSecret: "secret", RedirectURL: "https://example.com/cb"}, tokenEndpoint: server.URL}
		if _, err := u.ExchangeCode(context.Background(), "dummy-code"); err == nil {
			t.Fatalf("ExchangeCode() error = nil, want error（aud不一致は拒否想定）")
		}
	})

	t.Run("issがGoogleの発行者でないトークンは拒否される", func(t *testing.T) {
		idToken := fakeIDTokenWithClaims(t, map[string]any{
			"email": "taro.nutfes@gmail.com", "email_verified": true,
			"aud": "id", "iss": "https://evil.example.com",
		})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"id_token": idToken})
		}))
		defer server.Close()

		u := &manualUseCase{cfg: ManualConfig{ClientID: "id", ClientSecret: "secret", RedirectURL: "https://example.com/cb"}, tokenEndpoint: server.URL}
		if _, err := u.ExchangeCode(context.Background(), "dummy-code"); err == nil {
			t.Fatalf("ExchangeCode() error = nil, want error（iss不一致は拒否想定）")
		}
	})
}

// 署名鍵（ClientSecret）が空の場合に署名系が必ず失敗することを確認する。
// 空鍵のHMACはCookie偽造による認証バイパスにつながるため、fail-closedであることが必須。
func TestManualUseCase_EmptySecretFailsClosed(t *testing.T) {
	empty := &manualUseCase{cfg: ManualConfig{ClientID: "id", ClientSecret: "", RedirectURL: "https://example.com/cb"}, tokenEndpoint: googleTokenEndpoint}

	t.Run("空鍵ではstateもCookieも発行されない", func(t *testing.T) {
		if got := empty.MakeState("stamprally"); got != "" {
			t.Errorf("MakeState() = %q, want empty", got)
		}
		if got := empty.MakeSessionCookie("taro.nutfes@gmail.com"); got != "" {
			t.Errorf("MakeSessionCookie() = %q, want empty", got)
		}
	})

	t.Run("空鍵で偽造したCookieは検証を通らない", func(t *testing.T) {
		// 攻撃の再現: 鍵なし（空鍵）でpayloadと署名を自作する
		payload := base64.RawURLEncoding.EncodeToString([]byte(`{"email":"attacker.nutfes@gmail.com","exp":99999999999999}`))
		mac := hmac.New(sha256.New, []byte(""))
		mac.Write([]byte(payload))
		forged := payload + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

		if _, ok := empty.VerifySessionCookie(forged); ok {
			t.Fatalf("VerifySessionCookie(forged) ok = true, want false（空鍵の偽造Cookieは拒否必須）")
		}
		if _, ok := empty.VerifyState(forged); ok {
			t.Errorf("VerifyState(forged) ok = true, want false")
		}
	})
}

// VerifyUploadToken のテーブルテスト。UploadTokenが未設定の場合は正しい値を提示しても
// 常に拒否される（issue #444と同じfail-closed方針）ことを含めて確認する。
func TestManualUseCase_VerifyUploadToken(t *testing.T) {
	withToken := &manualUseCase{ //nolint:gosec // テスト用のダミートークンであり実資格情報ではない
		cfg: ManualConfig{UploadToken: "correct-upload-token"},
	}
	withoutToken := &manualUseCase{cfg: ManualConfig{UploadToken: ""}}

	cases := []struct {
		name   string
		u      *manualUseCase
		header string
		want   bool
	}{
		{
			name:   "正しいBearerトークンは許可",
			u:      withToken,
			header: "Bearer correct-upload-token",
			want:   true,
		},
		{
			name:   "スキームが小文字のbearerでも許可（大小無視）",
			u:      withToken,
			header: "bearer correct-upload-token",
			want:   true,
		},
		{
			name:   "トークン不一致は拒否",
			u:      withToken,
			header: "Bearer wrong-token",
			want:   false,
		},
		{
			name:   "スキームがBearerでない場合は拒否",
			u:      withToken,
			header: "Basic correct-upload-token",
			want:   false,
		},
		{
			name:   "空ヘッダは拒否",
			u:      withToken,
			header: "",
			want:   false,
		},
		{
			name:   "UploadToken未設定なら正しい値でも拒否（fail-closed）",
			u:      withoutToken,
			header: "Bearer correct-upload-token",
			want:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.u.VerifyUploadToken(c.header)
			if got != c.want {
				t.Errorf("VerifyUploadToken(%q) = %v, want %v", c.header, got, c.want)
			}
		})
	}
}

// SaveManual の正常系・異常系を確認する。
func TestManualUseCase_SaveManual(t *testing.T) {
	t.Run("正常なIDで保存するとファイル内容とバイト数が一致する", func(t *testing.T) {
		dir := t.TempDir()
		u := &manualUseCase{cfg: ManualConfig{ManualDir: dir}}

		content := "<html>hello</html>"
		n, err := u.SaveManual("summer-fes", strings.NewReader(content))
		if err != nil {
			t.Fatalf("SaveManual() error = %v, want nil", err)
		}
		if n != int64(len(content)) {
			t.Errorf("SaveManual() n = %d, want %d", n, len(content))
		}

		got, err := os.ReadFile(filepath.Join(dir, "summer-fes.html")) //nolint:gosec // G304: dirはt.TempDir()由来のテスト専用パスであり外部入力ではない
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}
		if string(got) != content {
			t.Errorf("saved content = %q, want %q", string(got), content)
		}
	})

	t.Run("不正なIDはErrInvalidManualIDを返す", func(t *testing.T) {
		dir := t.TempDir()
		u := &manualUseCase{cfg: ManualConfig{ManualDir: dir}}

		if _, err := u.SaveManual("../etc/passwd", strings.NewReader("dummy")); !errors.Is(err, ErrInvalidManualID) {
			t.Errorf("SaveManual() error = %v, want ErrInvalidManualID", err)
		}
	})

	t.Run("既存ファイルへの上書きで内容が置き換わる", func(t *testing.T) {
		dir := t.TempDir()
		u := &manualUseCase{cfg: ManualConfig{ManualDir: dir}}

		if _, err := u.SaveManual("summer-fes", strings.NewReader("old content")); err != nil {
			t.Fatalf("SaveManual() first write error = %v", err)
		}
		if _, err := u.SaveManual("summer-fes", strings.NewReader("new content")); err != nil {
			t.Fatalf("SaveManual() second write error = %v", err)
		}

		got, err := os.ReadFile(filepath.Join(dir, "summer-fes.html")) //nolint:gosec // G304: dirはt.TempDir()由来のテスト専用パスであり外部入力ではない
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}
		if string(got) != "new content" {
			t.Errorf("saved content = %q, want %q", string(got), "new content")
		}
	})
}

// PublicManualURL がRedirectURLから公開URLを正しく導出することを確認する。
func TestManualUseCase_PublicManualURL(t *testing.T) {
	u := &manualUseCase{cfg: ManualConfig{RedirectURL: "https://example.com/manuals/oauth/callback"}}

	got := u.PublicManualURL("summer-fes")
	want := "https://example.com/manuals/summer-fes"
	if got != want {
		t.Errorf("PublicManualURL() = %q, want %q", got, want)
	}
}
