# Changelog 

## [v0.3.0] - 07-10-2025
### Added 
- Logging and error handling in `unjoin` script (`unjoin_to_workgroup_wmi.ps1`).
- `--restart` support for `unjoin` command to reboot automatically after leaving a domain.

### Changed
- No changes made in this version

### Fixed
- Exit codes now properly reflect success/failure in `unjoin` operations (no longer always return 0).

## [v0.2.0] - 07-10-2025
### Added
- `hostname` command to change Windows hostname.

### Changed 
- Improved join script with logging and error checking.

### Fixed
- Removed WMIC dependency for `unjoin`.

## [v0.1.0] - 06-10-2025
### Added
- Initial release with `join` ans `unjoin` commands.