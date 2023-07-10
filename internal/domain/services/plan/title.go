package plan

import (
	"fmt"
	"strings"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
)

// GeneratePlanTitle プランのタイトルを生成する
// タイトルが生成できなかった場合は、nilを返す
func (s PlanService) GeneratePlanTitle(places []models.Place) (*string, error) {
	placeNames := make([]string, len(places))
	for i, place := range places {
		placeNames[i] = fmt.Sprintf("%s(%s)", place.Name, place.Category)
	}

	nGenerate := 3
	response, err := s.openaiChatCompletionClient.Complete(openai.ChatCompletionRequest{
		Model: openai.ModelGPT3Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: "あなたはコピーライトを生成するアシスタントです" +
					"例：相模原図書館（図書館）とスターバックスコーヒー（カフェ）を含むプラン" +
					"生成するコピーライト：新しい本を買って、カフェでゆっくり読書しませんか" +
					"要件：体験を想像させ、一目引くタイトルであること" +
					"最大文字数: 20文字",
			},
			{
				Role:    "system",
				Content: fmt.Sprintf("%sを含むプラン", strings.Join(placeNames, "と")),
			},
		},
		N: &nGenerate,
	})

	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("response.Choices is empty")
	}

	choices := filterByMessageLength(response.Choices, 15)
	if len(choices) == 0 {
		return nil, fmt.Errorf("response.Choices is empty")
	}

	title := choices[indexOfMaxMessageLength(choices)].Content
	replaceCharacters := []string{"\n", "「", "」", "'", "’", "\"", "”", "：", ":"}
	for _, character := range replaceCharacters {
		title = strings.ReplaceAll(title, character, "")
	}

	return &title, nil
}

func filterByMessageLength(messages []openai.ChatCompletionChoice, length int) []openai.ChatCompletionMessage {
	filteredMessages := make([]openai.ChatCompletionMessage, 0)
	for _, message := range messages {
		if len(message.Message.Content) <= length {
			filteredMessages = append(filteredMessages, message.Message)
		}
	}
	return filteredMessages
}

func indexOfMaxMessageLength(messages []openai.ChatCompletionMessage) int {
	maxLength := 0
	indexOfMaxLength := 0
	for i, message := range messages {
		if len(message.Content) > maxLength {
			maxLength = len(message.Content)
			indexOfMaxLength = i
		}
	}
	return indexOfMaxLength
}
