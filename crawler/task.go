package crawler

type task interface {
	taskType() TaskType
	call() interface{}
	retryCount() int
}


type commonTask struct {
	language	string
	since		string
	retry		int
}

type repoTask struct {
	commonTask
}

type developerTask struct {
	commonTask
}

type languageTask struct {
	retry		int
}

func (_ *repoTask) taskType() TaskType {
	return TaskTypeRepo
}

func (r *repoTask) call() (interface{}) {
	return parserFromRepos(r.language, r.since)
}

func (r *repoTask) retryCount() int  {
	return r.retry
}

func (_ *developerTask) taskType() TaskType {
	return TaskTypeDeveloper
}

func (r *developerTask) call() (interface{}) {
	return  parserFromDevelopers(r.language, r.since)
}

func (r *developerTask) retryCount() int  {
	return r.retry
}

func (_ *languageTask) taskType() TaskType {
	return TaskTypeLanguage
}

func (_ *languageTask) call() (interface{}) {
	return parserFromAllLanguage()
}

func (l *languageTask) retryCount() int  {
	return l.retry
}