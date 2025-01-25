package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	configreader "github.com/jiraconnector/internal/configReader"
	structures "github.com/jiraconnector/internal/structures"
)

type JiraConnector struct {
	maxRetry int
	cfg      *configreader.JiraConfig
	client   *http.Client
}

func NewJiraConnector(config configreader.Config, maxRetry int) *JiraConnector {
	return &JiraConnector{
		maxRetry: maxRetry,
		cfg:      &config.JiraCfg,
		client:   &http.Client{},
	}
}

func (con *JiraConnector) GetAllProjects() ([]structures.JiraProject, error) {
	url := fmt.Sprintf("%s/rest/api/2/project", con.cfg.Url)

	resp, err := con.retryRequest("GET", url)
	if err != nil {
		log.Printf("error while do request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while read get projects responce: %v", err)
		return nil, err
	}

	var projects []structures.JiraProject
	if err = json.Unmarshal(body, &projects); err != nil {
		log.Printf("error while unmarshal get projects body: %v", err)
		return nil, err
	}

	return projects, nil
}

func (con *JiraConnector) GetProjectIssues(projectId string) (*structures.JiraIssues, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s", con.cfg.Url, projectId)
	resp, err := con.retryRequest("GET", url)
	if err != nil {
		log.Printf("error while do request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while read get projects issues responce: %v", err)
		return nil, err
	}

	var issues structures.JiraIssues
	if err := json.Unmarshal(body, &issues); err != nil {
		log.Printf("error while unmarshal get projects issues body: %v", err)
		return nil, err
	}

	return &issues, nil
}

func (con *JiraConnector) retryRequest(method, url string) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	timeSleep := con.cfg.MinSleep

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("error while make request: %v", err)
		return nil, err
	}

	for i := 0; i < con.maxRetry; i++ {
		resp, err = con.client.Do(req)

		// if everything ok - return resp
		if err == nil && resp.StatusCode < 300 {
			return resp, nil
		}
		time.Sleep(time.Duration(timeSleep))
		timeSleep *= 2

		if timeSleep > con.cfg.MaxSleep {
			timeSleep = con.cfg.MaxSleep
		}
	}

	// if in cycle we didn't do response - return err
	return nil, err

}
