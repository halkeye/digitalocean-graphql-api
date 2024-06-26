package loaders

// import vikstrous/dataloadgen with your other imports
import (
	"context"
	"time"

	"github.com/vikstrous/dataloadgen"

	"github.com/halkeye/digitalocean-graphql-api/graph/model"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	DropletLoader      *dataloadgen.Loader[string, *model.Droplet]
	AppLoader          *dataloadgen.Loader[string, *model.App]
	DomainLoader       *dataloadgen.Loader[string, *model.Domain]
	ProjectLoader      *dataloadgen.Loader[string, *model.Project]
	VolumeLoader       *dataloadgen.Loader[string, *model.Volume]
	KubernetesLoader   *dataloadgen.Loader[string, *model.KubernetesCluster]
	DbaasLoader        *dataloadgen.Loader[string, *model.Dbaas]
	LoadBalancerLoader *dataloadgen.Loader[string, *model.LoadBalancer]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders() *Loaders {
	return &Loaders{
		ProjectLoader:      dataloadgen.NewLoader((&projectReader{}).getProjects, dataloadgen.WithWait(time.Millisecond)),
		DropletLoader:      dataloadgen.NewLoader((&dropletReader{}).getDroplets, dataloadgen.WithWait(time.Millisecond)),
		LoadBalancerLoader: dataloadgen.NewLoader((&loadbalancerReader{}).getLoadBalancers, dataloadgen.WithWait(time.Millisecond)),
		AppLoader:          dataloadgen.NewLoader((&appReader{}).getApps, dataloadgen.WithWait(time.Millisecond)),
		DomainLoader:       dataloadgen.NewLoader((&domainReader{}).getDomains, dataloadgen.WithWait(time.Millisecond)),
		VolumeLoader:       dataloadgen.NewLoader((&volumeReader{}).getVolumes, dataloadgen.WithWait(time.Millisecond)),
		KubernetesLoader:   dataloadgen.NewLoader((&kubernetesReader{}).getKubernetesClusters, dataloadgen.WithWait(time.Millisecond)),
		DbaasLoader:        dataloadgen.NewLoader((&dbaasReader{}).getDbaass, dataloadgen.WithWait(time.Millisecond)),
	}
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func WithContext(ctx context.Context) context.Context {
	loader := NewLoaders()
	return context.WithValue(ctx, loadersKey, loader)
}
