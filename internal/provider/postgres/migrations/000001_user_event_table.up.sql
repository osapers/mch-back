BEGIN;

CREATE SCHEMA IF NOT EXISTS mch;

CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE mch.user
(
    id          uuid primary key DEFAULT uuid_generate_v4(),
    first_name  text,
    last_name   text,
    middle_name text,
    email       text not null unique,
    photo       text,
    about       text,
    phone       text,
    password    text not null,
    tags        text[]
);

CREATE TABLE mch.event
(
    id                uuid primary key DEFAULT uuid_generate_v4(),
    name              text not null,
    image             text,
    category          text not null,
    date              bigint,
    short_description text,
    description       text,
    address           jsonb,
    website           text,
    email             text
);

CREATE TABLE mch.user_event
(
    user_id  uuid references mch.user (id) ON DELETE CASCADE,
    event_id uuid references mch.event (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, event_id)
);

CREATE TABLE mch.tag
(
    id   bigint primary key,
    name text not null unique
);

COMMIT;