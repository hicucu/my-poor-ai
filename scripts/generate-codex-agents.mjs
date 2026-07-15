#!/usr/bin/env node
/**
 * .codex/agents/*.toml 자동 생성기.
 *
 * agents/*.md(단일 소스)의 frontmatter와 본문을 Codex 에이전트 정의(toml)로 변환함.
 * .codex/agents/ 하위 toml은 전부 생성물 — 수동 편집 금지, 수정은 agents/*.md에.
 *
 * 변환 규칙:
 *   name                  = frontmatter name + "-agent" (이미 -agent로 끝나면 그대로)
 *   description           = frontmatter description (TOML basic string으로 이스케이프)
 *   developer_instructions = frontmatter를 제외한 .md 본문 verbatim ('''리터럴 블록)
 *   model/tools           = Codex 스키마 밖 정보 — 상단 주석으로만 보존
 *
 * 사용:
 *   node scripts/generate-codex-agents.mjs           # 전량 재생성
 *   node scripts/generate-codex-agents.mjs --check   # 디스크와 비교만 (CI용)
 * 종료 코드: 0 = 성공/일치, 1 = 오류/불일치
 */
import { readFileSync, writeFileSync, readdirSync, existsSync } from 'fs';
import { join } from 'path';

const ROOT = new URL('..', import.meta.url).pathname;
const AGENT_DIR = join(ROOT, 'agents');
const CODEX_DIR = join(ROOT, '.codex', 'agents');
const CHECK = process.argv.includes('--check');

const errors = [];

/** YAML 따옴표 값 → 실제 문자열 (JSON.parse 호환, 실패 시 따옴표만 제거) */
function unquote(v) {
  if (v.startsWith('"') && v.endsWith('"')) {
    try { return JSON.parse(v); } catch { return v.slice(1, -1); }
  }
  return v;
}

/** TOML basic string 이스케이프 */
const tomlEscape = (s) => s.replace(/\\/g, '\\\\').replace(/"/g, '\\"');

// ---------- agents/*.md → toml 내용 생성 ----------
const generated = new Map(); // toml 파일명 → 내용
for (const f of readdirSync(AGENT_DIR).filter((f) => f.endsWith('.md') && !f.startsWith('README')).sort()) {
  const src = `agents/${f}`;
  const text = readFileSync(join(AGENT_DIR, f), 'utf8');
  const m = text.match(/^---\n([\s\S]*?)\n---\n/);
  if (!m) { errors.push(`${src}: frontmatter 없음`); continue; }

  const fields = Object.fromEntries(m[1].split('\n').map((l) => {
    const i = l.indexOf(':');
    return [l.slice(0, i).trim(), l.slice(i + 1).trim()];
  }));
  const name = unquote(fields.name ?? '');
  const description = unquote(fields.description ?? '');
  if (!name || !description) { errors.push(`${src}: name/description 누락`); continue; }

  let body = text.slice(m[0].length).replace(/^\n+/, '');
  if (body.includes("'''")) {
    errors.push(`${src}: 본문에 ''' 포함 — TOML 리터럴 블록과 충돌, 본문 수정 필요`);
    continue;
  }
  if (!body.endsWith('\n')) body += '\n';

  const tomlName = name.endsWith('-agent') ? name : `${name}-agent`;
  const content = [
    `# AUTO-GENERATED from ${src} — 수동 편집 금지.`,
    `# 수정은 ${src}에 반영 후 \`node scripts/generate-codex-agents.mjs\` 재실행.`,
    `# model: ${fields.model ?? '-'} / tools: ${fields.tools ?? '-'}`,
    `name = "${tomlName}"`,
    `description = "${tomlEscape(description)}"`,
    `developer_instructions = '''`,
    body + `'''`,
    ``,
  ].join('\n');
  generated.set(`${tomlName}.toml`, content);
}

// ---------- 고아 toml 탐지 (대응 .md 없는 생성물) ----------
if (existsSync(CODEX_DIR)) {
  for (const f of readdirSync(CODEX_DIR).filter((f) => f.endsWith('.toml'))) {
    if (!generated.has(f)) errors.push(`.codex/agents/${f}: 대응하는 agents/*.md 없음 (고아 toml — 삭제 필요)`);
  }
}

// ---------- 쓰기 또는 비교 ----------
const stale = [];
for (const [f, content] of generated) {
  const path = join(CODEX_DIR, f);
  const current = existsSync(path) ? readFileSync(path, 'utf8') : null;
  if (current === content) continue;
  if (CHECK) stale.push(`.codex/agents/${f}${current === null ? ' (누락)' : ''}`);
  else writeFileSync(path, content);
}

// ---------- 결과 ----------
if (errors.length) {
  console.error(`✗ ${errors.length}건 오류:\n`);
  for (const e of errors) console.error(`  - ${e}`);
  process.exit(1);
}
if (CHECK && stale.length) {
  console.error(`✗ ${stale.length}개 toml이 agents/*.md와 불일치 — \`node scripts/generate-codex-agents.mjs\` 재실행 필요:\n`);
  for (const s of stale) console.error(`  - ${s}`);
  process.exit(1);
}
console.log(CHECK
  ? `✓ 일치 — .codex/agents/ toml ${generated.size}개가 agents/*.md와 동기화됨`
  : `✓ 생성 완료 — .codex/agents/ toml ${generated.size}개`);
