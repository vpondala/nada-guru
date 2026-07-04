package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vpondala/nada-guru/agents"
	"github.com/vpondala/nada-guru/knowledge"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

const (
	evalAppName    = "nada-guru-eval"
	evalUserID     = "eval-user"
	evalTimeout    = 90 * time.Second
	evalSessionFmt = "eval-session-%s"
)

type testCase struct {
	ID             string   `json:"id"`
	Description    string   `json:"description"`
	Input          string   `json:"input"`
	ExpectContains []string `json:"expect_contains"`
	AgentExpected  string   `json:"agent_expected"`
	SkipInCI       bool     `json:"skip_in_ci"`
}

func TestEvalSuite(t *testing.T) {
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("skipping eval suite: GEMINI_API_KEY not set")
	}

	data, err := os.ReadFile("test_cases.json")
	if err != nil {
		t.Fatalf("failed to read test_cases.json: %v", err)
	}
	var cases []testCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to parse test_cases.json: %v", err)
	}

	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}

	rootAgent, err := agents.New(store)
	if err != nil {
		t.Fatalf("agents.New() failed: %v", err)
	}

	r, err := runner.New(runner.Config{
		AppName:           evalAppName,
		Agent:             rootAgent,
		SessionService:    session.InMemoryService(),
		AutoCreateSession: true,
	})
	if err != nil {
		t.Fatalf("runner.New() failed: %v", err)
	}

	passed := 0
	failed := 0
	skipped := 0

	for _, tc := range cases {
		if tc.SkipInCI && os.Getenv("CI") != "" {
			t.Logf("SKIP %s: %s (CI)", tc.ID, tc.Description)
			skipped++
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), evalTimeout)
		sessionID := fmt.Sprintf(evalSessionFmt, tc.ID)
		response, err := invokeAgent(ctx, r, sessionID, tc.Input)
		cancel()
		if err != nil {
			t.Logf("FAIL %s: %s — agent error: %v", tc.ID, tc.Description, err)
			failed++
			continue
		}

		ok := true
		for _, substr := range tc.ExpectContains {
			if !strings.Contains(strings.ToLower(response), strings.ToLower(substr)) {
				t.Logf("FAIL %s: %s — expected %q in response", tc.ID, tc.Description, substr)
				ok = false
			}
		}
		if ok {
			t.Logf("PASS %s: %s", tc.ID, tc.Description)
			passed++
		} else {
			failed++
		}
	}

	total := passed + failed + skipped
	fmt.Printf("Pass rate: %d/%d (%d%%)\n", passed, total, passed*100/max(total, 1))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func invokeAgent(ctx context.Context, r *runner.Runner, sessionID, input string) (string, error) {
	msg := genai.NewContentFromText(input, genai.RoleUser)

	var responseText strings.Builder
	for event, err := range r.Run(ctx, evalUserID, sessionID, msg, agent.RunConfig{}) {
		if err != nil {
			return "", fmt.Errorf("runner.Run: %w", err)
		}
		if event == nil || event.LLMResponse.Content == nil {
			continue
		}
		for _, part := range event.LLMResponse.Content.Parts {
			if part == nil || part.Text == "" || part.Thought {
				continue
			}
			responseText.WriteString(part.Text)
		}
	}
	return responseText.String(), nil
}
