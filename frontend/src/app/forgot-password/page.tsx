'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { sendSMS, getCaptcha, forgotPassword } from '@/lib/services';

export default function ForgotPasswordPage() {
  const router = useRouter();
  const [phone, setPhone] = useState('');
  const [smsCode, setSmsCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [captcha, setCaptcha] = useState('');
  const [captchaKey, setCaptchaKey] = useState('');
  const [captchaImage, setCaptchaImage] = useState('');
  const [countdown, setCountdown] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const loadCaptcha = async () => {
    try {
      const res = await getCaptcha();
      setCaptchaKey(res.captcha_key);
      setCaptchaImage(res.captcha_image);
    } catch { /* ignore */ }
  };

  const handleSendSMS = async () => {
    if (!phone || !captcha || !captchaKey) { setError('请输入手机号和图形验证码'); return; }
    try {
      await sendSMS({ phone, captcha, captcha_key: captchaKey });
      setCountdown(60);
      const timer = setInterval(() => {
        setCountdown((c) => { if (c <= 1) { clearInterval(timer); return 0; } return c - 1; });
      }, 1000);
    } catch { setError('发送验证码失败'); loadCaptcha(); }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!phone || !smsCode || !newPassword) { setError('请填写所有字段'); return; }
    setLoading(true); setError('');
    try {
      await forgotPassword({ phone, sms_code: smsCode, new_password: newPassword });
      setSuccess(true);
      setTimeout(() => router.push('/login'), 1500);
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '重置密码失败';
      setError(msg);
    }
    setLoading(false);
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-sm bg-canvas rounded-lg p-8 border border-hairline shadow-sm">
        <div className="text-center mb-6">
          <h1 className="text-xl font-bold text-primary">重置密码</h1>
          <p className="text-sm text-mute mt-1">通过手机短信验证重置密码</p>
        </div>
        {success ? (
          <p className="text-center text-sm text-green-600">密码重置成功，正在跳转登录页...</p>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <input type="tel" value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="手机号" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
            <div className="flex gap-2">
              <input type="text" value={captcha} onChange={(e) => setCaptcha(e.target.value)} placeholder="图形验证码" className="flex-1 px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
              {captchaImage ? (
                <img src={captchaImage} alt="验证码" onClick={loadCaptcha} className="h-11 rounded-card cursor-pointer" />
              ) : (
                <button type="button" onClick={loadCaptcha} className="px-4 py-3 bg-surface rounded-card text-xs font-semibold text-primary whitespace-nowrap">获取图形码</button>
              )}
            </div>
            <div className="flex gap-2">
              <input type="text" value={smsCode} onChange={(e) => setSmsCode(e.target.value)} placeholder="短信验证码" className="flex-1 px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
              <button type="button" onClick={handleSendSMS} disabled={countdown > 0} className="px-4 py-3 bg-surface rounded-card text-xs font-semibold text-primary whitespace-nowrap disabled:opacity-50">
                {countdown > 0 ? `${countdown}s` : '发送验证码'}
              </button>
            </div>
            <input type="password" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} placeholder="新密码" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
            {error && <p className="text-xs text-red-600">{error}</p>}
            <button type="submit" disabled={loading} className="w-full py-3 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed disabled:opacity-50">
              重置密码
            </button>
          </form>
        )}
        <div className="mt-4 text-center text-sm">
          <Link href="/login" className="text-primary font-semibold">返回登录</Link>
        </div>
      </div>
    </div>
  );
}
