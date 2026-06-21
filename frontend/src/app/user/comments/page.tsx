'use client';

import { useEffect, useState } from 'react';
import Pagination from '@/components/Pagination';
import { getUserComments } from '@/lib/services';
import type { CommentItem, PaginatedList } from '@/types';

const targetLabels: Record<string, string> = {
  course: '课程',
};

export default function UserCommentsPage() {
  const [data, setData] = useState<PaginatedList<CommentItem>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });

  const fetchData = (page = 1) => {
    getUserComments({ page, pageSize: 10 }).then(setData).catch(() => {});
  };

  useEffect(() => { fetchData(); }, []);

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">我的评论</h2>
      <div className="space-y-3">
        {data.items.map((item) => (
          <div key={item.id} className="p-4 bg-surface rounded-card">
            <div className="flex items-center justify-between gap-3">
              <p className="text-xs text-mute">
                {targetLabels[item.target_type] || item.target_type} #{item.target_id} · {new Date(item.created_at).toLocaleString()}
              </p>
              {item.status !== 'visible' && <span className="text-xs px-2 py-0.5 rounded-full bg-secondary-bg text-mute">已隐藏</span>}
            </div>
            <p className="mt-2 text-sm text-ink whitespace-pre-wrap">{item.content}</p>
          </div>
        ))}
        {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无评论</p>}
      </div>
      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
