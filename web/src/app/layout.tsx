import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";

export const metadata: Metadata = {
  title: "Opensource Curator - AI Agent Library Rankings",
  description:
    "Curated open-source library rankings optimized for AI agent usage. Score libraries by maintenance, API clarity, documentation, security, and community signals.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <nav className="sticky top-0 z-50 border-b border-[var(--card-border)] bg-[var(--surface)]/80 backdrop-blur-md px-6 py-4">
          <div className="mx-auto flex max-w-6xl items-center justify-between">
            <Link
              href="/"
              className="flex items-center gap-2 text-lg font-bold text-[var(--foreground)] hover:text-[var(--primary)] transition-colors"
            >
              <span className="inline-flex h-7 w-7 items-center justify-center rounded-lg bg-[var(--primary)] text-sm font-bold text-white">
                C
              </span>
              Opensource Curator
            </Link>
            <div className="flex items-center gap-1">
              <Link
                href="/categories"
                className="rounded-lg px-3 py-2 text-sm text-[var(--muted)] transition-colors hover:bg-[var(--card)] hover:text-[var(--foreground)]"
              >
                Categories
              </Link>
              <Link
                href="/recommend"
                className="rounded-lg px-3 py-2 text-sm text-[var(--muted)] transition-colors hover:bg-[var(--card)] hover:text-[var(--foreground)]"
              >
                Recommend
              </Link>
            </div>
          </div>
        </nav>
        <main className="mx-auto max-w-6xl px-6 py-8">{children}</main>
      </body>
    </html>
  );
}
