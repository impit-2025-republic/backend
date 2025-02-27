package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Document struct {
	DocID   int    `json:"doc_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type VLLMRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

type VLLMResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

type StreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Text         string  `json:"text"`
		LogProbs     *string `json:"logprobs"`
		FinishReason *string `json:"finish_reason"`
		StopReason   *string `json:"stop_reason"`
	} `json:"choices"`
	Usage *string `json:"usage"`
}

type Vllm struct {
	vllmURL string
	apiKey  string
}

func NewVllm(vllmURL, apiKey string) Vllm {
	return Vllm{
		vllmURL: vllmURL,
		apiKey:  apiKey,
	}
}

func (s *Vllm) MakeVLLMIndexes(messages []Message, temperature float64) (string, error) {
	var prompt string
	for _, msg := range messages {
		prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	vllmReq := VLLMRequest{
		Model:       "Vikhrmodels/Vikhr-Nemo-12B-Instruct-R-21-09-24",
		MaxTokens:   2048,
		Prompt:      prompt,
		Temperature: temperature,
		Stream:      false,
	}

	jsonData, err := json.Marshal(vllmReq)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", s.vllmURL+"/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token-abc123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var vllmResp VLLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&vllmResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(vllmResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return vllmResp.Choices[0].Text, nil
}

func (s *Vllm) MakeVLLMRequest(messages []Message, temperature float64) (*http.Response, error) {
	var prompt string
	for _, msg := range messages {
		prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	vllmReq := VLLMRequest{
		Model:       "Vikhrmodels/Vikhr-Nemo-12B-Instruct-R-21-09-24",
		Prompt:      prompt,
		MaxTokens:   2048,
		Temperature: temperature,
		Stream:      true,
	}

	jsonData, err := json.Marshal(vllmReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", s.vllmURL+"/v1/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token-abc123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}
