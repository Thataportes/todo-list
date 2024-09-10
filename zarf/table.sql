CREATE TABLE task (
    id int AUTO_INCREMENT primary key,
    title varchar(255),
    description text,
    created_at DATETIME null,
    finished_at DATETIME null
);