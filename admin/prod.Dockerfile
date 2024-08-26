# ベースイメージの作成
FROM node:16.13.0
# コンテナ内で作業するディレクトリを指定
WORKDIR /app/next-project/seeft-admin
COPY ./ /app

ENV NEXT_PUBLIC_APP_ENV production