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

func (hd *HeapDump) Dump() string {
	dumpFile := fmt.Sprintf("%s%s", os.Getenv("LOG_FILE_PATH"), "hd.hprof")
	fileName := path.Join(dumpFile)

	destination := uuid.New().String()

	_, err := heapDump(fileName)
	if err != nil {
		log.Fatalf("heap dump generation failed. Error: %v", err)
	}

	log.Printf("compressing file. filename: %s", fileName)
	// util.Compress(fileName, data)

	log.Printf("uploading file to S3. filename: %s", fileName)
	// _, err = hd.aws.Upload(fileName, destination)
	// if err != nil {
	// 	log.Fatalf("heap dump s3 upload failed. Error: %v", err)
	// }

	return destination
}

func heapDump(filename string) ([]byte, error) {
	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return nil, err
	}

	// map ages
	for x := range processList {
		process := processList[x]

		if process.Executable() == "java" {
			log.Printf("generating heap dump. process_id: %d", process.Pid())
			_, err := exec.Command("jmap", fmt.Sprintf("-dump:format=b,file=%s", filename), strconv.Itoa(process.Pid())).Output()
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}
