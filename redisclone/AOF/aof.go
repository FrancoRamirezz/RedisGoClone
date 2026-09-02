package aofgo

import (
	"Pokemonscraper/redisclone/resp"
	"bufio"
	"fmt"
	"strings"

	"os"
	"sync"
	"time"
)

// the append only file  is one aspect of redis data presitence
// / the data is stored in the memory then must pass to the disk
type AOF struct {
	file   *os.File      // refers to a file
	rd     *bufio.Reader // the bufio reader reads from the file
	mu     sync.Mutex
	syncPd time.Duration
	done   chan struct{}
}

// the AOF creates a file to write to
// a way to read the file on startup
// We need to handle concurrent request so we will use a mutex
// a sync period. the redis AOF allows use to set the flush interval, or we can make a function for the flush interval
func NewAOFPath(path string) (*AOF, error) {
	// create the file it it does not exist
	// O_Create creates a file if it dosent exisit
	//O_Wronly we only write here to its own handle
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("error opening file:%v ", err)
	}
	aof := &AOF{file: f, rd: bufio.NewReader(f), syncPd: time.Second}
	// here we will pass the go routine, that moves the disk for every 1 second
	go aof.Flush()
	return aof, nil
}

// the flush is the go routine to sync AOF to disk for every second
// we are writing each command to the file: like SET and DEL
func (a *AOF) Flush() {
	for {
		a.mu.Lock()
		a.file.Sync()
		a.mu.Unlock()
		time.Sleep(a.syncPd)
	}
}

// we need to write to the Aof, when we recive commands like SET or DEL
// this will write
func (a *AOF) Write(resp resp.RespValue) error {
	close(a.done)
	a.mu.Lock()
	// defer closes the connection
	defer a.mu.Unlock()
	// resp the array header
	// read each arguement as a bulk string
	// Marshal() to write the command to the file in the resp format that we receive
	_, err := a.file.Write(resp.Marshal())
	if err != nil {
		return err
	}
	return nil
}

// Close the file and ensure the flush is closed when the served is closed
// call this when the server must gracefully shutdown, just in case
func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.file.Close()

}

// when the power go down, and we need to reboot the server
// to make sure the commands are not lost in restart we will write the data to a ssd usiing aof

func ReadAof(path string, handler func(args []string) error) error {
	// open the file name
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error reading aof to disk:")
	}
	// always close the file
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var args []string
	var expected int
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "*"):
			// for every array header takes *3 means 3 bulk strings form one command
			fmt.Sscanf(line[1:], "%d", &expected)
			args = make([]string, 0, expected)
		// Bulk string length prefix which we usuall skip
		case strings.HasPrefix(line, "$"):
		default:
			// collect the payload from the bulk
			args = append(args, line)
			if len(args) == expected && expected > 0 {
				// Complete command assembled — re-execute it.
				if err := handler(args); err != nil {
					fmt.Fprintf(os.Stderr, "aof: replay skipping bad entry: %v\n", err)
				}
				args = nil
				expected = 0
			}
		}
	}
	return nil
}
