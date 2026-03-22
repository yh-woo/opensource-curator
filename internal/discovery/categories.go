package discovery

// CategoryKeywords maps category slugs to npm search keywords.
// Each category has multiple keyword sets to maximize coverage.
var CategoryKeywords = map[string][]string{
	"http-client":      {"http client", "fetch request", "ajax http"},
	"orm-database":     {"orm database", "sql query builder", "database driver postgres mysql"},
	"auth":             {"authentication jwt", "oauth login", "password hash session"},
	"testing":          {"test framework", "testing assertion", "e2e test browser"},
	"cli-framework":    {"cli command line", "terminal cli framework", "cli argument parser"},
	"logging":          {"logger logging", "log transport", "structured logging"},
	"file-processing":  {"file parser csv", "pdf excel document", "image processing sharp"},
	"validation":       {"validation schema", "type check validate", "json schema validator"},
	"date-time":        {"date time format", "datetime timezone", "date manipulation"},
	"state-management": {"state management store", "reactive state", "state machine"},
	"web-framework":    {"web framework server", "http server router", "rest api framework"},
	"api-client":       {"api sdk client", "cloud sdk aws gcp", "stripe twilio sdk"},
	"caching":          {"cache redis memory", "cache lru ttl", "key value store"},
	"queue-messaging":  {"message queue", "job queue worker", "event bus pubsub"},
	"ai-ml-sdk":        {"ai sdk llm", "openai anthropic", "machine learning embedding vector"},

	// New categories to add
	"email":            {"email smtp send", "email template transactional"},
	"websocket":        {"websocket realtime", "socket.io ws"},
	"graphql":          {"graphql server client", "graphql schema"},
	"markdown-templating": {"markdown parser", "template engine handlebars"},
	"crypto-security":  {"encryption crypto", "security helmet csrf"},
	"config-env":       {"config environment dotenv", "configuration management"},
	"monitoring":       {"monitoring apm metrics", "error tracking sentry"},
	"bundler-build":    {"bundler build tool", "webpack vite esbuild rollup"},
}
