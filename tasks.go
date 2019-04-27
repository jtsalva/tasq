package tasq

import (
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
	"time"
)

type QTaskCallContext struct {
	service    *QTasksService
	tasklistid string
}

type QTask struct {
	*tasks.Task

	ctx      *QTaskCallContext
	Children []*QTask
}

func (task *QTask) InitNewService(tokenString []byte) error {
	var err error
	task.ctx.service, err = newQTasksService(tokenString)
	return err
}

func (task *QTask) Delete() error {
	return task.ctx.service.Delete(task.ctx.tasklistid, task.Id).Do()
}

func (task *QTask) Insert(tasklistid string) (*QTask, error) {
	return task.ctx.service.Insert(task.ctx.tasklistid, task).Do()
}

func (task *QTask) MoveToParent(parent string) (*QTask, error) {
	return task.ctx.service.Move(task.ctx.tasklistid, task.Id).Parent(parent).Do()
}

func (task *QTask) MoveToPrevious(previous string) (*QTask, error) {
	return task.ctx.service.Move(task.ctx.tasklistid, task.Id).Previous(previous).Do()
}

func (task *QTask) MoveToBeginning() (*QTask, error) {
	return task.ctx.service.Move(task.ctx.tasklistid, task.Id).Do()
}

func (task *QTask) Patch() (*QTask, error) {
	return task.ctx.service.Patch(task.ctx.tasklistid, task.Id, task).Do()
}

func (task *QTask) Update() (*QTask, error) {
	return task.ctx.service.Update(task.ctx.tasklistid, task.Id, task).Do()
}

func (task *QTask) Time() (time.Time, error) {
	return Time(task.Updated)
}

func (task *QTask) Refresh() error {
	id := task.Id
	entityTag := task.Etag

	updated, err := task.ctx.service.Get(task.ctx.tasklistid, id).IfNoneMatch(entityTag).Do()
	if err != nil {
		return err
	}

	task.Task = updated.Task
	return nil
}

type QTasks struct {
	*tasks.Tasks

	ctx   *QTaskCallContext
	Items []*QTask
}

func (tasks *QTasks) InitNewService(tokenString []byte) error {
	var err error
	tasks.ctx.service, err = newQTasksService(tokenString)
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

func (tasks *QTasks) Refresh() error {
	entityTag := tasks.Etag

	updated, err := tasks.ctx.service.List(tasks.ctx.tasklistid).IfNoneMatch(entityTag).Do()
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

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Clear(tasklistid string) *QTasksClearCall {
	return &QTasksClearCall{
		TasksClearCall: tasks.TasksService.Clear(tasklistid),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
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

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Delete(tasklistid string, taskid string) *QTasksDeleteCall {
	return &QTasksDeleteCall{
		TasksDeleteCall: tasks.TasksService.Delete(tasklistid, taskid),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
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

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Get(tasklistid string, taskid string) *QTasksGetCall {
	return &QTasksGetCall{
		TasksGetCall: tasks.TasksService.Get(tasklistid, taskid),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
	}
}

func (call *QTasksGetCall) Context(ctx context.Context) *QTasksGetCall {
	call.TasksGetCall.Context(ctx)
	return call
}

func (call *QTasksGetCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksGetCall.Do(opts...)
	return &QTask{
		Task: result,
		ctx:  call.ctx}, err
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

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Insert(tasklistid string, task *QTask) *QTasksInsertCall {
	return &QTasksInsertCall{
		TasksInsertCall: tasks.TasksService.Insert(tasklistid, task.Task),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
	}
}

func (call *QTasksInsertCall) Context(ctx context.Context) *QTasksInsertCall {
	call.TasksInsertCall.Context(ctx)
	return call
}

func (call *QTasksInsertCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksInsertCall.Do(opts...)
	return &QTask{
		Task: result,
		ctx:  call.ctx}, err
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

	ctx    *QTaskCallContext
	filter string
	sort   string
}

func (tasks *QTasksService) List(tasklistid string) *QTasksListCall {
	return &QTasksListCall{
		TasksListCall: tasks.TasksService.List(tasklistid),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
	}
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

// TODO: Filter for deleted tasks
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
		Tasks: result,
		ctx:   call.ctx,
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
		items = append(items, &QTask{
			Task: item,
			ctx:  call.ctx,
		})
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

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Move(tasklistid string, taskid string) *QTasksMoveCall {
	return &QTasksMoveCall{
		TasksMoveCall: tasks.TasksService.Move(tasklistid, taskid),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
	}
}

func (call *QTasksMoveCall) Context(ctx context.Context) *QTasksMoveCall {
	call.TasksMoveCall.Context(ctx)
	return call
}

func (call *QTasksMoveCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksMoveCall.Do(opts...)
	return &QTask{
		Task: result,
		ctx:  call.ctx}, err
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

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Patch(tasklistid string, taskid string, task *QTask) *QTasksPatchCall {
	return &QTasksPatchCall{
		TasksPatchCall: tasks.TasksService.Patch(tasklistid, taskid, task.Task),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
	}
}

func (call *QTasksPatchCall) Context(ctx context.Context) *QTasksPatchCall {
	call.TasksPatchCall.Context(ctx)
	return call
}

func (call *QTasksPatchCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksPatchCall.Do(opts...)
	return &QTask{
		Task: result,
		ctx:  call.ctx}, err
}

func (call *QTasksPatchCall) Fields(s ...googleapi.Field) *QTasksPatchCall {
	call.TasksPatchCall.Fields(s...)
	return call
}

type QTasksUpdateCall struct {
	*tasks.TasksUpdateCall

	ctx *QTaskCallContext
}

func (tasks *QTasksService) Update(tasklistid string, taskid string, task *QTask) *QTasksUpdateCall {
	return &QTasksUpdateCall{
		TasksUpdateCall: tasks.TasksService.Update(tasklistid, taskid, task.Task),
		ctx: &QTaskCallContext{
			service:    tasks,
			tasklistid: tasklistid,
		},
	}
}

func (call *QTasksUpdateCall) Context(ctx context.Context) *QTasksUpdateCall {
	call.TasksUpdateCall.Context(ctx)
	return call
}

func (call *QTasksUpdateCall) Do(opts ...googleapi.CallOption) (*QTask, error) {
	result, err := call.TasksUpdateCall.Do(opts...)
	return &QTask{
		Task: result,
		ctx:  call.ctx}, err
}

func (call *QTasksUpdateCall) Fields(s ...googleapi.Field) *QTasksUpdateCall {
	call.TasksUpdateCall.Fields(s...)
	return call
}
