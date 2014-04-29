package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb"
	"log"
	"os"
	"os/exec"
	"strings"
)

var Repos []GHRepo

func main() {
	Username := flag.String("username", "", "The username that you are targetting")
	update := flag.Bool("update", false, "rather than get, update packages")
	flag.Parse()

	if os.Getenv("GOPATH") == "" {
		fmt.Println("You don't have a GOPATH set, I don't know where to clone to! Please set one.")
		os.Exit(1)
	}

	if *Username == "" {
		fmt.Println("Please give a username to auto clone all their go stars to the gopath to.")
		os.Exit(1)
	}

	fmt.Print("Grabbing all stars off that user")
	GHUrl := fmt.Sprintf("https://api.github.com/users/%s/starred", *Username)
	s, e := ExpectGithubToBreak(GHUrl)
	if e != nil {
		fmt.Println("Cannot get the first set, Not going to attempt to get others.")
	}

	Repos = make([]GHRepo, 0)
	CastData := make([]GHRepo, 0)
	e = json.Unmarshal([]byte(s), &CastData)
	Repos = FilterForGoRepo(CastData, Repos)
	if e != nil {
		fmt.Println("Cannot decode the first set, Not going to attempt to get others.")
	}
	var PageCount int = 2
	var TripCount int = 0
	for {
		GHUrl = fmt.Sprintf("https://api.github.com/users/%s/starred?page=%d", *Username, PageCount)
		s, e = ExpectGithubToBreak(GHUrl)
		if e == nil && TripCount < 2 {
			CastData = make([]GHRepo, 0)
			e = json.Unmarshal([]byte(s), &CastData)
			if len(CastData) == 0 {
				// We have found it all!
				break
			}
			Repos = FilterForGoRepo(CastData, Repos)
			if e != nil {
				fmt.Println("Cannot decode the first set, Not going to attempt to get others.")
			}
			PageCount++
			fmt.Print(".")
		} else {
			TripCount++
			if TripCount < 2 {
				fmt.Println("API errors stopped the program from running.", e.Error())
			}
		}
	}

	bar := pb.StartNew(len(Repos))
	for _, repo := range Repos {
		GoGet(strings.Replace(repo.HtmlURL, "https://", "", -1), update)
		bar.Increment()
	}
}

func GoGet(url string, update *bool) {
	buf := make([]byte, 1024)
	var cmd *exec.Cmd
	if *update {
		cmd = exec.Command("go", "get", "-u", url)
	} else {
		cmd = exec.Command("go", "get", url)
	}
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	for {
		_, e := stdout.Read(buf)
		if e != nil {
			break
		}
	}
}
