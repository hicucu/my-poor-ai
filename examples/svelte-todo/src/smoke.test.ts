import { describe, it, expect } from 'vitest';

describe('vitest pipeline', () => {
  it('runs a trivial assertion', () => {
    expect(1 + 1).toBe(2);
  });
});
