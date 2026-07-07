# Pipeline State Tracking Test

## 테스트 목적

`pipeline-state.md`가 각 Phase 완료 시 정확하게 업데이트되고, 부분 재실행 시 올바른 Phase부터 재시작하는지 확인.

## 테스트 시나리오

### 시나리오 1: Phase 순서 추적

전체 파이프라인을 실행하면서 각 Phase 완료 후 `pipeline-state.md` 상태 확인:

```markdown
Phase 0 완료 후:

- [x] Phase 0: project-context
- [ ] Phase 1: brainstorming
- [ ] Phase 2: planning
- [ ] Phase 3: development
- [ ] Phase 4: review
- [ ] Phase 5: finishing

Phase 1 완료 후:

- [x] Phase 0: project-context
- [x] Phase 1: brainstorming (design.md 승인 완료)
- [ ] Phase 2: planning
      ...
```

검증: 각 Phase 완료 직후 해당 항목이 `[x]`로 표시됨.

### 시나리오 2: 부분 재실행

Phase 3 실행 중 중단 상황 시뮬레이션:

```
pipeline-state.md 상태:
- [x] Phase 0: project-context
- [x] Phase 1: brainstorming
- [x] Phase 2: planning
- [ ] Phase 3: development   ← 여기서 중단
- [ ] Phase 4: review
- [ ] Phase 5: finishing
```

이후 "이어서 해줘" 요청 시:

- [ ] pipeline-state.md 확인 후 Phase 3부터 재시작
- [ ] Phase 0-2 재실행 없음
- [ ] project-context.md 24시간 이내 → 캐시 재사용

### 시나리오 3: 새 요청으로 재시작

"새 기능으로 다시 시작해줘" 요청 시:

- [ ] 기존 `_workspaces/feature-tags/` → `_workspaces/feature-tags_prev/`로 백업
- [ ] 새 `_workspaces/` 생성 후 Phase 0부터 시작
- [ ] 기존 pipeline-state.md 덮어쓰지 않음 (백업에 보존)

## 실패 판정

- Phase 완료 후 `pipeline-state.md` 미업데이트 → 오케스트레이터 로직 오류
- "이어서" 요청 시 Phase 0부터 재실행 → 부분 재실행 로직 오류
- 백업 없이 기존 작업 덮어씀 → 데이터 손실 위험
