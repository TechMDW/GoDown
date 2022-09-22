package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/TechMDW/GoDown/internal/attacks"
	"github.com/TechMDW/GoDown/internal/terminal"
)

type History struct {
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: godown <command>")
		fmt.Println("Commands: httpflood, history")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "httpflood":
		httpflood := flag.NewFlagSet("httpflood", flag.ExitOnError)

		url := httpflood.String("url", "", "URL/IP/HOSTNAME to use")
		t := httpflood.Int64("t", 2, "timeout in sec")
		g := httpflood.Int("g", 1, "amount of goroutines")
		w := httpflood.Int("w", 50, "Time to wait between requests in ms")

		os.Stdout.Write([]byte("\033[2J\033[1;1H"))
		handleHttpFlood(httpflood, url, t, g, w)
	case "history":
		history := flag.NewFlagSet("history", flag.ExitOnError)

		os.Stdout.Write([]byte("\033[2J\033[1;1H"))
		handleHistory(history)
	}
}

func handleHttpFlood(c *flag.FlagSet, url *string, t *int64, g *int, w *int) {
	c.Parse(os.Args[2:])
	if *url == "" {
		fmt.Println("Usage: -url <url/ip/hostname>")
		os.Exit(1)
	}

	args := []string{
		fmt.Sprintf(`-url "%s"`, *url),
		fmt.Sprintf("-t %d", *t),
		fmt.Sprintf("-g %d", *g),
		fmt.Sprintf("-w %d", *w),
	}

	storeHistory("httpflood", args)
	channel := attacks.HttpFlood(*url, *t, *g, *w)

	sucess := 0
	errorGet := 0
	rateLimit := 0
	timeout := 0
	forbidden := 0

	fmt.Printf("\n")

	fmt.Println("Starting dos...")

	code := int64(200)
	start := time.Now()
	requests := int64(0)

	// Update terminal every second
	go func() {
		interval := time.NewTicker(time.Second)

		for range interval.C {
			// http.Client uses go routines
			numGoroutine := (runtime.NumGoroutine() - 1) / 2

			rps := float64(requests) / float64(time.Since(start).Seconds())
			rpm := rps * 60

			hoursSinceStart := time.Since(start).Hours()
			minutesSinceStart := time.Since(start).Minutes()
			secondsSinceStart := time.Since(start).Seconds()

			formatedTime := fmt.Sprintf("%02dh:%02dm:%02ds", int(hoursSinceStart), int(minutesSinceStart), int(secondsSinceStart))

			fmt.Printf("\033[2J\033[1;1H")

			fmt.Println(terminal.Output(terminal.BIYellow, "Press Ctrl+C to stop"))

			fmt.Printf("\n")

			fmt.Printf("Url: %s \n", terminal.Output(terminal.Blue, *url))
			fmt.Printf("Concurent requests: %s / %s \n", terminal.Output(terminal.Blue, numGoroutine), terminal.Output(terminal.Blue, *g))
			fmt.Printf("Timeout: %s sec\n", terminal.Output(terminal.Blue, *t))
			fmt.Printf("Time between requests: %s ms\n", terminal.Output(terminal.Blue, *w))
			fmt.Printf("Time since start: %s\n", terminal.Output(terminal.Blue, formatedTime))

			fmt.Printf("\n")

			fmt.Printf("%s \n", terminal.Output(terminal.Red, "-------Request data-------"))

			fmt.Printf("Requests: %s \n", terminal.Output(terminal.Blue, requests))
			fmt.Printf("Success: %s \n", terminal.Output(terminal.Blue, sucess))
			fmt.Printf("Error: %s \n", terminal.Output(terminal.Blue, errorGet))
			fmt.Printf("Rate limit: %s \n", terminal.Output(terminal.Blue, rateLimit))
			fmt.Printf("Timeout: %s \n", terminal.Output(terminal.Blue, timeout))
			fmt.Printf("Forbidden: %s \n", terminal.Output(terminal.Blue, forbidden))

			fmt.Printf("\n")

			fmt.Printf("%s \n", terminal.Output(terminal.Red, "-------Debug-------"))

			fmt.Printf("Requests per minute: %s \n", terminal.Output(terminal.Blue, fmt.Sprintf("%.0f", rpm)))
			fmt.Printf("Requests per second: %s \n", terminal.Output(terminal.Blue, fmt.Sprintf("%.0f", rps)))
			fmt.Printf("Last error code: %s \n", terminal.Output(terminal.Blue, code))
		}
	}()

	for {
		value := <-channel

		switch value {
		case "error":
			errorGet++
		case "rate limit":
			rateLimit++
		case "success":
			sucess++
		case "timeout":
			timeout++
		case "forbidden":
			forbidden++
		case "request":
			requests++
		default:
			if strings.Contains(value, "lastErrorCode") {
				codeString := strings.Split(value, ":")[1]

				parsedCode, err := strconv.ParseInt(strings.Trim(codeString, " "), 10, 64)

				if err == nil {
					code = parsedCode
				}
			}
		}
	}
}

func handleHistory(c *flag.FlagSet) {
	path, err := os.UserConfigDir()

	if err != nil {
		return
	}

	goDownHistoryPath := filepath.Join(path, "TechMDW", "GoDown", "history.json")

	file, err := os.OpenFile(goDownHistoryPath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		c.PrintDefaults()
		return
	}

	var history []History

	if _, err := file.Seek(0, 0); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	if err := json.NewDecoder(file).Decode(&history); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	for _, h := range history {
		fmt.Printf("%s\n", terminal.Output(terminal.BMagenta, h.Timestamp))

		// Color the url to stand out more
		// TODO: Make this better...
		// first := h.Args[0]
		// if strings.Contains(first, "https://") || strings.Contains(first, ".com") {
		// 	h.Args = append(h.Args[:0], h.Args[1:]...)

		// 	fmt.Printf("- %s %s %s\n\n", terminal.Output(terminal.Red, h.Command), terminal.Output(terminal.Green, first), terminal.Output(terminal.Yellow, strings.Join(h.Args, " ")))

		// 	continue
		// }

		fmt.Printf("- %s %s\n\n", terminal.Output(terminal.Red, h.Command), terminal.Output(terminal.Yellow, strings.Join(h.Args, " ")))
	}

	fmt.Printf("\n%s=Timestamp %s=Command %s=Arguments\n", terminal.Output(terminal.BMagenta, "Magenta"), terminal.Output(terminal.Red, "Red"), terminal.Output(terminal.Yellow, "Yellow"))
}

func storeHistory(command string, args []string) {
	path, err := os.UserConfigDir()

	if err != nil {
		return
	}

	// Join path
	techMDWPath := filepath.Join(path, "TechMDW")

	if _, err := os.Stat(techMDWPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(techMDWPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	goDownPath := filepath.Join(techMDWPath, "GoDown")

	if _, err := os.Stat(goDownPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(goDownPath, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	goDownHistoryPath := filepath.Join(goDownPath, "history.json")

	file, err := os.OpenFile(goDownHistoryPath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	log.SetFlags(16)

	var history []History

	if _, err := file.Seek(0, 0); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	if err := json.NewDecoder(file).Decode(&history); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	history = append(history, History{
		Command:   command,
		Args:      args,
		Timestamp: time.Now(),
	})

	if len(history) > 50 {
		history = history[1:]
	}

	if _, err := file.Seek(0, 0); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	if err := json.NewEncoder(file).Encode(history); err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	if err := file.Close(); err != nil {
		log.Println(err)
	}
}
