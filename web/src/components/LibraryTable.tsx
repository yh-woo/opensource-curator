import Link from "next/link";
import type { Library } from "@/lib/api";
import { ScoreBadge } from "./ScoreBadge";

export function LibraryTable({ libraries }: { libraries: Library[] }) {
  const sorted = [...libraries].sort(
    (a, b) => (b.score?.overall ?? 0) - (a.score?.overall ?? 0)
  );

  return (
    <div className="space-y-3">
      {sorted.map((lib, i) => (
        <Link
          key={lib.id}
          href={`/library/${lib.registry}/${lib.packageName}`}
          className="group flex items-center gap-4 rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-4 shadow-sm shadow-[var(--card-shadow)] transition-all hover:border-[var(--primary)]/40 hover:shadow-md hover:-translate-y-0.5"
        >
          <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-[var(--surface)] text-sm font-bold text-[var(--muted)]">
            {i + 1}
          </span>
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2">
              <span className="font-semibold text-[var(--foreground)] group-hover:text-[var(--primary)] transition-colors">
                {lib.name}
              </span>
              {lib.deprecated && (
                <span className="rounded-md bg-[var(--score-poor)]/15 border border-[var(--score-poor)]/25 px-2 py-0.5 text-[10px] font-bold text-[var(--score-poor)]">
                  DEPRECATED
                </span>
              )}
            </div>
            <p className="mt-0.5 text-xs text-[var(--muted)] line-clamp-1">
              {lib.description}
            </p>
          </div>
          <div className="hidden sm:flex items-center gap-4 text-xs text-[var(--muted-dim)]">
            {lib.score && (
              <>
                <MetricPill label="Maint" value={lib.score.breakdown.maintenanceHealth} />
                <MetricPill label="API" value={lib.score.breakdown.apiClarity} />
                <MetricPill label="Docs" value={lib.score.breakdown.docQuality} />
              </>
            )}
          </div>
          <div className="shrink-0">
            {lib.score ? (
              <ScoreBadge score={lib.score.overall} />
            ) : (
              <span className="text-xs text-[var(--muted-dim)]">--</span>
            )}
          </div>
        </Link>
      ))}
    </div>
  );
}

function MetricPill({ label, value }: { label: string; value: number }) {
  return (
    <div className="flex flex-col items-center">
      <span className="text-[10px] text-[var(--muted-dim)]">{label}</span>
      <span className="font-medium text-[var(--muted)] tabular-nums">
        {value.toFixed(0)}
      </span>
    </div>
  );
}
