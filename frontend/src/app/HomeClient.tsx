'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import CourseCard from '@/components/CourseCard';
import { getCourses } from '@/lib/services';
import type { NavItem, Banner, CourseListItem, CategoryBrief } from '@/types';

interface Props {
  navItems: NavItem[];
  banners: Banner[];
  latest: CourseListItem[];
  categories: CategoryBrief[];
}

export default function HomeClient({ navItems, banners, latest, categories }: Props) {
  const [activeTab, setActiveTab] = useState<number | null>(null);
  const [tabCourses, setTabCourses] = useState<CourseListItem[]>(latest);
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    if (activeTab === null) {
      setTabCourses(latest);
      return;
    }
    getCourses({ categoryId: activeTab, pageSize: 8 }).then((res) => {
      setTabCourses(res.items || []);
    }).catch(() => setTabCourses([]));
  }, [activeTab, latest]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/courses?q=${encodeURIComponent(searchQuery.trim())}`);
    }
  };

  return (
    <div>
      {/* 金刚区 */}
      {navItems.length > 0 && (
        <section className="bg-surface py-6">
          <div className="max-w-container mx-auto px-4">
            <div className="flex flex-wrap gap-4 justify-center">
              {navItems.map((item) => (
                <Link key={item.id} href={item.url} className="flex flex-col items-center gap-2 px-4 py-3 rounded-card hover:bg-canvas transition-colors min-w-[80px]">
                  {item.icon ? <img src={item.icon} alt="" className="w-10 h-10" /> : <div className="w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold text-sm">{item.title[0]}</div>}
                  <span className="text-xs font-semibold text-ink">{item.title}</span>
                </Link>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* Banner */}
      {banners.length > 0 && (
        <section className="max-w-container mx-auto px-4 mt-6">
          <div className="rounded-card overflow-hidden aspect-[3/1] relative">
            {banners[0].link_url ? (
              <Link href={banners[0].link_url}>
                <img src={banners[0].image_url} alt={banners[0].title} className="w-full h-full object-cover" />
              </Link>
            ) : (
              <img src={banners[0].image_url} alt={banners[0].title} className="w-full h-full object-cover" />
            )}
          </div>
        </section>
      )}

      {/* 搜索区 */}
      <section className="py-12 text-center">
        <h1 className="text-3xl md:text-4xl font-bold text-ink tracking-tight">知猿课堂，学有所长</h1>
        <p className="text-mute mt-2 text-sm">以知为基，以猿为伴，打造优质网课资源课堂</p>
        <form onSubmit={handleSearch} className="mt-6 max-w-lg mx-auto px-4 flex gap-2">
          <input type="text" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} placeholder="搜索一下" className="flex-1 px-5 py-3 rounded-full bg-surface text-sm border border-hairline focus:border-primary outline-none" />
          <button type="submit" className="px-6 py-3 bg-primary text-white rounded-full text-sm font-bold hover:bg-primary-pressed">搜索</button>
        </form>
      </section>

      {/* 最新发布 */}
      <section className="max-w-container mx-auto px-4 pb-12">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-ink">最新发布</h2>
          <Link href="/courses?sort=latest" className="text-sm text-primary font-semibold hover:underline">查看更多 →</Link>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {latest.slice(0, 8).map((c) => <CourseCard key={c.id} course={c} />)}
        </div>
      </section>

      {/* 知猿课堂 */}
      <section className="max-w-container mx-auto px-4 pb-12">
        <h2 className="text-xl font-bold text-ink mb-4">知猿课堂</h2>
        <div className="flex flex-wrap gap-2 mb-6">
          <button onClick={() => setActiveTab(null)} className={`px-4 py-2 rounded-full text-sm font-bold transition-colors ${activeTab === null ? 'bg-ink text-white' : 'bg-surface text-ink hover:bg-secondary-bg'}`}>
            全部课堂
          </button>
          {categories.map((cat) => (
            <button key={cat.id} onClick={() => setActiveTab(cat.id)} className={`px-4 py-2 rounded-full text-sm font-bold transition-colors ${activeTab === cat.id ? 'bg-ink text-white' : 'bg-surface text-ink hover:bg-secondary-bg'}`}>
              {cat.name}
            </button>
          ))}
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {tabCourses.slice(0, 8).map((c) => <CourseCard key={c.id} course={c} />)}
        </div>
        <div className="text-center mt-6">
          <Link href={activeTab ? `/courses?categoryId=${activeTab}` : '/courses'} className="text-sm text-primary font-semibold hover:underline">
            查看更多 →
          </Link>
        </div>
      </section>

      {/* VIP */}
      <section className="bg-surface py-12">
        <div className="max-w-container mx-auto px-4 text-center">
          <h2 className="text-xl font-bold text-ink mb-6">关于VIP</h2>
          <div className="max-w-sm mx-auto bg-canvas rounded-card p-8 border border-hairline">
            <div className="text-primary text-sm font-bold mb-2">终身VIP</div>
            <div className="text-4xl font-bold text-ink">99 <span className="text-base font-normal text-mute">金币</span></div>
            <div className="text-sm text-mute mt-1">永久</div>
            <ul className="mt-4 text-sm text-body space-y-2 text-left">
              <li>· 持续每天更新资源</li>
              <li>· 原价699金币</li>
              <li>· 全站免费下载</li>
            </ul>
            <Link href="/vip" className="mt-6 block w-full py-3 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed text-center">
              立即升级
            </Link>
            <p className="text-xs text-mute mt-3">75%的人选择该套餐</p>
          </div>
        </div>
      </section>
    </div>
  );
}
