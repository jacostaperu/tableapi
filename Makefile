# Live Reload
export CGO_ENABLED=0
VERSION :="0.0.1"
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")


watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

bundle:
	# cd uiwrapper/ && cp properties.constants.ts vmo2-dx-bocc/vmo2-dx-bocc/
	# cd uiwrapper/vmo2-dx-bocc/vmo2-dx-bocc/ && nx build vmo2-dx-bocc  --baseHref=/new/ &&   rm -rf ../../../dist/* && cp dist/vmo2-dx-bocc/* -pr  ../../../dist/

build-server:
	@echo "Generating Templates"
	templ generate
	@echo "Building Server..."
	rm -f dist/server
	@go build -ldflags="-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)' -X 'main.mode=prod'"  -o dist/server cmd/main.go


build-dev:
	@echo "Generating Templates"
	templ generate
	@echo "Building Server..."
	rm -f tmp/server
	@go build -ldflags="-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)' -X 'main.mode=dev'"  -o tmp/server cmd/main.go

deploy:
	# scp main ossmlwweb001.unify.local:/home/acostaaguinagaj/boccplus
	# ssh ossmlwweb001.unify.local sudo service boccplus stop
	# ssh ossmlwweb001.unify.local sudo cp /home/acostaaguinagaj/boccplus /opt/boccplus/boccplus
	# ssh ossmlwweb001.unify.local sudo service boccplus start
pushoracle:
	goreleaser build --snapshot --clean
	scp dist/tableapi_linux_amd64_v1/server  opc@jjhub.duckdns.org:/home/opc/server
	ssh  opc@jjhub.duckdns.org sudo systemctl stop tableapi
	ssh opc@jjhub.duckdns.org ./install.sh

release:
	goreleaser build


setup:
	# install cron to update duckdns
	scp -r opc/* opc@jjhub.duckdns.org:/home/opc/
	scp dist/tableapi_linux_amd64_v1/server opc@jjhub.duckdns.org:/home/opc/server
	ssh  opc@jjhub.duckdns.org sudo systemctl stop tableapi
	ssh opc@jjhub.duckdns.org chmod 755 install.sh
	ssh opc@jjhub.duckdns.org ./install.sh
