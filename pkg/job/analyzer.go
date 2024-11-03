package job

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	model        = "llama3"
	systemPrompt = "You are now a professional headhunter who is very good at analyzing jobs."
)

type JobAnalyzer struct {
	llm llms.Model
}

func NewJobAnalyzer() *JobAnalyzer {
	llm, err := ollama.New(
		ollama.WithModel(model),
		ollama.WithSystemPrompt(systemPrompt),
	)
	if err != nil {
		panic(err)
	}
	return &JobAnalyzer{
		llm: llm,
	}
}

func (ja *JobAnalyzer) Analyze(j *Job) (string, error) {
	ctx := context.Background()
	resp, err := ja.llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: fmt.Sprintf(
							`I am analyzing the job: %s\n
							Here is the job description: %s\n
							Here are the job requirements: %s\n
							Please provide a summary of the job.
							`,
							j.Title,
							j.Contents["Job Description"],
							j.Contents["Requirements"],
						),
					},
				},
			},
		},
		llms.WithTemperature(0.5),
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Content, nil
}
