package ai

import (
	"os"
)

/*
OPEN_AI_BASE_URL=http://localhost:11434
OPEN_AI_COMPLETIONS_PATHNAME=/api/chat
OPEN_AI_TOKENIZE_URL=/tokenize
*/

var openaiApiKey = os.Getenv("OPEN_AI_API_KEY")
var openaiBaseUrl = os.Getenv("OPEN_AI_BASE_URL")
var openAiCompletionsUrl = os.Getenv("OPEN_AI_TOKENIZE_URL")

var completionsUrl = openaiBaseUrl + openAiCompletionsUrl

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Think    bool      `json:"think"`
	Format   any       `json:"format"`
}

var createTestCaseFormat = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"name": map[string]any{
			"type":        "string",
			"description": "Обязательное поле: название тест-кейса",
		},
		"pre_condition": map[string]any{
			"type":        "string",
			"description": "Необязательное поле (может быть пустой строкой): условия, которые должны быть выполнены до начала теста",
		},
		"actions": map[string]any{
			"type":        "string",
			"description": "Обязательное поле: последовательность действий пользователя в системе",
		},
		"expected_result": map[string]any{
			"type":        "string",
			"description": "Обязательное поле: результат, который система должна показать после действий",
		},
		"post_condition": map[string]any{
			"type":        "string",
			"description": "Необязательное поле (может быть пустой строкой): действия или состояние системы после теста",
		},
	},
	"required": []string{"name", "pre_condition", "actions", "expected_result", "post_condition"},
}

/*
func CreateTestCase() {
	reqBody := Request{
		Model: "qwen3:latest",
		Stream: false,
		Think:  false,
		Format: createTestCaseFormat,
		Messages: []Message{
			{
				Role:    "system",
				Content: prompts.CreateTestCaseSystem,
			},
			{
				Role:    "user",
				Content: "",
			},
		},
	}
}
*/
