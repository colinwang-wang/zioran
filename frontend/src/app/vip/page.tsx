'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { getVipPackages, purchaseVip } from '@/lib/services';
import VipCard from '@/components/VipCard';
import type { VipPackage } from '@/types';

export default function VipPage() {
  const { isLoggedIn } = useAuth();
  const [packages, setPackages] = useState<VipPackage[]>([]);
  const [loadingId, setLoadingId] = useState<number | null>(null);

  useEffect(() => {
    getVipPackages().then(setPackages).catch(() => {
      setPackages([{ id: 1, name: '终身VIP', price: 99, original_price: 699, duration: '永久', features: [] }]);
    });
  }, []);

  const handlePurchase = async (pkg: VipPackage) => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    setLoadingId(pkg.id);
    try {
      await purchaseVip(pkg.id);
      alert('购买成功！已成为VIP会员');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '购买失败';
      alert(msg);
    } finally {
      setLoadingId(null);
    }
  };

  return (
    <div className="max-w-[1280px] mx-auto px-6 py-16 text-center">
      <h1 className="text-3xl font-bold text-[#000] mb-2">成为会员</h1>
      <p className="text-[#91918c] mb-10">解锁全站资源，享受VIP专属权益</p>
      <div className="flex justify-center">
        {packages.map((pkg) => (
          <VipCard
            key={pkg.id}
            name={pkg.name}
            price={pkg.price}
            originalPrice={pkg.original_price}
            actionLabel="立即升级"
            onAction={() => handlePurchase(pkg)}
            loading={loadingId === pkg.id}
          />
        ))}
      </div>
    </div>
  );
}
