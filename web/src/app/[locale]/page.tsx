import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";

export default function Home() {
  const t = useTranslations("home");

  return (
    <div className="space-y-16">
      <section className="py-20 text-center">
        <div className="mx-auto inline-flex items-center gap-2 rounded-full border border-[var(--card-border)] bg-[var(--card)] px-4 py-1.5 text-xs text-[var(--muted)] mb-6">
          <span className="inline-block h-1.5 w-1.5 rounded-full bg-[var(--secondary)] animate-pulse" />
          {t("badge")}
        </div>
        <h1 className="mb-4 text-5xl font-extrabold leading-tight tracking-tight">
          {t("title1")}
          <br />
          <span className="text-[var(--primary)]">{t("title2")}</span>{" "}
          {t("title3")}
        </h1>
        <p className="mx-auto max-w-xl text-lg text-[var(--muted)] leading-relaxed">
          {t("subtitle")}
        </p>
        <div className="mt-10 flex justify-center gap-4">
          <Link
            href="/categories"
            className="rounded-xl bg-[var(--primary)] px-7 py-3 text-sm font-semibold text-white shadow-lg shadow-[var(--primary)]/20 transition-all hover:bg-[var(--primary-hover)] hover:shadow-xl hover:shadow-[var(--primary)]/30 hover:-translate-y-0.5"
          >
            {t("browseCategories")}
          </Link>
          <Link
            href="/recommend"
            className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] px-7 py-3 text-sm font-semibold text-[var(--foreground)] transition-all hover:border-[var(--primary)] hover:text-[var(--primary)] hover:-translate-y-0.5"
          >
            {t("getRecommendations")}
          </Link>
        </div>
      </section>

      <section className="grid gap-6 md:grid-cols-3">
        <FeatureCard
          icon="01"
          title={t("feature1Title")}
          description={t("feature1Desc")}
          href="/scoring"
          learnMore={t("learnMore")}
        />
        <FeatureCard
          icon="02"
          title={t("feature2Title")}
          description={t("feature2Desc")}
        />
        <FeatureCard
          icon="03"
          title={t("feature3Title")}
          description={t("feature3Desc")}
        />
      </section>
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  description,
  href,
  learnMore,
}: {
  icon: string;
  title: string;
  description: string;
  href?: string;
  learnMore?: string;
}) {
  const content = (
    <>
      <div className="mb-3 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-[var(--primary)]/10 text-sm font-bold text-[var(--primary)]">
        {icon}
      </div>
      <h3 className="mb-2 text-lg font-semibold text-[var(--foreground)]">
        {title}
      </h3>
      <p className="text-sm leading-relaxed text-[var(--muted)]">
        {description}
      </p>
      {href && learnMore && (
        <span className="mt-3 inline-block text-xs font-medium text-[var(--primary)]">
          {learnMore}
        </span>
      )}
    </>
  );

  if (href) {
    return (
      <Link
        href={href}
        className="group rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-md shadow-[var(--card-shadow)] transition-all hover:border-[var(--primary)]/40 hover:shadow-lg hover:-translate-y-1"
      >
        {content}
      </Link>
    );
  }

  return (
    <div className="group rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-md shadow-[var(--card-shadow)] transition-all hover:border-[var(--primary)]/40 hover:shadow-lg hover:-translate-y-1">
      {content}
    </div>
  );
}
