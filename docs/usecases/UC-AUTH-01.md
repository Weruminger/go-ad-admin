# UC-AUTH-01 Login

**Basic Flow**
1. Nutzer öffnet /login
2. Sendet Credentials
3. System legt Session an und leitet auf /

**Fehlerfälle**
- Falsches PW → 401, Counter++
- CSRF fehlt → 400
- Rate-Limit → 429

**Akzeptanz**
- Set-Cookie vorhanden (HttpOnly+SameSite)
- Audit-Entry geschrieben
