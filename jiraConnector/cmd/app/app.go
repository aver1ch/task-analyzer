package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	jirahandlers "github.com/jiraconnector/internal/apiJiraConnector/jiraHandlers"
	jiraservice "github.com/jiraconnector/internal/apiJiraConnector/jiraService"
	configreader "github.com/jiraconnector/internal/configReader"
	"github.com/jiraconnector/internal/connector"
)

type JiraApp struct {
	server        *http.Server
	jiraConnector *connector.JiraConnector
}

func NewApp(cfg configreader.Config) (*JiraApp, error) {
	con := connector.NewJiraConnector(cfg)
	log.Println("created jira connection")

	service, err := jiraservice.NewJiraService(cfg, *con)
	if err != nil {
		ansErr := fmt.Errorf("error create service: %w", err)
		log.Println(ansErr)
		return nil, ansErr
	}
	log.Println("created jira service")

	router := mux.NewRouter()
	jiraHandler := jirahandlers.NewHandler(service, router)
	log.Println("created jira handlers")

	server := &http.Server{
		Addr:    cfg.ServerCfg.Port,
		Handler: jiraHandler,
	}

	return &JiraApp{
		server:        server,
		jiraConnector: con,
	}, nil
}

func (a *JiraApp) Run() error {
	log.Println("run app")
	return fmt.Errorf("run app err: %v", a.server.ListenAndServe())
}

func (a *JiraApp) Close() {
	log.Println("close app")
	a.server.Close()
}
