package checkawsec2mainte

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ntrv/check-aws-ec2-mainte/lib/ec2api"
	"github.com/ntrv/check-aws-ec2-mainte/lib/events"
	"github.com/ntrv/check-aws-ec2-mainte/lib/metadata"
)

// Fetch ...
func (c Cli) Fetch(ctx context.Context) (evs events.Events, err error) {
	// The default configuration sources are:
	// * Environment Variables
	// * Shared Configuration and Shared Credentials files.
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return
	}

	switch {
	case c.Args.Region == "":
		evs, err = c.fetchEC2MetaEvents(ctx, cfg)
	case len(c.Args.InstanceIds) == 0 && !c.Args.IsAll:
		evs, err = c.fetchEC2MetaEvents(ctx, cfg)
	default: // len(c.Opts.InstanceIds) != 0 || c.Opts.IsAll
		cfg.Region = c.Args.Region // Set Region from --region
		evs, err = c.fetchEC2Events(ctx, cfg)
	}

	// Remove already completed events
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceStatusEvent.html
	evs = evs.Filter(events.StateCompleted, events.StateCanceled)
	return
}

// fetchEC2Events ... Get EC2Events from Real EC2 API
func (c Cli) fetchEC2Events(ctx context.Context, cfg aws.Config) (evs events.Events, err error) {
	mt := ec2api.Mainte{
		Client:      ec2.New(cfg),
		InstanceIds: c.Args.InstanceIds, // If fetch events for all instances, instanceId must empty
	}
	return mt.Fetch(ctx)
}

// fetchEC2MetaEvents ... Get EC2Events from EC2 Metadata
// If Region or Instance ID is empty or not --all specified
func (c Cli) fetchEC2MetaEvents(ctx context.Context, cfg aws.Config) (evs events.Events, err error) {
	mt := metadata.Mainte{
		Client: ec2metadata.New(cfg),
	}
	return mt.Fetch(ctx)
}
