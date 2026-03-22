/**
 * Renders plain text with markdown-style links [text](url)
 * converted to clickable <a> tags. Also handles **bold** text,
 * including nested **[link](url)** patterns.
 */
export function MarkdownText({
  text,
  className,
}: {
  text: string;
  className?: string;
}) {
  return <span className={className}>{parseMarkdown(text, true)}</span>;
}

function parseMarkdown(text: string, renderLinks: boolean) {
  // Match **bold content** (may contain links) and [text](url)
  const parts = text.split(/(\*\*.*?\*\*|\[.*?\]\(.*?\))/g);

  return parts.map((part, i) => {
    // Bold: **text** — may contain [link](url) inside
    const boldMatch = part.match(/^\*\*(.*?)\*\*$/);
    if (boldMatch) {
      const inner = boldMatch[1];
      // Check if bold content contains a link
      const innerLinkMatch = inner.match(/^\[(.*?)\]\((.*?)\)$/);
      if (innerLinkMatch && renderLinks) {
        return (
          <a
            key={i}
            href={innerLinkMatch[2]}
            target="_blank"
            rel="noopener noreferrer"
            className="font-semibold text-[var(--primary)] hover:underline"
          >
            {innerLinkMatch[1]}
          </a>
        );
      }
      if (innerLinkMatch && !renderLinks) {
        return (
          <strong key={i} className="font-semibold">
            {innerLinkMatch[1]}
          </strong>
        );
      }
      return (
        <strong key={i} className="font-semibold">
          {inner}
        </strong>
      );
    }

    // Markdown link: [text](url)
    const linkMatch = part.match(/^\[(.*?)\]\((.*?)\)$/);
    if (linkMatch) {
      if (renderLinks) {
        return (
          <a
            key={i}
            href={linkMatch[2]}
            target="_blank"
            rel="noopener noreferrer"
            className="text-[var(--primary)] hover:underline"
          >
            {linkMatch[1]}
          </a>
        );
      }
      return <span key={i}>{linkMatch[1]}</span>;
    }

    return <span key={i}>{part}</span>;
  });
}

/**
 * Strips markdown syntax and renders as plain text.
 * Use inside <a> or other elements where nested <a> is invalid.
 */
export function PlainMarkdownText({
  text,
  className,
}: {
  text: string;
  className?: string;
}) {
  const plain = text
    .replace(/\*\*\[(.*?)\]\([^)]*\)\*\*/g, "$1")
    .replace(/\[(.*?)\]\([^)]*\)/g, "$1")
    .replace(/\*\*(.*?)\*\*/g, "$1");
  return <span className={className}>{plain}</span>;
}
