# IA/AI development (assistive, internal-only)

We may add IA/AI helper programs to provide additional internal tools for maintainers (debugging, testing, documentation, performance monitoring, and uncertainty calibration), while remaining strictly within the Hard Constraints below.

- developer assistance (summaries, documentation, tests, review)
- reliability/safety tooling (crash triage, performance regressions, log anomaly detection on non-personal telemetry)
- explainability/uncertainty

## Hard Constraints
- no screening/threat scoring
- no identity/intent/deception/suspiciousness inference
- no law-enforcement integrations
- no use on public/semi-public surveillance feeds
- consent-first
- local-first
- data minimization

## Notifications policy (safe-only)

Push notifications MUST NOT be triggered by detection/classification of people, emotions, intent, or “suspiciousness”.
Notifications may only be triggered by system/security/reliability signals, such as:

- System health: camera offline, packet loss, storage full, overheating, time drift, low battery.
- Security: failed login attempts, suspicious admin/API key activity, integrity check failures.
- Reliability: crash loops, queue backlog, dropped frames, service/model errors (without analyzing humans).
- Privacy/consent: capture disabled until consent is confirmed; retention limit reached; telemetry opt-in/out changes.
