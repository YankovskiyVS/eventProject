		CREATE TABLE IF NOT EXISTS tickets (
			id SERIAL PRIMARY KEY,
			uuid VARCHAR(36) UNIQUE NOT NULL,
			event_id INTEGER NOT NULL,
			order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL