'use client';

import { useState, useEffect } from 'react';
import { getUserOrders } from '@/lib/services';
import Pagination from '@/components/Pagination';
import type { PaginatedList, OrderItem } from '@/types';

export default function OrdersPage() {
  const [data, setData] = useState<PaginatedList<OrderItem>>({ items: [], total: 0, page: 1, pageSize: 10, totalPages: 0 });

  const fetchData = (page = 1) => { getUserOrders({ page }).then(setData).catch(() => {}); };
  useEffect(() => { fetchData(); }, []);

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">我的订单</h2>
      <div className="space-y-3">
        {data.items.map((o) => (
          <div key={o.id} className="p-4 bg-surface rounded-card flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold">{o.target_name}</p>
              <p className="text-xs text-mute mt-1">{o.order_no} · {new Date(o.created_at).toLocaleDateString()}</p>
            </div>
            <div className="text-right">
              <p className="text-sm font-bold text-primary">{o.amount} 金币</p>
              <p className="text-xs text-mute">{o.status === 'paid' ? '已支付' : o.status === 'pending' ? '待支付' : o.status}</p>
            </div>
          </div>
        ))}
        {data.items.length === 0 && <p className="text-sm text-mute text-center py-8">暂无订单</p>}
      </div>
      <Pagination page={data.page} totalPages={data.totalPages} onChange={fetchData} />
    </div>
  );
}
