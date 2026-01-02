import { expect, test } from '@jest/globals'
import * as process from 'process'
import * as cp from 'child_process'
import * as path from 'path'

test('test regex', () => {
  const re = new RegExp(`\\[([A-Z]{2,}-\\d+)\\]`, 'g')
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
    console.log(error)
  }
})
