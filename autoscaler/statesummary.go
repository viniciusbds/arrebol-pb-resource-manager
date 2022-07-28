package autoscaler

type QueueState struct {
	QueueID           string `json:"QueueID"`
	NumWorkers        int    `json:"NumWorkers"`
	NumHealthyWorkers int    `json:"NumHealthyWorkers"`
	NumReadToRunTasks int    `json:"NumReadToRunTasks"`
}
