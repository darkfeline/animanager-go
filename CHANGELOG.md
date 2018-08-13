# Changelog

## 0.14.0

### Changed

* Schema changes (version 4):
  * `anime.title` is no longer unique.
  * Added `episode_file` table.
  * Dropped `episode_type` table.
  * Dropped `file_priority` table.
  * Dropped `cache_anime` table if it exists.
