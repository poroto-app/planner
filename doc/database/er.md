### User

```mermaid
---
title: user
---
erDiagram
    USER {
        string id
        string firebase_uid
        string name
        string email
        string photo_url
    }
```

### Place

```mermaid
---
title: place
---
erDiagram
    places {
        char(36) id PK
        string name
    }

    google_places {
        string google_place_id PK
        char(36) place_id FK
        string name
        string formatted_address
        string vicinity
        int price_level
        float rating
        int user_ratings_total
        double latitude
        double longitude
        point location
    }

    google_place_types {
        char(36) id PK
        string google_place_id
        string type
        int order
    }

    google_place_photo_references {
        string photo_reference PK
        string google_place_id FK
        int width
        int height
    }

    google_place_photos {
        char(36) id PK
        string google_place_id FK
        string photo_reference FK
        string url
        int width
        int height
    }

    google_place_photo_attributions {
        char(36) id PK
        char(36) google_place_id FK
        string photo_reference FK
        string html_attribution
    }

    google_place_reviews {
        char(36) id PK
        string google_place_id FK
        string author_name
        string author_url
        string author_profile_photo_url
        string language
        int rating
        string text
        int time
    }

    google_place_opening_periods {
        char(36) id PK
        string google_place_id FK
        int open_day
        int open_time
        int close_day
        int close_time
    }

    places ||..|| google_places: "1:1"
    google_places ||..o{ google_place_types: "1:N"
    google_places ||..o{ google_place_photo_references: "1:N"
    google_places ||..o{ google_place_opening_periods: "1:N"
    google_place_photo_references ||..o{ google_place_photo_attributions: "1:N"
    google_place_photo_references ||..o{ google_place_photos: "1:N"
    google_places ||..o{ google_place_reviews: "1:N"
```

- types
    - 並び替えが発生しないため、単純な`order`カラムで順番を管理

### Place Photos

```mermaid
---
title: place photos
---
erDiagram
	place_photos {
		string id PK
		char(36) place_id FK
		string user_id FK
		string url
		int width
		int height
	}

	place ||..o{ place_photos: "1:N"
	user ||..o{ place_photos: "1:N"
```

### Plan Candidate

```mermaid
---
title: plan_candidate
---
erDiagram
    plan_candidate_sets {
        char(36) id PK
        timestamp expires_at
    }

    plan_candidate_set_meta_data {
        char(36) id PK
        char(36) plan_candidate_set_id FK
        double latitude_start
        double longitude_start
        int plan_duration_minutes
        bool is_created_from_current_location
    }

    plan_candidate_set_searched_places {
        char(36) id PK
        char(36) plan_candidate_set_id FK
        char(36) place_id FK
    }

    plan_candidate_set__meta_data_categories {
        char(36) id PK
        char(36) plan_candidate_set_id FK
        string category
        bool is_selected
    }

    plan_candidates {
        char(36) id PK
        char(36) plan_candidate_set_id FK
        VARCHAR(255) name
        int sort_order
    }

    plan_candidate_places {
        char(36) id PK
        char(36) plan_candidate_id FK
        char(36) plan_candidate_set_id FK
        char(36) place_id FK
        int sort_order
    }

    plan_candidate_sets ||..o{ plan_candidates: "1:N"
    plan_candidate_sets ||..|| plan_candidate_set_meta_data: "1:1"
    plan_candidate_sets ||..o{ plan_candidate_set_categories: "1:N"
    plan_candidate_sets ||..o{ plan_candidate_set_searched_places: "1:N"
    plan_candidates ||..o{ plan_candidate_places: "1:N"
    plan_candidate_places ||..|| places: "1:1"
    plan_candidate_set_searched_places ||..|| places: "1:1"
```

### Plan

```mermaid
---
title: plan
---
erDiagram
    plans {
        char(36) id PK
        char(36) user_id FK
        string name
        point location "プランの大まかな場所"
    }

    plan_places {
        char(36) id PK
        char(36) plan_id FK
        char(36) place_id FK
        int sort_order
    }

    plans ||..o{ plan_places: "1:N"
    plans ||..|| users: "1:1"
    plan_places ||..|| places: "1:1"
```

### Like Place

```mermaid
---
title: like_place
---
erDiagram
    plan_candidate_set_like_places {
        char(36) id PK
        char(36) plan_candidate_set_id FK "UNIQUE(plan_candidate_set_id, place_id)"
        char(36) place_id FK "UNIQUE(plan_candidate_set_id, place_id)"
    }
    
    user_like_places {
        char(36) id PK
        char(36) user_id FK "UNIQUE(user_id, place_id)"
        char(36) place_id FK "UNIQUE(user_id, place_id)"
    }

    plan_candidate_set_like_places }o..|| plan_candidate_sets: "N:1"
    plan_candidate_set_like_places }o..|| places: "N:1"
    user_like_places o|..|| places: "N:1"
    user_like_places o|..|| users: "N:1"
```