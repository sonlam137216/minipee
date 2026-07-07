import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { PublicProduct, listPublicProducts } from "../../shared/api";

export function PublicProductListPage() {
  const [products, setProducts] = useState<PublicProduct[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    async function load() {
      setLoading(true);
      setError(null);
      try {
        const response = await listPublicProducts();
        if (active) {
          setProducts(response.products);
        }
      } catch (err) {
        if (active) {
          setError(err instanceof Error ? err.message : "Failed to load catalog");
        }
      } finally {
        if (active) {
          setLoading(false);
        }
      }
    }
    void load();
    return () => {
      active = false;
    };
  }, []);

  return (
    <section className="panel">
      <h1>Catalog</h1>
      {loading ? <p>Loading catalog...</p> : null}
      {error !== null ? <p role="alert">{error}</p> : null}
      {!loading && error === null && products.length === 0 ? <p>No published products yet.</p> : null}
      {products.length > 0 ? (
        <ul className="product-list">
          {products.map((product) => (
            <li key={product.id}>
              <Link to={`/catalog/${product.id}`}>{product.name}</Link>
              <span>{product.status}</span>
            </li>
          ))}
        </ul>
      ) : null}
    </section>
  );
}
