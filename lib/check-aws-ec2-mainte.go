package checkawsec2mainte

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/mackerelio/checkers"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version    = "indev"
	goversion  = ""
	commitHash = ""
	buildDate  = ""
	revision   = fmt.Sprintf(
		"GoVer: %v\tCommitHash: %v\tBuildDate: %v",
		goversion,
		commitHash,
		buildDate,
	)
)

var (
	app = kingpin.New("check-aws-ec2-mainte", revision).Version(version).
		Author("ntrv")
	region = app.Flag("region", "AWS Region").Short('r').
		OverrideDefaultFromEnvar("AWS_REGION").Required().String()
	warnDuration = app.Flag("warning-duration", "Warning while duration").Short('w').
			Default("240h").Duration()
	critDuration = app.Flag("critical-duration", "Critical while duration").Short('c').
			Default("120h").Duration()
	instanceIds = app.Flag("instance-ids", "Available to specify multiple time").Short('i').
			Strings()
)

func Do() {
	var ckr *checkers.Checker

	events, err := fetchEvents(os.Args[1:])
	if err != nil {
		ckr = checkers.Unknown(err.Error())
	} else {
		ckr = run(events, time.Now())
	}

	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func fetchEvents(args []string) (EC2Events, error) {

	_, err := app.Parse(args)
	if err != nil {
		return nil, err
	}

	// The default configuration sources are:
	// * Environment Variables
	// * Shared Configuration and Shared Credentials files.
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	if *region != "" {
		cfg.Region = *region
	}

	// Default instanceId is from EC2 metadata
	if len(*instanceIds) == 0 {
		instanceId, err := getInstanceIdFromMetadata(cfg)
		if err != nil {
			return nil, err
		}
		*instanceIds = append(*instanceIds, instanceId)
	}

	mt := EC2Mainte{
		Client:      ec2.New(cfg),
		InstanceIds: *instanceIds,
	}

	events, err := mt.GetMainteInfo()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func run(events EC2Events, now time.Time) *checkers.Checker {
	if events.Len() != 0 {
		event := events.GetCloseEvent()

		if event.IsTimeOver(now, *critDuration) {
			return checkers.Critical(fmt.Sprintf("%+v", event))
		}
		return checkers.Warning(fmt.Sprintf("%+v", event))
	}

	return checkers.Ok("Not coming EC2 instance events")
}
