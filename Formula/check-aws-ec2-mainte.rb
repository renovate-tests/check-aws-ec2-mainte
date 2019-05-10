class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.0.11'
  sha256 'f2b80b971c68ba87b8deaf03cfb743781e945aeee82016c3b3c8e2669207c9a4'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.1.1/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
