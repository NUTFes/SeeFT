FROM ubuntu:20.04
RUN ln -fs /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# 必要なツールをインストールします。
RUN apt-get update && apt-get install -y \
    git \
    wget \
    unzip \
    xz-utils \
    libglu1-mesa \
    curl \
    python3 \
    python3-pip

# Flutter SDKをダウンロードしてパスに追加します。
ARG FLUTTER_SDK=/flutter
ENV FLUTTER_VERSION=3.10.6
RUN git clone https://github.com/flutter/flutter.git $FLUTTER_SDK
RUN cd $FLUTTER_SDK && git fetch && git checkout $FLUTTER_VERSION
ENV PATH="${PATH}:$FLUTTER_SDK/bin:$FLUTTER_SDK/bin/cache/dart-sdk/bin"
RUN chmod -R 777 $FLUTTER_SDK

# 新しいユーザーを作成します。
RUN useradd -ms /bin/bash flutteruser

# アプリのソースコードをコピーします。
COPY . /home/flutteruser/seeft_mobile
RUN chown -R flutteruser:flutteruser /home/flutteruser/seeft_mobile
RUN chmod -R 777 /home/flutteruser/seeft_mobile
RUN chown -R flutteruser:flutteruser $FLUTTER_SDK


USER flutteruser
WORKDIR /home/flutteruser/seeft_mobile


# Flutterの依存関係を取得し、アプリをビルドします。
RUN flutter doctor
# RUN flutter clean
RUN flutter pub get
RUN flutter pub upgrade
# webや他のプラットフォームに対してビルドする場合、適切なコマンドを使用してください。
RUN flutter build web

# FROM cirrusci/flutter:2.10.4 AS build
# WORKDIR /app
# COPY . /app
# RUN flutter pub get
# RUN flutter pub upgrade
# # webや他のプラットフォームに対してビルドする場合、適切なコマンドを使用してください。
# RUN flutter build web

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