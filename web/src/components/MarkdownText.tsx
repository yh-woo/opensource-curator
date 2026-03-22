/**
 * Renders plain text with markdown-style links [text](url)
 * converted to clickable <a> tags. Also handles **bold** text.
 */
export function MarkdownText({
  text,
  className,
}: {
  text: string;
  className?: string;
}) {
  // Match [text](url) and **bold**
  const parts = text.split(/(\[.*?\]\(.*?\)|\*\*.*?\*\*)/g);

  return (
    <span className={className}>
      {parts.map((part, i) => {
        // Markdown link: [text](url)
        const linkMatch = part.match(/^\[(.*?)\]\((.*?)\)$/);
        if (linkMatch) {
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

        // Bold: **text**
        const boldMatch = part.match(/^\*\*(.*?)\*\*$/);
        if (boldMatch) {
          return (
            <strong key={i} className="font-semibold">
              {boldMatch[1]}
            </strong>
          );
        }

        return <span key={i}>{part}</span>;
      })}
    </span>
  );
}
