package plan

import (
	"fmt"

	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
)

// GeneratePlanTitle プランのタイトルを生成する
// タイトルが生成できなかった場合は、nilを返す
func (s PlanService) GeneratePlanTitle() (title *string, err error) {
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
			Content: "ラーメン屋とゲームセンターを含むプラン",
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
