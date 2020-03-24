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

	defaultUnit = UnitMByte
	unitString  = []string{"Byte", "KByte", "MByte", "GByte"}
)

const (
	prefix   = "gomemanalysis_"
	suffix   = ".dat"
	UnitByte = iota
	UnitKByte
	UnitMByte
	UnitGByte
)

func main() {
	dir = flag.String("dir", "/tmp/gomemanalysis/", "dir of gomemanalysis file")
	addr = flag.String("addr", ":80", "webui addr")
	uri = flag.String("uri", "/", "web uri")
	flag.Parse()
	fmt.Println("start....")

	http.HandleFunc(*uri, handler)
	err := http.ListenAndServe(*addr, nil)
	fmt.Println(err)
}

func handler(writer http.ResponseWriter, _ *http.Request) {
	infos, err := getInfos()
	if err != nil {
		return
	}

	var (
		t            []string
		sys          []float64
		heapSys      []float64
		heapAlloc    []float64
		heapInuse    []float64
		heapReleased []float64
		heapIdle     []float64
		VMS          []float64
		RSS          []float64
	)

	for _, info := range infos {
		t = append(t, time.Unix(info.Timestamp, 0).Format("2006-01-02 15:04:05"))
		sys = append(sys, formatMemWithUnit(info.Sys, defaultUnit))
		heapSys = append(heapSys, formatMemWithUnit(info.HeapSys, defaultUnit))
		heapAlloc = append(heapAlloc, formatMemWithUnit(info.HeapAlloc, defaultUnit))
		heapInuse = append(heapInuse, formatMemWithUnit(info.HeapInuse, defaultUnit))
		heapReleased = append(heapReleased, formatMemWithUnit(info.HeapReleased, defaultUnit))
		heapIdle = append(heapIdle, formatMemWithUnit(info.HeapIdle, defaultUnit))
		VMS = append(VMS, formatMemWithUnit(info.VMS, defaultUnit))
		RSS = append(RSS, formatMemWithUnit(info.RSS, defaultUnit))
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.InitOpts{PageTitle: "gomemanalysis", Theme: charts.ThemeType.Infographic},
	)
	line.Title = "unit：" + unitString[defaultUnit]
	line.AddXAxis(t)
	opts := charts.LineOpts{Smooth: true}
	line.AddYAxis("sys", sys, opts)
	line.AddYAxis("heapSys", heapSys, opts)
	line.AddYAxis("heapAlloc", heapAlloc, opts)
	line.AddYAxis("heapInuse", heapInuse, opts)
	line.AddYAxis("heapReleased", heapReleased, opts)
	line.AddYAxis("heapIdle", heapIdle, opts)
	line.AddYAxis("VMS", VMS, opts)
	line.AddYAxis("RSS", RSS, opts)
	line.Render(writer)
}

func getInfos() ([]core.Msg, error) {
	var filename string
	var latest int64

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
		if t.Unix() > latest {
			latest = t.Unix()
			filename = path
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return readFile(filename)
}

func readFile(file string) ([]core.Msg, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(content, []byte{'\n'})
	var ret []core.Msg
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var msg core.Msg
		if err := json.Unmarshal(line, &msg); err != nil {
			return nil, err
		}
		ret = append(ret, msg)
	}
	return ret, nil
}

func formatMemWithUnit(nByte uint64, unit int) float64 {
	switch unit {
	case UnitByte:
		return float64(nByte)
	case UnitKByte:
		return float64(nByte) / 1024
	case UnitMByte:
		return float64(nByte) / 1024 / 1024
	case UnitGByte:
		return float64(nByte) / 1024 / 1024
	default:
		return 0
	}
}
