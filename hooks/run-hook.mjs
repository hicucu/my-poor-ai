#!/usr/bin/env node
/**
 * Cross-platform hook runner.
 * Finds bash (Git Bash on Windows, system bash on Unix) and runs the named
 * hook script in the same directory.
 *
 * Usage: node run-hook.mjs <script-name>
 */
import { spawnSync } from 'child_process';
import { fileURLToPath } from 'url';
import { dirname, join } from 'path';
import { existsSync } from 'fs';

const HOOK_DIR = dirname(fileURLToPath(import.meta.url));
const scriptName = process.argv[2];

if (!scriptName) {
  process.stderr.write('run-hook.mjs: missing script name\n');
  process.exit(1);
}

// Only allow plain filenames — reject path separators and dot-prefixed names
// so the runner can never execute anything outside HOOK_DIR.
if (!/^[A-Za-z0-9][A-Za-z0-9._-]*$/.test(scriptName)) {
  process.stderr.write(`run-hook.mjs: invalid script name: ${scriptName}\n`);
  process.exit(1);
}

const scriptPath = join(HOOK_DIR, scriptName);

function findBash() {
  if (process.platform !== 'win32') return 'bash';

  const candidates = [
    'C:\\Program Files\\Git\\bin\\bash.exe',
    'C:\\Program Files (x86)\\Git\\bin\\bash.exe',
  ];
  for (const p of candidates) {
    if (existsSync(p)) return p;
  }
  return 'bash'; // fall back to PATH (e.g. MSYS2, Cygwin)
}

const result = spawnSync(findBash(), [scriptPath], {
  stdio: 'inherit',
  env: process.env,
});

process.exit(result.status ?? 0);
