import Link from "next/link";
import { getCategories } from "@/lib/api";

export default async function CategoriesPage() {
  let categories;
  try {
    const res = await getCategories();
    categories = res.data;
  } catch {
    categories = null;
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-extrabold tracking-tight">Categories</h1>
        <p className="mt-2 text-[var(--muted)]">
          Browse open-source libraries by use case, ranked for AI agent usage.
        </p>
      </div>

      {categories ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map((cat) => (
            <Link
              key={cat.id}
              href={`/categories/${cat.slug}`}
              className="group rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-5 shadow-sm shadow-[var(--card-shadow)] transition-all hover:border-[var(--primary)]/40 hover:shadow-md hover:-translate-y-1"
            >
              <h2 className="font-semibold text-[var(--foreground)] group-hover:text-[var(--primary)] transition-colors">
                {cat.name}
              </h2>
              <p className="mt-1.5 text-sm text-[var(--muted)] line-clamp-2 leading-relaxed">
                {cat.description}
              </p>
              {cat.libraryCount != null && (
                <div className="mt-3 inline-flex items-center rounded-md bg-[var(--primary)]/10 px-2.5 py-1 text-xs font-semibold text-[var(--primary)]">
                  {cat.libraryCount} libraries
                </div>
              )}
            </Link>
          ))}
        </div>
      ) : (
        <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-10 text-center shadow-sm">
          <p className="text-[var(--muted)]">
            Unable to load categories. Is the API server running?
          </p>
          <p className="mt-2 text-sm text-[var(--muted-dim)]">
            Run{" "}
            <code className="rounded-md bg-[var(--surface)] px-2 py-0.5 text-[var(--primary)]">
              make dev-api
            </code>{" "}
            to start the backend.
          </p>
        </div>
      )}
    </div>
  );
}
