CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    last_name VARCHAR(255),
    work_position VARCHAR(255),
    salary NUMERIC(12, 2),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (id, email, name, last_name, work_position, salary, created_at, updated_at)
VALUES
    ('8093c97f-b943-48be-a1a4-528cb355bb01', 'admin@rbac.local', 'Admin', 'User', 'System Administrator', 95000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('c50a5d1e-d594-4365-8992-61ca476e19d0', 'jane.doe@rbac.local', 'Jane', 'Doe', 'Operations Manager', 72000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('59659bd6-83cc-4ee8-b277-4a5960e365ad', 'john.smith@rbac.local', 'John', 'Smith', 'Support Analyst', 48000.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO NOTHING;

