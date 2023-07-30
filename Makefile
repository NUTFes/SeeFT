RUN = api ls
GET = none

.PHONY: up
up:
	docker compose up 

.PHONY: up-db
up-db:
	docker compose up db

.PHONY: up-api
up-api:
	docker compose up -d db
	sleep 10
	docker compose up api

.PHONY: build
build:
	docker compose build

.PHONY: run
run:
	docker compose run --rm ${RUN}

.PHONY: down
down:
	docker compose down

.PHONY: exec
exec:
	docker compose exec api bash

.PHONY: migrate
migrate:
	docker compose run --rm api go run migrations/migrate.go up

.PHONY: tidy
tidy:
	docker compose run --rm api go mod tidy

.PHONY: go-init
go-init:
	docker compose run --rm api go mod init github.com/NUTFes/SeeFT/api

.PHONY: get
get:
	docker compose run --rm api go get -u ${GET}

.PHONY: vendor
vendor:
	docker compose run --rm api go mod vendor

.PHONY: storybook
storybook:
	docker compose run -p 6006:6006 --rm admin sh -c "npm run storybook"

.PHONY: seed
seed:
	docker compose run --rm api go mod tidy
	docker compose up -d db
	sleep 10
	docker compose run --rm api go run /app/seeds/seeds.go

# .PHONY: seed

# seed:
# 	docker compose run --rm server dart run ./sql/sql.dart seed
# 	docker compose run --rm server dart run ./sql/sql.dart user --csv ./sql/user.csv
#   docker compose run --rm server dart run ./sql/sql.dart task --csv ./sql/task.csv
#   docker compose run --rm server dart run ./sql/sql.dart task --csv ./sql/task_kikaku.csv
#   docker compose run --rm server dart run ./sql/sql.dart shift --csv ./sql/41st_shift_pre_sunny.csv -y41 -dpre -wsunny
#   docker compose run --rm server dart run ./sql/sql.dart shift --csv ./sql/41st_shift_1_sunny.csv -y41 -dpre -wsunny
#   docker compose run --rm server dart run ./sql/sql.dart shift --csv ./sql/41st_shift_2_sunny.csv -y41 -dpre -wsunny

