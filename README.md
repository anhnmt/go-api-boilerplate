# go-api-boilerplate

Golang API boilerplate

## Features

- Manage sessions & devices
- Access token & Refresh token (JWT)
- Revoke token & Revoke sessions
- Detect leaked token & block
- Encrypt payload with security route (RSA & AES)
- Recording trace telemetry (OpenTelemetry)
- ...

## Dependencies & Tools

- [Uber Fx](https://github.com/uber-go/fx)
- [Buf build](https://github.com/bufbuild/buf)
- [GORM gen](https://gorm.io/gen/index.html)
- [Fingerprint](https://github.com/anhnmt/go-fingerprint)
- [Vanguard](https://github.com/connectrpc/vanguard-go)
- [Protovalidate](https://github.com/bufbuild/protovalidate-go)
- [RSA Generator](https://www.csfieldguide.org.nz/en/interactives/rsa-key-generator/)