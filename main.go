package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	dir = flag.String("directory", "", "The root folder of the git project you would like to update the submodules of")

	remote = flag.String("remote", "origin", "The remote of the git repository which you want to push to")
	branch = flag.String("branch", "master", "The branch of the remote git repository which you want to push to")

	forkOwner = flag.String("fork-owner", "", "The repository which will be used to send a Pull Request GitHub")
	repo      = flag.String("repo", "", "The repository name of the repository on GitHub")

	updateSubmodules = exec.Command("git", "submodule", "update", "--remote", "--merge")
	checkDiff        = exec.Command("git", "status", "--porcelain")
	addChanges       = exec.Command("git", "add", "-A")
	commitChanges    = exec.Command("git", "commit", "-m Update submodules")
	pushChanges      *exec.Cmd // pushChanges is set later as it relies on the remote and branch flags. takes the form: `git push <remote> <branch>`
)

func main() {
	flag.Parse()

	pushChanges = exec.Command("git", "push", *remote, *branch)

	fmt.Println("Welcome... To the submodule genie!")
	fmt.Println("I'm going to update the submodules of the repository in the directory", *dir)

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("This is the updated absolute directory:", absDir)

	updateSubmodules.Dir = absDir
	checkDiff.Dir = absDir
	addChanges.Dir = absDir
	commitChanges.Dir = absDir
	pushChanges.Dir = absDir

	updateSubmodules.Stdin = os.Stdin
	updateSubmodules.Stdout = os.Stdout
	updateSubmodules.Stderr = os.Stderr

	pushChanges.Stdin = os.Stdin
	pushChanges.Stdout = os.Stdout
	pushChanges.Stderr = os.Stderr

	err = updateSubmodules.Run()
	if err != nil {
		fmt.Println("Could not update submodules:", err)
		return
	}

	out, err := checkDiff.Output()
	if err != nil {
		fmt.Println("Could not check diff:", err)
		return
	}

	fmt.Printf("`%s`", out)

	if string(out) == "" {
		fmt.Println("No updates were made, exiting now")
		return
	}

	err = addChanges.Run()
	if err != nil {
		fmt.Println("Could not add changes:", err)
		return
	}

	err = commitChanges.Run()
	if err != nil {
		fmt.Println("Could not commit changes:", err)
		return
	}

	err = pushChanges.Run()
	if err != nil {
		fmt.Println("Could not push commit:", err)
		return
	}

	pr := pullRequest{
		Title: "Update submodules",
		Body:  "Update submodules",
		Head:  "",
		Base:  "",
	}

	var buf bytes.Buffer

	info := json.NewEncoder(&buf).Encode()
}

type pullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}
