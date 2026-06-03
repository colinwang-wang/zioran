import HomeClient from './HomeClient';
import api from '@/lib/api';

async function getData() {
  try {
    const [navRes, bannerRes, latestRes, categoriesRes] = await Promise.all([
      fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/home/nav-items`, { next: { revalidate: 60 } }),
      fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/home/banners`, { next: { revalidate: 60 } }),
      fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/courses/latest`, { next: { revalidate: 60 } }),
      fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}/categories`, { next: { revalidate: 60 } }),
    ]);
    const [navItems, banners, latest, categories] = await Promise.all([
      navRes.ok ? navRes.json() : [],
      bannerRes.ok ? bannerRes.json() : [],
      latestRes.ok ? latestRes.json() : [],
      categoriesRes.ok ? categoriesRes.json() : [],
    ]);
    return { navItems, banners, latest, categories };
  } catch {
    return { navItems: [], banners: [], latest: [], categories: [] };
  }
}

export default async function HomePage() {
  const data = await getData();
  return <HomeClient {...data} />;
}
