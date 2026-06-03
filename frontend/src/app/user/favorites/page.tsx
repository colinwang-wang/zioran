'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { getUserFavorites, removeFavorite } from '@/lib/services';
import Pagination from '@/components/Pagination';
import type { PaginatedList, FavoriteItem } from '@/types';

export default function FavoritesPage() {
  const [data, setData] = useState<PaginatedList<FavoriteItem>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });

  const fetchData = (page = 1) => { getUserFavorites({ page }).then(setData).catch(() => {}); };
  useEffect(() => { fetchData(); }, []);

  const handleRemove = async (courseId: number) => {
    try {
      await removeFavorite(courseId);
      fetchData(data.page);
    } catch { /* ignore */ }
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">我的收藏</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        {data.items.map((f) => (
          <div key={f.course_id} className="rounded-card overflow-hidden bg-surface relative group">
            <Link href={`/courses/${f.slug}`}>
              {f.cover && <img src={f.cover} alt={f.title} className="w-full aspect-[4/3] object-cover" loading="lazy" />}
              <div className="p-3">
                <p className="text-sm font-semibold line-clamp-2">{f.title}</p>
              </div>
            </Link>
            <button onClick={() => handleRemove(f.course_id)} className="absolute top-2 right-2 w-7 h-7 bg-canvas rounded-full flex items-center justify-center text-xs text-mute opacity-0 group-hover:opacity-100 transition-opacity">✕</button>
          </div>
        ))}
      </div>
      {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无收藏</p>}
      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
