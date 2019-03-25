package checkawsec2mainte

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
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

type Checker struct {
	Opts Options
	Now  time.Time
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

func (c Checker) Run(events Events) *checkers.Checker {
	if events.Len() != 0 {
		event := events.GetCloseEvent()

		if event.IsTimeOver(c.Now, c.Opts.CritDuration) {
			return checkers.Critical(event.CreateMessage())
		}
		return checkers.Warning(event.CreateMessage())
	}

	return checkers.Ok("Not coming EC2 instance events")
}
