env-example:
	@echo "Creating .env.example file"
	@sed 's/=.*/=/' .env > .env.example
	@echo ".env.example file created."