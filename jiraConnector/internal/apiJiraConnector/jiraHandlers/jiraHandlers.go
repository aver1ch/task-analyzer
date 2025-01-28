package jirahandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	datatransformer "github.com/jiraconnector/internal/dataTransformer"
	"github.com/jiraconnector/internal/structures"
)

type JiraServiceInterface interface {
	GetProjectsPage(search string, limit, page int) (*structures.ResponseProject, error)
	UpdateProjects(projectId string) ([]structures.JiraIssue, error)

	PushDataToDb(project string, issues []structures.JiraIssue) error
	TransformDataToDb(project string, issues []structures.JiraIssue) []datatransformer.DataTransformer
}

type handler struct {
	service JiraServiceInterface
}

func NewHandler(service JiraServiceInterface, router *mux.Router) *mux.Router {
	h := handler{service: service}

	router.HandleFunc("/projects", h.projects).Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/updateProject", h.updateProject).Methods(http.MethodOptions, http.MethodPost)

	return router
}

func (h *handler) projects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	limit, page, search, err := getProjectParams(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest) //TODO: norm errors
		return
	}

	projects, err := h.service.GetProjectsPage(search, limit, page)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) //TODO: norm errors
		return
	}

	if err = json.NewEncoder(w).Encode(projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //TODO: norm errors
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h *handler) updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	project := vars["project"]
	if project == "" {
		//log.Println(err)
		//http.Error(w, err.Error(), http.StatusBadRequest) //TODO: norm errors
		return
	}

	issues, err := h.service.UpdateProjects(project)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError) //TODO: norm errors
		return
	}

	h.service.PushDataToDb(project, issues)

	w.WriteHeader(http.StatusOK)
}

func getProjectParams(r *http.Request) (int, int, string, error) {
	var err error
	limit := 20
	page := 1
	search := ""

	vars := mux.Vars(r)

	if vars["limit"] != "" {
		limit, err = strconv.Atoi(vars["limit"])
		if err != nil {
			log.Printf("incorrect limit param: %v", err)
			return 0, 0, "", err
		}
	}

	if vars["page"] != "" {
		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			log.Printf("incorrect page param: %v", err)
			return 0, 0, "", err
		}
	}

	if vars["search"] != "" {
		search = vars["search"]
	}

	return limit, page, search, nil
}
