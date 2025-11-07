# Requirement: AUTH (Sessions & Login)

## Ziel
Sichere Authentifizierung ohne externe JS-Toolchains, SSR-only.

## Funktional
- Login/Logout, Session mit HttpOnly/SameSite
- Rate-Limit pro IP & User
- CSRF-Token für state-changing Ops

## Prozess & KPIs
- 95% Logins < 300ms lokal
- Fehlversuche/Std, Lockouts/Std, aktive Sessions

## Risiken & Mitigation
- Brute-Force → Rate-Limit + Lockout
- Session Theft → HttpOnly + SameSite + Rotation
- Audit-Pflicht → JSONL Append-only

## Nicht-funktional
- DSGVO: Minimaldaten, Pseudonymisierung im High-Privacy-Modus
- WCAG 2.1 AA (Form-Labels, Fokus, Tastaturbedienung)
