import { cleanup, render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";

import { App } from "./App";

const seller = {
  id: "00000000-0000-4000-8000-000000000001",
  email: "seller@example.com",
  displayName: "Seller",
  createdAt: "2026-06-30T00:00:00Z",
  updatedAt: "2026-06-30T00:00:00Z"
};

const draftProduct = {
  id: "00000000-0000-4000-8000-000000000101",
  sellerId: seller.id,
  name: "Draft product",
  description: "Draft description",
  status: "draft",
  createdAt: "2026-06-30T00:00:00Z",
  updatedAt: "2026-06-30T00:00:00Z"
};

const publishedSellerProduct = {
  ...draftProduct,
  name: "Published product",
  status: "published",
  updatedAt: "2026-06-30T00:05:00Z"
};

const publicProduct = {
  id: publishedSellerProduct.id,
  name: publishedSellerProduct.name,
  description: publishedSellerProduct.description,
  status: "published",
  createdAt: publishedSellerProduct.createdAt,
  updatedAt: publishedSellerProduct.updatedAt
};

describe("seller app", () => {
  beforeEach(() => {
    localStorage.clear();
    window.history.pushState({}, "", "/");
    vi.stubGlobal("fetch", vi.fn());
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
    localStorage.clear();
  });

  test("login form submission stores token and calls login API", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ seller, accessToken: "token-1" }), {
          status: 200,
          headers: { "Content-Type": "application/json" }
        })
      )
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ products: [] }), {
          status: 200,
          headers: { "Content-Type": "application/json" }
        })
      )
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ products: [] }), {
          status: 200,
          headers: { "Content-Type": "application/json" }
        })
      );
    window.history.pushState({}, "", "/login");

    render(<App />);

    await user.type(screen.getByLabelText(/email/i), "seller@example.com");
    await user.type(screen.getByLabelText(/password/i), "password123");
    await user.click(screen.getByRole("button", { name: /log in/i }));

    await waitFor(() => expect(localStorage.getItem("marketplace.accessToken")).toBe("token-1"));
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe("http://localhost:8080/api/v1/auth/login");
    expect(init?.method).toBe("POST");
    expect(init?.body).toBe('{"email":"seller@example.com","password":"password123"}');
    expect(localStorage.getItem("marketplace.seller")).toBe(JSON.stringify(seller));
  });

  test("seller pages require authentication", async () => {
    window.history.pushState({}, "", "/products");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /seller login/i })).toBeInTheDocument();
  });

  test("malformed stored seller JSON is cleared and redirects to login", async () => {
    const fetchMock = vi.mocked(fetch);
    localStorage.setItem("marketplace.accessToken", "test-token");
    localStorage.setItem("marketplace.seller", "{malformed");
    window.history.pushState({}, "", "/products");

    expect(() => render(<App />)).not.toThrow();

    expect(await screen.findByRole("heading", { name: /seller login/i })).toBeInTheDocument();
    expect(localStorage.getItem("marketplace.accessToken")).toBeNull();
    expect(localStorage.getItem("marketplace.seller")).toBeNull();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  test("stored auth without a token is cleared and redirects to login", async () => {
    const fetchMock = vi.mocked(fetch);
    localStorage.setItem("marketplace.seller", JSON.stringify(seller));
    window.history.pushState({}, "", "/products");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /seller login/i })).toBeInTheDocument();
    expect(localStorage.getItem("marketplace.accessToken")).toBeNull();
    expect(localStorage.getItem("marketplace.seller")).toBeNull();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  test("stored auth without seller data is cleared and redirects to login", async () => {
    const fetchMock = vi.mocked(fetch);
    localStorage.setItem("marketplace.accessToken", "test-token");
    window.history.pushState({}, "", "/products");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /seller login/i })).toBeInTheDocument();
    expect(localStorage.getItem("marketplace.accessToken")).toBeNull();
    expect(localStorage.getItem("marketplace.seller")).toBeNull();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  test("stored auth with invalid seller shape is cleared and redirects to login", async () => {
    const fetchMock = vi.mocked(fetch);
    localStorage.setItem("marketplace.accessToken", "test-token");
    localStorage.setItem(
      "marketplace.seller",
      JSON.stringify({ email: "seller@example.com", displayName: "Seller", createdAt: "2026-06-30T00:00:00Z" })
    );
    window.history.pushState({}, "", "/products");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /seller login/i })).toBeInTheDocument();
    expect(localStorage.getItem("marketplace.accessToken")).toBeNull();
    expect(localStorage.getItem("marketplace.seller")).toBeNull();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  test("valid stored auth renders protected seller page and loads products", async () => {
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ products: [] }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      })
    );
    localStorage.setItem("marketplace.accessToken", "test-token");
    localStorage.setItem("marketplace.seller", JSON.stringify(seller));
    window.history.pushState({}, "", "/products");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /your products/i })).toBeInTheDocument();
    expect(await screen.findByText(/no products yet/i)).toBeInTheDocument();
    expect(localStorage.getItem("marketplace.accessToken")).toBe("test-token");
    expect(localStorage.getItem("marketplace.seller")).toBe(JSON.stringify(seller));
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe("http://localhost:8080/api/v1/seller/products");
    expect(init?.method).toBe("GET");
    expect((init?.headers as Headers).get("Authorization")).toBe("Bearer test-token");
  });

  test("logout clears stored auth and returns to unauthenticated navigation", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ products: [] }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      })
    );
    localStorage.setItem("marketplace.accessToken", "test-token");
    localStorage.setItem("marketplace.seller", JSON.stringify(seller));
    window.history.pushState({}, "", "/products");

    render(<App />);

    await screen.findByRole("heading", { name: /your products/i });
    await user.click(screen.getByRole("button", { name: /log out/i }));

    expect(localStorage.getItem("marketplace.accessToken")).toBeNull();
    expect(localStorage.getItem("marketplace.seller")).toBeNull();
    expect(screen.queryByRole("button", { name: /log out/i })).not.toBeInTheDocument();
  });

  test("create product form validates product name before API request", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    localStorage.setItem("marketplace.accessToken", "token-1");
    localStorage.setItem("marketplace.seller", JSON.stringify(seller));
    window.history.pushState({}, "", "/products/new");

    render(<App />);

    await user.type(screen.getByLabelText(/name/i), "ab");
    await user.click(screen.getByRole("button", { name: /create draft/i }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Product name must contain between 3 and 200 characters");
    expect(fetchMock).not.toHaveBeenCalled();
  });

  test("seller can publish a draft product from the detail page", async () => {
    const user = userEvent.setup();
    const fetchMock = vi.mocked(fetch);
    fetchMock
      .mockResolvedValueOnce(
        new Response(JSON.stringify(draftProduct), {
          status: 200,
          headers: { "Content-Type": "application/json" }
        })
      )
      .mockResolvedValueOnce(
        new Response(JSON.stringify(publishedSellerProduct), {
          status: 200,
          headers: { "Content-Type": "application/json" }
        })
      );
    localStorage.setItem("marketplace.accessToken", "test-token");
    localStorage.setItem("marketplace.seller", JSON.stringify(seller));
    window.history.pushState({}, "", `/products/${draftProduct.id}`);

    render(<App />);

    expect(await screen.findByRole("heading", { name: /draft product/i })).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: /publish product/i }));

    await waitFor(() => expect(screen.getByText("published")).toBeInTheDocument());
    const [url, init] = fetchMock.mock.calls[1];
    expect(url).toBe(`http://localhost:8080/api/v1/seller/products/${draftProduct.id}/publish`);
    expect(init?.method).toBe("POST");
    expect((init?.headers as Headers).get("Authorization")).toBe("Bearer test-token");
  });

  test("catalog lists public published products without requiring auth", async () => {
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify({ products: [publicProduct] }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      })
    );
    window.history.pushState({}, "", "/catalog");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /catalog/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /published product/i })).toHaveAttribute(
      "href",
      `/catalog/${publicProduct.id}`
    );
    expect(fetchMock).toHaveBeenCalledWith(
      "http://localhost:8080/api/v1/products",
      expect.objectContaining({ method: "GET" })
    );
    expect(screen.queryByText(/draft product/i)).not.toBeInTheDocument();
  });

  test("catalog detail shows a public published product", async () => {
    const fetchMock = vi.mocked(fetch);
    fetchMock.mockResolvedValueOnce(
      new Response(JSON.stringify(publicProduct), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      })
    );
    window.history.pushState({}, "", `/catalog/${publicProduct.id}`);

    render(<App />);

    expect(await screen.findByRole("heading", { name: /published product/i })).toBeInTheDocument();
    expect(screen.getByText("Draft description")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith(
      `http://localhost:8080/api/v1/products/${publicProduct.id}`,
      expect.objectContaining({ method: "GET" })
    );
  });
});
