import { BrowserRouter, Link, Navigate, Route, Routes } from "react-router-dom";

import { ProtectedRoute } from "./ProtectedRoute";
import { useAuth, AuthProvider } from "../features/auth/AuthContext";
import { LoginPage } from "../features/auth/LoginPage";
import { RegisterPage } from "../features/auth/RegisterPage";
import { CreateProductPage } from "../features/products/CreateProductPage";
import { ProductDetailPage } from "../features/products/ProductDetailPage";
import { ProductListPage } from "../features/products/ProductListPage";
import { PublicProductDetailPage } from "../features/products/PublicProductDetailPage";
import { PublicProductListPage } from "../features/products/PublicProductListPage";

export function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppShell />
      </BrowserRouter>
    </AuthProvider>
  );
}

export function AppShell() {
  const auth = useAuth();

  return (
    <main className="app-shell">
      <header className="topbar">
        <Link to="/catalog" className="brand">
          Marketplace
        </Link>
        <nav aria-label="Catalog navigation">
          <Link to="/catalog">Catalog</Link>
          {auth.token === null ? <Link to="/login">Seller login</Link> : null}
        </nav>
        {auth.token !== null ? (
          <nav aria-label="Seller navigation">
            <Link to="/products">Products</Link>
            <Link to="/products/new">Create product</Link>
            <button type="button" onClick={auth.logout}>
              Log out
            </button>
          </nav>
        ) : null}
      </header>
      <Routes>
        <Route path="/" element={<Navigate to="/catalog" replace />} />
        <Route path="/catalog" element={<PublicProductListPage />} />
        <Route path="/catalog/:productID" element={<PublicProductDetailPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/products" element={<ProductListPage />} />
          <Route path="/products/new" element={<CreateProductPage />} />
          <Route path="/products/:productID" element={<ProductDetailPage />} />
        </Route>
      </Routes>
    </main>
  );
}
