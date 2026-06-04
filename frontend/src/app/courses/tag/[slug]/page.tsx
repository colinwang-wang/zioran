import { Suspense } from 'react';
import CoursesClient from '../../CoursesClient';

async function getData(slug: string) {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  try {
    const tagsRes = await fetch(`${baseUrl}/tags`, { next: { revalidate: 300 } });
    const tags = tagsRes.ok ? await tagsRes.json().then(r => r.data || []) : [];
    const tag = tags.find((t: { slug: string }) => t.slug === slug);
    if (!tag) {
      const categoriesRes = await fetch(`${baseUrl}/categories`, { next: { revalidate: 300 } });
      const categories = categoriesRes.ok ? await categoriesRes.json().then(r => r.data || []) : [];
      return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories, tagId: null };
    }
    const tagId = tag.id;
    const [coursesRes, categoriesRes] = await Promise.all([
      fetch(`${baseUrl}/courses?page=1&pageSize=16&tagId=${tagId}`, { next: { revalidate: 60 } }),
      fetch(`${baseUrl}/categories`, { next: { revalidate: 300 } }),
    ]);
    const [courses, categories] = await Promise.all([
      coursesRes.ok ? coursesRes.json().then(r => r.data || { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }) : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 },
      categoriesRes.ok ? categoriesRes.json().then(r => r.data || []) : [],
    ]);
    return { courses, categories, tagId };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [], tagId: null };
  }
}

export default async function TagPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const data = await getData(slug);
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} initialTagId={data.tagId} />
    </Suspense>
  );
}
