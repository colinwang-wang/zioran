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
    if (activeTab === null) { setTabCourses(latest); return; }
    getCourses({ categoryId: activeTab, pageSize: 8 }).then((res) => {
      setTabCourses(Array.isArray(res) ? res : res?.items || []);
    }).catch(() => setTabCourses([]));
  }, [activeTab, latest]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) router.push(`/courses?q=${encodeURIComponent(searchQuery.trim())}`);
  };

  return (
    <div>
      {/* 金刚区 */}
      {navItems.length > 0 && (
        <section className="bg-[#f6f6f3] py-5">
          <div className="max-w-[1280px] mx-auto px-6">
            <div className="flex flex-wrap gap-4 justify-center">
              {navItems.map((item) => (
                <Link key={item.id} href={item.url || '#'} className="flex flex-col items-center gap-2 px-5 py-3 rounded-[16px] bg-white hover:shadow-sm transition min-w-[100px]">
                  {item.icon ? <img src={item.icon} alt="" className="w-12 h-12 rounded-xl object-cover" /> : <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#ff0036] to-[#ff6b6b] flex items-center justify-center text-white font-bold text-sm">{item.title[0]}</div>}
                  <span className="text-xs font-semibold text-[#000]">{item.title}</span>
                </Link>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* Banner + 搜索（合并为一个深色区域） */}
      <section className="py-8">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="relative rounded-[32px] overflow-hidden py-16 px-8 text-center text-white">
            {/* 背景：优先显示Banner图片，否则默认渐变 */}
            {banners.length > 0 && banners[0].image_url ? (
              <div className="absolute inset-0">
                <img src={banners[0].image_url} alt="" className="w-full h-full object-cover" />
                <div className="absolute inset-0 bg-black/50" />
              </div>
            ) : (
              <div className="absolute inset-0 bg-gradient-to-br from-[#1a1a2e] via-[#16213e] to-[#0f3460]" />
            )}
            {/* 装饰光效 */}
            <div className="absolute top-[-50%] right-[-10%] w-[400px] h-[400px] bg-[radial-gradient(circle,rgba(255,0,54,0.15),transparent_70%)]" />
            <h1 className="text-3xl md:text-[44px] font-bold tracking-tight relative">知猿课堂，学有所长</h1>
            <p className="text-white/60 mt-3 text-sm relative">以知为基，以猿为伴，打造优质网课资源课堂</p>
            <form onSubmit={handleSearch} className="mt-8 max-w-[560px] mx-auto flex gap-2 relative">
              <input type="text" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} placeholder="输入关键字搜索课程..." className="flex-1 px-5 py-3.5 rounded-full bg-white/95 text-[#333] text-sm border-none outline-none placeholder:text-[#999]" />
              <button type="submit" className="px-7 py-3.5 bg-[#ff0036] text-white rounded-full text-sm font-bold hover:bg-[#e6002f] transition">搜索</button>
            </form>
          </div>
        </div>
      </section>

      {/* 最新发布 */}
      <section className="max-w-[1280px] mx-auto px-6 pb-12">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-xl font-bold text-[#000]">最新发布</h2>
          <Link href="/courses?sort=latest" className="text-sm text-[#ff0036] font-semibold hover:underline">查看更多 →</Link>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {latest.slice(0, 8).map((c) => <CourseCard key={c.id} course={c} />)}
        </div>
      </section>

      {/* 知猿课堂（Tab切换） */}
      <section className="bg-[#f6f6f3] py-12">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="flex items-center justify-between mb-5">
            <h2 className="text-xl font-bold text-[#000]">知猿课堂</h2>
            <Link href={activeTab ? `/courses?categoryId=${activeTab}` : '/courses'} className="text-sm text-[#ff0036] font-semibold hover:underline">查看更多 →</Link>
          </div>
          <div className="flex flex-wrap gap-2 mb-6">
            <button onClick={() => setActiveTab(null)} className={`px-4 py-2 rounded-full text-sm font-bold transition ${activeTab === null ? 'bg-[#000] text-white' : 'bg-white text-[#000] hover:bg-[#e5e5e0]'}`}>全部课堂</button>
            {categories.map((cat) => (
              <button key={cat.id} onClick={() => setActiveTab(cat.id)} className={`px-4 py-2 rounded-full text-sm font-bold transition ${activeTab === cat.id ? 'bg-[#000] text-white' : 'bg-white text-[#000] hover:bg-[#e5e5e0]'}`}>{cat.name}</button>
            ))}
          </div>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {tabCourses.slice(0, 8).map((c) => <CourseCard key={c.id} course={c} />)}
          </div>
        </div>
      </section>

      {/* VIP */}
      <section className="py-16">
        <div className="max-w-[1280px] mx-auto px-6 text-center">
          <h2 className="text-xl font-bold text-[#000] mb-8">关于VIP</h2>
          <div className="max-w-[400px] mx-auto bg-white rounded-[32px] p-10 border-2 border-[#ff0036] relative shadow-[0_20px_60px_rgba(255,0,54,0.08)]">
            <div className="absolute top-0 left-0 right-0 h-1 bg-[#ff0036] rounded-t-[32px]" />
            <span className="inline-block bg-[#ff0036] text-white px-4 py-1 rounded-full text-xs font-bold mb-6">推荐</span>
            <h3 className="text-2xl font-bold text-[#000]">终身VIP</h3>
            <div className="text-[44px] font-bold text-[#ff0036] mt-2">99 <span className="text-base font-normal text-[#62625b]">金币</span></div>
            <div className="text-sm text-[#91918c] line-through mt-1">原价 699 金币</div>
            <div className="text-base text-[#62625b] mt-4 pb-4 border-b border-[#dadad3]">永久有效</div>
            <p className="mt-6 text-sm text-[#62625b] text-center">持续更新课堂资源，终身会员永久有效</p>
            <Link href="/vip" className="mt-8 block w-full py-3.5 bg-[#ff0036] text-white text-sm font-bold rounded-[16px] hover:bg-[#e6002f] text-center transition">立即升级</Link>
            <p className="text-xs text-[#91918c] mt-4">75% 的人选择该套餐</p>
          </div>
        </div>
      </section>
    </div>
  );
}
