package server

import (
	"encoding/json"
	"proj2/feed"
	"proj2/queue"
	"sync"
)

type Config struct {
    Encoder *json.Encoder // Represents the buffer to encode Responses
    Decoder *json.Decoder // Represents the buffer to decode Requests
    Mode    string        // Represents whether the server should execute
    // sequentially or in parallel
    // If Mode == "s"  then run the sequential version
    // If Mode == "p"  then run the parallel version
    // These are the only values for Version
    ConsumersCount int // Represents the number of consumers to spawn
}

// Run starts up the twitter server based on the configuration
// information provided and only returns when the server is fully
// shutdown.
func Run(config Config) {
    switch config.Mode {
    case "s":
        runSequential(config)
    case "p":
        runParallel(config)
    default:
        panic("Invalid mode")
    }
}

func runSequential(config Config) {
    f := feed.NewFeed()
    for {
        var req queue.Request
        if err := config.Decoder.Decode(&req); err != nil { // Decode will decode and store in req
            break
        }
        if req.Command == "DONE" {
            break
        }
        processRequest(f, config.Encoder, &req)
    }
}

func runParallel(config Config) {
    f := feed.NewFeed()
    taskQueue := queue.NewLockFreeQueue()
    var wg sync.WaitGroup
    mu := sync.Mutex{}
    cond := sync.NewCond(&mu)
    done := false

    for i := 0; i < config.ConsumersCount; i++ {
        wg.Add(1)
        go consumer(f, config.Encoder, taskQueue, &done, &mu, cond, &wg)
    }
    producer(config.Decoder, taskQueue, &done, &mu, cond)

    wg.Wait()
}

func producer(dec *json.Decoder, taskQueue *queue.LockFreeQueue, done *bool, mu *sync.Mutex, cond *sync.Cond) {
    for {
        var req queue.Request
        if err := dec.Decode(&req); err != nil {
            break
        }
        if req.Command == "DONE" {
            mu.Lock()
            *done = true
            cond.Broadcast()
            mu.Unlock()
            return
        }
        taskQueue.Enqueue(&req)
        cond.Signal()
    }
}

func consumer(f feed.Feed, encoder *json.Encoder, taskQueue *queue.LockFreeQueue, done *bool, mu *sync.Mutex, cond *sync.Cond, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        mu.Lock()
        for taskQueue.IsEmpty() && !*done {
            cond.Wait() 
        }
        if *done && taskQueue.IsEmpty() {
            mu.Unlock()
            return
        }
        req := taskQueue.Dequeue()
        mu.Unlock()
        if req != nil {
            processRequest(f, encoder, req)
        }
    }
}

func processRequest(f feed.Feed, encoder *json.Encoder, req *queue.Request) {
    response := map[string]interface{}{"id": req.ID}
    switch req.Command {
    case "ADD":
        f.Add(req.Body, req.Timestamp)
        response["success"] = true

    case "REMOVE":
        response["success"] = f.Remove(req.Timestamp)
    case "CONTAINS":
        response["success"] = f.Contains(req.Timestamp)
    case "FEED":
        var feedPosts []map[string]interface{}
        posts := f.GetPosts()
        for _, post := range posts {
            feedPosts = append(feedPosts, map[string]interface{}{
                "body":      post.GetBody(),
                "timestamp": post.GetTimestamp(),
            })
        }
        response["feed"] = feedPosts
    }
    encoder.Encode(response)
}