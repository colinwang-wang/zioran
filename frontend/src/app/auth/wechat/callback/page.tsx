'use client';

import { Suspense, useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { wechatLoginCallback } from '@/lib/services';

function CallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setAuth } = useAuth();
  const handledRef = useRef(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (handledRef.current) return;
    handledRef.current = true;

    const code = searchParams.get('code');
    const wechatError = searchParams.get('errmsg') || searchParams.get('error');
    if (!code) {
      setError(wechatError || '微信授权未完成');
      return;
    }

    wechatLoginCallback(code)
      .then((res) => {
        setAuth(res.token, res.user);
        router.replace('/');
      })
      .catch((err: unknown) => {
        const msg = (err as { response?: { data?: { message?: string; error?: string } } })?.response?.data;
        setError(msg?.message || msg?.error || '微信登录失败');
      });
  }, [router, searchParams, setAuth]);

  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4">
      <div className="w-full max-w-sm bg-canvas rounded-lg p-8 border border-hairline shadow-sm text-center">
        <h1 className="text-lg font-bold text-primary">微信登录</h1>
        {error ? (
          <>
            <p className="mt-3 text-sm text-red-600">{error}</p>
            <Link href="/login" className="mt-5 inline-flex px-4 py-2 rounded-card bg-primary text-white text-sm font-semibold">
              返回登录
            </Link>
          </>
        ) : (
          <p className="mt-3 text-sm text-mute">正在完成登录...</p>
        )}
      </div>
    </div>
  );
}

export default function WechatCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-[70vh] flex items-center justify-center px-4">
          <div className="w-full max-w-sm bg-canvas rounded-lg p-8 border border-hairline shadow-sm text-center">
            <h1 className="text-lg font-bold text-primary">微信登录</h1>
            <p className="mt-3 text-sm text-mute">正在完成登录...</p>
          </div>
        </div>
      }
    >
      <CallbackContent />
    </Suspense>
  );
}
