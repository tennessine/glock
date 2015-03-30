// Copyright © 2015 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package glock

import (
	"log"
	"sync"

	"github.com/coreos/go-etcd/etcd"
)

const (
	defaultTTL = 60
	defaultTry = 3

	notFound      = 100
	compareFailed = 101
	alreadyExists = 105
)

// You must call the NewMutex to get a Mutex pointer.
type Mutex struct {
	key    string
	id     string
	client *etcd.Client
	state  int32
	mutex  *sync.Mutex
	// Set to expire after a specified number of seconds, default is 60s.
	ttl uint64
}

func NewMutex(key string, id string, ttl uint64, machines []string) *Mutex {
	if ttl < 1 {
		ttl = defaultTTL
	}
	return &Mutex{
		key:    key,
		ttl:    ttl,
		client: etcd.NewClient(machines),
		mutex:  new(sync.Mutex),
	}
}

// Get a lock. When you set the ttl, if the timeout the lock will be automatically released.
func (m *Mutex) Lock() (err error) {
	m.mutex.Lock()
	for try := 1; try <= 3; try++ {
		_, err = m.client.Create(m.key, m.id, m.ttl)
		if err != nil {
			if e, ok := err.(*etcd.EtcdError); ok {
				if e.ErrorCode == alreadyExists {
				wait:
					// Get the already node's value.
					resp, err := m.client.Get(m.key, false, false)
					if err != nil {
						// Always try.
						try--
						continue
					}
					record("get ok")
					value := resp.Node.Value
					// Watch the lock node.
					record("will watch")
					_, err = m.client.Watch(m.key, 0, false, nil, nil)
					if err != nil {
						// Always try.
						try--
						continue
					}
					record("watch done.")
					// election.
					resp, err = m.client.CompareAndSwap(m.key, m.id, m.ttl, value, 0)
					if err != nil {
						goto wait
					}
					record("swap ok")
				} else {
					continue
				}
			}
		}
		record("create ok")
		// Get the lock.
		break
	}
	return
}

// Unlock.
func (m *Mutex) Unlock() (err error) {
	defer m.mutex.Unlock()
	for i := 1; i <= 3; i++ {
		_, err = m.client.Delete(m.key, false)
		if err != nil {
			if _, ok := err.(*etcd.EtcdError); !ok {
				// retry.
				continue
			}
		}
		break
	}
	return
}

var debug = false
var logger *log.Logger

func record(str string) {
	if debug {
		logger.Println(str)
	}
}
