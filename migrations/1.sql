USE userdb;

CREATE TABLE profiles (
     id INT(6) AUTO_INCREMENT PRIMARY KEY,
     name VARCHAR(100),
     address varchar(100),
     block varchar(100),
     zip_code varchar(100),
     city varchar(100),
     state varchar(100),
     user_id INT(6),
     FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE stores (
    id INT(6) AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100),
    address varchar(100),
    block varchar(100),
    city varchar(100),
    state varchar(100),
    photo_path varchar(100),
    user_id INT(6),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
