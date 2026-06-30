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

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(tokenKey));
  const [seller, setSeller] = useState<Seller | null>(() => {
    const stored = localStorage.getItem(sellerKey);
    return stored === null ? null : (JSON.parse(stored) as Seller);
  });

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
        localStorage.removeItem(tokenKey);
        localStorage.removeItem(sellerKey);
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
