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
      setPackages([{ id: 1, name: '终身VIP', price: 99, original_price: 699, duration: '永久', features: [] }]);
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
    <div className="max-w-[1280px] mx-auto px-6 py-16 text-center">
      <h1 className="text-3xl font-bold text-[#000] mb-2">成为会员</h1>
      <p className="text-[#91918c] mb-10">解锁全站资源，享受VIP专属权益</p>
      <div className="flex justify-center">
        {packages.map((pkg) => (
          <div key={pkg.id} className="max-w-[400px] w-full bg-white rounded-[32px] p-10 border-2 border-[#ff0036] relative shadow-[0_20px_60px_rgba(255,0,54,0.08)]">
            <div className="absolute top-0 left-0 right-0 h-1 bg-[#ff0036] rounded-t-[32px]" />
            <span className="inline-block bg-[#ff0036] text-white px-4 py-1 rounded-full text-xs font-bold mb-6">推荐</span>
            <h3 className="text-2xl font-bold text-[#000]">{pkg.name}</h3>
            <div className="text-[44px] font-bold text-[#ff0036] mt-2">{pkg.price} <span className="text-base font-normal text-[#62625b]">金币</span></div>
            <div className="text-sm text-[#91918c] line-through mt-1">原价 {pkg.original_price} 金币</div>
            <div className="text-base text-[#62625b] mt-4 pb-4 border-b border-[#dadad3]">{pkg.duration === '永久' ? '永久有效' : pkg.duration}</div>
            <p className="mt-6 text-sm text-[#62625b] text-center">持续更新课堂资源，终身会员永久有效</p>
            <button onClick={() => handlePurchase(pkg)} className="mt-8 block w-full py-3.5 bg-[#ff0036] text-white text-sm font-bold rounded-[16px] hover:bg-[#e6002f] transition">立即升级</button>
            <p className="text-xs text-[#91918c] mt-4">75% 的人选择该套餐</p>
          </div>
        ))}
      </div>
    </div>
  );
}
