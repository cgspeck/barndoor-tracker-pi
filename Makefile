build:
	@-rm -rf pi/html
	@cd backend && GOOS=linux GOARCH=arm go build .
	@cp backend/backend pi/
	@cd frontend && npm run build
	@cp -r frontend/build pi/html
