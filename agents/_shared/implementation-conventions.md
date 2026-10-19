# 공통 구현 규약 (공유 모듈)

> 이 파일은 실행 에이전트가 아니라 **구현 워커의 지침 모듈**임. `developer-agent`가 정의에서 이 파일을 `@include`함. 수정 후에는 `node scripts/generate-codex-agents.mjs` 재실행 필수 (Codex 미러에 인라인되는 내용이므로).

구현 워커는 아래 규율과 스택 컨벤션을 준수함. 워커 고유 계약(TDD·커밋·HANDOFF·출력 프로토콜 등)은 에이전트 정의에 따로 명시됨.

---

## 코딩 규율

1. **YAGNI** — 명세/스펙 범위 밖 기능·리팩터링 추가 금지. "나중에 필요할 것 같아서" 식의 확장 금지.
2. **단일 책임** — 지정된 대상(스펙)만 구현. 담당 범위 밖의 다른 파일 수정 금지 (오케스트레이터가 팬아웃으로 관리). 스펙이 지정한 테스트 파일은 담당 범위에 포함됨.
3. **기존 패턴 준수** — 인근 파일의 import 방식·명명·폴더 구조·의존 주입 패턴을 먼저 파악하고 그대로 따름. 새 컨벤션 도입 금지.
4. **의존 파일 먼저 읽기** — `dependencies`/의존 스펙에 명시된 파일을 모두 읽어 인터페이스를 파악한 후 구현.
5. **정적 타입 규율** — 정적 타입 언어는 명시적 타입 사용, `any` 등 회피 타입 금지. 불가피할 경우 명시적 주석으로 사유 기록 후 사용.
6. **기존 파일 최소 변경** — `modify` 작업은 기존 파일을 먼저 읽고 최소 변경 원칙 적용.
7. **상대 경로만** — 모든 경로는 CWD(프로젝트 루트) 기준 상대 경로. 절대 경로 / `~/` 사용 금지.

---

## 스택 컨벤션 매트릭스

스택 판정 기준:

- `_workspaces/stack-profile.json`이 **있으면** `profile.primary`/`profile.subtype`으로 판정.
- **없으면** 프로젝트 마커(매니페스트 파일·소스 레이아웃·인근 파일)로 판정.

판정한 스택에 따라 아래 표의 컨벤션을 적용함.

| primary | language    | 명명                  | import 스타일                              | 모듈/네임스페이스                 |
| ------- | ----------- | --------------------- | ------------------------------------------ | --------------------------------- |
| node    | typescript  | camelCase / PascalCase | esm 우선, 경로 alias 존재 시 활용          | export 명시                       |
| node    | javascript  | camelCase / PascalCase | esm 또는 cjs (profile 따름)                | module.exports / export           |
| python  | python      | snake_case            | 상대 import 신중 사용, 절대 import 우선    | `__init__.py` 갱신 필요 시 검토   |
| dotnet  | csharp      | PascalCase            | `using` 정렬, file-scoped namespace 우선   | namespace 일치 (폴더 구조 매칭)   |
| go      | go          | camelCase (내부) / PascalCase (export) | 표준 import 블록 (stdlib/3rd/local 그룹) | package 선언 일치                 |
| jvm     | java/kotlin | camelCase / PascalCase | package 선언 + import                      | 패키지 경로 = 디렉토리 경로       |
| rust    | rust        | snake_case            | `use` 선언 정리                            | mod 트리 갱신 필요 시 검토        |
| php     | php         | PascalCase (클래스)   | PSR-4 autoload 준수                        | namespace = 폴더 경로             |

미지원 스택(`primary: unknown` 또는 매트릭스에 없음)은 가장 유사한 행의 패턴을 차용하고 보고서에 명시.

---

## 테스트 컨벤션 (testFramework별)

테스트 파일 명명은 stack-profile의 `conventions.testPattern`, 위치는 `sourceLayout.tests`를 우선함. 둘 다 없으면 인근 기존 테스트의 위치·명명을 따름. 기존 테스트 파일이 있으면 먼저 읽고 패턴을 통일함.

| testFramework | 언어        | 테스트 파일 명명                      | 모킹 라이브러리                 | AAA 예시 패턴                              |
| ------------- | ----------- | ------------------------------------- | ------------------------------- | ------------------------------------------ |
| `vitest`      | ts/js       | `*.test.ts` 또는 `*.spec.ts`          | `vi.mock`, `vi.fn`              | `describe → it → expect`                   |
| `jest`        | ts/js       | `*.test.ts` 또는 `*.spec.ts`          | `jest.mock`, `jest.fn`          | `describe → it → expect`                   |
| `pytest`      | python      | `test_*.py` 또는 `*_test.py`          | `unittest.mock`, `pytest-mock`  | `def test_*: arrange/act/assert`           |
| `unittest`    | python      | `test_*.py`                           | `unittest.mock`                 | `class TestX(unittest.TestCase)`           |
| `xunit`       | csharp      | `*Tests.cs`                           | `Moq`, `NSubstitute`            | `[Fact] / [Theory]`                        |
| `nunit`       | csharp      | `*Tests.cs`                           | `Moq`                           | `[Test]`                                   |
| `mstest`      | csharp      | `*Tests.cs`                           | `Moq`                           | `[TestMethod]`                             |
| `go-test`     | go          | `*_test.go`                           | `testify/mock`, 인터페이스 더미 | `func TestX(t *testing.T)` + 테이블 드리븐 |
| `junit5`      | java/kotlin | `*Test.java` / `*Test.kt`             | `Mockito`, `MockK`              | `@Test`                                    |
| `cargo-test`  | rust        | 동일 파일 내 `#[cfg(test)] mod tests` | `mockall`                       | `#[test] fn name()`                        |
| `phpunit`     | php         | `*Test.php`                           | PHPUnit 내장 + `Prophecy`       | `public function testX(): void`            |
| `pest`        | php         | `*Test.php`                           | PHPUnit 기반                    | `it('does X', function () { ... })`        |

미지원 framework는 가장 유사한 행을 차용하고 보고서에 명시. 모킹은 외부 시스템(DB·HTTP·시간)에만 쓰고 도메인 로직은 실제 구현으로 테스트함.

---

## 파일 유형별 가이드 (subtype 분기)

파일 유형(스펙이 밝힌 파일 역할 또는 소스 레이아웃 위치) + 스택 subtype 조합으로 결정.

### UI 계층 (`profile.uiMarkers` 경로 또는 `type: component/view/page`)

- 입력/출력 인터페이스 명시 (Props·ViewModel·DTO)
- 이벤트 핸들러 명명 일관성 (`handle{Event}` 또는 `on{Event}` — 프로젝트 패턴 따름)
- 사이드 이펙트는 별도 계층(훅·서비스·use case)으로 위임
- 스타일링: 프로젝트 기존 방식 확인 후 적용

### 비즈니스 로직 (`profile.businessLogicMarkers` 경로 또는 `type: service/use-case/domain`)

- 입출력 타입 명시
- 함수형 우선 (필요할 때만 클래스)
- 외부 시스템 호출 시 에러 throw, 호출 측에서 처리
- 부수효과 최소화, 의존 주입 가능한 구조
- 이 계층 파일은 이후 단위테스트 대상이 되므로, 테스트가 가능하도록 순수 함수/명시적 인터페이스로 구현함

### 컨트롤러/라우터/핸들러 (`type: controller/router/handler`)

- 요청 파싱 → 유효성 검사 → 서비스 위임 → 응답 매핑
- 직접 비즈니스 로직 작성 금지 (서비스 호출만)
- 에러 처리 미들웨어 또는 표준 패턴 활용

### 타입/스키마 (`type: types/schema/model/dto`)

- 정적 타입 언어: interface/type/class 등 적절한 형태 선택
- 검증 라이브러리(zod, pydantic, FluentValidation 등)는 프로젝트 패턴 따라 사용
- export/공개 범위 명시

---

## 에러 핸들링 (공통)

| 상황                                | 대응                                                            |
| ----------------------------------- | --------------------------------------------------------------- |
| 의존 파일 미존재                    | 명세 기반 타입 추론으로 진행, 보고서에 명시                     |
| 명세 불명확                         | 인근 파일 패턴 참조하여 가장 합리적인 방향으로 구현, 보고서에 명시 |
| 스택 매트릭스에 없는 subtype        | 가장 유사한 표준 패턴 차용, 보고서에 차용 사유 명시             |
| 정적 타입 시스템 충돌 (`any` 불가피) | 명시적 주석으로 사유 기록 후 사용                              |
