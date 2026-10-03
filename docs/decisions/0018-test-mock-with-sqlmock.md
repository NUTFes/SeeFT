# 0018: テストではdb.Clientにgo-sqlmockの偽のDBを差し込み、repositoryの戻り値は変えない

- 状態：採用
- 決めた日：2026-07-07（PR #392をマージした日）
- 決めた人：45thのPM
- 決定の信頼度：記録なし
  - 根拠：#391・PR #392・`docs/development/test-roadmap.md`には、決めたことと、案Bを選ばなかった理由は書かれているが、どれだけ確かだと思っていたかは書かれていない
- 出典：#391、PR #392、`docs/development/test-roadmap.md`（「前提: コード構造がテスト戦略を規定する」）

## 背景

2026-07-06のMTで、SeeFT本体は保守性を上げる方向に決まり、直す前にテストを書くことになった（[0003](0003-maintenance-over-features.md)）。

apiのrepositoryは、ほとんどが`database/sql`の`*sql.Rows`や`*sql.Row`を返し、値の読み取り（Scan）はusecaseの側で行う作りになっている。`*sql.Row`はテストのコードから作れない。そのため、repositoryのインターフェースを偽物に差し替えても、usecaseのテストが書けなかった。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| A：DBへの接続（`db.Client`）に、go-sqlmockの偽の`*sql.DB`を差し込む | 今のコードを直さなくてよい。repositoryの実装ごとテストできる | テストにSQLの文字列を書く必要がある |
| B：repositoryの戻り値をentityに変える | テストが書きやすくなる | 141のメソッドに手を入れる大きな作業になる |

## 決定

案Aをとる。テストでは、`db.Client`にgo-sqlmockの偽の`*sql.DB`を差し込む。repositoryの戻り値は変えない。

## 理由

案Bは、141のメソッドに手を入れる大きな作業で、保守を優先する方針（[0003](0003-maintenance-over-features.md)）の中で手を付けられる大きさを超えていた。案Aなら、今のコードを変えずにテストを書き始められる。

## 前提

- 差し込む先のインターフェースは`api/lib/externals/db/db.go#Client`である
- テストで使う偽物は`api/lib/internals/repository/sqlmock_helper_test.go#fakeDBClient`である
- テストの進め方は`docs/development/test-roadmap.md`にある
- repositoryが`*sql.Rows`や`*sql.Row`を返す作りであること。戻り値をentityに変える作業をするなら、このADRを見直す

## 結果

2026-09-30の時点で、`api/`のテストのファイルは18本ある。そのうち9本がgo-sqlmockを使うテストで、ほかに偽物を作る補助のファイルが1本ある。9本はファイル名が`_sqlmock_test.go`のもので、repositoryとusecaseの両方にある。

## 追記
