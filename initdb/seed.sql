-- Insert users
INSERT INTO "users" ("id", "name", "cpf", "email", "password", "created_at", "updated_at") VALUES
(1, 'Alice Silva', '123.456.789-00', 'alice@example.com', 'senha123', NOW(), NOW()),
(2, 'Bob Santos', '987.654.321-00', 'bob@example.com', 'senha456', NOW(), NOW());

-- Insert accounts
INSERT INTO "accounts" ("id", "owner_id", "owner_type", "balance", "can_send", "can_receive", "active", "created_at", "updated_at") VALUES
(1, 1, 'user', 1000.00, TRUE, TRUE, TRUE, NOW(), NOW()),
(2, 2, 'user', 500.00, TRUE, TRUE, TRUE, NOW(), NOW());


INSERT INTO "stores" ("id", "name", "cnpj", "email", "password", "created_at", "updated_at") VALUES
(1, 'Loja Alpha', '12.345.678/0001-90', 'alpha@loja.com', 'senhaAlpha', NOW(), NOW()),
(2, 'Loja Beta', '98.765.432/0001-09', 'beta@loja.com', 'senhaBeta', NOW(), NOW());

INSERT INTO "accounts" ("id", "owner_id", "owner_type", "balance", "can_send", "can_receive", "active", "created_at", "updated_at") VALUES
(3, 1, 'store', 1200.00, TRUE, TRUE, TRUE, NOW(), NOW()),
(4, 2, 'store', 1500.00, TRUE, TRUE, TRUE, NOW(), NOW());