# GoLang MongoDB log processor
A simple GoLang Script that reads logs from the Redis queue (list) and stores them in MongoDB. It can recover the uncompleted tasks and process them.

It simply fetches the tasks from the Redis queue, puts them in `another HashSet (processing queue)`, and removes them from the processing queue when `finished`.


## Before Running
- You need to create MongoDB indexing for performance-enhancing on some fields like: `timestamps`, `method`, `path`, and so on.
- Edit the .env file with your config.

## To run using Makefile
- make build
- make up
- make help # for help


## TODO
- [x] Basic Fetching and Inserting into Mongo
- [x] Concurrent Processing
- [x] Recover running (not finished tasks)
- [ ] Retry task if error occurred
- [ ] Add Web API to fetch logs
