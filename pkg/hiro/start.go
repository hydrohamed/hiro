package hiro

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var Queues []Queue

func Start(queueName string, lastTask bool) {

	if queueName == "" {
		queueName = "main"
	}

	dataFile, err := os.ReadFile(DataFileLocation)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal([]byte(dataFile), &Queues)

	if lastTask {
		tasks := Queues[0].Tasks
		Run(0, len(tasks)-1)
		return
	}

	for i := range Queues {
		if Queues[i].Name == queueName {
			fmt.Printf("Queue %v has been started...\n\n", Queues[i].Name)
			startTime := time.Now()
			for j, task := range Queues[i].Tasks {

				if task.Completed {
					continue
				}

				LastPrintedTime = time.Now()
				LastPrintedStatus = 0
				Run(i, j)
			}
			fmt.Printf("Queue %v has been completed in %v seconds!\n", Queues[i].Name, int(time.Since(startTime).Seconds()))
			break
		} else {
			if i == len(Queues)-1 {
				fmt.Println("There is no such queue!")
			}
		}

	}

}
