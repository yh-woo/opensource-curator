import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  rewrites: async () => [
    {
      source: "/v1/:path*",
      destination: "http://localhost:8080/v1/:path*",
    },
  ],
};

export default nextConfig;
