import { Suspense } from 'react';
import CoursesClient from './CoursesClient';
import { normalizeAssetUrls } from '@/lib/assets';

export const dynamic = 'force-dynamic';

async function getData() {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  try {
    const [coursesRes, categoriesRes, tagsRes] = await Promise.all([
      fetch(`${baseUrl}/courses?page=1&pageSize=16`, { cache: 'no-store' }),
      fetch(`${baseUrl}/categories`, { cache: 'no-store' }),
      fetch(`${baseUrl}/tags`, { cache: 'no-store' }),
    ]);
    const [courses, categories, tags] = await Promise.all([
      coursesRes.ok ? coursesRes.json().then(r => normalizeAssetUrls(r.data || { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 })) : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 },
      categoriesRes.ok ? categoriesRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
      tagsRes.ok ? tagsRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
    ]);
    return { courses, categories, tags };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [], tags: [] };
  }
}

export default async function CoursesPage() {
  const data = await getData();
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} tags={data.tags} />
    </Suspense>
  );
}
