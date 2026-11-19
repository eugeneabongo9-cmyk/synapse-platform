package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AgentRepository struct {
	pool *pgxpool.Pool
}

func NewAgentRepository(pool *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{pool: pool}
}

func (r *AgentRepository) Save(ctx context.Context, agent *domain.Agent) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Upsert agent
	_, err = tx.Exec(ctx, `
		INSERT INTO agents (id, name, status, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at,
			version = EXCLUDED.version
	`, agent.ID, agent.Name, agent.Status, agent.CreatedAt, agent.UpdatedAt, agent.Version)
	if err != nil {
		return fmt.Errorf("save agent: %w", err)
	}

	// Write domain events to outbox in the same transaction
	for _, event := range agent.Events() {
		payload, _ := json.Marshal(event)
		_, err = tx.Exec(ctx, `
			INSERT INTO outbox (id, aggregate_id, event_type, payload, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`, uuid.New(), agent.ID, event.EventType(), payload, time.Now().UTC())
		if err != nil {
			return fmt.Errorf("write outbox: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *AgentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Agent, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, status, created_at, updated_at, version
		FROM agents WHERE id = $1
	`, id)

	var a domain.Agent
	var status string
	err := row.Scan(&a.ID, &a.Name, &status, &a.CreatedAt, &a.UpdatedAt, &a.Version)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent not found")
	}
	if err != nil {
		return nil, err
	}
	a.Status = domain.AgentStatus(status)
	return &a, nil
}

func (r *AgentRepository) List(ctx context.Context) ([]*domain.Agent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, status, created_at, updated_at, version
		FROM agents ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*domain.Agent
	for rows.Next() {
		var a domain.Agent
		var status string
		if err := rows.Scan(&a.ID, &a.Name, &status, &a.CreatedAt, &a.UpdatedAt, &a.Version); err != nil {
			return nil, err
		}
		a.Status = domain.AgentStatus(status)
		agents = append(agents, &a)
	}
	return agents, rows.Err()
}
