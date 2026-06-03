import { Suspense } from 'react';
import CoursesClient from './CoursesClient';

async function getData() {
  const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
  try {
    const [coursesRes, categoriesRes] = await Promise.all([
      fetch(`${baseUrl}/courses?page=1&pageSize=16`, { next: { revalidate: 60 } }),
      fetch(`${baseUrl}/categories`, { next: { revalidate: 300 } }),
    ]);
    const [courses, categories] = await Promise.all([
      coursesRes.ok ? coursesRes.json() : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 },
      categoriesRes.ok ? categoriesRes.json() : [],
    ]);
    return { courses, categories };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [] };
  }
}

export default async function CoursesPage() {
  const data = await getData();
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} />
    </Suspense>
  );
}
