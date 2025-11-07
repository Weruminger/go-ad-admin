# Requirement: AD-USERS

## Ziel
CRUD für AD-Benutzer mit OU-Whitelist, LDAPS only in prod.

## Funktional
- Suche (Filter/Paging), Create, Disable, Reset PW
- Policy-Checks (Passwort, Pflichtfelder)
- Dry-Run & Audit-Preview

## Prozess & KPIs
- Suche ≤ 200 Treffer in <1s
- Fehlerquote „DN-Konflikt“ < 2%/Monat

## Risiken
- DN-Kollision, Timeout, Schema-Änderungen
