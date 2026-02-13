class CodexMp < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v0.1.1.tar.gz"
  sha256 "7502f64a5f3a75efe51a9a89a271e01fd67bd5519d8d355758e8a0a529cca698"
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
