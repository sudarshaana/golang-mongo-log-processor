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


## Demo Log format
```
{
  "timestamp": "2024-10-18T12:11:50.535353+00:00",
  "method": "GET",
  "path": "/home/",
  "full_path": "/home/?order_by=-priority,-created_at",
  "GET": {
    "order_by": "-priority,-created_at"
  },
  "user": "sudarshaana@gmail.com",
  "status_code": 200,
  "duration": 0.059896,
  "request_body": {},
  "headers": {
    "Content-Length": "",
    "Content-Type": "text/plain",
    "Host": "127.0.0.1:8000",
    "Connection": "keep-alive",
    "Sec-Ch-Ua-Platform": "\"macOS\"",
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
    "Accept": "application/json, text/plain, */*",
    "Sec-Ch-Ua": "\"Google Chrome\";v=\"129\", \"Not=A?Brand\";v=\"8\", \"Chromium\";v=\"129\"",
    "Dnt": "1",
    "Sec-Ch-Ua-Mobile": "?0",
    "Origin": "http://localhost:3000",
    "Sec-Fetch-Site": "cross-site",
    "Sec-Fetch-Mode": "cors",
    "Sec-Fetch-Dest": "empty",
    "Referer": "http://localhost:3000/",
    "Accept-Encoding": "gzip, deflate, br, zstd",
    "Accept-Language": "en-BD,en-GB;q=0.9,en-US;q=0.8,en;q=0.7,bn;q=0.6"
  }
}```
