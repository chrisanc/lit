package cli

import (
	"CLI_App/internal/domain"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// osClear - > Takes the name of the os and returns the command lines to clear the stdout.
var osClear = map[string][]string{
	"windows": {"cmd", "/c", "cls"},
	"linux":   {"clear"},
}

// GetNamingConvention - > Asks the user to select their favorite naming convention for the variables.
func GetNamingConvention() int8 {
	// Variable to save the value the user selects.
	var selectedConvention int8

	// Ask the user for a valid input n times
	for selectedConvention > 4 || selectedConvention < 1 {
		clearScreen(osClear[runtime.GOOS][0], osClear[runtime.GOOS][1:]...)
		fmt.Println("Select a valid naming convention.")
		fmt.Println("[1] camelCase")
		fmt.Println("[2] CamelCase")
		fmt.Println("[3] CamelCase/camelCase (for languages like Go)")
		fmt.Println("[4] snake_case")
		fmt.Scanf("%d", &selectedConvention)
	}

	// Return the pattern to use
	return selectedConvention
}

func GetAlertsConfig() domain.Alerts {
	clearScreen(osClear[runtime.GOOS][0], osClear[runtime.GOOS][1:]...)
	var res []domain.FeedbackValues
	configType := []string{
		"Type the minimum parameters of the method to mark it as ",
		"Type the minimum complexity of the method to mark it as ",
		"Type the minimum length of the method to mark it as ",
	}
	feedbackType := []string{
		"info\n", "warning\n", "error\n",
	}

	for _, cfg := range configType {
		var inputs []uint
		for _, fdb := range feedbackType {
			inputs = append(inputs, readValue(cfg+fdb, 1, uint(300)))
		}
		res = append(res, domain.FeedbackValues{
			Info:    inputs[0],
			Warning: inputs[1],
			Error:   inputs[2],
		})
	}

	return domain.Alerts{
		Parameters: res[0],
		Complexity: res[1],
		MethodSize: res[2],
	}
}

func readValue[T uint | int8](message string, lowRange, highRange T) T {
	input := lowRange - 1
	for input < lowRange || input > highRange {
		fmt.Printf(message)
		fmt.Scan(&input)
		clearScreen(osClear[runtime.GOOS][0], osClear[runtime.GOOS][1:]...)
		if input < lowRange || input > highRange {
			fmt.Fprintln(os.Stderr, "Invalid input...")
		}
	}

	return input
}

// Function used to clear the screen in Windows or Linux...
func clearScreen(name string, args ...string) {
	// Clear the screen
	cmd := exec.Command(name, args...)
	// Set the stdout we want to clear
	cmd.Stdout = os.Stdout
	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running the command to clear your screen.")
		os.Exit(1)
	}
}
