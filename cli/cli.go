package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jcondeco207/docker-helper/containers"
	"github.com/jcondeco207/docker-helper/images"
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

func pickImages() []string {
	images := images.GetAllImages()

	imagesOptions := make([]string, len(images))

	for i, image := range images {
		imagesOptions[i] = fmt.Sprintf("%s %s", image.RepoTags[0], image.ID)
	}

	answers := Checkboxes(
		"Images:",
		imagesOptions,
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

func pickAndDeleteContainer() {
	answers := pickStoppedContainers()
	for _, option := range answers {
		var id string
		var name string

		_, err := fmt.Sscanf(option, "%s %s", &name, &id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		containers.DeleteContainer(id)
	}
}

func pickAndDeleteImage() {
	answers := pickImages()
	for _, option := range answers {
		var id string
		var name string

		_, err := fmt.Sscanf(option, "%s %s", &name, &id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		images.DeleteImage(id)
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

func pickContainersAction() bool {
	prompt := promptui.Select{
		Size:  10,
		Label: "Select action",
		Items: []string{"Start selected containers",
			"Show containers",
			"Show images",
			"Exec",
			"Stop selected containers",
			"Delete selected containers",
			"Show running containers",
			"Attach to container",
			"Return"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	switch result {

	case "Start selected containers":
		pickAndStartContainer()

	case "Show images":
		images.ShowAllImages()

	case "Show containers":
		containers.ShowAllContainers()

	case "Stop selected containers":
		pickAndStopContainer()

	case "Delete selected containers":
		pickAndDeleteContainer()

	case "Exec":
		pickAndExec()

	case "Attach to container":
		attachToContainer()

	case "Show running containers":
		containers.ShowRunning()

	case "Return":
		return true
	}

	return false
}

func pickImagesAction() bool {
	prompt := promptui.Select{
		Size:  10,
		Label: "Select action",
		Items: []string{
			"Show images",
			"Delete images",
			"Return"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}

	switch result {

	case "Show images":
		images.ShowAllImages()

	case "Delete images":
		pickAndDeleteImage()

	case "Return":
		return true
	}

	return false
}

func pickContext() bool {
	cont := true
	for cont {
		prompt := promptui.Select{
			Size:  10,
			Label: "Select context",
			Items: []string{"Images",
				"Containers",
				"Exit"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}

		switch result {

		case "Images":
			cont = pickImagesAction()

		case "Containers":
			cont = pickContainersAction()

		case "Exit":
			return false
		}
	}
	return true
}

func main() {
	header := `
 ____             _               _   _      _
|  _ \  ___   ___| | _____ _ __  | | | | ___| |_ __   ___ _ __
| | | |/ _ \ / __| |/ / _ \ '__| | |_| |/ _ \ | '_ \ / _ \ '__|
| |_| | (_) | (__|   <  __/ |    |  _  |  __/ | |_) |  __/ |
|____/ \___/ \___|_|\_\___|_|    |_| |_|\___|_| .__/ \___|_|
                                              |_|

                                                by João Condeço
	
	`
	fmt.Print(header)
	cont := true
	for cont {
		cont = pickContext()
		fmt.Printf("\n\n")
	}
}
