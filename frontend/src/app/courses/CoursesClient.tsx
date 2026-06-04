'use client';

import { useState } from 'react';
import { useSearchParams } from 'next/navigation';
import CourseCard from '@/components/CourseCard';
import Pagination from '@/components/Pagination';
import { getCourses, searchCourses } from '@/lib/services';
import type { PaginatedList, CourseListItem, CategoryBrief } from '@/types';

interface Props {
  initialCourses: PaginatedList<CourseListItem>;
  categories: CategoryBrief[];
  initialCategoryId?: number | null;
  initialTagId?: number | null;
}

export default function CoursesClient({ initialCourses, categories, initialCategoryId = null, initialTagId = null }: Props) {
  const searchParams = useSearchParams();
  const queryFromUrl = searchParams.get('q') || '';
  const catFromUrl = searchParams.get('categoryId');

  const [data, setData] = useState(initialCourses);
  const [page, setPage] = useState(initialCourses.page);
  const [activeCategory, setActiveCategory] = useState<number | null>(catFromUrl ? Number(catFromUrl) : initialCategoryId);
  const [activeTagId, setActiveTagId] = useState<number | null>(initialTagId);
  const [loading, setLoading] = useState(false);

  const fetchData = async (p: number, catId: number | null, tagId: number | null = activeTagId) => {
    setLoading(true);
    try {
      if (queryFromUrl) {
        const res = await searchCourses({ q: queryFromUrl, page: p, pageSize: 16 });
        setData(res);
      } else {
        const res = await getCourses({ page: p, pageSize: 16, categoryId: catId || undefined, tagId: catId ? undefined : tagId || undefined });
        setData(res);
      }
      setPage(p);
    } catch { /* ignore */ }
    setLoading(false);
  };

  const handleCategoryChange = (catId: number | null) => {
    setActiveCategory(catId);
    setActiveTagId(null);
    fetchData(1, catId, null);
  };

  return (
    <div className="max-w-container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="text-sm text-mute mb-6">
        <span>首页</span> &gt; <span>知猿课堂</span>
        {activeCategory && categories.find(c => c.id === activeCategory) && (
          <> &gt; <span>{categories.find(c => c.id === activeCategory)?.name}</span></>
        )}
      </div>

      {queryFromUrl && (
        <h1 className="text-xl font-bold text-ink mb-4">搜索: {queryFromUrl}</h1>
      )}

      {/* Filter chips */}
      {!queryFromUrl && (
        <div className="flex flex-wrap gap-2 mb-6">
          <button onClick={() => handleCategoryChange(null)} className={`px-4 py-2 rounded-full text-sm font-bold ${activeCategory === null ? 'bg-ink text-white' : 'bg-surface text-ink'}`}>
            全部
          </button>
          {categories.map((cat) => (
            <button key={cat.id} onClick={() => handleCategoryChange(cat.id)} className={`px-4 py-2 rounded-full text-sm font-bold ${activeCategory === cat.id ? 'bg-ink text-white' : 'bg-surface text-ink'}`}>
              {cat.name}
            </button>
          ))}
        </div>
      )}

      {/* Grid */}
      <div className={`grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 ${loading ? 'opacity-50' : ''}`}>
        {data.items.map((c) => <CourseCard key={c.id} course={c} />)}
      </div>

      {data.items.length === 0 && !loading && (
        <div className="text-center py-20 text-mute">暂无课程</div>
      )}

      <Pagination page={page} totalPages={data.totalPages} onChange={(p) => fetchData(p, activeCategory)} />
    </div>
  );
}
