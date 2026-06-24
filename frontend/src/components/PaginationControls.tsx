export default function PaginationControls({
  page,
  totalPages,
  totalItems,
  onPage,
}: {
  page: number;
  totalPages: number;
  totalItems: number;
  onPage: (page: number) => void;
}) {
  if (totalItems === 0) return null;

  return (
    <div className="flex flex-wrap items-center justify-between gap-3 border-t border-slate-100 px-4 py-3 text-sm dark:border-slate-800">
      <span className="text-slate-500">
        Page {page} of {totalPages} • {totalItems} records
      </span>
      <div className="flex gap-2">
        <button
          onClick={() => onPage(Math.max(1, page - 1))}
          disabled={page === 1}
          className="btn-secondary"
        >
          Previous
        </button>
        <button
          onClick={() => onPage(Math.min(totalPages, page + 1))}
          disabled={page === totalPages}
          className="btn-secondary"
        >
          Next
        </button>
      </div>
    </div>
  );
}
