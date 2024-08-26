CREATE TABLE tasks (
    id int AUTO_INCREMENT primary key,
    title varchar(255),
    description text,
    status ENUM('pending', 'completed') DEFAULT 'pending'
);