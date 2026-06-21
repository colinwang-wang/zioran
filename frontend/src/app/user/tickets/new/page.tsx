'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { createTicket } from '@/lib/services';

export default function NewTicketPage() {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;
    setLoading(true);
    try {
      await createTicket({ title: title.trim(), content: content.trim() });
      router.push('/user/tickets');
    } catch {
      alert('提交失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-bold">提交工单</h2>
        <Link href="/user/tickets" className="text-sm font-semibold text-primary hover:underline">我的工单</Link>
      </div>
      <form onSubmit={handleCreate} className="p-4 bg-canvas rounded-lg border border-hairline space-y-4">
        <div>
          <label className="block text-sm font-semibold text-ink mb-2">工单主题</label>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="例如：下载链接失效"
            className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none"
          />
        </div>
        <div>
          <label className="block text-sm font-semibold text-ink mb-2">问题说明</label>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="请写清课程名称、遇到的问题、错误链接或截图说明。"
            rows={6}
            className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none resize-none"
          />
        </div>
        <button type="submit" disabled={loading || !title.trim() || !content.trim()} className="px-5 py-3 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">
          提交工单
        </button>
      </form>
    </div>
  );
}
