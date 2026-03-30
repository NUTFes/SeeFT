package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

// ---------- テスト用 db.Client 実装 ----------

type testClient struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

func (c *testClient) DB() *sql.DB    { return c.sqlDB }
func (c *testClient) GormDB() *gorm.DB { return c.gormDB }
func (c *testClient) CloseDB()       { c.sqlDB.Close() }

func newTestClient(t testing.TB) db.Client {
	t.Helper()

	dbUser := envOrDefault("NUTMEG_DB_USER", "seeft")
	dbPassword := envOrDefault("NUTMEG_DB_PASSWORD", "password")
	dbHost := envOrDefault("NUTMEG_DB_HOST", "localhost")
	dbPort := envOrDefault("NUTMEG_DB_PORT", "5432")
	dbName := envOrDefault("NUTMEG_DB_NAME", "seeft_db")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		t.Skipf("DB not available: %v (set NUTMEG_DB_* env vars)", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		sqlDB.Close()
		t.Fatalf("gorm.Open: %v", err)
	}

	return &testClient{sqlDB: sqlDB, gormDB: gormDB}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ---------- テスト用 ShiftUseCase 組み立て ----------

func newTestShiftUseCase(t testing.TB, client db.Client) ShiftUseCase {
	t.Helper()
	crud := abstract.NewCrud(client)

	return NewShiftUseCase(
		repository.NewShiftRepository(client, crud),
		repository.NewShiftCardRepository(client),
		repository.NewTaskRepository(client, crud),
		repository.NewUserRepository(client, crud),
		repository.NewYearRepository(client, crud),
		repository.NewDateRepository(client, crud),
		repository.NewTimeRepository(client, crud),
		repository.NewWeatherRepository(client, crud),
		repository.NewPlaceRepository(client, crud),
		repository.NewGradeRepository(client, crud),
		repository.NewBureauRepository(client, crud),
	)
}

// ---------- テストデータ投入 ----------

const (
	benchNumUsers = 100  // N+1の差が出る規模（ユーザー数）
	benchNumTimes = 8    // 連続する時間枠数
	benchYearID   = 43   // seed.sql の year_id=43 (2024)
	benchDateID   = 2    // seed.sql の date_id=2 (1日目)
	benchWeatherID = 1   // seed.sql の weather_id=1 (晴れ)
	benchTaskID   = 4    // seed.sql の task_id=4 (テスト1)
	benchStartTimeID = 41 // time_id=41 -> 10:00
)

func seedBenchData(t testing.TB, sqlDB *sql.DB) (targetUserID int) {
	t.Helper()
	ctx := context.Background()

	// クリーンアップ: ベンチ用データを削除（既存seed dataは残す）
	t.Cleanup(func() {
		sqlDB.ExecContext(ctx, "DELETE FROM shifts WHERE user_id > 3")
		sqlDB.ExecContext(ctx, "DELETE FROM users WHERE id > 3")
	})

	// ベンチ用ユーザーを大量投入
	for i := 0; i < benchNumUsers; i++ {
		var userID int
		err := sqlDB.QueryRowContext(ctx,
			`INSERT INTO users (name, mail, grade_id, department_id, bureau_id, role_id, student_number, tel, password)
			 VALUES ($1, $2, $3, 1, $4, 1, $5, '00000000000', 'dummy')
			 RETURNING id`,
			fmt.Sprintf("bench_user_%04d", i),
			fmt.Sprintf("bench%04d@example.com", i),
			(i%10)+1, // grade_id: 1-10
			(i%8)+1,  // bureau_id: 1-8
			90000000+i,
		).Scan(&userID)
		if err != nil {
			t.Fatalf("insert user %d: %v", i, err)
		}

		// 最初のユーザーをターゲットにする
		if i == 0 {
			targetUserID = userID
		}

		// 各ユーザーに対して、同じ task/year/date/weather で複数の時間枠にシフトを作成
		for j := 0; j < benchNumTimes; j++ {
			timeID := benchStartTimeID + j
			_, err := sqlDB.ExecContext(ctx,
				`INSERT INTO shifts (task_id, user_id, year_id, date_id, time_id, weather_id, is_attendance)
				 VALUES ($1, $2, $3, $4, $5, $6, true)`,
				benchTaskID, userID, benchYearID, benchDateID, timeID, benchWeatherID,
			)
			if err != nil {
				t.Fatalf("insert shift user=%d time=%d: %v", userID, timeID, err)
			}
		}
	}

	t.Logf("Seeded %d users x %d time slots = %d shifts (target user_id=%d)",
		benchNumUsers, benchNumTimes, benchNumUsers*benchNumTimes, targetUserID)
	return targetUserID
}

// ---------- ベンチマーク ----------

func BenchmarkGetShiftCardsByUserAndDateAndWeather(b *testing.B) {
	// SQLログを無効化
	os.Setenv("DEBUG_SQL", "0")

	client := newTestClient(b)
	defer client.CloseDB()

	uc := newTestShiftUseCase(b, client)
	targetUserID := seedBenchData(b, client.DB())

	userIDStr := fmt.Sprintf("%d", targetUserID)
	dateIDStr := fmt.Sprintf("%d", benchDateID)
	weatherIDStr := fmt.Sprintf("%d", benchWeatherID)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := uc.GetShiftCardsByUserAndDateAndWeather(ctx, userIDStr, dateIDStr, weatherIDStr)
		if err != nil {
			b.Fatalf("GetShiftCardsByUserAndDateAndWeather: %v", err)
		}
	}
}

// ---------- 整合性テスト（冪等性チェック） ----------

// canonicalJSON は ShiftCard のスライスを正規化してJSON化する
// メンバー配列をソートして順序依存を排除
func canonicalJSON(cards []entity.ShiftCard) ([]byte, error) {
	// ShiftCards自体を開始時刻でソート
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].StartTime < cards[j].StartTime
	})

	for i := range cards {
		// 各時間枠のメンバーをソート
		for j := range cards[i].ShiftMembers {
			sortMembers(cards[i].ShiftMembers[j].Members)
		}
		// BeforeMembers/AfterMembersもソート
		sortMembers(cards[i].BeforeMembers.Members)
		sortMembers(cards[i].AfterMembers.Members)
	}

	return json.MarshalIndent(cards, "", "  ")
}

func sortMembers(members []entity.ShiftMember) {
	sort.Slice(members, func(i, j int) bool {
		if members[i].Name != members[j].Name {
			return members[i].Name < members[j].Name
		}
		if members[i].Grade != members[j].Grade {
			return members[i].Grade < members[j].Grade
		}
		return members[i].Bureau < members[j].Bureau
	})
}

// TestShiftCardsIntegrity は同一データに対して2回実行し、結果が一致することを検証する。
// これにより、修正が冪等性を保っていること（同じ入力→同じ出力）を確認する。
// さらに、結果の構造的な妥当性（カード数、メンバー数、名前の一意性）もチェックする。
//
//	DEBUG_SQL=0 go test -run TestShiftCardsIntegrity -v ./lib/usecase/
func TestShiftCardsIntegrity(t *testing.T) {
	os.Setenv("DEBUG_SQL", "0")

	client := newTestClient(t)
	defer client.CloseDB()

	uc := newTestShiftUseCase(t, client)
	targetUserID := seedBenchData(t, client.DB())

	userIDStr := fmt.Sprintf("%d", targetUserID)
	dateIDStr := fmt.Sprintf("%d", benchDateID)
	weatherIDStr := fmt.Sprintf("%d", benchWeatherID)
	ctx := context.Background()

	// 1回目の実行
	cards1, err := uc.GetShiftCardsByUserAndDateAndWeather(ctx, userIDStr, dateIDStr, weatherIDStr)
	if err != nil {
		t.Fatalf("1st call: %v", err)
	}

	// 2回目の実行（同じデータに対して）
	cards2, err := uc.GetShiftCardsByUserAndDateAndWeather(ctx, userIDStr, dateIDStr, weatherIDStr)
	if err != nil {
		t.Fatalf("2nd call: %v", err)
	}

	json1, err := canonicalJSON(cards1)
	if err != nil {
		t.Fatalf("canonicalJSON 1st: %v", err)
	}
	json2, err := canonicalJSON(cards2)
	if err != nil {
		t.Fatalf("canonicalJSON 2nd: %v", err)
	}

	// 冪等性チェック: 2回の実行結果が同一
	if string(json1) != string(json2) {
		os.MkdirAll("testdata", 0755)
		os.WriteFile("testdata/run1.json", json1, 0644)
		os.WriteFile("testdata/run2.json", json2, 0644)
		t.Fatal("Idempotency check FAILED: two calls with same data produced different results")
	}
	t.Log("Idempotency: PASS (two calls produce identical results)")

	// 構造チェック
	if len(cards1) == 0 {
		t.Fatal("No shift cards returned")
	}
	t.Logf("Cards returned: %d", len(cards1))

	for i, card := range cards1 {
		if card.TaskName == "" {
			t.Errorf("Card %d: empty task name", i)
		}
		if card.StartTime == "" || card.EndTime == "" {
			t.Errorf("Card %d: empty start/end time", i)
		}
		if len(card.ShiftMembers) == 0 {
			t.Errorf("Card %d: no shift members", i)
		}

		for j, sm := range card.ShiftMembers {
			if len(sm.Members) == 0 {
				t.Errorf("Card %d, TimeSlot %d: no members", i, j)
				continue
			}

			// メンバーの一意性チェック（同一メンバーが重複していないか）
			seen := make(map[string]bool)
			for _, m := range sm.Members {
				if m.Name == "" {
					t.Errorf("Card %d, TimeSlot %d: member with empty name", i, j)
				}
				if m.Grade == "" {
					t.Errorf("Card %d, TimeSlot %d: member %s has empty grade", i, j, m.Name)
				}
				if m.Bureau == "" {
					t.Errorf("Card %d, TimeSlot %d: member %s has empty bureau", i, j, m.Name)
				}
				key := m.Name + "|" + m.Grade + "|" + m.Bureau
				if seen[key] {
					t.Errorf("Card %d, TimeSlot %d: duplicate member %s", i, j, m.Name)
				}
				seen[key] = true
			}

			t.Logf("Card %d, TimeSlot %s-%s: %d unique members", i, sm.STime, sm.ETime, len(sm.Members))
		}
	}

	// ゴールデンファイルとして保存（将来のブランチ間比較用）
	os.MkdirAll("testdata", 0755)
	os.WriteFile("testdata/shift_cards_golden.json", json1, 0644)
	t.Logf("Golden file saved: testdata/shift_cards_golden.json (%d bytes)", len(json1))
}
