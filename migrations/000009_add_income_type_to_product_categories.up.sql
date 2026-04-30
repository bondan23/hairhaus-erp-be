ALTER TABLE product_categories ADD COLUMN income_type VARCHAR(255);

-- Backfill existing categories
UPDATE product_categories SET income_type = 'HAIRCUT' WHERE code = 'HC';
UPDATE product_categories SET income_type = 'TREATMENT' WHERE code != 'HC' AND code NOT IN (SELECT DISTINCT code FROM product_categories WHERE income_type IS NOT NULL);
