# Welle welle-v1-abschluss: v1-Abschluss: Sessions, Protokollrand, Betrieb, Einspielen, Zeitangaben, SQLite, Container

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

Die v1-Anforderungen sind bis auf zwei Nachträge umgesetzt: parallele Sessions im Record, definierter Protokollrand, TLS zum Client, kontrolliertes Herunterfahren mit atomarem Schreiben, nicht-interaktive Konfiguration, Einspielen einer Aufzeichnung (auch zeitgetreu und mit Antwortvergleich), die Aufzeichnung als Datenbankdatei, ein Container-Image mit Betriebsdokumentation. Die Abnahmeszenarien 1 bis 10 und 12 bis 16 des Lastenhefts sind nachweisbar, und alle Anforderungen (MUSS und SOLL) außer den zwei Nachträgen sind umgesetzt ([`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)). Die Nachträge — die Homebrew-Formel mit Abnahmeszenario 11 und der `CancelRequest` ohne Schlüssel — gehören seit der Entscheidung des Nutzers vom 2026-10-11 zu [welle-v1-nachtrag](welle-v1-nachtrag.md).

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
- Die Abnahmeszenarien 1 bis 10 und 12 bis 16 des Lastenhefts laufen automatisiert; Szenario 11 ist keine Bedingung der Welle und wird in `welle-v1-nachtrag` nachgewiesen — das *Mehr* gegenüber den Slice-DoDs.
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
| slice-v1-abschluss-postgres-versionen | Unterstützte PostgreSQL-Versionen | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
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
| slice-v1-abschluss-zeitangaben | Zeitangaben beim Aufzeichnen und Einspielen | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-zeitangaben-doku | Handbuch und README zu: Zeitangaben beim Aufzeichnen und Einspielen | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-antwortvergleich | Antwortvergleich beim Einspielen | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| slice-v1-abschluss-antwortvergleich-doku | Handbuch und README zu: Antwortvergleich beim Einspielen | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| slice-v1-abschluss-sqlite-format | SQLite als Aufzeichnungsformat | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-sqlite-format-lesen | SQLite-Aufzeichnungen lesen | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats) |
| slice-v1-abschluss-sqlite-format-doku | Handbuch und README zu: SQLite als Aufzeichnungsformat | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-container | Container-Image und Betriebsdokumentation | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-erster-release](welle-erster-release.md) und mit ihr [welle-v1-nachtrag](welle-v1-nachtrag.md).
- Wird blockiert von: Welle [welle-replay-semantik](done/welle-replay-semantik/welle-replay-semantik.md).
- Innerhalb der Welle: `slice-v1-abschluss-konfiguration` setzt `slice-v1-abschluss-herunterfahren` voraus (der allgemeine Leser übernimmt dessen Umgebungsvariable); `slice-v1-abschluss-konfigurationsdatei` setzt `slice-v1-abschluss-konfiguration` voraus (die Schlüssel folgen aus der Anmeldung am allgemeinen Leser); `slice-v1-abschluss-verbindungen-platzhalter` setzt `slice-v1-abschluss-konfigurationsdatei` voraus (Laden der Datei, Form von `connections:`); `slice-v1-abschluss-upstream-verbinden` setzt `slice-v1-abschluss-verbindungen-platzhalter` voraus (eingesetzt wird, was das Laden zerlegt hat); `slice-v1-abschluss-schreiben` setzt `slice-v1-abschluss-konfiguration` und `slice-v1-abschluss-konfigurationsdatei` voraus (`--force` über den allgemeinen Leser, auch aus der Datei); `slice-v1-abschluss-einspielen` setzt `slice-v1-abschluss-herunterfahren`, `slice-v1-abschluss-konfiguration`, `slice-v1-abschluss-konfigurationsdatei`, `slice-v1-abschluss-verbindungen-platzhalter` und `slice-v1-abschluss-upstream-verbinden` voraus (Signalbehandlung, allgemeiner Leser, Konfigurationsdatei, Prüfen, Auflösen und Einsetzen einer benannten Verbindung für `--upstream` bei `play`); `slice-v1-abschluss-einspielen-laufsteuerung` setzt `slice-v1-abschluss-einspielen` voraus (Kommando `play`, Abbruch bei jedem Fehler, Abbruchsignal); `slice-v1-abschluss-einspielen-extended` setzt `slice-v1-abschluss-einspielen` und `slice-v1-abschluss-einspielen-laufsteuerung` voraus (die Randformen [L·E]: Fortsetzung und erwarteter Fehler innerhalb einer Extended-Interaktion); `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` setzen je `slice-v1-abschluss-einspielen` voraus (Kommando `play`, Verbindungsaufbau, Einstufung der Fehler im Aufbau), nicht einander und nicht die Laufsteuerung oder Extended; `slice-v1-abschluss-sqlite-format-lesen` setzt `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-einspielen` voraus (das Schreiben und der Schema-Text; `play`); `slice-v1-abschluss-sqlite-format-doku` setzt `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` voraus; `slice-v1-abschluss-zeitangaben` setzt `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` voraus; `slice-v1-abschluss-antwortvergleich` setzt `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-einspielen-laufsteuerung` und `slice-v1-abschluss-einspielen-extended` voraus (Optionen der Laufsteuerung mit Vergleich, Vergleich von Extended-Interaktionen); `slice-v1-abschluss-postgres-versionen` setzt `slice-extended-query-lebendpruefung` aus welle-extended-query voraus (Erkennung der Lebendprüfung) und `slice-v1-abschluss-anmeldung` (die Matrix fährt dessen Fall SCRAM-SHA-256 gegen jede Version); `slice-v1-abschluss-tls-client` setzt `slice-v1-abschluss-sessions` voraus (mehrere Verbindungen); jeder Doku-Folge-Slice `<Kennung>-doku` setzt seinen Code-Slice `<Kennung>` voraus und steht in der Reihenfolge direkt dahinter (Entscheidung des Nutzers vom 2026-10-10, `AGENTS.md` §3.11, §3.13: Handbuch- und README-Teil als eigene Schicht; die Welle schließt erst mit beiden, ein Release liegt nicht dazwischen); `slice-v1-abschluss-container` setzt die übrigen Slices voraus (Image und Doku bilden den Endstand ab), einschließlich der Doku-Folge-Slices. `slice-v1-abschluss-cancel-ohne-schluessel` und `slice-v1-abschluss-homebrew` stehen seit der Entscheidung des Nutzers vom 2026-10-11 in [welle-v1-nachtrag](welle-v1-nachtrag.md) (dort §5); sie hängen weiter an `slice-v1-abschluss-protokollrand` bzw. `slice-v1-abschluss-container` aus dieser Welle.
- Außerhalb der Welle, nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness): `slice-v1-abschluss-herunterfahren` setzt `slice-harness-integration-wait` voraus (dessen §1 gibt die Tests mit Signal und Frist an diesen Slice, geschrieben über die Helfer, die jener umbaut); `slice-v1-abschluss-container` setzt `slice-harness-meldungskatalog-gate` voraus (dessen §1 nennt das Gate als vorher geliefert). Beide bleiben wellenlos; sie sind die einzigen Harness-Slices in der Reihe unten.

**Reihenfolge** (Entscheidungen des Nutzers vom 2026-10-08: Wellen vor Harness, die beiden Harness-Vorläufer, der Schnitt von `slice-v1-abschluss-betrieb`, der Schnitt von `slice-v1-abschluss-konfiguration`; dazu die vorab benannten Schnitte von `slice-v1-abschluss-konfigurationsdatei` und `slice-v1-abschluss-verbindungen-platzhalter` vom 2026-10-09 die beiden Schnitte von `slice-v1-abschluss-einspielen` vom 2026-10-09 und der Schnitt von `slice-v1-abschluss-sqlite-format` vom 2026-10-11 (*Schreiben | Lesen*) und die Umplanung von `slice-v1-abschluss-cancel-ohne-schluessel` und `slice-v1-abschluss-homebrew` nach `welle-v1-nachtrag` vom 2026-10-11 (die Schritte dahinter sind neu nummeriert, 33 statt 35); WIP-Limit 1; innerhalb der Abhängigkeiten nach Lieferwert für die Abnahmeszenarien und M3; die zehn Doku-Folge-Slices `-doku` stehen seit der Entscheidung des Nutzers vom 2026-10-10 direkt hinter ihrem Code-Slice):

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
19. `slice-v1-abschluss-zeitangaben` — Abnahmeszenario 13 (setzt 8, 16 und 17 voraus).
20. `slice-v1-abschluss-zeitangaben-doku` — Handbuch und README dazu (setzt 19 voraus).
21. `slice-v1-abschluss-antwortvergleich` — Abnahmeszenario 16 (setzt 8 bis 10 voraus).
22. `slice-v1-abschluss-antwortvergleich-doku` — Handbuch und README dazu (setzt 21 voraus).
23. `slice-v1-abschluss-sessions` — Anforderung parallele Sessions; Voraussetzung von Schritt 25.
24. `slice-v1-abschluss-sessions-doku` — Handbuch und README dazu (setzt 23 voraus).
25. `slice-v1-abschluss-tls-client` — Abnahmeszenario 15.
26. `slice-v1-abschluss-tls-client-doku` — Handbuch und README dazu (setzt 25 voraus).
27. `slice-v1-abschluss-anmeldung` — Passwort-Anmeldung im Record; Voraussetzung von Schritt 29.
28. `slice-v1-abschluss-anmeldung-doku` — Handbuch und README dazu (setzt 27 voraus).
29. `slice-v1-abschluss-postgres-versionen` — Versionsmatrix der Szenarien 1, 2 und 7.
30. `slice-v1-abschluss-protokollrand` — Protokollrand; Voraussetzung von `slice-v1-abschluss-cancel-ohne-schluessel` (welle-v1-nachtrag).
31. `slice-v1-abschluss-protokollrand-doku` — Handbuch und README dazu (setzt 30 voraus).
32. `slice-harness-meldungskatalog-gate` — wellenlos, Voraussetzung von Schritt 33; misst die Codes aller vorigen Schritte, bevor die Betriebsdokumentation den Katalog schreibt.
33. `slice-v1-abschluss-container` — Abnahmeszenario 9 (setzt 1 bis 32 voraus).
Nach dieser Welle folgen wellenlos `slice-harness-abdeckung-gate` und `slice-harness-coverage` (Abnahmeszenario 17, damit M3), dann [welle-erster-release](welle-erster-release.md), dann [welle-v1-nachtrag](welle-v1-nachtrag.md) (Homebrew und `CancelRequest` ohne Schlüssel, Entscheidung des Nutzers vom 2026-10-11), dann die übrigen Harness-Slices; deren Reihenfolge steht in ihrem §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- TLS zum Upstream im Record-Modus, Prüfung von Client-Zertifikaten, `COPY`, Replikationsprotokoll — nicht Teil des Produkts in dieser Welle.
- Gate für Code-Tabelle und Katalog — `slice-harness-meldungskatalog-gate` (wellenlos, Entscheidung des Nutzers vom 2026-10-07); er steht in der Reihenfolge (§5) direkt vor `slice-v1-abschluss-container` und zählt nicht zu den Slices dieser Welle.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-v1-abschluss-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
