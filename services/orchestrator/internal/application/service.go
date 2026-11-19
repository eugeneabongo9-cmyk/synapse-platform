package application

import (
	"context"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type AgentRepository interface {
	Save(ctx context.Context, agent *domain.Agent) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Agent, error)
	List(ctx context.Context) ([]*domain.Agent, error)
}

type EventBus interface {
	Publish(ctx context.Context, events []domain.DomainEvent) error
}

type OrchestratorService struct {
	repo     AgentRepository
	eventBus EventBus
	logger   zerolog.Logger
}

func NewOrchestratorService(repo AgentRepository, bus EventBus, logger zerolog.Logger) *OrchestratorService {
	return &OrchestratorService{repo: repo, eventBus: bus, logger: logger}
}

func (s *OrchestratorService) CreateAgent(ctx context.Context, name string) (*domain.Agent, error) {
	agent, err := domain.NewAgent(name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, agent); err != nil {
		return nil, err
	}
	_ = s.eventBus.Publish(ctx, agent.Events())
	agent.ClearEvents()
	s.logger.Info().Str("id", agent.ID.String()).Str("name", name).Msg("agent created")
	return agent, nil
}

func (s *OrchestratorService) StartAgent(ctx context.Context, id uuid.UUID) error {
	agent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := agent.Start(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, agent); err != nil {
		return err
	}
	_ = s.eventBus.Publish(ctx, agent.Events())
	agent.ClearEvents()
	s.logger.Info().Str("id", id.String()).Msg("agent started")
	return nil
}

func (s *OrchestratorService) CompleteAgent(ctx context.Context, id uuid.UUID) error {
	agent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := agent.Complete(); err != nil {
		return err
	}
	if err := s.repo.Save(ctx, agent); err != nil {
		return err
	}
	_ = s.eventBus.Publish(ctx, agent.Events())
	agent.ClearEvents()
	s.logger.Info().Str("id", id.String()).Msg("agent completed")
	return nil
}

func (s *OrchestratorService) GetAgent(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrchestratorService) ListAgents(ctx context.Context) ([]*domain.Agent, error) {
	return s.repo.List(ctx)
}
