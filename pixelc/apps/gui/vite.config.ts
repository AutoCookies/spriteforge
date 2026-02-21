import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'node:path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@renderer': resolve(__dirname, 'renderer/src'),
      '@shared': resolve(__dirname, 'electron/shared')
    }
  },
  root: 'renderer',
  build: {
    outDir: '../dist/renderer',
    emptyOutDir: true
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./renderer/src/test-setup.ts']
  }
});
