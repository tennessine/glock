glock
=====

A distributed lock service in Go using etcd . It easy to use like sync.Mutex. 

## Usage
```
// New lock
m1 := NewMutex("id", "object id", 10, []string{"http://127.0.0.1:4001"})
// Block
m1.Lock()
// Released
m1.Unlock()
```
So eazy.
