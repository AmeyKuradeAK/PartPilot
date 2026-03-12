package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Normalizer struct {
	client *openai.Client
}

func NewNormalizer(apiKey string) *Normalizer {
	return &Normalizer{
		client: openai.NewClient(apiKey),
	}
}

// Normalize takes a human-readable component description and attempts to
// convert it into a standardized, searchable Manufacturer Part Number (MPN)
// or standard identifier.
func (n *Normalizer) Normalize(ctx context.Context, description string) (string, error) {
	if n.client == nil {
		return "", fmt.Errorf("openai client not configured")
	}

	prompt := fmt.Sprintf(`You are a hardware engineering assistant.
Your job is to take a human-written component description and return ONLY the most likely Manufacturer Part Number (MPN) or standard component identifier.

Rules:
1. Return ONLY the part number string. No explanation, no markdown, no surrounding quotes.
2. If the description is a standard jellybean part (e.g., '10k 0603 resistor'), return a standard recognizable representation or a generic MPN (e.g., 'RC0603FR-0710KL').
3. If you cannot determine a part number, return the original description cleaned up.

Description: "%s"`, description)

	resp, err := n.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a specialized hardware component normalizer.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.1, // low temp for deterministic output
		},
	)

	if err != nil {
		return "", fmt.Errorf("openai completion error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned from openai")
	}

	normalized := strings.TrimSpace(resp.Choices[0].Message.Content)
	return normalized, nil
}
