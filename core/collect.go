package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/process"
)

type collect struct {
	intervalSec int
	file        *os.File
}

func NewCollect(intervalSec int, dir string) (*collect, error) {
	if err := os.MkdirAll(dir, 0666); err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("gomemanalysis_%d_%s.dat",
		os.Getpid(), time.Now().Format("20060102150405"))
	filename = path.Join(dir, filename)

	fp, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &collect{
		intervalSec: intervalSec, file: fp,
	}, nil
}

func (c *collect) collect() {
	go func() {
		p := process.Process{
			Pid: int32(os.Getpid()),
		}
		t := time.Tick(time.Second * time.Duration(c.intervalSec))
		errChan := make(chan error, 1)
		for {
			select {
			case <-t:
				e := c.do(p)
				if e != nil {
					errChan <- e
				}

			case err := <-errChan:
				fmt.Println(err)
			}
		}
	}()
}

func (c *collect) do(p process.Process) error {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	mis, _ := p.MemoryInfo()

	info := &Info{
		Timestamp:    time.Now().Unix(),
		Sys:          ms.Sys,
		HeapSys:      ms.HeapSys,
		HeapAlloc:    ms.HeapAlloc,
		HeapInuse:    ms.HeapInuse,
		HeapReleased: ms.HeapReleased,
		HeapIdle:     ms.HeapIdle,
		VMS:          mis.VMS,
		RSS:          mis.RSS,
	}
	raw, _ := json.Marshal(&info)
	_, err := c.file.Write(raw)
	if err != nil {
		return err
	}
	_, err = c.file.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return c.file.Sync()
}
