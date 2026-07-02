const assetOrigin = (() => {
  const configured = process.env.NEXT_PUBLIC_ASSET_ORIGIN;
  if (configured) return configured.replace(/\/$/, '');

  const apiURL = process.env.NEXT_PUBLIC_API_URL;
  if (apiURL) {
    try {
      return new URL(apiURL).origin;
    } catch {
      // Fall through to browser origin.
    }
  }

  if (typeof window !== 'undefined') {
    return ['localhost', '127.0.0.1'].includes(window.location.hostname)
      ? window.location.origin
      : 'https://api.zioran.com';
  }
  return '';
})();

export const assetUrl = (url?: string) => {
  if (!url) return '';
  if (/^(https?:)?\/\//.test(url) || url.startsWith('data:') || url.startsWith('blob:')) return url;
  if (url.startsWith('/uploads/')) return `${assetOrigin}${url}`;
  return url;
};

export const normalizeAssetUrls = <T>(value: T): T => {
  if (typeof value === 'string') return assetUrl(value) as T;
  if (Array.isArray(value)) return value.map((item) => normalizeAssetUrls(item)) as T;
  if (value && typeof value === 'object') {
    return Object.fromEntries(
      Object.entries(value as Record<string, unknown>).map(([key, val]) => [key, normalizeAssetUrls(val)])
    ) as T;
  }
  return value;
};
