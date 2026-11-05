import assert from 'node:assert/strict'
import { createHash } from 'node:crypto'
import { readFile } from 'node:fs/promises'
import path from 'node:path'

const root = path.resolve(import.meta.dirname, '..')
const manifest = JSON.parse(await readFile(path.join(root, 'data/manifest.json'), 'utf8'))
const ids = new Set()
let count = 0
const sha = (data) => createHash('sha256').update(data).digest('hex')
const safe = (relative) => {
  assert.ok(relative && !path.isAbsolute(relative) && !relative.split('/').includes('..') && !relative.includes('\\'))
  return path.join(root, relative)
}
for (const descriptor of manifest.catalogs) {
  const bytes = await readFile(safe('data/' + descriptor.path))
  const catalog = JSON.parse(bytes)
  assert.equal(sha(bytes), descriptor.sha256, descriptor.type + ' digest')
  assert.equal(catalog.version, descriptor.version)
  assert.equal(catalog.resources.length, descriptor.count)
  for (const item of catalog.resources) {
    assert.ok(!ids.has(item.id), item.id); ids.add(item.id)
    for (const [file, size, digest] of [['assetPath', 'byteSize', 'sha256'], ['previewPath', 'previewByteSize', 'previewSha256']]) {
      const data = await readFile(safe(item[file]))
      assert.equal(data.length, item[size], item.id + ' ' + size)
      assert.equal(sha(data), item[digest], item.id + ' ' + digest)
    }
    if (item.type === 'html-template') {
      const payload = JSON.parse(await readFile(safe(item.assetPath), 'utf8'))
      assert.equal(payload.id, item.id)
      assert.equal(payload.name, item.name)
      assert.ok(payload.html && payload.layoutPrompt, item.id)
      assert.ok(!/<h1\b/i.test(payload.html), item.id + ': body must omit metadata title')
      const document = await readFile(safe(item.details.documentPath), 'utf8')
      assert.ok(document.includes('data-no-dark'), item.id + ': missing source protection')
    }
    count++
  }
}
console.log(`Template market validated: ${count} resources; catalog/asset/preview hashes match.`)
