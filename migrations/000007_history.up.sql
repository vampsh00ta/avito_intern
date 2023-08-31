create table if not exists  history (

      id serial primary key ,
      user_id int,
      slug varchar(255) ,
      operation varchar(20),
      update_time timestamp

);

