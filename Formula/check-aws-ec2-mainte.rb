class CheckAwsEc2Mainte < Formula
  desc 'Mackerel plugin to detect EC2 Retirement'
  homepage 'https://github.com/ntrv/check-aws-ec2-mainte'
  version '0.1.1'
  sha256 'ad6968b57e236907f8b2ebea0de97205a98f930187a104ef292632a83272a47b'
  url 'https://github.com/ntrv/check-aws-ec2-mainte/releases/download/v0.1.1/check-aws-ec2-mainte_darwin_amd64.zip'
  head 'https://github.com/ntrv/check-aws-ec2-mainte.git'

  def install
    bin.install 'check-aws-ec2-mainte'
  end

  test do
    system "#{bin}/check-aws-ec2-mainte", "--version"
  end
end
