CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TEXT,
    username VARCHAR(255),
    action TEXT,
    status TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO logs (timestamp, username, action, status, created_at)
VALUES
    ('2026-09-14 08:00:00', 'admin@rbac.local', 'login', 'success', CURRENT_TIMESTAMP),
    ('2026-09-14 08:15:00', 'jane.doe@rbac.local', 'create_user', 'success', CURRENT_TIMESTAMP),
    ('2026-09-14 08:30:00', 'john.smith@rbac.local', 'update_profile', 'success', CURRENT_TIMESTAMP);

