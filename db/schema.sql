CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(20) NOT NULL,
    last_name  VARCHAR(20),
    email      VARCHAR(50) UNIQUE NOT NULL,
    pass_hash  VARCHAR(60) NOT NULL
);

CREATE TABLE IF NOT EXISTS lists
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(50) NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(20) NOT NULL,
    list_id INT NOT NULL,
    CONSTRAINT fk_list
    FOREIGN KEY (list_id) REFERENCES lists(id) ON DELETE CASCADE
);