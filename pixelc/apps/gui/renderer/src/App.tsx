import { useEffect, useState } from 'react';
import type { CompileEvent } from '@shared/types';
import { toCompileArgs, validateForm, type FormState } from './lib/config';

const defaultState: FormState = {
  inputPath: '',
  outDir: '',
  preset: 'unity',
  connectivity: 4,
  padding: 0,
  pivot: 'center',
  power2: false,
  fps: 12,
  batch: false,
  ignoreRaw: '',
  report: true,
  noOverwrite: false,
  dryRun: false
};

export function App(): JSX.Element {
  const [form, setForm] = useState<FormState>(defaultState);
  const [logs, setLogs] = useState<string[]>([]);
  const [progress, setProgress] = useState(0);
  const [atlasJson, setAtlasJson] = useState('');
  const [atlasPath, setAtlasPath] = useState('');
  const [reportJson, setReportJson] = useState('');
  const [doctorText, setDoctorText] = useState('Running doctor...');
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    void window.pixelc.doctor().then((r) => setDoctorText(r.text));
  }, []);

  const errs = validateForm(form);

  async function loadOutputs(outDir: string): Promise<void> {
    const atlasJsonPath = `${outDir}/atlas.json`;
    const atlasPngPath = `${outDir}/atlas.png`;
    const reportPath = `${outDir}/report.json`;
    if (await window.pixelc.fileExists(atlasJsonPath)) setAtlasJson(await window.pixelc.readTextFile(atlasJsonPath));
    if (await window.pixelc.fileExists(reportPath)) setReportJson(await window.pixelc.readTextFile(reportPath));
    setAtlasPath(atlasPngPath);
  }

  async function onCompile(): Promise<void> {
    if (errs.length) return;
    setLogs([]);
    setBusy(true);
    setProgress(0);
    try {
      await window.pixelc.compile(toCompileArgs(form), async (ev: CompileEvent) => {
        if (ev.type === 'log') setLogs((prev) => [...prev, ev.line]);
        if (ev.type === 'progress') setProgress(ev.value);
        if (ev.type === 'done') {
          setProgress(100);
          await loadOutputs(ev.outDir);
        }
        if (ev.type === 'error') setLogs((prev) => [...prev, `ERROR: ${ev.message}`]);
      });
    } catch (e) {
      setLogs((prev) => [...prev, `ERROR: ${(e as Error).message}`]);
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="layout">
      <section className="left">
        <h1>Pixel Asset Compiler</h1>
        <p className="doctor">Doctor: {doctorText}</p>
        <button onClick={async () => setForm({ ...form, inputPath: (await window.pixelc.pickInput()) ?? form.inputPath })}>Browse Input</button>
        <input value={form.inputPath} onChange={(e) => setForm({ ...form, inputPath: e.target.value })} placeholder="Input path" />
        <button onClick={async () => setForm({ ...form, outDir: (await window.pixelc.pickOutDir()) ?? form.outDir })}>Browse Output</button>
        <input value={form.outDir} onChange={(e) => setForm({ ...form, outDir: e.target.value })} placeholder="Output dir" />
        <label>Connectivity<select value={form.connectivity} onChange={(e) => setForm({ ...form, connectivity: Number(e.target.value) as 4 | 8 })}><option value={4}>4</option><option value={8}>8</option></select></label>
        <label>Padding<input type="number" value={form.padding} onChange={(e) => setForm({ ...form, padding: Number(e.target.value) })} /></label>
        <label>Pivot<select value={form.pivot} onChange={(e) => setForm({ ...form, pivot: e.target.value as 'center' | 'bottom-center' })}><option value="center">center</option><option value="bottom-center">bottom-center</option></select></label>
        <label>FPS<input type="number" value={form.fps} onChange={(e) => setForm({ ...form, fps: Number(e.target.value) })} /></label>
        <label>Ignore<input value={form.ignoreRaw} onChange={(e) => setForm({ ...form, ignoreRaw: e.target.value })} placeholder="*.tmp,*.psd" /></label>
        <label><input type="checkbox" checked={form.power2} onChange={(e) => setForm({ ...form, power2: e.target.checked })} />power2</label>
        <label><input type="checkbox" checked={form.batch} onChange={(e) => setForm({ ...form, batch: e.target.checked })} />batch</label>
        <label><input type="checkbox" checked={form.report} onChange={(e) => setForm({ ...form, report: e.target.checked })} />report</label>
        <button disabled={busy || errs.length > 0} onClick={onCompile}>Compile</button>
        {errs.length > 0 && <ul>{errs.map((e) => <li key={e}>{e}</li>)}</ul>}
      </section>
      <section className="right">
        <progress value={progress} max={100} />
        <button onClick={() => navigator.clipboard.writeText(logs.join('\n'))}>Copy logs</button>
        <button onClick={() => window.pixelc.openPath(form.outDir)}>Open output folder</button>
        <pre>{logs.join('\n')}</pre>
        <h3>Atlas Preview</h3>
        {atlasPath && <img src={`file://${atlasPath}`} alt="atlas" style={{ maxWidth: 280 }} />}
        <h3>Atlas JSON</h3>
        <pre>{atlasJson}</pre>
        <h3>Report</h3>
        <pre>{reportJson}</pre>
      </section>
    </div>
  );
}
