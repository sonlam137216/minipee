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
  });

  test("seller pages require authentication", async () => {
    window.history.pushState({}, "", "/products");

    render(<App />);

    expect(await screen.findByRole("heading", { name: /seller login/i })).toBeInTheDocument();
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
});
