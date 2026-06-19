package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tphummel/dice-collector/internal/dice"
)

func main() {
	dbPath := "dice.db"
	if v := os.Getenv("DICE_DB_PATH"); v != "" {
		dbPath = v
	}

	store, err := dice.Open(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening database:", err)
		os.Exit(1)
	}
	defer store.Close()

	shooters, err := store.Shooters()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading shooters:", err)
		os.Exit(1)
	}
	locations, err := store.Locations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading locations:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Craps Session Entry System")
	fmt.Println()

	for _, loc := range locations {
		fmt.Printf("%d. %s (%s, %s)\n", loc.ID, loc.Name, loc.City, loc.State)
	}
	sessionLocation := chooseByID(reader, "Choose Location: ", locations, func(l dice.Location) int64 { return l.ID })

	sessionID, err := store.InsertSession(sessionLocation.ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating session:", err)
		os.Exit(1)
	}
	fmt.Printf("\nSession Created: id=%d\n", sessionID)
	printThrowEncodingReference()

	var summary strings.Builder
	anotherTurn := "Y"
	for anotherTurn == "Y" {
		fmt.Println()
		for _, shooter := range shooters {
			fmt.Printf("%d. %s\n", shooter.ID, shooter.Name)
		}
		currentShooter := chooseByID(reader, "Choose Shooter: ", shooters, func(s dice.Shooter) int64 { return s.ID })

		turnID, err := store.InsertTurn(sessionID, currentShooter.ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating turn:", err)
			os.Exit(1)
		}
		fmt.Printf("\nTurn Created: %s at %s (%d)\n", currentShooter.Name, sessionLocation.Name, turnID)

		throwsInput := promptLine(reader, "Enter all throws , no spaces, no commas: ")
		throws, err := dice.ProcessTurnThrows(throwsInput)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error processing throws:", err)
			os.Exit(1)
		}

		for _, t := range throws {
			if err := store.InsertThrow(sessionID, currentShooter.ID, turnID, t.Value, t.Result, t.PropComeout); err != nil {
				fmt.Fprintln(os.Stderr, "error saving throw:", err)
				os.Exit(1)
			}
		}

		fmt.Printf("Turn Saved. %s: %s\n", currentShooter.Name, throwsInput)

		summary.WriteString(currentShooter.Name)
		summary.WriteString(": ")
		values := make([]string, len(throws))
		for i, t := range throws {
			values[i] = strconv.Itoa(t.Value)
		}
		summary.WriteString(strings.Join(values, ","))
		summary.WriteString("\n")

		anotherTurn = strings.ToUpper(promptLine(reader, "Would you like to enter the next turn? (Y|N): "))
	}

	fmt.Println("Session Complete")
	fmt.Println("Summary--", sessionLocation.Name)
	fmt.Println(strings.TrimRight(summary.String(), "\n"))
}

func printThrowEncodingReference() {
	fmt.Println()
	fmt.Println("Throw entry reference:")
	fmt.Println("  2-9 : Face value")
	fmt.Println("  T   : 10")
	fmt.Println("  E   : 11 (yo-leven)")
	fmt.Println("  B   : 12 (boxcars)")
	fmt.Println("  O   : Off table")
	fmt.Println("  M   : Misc / invalid roll")
}

func promptLine(r *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := r.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

// chooseByID re-prompts until the user enters the id of one of options.
func chooseByID[T any](r *bufio.Reader, prompt string, options []T, idOf func(T) int64) T {
	for {
		input := promptLine(r, prompt)
		choice, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if err == nil {
			for _, opt := range options {
				if idOf(opt) == choice {
					return opt
				}
			}
		}
		fmt.Println("Please enter one of the listed numbers.")
	}
}
