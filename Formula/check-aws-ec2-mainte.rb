class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.2.0'
  sha256 'aedfa5785548be2259b94e23548d661a29607aa69a9a235735aa1d5d202798d0'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.2.0/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
