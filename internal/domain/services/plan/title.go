package plan

import (
	"fmt"
	"strings"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
)

// GeneratePlanTitle プランのタイトルを生成する
// タイトルが生成できなかった場合は、nilを返す
func (s PlanService) GeneratePlanTitle(places []models.Place) (title *string, err error) {
	placeNames := make([]string, len(places))
	for i, place := range places {
		// TODO: 場所の名前だけでなく、その場所の特徴も含める
		// ex: スターバックスコーヒー（カフェ）
		placeNames[i] = place.Name
	}

	response, err := s.openaiChatCompletionClient.Complete([]openai.ChatCompletionMessage{
		{
			Role: "system",
			Content: "あなたはコピーライトを生成するアシスタントです" +
				"例：図書館とカフェを含むプラン" +
				"生成するコピーライト：新しい本を買って、カフェでゆっくり読書" +
				"要件：20文字以内の文章になっていること。　体験を想像させるような文章であること",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("%sを含むプラン", strings.Join(placeNames, "と")),
		},
	})

	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("response.Choices is empty")
	}

	return &response.Choices[0].Message.Content, nil
}
