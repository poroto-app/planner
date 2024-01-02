# Planner

指定された場所を巡るプランを作成するAPI

## 環境構築

### goのインストール
- 特定のバージョンのgoを使用するために[goenv](https://github.com/syndbg/goenv)の利用します
- [インストール方法はこちらを参考にしてください](https://github.com/syndbg/goenv/blob/master/INSTALL.md)

```shell
goenv install 1.19.6
```
- バージョンを指定
```shell
goenv global 1.19.6
```
- バージョンを確認
```shell
goenv version
# 1.19.6 (set by /your/path/to/planner/.go-version)

go version
# go version go1.19.6 linux/amd64
```

### ライブラリのインストール
```shell
go mod tidy
```

### MySQLの起動（Docker）
```shell
docker compose up -d
```

### IntelliJ IDEAの設定
1. `go env`を実行し、`GOROOT`を取得する
2. `Languages & Frameworks` → `GO`→ `GOROOT` を開く
3. `GOROOT`を入力する

### シークレットの復元

`plannnr API`で使用するシークレットは[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)で管理されています

1. gcloudコマンドをインストール
  参照： [gcloud CLI をインストールする](/https://cloud.google.com/sdk/docs/install)
2. Google Cloud porotoプロジェクトを操作する権利があるアカウントにログイン
   ```sh
   gcloud auth login
   ```
3. 復号化スクリプトを実行
    ```sh
    scripts/decrypt.sh
    ```
### シークレット（.env.local等）変更時
暗号化し、[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)で管理してください

1. [poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)を以下の場所にclone
    ```sh
    - your_dir_of_poroto
      - planner
      - infrastructure
    ```
2. Google Cloud porotoプロジェクトを操作する権利があるアカウントにログイン
   ```sh
   gcloud auth login
   ```
3. 暗号化スクリプトを実行
    ```sh
    scripts/encrypt.sh
    ```
4. infrastructureリポジトリの変更をコミット

## 開発思想

### ディレクトリ構成

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)に従います。

### 設計思想

- ドメイン駆動設計

### Linting

- linter に golangci-lint を使用しています
- 以下を実行して、警告があれば対応してください
    - see: https://golangci-lint.run/usage/install/#docker

```shell
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.46.2 golangci-lint run -v
```

## GraphQL
### コード生成
```shell
go generate ./...
```

## テストの実行
```shell
# sqlboiler.tomlをコピーして、テスト用のDBに接続するように変更
cp db/sqlboiler.toml sqlboiler.toml && sed -i 's/3306/3307/' sqlboiler.toml
# テストを実行
go test ./...
```

## Database
### Gooseのインストール
https://pressly.github.io/goose/installation/
```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### マイグレーションの作成
```shell
goose -dir db/migrations create create_user_table sql go
```

### マイグレーションの実行
```shell
DB_USER=root \
DB_PASSWORD=password \
DB_HOST=localhost:3306 \
DB_NAME=poroto \
goose -dir db/migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST)/$DB_NAME?parseTime=true&loc=Asia%2FTokyo" up
```

### SQLBoilerをインストール
[volatiletech/sqlboiler #Download](https://github.com/volatiletech/sqlboiler?tab=readme-ov-file#download)
```shell
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
````

### SQLBoilerでモデルコードの生成
```shell
sqlboiler mysql -c db/sqlboiler.toml 
```

## Trouble Shooting
### MySQLをアップグレード・ダウングレードしたら起動できなくなった
※ 本番環境ではデータを移行することが必要です

ローカルでデータを削除する場合は以下のコマンドを利用します
```shell
docker compose down
docker volume rm planner_mysql-data
```