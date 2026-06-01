package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type OpenAIConfig struct {
	Model          string `json:"model"`
	SystemPrompt   string `json:"system_prompt"`
	PromptTemplate string `json:"prompt_template"`
}

type OpenAIExecutor struct {
	apiKey string
	logger *zap.SugaredLogger
	client *http.Client
}

func NewOpenAIExecutor(apiKey string, logger *zap.SugaredLogger) Executor {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	return &OpenAIExecutor{
		apiKey: apiKey,
		logger: logger,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (e *OpenAIExecutor) Execute(ctx context.Context, agentID, userID string, configStr string, input interface{}) (map[string]interface{}, error) {
	if e.apiKey == "" {
		return nil, errors.New("openai api key not configured")
	}

	var cfg OpenAIConfig
	if configStr != "" {
		_ = json.Unmarshal([]byte(configStr), &cfg)
	}
	model := cfg.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	system := cfg.SystemPrompt
	prompt := cfg.PromptTemplate
	if prompt == "" {
		// fallback: use input as prompt
		if b, err := json.Marshal(input); err == nil {
			prompt = string(b)
		} else {
			prompt = ""
		}
	}

	reqBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": system},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}

	reqBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(reqBytes))
	if err != nil {
		e.logger.Errorf("openai request create error: %v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		e.logger.Errorf("openai request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		e.logger.Errorf("openai non-2xx: %d %s", resp.StatusCode, string(body))
		return nil, errors.New("openai api error")
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		e.logger.Errorf("openai response unmarshal error: %v", err)
		return nil, err
	}

	var text string
	if len(parsed.Choices) > 0 {
		text = parsed.Choices[0].Message.Content
	}

	return map[string]interface{}{
		"agent_id": agentID,
		"user_id":  userID,
		"status":   "executed",
		"result": map[string]interface{}{
			"text": text,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}, nil
}
