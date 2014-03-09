package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func ExpectGithubToBreak(url string) (out string, e error) {
	r, e := http.Get(url)
	if e != nil {
		fmt.Print("!")
	}

	if r.StatusCode != 200 {
		if r.StatusCode == 403 {
			fmt.Println("Github Rate Limit reached.")
			rlreset := r.Header.Get("X-RateLimit-Reset")
			if rlreset != "" {
				i, e := strconv.ParseInt(rlreset, 10, 64)
				if e == nil {
					// I don't really understand how this would happen but its worth checking.
					fmt.Printf("You can expect the rate limit to reset on %d\n", i)
				}
			}
			os.Exit(1)
		}
		return "", fmt.Errorf("Cannot get GH page")
	} else {
		return ioutil.ReadAll(r.Body), e
	}
}
