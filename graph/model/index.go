package model

import (
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/google/uuid"
)

func AppFromGodo(app *godo.App) *App {
	return &App{
		ID:                     app.URN(),
		Owner:                  &Team{UUID: uuid.MustParse(app.OwnerUUID)},
		LastDeploymentActiveAt: &app.LastDeploymentActiveAt,
		DefaultIngress:         &app.DefaultIngress,
		CreatedAt:              &app.CreatedAt,
		UpdatedAt:              &app.UpdatedAt,
	}
}

func ProjectFromGodo(project *godo.Project) *Project {
	parsedUUID, err := uuid.Parse(project.OwnerUUID)
	if err != nil {
		panic(fmt.Errorf("unable to parse uuid: %w", err))
	}

	createdAt, err := time.Parse(time.RFC3339, project.CreatedAt)
	if err != nil {
		panic(fmt.Errorf("unable to parse created at: %w", err))
	}

	updatedAt, err := time.Parse(time.RFC3339, project.UpdatedAt)
	if err != nil {
		panic(fmt.Errorf("unable to parse updated at: %w", err))
	}

	id := fmt.Sprintf("do:project:%s", project.ID)

	return &Project{
		ID:          id,
		Owner:       &Team{ID: fmt.Sprintf("do:team:%s", project.OwnerUUID), UUID: parsedUUID},
		Name:        project.Name,
		Description: &project.Description,
		Purpose:     project.Purpose,
		Environment: project.Environment,
		IsDefault:   project.IsDefault,
		CreatedAt:   &createdAt,
		UpdatedAt:   &updatedAt,
	}
}
