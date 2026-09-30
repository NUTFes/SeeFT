# 0006: マニュアルは SeeFT の API から、nutfes の Google アカウントに限って配信する

- 状態：採用
- 決めた日：2026-08-20（GAS での配信は 2026-08-19、Cloudflare Access は 2026-08-20 に見送った）
- 決めた人：45th の PM
- 確信度：記録なし
- 出典：#444、#445、#448、#480

## 背景

45th では、タスクごとのマニュアルを、読みやすい HTML に変換して配っていた。8月の途中までは、生成した HTML をリンクの共有で配っていた。

2026-08-20 に、情報局長から「アプリの中で、閲覧してよい人だけに配信してほしい」と正式な要望が出た。マニュアルには技大祭の内部の段取りが書かれているので、技大祭の関係者だけに見せたい。使い方の説明会（8/28）までに本番に出す必要があった。

技大祭の関係者は、全員が `○○.nutfes@gmail.com` の Google アカウントを持っている。マニュアルを見せる相手を絞るための認証は、別に用意する必要があった。

## 候補

| 候補 | 良い点 | 悪い点 |
| --- | --- | --- |
| リンクの共有で配り続ける | 手間がない | リンクを知っている人なら誰でも読める |
| GitHub Pages で公開する | 置くのが簡単 | 公開の範囲を絞れない |
| GAS の Web アプリで配信する（Google のログインで絞る） | 試作は動いた（8/17） | 複数の Google アカウントでログインしているブラウザやスマホでは、GAS のコードが動く前に Google 側で 404 になる。Google 側の長年の不具合で、デプロイを作り直しても避けられなかった。アプリの中に埋め込むこともできない |
| Cloudflare Access を前に置く | 認証を外のサービスに任せられる | 外部の契約が要り、費用が利用する人数（300 人超）に比例する |
| SeeFT の API に、Google のログインで絞る入口を足す | 新しいサーバーや契約が要らない。既存の API のホストに道を足すだけで済む | 自前で実装・保守する必要がある |

## 決定

SeeFT の API に、Google のログイン（OAuth の認可コードフロー）を足す。`○○.nutfes@gmail.com` のアカウント（と、例外として許可したアドレス）でログインした人にだけ、`GET /manuals/:id` でマニュアルを見せる。

- Google に求める権限は、メールアドレスを知るための最小限（`openid` と `email`）だけにする
- ログインに要る設定（環境変数）が欠けているときは、`/manuals` の道ごと出さない。設定の漏れで誰でも読める状態にならないようにする
- 生成した HTML は、トークンを持つ人が `PUT /manuals/:id` で置く。本番のサーバーに SSH で入らなくても差し替えられる（#448）

## 理由

GAS での配信は、複数の Google アカウントを使う人が多い環境で確実に 404 になり、実機で避けられなかった。Cloudflare Access は、300 人を超える委員に使わせると費用がかかる。

SeeFT の API に入口を足すなら、新しいサーバーや契約は要らず、説明会までに間に合った。見せる相手を「SeeFT のアカウント」ではなく「nutfes の Google アカウント」で区切ったのは、委員が全員すでに持っているアカウントで、新しく配るものが要らないからである。

Google に求める権限を最小限にしたのは、この範囲なら、Google の審査を受けていないアプリでもログインできる人数に上限が掛からない、と判断したためである（Google の方針による。コードからは確かめられない）。

## 前提

- 閲覧してよいかの判定：`api/lib/usecase/manual_usecase.go#manualUseCase.IsAllowed`、許可するアドレスの形：`api/lib/usecase/manual_usecase.go#manualAllowedDomainRe`
- ログインの入口：`api/lib/usecase/manual_usecase.go#manualUseCase.BuildAuthURL`
- 閲覧と配置：`api/lib/internals/controller/manual_controller.go#manualController.ShowManual`、`api/lib/internals/controller/manual_controller.go#manualController.UploadManual`
- 道を出すかの判断：`api/lib/router/router.go#router.ProvideRouter`、組み立ては `api/lib/di/di.go#InitializeServer`
- 技大祭の関係者が `○○.nutfes@gmail.com` のアカウントを持つこと。この形が変わったら、許可するアドレスの形を直す
- ログインに使う Google の OAuth クライアントを持つアカウントが残っていること。持ち主のアカウントが消えると、マニュアルが開けなくなる（#557）

## 結果

2026-08-25 に #445 をマージし、45th の技大祭の期間中は、この形でマニュアルを配信した。

## 追記
