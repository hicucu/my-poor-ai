---
description: my-poor-ai SessionStart 훅을 ~/.claude/settings.json에 자동 등록
allowed-tools: [Read, Edit, Bash]
---

# my-poor-ai:setup — SessionStart 훅 자동 등록

my-poor-ai 플러그인의 SessionStart 훅을 `~/.claude/settings.json`에 자동 등록함.

## 실행 절차

**1단계 — my-poor-ai 설치 경로 탐지**

`Bash`로 아래를 실행해 설치 버전 확인:

```
ls ~/.claude/plugins/cache/hicucu/my-poor-ai/
```

**2단계 — 훅 중복 확인**

`~/.claude/settings.json`을 읽어 `hooks.SessionStart` 배열에 `my-poor-ai` 관련 항목이 이미 있으면 중단하고 "이미 등록됨" 메시지를 출력함.

**3단계 — SessionStart 훅 추가**

중복이 없다면 `hooks.SessionStart` 배열에 아래 항목을 추가함 (섹션이 없으면 생성):

```json
{
  "hooks": [
    {
      "type": "command",
      "command": "node",
      "args": ["<PLUGIN_ROOT>/hooks/run-hook.mjs", "session-start"],
      "timeout": 10
    }
  ]
}
```

`<PLUGIN_ROOT>`는 `C:\Users\<USERNAME>\.claude\plugins\cache\hicucu\my-poor-ai\<VERSION>` 형태로 실제 경로로 치환함.

**4단계 — JSON 유효성 검증**

```bash
python3 -c "import json; json.load(open('$HOME/.claude/settings.json'))" && echo "JSON valid"
```

**5단계 — 완료 메시지**

"my-poor-ai SessionStart 훅 등록 완료. `/reload-plugins` 실행 후 적용됨."
