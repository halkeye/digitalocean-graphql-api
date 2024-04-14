package model_helpers

import (
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/google/uuid"

	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

func AppFromGodo(app *godo.App) *model.App {
	return &model.App{
		ID:                     app.URN(),
		Name:                   app.Spec.Name,
		Owner:                  &model.Team{UUID: uuid.MustParse(app.OwnerUUID)},
		LastDeploymentActiveAt: &app.LastDeploymentActiveAt,
		DefaultIngress:         &app.DefaultIngress,
		CreatedAt:              app.CreatedAt,
		UpdatedAt:              &app.UpdatedAt,
	}
}

func ProjectFromGodo(project *godo.Project) *model.Project {
	createdAt, err := time.Parse(time.RFC3339, project.CreatedAt)
	if err != nil {
		panic(fmt.Errorf("unable to parse created at: %w", err))
	}

	updatedAt, err := time.Parse(time.RFC3339, project.UpdatedAt)
	if err != nil {
		panic(fmt.Errorf("unable to parse updated at: %w", err))
	}

	id := fmt.Sprintf("do:project:%s", project.ID)

	return &model.Project{
		ID:          id,
		Owner:       &model.Team{ID: fmt.Sprintf("do:team:%s", project.OwnerUUID), UUID: uuid.MustParse(project.OwnerUUID)},
		Name:        project.Name,
		Description: &project.Description,
		Purpose:     project.Purpose,
		Environment: project.Environment,
		IsDefault:   project.IsDefault,
		CreatedAt:   createdAt,
		UpdatedAt:   &updatedAt,
	}
}

func AccountFromGodo(account *godo.Account) *model.Account {
	id := fmt.Sprintf("do:user:%s", account.UUID)

	return &model.Account{
		ID: id,

		Email:         account.Email,
		EmailVerified: account.EmailVerified,
		Status:        account.Status,
		UUID:          account.UUID,
	}
}
