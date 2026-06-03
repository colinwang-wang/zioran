'use client';

interface Props {
  page: number;
  totalPages: number;
  onChange: (page: number) => void;
}

export default function Pagination({ page, totalPages, onChange }: Props) {
  if (totalPages <= 1) return null;

  const pages: (number | string)[] = [];
  for (let i = 1; i <= totalPages; i++) {
    if (i === 1 || i === totalPages || (i >= page - 2 && i <= page + 2)) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...');
    }
  }

  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      <button onClick={() => onChange(page - 1)} disabled={page <= 1} className="px-3 py-2 rounded-card text-sm disabled:text-stone disabled:cursor-not-allowed hover:bg-surface">
        &lt;
      </button>
      {pages.map((p, i) =>
        typeof p === 'number' ? (
          <button key={i} onClick={() => onChange(p)} className={`w-9 h-9 rounded-card text-sm font-semibold ${p === page ? 'bg-primary text-white' : 'hover:bg-surface'}`}>
            {p}
          </button>
        ) : (
          <span key={i} className="px-1 text-mute">...</span>
        )
      )}
      <button onClick={() => onChange(page + 1)} disabled={page >= totalPages} className="px-3 py-2 rounded-card text-sm disabled:text-stone disabled:cursor-not-allowed hover:bg-surface">
        &gt;
      </button>
    </div>
  );
}
