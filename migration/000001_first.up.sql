CREATE TABLE IF NOT EXISTS table_first
(
    id serial unique primary key,
    text varchar(255) not null unique ,
    order_id serial not null
);
ALTER SEQUENCE table_first_order_id_seq RESTART WITH 4615793;