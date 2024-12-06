package aisuite

type FunctionCall struct {
	Name string
	Args string
}

type ToolCall struct {
	ID       string
	Tool     string
	Function FunctionCall
}

// ChatCompletionMessage is a message in a chat completion request.

type Role string

const (
	User      Role = "user"
	System    Role = "system"
	Assistant Role = "assistant"
)

type ChatCompletionMessage struct {
	Role    Role
	Content string
}

type ChatCompletionRequest struct {
	Model     string
	Messages  []ChatCompletionMessage
	MaxTokens int
	Stream    bool
}

type ChatCompletionChoice struct {
	Message ChatCompletionMessage
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice
}

// ChatCompletionStreamResponse is the response from a chat completion stream.

type ChatCompletionStreamChoiceDelta struct {
	Content      string
	Role         string
	FunctionCall *FunctionCall
	ToolCalls    []ToolCall
	Refusal      string
}

type ChatCompletionStreamChoice struct {
	Delta        ChatCompletionStreamChoiceDelta
	FinishReason string
}

type ChatCompletionStreamResponse struct {
	ID      string
	Model   string
	Choices []ChatCompletionStreamChoice
}

type ChatCompletionStream interface {
	Recv() (ChatCompletionStreamResponse, error)
	Close() error
}
