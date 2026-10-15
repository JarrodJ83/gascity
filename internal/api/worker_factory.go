package api

import (
	"github.com/gastownhall/gascity/internal/beads"
	"github.com/gastownhall/gascity/internal/session"
	"github.com/gastownhall/gascity/internal/worker"
)

func (s *Server) workerFactory(store beads.Store) (*worker.Factory, error) {
	cfg := s.state.Config()
	var resolveTransport func(template, provider string) string
	if cfg != nil {
		resolveTransport = func(template, provider string) string {
			return configuredSessionTransport(cfg, template, provider)
		}
	}
	return worker.NewFactory(worker.FactoryConfig{
		Store:                 store,
		Provider:              s.state.SessionProvider(),
		CityPath:              s.state.CityPath(),
		SearchPaths:           s.sessionLogPaths(),
		Recorder:              s.state.EventProvider(),
		UsageSink:             s.state.UsageSink(),
		ResolveTransport:      resolveTransport,
		ResolveSessionRuntime: s.resolveWorkerSessionRuntimeWithMetadata,
		Pricing:               cfg.PricingRegistry(),
		ActivityMemo:          s.activityMemo,
		RuntimeLeaseTTL:       session.RuntimeLeaseTTLFor(cfg),
	})
}

func (s *Server) workerSessionCatalog(store beads.Store) (*worker.SessionCatalog, error) {
	factory, err := s.workerFactory(store)
	if err != nil {
		return nil, err
	}
	return factory.Catalog()
}

// formulaBeadResolver is the worker's formula-attribution lookup: the acting
// bead and its run root resolved over this server's residency topology.
func (s *Server) formulaBeadResolver() func(id string) (beads.Bead, error) {
	return func(id string) (beads.Bead, error) {
		_, b, err := s.resolveBeadOwner(id)
		return b, err
	}
}
