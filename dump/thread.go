package dump

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/chandresh-pancholi/hdtd/communication"
	"github.com/chandresh-pancholi/hdtd/sink"
	"github.com/google/uuid"
	ps "github.com/mitchellh/go-ps"
)

type ThreadDump struct {
	aws   sink.AWS
	slack communication.Slack
}

func NewThreadDump(aws sink.AWS, slack communication.Slack) *ThreadDump {
	return &ThreadDump{
		aws:   aws,
		slack: slack,
	}
}

func (td *ThreadDump) Dump() string {
	dumpFile := fmt.Sprintf("%s%s", os.Getenv("LOG_FILE_PATH"), "td.out")
	fileName := path.Join(dumpFile)
	destination := uuid.New().String()

	threadDumpData, err := theadDump()
	if err != nil {
		log.Fatalf("thread dump failed. filename: %s, Error: %v", fileName, err)
		return ""
	}

	if threadDumpData != nil {
		log.Printf("compressing thread dump file. filename: %s", fileName)

		err = ioutil.WriteFile(fileName, threadDumpData, 0777)

		// util.Compress(fileName, threadDumpData)

		return destination
	}

	// _, err = td.aws.Upload(fileName, destination)
	// if err != nil {
	// 	log.Fatalf("heap dump s3 upload failed. Error: %v", err)
	// }

	return ""

}

func theadDump() ([]byte, error) {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return nil, err
	}

	// map ages
	for x := range processList {
		process := processList[x]
		if process.Executable() == "java" {
			log.Printf("generating thread dump. process_id: %d", process.Pid())
			out, err := exec.Command("jstack", "-l", strconv.Itoa(process.Pid())).Output()
			if err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, nil
}
