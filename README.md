# Planner
指定された場所を巡るプランを作成するAPI

## 環境構築
### .env.local
- シークレットを含んだ.env.localは[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)で管理されています
- Google Cloud KMSで管理されている鍵を利用し、暗号化されている.env.localを複合します
- 事前に[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)をクローンしておいてください。
```shell
# .env.localを復号
gcloud kms decrypt \
  --location "asia-northeast1" \
  --keyring "planner_api_key_ring" \
  --key "planner_api_crypt_key" \
  --plaintext-file ./.env.local \
  --ciphertext-file ./your_cloned_path/infrastructure/role/app/files/development/planner/.env.local.enc
```

### .env.local変更時
- 暗号化を行い、[poroto-app/infrastructure](https://github.com/poroto-app/infrastructure)で管理してください。
```shell
gcloud kms decrypt \
  --location "asia-northeast1" \
  --keyring "planner_api_key_ring" \
  --key "planner_api_crypt_key" \
  --plaintext-file ./.env.local \
  --ciphertext-file ./your_cloned_path/infrastructure/role/app/files/development/planner/.env.local.enc
```
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
