# Changelog

## 0.21.0

## Changed

* Default database path has been moved to
  `$XDG_STATE_HOME/animanager/database.db`.

## 0.20.0 - 2023-11-23

WARNING: If you're using a older version, you **must** update to this
version and run `add -no-eid` to populate the new `eid` field, which
will replace the old episode `id` field in future versions.  You will
need to run this once every 24 hours until it fills in all the fields
(there is a limit per run to not get banned).

### Added

* Added `-no-eid` flag to `add` command.
* Added `eid` field.

## Changed

* `add` command no longer prints AIDs to stdout since it already
  prints the AIDs to stderr logs.

## 0.19.0 - 2021-04-01

### Changed

* Changed default config file path to respect `XDG_CONFIG_HOME`.
* Changed default database path to respect `XDG_DATA_HOME`.

## 0.18.0 - 2020-01-02

### Changed

* End date is now considered by `add -incomplete`.
* The `unregister` command `-watched` flag is renamed to `-finished`.
* Changed `stats` output a bit.

### Fixed

* The `unregister` command `-finished` flag (renamed from `-watched`)
  actually works now.

## 0.17.0 - 2019-10-06

### Added

* Added `-watched` to `unregister` command.
* Added `update-titles` command.

### Changed

* `add` and `unregister` now print the affected AIDs to stdout.
* `add` now deletes episodes that were removed from AniDB.
* `register` now checks if the pattern is valid.

### Removed

* `-skipcache` option for `search` command.

## 0.16.0 - 2018-12-02

### Added

* Added `-missing` flag to `watchable` command.
* Added `stats` command.
* Added `unfinished` command.

### Changed

* `watchable` only prints the next episode instead of the next three.
* `watchable` now skips credit and trailer episode types (OP/ED/PV).
* `titlesearch` command renamed to `search`.
* `unregister` now accepts multiple aids.

## 0.15.0 - 2018-09-16

### Added

* `add` command now waits two seconds between queries per AniDB API
  requirements.
* New `-incomplete` flag for `add` command.

### Fixed

* `watch` command handles invalid input now.

## 0.14.0 - 2018-09-02

### Added

* New database schema (version 5):
  * Added `offset` to `watching` table.  This is used to adjust for
    different episode numbers in filenames compared to AniDB.
* Database is backed up before migrating.
* `register` command now has an `-offset` flag for setting file
  pattern episode offset.

### Fixed

* `watch` `-episode` flag now works.

### Changed

* `player` configuration option is now a list.
* File pattern now matches only against the base filename.
* `showfiles` now works on AIDs by default.
* `watch` doesn't prompt to mark done if already done.

## 0.13.0 - 2018-08-20

This is the first version of the Go version of Animanager.

### Changed

* Schema changes (version 4):
  * `anime.title` is no longer unique.
  * Added `episode_file` table.
  * Dropped `episode_type` table.
  * Dropped `file_priority` table.
  * Dropped `cache_anime` table if it exists.
* Animanager UI has been greatly changed:
  * Animanager now presents a command interface instead of a CLI.
  * Animanager is now aware of individual episodes, instead of simply
    tracking anime episodes by count.
  * Animanager now has a separate command `findfiles` for associating
    files with episodes.
  * File matching patterns now use Go regular expressions.
  * File matching patterns now use the first capturing group as the
    episode number, rather than the group named `ep`.

## 0.12.0 - 2018-08-20

See [PyPI Animanager](https://pypi.org/project/animanager/) for
previous versions, which are for the Python version of Animanager.
