"use client";

import { useState } from "react";
import Link from "next/link";
import type { Recommendation, RecommendResponse } from "@/lib/api";
import { ScoreBadge } from "@/components/ScoreBadge";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "";

export default function RecommendPage() {
  const [task, setTask] = useState("");
  const [prefer, setPrefer] = useState("");
  const [results, setResults] = useState<RecommendResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!task.trim()) return;

    setLoading(true);
    setError("");
    setResults(null);

    try {
      let url = `${API_BASE}/v1/recommend?task=${encodeURIComponent(task)}`;
      if (prefer) url += `&prefer=${encodeURIComponent(prefer)}`;
      const res = await fetch(url);
      if (!res.ok) throw new Error(`API error: ${res.status}`);
      const data: RecommendResponse = await res.json();
      setResults(data);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to get recommendations"
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-extrabold tracking-tight">
          Get Recommendations
        </h1>
        <p className="mt-2 text-[var(--muted)]">
          Describe what you need and get AI-agent-optimized library
          recommendations.
        </p>
      </div>

      <form
        onSubmit={handleSubmit}
        className="space-y-5 rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]"
      >
        <div>
          <label
            htmlFor="task"
            className="mb-1.5 block text-sm font-medium text-[var(--foreground)]"
          >
            What do you need?
          </label>
          <input
            id="task"
            type="text"
            value={task}
            onChange={(e) => setTask(e.target.value)}
            placeholder="e.g., make HTTP requests in TypeScript"
            className="w-full rounded-xl border border-[var(--card-border)] bg-[var(--input-bg)] px-4 py-3 text-[var(--foreground)] placeholder-[var(--muted-dim)] transition-colors focus:border-[var(--primary)] focus:outline-none focus:ring-1 focus:ring-[var(--primary)]/30"
          />
        </div>
        <div className="flex items-end gap-4">
          <div className="flex-1">
            <label
              htmlFor="prefer"
              className="mb-1.5 block text-sm font-medium text-[var(--foreground)]"
            >
              Preference
            </label>
            <select
              id="prefer"
              value={prefer}
              onChange={(e) => setPrefer(e.target.value)}
              className="w-full rounded-xl border border-[var(--card-border)] bg-[var(--input-bg)] px-4 py-3 text-[var(--foreground)] transition-colors focus:border-[var(--primary)] focus:outline-none focus:ring-1 focus:ring-[var(--primary)]/30"
            >
              <option value="">Best overall</option>
              <option value="lightweight">Lightweight</option>
              <option value="stable">Stable</option>
              <option value="secure">Secure</option>
              <option value="popular">Popular</option>
            </select>
          </div>
          <button
            type="submit"
            disabled={loading || !task.trim()}
            className="shrink-0 rounded-xl bg-[var(--primary)] px-7 py-3 text-sm font-semibold text-white shadow-lg shadow-[var(--primary)]/20 transition-all hover:bg-[var(--primary-hover)] hover:shadow-xl disabled:opacity-40 disabled:shadow-none disabled:hover:bg-[var(--primary)]"
          >
            {loading ? "Searching..." : "Search"}
          </button>
        </div>
      </form>

      {error && (
        <div className="rounded-xl border border-[var(--score-poor)]/30 bg-[var(--score-poor)]/10 p-4 text-sm text-[var(--score-poor)]">
          {error}
        </div>
      )}

      {results && (
        <div className="space-y-4">
          {results.query.matchedCategories.length > 0 && (
            <div className="flex flex-wrap gap-2">
              <span className="text-sm text-[var(--muted)]">Matched:</span>
              {results.query.matchedCategories.map((cat) => (
                <span
                  key={cat}
                  className="rounded-md bg-[var(--secondary)]/10 border border-[var(--secondary)]/20 px-2.5 py-0.5 text-xs font-medium text-[var(--secondary)]"
                >
                  {cat}
                </span>
              ))}
            </div>
          )}

          {results.data.length === 0 ? (
            <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-8 text-center">
              <p className="text-[var(--muted)]">
                No matching libraries found. Try a different description.
              </p>
            </div>
          ) : (
            <div className="space-y-3">
              {results.data.map((rec: Recommendation) => (
                <Link
                  key={rec.library.id}
                  href={`/library/${rec.library.registry}/${rec.library.packageName}`}
                  className="group flex items-center gap-4 rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-4 shadow-sm shadow-[var(--card-shadow)] transition-all hover:border-[var(--primary)]/40 hover:shadow-md hover:-translate-y-0.5"
                >
                  <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-[var(--primary)]/10 text-sm font-bold text-[var(--primary)]">
                    {rec.rank}
                  </span>
                  <div className="min-w-0 flex-1">
                    <span className="font-semibold text-[var(--foreground)] group-hover:text-[var(--primary)] transition-colors">
                      {rec.library.name}
                    </span>
                    <p className="mt-0.5 text-xs text-[var(--muted)] line-clamp-1">
                      {rec.matchReason}
                    </p>
                  </div>
                  <div className="shrink-0">
                    {rec.library.score && (
                      <ScoreBadge score={rec.library.score.overall} />
                    )}
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
