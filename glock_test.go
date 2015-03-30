// Copyright © 2015 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package glock

import (
	"log"
	"os"
	"runtime"
	"testing"
	"time"
)

func init() {
	logger = log.New(os.Stdout, "glock", log.Lshortfile|log.Ltime)
	debug = true
}

func TestLockNormal(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	slice := make([]int, 0, 2)
	m1 := NewMutex("key", ":10", 5, []string{"http://127.0.0.1:4001"})
	m2 := NewMutex("key", ":10", 5, []string{"http://127.0.0.1:4001"})
	m1.Lock()
	ch := make(chan bool)
	go func() {
		m2.Lock()
		slice = append(slice, 1)
		m2.Unlock()
		ch <- true
	}()
	slice = append(slice, 0)
	time.Sleep(2 * time.Second)
	m1.Unlock()
	<-ch
	for n, i := range slice {
		println("n,i:", n, i)
		if n != i {
			t.Fail()
		}
	}
}

func TestLockTimeout(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	slice := make([]int, 0, 2)
	m1 := NewMutex("key", ":10", 10, []string{"http://127.0.0.1:4001"})
	m2 := NewMutex("key", ":10", 10, []string{"http://127.0.0.1:4001"})
	m1.Lock()
	ch := make(chan bool)
	go func() {
		m2.Lock()
		slice = append(slice, 1)
		m2.Unlock()
		ch <- true
	}()
	slice = append(slice, 0)
	// m1.Unlock()
	<-ch
	for n, i := range slice {
		println("n,i:", n, i)
		if n != i {
			t.Fail()
		}
	}
}
