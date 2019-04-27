# tasQ
Built on top of [google.golang.org/api/tasks/v1](https://google.golang.org/api/tasks/v1) to add extra functionality to make it easier to get started working with the Google Tasks API in Go.

## Getting Started
1. Enable Google Tasks API from [API Console](https://console.developers.google.com/)
2. Create a new OAuth Client ID credential and download it as JSON
3. Configure your OAuth consent screen
4. Understand the concepts [developers.google.com/tasks/concepts](https://developers.google.com/tasks/concepts) 
5. Get tasq `go get -u github.com/jtsalva/tasq`

```Go
tasq.Init(&tasq.QConfig{
  Credentials: "/path/to/credentials.json",
  
  // Either QTasksReadWriteScope or QTasksReadOnlyScope
  Scope:       tasq.QTasksReadWriteScope,
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

## Iterating Tasklists
```Go
tasklists, err := svc.Tasklists.List().Do()

for _, tasklist := range tasklists.Items {
  fmt.Println(tasklist.Id, tasklist.Name)
}
```

## Iterating Tasks
```Go
tasks, err := svc.Tasks.List().Do()

for _, task := range tasks.Items {
  fmt.Println(task.Id, task.Title, task.Notes)
  
  // Iterate sub-tasks
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
