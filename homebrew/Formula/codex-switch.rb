class CodexSwitch < Formula
  desc "CLI for switching Codex auth profiles"
  homepage "https://github.com/quinn/multidex"
  url "https://github.com/quinn/multidex/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "efbd18aab342644d8a2e79b5593ce39809119cb666290b8fda885d0c9a467721"
  license "MIT"

  head "https://github.com/quinn/multidex.git", branch: "main"

  depends_on "bash"
  depends_on "fzf" => :optional

  def install
    bin.install "cli/codex-switch"
  end

  test do
    system "#{bin}/codex-switch", "help"
    system "#{bin}/codex-switch", "version"
  end
end
