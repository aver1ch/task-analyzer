package dbpusher

import (
	"database/sql"
	"fmt"
	"log"

	configreader "github.com/jiraconnector/internal/configReader"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
	_ "github.com/lib/pq"
)

type DbPusher struct {
	db *sql.DB
}

func NewDbPusher(cfg configreader.Config) *DbPusher {
	connStr := buildConnectionstring(&cfg.DBCfg)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("error open database: %v", err)
		return nil
	}

	return &DbPusher{
		db: db,
	}
}

func (dbp *DbPusher) Close() {
	dbp.db.Close()
}

func (dbp *DbPusher) PushProject(project structures.DBProject) (int, error) {
	var projectId int
	query := "INSERT INTO project (title) VALUES ($1) RETURNIND id"

	if err := dbp.db.QueryRow(query, project.Title).Scan(&projectId); err != nil {
		return 0, fmt.Errorf("error insert project %s: %v", project.Title, err)
	}

	log.Printf("Project %s was added", project.Title)
	return projectId, nil
}

func (dbp *DbPusher) getProjectId(projectName string) (int, error) {
	var projectId int
	var err error
	query := "SELECT id FROM project WHERE title=$1"

	_ = dbp.db.QueryRow(query, projectName).Scan(&projectId)
	if projectId == 0 {
		projectId, err = dbp.PushProject(structures.DBProject{Title: projectName})
		if err != nil {
			return 0, fmt.Errorf("error select project %s id: %v", projectName, err)
		}
	}

	return projectId, nil
}

func (dbp *DbPusher) PushProjects(projects []structures.DBProject) error {
	tx, err := dbp.db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin error: %v", err)
	}

	for _, project := range projects {
		_, err = dbp.PushProject(project)
		if err != nil {
			log.Printf("error while push project %s. Rollback.", project.Title)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error : %v", err)
	}
	log.Println("All projects were saved")
	return nil

}

func (dbp *DbPusher) PushAuthor(author structures.DBAuthor) (int, error) {
	var authorId int
	query := "INSERT INTO author (name) VALUES ($1) RETURNING id"

	if err := dbp.db.QueryRow(query, author.Name).Scan(&authorId); err != nil {
		return 0, fmt.Errorf("error insert author %s: %v", author.Name, err)
	}
	log.Printf("Author %s was added", author.Name)
	return authorId, nil

}

func (dbp *DbPusher) getAuthorId(author structures.DBAuthor) (int, error) {
	var authorId int
	var err error
	query := "SELECT id FROM author WHERE name=$1"

	_ = dbp.db.QueryRow(query, author.Name).Scan(&authorId)
	if authorId == 0 {
		authorId, err = dbp.PushAuthor(author)
		if err != nil {
			return 0, fmt.Errorf("error select author %s id: %v", author.Name, err)
		}
	}

	return authorId, nil
}

func (dbp *DbPusher) PushIssue(projectName string, issue datatransformer.DataTransformer) (int, error) {
	projectId, err := dbp.getProjectId(projectName)
	if err != nil {
		return 0, err
	}

	authorId, err := dbp.getAuthorId(issue.Author)
	if err != nil {
		return 0, err
	}

	assegneeId, err := dbp.getAuthorId(issue.Assignee)
	if err != nil {
		return 0, err
	}

	query := `
	INSER INTO issue 
		(projectId, authorId, assigneeId, key, summary, description, type, priority, status, createdTime, closedTime, updatedTume, timeSpent)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	var issueId int
	iss := issue.Issue
	iss.ProjectId = projectId
	iss.AuthorId = authorId
	iss.AssigneeId = assegneeId

	if err := dbp.db.QueryRow(
		query, iss.ProjectId, iss.AuthorId, iss.AssigneeId,
		iss.Key, iss.Summary, iss.Description, iss.Type,
		iss.Priority, iss.Status, iss.CreatedTime,
		iss.ClosedTime, iss.UpdatedTime, iss.TimeSpent).Scan(&issueId); err != nil {
		return 0, err
	}

	return issueId, nil
}

func (dbp *DbPusher) PushIssues(projectName string, issues []datatransformer.DataTransformer) error {
	tx, err := dbp.db.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin error: %v", err)
	}

	for _, issue := range issues {
		_, err := dbp.PushIssue(projectName, issue)
		if err != nil {
			log.Printf("error while push issues for project %s. Rollback.", projectName)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error : %v", err)
	}
	log.Println("All issues were saved")
	return nil
}

func buildConnectionstring(cfg *configreader.DBConfig) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
}
