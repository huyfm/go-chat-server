up:
	@docker compose up --build -d

log:
	@docker compose logs -f --tail 100 chat-server 

down:
	@docker compose down -v

reset: down, up

clean:
	@find . -type f -perm +111 -delete

.PHONY: up, log, down, reset, clean
