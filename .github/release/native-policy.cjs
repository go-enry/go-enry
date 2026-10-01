const {ManifestPlugin} = require('release-please/build/src/plugin');

// Only inputs shipped in the Python native library. CLI, generator sources,
// tests, Java, docs and workflow changes do not request a native minor bump.
function affectsNativeLibrary(path) {
  if (path === 'go.mod' || path === 'go.sum') return true;
  if (path.startsWith('internal/code-generator/') || path.startsWith('cmd/')) return false;
  if (path.endsWith('_test.go')) return false;
  if (!/\.(go|c|h)$/.test(path)) return false;
  return !path.includes('/') || /^(shared|data|regex|internal)\//.test(path);
}

class NativeMinorRelease extends ManifestPlugin {
  async preconfigure(strategies, commitsByPath) {
    const commits = (commitsByPath['.'] || []).filter(commit =>
      (commit.files || []).some(path => affectsNativeLibrary(path) ||
        path.startsWith('python/') || path === '.github/workflows/python-wheels.yml'));
    commitsByPath['.'] = commits;
    const nativeCommit = commits.find(commit => (commit.files || []).some(affectsNativeLibrary));
    if (nativeCommit) {
      // This is release-planning metadata, not a new Git commit. It makes
      // non-conventional upstream merges and chore(data) updates releasable,
      // while retaining original breaking-change notes and normal versioning.
      commits.push({
        sha: nativeCommit.sha,
        message: 'feat(native): refresh bundled detector and language data',
        files: nativeCommit.files,
      });
    }
    return strategies;
  }
}
module.exports = {affectsNativeLibrary, NativeMinorRelease};
