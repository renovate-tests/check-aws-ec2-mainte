# check-aws-ec2-mainte

Mackerel plugin to detect EC2 Retirement

[![GoDoc](https://godoc.org/github.com/ntrv/check-aws-ec2-mainte?status.svg)](https://godoc.org/github.com/ntrv/check-aws-ec2-mainte/lib)
[![CI](https://travis-ci.org/ntrv/check-aws-ec2-mainte.svg?branch=master)](https://travis-ci.org/ntrv/check-aws-ec2-mainte)
[![GitHub release](https://img.shields.io/github/release/ntrv/check-aws-ec2-mainte.svg)](https://github.com/ntrv/check-aws-ec2-mainte/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/84698652a28fdaf6bf5e/maintainability)](https://codeclimate.com/github/ntrv/check-aws-ec2-mainte/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/84698652a28fdaf6bf5e/test_coverage)](https://codeclimate.com/github/ntrv/check-aws-ec2-mainte/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/ntrv/check-aws-ec2-mainte)](https://goreportcard.com/report/github.com/ntrv/check-aws-ec2-mainte)
![Renovate](https://img.shields.io/badge/Renovate-Bot-yellow?style=flat&logo=dependabot)
![Go1.13](https://img.shields.io/badge/Go-1.13-blue.svg)
![GitHub](https://img.shields.io/github/license/ntrv/check-aws-ec2-mainte.svg)


## Install

### mkr

```bash
mkr plugin install ntrv/check-aws-ec2-mainte
```

### HomeBrew

```bash
brew install https://raw.githubusercontent.com/ntrv/check-aws-ec2-mainte/master/Formula/check-aws-ec2-mainte.rb
```

## Usage

```bash
$ TZ=Asia/Tokyo AWS_PROFILE=foo AWS_REGION=ap-northeast-1 ./check-aws-ec2-mainte --all
EC2 Mainte WARNING: Code: system-reboot, InstanceId: i-1234abcd, Date: 2019-03-29T21:00:00+09:00 - 2019-03-29T23:00:00+09:00, Description: scheduled reboot
```

## Example of mackerel-agent.conf

```toml
[plugin.check.ec2mainte]
command = ["/path/to/check-aws-ec2-mainte", "-c", "48h"]
env = { AWS_REGION = "ap-northeast-1", TZ = "US/Pacific" }
max_check_attempts = 3
timeout_seconds = 15
```

## ToDo

* Detect Spot Instance Interruptions
