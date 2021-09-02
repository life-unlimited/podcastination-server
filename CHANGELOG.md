## [1.0.4](https://github.com/life-unlimited/podcastination-server/compare/v1.0.3...v1.0.4) (2021-09-02)


### Bug Fixes

* fixed missing podcast xml details ([34aa30d](https://github.com/life-unlimited/podcastination-server/commit/34aa30dbc54d30155f1f7c583e215747b497a2b6))

## [1.0.3](https://gitlab.com/life-unlimited/podcastination-server/compare/v1.0.2...v1.0.3) (2021-05-12)


### Bug Fixes

* fixed unwanted image and pdf path for episodes in import job ([709b3fc](https://gitlab.com/life-unlimited/podcastination-server/commit/709b3fc297e868d8535340f97680124fd6552b02))

## [1.0.2](https://gitlab.com/life-unlimited/podcastination-server/compare/v1.0.1...v1.0.2) (2021-05-12)


### Bug Fixes

* fixed broken podcast xml generation ([b8e00d6](https://gitlab.com/life-unlimited/podcastination-server/commit/b8e00d61d8ea405ef6cf4df3e20c63a30851d605))
* fixed broken season retrieval by key ([46bd65e](https://gitlab.com/life-unlimited/podcastination-server/commit/46bd65ee8f2273169c9a50d991486963a1ccf140))

## [1.0.1](https://gitlab.com/life-unlimited/podcastination-server/compare/v1.0.0...v1.0.1) (2021-05-11)


### Bug Fixes

* fixed broken pdf file extension check for import job ([cd0c933](https://gitlab.com/life-unlimited/podcastination-server/commit/cd0c933cd3b341478ade3704441abf4871f91a47))

# 1.0.0 (2021-05-11)


### Bug Fixes

* added image check before trying to move ([77cd52a](https://gitlab.com/life-unlimited/podcastination-server/commit/77cd52a88745a20c307a9b648197a9b87ad25091))
* added missing error message for config file close ([f17fecc](https://gitlab.com/life-unlimited/podcastination-server/commit/f17fecca597944384abebea770ed1bb73ff9d3f5))
* added missing get methods by key for podcast and season store ([3246f04](https://gitlab.com/life-unlimited/podcastination-server/commit/3246f04022bbad3569b1e10b593b3321bdde0d73))
* added missing podcast id constraint for getting seasons by key ([8ca4e09](https://gitlab.com/life-unlimited/podcastination-server/commit/8ca4e090853d45407b37f9c44f7a48ac82a36e7d))
* added missing sort for seasons and episodes ([fe1c6fa](https://gitlab.com/life-unlimited/podcastination-server/commit/fe1c6fa8d0f373cd286c878f2af44b61d3754098))
* fixed broken episode parsing for db query ([93bd046](https://gitlab.com/life-unlimited/podcastination-server/commit/93bd046cf92c666b77c5a2581781d654caaac546))
* fixed faulty image importing ([b9a1535](https://gitlab.com/life-unlimited/podcastination-server/commit/b9a15353ff11bd70fcdedc1454270395f1feb8e1))
* fixed malformed length in podcast xml ([da6c74b](https://gitlab.com/life-unlimited/podcastination-server/commit/da6c74b3df3ae6e440a082bc03c7ed4a724172d0))
* fixed move errors for inter volume operations ([dd63e05](https://gitlab.com/life-unlimited/podcastination-server/commit/dd63e05d702de7d33d266b73fa6097298d8fa616))
* fixed null errors for episode store ([162ae63](https://gitlab.com/life-unlimited/podcastination-server/commit/162ae63704fff5f3045fdaac50ddabb9a4efb5b7))
* fixed null errors for owner store ([7945d3d](https://gitlab.com/life-unlimited/podcastination-server/commit/7945d3de03783075f5846d32011fac4ce9a9282f))
* fixed null errors for podcast store ([0f8a8f0](https://gitlab.com/life-unlimited/podcastination-server/commit/0f8a8f0f310818ea59441615008680b22e925fdc))
* fixed null errors for season store ([52aa94f](https://gitlab.com/life-unlimited/podcastination-server/commit/52aa94f14349f10955fdce384f6d967620d2a35b))
* fixed unwanted caching of api responses ([b550389](https://gitlab.com/life-unlimited/podcastination-server/commit/b550389bf4ea4a0f80f9a6609f1cb1bd4f20c4d3))
* fixed unwanted characters in filenames ([55e5bfd](https://gitlab.com/life-unlimited/podcastination-server/commit/55e5bfde602d7776128e3ba81e1722f65e9fac03))
* fixed wrong date sort order for multiple import tasks ([149895e](https://gitlab.com/life-unlimited/podcastination-server/commit/149895e99d254106136d20131c9828ae5714400f))
* fixed wrong json field for subtitle in import task details ([2f43f4e](https://gitlab.com/life-unlimited/podcastination-server/commit/2f43f4ea8c165a386cbe1aa0ae347ae8d90593ec))
* fixed wrong return type for owner and podcast store methods ([d55e149](https://gitlab.com/life-unlimited/podcastination-server/commit/d55e14907146ec192d55855d0d90917555b444fd))
* fixed wrong wildcard in sql queries ([1443e96](https://gitlab.com/life-unlimited/podcastination-server/commit/1443e961743f775230f59af57836ed1ad52438d9))
* minor fixes for stores ([fab1c5c](https://gitlab.com/life-unlimited/podcastination-server/commit/fab1c5cb9f06e7dd30721e10b87a1684571c4d0b))
* removed unwanted fields in episode model ([97cf583](https://gitlab.com/life-unlimited/podcastination-server/commit/97cf58340ffa3c79ddad3e17f4e182941956c018))


### Features

* add support for pdf attachments ([903c667](https://gitlab.com/life-unlimited/podcastination-server/commit/903c6679eb248f5c7df70ac75965efe30b0d1a98))
* added /podcasts endpoint ([e492a0e](https://gitlab.com/life-unlimited/podcastination-server/commit/e492a0e28f7c1845e7b0864952d8b278ea6d0c8b))
* added config reading ([85d9b83](https://gitlab.com/life-unlimited/podcastination-server/commit/85d9b838e1b06c8fb75383b575f42585b6aa00ed))
* added creation and update support for episode store ([0d3fdc0](https://gitlab.com/life-unlimited/podcastination-server/commit/0d3fdc019dd51468a969628f29c40230c1c0a89a))
* added endpoint for getting podcast by id ([2b84930](https://gitlab.com/life-unlimited/podcastination-server/commit/2b849305dbfb07c66fa68f347fccd349033fd591))
* added error message for podcast xml write fail ([5b59560](https://gitlab.com/life-unlimited/podcastination-server/commit/5b59560ab38c06728337f35689b401d94dd934b7))
* added feed link attribute to podcast model ([e53b36b](https://gitlab.com/life-unlimited/podcastination-server/commit/e53b36b4fae5f210267b2afb943fe676aa529491))
* added image location to episodes ([ecbfd44](https://gitlab.com/life-unlimited/podcastination-server/commit/ecbfd44da544ba40b1e853ab35a844b0695ba1ed))
* added initial app ([49b219e](https://gitlab.com/life-unlimited/podcastination-server/commit/49b219eb61c2a5e6091f48deda386cf0c93afb1c))
* added isAvailable field to episode model and store ([7c795b3](https://gitlab.com/life-unlimited/podcastination-server/commit/7c795b33f0f97deb917d58a99d41bb47fb5957dc))
* added job scheduling ([cf904b9](https://gitlab.com/life-unlimited/podcastination-server/commit/cf904b92790319b51d9031a5132bd2c0dc4f48c5))
* added key to podcasts ([cc62f25](https://gitlab.com/life-unlimited/podcastination-server/commit/cc62f25c2bcf003c9351cd89f52a86d0e53d6115))
* added key to seasons ([bbdb8e1](https://gitlab.com/life-unlimited/podcastination-server/commit/bbdb8e10040c9f658196860c4a8307fc71567bc4))
* added logs for internal errors by rest calls ([57f5b62](https://gitlab.com/life-unlimited/podcastination-server/commit/57f5b621dfb013fa89a490b14cd7c589f2494e37))
* added mp3 length attribute to podcast model ([ee8f4ee](https://gitlab.com/life-unlimited/podcastination-server/commit/ee8f4ee8a22a6d9edcccf5cdc7ef3827a2c4295f))
* added podcast and episode models ([9118f52](https://gitlab.com/life-unlimited/podcastination-server/commit/9118f528d16ed13f43f3347c11c96ea02eb2ec73))
* added podcast xml model and generation ([4047a6d](https://gitlab.com/life-unlimited/podcastination-server/commit/4047a6de71e87e048a964e6f3a1c3beb341ef50b))
* added static file server ([1b5bd1f](https://gitlab.com/life-unlimited/podcastination-server/commit/1b5bd1f9018cc2ded78f6609b28e92fd833399ab))
* added stores and some minor changes ([3e6462a](https://gitlab.com/life-unlimited/podcastination-server/commit/3e6462a8fa311d24843523ca223e6acb2fac0ca6))
* added yt url to episodes ([48a5a2a](https://gitlab.com/life-unlimited/podcastination-server/commit/48a5a2a08d392c3c59e9d3d184c05fb36b693388))
* created all models ([2abf9e4](https://gitlab.com/life-unlimited/podcastination-server/commit/2abf9e42c81a179f76733c9533fcb1d42939b3a7))
* finished first iteration ([a1edc75](https://gitlab.com/life-unlimited/podcastination-server/commit/a1edc75fb735d984a5272b00c199bab2df54ea49))
* finished first try ([686ba19](https://gitlab.com/life-unlimited/podcastination-server/commit/686ba19317b9c29a71e87a6292f4d035908673ae))
* finished import job for episodes ([bbe0629](https://gitlab.com/life-unlimited/podcastination-server/commit/bbe0629e369cf0b9635d17fb41f1e3e0f8183bb3))
* interval for import now taken from config ([b94a84b](https://gitlab.com/life-unlimited/podcastination-server/commit/b94a84be7d11403cd9ca05e5d98e379042512e31))
* now bootable ([7d9ae73](https://gitlab.com/life-unlimited/podcastination-server/commit/7d9ae7346f479eab5eca843a59cf2b44961ce26c))
