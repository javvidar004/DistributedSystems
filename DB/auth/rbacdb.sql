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

INSERT INTO users (id, email, password_hash, is_active, role, "group", created_at, updated_at)
VALUES
    ('8093c97f-b943-48be-a1a4-528cb355bb01', 'admin@rbac.local', 'admin123', TRUE, 'admin', 'core', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('c50a5d1e-d594-4365-8992-61ca476e19d0', 'jane.doe@rbac.local', 'password123', TRUE, 'manager', 'operations', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('59659bd6-83cc-4ee8-b277-4a5960e365ad', 'john.smith@rbac.local', 'secret456', TRUE, 'user', 'support', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

