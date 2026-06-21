'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { login, getCaptcha, getWechatAuthURL } from '@/lib/services';

export default function LoginPage() {
  const { setAuth } = useAuth();
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [captcha, setCaptcha] = useState('');
  const [captchaKey, setCaptchaKey] = useState('');
  const [captchaImage, setCaptchaImage] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const loadCaptcha = async () => {
    try {
      const res = await getCaptcha();
      setCaptchaKey(res.captcha_key);
      setCaptchaImage(res.captcha_image);
    } catch { /* ignore */ }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) { setError('请填写邮箱和密码'); return; }
    if (!captcha || !captchaKey) { setError('请先获取并输入验证码'); return; }
    setLoading(true);
    setError('');
    try {
      const res = await login({ email, password, captcha, captcha_key: captchaKey });
      setAuth(res.token, res.user);
      router.push('/');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || '登录失败';
      setError(msg);
      loadCaptcha();
    }
    setLoading(false);
  };

  const handleWechatLogin = async () => {
    try {
      const authURL = await getWechatAuthURL();
      window.location.href = authURL;
    } catch {
      setError('微信登录暂不可用');
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-sm bg-canvas rounded-lg p-8 border border-hairline shadow-sm">
        <div className="text-center mb-6">
          <h1 className="text-xl font-bold text-primary">知猿</h1>
          <p className="text-sm text-mute mt-1">欢迎回来</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="邮箱" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="密码" className="w-full px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
          <div className="flex gap-2">
            <input type="text" value={captcha} onChange={(e) => setCaptcha(e.target.value)} placeholder="验证码" className="flex-1 px-4 py-3 rounded-card bg-surface border border-hairline text-sm focus:border-primary outline-none" />
            {captchaImage ? (
              <img src={captchaImage} alt="验证码" onClick={loadCaptcha} className="h-11 rounded-card cursor-pointer" />
            ) : (
              <button type="button" onClick={loadCaptcha} className="px-4 py-3 bg-surface rounded-card text-xs font-semibold text-primary whitespace-nowrap">显示验证码</button>
            )}
          </div>
          {error && <p className="text-xs text-red-600">{error}</p>}
          <button type="submit" disabled={loading} className="w-full py-3 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed disabled:opacity-50">
            登 录
          </button>
        </form>
        <div className="mt-4 text-center text-sm">
          <span className="text-mute">没有账号？</span>
          <Link href="/register" className="text-primary font-semibold ml-1">注册</Link>
          <span className="text-mute ml-3">·</span>
          <Link href="/forgot-password" className="text-primary font-semibold ml-1">忘记密码</Link>
        </div>
        <div className="mt-6 border-t border-hairline pt-4 text-center">
          <p className="text-xs text-mute mb-3">社交账号快速登录</p>
          <button type="button" onClick={handleWechatLogin} className="px-4 py-2 bg-[#07c160] text-white text-xs font-bold rounded-card">微信登录</button>
        </div>
      </div>
    </div>
  );
}
