package projectapp

import (
	"TODO-list/business/domain/projectbus"
	"encoding/json"
	"time"
)

// NewProject represents the input data required to create a new project.
type NewProject struct {
	Name      string `json:"name"`
	CreatedBy int    `json:"created_by"`
}

// Decode decodes a JSON byte slice into a NewProject struct.
func (np *NewProject) Decode(data []byte) error {
	return json.Unmarshal(data, &np)
}

// toBusNewProject converts a NewProject struct from the application layer to the business layer representation.
func toBusNewProject(np NewProject) projectbus.NewProject {
	return projectbus.NewProject{
		Name:      np.Name,
		CreatedBy: np.CreatedBy,
	}
}

// Project represents a project entity in the application layer with all relevant fields.
type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
}

// Encode encodes the Project struct into a JSON byte slice.
func (p Project) Encode() ([]byte, string, error) {
	data, err := json.Marshal(p)
	return data, "application/json", err
}

// toAppProject converts a Project struct from the business layer to the application layer representation.
func toAppProject(projectbus projectbus.Project) Project {
	return Project{
		ID:        projectbus.ID,
		Name:      projectbus.Name,
		Active:    projectbus.Active,
		CreatedAt: projectbus.CreatedAt,
		CreatedBy: projectbus.CreatedBy,
	}
}

// Projects represents a collection of Project entities.
type Projects []Project

// Encode encodes the Projects slice into a JSON byte slice.
func (ps Projects) Encode() ([]byte, string, error) {
	data, err := json.Marshal(ps)
	return data, "application/json", err
}

// toAppProjects converts a slice of Project structs from the business layer to the application layer representation.
func toAppProjects(projectsBus []projectbus.Project) Projects {
	projectsApp := make(Projects, len(projectsBus))
	for i, projectBus := range projectsBus {
		projectsApp[i] = toAppProject(projectBus)
	}
	return projectsApp
}

// UpdateProject represents the input data required to update an existing project.
type UpdateProject struct {
	Name string `json:"name"`
}

// Decode decodes a JSON byte slice into an UpdateProject struct.
func (up *UpdateProject) Decode(data []byte) error {
	return json.Unmarshal(data, &up)
}

// toBusUpdateProject converts an UpdateProject struct from the application layer to the business layer representation.
func toBusUpdateProject(up UpdateProject) projectbus.UpdateProject {
	return projectbus.UpdateProject{
		Name: up.Name,
	}
}
