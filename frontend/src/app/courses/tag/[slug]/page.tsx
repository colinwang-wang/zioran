import { Suspense } from 'react';
import CoursesClient from '../../CoursesClient';
import { normalizeAssetUrls } from '@/lib/assets';

export const dynamic = 'force-dynamic';

async function getData(slug: string) {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  try {
    const tagsRes = await fetch(`${baseUrl}/tags`, { cache: 'no-store' });
    const tags = tagsRes.ok ? await tagsRes.json().then(r => normalizeAssetUrls(r.data || [])) : [];
    const tag = tags.find((t: { slug: string }) => t.slug === slug);
    if (!tag) {
      const categoriesRes = await fetch(`${baseUrl}/categories`, { cache: 'no-store' });
      const categories = categoriesRes.ok ? await categoriesRes.json().then(r => normalizeAssetUrls(r.data || [])) : [];
      return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories, tags, tagId: null };
    }
    const tagId = tag.id;
    const [coursesRes, categoriesRes] = await Promise.all([
      fetch(`${baseUrl}/courses?page=1&pageSize=16&tagId=${tagId}`, { cache: 'no-store' }),
      fetch(`${baseUrl}/categories`, { cache: 'no-store' }),
    ]);
    const [courses, categories] = await Promise.all([
      coursesRes.ok ? coursesRes.json().then(r => normalizeAssetUrls(r.data || { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 })) : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 },
      categoriesRes.ok ? categoriesRes.json().then(r => normalizeAssetUrls(r.data || [])) : [],
    ]);
    return { courses, categories, tags, tagId };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [], tags: [], tagId: null };
  }
}

export default async function TagPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const data = await getData(slug);
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} tags={data.tags} initialTagId={data.tagId} />
    </Suspense>
  );
}
