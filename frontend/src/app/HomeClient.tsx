'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import CourseCard from '@/components/CourseCard';
import VipCard from '@/components/VipCard';
import { getCourses } from '@/lib/services';
import type { NavItem, Banner, CourseListItem, CategoryBrief, VipPackage } from '@/types';

interface Props {
  navItems: NavItem[];
  banners: Banner[];
  latest: CourseListItem[];
  categories: CategoryBrief[];
  vipPackages: VipPackage[];
}

const fallbackCategories: CategoryBrief[] = [
  { id: 1, name: 'AIGC课堂', slug: 'aigc' },
  { id: 2, name: 'Blender课堂', slug: 'blender' },
  { id: 3, name: 'C4D课程', slug: 'c4d' },
  { id: 4, name: '手绘课程', slug: 'painting' },
  { id: 5, name: 'AE课程', slug: 'ae' },
  { id: 6, name: 'UI课程', slug: 'ui' },
];

const fallbackNavItems = [
  { id: 1, title: 'AIGC课堂', icon: 'AI', subtitle: 'AI绘画生成', url: '/courses?categoryId=1' },
  { id: 2, title: 'Blender课堂', icon: 'B', subtitle: '3D建模渲染', url: '/courses?categoryId=2' },
  { id: 3, title: 'C4D课程', icon: 'C4', subtitle: '三维动画', url: '/courses?categoryId=3' },
  { id: 4, title: '手绘课程', icon: 'PS', subtitle: '插画绘画', url: '/courses?categoryId=4' },
  { id: 5, title: 'AE课程', icon: 'AE', subtitle: '动效合成', url: '/courses?categoryId=5' },
  { id: 6, title: 'UI课程', icon: 'UI', subtitle: '界面设计', url: '/courses?categoryId=6' },
  { id: 7, title: '摄影课程', icon: '📷', subtitle: '摄影后期', url: '/courses' },
  { id: 8, title: '室内设计', icon: '🏠', subtitle: '空间规划', url: '/courses' },
];

const navSubtitles: Record<string, string> = Object.fromEntries(fallbackNavItems.map((item) => [item.title, item.subtitle]));
const navIcons: Record<string, string> = Object.fromEntries(fallbackNavItems.map((item) => [item.title, item.icon]));

function makeCourse(id: number, category: string, title: string, relativeTime: string): CourseListItem {
  return {
    id,
    title,
    subtitle: '',
    slug: `prototype-course-${id}`,
    cover: '',
    category: { id, name: category, slug: category.toLowerCase() },
    price: 2,
    vip_price: 0,
    relative_time: relativeTime,
    published_at: null,
  };
}

const fallbackLatestCourses = [
  makeCourse(101, '手绘课程', '塵蒲2026唯美古风半厚涂第5期基础课', '2小时前'),
  makeCourse(102, 'AIGC课程', 'Midjourney商业插画全流程实战2026', '5小时前'),
  makeCourse(103, 'Blender课程', 'Blender 4.0写实场景全流程教学', '1天前'),
  makeCourse(104, 'C4D课程', 'C4D电商产品渲染高级班第3期', '1天前'),
  makeCourse(105, '手绘课程', '角色设计商业插画系统训练营', '2天前'),
  makeCourse(106, 'AE课程', 'AE高级动效设计与MG动画实战', '2天前'),
  makeCourse(107, 'UI课程', 'UI/UX产品设计全流程进阶课', '3天前'),
  makeCourse(108, '平面设计', '品牌视觉设计与商业提案实战', '3天前'),
];

const fallbackClassroomCourses = [
  makeCourse(201, 'AIGC课程', 'Stable Diffusion商用级出图全攻略', '1天前'),
  makeCourse(202, 'Blender课程', 'Blender程序化建模进阶教程', '2天前'),
  makeCourse(203, 'C4D课程', 'C4D+OC渲染器产品级渲染', '3天前'),
  makeCourse(204, '手绘课程', '日系色彩速涂团练第6期', '3天前'),
  makeCourse(205, 'AE课程', 'AE表达式与脚本高级应用', '4天前'),
  makeCourse(206, '摄影课程', '人像摄影与后期调色系统课', '4天前'),
  makeCourse(207, '电商设计', '电商详情页设计与视觉营销', '5天前'),
  makeCourse(208, '室内设计', '3Dmax室内全景效果图实战', '5天前'),
];

function normalizeHref(url?: string) {
  if (!url) return '/courses';
  return url
    .replace('index.html', '/')
    .replace('courses.html', '/courses')
    .replace('guestbook.html', '/guestbook')
    .replace('vip.html', '/vip')
    .replace('login.html', '/login');
}

export default function HomeClient({ navItems, banners, latest, categories, vipPackages }: Props) {
  const [activeTab, setActiveTab] = useState<number | null>(null);
  const [tabCourses, setTabCourses] = useState<CourseListItem[]>(latest.length > 0 ? latest : fallbackClassroomCourses);
  const [tabLoadFailed, setTabLoadFailed] = useState(false);
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const visibleNavItems = navItems.length > 0
    ? navItems.slice(0, 10).map((item) => ({
        id: item.id,
        title: item.title,
        icon: item.icon || navIcons[item.title] || item.title.slice(0, 2),
        subtitle: item.subtitle || navSubtitles[item.title] || '精选课程',
        url: item.category_id ? `/courses?categoryId=${item.category_id}` : normalizeHref(item.url),
      }))
    : fallbackNavItems;
  const visibleCategories = categories.length > 0 ? categories : fallbackCategories;
  const visibleLatest = latest.length > 0 ? latest.slice(0, 8) : fallbackLatestCourses;
  const visibleTabCourses = tabCourses.slice(0, 8);
  const visibleVipPackage = vipPackages[0];

  useEffect(() => {
    if (activeTab === null) {
      setTabCourses(latest.length > 0 ? latest : fallbackClassroomCourses);
      setTabLoadFailed(latest.length === 0);
      return;
    }
    getCourses({ categoryId: activeTab, pageSize: 8 }).then((res) => {
      const courses = Array.isArray(res) ? res : res?.items || [];
      setTabCourses(courses);
      setTabLoadFailed(false);
    }).catch(() => {
      setTabCourses([]);
      setTabLoadFailed(true);
    });
  }, [activeTab, latest]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) router.push(`/courses?q=${encodeURIComponent(searchQuery.trim())}`);
  };

  return (
    <div>
      {/* 金刚区 */}
      <section className="bg-surface-soft py-8">
        <div className="mx-auto flex max-w-[1280px] gap-4 overflow-x-auto px-6">
          {visibleNavItems.map((item) => (
            <Link key={item.id} href={item.url} className="min-w-[140px] flex-shrink-0 rounded-[16px] bg-white px-6 py-5 text-center transition hover:-translate-y-0.5">
              {item.icon.startsWith('http') || item.icon.startsWith('/') ? (
                <img src={item.icon} alt="" className="mx-auto mb-3 h-12 w-12 rounded-xl object-cover" />
              ) : (
                <div className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-[#ff0036] to-[#ff6b6b] text-xl font-bold text-white">{item.icon}</div>
              )}
              <div className="text-sm font-semibold text-ink">{item.title}</div>
              <div className="mt-1 text-xs text-mute">{item.subtitle}</div>
            </Link>
          ))}
        </div>
      </section>

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
              <div className="absolute inset-0 bg-gradient-to-br from-[#1a1a2e] via-[#16213e] to-[#0f3460]" style={{ background: banners[0]?.background_color || undefined }} />
            )}
            {/* 装饰光效 */}
            <div className="absolute right-[-20%] top-[-50%] h-[400px] w-[400px] bg-[radial-gradient(circle,rgba(255,0,54,0.2),transparent_70%)]" />
            <h1 className="relative mb-3 text-[28px] font-bold tracking-[-0.8px] md:text-[44px]">知猿课堂，学有所长</h1>
            <p className="relative mb-8 text-base text-white/70">以知为基，以猿为伴，打造优质网课资源课堂</p>
            <form onSubmit={handleSearch} className="relative mx-auto max-w-[560px]">
              <input type="text" value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} placeholder="输入关键字搜索课程..." className="h-12 w-full rounded-full border-none bg-white/95 px-5 pr-28 text-base text-ink outline-none placeholder:text-[#999]" />
              <button type="submit" className="absolute right-1 top-1 h-10 rounded-full bg-[#ff0036] px-5 text-sm font-bold text-white transition hover:bg-[#e6002f]">搜索</button>
            </form>
          </div>
        </div>
      </section>

      {/* 最新发布 */}
      <section className="py-12">
        <div className="mx-auto max-w-[1280px] px-6">
          <div className="mb-6 flex items-center justify-between">
            <h2 className="text-[28px] font-bold tracking-[-1.2px] text-[#000]">最新发布</h2>
            <Link href="/courses?sort=latest" className="text-sm font-semibold text-mute hover:text-primary">查看更多 →</Link>
          </div>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-4">
            {visibleLatest.map((c) => <CourseCard key={c.id} course={c} />)}
          </div>
        </div>
      </section>

      {/* 知猿课堂（Tab切换） */}
      <section className="py-12 bg-[#f8f8f5]">
        <div className="max-w-[1280px] mx-auto px-6">
          <div className="mb-6 flex items-center justify-between">
            <h2 className="text-[28px] font-bold tracking-[-1.2px] text-[#000]">知猿课堂</h2>
            <Link href={activeTab ? `/courses?categoryId=${activeTab}` : '/courses'} className="text-sm font-semibold text-mute hover:text-primary">查看更多 →</Link>
          </div>
          <div className="flex flex-wrap gap-2 mb-6">
            <button onClick={() => setActiveTab(null)} className={`px-4 py-2 rounded-full text-sm font-bold transition ${activeTab === null ? 'bg-[#000] text-white' : 'bg-white text-[#000] hover:bg-[#e5e5e0]'}`}>全部课堂</button>
            {visibleCategories.map((cat) => (
              <button key={cat.id} onClick={() => setActiveTab(cat.id)} className={`px-4 py-2 rounded-full text-sm font-bold transition ${activeTab === cat.id ? 'bg-[#000] text-white' : 'bg-white text-[#000] hover:bg-[#e5e5e0]'}`}>{cat.name}</button>
            ))}
          </div>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {visibleTabCourses.map((c) => <CourseCard key={c.id} course={c} />)}
          </div>
          {visibleTabCourses.length === 0 && (
            <div className="rounded-[16px] bg-white py-16 text-center text-sm text-mute">
              {tabLoadFailed ? '课程加载失败，请稍后重试' : '暂无该分类课程'}
            </div>
          )}
        </div>
      </section>

      {/* VIP */}
      <section className="py-16">
        <div className="max-w-[1280px] mx-auto px-6 text-center">
          <div className="mb-6 flex justify-center">
            <h2 className="text-[28px] font-bold tracking-[-1.2px] text-[#000]">关于VIP</h2>
          </div>
          <VipCard
            name={visibleVipPackage?.name}
            price={visibleVipPackage?.price}
            originalPrice={visibleVipPackage?.original_price}
            actionHref="/vip"
          />
        </div>
      </section>
    </div>
  );
}
