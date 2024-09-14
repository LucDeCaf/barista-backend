-- +goose Up
PRAGMA foreign_keys=OFF;

CREATE TABLE blogs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_username TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(owner_username) REFERENCES users(username)
);

INSERT INTO blogs_new (owner_username,title,content,created_at,updated_at)
SELECT owner_username,title,content,created_at,updated_at FROM blogs;

DROP TABLE blogs;
ALTER TABLE blogs_new RENAME TO blogs;

-- +goose Down
PRAGMA foreign_keys=OFF;

CREATE TABLE blogs_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_username TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY(owner_username) REFERENCES users(username)
);

INSERT INTO blogs_new (owner_username,title,content,created_at,updated_at)
SELECT owner_username,title,content,created_at,updated_at FROM blogs;

DROP TABLE blogs;
ALTER TABLE blogs_new RENAME TO blogs;
