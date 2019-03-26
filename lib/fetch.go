package checkawsec2mainte

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/ntrv/check-aws-ec2-mainte/lib/ec2api"
	"github.com/ntrv/check-aws-ec2-mainte/lib/metadata"
)

func (c Checker) FetchEvents(ctx context.Context) (events Events, err error) {
	// The default configuration sources are:
	// * Environment Variables
	// * Shared Configuration and Shared Credentials files.
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return
	}

	switch {
	case c.Opts.Region == "":
		events, err = c.FetchEC2MetaEvents(ctx, cfg)
	case len(c.Opts.InstanceIds) == 0 && !c.Opts.IsAll:
		events, err = c.FetchEC2MetaEvents(ctx, cfg)
	default: // len(c.Opts.InstanceIds) != 0 || c.Opts.IsAll
		cfg.Region = c.Opts.Region // Set Region from --region
		events, err = c.FetchEC2Events(ctx, cfg)
	}

	// Remove already completed events
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceStatusEvent.html
	events = events.Filter(StateCompleted, StateCanceled)
	return
}

// Get EC2Events from Real EC2 API
func (c Checker) FetchEC2Events(ctx context.Context, cfg aws.Config) (events Events, err error) {
	mt := ec2api.Mainte{
		Client:      ec2.New(cfg),
		InstanceIds: c.Opts.InstanceIds, // If fetch events for all instances, instanceId must empty
	}

	evs, err := mt.GetEvents(ctx)
	if err != nil {
		return
	}
	events.SetEC2APIEvents(evs)
	return
}

// Get EC2Events from EC2 Metadata
// If Region or Instance ID is empty or not --all specified
func (_ Checker) FetchEC2MetaEvents(ctx context.Context, cfg aws.Config) (events Events, err error) {
	mt := metadata.Mainte{
		Client: ec2metadata.New(cfg),
	}

	evs, err := mt.GetEvents(ctx)
	if err != nil {
		return
	}
	events.SetMetadataEvents(evs)
	return
}
