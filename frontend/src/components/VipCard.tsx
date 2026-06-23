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

const benefits = ['全站课程免费下载', '持续每天更新资源', '专属客服支持'];

export default function VipCard({
  name = '终身VIP',
  price = 99,
  originalPrice = 699,
  actionLabel = '立即升级',
  onAction,
  actionHref,
  loading = false,
}: VipCardProps) {
  const actionClass = 'block w-full rounded-[16px] bg-[#ff0036] p-3.5 text-center text-base font-bold text-white transition hover:bg-[#e6002f] disabled:opacity-50';

  return (
    <div className="relative mx-auto w-full max-w-[400px] overflow-hidden rounded-[32px] border-2 border-[#ff0036] bg-white px-8 py-12">
      <div className="absolute left-0 right-0 top-0 h-1 bg-[#ff0036]" />
      <span className="inline-block bg-[#ff0036] text-white px-4 py-1 rounded-full text-xs font-bold mb-6">推荐</span>
      <h3 className="mb-2 text-[28px] font-bold text-[#000]">{name}</h3>
      <div className="mb-1 text-[44px] font-bold leading-tight text-[#ff0036]">
        {price} <span className="text-base font-normal text-[#62625b]">金币</span>
      </div>
      <div className="mb-6 text-sm text-[#91918c] line-through">原价 {originalPrice} 金币</div>
      <div className="mb-6 text-base text-[#62625b]">永久有效</div>
      <div className="mb-8 text-left">
        {benefits.map((benefit) => (
          <div key={benefit} className="border-b border-[#dadad3] py-2 text-sm text-[#33332e]">
            <span className="mr-2 font-bold text-[#ff0036]">✓</span>
            <span>{benefit}</span>
          </div>
        ))}
      </div>
      {actionHref ? (
        <Link href={actionHref} className={actionClass}>{actionLabel}</Link>
      ) : (
        <button type="button" onClick={onAction} disabled={loading} className={actionClass}>{actionLabel}</button>
      )}
      <div className="mt-3 text-xs text-[#62625b]">75% 的人选择该套餐</div>
    </div>
  );
}
