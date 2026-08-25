CREATE TABLE IF NOT EXISTS offices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workspaces (
    id SERIAL PRIMARY KEY,
    number VARCHAR(20) NOT NULL,
    level INTEGER NOT NULL,
    office_id INTEGER NOT NULL REFERENCES offices(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workspace_id INTEGER NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    booking_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO offices (id, name) VALUES
    (1, 'Centro de Operaciones'),
    (2, 'Oficina Norte')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, name, last_name, username, password, email)
VALUES
    (1, 'Ana', 'García', 'ana', 'ana123', 'ana@company.com'),
    (2, 'Luis', 'Pérez', 'luis', 'luis123', 'luis@company.com'),
    (3, 'María', 'López', 'maria', 'maria123', 'maria@company.com')
ON CONFLICT (id) DO NOTHING;

INSERT INTO workspaces (id, number, level, office_id)
VALUES
    (1, 'A1', 1, 1), (2, 'A2', 1, 1), (3, 'A3', 1, 1), (4, 'A4', 1, 1),
    (5, 'B1', 2, 1), (6, 'B2', 2, 1), (7, 'B3', 2, 1), (8, 'B4', 2, 1),
    (9, 'C1', 3, 1), (10, 'C2', 3, 1), (11, 'C3', 3, 1), (12, 'C4', 3, 1),
    (13, 'N1', 1, 2), (14, 'N2', 1, 2), (15, 'N3', 1, 2), (16, 'N4', 1, 2),
    (17, 'N5', 2, 2), (18, 'N6', 2, 2), (19, 'N7', 2, 2), (20, 'N8', 2, 2)
ON CONFLICT (id) DO NOTHING;

INSERT INTO bookings (id, user_id, workspace_id, booking_date)
VALUES
    (1, 1, 3, CURRENT_DATE + INTERVAL '2 day'),
    (2, 2, 7, CURRENT_DATE + INTERVAL '4 day'),
    (3, 3, 12, CURRENT_DATE + INTERVAL '6 day'),
    (4, 1, 17, CURRENT_DATE + INTERVAL '9 day')
ON CONFLICT (id) DO NOTHING;
