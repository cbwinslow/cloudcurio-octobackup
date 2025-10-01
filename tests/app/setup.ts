import { expect, afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'
import { vi } from 'vitest'
import matchers from '@testing-library/jest-dom/matchers'
import '@testing-library/jest-dom'

// Extend Vitest expect with jest-dom matchers
expect.extend(matchers)

// Cleanup after each test
afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})
