CREATE TABLE IF NOT EXISTS public.users (
    id uuid NOT NULL,
    username text NULL,
    enabled bool NULL,
    verified bool NULL,
    anonymous bool NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_users_id ON public.users USING btree (id);
CREATE INDEX IF NOT EXISTS idx_users_username ON public.users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON public.users USING btree (deleted_at);

CREATE TABLE IF NOT EXISTS public.user_roles (
    user_id uuid NOT NULL,
    name text NOT NULL,
    created_at timestamptz NULL,
    deleted_at timestamptz NULL,
    PRIMARY KEY (user_id, name)
);
CREATE INDEX IF NOT EXISTS idx_user_roles_deleted_at ON public.user_roles USING btree (deleted_at);

CREATE TABLE IF NOT EXISTS public.credentials (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL,
    hash text NULL,
    credential_type text NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT credentials_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id)   
);
CREATE INDEX IF NOT EXISTS idx_credentials_deleted_at ON public.credentials USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_credentials_user_id ON public.credentials USING btree (user_id);

CREATE TABLE IF NOT EXISTS public.clients (
    client_type text NOT NULL,
    client_secret uuid NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT clients_pkey PRIMARY KEY (client_type)
);
CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON public.clients USING btree (deleted_at);
CREATE INDEX IF NOT EXISTS idx_clients_client_secret ON public.clients USING btree (client_secret);

CREATE TABLE IF NOT EXISTS public.tokens (
    id SERIAL PRIMARY KEY,
    token text NOT NULL,
    token_type text NOT NULL,
    user_id uuid NOT NULL,
    revoked bool DEFAULT false NULL,
    client_type text NOT NULL,
    expires_at timestamptz NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL
);
CREATE INDEX IF NOT EXISTS idx_tokens_token ON public.tokens USING btree (token);                   
ALTER TABLE public.tokens ADD CONSTRAINT tokens_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id);
ALTER TABLE public.tokens ADD CONSTRAINT tokens_clients_fk FOREIGN KEY (client_type) REFERENCES public.clients(client_type);