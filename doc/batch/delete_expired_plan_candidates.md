# 有効期限切れプラン削除バッチ

## 概要

| 項目 | 内容 |
|------|------|
| バッチ名 | delete_expired_plan_candidates |
| 実行場所 | TODO |
| 実行時間 | TODO |
| 実行間隔 | TODO |

## 実行目的

有効期限切れのプラン候補を削除することにより、不要なデータを削除する。

## 影響範囲

- Firestore
  - `plan_candidates` コレクション
  - `plan_search_results` コレクション