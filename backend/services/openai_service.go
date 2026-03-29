package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAIFactCheckResult struct {
	IsValid             bool     `json:"is_valid"`
	Status              string   `json:"status"` // "TRUE", "FALSE", "PARTIALLY_TRUE"
	Reasoning           string   `json:"reasoning"`
	DetailedExplanation string   `json:"detailed_explanation"`
	Sources             []string `json:"sources"`
	ToxicityScore       float32  `json:"toxicity_score"`
	SpamScore           float32  `json:"spam_score"`
}

func VerifyStatementOpenAI(ctx context.Context, text string) (*OpenAIFactCheckResult, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not found in environment")
	}

	client := openai.NewClient(apiKey)

	prompt := fmt.Sprintf(`
Analyze the following statement for truthfulness, toxicity, and spam.
Statement: "%%s"

You MUST provide your response in valid JSON format ONLY. 
JSON Structure:
{
  "is_valid": boolean,
  "status": "TRUE" | "FALSE" | "PARTIALLY_TRUE",
  "reasoning": "A concise overview of whether it's true/false/spam/toxic",
  "detailed_explanation": "A deep analysis of the claim and identifying what's specifically wrong or right",
  "sources": ["List of authoritative sources or general knowledge domains derived from"],
  "toxicity_score": number (0.0 to 1.0),
  "spam_score": number (0.0 to 1.0)
}
`, text)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a professional forensic reality auditor. Always respond with JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	jsonStr := resp.Choices[0].Message.Content
	var result OpenAIFactCheckResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %%v", err)
	}

	return &result, nil
}
