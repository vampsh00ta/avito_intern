create index if not  exists concurrently segment_slug_idx on segments using btree (slug) where slug is not null;