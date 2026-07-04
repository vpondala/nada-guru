// Package agents implements specialist agents for the Nāda Guru system.
package agents

import (
	"context"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const searchModelName = "gemini-2.5-flash"

const searchInstruction = `You are Nāda Guru's Web Search Specialist.
Your sole job is to search the web using Google Search to answer the user's query when requested.
Use the Google Search tool to find relevant, up-to-date information.
Always clearly state that the source of the information is web search (e.g., "Based on web search: ...") and include the source URL(s).`

// NewSearchAgent creates a new dedicated search agent with a unique name.
func NewSearchAgent(name string) (agent.Agent, error) {
	ctx := context.Background()
	model, err := gemini.NewModel(ctx, searchModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        name,
		Model:       model,
		Description: "Queries the web using Google Search to find external information about Carnatic music when the knowledge base is insufficient.",
		Instruction: searchInstruction,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
