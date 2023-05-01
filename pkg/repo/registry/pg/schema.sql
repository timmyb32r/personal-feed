-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS source
(
    id serial primary key,
    description text,
    crawler_id int,
    crawler_meta json,
    schedule text
);

CREATE TABLE IF NOT EXISTS tags
(
    id serial primary key,
    name text
);

CREATE TABLE IF NOT EXISTS tags_link
(
    id serial primary key,
    source_id int,
    tag_id int
);

CREATE TABLE IF NOT EXISTS crawlers
(
    id serial primary key,
    name text
);

CREATE TABLE IF NOT EXISTS events
(
    id bigserial primary key,
    source_id int,
    depth int,
    parent_full_key text,
    current_node_json json,
    insert_date timestamp
);
CREATE INDEX IF NOT EXISTS events_by_source_id ON events (source_id, id);

CREATE TABLE IF NOT EXISTS users
(
    id serial primary key,
    email text,
    tg_chat_id bigint,
    nickname text,
    pass_hash text
);

CREATE TABLE IF NOT EXISTS subscribes
(
    id serial primary key,
    user_id int,
    source_id int,
    details json
        -- root bool
        -- full_paths: []
);

CREATE TABLE IF NOT EXISTS feed
(
    id bigserial primary key,
    user_id int,
    state int,
    event_id bigint
);
CREATE INDEX IF NOT EXISTS feed_by_user_id ON feed (user_id, id);
CREATE INDEX IF NOT EXISTS feed_by_user_id ON feed (user_id, state, id);

CREATE TABLE IF NOT EXISTS cron (
    id integer DEFAULT(1) UNIQUE NOT NULL,
    last_run_time TIMESTAMP NOT NULL,
    CONSTRAINT chk_cron_time_one_row CHECK (id = 1)
);
