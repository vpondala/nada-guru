package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/vpondala/nada-guru/agents"
	"github.com/vpondala/nada-guru/knowledge"
)

type testCase struct {
	ID             string   `json:"id"`
	Description    string   `json:"description"`
	Input          string   `json:"input"`
	ExpectContains []string `json:"expect_contains"`
	AgentExpected  string   `json:"agent_expected"`
	SkipInCI      bool     `json:"skip_in_ci"`
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

	passed := 0
	failed := 0
	skipped := 0

	for _, tc := range cases {
		if tc.SkipInCI && os.Getenv("CI") != "" {
			t.Logf("SKIP %s: %s (CI)", tc.ID, tc.Description)
			skipped++
			continue
		}

		ctx := context.Background()
		response, err := invokeAgent(ctx, rootAgent, tc.Input)
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

func invokeAgent(ctx context.Context, agent interface{}, input string) (string, error) {
	_ = ctx
	_ = agent
	_ = input
	return "[agent invocation not yet wired]", nil
}
