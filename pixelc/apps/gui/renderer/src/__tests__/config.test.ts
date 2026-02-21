import { describe, expect, it } from 'vitest';
import { toCompileArgs, validateForm, type FormState } from '../lib/config';

const base: FormState = {
  inputPath: '/tmp/in',
  outDir: '/tmp/out',
  preset: 'unity',
  connectivity: 4,
  padding: 0,
  pivot: 'center',
  power2: false,
  fps: 12,
  batch: false,
  ignoreRaw: '*.tmp,*.bak',
  report: true,
  noOverwrite: false,
  dryRun: false
};

describe('form config', () => {
  it('serializes ignore list into CompileArgs', () => {
    const args = toCompileArgs(base);
    expect(args.ignore).toEqual(['*.tmp', '*.bak']);
  });

  it('validates required fields', () => {
    const errors = validateForm({ ...base, inputPath: '', outDir: '' });
    expect(errors).toContain('Input path is required');
    expect(errors).toContain('Output directory is required');
  });
});
