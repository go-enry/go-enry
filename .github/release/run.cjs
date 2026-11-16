const fs = require('node:fs');
const {GitHub} = require('release-please/build/src/github');
const {Manifest} = require('release-please/build/src/manifest');
const {registerPlugin} = require('release-please/build/src/factories/plugin-factory');
const {NativeMinorRelease} = require('./native-policy.cjs');

async function main() {
  const [owner, repo, extra] = (process.env.GITHUB_REPOSITORY || '').split('/');
  if (!owner || !repo || extra || !process.env.GH_TOKEN) throw new Error('Repository and token are required');
  registerPlugin('native-minor', options => new NativeMinorRelease(
    options.github, options.targetBranch, options.repositoryConfig));
  const github = await GitHub.create({owner, repo, token: process.env.GH_TOKEN});
  const load = () => Manifest.fromManifest(github, 'master',
    '.github/release-please-config.json', '.github/.release-please-manifest.json',
    {plugins: ['native-minor']});
  if (process.argv.includes('--dry-run')) {
    const prs = await (await load()).buildPullRequests();
    console.log(JSON.stringify(prs.map(pr => ({title: pr.title.toString(),
      version: pr.version.toString(), files: pr.updates.map(update => update.path)}))));
    return;
  }
  if (process.env.GITHUB_REF !== 'refs/heads/master') throw new Error('Releases require master');
  await runReleaseCycle(load, process.env.GITHUB_OUTPUT);
}
async function runReleaseCycle(load, outputFile) {
  const releases = (await (await load()).createReleases()).filter(Boolean);
  writeReleaseOutputs(releases, outputFile);
  // Reload after creating tags so newer commits are compared with the new release.
  // Always reconcile the next PR, including when queued pushes were coalesced.
  await (await load()).createPullRequests();
}
function writeReleaseOutputs(releases, outputFile) {
  if (releases.length > 1) throw new Error('Expected at most one Python release');
  if (!releases.length) return false;
  const tag = releases[0].tagName;
  if (!/^python-v\d+\.\d+\.\d+$/.test(tag)) throw new Error('Unexpected release tag');
  fs.appendFileSync(outputFile, `release_created=true\ntag=${tag}\n`);
  return true;
}
module.exports = {writeReleaseOutputs, runReleaseCycle};
if (require.main === module) {
  main().catch(error => { console.error(error); process.exitCode = 1; });
}
