package tasq

import (
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
	"time"
)

type QTaskList struct {
	*tasks.TaskList

	service *QTasklistsService
}

func (taskList *QTaskList) InitNewService(tokenString []byte) error {
	var err error
	taskList.service, err = newQTasklistsService(tokenString)
	return err
}

func (taskList *QTaskList) Delete() error {
	return taskList.service.Delete(taskList.Id).Do()
}

func (taskList *QTaskList) Patch() (*QTaskList, error) {
	return taskList.service.Patch(taskList.Id, taskList).Do()
}

func (taskList *QTaskList) Update() (*QTaskList, error) {
	return taskList.service.Update(taskList.Id, taskList).Do()
}

func (taskList *QTaskList) Time() (time.Time, error) {
	return Time(taskList.Updated)
}

func (taskList *QTaskList) Refresh() error {
	id := taskList.Id
	entityTag := taskList.Etag

	updated, err := taskList.service.Get(id).IfNoneMatch(entityTag).Do()
	if err != nil {
		return err
	}

	taskList.TaskList = updated.TaskList
	return nil
}

type QTaskLists struct {
	*tasks.TaskLists

	service *QTasklistsService
	Items   []*QTaskList
}

func (taskLists *QTaskLists) InitNewService(tokenString []byte) error {
	var err error
	taskLists.service, err = newQTasklistsService(tokenString)
	return err
}

func (taskLists *QTaskLists) Time() (time.Time, error) {
	latest, _ := time.Parse(time.RFC3339, time.RFC3339)

	for _, taskList := range taskLists.Items {
		t, err := time.Parse(time.RFC3339, taskList.Updated)
		if err != nil {
			return t, err
		}
		if t.After(latest) {
			latest = t
		}
	}

	return latest, nil
}

func (taskLists *QTaskLists) Refresh() error {
	entityTag := taskLists.Etag

	updated, err := taskLists.service.List().IfNoneMatch(entityTag).Do()
	if err != nil {
		return err
	}

	taskLists.TaskLists = updated.TaskLists
	return nil
}

type QTasklistsService struct{ *tasks.TasklistsService }

func newQTasklistsService(tokenString []byte) (*QTasklistsService, error) {
	ctx := context.Background()
	tokenSource, err := Auth.getTokenSource(ctx, tokenString)
	if err != nil {
		return &QTasklistsService{}, err
	}

	opt := option.WithTokenSource(tokenSource)
	service, err := tasks.NewService(ctx, opt)
	if err != nil {
		return &QTasklistsService{}, err
	}

	return &QTasklistsService{TasklistsService: service.Tasklists}, nil
}

type QTasklistsDeleteCall struct {
	*tasks.TasklistsDeleteCall

	service *QTasklistsService
}

func (lists *QTasklistsService) Delete(tasklistid string) *QTasklistsDeleteCall {
	return &QTasklistsDeleteCall{
		TasklistsDeleteCall: lists.TasklistsService.Delete(tasklistid),
		service:             lists,
	}
}

func (call *QTasklistsDeleteCall) Context(ctx context.Context) *QTasklistsDeleteCall {
	call.TasklistsDeleteCall.Context(ctx)
	return call
}

func (call *QTasklistsDeleteCall) Do(opts ...googleapi.CallOption) error {
	return call.TasklistsDeleteCall.Do(opts...)
}

func (call *QTasklistsDeleteCall) Fields(s ...googleapi.Field) *QTasklistsDeleteCall {
	call.Fields(s...)
	return call
}

type QTasklistsGetCall struct {
	*tasks.TasklistsGetCall

	service *QTasklistsService
}

func (lists *QTasklistsService) Get(tasklistid string) *QTasklistsGetCall {
	return &QTasklistsGetCall{
		TasklistsGetCall: lists.TasklistsService.Get(tasklistid),
		service:          lists,
	}
}

func (call *QTasklistsGetCall) Context(ctx context.Context) *QTasklistsGetCall {
	call.TasklistsGetCall.Context(ctx)
	return call
}

func (call *QTasklistsGetCall) Do(opts ...googleapi.CallOption) (*QTaskList, error) {
	result, err := call.TasklistsGetCall.Do(opts...)
	taskList := &QTaskList{
		TaskList: result,
		service:  call.service,
	}

	return taskList, err
}

func (call *QTasklistsGetCall) Fields(s ...googleapi.Field) *QTasklistsGetCall {
	call.TasklistsGetCall.Fields(s...)
	return call
}

func (call *QTasklistsGetCall) IfNoneMatch(entityTag string) *QTasklistsGetCall {
	call.TasklistsGetCall.IfNoneMatch(entityTag)
	return call
}

type QTasklistsInsertCall struct {
	*tasks.TasklistsInsertCall

	service *QTasklistsService
}

func (lists *QTasklistsService) Insert(tasklist *QTaskList) *QTasklistsInsertCall {
	return &QTasklistsInsertCall{TasklistsInsertCall: lists.TasklistsService.Insert(tasklist.TaskList)}
}

func (call *QTasklistsInsertCall) Context(ctx context.Context) *QTasklistsInsertCall {
	call.TasklistsInsertCall.Context(ctx)
	return call
}

func (call *QTasklistsInsertCall) Do(opts ...googleapi.CallOption) (*QTaskList, error) {
	result, err := call.TasklistsInsertCall.Do(opts...)
	taskList := &QTaskList{
		TaskList: result,
		service:  call.service,
	}

	return taskList, err
}

func (call *QTasklistsInsertCall) Fields(s ...googleapi.Field) *QTasklistsInsertCall {
	call.TasklistsInsertCall.Fields(s...)
	return call
}

type QTasklistsListCall struct {
	*tasks.TasklistsListCall

	service *QTasklistsService
}

func (lists *QTasklistsService) List() *QTasklistsListCall {
	return &QTasklistsListCall{
		TasklistsListCall: lists.TasklistsService.List(),
		service:           lists,
	}
}

func (call *QTasklistsListCall) Context(ctx context.Context) *QTasklistsListCall {
	call.TasklistsListCall.Context(ctx)
	return call
}

func (call *QTasklistsListCall) Do(opts ...googleapi.CallOption) (*QTaskLists, error) {
	result, err := call.TasklistsListCall.Do(opts...)
	if err != nil {
		return &QTaskLists{}, err
	}

	items := make([]*QTaskList, 0)
	for _, item := range result.Items {
		items = append(items, &QTaskList{
			TaskList: item,
			service:  call.service,
		})
	}

	return &QTaskLists{
		TaskLists: result,
		Items:     items,
		service:   call.service}, err
}

func (call *QTasklistsListCall) Fields(s ...googleapi.Field) *QTasklistsListCall {
	call.TasklistsListCall.Fields(s...)
	return call
}

func (call *QTasklistsListCall) IfNoneMatch(entityTag string) *QTasklistsListCall {
	call.TasklistsListCall.IfNoneMatch(entityTag)
	return call
}

func (call *QTasklistsListCall) MaxResults(maxResults int64) *QTasklistsListCall {
	call.TasklistsListCall.MaxResults(maxResults)
	return call
}

func (call *QTasklistsListCall) PageToken(pageToken string) *QTasklistsListCall {
	call.TasklistsListCall.PageToken(pageToken)
	return call
}

type QTasklistsPatchCall struct {
	*tasks.TasklistsPatchCall

	service *QTasklistsService
}

func (lists *QTasklistsService) Patch(tasklistid string, tasklist *QTaskList) *QTasklistsPatchCall {
	return &QTasklistsPatchCall{
		TasklistsPatchCall: lists.TasklistsService.Patch(tasklistid, tasklist.TaskList),
		service:            lists,
	}
}

func (call *QTasklistsPatchCall) Context(ctx context.Context) *QTasklistsPatchCall {
	call.TasklistsPatchCall.Context(ctx)
	return call
}

func (call *QTasklistsPatchCall) Do(opts ...googleapi.CallOption) (*QTaskList, error) {
	result, err := call.TasklistsPatchCall.Do(opts...)
	taskList := &QTaskList{
		TaskList: result,
		service:  call.service,
	}

	return taskList, err
}

func (call *QTasklistsPatchCall) Fields(s ...googleapi.Field) *QTasklistsPatchCall {
	call.TasklistsPatchCall.Fields(s...)
	return call
}

type QTasklistsUpdateCall struct {
	*tasks.TasklistsUpdateCall

	service *QTasklistsService
}

func (lists *QTasklistsService) Update(taskslistid string, tasklist *QTaskList) *QTasklistsUpdateCall {
	return &QTasklistsUpdateCall{TasklistsUpdateCall: lists.TasklistsService.Update(taskslistid, tasklist.TaskList)}
}

func (call *QTasklistsUpdateCall) Context(ctx context.Context) *QTasklistsUpdateCall {
	call.TasklistsUpdateCall.Context(ctx)
	return call
}

func (call *QTasklistsUpdateCall) Do(opts ...googleapi.CallOption) (*QTaskList, error) {
	result, err := call.TasklistsUpdateCall.Do(opts...)
	taskList := &QTaskList{
		TaskList: result,
		service:  call.service,
	}

	return taskList, err
}

func (call *QTasklistsUpdateCall) Fields(s ...googleapi.Field) *QTasklistsUpdateCall {
	call.TasklistsUpdateCall.Fields(s...)
	return call
}
