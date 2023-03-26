FROM cirrusci/flutter:2.10.4 AS build
WORKDIR /app
COPY . /app
RUN flutter pub get
RUN flutter build web

FROM python:3 AS production
COPY --from=build /app/build/web /app
WORKDIR /app
CMD ["python", "-m", "http.server", "45029"]
