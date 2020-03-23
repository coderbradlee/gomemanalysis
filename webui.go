package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/lzxm160/gomemanalysis/core"
)

var (
	dir  *string
	addr *string
	uri  *string

	memUnit = MemUnit(MemUnitMByte)
)

type MemUnit int

const (
	prefix              = "gomemanalysis_"
	suffix              = ".dat"
	MemUnitByte MemUnit = iota + 1
	MemUnitKByte
	MemUnitMByte
	MemUnitGByte
)

var memUnitBrief []string

func main() {
	dir = flag.String("dir", "/tmp/gomemanalysis/", "dir of pprofplus dump file")
	addr = flag.String("addr", ":80", "dashboard addr")
	uri = flag.String("uri", "/", "web uri")
	flag.Parse()
	fmt.Println("start....")

	http.HandleFunc(*uri, requestHandler)
	err := http.ListenAndServe(*addr, nil)
	fmt.Println(err)
}

func getInfos() ([]core.Info, error) {
	var filename string
	var newestUnix int64

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && path != *dir {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(info.Name(), prefix) || !strings.HasSuffix(info.Name(), suffix) {
			return nil
		}
		pidTime := strings.TrimSuffix(strings.TrimPrefix(info.Name(), prefix), suffix)
		pid_time := strings.Split(pidTime, "_")
		if len(pid_time) != 2 {
			return nil
		}
		t, err := time.Parse("20060102150405", pid_time[1])
		if err != nil {
			return nil
		}
		if t.Unix() > newestUnix {
			newestUnix = t.Unix()
			filename = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("filename=%s\n", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(content, []byte{'\n'})
	var ret []core.Info
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var info core.Info
		if err := json.Unmarshal(line, &info); err != nil {
			return nil, err
		}
		ret = append(ret, info)
	}
	return ret, nil
}

func requestHandler(writer http.ResponseWriter, _ *http.Request) {
	infos, err := getInfos()
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	var x []string
	var Sys []float64
	var HeapSys []float64
	var HeapAlloc []float64
	var HeapInuse []float64
	var HeapReleased []float64
	var HeapIdle []float64
	//var HeapIdleMinusRleased []float64
	var VMS []float64
	var RSS []float64
	for _, info := range infos {
		x = append(x, time.Unix(info.Timestamp, 0).Format("01-02 15:04:05"))
		Sys = append(Sys, calcMemWithUnit(info.Sys, memUnit))
		HeapSys = append(HeapSys, calcMemWithUnit(info.HeapSys, memUnit))
		HeapAlloc = append(HeapAlloc, calcMemWithUnit(info.HeapAlloc, memUnit))
		HeapInuse = append(HeapInuse, calcMemWithUnit(info.HeapInuse, memUnit))
		HeapReleased = append(HeapReleased, calcMemWithUnit(info.HeapReleased, memUnit))
		HeapIdle = append(HeapIdle, calcMemWithUnit(info.HeapIdle, memUnit))
		//HeapIdleMinusRleased = append(HeapIdleMinusRleased, calcMemWithUnit(info.ms.HeapIdle-info.ms.HeapReleased, option.MemUint))
		VMS = append(VMS, calcMemWithUnit(info.VMS, memUnit))
		RSS = append(RSS, calcMemWithUnit(info.RSS, memUnit))
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "gomemanalysis", Theme: charts.ThemeType.Infographic},
	)
	line.Title = "unit：" + memUnitBrief[memUnit]

	line.AddXAxis(x)
	opts := charts.LineOpts{Smooth: true}
	line.AddYAxis("Sys", Sys, opts)
	line.AddYAxis("HeapSys", HeapSys, opts)
	line.AddYAxis("HeapAlloc", HeapAlloc, opts)
	line.AddYAxis("HeapInuse", HeapInuse, opts)
	line.AddYAxis("HeapReleased", HeapReleased, opts)
	line.AddYAxis("HeapIdle", HeapIdle, opts)
	//line.AddYAxis("HeapIdleMinusRleased", HeapIdleMinusRleased, opts)
	line.AddYAxis("VMS", VMS, opts)
	line.AddYAxis("RSS", RSS, opts)
	line.Render(writer)
}

func calcMemWithUnit(nByte uint64, unit MemUnit) float64 {
	switch unit {
	case MemUnitByte:
		return float64(nByte)
	case MemUnitKByte:
		return float64(nByte) / 1024
	case MemUnitMByte:
		return float64(nByte) / 1024 / 1024
	case MemUnitGByte:
		return float64(nByte) / 1024 / 1024
	}
	panic("never reach here.")
}

func init() {
	memUnitBrief = []string{"wrong", "Byte", "KByte", "MByte", "GByte"}
}
