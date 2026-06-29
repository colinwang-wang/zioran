import { Suspense } from 'react';
import CoursesClient from '../../CoursesClient';
import { normalizeAssetUrls } from '@/lib/assets';

export const dynamic = 'force-dynamic';

async function getData(slug: string) {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  try {
    const [categoriesRes, tagsRes] = await Promise.all([
      fetch(`${baseUrl}/categories`, { cache: 'no-store' }),
      fetch(`${baseUrl}/tags`, { cache: 'no-store' }),
    ]);
    const [categories, tags] = await Promise.all([
      categoriesRes.ok ? categoriesRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
      tagsRes.ok ? tagsRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
    ]);
    const cat = categories.find((c: { slug: string }) => c.slug === slug);
    if (!cat) {
      return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories, tags, categoryId: null };
    }
    const catId = cat.id;
    const coursesRes = await fetch(`${baseUrl}/courses?page=1&pageSize=16&categoryId=${catId}`, { cache: 'no-store' });
    const courses = coursesRes.ok ? await coursesRes.json().then(r => normalizeAssetUrls(r.data || { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 })) : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 };
    return { courses, categories, tags, categoryId: catId };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [], tags: [], categoryId: null };
  }
}

export default async function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const data = await getData(slug);
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} tags={data.tags} initialCategoryId={data.categoryId} />
    </Suspense>
  );
}
