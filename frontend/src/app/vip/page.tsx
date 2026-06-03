'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getVipPackages, purchaseVip } from '@/lib/services';
import type { VipPackage } from '@/types';

export default function VipPage() {
  const { isLoggedIn } = useAuth();
  const [packages, setPackages] = useState<VipPackage[]>([]);

  useEffect(() => {
    getVipPackages().then(setPackages).catch(() => {
      setPackages([{ id: 1, name: '终身VIP', price: 99, original_price: 699, duration: '永久', features: ['持续每天更新资源', '原价699金币', '全站免费下载'] }]);
    });
  }, []);

  const handlePurchase = async (pkg: VipPackage) => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      await purchaseVip(pkg.id);
      alert('购买成功！已成为VIP会员');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '购买失败';
      alert(msg);
    }
  };

  return (
    <div className="max-w-container mx-auto px-4 py-12">
      <div className="text-center mb-10">
        <h1 className="text-3xl font-bold text-ink">成为会员</h1>
        <p className="text-mute mt-2">解锁全站资源，享受VIP专属权益</p>
      </div>
      <div className="flex flex-wrap justify-center gap-6">
        {packages.map((pkg) => (
          <div key={pkg.id} className="w-full max-w-sm bg-canvas rounded-card border-2 border-primary p-8 text-center">
            <div className="text-primary text-sm font-bold mb-2">{pkg.name}</div>
            <div className="text-5xl font-bold text-ink">{pkg.price} <span className="text-base font-normal text-mute">金币</span></div>
            <div className="text-sm text-mute mt-1">{pkg.duration}</div>
            <ul className="mt-6 text-sm text-body space-y-3 text-left">
              {pkg.features.map((f, i) => <li key={i} className="flex items-start gap-2"><span className="text-primary">✓</span>{f}</li>)}
            </ul>
            <button onClick={() => handlePurchase(pkg)} className="mt-8 w-full py-3 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed">
              立即升级
            </button>
            <p className="text-xs text-mute mt-3">75%的人选择该套餐</p>
          </div>
        ))}
      </div>
    </div>
  );
}
