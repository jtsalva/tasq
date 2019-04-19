package tasq

import (
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
	"time"
)

type QTask struct {
	*tasks.Task

	service  *QTasksService
	Children []*QTask
}

func (task *QTask) InitNewService(tokenString []byte) error {
	var err error
	task.service, err = newQTasksService(tokenString)
	return err
}

func (task *QTask) Delete(tasklistid string) error {
	return task.service.Delete(tasklistid, task.Id).Do()
}

func (task *QTask) Insert(tasklistid string) (*QTask, error) {
	return task.service.Insert(tasklistid, task).Do()
}

func (task *QTask) MoveToParent(tasklistid string, parent string) (*QTask, error) {
	return task.service.Move(tasklistid, task.Id).Parent(parent).Do()
}

func (task *QTask) MoveToPrevious(tasklistid string, previous string) (*QTask, error) {
	return task.service.Move(tasklistid, task.Id).Previous(previous).Do()
}

func (task *QTask) MoveToBeginning(tasklistid string) (*QTask, error) {
	return task.service.Move(tasklistid, task.Id).Do()
}

func (task *QTask) Patch(tasklistid string) (*QTask, error) {
	return task.service.Patch(tasklistid, task.Id, task).Do()
}

func (task *QTask) Update(tasklistid string) (*QTask, error) {
	return task.service.Update(tasklistid, task.Id, task).Do()
}

func (task *QTask) Time() (time.Time, error) {
	return Time(task.Updated)
}

func (task *QTask) Refresh(tasklistid string) error {
	id := task.Id
	entityTag := task.Etag

	updated, err := task.service.Get(tasklistid, id).IfNoneMatch(entityTag).Do()
	if err != nil {
		return err
	}

	task.Task = updated.Task
	return nil
}

type QTasks struct {
	*tasks.Tasks

	service *QTasksService
	Items   []*QTask
}

func (tasks *QTasks) InitNewService(tokenString []byte) error {
	var err error
	tasks.service, err = newQTasksService(tokenString)
	return err
}

func (tasks *QTasks) Time() (time.Time, error) {
	latest, _ := time.Parse(time.RFC3339, time.RFC3339)

	for _, task := range tasks.Items {
		t, err := time.Parse(time.RFC3339, task.Updated)
		if err != nil {
			return t, err
		}
		if t.After(latest) {
			latest = t
		}
	}

	return latest, nil
}

func (tasks *QTasks) Refresh(tasklistid string) error {
	entityTag := tasks.Etag

	updated, err := tasks.service.List(tasklistid).IfNoneMatch(entityTag).Do()
	if err != nil {
		return err
	}

	tasks.Tasks = updated.Tasks
	return nil
}

type QTasksService struct{ *tasks.TasksService }

func newQTasksService(tokenString []byte) (*QTasksService, error) {
	ctx := context.Background()
	tokenSource, err := Auth.getTokenSource(ctx, tokenString)
	if err != nil {
		return &QTasksService{}, err
	}

	opt := option.WithTokenSource(tokenSource)
	service, err := tasks.NewService(ctx, opt)
	if err != nil {
		return &QTasksService{}, err
	}

	return &QTasksService{TasksService: service.Tasks}, nil
}

type QTasksClearCall struct {
	*tasks.TasksClearCall

	service *QTasksService
}

func (tasks *QTasksService) Clear(tasklistid string) *QTasksClearCall {
	return &QTasksClearCall{
		TasksClearCall: tasks.TasksService.Clear(tasklistid),
		service:        tasks,
	}
}

func (call *QTasksClearCall) Context(ctx context.Context) *QTasksClearCall {
	call.TasksClearCall.Context(ctx)
	return call
}

func (call *QTasksClearCall) Do(opts ...googleapi.CallOption) error {
	return call.TasksClearCall.Do(opts...)
}

func (call *QTasksClearCall) Fields(s ...googleapi.Field) *QTasksClearCall {
	call.TasksClearCall.Fields(s...)
	return call
}

type QTasksDeleteCall struct {
	*tasks.TasksDeleteCall

	service *QTasksService
}

func (tasks *QTasksService) Delete(tasklistid string, taskid string) *QTasksDeleteCall {
	return &QTasksDeleteCall{
		TasksDeleteCall: tasks.TasksService.Delete(tasklistid, taskid),
		service:         tasks,
	}
}

func (call *QTasksDeleteCall) Context(ctx context.Context) *QTasksDeleteCall {
	call.TasksDeleteCall.Context(ctx)
	return call
}

func (call *QTasksDeleteCall) Do(opts ...googleapi.CallOption) error {
	return call.TasksDeleteCall.Do(opts...)
}

func (call *QTasksDeleteCall) Fields(s ...googleapi.Field) *QTasksDeleteCall {
	call.TasksDeleteCall.Fields(s...)
	return call
}

type QTasksGetCall struct {
	*tasks.TasksGetCall

	service *QTasksService
}

func (tasks *QTasksService) Get(tasklistid string, taskid string) *QTasksGetCall {
	return &QTasksGetCall{
		TasksGetCall: tasks.TasksService.Get(tasklistid, taskid),
		service:      tasks,
	}
}

func (call *QTasksGetCall) Context(ctx context.Context) *QTasksGetCall {
	call.TasksGetCall.Context(ctx)
	return call
}

func (call *QTasksGetCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksGetCall.Do(opts...)
	return &QTask{Task: result, service: call.service}, err
}

func (call *QTasksGetCall) Fields(s ...googleapi.Field) *QTasksGetCall {
	call.TasksGetCall.Fields(s...)
	return call
}

func (call *QTasksGetCall) IfNoneMatch(entityTag string) *QTasksGetCall {
	call.TasksGetCall.IfNoneMatch(entityTag)
	return call
}

type QTasksInsertCall struct {
	*tasks.TasksInsertCall

	service *QTasksService
}

func (tasks *QTasksService) Insert(tasklistid string, task *QTask) *QTasksInsertCall {
	return &QTasksInsertCall{TasksInsertCall: tasks.TasksService.Insert(tasklistid, task.Task)}
}

func (call *QTasksInsertCall) Context(ctx context.Context) *QTasksInsertCall {
	call.TasksInsertCall.Context(ctx)
	return call
}

func (call *QTasksInsertCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksInsertCall.Do(opts...)
	return &QTask{Task: result, service: call.service}, err
}

func (call *QTasksInsertCall) Fields(s ...googleapi.Field) *QTasksInsertCall {
	call.TasksInsertCall.Fields(s...)
	return call
}

func (call *QTasksInsertCall) Parent(parent string) *QTasksInsertCall {
	call.TasksInsertCall.Parent(parent)
	return call
}

func (call *QTasksInsertCall) Previous(previous string) *QTasksInsertCall {
	call.TasksInsertCall.Previous(previous)
	return call
}

type QTasksListCall struct {
	*tasks.TasksListCall

	service *QTasksService
	filter  string
	sort    string
}

func (tasks *QTasksService) List(tasklistid string) *QTasksListCall {
	return &QTasksListCall{TasksListCall: tasks.TasksService.List(tasklistid)}
}

func (call *QTasksListCall) CompletedMax(completedMax string) *QTasksListCall {
	call.TasksListCall.CompletedMax(completedMax)
	return call
}

func (call *QTasksListCall) CompletedMin(completedMin string) *QTasksListCall {
	call.TasksListCall.CompletedMin(completedMin)
	return call
}

func (call *QTasksListCall) Context(ctx context.Context) *QTasksListCall {
	call.TasksListCall.Context(ctx)
	return call
}

func (call *QTasksListCall) Do(opts ...googleapi.CallOption) (*QTasks, error) {
	switch call.filter {
	case QOverdueFilter:
		call.DueMax(time.Now().Format(time.RFC3339))
	case QCompletedFilter:
		call.ShowHidden(true).ShowCompleted(true)
	case QNeedsActionFilter:
		call.ShowHidden(false).ShowCompleted(false)
	}

	result, err := call.TasksListCall.Do(opts...)
	if err != nil {
		return &QTasks{}, err
	}

	list := &QTasks{
		Tasks:   result,
		service: call.service,
	}
	if len(list.Items) > 0 {
		switch call.filter {
		case QCompletedFilter:
			statusFilter(list, "completed")
		case QNeedsActionFilter:
			statusFilter(list, "needsAction")
		}

		if len(list.Items) > 1 {
			switch call.sort {
			case QPositionSort:
				positionalSort(list.Items)
			case QLatestFirstSort:
				chronologicalSort(list.Items)
			case QOldestFirstSort:
				reverseChronologicalSort(list.Items)
			}
		}
	}

	items := make([]*QTask, 0)
	for _, item := range result.Items {
		items = append(items, &QTask{Task: item})
	}
	list.Items = raiseTasks(items)

	return list, err
}

func (call *QTasksListCall) DueMax(dueMax string) *QTasksListCall {
	call.TasksListCall.DueMax(dueMax)
	return call
}

func (call *QTasksListCall) DueMin(dueMin string) *QTasksListCall {
	call.TasksListCall.DueMin(dueMin)
	return call
}

func (call *QTasksListCall) Fields(s ...googleapi.Field) *QTasksListCall {
	call.TasksListCall.Fields(s...)
	return call
}

func (call *QTasksListCall) Filter(filter string) *QTasksListCall {
	call.filter = filter
	return call
}

func (call *QTasksListCall) IfNoneMatch(entityTag string) *QTasksListCall {
	call.TasksListCall.IfNoneMatch(entityTag)
	return call
}

func (call *QTasksListCall) MaxResults(maxResults int64) *QTasksListCall {
	call.TasksListCall.MaxResults(maxResults)
	return call
}

func (call *QTasksListCall) PageToken(pageToken string) *QTasksListCall {
	call.TasksListCall.PageToken(pageToken)
	return call
}

//func (call *QTasksListCall) Raise(raise bool) *QTasksListCall {
//	call.raise = raise
//	return call
//}

func (call *QTasksListCall) ShowCompleted(showCompleted bool) *QTasksListCall {
	call.TasksListCall.ShowCompleted(showCompleted)
	return call
}

func (call *QTasksListCall) ShowDeleted(showDeleted bool) *QTasksListCall {
	call.TasksListCall.ShowDeleted(showDeleted)
	return call
}

func (call *QTasksListCall) ShowHidden(showHidden bool) *QTasksListCall {
	call.TasksListCall.ShowHidden(showHidden)
	return call
}

func (call *QTasksListCall) Sort(sort string) *QTasksListCall {
	call.sort = sort
	return call
}

func (call *QTasksListCall) UpdateMin(updatedMin string) *QTasksListCall {
	call.TasksListCall.UpdatedMin(updatedMin)
	return call
}

type QTasksMoveCall struct {
	*tasks.TasksMoveCall

	service *QTasksService
}

func (tasks *QTasksService) Move(tasklistid string, taskid string) *QTasksMoveCall {
	return &QTasksMoveCall{
		TasksMoveCall: tasks.TasksService.Move(tasklistid, taskid),
		service:       tasks,
	}
}

func (call *QTasksMoveCall) Context(ctx context.Context) *QTasksMoveCall {
	call.TasksMoveCall.Context(ctx)
	return call
}

func (call *QTasksMoveCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksMoveCall.Do(opts...)
	return &QTask{Task: result}, err
}

func (call *QTasksMoveCall) Fields(s ...googleapi.Field) *QTasksMoveCall {
	call.TasksMoveCall.Fields(s...)
	return call
}

func (call *QTasksMoveCall) Parent(parent string) *QTasksMoveCall {
	call.TasksMoveCall.Parent(parent)
	return call
}

func (call *QTasksMoveCall) Previous(previous string) *QTasksMoveCall {
	call.TasksMoveCall.Previous(previous)
	return call
}

type QTasksPatchCall struct {
	*tasks.TasksPatchCall

	service *QTasksService
}

func (tasks *QTasksService) Patch(tasklistid string, taskid string, task *QTask) *QTasksPatchCall {
	return &QTasksPatchCall{TasksPatchCall: tasks.TasksService.Patch(tasklistid, taskid, task.Task)}
}

func (call *QTasksPatchCall) Context(ctx context.Context) *QTasksPatchCall {
	call.TasksPatchCall.Context(ctx)
	return call
}

func (call *QTasksPatchCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksPatchCall.Do(opts...)
	return &QTask{Task: result}, err
}

func (call *QTasksPatchCall) Fields(s ...googleapi.Field) *QTasksPatchCall {
	call.TasksPatchCall.Fields(s...)
	return call
}

type QTasksUpdateCall struct {
	*tasks.TasksUpdateCall

	service *QTasksService
}

func (tasks *QTasksService) Update(tasklistid string, taskid string, task *QTask) *QTasksUpdateCall {
	return &QTasksUpdateCall{TasksUpdateCall: tasks.TasksService.Update(tasklistid, taskid, task.Task)}
}

func (call *QTasksUpdateCall) Context(ctx context.Context) *QTasksUpdateCall {
	call.TasksUpdateCall.Context(ctx)
	return call
}

func (call *QTasksUpdateCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksUpdateCall.Do(opts...)
	return &QTask{Task: result}, err
}

func (call *QTasksUpdateCall) Fields(s ...googleapi.Field) *QTasksUpdateCall {
	call.TasksUpdateCall.Fields(s...)
	return call
}
