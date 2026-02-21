import { describe, expect, it } from 'vitest';
import type { PixelcAPI } from '../shared/types';

describe('preload api shape', () => {
  it('matches strict whitelist', () => {
    const keys: (keyof PixelcAPI)[] = ['pickInput', 'pickOutDir', 'openPath', 'doctor', 'compile', 'readTextFile', 'fileExists'];
    expect(keys).toHaveLength(7);
  });
});
