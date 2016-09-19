thinkrus
===
[logrus](https://github.com/Sirupsen/logrus) hook for [rethinkdb](https://rethinkdb.com/)

usage
---
###Installation
```
go get github.com/HeavyHorst/thinkrus
```

###Example
```go
import (
	"github.com/HeavyHorst/thinkrus"
	"github.com/Sirupsen/logrus"
)

func main() {
	log := logrus.New()
	hook, err := thinkrus.New("localhost:28015", "test", "logs", thinkrus.WithBatchInterval(5), thinkrus.WithBatchSize(500))
	if err == nil {
		log.Hooks.Add(hook)
	}
	defer hook.Close()
}
```

Details
---
###Buffer
All log messages are buffered by default to increase the performance. You need to run hook.Close() on program exit to get sure that all messages are written to the database.

### Configuration
The constructor takes three optional configuration options:

 1. WithBatchInterval(int):  flush the buffer at the given interval in seconds. (default is 5s)
 2. WithBatchSize(int): Max size of the buffer. The buffer gets flushed to the db if this threshold is reached. (default is 200)
 3. WithSession(s *gorethink.Session): A custom rethinkdb session.

### Message Field
We will insert your message into InfluxDB with the field message.
