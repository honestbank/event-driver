package storage

import (
	"github.com/KeluDiao/event-driver/event"
)

// TODO: implement me
// requirement:
// 1. key & source should be compose key
// 2. should use connection pool for scalability
// 3. try to make the DB care-free for users
type EventStore interface {
	Persist(key, source, content string) error
	LookUp(key, source string) (event.Message, error)
	LookUpByKey(key string) ([]event.Message, error)
}
