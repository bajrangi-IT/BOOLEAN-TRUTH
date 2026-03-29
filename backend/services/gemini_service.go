package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// --- Request/Response Structures ---

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

// --- Result Structure (returned to API caller) ---

type FactCheckResult struct {
	IsValid             bool     `json:"is_valid"`
	Status              string   `json:"status"` // "TRUE", "FALSE", "PARTIALLY_TRUE"
	Reasoning           string   `json:"reasoning"`
	DetailedExplanation string   `json:"detailed_explanation"`
	Sources             []string `json:"sources"`
	ToxicityScore       float32  `json:"toxicity_score"`
	SpamScore           float32  `json:"spam_score"`
}

// VerifyStatement calls the Gemini REST API directly (no SDK).
func VerifyStatement(ctx context.Context, text string) (*FactCheckResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not found in environment")
	}

	// Using gemini-1.5-flash via v1beta REST endpoint directly
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=%s",
		apiKey,
	)

	prompt := fmt.Sprintf(`Analyze the following statement for truthfulness, toxicity, and spam.
Statement: "%s"

You MUST respond with ONLY valid JSON. No markdown, no explanation, just the JSON.
Required JSON format:
{
  "is_valid": true or false,
  "status": "TRUE" or "FALSE" or "PARTIALLY_TRUE",
  "reasoning": "Short summary of the verdict",
  "detailed_explanation": "In-depth analysis of what is correct or incorrect, with evidence",
  "sources": ["Source 1", "Source 2"],
  "toxicity_score": 0.0,
  "spam_score": 0.0
}`, text)

	requestBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error %d: %s", resp.StatusCode, string(respBody))
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(respBody, &gemResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response envelope: %v", err)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	rawJSON := gemResp.Candidates[0].Content.Parts[0].Text

	var result FactCheckResult
	if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
		// Try to strip markdown wrappers and retry
		cleaned := cleanJSON(rawJSON)
		if err2 := json.Unmarshal([]byte(cleaned), &result); err2 != nil {
			return nil, fmt.Errorf("failed to parse AI JSON: %v\nRaw: %s", err2, rawJSON)
		}
	}

	return &result, nil
}

// cleanJSON strips markdown code fences from a JSON string robustly.
func cleanJSON(s string) string {
	s = trimSpace(s)

	// Find opening fence: ```json or ```
	openFences := []string{"```json", "```"}
	startIdx := -1
	fenceLen := 0
	for _, fence := range openFences {
		idx := strings.Index(s, fence)
		if idx >= 0 {
			startIdx = idx
			fenceLen = len(fence)
			break
		}
	}

	if startIdx >= 0 {
		// Move past the opening fence
		s = s[startIdx+fenceLen:]
		// Skip optional newline right after the fence
		if len(s) > 0 && s[0] == '\n' {
			s = s[1:]
		}
		// Find closing fence
		closeIdx := strings.Index(s, "```")
		if closeIdx >= 0 {
			s = s[:closeIdx]
		}
	}

	return trimSpace(s)
}

func trimSpace(s string) string {
	start := 0
	end := len(s) - 1
	for start <= end && (s[start] == ' ' || s[start] == '\n' || s[start] == '\r' || s[start] == '\t') {
		start++
	}
	for end >= start && (s[end] == ' ' || s[end] == '\n' || s[end] == '\r' || s[end] == '\t') {
		end--
	}
	if start > end {
		return ""
	}
	return s[start : end+1]
}
