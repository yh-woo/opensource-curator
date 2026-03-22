import { useTranslations } from "next-intl";
import { Link } from "@/i18n/navigation";

export default function ScoringPage() {
  const t = useTranslations("scoring");

  return (
    <div className="space-y-10">
      <div>
        <Link
          href="/"
          className="inline-flex items-center gap-1 text-sm text-[var(--muted)] transition-colors hover:text-[var(--primary)]"
        >
          {t("back")}
        </Link>
        <h1 className="mt-3 text-3xl font-extrabold tracking-tight">
          {t("title")}
        </h1>
        <p className="mt-2 text-[var(--muted)] leading-relaxed max-w-2xl">
          {t("subtitle")}
        </p>
      </div>

      <div className="space-y-6">
        <MetricSection
          number="01"
          title={t("m1Title")}
          weight={25}
          weightLabel={t("weight", { weight: 25 })}
          color="var(--score-excellent)"
          description={t("m1Desc")}
          factors={[
            t("m1f1"),
            t("m1f2"),
            t("m1f3"),
            t("m1f4"),
            t("m1f5"),
          ]}
          whyItMatters={t("m1Why")}
          whatWeMeasure={t("whatWeMeasure")}
          whyItMattersLabel={t("whyItMatters")}
        />

        <MetricSection
          number="02"
          title={t("m2Title")}
          weight={20}
          weightLabel={t("weight", { weight: 20 })}
          color="var(--score-good)"
          description={t("m2Desc")}
          factors={[
            t("m2f1"),
            t("m2f2"),
            t("m2f3"),
            t("m2f4"),
            t("m2f5"),
          ]}
          whyItMatters={t("m2Why")}
          whatWeMeasure={t("whatWeMeasure")}
          whyItMattersLabel={t("whyItMatters")}
        />

        <MetricSection
          number="03"
          title={t("m3Title")}
          weight={15}
          weightLabel={t("weight", { weight: 15 })}
          color="var(--primary)"
          description={t("m3Desc")}
          factors={[
            t("m3f1"),
            t("m3f2"),
            t("m3f3"),
            t("m3f4"),
            t("m3f5"),
          ]}
          whyItMatters={t("m3Why")}
          whatWeMeasure={t("whatWeMeasure")}
          whyItMattersLabel={t("whyItMatters")}
        />

        <MetricSection
          number="04"
          title={t("m4Title")}
          weight={15}
          weightLabel={t("weight", { weight: 15 })}
          color="var(--score-fair)"
          description={t("m4Desc")}
          factors={[
            t("m4f1"),
            t("m4f2"),
            t("m4f3"),
            t("m4f4"),
            t("m4f5"),
          ]}
          whyItMatters={t("m4Why")}
          whatWeMeasure={t("whatWeMeasure")}
          whyItMattersLabel={t("whyItMatters")}
        />

        <MetricSection
          number="05"
          title={t("m5Title")}
          weight={15}
          weightLabel={t("weight", { weight: 15 })}
          color="var(--secondary)"
          description={t("m5Desc")}
          factors={[
            t("m5f1"),
            t("m5f2"),
            t("m5f3"),
            t("m5f4"),
            t("m5f5"),
          ]}
          whyItMatters={t("m5Why")}
          whatWeMeasure={t("whatWeMeasure")}
          whyItMattersLabel={t("whyItMatters")}
        />

        <MetricSection
          number="06"
          title={t("m6Title")}
          weight={10}
          weightLabel={t("weight", { weight: 10 })}
          color="var(--score-poor)"
          description={t("m6Desc")}
          factors={[
            t("m6f1"),
            t("m6f2"),
            t("m6f3"),
            t("m6f4"),
            t("m6f5"),
          ]}
          whyItMatters={t("m6Why")}
          whatWeMeasure={t("whatWeMeasure")}
          whyItMattersLabel={t("whyItMatters")}
        />
      </div>

      <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
        <h2 className="text-lg font-semibold">{t("howTitle")}</h2>
        <div className="mt-4 space-y-3 text-sm text-[var(--muted)] leading-relaxed">
          <p>{t("howP1")}</p>
          <p>
            {t.rich("howP2", {
              command: () => (
                <code className="rounded-md bg-[var(--surface)] px-1.5 py-0.5 text-[var(--primary)]">
                  make collect
                </code>
              ),
            })}
          </p>
          <p>{t("howP3")}</p>
        </div>
      </div>

      <div className="flex justify-center">
        <Link
          href="/categories"
          className="rounded-xl bg-[var(--primary)] px-7 py-3 text-sm font-semibold text-white shadow-lg shadow-[var(--primary)]/20 transition-all hover:bg-[var(--primary-hover)] hover:shadow-xl hover:-translate-y-0.5"
        >
          {t("browseScoredLibraries")}
        </Link>
      </div>
    </div>
  );
}

function MetricSection({
  number,
  title,
  weight,
  weightLabel,
  color,
  description,
  factors,
  whyItMatters,
  whatWeMeasure,
  whyItMattersLabel,
}: {
  number: string;
  title: string;
  weight: number;
  weightLabel: string;
  color: string;
  description: string;
  factors: string[];
  whyItMatters: string;
  whatWeMeasure: string;
  whyItMattersLabel: string;
}) {
  return (
    <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
      <div className="flex items-start gap-4">
        <div
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-sm font-bold"
          style={{ backgroundColor: `${color}18`, color }}
        >
          {number}
        </div>
        <div className="flex-1 space-y-4">
          <div>
            <div className="flex items-center gap-3">
              <h2 className="text-xl font-semibold">{title}</h2>
              <span
                className="rounded-md px-2.5 py-0.5 text-xs font-bold"
                style={{
                  backgroundColor: `${color}18`,
                  color,
                  border: `1px solid ${color}30`,
                }}
              >
                {weightLabel}
              </span>
            </div>
            <p className="mt-2 text-[var(--muted)] leading-relaxed">
              {description}
            </p>
          </div>

          <div>
            <h3 className="text-sm font-semibold text-[var(--foreground)]">
              {whatWeMeasure}
            </h3>
            <ul className="mt-2 space-y-1">
              {factors.map((f, i) => (
                <li
                  key={i}
                  className="flex items-start gap-2 text-sm text-[var(--muted)]"
                >
                  <span
                    className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full"
                    style={{ backgroundColor: color }}
                  />
                  {f}
                </li>
              ))}
            </ul>
          </div>

          <div className="rounded-lg bg-[var(--surface)] p-4">
            <h3 className="text-sm font-semibold text-[var(--foreground)]">
              {whyItMattersLabel}
            </h3>
            <p className="mt-1 text-sm text-[var(--muted)] leading-relaxed">
              {whyItMatters}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
