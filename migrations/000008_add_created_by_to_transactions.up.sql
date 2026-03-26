ALTER TABLE transactions ADD COLUMN created_by_id UUID NOT NULL;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_created_by FOREIGN KEY (created_by_id) REFERENCES users(id);
