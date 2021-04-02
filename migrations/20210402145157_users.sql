-- +goose Up
CREATE TABLE users(
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(200) NOT NULL,
    nick_name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    password VARCHAR(200) NOT NULL,
    role VARCHAR(200) NOT NULL,

    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;