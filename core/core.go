package core

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type Option struct {
	IntervalSec int
	Dir         string
}

var defaultOption = Option{
	IntervalSec: 10,
	Dir:         "/tmp/gomemanalysis/",
}

type Info struct {
	Timestamp int64

	Sys          uint64 `json:"sys"`
	HeapSys      uint64 `json:"heapsys"`
	HeapAlloc    uint64 `json:"heapalloc"`
	HeapInuse    uint64 `json:"heapinuse"`
	HeapReleased uint64 `json:"heapreleased"`
	HeapIdle     uint64 `json:"heapidle"`

	VMS uint64 `json:"vms"`
	RSS uint64 `json:"rss"`
}

type WithOption func(option *Option)

var once sync.Once

func Start(opts ...WithOption) error {
	var err error
	once.Do(func() {
		err = start(opts...)
	})
	http.ListenAndServe(":8081", nil)
	return err
}

func start(opts ...WithOption) error {
	for _, mo := range opts {
		mo(&defaultOption)
	}
	c, err := NewCollect(defaultOption.IntervalSec, defaultOption.Dir)
	if err != nil {
		return err
	}
	go func() {
		c.collect()
	}()

	return nil
}
