ALTER TABLE products DROP CONSTRAINT products_status_check;

ALTER TABLE products
  ADD CONSTRAINT products_status_check CHECK (status IN ('draft', 'published'));

CREATE INDEX products_published_created_idx
  ON products (created_at DESC, id DESC)
  WHERE status = 'published';
