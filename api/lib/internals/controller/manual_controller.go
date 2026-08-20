package controller

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

// マニュアル閲覧セッションのCookie名（issue #444）
const manualSessionCookieName = "seeft_manual_session"

// Cookieの有効期間（usecase側のセッション署名TTLと揃えている）
const manualSessionCookieMaxAge = 30 * 24 * time.Hour

type manualController struct {
	u usecase.ManualUseCase
}

// ManualController はnutfes限定のマニュアルHTML配信（issue #444）を扱う。
type ManualController interface {
	ShowManual(echo.Context) error
	OAuthCallback(echo.Context) error
}

func NewManualController(u usecase.ManualUseCase) ManualController {
	return &manualController{u}
}

// ShowManual はマニュアルHTMLを返す。有効なセッションCookieがなければログイン誘導ページを返す。
func (mc *manualController) ShowManual(c echo.Context) error {
	id := c.Param("id")

	if cookie, err := c.Cookie(manualSessionCookieName); err == nil && cookie.Value != "" {
		if email, ok := mc.u.VerifySessionCookie(cookie.Value); ok {
			path, exists := mc.u.ManualPath(id)
			if !exists {
				log.Printf("manual_gate: not_found id=%s email=%s", id, email)
				return c.HTML(http.StatusNotFound, mc.renderNotFoundPage(id))
			}
			log.Printf("manual_gate: allowed id=%s email=%s", id, email)
			return c.File(path)
		}
	}

	log.Printf("manual_gate: access id=%s", id)
	return c.HTML(http.StatusOK, mc.renderEntryPage(id, ""))
}

// OAuthCallback はGoogleからの認可コールバックを処理し、許可対象なら閲覧セッションを発行する。
func (mc *manualController) OAuthCallback(c echo.Context) error {
	// 同意画面でユーザーがキャンセルした場合、Googleは code の代わりに error を付与する
	if errParam := c.QueryParam("error"); errParam != "" {
		manualID, ok := mc.u.VerifyState(c.QueryParam("state"))
		notice := "ログインがキャンセルされました。もう一度お試しください。"
		if ok {
			log.Printf("manual_gate: access id=%s reason=cancelled", manualID)
			return c.HTML(http.StatusOK, mc.renderEntryPage(manualID, notice))
		}
		log.Printf("manual_gate: access reason=cancelled")
		return c.HTML(http.StatusOK, mc.renderMessagePage("ログインがキャンセルされました", notice))
	}

	manualID, ok := mc.u.VerifyState(c.QueryParam("state"))
	if !ok {
		notice := "ログインの有効期限が切れました。もう一度お試しください。"
		if manualID != "" {
			log.Printf("manual_gate: state_invalid id=%s", manualID)
			return c.HTML(http.StatusOK, mc.renderEntryPage(manualID, notice))
		}
		log.Printf("manual_gate: state_invalid")
		return c.HTML(http.StatusOK, mc.renderMessagePage("ログインの有効期限切れ", notice))
	}

	email, err := mc.u.ExchangeCode(c.Request().Context(), c.QueryParam("code"))
	if err != nil {
		log.Printf("manual_gate: exchange_fail id=%s err=%v", manualID, err)
		notice := "ログインの有効期限が切れました。もう一度お試しください。"
		return c.HTML(http.StatusOK, mc.renderEntryPage(manualID, notice))
	}

	if !mc.u.IsAllowed(email) {
		log.Printf("manual_gate: denied id=%s email=%s", manualID, email)
		return c.HTML(http.StatusForbidden, mc.renderDeniedPage(email))
	}

	c.SetCookie(&http.Cookie{
		Name:     manualSessionCookieName,
		Value:    mc.u.MakeSessionCookie(email),
		Path:     "/manuals",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(manualSessionCookieMaxAge.Seconds()),
	})
	log.Printf("manual_gate: allowed id=%s email=%s", manualID, email)
	return c.Redirect(http.StatusFound, "/manuals/"+manualID)
}

// manualPageStyle は本機能のページ共通スタイル（白背景・ティール#009688・システムの日本語フォント・幅480px）。
const manualPageStyle = `
    body {
      background: #ffffff;
      color: #212121;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Hiragino Kaku Gothic ProN", "Hiragino Sans", Meiryo, sans-serif;
      max-width: 480px;
      margin: 48px auto;
      padding: 0 24px;
      line-height: 1.7;
    }
    h1 {
      color: #009688;
      font-size: 1.3rem;
    }
    p {
      font-size: 0.95rem;
    }
    .notice {
      background: #f1f8f7;
      border-left: 4px solid #009688;
      padding: 8px 12px;
    }
    .btn {
      display: inline-block;
      margin-top: 16px;
      padding: 10px 24px;
      background: #009688;
      color: #ffffff;
      text-decoration: none;
      border-radius: 4px;
      font-weight: bold;
    }
    .account {
      font-weight: bold;
      word-break: break-all;
    }
`

// manualPageHTML はタイトルと本文HTMLから完結したHTMLページを組み立てる。
// title・bodyHTML はいずれも呼び出し側で必要な範囲をエスケープ済みとして渡すこと。
func manualPageHTML(title, bodyHTML string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>%s</style>
</head>
<body>
%s
</body>
</html>`, html.EscapeString(title), manualPageStyle, bodyHTML)
}

// renderEntryPage はログイン誘導ページを組み立てる。notice が空でなければ通知文を併記する。
func (mc *manualController) renderEntryPage(manualID, notice string) string {
	authURL := mc.u.BuildAuthURL(manualID)

	var noticeHTML string
	if notice != "" {
		noticeHTML = fmt.Sprintf(`<p class="notice">%s</p>`, html.EscapeString(notice))
	}

	title := "マニュアル閲覧にはログインが必要です"
	body := fmt.Sprintf(`<h1>%s</h1>
%s
<p>このマニュアルは技大祭実行委員会のメンバー限定です。「example.nutfes@gmail.com」のような実行委員会のGoogleアカウントでログインしてください。</p>
<p><a class="btn" href="%s">Googleでログイン</a></p>`, html.EscapeString(title), noticeHTML, html.EscapeString(authURL))

	return manualPageHTML(title, body)
}

// renderMessagePage はログインボタンを持たない汎用メッセージページを組み立てる（manualIDが復元できない場合用）。
func (mc *manualController) renderMessagePage(title, message string) string {
	body := fmt.Sprintf(`<h1>%s</h1>
<p>%s</p>`, html.EscapeString(title), html.EscapeString(message))
	return manualPageHTML(title, body)
}

// renderNotFoundPage は存在しないマニュアルID用の404ページを組み立てる。
func (mc *manualController) renderNotFoundPage(manualID string) string {
	title := "マニュアルが見つかりません"
	body := fmt.Sprintf(`<h1>%s</h1>
<p>指定されたマニュアル（%s）は存在しないか、削除された可能性があります。</p>`, html.EscapeString(title), html.EscapeString(manualID))
	return manualPageHTML(title, body)
}

// renderDeniedPage は許可対象外アカウントによるアクセス時の403ページを組み立てる。
func (mc *manualController) renderDeniedPage(email string) string {
	title := "閲覧には委員アカウントが必要です"
	body := fmt.Sprintf(`<h1>%s</h1>
<p>現在ログイン中のアカウント: <span class="account">%s</span></p>
<p>このマニュアルは技大祭実行委員会のメンバー限定です。ブラウザのGoogleアカウントを実行委員会用のアカウント（example.nutfes@gmail.com など）に切り替えてから、もう一度お試しください。</p>`, html.EscapeString(title), html.EscapeString(email))
	return manualPageHTML(title, body)
}
