package main

import (
	"fmt"
	"github.com/drone/drone-go/drone"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
	"strconv"
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
	var host = getEnv("PLUGIN_DRONE_HOST")
	var orgName = getEnv("DRONE_REPO_OWNER")
	var repoName = getEnv("DRONE_REPO_NAME")

	var buildType = getEnv("DRONE_BUILD_EVENT")
	var tag = getEnv("DRONE_TAG")
	var branch = getEnv("DRONE_BRANCH")
	var pr = getEnv("DRONE_PULL_REQUEST")
	var deployTarget = getEnv("DRONE_DEPLOY_TO")

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

	var options = drone.ListOptions{Page: 0, Size: 50}
	var recentBuilds, _ = client.BuildList(orgName, repoName, options)

	switch buildType {
	case "push":
		test := func(b *drone.Build) bool { return b.Event == "push" && b.Ref == "refs/heads/" + branch }
		buildsToKill := filter(recentBuilds, test)
		killBuilds(client, orgName, repoName, buildsToKill)
	case "pull_request":
		test := func(b *drone.Build) bool { return b.Event == "pull_request" && b.Ref == "refs/pull/"+pr+"/head" }
		buildsToKill := filter(recentBuilds, test)
		killBuilds(client, orgName, repoName, buildsToKill)
	case "tag":
		test := func(b *drone.Build) bool { return b.Event == "tag" && b.Ref == "refs/tags/"+tag }
		buildsToKill := filter(recentBuilds, test)
		killBuilds(client, orgName, repoName, buildsToKill)
	case "deployment":
		test := func(b *drone.Build) bool { return b.Event == "deployment" && b.Deploy == deployTarget }
		buildsToKill := filter(recentBuilds, test)
		killBuilds(client, orgName, repoName, buildsToKill)
	default:
		fmt.Println("Unknown build type", buildType)
	}

	fmt.Println("Donezo")
}

func filter(builds []*drone.Build, test func(*drone.Build) bool) (ret []int64) {
	for _, build := range builds {
		if test(build) && isRunning(build) && isKillable(build) {
			ret = append(ret, build.Number)
		}
	}
	return
}

func killBuilds(client drone.Client, org string, repo string, buildNumbers []int64) {
	fmt.Println("Cancelling", len(buildNumbers), "builds")
	for _, b := range buildNumbers {
		fmt.Println("Killing build", b)
		err := client.BuildCancel(org, repo, int(b))
		if err != nil {
			fmt.Println("Failed to cancel", org, repo, b)
		}
	}
}

func isRunning(build *drone.Build) bool {
	return build.Status == "running" || build.Status == "pending"
}

func isKillable(build *drone.Build) bool {
	var currentBuildNumberEnv = getEnv("DRONE_BUILD_NUMBER")
	var currentBuildNumber, _ = strconv.Atoi(currentBuildNumberEnv)

	return int(build.Number) < currentBuildNumber
}
