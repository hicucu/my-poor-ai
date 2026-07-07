---
description: 역할 프리셋 카탈로그 — 역할명으로 시작하면 해당 스킬 번들의 진입 스킬로 라우팅. /my-poor-ai:roles {역할} 형태로 실행.
---

# my-poor-ai:roles — 역할 프리셋

어떤 스킬로 시작할지 모를 때 역할 이름으로 진입하는 얇은 라우팅 레이어. 새 기능이 아니라 기존 스킬 번들에 대한 명명된 진입점임.

## 사용법

```
/my-poor-ai:roles              # 프리셋 목록 표시
/my-poor-ai:roles architect    # 해당 역할의 진입 스킬 호출
```

## 프리셋 목록

| 역할          | 스킬 번들 (순서대로)                                                             | 진입 스킬                    |
| ------------- | -------------------------------------------------------------------------------- | ---------------------------- |
| **architect** | brainstorming → writing-plans → socratic-plan-review                              | `my-poor-ai:brainstorming`        |
| **builder**   | test-driven-development → subagent-driven-development → finishing-a-development-branch | `my-poor-ai:test-driven-development` |
| **debugger**  | systematic-debugging → verification-before-completion                             | `my-poor-ai:systematic-debugging` |
| **reviewer**  | requesting-code-review / receiving-code-review / `/my-poor-ai:code-review`             | `my-poor-ai:requesting-code-review` |
| **docs**      | sync-docs-from-diff / generate-claude-instructions                                | `my-poor-ai:sync-docs-from-diff`  |

## 실행 절차

1. 인자 없이 호출되면 위 프리셋 표를 그대로 표시하고 종료.
2. 역할명이 주어지면 해당 행의 **진입 스킬**을 Skill 도구로 호출. 이후 흐름은 스킬 자체의 지침을 따름 (번들의 후속 스킬은 진입 스킬이 자연스럽게 연결함).
3. 역할명이 표에 없으면 표를 표시하고 가장 가까운 역할을 제안.
4. reviewer 역할에서 "내 코드 리뷰 받기"는 `my-poor-ai:requesting-code-review`, "리뷰 피드백 처리"는 `my-poor-ai:receiving-code-review`, "브랜치 diff 4종 병렬 리뷰"는 `/my-poor-ai:code-review` 커맨드로 세분화하여 사용자 의도에 맞게 선택.

## 주의

- 프리셋은 using-my-poor-ai의 자동 분류(DEBUG/SIMPLE/FULL)를 대체하지 않음 — 사용자가 역할을 명시적으로 지정했을 때만 사용.
- 번들에 새 스킬을 추가할 때는 이 표와 README.md의 Role Presets 표를 함께 갱신.
