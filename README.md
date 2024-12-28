# Time scheduler

### Problem
Have a stream of tasks with specified time intervals. Don't wanna spawn a separate goroutine for each.

### Solution
Use the slice of tasks. Maintain order of entries from recent to latest. Drain entries until remained will be in the future. Be carefull, this package make calls in one goroutine, so if your routine is heavy - think about separate goroutine.

### Example
```golang
sched := NewScheduler(ctx)
sched.Add(someTimeAtFuture1, func() { fmt.Println("called")})
sched.Add(someTimeAtFuture2, func() { fmt.Println("called")})
// no need to stop, when ctx is done it will closed automatically
```