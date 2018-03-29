package crawler

import "fmt"

type work struct {
	id 			int
	jobQueue	chan task
	resultQueue chan interface{}
}

func NewWork() *work {
	return &work{
		jobQueue: make(chan task, 50),
		resultQueue: make(chan interface{}, 100),
	}
}

func (w *work) start()  {
	fmt.Printf("work id: %d start\n", w.id)
	for task := range w.jobQueue {
		retryCount := task.retryCount()
		des := ""
		switch task.taskType() {
		case TaskTypeLanguage:
			des = "fetch language"
		case TaskTypeRepo:
			des = "fetch repo"
		default:
			des = "fetch developer"
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
		w.resultQueue <- result
	}
	fmt.Printf("work id: %d stop\n",w.id)
}

func (w *work) stop()  {
	close(w.jobQueue)
	close(w.resultQueue)
}