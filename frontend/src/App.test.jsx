import { describe, it, expect, vi, afterEach } from "vitest";
import { render, screen, waitFor, cleanup } from "@testing-library/react";
import App from "./App.jsx";

function mockFetch(body, ok = true, status = 200) {
  return vi.fn().mockResolvedValue({
    ok,
    status,
    json: async () => body,
  });
}

describe("App", () => {
  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
  });

  it("shows backend message + green bg when safe", async () => {
    vi.stubGlobal(
      "fetch",
      mockFetch({ country: "ID", safe: true, reason: "all clear", message: "YES 🚀 go" })
    );
    const { container } = render(<App />);
    await waitFor(() => screen.getByText(/YES/));
    expect(container.querySelector(".screen").style.backgroundColor).toBe("rgb(22, 163, 74)");
  });

  it("shows red bg when not safe", async () => {
    vi.stubGlobal(
      "fetch",
      mockFetch({ country: "ID", safe: false, reason: "it's the weekend", message: "NO 🏖️" })
    );
    const { container } = render(<App />);
    await waitFor(() => screen.getByText(/NO/));
    expect(container.querySelector(".screen").style.backgroundColor).toBe("rgb(220, 38, 38)");
  });

  it("shows error on failed fetch", async () => {
    vi.stubGlobal("fetch", mockFetch(null, false, 500));
    render(<App />);
    await waitFor(() => screen.getByText(/HTTP 500/));
  });
});
