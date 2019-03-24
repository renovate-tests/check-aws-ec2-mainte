package checkawsec2mainte

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

// Set variable from -X option
var (
	version   = "indev"                     // git describe --tags
	buildDate = "1970-01-01 09:00:00+09:00" // date --rfc-3339=seconds
)

// Options ... Commandline Options
type Options struct {
	Region       string        `short:"r" long:"region" env:"AWS_REGION" description:"AWS Region"`
	CritDuration time.Duration `short:"c" long:"critical-duration" default:"72h" description:"Critical while duration"`
	InstanceIds  []string      `short:"i" long:"instance-id" description:"Filter as EC2 Instance Ids"`
	IsAll        bool          `short:"a" long:"all" description:"Fetch all instances events"`
	Version      func()        `short:"v" long:"version" description:"Print Build Information"`
}

type Checker struct {
	Opts Options
	Now  time.Time
}

func Do() {
	var ckr *checkers.Checker
	defer func() {
		if ckr != nil {
			ckr.Name = "EC2 Mainte"
			ckr.Exit()
		}
	}()

	c, err := NewChecker(os.Args)
	if err != nil {
		ckr = checkers.Unknown(err.Error())
		return
	}

	events, err := c.FetchEvents(context.Background())
	if err != nil {
		ckr = checkers.Unknown(err.Error())
		return
	}

	ckr = c.Run(events)
}

func NewChecker(args []string) (*Checker, error) {
	opts := Options{}

	opts.Version = func() {
		fmt.Fprintf(
			os.Stderr,
			"Version: %v\nGoVer: %v\nAwsSDKVer: %v\nBuildDate: %v\n",
			version,
			runtime.Version(),
			aws.SDKVersion,
			buildDate,
		)
		os.Exit(1)
	}

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return nil, err
	}

	return &Checker{
		Opts: opts,
		Now:  time.Now(),
	}, nil
}

func (c Checker) FetchEvents(ctx context.Context) (events EC2Events, err error) {
	// The default configuration sources are:
	// * Environment Variables
	// * Shared Configuration and Shared Credentials files.
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return
	}

	switch {
	case c.Opts.Region == "":
		events, err = c.FetchEC2MetaMainteEvents(ctx, cfg)
	case len(c.Opts.InstanceIds) == 0 && !c.Opts.IsAll:
		events, err = c.FetchEC2MetaMainteEvents(ctx, cfg)
	default: // len(c.Opts.InstanceIds) != 0 || c.Opts.IsAll
		cfg.Region = c.Opts.Region // Set Region from --region
		events, err = c.FetchEC2MainteEvents(ctx, cfg)
	}

	// Remove already completed events
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceStatusEvent.html
	events = events.Filter(StateCompleted, StateCanceled)
	return
}

func (c Checker) Run(events EC2Events) *checkers.Checker {
	if events.Len() != 0 {
		event := events.GetCloseEvent()

		if event.IsTimeOver(c.Now, c.Opts.CritDuration) {
			return checkers.Critical(event.CreateMessage())
		}
		return checkers.Warning(event.CreateMessage())
	}

	return checkers.Ok("Not coming EC2 instance events")
}

// Get EC2Events from Real EC2 API
func (c Checker) FetchEC2MainteEvents(ctx context.Context, cfg aws.Config) (EC2Events, error) {
	mt := EC2Mainte{
		Client:      ec2.New(cfg),
		InstanceIds: c.Opts.InstanceIds, // If fetch events for all instances, instanceId must empty
	}

	events, err := mt.GetEvents(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// Get EC2Events from EC2 Metadata
// If Region or Instance ID is empty or not --all specified
func (_ Checker) FetchEC2MetaMainteEvents(ctx context.Context, cfg aws.Config) (EC2Events, error) {
	mt := EC2MetaMainte{
		Client: ec2metadata.New(cfg),
	}

	events, err := mt.GetEvents(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}
