package orchestrate

import "context"

type LLMInvoker interface {
	Invoke(ctx context.Context, dir, agent, prompt string) (output string, exitCode int, err error)
}

var DefaultLLMInvoker LLMInvoker = openCodeInvoker{}

type openCodeInvoker struct{}

func (openCodeInvoker) Invoke(ctx context.Context, dir, agent, prompt string) (string, int, error) {
	return InvokeOpenCode(ctx, dir, agent, prompt)
}
