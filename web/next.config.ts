import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    // Default API URL - can be overridden via NEXT_PUBLIC_API_URL env var
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  },
};

export default nextConfig;
