package plangen

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"poroto.app/poroto/planner/internal/domain/models"
	"poroto.app/poroto/planner/internal/infrastructure/api/openai"
)

// GeneratePlanTitle プランのタイトルを生成する
// タイトルが生成できなかった場合は、nilを返す
func (s Service) GeneratePlanTitle(places []models.Place) (*string, error) {
	placeNames := make([]string, len(places))
	for i, place := range places {
		var categoryNames []string
		for _, category := range models.GetCategoriesFromSubCategories(place.Google.Types) {
			categoryNames = append(categoryNames, category.Name)
		}

		placeNames[i] = fmt.Sprintf("%s(%s)", place.Google.Name, strings.Join(categoryNames, ","))
	}

	nGenerate := 5
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

	choices := make([]openai.ChatCompletionMessage, len(response.Choices))
	for i, choice := range response.Choices {
		choices[i] = choice.Message
	}

	choices = replaceMessageContent(choices)

	choices = filterByMessageLength(choices, 30)
	if len(choices) == 0 {
		return nil, fmt.Errorf("response.Choices is empty")
	}

	title := choices[indexOfMaxMessageLength(choices)].Content

	return &title, nil
}

// replaceMessageContent メッセージの内容から不要な文字を削除する
func replaceMessageContent(choices []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	deleteCharacters := []string{"\n", "「", "」", "'", "’", "\"", "”", "：", ":"}
	for i, choice := range choices {
		content := choice.Content
		for _, deleteCharacter := range deleteCharacters {
			content = strings.ReplaceAll(content, deleteCharacter, "")
		}
		choices[i].Content = content
	}

	return choices
}

func filterByMessageLength(choices []openai.ChatCompletionMessage, length int) []openai.ChatCompletionMessage {
	filteredChoices := make([]openai.ChatCompletionMessage, 0)
	for _, choice := range choices {
		// MEMO: len(string) で得られるのはバイト数であり、文字数ではない
		messageLength := utf8.RuneCountInString(choice.Content)
		if messageLength <= length {
			filteredChoices = append(filteredChoices, choice)
		}
	}
	return filteredChoices
}

func indexOfMaxMessageLength(choices []openai.ChatCompletionMessage) int {
	maxLength := 0
	indexOfMaxLength := 0
	for i, choice := range choices {
		// MEMO: len(string) で得られるのはバイト数であり、文字数ではない
		messageLength := utf8.RuneCountInString(choice.Content)
		if messageLength > maxLength {
			maxLength = messageLength
			indexOfMaxLength = i
		}
	}
	return indexOfMaxLength
}
