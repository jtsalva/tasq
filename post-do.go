package tasq

import (
	"time"
)

const (
	QCompletedFilter   = "filter.completed"
	QNeedsActionFilter = "filter.needs_action"
	QOverdueFilter     = "filter.overdue"

	QPositionSort    = "sort.position"
	QLatestFirstSort = "sort.latest_first"
	QOldestFirstSort = "sort.oldest_first"
)

func raiseTasks(list []*QTask) []*QTask {
	raisedTasks := make([]*QTask, 0)
	positionalSort(list)

	topLevelTasks := make(map[string]int)
	for _, task := range list {
		if task.Parent == "" {
			topLevelTasks[task.Id] = len(raisedTasks)
			raisedTasks = append(raisedTasks, task)
		} else {
			parentIdx := topLevelTasks[task.Parent]
			raisedTasks[parentIdx].Children = append(raisedTasks[parentIdx].Children, task)
		}
	}

	return raisedTasks
}

func statusFilter(list *QTasks, status string) {
	matchingTasks := make([]*QTask, 0)

	for _, task := range list.Items {
		if task.Status == status {
			matchingTasks = append(matchingTasks, task)
		}
	}

	list.Items = matchingTasks
}

func positionalSort(list []*QTask) {
	length := len(list)

	for i := 1; i < length; i++ {
		j := i
		for j > 0 && list[j].Position < list[j-1].Position {
			list[j], list[j-1] = list[j-1], list[j]
			j -= 1
		}
	}
}

func chronologicalSort(list []*QTask) {
	length := len(list)

	for i := 1; i < length; i++ {
		j := i

		tj0, _ := time.Parse(time.RFC3339, list[j-1].Updated)
		tj1, _ := time.Parse(time.RFC3339, list[j].Updated)
		for j > 0 && tj1.After(tj0) {
			list[j], list[j-1] = list[j-1], list[j]
			j -= 1
		}
	}
}

func reverseChronologicalSort(list []*QTask) {
	chronologicalSort(list)

	i := 0
	j := len(list) - 1
	for i < j {
		list[i], list[j] = list[j], list[i]
		i++
		j--
	}
}
