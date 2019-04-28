# tasQ
Built on top of [google.golang.org/api/tasks/v1](https://google.golang.org/api/tasks/v1) to add extra functionality to make it easier to get started working with the Google Tasks API in Go.

* [Getting Started](#getting-started)
* [Listing Tasklists](#listing-tasklists)
* [Listing Tasks](#listing-tasks)
* [Filter and Sort Tasks](#filter-and-sort-tasks)
* [Interacting with Tasks](#interacting-with-tasks)

## Getting Started
1. Enable Google Tasks API from [API Console](https://console.developers.google.com/)
2. Create a new OAuth Client ID credential and download it as JSON
3. Configure your OAuth consent screen
4. Understand the concepts [developers.google.com/tasks/concepts](https://developers.google.com/tasks/concepts)
5. Get tasq `go get -u github.com/jtsalva/tasq`
6. Import tasq `import "github.com/jtsalva/tasq"`

```Go
tasq.Init(&tasq.QConfig{
  // Either QTasksReadWriteScope or QTasksReadOnlyScope
  Scope:       tasq.QTasksReadWriteScope,
  Credentials: "/path/to/credentials.json",
})

// Direct users here to grant access to your
// application from their Google accounts
authURL := tasq.Auth.GetAuthCodeURL()

// Once the user grants access and is
// redirected to your specified URL, grab
// the code from the query string
authCode := "4/OPHKKL4qxLFhalqyX740oPGAAKyS79Lm3sPgFgqFQHtATBnAJW2aLARa2kABuJJhgDOciv-LAT7p4MULMaP9C1"

// The auth code can only be used once
// to generate a token, the token is
// reusable, store it somewhere safe
token, err := tasq.Auth.GetToken(authCode)


svc, err := tasq.NewService(token)

// List all Tasklists
tasklists, err := svc.Tasklists.List().Do()

// List Tasks from Tasklist
tasks, err := svc.Tasks.List(tasklistId).Do()
```

## Listing Tasklists
```Go
// tasklists is of type QTaskLists
tasklists, err := svc.Tasklists.List().Do()

// tasklists.Items is of type []*QTaskList
for _, tasklist := range tasklists.Items {
  fmt.Println(tasklist.Id, tasklist.Name)
}
```

## Listing Tasks
```Go
// tasks is of type QTasks
tasks, err := svc.Tasks.List().Do()

// tasks.Items is of type []*QTask
for _, task := range tasks.Items {
  fmt.Println(task.Id, task.Title, task.Notes)

  // List sub-tasks
  for _, child := range task.Children {
    fmt.Println("\t", child.Id, child.Title, child.Notes)
  }
}
```

## Filter and Sort Tasks
Fllter by either
* `QCompletedFilter` - show only completed tasks
* `QNeedsActionFilter` - show only tasks needing action
* `QOverdueFilter` - show only tasks needing action where the datetime now is more than the due datetime
```Go
filteredTasks, err := svc.Tasks.List().Filter(tasq.QOverdueFilter).Do()
```
Additionally, you can sort your items either by
* `QPositionSort` - sort in the way the user positioned the tasks
* `QLatestFirstSort` - newly updated tasks first
* `QOldestFirstSort` - oldest updated tasks first
```Go
sortedTasks, err := svc.Tasks.List().Sort(tasq.QPositionSort).Do()
```
You can combine filter and sort
```Go
filterdAndSortedTasks, err := svc.Tasks.List().Filter(filter).Sort(sort).Do()
```

## Interacting with Tasks
You can directly manipulate and perform actions on a `QTaskList` and `QTask`.
```Go
// tasklist is of type QTaskList
tasklist, err := svc.Tasklists.Get(tasklistid).Do()

// task is of type QTask
task, err := svc.Tasks.Get(tasklistid, taskid).Do()
```

1. [Deleting](#deleting)
2. [Inserting](#inserting)
3. [Updating](#updating)
4. [Refreshing](#refreshing)
5. [Move to Parent](#move-to-parent)
6. [Move to Previous](#move-to-previous)
7. [Move to Beginning](#move-to-beginning)
8. [Get Time of Last Update](#get-time-of-last-update)

### Deleting
```Go
// Delete a list, including the tasks and subtasks within it
err := tasklist.Delete()

// Delete a single task
err := task.Delete()
```

### Inserting
Insert a task into another list
```Go
insertedTask, err := task.Insert(anotherTasklistid)
```

### Updating
```Go
tasklist.Title = "change tasklist title"
updatedTasklist, err := tasklist.Update()

task.Title = "change task title"
updatedTask, err := task.Update()
```

### Refreshing
If there have been remote changes, update the data currently stored in memory
```Go
err := tasklist.Refresh()
err := task.Refresh()
```

### Move to Parent
Make task a subtask to the given parent task id
```Go
movedTask, err := task.MoveToParent(parentTaskid)
```

### Move to Previous
Move task after given previous task id
```Go
movedTask, err := task.MoveToPrevious(previousTaskid)
```

### Move to Beginning
Moves task to beginning of list
```Go
movedTask, err := task.MoveToBeginning()
```

### Get Time of Last Update
Returns time of last update as type `time.Time`
```Go
tasklistUpdatedTime, err := tasklist.Time()
taskUpdatedTime, err := task.Time()
```