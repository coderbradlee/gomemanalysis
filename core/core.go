package core

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type Option struct {
	CaptureIntervalSec int
	DumpDir            string
	ServiceName        string
}

var option = Option{
	CaptureIntervalSec: 5,
	DumpDir:            "/tmp/gomemanalysis/",
	ServiceName:        "gomemanalysis",
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

type ModOption func(option *Option)

var once sync.Once

func Start(modOptions ...ModOption) error {
	var err error
	once.Do(func() {
		err = start(modOptions...)
	})
	http.ListenAndServe(":8081", nil)
	return err
}

func start(modOptions ...ModOption) error {
	for _, mo := range modOptions {
		mo(&option)
	}

	c, err := NewCollect(option.CaptureIntervalSec, option.DumpDir, option.ServiceName)
	if err != nil {
		return err
	}

	go func() {
		infoC := c.doAsync()
		for {
			info := <-infoC
			d.do(info)
		}
	}()

	return nil
}
