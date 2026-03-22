import { getTranslations } from "next-intl/server";
import { getLibraryBySlug } from "@/lib/api";
import { ScoreBadge, ScoreBar } from "@/components/ScoreBadge";
import { Link } from "@/i18n/navigation";

export default async function LibraryPage({
  params,
}: {
  params: Promise<{ registry: string; name: string }>;
}) {
  const { registry, name: rawName } = await params;
  const name = decodeURIComponent(rawName);
  const t = await getTranslations("library");

  let library;
  try {
    const res = await getLibraryBySlug(registry, name);
    library = res.data;
  } catch {
    library = null;
  }

  if (!library) {
    return (
      <div className="py-20 text-center">
        <h1 className="text-2xl font-bold">{t("notFound")}</h1>
        <p className="mt-2 text-[var(--muted)]">
          {t("notFoundDesc", { registry, name })}
        </p>
        <Link
          href="/categories"
          className="mt-6 inline-flex rounded-xl bg-[var(--primary)] px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-[var(--primary-hover)]"
        >
          {t("browseCategories")}
        </Link>
      </div>
    );
  }

  const b = library.score?.breakdown;

  return (
    <div className="space-y-8">
      <div>
        <Link
          href="/categories"
          className="inline-flex items-center gap-1 text-sm text-[var(--muted)] transition-colors hover:text-[var(--primary)]"
        >
          {t("back")}
        </Link>
        <div className="mt-3 flex items-center gap-4">
          <h1 className="text-3xl font-extrabold tracking-tight">
            {library.name}
          </h1>
          {library.score && <ScoreBadge score={library.score.overall} />}
          {library.deprecated && (
            <span className="rounded-lg bg-[var(--score-poor)]/15 border border-[var(--score-poor)]/25 px-3 py-1 text-xs font-bold text-[var(--score-poor)]">
              {t("deprecated")}
            </span>
          )}
        </div>
        <p className="mt-2 text-[var(--muted)] leading-relaxed">
          {library.description}
        </p>
        <div className="mt-4 flex flex-wrap gap-3 text-sm">
          <span className="rounded-lg bg-[var(--surface)] px-3 py-1.5 text-[var(--muted)]">
            {library.registry}/{library.packageName}
          </span>
          {library.latestVersion && (
            <span className="rounded-lg bg-[var(--surface)] px-3 py-1.5 text-[var(--muted)]">
              v{library.latestVersion}
            </span>
          )}
          {library.githubRepo && (
            <a
              href={`https://github.com/${library.githubRepo}`}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-lg bg-[var(--surface)] px-3 py-1.5 text-[var(--muted)] transition-colors hover:bg-[var(--card-hover)] hover:text-[var(--foreground)]"
            >
              {t("github")}
            </a>
          )}
        </div>
      </div>

      {b && (
        <div className="grid gap-6 md:grid-cols-2">
          <div className="space-y-5 rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
            <h2 className="text-lg font-semibold">{t("scoreBreakdown")}</h2>
            <ScoreBar
              label={t("maintenanceHealth")}
              score={b.maintenanceHealth}
              weight={0.25}
            />
            <ScoreBar
              label={t("apiClarity")}
              score={b.apiClarity}
              weight={0.2}
            />
            <ScoreBar
              label={t("docQuality")}
              score={b.docQuality}
              weight={0.15}
            />
            <ScoreBar
              label={t("securityPosture")}
              score={b.securityPosture}
              weight={0.15}
            />
            <ScoreBar
              label={t("communitySignal")}
              score={b.communitySignal}
              weight={0.15}
            />
            <ScoreBar
              label={t("deprecationSafety")}
              score={b.deprecationSafety}
              weight={0.1}
            />
          </div>

          <div className="space-y-6">
            <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
              <h2 className="text-lg font-semibold">{t("aiSuitability")}</h2>
              <div className="mt-3">
                {library.score!.overall >= 70 ? (
                  <div className="rounded-lg bg-[var(--score-excellent)]/10 border border-[var(--score-excellent)]/20 p-4">
                    <p className="text-sm font-medium text-[var(--score-excellent)]">
                      {t("aiExcellent")}
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      {t("aiExcellentDesc")}
                    </p>
                  </div>
                ) : library.score!.overall >= 50 ? (
                  <div className="rounded-lg bg-[var(--score-fair)]/10 border border-[var(--score-fair)]/20 p-4">
                    <p className="text-sm font-medium text-[var(--score-fair)]">
                      {t("aiModerate")}
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      {t("aiModerateDesc")}
                    </p>
                  </div>
                ) : (
                  <div className="rounded-lg bg-[var(--score-poor)]/10 border border-[var(--score-poor)]/20 p-4">
                    <p className="text-sm font-medium text-[var(--score-poor)]">
                      {t("aiPoor")}
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      {t("aiPoorDesc")}
                    </p>
                  </div>
                )}
              </div>
            </div>

            {library.alternatives && library.alternatives.length > 0 && (
              <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
                <h2 className="text-lg font-semibold">{t("alternatives")}</h2>
                <div className="mt-3 space-y-2">
                  {library.alternatives.map((alt) => (
                    <Link
                      key={`${alt.registry}-${alt.packageName}`}
                      href={`/library/${alt.registry}/${encodeURIComponent(alt.packageName)}`}
                      className="flex items-center justify-between rounded-lg border border-[var(--card-border)] bg-[var(--surface)] p-3 transition-all hover:border-[var(--primary)]/40 hover:bg-[var(--card-hover)]"
                    >
                      <div>
                        <span className="font-medium text-[var(--foreground)]">
                          {alt.name}
                        </span>
                        <p className="text-xs text-[var(--muted)]">
                          {alt.reason}
                        </p>
                      </div>
                      <ScoreBadge score={alt.overallScore} />
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
