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
- Die Längen 8 bis 11, 13 bis 15 und 269 bis zur Höchstlänge bleiben `PGR-E6001` (der Adapter lehnt sie jetzt aus dem Kopf ab, Zeitpunkt vor dem Lesen des Rests), die Längen 16 bis 268 bleiben bei der Bibliothek — Bestand bleibt bewusst stehen: Kein Client sendet sie (die Spezifikation entscheidet sie jetzt), und eine eigene Auswertung brächte kein Verhalten, das ein Test heute verlangt.
- Lesen oder Verwerfen des ungelesenen Rests bei einer Ablehnung aus dem Kopf — die Zustellung der `ErrorResponse` ist dann nicht zugesagt (Spezifikation `LH-FA-05.e`, geprüft nur unter Linux), außer bei den Längen 9 bis 15, wenn der Client genau die Bytes der Länge (höchstens 7) mit dem Kopf sendet, die der Recorder mit dem Kopf liest; mit mehr gesendeten Bytes ist sie nicht zugesagt (Grenze ohne Test); der Recorder würde sonst auf Bytes eines Clients warten, den er gerade ablehnt, und eine Grenze für das Lesen wäre eine neue Entscheidung ohne Anforderung dahinter.
- Das Verbindungsende einer Hauptverbindung, die ein Client bei Kontext-Abbruch mitten in einer Anfrage schließt — bleibt `PGR-E4003` (Exit-Code `4`) nach dem Abschnitt *Verbindungsende* der Spezifikation; es ist der Abbruch der Anfrage, nicht der `CancelRequest`, und eine andere Anforderung wäre ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Ein `CancelRequest` ohne Schlüssel (Länge 12) schließt die Verbindung mit `PGR-W3001` ohne Verbindungsfehler, in `record` und `replay` (Unit-Test mit 12 Byte). Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung).
- [x] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): pgx mit Kontext-Abbruch über `record` lässt den Lauf mit Exit-Code 0 enden, wenn sonst kein Fehler auftrat (Integrationstest); der Abbruch trifft eine Anfrage, die der Server vor der Frist des Clients beantwortet, sodass die Hauptverbindung regulär endet (Szenario in §3).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (Review: [Bericht](../../../reviews/2026-10-20-review-slice-v1-abschluss-cancel-ohne-schluessel.md), Verifikation: [Bericht](../../../reviews/2026-10-20-verifikation-slice-v1-abschluss-cancel-ohne-schluessel.md), Nach-Verifikation: [Bericht](../../../reviews/2026-10-20-nachverifikation-slice-v1-abschluss-cancel-ohne-schluessel.md))
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/pgwire` (`startkopf`, `cancelOhneSchluessel`; die Komplexitätsgrenze von `make lint` hat den Zweig aus `startup` in `startkopf` gelegt) | update | Der Adapter prüft die Länge des `CancelRequest` aus dem Kopf (8 Bytes, die er schon gelesen hat), vor dem Lesen des Rests: Länge 12 liest er selbst (die Prozess-ID, 4 Bytes; die Bibliothek lehnt 12 ab), meldet `PGR-W3001` und schließt, fehlen die 4 Bytes, ist es `PGR-E6001`; jede Länge außer 12 und 16 bis 268 (8 bis 11, 13 bis 15, 269 bis zur Höchstlänge) lehnt er sofort mit `PGR-E6001` und `ErrorResponse` ab, ohne den Rest abzuwarten (gemessen: eine Probe mit Kopf der Länge 9, 13, 15 und 269 und offener Verbindung ohne Rest, ohne Prüfung im Kopf ohne Antwort, mit ihr sofort `ErrorResponse`, in beiden Modi, auch über TLS; ohne Fork der Bibliothek); 16 bis 268 liest die Bibliothek (`PGR-W3001`, ein fehlender Rest wartet bis zum Ende der Verbindung) |
| `internal/adapters/driving/pgwire/server_cancel_test.go` (neu), `server_protokollrand_test.go` (`randServer` hält das Log) | update | Negative nach `LH-FA-05`: ein Fall je Länge 8 bis 11, 13 bis 15, 269, 270 und 10000 (`PGR-E6001`) und 12, 16, 17, 20 und 268 (`PGR-W3001`); Länge 12, 16 und 268 mit 0, 1 und 3 Bytes Rest und geschlossener Schreibseite des Clients (`PGR-E6001`); Länge 9, 11, 13, 15, 269 und 10000 mit nur dem Kopf und **offener** Verbindung: `ErrorResponse` und Schließen vor einer Frist (Reihenfolge Länge vor Rest, Zeitpunkt der Ablehnung); Länge 12, 16 und 268 mit fehlendem Rest und offener Verbindung: keine Antwort, bis der Client schließt (Grenze ohne Frist); Prozess-ID und Schlüssel 0, `0xff`, `0x2a`; Länge 0, 4, 7 und 10001 mit dem Code des `CancelRequest` (`PGR-W3003`, Reihenfolge der Gründe); jeweils in `record` und `replay` und an sechs Stellen (erste Nachricht, nach `SSLRequest` mit `N`, nach `GSSENCRequest`, über TLS nach `S`, unverschlüsselt bei TLS-Pflicht, nach `GSSENCRequest` bei TLS-Pflicht); eine zweite Verbindung bleibt unberührt und die nächste öffnet die zweite Session; ein Fall im `replay` mit begonnener Extended-Folge der zweiten Verbindung (F-V1; Länge 12 und 16: Parse, Bind, Execute ohne `Sync`, `CancelRequest` auf anderer Verbindung, danach `Sync` und die vollständigen aufgezeichneten Antworten, der Use Case sieht genau die vier Nachrichten und beendet, schließt oder öffnet nichts); die Fälle mit vollem Rest (269, 270, 10000) prüfen Ende und Log, nicht die `ErrorResponse` (F-V3), der Fall `cancelSenden` ist damit gegen `connection reset` stabil Die Fälle der Längen 9, 10, 11, 13, 14, 15 mit genau den Bytes der Länge (Länge minus 8, bis 7 Bytes, mit dem Kopf gesendet) verlangen die `ErrorResponse` (Spezifikation `LH-FA-05.e`); die Fälle 269, 270 und 10000 verlangen sie nicht; bei Länge 9 bis 15 mit mehr gesendeten Bytes gibt es keinen Fall, die Spezifikation sagt sie nicht zu. |
| `test/integration/cancel_e2e_test.go` (neu) | update | pgx mit Kontext-Abbruch, Exit-Code: Der Kontext einer Anfrage mit `pg_sleep(0.5)` endet nach 100 ms; `BuildContextWatcherHandler` setzt `CancelRequestContextWatcherHandler` mit `DeadlineDelay` 20 s (der Standard von pgx sendet bei laufender Anfrage keinen `CancelRequest`); die Anfrage endet regulär, die Hauptverbindung sendet `Terminate`, `SIGTERM` an den Recorder, Exit-Code 0, im Log genau einmal `PGR-W3001` und weder `PGR-E6001` noch `PGR-E4003` |
| `docs/user/benutzerhandbuch.md`, `README.md`, `docs/user/abdeckung-e2e.md`, `docs/user/abdeckung-unit.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung); die Abdeckungstabellen schreibt `make abdeckung` |

**Zusagen und Mutationen** (`AGENTS.md` §3.10; der Bericht nennt je Zeile den roten Test): Länge 12 ist `PGR-W3001` (Mutation: Zweig für 12 entfernt, oder 12 bis 16 angenommen) · Länge 11 und 13 sind `PGR-E6001` (Mutation: Untergrenze 11 oder Bereich 12 bis 15) · 16 und 268 sind `PGR-W3001`, 269 `PGR-E6001` (Fälle pinnen die Bibliothek) · gleich in `record` und `replay` (Mutation: Zweig nur im `record`) · mit TLS-Pflicht kein `PGR-E6003` (Mutation: Prüfung der Pflicht vor dem `CancelRequest`) · nach `SSLRequest` mit `N` und über TLS nach `S` · Rest nicht geliefert ist `PGR-E6001` (Mutation: Lesefehler verworfen) · die Länge wird vor dem Lesen des Rests abgelehnt (Mutation: Prüfung im Kopf entfernt, sodass die Bibliothek den Rest zuerst liest; rot im Fall mit offener Verbindung, der Mutant ist nicht äquivalent, weil er den Zeitpunkt der Ablehnung ändert) · die Verbindung wird geschlossen (Mutation: sie bleibt offen) · Exit-Code 0 (Integrationstest, Mutation: Länge 12 wie bisher abgelehnt, Lauf endet mit 6) · eine laufende Interaktion der zweiten Verbindung bleibt unberührt, auch im `replay` (Extended-Folge ohne `Sync` offen, `CancelRequest` auf anderer Verbindung, danach `Sync`: vollständige aufgezeichnete Antworten; Mutation: der Zweig beendet im `replay` alle Verbindungen oder setzt den Cursor zurück) · Länge 269 bis 10000 mit Rest: Verbindung geschlossen und `PGR-E6001` im Log, die `ErrorResponse` wird nicht verlangt (Mutation: Schließen unterbleibt; der Test liest bis zum Ende oder Rücksetzen und prüft nur Ende und Log); Länge 9, 11, 13, 15, 269, 10000 ohne Rest verlangen die `ErrorResponse` weiter.

**Größe** (`AGENTS.md` §3.14): Schätzung in der Metrik von §4, mit dem Doppelten für Tests: Produktcode rund 25, Unit-Tests rund 200, Integrationstest rund 100, Handbuch und README rund 30, zusammen **rund 350**. Gemessen (Basis `f7aeea92`, Stand `0ee53a97`, `git diff --numstat f7aeea92..0ee53a97 -- internal test docs/user README.md`): 560 (546 Einfügungen, 14 Löschungen), Marke 700; die Nacharbeit nach der Verifikation (F-V1, F-V3) fügte 83 hinzu (477 auf 560), Handbuch und README blieben unverändert. Zwei Liefer-Punkte, zwei Schichten (§8).

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
- Länge 8 bis 11, 13 bis 15, 269 bis zur Höchstlänge — `PGR-E6001` mit `ErrorResponse` (Exit-Code 6), Punkt; die `ErrorResponse` ist nur ohne ungelesene Bytes zugesagt, mit Rest gilt die Randform *Ablehnung aus dem Kopf bei ungelesenem Rest* weiter unten (Schließen, `PGR-E6001` und Exit-Code 6); unter 8 und über der Höchstlänge `PGR-W3003`, Tabelle *Protokollrand*.
- Reihenfolge der Gründe: Länge außerhalb von 8 bis zur Höchstlänge, dann Länge des `CancelRequest`, dann Rest nicht geliefert — Punkt; die ersten beiden stehen im Kopf und werden sofort abgelehnt, auch bei offener Verbindung ohne Rest, der dritte ist erst mit dem Ende der Verbindung feststellbar (Entscheidung des Architect nach dem Review, Option: Code prüft vor dem Lesen des Rests).
- Rest der Nachricht nicht geliefert (Abbruch, Ende mitten in der Nachricht) — `PGR-E6001`, Punkt und Zeile *Startnachricht mit vollständigem Kopf*; Verbindung offen und kein Rest bei Länge 12 bis 268: der Recorder wartet auf Rest oder Ende, eine Frist gibt es nicht (Grenze im Punkt; ein Client, der so wartet, hält nur seine eigene Verbindung).
- Stelle: erste Nachricht und nach `SSLRequest` (`N`, oder `S` und TLS) oder `GSSENCRequest` — Punkt; mit TLS-Pflicht ohne `PGR-E6003` — Abschnitt *Unverschlüsselte Clients* von `LH-FA-23.a`.
- Modus: `record` und `replay` gleich, keine Session, keine Aufzeichnung, keine Verbindung zum Server, eine zweite Verbindung unberührt — Punkt und *Modus*.
- Prozess-ID und Schlüssel — nicht ausgewertet, Punkt.
- Exit-Code — die Warnung ändert ihn nicht (`LH-FA-13`); eine Hauptverbindung, die der Client bei Kontext-Abbruch mitten in der Anfrage schließt, ist `PGR-E4003` (Exit-Code 4) nach *Verbindungsende*, Punkt.
- Laufende Interaktion der zweiten Verbindung (Entscheidung des Architect nach der Verifikation, F-V1): `record` — Anfrage gesendet, `ReadyForQuery` fehlt; `replay` — Extended-Folge bis vor `Sync` begonnen; in beiden Modi bekommt die zweite Verbindung danach die vollständigen Antworten, im `replay` die aufgezeichneten; beide Modi bekommen einen Fall, Punkt.
- Ablehnung aus dem Kopf bei ungelesenem Rest (Entscheidung des Architect nach der Verifikation, F-V3): kein Code zum Lesen oder Verwerfen des Rests; die `ErrorResponse` ist nur ohne ungelesene Bytes zugesagt, mit Rest sind Schließen, `PGR-E6001` und Exit-Code 6 zugesagt, Punkt. Ausnahme (Entscheidung des Architect, Rückgabe des Implementers): Länge 9, 10, 11, 13, 14, 15 mit genau den Bytes der Länge (höchstens 7), in einer Sendung mit dem Kopf — dort ist die `ErrorResponse` zugesagt, je Länge ein Fall in `TestCancelRequestLaenge` (vorhanden); „Rest“ sind die gesendeten Bytes nach dem Kopf (Entscheidung des Architect, Befund L-2): mehr gesendete Bytes bei Länge 9 bis 15 und der Rest der Längen 269 bis 10000 bleiben unzugesagt, ohne Fall (Grenze, kein Test).
- Plattform (`LH-QA-03`): Linux im Container geprüft (Unit- und Integrationstest); macOS und Windows — Grenze an der Zusage der Zustellung der `ErrorResponse` bei ungelesenem Rest, im Punkt als „geprüft nur unter Linux“; sonst hängt die Eigenschaft nicht vom Betriebssystem ab (Bytes auf einer TCP-Verbindung).
- Aufruf-Voraussetzung und Aufrufarten des Treibers: nicht zutreffend (keine Ordnung von Werten über Aufrufe, kein Stand, kein Nachweis).

Risiken:

- Ein Client sendet einen `CancelRequest` in einer Länge, die die Spezifikation als `PGR-E6001` führt (zum Beispiel 13 bis 15 oder über 268) — **Ausgang (Closure):** weiter offen → `BEO-REPO/cancelrequest-laenge-ohne-client-beleg` (1×): Die Behandlung ist je Länge entschieden und von Fällen gepinnt, ob ein Client eine dieser Längen sendet, belegt kein Test und keine Probe.
- Der Integrationstest trifft die Zeitfolge (Anfrage endet vor der Frist des Clients) nicht stabil — **Ausgang:** offen bis Closure; Gegenmittel: die Frist des Clients um ein Vielfaches über der Dauer der Anfrage. **Ausgang (Closure):** entfallen — die Frist ist 20 s gegen eine Anfrage von 0,5 s, der Test wird nur rot, wenn der Client die Antwort 20 s nach der Anfrage nicht liest; `TestE2ERecordCancelOhneSchluessel` lief grün im Lauf des Implementers, im Lauf des Verifiers (0,73 s) und in jedem `make gates` des Slice (`f3ba1835`, `cb8dfd49`, `c30b830a`, je ein Lauf ohne Wiederholung), ohne Rot. Nicht unter künstlicher Last gefahren; tritt ein Rot auf, ist es ein Auftreten von `BEO-REPO/test-faellt-unter-cpu-last-aus`.

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
| Länge aus dem Kopf entschieden, vor dem Lesen des Rests (Länge 9, 11, 13, 15, 269, 10000 und je Element 8, 10, 14, 270) bei offener Verbindung: `ErrorResponse` und Schließen | K0: Prüfung im Kopf entfernt (die Bibliothek liest den Rest zuerst); K9, K10, K11, K13, K14, K15, K269, K270, K10000: die Prüfung lässt genau diese Länge durch | `TestCancelRequestKopf/…/Länge_<n>_nur_Kopf,_Verbindung_offen` (K0: 108, je Länge 12; Meldung „keine Ablehnung aus dem Kopf“ nach 2 s) |
| Länge 12, 16, 17 und 268 mit fehlendem Rest bei offener Verbindung: keine Antwort bis zum Ende der Verbindung (Grenze ohne Frist) | W16, W17, W268: genau diese Länge wird aus dem Kopf abgelehnt; W12: Länge 12 wird vor dem Zweig der Länge 12 aus dem Kopf abgelehnt; K12: Länge 12 sofort `PGR-E6001` | `TestCancelRequestKopf/…/Länge_<n>,_<g>_Bytes_Rest,_Verbindung_offen` („Antwort vor dem Ende der Verbindung“), dazu `TestCancelRequestLaenge` |
| eine laufende Interaktion einer zweiten Verbindung wird nicht abgebrochen, `record` | M13: der Zweig beendet alle Verbindungen (neuer Stand: mit laufender Anfrage) | `TestCancelRequestLaufendeInteraktionUnberuehrt/Länge_12`, `_16` |
| dasselbe im `replay` (Extended-Folge bis vor `Sync`, `CancelRequest` auf anderer Verbindung, danach `Sync`: vollständige aufgezeichnete Antworten; F-V1) | MV1: `cancelGeschlossen` ruft `Zwangsende()` (beendet alle Verbindungen, `Forced` und `CloseConnection` am Use Case); MV2: `cancelGeschlossen` ruft im `replay` `CloseConnection` (setzt die Session zurück); MV3: `cancelGeschlossen` ruft `OpenConnection` (öffnet eine Session) | `TestCancelRequestLaufendeInteraktionUnberuehrtReplay/Länge_12`, `_16` (je Mutant 2): „nach dem CancelRequest: geöffnet 1, geschlossen 1, erzwungen 1, Shutdown 0; erwartet 1, 0, 0, 0“ (MV1), „geschlossen 1“ (MV2), „geöffnet 2“ (MV3); MV1 färbt zusätzlich `TestCancelRequestLaufendeInteraktionUnberuehrt/Länge_12`, `_16` (`record`) |
| Länge 269, 270 und 10000 mit vollem Rest: Verbindung geschlossen, `PGR-E6001` im Log und als gemerkter Fehler, keine Warnung `PGR-W3001`, keine Session; die `ErrorResponse` wird nicht verlangt (F-V3) | MV4: ab Länge 269 bleibt die Verbindung nach der Fehlerantwort 30 s offen; MV5: ab Länge 269 weder `ErrorResponse` noch Eintrag noch gemerkter Fehler; MV6: ab Länge 269 Warnung `PGR-W3001` statt Fehler; M2f und M2g wie oben (Länge 269 und 10000 als Länge 12 behandelt) | `TestCancelRequestLaenge/replay=false/erste_Nachricht/Länge_269`, `_270`, `_10000`: MV4 „Verbindung nicht geschlossen“ (3, je 10 s Frist), MV5 und MV6 „erwartet PGR-E6001 im Log und als gemerkter Fehler“ (3 je Mutant); M2f und M2g je 12 (alle Stellen, beide Modi) |
| Länge 9, 11, 13, 15, 269, 10000 mit nur dem Kopf: `ErrorResponse` bleibt verlangt (F-V3) | MV7: bei Länge außerhalb der Zusage nur gemerkt, keine `ErrorResponse` | `TestCancelRequestKopf/replay=false/erste_Nachricht/Länge_9`, `_13`, `_269`, `_10000` (`…_nur_Kopf,_Verbindung_offen`) |
| keine Session | M14: der Zweig öffnet eine Session | `…/Länge_12` (74) |
| eine zweite Verbindung bleibt unberührt, die laufende Anfrage wird nicht abgebrochen | M13: der Zweig beendet alle Verbindungen | `TestCancelRequestZweiteVerbindungUnberuehrt/…` |
| Exit-Code 0 (`LH-FA-13`) | M1 im Binary: Länge 12 wie bisher abgelehnt | `TestE2ERecordCancelOhneSchluessel`: „Exit-Code 6, erwartet 0“ |

Die Mutanten M2f bis M2i ersetzen das Verhalten der Bibliothek (v5.11.0) durch einen eigenen Zweig und zeigen, dass die Fälle für 16, 268, 269 und 10000 die Längen pinnen; ein Mutant im Quelltext der Bibliothek ist nicht möglich. Die Schlüssel von 4 bis 256 Bytes (16 bis 268) wertet allein die Bibliothek aus; dass der Recorder sie nicht auswertet, belegt der Fall `Länge_16,_Schlüssel_0xff` nur gegen die Bibliothek in dieser Version.

**Grüne Mutanten.** Einer, äquivalent mit Invariante, Ort und Sonde:

- K8: Länge 8 wird nicht aus dem Kopf abgelehnt (die Bibliothek entscheidet). Invariante: Eine Nachricht der Länge 8 besteht nur aus dem Kopf, es gibt keinen Rest, auf den der Recorder warten könnte; die Bibliothek liest 0 Bytes und lehnt „cancel request too short“ sofort ab. Ort: das Nachrichtenformat selbst (Länge zählt Kopf und Rest) und `ReceiveStartupMessage` ohne weiteres Lesen. Sonde: der Fall `Länge_8_nur_Kopf,_Verbindung_offen` mit dem Mutanten K8 ist grün (Lauf des Pakets: `rc=0`), also trennt auch der Fall mit offener Verbindung Original und Mutant nicht; der Fall bleibt als Pin im Test.
- Der frühere Mutant „Reihenfolge *Länge des `CancelRequest`* und *Rest nicht geliefert* vertauscht“ ist **nicht äquivalent** (Review F-1, F-2): Er ändert den Zeitpunkt der Ablehnung. Er ist K0 und für Länge 9 bis 15, 269 bis 10000 die Mutanten K9 bis K10000; die Sonde des Reviewers (Länge 13, nur der Kopf, Verbindung offen) ist der Fall `Länge_13_nur_Kopf,_Verbindung_offen`, rot unter K0 und K13. Die Sonde gegen das gebaute Binary (`/tmp/probe/probe4.sh`, unten): Länge 9, 11, 13, 15, 269 und 10000 mit nur dem Kopf werden bei offener Verbindung sofort mit `PGR-E6001` und `ErrorResponse` beantwortet, Länge 12, 16 und 268 antworten nicht.

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

**Läufe der Nacharbeit nach dem Review** (Basis `arch-cancel2`, Commit `cb8dfd49`):

- Engster Sensor: Bind-Mount-Lauf des Pakets `pgwire`: grün (21 s); `make lint`: Exit 0; `make abdeckung` schrieb `abdeckung-unit.md`.
- Mutanten K0, K8 bis K10000, W12 bis W268, K12: siehe Tabelle und Grüne Mutanten; `server.go` danach byte-gleich (`cmp`).
- Metrik der Marke: `git diff --numstat f7aeea92..cb8dfd49 -- internal test docs/user README.md`: 463 Einfügungen, 14 Löschungen, Summe **477** (Marke 700).
- `make gates` auf `cb8dfd49`: Exit 0, ein Lauf, kein Flake; Abdeckung der Unit-Tests 94,41 % bei Schwelle 94 %; alle Gates der Kette grün. Gegenproben nicht gefahren: kein Gate und kein Vertrag eines Werkzeugs geändert.
- F-4: `TestCancelRequestLaufendeInteraktionUnberuehrt` (record, Längen 12 und 16) hält die Anfrage der zweiten Verbindung im Recorder (`queryHalt`), sendet den `CancelRequest`, lässt die Anfrage danach enden; die Deklaration von `TestCancelRequestZweiteVerbindungUnberuehrt` sagt nur noch „sendet nach dem CancelRequest ihre Anfrage“.

**Läufe der Nacharbeit nach der Verifikation** (Zweig `impl-cancel3`, Basis `arch-cancel4` = `5184dbeb`, Code-Commit `0ee53a97`):

- F-V1: Die Grenze im `replay` entfällt. `TestCancelRequestLaufendeInteraktionUnberuehrtReplay` (Länge 12 und 16) sendet Parse, Bind, Execute ohne `Sync` auf der zweiten Verbindung, wartet, bis der Use Case drei Nachrichten sieht, sendet den `CancelRequest` auf einer anderen Verbindung, prüft am Fake (`geoeffnet` 1, `closed`, `forced`, `shutdowns` je 0), sendet dann `Sync` und `Terminate` und verlangt in der Antwort die Folge `1`, `2`, `D`, `C`, `Z` am Ende ohne `E`, dazu genau Parse, Bind, Execute, Sync beim Use Case. Der Fake hält keinen Cursor; geprüft sind seine Zähler `geoeffnet`, `closed`, `forced` und `shutdowns` (MV1 bis MV3) und genau die vier Nachrichten `Parse`, `Bind`, `Execute`, `Sync` beim Use Case. Aufrufe von `Query` und `Sent` am Use Case zählt der Test nicht: Mutanten, die sie in `cancelGeschlossen` im `replay` aufrufen, bleiben grün (Nach-Verifikation L-1), und „ohne Veränderung ihres Cursors“ aus der Spezifikation (`LH-FA-05.e`) ist damit von keinem Test geprüft (Rückgabe an den Architect, Closure-Notiz). Mutanten MV1 bis MV3 in der Tabelle, rot gesehen. 20 Wiederholungen des Tests grün.
- F-V3: Reproduziert am Stand vor der Nacharbeit: `-run 'TestCancelRequestLaenge$/.*/.*/L.nge_(269|270|10000)$' -count=100`, dreimal; ein Lauf rot mit `replay=false/über_TLS_nach_S/Länge_10000: write tcp …: write: connection reset by peer` (Schreiben der 10000 Bytes, nicht das Lesen der Antwort). Nach der Nacharbeit fünfmal 100 Läufe der 36 Unterfälle (Längen 269, 270, 10000 an sechs Stellen in beiden Modi) ohne Befund. `cancelSenden` toleriert `ECONNRESET` und `EPIPE` beim Schreiben, `connection reset` beim Lesen war schon kein Fehlschlag; ein Warten bis zur Frist bleibt einer. Mutanten MV4 bis MV7 in der Tabelle, rot gesehen. Bestehende Tests der Sonderanfragen (`TestProtokollrandStartnachricht`, `TestProtokollrandPruefreihenfolge`: unbekannte Sonderanfrage, Version, `SSLRequest`/`GSSENCRequest`) senden keinen ungelesenen Rest (Länge 8 oder Länge 40 ohne Rest); sie verlangen die `ErrorResponse` nur ohne ungelesene Bytes und blieben unverändert.
- Engster Sensor: Bind-Mount-Lauf des Pakets (`/home/vscode/fast.sh`, Image `pgwire-recorder:test`): `-run TestCancelRequest` grün. `make abdeckung` schrieb `abdeckung-unit.md`. Mutanten liefen gegen Bind-Mount; `server.go` danach byte-gleich (`cmp`).
- `make gates` auf `c30b830a` (Code unverändert seit `0ee53a97`, danach nur dieser Plan): Exit 0, ein Lauf (10 min 55 s), kein Flake; Abdeckung der Unit-Tests 94,43 % bei Schwelle 94 %; `test-integration`, `test-postgres-versionen`, `test-container`, `lint`, `abdeckung-check`, `kopf-check`, `workflows-check` grün. Gegenproben nicht gefahren: kein Gate und kein Vertrag eines Werkzeugs geändert.
- Metrik der Marke (§4): `git diff --numstat f7aeea92..0ee53a97 -- internal test docs/user README.md`: 546 Einfügungen, 14 Löschungen, Summe **560** (Marke 700).

**Beobachtungen.**

- F-V2 (Verifikation, LOW, nicht entschieden): `TestCancelRequestKopf` hat rund 100 Unterfälle, die je bis zu 2 s auf „keine Antwort“ warten und mit `t.Parallel()` laufen; die Laufzeit des Pakets hängt an `GOMAXPROCS` (Verifier: rund 20 s mit 20 Kernen, rund 96 s je Lauf mit `GOMAXPROCS=2`). Keine Zeitgrenze des Gates ist dafür belegt. Entscheidung des Planners (Closure): kein Folge-Slice, Eintrag `BEO-REPO/testlaufzeit-haengt-an-gomaxprocs` (1×); die Grenze des Jobs `ci-gates` ist 30 Minuten (`timeout-minutes` in `.github/workflows/ci.yml`), die Zeit des Pakets wächst auf zwei Kernen um rund 76 s gegen 20 Kerne, ob `make gates` auf dem Runner darunter bleibt, ist nicht gemessen.

- Der Mutant „Untergrenze 11, liest 4 Bytes“ hängt (der Server wartet auf das vierte Byte), die Tests färben ihn über die Frist des Clients (10 s, „Verbindung nicht geschlossen“) statt über die Antwort; ein Mutant, der `laenge-8` Bytes liest (M2d), färbt über die Antwort. Beide stehen in der Tabelle.
- F-V3 (Verifikation, LOW, behoben): ein Lauf von K13 färbte zusätzlich `Länge_10000` rot (Schreiben von 10000 Bytes an die vom Recorder schon geschlossene Verbindung: `connection reset by peer`); Ursache im Test (`cancelSenden`), nicht im Code; Eintrag `BEO-REPO/test-schreibt-an-vom-server-geschlossene-verbindung` (1×).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: wieder aufgetreten (Urteil des Planners, Closure; der Satz des Implementers „nicht wieder aufgetreten“ gilt nicht): F-1 des Reviews (Reihenfolge der Gründe, Zeitpunkt der Ablehnung, kein Fall), F-V1 der Verifikation (laufende Interaktion im `replay`, kein Fall) und L-1 der Nach-Verifikation (`Query`, `Sent`); Beleg `evidence/slice-v1-abschluss-cancel-ohne-schluessel.md`.

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

`/tmp/probe/probe4.sh` (nur der Kopf, Verbindung offen):

```bash
#!/bin/bash
# probe4.sh — nur der Kopf eines CancelRequest, Verbindung bleibt offen, socat -T1 beendet nach 1 s ohne Daten (der Client schliesst vorher nicht)
B=/tmp/probe/pgwire-recorder; W=/tmp/probe/w4; rm -rf $W; mkdir -p $W; cd $W
be32() { printf "\\x$(printf %02x $(($1>>24&255)))\\x$(printf %02x $(($1>>16&255)))\\x$(printf %02x $(($1>>8&255)))\\x$(printf %02x $(($1&255)))"; }
port=$((16000 + RANDOM % 1000))
$B record --listen 127.0.0.1:$port --upstream 127.0.0.1:1 --output $W/o.yaml 2>$W/err & pid=$!
sleep 1
t() { printf "Laenge %s, %s Bytes Rest: " "$1" "$2"; { be32 $1; be32 80877102; head -c $2 /dev/zero; sleep 3; } | socat -T1 - TCP:127.0.0.1:$port | strings | tr '\n' ' '; echo; }
for l in 9 11 13 15 269 10000; do t $l 0; done
for l in 12 16 268; do t $l 0; done
kill -TERM $pid; wait $pid; echo exit=$?
```

### Closure-Notiz

- **Was hat funktioniert:** Ein `CancelRequest` der Länge 12 wird in `record` und `replay`, mit und ohne TLS, wie jeder andere behandelt (`PGR-W3001`, Verbindung geschlossen, Exit-Code 0, `LH-FA-05`, `LH-FA-13`); die Länge wird aus dem Kopf entschieden, ohne Fork der Bibliothek. Die Randformen standen vor dem Code in der Spezifikation und im Plan §6; Review, Verifikation und Nach-Verifikation fanden sie je Länge, Stelle und Modus bestätigt oder präzisiert, kein HIGH. Endstand: 560 in der Metrik von §4 (Marke 700), `make gates` grün auf `c30b830a` (Code unverändert seitdem), `make gates-schnell` am Endstand grün (Nach-Verifikation Abschnitt 6).
- **Was ging anders als geplant:** (1) Die Reihenfolge der Gründe sagte den Zeitpunkt nicht, und der Code las den Rest vor der Längenprüfung der Bibliothek (Review F-1, F-2); der Architect entschied nach dem Review (Prüfung im Kopf). (2) Die Zusage zur laufenden Interaktion galt für `replay` ohne Fall (Verifikation F-V1); der Architect entschied „laufend“ je Modus, der Fall kam hinzu. (3) Länge 269 bis 10000 mit Rest: ein Flackern im Test (F-V3), die `ErrorResponse` ist dort nicht mehr zugesagt (Architect). (4) Die Schätzung 350 stimmte für den ersten Stand (354); die Nacharbeit zu Review und Verifikation brachte 560. (5) Die Nacharbeit nach dem Review entfernte §8 des Plans (Commit `d8a72c11`); weder die Verifikation noch die Nach-Verifikation merkten es, der Planner stellte §8 bei der Closure her (Eintrag im Register).
- **Steering-Loop-Eintrag:** benannte Spec-Lücke: `spec/spezifikation.md` `LH-FA-05.e` sagt, die zweite Verbindung bekomme im `replay` die aufgezeichneten Antworten „ohne Veränderung ihres Cursors“; kein Test prüft den Cursor (der Fake des Unit-Tests hält keinen, `test/integration/cancel_e2e_test.go` läuft nur im `record`), die Unit-Prüfung zählt Aufrufe am Fake und lässt `Query` und `Sent` im `replay`-Zweig von `cancelGeschlossen` durch (Nach-Verifikation L-1, zwei Mutanten grün). Die Entscheidung, die Zusage enger zu fassen oder einen Fall mit echtem Use Case zu verlangen, trifft der Architect (Rückgabe unten). **Empfehlung des Planners, keine Entscheidung:** `AGENTS.md` §3.12 (Reihenfolge der Gründe) um die Pflicht ergänzen, je Grund den Zeitpunkt zu nennen, zu dem er feststellbar ist (am Kopf, nach dem Rest, mit dem Ende der Verbindung), auch bei offener Verbindung; die Reihenfolge allein ließ hier den Zeitpunkt offen (`BEO-REPO/spec-randform-erst-im-review-entschieden`, verkörpert; die Schärfung entscheidet der Nutzer).
- **Beobachtungs-Register (`../observations/`):** neu `BEO-REPO/cancelrequest-laenge-ohne-client-beleg/` (Risiko 1, 1×), `BEO-REPO/testlaufzeit-haengt-an-gomaxprocs/` (F-V2, 1×) und `BEO-REPO/test-schreibt-an-vom-server-geschlossene-verbindung/` (F-V3, 1×); weitere Belege `evidence/slice-v1-abschluss-cancel-ohne-schluessel.md` in `BEO-REPO/gruener-mutant-faelschlich-aequivalent/` (F-2; 10), `BEO-REPO/spec-randform-erst-im-review-entschieden/` (F-1, F-V1; 46), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (F-1, F-V1, L-1; 63), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (F-4, L-1; 88), `BEO-REPO/plan-folgt-korrektur-nicht/` (F-3, F-V4, L-3; 75), `BEO-REPO/liste-gruener-mutanten-unvollstaendig/` (L-1; 4) und `BEO-REPO/implementer-ueberschreibt-plan-abschnitte/` (§8; 2); Vermerke ohne Evidence in `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf/` (Gegenbeleg: Verifikation und Nach-Verifikation fuhren jede Behebung von Hand, der Bericht nennt den Anlass; Zähler bleibt 3) und `BEO-REPO/verhalten-nur-unter-linux-geprueft/` (Gegenbeleg: §6 nennt die Plattform vor dem Code, die Grenze steht an der Zusage in der Spezifikation). Kein Eintrag erreichte mit diesem Slice 3× neu; `BEO-REPO/liste-gruener-mutanten-unvollstaendig/` (4×) und `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf/` (3×) stehen offen ohne Zuweisung, die der Nutzer trifft; der Lese-Schritt der Closure von `welle-v1-nachtrag` sichtet sie. F-V4 (Größe) und F-3 (Ort) sind behoben; `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: nicht aufgetreten (Schätzung 350, erster Stand 354).
- **Folge-Slices:** keiner angelegt. Die Rückgaben unten sind Entscheidungen des Architect oder des Nutzers, keine Sendung an einen Slice; `slice-v1-abschluss-protokollrand` (§1, Bestand, in `done/`) bekam keine Sendung.
- **Risiken aus §6:** *Ein Client sendet eine Länge, die die Spezifikation als `PGR-E6001` führt* — weiter offen → `BEO-REPO/cancelrequest-laenge-ohne-client-beleg` (kein Auftreten, keine Probe mit einem Client; die Behandlung ist je Länge gepinnt). *Integrationstest trifft die Zeitfolge nicht stabil* — entfallen: Frist 20 s gegen 0,5 s, grün im Lauf des Implementers, des Verifiers und in drei `make gates` ohne Wiederholung; nicht unter künstlicher Last gefahren (§6).
- **Rückgaben:** An den Architect: Die Zusage „ohne Veränderung ihres Cursors“ (`LH-FA-05.e`) geht weiter als jeder Test (Steering-Loop-Eintrag); L-2 (Rest bei Länge 9 bis 15) ist mit `f69478f2` entschieden. An den Nutzer: die Zuweisung von `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (4×) und `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (3×); die Empfehlung zu `AGENTS.md` §3.12 oben; `spec/` liegt nicht auf der Positivliste des Stop-Hooks (Nach-Verifikation Abschnitt 6), das ist beim Merge nach `main` zu beachten.
- **Drei Paarungen:** Anker: Der Lerneintrag verkörpert nichts (`liegt in` entfällt, die Lücke ist benannt). Folge-Slice: keiner genannt, keine Adresse zu prüfen. Register: die drei neuen und sieben fortgeschriebenen Einträge stehen als Verzeichnisse mit nicht leerem `evidence/` (`BEO-REPO/cancelrequest-laenge-ohne-client-beleg`, `BEO-REPO/testlaufzeit-haengt-an-gomaxprocs`, `BEO-REPO/test-schreibt-an-vom-server-geschlossene-verbindung`, dazu die sieben oben).
- **Nachzug des Planners (`AGENTS.md` §3.9, Closure, 2026-10-20):** L-1 und L-3 der Nach-Verifikation sind im Plan nachgezogen (§6 Randform 2, §7 Zeile zu F-V1); die Zusage der Spezifikation selbst blieb unberührt. §8 stand nicht mehr in der Datei und ist aus dem Stand `1ccfa245` hergestellt, mit der Sichtung bei Closure.
- **Welle:** Der Slice trägt `Welle: welle-v1-nachtrag`; er bleibt als Volltext in `done/`, bis die Welle schließt und `make archive-welle` ihn archiviert; die Roadmap bleibt unverändert (keine Umplanung, der Zustand steht im Verzeichnis).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — der Fall entstand, weil der Test nur den `CancelRequest` mit Schlüssel prüfte. Bei der Closure (2026-10-20) wieder durchgesehen: dieser Slice hat ihn, `BEO-REPO/spec-randform-erst-im-review-entschieden`, `BEO-REPO/gruener-mutant-faelschlich-aequivalent`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und `BEO-REPO/plan-folgt-korrektur-nicht` berührt (Belege in §7); alle fünf sind `verkörpert`, die Regel bleibt.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (PGWire-Adapter, Dokumentation). Nachgezählt beim Eintragen der Randformen (2026-10-20, Architect, `AGENTS.md` §3.13): zwei Liefer-Punkte (`LH-FA-05`, `LH-FA-13`), zwei Schichten; der Bootstrap bleibt unberührt, der Integrationstest zählt als Test des Binaries; die Spezifikation trägt der Architect vor dem Code, sie zählt nicht in die Marke. Nachtrag nach dem Review (2026-10-20, Architect): die Prüfung der Länge im Kopf ist ein Zweig in `startkopf` und ein Fall mit offener Verbindung; keine neue Schicht, kein neuer Liefer-Punkt. Bei der Closure (2026-10-20, Planner) nachgezählt: unverändert zwei Liefer-Punkte, zwei Schichten; die Nacharbeit nach Review und Verifikation (F-V1, F-V3) fügte Tests hinzu, keine Schicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
