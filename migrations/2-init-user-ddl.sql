CREATE TABLE IF NOT EXISTS app.users (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    username text NULL,
    email text UNIQUE NOT NULL,
    enabled bool NOT NULL DEFAULT false,
    verified bool NOT NULL DEFAULT false,
    anonymous bool NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_users_id ON app.users USING btree (id);
CREATE INDEX IF NOT EXISTS idx_users_username ON app.users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_users_enabled ON app.users USING btree (enabled);
CREATE INDEX IF NOT EXISTS idx_users_verified ON app.users USING btree (verified);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON app.users USING btree (deleted_at);

CREATE TABLE IF NOT EXISTS app.user_roles (
    user_id uuid NOT NULL,
    name text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (user_id, name),
    CONSTRAINT user_roles_users_fk FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_user_roles_deleted_at ON app.user_roles USING btree (deleted_at);

CREATE TABLE IF NOT EXISTS app.credentials (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL,
    hash text NOT NULL,
    credential_type text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT credentials_users_fk FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_credentials_deleted_at ON app.credentials USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_credentials_user_id ON app.credentials USING btree (user_id);

CREATE TABLE IF NOT EXISTS app.clients (
    client_type text NOT NULL,
    client_secret uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT clients_pkey PRIMARY KEY (client_type)
);
CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON app.clients USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_clients_client_secret ON app.clients USING btree (client_secret);

CREATE TABLE IF NOT EXISTS app.user_confirmations (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT user_confirmations_users_fk FOREIGN KEY (user_id) REFERENCES app.users(id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_user_confirmations_deleted_at ON app.user_confirmations USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_confirmations_user_id ON app.user_confirmations USING btree (user_id);

-- Insert default clients with consistent client types
INSERT INTO app.clients (client_type, client_secret, created_at) 
VALUES 
    ('web', gen_random_uuid(), CURRENT_TIMESTAMP),
    ('android', gen_random_uuid(), CURRENT_TIMESTAMP),
    ('ios', gen_random_uuid(), CURRENT_TIMESTAMP);