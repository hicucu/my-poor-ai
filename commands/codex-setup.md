---
description: my-poor-ai 에이전트와 multi-agent 기능을 ~/.codex/config.toml에 등록
allowed-tools: [Read, Edit, Write, Bash]
---

# my-poor-ai:codex-setup — Codex 설정 등록

my-poor-ai 플러그인의 에이전트와 multi-agent 기능을 `~/.codex/config.toml`에 등록함.

## 실행 절차

**1단계 — my-poor-ai 설치 경로 탐지**

이 파일(`commands/codex-setup.md`)이 위치한 디렉토리의 부모 디렉토리가 my-poor-ai 루트임.
플러그인 캐시 경로를 확인함:

Windows (PowerShell):

```powershell
$pluginRoot = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
# 또는 직접 확인:
ls "$HOME\.claude\plugins\cache\hicucu\my-poor-ai\"
```

Linux/macOS (Bash):

```bash
ls ~/.claude/plugins/cache/hicucu/my-poor-ai/
```

설치된 버전 디렉토리(예: `1.2.0/`)를 확인하고 `.codex/agents/` 경로를 파악함.

**2단계 — ~/.codex/config.toml 확인**

`~/.codex/config.toml` 파일을 읽어 `[features]` 섹션에 `multi_agent = true`가 이미 있는지 확인함.

이미 있으면 "multi_agent 이미 활성화됨" 메시지를 출력하고 3단계로 건너뜀.

**3단계 — multi_agent 기능 활성화**

파일이 없거나 `[features]` 섹션이 없으면 파일 끝에 추가:

```toml
[features]
multi_agent = true
```

`[features]` 섹션은 있지만 `multi_agent`가 없으면 해당 섹션에 `multi_agent = true` 추가.

**4단계 — 에이전트 전역 등록**

my-poor-ai 루트의 `.codex/agents/` 디렉토리에서 `*.toml` 파일 목록을 확인함.

`~/.codex/agents/` 디렉토리가 없으면 생성함.

각 `.toml` 파일을 `~/.codex/agents/my-poor-ai-{파일명}`으로 복사함.
(`my-poor-ai-` 접두사로 다른 플러그인과 충돌 방지)

실제 복사 (Windows PowerShell):

```powershell
$pluginAgentsDir = "<PLUGIN_ROOT>\.codex\agents"
$codexAgentsDir = "$HOME\.codex\agents"
New-Item -ItemType Directory -Force -Path $codexAgentsDir | Out-Null
Get-ChildItem "$pluginAgentsDir\*.toml" | ForEach-Object {
    Copy-Item $_.FullName "$codexAgentsDir\my-poor-ai-$($_.Name)"
}
$count = (Get-ChildItem "$codexAgentsDir\my-poor-ai-*.toml").Count
Write-Host "my-poor-ai 에이전트 $count개 등록 완료"
```

실제 복사 (Linux/macOS Bash):

```bash
plugin_agents="<PLUGIN_ROOT>/.codex/agents"
codex_agents="$HOME/.codex/agents"
mkdir -p "$codex_agents"
count=0
for f in "$plugin_agents"/*.toml; do
    cp "$f" "$codex_agents/my-poor-ai-$(basename "$f")"
    count=$((count + 1))
done
echo "my-poor-ai 에이전트 ${count}개 등록 완료"
```

**5단계 — 설정 검증**

```bash
python3 -c "
import tomllib, pathlib
cfg = pathlib.Path.home() / '.codex' / 'config.toml'
if cfg.exists():
    data = tomllib.loads(cfg.read_text())
    assert data.get('features', {}).get('multi_agent') == True, 'multi_agent not set'
    print('config.toml valid: multi_agent =', data['features']['multi_agent'])
else:
    print('WARNING: ~/.codex/config.toml not found')
"
```

에이전트 등록 확인:

Windows (PowerShell):

```powershell
$agentCount = (Get-ChildItem "$HOME\.codex\agents\my-poor-ai-*.toml" -ErrorAction SilentlyContinue).Count
Write-Host "~/.codex/agents/에 my-poor-ai 에이전트 $agentCount개 확인"
```

Linux/macOS (Bash):

```bash
count=$(ls ~/.codex/agents/my-poor-ai-*.toml 2>/dev/null | wc -l)
echo "~/.codex/agents/에 my-poor-ai 에이전트 ${count}개 확인"
```

**6단계 — 완료 메시지**

등록된 에이전트 수를 출력함:

```
my-poor-ai Codex 설정 완료.
- ~/.codex/config.toml: multi_agent = true
- ~/.codex/agents/: my-poor-ai 에이전트 등록 완료 (4단계 출력 참조)

codex 명령 실행 시 my-poor-ai 에이전트가 활성화됩니다.
서브에이전트 스킬(dispatching-parallel-agents 등) 사용을 위해 multi_agent 기능이 필요합니다.
```
