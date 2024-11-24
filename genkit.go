package main

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/googleai"
	"github.com/invopop/jsonschema"
)

type SummarizeLLM struct{}
type promptInput struct {
	URL string `json:"url"`
}

func SummarizeBlog(ctx context.Context, url string) (string, error) {
	// Intialize LLM API (Gemini)
	if err := googleai.Init(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// Define the webLoader tool
	webLoader := ai.DefineTool(
		"webLoader",
		"Loads a webpage and returns the textual content.",
		func(ctx context.Context, input struct {
			URL string `json:"url"`
		}) (string, error) {
			return fetchWebContent(input.URL)
		},
	)

	// Select Model
	model := googleai.Model("gemini-1.5-flash")

	// Define the Prompt
	summarizePrompt, err := dotprompt.Define("summarizePrompt",
		"First, fetch this link: {{url}}. Then, summarize the content within 20 words.",
		dotprompt.Config{
			Model: model,
			Tools: []ai.Tool{webLoader},
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 1,
			},
			InputSchema:  jsonschema.Reflect(promptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)
	if err != nil {
		return "", err
	}

	// Define the flow
	flow := genkit.DefineFlow("summarizeFlow", func(ctx context.Context, input string) (string, error) {
		resp, err := summarizePrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: &promptInput{
					URL: input,
				},
			},
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate summary: %w", err)
		}
		return resp.Text(), nil
	})

	// Run the flow
	result, err := flow.Run(ctx, url)
	if err != nil {
		return "", fmt.Errorf("failed to run the flow: %w", err)
	}
	return result, nil
}
