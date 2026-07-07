DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM products WHERE status = 'published') THEN
    RAISE EXCEPTION 'cannot rollback publish products migration while published products exist';
  END IF;
END
$$;

DROP INDEX IF EXISTS products_published_created_idx;

ALTER TABLE products DROP CONSTRAINT products_status_check;

ALTER TABLE products
  ADD CONSTRAINT products_status_check CHECK (status IN ('draft'));
