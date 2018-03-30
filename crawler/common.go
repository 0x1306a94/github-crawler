package crawler

const (
	Host		 			= "https://github.com"
	TrendingRepoPath 		= "https://github.com/trending"
	TrendingDevelopersPath	= "https://github.com/trending/developers"
)

const (
	AllLanguage = "All Language"
)

type TaskType string
const (
	TaskTypeLanguage 	= "language"
	TaskTypeRepo		= "repo"
	TaskTypeDeveloper	= "developer"
)

type ResultType string
const (
	ResultTypeLanguage 	= "language"
	ResultTypeRepo		= "repo"
	ResultTypeDeveloper	= "developer"
)

const (
	SinceToDay	= "daily"
	SinceWeek	= "weekly"
	SinceMonth	= "monthly"
)