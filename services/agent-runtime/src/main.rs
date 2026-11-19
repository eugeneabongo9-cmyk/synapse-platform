use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::sync::RwLock;
use uuid::Uuid;
use tracing::{info, error};

#[derive(Debug, Clone, Serialize, Deserialize)]
struct AgentTask {
    id: Uuid,
    agent_id: Uuid,
    payload: serde_json::Value,
    status: TaskStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
enum TaskStatus {
    Pending,
    Running,
    Completed,
    Failed,
}

struct Runtime {
    tasks: RwLock<HashMap<Uuid, AgentTask>>,
}

impl Runtime {
    fn new() -> Self {
        Self {
            tasks: RwLock::new(HashMap::new()),
        }
    }

    async fn submit_task(&self, agent_id: Uuid, payload: serde_json::Value) -> Uuid {
        let task = AgentTask {
            id: Uuid::new_v4(),
            agent_id,
            payload,
            status: TaskStatus::Pending,
        };

        let id = task.id;
        self.tasks.write().await.insert(id, task);

        info!(task_id = %id, "task submitted");
        id
    }

    async fn execute_task(&self, task_id: Uuid) -> Result<(), String> {
        let mut tasks = self.tasks.write().await;
        let task = tasks.get_mut(&task_id).ok_or("task not found")?;

        if task.status != TaskStatus::Pending {
            return Err("task not in pending state".into());
        }

        task.status = TaskStatus::Running;
        drop(tasks);

        // Simulate high-performance agent execution
        // In production this would dispatch to specialized executors
        tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;

        let mut tasks = self.tasks.write().await;
        if let Some(task) = tasks.get_mut(&task_id) {
            task.status = TaskStatus::Completed;
        }

        info!(task_id = %task_id, "task completed");
        Ok(())
    }
}

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let runtime = Runtime::new();

    info!("Synapse Agent Runtime started");

    // Demonstration of concurrent task execution
    let agent_id = Uuid::new_v4();
    let task_id = runtime
        .submit_task(agent_id, serde_json::json!({"action": "analyze", "data": "sample"}))
        .await;

    if let Err(e) = runtime.execute_task(task_id).await {
        error!("execution failed: {}", e);
    }
}
