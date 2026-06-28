const apiBase = process.env.NEXT_PUBLIC_API_URL || 'http://127.0.0.1:8080/api/v1';
const apiOrigin = apiBase.replace(/\/api\/v1\/?$/, '').replace(/\/$/, '');

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  images: {
    remotePatterns: [{ protocol: 'https', hostname: '**' }],
  },
  async rewrites() {
    return [
      { source: '/uploads/:path*', destination: `${apiOrigin}/uploads/:path*` },
    ];
  },
};

module.exports = nextConfig;
