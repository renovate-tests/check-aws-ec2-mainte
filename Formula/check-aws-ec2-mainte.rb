class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.1.1'
  sha256 '515e7fd5a53f5db4f6b2f789965317837beaec19c1ba122cbc85c8c5f5f71c00'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.1.1/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
