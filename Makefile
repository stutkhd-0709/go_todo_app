# phony -> 偽
# 通常ファイル名を登録するのがmakeだが、コマンドだけの場合もある (alias)
# aliasと実際のファイル名が衝突するのを防いでくれるらしい？
# あとパフォーマンスも上がる？
.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build:
	docker build -t stutkhd/gotodo:${DOCKER_TAG} \
		--target deploy ./

build-local: # Build docker image to local development
	docker compose build --no-cache

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

test:
	go test -race -shuffle=on ./..

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'