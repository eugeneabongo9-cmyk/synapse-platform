package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
)

type AgentRepository struct {
	mu     sync.RWMutex
	agents map[uuid.UUID]*domain.Agent
}

func NewAgentRepository() *AgentRepository {
	return &AgentRepository{agents: make(map[uuid.UUID]*domain.Agent)}
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
	a, ok := r.agents[id]
	if !ok {
		return nil, fmt.Errorf("agent not found")
	}
	return a, nil
}

func (r *AgentRepository) List(ctx context.Context) ([]*domain.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Agent, 0, len(r.agents))
	for _, a := range r.agents {
		result = append(result, a)
	}
	return result, nil
}
