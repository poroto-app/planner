# Create Plan

### 現在地からプランを作成
```mermaid
sequenceDiagram
    poroto->>planner: session_idを用いてプラン作成セッションを管理
    planner->>DB: session_idを保存
    
    poroto->>planner:　現在地を送信
    planner->>DB: 現在地をsession_idと紐付けて保存
    # この状態で現在地が保存されるため、リロードしても問題ない
    
    planner->>Places API: 現在地周辺の場所を取得
    planner->>DB: 取得した場所一覧の情報を保存
    planner->>poroto: 現在地に基づいて場所のカテゴリを提案
    
    poroto->>planner: 選択したカテゴリや時間をもとにプランを作成
    planner->>DB: カテゴリや時間を保存
    planner->>DB: 作成したプランを保存
    
    planner->>poroto: 作成したプランを提示
```

### 場所を指定してプランを作成
```mermaid
sequenceDiagram
    poroto->>planner: session_idを用いてプラン作成セッションを管理
    planner->>DB: session_idを保存
    
    poroto->>planner:　指定した場所を送信
    planner->>DB: 位置情報をsession_idと紐付けて保存
    # この状態で現在地が保存されるため、リロードしても問題ない
    
    planner->>Places API: 現在地周辺の場所を取得
    planner->>DB: 取得した場所一覧の情報を保存
    planner->>poroto: 位置情報に基づいて場所のカテゴリを提案
    
    poroto->>planner: 選択したカテゴリや時間をもとにプランを作成
    planner->>DB: カテゴリや時間を保存
    planner->>DB: 作成したプランを保存
    
    planner->>poroto: 作成したプランを提示
```

## 例外対応
### `/plans/interest`でリロードされる
| 発生する問題                       | 対応方法         |
|------------------------------|--------------|
| ユーザーが指定した位置情報をクライアント側で保持できない | `session_id`をもとに情報を復元 |
| 同じ位置情報に基づいて重複して周辺の場所を検索してしまう | 検索結果をキャッシュする |

