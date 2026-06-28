'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { getTickets, getTicketDetail, replyTicket } from '@/lib/services';
import type { Ticket } from '@/types';

const statusLabel: Record<string, string> = {
  open: '待处理',
  processing: '处理中',
  replied: '已回复',
  closed: '已关闭',
};

const formatDateTime = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleString('zh-CN', { hour12: false });
};

export default function TicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [selected, setSelected] = useState<Ticket | null>(null);
  const [replyContent, setReplyContent] = useState('');
  const [loading, setLoading] = useState(false);

  const loadTickets = async () => {
    try { const res = await getTickets(); setTickets(res.items); } catch { /* ignore */ }
  };

  useEffect(() => { loadTickets(); }, []);

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
          <p className="text-xs text-mute mt-1">状态: {statusLabel[selected.status] || selected.status} · {formatDateTime(selected.created_at)}</p>
          <p className="mt-3 text-sm text-ink">{selected.content}</p>
          {selected.attachments && selected.attachments.length > 0 && (
            <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-3">
              {selected.attachments.map((url, index) => (
                <a key={`${url}-${index}`} href={url} target="_blank" rel="noreferrer" className="block overflow-hidden rounded-card border border-hairline bg-surface">
                  <img src={url} alt="" className="h-28 w-full object-cover" />
                </a>
              ))}
            </div>
          )}
        </div>
        {selected.replies && selected.replies.length > 0 && (
          <div className="mt-4 space-y-3">
            {selected.replies.map((r) => (
              <div key={r.id} className={`p-4 rounded-lg border text-sm ${r.is_admin ? 'bg-blue-50 border-blue-200' : 'bg-canvas border-hairline'}`}>
                <p className="font-semibold text-xs text-mute">{r.is_admin ? '客服' : r.username} · {formatDateTime(r.created_at)}</p>
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
        <Link href="/user/tickets/new" className="px-4 py-2 bg-primary text-white text-sm font-bold rounded-card">
          提交工单
        </Link>
      </div>
      {tickets.length === 0 ? (
        <p className="text-sm text-mute text-center py-8">暂无工单</p>
      ) : (
        <div className="space-y-2">
          {tickets.map((t) => (
            <div key={t.id} className="p-4 bg-canvas rounded-lg border border-hairline hover:border-primary">
              <div className="flex items-center justify-between gap-3">
                <h3 className="font-semibold text-sm text-ink">{t.title}</h3>
                <span className={`text-xs px-2 py-0.5 rounded ${t.status === 'closed' ? 'bg-gray-100 text-mute' : 'bg-green-100 text-green-700'}`}>{t.status === 'closed' ? '已关闭' : '进行中'}</span>
              </div>
              <div className="mt-3 flex items-center justify-between gap-3">
                <p className="text-xs text-mute">{formatDateTime(t.created_at)}</p>
                <button onClick={() => handleSelect(t.id)} className="rounded-card bg-[#ff0036] px-3 py-1.5 text-xs font-bold text-white">查看详情</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
