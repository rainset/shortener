create table if not exists urls
(
    id       serial,
    hash     varchar(255) not null primary key,
    original varchar(255) not null
);