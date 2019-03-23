package checkawsec2mainte

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/mackerelio/checkers"

	"github.com/jessevdk/go-flags"
)

// Set variable from -X option
var (
	version   = "indev"                     // git describe --tags
	buildDate = "1970-01-01 09:00:00+09:00" // date --rfc-3339=seconds
)

// Options ... Commandline Options
type Options struct {
	Region       string        `short:"r" long:"region" required:"true" env:"AWS_REGION" description:"AWS Region"`
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
	ctx := context.Background()

	c, err := NewChecker(os.Args)
	if err != nil {
		os.Exit(1)
	}

	events, err := c.FetchEvents(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ckr := c.Run(events)

	ckr.Name = "EC2 Mainte"
	ckr.Exit()
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

func (c Checker) FetchEvents(ctx context.Context) (EC2Events, error) {
	// The default configuration sources are:
	// * Environment Variables
	// * Shared Configuration and Shared Credentials files.
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	// Set Region from --region
	if c.Opts.Region != "" {
		cfg.Region = c.Opts.Region
	}

	// Default instanceId is from EC2 metadata
	// If fetch events for all instances, instanceId must empty
	instanceIds := c.Opts.InstanceIds

	// Get Instance ID from EC2 Metadata
	if len(c.Opts.InstanceIds) == 0 && !c.Opts.IsAll {
		instanceId, err := GetInstanceIdFromMetadata(cfg)
		if err != nil {
			return nil, err
		}
		instanceIds = []string{instanceId}
	}

	mt := EC2Mainte{
		Client:      ec2.New(cfg),
		InstanceIds: instanceIds,
	}

	events, err := mt.GetMainteInfo(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
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
