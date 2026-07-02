import { Suspense } from 'react';
import CoursesClient from './CoursesClient';
import { normalizeAssetUrls } from '@/lib/assets';

export const dynamic = 'force-dynamic';

const emptyCourses = { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 };

type CoursesSearchParams = {
  q?: string | string[];
  categoryId?: string | string[];
  tagId?: string | string[];
};

function firstParam(value?: string | string[]) {
  return Array.isArray(value) ? value[0] : value;
}

async function getData(searchParams: CoursesSearchParams = {}) {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  const keyword = (firstParam(searchParams.q) || '').trim();
  const categoryId = firstParam(searchParams.categoryId);
  const tagId = firstParam(searchParams.tagId);
  const coursesUrl = keyword
    ? `${baseUrl}/search?q=${encodeURIComponent(keyword)}&page=1&pageSize=16`
    : `${baseUrl}/courses?${new URLSearchParams({
        page: '1',
        pageSize: '16',
        ...(categoryId ? { categoryId } : {}),
        ...(tagId && !categoryId ? { tagId } : {}),
      }).toString()}`;

  try {
    const [coursesRes, categoriesRes, tagsRes] = await Promise.all([
      fetch(coursesUrl, { cache: 'no-store' }),
      fetch(`${baseUrl}/categories`, { cache: 'no-store' }),
      fetch(`${baseUrl}/tags`, { cache: 'no-store' }),
    ]);
    const [courses, categories, tags] = await Promise.all([
      coursesRes.ok ? coursesRes.json().then(r => normalizeAssetUrls(r.data || emptyCourses)) : emptyCourses,
      categoriesRes.ok ? categoriesRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
      tagsRes.ok ? tagsRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
    ]);
    return { courses, categories, tags };
  } catch {
    return { courses: emptyCourses, categories: [], tags: [] };
  }
}

export default async function CoursesPage({ searchParams }: { searchParams?: CoursesSearchParams }) {
  const data = await getData(searchParams);
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} tags={data.tags} />
    </Suspense>
  );
}
