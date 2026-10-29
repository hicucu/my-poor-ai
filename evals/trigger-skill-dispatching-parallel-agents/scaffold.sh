#!/usr/bin/env bash
# 공용 베이스 + 프롬프트가 가리키는 실패 테스트 4개 (서로 다른 모듈)
set -euo pipefail
bash "$(dirname "$0")/../_fixtures/node-webapp.sh"
mkdir -p tests/auth tests/api tests/components tests/utils
cat > tests/auth/login.test.ts <<'T'
import { login } from '../../src/auth/login';
test('should redirect after login', async () => {
  const res = await login({ email: 'a@example.com', password: 'pw' });
  expect(res.redirect).toBe('/dashboard'); // 실제: undefined
});
T
cat > tests/api/users.test.ts <<'T'
import { app } from '../../src/app';
test('should return user list', async () => {
  const res = await app.inject({ method: 'GET', url: '/api/users' });
  expect(res.statusCode).toBe(200); // 실제: 500
});
T
cat > tests/components/Button.test.tsx <<'T'
import { render } from '@testing-library/react';
import { Button } from '../../src/components/Button';
test('matches snapshot', () => {
  expect(render(<Button label="저장" />).container).toMatchSnapshot(); // snapshot mismatch
});
T
cat > tests/utils/date.test.ts <<'T'
import { formatDate } from '../../src/utils/date';
test('timezone handling', () => {
  expect(formatDate('2026-01-01T00:00:00Z', 'Asia/Seoul')).toBe('2026-01-01 09:00'); // 실제: 2026-01-01 00:00
});
T
git add -A
git -c user.name=eval -c user.email=eval@example.com commit -qm "add failing tests"
