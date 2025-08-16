-- Insert users
INSERT INTO "users" ("id", "name", "cpf", "email", "password", "created_at", "updated_at") VALUES
(1, 'Alice Silva', '123.456.789-00', 'alice@example.com', 'senha123', NOW(), NOW()),
(2, 'Bob Santos', '987.654.321-00', 'bob@example.com', 'senha456', NOW(), NOW());

-- Insert accounts
INSERT INTO "accounts" ("id", "owner_id", "owner_type", "balance", "can_send", "can_receive", "active", "created_at", "updated_at") VALUES
(1, 1, 'user', 1000.00, TRUE, TRUE, TRUE, NOW(), NOW()),
(2, 2, 'user', 500.00, TRUE, TRUE, TRUE, NOW(), NOW());
