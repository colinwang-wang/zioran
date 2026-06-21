'use client';

import Link from 'next/link';

interface VipCardProps {
  name?: string;
  price?: number;
  originalPrice?: number;
  actionLabel?: string;
  onAction?: () => void;
  actionHref?: string;
  loading?: boolean;
}

const benefits = ['持续更新课堂资源', '终身会员永久有效'];

export default function VipCard({
  name = '终身VIP',
  price = 99,
  originalPrice = 699,
  actionLabel = '立即升级',
  onAction,
  actionHref,
  loading = false,
}: VipCardProps) {
  const actionClass = 'mt-8 block w-full py-3.5 bg-[#ff0036] text-white text-sm font-bold rounded-[16px] hover:bg-[#e6002f] text-center transition disabled:opacity-50';

  return (
    <div className="max-w-[400px] w-full mx-auto bg-white rounded-[32px] p-10 border-2 border-[#ff0036] relative shadow-[0_20px_60px_rgba(255,0,54,0.08)]">
      <div className="absolute top-0 left-0 right-0 h-1 bg-[#ff0036] rounded-t-[32px]" />
      <span className="inline-block bg-[#ff0036] text-white px-4 py-1 rounded-full text-xs font-bold mb-6">推荐</span>
      <h3 className="text-2xl font-bold text-[#000]">{name}</h3>
      <div className="text-[44px] font-bold text-[#ff0036] mt-2">
        {price} <span className="text-base font-normal text-[#62625b]">金币</span>
      </div>
      <div className="text-sm text-[#91918c] line-through mt-1">原价 {originalPrice} 金币</div>
      <div className="mt-7 space-y-3 border-t border-b border-[#dadad3] py-5">
        {benefits.map((benefit) => (
          <div key={benefit} className="flex items-center justify-center gap-2 text-sm font-semibold text-[#33332e]">
            <span className="flex h-5 w-5 items-center justify-center rounded-full bg-[#fff0f3] text-[#ff0036] text-xs">✓</span>
            <span>{benefit}</span>
          </div>
        ))}
      </div>
      {actionHref ? (
        <Link href={actionHref} className={actionClass}>{actionLabel}</Link>
      ) : (
        <button type="button" onClick={onAction} disabled={loading} className={actionClass}>{actionLabel}</button>
      )}
    </div>
  );
}
