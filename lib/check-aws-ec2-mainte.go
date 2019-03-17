package checkawsec2mainte

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/mackerelio/checkers"

	"github.com/jessevdk/go-flags"
)

var (
	version    = "indev"
	commitHash = ""
	buildDate  = ""
)

type options struct {
	Region       string        `short:"r" long:"region" required:"true" env:"AWS_REGION" description:"AWS Region"`
	CritDuration time.Duration `short:"c" long:"critical-duration" default:"72h" description:"Critical while duration"`
	InstanceIds  []string      `short:"i" long:"instance-id" description:"Filter as EC2 Instance Ids"`
	IsAll        bool          `short:"a" long:"all" description:"Fetch all instances events"`
	TimeZone     string        `long:"tz" env:"TZ" default:"Asia/Tokyo" description:"TimeZone to create message"`
	Version      func()        `short:"v" long:"version" description:"Print Build Information"`
}

type Checker struct {
	Opts options
	Now  time.Time
}

func Do() {

	c, err := NewChecker(os.Args)
	if err != nil {
		os.Exit(1)
	}

	events, err := c.fetchEvents()
	if err != nil {
		os.Exit(1)
	}

	ckr := c.run(events)

	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func NewChecker(args []string) (*Checker, error) {

	opts := options{}

	opts.Version = func() {
		fmt.Fprintf(
			os.Stderr,
			"Version: %v\nGoVer: %v\nCommitHash: %v\nBuildDate: %v\n",
			version,
			runtime.Version(),
			commitHash,
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

func (c Checker) fetchEvents() (EC2Events, error) {
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

	if len(c.Opts.InstanceIds) == 0 && !c.Opts.IsAll {
		instanceId, err := getInstanceIdFromMetadata(cfg)
		if err != nil {
			return nil, err
		}
		instanceIds = []string{instanceId}
	}

	mt := EC2Mainte{
		Client:      ec2.New(cfg),
		InstanceIds: instanceIds,
	}

	events, err := mt.GetMainteInfo()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c Checker) run(events EC2Events) *checkers.Checker {
	if events.Len() != 0 {
		event := events.GetCloseEvent()

		if event.IsTimeOver(c.Now, c.Opts.CritDuration) {
			return checkers.Critical(event.CreateMessage(c.Opts.TimeZone))
		}
		return checkers.Warning(event.CreateMessage(c.Opts.TimeZone))
	}

	return checkers.Ok("Not coming EC2 instance events")
}
