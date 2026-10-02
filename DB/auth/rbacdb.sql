CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    role VARCHAR(100),
    "group" VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (email, password_hash, is_active, role, "group", created_at, updated_at)
VALUES
    ('admin@rbac.local', 'admin123', TRUE, 'admin', 'core', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('jane.doe@rbac.local', 'password123', TRUE, 'manager', 'operations', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('john.smith@rbac.local', 'secret456', TRUE, 'user', 'support', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

