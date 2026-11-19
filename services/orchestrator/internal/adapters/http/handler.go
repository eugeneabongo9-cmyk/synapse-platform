package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/application"
	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	service *application.OrchestratorService
	logger  zerolog.Logger
}

func NewHandler(service *application.OrchestratorService, logger zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/agents", h.createAgent)
	mux.HandleFunc("GET /api/v1/agents", h.listAgents)
	mux.HandleFunc("GET /api/v1/agents/{id}", h.getAgent)
	mux.HandleFunc("POST /api/v1/agents/{id}/start", h.startAgent)
	mux.HandleFunc("POST /api/v1/agents/{id}/complete", h.completeAgent)
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("/", h.serveDemo)

	return mux
}

func (h *Handler) createAgent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
		return
	}

	agent, err := h.service.CreateAgent(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentResponse(agent))
}

func (h *Handler) listAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.service.ListAgents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]map[string]any, len(agents))
	for i, a := range agents {
		resp[i] = agentResponse(a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getAgent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	agent, err := h.service.GetAgent(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentResponse(agent))
}

func (h *Handler) startAgent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.StartAgent(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) completeAgent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.service.CompleteAgent(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "service": "synapse-orchestrator"})
}

func (h *Handler) serveDemo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(demoHTML))
}

func agentResponse(a *domain.Agent) map[string]any {
	return map[string]any{
		"id":         a.ID.String(),
		"name":       a.Name,
		"status":     a.Status,
		"created_at": a.CreatedAt.Format(time.RFC3339),
		"updated_at": a.UpdatedAt.Format(time.RFC3339),
		"version":    a.Version,
	}
}

const demoHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Synapse Platform - Live Demo</title>
<style>
  :root { --bg: #0f172a; --card: #1e293b; --accent: #38bdf8; --text: #e2e8f0; --muted: #94a3b8; }
  * { box-sizing: border-box; }
  body { font-family: system-ui, -apple-system, sans-serif; background: var(--bg); color: var(--text); margin: 0; padding: 2rem; line-height: 1.5; }
  h1 { font-size: 2.2rem; margin-bottom: 0.25rem; }
  .subtitle { color: var(--muted); margin-bottom: 2rem; }
  .card { background: var(--card); border-radius: 12px; padding: 1.5rem; margin-bottom: 1.5rem; }
  input, button { padding: 0.6rem 1rem; border-radius: 8px; border: none; font-size: 1rem; }
  input { background: #334155; color: white; width: 280px; margin-right: 0.5rem; }
  button { background: var(--accent); color: #0f172a; font-weight: 600; cursor: pointer; }
  button.secondary { background: #475569; color: white; }
  .agents { display: grid; gap: 1rem; }
  .agent { background: #334155; padding: 1rem; border-radius: 8px; display: flex; justify-content: space-between; align-items: center; }
  .status { font-size: 0.8rem; padding: 0.2rem 0.6rem; border-radius: 999px; text-transform: uppercase; }
  .pending { background: #f59e0b; color: #000; }
  .running { background: #3b82f6; color: white; }
  .completed { background: #10b981; color: white; }
</style>
</head>
<body>
  <h1>Synapse Platform</h1>
  <p class="subtitle">Production-grade AI Agent Orchestrator - Live Demo</p>

  <div class="card">
    <h2>Create Agent</h2>
    <input id="name" placeholder="Agent name" />
    <button onclick="createAgent()">Create</button>
  </div>

  <div class="card">
    <h2>Active Agents</h2>
    <div id="agents">Loading...</div>
  </div>

  <footer style="margin-top:3rem;color:#94a3b8;text-align:center;font-size:0.85rem">
    Hexagonal architecture - Domain-Driven Design - Event-driven<br>
    <a href="/health" style="color:#38bdf8">Health endpoint</a>
  </footer>

<script>
async function loadAgents() {
  const res = await fetch('/api/v1/agents');
  const agents = await res.json();
  const el = document.getElementById('agents');
  if (!agents || agents.length === 0) {
    el.innerHTML = '<p style="color:#94a3b8">No agents yet. Create one above.</p>';
    return;
  }
  let html = '';
  for (const a of agents) {
    let actions = '';
    if (a.status === 'pending') {
      actions = '<button onclick="start(\'' + a.id + '\')">Start</button>';
    } else if (a.status === 'running') {
      actions = '<button class="secondary" onclick="complete(\'' + a.id + '\')">Complete</button>';
    }
    html += '<div class="agent"><div><strong>' + a.name + '</strong><br><small style="color:#94a3b8">' + a.id.slice(0,8) + ' v' + a.version + '</small></div><div><span class="status ' + a.status + '">' + a.status + '</span> ' + actions + '</div></div>';
  }
  el.innerHTML = html;
}

async function createAgent() {
  const name = document.getElementById('name').value.trim();
  if (!name) return alert('Enter a name');
  await fetch('/api/v1/agents', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({name: name})
  });
  document.getElementById('name').value = '';
  loadAgents();
}

async function start(id) {
  await fetch('/api/v1/agents/' + id + '/start', {method: 'POST'});
  loadAgents();
}

async function complete(id) {
  await fetch('/api/v1/agents/' + id + '/complete', {method: 'POST'});
  loadAgents();
}

loadAgents();
setInterval(loadAgents, 3000);
</script>
</body>
</html>`
