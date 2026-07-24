## [0.6.1] - 2026-07-24

### 🚀 Features

- *(tz)* Embed Dhaka timezone data [`c754eaa`](https://github.com/AtifChy/aiub-companion/commit/c754eaa4f7e05c73adcf279cea0156e92ddad368)

### 🎨 Styling

- *(about)* Format build time display [`1d7b6f9`](https://github.com/AtifChy/aiub-companion/commit/1d7b6f900c470ea0574ebdb74b838731cc9deb68)

### ⚙️ Miscellaneous Tasks

- *(release)* Bump version to v0.6.1 [`2c669bb`](https://github.com/AtifChy/aiub-companion/commit/2c669bb95ec5da01916d761788856cc07db41702)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.6.0...v0.6.1

## [0.6.0] - 2026-07-23

### 🚀 Features

- *(deps)* Add tanstack query eslint plugin [`a9f2f6d`](https://github.com/AtifChy/aiub-companion/commit/a9f2f6da38c218589faa1379aa9fa3d748a68d1c)
- *(fetcher)* Add html fetching and dom parsing utilities [`24dc8b0`](https://github.com/AtifChy/aiub-companion/commit/24dc8b00c628c20c45bc22345779b30a56c07413)
- *(tz)* Add Asia/Dhaka timezone [`a627db1`](https://github.com/AtifChy/aiub-companion/commit/a627db1daefb81337371e9e7f41965cddeda54ec)
- *(fuzzy)* Add tie-breaker threshold for search [`1aec0c7`](https://github.com/AtifChy/aiub-companion/commit/1aec0c749e54db6da9b81026277be5c58cd215a0)
- *(calendar)* Get current or next exam [`c8205dc`](https://github.com/AtifChy/aiub-companion/commit/c8205dc367b5fb134c4fcdc97dd52adbc0602a16)
- *(semester)* Add ongoing exam status detection [`2edfefc`](https://github.com/AtifChy/aiub-companion/commit/2edfefcdc0b0fc3c08dc7e4452a0dc08029109fa)

### 🐛 Bug Fixes

- *(fetcher)* Handle nil node inputs gracefully [`0084b45`](https://github.com/AtifChy/aiub-companion/commit/0084b4582f149066138dc2191bb31711962dc335)
- *(notices)* Clear selection on selectionMode [`c69cc4a`](https://github.com/AtifChy/aiub-companion/commit/c69cc4a39a58588696eb08bed26a0742fd0f8c14)
- *(settings)* Unable to change settings [`551dc88`](https://github.com/AtifChy/aiub-companion/commit/551dc88ff58c8f654848b9b6e00f5e72e7ccc1b3)
- *(tz)* Embed `tzdata` for time zones [`ff336ff`](https://github.com/AtifChy/aiub-companion/commit/ff336ffa3495cb8feed26f30a4a8ef644135caf9)

### 🚜 Refactor

- *(calendar)* Use common fetcher utilities [`b9b4357`](https://github.com/AtifChy/aiub-companion/commit/b9b4357745675d915f973ef90e04394aa33f9327)
- *(notice)* Use common fetcher package [`ea00700`](https://github.com/AtifChy/aiub-companion/commit/ea0070068ce0cf1bddec36da83fc40c849b9c9ea)
- *(notices)* Rewrite read logic effects [`0e93d7b`](https://github.com/AtifChy/aiub-companion/commit/0e93d7b8c315d8c378fc2c785071f19ec4e385cc)
- *(db)* Update schema to use datetime and date types [`8cdf52b`](https://github.com/AtifChy/aiub-companion/commit/8cdf52bc1122e23ee885131e3692e0eebde7da8f)
- Use tz.Dhaka for timezone [`c425907`](https://github.com/AtifChy/aiub-companion/commit/c42590732d7e207de4ec17e146fc6d9f2fe518d3)
- Use time.Time for dates [`0db9690`](https://github.com/AtifChy/aiub-companion/commit/0db9690509c99543ec6ba4c6c73faa6e54c540b9)
- Use `toLocaleDateString` for date formatting [`64868f1`](https://github.com/AtifChy/aiub-companion/commit/64868f12b36d9eeb1243e96fb984e41c4774b10a)
- Remove unused date formatter and tests [`090b657`](https://github.com/AtifChy/aiub-companion/commit/090b6570c878a99defc2b4ce48f5a46278ed3285)
- *(fuzzy)* Adjust tie breaker threshold [`c97de01`](https://github.com/AtifChy/aiub-companion/commit/c97de01223b1d5d254a1ff493163565053110f3c)
- *(semester)* Use calendarKey helper function [`85734f9`](https://github.com/AtifChy/aiub-companion/commit/85734f94c33382053c80b35745ecc846f34bd3d9)
- *(notices)* Make the download btn container clickable [`715a44d`](https://github.com/AtifChy/aiub-companion/commit/715a44d3348af9276e80affeb74af45d4f5e8f66)

### 📚 Documentation

- Update changelog [`4099b4d`](https://github.com/AtifChy/aiub-companion/commit/4099b4d4eb90c13b5ceddecbaa4b1f2e58ceb767)

### 🎨 Styling

- *(semester)* Update exam date display format [`1b0eb6f`](https://github.com/AtifChy/aiub-companion/commit/1b0eb6f79bfebb03ff9229f42e015ff9c892bc27)

### 🧪 Testing

- Mock time in get next exam test [`5dd69e4`](https://github.com/AtifChy/aiub-companion/commit/5dd69e450c4e9db2a8e8563f13cddad0508d9f97)
- *(calendar)* Use fetcher in calendar scraper tests [`b56e0b3`](https://github.com/AtifChy/aiub-companion/commit/b56e0b38cdbab984acfd5ed530ee05f8e8322a4d)
- *(notice)* Refactor scraper tests to use fetcher [`ad0355f`](https://github.com/AtifChy/aiub-companion/commit/ad0355f014034fdcc20c0b2d7a15a1633721b5ec)
- *(fetcher)* Add fetcher and dom utility tests [`20a7244`](https://github.com/AtifChy/aiub-companion/commit/20a7244ceeaa326812c56e9d3ff9023ce15d6f48)

### ⚙️ Miscellaneous Tasks

- *(release)* Bump version to v0.6.0 [`d04b02a`](https://github.com/AtifChy/aiub-companion/commit/d04b02a84cec8fbdce3afdead06b259fb4c890f9)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.5.15...v0.6.0

## [0.5.15] - 2026-07-22

### 🚀 Features

- *(db:routine)* [**breaking**] Add course details to user routine [`710b5bb`](https://github.com/AtifChy/aiub-companion/commit/710b5bb1afa9805256eff38b88743a946e0cfc09)
- *(routine)* Handle adding existing course in routine [`1200ebb`](https://github.com/AtifChy/aiub-companion/commit/1200ebb89d19f16fdbbd4e1698bfeed4e4955ec3)

### 🚜 Refactor

- Remove unused cgpa and gpa-trend pages [`cbd23c7`](https://github.com/AtifChy/aiub-companion/commit/cbd23c783bbc521b4535ca8bd4f941c22587a856)

### 📚 Documentation

- Update changelog [`2edcc63`](https://github.com/AtifChy/aiub-companion/commit/2edcc630ab7b8c58b6f9dc04b2b2e87c7b688fd3)

### 🎨 Styling

- *(lint)* Fix linting issues in the codebase [`b43529c`](https://github.com/AtifChy/aiub-companion/commit/b43529cd51e243b1b8dfb4cf2ce13f0b57ca46cb)
- *(routine)* Update badge text size and remove ongoing border [`4078fc1`](https://github.com/AtifChy/aiub-companion/commit/4078fc13df120348cb386de0f4222fb0a5e2b542)

### ⚙️ Miscellaneous Tasks

- *(lint)* Replace `you-might-not-need-an-effect` with `react-doctor` [`e9197e6`](https://github.com/AtifChy/aiub-companion/commit/e9197e6319b5f8df492d9c56d8870a28c040f00d)
- *(release)* Bump version to v0.5.15 [`51b6835`](https://github.com/AtifChy/aiub-companion/commit/51b683503c89c805e2f8f28e6d391e8dee6f0755)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.5.12...v0.5.15

## [0.5.12] - 2026-07-21

### 🚀 Features

- *(settings)* Migrate to zustand store [`8cc3726`](https://github.com/AtifChy/aiub-companion/commit/8cc3726d13b27e9107b9fc705015c9908b5c522a)

### 🐛 Bug Fixes

- *(routine)* Handle file selection cancellation [`c483cec`](https://github.com/AtifChy/aiub-companion/commit/c483cece260586994db7ffbbabbd2006d2315e3f)

### 🚜 Refactor

- *(notices)* Update notice item click handler [`bbedc5e`](https://github.com/AtifChy/aiub-companion/commit/bbedc5e67fb2aa3aca99eb2f99b2807c2337469e)
- *(routine)* Improve course search popover logic [`43d17ae`](https://github.com/AtifChy/aiub-companion/commit/43d17ae4a2396a82f2637cf67668a4d73b87a263)

### 📚 Documentation

- Update changelog [`3b6e353`](https://github.com/AtifChy/aiub-companion/commit/3b6e353a111caf5305a2a676d771a8a2519ee659)

### ⚙️ Miscellaneous Tasks

- *(release)* Bump version to v0.5.12 [`5e5b5a3`](https://github.com/AtifChy/aiub-companion/commit/5e5b5a33384f51ee06aea0afdc158158d3800f00)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.5.10...v0.5.12

## [0.5.10] - 2026-07-21

### 🚀 Features

- *(notice)* Add pending notice handling [`8ebb700`](https://github.com/AtifChy/aiub-companion/commit/8ebb700e4a46f032f48f42e9f5488615a5fc9a93)
- *(notices)* Consume pending notice on load [`c5f11a8`](https://github.com/AtifChy/aiub-companion/commit/c5f11a83544e23dce832e4953ad68793b6e38203)
- *(notices)* Migrate state to zustand store [`36a3805`](https://github.com/AtifChy/aiub-companion/commit/36a38053ed49783293f84f0f12ee67e18744e3df)
- Escape key blurs search input [`3808f5d`](https://github.com/AtifChy/aiub-companion/commit/3808f5da5d47dfd5221fb1821ea7a57199b71db2)
- *(routine)* Improve course search ux and behavior [`7757fbc`](https://github.com/AtifChy/aiub-companion/commit/7757fbcdd7e7ae175de49d976ed2e85338594111)

### 🚜 Refactor

- *(notices)* Rename notice provider to selection provider [`bc61b9a`](https://github.com/AtifChy/aiub-companion/commit/bc61b9a93d17eb68e8e6853a0fadc853c8095a1e)
- *(notices)* Update import paths [`e309902`](https://github.com/AtifChy/aiub-companion/commit/e3099024ce67fa21eb18a73bd14a23c2a1325d04)

### 📚 Documentation

- Update changelog [`8d5f721`](https://github.com/AtifChy/aiub-companion/commit/8d5f7214e5365c8bf651dffae83469635867914b)

### 🎨 Styling

- Shorten arrow function bodies [`aff8144`](https://github.com/AtifChy/aiub-companion/commit/aff81448aaba7ed82268d3a46a78cc1e4ce54a45)

### ⚙️ Miscellaneous Tasks

- *(release)* Bump version to v0.5.10 [`bd7237f`](https://github.com/AtifChy/aiub-companion/commit/bd7237f4246605a99a5ce6e2a1c38184bc8590d4)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.5.5...v0.5.10

## [0.5.5] - 2026-07-20

### 🚀 Features

- *(notices)* Implement notice filters context [`4af599d`](https://github.com/AtifChy/aiub-companion/commit/4af599d25c8060077d0d5827c9217b50a9a13db3)

### 🐛 Bug Fixes

- *(notices)* Fix list not showing more than 50 notices [`0647632`](https://github.com/AtifChy/aiub-companion/commit/064763268b699c5163753bf71a4e772a1b88b78e)

### 🚜 Refactor

- *(notices)* Migrate notice selection to provider pattern [`d5bc07b`](https://github.com/AtifChy/aiub-companion/commit/d5bc07b912389d8e56f8e32e822189f954af6365)
- *(notices)* Move notice filters to provider [`611de11`](https://github.com/AtifChy/aiub-companion/commit/611de11b2bb5fa449d281df0182292aec7aa559d)
- *(notices)* Improve notice selection logic [`0ad98d6`](https://github.com/AtifChy/aiub-companion/commit/0ad98d629c0c64bc558f0f3006026a07c64122a0)

### 📚 Documentation

- Update changelog [`aba9eb8`](https://github.com/AtifChy/aiub-companion/commit/aba9eb81c72c450347958d4350ec41726b50945f)

### 🎨 Styling

- *(notices)* Update notice action bar styles [`9a0c2b6`](https://github.com/AtifChy/aiub-companion/commit/9a0c2b60e08895a65f6a611183a6da821f631421)

### ⚙️ Miscellaneous Tasks

- *(vite)* Update vendor chunking for ui and react [`127b7bc`](https://github.com/AtifChy/aiub-companion/commit/127b7bc54449c3833f4399b5ea0908d255d36762)
- *(release)* Bump version to v0.5.5 [`7bcc564`](https://github.com/AtifChy/aiub-companion/commit/7bcc564a6890e23578832fe13d5ae8a9088a6f06)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.5.0...v0.5.5

## [0.5.0] - 2026-07-20

### 🚀 Features

- *(frontend)* Add oxlint configuration [`b43505e`](https://github.com/AtifChy/aiub-companion/commit/b43505ec2f577684bb634a92310ec5abffcc98a6)
- *(frontend)* Add vitest for unit tests [`dc917ce`](https://github.com/AtifChy/aiub-companion/commit/dc917ce59882faba7cdd6bf1622af94e4e1a4345)
- *(db)* [**breaking**] Add new tables and modify existing ones [`93e9934`](https://github.com/AtifChy/aiub-companion/commit/93e993477afe542de2b52effdc51438054814443)
- *(db)* [**breaking**] Split class schedule from offered courses [`b3671bb`](https://github.com/AtifChy/aiub-companion/commit/b3671bb1188c6b10f0d049b0a675361bad6c3baf)
- *(routine)* [**breaking**] Add class schedule insertion [`07897d5`](https://github.com/AtifChy/aiub-companion/commit/07897d53d70570c0713db5f8304a06358bd6e876)
- *(tooltip)* Add side prop [`a65d7de`](https://github.com/AtifChy/aiub-companion/commit/a65d7de84dabef5b9d5301660191ec00f63f24dc)
- *(notices)* Extract notice mutations hook [`73299d0`](https://github.com/AtifChy/aiub-companion/commit/73299d03ef3f27abc575726bb317fd9a9cc0cb59)
- *(notices)* Add bulk selection and actions [`75c880b`](https://github.com/AtifChy/aiub-companion/commit/75c880b164eb7fda900bc1b337c4c872876f517f)

### 🐛 Bug Fixes

- *(db:notice)* Update notice on conflict [`1b42e58`](https://github.com/AtifChy/aiub-companion/commit/1b42e58fcfb15a979ecd99f069bda66ac8aa774f)

### 💼 Other

- Remove tools.go; use go.mod tool directive [`f6cf34d`](https://github.com/AtifChy/aiub-companion/commit/f6cf34d1e118ab7c58dbbc767a66dcf60958f9c4)

### 🚜 Refactor

- *(lint)* Fix linting issues in frontend components and hooks [`adc28f8`](https://github.com/AtifChy/aiub-companion/commit/adc28f8f6838f45faf3d78b9d4c895488e5494f6)
- *(notices)* Move notice types and utils to lib [`792d656`](https://github.com/AtifChy/aiub-companion/commit/792d6566200159987e1a916e0c57b1a8c99c4fa8)
- *(routine)* Extract course status logic [`6a5853e`](https://github.com/AtifChy/aiub-companion/commit/6a5853ef005e955e7deee4bbc67a1cb5493cb7ea)
- *(semester)* Move format event date to lib [`f6054e5`](https://github.com/AtifChy/aiub-companion/commit/f6054e5b822eedd51e8a0f8d05e48279db245383)
- *(db)* Rename `GetUserRoutine` to `ListUserRoutine` [`5a71f5c`](https://github.com/AtifChy/aiub-companion/commit/5a71f5c6aee3f0cbccecf750f5abf62f44a7e430)
- *(routine)* Split course and schedule data [`7e32fff`](https://github.com/AtifChy/aiub-companion/commit/7e32fffc79cb85158fe79e70017da7184780eb9e)
- *(routine)* Update routine logic for multiple schedules [`656d462`](https://github.com/AtifChy/aiub-companion/commit/656d462eb2da10efcf2535b7c7263532b95c687e)
- *(sql)* Rewrite and add class schedule index [`6366405`](https://github.com/AtifChy/aiub-companion/commit/63664053e1cd8dabc61dab97a8a80271d64885a6)
- *(db)* [**breaking**] Rename sql schema to migrations [`b9cc229`](https://github.com/AtifChy/aiub-companion/commit/b9cc2292b0749f53f66e9e8f55d654a89c7b09d5)
- *(db)* Rename conn to db [`fe11220`](https://github.com/AtifChy/aiub-companion/commit/fe11220eaa2c84d41513fb0b11388d9686b275ce)
- *(frontend:routine)* Move days array to routine lib [`2511ee3`](https://github.com/AtifChy/aiub-companion/commit/2511ee3c778d812ab93960b275345dc57366067f)
- *(routine)* Remove unused commented code [`ba45ba2`](https://github.com/AtifChy/aiub-companion/commit/ba45ba22aca68870786bf66790b960655a13d103)
- *(frontend)* Replace horizontal fade scroll with horizontal scroll [`616087b`](https://github.com/AtifChy/aiub-companion/commit/616087b1a0722b8024df2ad02f02d6f671935717)
- *(frontend:error)* Remove redundant error parsing utility [`5416dc1`](https://github.com/AtifChy/aiub-companion/commit/5416dc14b6a9f1e5a2e62a415a940ea6ddecc938)
- *(notices)* Centralize notice selection logic [`8e59ba5`](https://github.com/AtifChy/aiub-companion/commit/8e59ba5f364e42e55719e93d229cc0ceecf9f273)
- *(notices)* Simplify notice detail selection logic [`8048c25`](https://github.com/AtifChy/aiub-companion/commit/8048c250cc75fc67cb2f30fb28776fa81fb61786)

### 📚 Documentation

- Update changelog [`07e6edf`](https://github.com/AtifChy/aiub-companion/commit/07e6edf76b2ac6607daab674219ab8a78a6b9d8e)

### 🎨 Styling

- *(lint)* Disable confusing void expression rule [`579b49c`](https://github.com/AtifChy/aiub-companion/commit/579b49cbb6a4e506caa4941b1aa7ec3cf3fcb0dc)
- *(fmt)* Simplify onClick and onPressedChange props [`aba679d`](https://github.com/AtifChy/aiub-companion/commit/aba679d0722ca9d92770153a41d47dd480cc6ec2)

### 🧪 Testing

- *(lib)* Add new utility tests [`b758b33`](https://github.com/AtifChy/aiub-companion/commit/b758b330a9c9c9b50189c33dbf5d499dfa7427d3)
- *(routine)* Add parse course tests [`df4db1c`](https://github.com/AtifChy/aiub-companion/commit/df4db1c41033d33f18ced06812d07e10b14ff6bc)
- *(routine)* Update course to schedule in routine [`41cd8ef`](https://github.com/AtifChy/aiub-companion/commit/41cd8efdf6ea757fd9b6468826e53e5e15c65ca8)

### ⚙️ Miscellaneous Tasks

- *(script)* Update linter and formatter names [`78b8244`](https://github.com/AtifChy/aiub-companion/commit/78b8244bab68ecd4b8ea3f42d01d3e8fd2193ac8)
- *(lint)* Add vitest to oxlint config [`8eb15cd`](https://github.com/AtifChy/aiub-companion/commit/8eb15cd4927c19adbb9b5ac13674620456644de7)
- *(frontend:routine)* Remove redundant log statement [`57f45bb`](https://github.com/AtifChy/aiub-companion/commit/57f45bb28849a417492ff7cd566921d46a5d0600)
- *(git)* Typecheck already done by oxlint; skip from pre-commit hook [`0255c38`](https://github.com/AtifChy/aiub-companion/commit/0255c387db43e6be719d22a306e0fa28f7953621)
- *(release)* Bump version to v0.5.0 [`b4b3f04`](https://github.com/AtifChy/aiub-companion/commit/b4b3f04c27f95e80f1064e6c321e5cc092d1c8c2)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.4.0...v0.5.0

## [0.4.0] - 2026-07-17

### 🚀 Features

- *(event)* Add notice events [`1c1d561`](https://github.com/AtifChy/aiub-companion/commit/1c1d561e785df7e333c72abdf587f09dbf9ab19c)
- *(worker)* Rich notice notifications [`007e6f9`](https://github.com/AtifChy/aiub-companion/commit/007e6f94483275d1f621d072f722e98eca6ffde1)

### 🚜 Refactor

- *(routes)* Rename sections to SECTIONS constant [`b4e01d0`](https://github.com/AtifChy/aiub-companion/commit/b4e01d0107cadcbfed53eb2643b29a21c0bd2f6e)
- *(desktop)* Export window name constants [`c55122c`](https://github.com/AtifChy/aiub-companion/commit/c55122c0d6da2b08e5166506295853f1e579a064)
- *(db)* Update upsert and insert notice queries [`3a8ad0e`](https://github.com/AtifChy/aiub-companion/commit/3a8ad0e6c43530778f0ccb7a8841c9143d7cd41a)
- *(notice)* Use upsert notice for syncing [`62ee4df`](https://github.com/AtifChy/aiub-companion/commit/62ee4df55596d01aff98c937a05c9d906fab1862)
- *(notice)* Update notice event and selection logic [`0cd3114`](https://github.com/AtifChy/aiub-companion/commit/0cd3114163ff5578a7cc90e91c387b7a81037e10)
- *(autostart)* Remove custom implementation and use wails library [`470f5e0`](https://github.com/AtifChy/aiub-companion/commit/470f5e0ad4dbcf136373d73647648c75edd467a8)
- *(config)* Remove redundant log message [`19716f7`](https://github.com/AtifChy/aiub-companion/commit/19716f723940ce0d844f55ffa4b8b463b6b4f4bd)
- *(event)* Rename main window show event [`72b56a0`](https://github.com/AtifChy/aiub-companion/commit/72b56a0ce2ecf3bd9eb6cf35b42d08ec6500ddca)
- *(notices)* Split up notice list and toolbar components [`d6857af`](https://github.com/AtifChy/aiub-companion/commit/d6857af86e8b093dfe42b838c3429398c1248d7f)
- *(notices)* Rename category styles constant [`f90df96`](https://github.com/AtifChy/aiub-companion/commit/f90df963f4a5dddf851f37ef5aa61e1bd8362341)
- *(nav)* Extract nav item component [`df326c0`](https://github.com/AtifChy/aiub-companion/commit/df326c0b76e54529815607c888a0456af7d7b258)

### 📚 Documentation

- *(README)* Update `README.md` with latest features [`510c4d2`](https://github.com/AtifChy/aiub-companion/commit/510c4d2256a439f9edd56da2350d0a0364666b48)
- Update changelog [`c94ca0a`](https://github.com/AtifChy/aiub-companion/commit/c94ca0a10d2f1b3870a3d485100eaadadddb45c9)

### ⚙️ Miscellaneous Tasks

- Update prepare-release workflow [`bd42fda`](https://github.com/AtifChy/aiub-companion/commit/bd42fdacfb8f67d95a3d255e12f487dc72e076a7)
- *(cliff)* Remove skip ci commit parser [`5dc59f5`](https://github.com/AtifChy/aiub-companion/commit/5dc59f5ee53f180e108d25006c29f420063e48b4)
- *(release)* Bump version to v0.4.0 [`149315e`](https://github.com/AtifChy/aiub-companion/commit/149315e6950544e346ea9b217e7914ca5cfc6091)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.60...v0.4.0

## [0.3.60] - 2026-07-16

### 💼 Other

- Remove redundant version variable [`7bd8d49`](https://github.com/AtifChy/aiub-companion/commit/7bd8d4970a4b354945e9d6434e25d56be0312be8)

### 🚜 Refactor

- *(routes)* Remove unused routes and icons [`af6ec70`](https://github.com/AtifChy/aiub-companion/commit/af6ec70c71586a6b6cae305d85d67d085ea02431)
- Remove tools section from sidebar [`27c6c1c`](https://github.com/AtifChy/aiub-companion/commit/27c6c1c53cd5699001c2fb22b133e7bec1d9af95)
- *(route)* Update routes object structure [`2fbc179`](https://github.com/AtifChy/aiub-companion/commit/2fbc17941308a62c17f8df9cc45fb3a6a60ac91b)
- *(layout)* Simplify header breadcrumb logic [`7bb8590`](https://github.com/AtifChy/aiub-companion/commit/7bb8590b7c88acb7bf061434a6c1e74f43070b63)

### ⚙️ Miscellaneous Tasks

- *(release)* Bump version to v0.3.60 [`e22133d`](https://github.com/AtifChy/aiub-companion/commit/e22133d67b5c39957468f17c4b165ce2a42586bb)
- Update changelog [`5248973`](https://github.com/AtifChy/aiub-companion/commit/52489736989ffafd7ece7f3d653ac0bcc865c656)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.50...v0.3.60

## [0.3.50] - 2026-07-16

### 🚀 Features

- *(desktop)* Add main window show on notice click [`35953de`](https://github.com/AtifChy/aiub-companion/commit/35953dedf52679ba6bcc2a5265de54f62ed5cf05)
- *(autostart)* Add cross-platform autostart functionality [`cc5753a`](https://github.com/AtifChy/aiub-companion/commit/cc5753ad98e466ce6fcba9587cae296bc15d792f)
- *(config)* Add autostart option [`50a5ee4`](https://github.com/AtifChy/aiub-companion/commit/50a5ee4ac69ff84d3e21d1ee74b65c16931a0ea3)
- *(config)* Implement autostart functionality [`44acf17`](https://github.com/AtifChy/aiub-companion/commit/44acf17d58736f04def9676701ee53c143d8c9d4)

### 🐛 Bug Fixes

- *(config)* Garbled config validation error message [`c97dffb`](https://github.com/AtifChy/aiub-companion/commit/c97dffb643a2261628803e471c195e98f7e742de)

### 💼 Other

- *(docker)* Update golang base image to trixie [`0c8690a`](https://github.com/AtifChy/aiub-companion/commit/0c8690a9c7de3b2485989ff6fcbe30652909d814)
- *(meta)* Update ldflags module path [`58a8caa`](https://github.com/AtifChy/aiub-companion/commit/58a8caae7f9a083c0a3728f55e69211149b6289b)

### 🚜 Refactor

- *(updater)* Reuse updater not initialized error [`37d3856`](https://github.com/AtifChy/aiub-companion/commit/37d385631ab2aaf5b317640574d2f213ba2e993c)
- *(service)* Extract asset matcher function [`151641d`](https://github.com/AtifChy/aiub-companion/commit/151641db7cda30f01043e10406f39dc6293dcb4d)
- *(worker)* Simplify notification body message [`edab442`](https://github.com/AtifChy/aiub-companion/commit/edab442aeadcdadca46cb949b68f8fa0f40df136)
- *(autostart)* Remove service struct [`6e0857c`](https://github.com/AtifChy/aiub-companion/commit/6e0857cb6635322a72d12f781d62c74386cd0dde)
- Move build info to meta package [`ee79297`](https://github.com/AtifChy/aiub-companion/commit/ee792971ae2cd5a134fac90c93172c61d6cf41e7)
- *(config)* Reorder and add auto_start property [`370225c`](https://github.com/AtifChy/aiub-companion/commit/370225c62cc7ddeb550728c749066f88a31d2584)
- Rename config service to meta [`c7cf370`](https://github.com/AtifChy/aiub-companion/commit/c7cf3703c9905f440761e61a6db18dec46cbaeb7)
- *(config)* Extract config changed handler [`6ce0380`](https://github.com/AtifChy/aiub-companion/commit/6ce0380485992261cae9986746db1da0193a86fa)
- *(log)* Extract config changed handler [`80faf9b`](https://github.com/AtifChy/aiub-companion/commit/80faf9b20c352cb93dd7880b5acbb6fcb041c6b2)
- *(log)* Use meta app name [`6b4dc52`](https://github.com/AtifChy/aiub-companion/commit/6b4dc520d1aeb7ac7c2b6832c4127e2578a30d34)

### 🎨 Styling

- *(lint)* Fix lint issues [`3f34965`](https://github.com/AtifChy/aiub-companion/commit/3f349657d6d0d96e93f37850a1ae48de075a22ec)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`68b3138`](https://github.com/AtifChy/aiub-companion/commit/68b3138802fd74adfb62329a7244ad320f900b71)
- *(lint)* Add golangci-lint config [`9214915`](https://github.com/AtifChy/aiub-companion/commit/9214915feaab1db9f6bfb9adda746478adeb7782)
- *(release)* Remove changelog generation [`18e825f`](https://github.com/AtifChy/aiub-companion/commit/18e825f4be4a39824581ba8111963a31c0c07041)
- Add prepare release workflow [`bd2d2e6`](https://github.com/AtifChy/aiub-companion/commit/bd2d2e68c25fbe26f035d43ee4daffb0d8d7d552)
- Update setup-go and setup-node actions [`4d29716`](https://github.com/AtifChy/aiub-companion/commit/4d297165fc7422fd82976b09009967de721d09cc)
- *(release)* Bump version to v0.3.50 [`748f273`](https://github.com/AtifChy/aiub-companion/commit/748f273dd3c12329cf24085d5063e18fc4f7353b)
- Update changelog [`9bdcaca`](https://github.com/AtifChy/aiub-companion/commit/9bdcaca474c976ee483390ac79f09f71a80ab5b4)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.30...v0.3.50

## [0.3.30] - 2026-07-14

### 🚀 Features

- *(updater)* Add scheduled update checks [`3ec3c39`](https://github.com/AtifChy/aiub-companion/commit/3ec3c39c66c5da24ba8670b5c564bb3ee9e08c6d)

### 🐛 Bug Fixes

- *(database)* Add panic recovery to transactions [`30d4d98`](https://github.com/AtifChy/aiub-companion/commit/30d4d98ea45c42822697ff3e1848c42c9672a64c)
- *(theme)* Theme not applied globally [`1b71ba4`](https://github.com/AtifChy/aiub-companion/commit/1b71ba4ee50c4a99b4dacd145a1d232dbd685fb0)

### 🚜 Refactor

- *(desktop/window)* Improve window state management [`3a09286`](https://github.com/AtifChy/aiub-companion/commit/3a09286b231dfbe2319cbcdb40294673362e81b2)
- *(database)* Remove fts query sanitization [`41c7f0e`](https://github.com/AtifChy/aiub-companion/commit/41c7f0ea061d488cdc692701cedeb791e7fe9758)
- *(config)* Restructure config for clarity and consistency [`cb660b9`](https://github.com/AtifChy/aiub-companion/commit/cb660b949df762513f3a86783f3c9bde79add154)
- *(config)* Better align fields [`2a24995`](https://github.com/AtifChy/aiub-companion/commit/2a24995143bfbba9f7d3dc5d65949563c34c9507)
- *(database)* Rename sqlc generated package to sqlc [`a0fe7db`](https://github.com/AtifChy/aiub-companion/commit/a0fe7db30564706d76f489458b32ff0e853a35da)
- *(database)* Simplify transaction handling [`cb4c976`](https://github.com/AtifChy/aiub-companion/commit/cb4c9763d9a1826a40d5e2758ca0b517cc277f53)
- Rename db package to sqlc [`4c1ddfc`](https://github.com/AtifChy/aiub-companion/commit/4c1ddfc97bc81d2f4e6382e0b574c3ca30edc453)
- *(calendar)* Use testable time.now() [`29a0b10`](https://github.com/AtifChy/aiub-companion/commit/29a0b10743a41e91325db72d8fe4d49b1440498d)
- *(app)* Move layout and providers into app layout component [`d722352`](https://github.com/AtifChy/aiub-companion/commit/d722352f5a8c01c54d7c544f28f7b92dfe3c7231)

### 🎨 Styling

- *(routine)* Remove blank line at EOF [`82476f4`](https://github.com/AtifChy/aiub-companion/commit/82476f42efb30c8569240d9a468fa8d3d355877f)

### 🧪 Testing

- *(config)* Update config schema tests [`9b40cd1`](https://github.com/AtifChy/aiub-companion/commit/9b40cd13a956eed62c897817e715acf329727e2d)
- *(calendar)* Mock time for real world html tests [`e5550f1`](https://github.com/AtifChy/aiub-companion/commit/e5550f156659eb982d50bb25778f852a8239cfcd)

### ⚙️ Miscellaneous Tasks

- *(fmt)* Apply formatter [`e15be51`](https://github.com/AtifChy/aiub-companion/commit/e15be5190b1ad1b8d374087e35d0752f440ad7c5)
- *(settings)* Rename updateConfig to setConfig [`d46294f`](https://github.com/AtifChy/aiub-companion/commit/d46294f294cf299d3b2d273d8647aa88f09b398f)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.21...v0.3.30

## [0.3.21] - 2026-07-13

### 🚀 Features

- *(update/dialog)* Add dompurify to sanitize markdown output [`77262b0`](https://github.com/AtifChy/aiub-companion/commit/77262b06a57aabb7d09fc1b3b0453c5b1af8685b)

### 🐛 Bug Fixes

- *(updater)* Sanitizing markdown [`67baf64`](https://github.com/AtifChy/aiub-companion/commit/67baf64a230d09b62ef88c7381da29551920eb5e)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`87ca849`](https://github.com/AtifChy/aiub-companion/commit/87ca8493b3bbf73760db76f166481a7364018f58)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.20...v0.3.21

## [0.3.20] - 2026-07-13

### 🚜 Refactor

- *(routine)* Extract search result item component [`167562d`](https://github.com/AtifChy/aiub-companion/commit/167562dc9bbdd6163feeb77caa6646729c010ea7)
- *(routine)* Destructure useMutation return values [`6839bad`](https://github.com/AtifChy/aiub-companion/commit/6839bad062b58fabf5445481ace9ea31a17cb23c)
- *(help)* Extract faq section into component [`9666a66`](https://github.com/AtifChy/aiub-companion/commit/9666a6632caa0aa3cabea7ef477777029a1d48c6)

### ⚡ Performance

- *(routine)* Move search logic to component [`3c0b4a9`](https://github.com/AtifChy/aiub-companion/commit/3c0b4a9b639807c64728f7cbd3304327e55e326e)
- *(notices)* Split notice hooks and components [`8a8e90b`](https://github.com/AtifChy/aiub-companion/commit/8a8e90b0bfbdf744b3bf83fe65f5b12d61cc3125)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`c9fcfc9`](https://github.com/AtifChy/aiub-companion/commit/c9fcfc9e99e6ece282c7cd9569b4bee3032284e5)
- *(ci)* Update cliff config [`65f3f4e`](https://github.com/AtifChy/aiub-companion/commit/65f3f4e160181b4f92882481b4062c3da9f7c8c2)
- *(changelog)* Add new contributors and full changelog [`f8655d8`](https://github.com/AtifChy/aiub-companion/commit/f8655d893264ee5e3a1f312e96a86346b4136f49)
- *(scripts)* Rename pre-commit hook script [`a2130e9`](https://github.com/AtifChy/aiub-companion/commit/a2130e9090bb4d4e98ff94db968dbc509c611e93)
- *(lint)* Add eslint cache and concurrency flags [`13e4f20`](https://github.com/AtifChy/aiub-companion/commit/13e4f20c066ed780b80735a08bc3da816372c2a7)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.15...v0.3.20

## [0.3.15] - 2026-07-13

### 🚀 Features

- *(frontend)* Create update provider component [`d12c994`](https://github.com/AtifChy/aiub-companion/commit/d12c994509b29ec600c6ae328693269ef20a1c21)
- *(app)* Add update provider and dialog [`62141f4`](https://github.com/AtifChy/aiub-companion/commit/62141f43010e4fd55cbe06cd2b5d5b7e6d1fcb12)

### 🐛 Bug Fixes

- *(settings)* Fix double save on settings change [`90c8130`](https://github.com/AtifChy/aiub-companion/commit/90c8130de2cda906f04810516936803af942a37f)
- *(update/dialog)* Download & install button not working [`e133f4d`](https://github.com/AtifChy/aiub-companion/commit/e133f4d73a353a7663aa308757d5f45ffab338dd)

### 🚜 Refactor

- *(frontend)* Move providers to new directory [`2bd1177`](https://github.com/AtifChy/aiub-companion/commit/2bd1177809b7fa617369addfec8b28ffcf3dce8a)
- *(updater)* Simplify release handling [`ba2e2b0`](https://github.com/AtifChy/aiub-companion/commit/ba2e2b0a4f350cb3fd2e7c7b76914db3795f08bd)
- *(notices)* Update settings provider import path [`74414a5`](https://github.com/AtifChy/aiub-companion/commit/74414a577544a3ee4ae0c557d5bcd87be795b883)
- *(settings)* Remove update dialog and updater hook [`202efd8`](https://github.com/AtifChy/aiub-companion/commit/202efd8c06b5c154a1b4cebabe6550dd6c0eeee7)
- *(dialog)* Use update provider [`0b40e2c`](https://github.com/AtifChy/aiub-companion/commit/0b40e2cd52c8a393d7d95f50f5438dc543089b5e)
- *(updater)* Update release struct and handling [`3374b27`](https://github.com/AtifChy/aiub-companion/commit/3374b27a3bbeb252e109a7895b71eb6f4c5994f6)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`f3b65a5`](https://github.com/AtifChy/aiub-companion/commit/f3b65a5b7069914924dd70c3064d46e8317e9778)
- *(fmt)* Run formatter on frontend codebase [`8a838ff`](https://github.com/AtifChy/aiub-companion/commit/8a838ff2b44f5bacb7d7182112bb8ee991dc7659)
- *(git)* Add git pre-commit hooks [`2998e3f`](https://github.com/AtifChy/aiub-companion/commit/2998e3f5ff8a565020ab48046d41659c7cc4039a)
- *(github/workflow)* Pin pnpm version [`37a55cd`](https://github.com/AtifChy/aiub-companion/commit/37a55cdf8f30e09119abfd0cb2e8412425777d16)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.2...v0.3.15

## [0.3.2] - 2026-07-12

### 🚀 Features

- *(updater)* Add asset matcher to exclude installers [`5b798d1`](https://github.com/AtifChy/aiub-companion/commit/5b798d110bc6614b82b7e2711f539180df9e2771)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`2131721`](https://github.com/AtifChy/aiub-companion/commit/21317217fd659ab4833a1ab7ec09069eb5a0a891)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.3.0...v0.3.2

## [0.3.0] - 2026-07-12

### 🚀 Features

- *(updater)* Extract updater logic into a hook [`1458663`](https://github.com/AtifChy/aiub-companion/commit/1458663429ffb33b63619c06a11af5bb7a4dcd75)
- *(updater)* Print message when using github token [`42fb8ac`](https://github.com/AtifChy/aiub-companion/commit/42fb8ace897e11932b14b18ad2e1bc43a28ade50)
- *(notices)* Add fuzzy search [`da0afbf`](https://github.com/AtifChy/aiub-companion/commit/da0afbf51613f52851017b673be70a9e7d3ab78f)
- *(search)* Implement generic fuzzy search [`33484c5`](https://github.com/AtifChy/aiub-companion/commit/33484c5debb76d589962ce4432d6a2e6521fc662)

### 🚜 Refactor

- *(updater)* Split update download and install methods [`4ad0050`](https://github.com/AtifChy/aiub-companion/commit/4ad005082215be25ac61e58387c3198f40778c5d)
- *(settings)* Use tanstack query for settings provider [`fe268ac`](https://github.com/AtifChy/aiub-companion/commit/fe268ac1be9b27d30fcc0cbc79d1e9276049ddda)
- *(help)* Remove redundant diagnostics card [`84b8f61`](https://github.com/AtifChy/aiub-companion/commit/84b8f6172344fe6bc74bfd778c194534e2bfdfc2)
- *(settings)* Simplify settings initialization logic [`542ee75`](https://github.com/AtifChy/aiub-companion/commit/542ee75e99cb611a99fae895ae3d2c55ffaee637)
- *(notice)* Rename sqlite repository to repository [`fa34b65`](https://github.com/AtifChy/aiub-companion/commit/fa34b65025d2a2cfb6dcd5d9a9accde99dcbbb04)
- *(notice)* Use `time.DateOnly` constant [`cba37eb`](https://github.com/AtifChy/aiub-companion/commit/cba37ebf1bf6a1d3c9631fe06f694f7d6d9cc6ff)
- *(db)* Remove notices fts tables and triggers [`5b32e81`](https://github.com/AtifChy/aiub-companion/commit/5b32e814dc5bf152f361e9af4d3dad1a4b69aa6c)
- *(db)* Rename search to list notices [`39281fb`](https://github.com/AtifChy/aiub-companion/commit/39281fbe4301f9ca86649763d250fafbb3feb4cd)
- *(routine)* Rename sqlite repository to repository [`ef23f85`](https://github.com/AtifChy/aiub-companion/commit/ef23f85ebf0b3007e8f1e4f1c64907468b9898ba)
- *(db)* Remove fts schema and triggers [`e9279c0`](https://github.com/AtifChy/aiub-companion/commit/e9279c00b53a7eda6b79e88305e7cfa9f5eb9e23)
- *(routine)* Rename search to list offered courses [`073593f`](https://github.com/AtifChy/aiub-companion/commit/073593f33aa25e09daa04f5366e8cb5d3dd2a50a)
- *(routine)* Fuzzy search offered courses [`2f15a74`](https://github.com/AtifChy/aiub-companion/commit/2f15a74fcd2b1d59c69b49afc76997748c2f944e)
- *(sql)* Rename index schema file [`8e9d66b`](https://github.com/AtifChy/aiub-companion/commit/8e9d66bba2e96d41d1bd52fb0ba3f6b8bc360b1a)
- *(notice)* Rename rawQuery to query [`98ac7a3`](https://github.com/AtifChy/aiub-companion/commit/98ac7a34013983bc34919bc231d638dd314a2657)

### ⚡ Performance

- *(notice)* Limit search notices to 50 [`ea0461a`](https://github.com/AtifChy/aiub-companion/commit/ea0461a5d28401a06628d192e3aca19311dfd48c)

### 🎨 Styling

- *(updater)* Update markdown prose styles [`8219858`](https://github.com/AtifChy/aiub-companion/commit/8219858b9f9cb95cea05de83e060f74c97b98204)
- Adjust margin on scrollable containers [`7931219`](https://github.com/AtifChy/aiub-companion/commit/79312196634bb57c19a7ab6dcc9812d14d9906b3)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`543d4c3`](https://github.com/AtifChy/aiub-companion/commit/543d4c31662c1fbbb69c63411c0d03cdc44a37af)
- *(gitignore)* Add `bin/` [`388ed20`](https://github.com/AtifChy/aiub-companion/commit/388ed20f61e2c729777038e58a55f637f1b4a289)
- *(lint)* Switch back to `eslint` [`76462fc`](https://github.com/AtifChy/aiub-companion/commit/76462fc854aef1b0e88c6bb5ca0edf6cf6a07774)
- *(lint)* Fix linting issues in frontend components and hooks [`944ba3d`](https://github.com/AtifChy/aiub-companion/commit/944ba3d77b0c64e503ef3525be37daba300d2c54)
- *(lint)* Suppress `react-doctor/dangerous-html-sink` warning [`18f1cb0`](https://github.com/AtifChy/aiub-companion/commit/18f1cb0a384c64983daca3bb6683b5ef7d665a61)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.2.75...v0.3.0

## [0.2.75] - 2026-07-11

### 🚀 Features

- *(updater)* Add updater service [`c5f406f`](https://github.com/AtifChy/aiub-companion/commit/c5f406f1b978a54f368a1dc8650b4501e52a444b)
- *(main)* Initialize updater service [`0798665`](https://github.com/AtifChy/aiub-companion/commit/0798665f49703c37f2b257b3144db41fe2a03c81)
- *(frontend:updater)* Add update dialog and functionality [`1906203`](https://github.com/AtifChy/aiub-companion/commit/19062034d07f9247c1a3fa775ab1a8053fccdeea)

### 🚜 Refactor

- Apply linting rules to frontend [`f00184c`](https://github.com/AtifChy/aiub-companion/commit/f00184c3d47ffba8d68a306e4f8cda021af75671)
- *(notices)* Rename event unsubscribe variable [`f2a668c`](https://github.com/AtifChy/aiub-companion/commit/f2a668ccce5153b03ad6fe8fcd44d351a2d9d00f)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`a905fba`](https://github.com/AtifChy/aiub-companion/commit/a905fba94f59c8cb80bfdca2ac40731a3ba86c16)
- Switch to `oxfmt` for formatting [`45c4835`](https://github.com/AtifChy/aiub-companion/commit/45c48350c9257ed9d09b3b042011366b8d012206)
- *(fmt)* Apply `oxfmt` formatting on frontend files [`122dac8`](https://github.com/AtifChy/aiub-companion/commit/122dac86e72367013a87a9e0de136863c3036e03)
- *(lint)* Switch from `eslint` to `oxlint` [`3f7901a`](https://github.com/AtifChy/aiub-companion/commit/3f7901a2f70be86b392e8eff82d81aada417b3ef)
- *(lint)* Add oxlint plugins and rules [`1e8e8c6`](https://github.com/AtifChy/aiub-companion/commit/1e8e8c66c84e1689f6045fee04347ace435185b0)
- *(lint)* Update linting rules [`e8ea743`](https://github.com/AtifChy/aiub-companion/commit/e8ea743d64b0f196bfb152683f914b1b0e7aa930)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.2.61...v0.2.75

## [0.2.61] - 2026-07-09

### 🚜 Refactor

- *(routine)* Extract course parsing logic [`8a6163d`](https://github.com/AtifChy/aiub-companion/commit/8a6163d79b6199b4f75a4ab58f05cccffdde3971)
- Update launch time format [`ee5a1c7`](https://github.com/AtifChy/aiub-companion/commit/ee5a1c73e1453124b4bcbdab5ac55bf1d52ee431)
- Use http method constant [`4a3e415`](https://github.com/AtifChy/aiub-companion/commit/4a3e415da2cbfc2e026d126ef9ea75b3c613c394)
- Simplify window maximize logic [`274cb91`](https://github.com/AtifChy/aiub-companion/commit/274cb91074ea81f1ae3a4f0b687952601209cbfb)
- Use lowercase build info variables [`84812f1`](https://github.com/AtifChy/aiub-companion/commit/84812f16dae24a1a5398360b7c4db7fcea809712)

### 🎨 Styling

- Reorder struct fields [`1e82eda`](https://github.com/AtifChy/aiub-companion/commit/1e82eda624da4582128cd32620d1f3ac269770a5)
- Format defer calls concisely [`0de0a05`](https://github.com/AtifChy/aiub-companion/commit/0de0a05f7ddb8d9a7928dd716bfc39b0dcae1eaa)

### ⚙️ Miscellaneous Tasks

- Update changelog [skip ci] [`e55d86c`](https://github.com/AtifChy/aiub-companion/commit/e55d86c7c67eb4716ceb2163e26ed6e7eb92772b)
- *(ci)* Rename release artifacts by os and arch [`f7facb7`](https://github.com/AtifChy/aiub-companion/commit/f7facb727fe87e36fef2bc5e7a524aae8eae442c)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.2.50...v0.2.61

## [0.2.50] - 2026-07-07

### 🚜 Refactor

- Simplify theme provider logic [`9255a7e`](https://github.com/AtifChy/aiub-companion/commit/9255a7ea605d6d3a93ce7ef3288aa01c226d87d1)

### ⚙️ Miscellaneous Tasks

- *(ci)* Add github token to release workflow [`681ddff`](https://github.com/AtifChy/aiub-companion/commit/681ddff88a67724198345923b6527eb97907ef19)
- *(ci)* Add changelog generation workflow [`c1ea487`](https://github.com/AtifChy/aiub-companion/commit/c1ea48759e481e4204de25d8ee3bc5787b9d9052)
- Update changelog [skip ci] [`0408c54`](https://github.com/AtifChy/aiub-companion/commit/0408c541bd88c90b1ccff96bebec2efc369737f7)
- Comment out unsupported platforms [`2785700`](https://github.com/AtifChy/aiub-companion/commit/27857005a23eae7a8ca9c347d4cf1e5b7ce8cd44)
- *(ci)* Define github repo for release workflow [`e956f8c`](https://github.com/AtifChy/aiub-companion/commit/e956f8cfa26550bc80d6b26fc99479fb55d21170)

### ◀️ Revert

- *(ci)* Add github token to release workflow [`839410b`](https://github.com/AtifChy/aiub-companion/commit/839410b5b240d7bdd8ac954f165ba6f74c385f30)


### 👥 New Contributors
* @github-actions[bot] made their first contribution

> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.2.42...v0.2.50

## [0.2.42] - 2026-07-06

### 🚀 Features

- Add keyboard navigation to notice list items [`ecae2b2`](https://github.com/AtifChy/aiub-companion/commit/ecae2b23718e0c333e996fc1ee779baa75cbb0fe)
- *(routine)* Add keyboard navigation for course search [`0b2d9ea`](https://github.com/AtifChy/aiub-companion/commit/0b2d9eac5adf89bbf38b1693124aab4e178a5615)
- Add week parsing and display logic [`98380a5`](https://github.com/AtifChy/aiub-companion/commit/98380a58401d8195e958a3265e9344688b49293d)
- *(frontend)* Add react-scan and react-doctor [`06a1bb2`](https://github.com/AtifChy/aiub-companion/commit/06a1bb2aecc03224a04990c0911ed1bbcfd224cc)

### 🐛 Bug Fixes

- *(semester)* Use Date type for exam dates [`2899e11`](https://github.com/AtifChy/aiub-companion/commit/2899e11b1c3abffa2843afa315aa7d667dee2e55)

### 💼 Other

- Update Android and iOS build files [`aefdeec`](https://github.com/AtifChy/aiub-companion/commit/aefdeec0aefe3ab0982951dac82db15deaabf65a)
- Update build scripts [`3b16fa3`](https://github.com/AtifChy/aiub-companion/commit/3b16fa357b9d2da4f2a5811f7acbba206ffebe53)
- Update wails build files [`30693d1`](https://github.com/AtifChy/aiub-companion/commit/30693d1d1ba97d89e07dff499fc8316ba3fbdb46)

### 🚜 Refactor

- Split settings components into separate files [`bb72ec7`](https://github.com/AtifChy/aiub-companion/commit/bb72ec7fcf8bae70b95c8bedcd460c484e949ed1)
- Use react.use hook [`6bddd0c`](https://github.com/AtifChy/aiub-companion/commit/6bddd0cc02103942d98b4f5764418194fdebab29)
- Update use-mobile hook to use useSyncExternalStore [`0984fbe`](https://github.com/AtifChy/aiub-companion/commit/0984fbe71392337dd1192103e52567c167005126)
- *(routine)* Simplify useQuery and routine data access [`46f76a8`](https://github.com/AtifChy/aiub-companion/commit/46f76a87dabc7501d880b00acaacd41dc12de707)
- Optimize notice invalidation logic [`35e3c04`](https://github.com/AtifChy/aiub-companion/commit/35e3c04ff76456b19088115ca48d6f8688d9897b)
- Use functional useState update [`d8b95c5`](https://github.com/AtifChy/aiub-companion/commit/d8b95c547b8d3824a5bf525ab0e531b8fafe9e18)
- Update search input to functional component [`88cf0af`](https://github.com/AtifChy/aiub-companion/commit/88cf0af6a9f6f4b67b08ad9d923fed71ac9ff510)
- Simplify delayed loading logic [`10aa68e`](https://github.com/AtifChy/aiub-companion/commit/10aa68e11634f3ebb4a85244d3958594b12e322f)
- Improve notice read logic [`2bf842c`](https://github.com/AtifChy/aiub-companion/commit/2bf842cc2b9bfe3f8b98fcbc69395e2e766f5b0c)
- Cache current year in variable [`599c384`](https://github.com/AtifChy/aiub-companion/commit/599c384d464d7e26ecedc13cea706f01eff9e743)
- Improve a11y with semantic HTML [`c20f0ea`](https://github.com/AtifChy/aiub-companion/commit/c20f0ea5317a56fa6b1cfe15454600b549795af1)
- Remove unused env file [`f0231cb`](https://github.com/AtifChy/aiub-companion/commit/f0231cb358ca1f6f8cab42922da398b84c9279e6)
- *(frontend)* Remove unnecessary useCallback hook [`9358927`](https://github.com/AtifChy/aiub-companion/commit/93589272dd7cbd0db7db4fba93b855087d8bb34b)

### ⚙️ Miscellaneous Tasks

- *(lint)* Add doctor config file [`f54b644`](https://github.com/AtifChy/aiub-companion/commit/f54b6445d7ba154653c9b8b2a5c2964970877d78)
- *(package)* Add minimum age requirement and trustPolicy [`cfa3d22`](https://github.com/AtifChy/aiub-companion/commit/cfa3d2259dc4395ef16b0af6597790a29c013131)
- Update pnpm workspace config [`ea17562`](https://github.com/AtifChy/aiub-companion/commit/ea17562bd1cd4c08eed337e76e0e82e3fc241249)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.2.15...v0.2.42

## [0.2.15] - 2026-07-04

### 🚀 Features

- Implement lazy preload for routes [`bdbcb17`](https://github.com/AtifChy/aiub-companion/commit/bdbcb176eb74d7a59826c10ce176a8d5f4d04770)

### 🚜 Refactor

- Move loading component [`3c73f7a`](https://github.com/AtifChy/aiub-companion/commit/3c73f7ad37da3f7e96a363b9a617f4f7b6888eca)
- Use initial route constant [`9d09b01`](https://github.com/AtifChy/aiub-companion/commit/9d09b0168b7f39a38d5568da3ae0e05eae2a9cfb)
- *(calendar)* Simplify test parse function [`e5438fe`](https://github.com/AtifChy/aiub-companion/commit/e5438feb6031e57b183e800038cdfdb415d4581d)
- Remove unused header info extraction [`a9e33ec`](https://github.com/AtifChy/aiub-companion/commit/a9e33ecc5e1746944619ac3da9b2af982181c341)

### 🎨 Styling

- Improve scrollbar and layout consistency [`c31ad4b`](https://github.com/AtifChy/aiub-companion/commit/c31ad4bbf9d55e87ebd3486ca056aff7f78df3ff)
- Reorder imports [`9b28886`](https://github.com/AtifChy/aiub-companion/commit/9b28886cc682d08d440b6a2641062675fa013774)

### 🧪 Testing

- *(calendar)* Add scraper tests [`aa066cc`](https://github.com/AtifChy/aiub-companion/commit/aa066ccdd1b0c59e3409458ed8c3c30633b4a50d)

### ⚙️ Miscellaneous Tasks

- Update node version to 26.4.0 [`6cfe81c`](https://github.com/AtifChy/aiub-companion/commit/6cfe81c0ce443e2942c189c2632324ebab6b5cb7)
- Add typecheck script [`b715058`](https://github.com/AtifChy/aiub-companion/commit/b71505820f33fcbb464ee50600227e805bc354e9)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.2.5...v0.2.15

## [0.2.5] - 2026-07-04

### 🚀 Features

- *(db)* Add calendar cache schema and queries [`d52379d`](https://github.com/AtifChy/aiub-companion/commit/d52379d8c9132e2a3dac8f1780038dee36e57ae4)
- *(calendar)* Implement calendar scraping and caching [`56d55d2`](https://github.com/AtifChy/aiub-companion/commit/56d55d248d15eb168df1bff8c8b37ac19b2b75ff)

### 🚜 Refactor

- Remove semester column from calendar cache [`217bd04`](https://github.com/AtifChy/aiub-companion/commit/217bd04703c4453e0e8bf9fa5c29fcaedd7b0bc9)
- Remove redundant type assertion [`2800bb9`](https://github.com/AtifChy/aiub-companion/commit/2800bb96c679198f557b5a5de5fc87af7c29ad5a)
- Update semester page loading and display [`c941ceb`](https://github.com/AtifChy/aiub-companion/commit/c941ceb0f45bfc5c8565490e435d4d85ce73ab74)

### 🧪 Testing

- *(calendar)* Extract parse method for testing [`956d17d`](https://github.com/AtifChy/aiub-companion/commit/956d17dbae9b5380a7775bccfb60477dd62ea837)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.108...v0.2.5

## [0.1.108] - 2026-07-03

### 🚀 Features

- *(tray)* Add about menu item [`b1a8465`](https://github.com/AtifChy/aiub-companion/commit/b1a8465035b757bf7c30bd5afbfafa798d152cda)

### 🚜 Refactor

- Rename constant and update permissions [`a79a7b5`](https://github.com/AtifChy/aiub-companion/commit/a79a7b51a2b18ba866444ce914a676f1677c9459)
- *(db)* Remove unused querier interface [`85e777b`](https://github.com/AtifChy/aiub-companion/commit/85e777b37c234ea06bed322c51caf607d9f82863)
- Simplify error messages [`951fe41`](https://github.com/AtifChy/aiub-companion/commit/951fe413d4c40ba49cc8268f66975a4154d19dd2)

### ⚡ Performance

- Reduce quit delay [`353bfca`](https://github.com/AtifChy/aiub-companion/commit/353bfcaff8cf36ced5920ad6907b8cb3f5912cea)

### ⚙️ Miscellaneous Tasks

- Update nsis-install action [`94e9d3e`](https://github.com/AtifChy/aiub-companion/commit/94e9d3e8e551d08b986b1ceecfef39f0ee020cbf)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.102...v0.1.108

## [0.1.102] - 2026-07-02

### 🚀 Features

- Add calendar service and parser for semester scheduling [`db561fe`](https://github.com/AtifChy/aiub-companion/commit/db561fe793748b043939ece10e97f5b1db67f252)
- Add calendar service initialization [`aa85db2`](https://github.com/AtifChy/aiub-companion/commit/aa85db2f7c12b16f77d52181ab90a3a1a47a471e)
- Add deep equality check for settings save [`3b19971`](https://github.com/AtifChy/aiub-companion/commit/3b199719a1dc0a49e5f9bb6b5e3c2e6f85bedf98)
- Integrate next-themes for theme management [`a992292`](https://github.com/AtifChy/aiub-companion/commit/a99229201803de9649def5ba61d3075dea598715)
- Move toaster to app root [`38ce937`](https://github.com/AtifChy/aiub-companion/commit/38ce937558ae938e25ccc5a01437da59679dcd37)
- *(worker)* Add service shutdown method [`fc94684`](https://github.com/AtifChy/aiub-companion/commit/fc94684e7fa3b646b2131eed7cb43684dbab19d1)
- Attach to parent console on Windows [`1699853`](https://github.com/AtifChy/aiub-companion/commit/169985389aa6bcb22993767663e1cd7595bfe83b)
- Add semester dashboard page [`e93041a`](https://github.com/AtifChy/aiub-companion/commit/e93041a92ec263b70c5941d84ece3935aace7d11)
- Add generic search input component [`c864fad`](https://github.com/AtifChy/aiub-companion/commit/c864fadc299bc9e3160548338e0f0ebce55f17a7)

### 🐛 Bug Fixes

- Update shadcn css path [`75d597b`](https://github.com/AtifChy/aiub-companion/commit/75d597bae9680885f8a56807835d71cb1b719051)
- Move skipNextSave to after setConfig [`0999852`](https://github.com/AtifChy/aiub-companion/commit/099985237f95eea9dfd2c139e91fdb2f9db4c461)

### 🚜 Refactor

- *(calendar)* Reorder category rules [`e71250a`](https://github.com/AtifChy/aiub-companion/commit/e71250a0b98759db6b1017f8e8bdeb6aeb17c7e0)
- Use react-query mutations for notice actions [`f8027be`](https://github.com/AtifChy/aiub-companion/commit/f8027be6ef40394d1d543852d457e168af58787f)
- Remove separator, add border to card header [`7f7580c`](https://github.com/AtifChy/aiub-companion/commit/7f7580c7ca3050970e085d3425564d7980a1a933)
- Remove unused theme context [`2274a93`](https://github.com/AtifChy/aiub-companion/commit/2274a93505b35d5169436d7c3bca11f612924b80)
- Improve routine page ui and logic [`5b3b292`](https://github.com/AtifChy/aiub-companion/commit/5b3b29263f750bc03205b4958630ac3dc53319a0)
- Improve auto mark notices as read logic [`71321ca`](https://github.com/AtifChy/aiub-companion/commit/71321ca61cdd7876ff4720bf90dbf875112cd215)

### 🎨 Styling

- Add margin to avoid overlap of scrollbar and resizing window [`a5e2c4b`](https://github.com/AtifChy/aiub-companion/commit/a5e2c4b3204e41012d92c5b2ab355dd44f50b4c3)

### 🧪 Testing

- *(scraper)* Update expected attempts in test [`ba5c1a0`](https://github.com/AtifChy/aiub-companion/commit/ba5c1a00332746491e6ea6c4413d438c46a53951)

### ⚙️ Miscellaneous Tasks

- Update go dependencies [`1c11249`](https://github.com/AtifChy/aiub-companion/commit/1c112496aaea1ae73b734b1a6ab5bb1717a18763)
- Update shadcn components [`4876fb3`](https://github.com/AtifChy/aiub-companion/commit/4876fb33ac1b67bfdec546c65156c1e768a2b44e)
- *(frontend)* Update dependencies [`bfe13ea`](https://github.com/AtifChy/aiub-companion/commit/bfe13ea28f7b3d3f1d878d62c85b6b09200d2cc9)
- Update wails to v3.0.0-alpha2.111 [`afeb565`](https://github.com/AtifChy/aiub-companion/commit/afeb56508ce2ca1a57c8707958cee1a989d7581c)
- Update tsconfig paths and options [`8abc6be`](https://github.com/AtifChy/aiub-companion/commit/8abc6bee9f747f69f70274f86c54d06d17beae5c)

### ◀️ Revert

- Feat: add deep equality check for settings save [`91ae7a5`](https://github.com/AtifChy/aiub-companion/commit/91ae7a520935b12e8a6045fd6e4219970d783bbb)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.76...v0.1.102

## [0.1.76] - 2026-06-28

### 🚀 Features

- *(settings)* Implement reset config [`9a9be37`](https://github.com/AtifChy/aiub-companion/commit/9a9be371a97a8aaa69690c81d6bdc21010a34c20)
- Adjust sidebar width [`d2321d5`](https://github.com/AtifChy/aiub-companion/commit/d2321d523f390ed369d9e3685ee006ece86114bc)

### 🐛 Bug Fixes

- *(window)* Regression with restore window with config fixed [`6fe8b18`](https://github.com/AtifChy/aiub-companion/commit/6fe8b18b2554cd037a19b00c68fecf6be7434cac)
- Control alert dialog open state [`ea4acac`](https://github.com/AtifChy/aiub-companion/commit/ea4acaccef21bb28899afd2bbe2f853b5d8fa283)
- *(worker)* Return type of emitted event [`908fe3c`](https://github.com/AtifChy/aiub-companion/commit/908fe3c1193a5f9eb8329560fe5445f82d631074)
- Increase fetch html max attempts [`0820b03`](https://github.com/AtifChy/aiub-companion/commit/0820b03c89ef6f72208f7802fef852f26c16470b)
- Unsubscribe from notices synced [`03bd492`](https://github.com/AtifChy/aiub-companion/commit/03bd49298d7bead3027e28913cbfe0bbadb9313e)

### 🚜 Refactor

- Remove redundant DNS not found [`4119804`](https://github.com/AtifChy/aiub-companion/commit/41198046d07e5613ffafa6b5ede87d9323f8c031)

### 🧪 Testing

- Update config schema [`5adbf8b`](https://github.com/AtifChy/aiub-companion/commit/5adbf8ba8c92f99e083bee980c6ebc82718db407)
- *(window)* Add window [`36259b8`](https://github.com/AtifChy/aiub-companion/commit/36259b884d3f92608891243b7dc936989799a07e)

### ⚙️ Miscellaneous Tasks

- *(window)* Comment what SetRestoreWindow does [`7929598`](https://github.com/AtifChy/aiub-companion/commit/7929598a4e8ec67b7c1a1255a3de1ac4eef38bea)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.65...v0.1.76

## [0.1.65] - 2026-06-28

### 🚀 Features

- Add git-cliff configuration [`d12df9d`](https://github.com/AtifChy/aiub-companion/commit/d12df9d6ccbab301af7aedf84a8d694934cf92cf)
- *(app)* Add about page [`1f15dd8`](https://github.com/AtifChy/aiub-companion/commit/1f15dd8cd5a63670b2b744650ce5147615888cf1)
- Introduce window management package [`aa23548`](https://github.com/AtifChy/aiub-companion/commit/aa23548a09aa6e7deee183a6866bb82e6c78aa30)
- *(settings)* Add keep alive option [`aac79c9`](https://github.com/AtifChy/aiub-companion/commit/aac79c9b54a9d35f8846989278974d008f65c69d)
- *(error)* Introduce wailsError [`1148c9f`](https://github.com/AtifChy/aiub-companion/commit/1148c9f4f234771a25964ff25f558e105f77b454)
- *(ui)* Add focus event styling for window-controls [`b8be8c9`](https://github.com/AtifChy/aiub-companion/commit/b8be8c914a2ec730fcd258ac58111854ec61ef16)

### 🚜 Refactor

- *(log)* Extract path logic and expose `OpenLogFile` [`5f5dd18`](https://github.com/AtifChy/aiub-companion/commit/5f5dd1880afd20ea3be94ef973fd9f12403e1269)
- *(config)* Remove redundent window state [`3f13ce1`](https://github.com/AtifChy/aiub-companion/commit/3f13ce174ea0714dd53b66f54693f457e90fb312)
- Code tuning and cleanup [`0364af6`](https://github.com/AtifChy/aiub-companion/commit/0364af608cf3d0bf64727f21d3254ca3a53bf6a9)
- Wails already sends typed event constant [`d1b41e6`](https://github.com/AtifChy/aiub-companion/commit/d1b41e6b8e2d0a9c3f9215fd5786c83a854633ed)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.55...v0.1.65

## [0.1.55] - 2026-06-26

### 🚜 Refactor

- *(config)* Setup path using ServiceStartup [`03323de`](https://github.com/AtifChy/aiub-companion/commit/03323defa3e74d5686979fb909fdd0b7c7a6a706)
- Move config path assignment [`0a1226d`](https://github.com/AtifChy/aiub-companion/commit/0a1226df7d3bb59e965828f33e0272c7f1842557)
- *(database)* Make path a Service parameter [`03f6e84`](https://github.com/AtifChy/aiub-companion/commit/03f6e8488a8ff3fc375e68e13edd1c17fa642177)

### ⚙️ Miscellaneous Tasks

- Remove unused ci file [`b99d9da`](https://github.com/AtifChy/aiub-companion/commit/b99d9da4ebe05c31b4d29bd9b6272347c1852f01)
- Update release workflow artifact [`d7c33f4`](https://github.com/AtifChy/aiub-companion/commit/d7c33f40887b7339f8723bdb01bbe45b483e06e0)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.50...v0.1.55

## [0.1.50] - 2026-06-25

### 🚀 Features

- Centralize event constants [`3c93b44`](https://github.com/AtifChy/aiub-companion/commit/3c93b449a2cac8aa8c595818e7fe52c79816859d)

### 💼 Other

- Update wails dependencies [`a20fdc7`](https://github.com/AtifChy/aiub-companion/commit/a20fdc75bc6a91464e7eea863d5a17238f38da4f)

### 📚 Documentation

- Update project readme [`0eadff1`](https://github.com/AtifChy/aiub-companion/commit/0eadff1432be2faf3872310e009302372bfbe123)

### ⚙️ Miscellaneous Tasks

- Add github actions release workflow [`7271886`](https://github.com/AtifChy/aiub-companion/commit/7271886516c3d4111fe015343328871efc02996b)
- Add github release workflow configuration [`353088d`](https://github.com/AtifChy/aiub-companion/commit/353088d42217543f7302e4499b5bb52fca96af9e)
- Add changelog github action [`cb7500b`](https://github.com/AtifChy/aiub-companion/commit/cb7500b10ae02dd9fa99d3467dd3f57026b86a37)
- Remove changelog workflow [`c63ec4a`](https://github.com/AtifChy/aiub-companion/commit/c63ec4ad1df1a8531151a12fe6a91e139f47a9fd)
- Update gtk and webkitgtk [`da9611c`](https://github.com/AtifChy/aiub-companion/commit/da9611ce5d9516f4dce34d2419c0fc27fa84f0d2)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.42...v0.1.50

## [0.1.42] - 2026-06-25

### 🚀 Features

- Emit settings changed event [`20ef7e7`](https://github.com/AtifChy/aiub-companion/commit/20ef7e75ecf6f54b14da8dbeafa3130b1ba35fa8)
- *(log)* Integrate custom logger and startup service [`00a4bbe`](https://github.com/AtifChy/aiub-companion/commit/00a4bbe55a227e77aed2fa6ef0fc20c42ad9cabf)
- *(worker)* Add worker [`5e017ce`](https://github.com/AtifChy/aiub-companion/commit/5e017ced4b8f9be8dd4055d4ad41d99998094ed8)
- Adjust settings save logic on startup [`304636f`](https://github.com/AtifChy/aiub-companion/commit/304636f9e977e368b778b3457c110da4d9ed3f42)
- Add `BuildTime` build variable [`4353066`](https://github.com/AtifChy/aiub-companion/commit/4353066855858c28d632773d5dbad703229476dd)

### 🐛 Bug Fixes

- Improve validation error output [`c9362ed`](https://github.com/AtifChy/aiub-companion/commit/c9362ed3b9ea4bef9b0a2c7e58571fa6a59eead2)
- *(layout)* Double click maximize [`aaf8004`](https://github.com/AtifChy/aiub-companion/commit/aaf80048fc8d3363ab7e0a8732a3625ef50e3119)
- Production log use config instead [`f28d6b5`](https://github.com/AtifChy/aiub-companion/commit/f28d6b51440b63ac181152acb0bdfb87da212752)

### 💼 Other

- Update build time variables [`096ea36`](https://github.com/AtifChy/aiub-companion/commit/096ea36ee8130a08ee4f859a1b2c7d13b1bf839a)
- Update ldflags for build info [`37cc5d4`](https://github.com/AtifChy/aiub-companion/commit/37cc5d47345493225d04234ae0a074ec89dd772b)

### 🚜 Refactor

- Use proper name for app [`be5f077`](https://github.com/AtifChy/aiub-companion/commit/be5f077dac586a2809686e721eed28cbd16b86d2)
- Move ui management to desktop service [`29952b8`](https://github.com/AtifChy/aiub-companion/commit/29952b8e12df3d29badf27171933b98b762f35a3)
- Convert database init into a service [`469d743`](https://github.com/AtifChy/aiub-companion/commit/469d74344fb43f9941ed373896119a91742df7d9)
- Initialize service components on startup [`d79b25c`](https://github.com/AtifChy/aiub-companion/commit/d79b25c27ce15d692204481708240aa72377a532)
- Restructure main application initialization [`62625d7`](https://github.com/AtifChy/aiub-companion/commit/62625d73003dec90728bfdfd6a95600dfa5b9a99)
- Update type-checked eslint [`bc3d790`](https://github.com/AtifChy/aiub-companion/commit/bc3d7904b3d294381638edf83b45568ef32ec4bb)
- Move embed schema variable [`ddbd2c6`](https://github.com/AtifChy/aiub-companion/commit/ddbd2c63ca82f53f27f4102e9ce29c0d0807150e)
- Remove build time variables from main [`a019dba`](https://github.com/AtifChy/aiub-companion/commit/a019dbac5e20575f6436a832366b240be9a77862)
- Rename settings package [`62c3c71`](https://github.com/AtifChy/aiub-companion/commit/62c3c712fa878c5c1edcb1a46ede9857591eded5)
- Rename settings to config [`87baf4a`](https://github.com/AtifChy/aiub-companion/commit/87baf4a383315165beedd1d2ace4f7dbc8f1667b)
- Rename meta service and move logic to config [`b64b630`](https://github.com/AtifChy/aiub-companion/commit/b64b630ef8a00794d76d5f83d292bfd2f957b6bd)
- Another batch of renaming settings to config [`bdf6a00`](https://github.com/AtifChy/aiub-companion/commit/bdf6a00a5c064b5ae4959af549f5fa06ef7be1b0)

### 🧪 Testing

- Update settings schema [`78bb19b`](https://github.com/AtifChy/aiub-companion/commit/78bb19b76f04fbfa2f68eddd36c01387658a859f)

### ⚙️ Miscellaneous Tasks

- Code cleanup [`974616b`](https://github.com/AtifChy/aiub-companion/commit/974616bfbc4da536a5eb2e302e5e30c199786473)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.28...v0.1.42

## [0.1.28] - 2026-06-23

### 🚀 Features

- *(settings)* Add settings provider component [`8d20f58`](https://github.com/AtifChy/aiub-companion/commit/8d20f58005ea9a6b028d159f0697e33babf76113)
- *(log)* Add frontend logger service [`213b1be`](https://github.com/AtifChy/aiub-companion/commit/213b1be5b5a1952a3e3fd9e70519e7a0a7ffb8b2)
- *(settings)* Add sidebar config option [`1ecd604`](https://github.com/AtifChy/aiub-companion/commit/1ecd604442a903786e89ab09559fe0e2dd46e1f7)

### 🚜 Refactor

- Update service binding imports [`0b89462`](https://github.com/AtifChy/aiub-companion/commit/0b8946229ae032c13304abf79ad4f5278e661f92)
- *(layout)* Remove async from `handleResize` [`c2150cb`](https://github.com/AtifChy/aiub-companion/commit/c2150cb8f521d56521593e709b2f7aa50d75d55b)
- *(settings)* Rename updateConfig to setConfig [`a8b7c50`](https://github.com/AtifChy/aiub-companion/commit/a8b7c50bae650e980db4b7834fa21b97fc46b200)
- Extract app routing into component [`d424f0f`](https://github.com/AtifChy/aiub-companion/commit/d424f0fa3c24d41c72bca015a72c8837f4aa0d0b)

### ⚙️ Miscellaneous Tasks

- Disable react-query retry [`c4af592`](https://github.com/AtifChy/aiub-companion/commit/c4af59211638bbf9a62e9a8e69dc63264dc9a7b6)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.20...v0.1.28

## [0.1.20] - 2026-06-23

### 🚀 Features

- *(settings)* Add window size restore [`f9ea75a`](https://github.com/AtifChy/aiub-companion/commit/f9ea75a7c31231002ad114820e8ede389dcf9fe1)
- *(hook)* Add useDelayedLoading [`f5a48d0`](https://github.com/AtifChy/aiub-companion/commit/f5a48d03089abaad94c1feadcd1adb4f6e1d9000)

### 🐛 Bug Fixes

- *(scroll)* Move wheel event listener [`1b7b3d8`](https://github.com/AtifChy/aiub-companion/commit/1b7b3d8cb310aaa617dc6dd1b941ce819af6de36)

### 💼 Other

- *(deps-dev)* Update typescript [`cf0bad9`](https://github.com/AtifChy/aiub-companion/commit/cf0bad9b016824aa7ecc8d79097fddfa577542de)

### 🚜 Refactor

- *(frontend)* Split up notices components [`a0a91fd`](https://github.com/AtifChy/aiub-companion/commit/a0a91fdffb16ee3eba8960560d613d8014098ad8)
- *(frontend)* Code organization and cleanup [`63cc890`](https://github.com/AtifChy/aiub-companion/commit/63cc89012d9daf1b9370f79869c46e4df80a5914)
- *(settings)* Streamline item [`1cf524c`](https://github.com/AtifChy/aiub-companion/commit/1cf524ccb28eab92fec97d7b7fdcd07f58581c06)
- Remove pagination from notice filters [`754213a`](https://github.com/AtifChy/aiub-companion/commit/754213ae40cb12210fb787edde75c516f0856f10)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.12...v0.1.20

## [0.1.12] - 2026-06-22

### 🚀 Features

- Improve fts query unicode support [`9f328bc`](https://github.com/AtifChy/aiub-companion/commit/9f328bce81787f9b02218bffac2060090a13c664)
- Implement fetch with retriable [`05f6092`](https://github.com/AtifChy/aiub-companion/commit/05f6092c7effe9a4a7d923c76b24445de5655328)
- Add Help page with FAQ and resources [`6a10021`](https://github.com/AtifChy/aiub-companion/commit/6a100212a98bb2674291eb139e1f8e1d94bfa094)
- Add eslint and update dependencies [`7721cf4`](https://github.com/AtifChy/aiub-companion/commit/7721cf48576d9d6abb6b0ef955439d5ed08d8ded)

### 💼 Other

- Update go dependencies and refactor [`abba6e0`](https://github.com/AtifChy/aiub-companion/commit/abba6e0ad6cf46c93863786918d47261c904d761)

### 🚜 Refactor

- Reorder fields in Filter struct [`f52f4be`](https://github.com/AtifChy/aiub-companion/commit/f52f4be727eb46faba11688fd04c0553dfa7edee)

### 🎨 Styling

- Adjust page animation classes [`2a9ab9c`](https://github.com/AtifChy/aiub-companion/commit/2a9ab9c92afad71fe7d552f32652aa37a53d5367)

### 🧪 Testing

- Add test cases for routine, settings, fts [`178b5b7`](https://github.com/AtifChy/aiub-companion/commit/178b5b75fb60c2ac9a7819b8d82177ccc3844626)

### ⚙️ Miscellaneous Tasks

- Cleanup [`7def93c`](https://github.com/AtifChy/aiub-companion/commit/7def93c81e6f1096b2db38ffdbb1aeda9ea43b5f)
- Update go dependencies [`aaf9295`](https://github.com/AtifChy/aiub-companion/commit/aaf929556b83ba8ca7ae1bdbcba7b5feb2ccd008)
- Dyanmically set app version [`f799832`](https://github.com/AtifChy/aiub-companion/commit/f799832aef4bf0374924d5a8df0c0218bb492ad3)
- Update taskfiles [`5eb8088`](https://github.com/AtifChy/aiub-companion/commit/5eb80882e7f1b5b28a141105c15555dae4dbca1c)


> **Full Changelog**: https://github.com/AtifChy/aiub-companion/compare/v0.1.0...v0.1.12

## [0.1.0] - 2026-06-16


### 👥 New Contributors
* @AtifChy made their first contribution

