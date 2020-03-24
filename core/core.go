package core

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type Cfg struct {
	Interval int
	Dir      string
}

var defaultCfg = Cfg{
	Interval: 10,
	Dir:      "/tmp/gomemanalysis/",
}

type Msg struct {
	Timestamp    int64  `json:"timestamp"`
	Sys          uint64 `json:"sys"`
	HeapSys      uint64 `json:"heapsys"`
	HeapAlloc    uint64 `json:"heapalloc"`
	HeapInuse    uint64 `json:"heapinuse"`
	HeapReleased uint64 `json:"heapreleased"`
	HeapIdle     uint64 `json:"heapidle"`

	VMS uint64 `json:"vms"`
	RSS uint64 `json:"rss"`
}

type WithCfg func(cfg *Cfg)

var once sync.Once

func Start(cfgs ...WithCfg) error {
	var err error
	once.Do(func() {
		err = start(cfgs...)
	})
	http.ListenAndServe(":8081", nil)
	return err
}

func start(cfgs ...WithCfg) error {
	for _, cfg := range cfgs {
		cfg(&defaultCfg)
	}
	c, err := NewCollect(defaultCfg.Interval, defaultCfg.Dir)
	if err != nil {
		return err
	}
	go func() {
		c.collect()
	}()
	return nil
}
