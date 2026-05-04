class CodexMp < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/releases/download/v0.2.0/codex-multipass-v0.2.0.tar.gz"
  sha256 "27b25ace8c0f0b995f0ba804ded2780b95027b97085db6e544deab912e9552ed"
  license "MIT"

  head "https://github.com/BigCactusLabs/codex-multipass.git", branch: "main"

  depends_on "go" => :build
  depends_on "bash"
  depends_on "fzf" => :optional

  def install
    system "make", "build"
    bin.install "codex-mp"
  end

  test do
    system "#{bin}/codex-mp", "help"
    system "#{bin}/codex-mp", "version"
  end
end
