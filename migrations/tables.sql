create table if not exists urls
(
    id       serial,
    hash     varchar(255) not null primary key,
    original varchar(255) not null
);

create unique index if not exists urls_original_uindex
    on urls (original);

create table if not exists user_history_urls
(
    id       serial,
    cookie_id     varchar(255) not null,
    hash varchar(255) not null,
    UNIQUE (cookie_id,hash)
);