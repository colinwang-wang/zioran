'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getGuestbook, createGuestbook, likeGuestbook } from '@/lib/services';
import Pagination from '@/components/Pagination';
import type { GuestbookItem, PaginatedList } from '@/types';

export default function GuestbookPage() {
  const { isLoggedIn } = useAuth();
  const [data, setData] = useState<PaginatedList<GuestbookItem>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchData = async (page = 1) => {
    try {
      const res = await getGuestbook({ page, pageSize: 10 });
      setData(res);
    } catch { /* ignore */ }
  };

  useEffect(() => { fetchData(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    if (!content.trim()) return;
    setLoading(true);
    try {
      await createGuestbook(content.trim());
      setContent('');
      fetchData();
    } catch { /* ignore */ }
    setLoading(false);
  };

  const handleLike = async (id: number) => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      await likeGuestbook(id);
      fetchData(data.page);
    } catch { /* ignore */ }
  };

  return (
    <div className="max-w-container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold text-ink">留言反馈</h1>
      <p className="text-sm text-mute mt-1">需要其他课程可以留言</p>

      <form onSubmit={handleSubmit} className="mt-6">
        <textarea value={content} onChange={(e) => setContent(e.target.value)} placeholder={isLoggedIn ? '发表你的留言...' : '登录后发表留言'} className="w-full p-4 rounded-card bg-surface border border-hairline text-sm resize-none h-28 focus:border-primary outline-none" />
        <button type="submit" disabled={loading || !content.trim()} className="mt-2 px-6 py-2 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">
          提交
        </button>
      </form>

      <div className="mt-8 space-y-4">
        {data.items.map((item) => (
          <div key={item.id} className="p-4 bg-surface rounded-card">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-secondary-bg flex items-center justify-center text-xs font-bold shrink-0">
                {item.avatar ? <img src={item.avatar} alt="" className="w-full h-full rounded-full" /> : item.username[0]}
              </div>
              <div>
                <span className="text-sm font-semibold">{item.username}</span>
                <span className="text-xs text-mute ml-2">{new Date(item.created_at).toLocaleDateString()}</span>
              </div>
              {item.is_pinned && <span className="ml-auto text-xs bg-primary/10 text-primary px-2 py-0.5 rounded-full">置顶</span>}
            </div>
            <p className="text-sm text-body mt-2">{item.content}</p>
            <button onClick={() => handleLike(item.id)} className={`mt-2 text-xs flex items-center gap-1 ${item.is_liked ? 'text-primary' : 'text-mute'}`}>
              👍 {item.like_count}
            </button>
          </div>
        ))}
        {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无留言</p>}
      </div>

      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
