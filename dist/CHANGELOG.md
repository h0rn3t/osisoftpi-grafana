# Changelog

## 1.0.0

- Initial release.

## 2.0.0

- Move to React based framework.

## 3.1.0

- Added calculation to PI Points
- Added PI point configuration (thanks to @TheFern2)
- Added option to use last value from PiWebAPI
- Updated to Grafana plugin SDK v9.3.6

## 4.0.0

- Added a new dataframe label format. It can be disabled in the configuration page for backward compatibility
- Added engineering units to Dataframe field. This can be globaly disabled in the configuration page
- Optimized queries using PIWebAPI batch endpoint
- Improved raw query processing
- Added variable support in raw query
- Fixed annotations support
- Updated to Grafana plugin SDK v9.4.7
- Fixed PI AF calculation
- Added plugin screenshots

## 4.1.0

- Modified the PI Webapi controller endpoints used when calculation is selected
- Allow calculation when last value option is selected
- When calculation is selected, change label from Interpolated to Interval
- Fixed issue with variable in Element Path

## 4.2.0

- Fixed issue that only odd attributes were been shown
- Fixed issue when fetching afServerWebId

## 5.0.0

- Migrated backend to Go language
- Changed the query editor layout
- Support Grafana version 11
- Drop support for Grafana 8.x and 9.x

## 5.1.0

- Add units and description to new format - issue #154
- Fixed digital state - issue #159
- Fixed summary data - issue #160
- Fixed an error in recorded max number of points - issue #162
- Fix issue with summary when migrating from previous versions - issue $160

- Updated the query editor layout
- Added boundary type support in recorded values
- Recognize partial usage of variables in elements
- Added configuration to hide API errors in panel
- Truncate time from grafana date time picker to seconds
- Fixed warnings during deploy
- Fixed LICENSE file

### 5.3.1 — 2026-05-26

Patch release: hardening of WebSocket streaming introduced in 5.3.0. No breaking JSON schema changes; operators only need to replace the plugin binary/artifact.

#### Added

- Bounded WebSocket reconnect with exponential backoff (max 5 attempts, 1s–30s) when the PI read loop or subscriber channel closes
- `sweepStaleChannelConstructs()` to prevent unbounded growth of the streaming channel registry
- Unit tests for stream eligibility, disconnect propagation, `Dispose()` cleanup, and concurrent subscribe/register (`go test -race`)

#### Changed

- Datasource streaming is enabled by **Enable Streaming Support** only; experimental mode is no longer required
- Stricter `isStreamable()`: PI points only; excludes summary, interpolated, recorded values, expression, and AF attributes
- All `channelConstruct` access is mutex-protected; `SubscribeStream` accepts `ds/<uid>/<uuid>` paths
- Sender dispatch: non-blocking send, then up to 50ms blocking send, then drop with `dropped=true` in logs
- Renamed `pkg/plugin/steam.go` → `pkg/plugin/stream.go`

#### Fixed

- Subscribers no longer hang indefinitely after WebSocket read-loop failure (sender channels are closed on disconnect)
- `Dispose()` now closes WebSockets, sender channels, and clears streaming maps
- Disabling experimental features no longer turns off streaming in the datasource config UI

#### Migration

- Enable **Enable Streaming Support** in datasource settings (experimental toggle not needed)
- Dashboards that relied on streaming for summary/interpolated/recorded/expression queries fall back to HTTP automatically; use streaming only on compatible PI point queries

### 5.3.0

- Added PI Web API streamsets/channel WebSocket streaming support (PR #203)
- All tags in a query batch share a single persistent WebSocket connection
- Enable Streaming Support toggle in datasource config and per-query streaming in Query Editor
- Unit tests for streaming URL builder, fan-out, subscriber lifecycle, and frame conversion

### 5.2.0

- Improved query performance to PiWebAPI by joing all queries in Panel into one batch request only
- Change the Query Editor layout
- Increased WebID cache from 1 hour to 12 hours and made it configurable

- Added experimental feature to cache latest response in case of request failure to PiWebAPI