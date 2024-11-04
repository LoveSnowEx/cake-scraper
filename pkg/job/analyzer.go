package job

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	model        = "llama3.2"
	systemPrompt = "You are an experienced Software Engineer who is very good at look for jobs."
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

func extractResponseJson(response string) string {
	beginPattern := "{\n"
	endPattern := "\n}"
	begin := strings.Index(response, beginPattern)
	end := strings.LastIndex(response, endPattern)
	if begin == -1 || end == -1 {
		return response
	}
	response = response[begin+len(beginPattern) : end]
	return strings.Trim(response, " \n")
}

func (ja *JobAnalyzer) Analyze(j *Job) (string, error) {
	provideData := j.Contents
	provideData["Company"] = j.Company
	provideData["Title"] = j.Title
	ctx := context.Background()
	jsonData, err := json.MarshalIndent(provideData, "", "    ")
	if err != nil {
		return "", err
	}
	resp, err := ja.llm.GenerateContent(
		ctx,
		[]llms.MessageContent{
			{
				Role: llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{
					llms.TextContent{
						Text: `
Organize the key details from the provided job information in JSON format as follows:

{
	"Company": "",
	"Title": "",
	"Programming Languages": [],
	"Required Skills": [],
	"Preferred Skills": []
}

If any information is missing, use "Unknown" as the placeholder. Do not add or modify any fields other than the ones listed above.
Job information:
` + "```json\n" + string(jsonData) + "\n```",
					},
				},
			},
		},
		llms.WithTemperature(0.5),
	)
	if err != nil {
		return "", err
	}
	return extractResponseJson(resp.Choices[0].Content), nil
}
