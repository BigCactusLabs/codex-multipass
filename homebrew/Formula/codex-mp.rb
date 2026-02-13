class CodexMp < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "153716dbcc2a1810096b159ebf5593d02d64b54b39c0ea15175a32ea4e990368"
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
