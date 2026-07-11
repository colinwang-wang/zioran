'use client';

import { useState } from 'react';

interface Resource {
  name: string;
  url: string;
  password?: string;
}

interface DownloadModalProps {
  open: boolean;
  resources: Resource[];
  onClose: () => void;
}

export default function DownloadModal({ open, resources, onClose }: DownloadModalProps) {
  const [copied, setCopied] = useState(false);

  if (!open) return null;

  const handleCopy = () => {
    const text = resources.map(r => `${r.name}: ${r.url}${r.password ? ` 提取码:${r.password}` : ''}`).join('\n');
    navigator.clipboard.writeText(text).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }).catch(() => alert('复制失败，请手动复制'));
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div className="relative mx-4 w-full max-w-md rounded-xl bg-canvas p-6 shadow-xl" onClick={e => e.stopPropagation()}>
        <h3 className="text-lg font-bold text-ink mb-4">资源下载</h3>
        <div className="space-y-3 max-h-60 overflow-y-auto">
          {resources.map((r, idx) => (
            <div key={idx} className="rounded-card bg-surface p-3 border border-hairline">
              <p className="text-sm font-semibold text-ink">{r.name}</p>
              <p className="mt-1 text-xs text-mute break-all">{r.url}</p>
              {r.password && <p className="mt-1 text-xs text-primary">提取码: {r.password}</p>}
            </div>
          ))}
          {resources.length === 0 && <p className="text-sm text-mute text-center py-4">暂无可用资源</p>}
        </div>
        <div className="mt-4 flex justify-end gap-3">
          <button onClick={handleCopy} className="px-4 py-2 text-sm font-bold rounded-card border border-primary text-primary hover:bg-primary/5">
            {copied ? '已复制 ✓' : '复制全部'}
          </button>
          <button onClick={onClose} className="px-4 py-2 text-sm font-bold rounded-card bg-primary text-white hover:bg-primary-pressed">
            确定
          </button>
        </div>
      </div>
    </div>
  );
}
