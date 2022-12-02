package hiro

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

func Add(source string, queueName string) {

	if queueName == "" {
		queueName = "main"
	}

	destination := path.Base(source)
	segmentsNumber := 16

	task := Task{Source: source, Destination: destination, SegmentsNumber: segmentsNumber}

	dataFile, err := os.ReadFile(DataFileLocation)
	if err != nil {
		log.Fatal(err)
	}

	var queues []Queue

	json.Unmarshal([]byte(dataFile), &queues)

	if len(queues) > 0 {
		for i, queue := range queues {
			if queue.Name == queueName {
				queues[i].Tasks = append(queue.Tasks, task)
				break
			} else {
				if i == len(queues)-1 {
					newQueue := Queue{Name: queueName, Tasks: []Task{task}}
					queues = append(queues, newQueue)
				}
			}

		}
	} else {
		newQueue := Queue{Name: queueName, Tasks: []Task{task}}
		queues = []Queue{newQueue}
	}

	data, _ := json.MarshalIndent(queues, "", "  ")

	os.WriteFile(DataFileLocation, data, os.ModePerm)
}
