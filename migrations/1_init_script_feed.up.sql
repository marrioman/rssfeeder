CREATE SCHEMA IF NOT EXISTS feeder;

CREATE TABLE IF NOT EXISTS feeder.sources (
    id serial,
    url text,
    sitename varchar(50) unique,
    created_at timestamptz not null default now(),
    primary key (id)
);

CREATE TABLE IF NOT EXISTS feeder.content (
    id serial,
    item_title text unique,
    link text,
    description text,
    source_id integer,
    created_at timestamptz not null default now(),
    primary key (id)
);

INSERT INTO feeder.sources (url,sitename) VALUES 
('https://pitchfork.com/feed/feed-album-reviews/rss','RSS: Album Reviews'),
('https://meduza.io/rss2/all','Meduza.io');