**Stand:** verkörpert

Zielort: `AGENTS.md §3.12`, `.claude/commands/plan-welle.md` Schritt 6, `.claude/agents/architect.md` und `.claude/commands/implement-slice.md` Schritte 12 und 13 und Randform-Rückgabe — Herkunfts-Anker `seit slice-harness-randformen-vor-code`.

Der Absatz *Aufruf-Voraussetzung gehört dazu* in `AGENTS.md §3.12` (Ordnung der Werte über mehrere Aufrufe einer Session: §6 nennt nacheinander oder gleichzeitig; ein Gegenbeispiel im Bericht ist keine Rückgabe) — Herkunfts-Anker `seit slice-v1-abschluss-zeitangaben-aufzeichnen` (Beleg `evidence/slice-v1-abschluss-zeitangaben-aufzeichnen.md`, F-647).

Die Ordnung zweier Startfehler, die das Review fand (F-653), schließt der Absatz *Auch eine Lesart ist eine Entscheidung* in `.claude/commands/implement-slice.md` Randform-Rückgabe — Herkunfts-Anker `seit slice-v1-abschluss-zeitangaben` (Beleg `evidence/slice-v1-abschluss-zeitangaben.md`, F-653).

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.12): `slice-v1-abschluss-sessions` (F-679, F-681) — eine Randform einer Erkennung per Zustand und der Zeitpunkt des Bemerkens waren vor dem Code nicht genannt (Beleg `evidence/slice-v1-abschluss-sessions.md`).

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.12): `slice-harness-gate-kette-schnell` (F-694) — die Aufrufarten des Treibers (Make-Flags `-n`, `-q`, `-t`, `-i`, auch aus `MAKEFLAGS`) standen in §6 nicht; die Regel nannte sie nicht (Beleg `evidence/slice-harness-gate-kette-schnell.md`).

Der Absatz *Aufrufarten des Treibers gehören dazu* in `AGENTS.md §3.12` (ein Werkzeug, dessen Ergebnis ein Hook oder Gate später glaubt, nennt in §6 die Aufrufarten außerhalb des Normalfalls und je eine Entscheidung) — Herkunfts-Anker `seit slice-harness-gate-kette-schnell` (Beleg `evidence/slice-harness-gate-kette-schnell.md`, F-694). Als Vorschlag des Planners eingetragen; der Nutzer kann ihn streichen.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.12): `slice-v1-abschluss-tls-client` (F-1) — die Länge von `SSLRequest` und `GSSENCRequest` war im Code entschieden, in §6 und Spezifikation nicht genannt (Beleg `evidence/slice-v1-abschluss-tls-client.md`).

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.12): `slice-v1-abschluss-anmeldung` (F-3, V-2) — das Rennfenster der Rücknahme der Lesefrist stand weder in §6 noch in der Spezifikation und wurde erst nach dem Review als akzeptiertes Negativ entschieden (Beleg `evidence/slice-v1-abschluss-anmeldung.md`).

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.12): `slice-v1-abschluss-protokollrand` (V-1) — die Reihenfolge zweier Ablehnungen (`_pq_.` vor `PGR-E6003`) war im Code entschieden und stand in keiner Stelle der Spezifikation; die Verifikation fand sie (Beleg `evidence/slice-v1-abschluss-protokollrand.md`). Keine Zeile in `AGENTS.md` oder `.claude/commands/`; ob §3.12 die Reihenfolge mehrerer Ablehnungsgründe an einem Rand als Randform nennen soll, entscheidet der Nutzer.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.12): `slice-v1-abschluss-verbindungsende-format` (F-1) — zwei Lesarten zur Stelle von `ready_for_query` und von `connection_closed` in der Nachrichtenfolge standen weder in §6 noch in der Spezifikation; der Implementer entschied sie im Code (Beleg `evidence/slice-v1-abschluss-verbindungsende-format.md`). Retirement-Check von §3.12 (Urteil des Planners, 2026-10-15): die Beobachtung ist seit `slice-harness-randformen-vor-code` wieder aufgetreten; die Regel wird nicht gelockert. Keine Zeile in `AGENTS.md` oder `.claude/`; ob §3.12 die Stelle einer Nachricht in der Nachrichtenfolge ausdrücklich als Randform nennen soll (Tabelle Nachricht × Stelle: Anfang, Mitte, direkt davor, direkt danach, letzte Gruppe, nicht letzte Gruppe), entscheidet der Nutzer.
