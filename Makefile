.PHONY: up
up:
	docker compose up 

.PHONY: mac up
mac-up:
	docker compose -f docker-compose.mac.yml up 

.PHONY: prod up
prod-up:
	docker compose -f docker-compose.prod.yml up -d

.PHONY: mobile up
mobile-up:
	cd mobile && fvm flutter run -d web-server --web-port 45029 --dart-define-from-file=env/.env

.PHONY: up-db
up-db:
	docker compose up db

.PHONY: up-api
up-api:
	docker compose up -d db
	sleep 15
	docker compose up api

.PHONY: up-admin
up-admin:
	docker compose up -d db
	sleep 15
	docker compose up -d api
	docker compose up admin


.PHONY: build
build:
	docker compose build

.PHONY: mac build
mac-build:
	docker compose -f docker-compose.mac.yml build

.PHONY: prod build
prod-build:
	docker compose -f docker-compose.prod.yml build

.PHONY: down
down:
	docker compose down

.PHONY: exec
exec:
	docker compose exec api bash


.PHONY: tidy
tidy:
	docker compose run --rm api go mod tidy

.PHONY: go-init
go-init:
	docker compose run --rm api go mod init github.com/NUTFes/SeeFT/api


.PHONY: vendor
vendor:
	docker compose run --rm api go mod vendor

.PHONY: seed
seed:
	docker compose run --rm api go mod tidy
	docker compose up -d db
	sleep 15
	docker compose run --rm api go run /app/seeds/seeds.go

.PHONY: prod seed
prod-seed:
	docker compose -f docker-compose.prod.yml run --rm api go mod tidy
	docker compose -f docker-compose.prod.yml run --rm api go run /app/seeds/seeds.go

.PHONY: mac seed
mac-seed:
	docker compose -f docker-compose.mac.yml run --rm api go mod tidy
	docker compose -f docker-compose.mac.yml up -d db
	sleep 15
	docker compose -f docker-compose.mac.yml run --rm api go run /app/seeds/seeds.go

.PHONY: schemaspy
schemaspy:
	mkdir -p api/docs/schemaspy
	- docker compose run --rm schemaspy
	@echo "Extracting ER diagrams..."
	mkdir -p api/docs/er-diagrams
	find api/docs/schemaspy/diagrams -name '*.png' -exec cp {} api/docs/er-diagrams/ \;
	rm -rf api/docs/schemaspy
	@echo "ER diagrams saved to api/docs/er-diagrams/"

.PHONY: mac-schemaspy
mac-schemaspy:
	mkdir -p api/docs/schemaspy
	- docker compose run --rm schemaspy
	@echo "Extracting ER diagrams..."
	mkdir -p api/docs/er-diagrams
	find api/docs/schemaspy/diagrams -name '*.png' -exec cp {} api/docs/er-diagrams/ \;
	rm -rf api/docs/schemaspy
	@echo "ER diagrams saved to api/docs/er-diagrams/"

# mobile/lib/assetsに512*512のアイコンを用意しておくこと(コマンドのファイル名も変更する)
# リサイズ用にImageMagickをインストールする（`sudo apt-get install imagemagick` or `brew install imagemagick`）
.PHONY: mobile-icons-init
mobile-icons-init:
	cp mobile/lib/assets/44th_app-icon.png mobile/web/icons/Icon-512.png
	cp mobile/lib/assets/44th_app-icon.png mobile/web/icons/Icon-maskable-512.png
	if command -v convert >/dev/null 2>&1; then \
		convert mobile/lib/assets/44th_app-icon.png -resize 192x192 mobile/web/icons/Icon-192.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 192x192 mobile/web/icons/Icon-maskable-192.png; \
		echo "[INFO] Resized 44th_app-icon.png to Icon-192.png (192x192)"; \
		convert mobile/lib/assets/44th_app-icon.png -resize 192x192 mobile/web/splash/img/light-1x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 384x384 mobile/web/splash/img/light-2x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 512x512 mobile/web/splash/img/light-3x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 768x768 mobile/web/splash/img/light-4x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 192x192 mobile/web/splash/img/dark-1x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 384x384 mobile/web/splash/img/dark-2x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 512x512 mobile/web/splash/img/dark-3x.png; \
		convert mobile/lib/assets/44th_app-icon.png -resize 768x768 mobile/web/splash/img/dark-4x.png; \
		echo "[INFO] splash/img 各種サイズも変換しました"; \
	else \
		echo "[ERROR] ImageMagick (convert) が必要です。インストールしてください。"; \
		exit 1; \
	fi
	cd mobile && fvm flutter build web

# .PHONY: seed
# seed:
# 	docker compose run --rm server dart run ./sql/sql.dart seed
# 	docker compose run --rm server dart run ./sql/sql.dart user --csv ./sql/user.csv
#   docker compose run --rm server dart run ./sql/sql.dart task --csv ./sql/task.csv
#   docker compose run --rm server dart run ./sql/sql.dart task --csv ./sql/task_kikaku.csv
#   docker compose run --rm server dart run ./sql/sql.dart shift --csv ./sql/41st_shift_pre_sunny.csv -y41 -dpre -wsunny
#   docker compose run --rm server dart run ./sql/sql.dart shift --csv ./sql/41st_shift_1_sunny.csv -y41 -dpre -wsunny
#   docker compose run --rm server dart run ./sql/sql.dart shift --csv ./sql/41st_shift_2_sunny.csv -y41 -dpre -wsunny

