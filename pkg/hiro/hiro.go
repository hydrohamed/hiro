package hiro

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Queue struct {
	Name  string
	Tasks []Task
}

type Task struct {
	Source         string
	Destination    string
	SegmentsNumber int
	Segments       []Segment
	Completed      bool
}

type Segment struct {
	Start  int
	End    int
	Status int
}

var Tasks []Task

var DataFileDir string
var DataFileLocation string

var LastPrintedTime time.Time
var LastPrintedStatus int

func init() {
	homeDir, err := os.UserHomeDir()
	DataFileDir = homeDir + "/.config/hiro/"
	DataFileLocation = DataFileDir + "data.json"
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(DataFileLocation); errors.Is(err, os.ErrNotExist) {

		if os.IsNotExist(err) {
			os.MkdirAll(DataFileDir, os.ModePerm)
			os.Create(DataFileLocation)
		}

		newQueue := Queue{Name: "main", Tasks: []Task{}}
		queues := []Queue{newQueue}
		data, _ := json.MarshalIndent(queues, "", "  ")
		os.WriteFile(DataFileLocation, data, os.ModePerm)
	}
}

func generateRandomStr(n int) string {
	seed := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	randomStr := make([]byte, n)

	for i := range randomStr {
		randomStr[i] = seed[rand.Intn(len(seed))]
	}
	return string(randomStr)
}

func Run(i int, j int) error {

	fmt.Printf("The download of %v has been started...\n", Queues[i].Tasks[j].Destination)

	task := Queues[i].Tasks[j]

	if len(task.Segments) == 0 {

		segments := createSegments(task)
		Queues[i].Tasks[j].Segments = segments
		data, err := json.MarshalIndent(Queues, "", "  ")

		if err != nil {
			log.Fatal(err)
			return err
		}

		os.WriteFile(DataFileLocation, data, os.ModePerm)
	}

	wg := new(sync.WaitGroup)

	rand.Seed(time.Now().UnixNano())
	tempFileName := generateRandomStr(16)

	destinationFile, _ := os.OpenFile(tempFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	for k, segment := range Queues[i].Tasks[j].Segments {
		wg.Add(1)
		go downloadSegment(i, j, k, segment, wg, destinationFile)
	}

	wg.Wait()

	Queues[i].Tasks[j].Completed = true

	data, _ := json.MarshalIndent(Queues, "", "  ")
	os.WriteFile(DataFileDir+"data.json", data, os.ModePerm)

	destinationFile.Close()

	fmt.Printf("The download of %v has been completed!\n\n", Queues[i].Tasks[j].Destination)
	os.Rename(tempFileName, task.Destination)
	return nil
}

func createSegments(task Task) []Segment {

	res, _ := http.Head(task.Source)
	header := res.Header
	length, _ := strconv.Atoi(header["Content-Length"][0])
	segmentSize := length / task.SegmentsNumber
	segment := Segment{Start: 0, End: segmentSize, Status: 0}
	segments := []Segment{segment}

	for i := 1; i < task.SegmentsNumber; i++ {
		segments = append(segments, Segment{i*segmentSize + 1, (i + 1) * segmentSize, i*segmentSize + 1})
	}

	segments[task.SegmentsNumber-1].End = length
	return segments
}

func downloadSegment(i int, j int, k int, segment Segment, wg *sync.WaitGroup, f *os.File) {
	defer wg.Done()
	task := Queues[i].Tasks[j]
	chunkStart := segment.Start

	if segment.Status != segment.Start {
		chunkStart = segment.Status
	}

	for chunkStart < segment.End {
		chunkEnd := chunkStart + 16000

		if chunkEnd > segment.End {
			chunkEnd = segment.End
		}

		req, _ := http.NewRequest("GET", task.Source, nil)
		range_header := "bytes=" + strconv.Itoa(chunkStart) + "-" + strconv.Itoa(chunkEnd)
		req.Header.Add("Range", range_header)
		resp, _ := http.DefaultClient.Do(req)
		responseBody, _ := io.ReadAll(resp.Body)

		f.WriteAt(responseBody, int64(chunkStart))

		Queues[i].Tasks[j].Segments[k].Status = chunkEnd
		updateStatus()
		go printStatus(i, j)
		chunkStart = chunkStart + 16001
	}

}

func updateStatus() {
	data, _ := json.MarshalIndent(Queues, "", "  ")
	err := os.WriteFile(DataFileDir+"data.json", data, 0644)

	if err != nil {
		log.Fatal(err)
		return
	}

}

func printStatus(i int, j int) {
	timeInterval := time.Since(LastPrintedTime).Seconds()

	if timeInterval < 1 {
		return
	}

	task := Queues[i].Tasks[j]

	downloaded := 0

	for _, segment := range task.Segments {
		downloaded = downloaded + segment.Status - segment.Start
	}

	fileLength := task.Segments[len(task.Segments)-1].End
	progress := int(float32(downloaded) / float32(fileLength) * 100)

	newDownloaded := downloaded - LastPrintedStatus
	transferRate := (float32(newDownloaded) / float32(timeInterval) / 1000)
	fmt.Printf("(%v%%): Downloaded %vKB from %vKB. (Transfer Rate: %vKB/s)\n", progress, downloaded/1000, fileLength/1000, int(transferRate))
	LastPrintedTime = time.Now()
	LastPrintedStatus = downloaded
}
