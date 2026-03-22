"use client";

import { useLocale } from "next-intl";
import { usePathname, useRouter } from "@/i18n/navigation";

export function LocaleSwitcher() {
  const locale = useLocale();
  const router = useRouter();
  const pathname = usePathname();

  function switchLocale(next: "en" | "ko") {
    router.replace(pathname, { locale: next });
  }

  return (
    <div className="flex items-center gap-0.5 rounded-lg border border-[var(--card-border)] bg-[var(--surface)] text-xs font-medium">
      <button
        onClick={() => switchLocale("ko")}
        className={`rounded-l-md px-2 py-1 transition-colors ${
          locale === "ko"
            ? "bg-[var(--primary)] text-white"
            : "text-[var(--muted)] hover:text-[var(--foreground)]"
        }`}
      >
        KO
      </button>
      <button
        onClick={() => switchLocale("en")}
        className={`rounded-r-md px-2 py-1 transition-colors ${
          locale === "en"
            ? "bg-[var(--primary)] text-white"
            : "text-[var(--muted)] hover:text-[var(--foreground)]"
        }`}
      >
        EN
      </button>
    </div>
  );
}
