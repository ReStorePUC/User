CREATE DATABASE IF NOT EXISTS userdb;

USE userdb;

CREATE TABLE users (
    id INT(6) AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(50) unique,
    password VARCHAR(200),
    is_admin BOOLEAN DEFAULT FALSE
)