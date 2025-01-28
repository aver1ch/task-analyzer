package jiraservice

import (
	"log"

	configreader "github.com/jiraconnector/internal/configReader"
	"github.com/jiraconnector/internal/connector"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	dbpusher "github.com/jiraconnector/internal/dbPusher"
	"github.com/jiraconnector/internal/structures"
)

type JiraService struct {
	jiraConnector   connector.JiraConnector
	dataTransformer datatransformer.DataTransformer
	dbPusher        dbpusher.DbPusher
}

func NewJiraService(config configreader.Config, jiraConnector connector.JiraConnector) JiraService {
	return JiraService{
		jiraConnector:   jiraConnector,
		dataTransformer: *datatransformer.NewDataTransformer(),
		dbPusher:        *dbpusher.NewDbPusher(config),
	}
}

func (js JiraService) GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error) {
	return js.jiraConnector.GetProjectsPage(search, limit, page)
}
func (js JiraService) UpdateProjects(projectId string) ([]structures.JiraIssue, error) {
	return js.jiraConnector.GetProjectIssues(projectId)
}

func (js JiraService) PushDataToDb(project string, issues []structures.JiraIssue) error {
	data := js.TransformDataToDb(project, issues)

	if err := js.dbPusher.PushIssues(project, data); err != nil {
		log.Printf("error while push issues: %v", err)
		return err
	}

	return nil

}

func (js JiraService) TransformDataToDb(project string, issues []structures.JiraIssue) []datatransformer.DataTransformer {
	var issuesDb []datatransformer.DataTransformer

	for _, issue := range issues {
		issuesDb = append(issuesDb, *js.dataTransformer.TransformToDbIssueSet(project, issue))
	}

	return issuesDb
}
