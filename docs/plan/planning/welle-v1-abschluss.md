# Welle welle-v1-abschluss: v1-Abschluss: Sessions, Protokollrand, Betrieb, Einspielen, Zeitangaben, SQLite, Verbindungsende, Container

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M3.

**Verantwortlich:** pt9912. **Datum:** 2026-10-03.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die v1-Anforderungen sind bis auf zwei Nachträge umgesetzt: parallele Sessions im Record, definierter Protokollrand, TLS zum Client, kontrolliertes Herunterfahren mit atomarem Schreiben, nicht-interaktive Konfiguration, Einspielen einer Aufzeichnung (auch zeitgetreu und mit Antwortvergleich), die Aufzeichnung als Datenbankdatei, das Verbindungsende durch den Server (aufgezeichnet, wiedergegeben und beim Einspielen erwartet, Formatversion 2 nur mit Ende; [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), seit der Entscheidung des Nutzers vom 2026-10-12, [ADR-0039](../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)), ein Container-Image mit Betriebsdokumentation. Die Abnahmeszenarien 1 bis 10, 12 bis 16 und 18 des Lastenhefts sind nachweisbar, und alle Anforderungen (MUSS und SOLL) außer den zwei Nachträgen sind umgesetzt ([`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)). Die Nachträge — die Homebrew-Formel mit Abnahmeszenario 11 und der `CancelRequest` ohne Schlüssel — gehören seit der Entscheidung des Nutzers vom 2026-10-11 zu [welle-v1-nachtrag](welle-v1-nachtrag.md).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-replay-semantik](done/welle-replay-semantik/welle-replay-semantik.md) done.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Die Abnahmeszenarien 1 bis 10, 12 bis 16 und 18 des Lastenhefts laufen automatisiert (Szenario 16 mit seinem dritten Satz zum Verbindungsende, Szenario 18 mit Aufzeichnen und Replay; beide nur gegen PostgreSQL 17, nicht in der Versionsmatrix, Entscheidung des Nutzers vom 2026-10-12); Szenario 11 ist keine Bedingung der Welle und wird in `welle-v1-nachtrag` nachgewiesen — das *Mehr* gegenüber den Slice-DoDs.
- Das Binary ist reproduzierbar gebaut (gleiche Prüfsumme bei zwei Builds); das Container-Image ist gebaut.
- `make gates` grün.
- Closure-Notiz in `welle-v1-abschluss-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-v1-abschluss-sessions | Mehrere Sessions und Verbindungsfehler | [`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-v1-abschluss-sessions-doku | Handbuch und README zu: Mehrere Sessions und Verbindungsfehler | [`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-v1-abschluss-protokollrand | Protokollversion und Protokollrand | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-protokollrand-doku | Handbuch und README zu: Protokollversion und Protokollrand | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-postgres-versionen-matrix | Versionsmatrix im Gate | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
| slice-v1-abschluss-postgres-versionen-zusage | Zusage der unterstützten PostgreSQL-Versionen | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-RB-02`](../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung) |
| slice-v1-abschluss-anmeldung | Anmeldung des Clients im Record-Modus vermitteln | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-anmeldung-doku | Handbuch und README zu: Anmeldung des Clients im Record-Modus vermitteln | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-herunterfahren | Herunterfahren mit Frist in record und replay | [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../spec/lastenheft.md#lh-fa-15--ci-eignung) |
| slice-v1-abschluss-konfiguration | Konfiguration über Kommandozeile und Umgebung | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-konfigurationsdatei | Konfigurationsdatei und config show | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-verbindungen-platzhalter | Benannte Verbindungen und Platzhalter laden | [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-upstream-verbinden | Benannte Verbindung benutzen | [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-schreiben | Atomares Schreiben, --output und --force | [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings) |
| slice-v1-abschluss-einspielen | Einspielen einer Aufzeichnung | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-einspielen-laufsteuerung | Laufsteuerung beim Einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-einspielen-extended | Extended-Interaktionen einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) |
| slice-v1-abschluss-einspielen-extended-doku | Handbuch und README zu: Extended-Interaktionen einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) |
| slice-v1-abschluss-einspielen-anmeldung | Anmeldung beim Einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-einspielen-anmeldung-doku | Handbuch und README zu: Anmeldung beim Einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-einspielen-tls | TLS zum Server beim Einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-einspielen-tls-doku | Handbuch und README zu: TLS zum Server beim Einspielen | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-tls-client | TLS zum Client bei record und replay | [`LH-FA-23`](../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client) |
| slice-v1-abschluss-tls-client-doku | Handbuch und README zu: TLS zum Client bei record und replay | [`LH-FA-23`](../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client) |
| slice-v1-abschluss-zeitangaben-pruefen | Zeitangaben prüfen: absteigende Werte beim Lesen | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-zeitangaben-uhr | Uhr-Port und Systemuhr | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-zeitangaben-aufzeichnen | Zeitangaben aufzeichnen (`--record-timing`) | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-zeitangaben | Zeitgetreues Einspielen (`--keep-timing`) | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-zeitangaben-doku | Handbuch und README zu: Zeitangaben beim Aufzeichnen und Einspielen | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-antwortvergleich | Antwortvergleich beim Einspielen | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| slice-v1-abschluss-antwortvergleich-doku | Handbuch und README zu: Antwortvergleich beim Einspielen | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| slice-v1-abschluss-sqlite-format | SQLite als Aufzeichnungsformat | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-sqlite-format-lesen | SQLite-Aufzeichnungen lesen | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats) |
| slice-v1-abschluss-sqlite-format-doku | Handbuch und README zu: SQLite als Aufzeichnungsformat | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-verbindungsende-format | Verbindungsende im Format: Modell und YAML-Leser (Version 2) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-verbindungsende-format-schreiber | Verbindungsende im Format: YAML-Schreiber und Roundtrip (Version 2) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-verbindungsende-sqlite | Verbindungsende in SQLite: Leser (`meta.version` 1 und 2) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-verbindungsende-sqlite-schreiber | Verbindungsende in SQLite: Schreiber (`meta.version` in der Transaktion der ersten Session mit Ende) und Roundtrip beider Formate | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-verbindungsende-replay | Verbindungsende im Replay | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus) |
| slice-v1-abschluss-verbindungsende-einspielen | Verbindungsende beim Einspielen ohne Vergleich (`--allow-recorded-errors`) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-verbindungsende-vergleich | Verbindungsende im Antwortvergleich (Abnahmeszenario 16, Satz 3) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| slice-v1-abschluss-verbindungsende-record | Verbindungsende im Record: Use Case und Zustellung an den Client | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus) |
| slice-v1-abschluss-verbindungsende-upstream | Verbindungsende aus dem Upstream erkennen (Abnahmeszenario 18) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus) |
| slice-v1-abschluss-verbindungsende-doku | Handbuch und README zu: Verbindungsende durch den Server | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server) |
| slice-v1-abschluss-container | Container-Image und Betriebsdokumentation | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-erster-release](welle-erster-release.md) und mit ihr [welle-v1-nachtrag](welle-v1-nachtrag.md).
- Wird blockiert von: Welle [welle-replay-semantik](done/welle-replay-semantik/welle-replay-semantik.md).
- Innerhalb der Welle: `slice-v1-abschluss-konfiguration` setzt `slice-v1-abschluss-herunterfahren` voraus (der allgemeine Leser übernimmt dessen Umgebungsvariable); `slice-v1-abschluss-konfigurationsdatei` setzt `slice-v1-abschluss-konfiguration` voraus (die Schlüssel folgen aus der Anmeldung am allgemeinen Leser); `slice-v1-abschluss-verbindungen-platzhalter` setzt `slice-v1-abschluss-konfigurationsdatei` voraus (Laden der Datei, Form von `connections:`); `slice-v1-abschluss-upstream-verbinden` setzt `slice-v1-abschluss-verbindungen-platzhalter` voraus (eingesetzt wird, was das Laden zerlegt hat); `slice-v1-abschluss-schreiben` setzt `slice-v1-abschluss-konfiguration` und `slice-v1-abschluss-konfigurationsdatei` voraus (`--force` über den allgemeinen Leser, auch aus der Datei); `slice-v1-abschluss-einspielen` setzt `slice-v1-abschluss-herunterfahren`, `slice-v1-abschluss-konfiguration`, `slice-v1-abschluss-konfigurationsdatei`, `slice-v1-abschluss-verbindungen-platzhalter` und `slice-v1-abschluss-upstream-verbinden` voraus (Signalbehandlung, allgemeiner Leser, Konfigurationsdatei, Prüfen, Auflösen und Einsetzen einer benannten Verbindung für `--upstream` bei `play`); `slice-v1-abschluss-einspielen-laufsteuerung` setzt `slice-v1-abschluss-einspielen` voraus (Kommando `play`, Abbruch bei jedem Fehler, Abbruchsignal); `slice-v1-abschluss-einspielen-extended` setzt `slice-v1-abschluss-einspielen` und `slice-v1-abschluss-einspielen-laufsteuerung` voraus (die Randformen [L·E]: Fortsetzung und erwarteter Fehler innerhalb einer Extended-Interaktion); `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` setzen je `slice-v1-abschluss-einspielen` voraus (Kommando `play`, Verbindungsaufbau, Einstufung der Fehler im Aufbau), nicht einander und nicht die Laufsteuerung oder Extended; `slice-v1-abschluss-sqlite-format-lesen` setzt `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-einspielen` voraus (das Schreiben und der Schema-Text; `play`); `slice-v1-abschluss-sqlite-format-doku` setzt `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` voraus; `slice-v1-abschluss-zeitangaben-pruefen` setzt `slice-v1-abschluss-sqlite-format-lesen` voraus (die Leser beider Formate); `slice-v1-abschluss-zeitangaben-uhr` setzt `slice-v1-abschluss-einspielen` voraus (der Port steht vor seinen Benutzern); `slice-v1-abschluss-zeitangaben-aufzeichnen` setzt `slice-v1-abschluss-sqlite-format`, `slice-v1-abschluss-sqlite-format-lesen`, `slice-v1-abschluss-zeitangaben-pruefen` und `slice-v1-abschluss-zeitangaben-uhr` voraus (Schreiben und Roundtrip in beiden Formaten über die geprüften Leser, der Port); `slice-v1-abschluss-zeitangaben` setzt `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-uhr` und `slice-v1-abschluss-zeitangaben-aufzeichnen` voraus (das Warten braucht den Port und geprüfte Werte, das E2E von Abnahmeszenario 13 zeichnet mit `--record-timing` auf); `slice-v1-abschluss-zeitangaben-doku` setzt die drei Slices mit sichtbarem Verhalten voraus (`-pruefen`, `-aufzeichnen` und `slice-v1-abschluss-zeitangaben`; der Doku-Teil ist ungeteilt, wie bei `slice-v1-abschluss-sqlite-format-doku`); `slice-v1-abschluss-antwortvergleich` setzt `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-einspielen-laufsteuerung` und `slice-v1-abschluss-einspielen-extended` voraus (Optionen der Laufsteuerung mit Vergleich, Vergleich von Extended-Interaktionen); `slice-v1-abschluss-postgres-versionen-matrix` setzt `slice-extended-query-lebendpruefung` aus welle-extended-query voraus (Erkennung der Lebendprüfung) und `slice-v1-abschluss-anmeldung` (die Matrix fährt dessen Fall SCRAM-SHA-256 gegen jede Version); `slice-v1-abschluss-postgres-versionen-zusage` setzt `slice-v1-abschluss-postgres-versionen-matrix` voraus (die Zusage sagt nur, was das Gate belegt, und nimmt dessen Liste, Teilmenge und Aufrufe) und `slice-v1-abschluss-anmeldung-doku` (dessen Handbuch-Sätze zur Anmeldung zieht sie auf den Bereich nach); `slice-v1-abschluss-tls-client` setzt `slice-v1-abschluss-sessions` voraus (mehrere Verbindungen); die acht Slices zum Verbindungsende (Entscheidung des Nutzers vom 2026-10-12, Option B, [ADR-0039](../adr/0039-aufzeichnung-kennt-das-verbindungsende.md); Zuschnitt: Planner) stehen hinter `slice-v1-abschluss-protokollrand-doku` und vor `slice-harness-meldungskatalog-gate`; `-format` ist am 2026-10-15 geteilt in `-format` (Modell, Formregeln, YAML-Leser) und `-format-schreiber` (YAML-Schreiber, Roundtrip; setzt `-format` voraus), die Welle führt damit neun Slices zum Verbindungsende: `-sqlite`, `-replay` und `-einspielen` setzen je `-format` voraus (Modell, Formregeln und YAML-Leser; die Fixtures sind handgeschriebene Aufzeichnungen, ohne `record`), `-sqlite` zusätzlich `-format-schreiber` (Roundtrip beider Formate), `-vergleich` setzt `-einspielen` voraus (Schließen ohne `Terminate`, Einstufung ohne Vergleich), `-record` setzt `-format` und `-replay` voraus (Formregeln; Abbildung von `connection_closed` im PGWire-Adapter) und steht hinter `slice-v1-abschluss-sessions` und `slice-v1-abschluss-anmeldung`, weil beide `record.go` und den Aufbau im Upstream-Adapter berühren, `-upstream` setzt `-record`, `-sqlite` und `-format-schreiber` voraus (der Schreiber für die Aufzeichnung in YAML; der Use Case nimmt auf, was der Adapter liefert; E2E von Abnahmeszenario 18 auch in SQLite), `-doku` setzt die acht Code-Slices voraus (gebündelt hinter dem letzten, wie `slice-v1-abschluss-zeitangaben-doku`; ein Ende kann ein Anwender erst mit `record` erzeugen); jeder Doku-Folge-Slice `<Kennung>-doku` setzt seinen Code-Slice `<Kennung>` voraus und steht in der Reihenfolge direkt dahinter (Entscheidung des Nutzers vom 2026-10-10, `AGENTS.md` §3.11, §3.13: Handbuch- und README-Teil als eigene Schicht; die Welle schließt erst mit beiden, ein Release liegt nicht dazwischen); `slice-v1-abschluss-container` setzt die übrigen Slices voraus (Image und Doku bilden den Endstand ab), einschließlich der Doku-Folge-Slices und der Slices zum Verbindungsende. `slice-v1-abschluss-cancel-ohne-schluessel` und `slice-v1-abschluss-homebrew` stehen seit der Entscheidung des Nutzers vom 2026-10-11 in [welle-v1-nachtrag](welle-v1-nachtrag.md) (dort §5); sie hängen weiter an `slice-v1-abschluss-protokollrand` bzw. `slice-v1-abschluss-container` aus dieser Welle.
- Außerhalb der Welle, nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness): `slice-v1-abschluss-herunterfahren` setzt `slice-harness-integration-wait` voraus (dessen §1 gibt die Tests mit Signal und Frist an diesen Slice, geschrieben über die Helfer, die jener umbaut); `slice-v1-abschluss-container` setzt `slice-harness-meldungskatalog-gate` voraus (dessen §1 nennt das Gate als vorher geliefert). Beide bleiben wellenlos; sie sind die einzigen Harness-Slices in der Reihe unten.

**Reihenfolge** (Entscheidungen des Nutzers vom 2026-10-08: Wellen vor Harness, die beiden Harness-Vorläufer, der Schnitt von `slice-v1-abschluss-betrieb`, der Schnitt von `slice-v1-abschluss-konfiguration`; dazu die vorab benannten Schnitte von `slice-v1-abschluss-konfigurationsdatei` und `slice-v1-abschluss-verbindungen-platzhalter` vom 2026-10-09 die beiden Schnitte von `slice-v1-abschluss-einspielen` vom 2026-10-09 und der Schnitt von `slice-v1-abschluss-sqlite-format` vom 2026-10-11 (*Schreiben | Lesen*) und die Umplanung von `slice-v1-abschluss-cancel-ohne-schluessel` und `slice-v1-abschluss-homebrew` nach `welle-v1-nachtrag` vom 2026-10-11 (die Schritte dahinter sind neu nummeriert, 33 statt 35) und der Schnitt von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (*Prüfen | Uhr | Aufzeichnen | Einspielen*, Zuschnitt: Planner; die Schritte dahinter sind neu nummeriert, 36 statt 33) und die acht Slices zum Verbindungsende vom 2026-10-12 (Option B, Zuschnitt: Planner; der gebündelte Doku-Slice statt acht Doku-Folge-Slices ist Entscheidung des Nutzers vom 2026-10-12; Schritte 35 bis 42, die Schritte dahinter sind neu nummeriert, 44 statt 36) und deren Schnitt von `slice-v1-abschluss-verbindungsende-format` vom 2026-10-15 (*Leser | Schreiber*, Zuschnitt: Planner; Schritt 35 wird 35a und 35b, die Nummern dahinter bleiben) und der Schnitt von `slice-v1-abschluss-verbindungsende-sqlite` vom 2026-10-15 (*Leser | Schreiber*, Zuschnitt: Planner; Schritt 36 wird 36a und 36b, die Nummern dahinter bleiben) und der Schnitt von `slice-v1-abschluss-postgres-versionen` vom 2026-10-14 (*Matrix | Zusage*, Matrix zuerst, Entscheidung des Nutzers; Schritt 32 wird 32a und 32b, die Nummern dahinter bleiben); WIP-Limit 1; innerhalb der Abhängigkeiten nach Lieferwert für die Abnahmeszenarien und M3; die zehn Doku-Folge-Slices `-doku` stehen seit der Entscheidung des Nutzers vom 2026-10-10 direkt hinter ihrem Code-Slice):

1. `slice-harness-integration-wait` — wellenlos, Voraussetzung von Schritt 2.
2. `slice-v1-abschluss-herunterfahren` — Abnahmeszenario 5 (Herunterfahren im automatisierten Lauf); Voraussetzung des Einspielens.
3. `slice-v1-abschluss-konfiguration` — allgemeiner Leser für Kommandozeile und Umgebung; Voraussetzung von Schritt 4.
4. `slice-v1-abschluss-konfigurationsdatei` — Datei und `config show`, Teil von Abnahmeszenario 10 für Record und Replay; Voraussetzung von Schritt 5 und des Schreibens.
5. `slice-v1-abschluss-verbindungen-platzhalter` — Laden der benannten Verbindungen und Platzhalter, schließt die Lücke beim Klartext-Passwort; Voraussetzung von Schritt 6.
6. `slice-v1-abschluss-upstream-verbinden` — Benutzen der Verbindung, schließt Abnahmeszenario 10 für Record und Replay; Voraussetzung des Einspielens.
7. `slice-v1-abschluss-schreiben` — atomares Schreiben, `--output` und `--force`.
8. `slice-v1-abschluss-einspielen` — Kern: Abnahmeszenario 12 gegen einen Server ohne Passwort und ohne TLS, schließt Szenario 10 (Einspiel-Lauf); Voraussetzung von Schritt 9, 10, 12 und 14.
9. `slice-v1-abschluss-einspielen-laufsteuerung` — `--continue-on-error`, `--allow-recorded-errors`, `--finish-session-on-interrupt` (Boundary und Negative von [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung)); Voraussetzung von Schritt 10.
10. `slice-v1-abschluss-einspielen-extended` — Extended-Interaktionen einspielen, mit den Randformen [L·E] (setzt 8 und 9 voraus).
11. `slice-v1-abschluss-einspielen-extended-doku` — Handbuch und README dazu (setzt 10 voraus).
12. `slice-v1-abschluss-einspielen-anmeldung` — Anmeldung beim Einspielen (Klartext, MD5, SCRAM-SHA-256).
13. `slice-v1-abschluss-einspielen-anmeldung-doku` — Handbuch und README dazu (setzt 12 voraus).
14. `slice-v1-abschluss-einspielen-tls` — TLS zum Server beim Einspielen mit eigener Zertifizierungsstelle; hängt nicht an Schritt 12, steht direkt dahinter.
15. `slice-v1-abschluss-einspielen-tls-doku` — Handbuch und README dazu (setzt 14 voraus).
16. `slice-v1-abschluss-sqlite-format` — Schreiben von `sqlite` (`record --format sqlite`), Abnahmeszenario 14 erster Teil; geschnitten am 2026-10-11, das Lesen steht in Schritt 17.
17. `slice-v1-abschluss-sqlite-format-lesen` — Lesen von `sqlite` für `replay` und `play`, Formaterkennung, Fehlerklassen; schließt Abnahmeszenario 14 (setzt 16 und 8 voraus).
18. `slice-v1-abschluss-sqlite-format-doku` — Handbuch und README dazu (setzt 16 und 17 voraus).
19. `slice-v1-abschluss-zeitangaben-pruefen` — Lesen prüft die Werte: absteigende `offset_ms` einer Session sind `PGR-E3003`, in YAML und SQLite (setzt 17 voraus); geschnitten am 2026-10-11.
20. `slice-v1-abschluss-zeitangaben-uhr` — Uhr-Port, Fake-Uhr, Systemuhr und die Architekturregel dazu (setzt 8 voraus).
21. `slice-v1-abschluss-zeitangaben-aufzeichnen` — `record --record-timing`, Roundtrip in beiden Formaten (setzt 16, 17, 19 und 20 voraus).
22. `slice-v1-abschluss-zeitangaben` — `play --keep-timing` mit den Optionen, `PGR-E2001` und `PGR-E2003`; schließt Abnahmeszenario 13 (setzt 8, 19, 20 und 21 voraus).
23. `slice-v1-abschluss-zeitangaben-doku` — Handbuch und README dazu (setzt 22 voraus).
24. `slice-v1-abschluss-antwortvergleich` — Abnahmeszenario 16 (setzt 8 bis 10 voraus).
25. `slice-v1-abschluss-antwortvergleich-doku` — Handbuch und README dazu (setzt 24 voraus).
26. `slice-v1-abschluss-sessions` — Anforderung parallele Sessions; Voraussetzung von Schritt 28.
27. `slice-v1-abschluss-sessions-doku` — Handbuch und README dazu (setzt 26 voraus).
28. `slice-v1-abschluss-tls-client` — Abnahmeszenario 15.
29. `slice-v1-abschluss-tls-client-doku` — Handbuch und README dazu (setzt 28 voraus).
30. `slice-v1-abschluss-anmeldung` — Passwort-Anmeldung im Record; Voraussetzung von Schritt 32a.
31. `slice-v1-abschluss-anmeldung-doku` — Handbuch und README dazu (setzt 30 voraus).
32a. `slice-v1-abschluss-postgres-versionen-matrix` — Versionsmatrix im Gate: Szenarien 1, 2 und 7, Fall SCRAM-SHA-256, Lebendprüfungs-Abgleich gegen jede Version; geschnitten am 2026-10-14 aus `slice-v1-abschluss-postgres-versionen` (setzt 30 voraus).
32b. `slice-v1-abschluss-postgres-versionen-zusage` — Zusage der Versionen in Spezifikation, Lastenheft, Handbuch und README, mit den Handbuch-Sätzen zur Anmeldung (setzt 31 und 32a voraus).
33. `slice-v1-abschluss-protokollrand` — Protokollrand; Voraussetzung von `slice-v1-abschluss-cancel-ohne-schluessel` (welle-v1-nachtrag).
34. `slice-v1-abschluss-protokollrand-doku` — Handbuch und README dazu (setzt 33 voraus).
35a. `slice-v1-abschluss-verbindungsende-format` — Modell und YAML-Leser: `connection_closed`, Formregeln, Version 1 und 2 laden; geschnitten am 2026-10-15 (Voraussetzung von Schritt 35b, 36, 37, 38, 40).
35b. `slice-v1-abschluss-verbindungsende-format-schreiber` — YAML-Schreiber: Version 2 nur mit Ende, Roundtrip; geschnitten am 2026-10-15 aus `slice-v1-abschluss-verbindungsende-format` (setzt 35a voraus; Voraussetzung von Schritt 36 und 41).
36a. `slice-v1-abschluss-verbindungsende-sqlite` — SQLite-Leser: `meta.version` 1 und 2, `connection_closed`, Formregeln; geschnitten am 2026-10-15 (setzt 35a voraus; Voraussetzung von Schritt 36b und 37).
36b. `slice-v1-abschluss-verbindungsende-sqlite-schreiber` — SQLite-Schreiber: `meta.version` in der Transaktion der ersten Session mit Ende, Roundtrip beider Formate; geschnitten am 2026-10-15 aus `slice-v1-abschluss-verbindungsende-sqlite` (setzt 35a, 35b und 36a voraus; Voraussetzung von Schritt 41).
37. `slice-v1-abschluss-verbindungsende-replay` — Replay gibt das Ende wieder und schließt die Verbindung; Abbildung im PGWire-Adapter (setzt 35a voraus; Voraussetzung von Schritt 40).
38. `slice-v1-abschluss-verbindungsende-einspielen` — `play` ohne Vergleich: Ende erwartet mit `--allow-recorded-errors`, sonst `PGR-E4003`, Schließen ohne `Terminate` (setzt 35a voraus).
39. `slice-v1-abschluss-verbindungsende-vergleich` — `play --compare-responses`: Ende erwartet oder `PGR-E5004`, Zähler, Hilfetext; schließt Abnahmeszenario 16, Satz 3 (setzt 38 voraus).
40. `slice-v1-abschluss-verbindungsende-record` — Record-Service und PGWire-Zustellung des Endes (setzt 35a und 37 voraus; nach Schritt 26 und 30).
41. `slice-v1-abschluss-verbindungsende-upstream` — Upstream-Adapter erkennt das Ende; schließt Abnahmeszenario 18 (setzt 35b, 36b und 40 voraus).
42. `slice-v1-abschluss-verbindungsende-doku` — Handbuch und README dazu (setzt 35a bis 41 voraus).
43. `slice-harness-meldungskatalog-gate` — wellenlos, Voraussetzung von Schritt 44; misst die Codes aller vorigen Schritte, bevor die Betriebsdokumentation den Katalog schreibt.
44. `slice-v1-abschluss-container` — Abnahmeszenario 9 (setzt 1 bis 43 voraus).
Vor Schritt 28 stehen wellenlos die drei Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-13); keiner ist Voraussetzung eines Schritts und kein Schritt der ihren, die Nummern bleiben. Hinter Schritt 42 und vor Schritt 43 steht wellenlos `slice-harness-gegenprobe-matrix-p7-flanke` (Entscheidung des Nutzers vom 2026-10-15); er ist keine Voraussetzung eines Schritts, die Nummern bleiben. Nach dieser Welle folgen wellenlos `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage` (Abnahmeszenario 17, damit M3), dann [welle-erster-release](welle-erster-release.md), dann [welle-v1-nachtrag](welle-v1-nachtrag.md) (Homebrew und `CancelRequest` ohne Schlüssel, Entscheidung des Nutzers vom 2026-10-11), dann die übrigen Harness-Slices; deren Reihenfolge steht in ihrem §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- TLS zum Upstream im Record-Modus, Prüfung von Client-Zertifikaten, `COPY`, Replikationsprotokoll — nicht Teil des Produkts in dieser Welle.
- Gate für Code-Tabelle und Katalog — `slice-harness-meldungskatalog-gate` (wellenlos, Entscheidung des Nutzers vom 2026-10-07); er steht in der Reihenfolge (§5) direkt vor `slice-v1-abschluss-container` und zählt nicht zu den Slices dieser Welle.
- Auslösen eines Verbindungsendes durch das Einspielen, Aufzeichnen eines Netzwerkabbruchs oder eines Endes im Verbindungsaufbau, Wiedergabe eines Verbindungsendes ohne Fehlerantwort — Out-of-Scope von [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server); sie bleiben Verbindungsfehler oder nicht unterstützt wie bisher.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-v1-abschluss-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
