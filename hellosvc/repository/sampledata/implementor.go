package sampledata

import "context"

type SampleDataClient struct {
	// TODO: uncomment line below with imported client implementor, 
	// cli <contract.ImportedClient>
}

func (c *SampleDataClient) HealthCheck(ctx context.Context) error {
	return nil
}
