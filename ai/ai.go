package ai

import (
	"encoding/json"
	"eikva.ru/eikva/ai/prompts"
	envvars "eikva.ru/eikva/env_vars"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/requests"
)

func GetCreateTestCaseFormat(amount int) map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"result": map[string]any{
				"type":     "array",
				"minItems": amount,
				"maxItems": amount,
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name":           map[string]any{"type": "string"},
						"description":    map[string]any{"type": "string"},
						"source_ref":     map[string]any{"type": "string"},
						"pre_condition":  map[string]any{"type": "string"},
						"post_condition": map[string]any{"type": "string"},
						"steps": map[string]any{
							"type": "array",
							"minItems": 1,
							"maxItems": 20,
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"description":     map[string]any{"type": "string"},
									"data":            map[string]any{"type": "string"},
									"expected_result": map[string]any{"type": "string"},
								},
								"required": []string{"description", "data", "expected_result"},
							},
						},
					},
					"required": []string{"name", "description", "source_ref", "pre_condition", "post_condition", "steps"},
				},
			},
		},
	}
}

func StartTestCaseListGeneration(lenght int, input *string) (*[]*models.CreateTestCaseOutputEntry, error) {
	completionsUrl := envvars.Get(envvars.OpenAiBaseUrl) + envvars.Get(envvars.OpenAiCompletionsPathName)
	reqBody := models.OpenAiRequest{
		Model:  "qwen3:latest",
		Stream: false,
		Think:  false,
		Format: GetCreateTestCaseFormat(lenght),
		Messages: []models.ModelMessage{
			{
				Role:    "system",
				Content: prompts.CreateTestCaseSystem,
			},
			{
				Role:    "user",
				Content: prompts.CreateTestCaseUserMessageTemplate(*input),
			},
		},
	}

	var response models.ModelReponse
	err := requests.Post(completionsUrl, reqBody, &response)

	if err != nil {
		return nil, err
	}

	var generated models.ModelMessageContent
	if err := json.Unmarshal([]byte(response.Message.Content), &generated); err != nil {
		return nil, err
	}

	 return &generated.Result, nil
}
