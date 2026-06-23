import Link from 'next/link';
import type { CourseListItem } from '@/types';

export default function CourseCard({ course }: { course: CourseListItem }) {
  return (
    <Link href={`/courses/${course.slug}`} className="group block overflow-hidden rounded-card bg-surface transition hover:-translate-y-1">
      <div className="relative aspect-[16/10] overflow-hidden">
        {course.cover ? (
          <img src={course.cover} alt={course.title} className="w-full h-full object-cover group-hover:scale-105 transition-transform" loading="lazy" />
        ) : (
          <div className="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#dddddd] to-[#eeeeee] text-sm text-ash">课程封面</div>
        )}
      </div>
      <div className="px-4 py-3">
        {course.category && (
          <div className="mb-1.5 text-xs font-bold text-primary">{course.category.name}</div>
        )}
        <h3 className="mb-2 line-clamp-2 text-sm font-semibold leading-[1.3] text-ink">{course.title}</h3>
        <div className="flex items-center justify-between">
          <span className="text-xs text-mute">{course.relative_time}</span>
          <span className="text-sm font-bold text-primary">{course.price > 0 ? `${course.price} 金币` : '免费'}</span>
        </div>
      </div>
    </Link>
  );
}
