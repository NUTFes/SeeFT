package usecase

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// GASのウェブアプリを真似る。POST /macros/s/{デプロイID}/exec → 302 → GET /macros/echo
func newFakeGAS(t *testing.T, postStatus, echoStatus int) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/macros/s/AKfycbSECRET/exec", func(w http.ResponseWriter, r *http.Request) {
		if postStatus != http.StatusFound {
			w.WriteHeader(postStatus)
			_, _ = w.Write([]byte("<html><head><title>ページが見つかりません</title></head></html>"))
			return
		}
		http.Redirect(w, r, "/macros/echo?user_content_key=ONE_TIME_KEY&lib=x", http.StatusFound)
	})
	mux.HandleFunc("/macros/echo", func(w http.ResponseWriter, r *http.Request) {
		if echoStatus != http.StatusOK {
			w.WriteHeader(echoStatus)
			_, _ = w.Write([]byte("<html><head><title>Google Drive - Page Not Found</title></head></html>"))
			return
		}
		_, _ = w.Write([]byte("Success"))
	})
	ts := httptest.NewTLSServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// postToGASを呼び、戻り値とログ出力を返す
func callPostToGAS(t *testing.T, ts *httptest.Server) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	orig := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(orig) })
	err := postToGAS(ts.Client(), ts.URL+"/macros/s/AKfycbSECRET/exec", []byte(`{"rescue_type":"question"}`))
	return buf.String(), err
}

func assertNoSecrets(t *testing.T, logs string) {
	t.Helper()
	for _, secret := range []string{"AKfycbSECRET", "ONE_TIME_KEY"} {
		if strings.Contains(logs, secret) {
			t.Errorf("ログに %s が出ている: %s", secret, logs)
		}
	}
}

func TestPostToGAS_302のあと200なら成功し両段階をログに出す(t *testing.T) {
	ts := newFakeGAS(t, http.StatusFound, http.StatusOK)
	logs, err := callPostToGAS(t, ts)
	if err != nil {
		t.Fatalf("成功のはずがエラー: %v", err)
	}
	for _, want := range []string{"POST ", "/macros/s/…/exec → 302", "GET ", "/macros/echo → 200", `body="Success"`} {
		if !strings.Contains(logs, want) {
			t.Errorf("ログに %q が無い: %s", want, logs)
		}
	}
	assertNoSecrets(t, logs)
}

func TestPostToGAS_リダイレクト先のGETが404なら段階とタイトルが分かる(t *testing.T) {
	ts := newFakeGAS(t, http.StatusFound, http.StatusNotFound)
	logs, err := callPostToGAS(t, ts)
	if err == nil || !strings.Contains(err.Error(), "404 (GET ") {
		t.Fatalf("GET段階の404エラーになるはず: %v", err)
	}
	for _, want := range []string{"/macros/s/…/exec → 302", "/macros/echo → 404", `title="Google Drive - Page Not Found"`} {
		if !strings.Contains(logs, want) {
			t.Errorf("ログに %q が無い: %s", want, logs)
		}
	}
	assertNoSecrets(t, logs)
}

func TestPostToGAS_POST自体が404ならリダイレクト前の段階と分かる(t *testing.T) {
	ts := newFakeGAS(t, http.StatusNotFound, http.StatusOK)
	logs, err := callPostToGAS(t, ts)
	if err == nil || !strings.Contains(err.Error(), "404 (POST ") {
		t.Fatalf("POST段階の404エラーになるはず: %v", err)
	}
	if strings.Contains(logs, "/macros/echo") {
		t.Errorf("リダイレクトしていないのにechoが出ている: %s", logs)
	}
	if !strings.Contains(logs, `title="ページが見つかりません"`) {
		t.Errorf("エラーページのタイトルが無い: %s", logs)
	}
	assertNoSecrets(t, logs)
}
