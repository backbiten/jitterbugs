# IA/AI development (assistive, internal-only)

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

Push notifications must not be triggered by detection/classification of people, emotions, intent, or suspiciousness.

Allowed trigger categories:
- System health (e.g., service errors, resource exhaustion)
- Security (e.g., unauthorized access attempts, certificate expiry)
- Reliability (e.g., crash reports, performance degradation)
- Privacy/consent (e.g., data retention expiry, consent revocation events)
