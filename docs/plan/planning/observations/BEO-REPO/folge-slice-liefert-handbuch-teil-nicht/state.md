**Stand:** offen

Teil geliefert, kein Auftreten (keine Datei in `evidence/`, der Zähler bleibt): `slice-v1-abschluss-einspielen-extended-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-einspielen-extended`; die Suche nach den ersetzten Sätzen findet nichts Verbliebenes.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-einspielen-anmeldung-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-einspielen-anmeldung`; die Suche nach den ersetzten Sätzen findet nichts Verbliebenes.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-einspielen-tls-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-einspielen-tls`; die Suche nach den ersetzten Sätzen (unverschlüsselt, `PGR-E2004`) findet nichts Verbliebenes außer `record`.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-sqlite-format-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen`; die Suche nach den ersetzten Grenz-Sätzen findet nichts Verbliebenes.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-zeitangaben-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben`; die Suche nach den ersetzten Sätzen findet nichts Verbliebenes.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-antwortvergleich-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-antwortvergleich` und `slice-v1-abschluss-umgebung-anderes-kommando`; die Suche nach den ersetzten Sätzen findet nichts Verbliebenes. Das aufgezeichnete Ende (`FATAL`/`PANIC` bei `--compare-responses`) steht als akzeptierte Lücke im Nehmer `slice-v1-abschluss-verbindungsende-doku` (Entscheidung des Nutzers vom 2026-10-12), kein Auftreten.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-sessions-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-sessions` und die Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` (Optionen und Variablen mit den Kommandos, bei denen sie gelten); V-310 war ein Bestandssatz, der zu weit gefasst war, kein nicht gelieferter Teil. Das aufgezeichnete Ende des Servers steht als akzeptierte Lücke im Nehmer `slice-v1-abschluss-verbindungsende-doku` (Entscheidung des Nutzers vom 2026-10-12), kein Auftreten.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-tls-client-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-tls-client` und die Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando`; die Lücke des Reviews (F-1) war eine Lücke des Plans, nicht der Lieferung.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-anmeldung-doku` — Verifikation bestätigt Handbuch- und README-Teil von `slice-v1-abschluss-anmeldung`; die Suche nach den ersetzten Sätzen (Datenbank ohne Passwort) findet nichts Verbliebenes. Die Lücke zur nicht lesbaren Nachricht (F-7, V-1) ist eine benannte Lücke und ein behobener zu weiter Satz, kein nicht gelieferter Teil.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-postgres-versionen-zusage` — die Verifikation bestätigt Handbuch- und README-Teil der Sendung aus `slice-v1-abschluss-anmeldung` und `slice-v1-abschluss-anmeldung-doku` (die Sätze „geprüft mit PostgreSQL 17“ zur Anmeldung sind je Verfahren nachgezogen).

Teil geliefert, kein Auftreten: `slice-v1-abschluss-protokollrand-doku` — Verifikation bestätigt Handbuch- und README-Teil der Sendung aus `slice-v1-abschluss-protokollrand` (alle fünf Inhalte gegen das gebaute Binary). V-1 (Klartext bei TLS-Pflicht stand nur für `first-request` im Satz der Zählung) fand die Verifikation vor der Closure, behoben in `ef8f5f2`; der Slice schloss nicht mit dem Fehlen. Der Fall, dass `replay` beim Beenden `PGR-W2001` meldet, wenn ein Client allein abbricht (P15 (a)), steht als eigener Eintrag `BEO-REPO/replay-meldet-w2001-beim-abbruch-des-clients` (1×, offen, ohne Adresse); hier kein Auftreten. V-1 (Klartext bei TLS-Pflicht stand nur für `first-request` im Satz der Zählung) war kein Auftreten (Entscheidung des Nutzers vom 2026-10-15), der Zähler bleibt 1×.

Teil geliefert, kein Auftreten: `slice-v1-abschluss-verbindungsende-doku` — Handbuch- und README-Teil der akzeptierten Lücke zu `FATAL`/`PANIC` bei `--compare-responses` ersetzt, am Binary von Verifikation bestätigt; kein Release lag in der Lücke (Entscheidung des Nutzers vom 2026-10-12). Der Zähler bleibt 1×.

Weiteres Auftreten, offen (Urteil des Planners, `slice-v1-abschluss-binary-namen`, 2026-10-20; Beleg `evidence/slice-v1-abschluss-binary-namen.md`): Der Abschnitt zu `v0.1.0` in Handbuch und `README.md` ist an `slice-v1-abschluss-release-v0-1-1` übergeben. Der Zähler steht bei 2, unter der Schwelle; `slice-v1-abschluss-release-v0-1-1` liefert den Teil erst nach dem Tag des Nutzers, bis dahin bleibt der Zustand offen.

Kein Auftreten, Gegenbeleg (`slice-erster-release-homebrew-nachweis`, Urteil des Planners, 2026-10-21; kein Zählerzuwachs, der Zähler bleibt bei 2): Der Handbuch-Abschnitt *Mit Homebrew installieren*, die Zeile in §1 *Voraussetzungen* und die Zeile in `README.md` liegen im Slice selbst, und jede Zeile ist gegen den Lauf belegt (die Verifikation hat sie gegen ihren eigenen Lauf und gegen die Formel im Tap geprüft). Die Zählung in §8 blieb bei zwei Schichten (Harness-Werkzeug, Nutzer-Doku), ein Doku-Folge-Slice `…-doku` war nicht nötig.

Geber geschlossen, Nehmer offen, kein Auftreten (`slice-play-feldvergleich`, Urteil des Planners, 2026-10-23): Handbuch und `README.md` sind unberührt; die Option `--diff-all` steht nur im Binary und in der Spezifikation. Der Nehmer `slice-play-feldvergleich-doku` liegt in `open/` und nimmt die Sendung in seinem §1 an; ein Release mit der Option entsteht erst nach ihm (`AGENTS.md` §3.11). Schließt er ohne Handbuchteil, ist es ein Auftreten.
