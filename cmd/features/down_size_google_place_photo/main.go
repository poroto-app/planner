package main

import (
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"poroto.app/poroto/planner/internal/env"
	"poroto.app/poroto/planner/internal/infrastructure/rdb"
	"poroto.app/poroto/planner/internal/infrastructure/rdb/generated"
	"regexp"
	"strings"
)

func main() {
	env.LoadEnv()

	db, err := rdb.InitDB(false)
	if err != nil {
		log.Fatalf("error while initializing db: %v", err)
	}

	ctx := context.Background()
	googlePlacePhotoSlice, err := generated.GooglePlacePhotos().All(ctx, db)
	if err != nil {
		log.Fatalf("error while getting google place photos: %v", err)
	}

	for _, googlePlacePhoto := range googlePlacePhotoSlice {
		newUrl, err := rewriteUrl(googlePlacePhoto.URL, 500, 500)
		if err != nil {
			log.Printf("error while rewriting URL: %v", err)
			continue
		}

		googlePlacePhoto.URL = *newUrl
		if _, err := googlePlacePhoto.Update(ctx, db, boil.Infer()); err != nil {
			log.Printf("error while updating google place photo: %v", err)
			continue
		}
	}
}

func rewriteUrl(url string, width int, height int) (*string, error) {
	// URLの最後の部分を取得
	parts := strings.Split(url, "=")
	if len(parts) < 2 {
		fmt.Println("Invalid URL format")
		return nil, fmt.Errorf("invalid URL format")
	}
	lastPart := parts[len(parts)-1]

	// 正規表現パターンを定義
	re := regexp.MustCompile(`w\d+-h\d+`)

	// 置換文字列を定義
	newSuffix := fmt.Sprintf("w%d-h%d", width, height)

	// URLの最後の部分の書き換え
	updatedLastPart := re.ReplaceAllString(lastPart, newSuffix)

	// 更新されたURLを組み立て
	updatedURL := strings.Join(parts[:len(parts)-1], "=") + "=" + updatedLastPart

	return &updatedURL, nil
}
