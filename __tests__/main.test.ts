import { expect, test } from '@jest/globals'

test('junk', () => {
  const re = new RegExp(`\\[([A-Z]{2,}-\\d{3,})\\]`, 'g')
  expect(re.test('foo')).toEqual(false)
  expect(re.test('FOO-123')).toEqual(false)
  expect(re.test('[FOO-123]')).toEqual(true)
})
