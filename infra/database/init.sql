CREATE TABLE IF NOT EXISTS users (
    id text primary key default gen_random_uuid(),
    username text unique not null,
    password text unique not null,
    created_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS bookmarks (
    id text primary key default gen_random_uuid(),
    owner_id text not null references users(id),
    title text unique not null,
    url text not null,
    tags text,
    description text,
    read boolean default false,
    created_at timestamp default now()
);

CREATE TABLE IF NOT EXISTS apiKeys (
    id text primary key default gen_random_uuid(),
    value text unique not null default gen_random_uuid(),
    active boolean default true,
    created_at timestamp default now()
);
