import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { App } from '../App';

const pixelc = {
  pickInput: vi.fn(async () => '/in'),
  pickOutDir: vi.fn(async () => '/out'),
  openPath: vi.fn(async () => undefined),
  doctor: vi.fn(async () => ({ ok: true, text: 'ok' })),
  compile: vi.fn(async (_args, onEvent) => onEvent({ type: 'done', code: 0, outDir: '/out' })),
  readTextFile: vi.fn(async () => '{"sprites":[]}'),
  fileExists: vi.fn(async () => false)
};

describe('app transitions', () => {
  it('disables compile when required fields are missing', () => {
    (window as unknown as { pixelc: typeof pixelc }).pixelc = pixelc;
    render(<App />);
    expect(screen.getByText('Compile')).toBeDisabled();
  });

  it('enables compile after selecting paths', async () => {
    (window as unknown as { pixelc: typeof pixelc }).pixelc = pixelc;
    render(<App />);
    fireEvent.click(screen.getByText('Browse Input'));
    fireEvent.click(screen.getByText('Browse Output'));
    await vi.waitFor(() => expect(screen.getByText('Compile')).not.toBeDisabled());
  });
});
