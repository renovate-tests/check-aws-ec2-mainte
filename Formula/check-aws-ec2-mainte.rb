class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.1.3'
  sha256 'e67b8ded88d90f12075d59a1ed7f9cb5f5f67a7a612b69384f07399a6dc8b849'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.1.3/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
