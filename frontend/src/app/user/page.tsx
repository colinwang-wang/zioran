'use client';

import { useState, useEffect } from 'react';
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
    <div>
      <div className="flex items-center gap-4 mb-6">
        <div className="w-16 h-16 rounded-full bg-surface flex items-center justify-center text-2xl font-bold text-primary">
          {user?.avatar ? <img src={user.avatar} alt="" className="w-full h-full rounded-full" /> : user?.username?.[0]}
        </div>
        <div>
          <h2 className="text-lg font-bold">{user?.username}</h2>
          <p className="text-sm text-mute">{user?.phone}</p>
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="p-4 bg-surface rounded-card">
          <p className="text-xs text-mute">VIP状态</p>
          <p className="text-lg font-bold text-primary mt-1">{vip?.is_vip ? 'VIP会员' : '普通用户'}</p>
        </div>
        <div className="p-4 bg-surface rounded-card">
          <p className="text-xs text-mute">金币余额</p>
          <p className="text-lg font-bold text-ink mt-1">{coins?.balance ?? 0}</p>
        </div>
        <div className="p-4 bg-surface rounded-card">
          <p className="text-xs text-mute">累计消费</p>
          <p className="text-lg font-bold text-ink mt-1">{coins?.total_spent ?? 0}</p>
        </div>
      </div>
    </div>
  );
}
