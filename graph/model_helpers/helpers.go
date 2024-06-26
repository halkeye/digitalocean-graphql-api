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

func AccountFromGodo(account *godo.Account) (*model.Account, error) {
	id := fmt.Sprintf("do:user:%s", account.UUID)
	uuidObj, err := uuid.Parse(account.Team.UUID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse updated at: %w", err)
	}

	return &model.Account{
		ID:            id,
		Email:         account.Email,
		EmailVerified: account.EmailVerified,
		Status:        account.Status,
		UUID:          account.UUID,
		Team: &model.Team{
			UUID: uuidObj,
			Name: account.Team.Name,
		},
	}, nil
}

func VolumeFromGodo(volume *godo.Volume) *model.Volume {
	id := fmt.Sprintf("do:volume:%s", volume.ID)

	return &model.Volume{
		ID:          id,
		Name:        volume.Name,
		Description: volume.Description,
	}
}

func DropletFromGodo(droplet *godo.Droplet) *model.Droplet {
	return &model.Droplet{
		ID:     fmt.Sprintf("do:droplet:%d", droplet.ID),
		Name:   droplet.Name,
		Memory: &droplet.Memory,
		Vcpus:  &droplet.Vcpus,
		Disk:   &droplet.Disk,
		Region: &model.Region{
			ID:   droplet.Region.Slug,
			Name: droplet.Region.Name,
		},
		SizeSlug:  &droplet.SizeSlug,
		BackupIDs: droplet.BackupIDs,
	}
}

func DbaasFromGodo(dbaas *godo.Database) *model.Dbaas {
	return &model.Dbaas{
		ID:   fmt.Sprintf("do:dbaas:%s", dbaas.ID),
		Name: dbaas.Name,
	}
}

func KubernetesClusterFromGodo(k8s *godo.KubernetesCluster) *model.KubernetesCluster {
	return &model.KubernetesCluster{
		ID:   fmt.Sprintf("do:kubernetes:%s", k8s.ID),
		Name: k8s.Name,
	}
}

func LoadBalancerFromGodo(lb *godo.LoadBalancer) *model.LoadBalancer {
	return &model.LoadBalancer{
		ID:   fmt.Sprintf("do:loadbalancer:%s", lb.ID),
		Name: lb.Name,
	}
}
