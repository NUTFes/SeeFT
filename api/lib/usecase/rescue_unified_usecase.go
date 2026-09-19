package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type rescueUnifiedUseCase struct {
	questionRescueRepository    repository.QuestionRescueRepository
	shorthandedRescueRepository repository.ShorthandedRescueRepository
	troubleRescueRepository     repository.TroubleRescueRepository
	userRepository              repository.UserRepository
	taskRepository              repository.TaskRepository
}

type RescueUnifiedUseCase interface {
	GetAllRescues(context.Context) ([]entity.RescueResponse, error)
	GetRescuesByUserID(context.Context, string) ([]entity.RescueResponse, error)
	SaveRescueToSpreadsheet(map[string]interface{}) error
}

func NewRescueUnifiedUseCase(
	qr repository.QuestionRescueRepository,
	sr repository.ShorthandedRescueRepository,
	tr repository.TroubleRescueRepository,
	ur repository.UserRepository,
	tar repository.TaskRepository,
) RescueUnifiedUseCase {
	return &rescueUnifiedUseCase{qr, sr, tr, ur, tar}
}

// 全件取得
func (ru *rescueUnifiedUseCase) GetAllRescues(c context.Context) ([]entity.RescueResponse, error) {
	var allRescues []entity.RescueResponse

	// Question Rescues取得
	questionRescues, err := ru.getQuestionRescues(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescues")
	}
	allRescues = append(allRescues, questionRescues...)

	// Shorthanded Rescues取得
	shorthandedRescues, err := ru.getShorthandedRescues(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues")
	}
	allRescues = append(allRescues, shorthandedRescues...)

	// Trouble Rescues取得
	troubleRescues, err := ru.getTroubleRescues(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues")
	}
	allRescues = append(allRescues, troubleRescues...)

	// 時刻でソート（降順）
	sort.Slice(allRescues, func(i, j int) bool {
		return allRescues[i].Time > allRescues[j].Time
	})

	return allRescues, nil
}

// ユーザーID別取得
func (ru *rescueUnifiedUseCase) GetRescuesByUserID(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var allRescues []entity.RescueResponse

	// Question Rescues取得
	questionRescues, err := ru.getQuestionRescues(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescues")
	}
	allRescues = append(allRescues, questionRescues...)

	// Shorthanded Rescues取得
	shorthandedRescues, err := ru.getShorthandedRescues(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues")
	}
	allRescues = append(allRescues, shorthandedRescues...)

	// Trouble Rescues取得
	troubleRescues, err := ru.getTroubleRescues(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues")
	}
	allRescues = append(allRescues, troubleRescues...)

	// 時刻でソート（降順）
	sort.Slice(allRescues, func(i, j int) bool {
		return allRescues[i].Time > allRescues[j].Time
	})

	return allRescues, nil
}

// Question Rescues取得
func (ru *rescueUnifiedUseCase) getQuestionRescues(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = ru.questionRescueRepository.FindByUserID(c, userID)
	} else {
		rows, err = ru.questionRescueRepository.All(c)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var rescues []entity.RescueResponse
	for rows.Next() {
		var questionRescue entity.QuestionRescueForGet
		var response sql.NullString
		err := rows.Scan(&questionRescue.ID, &questionRescue.UserID, &questionRescue.Question, &questionRescue.Status, &response, &questionRescue.Time, &questionRescue.CreatedAt, &questionRescue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questionRescue.Response = response.String

		// ユーザー名を取得
		userName, err := ru.getUserName(c, strconv.Itoa(questionRescue.UserID))
		if err != nil {
			userName = "不明なユーザー"
		}

		rescue := entity.NewQuestionRescueResponse(&questionRescue, userName)
		rescues = append(rescues, *rescue)
	}
	return rescues, nil
}

// Shorthanded Rescues取得
func (ru *rescueUnifiedUseCase) getShorthandedRescues(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = ru.shorthandedRescueRepository.FindByUserID(c, userID)
	} else {
		rows, err = ru.shorthandedRescueRepository.All(c)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var rescues []entity.RescueResponse
	for rows.Next() {
		var shorthandedRescue entity.ShorthandedRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &place, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shorthandedRescue.Place = place.String
		shorthandedRescue.Response = response.String

		// ユーザー名を取得
		userName, err := ru.getUserName(c, strconv.Itoa(shorthandedRescue.UserID))
		if err != nil {
			userName = "不明なユーザー"
		}

		// タスク名を取得
		taskName, err := ru.getTaskName(c, strconv.Itoa(shorthandedRescue.TaskID))
		if err != nil {
			taskName = "不明なタスク"
		}

		rescue := entity.NewShorthandedRescueResponse(&shorthandedRescue, userName, taskName)
		rescues = append(rescues, *rescue)
	}
	return rescues, nil
}

// Trouble Rescues取得
func (ru *rescueUnifiedUseCase) getTroubleRescues(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = ru.troubleRescueRepository.FindByUserID(c, userID)
	} else {
		rows, err = ru.troubleRescueRepository.All(c)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var rescues []entity.RescueResponse
	for rows.Next() {
		var troubleRescue entity.TroubleRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &place, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		troubleRescue.Place = place.String
		troubleRescue.Response = response.String

		// ユーザー名を取得
		userName, err := ru.getUserName(c, strconv.Itoa(troubleRescue.UserID))
		if err != nil {
			userName = "不明なユーザー"
		}

		// タスク名を取得
		taskName, err := ru.getTaskName(c, strconv.Itoa(troubleRescue.TaskID))
		if err != nil {
			taskName = "タスク外"
		}

		rescue := entity.NewTroubleRescueResponse(&troubleRescue, userName, taskName)
		rescues = append(rescues, *rescue)
	}
	return rescues, nil
}

// ユーザー名取得
func (ru *rescueUnifiedUseCase) getUserName(c context.Context, userID string) (string, error) {
	row, err := ru.userRepository.Find(c, userID)
	if err != nil {
		return "", err
	}

	var id int
	var name, mail, studentNumber, tel, password string
	var gradeID, departmentID, bureauID, roleID int
	var createdAt, updatedAt string
	var slackUserID sql.NullString
	err = row.Scan(&id, &name, &mail, &gradeID, &departmentID, &bureauID, &roleID, &studentNumber, &tel, &password, &createdAt, &updatedAt, &slackUserID)
	if err != nil {
		return "", err
	}

	return name, nil
}

// タスク名取得
func (ru *rescueUnifiedUseCase) getTaskName(c context.Context, taskID string) (string, error) {
	row, err := ru.taskRepository.Find(c, taskID)
	if err != nil {
		return "", err
	}

	var id int
	var task, url, manualUrl, color, remark string
	var placeID, bureauID, maxMember, yearID int
	var createdAt, updatedAt string
	err = row.Scan(&id, &task, &placeID, &url, &manualUrl, &bureauID, &maxMember, &color, &remark, &yearID, &createdAt, &updatedAt)
	if err != nil {
		return "", err
	}

	return task, nil
}

// スプシ保存用関数（Google Sheets APIのラッパーを想定）
func (ru *rescueUnifiedUseCase) SaveRescueToSpreadsheet(data map[string]interface{}) error {
	// Google Sheets APIで保存処理
	// GAS送信
	err := ru.SendRescueToGAS(data)
	if err != nil {
		fmt.Println("GAS送信失敗:", err)
		return err
	}
	return nil
}

// GAS送信関数
func (ru *rescueUnifiedUseCase) SendRescueToGAS(data map[string]interface{}) error {
	// GASのURLを環境変数から取得
	gasURL := os.Getenv("RESCUE_GAS_URL")
	if gasURL == "" {
		return errors.New("GAS URLが設定されていません (RESCUE_GAS_URL)")
	}
	parsedURL, err := url.Parse(gasURL)
	if err != nil || parsedURL.Scheme != "https" {
		return errors.New("RESCUE_GAS_URL が不正です（https のみ許可）")
	}

	// JSONに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "レスキューデータのJSON変換失敗")
	}

	return postToGAS(&http.Client{}, gasURL, jsonData)
}

// GASへPOSTし、リダイレクトの各段階をログに出す。
// GASのウェブアプリは「POST → 302」「リダイレクト先のGET → 200」の2段階で応答する。
// 2026-09-19に書き込み済みなのに404で失敗扱いになったとき、どちらの段階の404かが
// ログから分からなかったので、段階ごとのステータスと経過時間を残す（#547）
func postToGAS(client *http.Client, gasURL string, body []byte) error {
	req, err := http.NewRequest("POST", gasURL, bytes.NewReader(body)) //nolint:gosec // G704: gasURL はhttps スキームを検証済みの環境変数
	if err != nil {
		return errors.Wrap(err, "GASリクエスト作成失敗")
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	var hops []string
	traced := *client
	traced.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		// next.Response はこのリダイレクトを起こした直前の応答
		prev := via[len(via)-1]
		hops = append(hops, fmt.Sprintf("%s %s → %d (%.2fs)", prev.Method, gasLogURL(prev.URL), next.Response.StatusCode, time.Since(start).Seconds()))
		if client.CheckRedirect != nil {
			return client.CheckRedirect(next, via)
		}
		// CheckRedirectを設定すると既定の上限が外れるので、既定と同じ10回で止める
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}

	resp, err := traced.Do(req) //nolint:gosec // G704: gasURL はhttps スキームを検証済みの環境変数
	if err != nil {
		log.Printf("GAS送信: %s → エラー (%.2fs): %v", strings.Join(hops, " / "), time.Since(start).Seconds(), err)
		return errors.Wrap(err, "GASへの送信失敗")
	}
	defer func() { _ = resp.Body.Close() }()

	hops = append(hops, fmt.Sprintf("%s %s → %d (%.2fs)", resp.Request.Method, gasLogURL(resp.Request.URL), resp.StatusCode, time.Since(start).Seconds()))
	if resp.StatusCode != http.StatusOK {
		// エラーページのタイトルで、Googleのどの種類の404かを見分ける
		log.Printf("GAS送信: %s title=%q", strings.Join(hops, " / "), gasPageTitle(resp.Body))
		return errors.Errorf("GASが非OKステータスを返しました: %d (%s %s)", resp.StatusCode, resp.Request.Method, gasLogURL(resp.Request.URL))
	}
	// 200の本文はdoPostの戻り値（Success / Duplicate of N / Error: ...）
	head, _ := io.ReadAll(io.LimitReader(resp.Body, 100))
	log.Printf("GAS送信: %s body=%q", strings.Join(hops, " / "), string(head))
	return nil
}

// ウェブアプリのURLは知っていれば誰でもスプシに書き込めるので、パス中のデプロイIDを伏せる。
// クエリ（リダイレクト先のuser_content_key）も出さない
var gasDeploymentIDPattern = regexp.MustCompile(`/s/[^/]+`)

func gasLogURL(u *url.URL) string {
	return u.Host + gasDeploymentIDPattern.ReplaceAllString(u.Path, "/s/…")
}

var htmlTitlePattern = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

func gasPageTitle(body io.Reader) string {
	b, _ := io.ReadAll(io.LimitReader(body, 64*1024))
	m := htmlTitlePattern.FindSubmatch(b)
	if m == nil {
		return ""
	}
	return strings.TrimSpace(string(m[1]))
}
