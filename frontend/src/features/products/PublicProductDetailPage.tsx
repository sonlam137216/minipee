import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { PublicProduct, getPublicProduct } from "../../shared/api";

export function PublicProductDetailPage() {
  const { productID } = useParams();
  const [product, setProduct] = useState<PublicProduct | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    async function load() {
      if (productID === undefined) {
        return;
      }
      setLoading(true);
      setError(null);
      try {
        const response = await getPublicProduct(productID);
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
  }, [productID]);

  return (
    <section className="panel">
      <Link to="/catalog">Back to catalog</Link>
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
              <dt>Published</dt>
              <dd>{product.updatedAt}</dd>
            </div>
          </dl>
        </article>
      ) : null}
    </section>
  );
}
