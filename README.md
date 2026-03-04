■ todo-app(TODOリスト)

Go(標準ライブラリ`net/http`)+SQLiteで作ったシンプルなTODOリストです。
タスクの追加,削除,完了切り替えに加えて、タイトルと期限日のインライン編集ができます。

■ できること

- タスク追加(タイトル,期限日)
- タスク削除
- 完了/未完了の切り替え(チェックボックス)
- タイトルのインライン編集(クリック → 入力 → Enter/フォーカス外で保存)
- 期限日のインライン編集(クリック → date picker → Enter/フォーカス外で保存)

■ 使用技術

- Go(net/http, html/template)
- SQLite(github.com/mattn/go-sqlite3)
- HTML, CSS, Vanilla JavaScript(フレームワークなし)

■ 起動方法

1) 依存関係を取得

```bash
go mod tidy
```

2) サーバー起動

```bash
go run ./cmd/server
```

起動後、ブラウザでhttp://localhost:8080を開きます。

---

■ テスト実行

```bash
go test ./...
```

出力の見方(例):

- ? xxx [no test files] … そのパッケージ配下に *_test.go が無い(テスト未作成)
- ok xxx 0.743s … テストが実行され、すべて成功(所要時間も表示)
- ok xxx (cached) … 前回と同じ結果を再利用した(コードや依存が変わっていないので高速化のためにキャッシュを使用)

※ キャッシュを無視して毎回走らせたい場合は:

```bash
go test -count=1 ./...
```

■ Windowsの注意(sqlite3/cgo)

github.com/mattn/go-sqlite3はcgoを使うため、WindowsではCコンパイラ(例: gcc)が必要です。

- 例:MSYS2を入れてmingw-w64-ucrt-x86_64-gccを導入する
- gcc --versionで確認できます

■ アーキテクチャ(MVCを意識した層分割)

厳密なフレームワークではありませんが、MVCを意識した構成(Controller/Service/Repository)で実装しています。

- View:templates/(HTMLテンプレート)
- Controller:internal/controller/
  - HTTPリクエストを受け取る
  - 入力を読み取ってServiceを呼ぶ
  - レスポンス(リダイレクト/ステータス)を返す
- Service:internal/service/
  - アプリ側のルール(入力整形・簡易バリデーション)
  - Repositoryを呼んで永続化する
- Repository:internal/repository/
  - SQLを集約(DBへの読み書き)
- Model:internal/model/
  - Todo 構造体など

■ ディレクトリ構成

├─ cmd/server/                # エントリーポイント(main.go)
├─ internal/
│  ├─ controller/             # HTTP Handler(Controller)
│  ├─ service/                # ビジネスロジック(Service)
│  ├─ repository/             # DBアクセス(Repository)
│  ├─ model/                  # データ構造(Model)
│  └─ db/                     # DB初期化(接続・テーブル作成)
├─ templates/                 # View(HTMLテンプレート)
├─ static/
│  ├─ css/                    # CSS
│  └─ js/                     # JavaScript(fetch/DOM更新)
└─ todo.db                    # SQLite DB(自動生成/更新)

■ 処理の流れ(例:完了/未完了切り替え)

1. 画面でチェックボックスを操作
2. static/js/script.jsがfetch(/toggle)をPOST
3. Controller(Toggle)がリクエストを受け取る
4. Serviceが入力を整形してRepositoryを呼ぶ
5. RepositoryがUPDATE todos SET completed = ? を実行
6. 成功したらJS側でDOMを「未完了 ↔ 完了済み」へ移動して即反映

■ メモ

- CSS/JSを更新したのに反映されない場合、ブラウザキャッシュが原因のことがあります  
  - DevTools → Network → Disable cache(DevToolsを開いた状態でリロード)
  - もしくはCtrl+F5(ハードリロード)

■ GitHubに上げるときの注意

- todo.dbは実行すると更新されるので、リポジトリには入れず .gitignoreで除外するのがおすすめです。
