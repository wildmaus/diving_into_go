# Batch
The task was to implement a client to optimize interaction with the service.   
The service processes objects in batches (the maximum size of a batch is _n_). It is known that the service processes them for a certain time _p_ and is blocked for a long time if another batch arrives at this time.     
In implementation use mutex for accessing queue and last process time to avoid rewriting problems. Another mutex is used for allow only one process goroutine work with service at time.
### Commands
```bash
go run app/app.go
```