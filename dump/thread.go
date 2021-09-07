package dump

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strconv"

	ps "github.com/mitchellh/go-ps"
)

func ThreadDump() {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return
	}

	// map ages
	for x := range processList {
		process := processList[x]
		// log.Printf("%d\t%s\n", /\process.Pid(), process.Executable())

		if process.Executable() == "java" {
			log.Printf("process id: %d", process.Pid())
			out, err := exec.Command("jstack", "-l", strconv.Itoa(process.Pid())).Output()
			if err != nil {
				log.Fatal(err)
			}

			dumpFile := fmt.Sprintf("/%s/%s/%s", "tmp", strconv.Itoa(process.Pid()), "hd.hprof")
			fileName := path.Join(dumpFile)
			err = ioutil.WriteFile(fileName, out, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
