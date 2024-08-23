CREATE TABLE task (
    id int AUTO_INCREMENT primary key,
    title varchar(255),
    description text,
    status ENUM('peding', 'completed') DEFAULT 'peding'
);