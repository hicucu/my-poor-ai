#!/usr/bin/env bash
# 공용 베이스 + 사용자 목록 페이지 (필터 없음)
set -euo pipefail
bash "$(dirname "$0")/../_fixtures/node-webapp.sh"
cat > src/routes/users.js <<'T'
import { listUsers } from '../models/user.js';
// GET /users — 전체 사용자 목록
export function handleUsers(req, res) {
  res.writeHead(200).end(JSON.stringify(listUsers().map(({ id, email }) => ({ id, email }))));
}
T
cat >> src/models/user.js <<'T'

export function listUsers() {
  return [...users.values()];
}
T
cat > src/views/users.html <<'T'
<h1>사용자 목록</h1>
<table id="users"><thead><tr><th>ID</th><th>이메일</th></tr></thead><tbody></tbody></table>
<script>
  fetch('/users').then((r) => r.json()).then((users) => {
    document.querySelector('#users tbody').innerHTML =
      users.map((u) => `<tr><td>${u.id}</td><td>${u.email}</td></tr>`).join('');
  });
</script>
T
git add -A
git -c user.name=eval -c user.email=eval@example.com commit -qm "add user list page"
