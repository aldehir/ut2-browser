.PHONY: build
build: ut2-browser

ut2-browser: frontend
	CGO_ENABLED=0 go build .

.PHONY: frontend
frontend:
	cd web && npm run build

.PHONY: clean
clean:
	rm ut2-browser
	rm -rf web/dist
