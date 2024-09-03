CREATE TABLE tasks (
    id int AUTO_INCREMENT primary key,
    title varchar(255),
    description text,
    status BOOLEAN DEFAULT FALSE
);