package models

type ModelMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAiRequest struct {
	Model    string         `json:"model"`
	Messages []ModelMessage `json:"messages"`
	Stream   bool           `json:"stream"`
	Think    bool           `json:"think"`
	Format   any            `json:"response_format"`
}

type Choice struct {
	Message ModelMessage `json:"message"`
}

type ModelReponse struct {
	Model              string   `json:"model"`
	CreatedAt          string   `json:"created_at"`
	Choices            []Choice `json:"choices"`
	Done               bool     `json:"done"`
	DoneReason         string   `json:"done_reason"`
	TotalDuration      int      `json:"total_duration"`
	LoadDuration       int      `json:"load_duration"`
	PromptEvalCount    int      `json:"prompt_eval_count"`
	PromptEvalDuration int      `json:"prompt_eval_duration"`
	EvalCount          int      `json:"eval_count"`
	EvalDuration       int      `json:"eval_duration"`
}

type ModelMessageContent struct {
	Result []*CreateTestCaseOutputEntry `json:"result"`
}

type CreateTestCaseOutputStep struct {
	Data           string `json:"data"`
	Description    string `json:"description"`
	ExpectedResult string `json:"expected_result"`
}

type CreateTestCaseOutputEntry struct {
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	SourceRef     string                     `json:"source_ref"`
	PreCondition  string                     `json:"pre_condition"`
	PostCondition string                     `json:"post_condition"`
	Steps         []CreateTestCaseOutputStep `json:"steps"`
}
