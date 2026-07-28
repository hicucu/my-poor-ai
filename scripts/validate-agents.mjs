#!/usr/bin/env node
/**
 * my-poor-ai 저장소 정합성 검증기.
 *
 * 검사 항목:
 *  1. agents/*.md frontmatter — name=파일명, model 유효값, tools 화이트리스트,
 *     tools 알파벳순, 필드 순서·미지원 필드
 *  2. 참조 해소 — {팀_위치}/agents/X.md 경로 참조와 subagent_type="..." 참조가
 *     실재 에이전트 정의로 해소되는지
 *  3. 코드펜스 균형 — 중첩 펜스 조기 닫힘/미닫힘 (포매터 손상 회귀 방지)
 *
 * 사용: node scripts/validate-agents.mjs   (저장소 루트에서 실행)
 * 종료 코드: 0 = 통과, 1 = 위반 발견
 */
import { readFileSync, readdirSync, existsSync, statSync } from 'fs';
import { join } from 'path';

const ROOT = new URL('..', import.meta.url).pathname;

// 모델 별칭 + 풀 모델 ID(claude-opus-5 등). 별칭만 허용하면 특정 버전 고정이 불가능하다.
const MODEL_ALIASES = new Set(['haiku', 'sonnet', 'opus', 'fable', 'inherit']);
const MODEL_ID = /^claude-[a-z0-9][a-z0-9.-]*$/;

// 백그라운드 서브에이전트가 실제로 보유하는 빌트인 도구 집합.
// 플러그인 에이전트는 기본이 백그라운드 실행이므로 이 집합이 유효 상한이다.
// `Agent`는 포그라운드/깊이 제한 조건부라 별도 허용 (review-agent 오케스트레이터용).
const VALID_TOOLS = new Set([
  'Agent', 'Artifact', 'Bash', 'Edit', 'EnterWorktree', 'ExitWorktree',
  'Glob', 'Grep', 'Monitor', 'NotebookEdit', 'PowerShell', 'Read',
  'SendMessage', 'Skill', 'TaskStop', 'TodoWrite', 'ToolSearch',
  'WebFetch', 'WebSearch', 'Write',
]);

// 서브에이전트에서 항상 제거되는 도구 — tools에 적어도 바인딩되지 않는다.
// 전량 제거되면 0-tool 스폰 실패로 이어지므로 이름 단계에서 잡는다.
const NEVER_IN_SUBAGENT = new Set([
  'AskUserQuestion', 'EndConversation', 'EnterPlanMode', 'ExitPlanMode',
  'ScheduleWakeup', 'TaskOutput', 'WaitForMcpServers', 'Workflow',
]);

// 지원 frontmatter 필드의 정본 순서. 존재하는 키가 이 순서의 부분수열이어야 한다.
// (고정 4개 목록이던 시절엔 isolation·effort 등 신규 필드 채택이 CI에서 막혔다.)
const FIELD_ORDER = [
  'name', 'description', 'model', 'tools', 'disallowedTools',
  'skills', 'effort', 'isolation', 'background', 'maxTurns', 'memory', 'color',
];
const REQUIRED_FIELDS = ['name', 'description', 'model', 'tools'];

// 플러그인으로 배포되는 에이전트에서는 무시되는 필드 — 조용히 안 먹느니 반려한다.
const PLUGIN_IGNORED_FIELDS = new Set(['hooks', 'mcpServers', 'permissionMode']);

const errors = [];
const err = (file, msg) => errors.push(`${file}: ${msg}`);

// ---------- 1. agents/*.md frontmatter ----------
const agentDir = join(ROOT, 'agents');
const agentNames = new Set();
for (const f of readdirSync(agentDir).filter((f) => f.endsWith('.md') && !f.startsWith('README')).sort()) {
  const path = `agents/${f}`;
  const text = readFileSync(join(agentDir, f), 'utf8');
  const m = text.match(/^---\n([\s\S]*?)\n---\n/);
  if (!m) { err(path, 'frontmatter 없음'); continue; }

  const lines = m[1].split('\n');
  const keys = lines.map((l) => l.split(':', 1)[0].trim());
  const fields = Object.fromEntries(lines.map((l) => {
    const i = l.indexOf(':');
    return [l.slice(0, i).trim(), l.slice(i + 1).trim()];
  }));

  const slug = f.replace(/\.md$/, '');
  const name = (fields.name ?? '').replace(/^"|"$/g, '');
  if (name !== slug) err(path, `name(${name})이 파일명 slug(${slug})와 다름`);

  for (const k of REQUIRED_FIELDS) {
    if (!fields[k]) err(path, `${k} 없음`);
  }

  const model = (fields.model ?? '').replace(/^"|"$/g, '');
  if (model && !MODEL_ALIASES.has(model) && !MODEL_ID.test(model)) {
    err(path, `유효하지 않은 model: ${model} (별칭 ${[...MODEL_ALIASES].join('/')} 또는 claude-* 풀 ID)`);
  }

  if (fields.tools) {
    const tools = fields.tools.split(',').map((t) => t.trim());
    for (const t of tools) {
      if (VALID_TOOLS.has(t)) continue;
      let hint = ' (레지스트리 실명만 허용)';
      if (t === 'Task') hint = ' (v2.1.63부터 Agent 사용)';
      else if (NEVER_IN_SUBAGENT.has(t)) hint = ' (서브에이전트에서 항상 제거되는 도구 — 바인딩되지 않음)';
      err(path, `유효하지 않은 도구명: ${t}${hint}`);
    }
    const sorted = [...tools].sort();
    if (tools.join() !== sorted.join()) {
      err(path, `tools가 알파벳순이 아님: ${tools.join(', ')} → ${sorted.join(', ')}`);
    }
  }

  for (const k of keys) {
    if (PLUGIN_IGNORED_FIELDS.has(k)) {
      err(path, `플러그인 에이전트에서 무시되는 필드: ${k} (settings.json 또는 .claude/agents/로 이전 필요)`);
    } else if (!FIELD_ORDER.includes(k)) {
      err(path, `지원되지 않는 frontmatter 필드: ${k}`);
    }
  }

  const known = keys.filter((k) => FIELD_ORDER.includes(k));
  const expected = FIELD_ORDER.filter((k) => known.includes(k));
  if (known.join() !== expected.join()) {
    err(path, `필드 순서 위반: ${known.join(', ')} (기대: ${expected.join(', ')})`);
  }
  agentNames.add(slug);
}

// ---------- 마크다운 파일 수집 ----------
function* mdFiles(dir) {
  for (const e of readdirSync(join(ROOT, dir))) {
    if (e === 'node_modules' || e.startsWith('.')) continue;
    const rel = `${dir}/${e}`;
    if (statSync(join(ROOT, rel)).isDirectory()) yield* mdFiles(rel);
    else if (e.endsWith('.md')) yield rel;
  }
}
const allMd = ['AGENTS.md', 'CLAUDE.md', 'README.md'];
for (const d of ['agents', 'commands', 'skills', 'tests']) allMd.push(...mdFiles(d));

// ---------- 2. 참조 해소 ----------
for (const rel of allMd) {
  const text = readFileSync(join(ROOT, rel), 'utf8');
  for (const m of text.matchAll(/(?:\{팀_위치\}\/)?agents\/([a-z0-9-]+\.md)/g)) {
    if (!existsSync(join(agentDir, m[1]))) err(rel, `존재하지 않는 에이전트 파일 참조: agents/${m[1]}`);
  }
  for (const m of text.matchAll(/subagent_type[=:]\s*"([^"]+)"/g)) {
    const name = m[1].replace(/^my-poor-ai:/, '');
    if (name !== 'general-purpose' && !agentNames.has(name)) {
      err(rel, `해소되지 않는 subagent_type: ${m[1]}`);
    }
  }
}

// ---------- 2b. @include 대상 해소 ----------
// 워커가 공유 지침 모듈을 include하는 `@include: <path>` 지시가 실재 파일로 해소되는지 검사.
// path는 agents/ 기준 상대 (예: _shared/implementation-conventions.md).
for (const rel of allMd) {
  const text = readFileSync(join(ROOT, rel), 'utf8');
  for (const m of text.matchAll(/^@include:\s*(\S+)\s*$/gm)) {
    if (!existsSync(join(agentDir, m[1]))) err(rel, `해소되지 않는 @include: agents/${m[1]}`);
  }
}

// ---------- 3. 코드펜스 균형 ----------
for (const rel of allMd) {
  const lines = readFileSync(join(ROOT, rel), 'utf8').split('\n');
  let open = 0; // 현재 열린 펜스의 백틱 수 (0 = 닫힘)
  for (const line of lines) {
    const m = line.match(/^\s*(`{3,})(.*)$/);
    if (!m) continue;
    if (open === 0) open = m[1].length;
    else if (m[1].length >= open && m[2].trim() === '') open = 0; // 유효한 닫힘
    // 더 짧거나 info string이 있는 펜스는 열린 블록의 리터럴 내용
  }
  if (open !== 0) err(rel, `EOF에서 닫히지 않은 코드펜스 (${open}-backtick)`);
}

// ---------- 4. 커맨드 카탈로그 등록 ----------
// 모든 commands/*.md는 카탈로그(commands.md) 또는 루트 문서에서 소개되어야 함 —
// 이번 세션류의 "존재하지만 어디서도 안내되지 않는 커맨드" 재발 방지.
{
  const catalog = readFileSync(join(ROOT, 'commands/commands.md'), 'utf8')
    + readFileSync(join(ROOT, 'CLAUDE.md'), 'utf8')
    + readFileSync(join(ROOT, 'README.md'), 'utf8');
  for (const f of readdirSync(join(ROOT, 'commands')).filter((f) => f.endsWith('.md')).sort()) {
    const stem = f.replace(/\.md$/, '');
    if (stem === 'commands' || f.startsWith('README')) continue;
    if (!catalog.includes(stem)) {
      err(`commands/${f}`, '카탈로그(commands.md)·CLAUDE.md·README.md 어디에도 소개되지 않는 커맨드');
    }
  }
}

// ---------- 5. hooks 매니페스트 → 스크립트 실재 ----------
// hooks.json/hooks-cursor.json이 참조하는 hooks/ 스크립트가 실제로 존재해야 함 —
// 고아 훅 래퍼·깨진 진입점 재발 방지.
for (const manifest of ['hooks/hooks.json', 'hooks/hooks-cursor.json']) {
  if (!existsSync(join(ROOT, manifest))) continue;
  const text = readFileSync(join(ROOT, manifest), 'utf8');
  for (const m of text.matchAll(/([A-Za-z0-9_-]+\.(?:mjs|sh|cmd))|(?:hooks\/)([A-Za-z0-9_-]+)(?=["' ])/g)) {
    const name = m[1] ?? m[2];
    if (name && !existsSync(join(ROOT, 'hooks', name))) {
      err(manifest, `존재하지 않는 훅 스크립트 참조: hooks/${name}`);
    }
  }
}

// ---------- 결과 ----------
if (errors.length) {
  console.error(`✗ ${errors.length}건 위반:\n`);
  for (const e of errors) console.error(`  - ${e}`);
  process.exit(1);
}
console.log(`✓ 통과 — 에이전트 ${agentNames.size}개, 마크다운 ${allMd.length}개 검사`);
