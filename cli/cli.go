package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jcondeco207/docker-helper/containers"
	"github.com/manifoldco/promptui"
)

func Checkboxes(label string, opts []string) []string {
	res := []string{}
	prompt := &survey.MultiSelect{
		Message: label,
		Options: opts,
	}
	survey.AskOne(prompt, &res)

	return res
}

func pickRunningContainers() []string {
	runningContainers := containers.GetRunningContainers()

	containersOptions := make([]string, len(runningContainers))

	for i, container := range runningContainers {
		containersOptions[i] = fmt.Sprintf("%s %s", container.Names[0], container.ID)
	}

	answers := Checkboxes(
		"Running containers:",
		containersOptions,
	)

	return answers

}

func pickStoppedContainers() []string {
	runningContainers := containers.GetStoppedContainers()
	containersOptions := make([]string, len(runningContainers))

	for i, container := range runningContainers {
		containersOptions[i] = fmt.Sprintf("%s %s", container.Names[0], container.ID)
	}

	answers := Checkboxes(
		"Running containers:",
		containersOptions,
	)

	return answers

}

func pickAndStartContainer() {
	answers := pickStoppedContainers()
	for _, option := range answers {
		var id string
		var name string

		_, err := fmt.Sscanf(option, "%s %s", &name, &id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		containers.StartContainer(id)
	}
}

func pickAndStopContainer() {
	answers := pickRunningContainers()
	for _, option := range answers {
		var id string
		var name string

		_, err := fmt.Sscanf(option, "%s %s", &name, &id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		containers.StopContainer(id)
	}
}

func pickAndExec() {
	runningContainers := containers.GetRunningContainers()

	containersOptions := make([]string, len(runningContainers))

	for i, container := range runningContainers {
		containersOptions[i] = fmt.Sprintf("%s %s", container.Names[0], container.ID)
	}

	prompt := promptui.Select{
		Label: "Select Day",
		Items: containersOptions,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	var id string
	var name string

	_, err = fmt.Sscanf(result, "%s %s", &name, &id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var command string
	fmt.Print("Enter some text: ")
	_, err = fmt.Scanln(&command)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	cmd := strings.Fields(command)

	containers.ExecFunction(id, cmd)

}

func attachToContainer() {
	runningContainers := containers.GetRunningContainers()

	containersOptions := make([]string, len(runningContainers))

	for i, container := range runningContainers {
		containersOptions[i] = fmt.Sprintf("%s %s", container.Names[0], container.ID)
	}

	prompt := promptui.Select{
		Label: "Select container",
		Items: containersOptions,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	var id string
	var name string

	_, err = fmt.Sscanf(result, "%s %s", &name, &id)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	containers.AttachToContainer(id)
}

func pickAction() bool {
	prompt := promptui.Select{
		Label: "Select action:",
		Items: []string{"Start selected containers", "Exec", "Stop selected containers", "Show running containers", "Attach to container", "Quit"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	switch result {

	case "Start selected containers":
		pickAndStartContainer()

	case "Stop selected containers":
		pickAndStopContainer()

	case "Exec":
		pickAndExec()

	case "Attach to container":
		attachToContainer()

	case "Show running containers":
		containers.ShowRunning()

	case "Quit":
		return false
	}

	return true
}

func main() {
	cont := true
	for cont {
		cont = pickAction()
		fmt.Printf("\n\n")
	}
}
