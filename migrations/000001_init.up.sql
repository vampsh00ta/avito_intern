BEGIN;
create table users (
    id serial primary key,
    username varchar(255),
    last_update  timestamp
);
create table segments (
                       id serial primary key,
                       slug varchar(50)
);
commit;


