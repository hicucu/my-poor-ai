#!/usr/bin/env bash
# trigger 케이스 공용 픽스처 — 샌드박스 cwd에 소규모 Node 웹앱을 만들고 커밋 1개를 남김.
# 빈 샌드박스에서는 모델이 대상 파일을 찾다가 스킬 호출 전에 멈추므로, 프롬프트가 가리키는 대상을 실제로 둠.
# 사용처: 각 trigger 케이스의 case.yaml `context.scaffold_script` (실행 시 --scaffold 필요)
set -euo pipefail

mkdir -p src/routes src/models src/views docs/my-poor-ai/plans test

cat > package.json <<'EOF'
{
  "name": "board-app",
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "start": "node src/server.js",
    "test": "node --test"
  }
}
EOF

cat > src/config.js <<'EOF'
export const config = {
  port: 3000,
  timeout: 3000,
  db: { url: process.env.DATABASE_URL ?? 'postgres://localhost/board', poolSize: 10 },
};
EOF

cat > src/server.js <<'EOF'
import http from 'node:http';
import { config } from './config.js';
import { handleAuth } from './routes/auth.js';
import { handlePosts } from './routes/posts.js';

const server = http.createServer(async (req, res) => {
  if (req.url.startsWith('/auth')) return handleAuth(req, res);
  if (req.url.startsWith('/posts')) return handlePosts(req, res);
  res.writeHead(404).end();
});
server.setTimeout(config.timeout);
server.listen(config.port);
EOF

cat > src/models/user.js <<'EOF'
const users = new Map();

export function createUser({ email, password }) {
  if (users.has(email)) throw new Error('already registered');
  const user = { id: users.size + 1, email, password };
  users.set(email, user);
  return user;
}

export function findUser(email) {
  return users.get(email) ?? null;
}
EOF

cat > src/models/post.js <<'EOF'
const posts = [];

export function createPost({ authorId, title, body }) {
  const post = { id: posts.length + 1, authorId, title, body, createdAt: new Date() };
  posts.push(post);
  return post;
}

export function listPosts() {
  return posts;
}
EOF

cat > src/routes/auth.js <<'EOF'
import { createUser, findUser } from '../models/user.js';

// POST /auth/signup — 이메일·비밀번호로 즉시 가입
export async function handleAuth(req, res) {
  const body = await readJson(req);
  if (req.url === '/auth/signup') {
    const user = createUser(body);
    return res.writeHead(201).end(JSON.stringify({ id: user.id }));
  }
  if (req.url === '/auth/login') {
    const user = findUser(body.email);
    if (!user || user.password !== body.password) return res.writeHead(401).end();
    return res.writeHead(200).end(JSON.stringify({ token: `t-${user.id}` }));
  }
  res.writeHead(404).end();
}

async function readJson(req) {
  let raw = '';
  for await (const chunk of req) raw += chunk;
  return raw ? JSON.parse(raw) : {};
}
EOF

cat > src/routes/posts.js <<'EOF'
import { createPost, listPosts } from '../models/post.js';

export async function handlePosts(req, res) {
  if (req.method === 'GET') return res.writeHead(200).end(JSON.stringify(listPosts()));
  if (req.method === 'POST') {
    let raw = '';
    for await (const chunk of req) raw += chunk;
    const post = createPost(JSON.parse(raw));
    return res.writeHead(201).end(JSON.stringify(post));
  }
  res.writeHead(405).end();
}
EOF

cat > src/views/signup.html <<'EOF'
<form id="signup" method="post" action="/auth/signup">
  <input name="email" type="email" required>
  <input name="password" type="password" required>
  <button type="submit">제출</button>
</form>
EOF

cat > test/user.test.js <<'EOF'
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { createUser, findUser } from '../src/models/user.js';

test('createUser stores a user', () => {
  createUser({ email: 'a@example.com', password: 'pw' });
  assert.equal(findUser('a@example.com').email, 'a@example.com');
});
EOF

cat > docs/my-poor-ai/plans/2024-01-15-auth-system.md <<'EOF'
# Auth System Implementation Plan

**Goal:** Replace the plain-text password check with hashed passwords and session tokens.

## Task 1: Hash passwords on signup
- Modify `src/models/user.js`: store `crypto.scryptSync` hash + salt instead of the raw password
- Test: `test/user.test.js` — stored user has no `password` field, has `hash` and `salt`

## Task 2: Verify hashed passwords on login
- Modify `src/routes/auth.js`: compare with `crypto.timingSafeEqual`
- Test: wrong password returns 401, correct password returns 200

## Task 3: Signed session tokens
- Create `src/auth/token.js`: HMAC-signed token with expiry
- Modify `src/routes/auth.js`: return the signed token on login
- Test: tampered token fails verification
EOF

git add -A
git -c user.name=eval -c user.email=eval@example.com commit -qm "init board-app"
