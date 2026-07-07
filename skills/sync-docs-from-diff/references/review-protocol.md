# Phase 3: 사용자 검토 프로토콜

오케스트레이터(메인 Claude)가 직접 수행하는 검토·적용 절차. 별도 에이전트를 호출하지 않음.

## 목차

- [전제 조건](#전제-조건)
- [Step 1: 영역별 인덱스 요약](#step-1-영역별-인덱스-요약)
- [Step 2: 검토 모드 선택](#step-2-검토-모드-선택)
- [Step 3: 파일별 검토 루프](#step-3-파일별-검토-루프)
- [Step 4: 적용 (Edit 도구)](#step-4-적용-edit-도구)
- [Step 5: apply_log 기록](#step-5-apply_log-기록)
- [엣지 케이스](#엣지-케이스)

## 전제 조건

- Phase 2 완료 — `_workspaces/proposals/{inline,readme,docs}/_index.md` 존재
- 각 영역에 0개 이상의 `*.patch.md` 파일
- 사용자가 검토 후 적용 모드를 선택한 상태 (오케스트레이터의 기본 동작)

## Step 1: 영역별 인덱스 요약

3개 인덱스 파일을 모두 읽고 사용자에게 표 형태로 요약 제시:

```
검토할 패치가 준비되었습니다.

| 영역    | 패치 수 | 주요 변경                                         |
|---------|--------|--------------------------------------------------|
| README  | 1      | Features 1줄 추가, CLI 예시 1건 수정              |
| docs    | 7      | API 시그니처 3건, 가이드 코드 2건, 아키텍처 2건   |
| inline  | 4      | src/api/ 하위 component README 4건                |
| 합계    | 12     | -                                                |
```

총 패치 수가 0이면 "검토할 변경 없음"으로 보고하고 Phase 4로 직행.

## Step 2: 검토 모드 선택

`AskUserQuestion`으로 검토 진행 방식 묻기:

```
질문: 어떻게 검토할까요?
- 1. 파일별로 하나씩 검토 (Recommended) — 각 패치를 보고 승인/거부 결정
- 2. 영역별 일괄 (README 일괄 → docs 일괄 → inline 일괄) — 영역 단위 한 번 보고 일괄 승인/거부
- 3. 전부 일괄 승인 — 모든 패치를 묻지 않고 적용. 위험.
```

## Step 3: 파일별 검토 루프

선택된 모드에 따라 패치 파일들을 순회. 권장 순서: **inline → docs → README**.
- 이유: README는 보통 docs를 참조하므로 docs를 먼저 정리한 뒤 README의 표현을 맞추는 것이 안전.

각 패치마다:

1. 패치 파일을 Read로 읽음
2. 사용자에게 다음 형식으로 제시:

```
[1/12] 파일: README.md
변경 사유: searchUsers() API 신규 추가에 따른 Features 섹션 갱신

Before (line 34):
- User retrieval by ID

After:
- User retrieval by ID
- User search by query string (NEW)

승인하시겠습니까?
```

3. `AskUserQuestion`으로 4지선다:
   - **승인 후 적용** — 즉시 Edit 적용
   - **수정 제안** — 사용자가 "X 부분을 Y로 바꿔서" 같은 추가 지시. 메인 Claude가 재작성 후 다시 묻기
   - **거부 (이번 영역에서)** — 적용하지 않고 다음 패치로
   - **보류 (전체 중단)** — 검토를 멈추고 현재까지의 진행 상황을 사용자에게 보고

## Step 4: 적용 (Edit 도구)

승인 시:

1. 원본 파일을 Read로 다시 확인 (Phase 2 시점과 달라졌을 수 있음 — 외부 변경 감지)
2. 패치의 Before 텍스트가 현재 파일에 그대로 있는지 grep/Read로 검증
3. 일치하면 `Edit` 도구로 Before → After 치환
4. 일치하지 않으면(파일이 변경됨) — 적용 보류, apply_log에 "skipped: file diverged" 기록, 사용자에게 한 줄 보고
5. 한 패치에 여러 Before/After 블록이 있으면 위에서부터 순차 적용. 한 블록 실패 시 그 블록만 skip 처리, 나머지 블록은 계속.

> Edit 도구의 old_string 충돌(여러 군데 매칭)이 발생하면 더 넓은 컨텍스트를 포함해 재시도. 그래도 실패하면 사용자에게 "수동 수정 필요"로 보고하고 skip.

## Step 5: apply_log 기록

매 패치 처리 후 `_workspaces/03_apply_log.md`에 한 줄 추가:

```markdown
# Apply Log
- 시작: 2026-05-07 17:30
- 베이스: develop@a1b2c3d → HEAD@e4f5g6h

## 항목
- [APPLIED] README.md  (Features +1, CLI 1건)         2026-05-07 17:32
- [REJECTED] docs/api/users.md (사용자 거부)           2026-05-07 17:33
- [SKIPPED] docs/architecture/modules.md (file diverged) 2026-05-07 17:34
- [APPLIED] src/api/README.md (1건)                     2026-05-07 17:35
...

## 요약
- APPLIED: 9 / REJECTED: 2 / SKIPPED: 1 / HOLD: 0
```

이 로그는 doc-sync-validator의 입력으로 사용되며, 후속 부분 재실행 시 "이미 적용된 것"의 근거가 됨.

## 엣지 케이스

### 패치의 Before가 파일에 부분 일치
- 줄바꿈/공백만 다른 경우: 정규화 후 재시도. 그래도 실패하면 skip.
- 내용 자체가 다른 경우: skip + apply_log에 "diverged" 기록.

### 패치 적용 도중 사용자가 "전체 중단" 선택
- 즉시 루프 종료. apply_log에 "HOLD - user interrupted"로 마무리.
- Phase 4(검증)로 넘어갈지 사용자에게 묻기. 보통 부분 적용 상태이므로 검증해도 의미 있음.

### 동일 파일에 여러 영역이 패치 제안
- 영역 분리 원칙 위반 — 정상이라면 발생하지 않아야 함. 발생 시:
  1. 사용자에게 충돌 사실 보고
  2. 어느 영역의 패치를 우선 적용할지 물음 (보통 readme/docs/inline 순)
  3. 충돌 영역의 결과를 apply_log에 명시

### 패치 파일이 손상되어 Before/After를 파싱할 수 없음
- skip + apply_log에 "patch_corrupt" 기록. validator가 보고서에서 강조.

### 사용자가 검토 도중 "남은 거 전부 자동 승인"으로 변경 요청
- 모드 전환을 허용. 이후 패치는 묻지 않고 적용. 단, apply_log에 "[bulk-approved from this point]" 마커 기록.
