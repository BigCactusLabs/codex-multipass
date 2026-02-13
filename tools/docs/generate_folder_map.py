#!/usr/bin/env python3
"""Generate and verify a deterministic repository folder map."""

from __future__ import annotations

import argparse
import difflib
import os
import re
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, Iterator, Sequence


DEFAULT_OUTPUT = "docs/repo_folder_map.md"
DEFAULT_SUMMARY = "Regenerated folder map snapshot"
VERSION_RE = re.compile(r"^- map_version:\s*([0-9]+\.[0-9]+\.[0-9]+)\s*$", re.MULTILINE)
BUMP_RE = re.compile(r"^- bump:\s*(none|patch|minor|major)\s*$", re.MULTILINE)
SUMMARY_RE = re.compile(r"^- summary:\s*(.*)\s*$", re.MULTILINE)


@dataclass(frozen=True)
class Node:
    name: str
    path: Path
    is_dir: bool


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo-root", default=".", help="Repository root path.")
    parser.add_argument(
        "--output",
        default=DEFAULT_OUTPUT,
        help=f"Generated map output path, relative to repo root (default: {DEFAULT_OUTPUT}).",
    )
    parser.add_argument(
        "--bump",
        choices=("none", "patch", "minor", "major"),
        default="patch",
        help="Semantic version bump strategy for map metadata.",
    )
    parser.add_argument("--summary", default=None, help="Reason summary to store in metadata.")
    parser.add_argument("--check", action="store_true", help="Check for drift without writing.")
    return parser.parse_args()


def should_ignore(rel_path: str) -> bool:
    if rel_path in {".git", ".DS_Store"}:
        return True
    parts = rel_path.split("/")
    blocked_parts = {
        ".git",
        "__pycache__",
        ".pytest_cache",
        ".mypy_cache",
        ".ruff_cache",
        ".venv",
        "node_modules",
    }
    if any(part in blocked_parts for part in parts):
        return True
    return False


def load_current_metadata(output_path: Path) -> tuple[str, str, str]:
    if not output_path.is_file():
        return ("0.0.0", "none", DEFAULT_SUMMARY)

    text = output_path.read_text(encoding="utf-8")
    version_match = VERSION_RE.search(text)
    bump_match = BUMP_RE.search(text)
    summary_match = SUMMARY_RE.search(text)
    version = version_match.group(1) if version_match else "0.0.0"
    bump = bump_match.group(1) if bump_match else "none"
    summary = summary_match.group(1).strip() if summary_match else DEFAULT_SUMMARY
    return (version, bump, summary)


def bump_version(version: str, bump: str) -> str:
    major_str, minor_str, patch_str = version.split(".")
    major = int(major_str)
    minor = int(minor_str)
    patch = int(patch_str)

    if bump == "none":
        return f"{major}.{minor}.{patch}"
    if bump == "patch":
        patch += 1
    elif bump == "minor":
        minor += 1
        patch = 0
    elif bump == "major":
        major += 1
        minor = 0
        patch = 0
    return f"{major}.{minor}.{patch}"


def build_tree_lines(repo_root: Path, output_rel: Path) -> list[str]:
    lines = ["."]

    def walk(directory: Path, prefix: str) -> None:
        entries: list[Node] = []
        for child in directory.iterdir():
            rel = child.relative_to(repo_root)
            rel_text = str(rel).replace(os.sep, "/")
            if rel == output_rel:
                continue
            if should_ignore(rel_text):
                continue
            entries.append(Node(name=child.name, path=child, is_dir=child.is_dir()))

        entries.sort(key=lambda node: (not node.is_dir, node.name.lower(), node.name))

        for idx, entry in enumerate(entries):
            is_last = idx == len(entries) - 1
            branch = "└── " if is_last else "├── "
            suffix = "/" if entry.is_dir else ""
            lines.append(f"{prefix}{branch}{entry.name}{suffix}")
            if entry.is_dir:
                next_prefix = f"{prefix}{'    ' if is_last else '│   '}"
                walk(entry.path, next_prefix)

    walk(repo_root, "")
    return lines


def count_entries(lines: Sequence[str]) -> tuple[int, int]:
    dir_count = 0
    file_count = 0
    for line in lines[1:]:
        if line.endswith("/"):
            dir_count += 1
        else:
            file_count += 1
    return dir_count, file_count


def render_map(
    repo_root: Path, output_rel: Path, map_version: str, summary: str, bump: str
) -> str:
    tree_lines = build_tree_lines(repo_root, output_rel)
    dir_count, file_count = count_entries(tree_lines)

    out_lines = [
        "# Repository Folder Map",
        "",
        "Generated file. Do not edit manually.",
        "",
        "## Metadata",
        f"- map_version: {map_version}",
        f"- bump: {bump}",
        f"- summary: {summary}",
        f"- directories: {dir_count}",
        f"- files: {file_count}",
        "",
        "## Tree",
        "```text",
        *tree_lines,
        "```",
        "",
    ]
    return "\n".join(out_lines)


def print_diff(existing: str, expected: str, output_rel: Path) -> None:
    diff = difflib.unified_diff(
        existing.splitlines(),
        expected.splitlines(),
        fromfile=f"{output_rel} (current)",
        tofile=f"{output_rel} (expected)",
        lineterm="",
    )
    for line in diff:
        print(line)


def main() -> int:
    args = parse_args()
    repo_root = Path(args.repo_root).resolve()
    output_rel = Path(args.output)
    output_path = (repo_root / output_rel).resolve()

    if not repo_root.is_dir():
        print(f"Repository root does not exist: {repo_root}", file=sys.stderr)
        return 2

    if output_path.is_dir():
        print(f"Output path is a directory: {output_path}", file=sys.stderr)
        return 2

    if not args.check:
        output_path.parent.mkdir(parents=True, exist_ok=True)

    current_version, current_bump, current_summary = load_current_metadata(output_path)
    summary = args.summary if args.summary is not None else current_summary
    if not summary.strip():
        summary = DEFAULT_SUMMARY

    if args.check:
        bump = current_bump
        target_version = current_version
    else:
        bump = args.bump
        target_version = bump_version(current_version, bump)

    expected = render_map(repo_root, output_rel, target_version, summary, bump)

    if args.check:
        if not output_path.is_file():
            print(f"Map file missing: {output_rel}", file=sys.stderr)
            return 1

        existing = output_path.read_text(encoding="utf-8")
        if existing == expected:
            print(f"Repo map is up to date: {output_rel}")
            return 0

        print(f"Repo map drift detected: {output_rel}", file=sys.stderr)
        print_diff(existing, expected, output_rel)
        return 1

    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(expected, encoding="utf-8")
    print(f"Wrote folder map: {output_rel}")
    print(f"map_version={target_version}")
    print(f"summary={summary}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
