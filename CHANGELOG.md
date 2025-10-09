# Changelog 

## [v1.0.0] - 09-10-2025
### Breaking Changes
- All PowerShell scripts (`join`, `unjoin`, and `hostname`) are now **embedded** into the `machina.exe` binary.
- The external `scripts/` folder is no longer required for runtime.
- Commands that previously loaded `.ps1` files from disk now executes from temporary embedded copies.

### Changed
- `change-hostname.ps1`, `join-ad.ps1`, and `unjoin_to_workgroup.ps1` now run entirely from embedded resources.
- CLI commands now call `run.PS1Embedded()` internally instead of reading from disk.

### Fixed
- Removed dependecy on local script paths when running `machina.exe` directly.
- Eliminated `0xfff0000` errors caused by missing script files when executed standalone.

## [v0.3.1] - 09-10-2025
### Added
- `--registry` flag for `hostname` command to force hostname changes by directly editing registry keys.
- `--yes` flag to bypass confirmation when using `--registry`.
- Cobra-based confirmation prompt before applying forced hostname changes.

### Changed
- Removed in-script prompting from `change-hostname.ps1`; all confirmation is handled in Cobra.

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