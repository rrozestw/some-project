// Package gls implements goroutine-local storage.
package gls

import "sync"

// Values is simply a map of key types to value types. Used by SetValues to
// set multiple values at once.
type Values map[interface{}]interface{}

var (
	// dataLock protects access to the data map
	dataLock sync.RWMutex
	// data is a map of goroutine IDs that stores the key,value pairs
	data map[uint64]Values
)

func init() {
	data = map[uint64]Values{}
}

// With is a convenience function that stores the given values on this
// goroutine, calls the provided function (which will have access to the
// values) and then cleans up after itself.
func With(values Values, f func()) {
	SetValues(values)
	f()
	Cleanup()
}

// SetValues replaces all values for this goroutine.
func SetValues(values Values) {
	gid := curGoroutineID()
	dataLock.Lock()
	data[gid] = values
	dataLock.Unlock()
}

// Set sets the value by key and associates it with the current goroutine.
func Set(key string, value interface{}) {
	gid := curGoroutineID()
	dataLock.Lock()
	if data[gid] == nil {
		data[gid] = Values{}
	}
	data[gid][key] = value
	dataLock.Unlock()
}

// Get gets the value by key as it exists for the current goroutine.
func Get(key string) interface{} {
	gid := curGoroutineID()
	dataLock.RLock()
	if data[gid] == nil {
		dataLock.RUnlock()
		return nil
	}
	value := data[gid][key]
	dataLock.RUnlock()
	return value
}

// Go creates a new goroutine and runs the provided function in that new
// goroutine. It also associates any key,value pairs stored for the parent
// goroutine with the child goroutine. This function must be used if you wish
// to preserve the reference to any data stored in gls. This function
// automatically cleans up after itself. Do not call cleanup in the function
// passed to this function.
func Go(f func()) {
	parentData := getValues()
	go func() {
		linkGRs(parentData)
		f()
		unlinkGRs()
	}()
}

// Cleanup removes all data associated with this goroutine. If this is not
// called, the data may persist for the lifetime of your application. This
// must be called from the very first goroutine to invoke Set
func Cleanup() {
	gid := curGoroutineID()
	dataLock.Lock()
	delete(data, gid)
	dataLock.Unlock()
}

// getValues unlinks two goroutines
func getValues() Values {
	gid := curGoroutineID()
	dataLock.Lock()
	values := data[gid]
	dataLock.Unlock()
	return values
}

// linkGRs links two goroutines together, allowing the child to access the
// data present in the parent.
func linkGRs(parentData Values) {
	childID := curGoroutineID()
	dataLock.Lock()
	data[childID] = parentData
	dataLock.Unlock()
}

// unlinkGRs unlinks two goroutines
func unlinkGRs() {
	childID := curGoroutineID()
	dataLock.Lock()
	delete(data, childID)
	dataLock.Unlock()
}
