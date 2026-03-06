export function ScoreBadge({ score }: { score: number }) {
  const color =
    score >= 80
      ? "var(--score-excellent)"
      : score >= 60
        ? "var(--score-good)"
        : score >= 40
          ? "var(--score-fair)"
          : "var(--score-poor)";

  return (
    <span
      className="inline-flex items-center rounded-lg px-2.5 py-1 text-xs font-bold tabular-nums"
      style={{ backgroundColor: `${color}18`, color, border: `1px solid ${color}30` }}
    >
      {score.toFixed(1)}
    </span>
  );
}

export function ScoreBar({
  label,
  score,
  weight,
}: {
  label: string;
  score: number;
  weight?: number;
}) {
  const color =
    score >= 80
      ? "var(--score-excellent)"
      : score >= 60
        ? "var(--score-good)"
        : score >= 40
          ? "var(--score-fair)"
          : "var(--score-poor)";

  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between text-sm">
        <span className="text-[var(--foreground)]">
          {label}
          {weight != null && (
            <span className="ml-1.5 text-xs text-[var(--muted-dim)]">
              {(weight * 100).toFixed(0)}%
            </span>
          )}
        </span>
        <span className="font-semibold tabular-nums" style={{ color }}>
          {score.toFixed(1)}
        </span>
      </div>
      <div className="h-2.5 overflow-hidden rounded-full bg-[var(--surface)]">
        <div
          className="h-full rounded-full transition-all duration-500"
          style={{ width: `${score}%`, backgroundColor: color }}
        />
      </div>
    </div>
  );
}
