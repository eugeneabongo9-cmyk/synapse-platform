package memory

import (
	"context"
	"sync"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
)

// AgentRepository is an in-memory implementation of the repository port.
type AgentRepository struct {
	mu     sync.RWMutex
	agents map[uuid.UUID]*domain.Agent
}

func NewAgentRepository() *AgentRepository {
	return &AgentRepository{
		agents: make(map[uuid.UUID]*domain.Agent),
	}
}

func (r *AgentRepository) Save(ctx context.Context, agent *domain.Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.agents[agent.ID] = agent
	return nil
}

func (r *AgentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, ok := r.agents[id]
	if !ok {
		return nil, fmt.Errorf("agent not found")
	}
	return agent, nil
}
