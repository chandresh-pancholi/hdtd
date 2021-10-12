package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/process"
)

func Compress(name string, data []byte) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	// Setting the Header fields is optional.
	zw.Name = name
	zw.ModTime = time.Now()

	_, err := zw.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
}

func KillProcess(processId string) error {
	pid, _ := strconv.Atoi(processId)

	processes, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return err
		}
		fmt.Printf("process name: %s, PID: %d", n, p.Pid)
		fmt.Println()
		if p.Pid == int32(pid) {
			return p.Kill()
		}
	}
	return fmt.Errorf("process not found. procesId: %s", processId)

}
