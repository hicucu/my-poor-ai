---
name: preventing-github-actions-loops
description: "Use when writing, editing, or reviewing GitHub Actions workflows (.github/workflows/*.yml) to prevent or detect infinite loops / self-triggering cycles — cases where a workflow's outputs (commits/pushes, pull requests, issue comments, tags/releases, dispatch) re-trigger another workflow (action A triggers B, and B triggers A again). Triggers: writing CI/CD workflows that push commits or open PRs, \"create/write a workflow\", \"inspect workflows\", \"prevent infinite loops\", \"check for circular triggers\", \"audit workflows\", \"recursive workflow\", \"workflow_run loop\"."
---

# GitHub Actions 무한루프 방지·점검

## 개요

워크플로우가 만들어낸 결과(push·PR·comment·tag·dispatch)가 다시 어떤 워크플로우의 트리거가 되면 순환이 발생함. 핵심 판단 한 줄: **"이 잡이 생산하는 이벤트가 어떤 워크플로우의 트리거 목록에 들어가고, 그 이벤트가 GitHub 기본 보호를 우회하는가?"**

두 가지 모드로 사용함:

- **Mode A — 작성/수정 시점**: 워크플로우를 새로 만들거나 고칠 때, 루프 유발 패턴을 피하고 가드를 미리 삽입.
- **Mode B — 점검 시점**: 기존 `.github/workflows/*.yml`를 런타임 순환 관점에서 감사.

## 핵심: GitHub 기본 보호와 우회 조건

`GITHUB_TOKEN`으로 만든 push·PR·comment·tag 이벤트는 **새 워크플로우 실행을 트리거하지 않음** (GitHub 기본 무한루프 방지 장치). 따라서 루프는 **이 보호가 깨지는 경로**에서만 발생함:

| 우회 경로                                                     | 루프 발생                                           |
| ------------------------------------------------------------- | --------------------------------------------------- |
| PAT(개인 토큰)·Deploy key·GitHub App 토큰으로 push/PR/comment | 발생 — 보호 우회됨                                  |
| `workflow_run` 트리거                                         | 발생 가능 — GITHUB_TOKEN 워크플로우 완료로도 깨어남 |
| `repository_dispatch` / `gh workflow run` / dispatch API 호출 | 발생 가능 — 잡이 직접 다른 워크플로우 기동          |
| `GITHUB_TOKEN`만 사용 + push/PR/comment                       | 미발생 — 기본 보호가 차단                           |

## Mode A — 작성/수정 시점 체크리스트

워크플로우 YAML을 쓰기 전·후 점검:

1. 이 워크플로우가 **무언가를 생산하는가?** (자동 커밋/push, PR 생성, comment, tag/release, dispatch)
2. 생산하는 이벤트가 **이 레포의 어떤 트리거 목록에 들어가는가?** (자기 자신 포함)
3. 들어간다면, 생산 행위에 **PAT/Deploy key/App 토큰**을 쓰는가? 또는 트리거가 `workflow_run`·`repository_dispatch`인가? → 그렇다면 가드 필수.
4. 아래 차단 가드 중 최소 하나를 삽입해 순환 경로를 끊음.
5. 정말 재트리거가 필요 없다면 **PAT 대신 `GITHUB_TOKEN`** 으로 다운그레이드 (가장 단순한 차단).

## Mode B — 점검 절차

1. **트리거 인벤토리** — 모든 워크플로우의 `on:` 이벤트 수집 (push, pull_request, issue_comment, issues, release, create, workflow_run, repository_dispatch, schedule 등).
2. **생산 행위 인벤토리** — 각 잡이 만드는 이벤트:
   - `git push` / `git commit` 후 push
   - `stefanzweifel/git-auto-commit-action`, `EndBug/add-and-commit`, `ad-m/github-push-action`
   - `peter-evans/create-pull-request`, `gh pr create`
   - `gh pr comment` / `gh issue comment` / `actions/github-script`로 comment
   - `git tag` push / `gh release create` / `softprops/action-gh-release`
   - `gh workflow run`, `repository_dispatch`/`workflow_dispatch` API 호출
3. **토큰 확인** — 각 생산 행위의 토큰을 봄. `secrets.GITHUB_TOKEN`인가, `secrets.PAT_*`·`secrets.*_TOKEN`·deploy key·App 토큰인가? `actions/checkout`의 `token:`·`persist-credentials`도 확인.
4. **순환 매칭** — 생산 이벤트(2)가 트리거 목록(1)에 들고 보호 우회 조건(위 표)에 해당하면 **루프 후보**로 표시. 자기참조와 A→B→A 상호참조 모두 검사.
5. **가드 존재 확인** — 후보마다 아래 가드가 경로를 끊는지 봄.

## 차단 가드 (하나라도 순환을 끊으면 안전)

| 가드              | 형태                                                                                                                 |
| ----------------- | -------------------------------------------------------------------------------------------------------------------- |
| 봇 actor 제외     | `if: github.actor != 'github-actions[bot]'`                                                                          |
| 스킵 마커         | 자동 커밋 메시지에 `[skip ci]`/`[no ci]`, 또는 `if: ${{ !contains(github.event.head_commit.message, '[skip ci]') }}` |
| 경로/브랜치 필터  | 자동 생성물 경로를 `paths-ignore`로 제외                                                                             |
| workflow_run 조건 | `if: github.event.workflow_run.conclusion == 'success'` + 상호 트리거 구조 회피                                      |
| concurrency 취소  | `concurrency: { group: ..., cancel-in-progress: true }` (루프 자체는 못 막지만 폭주 완화)                            |
| 토큰 다운그레이드 | 재트리거 불필요 시 PAT 대신 `GITHUB_TOKEN` 사용                                                                      |

## 점검 보고 형식

루프 후보별로:

```
[위험] <workflow-a.yml> (on: push)
  생산: stefanzweifel/git-auto-commit-action, 토큰=secrets.PAT
  순환: push → workflow-a.yml 재트리거 (PAT가 기본 보호 우회)
  가드: 없음
  권장: actor 가드 또는 [skip ci] 마커 추가, 또는 GITHUB_TOKEN으로 다운그레이드
```

위험도: 가드 없는 PAT/dispatch 순환 = 높음 · 가드 있으나 취약 = 중간 · GITHUB_TOKEN 전용 = 정보성.

## 흔한 실수

- **GITHUB_TOKEN인데 루프로 오판** — 기본 보호로 push/PR/comment는 재트리거 안 됨. `workflow_run`은 예외(GITHUB_TOKEN이어도 깨어남).
- **토큰 출처를 안 봄** — `checkout`의 `token:`과 push 액션의 `github_token:`/`token:`을 모두 확인해야 PAT 사용을 잡음.
- **자기참조만 보고 A→B→A 상호참조를 놓침** — 워크플로우 쌍을 교차 검사.
- **concurrency만 믿음** — 중복 취소는 폭주를 줄일 뿐 순환 자체를 끊지 못함.
- **schedule/외부 dispatch 무시** — cron이나 외부 dispatch가 사이클의 시작점일 수 있음.
