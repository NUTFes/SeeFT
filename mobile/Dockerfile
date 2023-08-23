FROM cirrusci/flutter:2.10.4 AS build
WORKDIR /app
COPY . /app
RUN flutter pub get
RUN flutter pub upgrade
# webや他のプラットフォームに対してビルドする場合、適切なコマンドを使用してください。
RUN flutter build web

# エントリポイントを指定してコンテナを実行します。
# この例では、Flutter webサーバーをエントリポイントとして指定しています。
# CMD ["flutter", "run", "-d", "web-server", "--web-port=45029"]

# CMD ["python3", "-m", "http.server", "45029"]

# FROM cirrusci/flutter:3.7.7 AS build
# WORKDIR /app
# COPY . /app
# RUN flutter clean
# RUN flutter pub get
# RUN flutter build web

# FROM python:3 AS production
# COPY --from=build /app/build/web /app
# WORKDIR /app
# CMD ["python", "-m", "http.server", "45029"]