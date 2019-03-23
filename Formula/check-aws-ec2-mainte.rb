class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.0.6'
  sha256 'ab5824fe7bf032005d03285e5b2a0b6ebb3e436203cca8296035abf95a12ec96'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.0.6/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
