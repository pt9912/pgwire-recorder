# Slice slice-v1-abschluss-cancel-ohne-schluessel: CancelRequest ohne Schlüssel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)

**Berührte Spec-Stellen:** `LH-FA-05.e` (Spezifikation, Protokollrand, *Länge des `CancelRequest`*) · `LH-FA-13.b` · `ARC-006`

**Verantwortlich:** pt9912 (Zweig impl-cancel)

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein `CancelRequest` der Länge 12, den ein Client ohne Abbruchschlüssel sendet (pgx, weil der Recorder `BackendKeyData` nicht weitergibt), wird wie jeder `CancelRequest` behandelt: Verbindung geschlossen, Warnung `PGR-W3001`, kein Verbindungsfehler (`LH-FA-05.e`), in `record` und `replay`, mit und ohne TLS. Heute lehnt die PGWire-Bibliothek ihn als „cancel request too short“ ab (gemessen mit v5.11.0: sie nimmt die Längen 16 bis 268 an, die Länge 12 und alle darunter und dazwischen nicht), und der Recorder meldet `PGR-E6001` mit Klasse 6 (Review-Finding F-310 aus `docs/reviews/2026-10-05-review-slice-extended-query-record.md`, Sonde P3). Die Längen, ihre Folgen und die Reihenfolge der Gründe stehen in der Spezifikation (`LH-FA-05.e`, *Länge des `CancelRequest`*); der Slice setzt sie um.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Weitergabe von `BackendKeyData` und echtes Abbrechen am Upstream — `LH-FA-05.e` legt fest, dass ein `CancelRequest` nicht weitergeleitet wird; eine Änderung wäre eine neue Anforderung.
- Übriger Protokollrand — ist mit `slice-v1-abschluss-protokollrand` geliefert (Bestand, `done`); dieser Slice ändert nur die Länge 12 des `CancelRequest`.
- Die Längen 8 bis 11, 13 bis 15 und 269 bis zur Höchstlänge bleiben `PGR-E6001`, die Längen 16 bis 268 bleiben bei der Bibliothek — Bestand bleibt bewusst stehen: Kein Client sendet sie (die Spezifikation entscheidet sie jetzt), und eine eigene Auswertung brächte kein Verhalten, das ein Test heute verlangt.
- Das Verbindungsende einer Hauptverbindung, die ein Client bei Kontext-Abbruch mitten in einer Anfrage schließt — bleibt `PGR-E4003` (Exit-Code `4`) nach dem Abschnitt *Verbindungsende* der Spezifikation; es ist der Abbruch der Anfrage, nicht der `CancelRequest`, und eine andere Anforderung wäre ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Ein `CancelRequest` ohne Schlüssel (Länge 12) schließt die Verbindung mit `PGR-W3001` ohne Verbindungsfehler, in `record` und `replay` (Unit-Test mit 12 Byte). Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): pgx mit Kontext-Abbruch über `record` lässt den Lauf mit Exit-Code 0 enden, wenn sonst kein Fehler auftrat (Integrationstest); der Abbruch trifft eine Anfrage, die der Server vor der Frist des Clients beantwortet, sodass die Hauptverbindung regulär endet (Szenario in §3).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/pgwire` (`startup`) | update | Nur die Länge 12 ist neu: bei Startcode `CancelRequest` und Länge 12 liest der Adapter die Prozess-ID selbst (die Bibliothek lehnt 12 ab), meldet `PGR-W3001` und schließt; fehlen die 4 Bytes, ist es `PGR-E6001`. Länge 8 bis 11, 13 bis 15 und 269 bis zur Höchstlänge bleiben bei der Bibliothek (`PGR-E6001`, gemessen), 16 bis 268 ebenso (`PGR-W3001`) |
| `internal/adapters/driving/pgwire/server_cancel_test.go` (neu), `server_protokollrand_test.go` (`randServer` hält das Log) | update | Negative nach `LH-FA-05`: ein Fall je Länge 8 bis 11, 13 bis 15, 269, 270 und 10000 (`PGR-E6001`) und 12, 16, 17, 20 und 268 (`PGR-W3001`); Länge 12, 16 und 268 mit 0, 1 und 3 Bytes Rest (`PGR-E6001`); Prozess-ID und Schlüssel 0, `0xff`, `0x2a`; Länge 0, 4, 7 und 10001 mit dem Code des `CancelRequest` (`PGR-W3003`, Reihenfolge der Gründe); jeweils in `record` und `replay` und an sechs Stellen (erste Nachricht, nach `SSLRequest` mit `N`, nach `GSSENCRequest`, über TLS nach `S`, unverschlüsselt bei TLS-Pflicht, nach `GSSENCRequest` bei TLS-Pflicht); eine zweite Verbindung bleibt unberührt und die nächste öffnet die zweite Session |
| `test/integration/cancel_e2e_test.go` (neu) | update | pgx mit Kontext-Abbruch, Exit-Code: Der Kontext einer Anfrage mit `pg_sleep(0.5)` endet nach 100 ms; `BuildContextWatcherHandler` setzt `CancelRequestContextWatcherHandler` mit `DeadlineDelay` 20 s (der Standard von pgx sendet bei laufender Anfrage keinen `CancelRequest`); die Anfrage endet regulär, die Hauptverbindung sendet `Terminate`, `SIGTERM` an den Recorder, Exit-Code 0, im Log genau einmal `PGR-W3001` und weder `PGR-E6001` noch `PGR-E4003` |
| `docs/user/benutzerhandbuch.md`, `README.md`, `docs/user/abdeckung-e2e.md`, `docs/user/abdeckung-unit.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung); die Abdeckungstabellen schreibt `make abdeckung` |

**Zusagen und Mutationen** (`AGENTS.md` §3.10; der Bericht nennt je Zeile den roten Test): Länge 12 ist `PGR-W3001` (Mutation: Zweig für 12 entfernt, oder 12 bis 16 angenommen) · Länge 11 und 13 sind `PGR-E6001` (Mutation: Untergrenze 11 oder Bereich 12 bis 15) · 16 und 268 sind `PGR-W3001`, 269 `PGR-E6001` (Fälle pinnen die Bibliothek) · gleich in `record` und `replay` (Mutation: Zweig nur im `record`) · mit TLS-Pflicht kein `PGR-E6003` (Mutation: Prüfung der Pflicht vor dem `CancelRequest`) · nach `SSLRequest` mit `N` und über TLS nach `S` · Rest nicht geliefert ist `PGR-E6001` (Mutation: Lesefehler verworfen) · die Verbindung wird geschlossen (Mutation: sie bleibt offen) · Exit-Code 0 (Integrationstest, Mutation: Länge 12 wie bisher abgelehnt, Lauf endet mit 6).

**Größe** (`AGENTS.md` §3.14): Schätzung in der Metrik von §4, mit dem Doppelten für Tests: Produktcode rund 25, Unit-Tests rund 200, Integrationstest rund 100, Handbuch und README rund 30, zusammen **rund 350**. Gemessen (Basis `f7aeea92`, letzter Commit `f3ba1835`): 354 (340 Einfügungen, 14 Löschungen). Zwei Liefer-Punkte, zwei Schichten (§8).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-protokollrand` ist `done`, und `welle-erster-release` ist `done` (Start-Trigger von `welle-v1-nachtrag`; Entscheidung des Nutzers vom 2026-10-11).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Behandlung verlangt die Weitergabe von `BackendKeyData`, oder die **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user README.md`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist, liegt **über 700** (Schätzung in §3: rund 350).
- `in-progress` → `open` (blockiert — Carveout?): der Adapter kann die Länge 12 nicht vor der Bibliothek lesen (der Kopf wird heute vor der Bibliothek gelesen, `startup`; die Bedingung gilt, wenn das mit `LH-FA-23.a` nicht zusammen geht).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12), alle entschieden vor dem Code in der Spezifikation, `spec/spezifikation.md` `LH-FA-05.e`, Aufzählungspunkt *Länge des `CancelRequest`* (Nachfolgend: „Punkt“):

- Länge 12 (Prozess-ID, kein Schlüssel) und 16 bis 268 — `PGR-W3001`, Punkt.
- Länge 8 bis 11, 13 bis 15, 269 bis zur Höchstlänge — `PGR-E6001` mit `ErrorResponse` (Exit-Code 6), Punkt; unter 8 und über der Höchstlänge `PGR-W3003`, Tabelle *Protokollrand*.
- Reihenfolge der Gründe: Länge außerhalb von 8 bis zur Höchstlänge, dann Länge des `CancelRequest`, dann Rest nicht geliefert — Punkt.
- Rest der Nachricht nicht geliefert (Abbruch, Ende mitten in der Nachricht) — `PGR-E6001`, Punkt und Zeile *Startnachricht mit vollständigem Kopf*.
- Stelle: erste Nachricht und nach `SSLRequest` (`N`, oder `S` und TLS) oder `GSSENCRequest` — Punkt; mit TLS-Pflicht ohne `PGR-E6003` — Abschnitt *Unverschlüsselte Clients* von `LH-FA-23.a`.
- Modus: `record` und `replay` gleich, keine Session, keine Aufzeichnung, keine Verbindung zum Server, eine zweite Verbindung unberührt — Punkt und *Modus*.
- Prozess-ID und Schlüssel — nicht ausgewertet, Punkt.
- Exit-Code — die Warnung ändert ihn nicht (`LH-FA-13`); eine Hauptverbindung, die der Client bei Kontext-Abbruch mitten in der Anfrage schließt, ist `PGR-E4003` (Exit-Code 4) nach *Verbindungsende*, Punkt.
- Plattform (`LH-QA-03`): die Eigenschaft hängt weder vom Betriebssystem noch vom Dateisystem ab (Bytes auf einer TCP-Verbindung); Unit- und Integrationstest laufen unter Linux im Container, eine Grenze je Plattform gibt es nicht — Punkt.
- Aufruf-Voraussetzung und Aufrufarten des Treibers: nicht zutreffend (keine Ordnung von Werten über Aufrufe, kein Stand, kein Nachweis).

Risiken:

- Ein Client sendet einen `CancelRequest` in einer Länge, die die Spezifikation als `PGR-E6001` führt (zum Beispiel 13 bis 15 oder über 268) — **Ausgang:** offen bis Closure.
- Der Integrationstest trifft die Zeitfolge (Anfrage endet vor der Frist des Clients) nicht stabil — **Ausgang:** offen bis Closure; Gegenmittel: die Frist des Clients um ein Vielfaches über der Dauer der Anfrage.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

Regeln: `.claude/commands/implement-slice.md` Schritte 17 bis 19. Die Läufe stehen unten; der Bericht verweist hierher.

**Wege.** Die Mutanten liefen gegen einen Bind-Mount statt gegen den Build-Kontext (`docker run --network none -v $PWD/internal:/src/internal:ro pgwire-recorder:test go test ./internal/adapters/driving/pgwire -run 'TestCancelRequest|TestProtokollrand|TestTLS'`), der Mutant des Integrationstests mit `touch` auf `server.go` und `PGR_LAUF=pgr-it-cancel PGR_TESTS_1=TestE2ERecordCancelOhneSchluessel make test-integration`. Mutation und Zurücksetzen schrieben die Datei neu (kein `cp -p`); nach dem Lauf war `server.go` byte-gleich mit dem Stand vor der Mutation (`cmp`).

**Mutationen** (je Zeile: Zusage · Mutation · roter Test, Zahl der roten Unterfälle in Klammern; alle Mutanten färbten nur Tests des Pakets `pgwire` rot, nichts anderes lief):

| Zusage (Spezifikation `LH-FA-05.e`, *Länge des `CancelRequest`*) | Mutation | roter Test |
|---|---|---|
| Länge 12 ist `PGR-W3001`, in `record` und `replay` | M1: Bedingung des Zweigs `laenge == 0` (12 wie bisher von der Bibliothek abgelehnt) | `TestCancelRequestLaenge/…/Länge_12` (38), `TestCancelRequestZweiteVerbindungUnberuehrt/…/Länge_12` |
| Länge 13 bis 15 sind `PGR-E6001` (je Element) | M2a: Bereich 12 bis 15; M2c: Bereich 12 bis 16 | `TestCancelRequestLaenge/…/Länge_13`, `_14`, `_15` (36 je Mutant) |
| Länge 11 ist `PGR-E6001` | M2d: Bereich 11 bis 12, liest `laenge-8` Bytes | `TestCancelRequestLaenge/…/Länge_11` (12) |
| Länge 8, 9, 10 sind `PGR-E6001` | M2e: Bereich 8 bis 12, liest `laenge-8` Bytes | `…/Länge_8`, `_9`, `_10`, `_11` (52, dazu `TestProtokollrandStartnachricht`) |
| Länge 16 ist `PGR-W3001` | M2i: Länge 16 mit `PGR-E6001` abgelehnt | `…/Länge_16`, `Länge_16,_Schlüssel_0xff`, `TestCancelRequest`, `TestTLSRangfolge/CancelRequest` (27) |
| Länge 268 ist `PGR-W3001` | M2h: Länge 268 mit `PGR-E6001` abgelehnt | `…/Länge_268` (12) |
| Länge 269 ist `PGR-E6001` | M2f: Bereich `12 \|\| 269`, liest `laenge-8` Bytes | `…/Länge_269` (12) |
| Länge 270 und 10000 (Höchstlänge) sind `PGR-E6001` | M2g: Bereich `12 \|\| 10000`, liest `laenge-8` Bytes | `…/Länge_10000` (12) |
| gleich in `record` und `replay` | M3: Zweig nur im `record`; M3b: nur im `replay` | `…/replay=true/…/Länge_12` (19); `…/replay=false/…/Länge_12` (19) |
| mit TLS-Pflicht kein `PGR-E6003` | M4: Prüfung der Pflicht vor dem `CancelRequest` | `…/unverschlüsselt_bei_TLS-Pflicht/Länge_12` (24) |
| über TLS nach `S` | M11: Zweig nur ohne TLS | `…/über_TLS_nach_S/Länge_12` (6) |
| nach `SSLRequest` mit `N` und nach `GSSENCRequest` | M12: Zweig nur als erste Nachricht | `…/nach_SSLRequest_mit_N/…`, `…/nach_GSSENCRequest_mit_N/…` (24) |
| Rest nicht geliefert ist `PGR-E6001` mit `ErrorResponse` (Länge 12; 16 und 268 laufen über die Bibliothek) | M5: Lesefehler verworfen; M15: Fehler ohne `ErrorResponse` (nur gemerkt); M16: Warnung `PGR-W3003` statt Fehler | `…/Länge_12,_0_von_4_Bytes_Rest`, `_1_von_…`, `_3_von_…` (36 je Mutant) |
| die Verbindung wird geschlossen | M6: Verbindung bleibt 30 s offen | `…/Länge_12` (10 s Frist des Clients: „Verbindung nicht geschlossen“, 38) |
| Warnung `PGR-W3001`, Stufe WARN, genau eine | M7: keine Warnung; M7b: Stufe DEBUG; M7c: Code `PGR-W3003` | `…/Länge_12`, `Länge_16`, `Länge_268` (97 je Mutant), `TestCancelRequest`, `TestTLSRangfolge/CancelRequest` |
| keine Antwort nach der Warnung (verneinend) | M17: eine `ErrorResponse` nach der Warnung | `…/Länge_12` (38); der Test liest bis zum Ende des Stroms und prüft, dass keine Nachricht kommt |
| kein Verbindungsfehler (verneinend) | M18: Fehler gemerkt (Exit-Code 6) | `…/Länge_12` (36), geprüft über `FirstErrorCode` |
| keine Warnung `PGR-W3001` bei `PGR-E6001` (verneinend) | M19: Warnung auch bei fehlendem Rest | `…/Länge_12,_0_von_4_Bytes_Rest` (36) |
| Prozess-ID und Schlüssel werden nicht ausgewertet | M8: Prozess-ID mit erstem Byte `0xff` abgelehnt; M8b: Prozess-ID mit letztem Byte 0 abgelehnt | `…/Länge_12,_Prozess-ID_0xff` (12); `…/Länge_12,_Prozess-ID_und_Schlüssel_0` (12) |
| Reihenfolge: Länge außerhalb von 8 bis 10000 (`PGR-W3003`) vor der Länge des `CancelRequest` | M9: `CancelRequest` unter 8 mit `PGR-E6001`; M9b: über 10000 mit `PGR-E6001` | `TestCancelRequestLaengeAusserhalb/…/Länge_0`, `_4`, `_7` und `TestProtokollrandPruefreihenfolge` (14); `…/Länge_10001` (6) |
| keine Session | M14: der Zweig öffnet eine Session | `…/Länge_12` (74) |
| eine zweite Verbindung bleibt unberührt, die laufende Anfrage wird nicht abgebrochen | M13: der Zweig beendet alle Verbindungen | `TestCancelRequestZweiteVerbindungUnberuehrt/…` |
| Exit-Code 0 (`LH-FA-13`) | M1 im Binary: Länge 12 wie bisher abgelehnt | `TestE2ERecordCancelOhneSchluessel`: „Exit-Code 6, erwartet 0“ |

Die Mutanten M2f bis M2i ersetzen das Verhalten der Bibliothek (v5.11.0) durch einen eigenen Zweig und zeigen, dass die Fälle für 16, 268, 269 und 10000 die Längen pinnen; ein Mutant im Quelltext der Bibliothek ist nicht möglich. Die Schlüssel von 4 bis 256 Bytes (16 bis 268) wertet allein die Bibliothek aus; dass der Recorder sie nicht auswertet, belegt der Fall `Länge_16,_Schlüssel_0xff` nur gegen die Bibliothek in dieser Version.

**Grüne Mutanten.** Keiner offen. Eingeordnet:

- Reihenfolge *Länge des `CancelRequest`* vor *Rest nicht geliefert* (Spezifikation): Beide Gründe sind `PGR-E6001` mit `ErrorResponse` und Exit-Code 6; ein Mutant, der die Reihenfolge vertauscht, lässt Code, Antwortart und Exit-Code gleich. Invariante: beide Pfade enden in `fail(CodeUnsupported, …, "Startnachricht nicht lesbar")`. Ort: `server.go`, `cancelOhneSchluessel` und der Fehlerpfad nach `ReceiveStartupMessage`. Sonde (`/tmp/probe/probe3.sh` gegen das gebaute Binary, `record`): Länge 13 mit vollem Rest endet mit `…Startnachricht nicht lesbar: cancel request too short`, Länge 13 ohne Rest mit `…Startnachricht nicht lesbar: EOF`, ebenso Länge 12 ohne Rest und Länge 269 ohne Rest (`EOF`) gegen `cancel request too long` mit vollem Rest. **Die Sonde trennt sie im Text der Meldung**, nicht im Code: Bei Länge 13 und fehlendem Rest nennt der Recorder (Bibliothek) den Rest, nicht die Länge, entgegen der Reihenfolge der Spezifikation. Der Text gehört nicht zum Vertrag (Code, `ErrorResponse`, Exit-Code); der Punkt geht als Frage an den Architect (siehe Bericht), kein Code ist darauf entschieden.

**Probe-Tabelle Handbuch und README** (Schritt 17; Binary: `pgwire-recorder:dev` aus `make build`, extrahiert nach `/tmp/probe/`, Probe `/tmp/probe/probe.sh` mit `record --upstream 127.0.0.1:1` (nichts lauscht) und `replay --input test/integration/testdata/ende-terminate.yaml`, je ohne und mit `--tls-cert`, rohe Bytes über `/dev/tcp`; Unit-Test `T` = `TestCancelRequestLaenge` usw.):

| Satz (Handbuch, Bullet *Abbruchwunsch*; README gleich) | Probe | Ergebnis |
|---|---|---|
| leitet den Abbruchwunsch nicht an die Datenbank weiter | `record` mit unerreichbarem Upstream, Längen 12, 16, 268 | Exit 0, 3× `PGR-W3001`, kein `PGR-E4002`; Unit: `rec.opened == 0` (M14) |
| schließt seine Verbindung mit `PGR-W3001` | Probe, Unit mit Lesen bis zum Ende ohne Halbschließen | 3× `PGR-W3001`, keine Antwortbytes; M6 rot |
| ohne Fehler, Exit-Code unverändert | Probe `zugesagt` in 4 Läufen (`record`/`replay`, ohne/mit TLS); Integrationstest | Exit 0; M18 rot; `TestE2ERecordCancelOhneSchluessel` grün |
| andere Verbindungen laufen unberührt weiter | `TestCancelRequestZweiteVerbindungUnberuehrt` (Längen 12 und 16, beide Modi) | grün; M13 rot |
| öffnet keine Sitzung | Unit, alle Fälle `sessions == 0` | grün; M14 rot |
| ohne Schlüssel (12) und mit Schlüssel 4 bis 256 (16 bis 268) | Probe 12, 16, 268; Unit 12, 16, 17, 20, 268 | `PGR-W3001`; M1, M2h, M2i rot |
| in `record` und `replay` | Probe, Unit | gleich; M3, M3b rot |
| als erste Nachricht und nach einer Anfrage nach Verschlüsselung | Unit, 6 Stellen | grün; M11, M12 rot |
| mit `--tls-cert` auch unverschlüsselt, ohne `PGR-E6003` | Probe mit `--tls-cert` (Klartext-Verbindung), Unit „unverschlüsselt bei TLS-Pflicht“ | 3× `PGR-W3001`, kein `PGR-E6003`; M4 rot |
| jede andere Länge von 8 bis 10000 (8 bis 11, 13 bis 15, 269 bis 10000) ist `PGR-E6001`, Exit-Code 6 | Probe `E6001` mit 8, 9, 10, 11, 13, 14, 15, 269, 10000; Unit | 9× `PGR-E6001`, Exit 6; M2a, M2c bis M2g rot |
| ein unvollständiger Rest ist `PGR-E6001` | Probe `rest-fehlt` (Länge 12, 2 von 4 Bytes, Client schließt); Unit 0, 1, 3 Bytes | 1× `PGR-E6001`, Exit 6; M5 rot |
| unter 8 oder über 10000 ist `PGR-W3003` | Probe `W3003` mit 7 und 10001; Unit 0, 4, 7, 10001 | 2× `PGR-W3003`, Exit 0; M9, M9b rot |
| Tabellenzeile `PGR-E6001` nennt Längen und unvollständigen Rest | wie die beiden Zeilen darüber | siehe dort |
| Tabellenzeile `PGR-W3001` „auch bei einem Abbruchwunsch ohne Schlüssel“ | wie die Zeile für 12 | siehe dort |
| `replay` beginnt jede Verbindung, die es bedient, mit `AuthenticationOk` (einen Abbruchwunsch bedient es nicht) | Unit `replay=true`: keine Antwortbytes bei Länge 12; Bestand `TestE2EReplayAuswahl` für den ersten Teil | grün |

**Treffer der Durchsicht nach Quantoren** (Schritt 17, Merkmal „jede“, „alle“, „keine“, „nie“ mit Verbindung, Warnung, Exit-Code, Startnachricht, Sitzung; je Zeile und Absatz in `docs/user/benutzerhandbuch.md` und `README.md`):

- Handbuch *`replay` beginnt jede angenommene Verbindung mit `AuthenticationOk`* und README, gleicher Satz: eng gefasst (*die es bedient (einen Abbruchwunsch bedient es nicht)*); beide Orte gleich.
- Handbuch *bedient eine angenommene Verbindung mit jedem Benutzer und jeder Datenbank* (Abschnitt zu Sicherheit): unverändert; der Satz sagt, dass Benutzer und Datenbank nicht geprüft werden, nicht, dass jede Verbindung bedient wird; „bedient“ setzt die Startnachricht voraus.
- Handbuch *ein Abbruchwunsch … zählen bei keiner der beiden Zuordnungen als Verbindung* (Zählung bei `replay`): unverändert und belegt (Unit: `geoeffnet == 0` für alle Längen in `replay`).
- Handbuch *Eine Verbindung zählt erst mit ihrer Startnachricht … eine mit nur einem Abbruchwunsch (`PGR-W3001`)*: unverändert; die Längen mit `PGR-E6001` stehen im Bullet darunter (Bruch mitten in der Startnachricht) und im neuen.
- Handbuch *Beim Ablauf beendet das Werkzeug jede noch laufende Verbindung*: unverändert; Zwangsende nach `--shutdown-timeout`, ein anderer Vorgang als der `CancelRequest`.
- Handbuch Exit-Code-Tabelle (Zeile 6 *nicht unterstützte Funktion*) und *Ein Fehler, der nur eine Verbindung betrifft*: unverändert; `PGR-E6001` gilt weiter.
- Handbuch *Exit-Code ändert sich nicht* bei `PGR-W3002` (TLS) und bei der Anmeldung: unverändert, andere Warnung bzw. anderer Vorgang.
- Handbuch Tabellenzeilen `PGR-E6001` und `PGR-W3001`: verbunden (ergänzt, siehe oben).

Für den neuen Satz *Jede andere Länge von 8 bis 10000 Bytes … ist `PGR-E6001`* zählt dieselbe Seite die Ausnahmen für denselben Gegenstand auf: Länge unter 8 oder über 10000 (`PGR-W3003`, im Satz genannt, Probe `W3003`) und ein unvollständiger Rest (im Satz genannt, Probe `rest-fehlt`); der README-Satz nennt beide ebenfalls.

**Läufe.**

- Engster Sensor: `/tmp/fast.sh ./internal/adapters/driving/pgwire` (Bind-Mount): grün, 10,9 s, vor und nach dem Umbau von `startup` in `startkopf`.
- `make lint` nach dem Umbau: Exit 0 (der erste Stand von `startup` war über `gocognit` 20, `cyclop` 15; der Cancel-Zweig steht deshalb in `startkopf`).
- `make test-integration` mit `PGR_TESTS_1=TestE2ERecordCancelOhneSchluessel`: grün (0,73 s), Mutant M1 im Binary: rot (Exit-Code 6).
- `make abdeckung`: schrieb `abdeckung-e2e.md` und `abdeckung-unit.md`.
- Metrik der Marke (§4): `git diff --numstat f7aeea92..f3ba1835 -- internal test docs/user README.md`: 340 Einfügungen und 14 Löschungen, Summe **354**, Schätzung 350, Marke 700.
- `make gates` auf f3ba1835: Exit 0, ein Lauf, keine Wiederholung nötig (`TestE2EPlayVergleichAbweichung` grün); Abdeckung der Unit-Tests 94,43 % bei Schwelle 94 %; `test-integration` und `test-container` grün; `lint` grün. Die Gegenproben liefen nicht: Dieser Slice ändert weder ein Gate noch einen Vertrag eines Werkzeugs (`AGENTS.md` §6, Schritt 6); die Abdeckungstabellen schrieb `make abdeckung`, `make gates` prüft sie über `abdeckung-check`.

**Beobachtungen.**

- Der Mutant „Untergrenze 11, liest 4 Bytes“ hängt (der Server wartet auf das vierte Byte), die Tests färben ihn über die Frist des Clients (10 s, „Verbindung nicht geschlossen“) statt über die Antwort; ein Mutant, der `laenge-8` Bytes liest (M2d), färbt über die Antwort. Beide stehen in der Tabelle.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: in diesem Slice nicht wieder aufgetreten (je Zusage ein Fall, die Mutationen oben vom Implementer selbst gefahren).

**Proben zum Nachfahren** (Binary aus `make build`: `id=$(docker create pgwire-recorder:dev); docker cp $id:/pgwire-recorder /tmp/probe/pgwire-recorder`; die Aufnahme `ende-terminate.yaml` liegt im Repo; `socat` und `openssl` im Dev-Container):

`/tmp/probe/probe.sh`:

```bash
#!/bin/bash
# probe.sh — CancelRequest gegen das gebaute Binary (record, replay; ohne und mit --tls-cert)
B=/tmp/probe/pgwire-recorder
W=/tmp/probe/w; rm -rf $W; mkdir -p $W; cd $W
openssl req -x509 -newkey rsa:2048 -nodes -keyout k.pem -out c.pem -days 1 -subj /CN=127.0.0.1 >/dev/null 2>&1
be32() { printf "\\x$(printf %02x $(($1>>24&255)))\\x$(printf %02x $(($1>>16&255)))\\x$(printf %02x $(($1>>8&255)))\\x$(printf %02x $(($1&255)))"; }
cancel() { # laenge geliefert
  be32 $1; be32 80877102; head -c "$2" /dev/zero
}
sende() { # port laenge geliefert -> Ausgabe: Zahl der Antwort-Bytes und ob die Verbindung geschlossen wurde
  local port=$1
  exec 3<>/dev/tcp/127.0.0.1/$port
  cancel $2 $3 >&3
  [ "$4" = halb ] && true
  local n; n=$(timeout 5 cat <&3 | wc -c); local rc=${PIPESTATUS[0]}
  exec 3<&- 3>&-
  echo "len=$2 rest=$3 antwort_bytes=$n"
}
lauf() { # modus tls name laengen...
  local modus=$1 tls=$2 name=$3; shift 3
  local port=$((15000 + RANDOM % 1000)) args
  args=(--listen 127.0.0.1:$port)
  if [ $modus = record ]; then args+=(--upstream 127.0.0.1:1 --output $W/o-$RANDOM.yaml); else args+=(--input /workspaces/pgr-dev2/pgwire-recorder/test/integration/testdata/ende-terminate.yaml); fi
  [ $tls = 1 ] && args+=(--tls-cert $W/c.pem --tls-key $W/k.pem)
  $B $modus "${args[@]}" 2>$W/err & local pid=$!
  for _ in $(seq 50); do (exec 3<>/dev/tcp/127.0.0.1/$port) 2>/dev/null && break; sleep 0.1; done
  echo "== $modus tls=$tls $name"
  for l in "$@"; do if [ $l = trunc ]; then sende $port 12 2; else sende $port $l $((l-8)); fi; done
  kill -TERM $pid; wait $pid; echo "exit=$?"
  grep -o 'level=[A-Z]* .*code=PGR-[EW][0-9]*' $W/err | sed 's/ time=[^ ]*//; s/msg=.*code=/code=/' | sort | uniq -c
}
for modus in record replay; do
  for tls in 0 1; do
    lauf $modus $tls zugesagt 12 16 268
    lauf $modus $tls E6001 8 9 10 11 13 14 15 269 10000
    lauf $modus $tls W3003 7 10001
  done
done
lauf record 0 rest-fehlt trunc
```

`/tmp/probe/probe3.sh` (Text der Meldung bei Länge 13, 12, 16 und 269 mit und ohne Rest):

```bash
#!/bin/bash
B=/tmp/probe/pgwire-recorder; W=/tmp/probe/w3; rm -rf $W; mkdir -p $W; cd $W
be32() { printf "\\x$(printf %02x $(($1>>24&255)))\\x$(printf %02x $(($1>>16&255)))\\x$(printf %02x $(($1>>8&255)))\\x$(printf %02x $(($1&255)))"; }
port=$((16000 + RANDOM % 1000))
$B record --listen 127.0.0.1:$port --upstream 127.0.0.1:1 --output $W/o.yaml 2>$W/err & pid=$!
sleep 1
t() { printf "%s: " "$1"; { be32 $2; be32 80877102; head -c $3 /dev/zero; } | socat -t2 - TCP:127.0.0.1:$port | strings | tr '\n' ' '; echo; }
t "13 voll" 13 5
t "13 Rest fehlt" 13 0
t "13 Rest 2 von 5" 13 2
t "12 Rest fehlt" 12 0
t "12 Rest 2 von 4" 12 2
t "16 Rest 2 von 8" 16 2
t "269 voll" 269 261
t "269 Rest fehlt" 269 0
kill -TERM $pid; wait $pid; echo exit=$?
```

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — der Fall entstand, weil der Test nur den `CancelRequest` mit Schlüssel prüfte.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (PGWire-Adapter, Dokumentation). Nachgezählt beim Eintragen der Randformen (2026-10-20, Architect, `AGENTS.md` §3.13): zwei Liefer-Punkte (`LH-FA-05`, `LH-FA-13`), zwei Schichten; der Bootstrap bleibt unberührt, der Integrationstest zählt als Test des Binaries; die Spezifikation trägt der Architect vor dem Code, sie zählt nicht in die Marke.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
