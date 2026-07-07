---
complexity: complex
depends-on: none
estimated-tasks: 9
---

# Spec A: 프로젝트 초기 셋업 (Vite + Svelte 5 + Vitest)

## 목표

빈 프로젝트에 Vite `svelte-ts` 스캐폴드를 생성하고 데모 콘텐츠를 제거한 뒤,
Vitest + `@testing-library/svelte` + `jsdom` 테스트 파이프라인을 구축한다.
`npm run dev` / `npm run build` / `npm run check` / `npm test` 가 모두 무오류로
동작하는 상태를 만들어 이후 스펙들이 TDD로 진행될 기반을 마련한다.

> 이 스펙은 스캐폴딩 성격이라 RED-GREEN 사이클 대신 **검증 절차**(빌드/체크 통과 +
> placeholder 테스트 1개로 vitest 파이프라인 동작 확인) 중심으로 작성한다.

## 구현 범위

### 변경 파일

- 생성(스캐폴드): `package.json`, `vite.config.ts`, `svelte.config.js`,
  `tsconfig.json`, `tsconfig.node.json`, `tsconfig.app.json`, `index.html`,
  `src/main.ts`, `src/App.svelte`, `src/app.css`, `src/vite-env.d.ts`, `.gitignore`
- 제거(데모): `src/lib/Counter.svelte`, `src/assets/`, `public/vite.svg` 등 데모 자산
- 수정: `vite.config.ts` (test 필드/플러그인 추가), `package.json` (`test` 스크립트),
  `tsconfig.app.json` (테스트 타입 추가), `index.html` (타이틀), `src/App.svelte` (최소화),
  `src/app.css` (데모 스타일 정리)
- 생성(테스트 인프라): `vitest-setup.ts`, `src/smoke.test.ts`

### 태스크 목록

- [ ] **태스크 1: Vite svelte-ts 스캐폴드 생성**

  현재 디렉토리에는 이미 `_workspaces/`, `REQUIREMENTS.md`, `.git/` 이 있어
  `npm create vite . ` 가 대화형 프롬프트로 멈출 수 있으므로, 임시 하위 디렉토리에
  생성 후 루트로 복사한다.

  ```bash
  npm create vite@latest .vite-scaffold -- --template svelte-ts
  cp -a .vite-scaffold/. .
  rm -rf .vite-scaffold
  ```

- [ ] **태스크 2: 의존성 설치 및 Svelte 5 확인**

  ```bash
  npm install
  npx svelte --version   # 또는: node -p "require('./node_modules/svelte/package.json').version"
  ```

  예상: `svelte` 5.x 가 설치됨. Vite 8.x, `@sveltejs/vite-plugin-svelte` 존재.

- [ ] **태스크 3: 데모 콘텐츠 제거**

  ```bash
  rm -f src/lib/Counter.svelte
  rm -rf src/assets
  rm -f public/vite.svg
  ```

  `src/App.svelte` 를 아래 최소 콘텐츠로 교체:

  ```svelte
  <script lang="ts">
  </script>

  <main class="app">
    <h1>Todos</h1>
  </main>

  <style>
    .app {
      max-width: 480px;
      margin: 2rem auto;
      font-family: system-ui, sans-serif;
    }
  </style>
  ```

  `src/main.ts` 는 스캐폴드 기본(App 마운트)을 유지하되, 삭제된 자산을 import 하지
  않는지 확인. `src/app.css` 는 데모 스타일을 지우고 최소 리셋만 남긴다:

  ```css
  :root {
    color-scheme: light dark;
    font-family: system-ui, sans-serif;
  }

  body {
    margin: 0;
  }

  * {
    box-sizing: border-box;
  }
  ```

  `index.html` 의 `<title>` 을 `Svelte Todo` 로 변경하고, 삭제한 `vite.svg`
  favicon 링크 라인을 제거한다.

- [ ] **태스크 4: 테스트 도구 설치**

  ```bash
  npm install -D vitest jsdom @testing-library/svelte @testing-library/jest-dom
  ```

  예상: 4개 devDependency 가 `package.json` 에 추가됨.

- [ ] **태스크 5: `package.json` 에 test 스크립트 추가**

  `scripts` 에 다음 항목을 추가한다 (기존 dev/build/preview/check 는 유지):

  ```json
  "test": "vitest run"
  ```

- [ ] **태스크 6: Vitest 설정 (vite.config.ts 확장)**

  `vite.config.ts` 를 아래로 교체한다. `@testing-library/svelte/vite` 의
  `svelteTesting()` 플러그인이 브라우저 resolve 조건과 테스트 자동 cleanup 을 처리한다.

  ```ts
  /// <reference types="vitest/config" />
  import { defineConfig } from 'vite';
  import { svelte } from '@sveltejs/vite-plugin-svelte';
  import { svelteTesting } from '@testing-library/svelte/vite';

  export default defineConfig({
    plugins: [svelte(), svelteTesting()],
    test: {
      environment: 'jsdom',
      globals: true,
      setupFiles: ['./vitest-setup.ts'],
    },
  });
  ```

  `vitest-setup.ts` 생성 (jest-dom 매처 등록):

  ```ts
  import '@testing-library/jest-dom/vitest';
  ```

- [ ] **태스크 7: TypeScript 테스트 타입 등록**

  `tsconfig.app.json` 의 `compilerOptions` 에 테스트 전역/매처 타입을 추가하여
  `npm run check` 가 테스트 파일에서 오류를 내지 않게 한다. (템플릿 버전에 따라
  `tsconfig.app.json` 이 없으면 `tsconfig.json` 에 추가)

  ```jsonc
  "types": ["vitest/globals", "@testing-library/jest-dom"]
  ```

  그리고 `include` 에 `vitest-setup.ts` 가 포함되도록 확인(보통 루트 `.ts` 는
  기본 include 됨; 누락 시 `"include": ["src", "vitest-setup.ts"]`).

- [ ] **태스크 8: placeholder 테스트로 파이프라인 검증**

  `src/smoke.test.ts` 생성:

  ```ts
  import { describe, it, expect } from 'vitest';

  describe('vitest pipeline', () => {
    it('runs a trivial assertion', () => {
      expect(1 + 1).toBe(2);
    });
  });
  ```

  전체 검증 명령 실행:

  ```bash
  npm test
  npm run check
  npm run build
  npm run dev -- --port 5173 &   # 기동 확인 후 종료 (Ctrl+C / kill)
  ```

  예상:
  - `npm test`: 1 passed (smoke.test.ts)
  - `npm run check`: 0 errors
  - `npm run build`: `dist/` 생성, 오류 없음
  - `npm run dev`: 로컬 서버가 `http://localhost:5173` 에서 기동

- [ ] **태스크 9: 커밋**

  ```bash
  git add -A
  git commit -m "chore: scaffold Vite Svelte5 TS project with Vitest test pipeline"
  ```

## 완료 기준

- [ ] Vite `svelte-ts` 스캐폴드가 프로젝트 루트에 생성됨 (Svelte 5.x).
- [ ] 데모 콘텐츠(`Counter.svelte`, 데모 자산)가 제거되고 `App.svelte` 가 최소화됨.
- [ ] `vitest`, `jsdom`, `@testing-library/svelte`, `@testing-library/jest-dom` 설치됨.
- [ ] `package.json` 에 `"test": "vitest run"` 스크립트 존재.
- [ ] `npm test` 가 smoke 테스트 1개를 통과.
- [ ] `npm run check` / `npm run build` 무오류, `npm run dev` 정상 기동.
- [ ] `_workspaces/` 는 삭제/변경되지 않음.

## 제외 범위

- 실제 Todo 기능/타입/컴포넌트 구현 (spec-b 이후에서 진행).
- CSS 디자인 (기능 완성 후 spec-h 에서 최소 스타일만).
- E2E/Playwright 설정.
