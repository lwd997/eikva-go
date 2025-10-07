package ai

import (
	"encoding/json"
	"errors"
	"fmt"

	"eikva.ru/eikva/ai/prompts"
	envvars "eikva.ru/eikva/env_vars"
	"eikva.ru/eikva/models"
	"eikva.ru/eikva/requests"
)

func GetCreateTestCaseFormat(amount int) map[string]any {
	return map[string]any{
		"type": "json_schema",
		"json_schema": map[string]any{
			"name": "test_case_schema",
			"schema": map[string]any{
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
									"type":     "array",
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
			},
		},
	}
}

func StartTestCaseListGeneration(lenght int, input *string) (*[]*models.CreateTestCaseOutputEntry, error) {
	completionsUrl := envvars.Get(envvars.OpenAiBaseUrl) + envvars.Get(envvars.OpenAiCompletionsPathName)
	reqBody := models.OpenAiRequest{
		Model:  "qwen3-32b-awq",
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

	err := requests.Post(&requests.PostConfig{
		Url: completionsUrl,
		Headers: &map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", envvars.Get(envvars.OpenAiApiKey)),
		},
		ReqBody:  &reqBody,
		RespBody: &response,
	})

	if err != nil {
		return nil, err
	}

	if len(response.Choices) != 1 {
		return nil, errors.New("Не верный формат ответа модели")
	}

	var generated models.ModelMessageContent
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &generated); err != nil {
		return nil, err
	}

	return &generated.Result, nil
}

func StartTextCompression(content *string) (*string, error) {
	completionsUrl := envvars.Get(envvars.OpenAiBaseUrl) + envvars.Get(envvars.OpenAiCompletionsPathName)
	reqBody := models.OpenAiRequest{
		Model:  "qwen3-32b-awq",
		Stream: false,
		Think:  false,
		Messages: []models.ModelMessage{
			{
				Role:    "system",
				Content: prompts.CompressUploadSystem,
			},
			{
				Role:    "user",
				Content: prompts.CreateCompressUserMessageTemplate(*content),
			},
		},
	}

	var response models.ModelReponse

	err := requests.Post(&requests.PostConfig{
		Url: completionsUrl,
		Headers: &map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", envvars.Get(envvars.OpenAiApiKey)),
		},
		ReqBody:  &reqBody,
		RespBody: &response,
	})

	if err != nil {
		return nil, err
	}

	if len(response.Choices) != 1 {
		return nil, errors.New("Не верный формат ответа модели")
	}

	return &response.Choices[0].Message.Content, nil
}
