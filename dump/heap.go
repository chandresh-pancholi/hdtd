package dump

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	ps "github.com/mitchellh/go-ps"
)

func Dump() {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return
	}

	// map ages
	for x := range processList {
		process := processList[x]

		if process.Executable() == "java" {
			log.Printf("process id: %d", process.Pid())
			dumpFile := fmt.Sprintf("/%s/%s/%s", "tmp", strconv.Itoa(process.Pid()), "hd.hprof")

			_, err := exec.Command("jmap", fmt.Sprintf("-dump:format=b,file=%s", dumpFile), strconv.Itoa(process.Pid())).Output()
			if err != nil {
				log.Fatal(err)
			}

		}
	}
}
