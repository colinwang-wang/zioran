'use client';

import { useState } from 'react';
import { recharge } from '@/lib/services';
import type { RechargeResponse } from '@/types';

const amounts = [10, 50, 100, 200, 500, 1000];

export default function RechargePage() {
  const [amount, setAmount] = useState(100);
  const [method, setMethod] = useState('wechat');
  const [loading, setLoading] = useState(false);
  const [payment, setPayment] = useState<RechargeResponse | null>(null);

  const handleRecharge = async () => {
    setLoading(true);
    setPayment(null);
    try {
      const res = await recharge({ amount, pay_method: method });
      setPayment(res);
      if (res.pay_url.startsWith('mock://')) {
        alert('模拟充值已完成');
      } else if (method === 'alipay' && res.pay_url.startsWith('http')) {
        window.location.href = res.pay_url;
      }
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '充值失败';
      alert(msg);
    }
    setLoading(false);
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">在线充值</h2>
      <div className="grid grid-cols-3 gap-3 mb-6">
        {amounts.map((a) => (
          <button key={a} onClick={() => setAmount(a)} className={`py-3 rounded-card text-sm font-bold ${amount === a ? 'bg-primary text-white' : 'bg-surface text-ink'}`}>{a} 金币</button>
        ))}
      </div>
      <div className="flex gap-3 mb-6">
        <button onClick={() => setMethod('wechat')} className={`flex-1 py-3 rounded-card text-sm font-bold ${method === 'wechat' ? 'bg-[#07c160] text-white' : 'bg-surface'}`}>微信支付</button>
        <button onClick={() => setMethod('alipay')} className={`flex-1 py-3 rounded-card text-sm font-bold ${method === 'alipay' ? 'bg-[#1677ff] text-white' : 'bg-surface'}`}>支付宝</button>
      </div>
      <button onClick={handleRecharge} disabled={loading} className="w-full py-3 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">确认充值 {amount} 金币</button>
      {payment && !payment.pay_url.startsWith('mock://') && (
        <div className="mt-4 p-4 rounded-card bg-surface border border-hairline text-sm">
          <p className="font-semibold text-ink">订单已创建：{payment.order_no}</p>
          <p className="text-mute mt-1">请完成支付，到账以支付回调为准。</p>
          <a href={payment.pay_url} target="_blank" rel="noreferrer" className="block mt-3 break-all text-primary hover:underline">
            打开支付链接
          </a>
        </div>
      )}
    </div>
  );
}
