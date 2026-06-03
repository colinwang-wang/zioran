'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getComments, createComment } from '@/lib/services';
import type { CommentItem } from '@/types';

export default function CommentSection({ courseId }: { courseId: number }) {
  const { isLoggedIn } = useAuth();
  const [comments, setComments] = useState<CommentItem[]>([]);
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);

  const fetchComments = async () => {
    try {
      const res = await getComments({ target_type: 'course', target_id: courseId });
      setComments(res.items || []);
    } catch { /* ignore */ }
  };

  useEffect(() => { fetchComments(); }, [courseId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    if (!content.trim()) return;
    setLoading(true);
    try {
      await createComment({ target_type: 'course', target_id: courseId, content: content.trim() });
      setContent('');
      fetchComments();
    } catch { /* ignore */ }
    setLoading(false);
  };

  return (
    <div className="mt-8">
      <h3 className="font-bold text-ink mb-4">评论区</h3>
      <form onSubmit={handleSubmit} className="mb-6">
        <textarea value={content} onChange={(e) => setContent(e.target.value)} placeholder={isLoggedIn ? '发表评论...' : '登录后发表评论'} className="w-full p-4 rounded-card bg-surface border border-hairline text-sm resize-none h-24 focus:border-primary outline-none" />
        <button type="submit" disabled={loading || !content.trim()} className="mt-2 px-6 py-2 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">
          发表
        </button>
      </form>
      <div className="space-y-4">
        {comments.map((c) => (
          <div key={c.id} className="flex gap-3">
            <div className="w-8 h-8 rounded-full bg-surface flex items-center justify-center text-xs font-bold shrink-0">
              {c.avatar ? <img src={c.avatar} alt="" className="w-full h-full rounded-full" /> : c.username[0]}
            </div>
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <span className="text-sm font-semibold">{c.username}</span>
                <span className="text-xs text-mute">{new Date(c.created_at).toLocaleDateString()}</span>
              </div>
              <p className="text-sm text-body mt-1">{c.content}</p>
              {c.children && c.children.length > 0 && (
                <div className="ml-4 mt-3 space-y-3 border-l-2 border-hairline pl-4">
                  {c.children.map((child) => (
                    <div key={child.id}>
                      <span className="text-sm font-semibold">{child.username}</span>
                      <span className="text-xs text-mute ml-2">{new Date(child.created_at).toLocaleDateString()}</span>
                      <p className="text-sm text-body mt-0.5">{child.content}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        ))}
        {comments.length === 0 && <p className="text-sm text-mute text-center py-4">暂无评论</p>}
      </div>
    </div>
  );
}
