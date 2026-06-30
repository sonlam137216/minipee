import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { useAuth } from "../auth/AuthContext";
import { Product, getProduct } from "../../shared/api";

export function ProductDetailPage() {
  const auth = useAuth();
  const { productID } = useParams();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    async function load() {
      if (auth.token === null || productID === undefined) {
        return;
      }
      setLoading(true);
      setError(null);
      try {
        const response = await getProduct(auth.token, productID);
        if (active) {
          setProduct(response);
        }
      } catch (err) {
        if (active) {
          setError(err instanceof Error ? err.message : "Failed to load product");
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
  }, [auth.token, productID]);

  return (
    <section className="panel">
      <Link to="/products">Back to products</Link>
      {loading ? <p>Loading product...</p> : null}
      {error !== null ? <p role="alert">{error}</p> : null}
      {product !== null ? (
        <article>
          <p className="status">{product.status}</p>
          <h1>{product.name}</h1>
          <p>{product.description === "" ? "No description provided." : product.description}</p>
          <dl>
            <div>
              <dt>Product ID</dt>
              <dd>{product.id}</dd>
            </div>
            <div>
              <dt>Created</dt>
              <dd>{product.createdAt}</dd>
            </div>
          </dl>
        </article>
      ) : null}
    </section>
  );
}
