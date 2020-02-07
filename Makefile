build:
	@-rm -rf pi/html
	@cd backend && GOOS=linux GOARCH=arm GOARM=7 go build .
	@cp backend/barndoor-tracker-pi pi/
	@cd frontend && npm run build
	@cp -r frontend/build pi/html
