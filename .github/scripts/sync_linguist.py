"""Resolve a validated Linguist release before installing tools or cloning it."""
import os
from pathlib import Path
import re
import subprocess
import sys

TAG = re.compile(r"v[0-9]+\.[0-9]+\.[0-9]+")
SHA = re.compile(r"[0-9a-f]{40}")
GENERATOR = Path("internal/code-generator/generator/generator_test.go")


def validate_tag(tag):
    if not TAG.fullmatch(tag):
        raise ValueError("Expected a stable Linguist tag such as v9.5.0")
    return tag


def resolve_commit(refs, tag):
    refs = dict(line.split()[::-1] for line in refs.splitlines() if line.strip())
    commit = refs.get(f"refs/tags/{tag}^{{}}", refs.get(f"refs/tags/{tag}", ""))
    if not SHA.fullmatch(commit):
        raise ValueError("Release tag did not resolve to a commit")
    return commit


def check():
    tag = os.environ.get("LINGUIST_TAG", "") or subprocess.check_output(
        ["gh", "api", "repos/github/linguist/releases/latest", "--jq", ".tag_name"],
        text=True,
    ).strip()
    validate_tag(tag)
    refs = subprocess.check_output(
        ["git", "ls-remote", "https://github.com/github/linguist.git",
         f"refs/tags/{tag}", f"refs/tags/{tag}^{{}}"], text=True,
    )
    commit = resolve_commit(refs, tag)
    previous = re.search(r'\bcommit\s*=\s*"([0-9a-f]{40})"', GENERATOR.read_text()).group(1)
    with open(os.environ["GITHUB_OUTPUT"], "a") as output:
        output.write(f"tag={tag}\ncommit={commit}\nchanged={str(commit != previous).lower()}\n")


def update():
    tag = validate_tag(os.environ["LINGUIST_TAG"])
    commit = os.environ["LINGUIST_COMMIT"]
    if not SHA.fullmatch(commit):
        raise ValueError("Invalid Linguist commit")
    source, count = re.subn(r'(\bcommit\s*=\s*)"[0-9a-f]{40}"',
                            lambda m: f'{m[1]}"{commit}"', GENERATOR.read_text())
    if count != 1:
        raise ValueError("Expected exactly one generator fixture commit")
    readme = Path("README.md")
    description, count = re.subn(r'(version \*\*)v[0-9]+\.[0-9]+\.[0-9]+(\*\*\.)',
                                 lambda m: f"{m[1]}{tag}{m[2]}", readme.read_text())
    if count != 1:
        raise ValueError("Expected exactly one Linguist version in README")
    GENERATOR.write_text(source)
    readme.write_text(description)


if __name__ == "__main__":
    {"check": check, "update": update}[sys.argv[1]]()
