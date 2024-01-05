# Firestore Schema

## Place

```mermaid
erDiagram
    Place {
        id string
        googlePlaceId string
    }
    GooglePlace {
        googlePlaceId string
    }
    GooglePlacePhoto {
        googlePlaceId string
        photoReference string
        url *string
    }
    GooglePlaceReview {
        googlePlaceId string
    }

    Place ||..|{ GooglePlace: has
    Place ||..o{ GooglePlacePhoto: has
    Place ||..o{ GooglePlaceReview: has
```

## PlanCandidate

```mermaid
erDiagram
    PlanCandidate {
        placeIds string[]
    }
    Plan {
        placeIds string[]
    }
    Place {
        id string
        googlePlaceId string
    }

    PlanCandidate ||..|{ Place: related
    PlanCandidate ||..|{ Plan: has
```

PlanCandidateはプランの情報と、プランに含まれる場所の参照を持つ

### Google Places APIで取得したデータ

`Google Places API`
から取得された場所のデータは`Places`, `GooglePlace`、`GooglePlacePhotos`、`GooglePlaceReviews`に保存される。  