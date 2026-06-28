export const dynamic = 'force-dynamic';
import HomeClient from './HomeClient';

async function getData() {
  const API = 'http://127.0.0.1:8080/api/v1';
  try {
    const [navRes, bannerRes, latestRes, categoriesRes, vipRes] = await Promise.all([
      fetch(API + '/home/nav-items', { cache: 'no-store' }),
      fetch(API + '/home/banners', { cache: 'no-store' }),
      fetch(API + '/courses/latest', { cache: 'no-store' }),
      fetch(API + '/categories', { cache: 'no-store' }),
      fetch(API + '/vip/packages', { cache: 'no-store' }),
    ]);
    const extract = async (res: Response) => {
      if (!res.ok) return [];
      const json = await res.json();
      return json.data || [];
    };
    const [navItems, banners, latest, categories, vipPackages] = await Promise.all([
      extract(navRes), extract(bannerRes), extract(latestRes), extract(categoriesRes), extract(vipRes),
    ]);
    return { navItems, banners, latest, categories, vipPackages };
  } catch {
    return { navItems: [], banners: [], latest: [], categories: [], vipPackages: [] };
  }
}

export default async function HomePage() {
  const data = await getData();
  return <HomeClient {...data} />;
}
