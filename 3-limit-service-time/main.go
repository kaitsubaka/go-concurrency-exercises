//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

func (u *User) addSecond() {
	atomic.AddInt64(&u.TimeUsed, 1)
}

const freeQuota = 10

func (u *User) quotaExceeded() bool {
	return !u.IsPremium && u.TimeUsed >= freeQuota
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	ticker := time.Tick(time.Second)
	done := make(chan struct{})

	go func() {
		process()
		done <- struct{}{}
	}()

	for {
		select {
		case <-done:
			return true
		case <-ticker:
			u.addSecond()
			if u.quotaExceeded() {
				// process killed
				return false
			}
		}
	}
}

func main() {
	RunMockServer()
}
