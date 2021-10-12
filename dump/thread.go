package dump

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/chandresh-pancholi/hdtd/communication"
	"github.com/chandresh-pancholi/hdtd/sink"
	"github.com/google/uuid"
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

func (td *ThreadDump) Dump(processId string) string {
	dumpFile := fmt.Sprintf("%s%s", os.Getenv("LOG_FILE_PATH"), "thread.tdump")
	fileName := path.Join(dumpFile)
	destination := uuid.New().String()

	threadDumpData, err := theadDump(processId)
	if err != nil {
		log.Fatalf("thread dump failed. filename: %s, Error: %v", fileName, err)
		return ""
	}

	if threadDumpData != nil {
		log.Printf("compressing thread dump file. filename: %s", fileName)

		err = ioutil.WriteFile(fileName, threadDumpData, 0777)

		return destination
	}

	_, err = td.aws.Upload(fileName, destination)
	if err != nil {
		log.Fatalf("heap dump s3 upload failed. Error: %v", err)
	}

	return ""

}

func theadDump(processId string) ([]byte, error) {
	out, err := exec.Command("jstack", "-l", processId).Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
