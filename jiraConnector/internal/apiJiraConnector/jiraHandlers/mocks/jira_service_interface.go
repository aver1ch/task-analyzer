// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	datatransformer "github.com/jiraconnector/internal/dataTransformer"

	mock "github.com/stretchr/testify/mock"

	structures "github.com/jiraconnector/internal/structures"
)

// JiraServiceInterface is an autogenerated mock type for the JiraServiceInterface type
type JiraServiceInterface struct {
	mock.Mock
}

// GetProjectsPage provides a mock function with given fields: search, limit, page
func (_m *JiraServiceInterface) GetProjectsPage(search string, limit int, page int) (*structures.ResponseProject, error) {
	ret := _m.Called(search, limit, page)

	if len(ret) == 0 {
		panic("no return value specified for GetProjectsPage")
	}

	var r0 *structures.ResponseProject
	var r1 error
	if rf, ok := ret.Get(0).(func(string, int, int) (*structures.ResponseProject, error)); ok {
		return rf(search, limit, page)
	}
	if rf, ok := ret.Get(0).(func(string, int, int) *structures.ResponseProject); ok {
		r0 = rf(search, limit, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*structures.ResponseProject)
		}
	}

	if rf, ok := ret.Get(1).(func(string, int, int) error); ok {
		r1 = rf(search, limit, page)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PushDataToDb provides a mock function with given fields: project, issues
func (_m *JiraServiceInterface) PushDataToDb(project string, issues []structures.JiraIssue) error {
	ret := _m.Called(project, issues)

	if len(ret) == 0 {
		panic("no return value specified for PushDataToDb")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []structures.JiraIssue) error); ok {
		r0 = rf(project, issues)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransformDataToDb provides a mock function with given fields: project, issues
func (_m *JiraServiceInterface) TransformDataToDb(project string, issues []structures.JiraIssue) []datatransformer.DataTransformer {
	ret := _m.Called(project, issues)

	if len(ret) == 0 {
		panic("no return value specified for TransformDataToDb")
	}

	var r0 []datatransformer.DataTransformer
	if rf, ok := ret.Get(0).(func(string, []structures.JiraIssue) []datatransformer.DataTransformer); ok {
		r0 = rf(project, issues)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]datatransformer.DataTransformer)
		}
	}

	return r0
}

// UpdateProjects provides a mock function with given fields: projectId
func (_m *JiraServiceInterface) UpdateProjects(projectId string) ([]structures.JiraIssue, error) {
	ret := _m.Called(projectId)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProjects")
	}

	var r0 []structures.JiraIssue
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]structures.JiraIssue, error)); ok {
		return rf(projectId)
	}
	if rf, ok := ret.Get(0).(func(string) []structures.JiraIssue); ok {
		r0 = rf(projectId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]structures.JiraIssue)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(projectId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJiraServiceInterface creates a new instance of JiraServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJiraServiceInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *JiraServiceInterface {
	mock := &JiraServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
