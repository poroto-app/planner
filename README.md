# Planner

指定された場所を巡るプランを作成するAPI

## 環境構築

### goのインストール
- 特定のバージョンのgoを使用するために[goenv](https://github.com/syndbg/goenv)の利用します
- [インストール方法はこちらを参考にしてください](https://github.com/syndbg/goenv/blob/master/INSTALL.md)

```shell
goenv install 1.22.0
```
- バージョンを確認
```shell
goenv version
# 1.22.0 (set by /your/path/to/planner/.go-version)

go version
# go version go1.22.0 linux/amd64
```

### ライブラリのインストール
```shell
go mod tidy
```

### MySQLの起動（Docker）
```shell
cd docker
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
goose -dir db/migrations create <your migration name> sql
```

### マイグレーションの実行
```shell
DB_USER=root \
DB_PASSWORD=password \
DB_HOST=localhost \
DB_PORT=3306 \
DB_NAME=poroto \
goose -dir db/migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true&loc=Asia%2FTokyo" up
```

### SQLBoilerをインストール
[volatiletech/sqlboiler #Download](https://github.com/volatiletech/sqlboiler?tab=readme-ov-file#download)
```shell
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
````

### SQLBoiler Extentions(追加済み)
デフォルトのSQLBoilerにはBulk操作を行うための関数が含まれていないため、拡張テンプレートを追加する必要があります。  
[tiendc / sqlboiler-extensions](https://github.com/tiendc/sqlboiler-extensions)
```shell
git submodule add --name "db/extensions"  https://github.com/tiendc/sqlboiler-extensions.git "db/extensions"
git submodule update --init
```

### SQLBoilerでモデルコードの生成
```shell
cp db/sqlboiler_template.toml sqlboiler.toml && sed -i -e "s|\${GOPATH}|$(go env GOPATH)|g" sqlboiler.toml && sed -i -e "s|<sqlboiler-version>|$(grep "github.com/volatiletech/sqlboiler/v4" go.mod | awk '{print $2}')|g" sqlboiler.toml
sqlboiler mysql 
```

## Trouble Shooting
### MySQLをアップグレード・ダウングレードしたら起動できなくなった
※ 本番環境ではデータを移行することが必要です

ローカルでデータを削除する場合は以下のコマンドを利用します
```shell
docker compose down
docker volume rm docker_mysql-data
```

### `go version`と`goenv version`の結果が違う
以下の内容を`~/.zprofile`, `~/.bash_profile`, `~/.config/fish/config.fish`に含める
```sh
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
export PATH="$GOROOT/bin:$PATH"
export PATH="$PATH:$GOPATH/bin"
```