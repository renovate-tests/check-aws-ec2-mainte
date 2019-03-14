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
		OverrideDefaultFromEnvar("AWS_REGION").String()
	warnDuration = app.Flag("warning-duration", "Warning while duration").Short('w').
			Default("240h").Duration()
	critDuration = app.Flag("critical-duration", "Critical while duration").Short('c').
			Default("120h").Duration()
	instanceIds = app.Flag("instance-ids", "Available to specify multiple time").Short('i').
			Default(getInstanceIdFromMetadata()).Strings()
)

func Do() {
	ckr := run(os.Args[1:])
	ckr.Name = "EC2 Mainte"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {

	_, err := app.Parse(args)
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return checkers.Unknown(err.Error())
	}

	if *region == "" {
		cfg.Region = *region
	}

	mt, err := GetMainteInfo(ec2.New(cfg), *instanceIds...)
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
