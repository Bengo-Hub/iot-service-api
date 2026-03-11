# Security Policy

The IoT Service holds sensitive device and telemetry data. Please follow these guidelines to keep it secure.

## Supported Versions

| Version | Supported |
|---------|-----------|
| `main` branch | ✅ |
| Tagged releases (future) | ✅ |
| Older branches/forks | ❌ |

## Reporting Vulnerabilities

1. Email `security@bengobox.com` with a detailed description (do not open a public issue).
2. Include reproduction steps, impact assessment, and suggested mitigations if possible.
3. Encrypt communications when feasible; public PGP keys are available on request.

We will acknowledge reports within 48 hours and coordinate remediation and disclosure.

## Secure Development Practices

- Never commit credentials or production data.
- Use parameterised queries; Ent handles this by default.
- Validate input rigorously, especially for device provisioning endpoints.
- Enforce least privilege in database roles and service accounts.
- Run `govulncheck` / dependency scanners regularly.
- Always validate JWT tokens via auth-service JWKS endpoint.
- Secure device authentication with certificates and proper key management.

## Infrastructure Considerations

- Enable point-in-time recovery for PostgreSQL.
- Protect message brokers (NATS/Kafka) with TLS and authentication.
- Monitor audit logs for suspicious device access or unauthorized operations.
- Enforce RBAC at both API and database levels.
- Implement device certificate rotation and revocation.

Thank you for helping keep the BengoBox platform secure.

