'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { getBanners, getVipPackages, purchaseVip } from '@/lib/services';
import VipCard from '@/components/VipCard';
import type { Banner, VipPackage } from '@/types';

export default function VipPage() {
  const { isLoggedIn } = useAuth();
  const router = useRouter();
  const [packages, setPackages] = useState<VipPackage[]>([]);
  const [banners, setBanners] = useState<Banner[]>([]);
  const [loadingId, setLoadingId] = useState<number | null>(null);

  useEffect(() => {
    getVipPackages().then(setPackages).catch(() => {
      setPackages([{ id: 1, name: '终身VIP', price: 99, original_price: 699, duration: '永久', features: [] }]);
    });
    getBanners('vip').then(setBanners).catch(() => {});
  }, []);

  const handlePurchase = async (pkg: VipPackage) => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    setLoadingId(pkg.id);
    try {
      await purchaseVip(pkg.id);
      alert('购买成功！已成为VIP会员');
    } catch (err: unknown) {
      const data = (err as { response?: { data?: { message?: string } } })?.response?.data;
      const msg = data?.message || '购买失败';
      if (msg.includes('金币') || msg.includes('余额') || msg.includes('不足')) {
        router.push(`/user/recharge?returnTo=/vip&amount=${pkg.price}`);
      } else {
        alert(msg);
      }
    } finally {
      setLoadingId(null);
    }
  };

  return (
    <div>
      <section className="relative overflow-hidden text-white" style={{ background: banners[0]?.background_color || '#111827', minHeight: '240px' }}>
        {banners[0]?.image_url && <img src={banners[0].image_url} alt="" className="absolute inset-0 h-full w-full object-cover" />}
        <div className="absolute inset-0 bg-black/40" />
        <div className="relative mx-auto max-w-[1280px] px-6 py-16 text-center flex flex-col justify-center" style={{ minHeight: '240px' }}>
          <h1 className="text-4xl font-bold">成为知猿VIP会员</h1>
          <p className="mt-3 text-base text-white/80">全站课程免费下载，持续更新优质资源</p>
        </div>
      </section>
      <div className="max-w-[1280px] mx-auto px-6 py-10 pb-16 text-center">
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
    </div>
  );
}
