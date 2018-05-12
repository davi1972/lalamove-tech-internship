package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version

	var compareSlice [][]string
	// Remove versions that are lower than min and sort
	semver.Sort(releases)

	for _, versions := range releases {
		if versions.Major < minVersion.Major || versions.Minor < minVersion.Minor {
			releases = releases[1:]
		} else {
			compareSlice = append(compareSlice, strings.Split(versions.String(), "."))
		}
	}
	// Create a dummy variable for sorting
	point := compareSlice[0]
	prevVer := releases[0]
	for index, versions := range releases {
		for i := 0; i < (len(compareSlice[index])); i++ {
			// Convert the strings to int again
			a, err := strconv.ParseInt((compareSlice[index][i]), 10, 64)
			b, err := strconv.ParseInt(point[i], 10, 64)
			if err != nil {
				break
			}
			if a > b {
				if i < 2 {
					versionSlice = append(versionSlice, prevVer)
				}
				break
			}
		}
		point = compareSlice[index]
		prevVer = versions
		if index == len(releases)-1 {
			versionSlice = append(versionSlice, prevVer)
		}
	}
	for i := len(versionSlice)/2 - 1; i >= 0; i-- {
		opp := len(versionSlice) - 1 - i
		versionSlice[i], versionSlice[opp] = versionSlice[opp], versionSlice[i]
	}
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	return versionSlice
}

func parseFile(args string) [][]string {
	data, err := ioutil.ReadFile(args)
	if err != nil {
		panic(err)
	}
	dataString := string(data)
	dataStringArr := strings.Split(dataString, "\n")

	var toCheck [][]string
	i := dataStringArr[1:]
	for _, value := range i {
		arr := strings.Split(value, ",")
		toCheck = append(toCheck, arr)
	}
	return toCheck
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	cmd := os.Args
	commands := parseFile(cmd[1])
	fmt.Println(commands)

	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	for _, command := range commands {
		repo := strings.Split(command[0], "/")
		releases, _, err := client.Repositories.ListReleases(ctx, repo[0], repo[1], opt)
		if err != nil {
			panic(err) // is this really a good way?
		}
		minVersion := semver.New(command[1])
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of %s/%s: %s \n", repo[0], repo[1], versionSlice)
	}
}
