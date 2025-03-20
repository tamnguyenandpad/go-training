submit-greeting-jobs: ## Submit greeting jobs with optional names, fmt: http://localhost:<httpPort>/greeting?name=Alice&name=Bob
	@read -p "Enter names (comma separated, default: Joe,Lily,Bob,Kevin,Henry,Denies,Tom): " names; \
	names=$${names:-Joe,Lily,Bob,Kevin,Henry,Denies,Tom}; \
	url="http://localhost:$(httpPort)/greeting?$$(echo $$names | sed 's/[^,]*/name=&/g' | tr ',' '&')"; \
	curl --location --request POST "$$url"