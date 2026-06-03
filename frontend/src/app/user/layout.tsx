'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

const tabs = [
  { href: '/user', label: '我的资料' },
  { href: '/user/orders', label: '我的订单' },
  { href: '/user/downloads', label: '我的下载' },
  { href: '/user/favorites', label: '我的收藏' },
  { href: '/user/recharge', label: '充值' },
  { href: '/user/settings', label: '账号设置' },
];

export default function UserLayout({ children }: { children: React.ReactNode }) {
  const { isLoggedIn } = useAuth();
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    if (!isLoggedIn) router.push('/login');
  }, [isLoggedIn, router]);

  if (!isLoggedIn) return null;

  return (
    <div className="max-w-container mx-auto px-4 py-8">
      <h1 className="text-xl font-bold text-ink mb-6">个人中心</h1>
      <div className="flex flex-col md:flex-row gap-6">
        <aside className="w-full md:w-48 shrink-0">
          <nav className="flex md:flex-col gap-2 overflow-x-auto">
            {tabs.map((tab) => (
              <Link key={tab.href} href={tab.href} className={`px-4 py-2 rounded-card text-sm font-semibold whitespace-nowrap ${pathname === tab.href ? 'bg-primary text-white' : 'bg-surface text-ink hover:bg-secondary-bg'}`}>
                {tab.label}
              </Link>
            ))}
          </nav>
        </aside>
        <div className="flex-1 min-w-0">{children}</div>
      </div>
    </div>
  );
}
