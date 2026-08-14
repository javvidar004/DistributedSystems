CREATE TABLE IF NOT EXISTS public.users (
	id INTEGER SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	username VARCHAR(100) NOT NULL UNIQUE,
	password VARCHAR(100) NOT NULL
);

INSERT INTO public.users (id, name, username, password)
VALUES
	(1, 'Alice Johnson', 'alice', 'alice123'),
	(2, 'Bob Smith', 'bob', 'bob123'),
	(3, 'Carol Davis', 'carol', 'carol123'),
	(4, 'David Brown', 'david', 'david123'),
	(5, 'Eva Martinez', 'eva', 'eva123')
ON CONFLICT (id) DO NOTHING;
