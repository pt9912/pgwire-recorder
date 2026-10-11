# Review-Report: slice-v1-abschluss-einspielen-tls-doku — 2026-10-11

**Review-Art:** Code-Review (Dokumentation). Geprüft wird gegen Plan §1, §3, §6 und §8, gegen `AGENTS.md` §3.9, §3.11 (inklusive Regel für Handbuch und README) und §3.13 und gegen das gebaute Binary.

**Gegenstand:** `0f9daf6..7b2663a`: Entscheidungen des Architect `15a789a` (Plan §6), Handbuch und README `eea9569`, Belege des Implementers `7b2663a` (Plan §7). Geändert sind `docs/user/benutzerhandbuch.md`, `README.md` und der Plan.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-v1-abschluss-einspielen-tls-doku` ganz gelesen am Stand `7b2663a`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `AGENTS.md` §3.9, §3.11, §3.13; `LH-FA-20`, `LH-FA-17`, `LH-QA-05`, `LH-RB-01` über den Kopf des Plans.
- `slice-doku-ist-stand` §7 (Befehlsfolge).
- Vorheriger Report am selben Modul: `slice-v1-abschluss-einspielen-tls` (bis F-608); die Nummern dieses Laufs beginnen bei F-609.

**Ausgeführte Läufe** (alles Docker-only, Scratch-Präfix `rev-tlsdoku-`, nichts im Repo geändert; Container, Netz und Scratch-Verzeichnis danach entfernt):

- **Binary.** `make build` am Stand `7b2663a` (Stufe `runtime`, `pgwire-recorder:dev`).
- **Server.** `postgres:17-alpine` in einem eigenen Netz (IPv4 und IPv6): `pgssl` (`ssl=on`, `hostnossl … reject`, Blatt von `ca.pem` mit `DNS:pgssl, DNS:pgpw, DNS:localhost` und drei IP-Adressen), `pgexp` (Blatt gültig 2020-01-01 bis 2020-01-02), `pgplain` (`ssl=off`), `pgpw` (`ssl=on`, `password`, `hostnossl … reject`), `pgpwplain` (`ssl=off`, `password`), `pgchain` (Blatt von einem Zwischenzertifikat, Server sendet beide); Python-Fake mit `N`, `Z`, Ende ohne Antwort und `S` (Schweigen: gestartet, nicht ausgewertet). Zertifikate mit `openssl` 3.5 in einem alpine-Hilfscontainer.
- **Probe-Matrix (1) bis (15).** Alle Zeilen gefahren und mit den Aussagen von Handbuch und README verglichen: (1) `sslmode=require` plus `--upstream-ca` Exit 0, ohne CA Exit 4 `PGR-E4005`; (2) `--upstream-tls` mit `host:port` Exit 0; (3) IPv4 und IPv6 im Blatt Exit 0; (4) unbekannte CA Exit 4 `PGR-E4005`; (5) Name nicht im Blatt Exit 4 `PGR-E4005`; (6) abgelaufen Exit 4 `PGR-E4005`; (7) `N` Exit 4 `PGR-E4005`; (8) `Z` und Ende Exit 4 `PGR-E4002`; (9) `ssl=off` mit TLS-Wunsch Exit 4 `PGR-E4005`, Tabelle nicht angelegt; ohne TLS gegen `hostnossl reject` Exit 4 `PGR-E4005`; (10) `--upstream-ca` ohne TLS, auch mit `--upstream-tls=false` und `sslmode=require`, Exit 2 `PGR-E2001`; (11) fehlende Datei, Rechte 000, Verzeichnis, Textdatei, leere Datei, PEM nur mit Schlüssel: je Exit 2 `PGR-E2007`, Meldung nennt Option und Grund; (12) `--upstream-tls=false`, Umgebung `false` und Schlüssel `upstream_tls: false` gewinnen gegen `sslmode=require` (gegen `pgplain` Exit 0; `--upstream-tls=false` gegen `pgssl` Exit 4 mit „no encryption“); `sslmode=disable` plus `--upstream-tls` Exit 0; (13) `SSL_CERT_FILE` ersetzt `--upstream-ca`, `--upstream-ca` ergänzt `SSL_CERT_FILE`; (14) Klartext-Passwort: `pgpw` mit TLS Exit 0, Server-Log `method=password` und `SSL enabled`; `pgpwplain` ohne TLS Exit 0, kein `SSL enabled`; falsches Passwort Exit 4 `PGR-E4005`; (15) `record` mit `sslmode=require` Exit 2 `PGR-E2004`; `psql sslmode=require` gegen `replay`: „server does not support SSL“.
- **Weitere Proben.** Jede Sitzung baut TLS auf: Aufzeichnung mit zwei Sitzungen, im Log des Servers zwei Zeilen `SSL enabled` mehr. Schlüssel `upstream_tls` und `upstream_ca` im Abschnitt `play:` wirken, `config show` nennt beide. Das Beispiel `tls.yaml` aus Handbuch §4 wörtlich aus dem Handbuch herausgeschnitten und mit `--upstream-ca ./ca.pem` gegen `localhost` im Netzwerk-Namensraum von `pgssl` (Rolle `dev`, Datenbank `myapp`): Exit 0, Tabelle angelegt.
- **Meldungen.** `grep` über die Ausgaben aller Läufe auf Verzeichnisse, Dateinamen und Kopfzeile der CA-Dateien und das Passwort: kein Treffer.
- **Befehlsfolge aus Plan §3** (mit `grep -vx PGWIRE_RECORDER_PASSWORD`): Optionen des Handbuchs, die `--help` nicht kennt: `--h --help --name --rm --version`; Variablen: leer; Codes des Handbuchs außerhalb des Katalogs: leer; Gegenrichtung (Hilfe → Handbuch, Katalog → Handbuch): leer.
- **Verbotene Wörter.** `grep -i` im Handbuch auf Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review, „noch nicht“, „kommt“, „geplant“, „später“, „bisher“: fünf Treffer, vier im Bestand, einer neu (F-611). README: neue Zeilen ohne „sicher“, „geschützt“, Slice, Welle, „bisher“.
- **Namensabgleich.** `TestE2EPlayTLSZertifikatsspeicherDesSystems` und `test/integration/play_tls_e2e_test.go` existieren; die in §3 und §7 genannten Pläne (`slice-doku-ist-stand`, `slice-v1-abschluss-einspielen-anmeldung-doku`), `harness/mk/integration.mk` und `internal/hexagon/model/fehler.go` ebenfalls.
- **Gates.** `make docs-check`: 528 Dateien, 0 Befunde; `make kopf-check` und `make abdeckung-check`: Exit 0.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-609 | MEDIUM | DoD-Zeile 2 verlangt die Grenze *Zertifikatsspeicher nicht ladbar* und Abdeckungstabellen „über `make abdeckung` nachgezogen“; das Handbuch nennt die Grenze nicht, und kein Test ist geändert. §7 begründet das mit „Entscheidung des Nutzers, Option A, siehe §6“, aber §6 trägt keine Option A und keine Entscheidung, nur „gehört dem Nutzer (siehe Bericht des Architect)“; die DoD steht unverändert. Der Verifier findet die Entscheidung in keinem Artefakt des Repos. | `AGENTS.md` §3.9 (Plan folgt jeder Korrektur) | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-tls-doku.md` · „Entscheidung des Nutzers, Option A, siehe §6“ und „die Abdeckungstabellen sind über `make abdeckung` nachgezogen“ | nein (Urteil; `make abdeckung-check` ist grün, belegt aber nur, dass nichts zu schreiben war) | Entscheidung des Nutzers ohne Artefakt im Plan; DoD folgt der Abweichung nicht |
| F-610 | LOW | Handbuch §4 sagt, das Zertifikat gelte, wenn „die Zertifizierungsstelle, die es ausgestellt hat“ im Speicher oder in `--upstream-ca` steht, und der Katalog nennt als Ursache von `PGR-E4005`, sie stehe dort nicht. Ein Blatt, das ein Zwischenzertifikat ausgestellt hat und das der Server mit diesem sendet, gelingt mit nur der Wurzel in `--upstream-ca`, obwohl die ausstellende Stelle nicht in der Datei steht (Exit 0); Zwischenzertifikate sind in §7 als nicht gefahren geführt, die Sätze sagen über sie etwas. | `AGENTS.md` §3.11 (Satz sagt nur zu, was ein Test oder eine Probe prüft) | `docs/user/benutzerhandbuch.md` · „Das Zertifikat gilt, wenn die Zertifizierungsstelle, die es ausgestellt hat, im“ und Katalogzeile `PGR-E4005` · „Die Zertifizierungsstelle steht weder im Speicher des Systems noch in `--upstream-ca`“ | ja (Probe mit zwei Zertifikaten in einer Kette, wie oben) | Zusage über die Zertifizierungsstelle weiter als die Probe |
| F-611 | INFO | Der Beleg in §7 sagt, das `grep` auf „kommt“ finde keinen Treffer in den neuen Sätzen; Zeile „kommt. Ohne beides verbindet `play` ohne TLS.“ ist neu. Gemeint ist die Herkunft des Werts, nicht ein Ausblick; der Beleg ist ungenau, der Text trägt keine Zusage über Künftiges. Zuständig: Implementer (Wortlaut des Belegs). | Maintainability | `docs/user/benutzerhandbuch.md` · „kommt. Ohne beides verbindet `play` ohne TLS.“ | ja (`grep -n -i kommt`) | Beleg ungenauer als der Lauf |
| F-612 | INFO | Zur Frage der Reihenfolge: Die Optionstabelle in §5 führt `--config` vor `--log-level`, `play --help` führt `--log-level` vor `--config`. Das ist Bestand (nicht Teil des Diffs); die zwei neuen Zeilen stehen in der Reihenfolge von `play --help`. Kein Mangel dieses Diffs; ob die Tabelle der Hilfe folgen soll, entscheidet der Planner. | Maintainability | `docs/user/benutzerhandbuch.md` · „\| `--config` \| `record`, `replay`, `play`, `config show` \|“ | ja (`play --help` gegen die Tabelle) | Reihenfolge der Tabelle gegen die Hilfe |
| F-613 | INFO | Die Abweichung der Aussagen von DoD-Zeile 2 und der Zeile zu den Abdeckungstabellen trifft die Closure-Angleichung des Planners (siehe F-609); DoD-Konformität selbst prüft der Verifier. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-tls-doku.md` · „Der Planner gleicht die Zeile bei der Closure an.“ | nein | DoD-Angleichung bei Closure |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| (1) Regel Handbuch: keine Chronik, kein Zielstand, kein Verweis auf Spec, Lastenheft, ADR, Slice, Welle, Review im Diff; README ohne Stand-Aussagen, „sicher“ und „geschützt“ | geprüft, ohne Befund (F-611 betrifft nur den Beleg) |
| (2) Aussagen zu TLS gegen Binary und Server: alle Zeilen der Probe-Matrix (siehe oben), Beispiel `tls.yaml` wörtlich | geprüft, ohne Befund außer F-610 |
| (3) Nicht belegbare Aussagen (TLS-Version, Handshake-Frist, „nicht ladbar“, Zone, Fake `S` und Schweigen, MD5 und SCRAM mit TLS): im Text steht dazu nichts; der Klartext-Satz in §4 sagt nur, was der Server sieht | geprüft, ohne Befund |
| (4) Nachbarsätze: §1 (Voraussetzungen, „Verbindung zum Werkzeug läuft ohne Verschlüsselung“), §4 Hinweise, §4 *Ein Treiber anschließen* (`sslmode=prefer`), §5 Absatz zu `sslmode`, §6 *Rollen und Rechte*, §7 Fehlerbehebung (*Die Anwendung kann sich nicht verbinden*), Katalog (`PGR-E2001`, `E2004`, `E2007`, `E4002`, `E4005`), README („nur ohne TLS“) | geprüft, ohne Befund; „record verbindet ohne TLS“ und „unverschlüsselt“ stimmen mit Probe (15) und (9) |
| (5) Optionstabelle §5 gegen `play --help`: Namen, Umgebungsvariablen, Standardwerte der zwei neuen Zeilen | geprüft, ohne Befund; Reihenfolge Bestand siehe F-612 |
| (6) Option A: das Handbuch nennt die Grenze nicht und sagt stattdessen den beobachtbaren Fall (Zertifikat weder im Speicher noch in der Datei → `PGR-E4005`); der Ersatzsatz stimmt mit den Proben (4) und (13) | geprüft, ohne Befund am Handbuch; die Plan-Seite siehe F-609 |
| (7) `AGENTS.md` §3.13: der Diff nennt keinen anderen Slice neu als Adresse; „Planner“ in §7 ist keine Slice-Kennung | geprüft, ohne Befund |
| (7) `AGENTS.md` §3.9: Kopf, §1, §3 und §6 des Plans folgen dem Diff (Handbuch und README, keine Spec-Stelle); `make kopf-check` grün | geprüft, ohne Befund außer F-609 (DoD) |
| In §7 genannte Tests, Dateien und Pläne | geprüft, alle vorhanden |
| `//nolint`, Gates, Produkt-Code, Tests, `docs/user/abdeckung-*.md`, `spec/` | geprüft, ohne Befund (nicht im Diff) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Entscheidung des Nutzers ohne Artefakt im Plan · Zusage über die Zertifizierungsstelle weiter als die Probe · Beleg ungenauer als der Lauf · Reihenfolge der Tabelle gegen die Hilfe

## Verdikt

**Merge-blockierend:** ja, wegen F-609 (MEDIUM): Der Text von Handbuch und README ist an den Proben gedeckt und blockiert nicht; blockiert ist die Closure, solange die Entscheidung zur DoD-Zeile 2 in keinem Artefakt des Plans steht.

**Übergabe:**

- **F-609** an den **Planner** (und für die Entscheidung an den Nutzer): Artefakt dieser Übergabe ist dieser Report, Zeile F-609. Kein Widerspruch des Implementers liegt vor; der Konflikt-Pfad ist nicht eröffnet.
- **F-610** an den **Implementer**; ohne Blockade.
- **F-611** an den **Implementer**, **F-612** und **F-613** an den **Planner**; ohne Blockade, ohne erwartete Aktion im Diff.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
