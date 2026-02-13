class CodexSwitch < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/BigCactusLabs/codex-multipass"
  url "https://github.com/BigCactusLabs/codex-multipass/archive/refs/tags/v0.1.0.tar.gz" # Placeholder
  sha256 "REPLACE_WITH_ACTUAL_SHA256" # Placeholder
  license "MIT"

  depends_on "bash"
  depends_on "fzf" => :optional

  def install
    bin.install "cli/codex-switch"
  end

  test do
    system "#{bin}/codex-switch", "help"
  end
end
