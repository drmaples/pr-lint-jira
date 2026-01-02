import { expect, test } from '@jest/globals'
import * as process from 'process'
import * as cp from 'child_process'
import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url) // get the resolved path to the file
const __dirname = path.dirname(__filename) // get the name of the directory

test('test regex', () => {
  const re = new RegExp('\\[([A-Z]{2,}-\\d+)\\]', 'g')
  expect(re.test('foo')).toEqual(false)
  expect(re.test('FOO-123')).toEqual(false)
  expect(re.test('[FOO-123]')).toEqual(true)
})

test('test runs', () => {
  process.env['IS_CI'] = 'true'
  process.env['INPUT_TOKEN'] = 'xxxxxx'
  const ip = path.join(__dirname, '..', 'lib', 'main.js')
  try {
    cp.execSync(`node ${ip}`, { env: process.env })
  } catch (error) {
    console.info(error)
  }
})
