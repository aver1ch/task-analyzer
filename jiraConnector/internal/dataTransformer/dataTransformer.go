package datatransformer

import (
	"strconv"
	"time"

	"github.com/jiraconnector/internal/structures"
)

//add changelog
//func TransformStatusDB(jiraIssue structures.JiraIssue) structures.DBStatusChanges {}

func TransformAuthorDB(jiraAuthor structures.User) structures.DBAuthor {
	return structures.DBAuthor{
		Name: jiraAuthor.Name,
	}
}

func TransformProjectDB(jiraProject structures.JiraProject) structures.DBProject {
	return structures.DBProject{
		Title: jiraProject.Name,
	}
}

func TransformIssueDB(jiraIssue structures.JiraIssue) structures.DBIssue {
	id, _ := strconv.Atoi(jiraIssue.Id)
	prjId, _ := strconv.Atoi(jiraIssue.Fields.Project.Id)
	authorId, _ := strconv.Atoi(jiraIssue.Fields.Author.Key)
	assignedId, _ := strconv.Atoi(jiraIssue.Fields.Assignee.Key)

	createdTime, _ := time.Parse("2006-01-002T15:04:05.999-0700", jiraIssue.Fields.CreatedTime)
	closedTime, _ := time.Parse("2006-01-002T15:04:05.999-0700", jiraIssue.Fields.ClosedTime)
	updatedTime, _ := time.Parse("2006-01-002T15:04:05.999-0700", jiraIssue.Fields.UpdatedTime)
	timeSpent, _ := strconv.Atoi(jiraIssue.Fields.TimeSpent)

	return structures.DBIssue{
		Id:          id,
		ProjectId:   prjId,
		AuthorId:    authorId,
		AssigneeId:  assignedId,
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
