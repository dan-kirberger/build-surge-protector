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

	var buildType = getEnv("DRONE_BUILD_EVENT")
	//var tag = getEnv("DRONE_TAG")
	//var branch = getEnv("DRONE_BRANCH")
	//var pr = getEnv("DRONE_PULL_REQUEST")
	//var deployTarget = getEnv("DRONE_DEPLOY_TO")

	fmt.Println("Testing drone token...")
	config := new(oauth2.Config)
	auth := config.Client(
		context.Background(),
		&oauth2.Token{
			AccessToken: token,
		},
	)

	client := drone.NewClient(host, auth)

	user, _ := client.Self()
	fmt.Println("Authenticated as", user.Login)

	var options = drone.ListOptions{Page:0, Size:50}
	var recentBuilds, _ = client.BuildList(orgName, repoName, options)
	fmt.Println("Found", len(recentBuilds), "builds")



	if buildType == "push" {
		fmt.Println("Cancelling old build", recentBuilds[0].Number)
		client.BuildCancel(orgName, repoName, int(recentBuilds[0].Number))
		fmt.Println("Donezo")
	}
}