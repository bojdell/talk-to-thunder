CREATE TABLE snippets (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    state INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    seed_text TEXT NOT NULL,
    generated_text TEXT
)
CHARACTER SET = utf8mb4
COLLATE = utf8mb4_unicode_ci;
