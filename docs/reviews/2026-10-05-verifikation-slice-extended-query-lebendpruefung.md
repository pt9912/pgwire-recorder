# Verifikation: slice-extended-query-lebendpruefung — 2026-10-05

**Rolle:** Verifier (Modul 11). Ich prüfe, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-extended-query-lebendpruefung.md` (Kopf, §1 bis §4, §6) gegen den Gesamt-Diff `90f6fb1..aca39ce` bei HEAD `aca39ce`. Ausgenommen sind die Planer-Commits `c1f1e43` und `5047288`, die zu `slice-v1-abschluss-postgres-versionen` gehören. Code-Commits: `0841891` (Umsetzung), `98f4140` (F-350, F-355, F-356), `aca39ce` (F-360 bis F-362). Spec-, Lastenheft- und ADR-Commits: `f906c04`, `4abad8e`, `eb61782`, `7ffd029`, `e2871a3`, `4149080`, `c9ba569`.

**Eingang:**

- die DoD-Liefer-Punkte 1 bis 3 und die Commit-Messages
- der [Review-Report](2026-10-05-review-slice-extended-query-lebendpruefung.md) (F-350 bis F-359)
- das [Folge-Review](2026-10-05-folge-review-slice-extended-query-lebendpruefung.md) (F-360 bis F-363). `aca39ce` behebt F-360 bis F-362 und hat danach kein Review gesehen; ich habe ihn mit eigenen Mutationen geprüft (U18, U19, U07 gegen PostgreSQL 16).

Der Arbeitsbaum war bei meinem Start sauber.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen liefen nur in Kopien von `git archive aca39ce` im Scratchpad:

- Unit-Mutationen: je Mutation eine eigene Kopie, `go test -count=1` über `internal/hexagon/services` und `internal/adapters/driving/pgwire`, netzlos in einem Container des Images `pgr-verif-lp:deps` (Stufe `deps` des `Dockerfile`)
- E2E: Images `pgr-verif-lp:int-<name>` (Stufe `integration`) aus der unveränderten Kopie und aus sieben Mutanten, eigenes internes Docker-Netz `pgr-verif-lp-net`, eigene PostgreSQL-Container 14.24, 16.15, 17.11 und 18.6 aus den per Digest gepinnten Images. `server_version` habe ich in jedem Container mit `psql` abgefragt.

Eigene Container, das Netz, die eigenen Images und die Kopien habe ich danach gelöscht. Die PostgreSQL-Images des Hosts habe ich nicht angefasst. `git status --short` war vor dem Anlegen dieser Datei leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), drei Behauptungen (Unit-Tests). **Bestätigt.**

Gegen die Negative von LH-FA-09 (eingehend und ausbleibend aufgezeichnet), gegen LH-FA-09.a (*Definition*, *Serverversion*, *Antwort*, *Aufgezeichnete Lebendprüfungen*) und gegen [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) (*Erkennung*, *Eingehend*, *Aufgezeichnet*, *Grenze*). Für jede Teilbehauptung gilt: Die Unit-Tests sind bei `aca39ce` grün, und die Mutationen färben sie rot.

**1a — Eingehende Lebendprüfung zwischen zwei Interaktionen: `EmptyQueryResponse` und `ReadyForQuery`, Cursor bleibt stehen. Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| U01 | **Antwort außerhalb der Reihe abgeschaltet** (Zustand ohne den Fix für eingehende) | `TestReplayLebendpruefungAusserDerReihe`, `…Extended`, `…Serverversion`, `TestReplayExtendedAmCursor`. Grund: Die Lebendprüfung `-- ping` erhält keine `EmptyQueryResponse`, sondern wird mit der Interaktion am Cursor verglichen. |
| U20 | Lebendprüfung ordnet eine Session zu | `TestReplayLebendpruefungAusserDerReihe` |
| U04 | Status immer `I` | `…AusserDerReihe`, `…Extended` |
| U05 | Status nicht aus Extended-Antworten | `…Extended` |
| U24 | Status nicht aus einfachen Antworten | `…AusserDerReihe` |
| U06 | Status nach dem Handshake leer | `…AusserDerReihe`, `…Extended`, `TestReplayExtendedAmCursor` |

Die Teilbehauptungen „auch vor der ersten und nach der letzten Interaktion“ (LH-FA-10.a) und „ohne Zuordnung, auch ohne freie Session“ (LH-FA-03.a Schritt 4, LH-FA-12.a) trägt `TestReplayLebendpruefungAusserDerReihe`. Die Verbindungen `a` und `c` stellen dort Lebendprüfungen ohne freie Session.

**1b — Eine aufgezeichnete Lebendprüfung, die ausbleibt, ist keine Abweichung. Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| U02 | **Überspringen abgeschaltet** (`ohneLebendpruefungen` entfernt, Zustand ohne den Fix für aufgezeichnete) | `TestReplayLebendpruefungAufgezeichnet` (`erwartet "-- ping", empfangen "SELECT 2"`), `…NummerNachDemEnde`, `…Serverversion` |
| U22 | Filter nimmt auch Extended-Interaktionen | zehn `TestReplayExtended*` und zwei `TestReplayLebendpruefung*` |
| U16 | offener Blockkommentar gilt als Lebendprüfung | `…Erkennung`, `…VertikalerTabulator`, `…Aufgezeichnet` |
| U18 | Nummer nach dem Ende aus der Zahl der Interaktionen statt aus `sequence` (F-356) | `…Aufgezeichnet`, `…NummerNachDemEnde`, `TestReplayLetzteNummer` |
| U19 | `letzteNummer` ohne den Schutz für eine Session ohne Interaktion (F-362) | `TestReplayLetzteNummer` (Panik `index out of range [-1]`) |

`PGR-E3004` für eine Aufzeichnung nur aus Lebendprüfungen (LH-FA-03.a Schritt 2) prüft `TestReplayLebendpruefungAufgezeichnet` an ihrem Ende. Die Warnung „1 von 2“ nach LH-FA-03.b zählt die Lebendprüfungen nicht mit; auch das prüft dieser Test.

**1c — Jede andere einfache Anfrage außer der Reihe, auch mitten in einer Extended-Interaktion, bleibt `PGR-E5001`. Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| U03 | `mitten()` nicht geprüft: Lebendprüfung mitten in einer Extended-Interaktion beantwortet | `TestReplayLebendpruefungExtended` |
| U21 | nach einer Flush-Gruppe gilt der Cursor nicht als „mitten“ | `TestReplayLebendpruefungExtended` |
| U17 | `;` als Leerraum | `…Erkennung`, `…VertikalerTabulator`, `…Aufgezeichnet`, `TestReplayExtendedAmCursor` |

Die übrigen Fälle der falschen Protokollart nach LH-FA-18.a prüfen die Tests aus `slice-extended-query-replay` (`TestReplayExtendedFalscheArt`). Sie sind grün. Den Satz „ein `Parse` mit leerem SQL-Text ist keine Lebendprüfung“ habe ich durch Code-Lesung bestätigt: `ClientMessage` ruft `istLebendpruefung` nicht.

**1d — Erkennung (LH-FA-09.a *Definition*). Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| U14 | Blockkommentare nicht verschachtelt | `TestLebendpruefungErkennung` |
| U15 | `\r` beendet keinen Zeilenkommentar | `TestLebendpruefungErkennung` |
| U23 | `--` im Blockkommentar beginnt einen Zeilenkommentar | `TestLebendpruefungErkennung` |
| U16, U17 | siehe 1b, 1c | — |

Die Lexik habe ich durch Code-Lesung gegen die Regeln des PostgreSQL-Scanners geprüft (`xcstart`, `xcstop` `\*+\/`, `comment` bis `[\n\r]`, `space`). Dazu gehören die Fälle `/*/`, `/* a **/` und `/* a */*`. Ohne Abweichung.

**1e — Versionsgrenze (LH-FA-09.a *Serverversion*). Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| U07 | Grenze 16 | `TestLebendpruefungServerversion` (`16.15 = true`), `TestReplayLebendpruefungServerversion` (gegen 16 aufgezeichnetes `\v-- ping\v` übersprungen) |
| U08 | Grenze 18 | beide Tests (`17.0 = false`; gegen 17 aufgezeichnete Lebendprüfung nicht übersprungen) |
| U09 | `>` statt `>=` | beide Tests |
| U10 | `vt` bei der Zuordnung nicht gesetzt | `TestReplayLebendpruefungServerversion` |
| U11 | `\v` schon vor der Zuordnung Leerraum | `TestReplayLebendpruefungServerversion` |
| U12, U13 | aufgezeichnete Lebendprüfungen mit `\v` immer bzw. nie nach der Version | `TestReplayLebendpruefungServerversion` |

Ende zu Ende (siehe Punkt 2) habe ich `TestE2EReplayLebendpruefungWiePostgres` mit den Mutanten an beiden Seiten der Grenze gefahren:

| Image | Server | Ergebnis |
|---|---|---|
| unverändert | 14.24, 16.15, 17.11, 18.6 | grün |
| U08 (Grenze 18) | 17.11 | **rot:** `PGR-E5001 … erwartet "SELECT 1/0", empfangen "\v"` |
| U08 (Grenze 18) | 16.15 | grün (erwartet, die Grenze liegt dort ohnehin darüber) |
| U07 (Grenze 16) | 16.15 | **rot:** Das Replay beantwortet `\v` leer, PostgreSQL 16 antwortet `42601 syntax error`. F-360 ist damit Ende zu Ende geschlossen. |
| U07 (Grenze 16) | 17.11 | grün (erwartet, siehe V-33) |
| U10 | 17.11 | rot: `erwartet "SELECT 1/0", empfangen "\v"` |
| U04 (Status immer `I`) | 17.11 | rot: Sicht des Replay weicht ab (Transaktionsstatus `E` nach `SELECT 1/0`) |

**Alle 24 Unit-Mutationen sind rot.** Keine überlebt.

### Punkt 2 — [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [LH-QA-02](../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung): `pgxpool` und `database/sql` mit Pausen, Lagen S4 bis S7. **Bestätigt.**

- `TestE2EReplayLebendpruefungPgxpool` und `…DatabaseSQL` sind im eigenen `make gates` grün, je mit allen vier Lagen. Ebenso in meinen Läufen gegen 14.24, 16.15, 17.11 und 18.6.
- Am Treiber sind nur Host und Port umgestellt (`dsn(listen)`). `pgxpool` und `database/sql` laufen in der Default-Konfiguration, nur auf eine Verbindung begrenzt. Das entspricht LH-QA-02 und „nur Host und Port“ aus LH-FA-18.
- Dass die Pausen beim Aufzeichnen Lebendprüfungen auslösen, prüft der Test selbst: Die Aufzeichnung mit Pausen enthält mehr `-- ping` als die ohne. Dass sie es beim Wiedergeben tun, zeigt erst die Mutation U01 (unten). Der Test selbst behauptet es nicht.

**Bewusstes Brechen**, je Lage. E0 ist der Stand vor dem Slice: `replay.go` aus `90f6fb1`, die Tests von `aca39ce`.

| Lage | E0 ohne Fix | U01 Antwort aus | U02 Überspringen aus |
|---|---|---|---|
| pgxpool S4 ohne Pausen | grün | grün | grün |
| pgxpool S5 Pausen beim Wiedergeben | **rot** | **rot** | grün |
| pgxpool S6 Pausen beim Aufzeichnen | **rot** | grün | **rot** |
| pgxpool S7 Pausen bei beidem | grün | rot | rot |
| database/sql S4 | grün | rot | rot |
| database/sql S5 | **rot** | **rot** | rot |
| database/sql S6 | **rot** | rot | **rot** |
| database/sql S7 | grün | rot | rot |

Gründe der roten Lagen, aus stderr des Replay:

- **U01 und E0 in S5:** `PGR-E5001 … erwartet Client-Nachricht bind (Anweisung "SELECT $1::int + 1"), empfangen Anfrage …`. Der Pool sendet `-- ping` nach der Pause, und das Replay vergleicht es mit der nächsten Interaktion. Danach verbindet sich der Pool neu und erhält `PGR-E5003`.
- **U02 und E0 in S6:** `PGR-E5001 … Interaktion 3: erwartet Anfrage "-- ping", empfangen Client-Nachricht bind`. Die aufgezeichnete Lebendprüfung bleibt aus.

Damit färbt die abgeschaltete Antwort außerhalb der Reihe ihre Lage (S5) und das abgeschaltete Überspringen seine (S6). Bei `pgxpool` geschieht das je ohne die andere Lage. Bei `database/sql` färbt jede der beiden Mutationen alle vier Lagen. Grund: `database/sql` prüft eine Verbindung bei der ersten Wiederverwendung auch ohne Pause; schon die Aufzeichnung ohne Pausen enthält also Lebendprüfungen. So beschreibt es der Kommentar an `lebendLagen`.

S7 ist im Zustand ohne den Fix grün: Lebendprüfungen in Aufzeichnung und Wiedergabe fallen dort zusammen. Die Lage ist also ein Regressionsschutz und kein Beleg für den Fix. Den Fix belegen S5 und S6 (V-35).

### Punkt 3 — ADR angenommen und im Index; LH-FA-09.a und LH-FA-18.a §Replay nennen Ausnahme und falsche Protokollart; Handbuch nennt das Verhalten. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) `Accepted` | Kopf der ADR, Geschichte mit Zeile `Accepted` (`c9ba569`) | bestätigt |
| im ADR-Index | Indexzeile 0031 und Zeile in der Tabelle der ergänzenden ADRs („ergänzt 0007“, kein `Supersedes`) | bestätigt; [ADR-0007](../plan/adr/0007-strict-replay.md) ist unverändert (nicht im Diff) |
| LH-FA-09.a nennt die Ausnahme | Schritt 1, „einzige Ausnahme“, Absatz *Lebendprüfung* mit *Definition*, *Serverversion*, *Antwort*, *Aufgezeichnete Lebendprüfungen* | bestätigt; die Leerraum-Regel je Serverversion steht in *Definition* und *Serverversion* |
| LH-FA-18.a §Replay nennt die falsche Protokollart | Absatz „Eine Client-Nachricht der falschen Protokollart am Cursor …“, mit `Query` mitten in einer Extended-Interaktion und `Parse` mit leerem Text | bestätigt |
| Folgestellen | LH-FA-03.a (Schritte 2 und 4), LH-FA-03.b, LH-FA-10.a, LH-FA-12.a, `SPEC-011`, Historie | bestätigt; LH-FA-06.a unverändert, wie §1 sagt |
| Lastenheft LH-FA-09 | Die Negative nimmt eingehende und ausbleibende aufgezeichnete Lebendprüfungen aus; Status `Draft`, kein Versionssprung | bestätigt |
| Handbuch | „Mit einem Datenbanktreiber arbeiten“ (Pools, `database/sql`, `ShouldPing`, `\v` ab 17, vor der ersten Anfrage nie, mitten in einer Folge `PGR-E5001`), Wiedergabe, Fragen, `PGR-E3004` | bestätigt; F-361 ist eingearbeitet |

Spezifikation gegen Code, Satz für Satz:

- Leerraum: `" \t\n\r\f"`, dazu `\v` nach Version.
- Zeilenkommentar bis `\n` oder `\r`.
- Verschachtelte Blockkommentare; `--` darin bedeutungslos.
- Nicht-ASCII-Zeichen sind keine Lebendprüfung: byteweise Prüfung.
- Hauptversion aus den führenden Ziffern; ein Überlauf gilt als „nicht lesbar“.
- `vt` vor der Zuordnung `false`.
- Status nach dem Handshake `I`.

Ohne Abweichung.

`Schärft:` von [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) nennt LH-FA-03.a und LH-FA-10.a nicht (F-352). Der Reviewer hat das begründet offen gelassen. Für die DoD ist es ohne Folge, die Begründung gehört in die Closure-Notiz.

### Punkt 4 — `make gates` grün. **Bestätigt.**

- Eigener Lauf von `make gates` an `aca39ce`, **Exit 0**:
  - `abdeckung-check` und `abdeckung-gegenprobe` grün
  - `a-check` 0 Befunde, `a-check-negativ` grün
  - `baseline-verify` v6.13.0 OK, 54 Dateien
  - `docs-check` 178 Dateien, 0 Befunde
  - `commit-msg-gegenprobe` grün
  - `test-integration` in beiden Phasen grün (`run-integration-tests: gruen`), 31 `--- PASS`-Zeilen über beide Phasen, keine `--- FAIL`
- Die Stufe `test` kam im Gate-Lauf aus dem Cache. In der Kopie habe ich sie darum ohne Cache gefahren (`--no-cache-filter test`): `gofmt -l` leer, `go vet -tags integration` grün, alle sechs Unit-Pakete `ok`.

**Zusätzlich gegen PostgreSQL 16:** `make test-integration POSTGRES_IMAGE=postgres:16-alpine@sha256:721873c34ceb9f8d8fc265984940dc982404c105f19ad51be9fdc5970a6080ea` im Repo, **Exit 0**, beide Phasen grün, alle drei `TestE2EReplayLebendpruefung*` mit allen Lagen `PASS`.

Gegen 14.24 und 18.6 habe ich die erste Phase mit dem eigenen Image gefahren: je 27 `PASS`, 0 `FAIL`. Vier Tests wurden übersprungen, weil `PGR_DATEN` fehlte (zweite Phase und ihre Vorbereitung). Die Behauptung in §6 Risiko 5 („volle Suite gegen 14, 16, 17 und 18“) ist damit für 16 und 17 voll bestätigt, für 14 und 18 in der ersten Phase.

### Prozess-Punkte der DoD (nur vermerkt)

- **Review:** Reports zu `0841891` und `98f4140` liegen vor. `aca39ce` hat kein Review. Dieser Bericht prüft ihn funktional (U07 gegen 16, U18, U19). Ob ein weiteres Folge-Review nötig ist, entscheidet der Planner.
- **Offen, wie bei `in-progress` zu erwarten:** Closure-Notiz, Register, Ausgänge der Risiken in §6, die Paarungen. Alle fünf Risiken tragen „offen bis Closure“. Das ist kein Ausgang aus der geschlossenen Menge und muss bei der Closure ersetzt werden. Risiko 5 ist nach dem eigenen Text eingetreten; sein Kandidat ist `slice-v1-abschluss-postgres-versionen`.

---

## 2. Plan gegen Code

### Kopf

| Feld | Urteil |
|---|---|
| `Bezug` | konform: LH-FA-18, LH-FA-09, LH-QA-02, [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) |
| `Berührte Spec-Stellen` | konform; jede genannte Stelle ist im Diff berührt. `ARC-002` ist über den Cursor berührt (`status`, `vt`); `spec/architecture.md` ist unverändert, wie bei den Vorgänger-Slices. |

### §1 Ziel und Abgrenzung

| Plan | Code | Urteil |
|---|---|---|
| Lebendprüfung zwischen Interaktionen, `EmptyQueryResponse` + `ReadyForQuery`, Cursor bleibt | `Query`, Kurzschluss vor `zuordnen` | konform (U01, U20) |
| aufgezeichnete übersprungen, auch für Zuordnung, nicht verbrauchte, `PGR-E3004` | `ohneLebendpruefungen` in `NewReplayService` | konform (U02) |
| Status des letzten `ReadyForQuery` | `cursor.status`, `merke` | konform (U04 bis U06, U24) |
| jede andere Anfrage, auch mitten in Extended, `PGR-E5001` | `mitten()` | konform (U03, U21) |
| `\v` nach `server_version`, vor der Zuordnung nicht | `vtLeerraum`, `cursor.vt` | konform (U07 bis U13) |
| Record und `play` unverändert, kein Format | nicht im Diff | konform |
| `pgwire`-Adapter ohne Änderung | nicht im Diff | konform |
| Lastenheft Draft, kein Versionssprung | Kopf des Lastenhefts | konform |

### §3 Plan

Jede Zeile hat ihr Gegenstück im Diff. Die Folge-Slices sind in `4149080` und `aca39ce` nachgezogen: `slice-v1-abschluss-sessions` mit der Kopplung an `frei` und F-362, `slice-v1-abschluss-postgres-versionen` mit der Grenze. Abweichung: `go.mod` nimmt neben `puddle` auch `golang.org/x/sync` als indirekte Abhängigkeit auf; §3 nennt nur `puddle` (V-34).

### §4 Trigger

Keine Rückführung ist eingetreten: Das Recording-Format ist unverändert, und [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) ist angenommen.

### §6 Risiken

Die Texte entsprechen dem Code und den Tests. Die in Risiko 1 genannten Mutationen („Status immer `I`“, „nicht aus Extended-Antworten“, „nach dem Handshake leer“) habe ich als U04, U05 und U06 selbst rot gesehen. Zur Sonde in Risiko 5: Die Versionen 14.24, 16.15, 17.11 und 18.6 habe ich in eigenen Containern bestätigt (`show server_version`); die Grenze ist Ende zu Ende gegen 16 und 17 belegt (1e).

---

## 3. Befunde

| # | Schwere | Befund | Adressat |
|---|---|---|---|
| V-33 | INFO | Das Gate fährt die Integrationstests nur gegen PostgreSQL 17. Eine zu niedrige Grenze im Produktcode (U07, Grenze 16) ist dort Ende zu Ende grün und nur gegen 16 rot. Im Gate fangen sie die Unit-Tests (`TestLebendpruefungServerversion`, `TestReplayLebendpruefungServerversion`); die Messung am realen Server unter 17 liegt außerhalb des Gates. Das ist der Gegenstand der Versionsmatrix in `slice-v1-abschluss-postgres-versionen`, keine DoD-Lücke. | Planer `slice-v1-abschluss-postgres-versionen` |
| V-34 | INFO | §3 nennt für `go.mod` nur `puddle`; der Diff nimmt auch `golang.org/x/sync` (indirekt, über `puddle`) auf. Ohne Wirkung auf die DoD; Hard Rule 3.9 verlangt den Nachzug nur für Änderungen mit Plan-Bezug. | Implementer, bei der Closure |
| V-35 | INFO | Lage S7 („Pausen bei beidem“) ist im Stand vor dem Slice grün (E0), für beide Treiber. Sie belegt den Fix nicht, sondern schützt vor einer Regression (U01 und U02 färben sie rot). Die DoD-Behauptung „ohne Abweichung“ ist davon nicht berührt. Wer die Lagen später kürzt, darf S5 und S6 nicht streichen. | Closure-Notiz |

Keine Befunde der Schwere HIGH, MEDIUM oder LOW.

## 4. Urteil

| DoD-Punkt | Urteil |
|---|---|
| 1 — LH-FA-09, Unit-Tests | **bestätigt**, 24 von 24 Mutationen rot |
| 2 — LH-FA-18, LH-QA-02, Lagen S4 bis S7 | **bestätigt**; abgeschaltete Antwort färbt S5, abgeschaltetes Überspringen S6, der Stand ohne den Fix beide |
| 3 — ADR, Index, Spezifikation, Handbuch | **bestätigt** |
| 4 — `make gates` | **bestätigt**, Exit 0 an `aca39ce`; zusätzlich grün gegen 16 (volle Suite) und gegen 14 und 18 (erste Phase) |

Die Liefer-Punkte sind erfüllt. Für die Closure bleiben die Prozess-Punkte offen; die Risiko-Ausgänge in §6 brauchen dabei Werte aus der geschlossenen Menge.
