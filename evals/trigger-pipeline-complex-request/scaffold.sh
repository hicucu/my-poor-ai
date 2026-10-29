#!/usr/bin/env bash
# 공용 픽스처 호출 — scaffold_script는 케이스 디렉토리 밖 경로를 거부하므로 래퍼를 둠
exec bash "$(dirname "$0")/../_fixtures/node-webapp.sh"
