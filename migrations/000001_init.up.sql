BEGIN;
create table IF  not EXISTS users (
    id serial primary key,
    username varchar(255),
    last_update  timestamp
);
create table  IF  not EXISTS segments (
                       id serial primary key,
                       slug varchar(50)
);
commit;


