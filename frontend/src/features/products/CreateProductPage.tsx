import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";

import { useAuth } from "../auth/AuthContext";
import { createProduct } from "../../shared/api";

export function CreateProductPage() {
  const auth = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    if (trimmedName.length < 3 || trimmedName.length > 200) {
      setError("Product name must contain between 3 and 200 characters");
      return;
    }
    if (auth.token === null) {
      setError("Authentication required");
      return;
    }
    setError(null);
    setLoading(true);
    try {
      const product = await createProduct(auth.token, { name: trimmedName, description });
      navigate(`/products/${product.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create product");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="panel">
      <h1>Create draft product</h1>
      <form onSubmit={submit}>
        <label>
          Name
          <input value={name} onChange={(event) => setName(event.target.value)} required minLength={3} maxLength={200} />
        </label>
        <label>
          Description
          <textarea value={description} onChange={(event) => setDescription(event.target.value)} rows={5} />
        </label>
        {error !== null ? <p role="alert">{error}</p> : null}
        <button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create draft"}
        </button>
      </form>
    </section>
  );
}
