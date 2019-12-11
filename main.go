package main

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
)

func getEnv(varName string) string {
	value, _ := os.LookupEnv(varName)
	return value
}

func checkEnv(varName string) {
	_, exists := os.LookupEnv(varName)
	if !exists {
		fmt.Println("Must set environment variable: " + varName)
		os.Exit(1)
	}
}

func main() {
	var token = getEnv("DRONE_TOKEN")
	var host  = getEnv("PLUGIN_DRONE_HOST")
	var orgName = getEnv("DRONE_REPO_OWNER")
	var repoName = getEnv("DRONE_REPO_NAME")
	//DRONE_REPO_OWNER
	//DRONE_REPO_NAME
	println("here i am")
	config := new(oauth2.Config)
	auther := config.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)

	// create the drone client with authenticator
	client := drone.NewClient(host, auther)

	// gets the current user
	user, err := client.Self()
	fmt.Println(user, err)

	// gets the named repository information
	repo, err := client.Repo(orgName, repoName)
	fmt.Println(repo, err)
}