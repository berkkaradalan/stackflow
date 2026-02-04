package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		// Enable pgvector extension
		`CREATE EXTENSION IF NOT EXISTS vector;`,

		// 1. organizations table
		`CREATE TABLE IF NOT EXISTS organizations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			logo_url TEXT,
			settings JSONB DEFAULT '{}',
			subscription_tier VARCHAR(50) DEFAULT 'free',
			subscription_ends_at TIMESTAMP,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug);`,
		`CREATE INDEX IF NOT EXISTS idx_organizations_is_active ON organizations(is_active);`,

		// 2. users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'member',
			avatar_url TEXT,
			is_active BOOLEAN DEFAULT true,
			email_verified BOOLEAN DEFAULT false,
			last_login_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);`,

		// 3. refresh_tokens table
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(500) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			is_revoked BOOLEAN DEFAULT false,
			revoked_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);`,

		// 4. invitations table
		`CREATE TABLE IF NOT EXISTS invitations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			invited_by UUID NOT NULL REFERENCES users(id),
			email VARCHAR(255) NOT NULL,
			role VARCHAR(50) DEFAULT 'member',
			token VARCHAR(255) UNIQUE NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			expires_at TIMESTAMP NOT NULL,
			accepted_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_invitations_organization_id ON invitations(organization_id);`,
		`CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(email);`,
		`CREATE INDEX IF NOT EXISTS idx_invitations_token ON invitations(token);`,
		`CREATE INDEX IF NOT EXISTS idx_invitations_status ON invitations(status);`,

		// 5. projects table
		`CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			created_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			tech_stack JSONB DEFAULT '{}',
			workflow_template VARCHAR(100),
			status VARCHAR(50) DEFAULT 'active',
			settings JSONB DEFAULT '{}',
			code_repository_url TEXT,
			deployed_url TEXT,
			started_at TIMESTAMP,
			completed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_projects_organization_id ON projects(organization_id);`,
		`CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);`,
		`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);`,
		`CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC);`,

		// 6. workflows table
		`CREATE TABLE IF NOT EXISTS workflows (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT,
			workflow_type VARCHAR(50),
			definition JSONB NOT NULL,
			is_active BOOLEAN DEFAULT true,
			usage_count INT DEFAULT 0,
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_workflow_type ON workflows(workflow_type);`,
		`CREATE INDEX IF NOT EXISTS idx_workflows_is_active ON workflows(is_active);`,

		// 7. project_workflows table
		`CREATE TABLE IF NOT EXISTS project_workflows (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
			is_active BOOLEAN DEFAULT true,
			custom_config JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(project_id, workflow_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_project_workflows_project_id ON project_workflows(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_project_workflows_workflow_id ON project_workflows(workflow_id);`,

		// 8. agents table
		`CREATE TABLE IF NOT EXISTS agents (
			id VARCHAR(100) PRIMARY KEY,
			organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			role VARCHAR(100) NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			model_config JSONB DEFAULT '{}',
			status VARCHAR(50) DEFAULT 'available',
			max_concurrent_tasks INT DEFAULT 3,
			current_tasks_count INT DEFAULT 0,
			total_tasks_completed INT DEFAULT 0,
			total_tasks_failed INT DEFAULT 0,
			success_rate DECIMAL(5,4) DEFAULT 0.0000,
			avg_completion_time_hours DECIMAL(10,2),
			last_active_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_agents_organization_id ON agents(organization_id);`,
		`CREATE INDEX IF NOT EXISTS idx_agents_role ON agents(role);`,
		`CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status);`,

		// 9. tasks table
		`CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			created_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
			parent_task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
			title VARCHAR(500) NOT NULL,
			description TEXT,
			task_type VARCHAR(100) NOT NULL,
			priority VARCHAR(50) DEFAULT 'medium',
			status VARCHAR(50) DEFAULT 'created',
			assigned_to VARCHAR(100),
			assigned_at TIMESTAMP,
			requirements JSONB DEFAULT '{}',
			context JSONB DEFAULT '{}',
			result JSONB,
			feedback TEXT,
			estimated_time_hours DECIMAL(10,2),
			actual_time_hours DECIMAL(10,2),
			started_at TIMESTAMP,
			completed_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks(created_by);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_parent_task_id ON tasks(parent_task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);`,

		// 10. task_assignments table
		`CREATE TABLE IF NOT EXISTS task_assignments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			from_agent VARCHAR(100),
			to_agent VARCHAR(100) NOT NULL,
			assignment_reason TEXT,
			assigned_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_task_assignments_task_id ON task_assignments(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_task_assignments_to_agent ON task_assignments(to_agent);`,
		`CREATE INDEX IF NOT EXISTS idx_task_assignments_assigned_at ON task_assignments(assigned_at DESC);`,

		// 11. task_events table (Event Sourcing)
		`CREATE TABLE IF NOT EXISTS task_events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			event_type VARCHAR(100) NOT NULL,
			agent_id VARCHAR(100),
			user_id UUID REFERENCES users(id),
			event_data JSONB DEFAULT '{}',
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_task_events_task_id ON task_events(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_task_events_event_type ON task_events(event_type);`,
		`CREATE INDEX IF NOT EXISTS idx_task_events_agent_id ON task_events(agent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_task_events_created_at ON task_events(created_at DESC);`,

		// 12. task_conversations table
		`CREATE TABLE IF NOT EXISTS task_conversations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			sender_type VARCHAR(50) NOT NULL,
			sender_id VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			message_type VARCHAR(50) DEFAULT 'text',
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_task_conversations_task_id ON task_conversations(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_task_conversations_created_at ON task_conversations(created_at DESC);`,

		// 13. agent_memory table (with pgvector)
		`CREATE TABLE IF NOT EXISTS agent_memory (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			agent_id VARCHAR(100) NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
			task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
			project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
			memory_type VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			embedding vector(1536),
			importance_score DECIMAL(5,4) DEFAULT 0.5000,
			access_count INT DEFAULT 0,
			last_accessed_at TIMESTAMP,
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_memory_agent_id ON agent_memory(agent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_memory_task_id ON agent_memory(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_memory_project_id ON agent_memory(project_id);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_memory_memory_type ON agent_memory(memory_type);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_memory_importance_score ON agent_memory(importance_score DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_memory_created_at ON agent_memory(created_at DESC);`,
		// pgvector index - only create if not exists
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_indexes 
				WHERE indexname = 'idx_agent_memory_embedding'
			) THEN
				CREATE INDEX idx_agent_memory_embedding ON agent_memory 
				USING ivfflat (embedding vector_cosine_ops) 
				WITH (lists = 100);
			END IF;
		END
		$$;`,

		// 14. agent_notes table
		`CREATE TABLE IF NOT EXISTS agent_notes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			agent_id VARCHAR(100) NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			file_path TEXT NOT NULL,
			note_type VARCHAR(50),
			summary TEXT,
			file_size_bytes BIGINT,
			last_synced_at TIMESTAMP DEFAULT NOW(),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_notes_agent_id ON agent_notes(agent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_notes_project_id ON agent_notes(project_id);`,

		// 15. code_artifacts table
		`CREATE TABLE IF NOT EXISTS code_artifacts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			agent_id VARCHAR(100) NOT NULL REFERENCES agents(id),
			file_path TEXT NOT NULL,
			language VARCHAR(50),
			content TEXT NOT NULL,
			version INT DEFAULT 1,
			is_latest BOOLEAN DEFAULT true,
			diff_from_previous TEXT,
			lines_of_code INT,
			validation_status VARCHAR(50),
			validation_errors JSONB,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_code_artifacts_task_id ON code_artifacts(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_code_artifacts_agent_id ON code_artifacts(agent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_code_artifacts_is_latest ON code_artifacts(is_latest);`,
		`CREATE INDEX IF NOT EXISTS idx_code_artifacts_created_at ON code_artifacts(created_at DESC);`,

		// 16. test_results table
		`CREATE TABLE IF NOT EXISTS test_results (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			agent_id VARCHAR(100) NOT NULL REFERENCES agents(id),
			test_type VARCHAR(50),
			test_name VARCHAR(255),
			status VARCHAR(50),
			execution_time_ms INT,
			error_message TEXT,
			stack_trace TEXT,
			test_output JSONB,
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_task_id ON test_results(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_status ON test_results(status);`,
		`CREATE INDEX IF NOT EXISTS idx_test_results_created_at ON test_results(created_at DESC);`,

		// 17. cost_tracking table
		`CREATE TABLE IF NOT EXISTS cost_tracking (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
			task_id UUID REFERENCES tasks(id) ON DELETE SET NULL,
			agent_id VARCHAR(100) REFERENCES agents(id),
			model_name VARCHAR(100) NOT NULL,
			operation_type VARCHAR(50),
			input_tokens INT,
			output_tokens INT,
			total_tokens INT,
			cost_usd DECIMAL(10,6),
			metadata JSONB DEFAULT '{}',
			created_at TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE INDEX IF NOT EXISTS idx_cost_tracking_organization_id ON cost_tracking(organization_id);`,
		`CREATE INDEX IF NOT EXISTS idx_cost_tracking_task_id ON cost_tracking(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_cost_tracking_agent_id ON cost_tracking(agent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_cost_tracking_created_at ON cost_tracking(created_at DESC);`,

		// 18. agent_performance_metrics table
		`CREATE TABLE IF NOT EXISTS agent_performance_metrics (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			agent_id VARCHAR(100) NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
			metric_date DATE NOT NULL,
			tasks_completed INT DEFAULT 0,
			tasks_failed INT DEFAULT 0,
			avg_completion_time_hours DECIMAL(10,2),
			code_review_pass_rate DECIMAL(5,4),
			test_pass_rate DECIMAL(5,4),
			total_cost_usd DECIMAL(10,2),
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(agent_id, metric_date)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_performance_metrics_agent_id ON agent_performance_metrics(agent_id);`,
		`CREATE INDEX IF NOT EXISTS idx_agent_performance_metrics_metric_date ON agent_performance_metrics(metric_date DESC);`,

		// Create updated_at trigger function
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ language 'plpgsql';`,

		// Apply updated_at triggers to relevant tables
		`DO $$
		DECLARE
			t text;
		BEGIN
			FOR t IN 
				SELECT tablename 
				FROM pg_tables 
				WHERE schemaname = 'public' 
				AND tablename IN ('organizations', 'users', 'projects', 'tasks', 'agents', 'agent_memory', 'agent_notes', 'workflows')
			LOOP
				EXECUTE format('
					DROP TRIGGER IF EXISTS update_%I_updated_at ON %I;
					CREATE TRIGGER update_%I_updated_at
					BEFORE UPDATE ON %I
					FOR EACH ROW
					EXECUTE FUNCTION update_updated_at_column();
				', t, t, t, t);
			END LOOP;
		END;
		$$;`,
	}

	for i, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("migration failed at step %d: %w", i+1, err)
		}
	}

	return nil
}