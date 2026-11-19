package application

import (
	"context"
	"fmt"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// AgentRepository defines the port for agent persistence.
type AgentRepository interface {
	Save(ctx context.Context, agent *domain.Agent) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Agent, error)
}

// EventBus defines the port for publishing domain events.
type EventBus interface {
	Publish(ctx context.Context, events []domain.DomainEvent) error
}

// OrchestratorService is the application service coordinating use cases.
type OrchestratorService struct {
	repo     AgentRepository
	eventBus EventBus
	logger   zerolog.Logger
}

func NewOrchestratorService(repo AgentRepository, bus EventBus, logger zerolog.Logger) *OrchestratorService {
	return &OrchestratorService{
		repo:     repo,
		eventBus: bus,
		logger:   logger,
	}
}

// CreateAgent handles the create agent use case.
func (s *OrchestratorService) CreateAgent(ctx context.Context, name string) (*domain.Agent, error) {
	agent, err := domain.NewAgent(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	if err := s.repo.Save(ctx, agent); err != nil {
		return nil, fmt.Errorf("failed to persist agent: %w", err)
	}

	if err := s.eventBus.Publish(ctx, agent.Events()); err != nil {
		s.logger.Error().Err(err).Msg("failed to publish events")
		// Note: In production, this would use outbox pattern
	}

	agent.ClearEvents()

	s.logger.Info().
		Str("agent_id", agent.ID.String()).
		Str("name", name).
		Msg("agent created")

	return agent, nil
}

// StartAgent handles starting an agent.
func (s *OrchestratorService) StartAgent(ctx context.Context, id uuid.UUID) error {
	agent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("agent not found: %w", err)
	}

	if err := agent.Start(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, agent); err != nil {
		return fmt.Errorf("failed to save agent: %w", err)
	}

	return s.eventBus.Publish(ctx, agent.Events())
}
