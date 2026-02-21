import { describe, expect, it } from 'vitest';
import { validateCompileArgs } from './compile';

describe('compile validation', () => {
  it('rejects unsupported noOverwrite flag', () => {
    expect(() =>
      validateCompileArgs({
        inputPath: 'a',
        outDir: 'b',
        preset: 'unity',
        connectivity: 4,
        padding: 0,
        pivot: 'center',
        power2: false,
        fps: 12,
        batch: false,
        ignore: [],
        report: false,
        noOverwrite: true,
        dryRun: false
      })
    ).toThrowError(/noOverwrite/);
  });
});
