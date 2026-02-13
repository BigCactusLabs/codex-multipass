class CodexMp < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v0.1.3.tar.gz"
  sha256 "a3cbb2028de4dd3b02e48f837b6ea08a6221938d07bfd8f4c449c83ac980cf3a"
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
