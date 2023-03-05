# Planner
指定された場所を巡るプランを作成するAPI

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
