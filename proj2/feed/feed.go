package feed

import (
	"proj2/lock"
)

//Feed represents a user's twitter feed
// You will add to this interface the implementations as you complete them.
type Feed interface {
	Add(body string, timestamp float64)
	Remove(timestamp float64) bool
	Contains(timestamp float64) bool
	GetPosts() []*post
}

//feed is the internal representation of a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation. You can assume the feed will not have duplicate posts
type feed struct {
	start *post // a pointer to the beginning post
	rwLock *lock.RWLock
}

//post is the internal representation of a post on a user's twitter feed (hidden from outside packages)
// You CAN add to this structure but you cannot remove any of the original fields. You must use
// the original fields in your implementation.
type post struct {
	body      string // the text of the post
	timestamp float64  // Unix timestamp of the post
	next      *post  // the next post in the feed
}

//NewPost creates and returns a new post value given its body and timestamp
func newPost(body string, timestamp float64, next *post) *post {
	return &post{body, timestamp, next}
}

//NewFeed creates a empy user feed
func NewFeed() Feed {
	return &feed{start: nil, rwLock: lock.NewRWLock()}
}

// Add inserts a new post to the feed. The feed is always ordered by the timestamp where
// the most recent timestamp is at the beginning of the feed followed by the second most
// recent timestamp, etc. You may need to insert a new post somewhere in the feed because
// the given timestamp may not be the most recent.
func (f *feed) Add(body string, timestamp float64) {
	f.rwLock.Lock() // prevent multiple writing at the same time
	defer f.rwLock.Unlock()

	newPost := newPost(body, timestamp, nil)
	if f.start == nil {
		f.start = newPost
		return
	}

	if f.start.timestamp < timestamp {
		newPost.next = f.start
		f.start = newPost
		return
	}

	curr := f.start
	for curr.next != nil && curr.next.timestamp > timestamp {
		curr = curr.next
	}

	newPost.next = curr.next
	curr.next = newPost
}




// Remove deletes the post with the given timestamp. If the timestamp
// is not included in a post of the feed then the feed remains
// unchanged. Return true if the deletion was a success, otherwise return false
func (f *feed) Remove(timestamp float64) bool {
	f.rwLock.Lock() // prevent multiple writing at the same time
	defer f.rwLock.Unlock()

	if f.start == nil {
		return false
	}

	if f.start.timestamp == timestamp {
		f.start = f.start.next
		return true
	}

	curr := f.start
	for curr.next != nil && curr.next.timestamp != timestamp {
		curr = curr.next
	}

	if curr.next == nil {
		return false
	}

	curr.next = curr.next.next
	return true
}


// Contains determines whether a post with the given timestamp is
// inside a feed. The function returns true if there is a post
// with the timestamp, otherwise, false.
func (f *feed) Contains(timestamp float64) bool {
	f.rwLock.RLock()  // allow multi reading
	defer f.rwLock.RUnlock()

	curr := f.start
	for curr != nil {
		if curr.timestamp == timestamp {
			return true
		}
		curr = curr.next
	}

	return false
}


func (f *feed) GetPosts() []*post {
	f.rwLock.RLock() // using rlock to prevent writing when reading
	defer f.rwLock.RUnlock()

	var posts []*post
	curr := f.start
	for curr != nil {
		posts = append(posts, curr)
		curr = curr.next
	}

	return posts
}

func (p *post) GetBody() string {
    return p.body
}

func (p *post) GetTimestamp() float64 {
    return p.timestamp
}


