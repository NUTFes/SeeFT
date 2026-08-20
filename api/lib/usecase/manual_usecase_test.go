package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestManualUseCase() *manualUseCase {
	return &manualUseCase{
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

	u := &manualUseCase{cfg: ManualConfig{ManualDir: dir}, tokenEndpoint: "https://oauth2.googleapis.com/token"}

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

	header := map[string]string{"alg": "RS256", "typ": "JWT"}
	claims := map[string]any{
		"email":          email,
		"email_verified": emailVerified,
	}

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
}
