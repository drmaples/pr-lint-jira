import { expect, test } from '@jest/globals'
import { env } from 'process'
import { defaultTitleBodyRegex, run } from '../src/main'

test('test default title body regex', () => {
  const re = new RegExp(defaultTitleBodyRegex.source)
  expect(re.test('foo')).toEqual(false) // do better
  expect(re.test('FOO-123')).toEqual(false) // no brackets
  expect(re.test('[f-123]')).toEqual(false) // min 2 char slug
  expect(re.test('[hi-]')).toEqual(false) // need a digit after the slug
  expect(re.test('nothing in here')).toEqual(false) // no tix at all
  expect(re.test('what [is-going] on here')).toEqual(false) // no tix at all

  expect(re.test('[HI-7]')).toEqual(true)
  expect(re.test('[FOO-123]')).toEqual(true)
  expect(re.test('[zzz-123]')).toEqual(true)
  expect(re.test('[zZz-123]')).toEqual(true)
  expect(re.test('[XXXXXXXXXX-0000000000]')).toEqual(true)
  expect(re.test('start junk [zZz-123]')).toEqual(true)
  expect(re.test('[zZz-123] end junk')).toEqual(true)
  expect(re.test('start junk [zZz-123] end junk')).toEqual(true)
})

test('test runs', async () => {
  env['IS_CI'] = 'true'
  env['INPUT_TOKEN'] = 'xxxxxx'

  try {
    await run()
  } catch (e) {
    console.error('>>>>>', e)
  }
})
