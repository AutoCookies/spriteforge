import type { CompileArgs } from '../shared/types';

export function validateCompileArgs(args: CompileArgs): void {
  if (!args.inputPath || !args.outDir) throw new Error('inputPath and outDir are required');
  if (args.preset !== 'unity') throw new Error('Only unity preset is supported');
  if (![4, 8].includes(args.connectivity)) throw new Error('connectivity must be 4 or 8');
  if (args.noOverwrite) throw new Error('noOverwrite is not supported by current pixelc CLI');
}
