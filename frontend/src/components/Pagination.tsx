'use client';

import { useEffect, useState } from 'react';

interface Props {
  page: number;
  totalPages: number;
  onChange: (page: number) => void;
}

export default function Pagination({ page, totalPages, onChange }: Props) {
  const pageCount = Math.max(1, totalPages);
  const [jumpPage, setJumpPage] = useState('');

  useEffect(() => {
    setJumpPage('');
  }, [page, totalPages]);

  const pages: (number | string)[] = [];
  for (let i = 1; i <= pageCount; i++) {
    if (i === 1 || i === pageCount || (i >= page - 2 && i <= page + 2)) {
      pages.push(i);
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...');
    }
  }

  const goToPage = (nextPage: number) => {
    const safePage = Math.min(Math.max(1, nextPage), pageCount);
    if (safePage !== page) onChange(safePage);
  };

  const handleJump = (e: React.FormEvent) => {
    e.preventDefault();
    const nextPage = Number(jumpPage);
    if (Number.isFinite(nextPage) && nextPage >= 1) {
      goToPage(nextPage);
    }
  };

  return (
    <form onSubmit={handleJump} className="flex flex-wrap items-center justify-center gap-2 py-8">
      <button type="button" onClick={() => goToPage(page - 1)} disabled={page <= 1} className="inline-flex h-10 min-w-10 items-center justify-center rounded-card bg-surface px-3 text-sm font-semibold text-ink disabled:cursor-not-allowed disabled:text-ash">
        上一页
      </button>
      {pages.map((p, i) =>
        typeof p === 'number' ? (
          <button key={i} type="button" onClick={() => goToPage(p)} className={`inline-flex h-10 min-w-10 items-center justify-center rounded-card px-3 text-sm font-semibold ${p === page ? 'bg-primary text-white' : 'bg-surface text-ink'}`}>
            {p}
          </button>
        ) : (
          <span key={i} className="inline-flex h-10 min-w-10 items-center justify-center rounded-card bg-surface px-3 text-sm font-semibold text-ink">...</span>
        )
      )}
      <button type="button" onClick={() => goToPage(page + 1)} disabled={page >= pageCount} className="inline-flex h-10 min-w-10 items-center justify-center rounded-card bg-surface px-3 text-sm font-semibold text-ink disabled:cursor-not-allowed disabled:text-ash">
        下一页
      </button>
      <input
        type="number"
        min={1}
        max={pageCount}
        value={jumpPage}
        onChange={(e) => setJumpPage(e.target.value)}
        placeholder="页码"
        className="h-10 w-14 rounded-card border border-hairline bg-canvas text-center text-sm outline-none"
      />
      <button type="submit" className="h-10 rounded-card bg-surface px-3 text-sm font-semibold text-ink">
        跳转
      </button>
    </form>
  );
}
