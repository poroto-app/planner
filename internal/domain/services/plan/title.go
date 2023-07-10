package plan

import (
	"fmt"
	"strings"
	"unicode/utf8"

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

	choices := replaceMessageContent(response.Choices)

	choices = filterByMessageLength(choices, 30)
	if len(choices) == 0 {
		return nil, fmt.Errorf("response.Choices is empty")
	}

	title := choices[indexOfMaxMessageLength(choices)].Message.Content

	return &title, nil
}

// replaceMessageContent メッセージの内容から不要な文字を削除する
func replaceMessageContent(choices []openai.ChatCompletionChoice) []openai.ChatCompletionChoice {
	deleteCharacters := []string{"\n", "「", "」", "'", "’", "\"", "”", "：", ":"}
	for i, choice := range choices {
		for _, replaceCharacter := range deleteCharacters {
			choices[i].Message.Content = strings.ReplaceAll(choice.Message.Content, replaceCharacter, "")
		}
	}

	return choices
}

func filterByMessageLength(choices []openai.ChatCompletionChoice, length int) []openai.ChatCompletionChoice {
	filteredChoices := make([]openai.ChatCompletionChoice, 0)
	for _, choice := range choices {
		if len(choice.Message.Content) <= length {
			filteredChoices = append(filteredChoices, choice)
		}
	}
	return filteredChoices
}

func indexOfMaxMessageLength(choices []openai.ChatCompletionChoice) int {
	maxLength := 0
	indexOfMaxLength := 0
	for i, choice := range choices {
		// MEMO: len(string) で得られるのはバイト数であり、文字数ではない
		messageLength := utf8.RuneCountInString(choice.Message.Content)
		if messageLength > maxLength {
			maxLength = messageLength
			indexOfMaxLength = i
		}
	}
	return indexOfMaxLength
}
