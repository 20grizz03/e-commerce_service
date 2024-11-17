create table item(
    ID serial primary key,
    name varchar(250) not null,
    prise integer not null,
    quantity integer not null,
    category varchar(50) not null ,
    info text
)