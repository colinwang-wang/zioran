import Link from 'next/link';

export default function Footer() {
  const links = ['AIGC', 'Blender', 'C4D', 'AE', 'UI', '手绘', '摄影', '室内设计', '平面设计', '电商设计', '3Dmax', 'zbrush'];

  return (
    <footer className="border-t border-hairline bg-surface-soft py-8">
      <div className="mx-auto max-w-container px-6">
        <div className="flex flex-col items-start justify-between gap-4 md:flex-row md:items-center">
          <div className="flex flex-wrap gap-4">
            {links.map((label) => (
              <Link key={label} href="/courses" className="text-xs text-mute hover:text-primary">{label}</Link>
            ))}
          </div>
          <div className="text-xs text-ash">© 2026 知猿 zioran.com</div>
        </div>
      </div>
    </footer>
  );
}
