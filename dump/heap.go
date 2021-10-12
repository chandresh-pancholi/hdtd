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

type HeapDump struct {
	aws   sink.AWS
	slack communication.Slack
}

func NewHeapDump(aws sink.AWS, slack communication.Slack) *HeapDump {
	return &HeapDump{
		aws:   aws,
		slack: slack,
	}
}

func (hd *HeapDump) Dump(processId string) string {
	dumpFile := fmt.Sprintf("%s%s", os.Getenv("LOG_FILE_PATH"), "hd.hprof")
	fileName := path.Join(dumpFile)

	destination := uuid.New().String()

	_, err := heapDump(processId, fileName)
	if err != nil {
		log.Fatalf("heap dump generation failed. Error: %v", err)
	}

	log.Printf("uploading file to S3. filename: %s", fileName)
	_, err = hd.aws.Upload(fileName, destination)
	if err != nil {
		log.Fatalf("heap dump s3 upload failed. Error: %v", err)
	}

	return destination
}

func heapDump(processId, filename string) ([]byte, error) {
	log.Printf("generating heap dump. process_id: %s", processId)
	_, err := exec.Command("jmap", fmt.Sprintf("-dump:format=b,file=%s", filename), processId).Output()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}
