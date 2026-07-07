export type Seller = {
  id: string;
  email: string;
  displayName: string;
  createdAt: string;
  updatedAt: string;
};

export type ProductStatus = "draft" | "published";

export type SellerProduct = {
  id: string;
  sellerId: string;
  name: string;
  description: string;
  status: ProductStatus;
  createdAt: string;
  updatedAt: string;
};

export type PublicProduct = {
  id: string;
  name: string;
  description: string;
  status: ProductStatus;
  createdAt: string;
  updatedAt: string;
};

export type AuthResponse = {
  seller: Seller;
  accessToken: string;
};

export type ProductListResponse = {
  products: SellerProduct[];
};

export type PublicProductListResponse = {
  products: PublicProduct[];
};

type ApiErrorResponse = {
  error?: {
    code?: string;
    message?: string;
    fields?: Record<string, string>;
  };
};

export class ApiError extends Error {
  readonly status: number;
  readonly fields: Record<string, string>;

  constructor(status: number, message: string, fields: Record<string, string> = {}) {
    super(message);
    this.status = status;
    this.fields = fields;
  }
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

export async function registerSeller(input: {
  email: string;
  password: string;
  displayName: string;
}): Promise<AuthResponse> {
  return request<AuthResponse>("/api/v1/auth/register", { method: "POST", body: input });
}

export async function loginSeller(input: { email: string; password: string }): Promise<AuthResponse> {
  return request<AuthResponse>("/api/v1/auth/login", { method: "POST", body: input });
}

export async function createProduct(token: string, input: { name: string; description: string }): Promise<SellerProduct> {
  return request<SellerProduct>("/api/v1/seller/products", { method: "POST", token, body: input });
}

export async function listProducts(token: string): Promise<ProductListResponse> {
  return request<ProductListResponse>("/api/v1/seller/products", { method: "GET", token });
}

export async function getProduct(token: string, productID: string): Promise<SellerProduct> {
  return request<SellerProduct>(`/api/v1/seller/products/${encodeURIComponent(productID)}`, { method: "GET", token });
}

export async function publishProduct(token: string, productID: string): Promise<SellerProduct> {
  return request<SellerProduct>(`/api/v1/seller/products/${encodeURIComponent(productID)}/publish`, { method: "POST", token });
}

export async function listPublicProducts(): Promise<PublicProductListResponse> {
  return request<PublicProductListResponse>("/api/v1/products", { method: "GET" });
}

export async function getPublicProduct(productID: string): Promise<PublicProduct> {
  return request<PublicProduct>(`/api/v1/products/${encodeURIComponent(productID)}`, { method: "GET" });
}

async function request<T>(
  path: string,
  options: { method: "GET" | "POST"; token?: string; body?: unknown }
): Promise<T> {
  const headers = new Headers();
  headers.set("Accept", "application/json");
  if (options.body !== undefined) {
    headers.set("Content-Type", "application/json");
  }
  if (options.token !== undefined) {
    headers.set("Authorization", `Bearer ${options.token}`);
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method,
    headers,
    body: options.body === undefined ? undefined : JSON.stringify(options.body)
  });

  if (!response.ok) {
    throw await parseApiError(response);
  }

  return (await response.json()) as T;
}

async function parseApiError(response: Response): Promise<ApiError> {
  try {
    const body = (await response.json()) as ApiErrorResponse;
    return new ApiError(response.status, body.error?.message ?? "Request failed", body.error?.fields ?? {});
  } catch {
    return new ApiError(response.status, "Request failed");
  }
}
