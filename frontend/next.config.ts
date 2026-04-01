import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Allow hot module reloading when testing over the LAN network ip
  allowedDevOrigins: ['192.168.1.176'],
};

export default nextConfig;
