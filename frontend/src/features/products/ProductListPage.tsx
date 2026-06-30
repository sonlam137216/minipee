import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { useAuth } from "../auth/AuthContext";
import { Product, listProducts } from "../../shared/api";

export function ProductListPage() {
  const auth = useAuth();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    async function load() {
      if (auth.token === null) {
        return;
      }
      setLoading(true);
      setError(null);
      try {
        const response = await listProducts(auth.token);
        if (active) {
          setProducts(response.products);
        }
      } catch (err) {
        if (active) {
          setError(err instanceof Error ? err.message : "Failed to load products");
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
  }, [auth.token]);

  return (
    <section className="panel">
      <div className="section-header">
        <h1>Your products</h1>
        <Link to="/products/new">Create product</Link>
      </div>
      {loading ? <p>Loading products...</p> : null}
      {error !== null ? <p role="alert">{error}</p> : null}
      {!loading && error === null && products.length === 0 ? <p>No products yet.</p> : null}
      {products.length > 0 ? (
        <ul className="product-list">
          {products.map((product) => (
            <li key={product.id}>
              <Link to={`/products/${product.id}`}>{product.name}</Link>
              <span>{product.status}</span>
            </li>
          ))}
        </ul>
      ) : null}
    </section>
  );
}
