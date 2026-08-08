# Changelog

## [0.2.0](https://github.com/hexlet-basics/hexlet-basics/compare/hexlet_basics-v0.1.3...hexlet_basics-v0.2.0) (2026-08-08)


### ⚠ BREAKING CHANGES

* **api:** rename the member family and the course-readiness enum

### refactor

* **api:** rename the member family and the course-readiness enum ([fcff623](https://github.com/hexlet-basics/hexlet-basics/commit/fcff6235d70146662b49805199c3efaf0f110dbb)), closes [#766](https://github.com/hexlet-basics/hexlet-basics/issues/766)


### Features

* **admin:** ai lesson-review actions with openai-go worker ([fda0a91](https://github.com/hexlet-basics/hexlet-basics/commit/fda0a91f3755e8bd8fcde508caddfc46e60f483a))
* **admin:** blog post create/update/delete and related-courses set ([24d4e04](https://github.com/hexlet-basics/hexlet-basics/commit/24d4e04913a2467818762feff74148e17de59a96))
* **admin:** implement assistant messages list over ai_messages ([37ea6a9](https://github.com/hexlet-basics/hexlet-basics/commit/37ea6a9eb2cde8d91958a88cc774fa08ac9f51ae))
* **admin:** phase-3 admin screens — courses, blog posts, qna, permissions ([8eda7a3](https://github.com/hexlet-basics/hexlet-basics/commit/8eda7a3c106ac645d0cd6b1e54715ab4c0e87aed))
* **admin:** redirect admin root to the courses list ([46a8af9](https://github.com/hexlet-basics/hexlet-basics/commit/46a8af9ad867fa5553fc2b98435fa8dcd0fe2404))
* **admin:** wire all phase-2 resources into the frontend crud engine ([6fc2ddb](https://github.com/hexlet-basics/hexlet-basics/commit/6fc2ddb427c73573bde6b40683351f06123e9de3))
* **api:** design full frontend TypeSpec contract with zod client ([66bc928](https://github.com/hexlet-basics/hexlet-basics/commit/66bc928289ec6cfde5974e9792b11fe8b93cfffa))
* **auth:** implement email+password auth + port header/footer/login/register ([4126ea8](https://github.com/hexlet-basics/hexlet-basics/commit/4126ea809c9c52a1b7efa74a210436f36187d12b))
* **backend:** admin CourseCategory CRUD with integration tests ([3c77e5b](https://github.com/hexlet-basics/hexlet-basics/commit/3c77e5b0ef5b706d27b2b04b83917dbd01525ad6))
* **backend:** course loading via git fetch, build job, and GitHub webhook ([e5667e1](https://github.com/hexlet-basics/hexlet-basics/commit/e5667e1a64ff592156d3c7308cf63df97ba9027c))
* **backend:** enforce course-category uniqueness via DB constraints ([6107b1c](https://github.com/hexlet-basics/hexlet-basics/commit/6107b1c14086e1e907c30e9f94b00ab86427b917))
* **backend:** Go foundation, ent data layer, ogen API, goverter mapping ([2a7ff17](https://github.com/hexlet-basics/hexlet-basics/commit/2a7ff1758f5992fff521b565515eb9b0c38b6eb5))
* **backend:** implement admin attachments upload + blob read path ([deae432](https://github.com/hexlet-basics/hexlet-basics/commit/deae432e9719f09cd00d36a21a162848eb0b7799))
* **backend:** implement admin blog posts read (list/get) ([deb53ac](https://github.com/hexlet-basics/hexlet-basics/commit/deb53ac0e3b7f61a0e11f3e10b3ddcaf04a5341a))
* **backend:** implement admin course landing pages CRUD ([4421ad5](https://github.com/hexlet-basics/hexlet-basics/commit/4421ad5551674a90882055ca845a4f399b78c956))
* **backend:** implement admin course-lessons read lists ([2e2eae3](https://github.com/hexlet-basics/hexlet-basics/commit/2e2eae3a14b16c17d1bf8a47a712d33796f822d5))
* **backend:** implement admin courses list/get/create/update ([bebb9ac](https://github.com/hexlet-basics/hexlet-basics/commit/bebb9aced1297c6839ed924ecfb18fc0e7d600bb))
* **backend:** implement admin leads read-only list ([77e1cc3](https://github.com/hexlet-basics/hexlet-basics/commit/77e1cc3de4f85aa5a281a0e01f307da3bd0d3b42))
* **backend:** implement admin management users (list/get/update) ([31363a2](https://github.com/hexlet-basics/hexlet-basics/commit/31363a2ff5d612a1f830ed3bc2886036fefa1bd0))
* **backend:** implement admin QnA items CRUD (category + landing page) ([5e05fba](https://github.com/hexlet-basics/hexlet-basics/commit/5e05fba7bfa251f9204bc33b7ba0a01e71e6aaf3))
* **backend:** implement admin reviews CRUD ([ba2255e](https://github.com/hexlet-basics/hexlet-basics/commit/ba2255ea09e4e30275a0946f75ac51025d6bea52))
* **backend:** implement admin roles + role permissions ([66640c8](https://github.com/hexlet-basics/hexlet-basics/commit/66640c8e93f19ea218e5a19918485302fdcd939e))
* **backend:** implement admin staff members CRUD ([26e80e6](https://github.com/hexlet-basics/hexlet-basics/commit/26e80e6a5378b7589ffb52210e52428ccce1c7ab))
* **backend:** implement admin users CRUD + search ([b2cf00d](https://github.com/hexlet-basics/hexlet-basics/commit/b2cf00da2b31cdebea1b5f8957aa0993cfc27f42))
* **backend:** isolate async runtime and generate amocrm client ([bc1556f](https://github.com/hexlet-basics/hexlet-basics/commit/bc1556f8bfc10b696c28666c12665b66e4807dcd))
* **backend:** map every Course field from real data, drop goverter ignores ([2e8e90c](https://github.com/hexlet-basics/hexlet-basics/commit/2e8e90c15494274bdf87759db65800d2b29ebc45))
* **backend:** secure sessions and publish domain events ([3e5873b](https://github.com/hexlet-basics/hexlet-basics/commit/3e5873bd0964ceb6e40c295e16c520b2feadc2a1))
* **backend:** wire contract-first observability ([4824adb](https://github.com/hexlet-basics/hexlet-basics/commit/4824adb752912869d11eb55839adc770aea7b695))
* bootstrap Go + TypeSpec + generated-frontend stack; move Rails to legacy/ ([28b101c](https://github.com/hexlet-basics/hexlet-basics/commit/28b101c1d167bfae050f686634793b0efcedc226))
* **frontend:** add admin CRUD engine and wire course categories ([9ce8011](https://github.com/hexlet-basics/hexlet-basics/commit/9ce80114863191610dedf934b007da6cb091f6a6))
* **frontend:** forward JWT cookie over SSR and resolve current user ([90a0708](https://github.com/hexlet-basics/hexlet-basics/commit/90a070825a5a6767028ac47b7e218e7a442401a0))
* **frontend:** wire admin banners and add select/datetime to the CRUD engine ([92c65e4](https://github.com/hexlet-basics/hexlet-basics/commit/92c65e46a48f8799e3499631796d252819ea2765))
* localize backend responses ([20d373b](https://github.com/hexlet-basics/hexlet-basics/commit/20d373b5d0d40ae28f2174fa5d945b76a6af9560))


### Bug Fixes

* **api:** validate pagination query parameters ([f6b4028](https://github.com/hexlet-basics/hexlet-basics/commit/f6b4028cbdb9acfd5d135ca67f34a77bc0f00094))
* **app:** reload once on stale chunk import errors after deploy ([dcb801a](https://github.com/hexlet-basics/hexlet-basics/commit/dcb801afb5e548b3db75d494c9941aaf929276d3))
* **auth:** enforce contract-declared security ([acc8df1](https://github.com/hexlet-basics/hexlet-basics/commit/acc8df1386643a2542abdd53e9b2e4d36c6466d4))
* **auth:** resolve JWT sessions by user id ([2bb92b7](https://github.com/hexlet-basics/hexlet-basics/commit/2bb92b7f69cb5817b2ee6daf8fbade803cdca740))
* **backend:** re-hash atlas migrations to fix checksum mismatch ([aa2c34f](https://github.com/hexlet-basics/hexlet-basics/commit/aa2c34fdbf76eb44bf8dc1dbcf8a71dab429caf7))
* **ci:** approve @swc/core build script for pnpm ([329df5f](https://github.com/hexlet-basics/hexlet-basics/commit/329df5f8a0d7b64d72af4f6fb4961294eb583434))
* **ci:** re-sync lockfile with the pinned @hey-api/openapi-ts snapshot ([4bff502](https://github.com/hexlet-basics/hexlet-basics/commit/4bff502c07a2c9a3d80bfb0c94d06676357a2286))
* **courseloader:** confine image paths to locale root ([3022f14](https://github.com/hexlet-basics/hexlet-basics/commit/3022f1440f9a7c8530efee326fda42db5dde9e3b))
* **courseloader:** make content upserts atomic ([841e466](https://github.com/hexlet-basics/hexlet-basics/commit/841e466974e42fce1e213674356f072a5981396a))
* **courseloader:** reject unknown YAML fields ([2be3ab1](https://github.com/hexlet-basics/hexlet-basics/commit/2be3ab1fc55e7988bca04743e318ace88f79e6d3))
* **deps:** pin @hey-api/openapi-ts back to the next channel ([d4302df](https://github.com/hexlet-basics/hexlet-basics/commit/d4302df720bf5773f5086e4cf4f54c80ca85157b))
* **di:** return startup dependency errors ([35616d6](https://github.com/hexlet-basics/hexlet-basics/commit/35616d6816812756be13e7cc662cf8cf05cc0c3d))
* **events:** preserve legacy event fan-out semantics ([785634e](https://github.com/hexlet-basics/hexlet-basics/commit/785634e8074ea15fe03a7c5586fa8cc71f7e4be4))
* **exercise-loader:** reap course version builds stuck in building ([12254f6](https://github.com/hexlet-basics/hexlet-basics/commit/12254f688e03aae68938af0bcf2b33cff3d6e999))
* **jobs:** instrument River with OpenTelemetry ([0c6ceb1](https://github.com/hexlet-basics/hexlet-basics/commit/0c6ceb1be7674fa97a43edb8ccad813db9314569))
* **lessons:** quiet code-language Sentry noise and highlight lua ([fe101ee](https://github.com/hexlet-basics/hexlet-basics/commit/fe101eef7b5a189186830674b6fdf16f2dc5ec65))
* **lessons:** stop ShikiError on unmapped course code languages ([25964a9](https://github.com/hexlet-basics/hexlet-basics/commit/25964a9e74e11d5b1d948d8cddcb8a70aaa19d14))
* **localization:** validate required locale catalogs ([7482b26](https://github.com/hexlet-basics/hexlet-basics/commit/7482b264bfcf8ac0bac01881688fb643aac92d19))
* make course loader claims atomic ([8581c18](https://github.com/hexlet-basics/hexlet-basics/commit/8581c189ec89be1b61a852ae7aaedb19fd04007c))
* serve attachments with standard HTTP semantics ([73b3812](https://github.com/hexlet-basics/hexlet-basics/commit/73b38127f810409d96ba6867237dd3ce399b86a2))
* **telemetry:** capture HTTP panics with Sentry ([3223e19](https://github.com/hexlet-basics/hexlet-basics/commit/3223e1951ffaad4e7728a59a423087f5d25bfa89))


### Reverts

* **frontend:** keep explicit typed-selector admin form fields ([d2ce28c](https://github.com/hexlet-basics/hexlet-basics/commit/d2ce28cecab90ba48fbb6d27fced4355ba98ecd0))


### Miscellaneous

* add root Makefile with air + concurrently dev workflow ([3f45783](https://github.com/hexlet-basics/hexlet-basics/commit/3f457834f601a2dec107325c4bb814923c044567))
* **backend:** lay the river background-job queue backbone ([4c2a259](https://github.com/hexlet-basics/hexlet-basics/commit/4c2a259b057cc82919dac85e5cf24ce0907899ee))
* **backend:** narrow air watch to Go source dirs ([600062c](https://github.com/hexlet-basics/hexlet-basics/commit/600062c634bfda9366edf1488cdd03eb96da2b5e))
* **deps:** migrate the admin lists to TanStack Table v9 ([17b2522](https://github.com/hexlet-basics/hexlet-basics/commit/17b25229a97dc655be363a74d50f6ec9a5f7bd2f))
* **deps:** patch the flagged build-time JS dependencies ([5beca55](https://github.com/hexlet-basics/hexlet-basics/commit/5beca558cab9faa9c70d3eeac0163bd12c8a1fa7))
* **deps:** pin @hey-api/openapi-ts back to next channel ([b86bc14](https://github.com/hexlet-basics/hexlet-basics/commit/b86bc145422c454da6874c2286bbe6d72a7445d4))
* **deps:** update dependencies ([de3366a](https://github.com/hexlet-basics/hexlet-basics/commit/de3366a5c5fc84fa56db73b26cc5c1fc4ba666a6))
* **deps:** update frontend and Go dependencies ([5050065](https://github.com/hexlet-basics/hexlet-basics/commit/5050065b69f88a96e22466861f22fcc242c5243e))
* **deps:** update frontend and Go dependencies ([55c8dcd](https://github.com/hexlet-basics/hexlet-basics/commit/55c8dcd8dda8c7e4d049a69b14213729581fd693))
* **deps:** update frontend and Go dependencies ([8db0687](https://github.com/hexlet-basics/hexlet-basics/commit/8db0687916bb44247ef922773103a0cb14d0a1f1))
* **deps:** update frontend and Go dependencies ([fa5dd75](https://github.com/hexlet-basics/hexlet-basics/commit/fa5dd75d8f804173e8062f50be40740a16497a4a))
* **deps:** update to latest; TS 7 + hey-api next; consolidate api-spec ([14b4b86](https://github.com/hexlet-basics/hexlet-basics/commit/14b4b86f3c3ae71663debb058d6586c7486ee4af))
* **fixtures:** add legacy-&gt;Go testfixtures export tool ([51f30ce](https://github.com/hexlet-basics/hexlet-basics/commit/51f30ce35a68edf872f6ad3399f3742aba58e51c))
* **hooks:** manage hooks with lefthook + commitlint ([f0d6c81](https://github.com/hexlet-basics/hexlet-basics/commit/f0d6c813778a56a2e2b6a47674ede25e6258060d))
* **lint:** exclude generated route tree from biome, ignore .tanstack ([afbd600](https://github.com/hexlet-basics/hexlet-basics/commit/afbd600016d18454ca5b5da0026e916b6697b87d))
* **make:** add deps-update and update-skills targets ([6d7e5b2](https://github.com/hexlet-basics/hexlet-basics/commit/6d7e5b20b949431331e5ac6eca997707ac87116f))
* **sentry:** ignore errors coming only from third-party code ([075152d](https://github.com/hexlet-basics/hexlet-basics/commit/075152d71542cd9511846eaa5b542f4f22d91db8))
* **tooling:** pin golangci-lint via mise and clear lint findings ([645d5fc](https://github.com/hexlet-basics/hexlet-basics/commit/645d5fce1b20fcece0e29393b9e44a04c6406713))
* update dependencies and sync sorbet RBIs ([3e0fc5f](https://github.com/hexlet-basics/hexlet-basics/commit/3e0fc5f1d68e31b97dc96fdcd92c19f215bb1a27))
* update project infrastructure ([4eac22d](https://github.com/hexlet-basics/hexlet-basics/commit/4eac22d51df24f491961405b9c31ee0b7b8ebe29))
* update the agent skills to their latest upstream versions ([78e50d7](https://github.com/hexlet-basics/hexlet-basics/commit/78e50d7422fc52eaf538c5308703a1e1abb531d9))
* **web:** add oxlint and oxfmt ([25d822b](https://github.com/hexlet-basics/hexlet-basics/commit/25d822b33ec8beaa9c2f75acb20d55111ce40d7c))


### Build System

* **backend:** add dev-migrate target for dev + test DBs ([5a5c6fc](https://github.com/hexlet-basics/hexlet-basics/commit/5a5c6fc11df65beae25d3ae845ee7b3e6d4566a8))
* **backend:** express setup deps as make prerequisites ([b324673](https://github.com/hexlet-basics/hexlet-basics/commit/b324673c01c3a28f408957cf37ba67c4f594f22e))
* **backend:** load .env in config and add .env.example to setup ([e9e50b5](https://github.com/hexlet-basics/hexlet-basics/commit/e9e50b575846b725401d67932f4bb3472f024fae))
* **backend:** own the DB schema with atlas migrations, retire structure.sql ([f2579fe](https://github.com/hexlet-basics/hexlet-basics/commit/f2579feefd68b474e4066f280dc4a8ddb3642960))
* **backend:** split bootstrap into prepare (toolchain) + install (deps) ([91def48](https://github.com/hexlet-basics/hexlet-basics/commit/91def48226b060ac1114de33fe3936e83c6130e7))
* **backend:** split mise toolchain install into its own `make install` target ([0585bfa](https://github.com/hexlet-basics/hexlet-basics/commit/0585bfa743448d87a1a3d9aa47e87b99477ece74))
* ignore node_modules in go module walks ([93728ce](https://github.com/hexlet-basics/hexlet-basics/commit/93728cefca28b6daec86f710facc715955796717))
