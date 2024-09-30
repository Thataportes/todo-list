CREATE TABLE task (
    id int AUTO_INCREMENT primary key,
    title varchar(255),
    description text,
    created_at DATETIME null,
    finished_at DATETIME null,
    created_by INT NOT NULL,
    assigned_to INT NULL
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status BOOLEAN NOT NULL DEFAULT TRUE, 
    created_at DATETIME NULL,
    last_updated_at DATETIME NULL
);

