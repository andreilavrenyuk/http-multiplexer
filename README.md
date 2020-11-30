# HTTP Multiplexer

### Run server:
```shell
$ go run main.go
```

### Request example:
```shell
$ curl --silent --location --request GET 'http://localhost:3000/urls' \
--header 'Content-Type: application/json' \
--data-raw '[
	"http://jsonplaceholder.typicode.com/posts/1",
	"http://jsonplaceholder.typicode.com/posts/2",
	"http://jsonplaceholder.typicode.com/posts/3",
	"http://jsonplaceholder.typicode.com/posts/4",
	"http://jsonplaceholder.typicode.com/posts/5",
	"http://jsonplaceholder.typicode.com/posts/6",
	"http://jsonplaceholder.typicode.com/posts/7",
	"http://jsonplaceholder.typicode.com/posts/8",
	"http://jsonplaceholder.typicode.com/posts/9",
	"http://jsonplaceholder.typicode.com/posts/10",
	"http://jsonplaceholder.typicode.com/posts/11",
	"http://jsonplaceholder.typicode.com/posts/12",
	"http://jsonplaceholder.typicode.com/posts/13",
	"http://jsonplaceholder.typicode.com/posts/14",
	"http://jsonplaceholder.typicode.com/posts/15",
	"http://jsonplaceholder.typicode.com/posts/16",
	"http://jsonplaceholder.typicode.com/posts/17",
	"http://jsonplaceholder.typicode.com/posts/18",
	"http://jsonplaceholder.typicode.com/posts/19",
	"http://jsonplaceholder.typicode.com/posts/20"
]' | jq
```