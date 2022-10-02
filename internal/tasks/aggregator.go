package tasks

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/iagapie/rumors/internal/consts"
)

type Aggregator struct {
	Log asynq.Logger
}

func (a *Aggregator) Aggregate(group string, tasks []*asynq.Task) *asynq.Task {
	a.Log.Debug(fmt.Sprintf("aggregator group %s", group))

	result := []byte{byte(91)} // [
	for i, task := range tasks {
		if i != 0 {
			result = append(result, byte(44)) // ,
		}
		result = append(result, task.Payload()...)
	} // ,
	result = append(result, byte(93)) // ]

	return asynq.NewTask(consts.TaskFeedItemAggregated, result)
}
