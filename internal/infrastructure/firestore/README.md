# Firestore Schema

## PlanCandidate
```mermaid
erDiagram
    PlanCandidate
    Plan {
    %% TODO: PlaceRepositoryに保存した値を取得するようにする
        Places Place[]
    }
    Place {
        id string
        googlePlaceId string
    }
    GooglePlaceApiSearchResult {
        googlePlaceId string
    }
    GooglePlaceAPiPhotos {
        googlePlaceId string
    }
    GooglePlaceApiReviews {
        googlePlaceId string
    }

    PlanCandidate ||..|{ GooglePlaceApiSearchResult: has
    PlanCandidate ||..|{ GooglePlaceAPiPhotos: has
    PlanCandidate ||..|{ GooglePlaceApiReviews: has
    PlanCandidate ||..|{ Place: has
    PlanCandidate ||..|{ Plan: has
```
PlanCandidateはプランとプランに含まれる場所の情報を持っている  

### Google Places APIで取得したデータ
`Google Places API`から取得された場所のデータは`GooglePlaceApiSearchResult`、`GooglePlaceAPiPhotos`、`GooglePlaceApiReviews`に保存される。  
`PlanCandidateId`がわかっていれば、すべてのデータが取得できるようにすると最も取得効率が良いため、写真やレビューのサブコレクションを`PlanCandidate`の直下においている  
（写真やレビューのサブコレクションを`GooglePlacesApiSearchResult`のそれぞれのドキュメントの中に作成すると、2段階で取得する必要がある）