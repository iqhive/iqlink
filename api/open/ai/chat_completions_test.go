package ai_test

import (
	"encoding/json"
	"testing"

	"github.com/iqhive/iqlink/api/open/ai"
	"github.com/iqhive/iqlink/xyz"
)

// TestTools round-trips the Tool union, whose accessor is built when the
// package is initialised, so an invalid case panics on import.
func TestTools(t *testing.T) {
	fn := ai.Function{Name: "get_weather", Arguments: `{"city":"Auckland"}`}
	data, err := json.Marshal(ai.Tools.Function.As(fn))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"function":{"name":"get_weather","arguments":"{\"city\":\"Auckland\"}"},"type":"function"}` {
		t.Fatalf("a function tool encoded as %s", data)
	}
	var tool ai.Tool
	if err := json.Unmarshal(data, &tool); err != nil {
		t.Fatal(err)
	}
	if xyz.ValueOf(tool) != ai.Tools.Function || ai.Tools.Function.Get(tool) != fn {
		t.Fatalf("a function tool decoded as %v", xyz.ValueOf(tool))
	}
}

// TestMessages checks that each message decodes to the case its role
// names, including an assistant message carrying tool calls.
func TestMessages(t *testing.T) {
	var messages []ai.Message
	if err := json.Unmarshal([]byte(`[
		{"role":"system","content":"be brief"},
		{"role":"user","content":"weather?"},
		{"role":"assistant","content":"","tool_calls":[{"type":"function","function":{"name":"get_weather"}}]},
		{"role":"tool","tool_call_id":"call_1","content":"sunny"}
	]`), &messages); err != nil {
		t.Fatal(err)
	}
	if len(messages) != 4 {
		t.Fatalf("decoded %d messages", len(messages))
	}
	if xyz.ValueOf(messages[0]) != ai.Messages.System || ai.Messages.System.Get(messages[0]).Content != "be brief" {
		t.Fatalf("a system message decoded as %v", xyz.ValueOf(messages[0]))
	}
	if xyz.ValueOf(messages[1]) != ai.Messages.User || ai.Messages.User.Get(messages[1]).Content != "weather?" {
		t.Fatalf("a user message decoded as %v", xyz.ValueOf(messages[1]))
	}
	if xyz.ValueOf(messages[2]) != ai.Messages.Assistant {
		t.Fatalf("an assistant message decoded as %v", xyz.ValueOf(messages[2]))
	}
	calls := ai.Messages.Assistant.Get(messages[2]).ToolCalls
	if len(calls) != 1 || xyz.ValueOf(calls[0]) != ai.Tools.Function || ai.Tools.Function.Get(calls[0]).Name != "get_weather" {
		t.Fatalf("assistant tool calls decoded as %v", calls)
	}
	if xyz.ValueOf(messages[3]) != ai.Messages.Tool || ai.Messages.Tool.Get(messages[3]) != (ai.ToolOutput{ToolCallID: "call_1", Content: "sunny"}) {
		t.Fatalf("a tool message decoded as %v", xyz.ValueOf(messages[3]))
	}
}
