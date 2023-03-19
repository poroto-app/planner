# Planner

指定された場所を巡るプランを作成するAPI

## 環境構築

### goのインストール
- 特定のバージョンのgoを使用するために[goenv](https://github.com/syndbg/goenv)の利用します
- [インストール方法はこちらを参考にしてください](https://github.com/syndbg/goenv/blob/master/INSTALL.md)

```shell
goenv install 1.19.6

goenv version
# 1.19.6 (set by /your/path/to/planner/.go-version)
```

### シークレットの復元

- `plannner API`で使用するシークレットは[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)で管理されています
- `scipts/decrypt.sh`
  を実行することで復元できます（※ [事前に gcloud をインストールする必要があります](https://cloud.google.com/sdk/docs/install?hl=ja)）

### シークレット（.env.local等）変更時

- 暗号化し、[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)で管理してください
- `scripts/encrypt.sh` を実行することで暗号化できます

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
