import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

const nextConfig: NextConfig = {
  output: "standalone",
  rewrites: async () => [
    {
      source: "/v1/:path*",
      destination: "http://localhost:8080/v1/:path*",
    },
  ],
};

export default withNextIntl(nextConfig);
