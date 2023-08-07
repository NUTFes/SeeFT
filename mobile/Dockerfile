FROM ubuntu:latest

# 必要なツールをインストールします。
RUN apt-get update && apt-get install -y \
    git \
    wget \
    unzip \
    xz-utils \
    libglu1-mesa \
    python3 \
    python3-pip

# Flutter SDKをダウンロードしてパスに追加します。
WORKDIR /opt
RUN wget https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.10.6-stable.tar.xz
RUN tar xf flutter_linux_3.10.6-stable.tar.xz
ENV PATH="$PATH:/opt/flutter/bin"

# 新しいユーザーを作成します。
RUN useradd -ms /bin/bash flutteruser

# アプリのソースコードをコピーします。
COPY . /home/flutteruser/seeft_mobile
RUN chown -R flutteruser:flutteruser /home/flutteruser/seeft_mobile
RUN chmod -R 777 /home/flutteruser/seeft_mobile


USER flutteruser
WORKDIR /home/flutteruser/seeft_mobile


# Flutterの依存関係を取得し、アプリをビルドします。
RUN flutter doctor
RUN flutter clean
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