# check-aws-ec2-mainte

[![GoDoc](https://godoc.org/github.com/ntrv/check-aws-ec2-mainte?status.svg)](https://godoc.org/github.com/ntrv/check-aws-ec2-mainte/lib)
[![CI](https://travis-ci.org/ntrv/check-aws-ec2-mainte.svg?branch=master)](https://travis-ci.org/ntrv/check-aws-ec2-mainte)
[![GitHub release](https://img.shields.io/github/release/ntrv/check-aws-ec2-mainte.svg)](https://github.com/ntrv/check-aws-ec2-mainte/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/84698652a28fdaf6bf5e/maintainability)](https://codeclimate.com/github/ntrv/check-aws-ec2-mainte/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/84698652a28fdaf6bf5e/test_coverage)](https://codeclimate.com/github/ntrv/check-aws-ec2-mainte/test_coverage)
![Go1.12](https://img.shields.io/badge/Go-1.12-blue.svg)
![GitHub](https://img.shields.io/github/license/ntrv/check-aws-ec2-mainte.svg)


## Usage

```bash
$ TZ=Asia/Tokyo AWS_PROFILE=foo AWS_REGION=ap-northeast-1 ./check-aws-ec2-mainte --all
EC2 Mainte WARNING: Code: system-reboot, InstanceId: i-1234abcd, Date: 2019-03-29T21:00:00+09:00 - 2019-03-29T23:00:00+09:00, Description: scheduled reboot
```
