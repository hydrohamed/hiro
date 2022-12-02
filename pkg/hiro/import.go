package hiro

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func Import(fileLocation string, queueName string) {
	file, err := os.Open(fileLocation)

	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for i := 1; scanner.Scan(); i++ {
		Add(scanner.Text(), queueName)
		fmt.Printf("add %v to %v\n", scanner.Text(), queueName)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
