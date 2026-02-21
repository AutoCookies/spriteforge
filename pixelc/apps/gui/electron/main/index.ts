import { app, BrowserWindow, dialog, ipcMain, shell } from 'electron';
import { existsSync } from 'node:fs';
import { readFile } from 'node:fs/promises';
import { join } from 'node:path';
import { spawn } from 'node:child_process';
import type { CompileArgs } from '../shared/types';
const RENDERER_DIST = join(__dirname, '../../renderer/index.html');

export function validateCompileArgs(args: CompileArgs): void {
  if (!args.inputPath || !args.outDir) throw new Error('inputPath and outDir are required');
  if (args.preset !== 'unity') throw new Error('Only unity preset is supported');
  if (![4, 8].includes(args.connectivity)) throw new Error('connectivity must be 4 or 8');
  if (args.noOverwrite) throw new Error('noOverwrite is not supported by current pixelc CLI');
}

function resolvePixelcBinary(): string {
  const exe = process.platform === 'win32' ? 'pixelc.exe' : 'pixelc';
  const bundled = join(process.resourcesPath, 'bin', exe);
  if (existsSync(bundled)) return bundled;
  const devLocal = join(app.getAppPath(), '.bin', exe);
  if (existsSync(devLocal)) return devLocal;
  return 'pixelc';
}

async function createWindow(): Promise<void> {
  const win = new BrowserWindow({
    width: 1220,
    height: 760,
    webPreferences: {
      preload: join(__dirname, '../preload/index.cjs'),
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  if (process.env.VITE_DEV_SERVER_URL) {
    await win.loadURL(process.env.VITE_DEV_SERVER_URL);
  } else {
    await win.loadFile(RENDERER_DIST);
  }
}

function wireIpc(): void {
  ipcMain.handle('pick-input', async () => {
    const result = await dialog.showOpenDialog({ properties: ['openFile', 'openDirectory'], filters: [{ name: 'PNG', extensions: ['png'] }] });
    return result.canceled ? null : result.filePaths[0];
  });

  ipcMain.handle('pick-out-dir', async () => {
    const result = await dialog.showOpenDialog({ properties: ['openDirectory', 'createDirectory'] });
    return result.canceled ? null : result.filePaths[0];
  });

  ipcMain.handle('open-path', async (_e, path: string) => {
    await shell.openPath(path);
  });

  ipcMain.handle('doctor', async () => {
    const bin = resolvePixelcBinary();
    return await new Promise<{ ok: boolean; text: string }>((resolve) => {
      const child = spawn(bin, ['doctor']);
      let text = '';
      child.stdout.on('data', (d) => (text += d.toString()));
      child.stderr.on('data', (d) => (text += d.toString()));
      child.on('close', (code) => resolve({ ok: code === 0, text: text.trim() }));
      child.on('error', (err) => resolve({ ok: false, text: err.message }));
    });
  });

  ipcMain.handle('read-text-file', async (_e, path: string) => readFile(path, 'utf-8'));
  ipcMain.handle('file-exists', async (_e, path: string) => existsSync(path));

  ipcMain.handle('compile', async (e, args: CompileArgs) => {
    validateCompileArgs(args);
    const bin = resolvePixelcBinary();
    const cliArgs = [
      'compile',
      args.inputPath,
      '--out',
      args.outDir,
      '--preset',
      args.preset,
      '--connectivity',
      String(args.connectivity),
      '--padding',
      String(args.padding),
      '--pivot',
      args.pivot,
      '--fps',
      String(args.fps)
    ];
    if (args.power2) cliArgs.push('--power2');
    if (args.batch) cliArgs.push('--batch');
    if (args.report) cliArgs.push('--report');
    if (args.dryRun) cliArgs.push('--dry-run');
    for (const pattern of args.ignore) cliArgs.push('--ignore', pattern);

    await new Promise<void>((resolve, reject) => {
      const child = spawn(bin, cliArgs);
      let seen = 0;
      const emitLine = (line: string) => {
        if (!line.trim()) return;
        e.sender.send('compile-event', { type: 'log', line });
        seen += 1;
        e.sender.send('compile-event', { type: 'progress', value: Math.min(95, seen * 5) });
      };
      child.stdout.on('data', (d) => d.toString().split('\n').forEach(emitLine));
      child.stderr.on('data', (d) => d.toString().split('\n').forEach(emitLine));
      child.on('close', (code) => {
        const mapped = code === 0 ? 0 : 1;
        e.sender.send('compile-event', { type: 'done', code: mapped, outDir: args.outDir });
        resolve();
      });
      child.on('error', (err) => {
        e.sender.send('compile-event', { type: 'error', message: err.message });
        reject(err);
      });
    });
  });
}

app.whenReady().then(async () => {
  wireIpc();
  await createWindow();
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit();
});
