# check-aws-ec2-mainte

[![CI](https://travis-ci.org/ntrv/check-aws-ec2-mainte.svg?branch=master)](https://travis-ci.org/ntrv/check-aws-ec2-mainte)

[![Maintainability](https://api.codeclimate.com/v1/badges/84698652a28fdaf6bf5e/maintainability)](https://codeclimate.com/github/ntrv/check-aws-ec2-mainte/maintainability)

## Usage

```bash
$ TZ=Asia/Tokyo AWS_PROFILE=foo AWS_REGION=ap-northeast-1 ./check-aws-ec2-mainte --all
EC2 Mainte WARNING: Code: system-reboot, InstanceId: i-1234abcd, Date: 2019-03-29T21:00:00+09:00 - 2019-03-29T23:00:00+09:00, Description: scheduled reboot
```
