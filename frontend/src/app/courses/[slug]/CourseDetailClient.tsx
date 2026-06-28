'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { likeCourse, downloadCourse, purchaseCourse, addFavorite, removeFavorite } from '@/lib/services';
import CourseCard from '@/components/CourseCard';
import CommentSection from './CommentSection';
import type { CourseDetail, ResourceItem } from '@/types';

export default function CourseDetailClient({ course }: { course: CourseDetail }) {
  const { isLoggedIn } = useAuth();
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(course.like_count);
  const [favorited, setFavorited] = useState(course.user_access?.is_favorited || false);
  const [resources, setResources] = useState<ResourceItem[]>(course.resources || []);

  const handleLike = async () => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      await likeCourse(course.id);
      setLiked(!liked);
      setLikeCount(liked ? likeCount - 1 : likeCount + 1);
    } catch { /* ignore */ }
  };

  const handleFavorite = async () => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      if (favorited) { await removeFavorite(course.id); } else { await addFavorite(course.id); }
      setFavorited(!favorited);
    } catch { /* ignore */ }
  };

  const [purchased, setPurchased] = useState(course.user_access?.has_purchased || false);

  const handlePurchase = async () => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      await purchaseCourse(course.id);
      setPurchased(true);
      alert('购买成功！');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '购买失败，余额不足请先充值';
      if (msg.includes('金币') || msg.includes('余额') || msg.includes('不足')) {
        window.location.href = `/user/recharge?returnTo=/courses/${course.slug}&amount=${course.price}`;
      } else {
        alert(msg);
      }
    }
  };

  const handleDownload = async () => {
    if (!isLoggedIn) { window.location.href = '/login'; return; }
    try {
      const res = await downloadCourse(course.id);
      setResources(res.resources || []);
      if (res.resources.length > 0) {
        alert(`下载链接:\n${res.resources.map(r => `${r.name}: ${r.url}${r.password ? ` 密码:${r.password}` : ''}`).join('\n')}`);
      } else {
        alert('暂无可用资源，请联系管理员');
      }
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '下载失败';
      alert(msg);
    }
  };

  const canDownload = course.user_access?.can_download || purchased;
  const isVip = course.user_access?.is_vip || false;

  return (
    <div className="max-w-container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="text-sm text-mute mb-6">
        <Link href="/" className="hover:text-primary">首页</Link> &gt;{' '}
        <Link href="/courses" className="hover:text-primary">知猿课堂</Link> &gt;{' '}
        {course.category && <><Link href={`/courses/category/${course.category.slug}`} className="hover:text-primary">{course.category.name}</Link> &gt; </>}
        <span>正文</span>
      </div>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Main content */}
        <article className="flex-1 min-w-0">
          <h1 className="text-2xl font-bold text-ink">{course.title}</h1>
          {course.subtitle && <p className="text-mute mt-1">{course.subtitle}</p>}
          <div className="flex items-center gap-3 text-sm text-mute mt-3">
            {course.published_at && <span>{new Date(course.published_at).toLocaleDateString()}</span>}
            {course.category && <span>| {course.category.name}</span>}
          </div>

          {/* Cover */}
          {course.cover && (
            <div className="mt-6 rounded-card overflow-hidden">
              <img src={course.cover} alt={course.title} className="w-full" />
            </div>
          )}

          {/* Detail */}
          {course.detail_title && <h2 className="text-lg font-bold text-ink mt-8">{course.detail_title}</h2>}
          {course.detail_subtitle && <p className="text-mute mt-1">{course.detail_subtitle}</p>}
          {course.content && (
            <div className="mt-4 prose prose-sm max-w-none" dangerouslySetInnerHTML={{ __html: course.content }} />
          )}

          {/* Download box */}
          <div className="mt-8 p-6 bg-surface rounded-card border border-hairline">
            <h3 className="font-bold text-ink">资源下载</h3>
            <p className="text-sm text-mute mt-2">下载价格: <span className="text-primary font-bold">{course.price} 金币</span></p>
            {course.vip_price === 0 && <p className="text-sm text-primary mt-1">终身VIP免费</p>}
            {canDownload ? (
              <button onClick={handleDownload} className="mt-4 px-6 py-3 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed">
                立即下载
              </button>
            ) : (
              <div className="flex gap-3 mt-4">
                <button onClick={handlePurchase} className="px-6 py-3 bg-primary text-white text-sm font-bold rounded-card hover:bg-primary-pressed">
                  立即购买
                </button>
                {!isVip && (
                  <Link href="/vip" className="px-6 py-3 bg-ink text-white text-sm font-bold rounded-card">
                    升级终身VIP
                  </Link>
                )}
              </div>
            )}
            <p className="text-xs text-mute mt-3">提示：如遇到下载链接错误、失效等情况，可以直接在个人中心&gt;提交工单说明情况等待补发链接。</p>
            {resources.length > 0 && (
              <div className="mt-4 space-y-2">
                {resources.map((r) => (
                  <div key={r.id} className="rounded-card bg-canvas border border-hairline p-3 text-sm">
                    <a href={r.url} target="_blank" rel="noreferrer" className="font-semibold text-primary hover:underline">{r.name}</a>
                    {r.password && <span className="ml-2 text-mute">提取码: {r.password}</span>}
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Actions */}
          <div className="flex items-center gap-4 mt-6 py-4 border-t border-b border-hairline">
            <button onClick={handleFavorite} className={`flex items-center gap-1 text-sm ${favorited ? 'text-primary' : 'text-mute'}`}>
              ♡ {favorited ? '已收藏' : '收藏'}
            </button>
            <button onClick={handleLike} className={`flex items-center gap-1 text-sm ${liked ? 'text-primary' : 'text-mute'}`}>
              👍 {likeCount}
            </button>
            <button onClick={() => { navigator.clipboard.writeText(window.location.href); alert('链接已复制'); }} className="flex items-center gap-1 text-sm text-mute hover:text-primary">
              ✎ 分享
            </button>
          </div>

          {/* Tags */}
          {course.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-4">
              {course.tags.map((tag) => (
                <Link key={tag.id} href={`/courses/tag/${tag.slug}`} className="px-3 py-1 bg-surface rounded-full text-xs text-mute hover:text-primary">
                  {tag.name}
                </Link>
              ))}
            </div>
          )}

          {/* Prev/Next */}
          <div className="flex justify-between mt-8 py-4 border-t border-hairline text-sm">
            {course.prev_course ? <Link href={`/courses/${course.prev_course.slug}`} className="text-primary hover:underline">← {course.prev_course.title}</Link> : <span />}
            {course.next_course ? <Link href={`/courses/${course.next_course.slug}`} className="text-primary hover:underline">{course.next_course.title} →</Link> : <span />}
          </div>

          {/* Related */}
          {course.related_courses.length > 0 && (
            <div className="mt-8 pt-6 border-t border-hairline">
              <h3 className="font-bold text-ink mb-4">猜你喜欢</h3>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {course.related_courses.slice(0, 3).map((c) => <CourseCard key={c.id} course={c} />)}
              </div>
            </div>
          )}

          {/* Comments */}
          <CommentSection courseId={course.id} />
        </article>

        {/* Sidebar */}
        <aside className="w-full lg:w-72 shrink-0">
          <div className="sticky top-20 space-y-4">
            <div className="p-6 bg-surface rounded-card border border-hairline">
              <div className="text-2xl font-bold text-primary">{course.price} <span className="text-sm font-normal">金币</span></div>
              {course.vip_price === 0 && <p className="text-sm text-primary mt-1">终身VIP免费</p>}
              <Link href="/vip" className="mt-4 block w-full py-3 bg-ink text-white text-sm font-bold rounded-card text-center">升级终身VIP</Link>
              {canDownload ? (
                <button onClick={handleDownload} className="mt-2 block w-full py-3 bg-primary text-white text-sm font-bold rounded-card">立即下载</button>
              ) : (
                <button onClick={handlePurchase} className="mt-2 block w-full py-3 bg-primary text-white text-sm font-bold rounded-card">立即购买</button>
              )}
              <p className="text-xs text-mute mt-3 text-center">提示：如遇到下载链接错误、失效等情况，可以直接在个人中心&gt;提交工单说明情况等待补发链接。</p>
            </div>

            {/* Hot tags */}
            {course.tags.length > 0 && (
              <div className="p-4 bg-surface rounded-card">
                <h4 className="text-sm font-bold mb-3">热门标签</h4>
                <div className="flex flex-wrap gap-2">
                  {course.tags.map((tag) => (
                    <Link key={tag.id} href={`/courses/tag/${tag.slug}`} className="px-3 py-1 bg-canvas rounded-full text-xs text-mute hover:text-primary border border-hairline">
                      {tag.name}
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </div>
        </aside>
      </div>
    </div>
  );
}
