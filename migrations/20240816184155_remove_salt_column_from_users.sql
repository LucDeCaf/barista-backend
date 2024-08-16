-- +goose Up
ALTER TABLE users DROP salt; 
ALTER TABLE users RENAME password_hash TO password_hash_with_salt;

-- +goose Down
ALTER TABLE users ADD salt TEXT NOT NULL;
ALTER TABLE users RENAME password_hash_with_salt TO password_hash;
