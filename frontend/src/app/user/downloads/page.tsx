'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { getUserDownloads } from '@/lib/services';
import Pagination from '@/components/Pagination';
import type { PaginatedList, DownloadItem } from '@/types';

export default function DownloadsPage() {
  const [data, setData] = useState<PaginatedList<DownloadItem>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });

  const fetchData = (page = 1) => { getUserDownloads({ page }).then(setData).catch(() => {}); };
  useEffect(() => { fetchData(); }, []);

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">我的下载</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        {data.items.map((d) => (
          <div key={d.id} className="rounded-card overflow-hidden bg-surface">
            {d.cover && <img src={d.cover} alt={d.title} className="w-full aspect-[4/3] object-cover" loading="lazy" />}
            <div className="p-3">
              <p className="text-sm font-semibold line-clamp-2">{d.title}</p>
              <p className="text-xs text-mute mt-1">{new Date(d.created_at).toLocaleDateString()}</p>
            </div>
          </div>
        ))}
      </div>
      {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无下载记录</p>}
      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
