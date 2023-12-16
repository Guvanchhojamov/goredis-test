-- CREATE SEQUENCE table_first_order_id_seq START 1 INCREMENT BY 1;
CREATE TABLE IF NOT EXISTS table_first
(
    id serial unique primary key,
    text varchar(255) not null unique ,
    order_id int not null
);

