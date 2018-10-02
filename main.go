package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ints []int

func (i *ints) String() string {
	return fmt.Sprintf("%d", *i)
}
func (i *ints) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		*i = append(*i, -1)
	} else {
		*i = append(*i, tmp)
	}
	return nil
}

var (
	teams  ints
	source = flag.String("source", "", "github source organisation")
	target = flag.String("target", "", "github target organization")
	token  = flag.String("token", "", "oauth token")
)

type repository struct {
	Name string `json:"name"`
}

func main() {
	flag.Var(&teams, "teams", "optional team ids")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		return
	}
	fmt.Println("transfering repositories (see https://developer.github.com/v3/repos/#transfer-a-repository)...")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/orgs/%s/repos", *source), nil)
	if err != nil {
		log.Fatalf("failed to create github api request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", *token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("github api call to retrieve repos from organization failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("github api call to retrieve repos from organization failed: %s", resp.Status)
	}

	var repos []repository
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		log.Fatalf("failed to unmarshal github api response: %v", err)
	}

	// TODO: support pagination in api call
	for _, r := range repos {
		type body struct {
			Owner string `json:"new_owner"`
			Teams []int  `json:"team_ids"`
		}
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(&body{
			Owner: *target,
			Teams: teams,
		})
		if err != nil {
			log.Fatalf("failed to encode github api request: %v", err)
		}

		req, err := http.NewRequest("POST",
			fmt.Sprintf("https://api.github.com/repos/%s/%s/transfer", *source, r.Name),
			b)
		if err != nil {
			log.Fatalf("failed to create github api request: %v", err)
		}
		req.Header.Set("Accept", "application/vnd.github.nightshade-preview+json")
		req.Header.Set("Authorization", fmt.Sprintf("token %s", *token))
		log.Printf("transfering %s from %s to %s\n", r.Name, *source, *target)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("github api call to transfer repo failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusAccepted {
			log.Fatalf("transfering %s failed: %s", r.Name, resp.Status)
		}
	}
	fmt.Println("...done")
}
