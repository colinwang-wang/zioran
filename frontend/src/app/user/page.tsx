'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { getCoinBalance, getVipStatus } from '@/lib/services';
import type { CoinBalance, VipStatus } from '@/types';

export default function UserProfilePage() {
  const { user } = useAuth();
  const [coins, setCoins] = useState<CoinBalance | null>(null);
  const [vip, setVip] = useState<VipStatus | null>(null);

  useEffect(() => {
    getCoinBalance().then(setCoins).catch(() => {});
    getVipStatus().then(setVip).catch(() => {});
  }, []);

  return (
    <div className="space-y-6">
      <section className="overflow-hidden rounded-[18px] bg-ink text-white">
        <div className="flex flex-col gap-5 p-6 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-4">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl font-bold text-white">
              {user?.avatar ? <img src={user.avatar} alt="" className="h-full w-full rounded-full object-cover" /> : user?.username?.[0]}
            </div>
            <div>
              <div className="flex flex-wrap items-center gap-2">
                <h2 className="text-xl font-bold">{user?.username}</h2>
                <span className={`rounded-full px-2.5 py-1 text-xs font-bold ${vip?.is_vip ? 'bg-[#ff0036] text-white' : 'bg-white/10 text-white/70'}`}>
                  {vip?.is_vip ? 'VIP会员' : '普通用户'}
                </span>
              </div>
              {user?.email && <p className="mt-1 text-sm text-white/60">{user.email}</p>}
              {vip?.is_vip && vip.expires_at && <p className="mt-1 text-xs text-white/50">到期时间：{vip.expires_at}</p>}
            </div>
          </div>
          <div className="grid grid-cols-2 gap-3 sm:min-w-[280px]">
            <div className="rounded-card bg-white/10 p-4">
              <p className="text-xs text-white/50">金币余额</p>
              <p className="mt-1 text-2xl font-bold">{coins?.balance ?? 0}</p>
            </div>
            <div className="rounded-card bg-white/10 p-4">
              <p className="text-xs text-white/50">累计消费</p>
              <p className="mt-1 text-2xl font-bold">{coins?.total_spent ?? 0}</p>
            </div>
          </div>
        </div>
      </section>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Link href="/user/recharge" className="rounded-card border border-hairline bg-canvas p-4 transition hover:border-primary">
          <p className="text-sm font-bold text-ink">金币充值</p>
          <p className="mt-1 text-xs text-mute">充值后自动轮询到账状态</p>
        </Link>
        <Link href="/user/orders" className="rounded-card border border-hairline bg-canvas p-4 transition hover:border-primary">
          <p className="text-sm font-bold text-ink">订单记录</p>
          <p className="mt-1 text-xs text-mute">查看购买和充值明细</p>
        </Link>
        <Link href="/user/tickets" className="rounded-card border border-hairline bg-canvas p-4 transition hover:border-primary">
          <p className="text-sm font-bold text-ink">售后工单</p>
          <p className="mt-1 text-xs text-mute">提交问题并上传截图</p>
        </Link>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div className="rounded-card bg-surface p-4">
          <p className="text-xs text-mute">VIP状态</p>
          <p className="mt-1 text-lg font-bold text-primary">{vip?.is_vip ? 'VIP会员' : '普通用户'}</p>
        </div>
        <div className="rounded-card bg-surface p-4">
          <p className="text-xs text-mute">累计获得</p>
          <p className="mt-1 text-lg font-bold text-ink">{coins?.total_earned ?? 0}</p>
        </div>
        <div className="rounded-card bg-surface p-4">
          <p className="text-xs text-mute">套餐类型</p>
          <p className="mt-1 text-lg font-bold text-ink">{vip?.package || '-'}</p>
        </div>
      </div>
    </div>
  );
}
