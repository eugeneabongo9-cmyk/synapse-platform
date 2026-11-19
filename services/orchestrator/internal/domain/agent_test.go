package domain

import (
	"testing"
)

func TestNewAgent(t *testing.T) {
	agent, err := NewAgent("test-agent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if agent.Name != "test-agent" {
		t.Errorf("expected name test-agent, got %s", agent.Name)
	}

	if agent.Status != AgentStatusPending {
		t.Errorf("expected pending status, got %s", agent.Status)
	}

	if len(agent.Events()) != 1 {
		t.Errorf("expected 1 event, got %d", len(agent.Events()))
	}
}

func TestAgentLifecycle(t *testing.T) {
	agent, _ := NewAgent("lifecycle-test")

	if err := agent.Start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	if agent.Status != AgentStatusRunning {
		t.Errorf("expected running, got %s", agent.Status)
	}

	if err := agent.Complete(); err != nil {
		t.Fatalf("complete failed: %v", err)
	}

	if agent.Status != AgentStatusCompleted {
		t.Errorf("expected completed, got %s", agent.Status)
	}
}

func TestInvalidTransitions(t *testing.T) {
	agent, _ := NewAgent("invalid")

	// Cannot complete before start
	if err := agent.Complete(); err == nil {
		t.Error("expected error when completing pending agent")
	}
}
