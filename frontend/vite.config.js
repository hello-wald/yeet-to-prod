/// <reference types="vitest/config" />
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// Pages serves this project under /yeet-to-prod/. Dev + preview stay at / so
// local URLs are simple; only the production build gets the subpath.
export default defineConfig(({ command }) => ({
  base: command === "build" ? "/yeet-to-prod/" : "/",
  plugins: [react()],
  test: {
    environment: "jsdom",
  },
}));
