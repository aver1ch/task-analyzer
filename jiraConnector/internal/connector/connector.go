package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	configreader "github.com/jiraconnector/internal/configReader"
	"github.com/jiraconnector/internal/structures"
)

type JiraConnector struct {
	cfg    *configreader.JiraConfig
	client *http.Client
}

type JiraConnectorInterface interface {
	GetAllProjects() ([]structures.JiraProject, error)
	GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error)
	GetProjectIssues(projectId string) ([]structures.JiraIssue, error)
}

func NewJiraConnector(config configreader.Config) *JiraConnector {
	return &JiraConnector{
		cfg:    &config.JiraCfg,
		client: &http.Client{},
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

func (con *JiraConnector) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	allProjects, err := con.GetAllProjects()
	if err != nil {
		log.Printf("error while get all projects: %v", err)
		return nil, err
	}

	var pageProjects []structures.JiraProject
	for _, proj := range allProjects {
		if search == "" || containsSearchProject(proj.Name, search) {
			pageProjects = append(pageProjects, proj)
		}
	}

	totalProjects := len(pageProjects)
	start := (page - 1) * limit
	if start >= totalProjects {
		return nil, nil
	}
	end := start + limit
	if end > totalProjects {
		end = totalProjects
	}

	return &structures.ResponseProject{
			Projects: pageProjects[start:end],
			PageInfo: structures.PageInfo{
				PageCount:     int(math.Ceil(float64(totalProjects) / float64(limit))),
				CurrentPage:   page,
				ProjectsCount: totalProjects,
			},
		},
		nil
}

func (con *JiraConnector) GetProjectIssues(projectId string) ([]structures.JiraIssue, error) {

	//get all issues for this project
	totalIssues, err := con.getTotalIssues(projectId)
	if err != nil {
		log.Printf("error while get total issues: %v", err)
		return nil, err
	}

	if totalIssues == 0 {
		return nil, nil
	}

	//create common source = map for results
	var allIssues []structures.JiraIssue
	threadCount := con.cfg.ThreadCount
	issueReq := con.cfg.IssueInOneReq

	//create go routines
	var wg sync.WaitGroup
	var issuesMux sync.Mutex

	for i := 0; i < threadCount; i++ {
		wg.Add(1)
		//find start index for get issues
		issueStart := i*issueReq + 1
		if issueStart > totalIssues {
			issueStart = totalIssues
		}

		go func() {
			defer wg.Done()
			//TODO: should i add count of recquest like in example?

			issues, err := con.getIssuesForOneThread(issueStart, projectId)
			if err != nil {
				log.Printf("error in thread num #%d: %v", i, err)
				//TODO: stop all threads
			}

			issuesMux.Lock()
			defer issuesMux.Unlock()
			allIssues = append(allIssues, issues...)
		}()
	}

	wg.Wait()

	return allIssues, nil
}

func (con *JiraConnector) getIssuesForOneThread(startAt int, projectId string) ([]structures.JiraIssue, error) {
	url := fmt.Sprintf(
		"%s/rest/api/2/search?jql=project=%s&startAt=%d&maxResult=%d",
		con.cfg.Url, projectId, startAt, con.cfg.IssueInOneReq)

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

	return issues.Issues, nil
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

	for {
		resp, err = con.client.Do(req)

		// if everything ok - return resp
		if err == nil && resp.StatusCode < 300 {
			return resp, nil
		}
		time.Sleep(time.Duration(timeSleep))
		timeSleep *= 2

		if timeSleep > con.cfg.MaxSleep {
			break
		}
	}

	// if in cycle we didn't do response - return err
	return nil, err

}

func (con *JiraConnector) getTotalIssues(projectId string) (int, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&maxResults=0", con.cfg.Url, projectId)
	resp, err := con.retryRequest("GET", url)
	if err != nil {
		log.Printf("error while do request: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error while read get projects issues responce: %v", err)
		return 0, err
	}

	var issues structures.JiraIssues
	if err := json.Unmarshal(body, &issues); err != nil {
		log.Printf("error while unmarshal get projects issues body: %v", err)
		return 0, err
	}

	return issues.Total, nil
}

func containsSearchProject(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}
