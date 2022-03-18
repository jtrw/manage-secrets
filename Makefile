RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

.PHONY: run
run: ##run web service
	go run backend/app/main.go run

.PHONY: get
get: ##get value form kv storage
	go run backend/app/main.go kv get $(RUN_ARGS)

.PHONY: set
set: ##get value form kv storage
	go run backend/app/main.go kv set $(RUN_ARGS)