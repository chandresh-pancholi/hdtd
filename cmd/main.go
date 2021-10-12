package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chandresh-pancholi/hdtd/communication"
	"github.com/chandresh-pancholi/hdtd/dump"
	"github.com/chandresh-pancholi/hdtd/sink"
	"github.com/chandresh-pancholi/hdtd/util"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	// c1 := make(chan string)
	// c2 := make(chan string)

	s3 := sink.NewAwsS3Client()
	slack := communication.NewSlack()

	// go func() {
	// 	d := dump.NewHeapDump(*s3, *slack)
	// 	c1 <- d.Dump()
	// }()
	// go func() {
	// 	td := dump.NewThreadDump(*s3, *slack)
	// 	c2 <- td.Dump()
	// }()

	// for i := 0; i < 2; i++ {
	// 	select {
	// 	case msg1 := <-c1:
	// 		log.Printf("heap dump generated. state: %v", msg1)
	// 	case msg2 := <-c2:
	// 		log.Printf("thread dump generated. state: %v", msg2)
	// 	}
	// }

	d := dump.NewHeapDump(*s3, *slack)
	td := dump.NewThreadDump(*s3, *slack)

	processId := os.Args[1]
	hdDestination := d.Dump(processId)
	tdDestination := td.Dump(processId)

	slackMessage := fmt.Sprintf("Pod name: %s\n Heap dump: %s\n Thread dump: %s", os.Getenv("POD_NAME"), hdDestination, tdDestination)

	log.Println("sending slack message")
	err := slack.Publish(slackMessage)
	if err != nil {
		log.Fatalf("publishing message to slack failed. Error: %v", err)
	}

	util.KillProcess(processId)
}
