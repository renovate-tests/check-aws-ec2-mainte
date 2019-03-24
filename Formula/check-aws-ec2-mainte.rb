class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.0.7'
  sha256 '7d7c278405f5f76a71ecfaeba3f4309924a27a11ccb0646d262c96f9d3830f33'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.0.7/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
