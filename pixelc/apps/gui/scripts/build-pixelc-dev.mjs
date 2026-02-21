import { spawnSync } from 'node:child_process';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { mkdirSync } from 'node:fs';

const root = join(dirname(fileURLToPath(import.meta.url)), '..', '..', '..', '..');
const outDir = join(root, 'pixelc', 'apps', 'gui', '.bin');
mkdirSync(outDir, { recursive: true });
const exe = process.platform === 'win32' ? 'pixelc.exe' : 'pixelc';
const out = join(outDir, exe);
const res = spawnSync('go', ['build', '-o', out, './cmd/pixelc'], { cwd: root, stdio: 'inherit' });
if (res.status !== 0) {
  process.exit(res.status ?? 1);
}
