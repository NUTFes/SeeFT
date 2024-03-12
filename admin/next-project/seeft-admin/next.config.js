/** @type {import('next').NextConfig} */
const isProd = process.env.NEXT_PUBLIC_APP_ENV === 'production';

module.exports = {
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback.fs = false;
      config.resolve.fallback.child_process = false;
      config.resolve.fallback.net = false;
      config.resolve.fallback.dns = false;
      config.resolve.fallback.tls = false;
    }
    return config;
  },
  env: {
    SSR_API_URI: isProd ? 'https://seeft-api.nutfes.net' : 'http://nutfes-seeft-api:1234',
    CSR_API_URI: isProd ? 'https://seeft-api.nutfes.net' : 'http://localhost:1234'
  }
};
