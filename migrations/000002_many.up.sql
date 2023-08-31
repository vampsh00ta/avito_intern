create table user_segment (
      user_id int,
      segment_id int,
      CONSTRAINT pk_user_segment PRIMARY KEY (user_id, segment_id) ,
      CONSTRAINT fk_user_id
          FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
      CONSTRAINT fk_segment_id
          FOREIGN KEY (segment_id) REFERENCES segments (id) ON DELETE CASCADE ON UPDATE CASCADE
);


