# Eine Meldeschwelle der Größe steht in keinem Artefakt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Neben der Rückführungs-Marke in §4 eines Plans gilt eine Meldeschwelle darunter (der Implementer meldet,
bevor er sie überschreitet). Sie steht in keinem Artefakt des Repos (Plan, `AGENTS.md`, `.claude/commands`,
`harness/`); ob sie ohne Anhalten überschritten wurde, ist gegen keine Regel prüfbar (Review F-7).

Abgrenzung zu `BEO-REPO/groessenprobe-ohne-messbaren-stand`: Dort nennt die Marke ihre Metrik nicht. Hier
ist die Marke messbar, die Schwelle davor fehlt als Artefakt.

Erweitert am 2026-10-16 (Entscheidung des Nutzers): Die Beobachtung gilt auch für die Entscheidung nach der
Meldung. Meldet der Implementer das Erreichen der Meldeschwelle und entscheidet der Planner „weiter“ oder
„schneiden“, steht die Entscheidung nur als Satz im Bericht, ohne datierte Zeile in einem Artefakt des
Plans (Review F-5, Verifikation V-3 in `slice-v1-abschluss-verbindungsende-replay`).
