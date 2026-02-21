import type { CompileArgs } from '@shared/types';

export type FormState = {
  inputPath: string;
  outDir: string;
  preset: 'unity';
  connectivity: 4 | 8;
  padding: number;
  pivot: 'center' | 'bottom-center';
  power2: boolean;
  fps: number;
  batch: boolean;
  ignoreRaw: string;
  report: boolean;
  noOverwrite: boolean;
  dryRun: boolean;
};

export function validateForm(state: FormState): string[] {
  const errs: string[] = [];
  if (!state.inputPath) errs.push('Input path is required');
  if (!state.outDir) errs.push('Output directory is required');
  if (state.padding < 0) errs.push('Padding must be >= 0');
  if (state.fps <= 0) errs.push('FPS must be > 0');
  return errs;
}

export function toCompileArgs(state: FormState): CompileArgs {
  return {
    inputPath: state.inputPath,
    outDir: state.outDir,
    preset: state.preset,
    connectivity: state.connectivity,
    padding: state.padding,
    pivot: state.pivot,
    power2: state.power2,
    fps: state.fps,
    batch: state.batch,
    ignore: state.ignoreRaw
      .split(',')
      .map((v) => v.trim())
      .filter(Boolean),
    report: state.report,
    noOverwrite: state.noOverwrite,
    dryRun: state.dryRun
  };
}
