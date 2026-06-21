'use client';

import { useEffect, useState } from 'react';
import Pagination from '@/components/Pagination';
import { getCoinTransactions } from '@/lib/services';
import type { CoinTransaction, PaginatedList } from '@/types';

const typeLabels: Record<string, string> = {
  recharge: '充值',
  purchase: '购买资源',
  vip: '升级VIP',
};

export default function TransactionsPage() {
  const [data, setData] = useState<PaginatedList<CoinTransaction>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });

  const fetchData = (page = 1) => {
    getCoinTransactions({ page, pageSize: 10 }).then(setData).catch(() => {});
  };

  useEffect(() => { fetchData(); }, []);

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">充值记录</h2>
      <div className="space-y-3">
        {data.items.map((item) => (
          <div key={item.id} className="p-4 bg-surface rounded-card flex items-center justify-between gap-4">
            <div className="min-w-0">
              <p className="text-sm font-semibold text-ink">{item.description || typeLabels[item.type] || item.type}</p>
              <p className="text-xs text-mute mt-1">{typeLabels[item.type] || item.type} · {new Date(item.created_at).toLocaleString()}</p>
            </div>
            <div className="text-right shrink-0">
              <p className={`text-sm font-bold ${item.amount >= 0 ? 'text-primary' : 'text-ink'}`}>{item.amount >= 0 ? '+' : ''}{item.amount} 金币</p>
              <p className="text-xs text-mute mt-1">余额 {item.balance_after}</p>
            </div>
          </div>
        ))}
        {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无充值记录</p>}
      </div>
      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
