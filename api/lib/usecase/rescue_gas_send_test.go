package usecase

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const fakeGASPath = "/macros/s/AKfycbSECRET/exec"

// GASのPOSTを真似る。respond で応答（302の飛び先や、404・200のエラーページ）を決める
func newFakeGAS(t *testing.T, respond func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc(fakeGASPath, respond)
	ts := httptest.NewTLSServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func redirectTo(location string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, location, http.StatusFound)
	}
}

func errorPage(status int, title string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, "<html><head><title>%s</title></head></html>", title)
	}
}

// 送ったリクエストを数える。テスト用サーバー以外（本物のGoogle）へは出さない
type recordingTransport struct {
	base     http.RoundTripper
	testHost string
	requests []string
}

func (rt *recordingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	rt.requests = append(rt.requests, r.Method+" "+r.URL.Host+r.URL.Path)
	if r.URL.Host != rt.testHost {
		return nil, fmt.Errorf("テスト用サーバー以外へのリクエスト: %s", r.URL)
	}
	return rt.base.RoundTrip(r)
}

// postToGASを呼び、ログ出力・送ったリクエスト・戻り値を返す
func callPostToGAS(t *testing.T, ts *httptest.Server) (string, []string, error) {
	t.Helper()
	var buf bytes.Buffer
	orig := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(orig) })

	rt := &recordingTransport{base: ts.Client().Transport, testHost: strings.TrimPrefix(ts.URL, "https://")}
	err := postToGAS(&http.Client{Transport: rt}, ts.URL+fakeGASPath, []byte(`{"rescue_type":"question"}`))
	return buf.String(), rt.requests, err
}

func assertNoSecrets(t *testing.T, logs string) {
	t.Helper()
	for _, secret := range []string{"AKfycbSECRET", "ONE_TIME_KEY"} {
		if strings.Contains(logs, secret) {
			t.Errorf("ログに %s が出ている: %s", secret, logs)
		}
	}
}

func TestPostToGAS_結果置き場への302で成功とし取りに行かない(t *testing.T) {
	ts := newFakeGAS(t, redirectTo("https://script.googleusercontent.com/macros/echo?user_content_key=ONE_TIME_KEY&lib=x"))
	logs, requests, err := callPostToGAS(t, ts)
	if err != nil {
		t.Fatalf("成功のはずがエラー: %v", err)
	}
	if len(requests) != 1 || !strings.HasPrefix(requests[0], "POST ") {
		t.Errorf("POSTの1回だけのはず: %v", requests)
	}
	if !strings.Contains(logs, `/macros/s/…/exec" → 302 "script.googleusercontent.com/macros/echo"`) {
		t.Errorf("ログに段階と飛び先が無い: %s", logs)
	}
	assertNoSecrets(t, logs)
}

func TestPostToGAS_ログイン画面への302は失敗(t *testing.T) {
	// デプロイのアクセス権が「ログインが必要」に変わると、doPostは動かずログイン画面へ飛ばされる
	ts := newFakeGAS(t, redirectTo("https://accounts.google.com/ServiceLogin?continue=x"))
	logs, requests, err := callPostToGAS(t, ts)
	if err == nil || !strings.Contains(err.Error(), "302 accounts.google.com/ServiceLogin") {
		t.Fatalf("失敗になるはず: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("ログイン画面へは行かないはず: %v", requests)
	}
	assertNoSecrets(t, logs)
}

func TestPostToGAS_POSTが404なら失敗しタイトルを残す(t *testing.T) {
	ts := newFakeGAS(t, errorPage(http.StatusNotFound, "ページが見つかりません"))
	logs, _, err := callPostToGAS(t, ts)
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("失敗になるはず: %v", err)
	}
	if !strings.Contains(logs, `→ 404`) || !strings.Contains(logs, `title="ページが見つかりません"`) {
		t.Errorf("ログにステータスとタイトルが無い: %s", logs)
	}
	assertNoSecrets(t, logs)
}

func TestPostToGAS_POSTに200が直接返るのは失敗(t *testing.T) {
	// 9/19の計測で、エラーページが200で返ることがあった。200だけでは成功と判断しない
	ts := newFakeGAS(t, errorPage(http.StatusOK, "エラー"))
	logs, _, err := callPostToGAS(t, ts)
	if err == nil || !strings.Contains(err.Error(), "200") {
		t.Fatalf("失敗になるはず: %v", err)
	}
	if !strings.Contains(logs, `title="エラー"`) {
		t.Errorf("ログにタイトルが無い: %s", logs)
	}
	assertNoSecrets(t, logs)
}
