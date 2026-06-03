'use client';

import Link from 'next/link';
import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

export default function Header() {
  const { user, isLoggedIn, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const router = useRouter();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/courses?q=${encodeURIComponent(searchQuery.trim())}`);
      setSearchOpen(false);
      setSearchQuery('');
    }
  };

  return (
    <header className="sticky top-0 z-50 bg-canvas border-b border-hairline">
      <div className="max-w-container mx-auto px-4 h-16 flex items-center justify-between">
        {/* Logo */}
        <div className="flex items-center gap-8">
          <Link href="/" className="text-xl font-bold text-primary">知猿</Link>
          <nav className="hidden md:flex items-center gap-6">
            <Link href="/" className="text-sm font-semibold text-ink hover:text-primary">首页</Link>
            <Link href="/courses" className="text-sm font-semibold text-ink hover:text-primary">知猿课堂</Link>
            <Link href="/guestbook" className="text-sm font-semibold text-ink hover:text-primary">留言反馈</Link>
            <Link href="/vip" className="text-sm font-semibold text-ink hover:text-primary">成为会员</Link>
          </nav>
        </div>

        {/* Right */}
        <div className="flex items-center gap-3">
          <button onClick={() => setSearchOpen(!searchOpen)} className="p-2 rounded-full hover:bg-surface" aria-label="搜索">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          </button>
          {isLoggedIn ? (
            <div className="relative">
              <button onClick={() => setUserMenuOpen(!userMenuOpen)} className="flex items-center gap-2 px-3 py-1.5 rounded-card bg-surface text-sm font-semibold">
                {user?.avatar ? <img src={user.avatar} alt="" className="w-6 h-6 rounded-full" /> : <span className="w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center text-xs">{user?.username?.[0]}</span>}
                {user?.username}
              </button>
              {userMenuOpen && (
                <div className="absolute right-0 top-full mt-2 w-40 bg-canvas rounded-card shadow-lg border border-hairline py-2 z-50">
                  <Link href="/user" className="block px-4 py-2 text-sm hover:bg-surface" onClick={() => setUserMenuOpen(false)}>个人中心</Link>
                  <button onClick={() => { logout(); setUserMenuOpen(false); }} className="w-full text-left px-4 py-2 text-sm hover:bg-surface text-primary">退出登录</button>
                </div>
              )}
            </div>
          ) : (
            <Link href="/login" className="px-4 py-2 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed">登录</Link>
          )}
          {/* Mobile hamburger */}
          <button onClick={() => setMenuOpen(!menuOpen)} className="md:hidden p-2" aria-label="菜单">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" /></svg>
          </button>
        </div>
      </div>

      {/* Search overlay */}
      {searchOpen && (
        <div className="absolute inset-x-0 top-full bg-canvas border-b border-hairline p-4 shadow-lg">
          <form onSubmit={handleSearch} className="max-w-container mx-auto flex gap-2">
            <input type="text" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} placeholder="搜索课程..." className="flex-1 px-4 py-3 rounded-full bg-surface text-sm" autoFocus />
            <button type="submit" className="px-6 py-3 bg-primary text-white rounded-full text-sm font-bold">搜索</button>
          </form>
        </div>
      )}

      {/* Mobile menu */}
      {menuOpen && (
        <div className="md:hidden border-t border-hairline bg-canvas">
          <nav className="flex flex-col p-4 gap-3">
            <Link href="/" className="py-2 text-sm font-semibold" onClick={() => setMenuOpen(false)}>首页</Link>
            <Link href="/courses" className="py-2 text-sm font-semibold" onClick={() => setMenuOpen(false)}>知猿课堂</Link>
            <Link href="/guestbook" className="py-2 text-sm font-semibold" onClick={() => setMenuOpen(false)}>留言反馈</Link>
            <Link href="/vip" className="py-2 text-sm font-semibold" onClick={() => setMenuOpen(false)}>成为会员</Link>
          </nav>
        </div>
      )}
    </header>
  );
}
