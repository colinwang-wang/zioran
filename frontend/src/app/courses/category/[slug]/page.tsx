import { Suspense } from 'react';
import CoursesClient from '../../CoursesClient';

async function getData(slug: string) {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  try {
    const categoriesRes = await fetch(`${baseUrl}/categories`, { next: { revalidate: 300 } });
    const categories = categoriesRes.ok ? await categoriesRes.json().then(r => r.data || []) : [];
    const cat = categories.find((c: { slug: string }) => c.slug === slug);
    const catId = cat?.id || '';
    const coursesRes = await fetch(`${baseUrl}/courses?page=1&pageSize=16&categoryId=${catId}`, { next: { revalidate: 60 } });
    const courses = coursesRes.ok ? await coursesRes.json().then(r => r.data || { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }) : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 };
    return { courses, categories };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [] };
  }
}

export default async function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const data = await getData(slug);
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} />
    </Suspense>
  );
}
