'use client';

import { useState } from 'react';
import { recharge } from '@/lib/services';

const amounts = [10, 50, 100, 200, 500, 1000];

export default function RechargePage() {
  const [amount, setAmount] = useState(100);
  const [method, setMethod] = useState('wechat');
  const [loading, setLoading] = useState(false);

  const handleRecharge = async () => {
    setLoading(true);
    try {
      await recharge({ amount, pay_method: method });
      alert('充值成功');
    } catch { alert('充值失败'); }
    setLoading(false);
  };

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">充值</h2>
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
    </div>
  );
}
