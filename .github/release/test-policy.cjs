const assert = require('node:assert/strict');
const test = require('node:test');
const fs = require('node:fs');
const path = require('node:path');
const {Manifest} = require('release-please/build/src/manifest');
const {buildStrategy} = require('release-please/build/src/factory');
const {parseConventionalCommits} = require('release-please/build/src/commit');
const {Version} = require('release-please/build/src/version');
const {TagName} = require('release-please/build/src/util/tag-name');
const {affectsNativeLibrary, NativeMinorRelease} = require('./native-policy.cjs');
const root = path.resolve(__dirname, '../..');
const github = {repository: {owner: 'example', repo: 'enry'},
  getFileJson: async file => JSON.parse(fs.readFileSync(path.join(root, file)))};

async function propose(commits, version = '0.3.0') {
  const manifest = await Manifest.fromManifest(github, 'master',
    '.github/release-please-config.json', '.github/.release-please-manifest.json');
  const strategy = await buildStrategy({...manifest.repositoryConfig['.'], github,
    path: '.', targetBranch: 'master'});
  const plugin = new NativeMinorRelease(github, 'master', manifest.repositoryConfig);
  const byPath = {'.': commits};
  await plugin.preconfigure({'.': strategy}, byPath, {});
  return strategy.buildReleasePullRequest(parseConventionalCommits(byPath['.']),
    {tag: new TagName(Version.parse(version), 'python', '-', true), sha: 'previous', notes: ''});
}
const commit = (message, files) => ({sha: 'abc123', message, files});

test('native input boundaries', () => {
  for (const file of ['common.go','data/content.go','data/rule/rule.go','shared/enry.go',
    'regex/standard.go','internal/tokenizer/tokenizer.go','shared/helper.c','go.mod','go.sum']) {
    assert.equal(affectsNativeLibrary(file), true, file);
  }
  for (const file of ['common_test.go','data/heuristics_test.go','README.md','java/Enry.java',
    'cmd/enry/main.go','internal/code-generator/main.go','.github/workflows/goTest.yml',
    '_testdata/Go/main.go','python/setup.py']) {
    assert.equal(affectsNativeLibrary(file), false, file);
  }
});
for (const [name, message, files, expected] of [
  ['native fix','fix(detector): handle marker',['common.go'],'0.4.0'],
  ['Linguist sync','chore(data): update Linguist definitions',['data/content.go','README.md'],'0.4.0'],
  ['unconventional upstream merge','Merge upstream release',['common.go'],'0.4.0'],
  ['Python fix','fix(python): correct loader',['python/enry/definitions.py'],'0.3.1'],
  ['Python feature','feat(python): add API',['python/enry/__init__.py'],'0.4.0'],
  ['docs only','docs: explain bindings',['README.md'],null],
  ['CI only','fix(ci): improve tests',['.github/workflows/goTest.yml'],null],
  ['Java only','fix(java): repair ABI',['java/src/main/Enry.java'],null],
  ['CLI only','feat(cli): add flag',['cmd/enry/main.go'],null],
  ['generator only','fix(generator): adjust template',['internal/code-generator/main.go'],null],
  ['native tests only','fix(tests): improve assertion',['common_test.go'],null],
  ['release metadata','chore(python): release 0.4.0',['python/CHANGELOG.md','.github/.release-please-manifest.json'],null],
]) {
  test(name, async () => {
    const pr = await propose([commit(message,files)]);
    assert.equal(pr?.version.toString() ?? null, expected);
    if (pr) assert.deepEqual(pr.updates.map(update=>update.path),['python/CHANGELOG.md']);
  });
}
test('multiple native merges produce one minor bump; next release advances once', async () => {
  const commits = [commit('Updated data',['data/content.go']),commit('fix: correct parser',['common.go'])];
  assert.equal((await propose(commits)).version.toString(),'0.4.0');
  assert.equal((await propose(commits,'0.4.0')).version.toString(),'0.5.0');
  assert.equal(commits.length,2);
});
test('breaking changes retain precedence after 1.0', async () => {
  const pr=await propose([commit('feat!: remove API',['common.go'])],'1.2.3');
  assert.equal(pr.version.toString(),'2.0.0');
});

test('release outputs preserve the workflow dispatch contract', () => {
  const {writeReleaseOutputs} = require('./run.cjs');
  const directory=fs.mkdtempSync(path.join(require('node:os').tmpdir(),'enry-release-'));
  const output=path.join(directory,'output');
  try {
    assert.equal(writeReleaseOutputs([],output),false);
    assert.equal(fs.existsSync(output),false);
    assert.equal(writeReleaseOutputs([{tagName:'python-v0.4.0'}],output),true);
    assert.equal(fs.readFileSync(output,'utf8'),'release_created=true\ntag=python-v0.4.0\n');
    assert.throws(()=>writeReleaseOutputs([{tagName:'python-v0.4.0\nunsafe=true'}],output));
    assert.throws(()=>writeReleaseOutputs([{tagName:'v0.4.0'}],output));
    assert.throws(()=>writeReleaseOutputs([{},{}],output));
    const workflow=fs.readFileSync(path.join(root,'.github/workflows/python-release.yml'),'utf8');
    assert.ok(workflow.includes('created: ${{ steps.release.outputs.release_created }}'));
    assert.ok(workflow.includes('tag: ${{ steps.release.outputs.tag }}'));
  } finally {fs.rmSync(directory,{recursive:true});}
});
