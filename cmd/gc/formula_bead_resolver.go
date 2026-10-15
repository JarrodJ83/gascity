package main

import (
	"github.com/gastownhall/gascity/internal/beads"
	"github.com/gastownhall/gascity/internal/config"
	"github.com/gastownhall/gascity/internal/storeref"
)

// newLazyBeadResolver returns a by-id bead lookup over the topology build
// produces. build runs on the first lookup only, so a caller that never looks
// anything up never builds the topology.
func newLazyBeadResolver(_ func() storeref.Topology) func(id string) (beads.Bead, error) {
	return func(string) (beads.Bead, error) { return beads.Bead{}, beads.ErrNotFound }
}

// formulaBeadResolver is the worker's formula-attribution lookup for this
// runtime: the acting bead and its run root resolved over the serving stores.
// topology is nil in production and builds the serving frame; tests inject a
// counting builder.
func (cr *CityRuntime) formulaBeadResolver(cfg *config.City, topology func() storeref.Topology) func(id string) (beads.Bead, error) {
	if topology == nil {
		topology = func() storeref.Topology { return cr.residencyTopology(cr.routeRecoveryRigStores(cfg)) }
	}
	return newLazyBeadResolver(topology)
}
