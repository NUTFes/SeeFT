# 0018: テストでは db.Client に go-sqlmock の偽の DB を差し込み、repository の戻り値は変えない

- 状態：採用
- 決めた日：2026-07-07（PR #392 をマージした日）
- 決めた人：45th の PM
- 決定の信頼度：記録なし
  - 根拠：#391・PR #392・`docs/development/test-roadmap.md` には、決めたことと、案 B を選ばなかった理由は書かれているが、どれだけ確かだと思っていたかは書かれていない
- 出典：#391、PR #392、`docs/development/test-roadmap.md`（「前提: コード構造がテスト戦略を規定する」）

## 背景

2026-07-06 の MT で、SeeFT 本体は保守性を上げる方向に決まり、直す前にテストを書くことになった（[0003](0003-maintenance-over-features.md)）。

api の repository は、ほとんどが `database/sql` の `*sql.Rows` や `*sql.Row` を返し、値の読み取り（Scan）は usecase の側で行う作りになっている。`*sql.Row` はテストのコードから作れない。そのため、repository のインターフェースを偽物に差し替えても、usecase のテストが書けなかった。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| A：DB への接続（`db.Client`）に、go-sqlmock の偽の `*sql.DB` を差し込む | 今のコードを直さなくてよい。repository の実装ごとテストできる | テストに SQL の文字列を書く必要がある |
| B：repository の戻り値を entity に変える | テストが書きやすくなる | 141 のメソッドに手を入れる大きな作業になる |

## 決定

案 A をとる。テストでは、`db.Client` に go-sqlmock の偽の `*sql.DB` を差し込む。repository の戻り値は変えない。

## 理由

案 B は、141 のメソッドに手を入れる大きな作業で、保守に寄せる方針（[0003](0003-maintenance-over-features.md)）で手を付けられる大きさを超えていた。案 A なら、今のコードを変えずにテストを書き始められる。

## 前提

- 差し込む先のインターフェース：`api/lib/externals/db/db.go#Client`
- テストで使う偽物：`api/lib/internals/repository/sqlmock_helper_test.go#fakeDBClient`
- テストの進め方：`docs/development/test-roadmap.md`
- repository が `*sql.Rows` や `*sql.Row` を返す作りであること。戻り値を entity に変える作業をするなら、この ADR を見直す

## 結果

2026-09-30 の時点で、`api/` のテストのファイルは 18 本ある。そのうち 9 本が go-sqlmock を使うテストで（ファイル名が `_sqlmock_test.go` のもの。repository と usecase の両方）、ほかに偽物を作る補助のファイルが1本ある。

## 追記
