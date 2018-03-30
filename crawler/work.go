package crawler

import "fmt"

type work struct {
	id 			int
	jobQueue	chan task
	resultQueue chan WorkResult
}

func NewWork() *work {
	return &work{
		jobQueue: make(chan task, 50),
		resultQueue: make(chan WorkResult, 100),
	}
}

func (w *work) start()  {
	fmt.Printf("work id: %d start\n", w.id)
	for task := range w.jobQueue {
		retryCount := task.retryCount()
		des := ""
		var resultType ResultType
		switch task.taskType() {
		case TaskTypeLanguage:
			des = "fetch language"
			resultType = ResultTypeLanguage
		case TaskTypeRepo:
			des = "fetch repo"
			resultType = ResultTypeRepo
		default:
			des = "fetch developer"
			resultType = ResultTypeDeveloper
		}

		if retryCount <= 0 {
			retryCount = 1
		}
		var result interface{} = nil
		for i := 0; i < retryCount; i++ {
			fmt.Printf("work id: %d exec: %v\n", w.id, des)
			r := task.call()
			switch task.taskType() {
			case TaskTypeLanguage:
				if r != nil && len(r.([]string)) > 0 {
					result = r
				}
			default:
				if r != nil {
					result = r
				}
			}
			if result != nil {
				break
			}
		}
		r := WorkResult{
			ResultType: resultType,
			Since: task.sinceDesc(),
			Language: task.languageDesc(),
			Result: result,
		}
		w.resultQueue <- r
	}
	fmt.Printf("work id: %d stop\n",w.id)
}

func (w *work) stop()  {
	close(w.jobQueue)
	close(w.resultQueue)
}