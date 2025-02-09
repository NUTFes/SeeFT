package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// structが持つべき関数を定義
type RescueUseCase interface {
	CreateRescue(context.Context, string) error
}

// クラス
type rescueUseCase struct{}

// インスタンス化
func NewRescueUseCase() RescueUseCase {
	return &rescueUseCase{}
}

// structの中身を定義
func (u *rescueUseCase) CreateRescue(ctx context.Context, taskID string) error {
	url := "https://script.google.com/macros/s/AKfycbz6ueqUXhXX6EvFwBCXPBtD571CpGy8hHJysSRiUIBPaEesil8qUB3Q_fB1Oy0Fk3MGPg/exec"
	data := map[string]string{"task_id": taskID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send rescue request")
	}

	return nil
}
