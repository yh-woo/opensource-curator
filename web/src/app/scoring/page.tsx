import Link from "next/link";

export default function ScoringPage() {
  return (
    <div className="space-y-10">
      <div>
        <Link
          href="/"
          className="inline-flex items-center gap-1 text-sm text-[var(--muted)] transition-colors hover:text-[var(--primary)]"
        >
          &larr; Home
        </Link>
        <h1 className="mt-3 text-3xl font-extrabold tracking-tight">
          6-Metric Scoring System
        </h1>
        <p className="mt-2 text-[var(--muted)] leading-relaxed max-w-2xl">
          Every library is scored on 6 metrics that matter most for AI agent
          usability. Scores are computed from live GitHub and npm data, not
          static popularity lists.
        </p>
      </div>

      <div className="space-y-6">
        <MetricSection
          number="01"
          title="Maintenance Health"
          weight={25}
          color="var(--score-excellent)"
          description="How actively is this library being maintained? Abandoned libraries are a risk for any project, but especially for AI agents that need reliable, up-to-date tools."
          factors={[
            "Commit frequency in the last 90 days",
            "Time since the last release",
            "Open issue response rate",
            "Pull request merge cadence",
            "Number of active contributors",
          ]}
          whyItMatters="AI agents often run unattended. A library that stops receiving updates can introduce silent vulnerabilities or incompatibilities that no one is around to fix."
        />

        <MetricSection
          number="02"
          title="API Clarity"
          weight={20}
          color="var(--score-good)"
          description="How easy is it for an AI agent to correctly use this library's API? Clear, consistent APIs reduce hallucination and improve code generation accuracy."
          factors={[
            "TypeScript type definitions (built-in or @types/)",
            "Named exports vs default exports",
            "Consistent naming conventions",
            "Number of top-level exports (fewer = clearer)",
            "Predictable function signatures",
          ]}
          whyItMatters="AI agents generate code by pattern matching. Libraries with clear, typed APIs produce fewer errors and require less human review."
        />

        <MetricSection
          number="03"
          title="Documentation Quality"
          weight={15}
          color="var(--primary)"
          description="Is the documentation comprehensive enough for an AI agent to understand usage patterns? Good docs mean better generated code."
          factors={[
            "README length and structure",
            "Code examples in documentation",
            "API reference completeness",
            "Getting started / quickstart guide",
            "Changelog availability",
          ]}
          whyItMatters="AI agents learn usage patterns from documentation. Well-documented libraries with rich examples lead to higher quality code generation."
        />

        <MetricSection
          number="04"
          title="Security Posture"
          weight={15}
          color="var(--score-fair)"
          description="How secure is this library to depend on? Security matters for all projects, but especially when AI agents auto-install dependencies."
          factors={[
            "Known vulnerability count (CVEs)",
            "Dependency count (fewer = smaller attack surface)",
            "Time to patch security issues",
            "Security policy presence",
            "npm audit status",
          ]}
          whyItMatters="When an AI agent recommends a library, it should not introduce security risks. Fewer dependencies and fast patching reduce exposure."
        />

        <MetricSection
          number="05"
          title="Community Signal"
          weight={15}
          color="var(--secondary)"
          description="How strong is the community around this library? Community signals indicate long-term viability and ecosystem support."
          factors={[
            "GitHub stars and growth trend",
            "npm weekly download count",
            "Number of dependents (other packages using it)",
            "Stack Overflow question volume",
            "Contributor diversity",
          ]}
          whyItMatters="Libraries with strong communities get better support, more plugins, and are less likely to be abandoned. AI agents benefit from well-tested, battle-proven tools."
        />

        <MetricSection
          number="06"
          title="Deprecation Safety"
          weight={10}
          color="var(--score-poor)"
          description="Is this library at risk of deprecation? Using a deprecated library wastes time and creates migration burden."
          factors={[
            "Official npm deprecation flag",
            "GitHub repository archived status",
            "Successor library availability",
            "Migration guide existence",
            "Maintenance trajectory (declining activity)",
          ]}
          whyItMatters="AI agents should never recommend deprecated libraries. This metric catches libraries that are winding down before they officially announce deprecation."
        />
      </div>

      <div className="rounded-xl border border-[var(--card-border)] bg-[var(--card)] p-6 shadow-sm shadow-[var(--card-shadow)]">
        <h2 className="text-lg font-semibold">How Scores Are Computed</h2>
        <div className="mt-4 space-y-3 text-sm text-[var(--muted)] leading-relaxed">
          <p>
            Each metric produces a score from 0 to 100. The overall score is a
            weighted average using the weights shown above.
          </p>
          <p>
            Data is collected from the GitHub API and npm registry in real-time.
            Scores are refreshed during each collection run (daily when the
            worker is active, or manually via{" "}
            <code className="rounded-md bg-[var(--surface)] px-1.5 py-0.5 text-[var(--primary)]">
              make collect
            </code>
            ).
          </p>
          <p>
            Libraries flagged as deprecated (via npm or GitHub archived status)
            receive an overall score of 0, regardless of individual metric
            scores.
          </p>
        </div>
      </div>

      <div className="flex justify-center">
        <Link
          href="/categories"
          className="rounded-xl bg-[var(--primary)] px-7 py-3 text-sm font-semibold text-white shadow-lg shadow-[var(--primary)]/20 transition-all hover:bg-[var(--primary-hover)] hover:shadow-xl hover:-translate-y-0.5"
        >
          Browse Scored Libraries
        </Link>
      </div>
    </div>
  );
}

function MetricSection({
  number,
  title,
  weight,
  color,
  description,
  factors,
  whyItMatters,
}: {
  number: string;
  title: string;
  weight: number;
  color: string;
  description: string;
  factors: string[];
  whyItMatters: string;
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
                {weight}% weight
              </span>
            </div>
            <p className="mt-2 text-[var(--muted)] leading-relaxed">
              {description}
            </p>
          </div>

          <div>
            <h3 className="text-sm font-semibold text-[var(--foreground)]">
              What we measure
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
              Why it matters for AI agents
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
