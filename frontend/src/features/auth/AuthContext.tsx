import { createContext, ReactNode, useContext, useMemo, useState } from "react";

import { AuthResponse, Seller, loginSeller, registerSeller } from "../../shared/api";

type AuthState = {
  seller: Seller | null;
  token: string | null;
  login: (input: { email: string; password: string }) => Promise<void>;
  register: (input: { email: string; password: string; displayName: string }) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthState | null>(null);
const tokenKey = "marketplace.accessToken";
const sellerKey = "marketplace.seller";

function clearStoredAuth() {
  localStorage.removeItem(tokenKey);
  localStorage.removeItem(sellerKey);
}

function isSeller(value: unknown): value is Seller {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const candidate = value as Record<string, unknown>;
  return (
    typeof candidate.id === "string" &&
    typeof candidate.email === "string" &&
    typeof candidate.displayName === "string" &&
    typeof candidate.createdAt === "string" &&
    typeof candidate.updatedAt === "string"
  );
}

function readStoredAuth(): Pick<AuthState, "seller" | "token"> {
  const storedToken = localStorage.getItem(tokenKey);
  const storedSeller = localStorage.getItem(sellerKey);

  if (storedToken === null && storedSeller === null) {
    return { token: null, seller: null };
  }

  if (storedToken === null || storedToken.trim() === "" || storedSeller === null) {
    clearStoredAuth();
    return { token: null, seller: null };
  }

  try {
    const parsedSeller: unknown = JSON.parse(storedSeller);
    if (!isSeller(parsedSeller)) {
      clearStoredAuth();
      return { token: null, seller: null };
    }

    return { token: storedToken, seller: parsedSeller };
  } catch {
    clearStoredAuth();
    return { token: null, seller: null };
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [storedAuth] = useState(readStoredAuth);
  const [token, setToken] = useState<string | null>(storedAuth.token);
  const [seller, setSeller] = useState<Seller | null>(storedAuth.seller);

  function storeAuth(result: AuthResponse) {
    localStorage.setItem(tokenKey, result.accessToken);
    localStorage.setItem(sellerKey, JSON.stringify(result.seller));
    setToken(result.accessToken);
    setSeller(result.seller);
  }

  const value = useMemo<AuthState>(
    () => ({
      seller,
      token,
      login: async (input) => storeAuth(await loginSeller(input)),
      register: async (input) => storeAuth(await registerSeller(input)),
      logout: () => {
        clearStoredAuth();
        setToken(null);
        setSeller(null);
      }
    }),
    [seller, token]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthState {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
