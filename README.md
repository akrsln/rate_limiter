# Rate limiter with token bucket algorithm 
```go

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"net/http"
	"time"
)

func exampleUsage() {
	server := http.NewServeMux()
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	bucket := NewBucket(5, 1*time.Minute, redisClient)
	server.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		if !bucket.Allow("userName") {
			w.WriteHeader(429)
			fmt.Fprintf(w, "Too many request!")
			fmt.Fprintf(w, "Reached max allowed request count. Try again after %f seconds.. ", bucket.NextAllowed("userName").Seconds())
			return
		}
		fmt.Fprintf(w, "Remaining request count: %d \n", bucket.Remaining("userName"))
		return
	})
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		return
	}
}

```
