import { _electron as electron, expect, test } from '@playwright/test';
import { mkdirSync, mkdtempSync, readFileSync, writeFileSync, existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { PNG } from 'pngjs';
import { spawnSync } from 'node:child_process';

function writePng(filePath: string): void {
  const png = new PNG({ width: 2, height: 2 });
  png.data.fill(255);
  writeFileSync(filePath, PNG.sync.write(png));
}

test('compiles sample input through GUI', async () => {
  const temp = mkdtempSync(join(tmpdir(), 'pixelc-gui-'));
  const input = join(temp, 'sheet.png');
  const out = join(temp, 'out');
  mkdirSync(out, { recursive: true });
  writePng(input);

  mkdirSync('.bin', { recursive: true });
  const exe = process.platform === 'win32' ? 'pixelc.exe' : 'pixelc';
  const go = spawnSync('go', ['build', '-o', join('.bin', exe), '../../../../cmd/pixelc'], { cwd: process.cwd(), stdio: 'inherit' });
  expect(go.status).toBe(0);

  const app = await electron.launch({ args: ['.'], cwd: process.cwd() });
  const page = await app.firstWindow();
  await page.getByPlaceholder('Input path').fill(input);
  await page.getByPlaceholder('Output dir').fill(out);
  await page.getByText('Compile').click();

  await expect.poll(() => existsSync(join(out, 'atlas.json'))).toBeTruthy();
  await expect.poll(() => existsSync(join(out, 'atlas.png'))).toBeTruthy();
  const atlas = JSON.parse(readFileSync(join(out, 'atlas.json'), 'utf-8'));
  expect(Array.isArray(atlas.sprites)).toBeTruthy();
  await app.close();
});
