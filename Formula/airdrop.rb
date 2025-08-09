class Airdrop < Formula
  desc "Simple description of airdrop"
  homepage "https://github.com/xvoland/airdrop"
  url "https://github.com/xvoland/airdrop/releases/download/v0.2.5/airdrop_darwin_arm64.tar.gz"
  sha256 "e7a8990276b73f07d39dafa827fc95adbbce265ea253ad4c085e7ac643d1fab0"

  def install
    bin.install "airdrop"
  end

  test do
    system "#{bin}/airdrop", "--version"
  end
end
