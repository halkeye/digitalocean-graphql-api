// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/google/uuid"
)

// An object with an ID
type Node interface {
	IsNode()
	// The id of the object.
	GetID() string
}

type Resource interface {
	IsResource()
	GetID() string
	GetName() string
}

// Account Information
type Account struct {
	// Email address
	Email string `json:"email"`
	// Has email been verified
	EmailVerified bool `json:"emailVerified"`
	// The id of the account
	ID string `json:"id"`
	// Account Status
	Status string `json:"status"`
	// Account UUID
	UUID string `json:"uuid"`
	// Team
	Team *Team `json:"team"`
}

func (Account) IsNode() {}

// The id of the object.
func (this Account) GetID() string { return this.ID }

// Account Limits
type AccountLimits struct {
	// How many droplets can you have at once
	DropletLimit int `json:"dropletLimit"`
	// How many volumes can you have at once
	VolumeLimit int `json:"volumeLimit"`
}

type App struct {
	ID                     string     `json:"id"`
	Name                   string     `json:"name"`
	Owner                  *Team      `json:"owner"`
	LastDeploymentActiveAt *time.Time `json:"lastDeploymentActiveAt,omitempty"`
	DefaultIngress         *string    `json:"defaultIngress,omitempty"`
	CreatedAt              *time.Time `json:"createdAt,omitempty"`
	UpdatedAt              *time.Time `json:"updatedAt,omitempty"`
}

func (App) IsNode() {}

// The id of the object.
func (this App) GetID() string { return this.ID }

func (App) IsResource() {}

func (this App) GetName() string { return this.Name }

type Domain struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	TTL      int     `json:"ttl"`
	ZoneFile *string `json:"zoneFile,omitempty"`
}

func (Domain) IsNode() {}

// The id of the object.
func (this Domain) GetID() string { return this.ID }

func (Domain) IsResource() {}

func (this Domain) GetName() string { return this.Name }

type Droplet struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Memory    *int    `json:"memory,omitempty"`
	Vcpus     *int    `json:"vcpus,omitempty"`
	Disk      *int    `json:"disk,omitempty"`
	Region    *Region `json:"region,omitempty"`
	SizeSlug  *string `json:"sizeSlug,omitempty"`
	BackupIDs []int   `json:"backupIDs"`
}

func (Droplet) IsNode() {}

// The id of the object.
func (this Droplet) GetID() string { return this.ID }

func (Droplet) IsResource() {}

func (this Droplet) GetName() string { return this.Name }

// Information about pagination in a connection.
type PageInfo struct {
	// When paginating forwards, the cursor to continue.
	EndCursor *string `json:"endCursor,omitempty"`
	// When paginating forwards, are there more items?
	HasNextPage bool `json:"hasNextPage"`
	// When paginating backwards, are there more items?
	HasPreviousPage bool `json:"hasPreviousPage"`
	// When paginating backwards, the cursor to continue.
	StartCursor *string `json:"startCursor,omitempty"`
}

// Projects allow you to organize your resources into groups that fit the way you work. You can group resources (like Droplets, Spaces, load balancers, domains, and floating IPs) in ways that align with the applications you host on DigitalOcean.
type Project struct {
	// The id of the account
	ID          string     `json:"id"`
	Owner       *Team      `json:"owner"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Purpose     string     `json:"purpose"`
	Environment string     `json:"environment"`
	IsDefault   bool       `json:"isDefault"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	// Project Resources
	Resources *ProjectResourcesConnection `json:"resources,omitempty"`
}

func (Project) IsNode() {}

// The id of the object.
func (this Project) GetID() string { return this.ID }

type ProjectResource struct {
	ID         string    `json:"id"`
	AssignedAt time.Time `json:"assignedAt"`
	Resource   Resource  `json:"resource,omitempty"`
	Status     string    `json:"status"`
}

// ProjectResources Connection
type ProjectResourcesConnection struct {
	// Edges
	Edges []*ProjectResourcesEdge `json:"edges"`
	// Pagination info
	PageInfo *PageInfo `json:"pageInfo"`
}

// ProjectResources Edge
type ProjectResourcesEdge struct {
	// Cursor
	Cursor string `json:"cursor"`
	// Project Node
	Node *ProjectResource `json:"node,omitempty"`
}

// Projects Connection
type ProjectsConnection struct {
	// Edges
	Edges []*ProjectsEdge `json:"edges"`
	// Pagination info
	PageInfo *PageInfo `json:"pageInfo"`
}

// Project Edge
type ProjectsEdge struct {
	// Cursor
	Cursor string `json:"cursor"`
	// Project Resource Node
	Node *Project `json:"node,omitempty"`
}

// All the queries
type Query struct {
}

type Region struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Sizes     []string `json:"sizes"`
	Available *bool    `json:"available,omitempty"`
	Features  []string `json:"features"`
}

func (Region) IsNode() {}

// The id of the object.
func (this Region) GetID() string { return this.ID }

// Team information
type Team struct {
	// The id of the team
	ID string `json:"id"`
	// What is the teams limits
	Limits *AccountLimits `json:"limits,omitempty"`
	// Team UUID
	UUID uuid.UUID `json:"uuid"`
}

func (Team) IsNode() {}

// The id of the object.
func (this Team) GetID() string { return this.ID }

type Volume struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (Volume) IsNode() {}

// The id of the object.
func (this Volume) GetID() string { return this.ID }

func (Volume) IsResource() {}

func (this Volume) GetName() string { return this.Name }
