-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE IF NOT EXISTS source
(
    id serial primary key,
    description text,
    crawler_id int,
    crawler_meta json,
    schedule text,
    num_should_matched int
);

CREATE TABLE IF NOT EXISTS events_iteration
(
    id serial primary key,
    source_id int,
    insert_timestamp timestamp,
    body text
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

CREATE TABLE public.events (
	id bigserial NOT NULL,
	source_id int4 NULL,
	"depth" int4 NULL,
	current_node_json json NULL,
	insert_date timestamp NULL,
	business_time timestamp NULL,
	current_full_key text NULL,
	CONSTRAINT events_pkey PRIMARY KEY (id),
	CONSTRAINT events_un UNIQUE (source_id, current_full_key)
);
CREATE INDEX events_by_source_id ON public.events USING btree (source_id, id);

CREATE TABLE public.events_doc (
	current_full_key text NOT NULL PRIMARY KEY,
	body text NOT NULL
);

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
