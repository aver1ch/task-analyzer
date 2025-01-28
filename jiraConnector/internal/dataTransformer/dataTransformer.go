package datatransformer

import (
	"strconv"
	"time"

	"github.com/jiraconnector/internal/structures"
)

type DataTransformer struct {
	Project  structures.DBProject
	Issue    structures.DBIssue
	Author   structures.DBAuthor
	Assignee structures.DBAuthor
}

func NewDataTransformer() *DataTransformer {
	return &DataTransformer{}
}

//add changelog
//func TransformStatusDB(jiraIssue structures.JiraIssue) structures.DBStatusChanges {}

func (dt *DataTransformer) TransformAuthorDB(jiraAuthor structures.User) structures.DBAuthor {
	return structures.DBAuthor{
		Name: jiraAuthor.Name,
	}
}

func (dt *DataTransformer) TransformProjectDB(jiraProject structures.JiraProject) structures.DBProject {
	return structures.DBProject{
		Title: jiraProject.Name,
	}
}

func (dt *DataTransformer) TransformIssueDB(jiraIssue structures.JiraIssue) structures.DBIssue {
	layout := "2006-01-02T15:04:05.000-0700"
	timeSpent, _ := strconv.Atoi(jiraIssue.Fields.TimeSpent)
	createdTime, _ := time.Parse(layout, jiraIssue.Fields.CreatedTime)
	updatedTime, _ := time.Parse(layout, jiraIssue.Fields.UpdatedTime)
	closedTime, _ := time.Parse(layout, jiraIssue.Fields.ClosedTime)

	return structures.DBIssue{
		Key:         jiraIssue.Key,
		Summary:     jiraIssue.Fields.Summary,
		Description: jiraIssue.Fields.Description,
		Type:        jiraIssue.Fields.Type.Description,
		Priority:    jiraIssue.Fields.Project.Name,
		Status:      jiraIssue.Fields.Status.Name,
		CreatedTime: createdTime,
		ClosedTime:  closedTime,
		UpdatedTime: updatedTime,
		TimeSpent:   timeSpent,
	}
}

func (dt *DataTransformer) TransformToDbIssueSet(projectName string, jiraIssue structures.JiraIssue) *DataTransformer {
	return &DataTransformer{
		Project:  structures.DBProject{Title: projectName},
		Issue:    dt.TransformIssueDB(jiraIssue),
		Author:   dt.TransformAuthorDB(jiraIssue.Fields.Author),
		Assignee: dt.TransformAuthorDB(jiraIssue.Fields.Assignee),
	}
}
