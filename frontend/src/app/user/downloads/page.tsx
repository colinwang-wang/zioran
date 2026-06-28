'use client';

import { useState, useEffect } from 'react';
import { getUserDownloads, downloadCourse } from '@/lib/services';
import Pagination from '@/components/Pagination';
import type { PaginatedList, DownloadItem } from '@/types';

const formatDateTime = (value: string) => {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleString('zh-CN', { hour12: false });
};

export default function DownloadsPage() {
  const [data, setData] = useState<PaginatedList<DownloadItem>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });

  const fetchData = (page = 1) => { getUserDownloads({ page }).then(setData).catch(() => {}); };
  useEffect(() => { fetchData(); }, []);

  const handleDownload = async (courseId: number) => {
    try {
      const res = await downloadCourse(courseId);
      if (res.resources?.length > 0) {
        alert(`下载链接:\n${res.resources.map((r: { name: string; url: string; password?: string }) => `${r.name}: ${r.url}${r.password ? ` 密码:${r.password}` : ''}`).join('\n')}`);
      } else {
        alert('暂无可用资源');
      }
    } catch { alert('下载失败'); }
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">我的下载</h2>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-hairline text-left text-mute">
              <th className="py-3 font-semibold">课程名称</th>
              <th className="py-3 font-semibold">订单号</th>
              <th className="py-3 font-semibold">下载时间</th>
              <th className="py-3 font-semibold text-right">金币</th>
              <th className="py-3 font-semibold text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((d) => (
              <tr key={d.id} className="border-b border-hairline">
                <td className="py-3">
                  <div className="flex items-center gap-3">
                    {d.cover && <img src={d.cover} alt="" className="w-12 h-9 rounded object-cover shrink-0" />}
                    <span className="font-semibold text-ink line-clamp-1">{d.title}</span>
                  </div>
                </td>
                <td className="py-3 text-mute whitespace-nowrap">{d.order_no || '-'}</td>
                <td className="py-3 text-mute whitespace-nowrap">{formatDateTime(d.created_at)}</td>
                <td className="py-3 text-right text-mute whitespace-nowrap">{d.amount ?? 0} 金币</td>
                <td className="py-3 text-right">
                  <button onClick={() => handleDownload(d.course_id)} className="inline-block px-3 py-1 bg-primary text-white text-xs font-bold rounded-card">
                    下载
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无下载记录</p>}
      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
