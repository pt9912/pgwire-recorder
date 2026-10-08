# Welle welle-v1-abschluss: v1-Abschluss: Sessions, Protokollrand, Betrieb, Einspielen, Zeitangaben, SQLite, Container, Homebrew

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

Alle v1-Anforderungen sind umgesetzt: parallele Sessions im Record, definierter Protokollrand, TLS zum Client, kontrolliertes Herunterfahren mit atomarem Schreiben, nicht-interaktive Konfiguration, Einspielen einer Aufzeichnung (auch zeitgetreu und mit Antwortvergleich), die Aufzeichnung als Datenbankdatei, ein Container-Image mit Betriebsdokumentation und die Homebrew-Formel. Die Abnahmeszenarien 1 bis 10 und 12 bis 16 des Lastenhefts sind nachweisbar, Szenario 11 (Homebrew) weist `welle-erster-release` nach, und alle Anforderungen (MUSS und SOLL) sind umgesetzt ([`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)).

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
- Die Abnahmeszenarien 1 bis 10 und 12 bis 16 des Lastenhefts laufen automatisiert; Szenario 11 ist keine Bedingung der Welle und wird in `welle-erster-release` nachgewiesen — das *Mehr* gegenüber den Slice-DoDs.
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
| slice-v1-abschluss-protokollrand | Protokollversion und Protokollrand | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-postgres-versionen | Unterstützte PostgreSQL-Versionen | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-09`](../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) |
| slice-v1-abschluss-cancel-ohne-schluessel | CancelRequest ohne Schlüssel | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-v1-abschluss-anmeldung | Anmeldung des Clients im Record-Modus vermitteln | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-herunterfahren | Herunterfahren mit Frist in record und replay | [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../spec/lastenheft.md#lh-fa-15--ci-eignung) |
| slice-v1-abschluss-konfiguration | Konfiguration über Kommandozeile und Umgebung | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-konfigurationsdatei | Konfigurationsdatei, benannte Verbindungen und config show | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-schreiben | Atomares Schreiben, --output und --force | [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings) |
| slice-v1-abschluss-homebrew | Homebrew-Bereitstellung | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |
| slice-v1-abschluss-einspielen | Einspielen einer Aufzeichnung | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) |
| slice-v1-abschluss-tls-client | TLS zum Client bei record und replay | [`LH-FA-23`](../../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client) |
| slice-v1-abschluss-zeitangaben | Zeitangaben beim Aufzeichnen und Einspielen | [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-antwortvergleich | Antwortvergleich beim Einspielen | [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) |
| slice-v1-abschluss-sqlite-format | SQLite als Aufzeichnungsformat | [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| slice-v1-abschluss-container | Container-Image und Betriebsdokumentation | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-erster-release](welle-erster-release.md).
- Wird blockiert von: Welle [welle-replay-semantik](done/welle-replay-semantik/welle-replay-semantik.md).
- Innerhalb der Welle: `slice-v1-abschluss-konfiguration` setzt `slice-v1-abschluss-herunterfahren` voraus (der allgemeine Leser übernimmt dessen Umgebungsvariable); `slice-v1-abschluss-konfigurationsdatei` setzt `slice-v1-abschluss-konfiguration` voraus (die Schlüssel folgen aus der Anmeldung am allgemeinen Leser); `slice-v1-abschluss-schreiben` setzt `slice-v1-abschluss-konfiguration` und `slice-v1-abschluss-konfigurationsdatei` voraus (`--force` über den allgemeinen Leser, auch aus der Datei); `slice-v1-abschluss-einspielen` setzt `slice-v1-abschluss-herunterfahren`, `slice-v1-abschluss-konfiguration` und `slice-v1-abschluss-konfigurationsdatei` voraus (Signalbehandlung, allgemeiner Leser, Konfigurationsdatei); `slice-v1-abschluss-zeitangaben` setzt `slice-v1-abschluss-einspielen` und `slice-v1-abschluss-sqlite-format` voraus; `slice-v1-abschluss-antwortvergleich` setzt `slice-v1-abschluss-einspielen` voraus; `slice-v1-abschluss-cancel-ohne-schluessel` setzt `slice-v1-abschluss-protokollrand` voraus; `slice-v1-abschluss-postgres-versionen` setzt `slice-extended-query-lebendpruefung` aus welle-extended-query voraus (Erkennung der Lebendprüfung) und `slice-v1-abschluss-anmeldung` (die Matrix fährt dessen Fall SCRAM-SHA-256 gegen jede Version); `slice-v1-abschluss-tls-client` setzt `slice-v1-abschluss-sessions` voraus (mehrere Verbindungen); `slice-v1-abschluss-container` setzt die übrigen Slices außer `slice-v1-abschluss-homebrew` voraus (Image und Doku bilden den Endstand ab); `slice-v1-abschluss-homebrew` setzt `slice-v1-abschluss-container` voraus (Release-Artefakte).
- Außerhalb der Welle, nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness): `slice-v1-abschluss-herunterfahren` setzt `slice-harness-integration-wait` voraus (dessen §1 gibt die Tests mit Signal und Frist an diesen Slice, geschrieben über die Helfer, die jener umbaut); `slice-v1-abschluss-container` setzt `slice-harness-meldungskatalog-gate` voraus (dessen §1 nennt das Gate als vorher geliefert). Beide bleiben wellenlos; sie sind die einzigen Harness-Slices in der Reihe unten.

**Reihenfolge** (Entscheidungen des Nutzers vom 2026-10-08: Wellen vor Harness, die beiden Harness-Vorläufer, der Schnitt von `slice-v1-abschluss-betrieb`, der Schnitt von `slice-v1-abschluss-konfiguration`; WIP-Limit 1; innerhalb der Abhängigkeiten nach Lieferwert für die Abnahmeszenarien und M3):

1. `slice-harness-integration-wait` — wellenlos, Voraussetzung von Schritt 2.
2. `slice-v1-abschluss-herunterfahren` — Abnahmeszenario 5 (Herunterfahren im automatisierten Lauf); Voraussetzung des Einspielens.
3. `slice-v1-abschluss-konfiguration` — allgemeiner Leser für Kommandozeile und Umgebung; Voraussetzung von Schritt 4.
4. `slice-v1-abschluss-konfigurationsdatei` — Abnahmeszenario 10 für Record und Replay; Voraussetzung des Schreibens und des Einspielens.
5. `slice-v1-abschluss-schreiben` — atomares Schreiben, `--output` und `--force`.
6. `slice-v1-abschluss-einspielen` — Abnahmeszenario 12, schließt Szenario 10 (Einspiel-Lauf).
7. `slice-v1-abschluss-sqlite-format` — Abnahmeszenario 14; nach Schritt 6 ist auch dessen zweiter Satz (Einspielen der Datenbankdatei) prüfbar.
8. `slice-v1-abschluss-zeitangaben` — Abnahmeszenario 13 (setzt 6 und 7 voraus).
9. `slice-v1-abschluss-antwortvergleich` — Abnahmeszenario 16 (setzt 6 voraus).
10. `slice-v1-abschluss-sessions` — Anforderung parallele Sessions; Voraussetzung von Schritt 11.
11. `slice-v1-abschluss-tls-client` — Abnahmeszenario 15.
12. `slice-v1-abschluss-anmeldung` — Passwort-Anmeldung im Record; Voraussetzung von Schritt 13.
13. `slice-v1-abschluss-postgres-versionen` — Versionsmatrix der Szenarien 1, 2 und 7.
14. `slice-v1-abschluss-protokollrand` — Protokollrand; Voraussetzung von Schritt 15.
15. `slice-v1-abschluss-cancel-ohne-schluessel`.
16. `slice-harness-meldungskatalog-gate` — wellenlos, Voraussetzung von Schritt 17; misst die Codes aller vorigen Schritte, bevor die Betriebsdokumentation den Katalog schreibt.
17. `slice-v1-abschluss-container` — Abnahmeszenario 9 (setzt 1 bis 16 voraus).
18. `slice-v1-abschluss-homebrew` — Verfahren für Abnahmeszenario 11, Nachweis in welle-erster-release (setzt 17 voraus).

Nach dieser Welle folgen wellenlos `slice-harness-abdeckung-gate` und `slice-harness-coverage` (Abnahmeszenario 17, damit M3), dann [welle-erster-release](welle-erster-release.md), dann die übrigen Harness-Slices; deren Reihenfolge steht in ihrem §4.

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
