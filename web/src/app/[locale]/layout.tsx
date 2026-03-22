import type { Metadata } from "next";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, getTranslations } from "next-intl/server";
import { notFound } from "next/navigation";
import { routing } from "@/i18n/routing";
import { Link } from "@/i18n/navigation";
import { LocaleSwitcher } from "@/components/LocaleSwitcher";

export const metadata: Metadata = {
  title: "Opensource Curator - AI Agent Library Rankings",
  description:
    "Curated open-source library rankings optimized for AI agent usage. Score libraries by maintenance, API clarity, documentation, security, and community signals.",
};

export default async function LocaleLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  if (!routing.locales.includes(locale as "en" | "ko")) {
    notFound();
  }

  const messages = await getMessages();
  const t = await getTranslations("nav");

  return (
    <html lang={locale}>
      <body>
        <NextIntlClientProvider messages={messages}>
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
                  {t("categories")}
                </Link>
                <Link
                  href="/scoring"
                  className="rounded-lg px-3 py-2 text-sm text-[var(--muted)] transition-colors hover:bg-[var(--card)] hover:text-[var(--foreground)]"
                >
                  {t("scoring")}
                </Link>
                <Link
                  href="/recommend"
                  className="rounded-lg px-3 py-2 text-sm text-[var(--muted)] transition-colors hover:bg-[var(--card)] hover:text-[var(--foreground)]"
                >
                  {t("recommend")}
                </Link>
                <LocaleSwitcher />
              </div>
            </div>
          </nav>
          <main className="mx-auto max-w-6xl px-6 py-8">{children}</main>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
