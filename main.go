package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	dir = flag.String("directory", "", "The root folder of the git project you would like to update the submodules of")

	remote = flag.String("remote", "origin", "The remote of the git repository which you want to push to")
	branch = flag.String("branch", "master", "The branch of the remote git repository which you want to push to")

	forkOwner = flag.String("fork-owner", "", "The owner of the repository the Pull Request will be sent from")
	forkRepo  = flag.String("fork-repo", "", "The repository the Pull Request will be sent from of")

	owner    = flag.String("owner", "", "The owner of the repository which the Pull Request will be sent to")
	repo     = flag.String("repo", "", "The repository which the Pull Request will be sent to")
	toBranch = flag.String("to-branch", "master", "The branch which you would like your changes merged into")

	token = flag.String("token", "", "The API token which has access to both the fork and main repositories")

	updateSubmodules = exec.Command("git", "submodule", "update", "--remote", "--merge")
	checkDiff        = exec.Command("git", "status", "--porcelain")
	addChanges       = exec.Command("git", "add", "-A")
	commitChanges    = exec.Command("git", "commit", "-m Update submodules")
	pushChanges      *exec.Cmd // pushChanges is set later as it relies on the remote and branch flags. In the form: `git push <remote> <branch>`
)

func main() {
	flag.Parse()

	pushChanges = exec.Command("git", "push", *remote, *branch)

	fmt.Println("Welcome... To the submodule genie!")
	fmt.Println("I'm going to update the submodules of the repository in the directory", *dir)

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Println("Couldn't create an absolute filename", err)
		return
	}

	err = runGit(absDir)
	if err != nil {
		fmt.Println("Couldn't finish the git operations:", err)
		return
	}

	err = makePR()
	if err != nil {
		fmt.Println("Couldn't make PR:", err)
		return
	}

	fmt.Println("I've been waiting a while to say this: My work here... is done!")
}

func makePR() error {
	pr := pullRequest{
		Title: "Update submodules",
		Body:  "Hey! I'm just your friendly neighbourhood bot, updating your submodules.",
		Head:  fmt.Sprintf("%v:%v", *forkOwner, *branch),
		Base:  *toBranch,
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(pr)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls?access_token=%v", *owner, *repo, *token)
	resp, err := http.Post(endpoint, "application/json", &buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func runGit(dir string) error {
	updateSubmodules.Dir = dir
	checkDiff.Dir = dir
	addChanges.Dir = dir
	commitChanges.Dir = dir
	pushChanges.Dir = dir

	updateSubmodules.Stdin = os.Stdin
	updateSubmodules.Stdout = os.Stdout
	updateSubmodules.Stderr = os.Stderr

	pushChanges.Stdin = os.Stdin
	pushChanges.Stdout = os.Stdout
	pushChanges.Stderr = os.Stderr

	err := updateSubmodules.Run()
	if err != nil {
		return err
	}

	out, err := checkDiff.Output()
	if err != nil {
		return err
	}

	fmt.Printf("`%s`", out)

	if string(out) == "" {
		return err
	}

	err = addChanges.Run()
	if err != nil {
		return err
	}

	err = commitChanges.Run()
	if err != nil {
		return err
	}

	err = pushChanges.Run()
	if err != nil {
		return err
	}

	return nil
}

type pullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}
