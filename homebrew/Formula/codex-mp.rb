class CodexMp < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v0.1.6.tar.gz"
  sha256 "58c4c89df7881512d668e9077712c9af7340396f75ea96b27388127f3ae8e592"
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
