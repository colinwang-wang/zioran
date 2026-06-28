import { Suspense } from 'react';
import CoursesClient from '../courses/CoursesClient';

async function getData() {
  const baseUrl = 'http://127.0.0.1:8080/api/v1';
  try {
    const [coursesRes, categoriesRes, tagsRes] = await Promise.all([
      fetch(`${baseUrl}/courses?page=1&pageSize=16&sort=latest`, { next: { revalidate: 60 } }),
      fetch(`${baseUrl}/categories`, { next: { revalidate: 300 } }),
      fetch(`${baseUrl}/tags`, { next: { revalidate: 300 } }),
    ]);
    const [courses, categories, tags] = await Promise.all([
      coursesRes.ok ? coursesRes.json().then(r => r.data || { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }) : { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 },
      categoriesRes.ok ? categoriesRes.json().then(r => r.data || []) : [],
      tagsRes.ok ? tagsRes.json().then(r => r.data || []) : [],
    ]);
    return { courses, categories, tags };
  } catch {
    return { courses: { items: [], total: 0, page: 1, pageSize: 16, totalPages: 0 }, categories: [], tags: [] };
  }
}

export default async function NewPage() {
  const data = await getData();
  return (
    <Suspense fallback={<div className="max-w-container mx-auto px-4 py-8">加载中...</div>}>
      <CoursesClient initialCourses={data.courses} categories={data.categories} tags={data.tags} />
    </Suspense>
  );
}
