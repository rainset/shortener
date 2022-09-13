create table if not exists urls
(
    id       serial,
    hash     varchar(255) not null primary key,
    original varchar(255) not null
);

create unique index if not exists urls_original_uindex
    on urls (original);