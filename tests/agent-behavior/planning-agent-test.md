# Planning Agent Test

## 테스트 목적

`planning-agent`가 design.md를 읽고 플레이스홀더 없는 TDD 기반 스펙 파일을 생성하며 `file-manifest.json`을 포함하는지 확인.

## 사전 조건

`brainstorming-agent-test.md` 통과 후 생성된 `_workspaces/feature-tags/design.md` 사용.

## 호출 방법

```
Agent(subagent_type="planning-agent"):
  branch-slug: feature-tags
  design-path: _workspaces/feature-tags/design.md
  프로젝트 경로: {테스트용 프로젝트 경로}
```

## 검증 항목

- [ ] `_workspaces/feature-tags/specs/` 디렉토리 생성됨
- [ ] 스펙 파일 1개 이상 생성됨 (`spec-a.md`, ...)
- [ ] 각 스펙 파일에 실제 코드가 포함된 TDD 태스크 목록 존재
- [ ] "TBD", "TODO" 등 플레이스홀더 없음
- [ ] `file-manifest.json` 생성됨
- [ ] `file-manifest.json`에 `developmentOrder` 포함
- [ ] `depends-on` 필드가 의존 순서를 올바르게 반영
- [ ] STATUS: DONE 반환

## 플레이스홀더 스캔

```bash
grep -r "TBD\|TODO\|나중에\|implement later" _workspaces/feature-tags/specs/
# 결과 없어야 함
```
