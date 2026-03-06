import { getLibraryBySlug } from "@/lib/api";
import { ScoreBadge, ScoreBar } from "@/components/ScoreBadge";
import Link from "next/link";

export default async function LibraryPage({
  params,
}: {
  params: Promise<{ registry: string; name: string }>;
}) {
  const { registry, name } = await params;

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
        <h1 className="text-2xl font-bold">Library not found</h1>
        <p className="mt-2 text-[var(--muted)]">
          Could not load {registry}/{name}.
        </p>
        <Link
          href="/categories"
          className="mt-6 inline-flex rounded-xl bg-[var(--primary)] px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-[var(--primary-hover)]"
        >
          Browse categories
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
          &larr; Categories
        </Link>
        <div className="mt-3 flex items-center gap-4">
          <h1 className="text-3xl font-extrabold tracking-tight">
            {library.name}
          </h1>
          {library.score && <ScoreBadge score={library.score.overall} />}
          {library.deprecated && (
            <span className="rounded-lg bg-[var(--score-poor)]/15 border border-[var(--score-poor)]/25 px-3 py-1 text-xs font-bold text-[var(--score-poor)]">
              Deprecated
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
              GitHub &nearr;
            </a>
          )}
        </div>
      </div>

      {b && (
        <div className="grid gap-6 md:grid-cols-2">
          <div className="space-y-5 rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
            <h2 className="text-lg font-semibold">Score Breakdown</h2>
            <ScoreBar
              label="Maintenance Health"
              score={b.maintenanceHealth}
              weight={0.25}
            />
            <ScoreBar
              label="API Clarity"
              score={b.apiClarity}
              weight={0.2}
            />
            <ScoreBar
              label="Doc Quality"
              score={b.docQuality}
              weight={0.15}
            />
            <ScoreBar
              label="Security Posture"
              score={b.securityPosture}
              weight={0.15}
            />
            <ScoreBar
              label="Community Signal"
              score={b.communitySignal}
              weight={0.15}
            />
            <ScoreBar
              label="Deprecation Safety"
              score={b.deprecationSafety}
              weight={0.1}
            />
          </div>

          <div className="space-y-6">
            <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
              <h2 className="text-lg font-semibold">AI Agent Suitability</h2>
              <div className="mt-3">
                {library.score!.overall >= 70 ? (
                  <div className="rounded-lg bg-[var(--score-excellent)]/10 border border-[var(--score-excellent)]/20 p-4">
                    <p className="text-sm font-medium text-[var(--score-excellent)]">
                      Excellent fit
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      Well-suited for AI agent usage. Clear APIs, good
                      documentation, and active maintenance.
                    </p>
                  </div>
                ) : library.score!.overall >= 50 ? (
                  <div className="rounded-lg bg-[var(--score-fair)]/10 border border-[var(--score-fair)]/20 p-4">
                    <p className="text-sm font-medium text-[var(--score-fair)]">
                      Moderate fit
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      Usable by AI agents but may have gaps in documentation or
                      API clarity.
                    </p>
                  </div>
                ) : (
                  <div className="rounded-lg bg-[var(--score-poor)]/10 border border-[var(--score-poor)]/20 p-4">
                    <p className="text-sm font-medium text-[var(--score-poor)]">
                      Needs improvement
                    </p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      May present challenges for AI agent integration. Consider
                      alternatives.
                    </p>
                  </div>
                )}
              </div>
            </div>

            {library.alternatives && library.alternatives.length > 0 && (
              <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
                <h2 className="text-lg font-semibold">Alternatives</h2>
                <div className="mt-3 space-y-2">
                  {library.alternatives.map((alt) => (
                    <Link
                      key={`${alt.registry}-${alt.packageName}`}
                      href={`/library/${alt.registry}/${alt.packageName}`}
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
