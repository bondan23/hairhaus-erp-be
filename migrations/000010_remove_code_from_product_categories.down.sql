ALTER TABLE product_categories ADD COLUMN code VARCHAR(255);
UPDATE product_categories SET code = name; -- Fallback to name
ALTER TABLE product_categories ALTER COLUMN code SET NOT NULL;
CREATE UNIQUE INDEX idx_product_categories_code ON product_categories(code);
