class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.0.9'
  sha256 '3d86533af4438db9229d0f92fb621ec0e7edb9082a63085c0b9281ac56637f98'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.0.9/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
