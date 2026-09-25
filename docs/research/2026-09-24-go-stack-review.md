# Ревизия Go-стека: насколько удачно выбраны компоненты

**Вопрос.** Насколько хорошо выбран стек Go-переписывания (корень репозитория)?
Есть ли для каких-то компонентов решение лучше — с учётом ограничений проекта?

**Дата.** 2026-09-24. **Статус.** Только исследование, код не менялся. Итоговые
действия — в [«Если что-то менять»](#если-что-то-менять-top-3).

**Источники.** Первичные: сами ADR и файлы репозитория, исходники модулей из
`$GOMODCACHE` и `node_modules` (это опубликованное содержимое пакетов),
GitHub REST API (`gh api repos/...` — статус archived, дата последнего push,
лицензия, последний релиз), npm registry (`npm view`), официальная документация
(atlasgo.io, tanstack.com, yandex.cloud, README microsoft/kiota). Сторонние
обзоры и сравнительные статьи не использовались как основание ни для одного
утверждения. Там, где вывод сделан экспериментом, эксперимент описан и
воспроизводим (запускался в scratch-каталоге, репозиторий не трогал).

**Срез версий.** Исследование шло по рабочему дереву с незакоммиченным апдейтом
зависимостей (`go 1.26.4 → 1.27.1`, River `0.43 → 0.47`, go-pkgz/auth
`2.1.6 → 2.3.0`, TypeSpec `1.15 → 1.16`, hey-api `next → 0.99.0` и т. д.).
**Обновление 2026-09-24:** апдейт закоммичен (`5042ec67`, `d593d116`,
`ab473d55`, `1c908abe`); hey-api в нём оставлен на канале `next`
(`0.0.0-next-20260824173136`), так что проблема из §1.2 снята.

`docs/STACK.md`, на который ссылается `CLAUDE.md`, удалён коммитом `890ed208`
(«docs: remove obsolete stack plan»); источником истины теперь служат ADR
0001–0013. Ссылку в `CLAUDE.md` стоит убрать.

---

## Короткий вывод

Стек выбран хорошо. Ни один компонент не заархивирован, не сменил лицензию на
несовместимую и не имеет блокирующего ограничения, которое нельзя обойти без
смены библиотеки. Каждый выбор согласован с ключевыми ограничениями проекта:
contract-first (ADR-0001), жёсткий переход при паритете (ADR-0002) и «генерация
вместо ручного кода». Альтернативы, которые популярнее (GORM, oapi-codegen,
goose, fx), по этому критерию проигрывают: они требуют больше ручной обвязки.

Реальные проблемы лежат не в выборе библиотек, а в **расхождениях между ADR,
конфигами и кодом**:

1. ~~Переход на `@hey-api/openapi-ts@0.99.0` ломает `make gen-client`.~~
   Снято: в коммите `d593d116` пакет остался на `next` (см. §1.2).
2. Обход ogen для multipart не нужен: достаточно поправить одну строку в
   TypeSpec, и ogen сгенерирует операцию загрузки (проверено запуском).
3. ADR-0006 указывает не тот SES-клиент: Postbox реализует **SES v2 API**, а в
   ADR указан `aws-sdk-go-v2/service/ses` (v1).

---

## 1. Сводная таблица

Легенда вердиктов: **оставить** — менять незачем; **оставить, поправить** —
библиотека верная, но есть конкретный дефект в конфигурации или ADR;
**пересмотреть** — есть более удачный вариант, но не срочно; **заменить** —
в этом срезе таких нет.

| # | Слой | Выбрано (версия в дереве) | Реальные альтернативы 2026 | Вердикт |
|---|---|---|---|---|
| 1.1 | Контракт: описание | TypeSpec 1.16 → OpenAPI 3 | ручной OpenAPI YAML; Protobuf + Connect | **оставить** |
| 1.2 | Контракт: TS-клиент | hey-api openapi-ts (`0.0.0-next-20260824173136`) + TanStack Query + zod | orval 8.x; openapi-typescript + openapi-fetch | **оставить** |
| 1.3 | Контракт: Go-сервер | ogen 1.24.0 | oapi-codegen 2.8; Huma 2.39 (code-first); connect-go | **оставить, поправить** (multipart через TypeSpec) |
| 2 | ORM | ent 0.14.6 (+ свои шаблоны `adminput.tmpl`) | sqlc 1.31; bun 1.2; Jet 2.16; GORM 1.31 | **оставить** (следить за темпом релизов) |
| 3 | Драйвер БД | pgx v5 (`stdlib`); lib/pq — только тип `pq.StringArray` | — | **оставить** |
| 4 | Миграции | Atlas CLI 1.2.0 (mise) + `ariga.io/atlas` 1.3.0 (Go-библиотека в тестах) | goose 3.28; golang-migrate 4.20 | **оставить, поправить** (редакция бинаря, retention) |
| 5 | Конвертеры ent→api | goverter (1.9.4, есть 1.11.0) | ручной код; copier-подобные рефлексивные | **оставить** |
| 6 | DI | samber/do v2.1.0 + errgroup (ADR-0009/0010) | google/wire (архив); uber-go/fx; ручная сборка | **оставить** |
| 7 | Фоновые задачи | River 0.47 (riverdatabasesql) | asynq 0.26 (Redis); Temporal | **оставить** |
| 8 | Доменные события | Watermill 1.5.3 + watermill-sql v4 | «только River»; LISTEN/NOTIFY руками | **оставить** |
| 9 | Аутентификация | go-pkgz/auth v2.3.0 + golang-jwt v5 + x/crypto/bcrypt; go-webauthn — запланирован | markbates/goth; свой код поверх golang-jwt | **оставить** (поправить текст ADR) |
| 10 | Почта | Yandex Postbox через SES-совместимый API (не реализовано) | SMTP (go-mail) | **оставить, поправить** (`sesv2`, не `ses`) |
| 11 | Файлы | gocloud.dev/blob 0.46 | aws-sdk-go-v2/s3 напрямую; minio-go | **оставить** (поправить факт в ADR) |
| 12 | Ошибки/трейсинг | sentry-go 0.49 + OTel `otelconf` + samber/oops | — | **оставить** |
| 13 | Конфиг | caarlos0/env v11 + godotenv | knadh/koanf v2 | **оставить** |
| 14 | Логи | log/slog + lmittmann/tint | zap, zerolog | **оставить** |
| 15 | i18n | nicksnyder/go-i18n v2.6.1 | golang.org/x/text/message | **оставить** |
| 16 | Тесты | testcontainers-go 0.44 + testfixtures v3.19 + testify 1.12 | общая тестовая БД; txdb | **оставить** |
| 17 | Live reload | air 1.67.4 (`go tool`) | watchexec; reflex | **оставить** |
| 18 | Клиент amoCRM | kiota 1.34.1 (.NET-бинарь) + 3 рантайм-модуля kiota | ogen в режиме клиента | **пересмотреть** (низкий приоритет) |
| 19 | LLM-клиент | openai-go **v1.12.0** | openai-go v3.x | **пересмотреть** (два мажора позади) |
| 20 | Фронтенд SSR | TanStack Start 1.168.58 (статус RC) | React Router v7; Next.js | **оставить** (учитывать RC) |

Ниже — обоснование каждой строки с источниками.

---

## 1.1 TypeSpec как источник контракта — оставить

- Активно развивается: последний стабильный релиз `typespec-stable@1.16.0` от
  2026-09-09, репозиторий не в архиве, лицензия MIT
  ([releases](https://github.com/microsoft/typespec/releases/tag/typespec-stable@1.16.0)).
- Готового эмиттера Go-сервера у TypeSpec нет (есть только
  `@typespec/http-server-js` в `0.58.0-alpha`, `npm view @typespec/http-server-js`),
  поэтому схема «TypeSpec → OpenAPI → ogen» — единственный путь с генерацией
  Go-сервера из TypeSpec. Альтернатива «Protobuf + Connect»
  ([connect-go v1.21.0](https://github.com/connectrpc/connect-go/releases/tag/v1.21.0))
  несовместима с ADR-0002: публичные URL-маршруты Rails (`/languages/...`) нужно
  сохранить, а Connect навязывает пути вида `/package.Service/Method`.
- Минус один: TypeSpec эмитит конструкции, которые ogen не поддерживает
  (см. 1.3). Это решается правкой `.tsp`, как и предписывает ADR-0001
  («some OpenAPI constructs need the TypeSpec source adjusted»).

## 1.2 hey-api openapi-ts — оставить

**Обновление 2026-09-24.** Коммит `d593d116` оставил пакет на канале `next` и
поднял его до `0.0.0-next-20260824173136`; все шесть пинов `@hey-api/*` в
`minimumReleaseAgeExclude` сдвинуты вместе с ним. Контрольный прогон
`make gen-client` на этой сборке с TS 7 проходит, `src/client/` не меняется
(`git status` чистый). Ниже — исходный разбор, почему стабильная `0.99.0` не
годится; он остаётся в силе как предупреждение для будущих апдейтов.

**Факт в репозитории.** `CLAUDE.md` требует держать `@hey-api/openapi-ts` на
канале `next`, пока `typescript` 7.x. В незакоммиченном `package.json` пакет
переведён с `0.0.0-next-20260729142554` на стабильный `0.99.0`.

**Проверка.** Стабильная `0.99.0` опубликована 2026-06-22 (`npm view
@hey-api/openapi-ts time`), то есть она **старше** закреплённой next-сборки от
2026-07-29. Это даунгрейд. В собранном пакете стоит `import ts from
"typescript"` и обращения к `ts.SyntaxKind.*`
(`node_modules/.pnpm/@hey-api+openapi-ts@0.99.0_typescript@7.0.2/.../dist/init-D6Y8JFUS.mjs:8`).
Корневой экспорт `typescript@7.0.2` содержит только `version`/`versionMajorMinor`
(`node_modules/typescript/package.json`, поле `exports["."] → ./lib/version.cjs`).
Запуск `openapi-ts` на `api-spec/dist/openapi.yaml` (вывод в scratch-каталог)
падает так:

```
TypeError: Cannot read properties of undefined (reading 'AnyKeyword')
    at .../@hey-api+openapi-ts@0.99.0_typescript@7.0.2/.../init-D6Y8JFUS.mjs:4017:21
```

Ровно этот баг открыт в апстриме:
[hey-api/hey-api#4235](https://github.com/hey-api/hey-api/issues/4235) (open с
2026-07-09). PR с исправлением
[#4236](https://github.com/hey-api/hey-api/pull/4236) закрыт без слияния.
Заявленный peer-диапазон `typescript: ">=5.5.3 || ..."` формально пропускает
7.0.2, так что pnpm об этом не предупреждает.

**Контрольный прогон.** Ранее закреплённая `0.0.0-next-20260729142554` (она ещё
лежит в `node_modules/.pnpm`) с конфигом из `openapi-ts.config.ts` (вывод
перенаправлен в scratch-каталог) отрабатывает с кодом 0 на том же TS 7.0.2.
Результат совпадает с закоммиченным `src/client/` с точностью до `baseURL`,
который выводится из пути к входному файлу.

**Вывод.** Вернуть пин `0.0.0-next-20260729142554` — эта версия проверена.
Более свежая next-сборка `0.0.0-next-20260824173136` не объявляет
peer-зависимость от `typescript` (`npm view @hey-api/openapi-ts@next
peerDependencies` пусто), но это ещё не доказательство: `0.99.0` тоже
объявляла совместимый диапазон и всё равно падала. Поднимать до неё можно
только после такого же пробного прогона. Пин в `package.json` и шесть строк
`@hey-api/*@0.0.0-next-…` в `minimumReleaseAgeExclude` (`pnpm-workspace.yaml`)
нужно двигать вместе, как и сказано в `CLAUDE.md`.

**Сигнал на будущее.** Последний push в репозиторий hey-api был 2026-08-24, открытых
issues и PR — 639 (`open_issues_count`), последний стабильный релиз вышел 3 месяца назад
([репозиторий](https://github.com/hey-api/hey-api)). Проект жив, но TS 7 для
него пока не приоритет. Если к моменту паритета стабильной версии с поддержкой
TS 7 так и не будет, запасной вариант —
[orval v8.37.0](https://github.com/orval-labs/orval/releases/tag/v8.37.0)
(2026-09-23, MIT). Он тоже генерирует TanStack Query-хуки и zod-схемы. Но
миграция затронет `src/client/**` и CRUD-движок, который опирается на
сгенерированные `xxxOptions()` (ADR-0008), поэтому сейчас она не оправдана.

## 1.3 ogen — оставить; multipart чинится в TypeSpec, а не в Go

**Обновление 2026-09-24.** Сделано: `HttpPart<bytes>` в `api-spec/admin.tsp`,
`internal/apigen/ogen.yml` удалён, загрузка — сгенерированная операция
`AdminUploadAttachment`. Тесты в `internal/handlers/attachments_test.go` шлют настоящее multipart-тело
через `NewRouter` и сгенерированный декодер с `SecurityHandler` (in-process,
`httptest`; живой сервер и фронтенд не проверялись).
`internal/handlers/attachments.go` остался ради `GET /storage/{key}` — этого
маршрута нет в контракте. Ограничение размера тела держит `http.MaxBytesHandler`
в роутере. Изменилось поведение: форма без файла и слишком большое тело теперь
дают 400 `problem+json` от декодера, а не 422.

**Состояние.** ogen v1.24.0 (2026-08-07), Apache-2.0, активные коммиты после
релиза ([releases](https://github.com/ogen-go/ogen/releases/tag/v1.24.0)).
Задача про `requestBody.encoding` всё ещё открыта:
[ogen-go/ogen#1159](https://github.com/ogen-go/ogen/issues/1159) (с 2024-01).

**Откуда берётся трение.** `internal/apigen/ogen.yml` пропускает операции с
`form content encoding`, поэтому `POST /admin/attachments` обслуживается ручным
адаптером `internal/handlers/attachments.go` со своим дублированием
JWT/XSRF/admin-проверок (ADR-0011). Причина в исходниках ogen:
`gen/gen_contents.go:185-193` принимает для поля формы только encoding
`""`/form-urlencoded/JSON и на любой другой `contentType` возвращает
`ErrNotImplemented{"form content encoding"}`. TypeSpec для `HttpPart<File>`
эмитит `encoding.file.contentType: '*/*'` — отсюда и отказ.

**Эксперимент 1.** В копии `api-spec/dist/openapi.yaml` удалён блок `encoding`
у `/admin/attachments`. `go tool ogen` (v1.24.0) завершился с кодом 0,
операция `AdminUploadAttachment` сгенерирована. Для контроля: замена `*/*` на
`application/octet-stream` снова даёт отказ генерации.

**Эксперимент 2.** Минимальный `.tsp` скомпилирован TypeSpec 1.16:

| TypeSpec | Что попадает в OpenAPI |
|---|---|
| `file: HttpPart<File>` | `encoding.file.contentType: '*/*'` → ogen отказывает |
| `file: HttpPart<File<"image/png" \| "image/jpeg">>` | `contentType: image/png, image/jpeg` → ogen отказывает |
| `file: HttpPart<bytes>` | блока `encoding` нет, `type: string, format: binary` |

ogen на варианте `HttpPart<bytes>` генерирует поле `File ht.MultipartFile`.
В `MultipartFile` есть имя файла, размер и заголовки части, так что серверная
сторона ничего не теряет.

**Вывод.** Заменить в `api-spec/admin.tsp:39` `HttpPart<File>` на
`HttpPart<bytes>`, убрать `ignore_not_implemented` из `internal/apigen/ogen.yml` и удалить
ручной адаптер. Это сразу закрывает единственную операцию вне сгенерированного
слоя и единственное место, где политика безопасности задана вручную (ADR-0011).
Цена минимальна:

- **Схема поля не меняется.** В обоих вариантах это `type: string, format:
  binary`, меняется только наличие блока `encoding`. TS-клиент и сейчас типизирует
  поле как `file: Blob | File` (проверено генерацией next-сборкой hey-api), так
  что фронтенд не затронут.
- **Content-Type части не проверяется.** Без `encoding` part по умолчанию
  получает `application/octet-stream`, а браузер присылает `image/png` и т. п.
  Сгенерированный ogen-декодер (`decodeAdminUploadAttachmentRequest` в
  `oas_request_decoders_gen.go`) проверяет только внешний
  `Content-Type: multipart/form-data`, делает `r.ParseMultipartForm`, берёт
  `r.MultipartForm.File["file"]` и отдаёт `ht.MultipartFile{…, Header:
  fh.Header}`, не глядя на Content-Type самой части. Реальный тип файла
  доступен обработчику через `Header`.
- Сквозная проверка (реальный multipart-запрос к сгенерированному серверу)
  не делалась.

**Почему не менять генератор.** Альтернативы проигрывают по критериям ADR-0001:

- [oapi-codegen v2.8.0](https://github.com/oapi-codegen/oapi-codegen/releases/tag/v2.8.0)
  по-прежнему генерирует только типы и роутинг, а валидацию отдаёт рантайм-
  middleware (именно поэтому его отвергли в ADR-0001).
- [Huma v2.39.1](https://github.com/danielgtaylor/huma/releases/tag/v2.39.1) —
  code-first: OpenAPI выводится из Go-кода. Это переворачивает направление
  конвейера ADR-0001, и TypeSpec терял бы роль источника истины.

## 2. ent — оставить, но учитывать замедление релизов

- Последний релиз [v0.14.6](https://github.com/ent/ent/releases/tag/v0.14.6)
  вышел 2026-03-23. После него в `master` 8 коммитов за полгода (последний —
  2026-09-04; среди них есть содержательные: top-level UNION, CTE builder).
  Проект поддерживается, но медленно. Версия всё ещё 0.x, Apache-2.0.
- Проект глубоко завязан на кодоген ent: 32 схемы, собственный шаблон
  `adminput.tmpl` для create/update, goverter поверх ent-типов. Это ровно тот
  подход «генерация вместо ручного кода», который требует пользователь.
- Альтернативы:
  - [sqlc v1.31.1](https://github.com/sqlc-dev/sqlc/releases/tag/v1.31.1)
    (2026-04) требует вручную писать SQL на каждый запрос, включая
    admin-CRUD с фильтрами и пагинацией: больше ручного кода.
  - [Jet v2.16.0](https://github.com/go-jet/jet/releases/tag/v2.16.0)
    (2026-08) генерирует type-safe билдер из живой схемы. Идея близкая, но
    без мутаций/хуков/шаблонов ent.
  - [bun v1.2.18](https://github.com/uptrace/bun/releases/tag/v1.2.18) и
    [GORM v1.31.2](https://github.com/go-gorm/gorm/releases/tag/v1.31.2) —
    рефлексивные ORM без кодогенерации схемы.

  Ни одна не лучше для этого проекта. Замена означала бы переписать handlers,
  apiconv и шаблоны без выигрыша по паритету.
- Уже принятое решение «atlas владеет схемой, ent отвязан от миграций» (см.
  `atlas.hcl`) снимает главный риск ent — auto-migrate.

## 3. Драйвер PostgreSQL — оставить

Соединения открываются через `pgx/v5/stdlib` (`internal/store/store.go:12`,
`internal/testsupport/testdb/testdb.go:19`). `lib/pq` импортируется только ради
типа `pq.StringArray` в `internal/apiconv/apiconv.go:527`. lib/pq не в архиве,
релиз v1.12.3 от 2026-04-03 ([репозиторий](https://github.com/lib/pq)),
пометки о maintenance mode в текущем README нет. Удалять его как рискованный
незачем. При желании тип можно заменить на `pgtype`, но это косметика.

## 4. Atlas — оставить; поправить редакцию бинаря и пин

**Лицензии.** У Atlas две сборки CLI
([atlasgo.io/community-edition](https://atlasgo.io/community-edition),
[atlasgo.io/cli-reference](https://atlasgo.io/cli-reference)):

- стандартный бинарь — под **Atlas MSA** (проприетарная часть + OSS),
  бесплатен, Pro-функции открываются после `atlas login`;
- **Atlas Community** — собран только из OSS-репозитория, **Apache-2.0**.
  `migrate diff/apply/new` в нём есть. Нет `migrate lint/down/checkpoint/
  rebase/...`, `schema plan/lint`, а также поддержки views, триггеров, функций,
  sequences и RLS.

mise ставит Atlas через aqua-реестр по URL
`https://release.ariga.io/atlas/atlas-{OS}-{Arch}-{Version}`
([aqua-registry: pkgs/ariga/atlas](https://github.com/aquaproj/aqua-registry/blob/main/pkgs/ariga/atlas/registry.yaml)).
Это **стандартная (MSA)** сборка, а не `atlas-community-…`. Команды из
`Makefile` (`migrate new`, `migrate hash`, `migrate apply`) есть в обеих
редакциях. Go-библиотека `ariga.io/atlas` v1.3.0, которой тесты применяют
миграции, — Apache-2.0 ([репозиторий](https://github.com/ariga/atlas)).

**Retention.** Ariga удаляет бинарники релизов старше 6 месяцев с CDN и
Docker Hub, а поддерживает только две последние минорные версии
([cli-reference, раздел Distributed Binaries](https://atlasgo.io/cli-reference);
подтверждение мейнтейнера:
[ariga/atlas#3296](https://github.com/ariga/atlas/issues/3296#issuecomment-2574924728)).
Закреплённый `atlas = "1.2.0"` в `mise.toml` через полгода перестанет
скачиваться, и `mise install` у нового разработчика или в CI упадёт. Сейчас уже
вышла v1.3.0 ([release](https://github.com/ariga/atlas/releases/tag/v1.3.0)).

**Проверено утверждение из `mise.toml`.** `go list -m -versions
ariga.io/atlas/cmd/atlas` действительно заканчивается на `v0.13.1`, так что
пустить CLI через `go tool` нельзя.

**Вывод.** Atlas оставить: goose и golang-migrate — только раннеры
вручную написанных SQL, и у них нет `migrate diff` (планирования изменений из
желаемой схемы). Нужно одно из двух:

- (a) перейти на Community-сборку, раз Pro-функции не используются: чистая
  Apache-2.0 и нет вопроса о принятии MSA. Проверить, есть ли в aqua/mise
  отдельный пакет для community-бинаря; если нет — ставить его по
  `release.ariga.io/atlas/atlas-community-…` через `ubi`/`http`-бэкенд mise;
- (b) остаться на стандартной сборке, но зафиксировать в ADR, что она под MSA,
  и обновлять пин не реже раза в 6 месяцев.

Рекомендую (a): она дешевле по юридическим вопросам и ничего не отнимает у
текущего workflow.

## 5. goverter — оставить

Генерирует `internal/apiconv/apiconv.gen.go`, активный проект (релиз
[v1.11.0](https://github.com/jmattheis/goverter/releases/tag/v1.11.0) от
2026-09-11, MIT). В `go.mod` стоит 1.9.4 (как tool) — можно поднять вместе с
остальным апдейтом. Альтернатива — ручной маппинг, что противоречит правилу
«без ручного кода».

## 6. DI: samber/do v2 — оставить

- [samber/do v2.1.0](https://github.com/samber/do/releases/tag/v2.1.0)
  (2026-07-20, MIT), активный.
- [google/wire](https://github.com/google/wire) **заархивирован** (последний
  релиз v0.7.0, 2025-08-22) — этот вариант отпадает.
- [uber-go/fx](https://github.com/uber-go/fx) — последний релиз v1.24.0 от
  2025-05-13, берёт на себя ещё и жизненный цикл. ADR-0009/0010 сознательно
  оставили жизненный цикл за `errgroup`, так что fx дублировал бы его.
- Ручная сборка графа — больше ручного кода. Текущее разделение (do только
  конструирует, errgroup управляет жизненным циклом) — разумный минимум.

## 7. River — оставить

- [River v0.47.0](https://github.com/riverqueue/river/releases/tag/v0.47.0)
  (2026-08-31), лицензия **MPL-2.0** (файловый copyleft: обязательства
  касаются только изменённых файлов самого River, код приложения не
  затрагивается). Сам River 0.x, минорные версии могут ломать API.
- Главное достоинство для проекта — постановка задач в той же транзакции, что
  и бизнес-запись (ADR-0010: insert-only клиент в `cmd/server`). У
  [asynq v0.26.0](https://github.com/hibiken/asynq/releases/tag/v0.26.0) такого
  нет: он работает через Redis, а это новая инфраструктура.
  [Temporal](https://github.com/temporalio/sdk-go/releases/tag/v1.49.0) требует
  отдельного кластера — для джобов уровня «отправить письмо / собрать курс»
  это избыточно.

## 8. Watermill — оставить

[Watermill v1.5.3](https://github.com/ThreeDotsLabs/watermill/releases/tag/v1.5.3)
(2026-08-25, MIT). Вариант «только River» ADR-0004 уже разобрал по существу:
издатель стал бы знать всех подписчиков. Добавлю одно: второй механизм не
тянет новой инфраструктуры (оба работают на Postgres, таблицы создаёт Atlas),
поэтому цена двух систем — операционная, а не инфраструктурная. Менять не на
что.

## 9. Аутентификация — оставить; поправить текст ADR

- [go-pkgz/auth](https://github.com/go-pkgz/auth) — тег v2.3.0 от 2026-08-26,
  MIT, активный. Остаётся ли библиотека оправданной? Да: в коде используются
  её `token.Service` (JWT + XSRF double-submit) и `middleware.Authenticator`
  (`internal/handlers/auth.go:88-113`), то есть она избавляет от своей
  реализации cookie/XSRF/refresh.
- [go-webauthn/webauthn v0.18.2](https://github.com/go-webauthn/webauthn/releases/tag/v0.18.2)
  (2026-09-19, BSD-3) — живой, по-прежнему 0.x. В `go.mod` его **ещё нет**:
  passkey-операции существуют только в контракте (`docs/PARITY.md`, блокер 5).
- [markbates/goth](https://github.com/markbates/goth) (последний релиз v1.82.0
  от 2025-08-18) — только OAuth-провайдеры, без JWT/XSRF. Хуже текущего выбора.

**Расхождения ADR с кодом (исправить текст, не код):**

- ADR-0002 говорит, что bcrypt проверяется через `go-crypt/crypt`. На деле —
  `golang.org/x/crypto/bcrypt` (`internal/handlers/auth.go:465`,
  `internal/accounts/registrar.go:51`), а `go-crypt` в `go.mod` нет. Текущий
  вариант проще и правильнее, поправить нужно ADR.
- ADR-0003 говорит, что провайдер `direct` не используется. В
  `internal/handlers/auth.go:108-110` в `Authenticator` передаётся
  `provider.DirectHandler`. Судя по коду, он нужен как заглушка, чтобы
  `Authenticator` мог работать, а вход по паролю идёт мимо него. Стоит
  уточнить формулировку в ADR.

## 10. Почта через Postbox — оставить, но в ADR указан не тот клиент

Метод SendEmail в Postbox — `POST /v2/email/outbound-emails`
([Postbox API reference: SendEmail](https://yandex.cloud/en/docs/postbox/aws-compatible-api/api-ref/send-email)).
Это REST-API **SES v2**. В ADR-0006 указан `aws-sdk-go-v2/service/ses` —
клиент SES **v1** (Query API, `Action=SendEmail`). Писать надо на
`aws-sdk-go-v2/service/sesv2`. Кода почты пока нет (`docs/PARITY.md`, блокер 3),
поэтому достаточно поправить ADR до начала реализации. Само решение
«API, а не SMTP» разумно.

## 11. gocloud.dev/blob — оставить; поправить факт в ADR

[go-cloud v0.46.0](https://github.com/google/go-cloud/releases/tag/v0.46.0)
(2026-06-02, Apache-2.0, 2 открытых issue). Выбор верный: на разработке
`fileblob`, в проде `s3blob`, переключение одним URL.
Уточнение к ADR-0005: заархивирован **сервер** `minio/minio`
([репозиторий](https://github.com/minio/minio), archived), а клиент
[`minio/minio-go`](https://github.com/minio/minio-go) не в архиве (push
2026-09-15). На вывод это не влияет.

## 12–17. Наблюдаемость, конфиг, логи, i18n, тесты, air — оставить

Все активны и не в архиве (данные `gh api`):
[sentry-go v0.49.0](https://github.com/getsentry/sentry-go/releases/tag/v0.49.0),
[caarlos0/env v11.4.1](https://github.com/caarlos0/env/releases/tag/v11.4.1),
[tint v1.2.0](https://github.com/lmittmann/tint/releases/tag/v1.2.0),
[go-i18n v2.6.1](https://github.com/nicksnyder/go-i18n/releases/tag/v2.6.1),
[testcontainers-go v0.44.0](https://github.com/testcontainers/testcontainers-go/releases/tag/v0.44.0),
[testfixtures v3.19.0](https://github.com/go-testfixtures/testfixtures/releases/tag/v3.19.0),
[testify v1.12.1](https://github.com/stretchr/testify/releases/tag/v1.12.1),
[samber/oops v1.23.2](https://github.com/samber/oops/releases/tag/v1.23.2).
Замечания:

- **air** под **GPL-3.0** ([air-verse/air](https://github.com/air-verse/air)).
  Это инструмент разработки, в бинарь `bin/server` он не линкуется, так что на
  дистрибуцию лицензия не влияет. Но он записан как `tool` в `go.mod`, то есть
  попадает в граф модуля. Юридически это нейтрально; упоминаю для полноты.
- **koanf** ([v2.3.7](https://github.com/knadh/koanf/releases/tag/v2.3.7))
  понадобится, только если появятся файлы конфигурации. Пока конфиг задаётся
  через env, `caarlos0/env` проще.
- **testcontainers + Atlas-миграции + фикстуры** — сильное решение: тестовая
  схема не может разойтись с продовой, потому что применяются те же миграции.

## 18. Клиент amoCRM через kiota — пересмотреть (низкий приоритет)

- kiota жив ([v1.35.0](https://github.com/microsoft/kiota/releases/tag/v1.35.0)
  от 2026-09-04, в `mise.toml` — 1.34.1), поддержка Go в README отмечена как
  полная ✔ по генерации/абстракциям/HTTP
  ([README, Supported languages](https://github.com/microsoft/kiota#supported-languages)).
- Цена для проекта: отдельный .NET-бинарь в `mise.toml` и три рантайм-модуля
  (`kiota-abstractions-go`, `kiota-http-go`, `kiota-serialization-json-go`) ради
  **одной** операции `POST /api/v4/leads/unsorted/forms`. У
  [kiota-abstractions-go](https://github.com/microsoft/kiota-abstractions-go)
  19 звёзд — узкая база пользователей.
- ogen уже есть в toolchain и умеет то же самое: фича `paths/client`
  (`gen/features.go:96`) и фильтр операций `filters.path_regex` / `methods`
  (`gen/options.go:363-392`) — аналог `--include-path` у kiota. Спецификация
  amoCRM — собственный TypeSpec-вывод Hexlet (`AMOCRM_OPENAPI` в `Makefile`),
  то есть её можно поправить под ogen, если что-то не сгенерируется.
- **Рекомендация:** перейти на ogen-клиент, когда будут трогать
  `internal/amocrm` в следующий раз. Это минус один язык в toolchain и минус
  три зависимости. Не проверено: переварит ли ogen эту спецификацию
  (`application/hal+json` и т. п.) — нужен пробный прогон. Срочности нет:
  kiota не блокирует паритет.

## 19. openai-go — пересмотреть

В `go.mod` — `github.com/openai/openai-go v1.12.0`, а актуальная ветка —
[v3.66.0](https://github.com/openai/openai-go/releases/tag/v3.66.0)
(2026-09-23). Путь импорта меняется (`/v3`), так что обычный `go get -u` этого
не подтянет. Порт ассистента (in-app SSE) ещё впереди (`docs/PARITY.md`,
блокер 8), и лучше перейти до него, чем после.

## 20. TanStack Start — оставить, помня о статусе RC

Официальная документация всё ещё называет Start **Release Candidate**: API
считается стабильным, но это не v1
([overview](https://tanstack.com/start/latest/docs/framework/react/overview)).
Темп выпусков высокий — 10 версий `@tanstack/react-start` за последний месяц
(`npm view @tanstack/react-start time`). Это уже заметно в репозитории:
`@tanstack/react-start` и `react-router-ssr-query` закреплены точно, а в
`pnpm-workspace.yaml` около десятка исключений `minimumReleaseAgeExclude` для
пакетов TanStack. Альтернативы, разобранные в ADR-0008, не стали лучше: React
Router v7 и Next.js заменили бы роутер и модель загрузки данных. Вывод —
оставить и при апдейтах двигать связку `react-start`/`start-*`/`router-*`
атомарно.

---

## Если что-то менять (top-3)

1. ~~**Откатить `@hey-api/openapi-ts` на канал `next`**~~ — сделано в
   `d593d116` (`0.0.0-next-20260824173136`, `make gen-client` проверен).
   Стабильную `0.99.0` не брать, пока открыт
   [#4235](https://github.com/hey-api/hey-api/issues/4235).
2. ~~**Убрать multipart-обход ogen через TypeSpec**~~ — сделано (см. §1.3):
   `HttpPart<bytes>`, `ogen.yml` удалён, загрузка идёт через сгенерированный
   `SecurityHandler`. `attachments.go` остался только ради `GET /storage/{key}`.
3. **Привести ADR и тулчейн в соответствие с фактами**, пока соответствующий
   код не написан:
   - ADR-0006: `service/sesv2` вместо `service/ses`;
   - ADR-0002: `x/crypto/bcrypt` вместо `go-crypt`;
   - ADR-0003: формулировка про `DirectHandler`;
   - ADR-0005: какой именно MinIO заархивирован;
   - Atlas: перейти на Community-сборку (Apache-2.0) или записать в ADR, что
     используется MSA-сборка и пин нужно обновлять раньше, чем истекут
     6 месяцев хранения бинарников;
   - убрать ссылку на `docs/STACK.md` из `CLAUDE.md`.

## Итог

Стек менять не нужно. Он последовательно воплощает contract-first и
«генерацию вместо рукописного кода», все ключевые библиотеки живы, а у
альтернатив либо больше ручного кода (oapi-codegen, sqlc, goose, ручной DI),
либо лишняя инфраструктура (asynq/Redis, Temporal), либо они несовместимы с
сохранением URL (Connect). Из реальных рисков только два: зависимость от
hey-api next-канала, пока у него нет стабильной поддержки TS 7, и RC-статус
TanStack Start. Оба проявляются как трение при апдейтах, а не как ошибка
выбора. Более надёжной замены, которая устроила бы эти ограничения, сейчас
нет. Работа — в пунктах top-3, а не в смене библиотек.
