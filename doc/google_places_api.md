# Google Places API

### 料金体系

https://developers.google.com/maps/documentation/places/web-service/usage-and-billing?hl=ja

- [Places APIはリクエストに含まれるフィールド内で最上位のSKUに基づいて課金される](https://developers.google.com/maps/documentation/places/web-service/usage-and-billing?hl=ja)

### Place Detail APIの呼び出し方

- プラン作成時には複数回、Place Detailによる情報が必要な場面がある（開店時刻、写真、レビュー等）
- これらのリクエストを別々に行ってしまうと、`SKU料金*リクエスト回数`の料金が発生してしまう
- したがって、planner では、Place Detailによる情報を一度に取得するようにしている

### リクエストのフロー

```mermaid
sequenceDiagram
    participant P as planner
    participant NS as Nearby Search
    participant PD as Place Details
    P ->> NS: 周辺検索（10カテゴリ検索）
    NS ->> P: 周辺の場所（60件程度）
    alt ユーザーが場所を指定してプランを作成した場合
        P ->> PD: ユーザーが指定した場所を検索
        PD ->> P: 場所の詳細情報
    end
    P ->> PD: 場所の詳細情報（プランに含まれる12件程度）
    PD ->> P: 場所の詳細情報
```