package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/eiannone/keyboard"
)

var wheel = []string{
	"0", "32", "15", "19", "4", "21", "2", "25", "17", "34", "6",
	"27", "13", "36", "11", "30", "8", "23", "10", "5", "24", "16",
	"33", "1", "20", "14", "31", "9", "22", "18", "29", "7", "28",
	"12", "35", "3", "26",
}

func cleanupAndExit() {
	_ = keyboard.Close()
	fmt.Println("\n👋 Goodbye!")
	os.Exit(0)
}

func spinUntilStop(startIndex int, stopChan chan struct{}) int {
	index := startIndex
	for {
		select {
		case <-stopChan:
			return index
		default:
			fmt.Printf("\r⏳ %s ", wheel[index])
			index = (index + 1) % len(wheel)
			time.Sleep(80 * time.Millisecond)
		}
	}
}

func spinWithInertia(startIndex int) string {
	index := startIndex
	delay := 40 * time.Millisecond
	steps := rand.Intn(15) + 25

	for i := 0; i < steps; i++ {
		fmt.Printf("\r🎲 %s ", wheel[index])
		index = (index + 1) % len(wheel)
		time.Sleep(delay)
		delay += 10 * time.Millisecond // simulate slowdown
	}
	fmt.Println()
	return wheel[index]
}

func main() {
	rand.Seed(time.Now().UnixNano())

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
		var finalIndex int

		go func() {
			finalIndex = spinUntilStop(startIndex, stopChan)
		}()

		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeySpace {
				close(stopChan)
				break
			}
			if key == keyboard.KeyCtrlC {
				cleanupAndExit()
			}
		}

		fmt.Println("\n🌀 Slowing down...")
		result := spinWithInertia(finalIndex)
		fmt.Printf("✅ Landed on: %s\n\n", result)
	}
}
