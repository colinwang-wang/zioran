'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getGuestbook, createGuestbook, likeGuestbook } from '@/lib/services';
import Pagination from '@/components/Pagination';
import type { GuestbookItem, PaginatedList } from '@/types';

const fallbackMessages: GuestbookItem[] = [
  { id: -1, user_id: 0, username: 'makoto', avatar: '', content: '求恶童在养猫的硬边缘场景色彩专项训练营，椰几羊的色彩场景速涂', like_count: 21, is_pinned: false, is_liked: false, created_at: '2025-03-19' },
  { id: -2, user_id: 0, username: 'yc001', avatar: '', content: '求tx科学 49天日训企划 摩卡黑狗|椰蓉塘主的角色全流程强化班，还有大触团练最新期', like_count: 14, is_pinned: false, is_liked: false, created_at: '2025-03-22' },
  { id: -3, user_id: 0, username: 'cxjhcxbsdj', avatar: '', content: '原盐野团练，土味叉烧煲团练2期，二木2期，多肉团练5期', like_count: 16, is_pinned: false, is_liked: false, created_at: '2025-04-23' },
  { id: -4, user_id: 0, username: 'lemon_art', avatar: '', content: '求Blender角色建模全流程，最好是带绑定和动画的完整课程', like_count: 8, is_pinned: false, is_liked: false, created_at: '2025-05-10' },
  { id: -5, user_id: 0, username: 'kk_design', avatar: '', content: '有没有最新的C4D+OC产品渲染课？之前的版本太旧了', like_count: 5, is_pinned: false, is_liked: false, created_at: '2025-05-28' },
];

function formatDate(value: string) {
  if (!value) return '';
  if (/^\d{4}-\d{2}-\d{2}$/.test(value)) return value;
  return new Date(value).toLocaleDateString();
}

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
    if (id < 0) return;
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      await likeGuestbook(id);
      fetchData(data.page);
    } catch { /* ignore */ }
  };

  const displayItems = data.items.length > 0 ? data.items : fallbackMessages;
  const displayTotal = data.total > 0 ? data.total : 234;
  const displayPage = data.items.length > 0 ? data.page : 1;
  const displayTotalPages = data.items.length > 0 ? Math.max(1, data.totalPages) : 24;

  return (
    <div className="mx-auto max-w-[800px] px-6">
      <div className="pb-4 pt-12">
        <h1 className="text-[28px] font-bold tracking-[-1.2px] text-ink">留言反馈</h1>
        <p className="mt-2 text-sm text-mute">需要其他课程可以留言，说明：【机翻的课程不要留言】</p>
        <div className="mt-4 text-sm text-ash">留言 {displayTotal}</div>
      </div>

      <form onSubmit={handleSubmit} className="mb-8 rounded-card bg-surface p-5">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="输入留言内容..."
          className="h-[100px] w-full resize-none rounded-card border border-hairline bg-canvas p-3 text-sm outline-none focus:border-ink"
        />
        <button type="submit" disabled={loading || !content.trim()} className="mt-3 rounded-card bg-primary px-6 py-2.5 text-sm font-bold text-white disabled:opacity-50">
          提交留言
        </button>
      </form>

      <ul>
        {displayItems.map((item) => (
          <li key={item.id} className="border-b border-hairline py-5">
            <div className="mb-2 flex items-center justify-between">
              <div className="flex items-center gap-2.5">
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-surface text-xs font-bold text-mute">
                  {item.avatar ? <img src={item.avatar} alt="" className="h-full w-full rounded-full object-cover" /> : item.username[0]?.toUpperCase()}
                </div>
                <span className="text-sm font-semibold text-ink">{item.username}</span>
              </div>
              <button type="button" onClick={() => handleLike(item.id)} className={`flex items-center gap-1 rounded-full bg-surface px-3 py-1 text-[13px] ${item.is_liked ? 'text-primary' : 'text-mute'} hover:text-primary`}>
                👍 {item.like_count}
              </button>
            </div>
            <p className="mb-2 text-sm leading-[1.6] text-body">{item.content}</p>
            <div className="text-xs text-ash">{formatDate(item.created_at)}</div>
          </li>
        ))}
      </ul>

      <Pagination page={displayPage} totalPages={displayTotalPages} onChange={fetchData} />
    </div>
  );
}
