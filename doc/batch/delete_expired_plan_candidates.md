# 有効期限切れプラン削除バッチ

## 概要

| 項目 | 内容                                                                                                                                              |
|------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| バッチ名 | delete_expired_plan_candidates                                                                                                                  |
| 実行場所 | Cloud Functions(Github Actionsによりデプロイ)                                                                                                          |
| 実行時間 | [4:00AM(JST)](https://github.com/poroto-app/infrastructure/blob/0dc06438fc35f6c503d04e9bd963a8cc20b1400d/terraform/development/scheduler.tf#L5) |
| 実行間隔 | [毎日](https://github.com/poroto-app/infrastructure/blob/0dc06438fc35f6c503d04e9bd963a8cc20b1400d/terraform/development/scheduler.tf#L5)                                                                                                                                          |

## 実行目的

有効期限切れのプラン候補を削除することにより、不要なデータを削除する。

## 影響範囲

- Firestore
  - `plan_candidates` コレクション

## ユーザー影響

### プラン候補を表示できなくなる

削除されたプラン候補は表示できなくなる。  
ユーザーがプラン候補のURLを持っていた場合や共有していた場合には、404エラーが表示される。  
[削除までには１週間の猶予](https://github.com/poroto-app/planner/blob/develop/internal/domain/services/plan/plan_candidate.go#L16)があるため、この影響は許容される。
