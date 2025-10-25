comma := ,
project ?= go-lab
project_out_dir ?= out/$(service)/$(cmd)
service_path ?= services/$(service)
service_path_cmd ?= $(service_path)/cmd/$(cmd)
build_platforms ?= linux/arm64,linux/amd64
docker_registry ?= ko.local# replace with ghcr.io/brunoluiz/dev/ to push to the registry
docker_namespace ?= brunoluiz/$(project)
docker_repository ?= $(docker_registry)/$(docker_namespace)/dev/services/$(service)
docker_tag ?= $(shell git rev-parse HEAD)
docker_image ?= $(docker_repository)/$(cmd):$(docker_tag)
git_current_branch := $(shell git rev-parse --abbrev-ref HEAD)
git_base := $(if $(filter main,$(git_current_branch)),refs/remotes/origin/main~1,refs/remotes/origin/main)
OTEL_SERVICE_NAME=$(service)-$(cmd)

.PHONY: install
install:
	mise install
	go mod download

.PHONY: run
run:
	. $(service_path)/.env.default; \
	if test -f $(service_path)/.env; then . $(service_path)/.env; fi;
	export OTEL_SERVICE_NAME=$(service)-$(cmd); \
	docker compose -f ./$(service_path)/docker-compose.yaml up -d || true; \
	air --build.cmd "go build -o $(project_out_dir)/app ./$(service_path_cmd)" --build.bin "./$(project_out_dir)/app"

.PHONY: migrate
migrate:
	migrate -source file://./services/$(service)/internal/database/migration -database "$(DB_DSN)" up

.PHONY: format
format:
	golangci-lint fmt --enable gofumpt,goimports ./...
	prettier --write .

.PHONY: lint
lint:
	buf lint
	golangci-lint run --timeout 5m --color always --whole-files $(if $(files),$(files),./...)

.PHONY: scan
scan:
	trivy fs --exit-code 1 --no-progress --scanners vuln,misconfig,license .

.PHONY: test
test:
	go test -race ./...

.PHONY: monogo
monogo:
	@monogo detect --entrypoints $(shell find services -type d -name cmd -print0 \
	| xargs -0 -I {} find {} -maxdepth 1 -mindepth 1 -type d \
	| paste -sd ',' -) \
	--base-ref $(git_base) --compare-ref 'HEAD' --output github

.PHONY: docker-all
docker-all: docker-login docker-build docker-sign docker-scan

.PHONY: docker-login
docker-login:
	@echo $(docker_password) | docker login $(docker_registry) -u $(docker_user) --password-stdin

.PHONY: docker-build
docker-build:
	KO_DOCKER_REPO=$(docker_repository) \
	ko build ./$(service_path_cmd) \
		--base-import-paths \
		--tags $(docker_tag) \
		--platform $(build_platforms) \
		--sbom-dir $(project_out_dir) \
		--push

.PHONY: docker-sign
docker-sign:
	cosign sign --yes $(docker_image)

.PHONY: docker-scan
docker-scan:
	$(foreach platform,$(subst $(comma), ,$(build_platforms)),$(MAKE) docker-scan-platform platform=$(platform);)

.PHONY: docker-scan-platform
docker-scan-platform:
	trivy image --format sarif --platform $(platform) -o "$(project_out_dir)/$(cmd)-$(subst /,-,$(platform)).sarif" --scanners vuln,misconfig,license $(docker_image)

.PHONY: ci-debug
ci-debug:
	git show-ref
	git branch -a
	# env

ci-details:
		@echo "service=$$(echo "$$entrypoint" | cut -d'/' -f2)"
		@echo "cmd=$$(echo "$$entrypoint" | cut -d'/' -f4)"

kustomize-build:
	@for overlay in $$(find ./services/$(service)/kustomize/$(cmd)/overlays -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do \
		mkdir -p ./services/$(service)/manifests/$(cmd)/$$overlay/; \
		kustomize build ./services/$(service)/kustomize/$(cmd)/overlays/$$overlay > ./services/$(service)/manifests/$(cmd)/$$overlay/manifest.yaml; \
		echo "Generated manifests for ./services/$(service)/kustomize/$(cmd)/overlays/$$overlay"; \
	done

.PHONY: kustomize-patch-overlays
kustomize-patch-overlays:
	@for overlay in $$(find ./services/$(service)/kustomize/$(cmd)/overlays -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do \
		if [ -f ./services/$(service)/kustomize/$(cmd)/overlays/$$overlay/kustomization.yaml ]; then \
			sed -E -i.bak 's/newTag: latest/newTag: $(docker_tag)/g' ./services/$(service)/kustomize/$(cmd)/overlays/$$overlay/kustomization.yaml && rm ./services/$(service)/kustomize/$(cmd)/overlays/$$overlay/kustomization.yaml.bak; \
		fi \
	done

.PHONY: deploy-pr
deploy-pr: kustomize-patch-overlays kustomize-build
	git checkout -b deploy/$(service)/$(cmd)/$(docker_tag)
	@for overlay in $$(find ./services/$(service)/kustomize/$(cmd)/overlays -mindepth 1 -maxdepth 1 -type d -exec basename {} \;); do \
		git add ./services/$(service)/kustomize/$(cmd)/overlays/$$overlay/kustomization.yaml ./services/$(service)/manifests/$(cmd)/$$overlay/; \
		git commit -m "chore(deploy,$$overlay): $(service)/$(cmd)/$$overlay" || echo "No changes for $$overlay"; \
	done
	git push -f -u origin deploy/$(service)/$(cmd)/$(docker_tag)
	gh pr create --title "chore(deploy,$$overlay): $(service)/$(cmd)/$$overlay" --body "Automated deployment for commit $(docker_tag)" --base main
