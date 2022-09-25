package utils

var (
	Handlers            map[string]Handler
	CachedUsers         map[int64]CacheUser
	ChatLogIsLoaded     map[int64]bool
	ChatLogMessageCache map[int64]map[int64]*SomePlaceholder
)
