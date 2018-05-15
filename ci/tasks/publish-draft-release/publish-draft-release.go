package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	gitHubUrl                      string = "https://api.github.com/repos/"
	gitHubCredentialsVariable      string = "GITHUB_ACCESS_KEY"
	debugMessagesVariable          string = "DEBUG_MESSAGES_ENABLED"
	getReleasesMethodType          string = "GET"
	editReleaseMethodType          string = "PATCH"
	defaultVersion                 string = "latest"
	defaultBranch                  string = "master"
	minimumNumberOfCommandLineArgs int    = 3
)

var (
	owner         string
	repo          string
	branch        string
	version       string
	description   string
	credentials   string = os.Getenv(gitHubCredentialsVariable)
	debugMessages bool
)

type IdAndTag struct {
	Id      int64  `json:"id"`
	TagName string `json:"tag_name"`
}

type Release struct {
	TagName    string `json:"tag_name"`
	Branch     string `json:"target_commitish"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

func getReleaseApiUrl() string {
	return fmt.Sprintf("%s%s/%s/releases", gitHubUrl, owner, repo)
}

func getEditReleaseApiUrl(id int64) string {
	return fmt.Sprintf("%s/%d", getReleaseApiUrl(), id)
}

func sendRequest(methodType string, url string, body io.Reader, bodySize int64) (int, []byte, error) {
	req, err := http.NewRequest(methodType, url, body)

	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", credentials))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.ContentLength = bodySize

	if debugMessages {
		if body == nil {
			os.Stdout.WriteString(fmt.Sprintf("Sending request to %s: %+v\n", url, req))
		} else {
			os.Stdout.WriteString(fmt.Sprintf("Sending request to %s with data %s: %+v\n", url, body, req))
		}
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, bodyBytes, err
}

func getReleaseId() (int64, error) {
	url := getReleaseApiUrl()
	code, bodyBytes, err := sendRequest(getReleasesMethodType, url, bytes.NewBuffer([]byte{}), 0)

	if err != nil {
		return 0, err
	}

	if code == http.StatusOK {
		fmt.Printf(fmt.Sprintf("Got releases on branch %s: %s\n", branch, string(bodyBytes)))
	} else {
		return 0, fmt.Errorf("error getting releases for version %s with response code %d", url, code)
	}

	var releases []IdAndTag
	json.Unmarshal(bodyBytes, &releases)
	os.Stdout.WriteString(fmt.Sprintf("Response from %s: %v\n", url, releases))

	for _, release := range releases {
		if release.TagName == version {
			return release.Id, err
		}
	}

	return 0, fmt.Errorf("No ID found from %s using version %s", url, version)
}

func editRelease(id int64, body io.Reader, bodySize int64) (int, []byte, error) {
	return sendRequest(editReleaseMethodType, getEditReleaseApiUrl(id), body, bodySize)
}

func publishDraftRelease() error {
	id, err := getReleaseId()

	if err != nil {
		return fmt.Errorf("error getting releases: %s", err)
	}

	if id == 0 {
		return fmt.Errorf("error getting releases - no ID found")
	}

	release := Release{
		TagName:    version,
		Branch:     branch,
		Name:       version,
		Body:       description,
		Draft:      false,
		Prerelease: false,
	}

	releaseData, err := json.Marshal(release)

	if err != nil {
		return fmt.Errorf("error setting JSON data %s when publishing draft release to %s due to %s", releaseData, getReleaseApiUrl(), err)
	}

	fmt.Printf(fmt.Sprintf("Publishing draft release of version %s on branch %s\n", version, branch))

	releaseBuffer := bytes.NewBuffer(releaseData)

	code, _, err := editRelease(id, releaseBuffer, int64(releaseBuffer.Len()))

	if err != nil {
		return fmt.Errorf("error publishing draft release to %s", getReleaseApiUrl())
	}

	if code == http.StatusOK {
		fmt.Printf(fmt.Sprintf("Published draft release of id %d for version %s on branch %s\n", id, version, branch))
	} else {
		return fmt.Errorf("error publishing draft release to %s with response code %d", getReleaseApiUrl(), code)
	}

	return nil
}

func main() {
	if credentials == "" {
		os.Stderr.WriteString("Must provide GitHub credentials via GITHUB_ACCESS_KEY\n")
		os.Exit(1)
	}

	numberOfCommandLineArgs := len(os.Args)

	if numberOfCommandLineArgs < minimumNumberOfCommandLineArgs {
		os.Stderr.WriteString(fmt.Sprintf("Only found %d arguments - Must provide owner and repo\n", numberOfCommandLineArgs))
		os.Exit(1)
	}

	owner = os.Args[1]
	repo = os.Args[2]

	if owner == "" {
		os.Stderr.WriteString("Must provide owner as the first argument\n")
		os.Exit(1)
	}

	if repo == "" {
		os.Stderr.WriteString("Must provide repo as the second argument\n")
		os.Exit(1)
	}

	if numberOfCommandLineArgs < 4 {
		version = defaultVersion
	} else {
		version = os.Args[3]
	}

	if numberOfCommandLineArgs < 5 {
		branch = defaultBranch
	} else {
		branch = os.Args[4]
	}

	if numberOfCommandLineArgs < 6 {
		description = version
	} else {
		description = os.Args[5]
	}

	var err error

	debugMessages, err = strconv.ParseBool(os.Getenv(debugMessagesVariable))

	if err != nil {
		debugMessages = false
	}

	err = publishDraftRelease()

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Failed to publish draft release - %s\n", err))
	}
}
