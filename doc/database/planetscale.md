# PlanetScale

### 認証

```shell
pscale auth login
```

### ローカルからPlanetScaleに接続する

[PlanetScale CLI commands - connect](https://planetscale.com/docs/reference/connect)

```shell
pscale connect poroto branch_name --port 3309 --org poroto
```

### マイグレーションを実行

```shell
pscale connect poroto staging --port 3309 --org poroto
goose -dir db/migrations mysql "user:password@tcp(localhost:3309)/poroto?tls=true&interpolateParams=true&parseTime=true&loc=Asia%2FTokyo" up
```

### マイグレーションファイルの登録

- マイグレーションを適切に行うために、マイグレーション履歴が`goose_db_version`というテーブルで管理されている
- PlanetScaleではこのマイグレーション履歴を開発ブランチからproductionにコピーする機能がある。

1. Settings > General
2. `Automatically copy migration data`に以下を設定

   |項目名|値|
   |---|---|
   |Migration framework|Other|
   |Migration table name|goose_db_version|