import type { PixelcAPI } from '@shared/types';

declare global {
  interface Window {
    pixelc: PixelcAPI;
  }
}

export {};
