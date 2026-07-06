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
## [0.1.55] - 2026-06-26

### 🚜 Refactor

- *(config)* Setup path using ServiceStartup [`03323de`](https://github.com/AtifChy/aiub-companion/commit/03323defa3e74d5686979fb909fdd0b7c7a6a706)
- Move config path assignment [`0a1226d`](https://github.com/AtifChy/aiub-companion/commit/0a1226df7d3bb59e965828f33e0272c7f1842557)
- *(database)* Make path a Service parameter [`03f6e84`](https://github.com/AtifChy/aiub-companion/commit/03f6e8488a8ff3fc375e68e13edd1c17fa642177)

### ⚙️ Miscellaneous Tasks

- Remove unused ci file [`b99d9da`](https://github.com/AtifChy/aiub-companion/commit/b99d9da4ebe05c31b4d29bd9b6272347c1852f01)
- Update release workflow artifact [`d7c33f4`](https://github.com/AtifChy/aiub-companion/commit/d7c33f40887b7339f8723bdb01bbe45b483e06e0)
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
## [0.1.0] - 2026-06-16
