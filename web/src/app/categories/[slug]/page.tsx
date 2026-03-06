import { getCategory } from "@/lib/api";
import { LibraryTable } from "@/components/LibraryTable";
import Link from "next/link";

export default async function CategoryPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;

  let category;
  try {
    const res = await getCategory(slug);
    category = res.data;
  } catch {
    category = null;
  }

  if (!category) {
    return (
      <div className="py-20 text-center">
        <h1 className="text-2xl font-bold">Category not found</h1>
        <p className="mt-2 text-[var(--muted)]">
          Could not load category &ldquo;{slug}&rdquo;.
        </p>
        <Link
          href="/categories"
          className="mt-6 inline-flex rounded-xl bg-[var(--primary)] px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-[var(--primary-hover)]"
        >
          Back to categories
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <Link
          href="/categories"
          className="inline-flex items-center gap-1 text-sm text-[var(--muted)] transition-colors hover:text-[var(--primary)]"
        >
          &larr; Categories
        </Link>
        <h1 className="mt-2 text-3xl font-extrabold tracking-tight">
          {category.name}
        </h1>
        <p className="mt-1 text-[var(--muted)]">{category.description}</p>
      </div>

      {category.libraries && category.libraries.length > 0 ? (
        <LibraryTable libraries={category.libraries} />
      ) : (
        <p className="text-[var(--muted)]">
          No libraries in this category yet.
        </p>
      )}
    </div>
  );
}
