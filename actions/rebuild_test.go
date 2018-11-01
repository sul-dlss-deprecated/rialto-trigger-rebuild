package actions

import (
	"strings"
	"testing"

	"github.com/knakk/sparql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sul-dlss/rialto-derivatives/models"
	"github.com/sul-dlss/rialto-trigger-rebuild/runtime"
)

// MockedWriter is a mocked object that implements the Writer interface
type MockedWriter struct {
	mock.Mock
}

func (f *MockedWriter) Add(docs []models.Resource) error {
	return nil
}

func (f *MockedWriter) RemoveAll() error {
	return nil
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) QueryResources(f func(*sparql.Results) error) error {
	arg := m.Called(f)
	return f(arg.Get(0).(*sparql.Results))
}

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) Publish(subjects []string) error {
	m.Called(subjects)
	return nil
}

func TestRebuildRepository(t *testing.T) {
	fakeWriter := new(MockedWriter)
	fakeRepo := new(MockRepo)

	organizationJSON := strings.NewReader(`{
    "head" : {
  "vars" : [ "s" ]
},
"results" : {
  "bindings" : [ {
    "s" : {
      "type" : "uri",
      "value" : "http://example.com/record1"
    }
  },{
    "s" : {
      "type" : "uri",
      "value" : "http://example.com/record2"
    }
  }
	 ]
}
    }`)
	results, _ := sparql.ParseJSON(organizationJSON)
	fakeRepo.On("QueryResources", mock.Anything).Return(results)

	fakeMessage := new(MockMessageService)
	fakeMessage.On("Publish", []string{"http://example.com/record1", "http://example.com/record2"})

	reg := &runtime.Registry{
		Derivative: fakeWriter,
		Canonical:  fakeRepo,
		Topic:      fakeMessage,
	}
	action := NewRebuildAction(reg)
	err := action.Run()

	assert.Nil(t, err)
}
