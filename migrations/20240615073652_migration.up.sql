CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password VARCHAR(255) NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       gender VARCHAR(50),
                       age INT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL,
                       INDEX idx_email (email),
                       INDEX idx_name (name),
                       INDEX idx_age (age),
                       INDEX idx_created_at (created_at),
                       INDEX idx_updated_at (updated_at),
                       INDEX idx_deleted_at (deleted_at)
);