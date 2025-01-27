package jiraservice

import (
	configreader "github.com/jiraconnector/internal/configReader"
	"github.com/jiraconnector/internal/connector"
	"github.com/jiraconnector/internal/structures"
)

type JiraService struct {
	jiraConnector connector.JiraConnector
}

func NewJiraService(config configreader.Config, jiraConnector connector.JiraConnector) JiraService {
	return JiraService{
		jiraConnector: jiraConnector,
	}
}

func (js *JiraService) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	return js.jiraConnector.GetProjectsPage(search, limit, page)
}
func (js *JiraService) UpdateProjects(projectId string) ([]structures.JiraIssue, error) {
	return js.jiraConnector.GetProjectIssues(projectId)
}
