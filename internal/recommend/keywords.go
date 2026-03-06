package recommend

// categoryKeywords maps keywords/phrases to category slugs.
var categoryKeywords = map[string][]string{
	"http-client":      {"http", "request", "fetch", "api call", "rest client", "http client"},
	"orm-database":     {"database", "orm", "sql", "query builder", "db", "postgres", "mysql", "sqlite", "mongodb", "redis"},
	"auth":             {"auth", "authentication", "authorization", "login", "jwt", "oauth", "session", "password", "bcrypt"},
	"testing":          {"test", "testing", "unit test", "e2e", "integration test", "mock", "assert", "spec"},
	"cli-framework":    {"cli", "command line", "terminal", "argv", "args", "prompt"},
	"logging":          {"log", "logging", "logger", "debug"},
	"file-processing":  {"file", "csv", "excel", "pdf", "image", "upload", "parse file", "document"},
	"validation":       {"validation", "validate", "schema", "type check", "sanitize", "form validation"},
	"date-time":        {"date", "time", "datetime", "timestamp", "calendar", "timezone"},
	"state-management": {"state", "store", "state management", "global state", "reactive"},
	"web-framework":    {"web framework", "server", "web server", "routing", "middleware", "express", "api server"},
	"api-client":       {"sdk", "api client", "cloud", "aws", "stripe", "github api", "openai"},
	"caching":          {"cache", "caching", "redis", "memcache", "lru", "storage"},
	"queue-messaging":  {"queue", "message", "messaging", "pubsub", "event", "job", "worker", "task queue"},
	"ai-ml-sdk":        {"ai", "ml", "llm", "gpt", "claude", "openai", "langchain", "embedding", "vector", "rag"},
}

// preferSortFields maps prefer parameter values to sort fields.
var preferSortFields = map[string]string{
	"lightweight": "dependencies_count",
	"stable":      "maintenance_health",
	"secure":      "security_posture",
	"popular":     "community_signal",
}
