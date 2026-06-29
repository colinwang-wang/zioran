import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import CourseDetailClient from './CourseDetailClient';
import { normalizeAssetUrls } from '@/lib/assets';
import type { CourseDetail } from '@/types';

const baseUrl = 'http://127.0.0.1:8080/api/v1';

export const dynamic = 'force-dynamic';

async function getCourse(slug: string): Promise<CourseDetail | null> {
  try {
    const res = await fetch(`${baseUrl}/courses/${slug}`, { cache: 'no-store' });
    if (!res.ok) return null;
    return res.json().then(r => normalizeAssetUrls(r.data));
  } catch {
    return null;
  }
}

export async function generateMetadata({ params }: { params: { slug: string } }): Promise<Metadata> {
  const course = await getCourse(params.slug);
  return { title: course ? `${course.title} - 知猿课堂` : '课程详情 - 知猿课堂' };
}

export default async function CourseDetailPage({ params }: { params: { slug: string } }) {
  const course = await getCourse(params.slug);
  if (!course) notFound();
  return <CourseDetailClient course={course} />;
}
