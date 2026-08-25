# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	townbasketball/cmd/server	[no test files]
ok  	townbasketball/internal/auth	0.014s
?   	townbasketball/internal/config	[no test files]
ok  	townbasketball/internal/domain	0.002s
ok  	townbasketball/internal/guestbook	0.014s
ok  	townbasketball/internal/httpapi	0.028s
--- FAIL: TestScoreCorrectionKeepsAudit (0.01s)
    score_test.go:31: published response retained old score: {ID:game-score HomeTeamID:team-a AwayTeamID:team-b ScheduledAt:0001-01-01 00:00:00 +0000 UTC Venue:体育馆 Status:scheduled HomeScore:0 AwayScore:0 Published:true}
FAIL
FAIL	townbasketball/internal/league	0.025s
ok  	townbasketball/internal/media	0.010s
?   	townbasketball/internal/reporting	[no test files]
ok  	townbasketball/internal/store	0.015s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/server): exit `0`
