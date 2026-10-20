# Review-Report: slice-v1-abschluss-cancel-ohne-schluessel — 2026-10-20

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Spezifikation `LH-FA-05.e` *Länge des `CancelRequest`*, `AGENTS.md` §3.7, §3.9 bis §3.14).

**Gegenstand:** `slice-v1-abschluss-cancel-ohne-schluessel`, Diff `f7aeea92..impl-cancel` (Kopf `192b437a`), sechs Commits, neun Dateien.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-20

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, keine Lifecycle-Pfade, und Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill). Das `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-v1-abschluss-cancel-ohne-schluessel` (§1 bis §8, insbesondere §3, §6 *Randformen* und §7 *Belege des Implementers*)
- Spezifikation `LH-FA-05.e`, Aufzählungspunkt *Länge des `CancelRequest`* und die Zeile *Startnachricht mit vollständigem Kopf* der Tabelle *Protokollrand*
- berührte Anforderungen `LH-FA-05`, `LH-FA-13`
- keine ADR in Commit-Messages genannt; Hard Rules in `AGENTS.md` §3.7, §3.9 bis §3.14
- vorherige Findings am Modul: `BEO-REPO/gruener-mutant-faelschlich-aequivalent`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`

**Läufe des Reviewers** (nur engste Sensoren, nacheinander, im Dev-Container `pgr-dev3` gegen eine Wegwerfkopie von `internal/`): Paket `pgwire` Unit-Tests grün (10,9 s) auf dem Stand `192b437a`; eine Sonde-Mutation (Prüfung des Startcodes im Zweig der Länge 12 entfernt) und eine Sonde-Testdatei (Kopf ohne Rest, Verbindung offen), beide in der Wegwerfkopie, danach gelöscht. Kein `make gates`, kein Integrationstest.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Spezifikation ordnet die Gründe „Länge des `CancelRequest`“ (`PGR-E6001` mit `ErrorResponse`) vor „Rest nicht geliefert“. Der Code liest für die Längen 9 bis 11, 13 bis 15 und 269 bis 10000 den Rest zuerst (die Bibliothek liest ihn vor der Längenprüfung): Sonde, Kopf der Länge 13, 15, 269 oder 10000 mit Code des `CancelRequest`, kein weiteres Byte, Verbindung offen — nach 700 ms weder `ErrorResponse` noch Schließen (Länge 8 antwortet sofort). Die Reihenfolge zeigt sich im Text (`EOF` statt `cancel request too short`) und im Zeitpunkt der Ablehnung; kein Test pinnt sie, alle Fälle mit fehlendem Rest schließen die Schreibseite. | Spezifikation `LH-FA-05.e` *Länge des `CancelRequest`* („Die Reihenfolge der Gründe ist: …“) · `AGENTS.md` §3.10 (Zusage ohne Mutation) | `internal/adapters/driving/pgwire/server.go` · `msg, err := k.be.ReceiveStartupMessage()` (Rest wird vor der Längenprüfung der Bibliothek gelesen); Spezifikation `spec/spezifikation.md` · „dann Rest der Nachricht nicht geliefert (`PGR-E6001`)“ | ja — ein Unit-Test mit offener Verbindung und Kopf ohne Rest, erwartet Antwort vor der Frist | Reihenfolge der Gründe laut Spezifikation vom Code nicht eingehalten und nicht gepinnt |
| F-2 | MEDIUM | §7 führt die vertauschte Reihenfolge von „Länge des `CancelRequest`“ und „Rest nicht geliefert“ als grünen Mutanten, „Keiner offen“, mit Invariante „beide Pfade enden in `fail(CodeUnsupported, …)`“. Die Invariante sagt nichts über den Zeitpunkt; die Sonde des Autors trennt die Pfade im Text, und die Sonde des Reviewers (F-1) trennt sie am Zeitpunkt der Ablehnung. Der Mutant ist damit nicht äquivalent im Sinn von §3.10 (a) bis (c), ein Fall dazu fehlt. | `AGENTS.md` §3.10 *Ein grüner Mutant gilt als äquivalent erst mit Invariante, Ort und Sonde* | `docs/plan/planning/in-progress/slice-v1-abschluss-cancel-ohne-schluessel.md` · „Beide Gründe sind `PGR-E6001` mit `ErrorResponse` und Exit-Code 6; ein Mutant, der die Reihenfolge vertauscht, lässt Code, Antwortart und Exit-Code gleich.“ | ja — derselbe Test wie F-1 färbt den Mutanten | Grüner Mutant als gleichwertig geführt, Sonde trennt ihn (`BEO-REPO/gruener-mutant-faelschlich-aequivalent`) |
| F-3 | LOW | §3 des Plans nennt als Ort der Änderung `startup`; die Umsetzung liegt in `startkopf` (Signatur geändert: `startkopf(k, laenge, code)`), `cancelOhneSchluessel` und `cancelGeschlossen`. Der Grund (Komplexitätsgrenze von `make lint`) steht nur in §7. | `AGENTS.md` §3.9 (Der Slice-Plan folgt jeder Korrektur) | `docs/plan/planning/in-progress/slice-v1-abschluss-cancel-ohne-schluessel.md` · „`internal/adapters/driving/pgwire` (`startup`)“ | nein — Urteil des Reviews | Plan §3 nennt den Ort der Umsetzung nicht |
| F-4 | INFO | `TestCancelRequestZweiteVerbindungUnberuehrt` sendet die Anfrage der zweiten Verbindung nach dem `CancelRequest`; eine Interaktion, die zum Zeitpunkt des `CancelRequest` läuft, kommt nicht vor. Der Mutant „Zweig beendet alle Verbindungen“ (M13) färbt den Test trotzdem rot; die Zusage „eine laufende Interaktion wird nicht abgebrochen“ der Spezifikation ist enger geprüft als formuliert. Hinweis an den Verifier. | `AGENTS.md` §3.11 (ein Kommentar sagt nur zu, was ein Test prüft) | `internal/adapters/driving/pgwire/server_cancel_test.go` · „setzt ihre Anfrage fort und endet regulär“ | nein | Zusage der Abdeckungs-Deklaration weiter als der Fall |

### Bewertung der offenen Frage des Implementers (Reihenfolge der Gründe)

Ja, das ist ein Widerspruch zum Wortlaut der Spezifikation, nicht nur ein Unterschied des Meldungstextes: Der Code lehnt Längen, die die Bibliothek ablehnt, erst nach dem Lesen des Rests ab; ein Client, der den Rest nicht liefert und die Verbindung offen lässt, sieht keine Ablehnung (F-1). Ob die Reihenfolge Vertrag ist (Code anpassen, zum Beispiel Länge vor dem Lesen des Rests prüfen) oder die Spezifikation enger zu fassen ist (Reihenfolge nur für Code, Antwortart und Exit-Code, nicht für den Zeitpunkt), ist eine Randform-Entscheidung und **geht als Rückgabe an den Architect**; der Reviewer entscheidet sie nicht. Handbuch und README sagen zur Reihenfolge nichts und widersprechen nicht.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/pgwire/server.go` (Zweig der Länge 12, `startkopf`, `cancelOhneSchluessel`, Reihenfolge Länge außerhalb vor Länge 12, TLS-Pflicht ohne `PGR-E6003`, Core-Reinheit, keine Direktive `nolint`) | geprüft, ohne Befund außer F-1; Sonde-Mutant „Prüfung des Startcodes im Zweig der Länge 12 entfernt“ färbt `TestTLSSonderanfrageMitFalscherLaenge` rot (Bestand), die neuen Tests fangen ihn nicht — Hinweis, kein Finding, die Zusage (`SSLRequest` mit Länge außer 8) gehört zum Bestand |
| `internal/adapters/driving/pgwire/server_cancel_test.go`, `server_protokollrand_test.go` (je Aufzählungselement ein Fall: Längen 8 bis 11, 13 bis 15, 269, 270, 10000, 12, 16, 17, 20, 268; Rest 0, 1, 3; sechs Stellen; zwei Modi; Mutationstabelle in §7 gegen den Code gelesen) | geprüft, ohne Befund außer F-1 und F-4; Unit-Tests des Pakets grün |
| `test/integration/cancel_e2e_test.go` (Szenario, Frist 20 s gegen `pg_sleep(0.5)`, Zählung genau einer Warnung, Exit-Code 0) | geprüft, ohne Befund; nicht gelaufen (kein Integrationslauf nach Auftrag), der Lauf liegt beim Implementer (§7) und beim Verifier |
| `docs/user/benutzerhandbuch.md`, `README.md` (Ist-Zustand, kein Verweis auf Spezifikation oder Slices im Handbuch, Quantoren „jede“ und „keine“ gegen die Proben-Tabelle in §7 gelesen, Tabellenzeilen `PGR-E6001` und `PGR-W3001`) | geprüft, ohne Befund; kein Beispiel als Datei im Diff, daher kein Beispiel zu laufen |
| `docs/user/abdeckung-e2e.md`, `docs/user/abdeckung-unit.md` | geprüft, ohne Befund (von `make abdeckung` geschrieben, Deklarationen entsprechen den Testkommentaren) |
| Plan §1, §6, §8 (Randformen vor dem Code: der Diff trägt keine neue Randform in §6 ein; keine Adresse neu oder geändert; Zählung in §8 unverändert zwei Liefer-Punkte, zwei Schichten) | geprüft, ohne Befund |
| Größenmetrik (`AGENTS.md` §3.14): `git diff --numstat f7aeea92..192b437a -- internal test docs/user README.md` | nachgerechnet: 354 (340 Einfügungen, 14 Löschungen) gegen Marke 700 und Schätzung 350; Marke nicht erreicht |
| Kommentare in Code und Tests (`AGENTS.md` §3.7) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Reihenfolge der Gründe laut Spezifikation vom Code nicht eingehalten und nicht gepinnt · Grüner Mutant als gleichwertig geführt, Sonde trennt ihn · Plan §3 nennt den Ort der Umsetzung nicht · Zusage der Abdeckungs-Deklaration weiter als der Fall

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM) blockieren, bis der Architect die Randform *Reihenfolge der Gründe* entschieden hat (Rückgabe, siehe oben); danach Code und Test (F-1, F-2) oder die engere Fassung der Spezifikation samt Plan (F-3).

**Übergabe:** F-1 als Randform-Rückgabe an den Architect, F-2 bis F-4 an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation; DoD-Konformität (`make gates`, Integrationstest) prüft der Verifier.
