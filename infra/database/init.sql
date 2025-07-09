CREATE TABLE IF NOT EXISTS users (
    id uuid primary key default gen_random_uuid(),
    username text unique not null,
    password text unique not null,
    createdAt timestamp default now()
);

CREATE TABLE IF NOT EXISTS bookmarks (
    id uuid primary key default gen_random_uuid(),
    ownerId uuid not null references users(id);
    title text unique not null,
    description text,
    read boolean default false,
    createdAt timestamp default now()
);

CREATE TABLE IF NOT EXISTS apiKeys (
    id uuid primary key default gen_random_uuid(),
    value uuid unique not null default gen_random_uuid(),
    active boolean default true,
    createdAt timestamp default now()
);
