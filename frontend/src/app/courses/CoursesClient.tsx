'use client';

import { useState, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import CourseCard from '@/components/CourseCard';
import Pagination from '@/components/Pagination';
import { getCourses, searchCourses } from '@/lib/services';
import type { PaginatedList, CourseListItem, CategoryBrief, TagBrief } from '@/types';

interface Props {
  initialCourses: PaginatedList<CourseListItem>;
  categories: CategoryBrief[];
  tags: TagBrief[];
  initialCategoryId?: number | null;
  initialTagId?: number | null;
}

export default function CoursesClient({ initialCourses, categories, tags, initialCategoryId = null, initialTagId = null }: Props) {
  const searchParams = useSearchParams();
  const queryFromUrl = searchParams.get('q') || '';
  const catFromUrl = searchParams.get('categoryId');

  const [data, setData] = useState(initialCourses);
  const [page, setPage] = useState(initialCourses.page);
  const [activeCategory, setActiveCategory] = useState<number | null>(catFromUrl ? Number(catFromUrl) : initialCategoryId);
  const [activeTagId, setActiveTagId] = useState<number | null>(initialTagId);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (queryFromUrl) {
      fetchData(1, null);
      return;
    }
    if (catFromUrl) {
      const nextCategory = Number(catFromUrl);
      setActiveCategory(nextCategory);
      setActiveTagId(null);
      fetchData(1, nextCategory, null);
    }
  }, [queryFromUrl, catFromUrl]);

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

  const handleTagChange = (tagId: number | null) => {
    setActiveTagId(tagId);
    setActiveCategory(null);
    fetchData(1, null, tagId);
  };

  return (
    <div className="max-w-container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="text-sm text-mute mb-6">
        <Link href="/" className="hover:text-primary">首页</Link> &gt; <Link href="/courses" className="hover:text-primary">知猿课堂</Link>
        {activeCategory && categories.find(c => c.id === activeCategory) && (
          <> &gt; <span>{categories.find(c => c.id === activeCategory)?.name}</span></>
        )}
        {activeTagId && tags.find(t => t.id === activeTagId) && (
          <> &gt; <span>{tags.find(t => t.id === activeTagId)?.name}</span></>
        )}
      </div>

      {queryFromUrl && (
        <h1 className="text-xl font-bold text-ink mb-4">搜索: {queryFromUrl}</h1>
      )}

      {/* Filter chips */}
      {!queryFromUrl && (
        <div className="mb-6 space-y-3">
          <div className="flex flex-wrap gap-2">
            <button onClick={() => handleCategoryChange(null)} className={`px-4 py-2 rounded-full text-sm font-bold ${activeCategory === null && activeTagId === null ? 'bg-ink text-white' : 'bg-surface text-ink'}`}>
              全部
            </button>
            {categories.map((cat) => (
              <button key={cat.id} onClick={() => handleCategoryChange(cat.id)} className={`px-4 py-2 rounded-full text-sm font-bold ${activeCategory === cat.id ? 'bg-ink text-white' : 'bg-surface text-ink'}`}>
                {cat.name}
              </button>
            ))}
          </div>
          {tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              <span className="inline-flex items-center text-xs font-semibold text-mute">标签</span>
              {tags.map((tag) => (
                <button key={tag.id} onClick={() => handleTagChange(tag.id)} className={`px-3 py-1.5 rounded-full text-xs font-bold ${activeTagId === tag.id ? 'bg-primary text-white' : 'bg-surface text-ink'}`}>
                  {tag.name}
                </button>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Grid */}
      <div className={`grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 ${loading ? 'opacity-50' : ''}`}>
        {data.items.map((c) => <CourseCard key={c.id} course={c} />)}
      </div>

      {data.items.length === 0 && !loading && (
        <div className="text-center py-20 text-mute">{queryFromUrl ? '暂无相关内容' : '暂无课程'}</div>
      )}

      <Pagination page={page} totalPages={data.totalPages} onChange={(p) => fetchData(p, activeCategory)} />
    </div>
  );
}
