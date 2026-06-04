'use client';

import { useState, useEffect } from 'react';
import { getTickets, createTicket, getTicketDetail, replyTicket } from '@/lib/services';
import type { Ticket } from '@/types';

export default function TicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [selected, setSelected] = useState<Ticket | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [replyContent, setReplyContent] = useState('');
  const [loading, setLoading] = useState(false);

  const loadTickets = async () => {
    try { const res = await getTickets(); setTickets(res.items); } catch { /* ignore */ }
  };

  useEffect(() => { loadTickets(); }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || !content) return;
    setLoading(true);
    try { await createTicket({ title, content }); setTitle(''); setContent(''); setShowForm(false); loadTickets(); } catch { /* ignore */ }
    setLoading(false);
  };

  const handleSelect = async (id: number) => {
    try { const t = await getTicketDetail(id); setSelected(t); } catch { /* ignore */ }
  };

  const handleReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!replyContent || !selected) return;
    setLoading(true);
    try { await replyTicket(selected.id, replyContent); setReplyContent(''); handleSelect(selected.id); } catch { /* ignore */ }
    setLoading(false);
  };

  if (selected) {
    return (
      <div>
        <button onClick={() => setSelected(null)} className="text-sm text-primary mb-4">← 返回列表</button>
        <div className="bg-canvas rounded-lg border border-hairline p-6">
          <h2 className="font-bold text-ink">{selected.title}</h2>
          <p className="text-xs text-mute mt-1">状态: {selected.status} · {selected.created_at}</p>
          <p className="mt-3 text-sm text-ink">{selected.content}</p>
        </div>
        {selected.replies && selected.replies.length > 0 && (
          <div className="mt-4 space-y-3">
            {selected.replies.map((r) => (
              <div key={r.id} className={`p-4 rounded-lg border text-sm ${r.is_admin ? 'bg-blue-50 border-blue-200' : 'bg-canvas border-hairline'}`}>
                <p className="font-semibold text-xs text-mute">{r.is_admin ? '客服' : r.username} · {r.created_at}</p>
                <p className="mt-1 text-ink">{r.content}</p>
              </div>
            ))}
          </div>
        )}
        <form onSubmit={handleReply} className="mt-4 flex gap-2">
          <input value={replyContent} onChange={(e) => setReplyContent(e.target.value)} placeholder="输入回复..." className="flex-1 px-4 py-2 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
          <button type="submit" disabled={loading} className="px-4 py-2 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">回复</button>
        </form>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="font-bold text-ink">我的工单</h2>
        <button onClick={() => setShowForm(!showForm)} className="px-4 py-2 bg-primary text-white text-sm font-bold rounded-card">
          {showForm ? '取消' : '提交工单'}
        </button>
      </div>
      {showForm && (
        <form onSubmit={handleCreate} className="mb-6 p-4 bg-canvas rounded-lg border border-hairline space-y-3">
          <input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="工单主题" className="w-full px-4 py-2 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
          <textarea value={content} onChange={(e) => setContent(e.target.value)} placeholder="详细描述您的问题..." rows={4} className="w-full px-4 py-2 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none resize-none" />
          <button type="submit" disabled={loading} className="px-4 py-2 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">提交</button>
        </form>
      )}
      {tickets.length === 0 ? (
        <p className="text-sm text-mute text-center py-8">暂无工单</p>
      ) : (
        <div className="space-y-2">
          {tickets.map((t) => (
            <div key={t.id} onClick={() => handleSelect(t.id)} className="p-4 bg-canvas rounded-lg border border-hairline cursor-pointer hover:border-primary">
              <div className="flex justify-between items-center">
                <h3 className="font-semibold text-sm text-ink">{t.title}</h3>
                <span className={`text-xs px-2 py-0.5 rounded ${t.status === 'closed' ? 'bg-gray-100 text-mute' : 'bg-green-100 text-green-700'}`}>{t.status === 'closed' ? '已关闭' : '进行中'}</span>
              </div>
              <p className="text-xs text-mute mt-1">{t.created_at}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
