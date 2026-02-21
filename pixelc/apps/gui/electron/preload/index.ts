import { contextBridge, ipcRenderer } from 'electron';
import type { CompileArgs, CompileEvent, PixelcAPI } from '../shared/types';

const api: PixelcAPI = {
  pickInput: () => ipcRenderer.invoke('pick-input'),
  pickOutDir: () => ipcRenderer.invoke('pick-out-dir'),
  openPath: (path: string) => ipcRenderer.invoke('open-path', path),
  doctor: () => ipcRenderer.invoke('doctor'),
  compile: async (args: CompileArgs, onEvent: (ev: CompileEvent) => void) => {
    const listener = (_event: Electron.IpcRendererEvent, payload: CompileEvent) => onEvent(payload);
    ipcRenderer.on('compile-event', listener);
    try {
      await ipcRenderer.invoke('compile', args);
    } finally {
      ipcRenderer.removeListener('compile-event', listener);
    }
  },
  readTextFile: (path: string) => ipcRenderer.invoke('read-text-file', path),
  fileExists: (path: string) => ipcRenderer.invoke('file-exists', path)
};

contextBridge.exposeInMainWorld('pixelc', api);
