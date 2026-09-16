CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    role VARCHAR(100),
    "group" VARCHAR(100),
    name VARCHAR(255),
    last_name VARCHAR(255),
    work_position VARCHAR(255),
    salary NUMERIC(12, 2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TEXT,
    username VARCHAR(255),
    action TEXT,
    status TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (email, password_hash, is_active, role, "group", name, last_name, work_position, salary, created_at, updated_at)
VALUES
    ('admin@rbac.local', 'admin123', TRUE, 'admin', 'core', 'Admin', 'User', 'System Administrator', 95000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('jane.doe@rbac.local', 'password123', TRUE, 'manager', 'operations', 'Jane', 'Doe', 'Operations Manager', 72000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('john.smith@rbac.local', 'secret456', TRUE, 'user', 'support', 'John', 'Smith', 'Support Analyst', 48000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

INSERT INTO logs (timestamp, username, action, status, created_at)
VALUES
    ('2026-09-14 08:00:00', 'admin@rbac.local', 'login', 'success', CURRENT_TIMESTAMP),
    ('2026-09-14 08:15:00', 'jane.doe@rbac.local', 'create_user', 'success', CURRENT_TIMESTAMP),
    ('2026-09-14 08:30:00', 'john.smith@rbac.local', 'update_profile', 'success', CURRENT_TIMESTAMP);

