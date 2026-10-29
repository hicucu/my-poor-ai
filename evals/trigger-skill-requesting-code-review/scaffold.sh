#!/usr/bin/env bash
# 공용 베이스 + 인증 기능 커밋 2개. 프롬프트의 abc123·def456을 태그로 만들어 범위가 해석되게 함
set -euo pipefail
bash "$(dirname "$0")/../_fixtures/node-webapp.sh"
g() { git -c user.name=eval -c user.email=eval@example.com "$@"; }
g tag abc123
mkdir -p src/auth
cat > src/auth/token.js <<'T'
import crypto from 'node:crypto';
const SECRET = 'dev-secret';
export function sign(userId) {
  const payload = `${userId}.${Date.now() + 3600_000}`;
  return `${payload}.${crypto.createHmac('sha256', SECRET).update(payload).digest('hex')}`;
}
export function verify(token) {
  const [userId, exp, sig] = token.split('.');
  const expected = crypto.createHmac('sha256', SECRET).update(`${userId}.${exp}`).digest('hex');
  return sig === expected && Number(exp) > Date.now() ? Number(userId) : null;
}
T
g add -A; g commit -qm "feat(auth): signed session tokens"
sed -i.bak 's/token: `t-${user.id}`/token: sign(user.id)/' src/routes/auth.js && rm src/routes/auth.js.bak
sed -i.bak "1i\\
import { sign } from '../auth/token.js';
" src/routes/auth.js && rm src/routes/auth.js.bak
g add -A; g commit -qm "feat(auth): issue signed token on login"
g tag def456
