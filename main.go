package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

var finalIndex int
var wheel = []string{
	" 0", "32", "15", "19", " 4", "21", " 2", "25", "17", "34", " 6",
	"27", "13", "36", "11", "30", " 8", "23", "10", " 5", "24", "16",
	"33", " 1", "20", "14", "31", " 9", "22", "18", "29", " 7", "28",
	"12", "35", " 3", "26",
}

var redNumbers = map[string]bool{
	" 1": true, " 3": true, " 5": true, " 7": true, " 9": true, "12": true,
	"14": true, "16": true, "18": true, "19": true, "21": true, "23": true,
	"25": true, "27": true, "30": true, "32": true, "34": true, "36": true,
}

var blackNumbers = map[string]bool{
	" 2": true, " 4": true, " 6": true, " 8": true, "10": true, "11": true,
	"13": true, "15": true, "17": true, "20": true, "22": true, "24": true,
	"26": true, "28": true, "29": true, "31": true, "33": true, "35": true,
}

func colorize(num string) string {
	switch {
	case redNumbers[num]:
		return "\033[101m" + num + "\033[0m" // red
	case blackNumbers[num]:
		return "\033[100m" + num + "\033[0m" // white/gray
	default:
		return "\033[102m" + num + "\033[0m" // green (0)
	}
}

func cleanupAndExit() {
	_ = keyboard.Close()
	fmt.Println("\n👋 Goodbye!")
	os.Exit(0)
}

func spinUntilStop(debug bool, startIndex int, stopChan chan struct{}) {
	index := startIndex
	for {
		select {
		case <-stopChan:
			if debug {
				log.Printf("Stopped at %s", colorize(wheel[index]))
			}
			finalIndex = index
			return
		default:
			fmt.Printf("\r⏳ %s ", colorize(wheel[index]))
			index = (index + 1) % len(wheel)
			time.Sleep(80 * time.Millisecond)
		}
	}
}

func spinWithInertia(debug bool, startIndex int) string {
	index := startIndex
	delay := 40 * time.Millisecond
	steps := rand.Intn(15) + 25

	if debug {
		log.Printf("Slowing down for %d steps\n", steps)
	}

	for i := 0; i < steps; i++ {
		fmt.Printf("\r🎲 %s ", colorize(wheel[index]))
		index = (index + 1) % len(wheel)
		time.Sleep(delay)
		delay += 10 * time.Millisecond // simulate slowdown
	}
	fmt.Println()
	return wheel[index]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var debug bool
	debug, err := strconv.ParseBool(os.Getenv("EUROULETTE_DEBUG"))
	if err != nil {
		debug = false
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Setup OS signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		cleanupAndExit()
	}()

	fmt.Println("🎰 European Roulette CLI")
	fmt.Println("Press SPACE to start spinning, SPACE again to stop.")
	fmt.Println("Press Ctrl+C to exit.\n")

	for {
		fmt.Println("🎯 Ready to spin. Press SPACE to start.")

		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeySpace {
				break
			}
			if key == keyboard.KeyCtrlC {
				cleanupAndExit()
			}
		}

		startIndex := rand.Intn(len(wheel))
		stopChan := make(chan struct{})

		go func() {
			spinUntilStop(debug, startIndex, stopChan)
		}()

		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeySpace {
				close(stopChan)
				time.Sleep(80 * time.Millisecond)
				break
			}
			if key == keyboard.KeyCtrlC {
				cleanupAndExit()
			}
		}

		fmt.Println("\n🌀 Slowing down...")
		if debug {
			log.Printf("Picked up at %s", colorize(wheel[finalIndex]))
		}
		result := spinWithInertia(debug, finalIndex)
		fmt.Printf("✅ Landed on: %s\n\n", colorize(result))
	}
}
