export type CompileArgs = {
  inputPath: string;
  outDir: string;
  preset: 'unity';
  connectivity: 4 | 8;
  padding: number;
  pivot: 'center' | 'bottom-center';
  power2: boolean;
  fps: number;
  batch: boolean;
  ignore: string[];
  report: boolean;
  noOverwrite: boolean;
  dryRun: boolean;
};

export type CompileEvent =
  | { type: 'log'; line: string }
  | { type: 'progress'; value: number }
  | { type: 'done'; code: 0 | 1 | 2 | 3; outDir: string }
  | { type: 'error'; message: string };

export interface PixelcAPI {
  pickInput(): Promise<string | null>;
  pickOutDir(): Promise<string | null>;
  openPath(path: string): Promise<void>;
  doctor(): Promise<{ ok: boolean; text: string }>;
  compile(args: CompileArgs, onEvent: (ev: CompileEvent) => void): Promise<void>;
  readTextFile(path: string): Promise<string>;
  fileExists(path: string): Promise<boolean>;
}
