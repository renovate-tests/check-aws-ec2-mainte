package checkawsec2mainte

import (
	"fmt"
	"os"

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
	ckr := run(os.Args[1:])
	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func prepare(args []string) (EC2Maintes, error) {

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

	mt, err := GetMainteInfo(ec2.New(cfg), *instanceIds...)
	if err != nil {
		return nil, err
	}
	return mt, nil
}

func run(args []string) *checkers.Checker {

	mt, err := prepare(args)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	if mt.Len() != 0 {
		event := mt.GetCloseEvent()

		if event.IsTimeOver(*critDuration) {
			return checkers.Critical(fmt.Sprintf("%+v", event))
		}
		return checkers.Warning(fmt.Sprintf("%+v", event))
	}

	return checkers.Ok("Not coming EC2 instance events")
}
