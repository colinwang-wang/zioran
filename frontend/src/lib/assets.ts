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

// --- OSS 图片处理工具 ---

/**
 * 生成缩略图 URL（等比缩放，不超过指定宽高）
 */
export const thumbnailUrl = (url: string, width: number, height?: number) => {
  const resolved = assetUrl(url);
  if (!resolved || !isOssUrl(resolved)) return resolved;
  const h = height || width;
  return appendProcess(resolved, `image/resize,m_lfit,w_${width},h_${h}`);
};

/**
 * 生成居中裁剪 URL（精确裁剪到指定尺寸）
 */
export const cropUrl = (url: string, width: number, height: number) => {
  const resolved = assetUrl(url);
  if (!resolved || !isOssUrl(resolved)) return resolved;
  return appendProcess(resolved, `image/resize,m_fill,w_${width},h_${height}`);
};

/**
 * 生成 webp 格式 URL（可选 resize）
 */
export const webpUrl = (url: string, width?: number, height?: number) => {
  const resolved = assetUrl(url);
  if (!resolved || !isOssUrl(resolved)) return resolved;
  if (width && height) {
    return appendProcess(resolved, `image/resize,m_lfit,w_${width},h_${height}/format,webp`);
  }
  return appendProcess(resolved, 'image/format,webp');
};

/** 检查是否为 OSS/CDN 域名的 URL */
const isOssUrl = (url: string) => {
  // 匹配 OSS 直链或 CDN 域名
  return /\.aliyuncs\.com\//.test(url) || /img\.zioran\.com\//.test(url);
};

/** 为 URL 追加 x-oss-process 参数 */
const appendProcess = (url: string, process: string) => {
  const sep = url.includes('?') ? '&' : '?';
  return `${url}${sep}x-oss-process=${process}`;
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
