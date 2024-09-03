# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Подсчет покрытия кода тестами
В директории пректа нужно выполнить команды:

```
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
sed -i -e '/mock/d' profile.cov 
go tool cover -func profile.cov 
```

## Cборка проекта
```
cd cmd/shortener
BUILD_VERSION=v1.0.1 // указать актуальную версию
go build -ldflags "-X 'main.buildCommit=$(git rev-parse --short=8 HEAD)' -X 'main.buildVersion=$(echo $BUILD_VERSION)' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" .
```
