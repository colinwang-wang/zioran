import Link from 'next/link';
import type { CourseListItem } from '@/types';

export default function CourseCard({ course }: { course: CourseListItem }) {
  return (
    <Link href={`/courses/${course.slug}`} className="group block rounded-card overflow-hidden bg-surface">
      <div className="aspect-[4/3] relative overflow-hidden">
        {course.cover ? (
          <img src={course.cover} alt={course.title} className="w-full h-full object-cover group-hover:scale-105 transition-transform" loading="lazy" />
        ) : (
          <div className="w-full h-full bg-secondary-bg flex items-center justify-center text-mute text-sm">暂无图片</div>
        )}
      </div>
      <div className="p-3">
        {course.category && (
          <span className="text-xs text-primary font-medium">{course.category.name}</span>
        )}
        <h3 className="text-sm font-semibold text-ink mt-1 line-clamp-2 leading-snug">{course.title}</h3>
        <div className="flex items-center justify-between mt-2">
          <span className="text-xs text-mute">{course.relative_time}</span>
          <span className="text-sm font-bold text-primary">{course.price > 0 ? `${course.price}金币` : '免费'}</span>
        </div>
      </div>
    </Link>
  );
}
