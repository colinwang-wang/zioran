'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { QRCodeSVG } from 'qrcode.react';
import { getOrder, getRechargeConfig, recharge } from '@/lib/services';
import type { RechargeResponse } from '@/types';

const fallbackAmounts = [10, 50, 100, 200, 500, 1000];

export default function RechargePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const returnTo = searchParams.get('returnTo');
  const requestedAmount = Number(searchParams.get('amount') || 0);
  const [amount, setAmount] = useState(requestedAmount > 0 ? requestedAmount : 100);
  const [amounts, setAmounts] = useState(fallbackAmounts);
  const [ratio, setRatio] = useState(1);
  const [customAmount, setCustomAmount] = useState(requestedAmount > 0 ? String(requestedAmount) : '');
  const [isCustom, setIsCustom] = useState(false);
  const [method, setMethod] = useState('wechat');
  const [loading, setLoading] = useState(false);
  const [payment, setPayment] = useState<RechargeResponse | null>(null);
  const [payStatus, setPayStatus] = useState<'idle' | 'pending' | 'paid'>('idle');
  const coins = amount * ratio;

  useEffect(() => {
    getRechargeConfig()
      .then((config) => {
        const nextAmounts = config.amounts?.filter((item) => item > 0) || [];
        setAmounts(nextAmounts.length > 0 ? nextAmounts : fallbackAmounts);
        setRatio(config.ratio > 0 ? config.ratio : 1);
        if (requestedAmount > 0) {
          setAmount(requestedAmount);
          setIsCustom(!nextAmounts.includes(requestedAmount));
          setCustomAmount(String(requestedAmount));
        } else if (nextAmounts.length > 0) {
          setAmount(nextAmounts[0]);
        }
      })
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (!payment || payStatus !== 'pending') return;
    const timer = window.setInterval(async () => {
      try {
        const order = await getOrder(payment.order_id);
        if (order.status === 'paid') {
          setPayStatus('paid');
          window.clearInterval(timer);
          if (returnTo) {
            setTimeout(() => router.push(returnTo), 1200);
          }
        }
      } catch { /* keep polling */ }
    }, 3000);
    return () => window.clearInterval(timer);
  }, [payment, payStatus, returnTo, router]);

  const handleRecharge = async () => {
    setLoading(true);
    setPayment(null);
    setPayStatus('idle');
    try {
      const res = await recharge({ amount, pay_method: method });
      setPayment(res);
      if (res.pay_url.startsWith('mock://')) {
        setPayStatus('paid');
      } else if (method === 'alipay' && res.pay_url.startsWith('http')) {
        window.localStorage.setItem('pending_recharge_order_id', String(res.order_id));
        window.location.href = res.pay_url;
      } else {
        setPayStatus('pending');
      }
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '充值失败';
      alert(msg);
    }
    setLoading(false);
  };

  const isWechatQR = payment && method === 'wechat' && payment.pay_url && !payment.pay_url.startsWith('mock://');

  return (
    <div>
      <h2 className="text-lg font-bold mb-4">在线充值</h2>
      <div className="mb-4 rounded-card bg-surface px-4 py-3 text-sm text-mute">
        当前充值比例：<span className="font-bold text-primary">1 元 = {ratio} 金币</span>
      </div>
      <div className="grid grid-cols-3 gap-3 mb-6">
        {amounts.map((a) => (
          <button key={a} onClick={() => { setAmount(a); setIsCustom(false); }} className={`py-3 rounded-card text-sm font-bold ${!isCustom && amount === a ? 'bg-primary text-white' : 'bg-surface text-ink'}`}>
            <span className="block">¥{a}</span>
            <span className={`mt-1 block text-xs font-medium ${!isCustom && amount === a ? 'text-white/80' : 'text-mute'}`}>{a * ratio} 金币</span>
          </button>
        ))}
        <button onClick={() => setIsCustom(true)} className={`py-3 rounded-card text-sm font-bold ${isCustom ? 'bg-primary text-white' : 'bg-surface text-ink'}`}>自定义</button>
      </div>
      {isCustom && (
        <div className="mb-6">
          <input type="number" min="1" placeholder="输入充值金额（元）" value={customAmount} onChange={(e) => { setCustomAmount(e.target.value); setAmount(Number(e.target.value) || 0); }} className="w-full px-4 py-3 rounded-card border border-hairline bg-surface text-sm" />
          <p className="mt-2 text-xs text-mute">预计到账 {coins || 0} 金币</p>
        </div>
      )}
      <div className="flex gap-3 mb-6">
        <button onClick={() => setMethod('wechat')} className={`flex-1 py-3 rounded-card text-sm font-bold ${method === 'wechat' ? 'bg-[#07c160] text-white' : 'bg-surface'}`}>微信支付</button>
        <button onClick={() => setMethod('alipay')} className={`flex-1 py-3 rounded-card text-sm font-bold ${method === 'alipay' ? 'bg-[#1677ff] text-white' : 'bg-surface'}`}>支付宝</button>
      </div>
      <button onClick={handleRecharge} disabled={loading || amount < 1} className="w-full py-3 bg-primary text-white text-sm font-bold rounded-card disabled:opacity-50">确认支付 ¥{amount || 0}，到账 {coins || 0} 金币</button>
      {isWechatQR && (
        <div className="mt-4 p-6 rounded-card bg-surface border border-hairline flex flex-col items-center">
          <p className="font-semibold text-ink text-sm mb-1">订单：{payment.order_no}</p>
          <p className="text-mute text-xs mb-4">请使用微信扫码支付 ¥{payment.amount}，到账 {payment.coins} 金币</p>
          <QRCodeSVG value={payment.pay_url} size={200} />
          <p className="text-mute text-xs mt-4">支付完成后页面将自动更新</p>
        </div>
      )}
      {payment && payStatus === 'paid' && (
        <div className="mt-4 rounded-card border border-[#b7eb8f] bg-[#f6ffed] p-4 text-sm text-[#237804]">
          支付成功，已到账 {payment.coins} 金币。
        </div>
      )}
      {payment && !isWechatQR && !payment.pay_url.startsWith('mock://') && (
        <div className="mt-4 p-4 rounded-card bg-surface border border-hairline text-sm">
          <p className="font-semibold text-ink">订单已创建：{payment.order_no}</p>
          <p className="text-mute mt-1">请完成 ¥{payment.amount} 支付，到账 {payment.coins} 金币。</p>
        </div>
      )}
    </div>
  );
}
