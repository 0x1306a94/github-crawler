package crawler

const (
	Host		 			= "https://github.com"
	TrendingRepoPath 		= "https://github.com/trending"
	TrendingDevelopersPath	= "https://github.com/trending/developers"
)

const (
	AllLanguage = "All Language"
)

type TaskType int
const (
	TaskTypeLanguage = TaskType(iota + 1)
	TaskTypeRepo
	TaskTypeDeveloper
)

type ResultType int
const (
	ResultTypeLanguage = ResultType(iota + 1)
	ResultTypeRepo
	ResultTypeDeveloper
)

const (
	SinceToDay	= "daily"
	SinceWeek	= "weekly"
	SinceMonth	= "monthly"
)