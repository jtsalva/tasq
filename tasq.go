package tasq

import (
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
	"time"
)

const (
	QTasksReadOnlyScope  = tasks.TasksReadonlyScope
	QTasksReadWriteScope = tasks.TasksScope
)

type QConfig struct {
	Credentials string
	Scope       string
}

type QService struct {
	*tasks.Service

	Tasklists *QTasklistsService
	Tasks     *QTasksService
}

func Init(cfg *QConfig) error {
	return Auth.init(cfg)
}

func NewService(tokenString []byte) (*QService, error) {
	tasqService := &QService{
		Tasklists: &QTasklistsService{},
		Tasks:     &QTasksService{},
	}

	ctx := context.Background()
	tokenSource, err := Auth.getTokenSource(ctx, tokenString)
	if err != nil {
		return tasqService, err
	}

	opt := option.WithTokenSource(tokenSource)
	tasqService.Service, err = tasks.NewService(ctx, opt)
	if err != nil {
		return tasqService, err
	}

	tasqService.Tasklists.TasklistsService = tasqService.Service.Tasklists
	tasqService.Tasks.TasksService = tasqService.Service.Tasks

	return tasqService, nil
}

func Time(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}
