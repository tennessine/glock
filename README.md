glock
=====

A distributed lock service in Go using etcd . It easy to use like sync.Mutex. 

## Usage
```
// New lock

m1 := NewMutex("key", ":10", 10, []string{"http://127.0.0.1:4001"})
m1.Lock()
m1.Unlock()
```