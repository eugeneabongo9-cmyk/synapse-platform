package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AgentStatus represents the lifecycle state of an AI agent.
type AgentStatus string

const (
	AgentStatusPending   AgentStatus = "pending"
	AgentStatusRunning   AgentStatus = "running"
	AgentStatusCompleted AgentStatus = "completed"
	AgentStatusFailed    AgentStatus = "failed"
	AgentStatusCancelled AgentStatus = "cancelled"
)

// Agent is the aggregate root for agent orchestration.
type Agent struct {
	ID          uuid.UUID
	Name        string
	Status      AgentStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
	events      []DomainEvent
}

// NewAgent creates a new agent aggregate.
func NewAgent(name string) (*Agent, error) {
	if name == "" {
		return nil, errors.New("agent name cannot be empty")
	}

	now := time.Now().UTC()
	agent := &Agent{
		ID:        uuid.New(),
		Name:      name,
		Status:    AgentStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}

	agent.raise(AgentCreatedEvent{
		AgentID:   agent.ID,
		Name:      name,
		CreatedAt: now,
	})

	return agent, nil
}

// Start transitions the agent to running state.
func (a *Agent) Start() error {
	if a.Status != AgentStatusPending {
		return errors.New("agent can only be started from pending state")
	}

	a.Status = AgentStatusRunning
	a.UpdatedAt = time.Now().UTC()
	a.Version++

	a.raise(AgentStartedEvent{
		AgentID:   a.ID,
		StartedAt: a.UpdatedAt,
	})

	return nil
}

// Complete marks the agent as successfully completed.
func (a *Agent) Complete() error {
	if a.Status != AgentStatusRunning {
		return errors.New("agent can only be completed from running state")
	}

	a.Status = AgentStatusCompleted
	a.UpdatedAt = time.Now().UTC()
	a.Version++

	a.raise(AgentCompletedEvent{
		AgentID:     a.ID,
		CompletedAt: a.UpdatedAt,
	})

	return nil
}

func (a *Agent) raise(event DomainEvent) {
	a.events = append(a.events, event)
}

// Events returns uncommitted domain events.
func (a *Agent) Events() []DomainEvent {
	return a.events
}

// ClearEvents clears uncommitted events after persistence.
func (a *Agent) ClearEvents() {
	a.events = nil
}

// DomainEvent is the base interface for all domain events.
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

type AgentCreatedEvent struct {
	AgentID   uuid.UUID
	Name      string
	CreatedAt time.Time
}

func (e AgentCreatedEvent) EventType() string   { return "agent.created" }
func (e AgentCreatedEvent) OccurredAt() time.Time { return e.CreatedAt }

type AgentStartedEvent struct {
	AgentID   uuid.UUID
	StartedAt time.Time
}

func (e AgentStartedEvent) EventType() string   { return "agent.started" }
func (e AgentStartedEvent) OccurredAt() time.Time { return e.StartedAt }

type AgentCompletedEvent struct {
	AgentID     uuid.UUID
	CompletedAt time.Time
}

func (e AgentCompletedEvent) EventType() string   { return "agent.completed" }
func (e AgentCompletedEvent) OccurredAt() time.Time { return e.CompletedAt }
