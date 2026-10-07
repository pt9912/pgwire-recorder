# Slice slice-harness-blackbox-driven: Black-Box-Tests der getriebenen Adapter (PostgreSQL und Recording)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Er trägt zum
Abnahmeszenario 17 (M3) bei, das `slice-lastenheft-pruefbarkeit` anlegt und das mit
dem letzten Slice der Reihe nachweisbar wird. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 3; Messmethode 1 für die Testdateien dieser Pakete). Bindung: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Export-Test-Brücke, Entscheidung 4; Bereinigung vor dem Gate ohne Stufen, Entscheidung 5: die Befunde der Testdateien dieser Pakete behebt dieser Slice), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckungs-Deklarationen an den Tests bleiben unverändert).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (Punkt 7 und 8: Brücke und dauerhafte Ausnahmen, gegen die umgestellt wird. Geändert vom Architect unter der Kennung dieses Slice: Punkt 7, die Funktion, an die die Brücke zum Erzeugen weiterreicht, ruft auch der Produkt-Code)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Unit-Tests unter `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` laufen als Black-Box-Pakete
(`package <name>_test`) und prüfen über die exportierte Schnittstelle des Pakets;
keine Testdatei unter diesen Pfaden hat in `make lint` einen Befund. Dafür behebt der
Slice in den Testdateien, die er umschreibt, die 12 Befunde aus der Messung in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Stand `79f40e1`; Aufteilung je Paketgruppe in der Fassung `92d1b86`
der ADR): `testpackage` 3, dazu `errcheck`, `gochecknoglobals` und `unused`, zusammen 9. Wo ein Test einen unexportierten Teil braucht, geht
er über die Export-Test-Brücke nach `SPEC-049` Punkt 7 (`export_test.go`) oder wird
gegen die exportierte Schnittstelle umgeschrieben. Die Fälle und ihre Prüfungen
bleiben erhalten; der Slice ist ein Umbau der Tests, keine neue Prüfung.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: Die Tests werden auf
Black-Box-Pakete umgestellt; am 2026-10-06, nach der Messung: Bereinigung vor dem
Gate, ohne Stufen, und die vier Umstellungs-Slices beheben dabei auch die übrigen
Befunde ihrer Testdateien
(`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`). Geschnitten ist
nach Paketgruppe, je in einer Schicht des Hexagons, damit jeder Schnitt einzeln
lieferbar und in einer Review-Sitzung prüfbar bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Tests der übrigen Pakete — sie übernehmen `slice-harness-blackbox-kern`, `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`.
- Neue Exporte im Produkt-Code nur für Tests — sie verbreiterten die Schnittstelle des
  Pakets für einen Zweck außerhalb des Produkts; der Weg für unexportierte Teile ist
  die Brücke in einer `_test.go`-Datei (§6). Braucht ein Test mehr, ist das ein Befund
  für den Architect, keine Entscheidung im Umbau (§4).
- Neue Fälle oder geänderte Erwartungen — ein anderer Vorgang; der Umbau soll an
  derselben Testliste messbar sein (§2). Die Schwelle der Testabdeckung übernimmt
  `slice-harness-coverage`, gemessen erst nach allen vier Umstellungs-Slices.
- Das Scharfschalten am Gate und der White-Box-Fall für diese Pfade in der
  Lint-Gegenprobe — übernimmt `slice-harness-lint`: Vor dem Gate gibt es weder eine
  Stufe noch eine Gegenprobe (Entscheidung 5); dieser Slice belegt sein Ergebnis mit
  dem Werkzeug `make lint`.
- Befunde im Produkt-Code dieser Pakete — übernimmt `slice-lint-bestand-kern-driven` nach den
  Umstellungs-Slices.
- Produkt-Verhalten, Lastenheft und die Spezifikation außer `SPEC-049` Punkt 7 —
  Schicht-Abgrenzung: Der Slice ändert Testdateien und im Produkt-Code nur eine
  Stelle: `Open` in `upstream.go` legt seine Session über den unexportierten
  Konstruktor `newSession(conn)` an (§6 *Session über `net.Pipe`*, Entscheidung des
  Architect nach §4). Das ist kein Export und keine Verhaltensänderung, sondern
  der Weg, auf dem `TestCloseMitSchreibfrist` eine vom Produkt erzeugte Session
  bekommt. `.golangci.yml` und die Gegenprobe des Lint-Gates ändert er nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Die Testdateien unter `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` gehören zu `_test`-Paketen; die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau gleich, keine Prüfung ist
      entfallen, die Abdeckungs-Deklarationen nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) stehen unverändert und
      `make abdeckung-check` ist grün. Unexportierte Teile erreicht ein Test nur über
      die Brücke nach `SPEC-049` Punkt 7. — Verifikation, Abschnitt 1 Punkt 1 und
      Abschnitt 3: `package postgres_test` und `package recording_test`, im Paket des
      Codes nur `postgres/export_test.go`; `go test -list .` an `aff75ee` und `874b29d`
      gleich (`postgres` 14, `recording` 13), 108 Zeilen `--- PASS` an beiden Ständen,
      11 Deklarationen gleich und gleich zugeordnet, `make abdeckung-check` Exit 0; die
      Brücke reicht nur weiter und legt keinen Wert an (`SPEC-049` Punkt 7 in der
      Fassung `f0ea7be`).
- [x] `make lint` meldet in den Testdateien unter `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` keinen Befund: Die
      12 Befunde der Messung sind behoben, ohne `//nolint` und ohne Änderung an
      `.golangci.yml`; die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8 blendet das
      Profil aus. Die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem Umbau
      stehen im Bericht; Befunde im Produkt-Code bleiben für `slice-lint-bestand-kern-driven`.
      — Verifikation, Abschnitt 1 Punkt 2: die Zeilen vorher und nachher in §7
      *Belege des Implementers*, Wort für Wort bestätigt; 12 → 0 in `*_test.go`, 4 im
      Produkt-Code; modulweit 72 → 60 ohne neuen Befund (`comm -13` leer); kein
      `//nolint`, `.golangci.yml` unverändert, kein neues `_ =`.
- [x] `make gates` grün. — an `f0ea7be` (Go-Code gleich `874b29d`; Verifikation,
      Abschnitt 1 Punkt 3 und Abschnitt 7); `76bfe09` und diese Closure ändern nur
      Berichte, Pläne und das Register, `make docs-check` und `make kopf-check` grün am
      Stand dieser Closure.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      — `docs/reviews/2026-10-07-review-slice-harness-blackbox-driven.md` (bis
      `1c30e43`; F-444 bis F-447, alle INFO).
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
| `internal/adapters/driven/postgres/upstream_test.go`, `upstream_extended_test.go` | refactor | `package postgres_test` |
| `internal/adapters/driven/recording/yaml_test.go` | refactor | `package recording_test` |
| `export_test.go` je Paket, nur wo nötig | neu | Brücke zu unexportierten Teilen nach `SPEC-049` Punkt 7; in `postgres`: an `toFrontendMessage`, an `conn` einer übergebenen Session, an `newSession` |
| `internal/adapters/driven/postgres/upstream.go` | refactor | `newSession(conn net.Conn) *session` legt `conn` und `fe` an, `Open` ruft ihn statt des Literals; ohne Verhaltensänderung (§6 *Session über `net.Pipe`*) |
| dieselben Testdateien | update | übrige Befunde aus `make lint` (dazu `errcheck`, `gochecknoglobals` und `unused`, zusammen 9) in den Dateien, die der Umbau ohnehin umschreibt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-blackbox-kern` liegt in `done/` (WIP-Limit 1,
Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05 und 2026-10-06:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices,
`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`, `slice-harness-lint`,
`slice-harness-abdeckung-gate`, Coverage, Mutation). Gemessen wird mit dem Werkzeug
`make lint` aus `slice-harness-lint-werkzeug`. Vor dem ersten Code-Commit prüft der
Architect §6 gegen `SPEC-049` (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Umbau ist nicht in einer
  Review-Sitzung prüfbar, etwa weil mehr als eine Handvoll Tests umgeschrieben statt
  umgestellt werden muss; dann je Paket ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Ein Test erreicht einen Teil, den er
  prüfen muss, weder über die exportierte Schnittstelle noch über die Brücke, die die
  ADR zulässt, oder ein Befund in einer Testdatei lässt sich nur mit einer neuen
  Ausnahme beheben; dann zuerst eine Entscheidung des Architect (Brücke erweitern,
  Schnittstelle des Pakets ändern, Ausnahme mit dauerhaftem Grund).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, `make lint` ohne Befund in den Testdateien dieser
Pfade,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden vom Architect am 2026-10-06 in
`SPEC-049` Punkt 7 und 8 und [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5; was dort nicht
steht, gibt der Implementer an den Architect zurück.

- **Export-Test-Brücke** — entschieden (`SPEC-049` Punkt 7): einzige Datei
  `export_test.go` im Paket des Codes, nur Typ-Aliase, Konstanten und Funktionen oder
  Methoden, die an Unexportiertes weiterreichen; Zustand nur an einem übergebenen Wert,
  nie auf Paketebene (eine Frist für einen Test verkürzt sie nur an einem übergebenen
  Wert). Eine Variable oder eine Funktion `Test…` darin ist ein Befund. Einen Wert
  eines unexportierten Typs legt nur der Produkt-Code an, weder Brücke noch Test, auch
  nicht lokal oder über einen Alias; was nur an ihm zu sehen ist, prüft der Test über
  die exportierte Schnittstelle (Architect 2026-10-06, F-441/F-442 in
  `slice-harness-blackbox-kern`).
- **Übrige Befunde der Testdateien** — `contextcheck`: Kontext aus `t.Context()`, in einer Funktion für `t.Cleanup` aus `context.WithoutCancel(t.Context())` (`SPEC-049` Punkt 5); `errcheck` an `Close` einer Datei oder Datenbank: den Fehler prüfen (etwa in `t.Cleanup` mit `t.Error`); `gochecknoglobals`: Testdaten in Funktionen oder als Konstanten; `unused`: entfernen. Kein `_ =` vor einem Fehler, den
  `errcheck` meldet, und keine Ausnahme, die nur Bestand aussetzt (Entscheidung 5).
- **White-Box-Zugriffe im Bestand** — Namensabgleich per Suche am Stand `ce50a10`
  (ungemessen, kann Fehltreffer enthalten, wo ein Testhelfer gleich heißt): in `recording` keine; in `postgres` u. a. `session`, `toFrontendMessage`, `lies`.
  Je Zugriff: über die exportierte Schnittstelle prüfbar, über die Brücke, oder
  Befund. `toFrontendMessage` ist eine Übersetzung in `pgproto3`-Nachrichten; ob ein Test sie direkt oder über `Upstream` prüft, entscheidet die Brücke.
- **White-Box-Zugriffe, gemessen** (Implementer per AST vor dem Code, bestätigt vom
  Architect am 2026-10-07): in `recording` keine. In `postgres`: `session` als Typ in
  Testhelfern wird `driven.UpstreamSession`; `conn` einer übergebenen, von `Open`
  erzeugten Session für `SetReadDeadline` und `CloseWrite` über die Brücke (sie
  reicht nur an das Feld eines übergebenen Werts weiter); `toFrontendMessage` über
  eine Brücke, die nur weiterreicht; `lies` war ein Fehltreffer.
- **Session über `net.Pipe`** — entschieden vom Architect am 2026-10-07 (`SPEC-049`
  Punkt 7): `TestCloseMitSchreibfrist` legt kein `session`-Literal mehr an. Über TCP
  ist der Zustand „Schreibsperre frei, Terminate blockiert“ nicht deterministisch
  herzustellen, weil die Puffer des Kernels das Terminate aufnehmen; `net.Pipe` hat
  keinen Puffer. Weg: `upstream.go` bekommt den unexportierten Konstruktor
  `newSession(conn net.Conn) *session`, den `Open` für seine Session ruft; die Brücke
  reicht mit einer exportierten Funktion, die `driven.UpstreamSession` liefert, an
  ihn weiter; der Test stellt die Verbindung (`net.Pipe`) und ruft `Close` über die
  Schnittstelle. Damit erzeugt der Produkt-Code den Wert, und ein Feld, das
  `newSession` später setzt, hat auch der Test. Verworfen: TCP mit kleinen Puffern
  (nicht deterministisch), eine exportierte Dial-Option (neuer Export, §1), die
  Prüfung streichen (sie fängt als einzige das Entfernen von `SetWriteDeadline`).
  Rot sein muss danach: `SetWriteDeadline` in `Close` entfernt →
  `TestCloseMitSchreibfrist`. Dass `Open` die Session über `newSession` anlegt, prüft
  weder ein Test noch `unused` (die Brücke hält `newSession` in Gebrauch); das ist
  Urteil des Review (Grenze von `SPEC-049`).
- **Fakes der Ports** — die Ports sind exportiert (`internal/hexagon/ports/...`); ein
  Fake in einem `_test`-Paket implementiert sie unverändert. Ein Fake, der auf
  unexportierte Felder eines Produkt-Typs greift, fällt unter die Brücke.
- **Mutationstests auf unexportierte Teile** — eine Mutation, die bisher ein Test auf
  eine unexportierte Funktion fing, muss auch nach dem Umbau fangen; sonst fehlt eine
  Prüfung, obwohl die Testliste gleich ist.

**Risiken:**

- **Prüfung geht im Umbau verloren** — ein Test bleibt in der Liste, prüft aber
  weniger, weil ein unexportierter Vergleich wegfiel
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Je umgeschriebener (nicht nur
  umgestellter) Test eine Mutation, die er weiter fängt. — **Ausgang:** entfallen: Die
  Testliste ist gleich (`postgres` 14, `recording` 13, 108 Zeilen `--- PASS` an beiden
  Ständen), die 11 Deklarationen sind gleich zugeordnet; die acht Mutationen in §7 sind
  in der Verifikation rot aus dem richtigen Grund (Abschnitt 4), ebenso R2 bis R6 des
  Reviews (F-447). Der grüne Mutant R1 (`Open` gibt eine neue Session aus
  `newSession(conn)` zurück) war schon an `aff75ee` ungeprüft, weil das Literal
  `&session{conn: conn, fe: fe}` kein Test hielt (Verifikation, Abschnitt 2); er ist
  keine im Umbau verlorene Prüfung und geht als
  `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft` ins Register (§7).
- **Abdeckung sinkt** — Black-Box-Tests erreichen unexportierte Pfade seltener; die
  Zahl misst erst `slice-harness-coverage` nach allen vier Umstellungs-Slices.
  — **Ausgang:** entfallen: Der Umbau gibt keinen Pfad ab. Jeder frühere Zugriff auf
  Unexportiertes geht über die Brücke an dieselbe Funktion oder dasselbe Feld
  (`toFrontendMessage`, `conn` über `Verbindung`), die Session aus dem früheren Literal
  legt jetzt `newSession` mit denselben zwei Feldern an (AST-Messung in §7;
  Verifikation, Abschnitt 1 *Produkt-Code* und Abschnitt 3), und `recording` hatte
  keinen White-Box-Zugriff. Die Zahl für alle vier Umstellungs-Slices misst
  `slice-harness-coverage` nach seinem eigenen §1, nicht als Ausgang dieses Risikos.
- **Befund verdeckt statt behoben** — `_ =` vor einem ungeprüften Fehler, eine Globale
  als Funktion mit demselben geteilten Zustand: `make lint` ist grün, der Test nicht
  besser. Review prüft die Form, das Werkzeug nur die Zahl (Grenze von `SPEC-049`).
  — **Ausgang:** entfallen: Kein neues `_ =` (94 Blank-Zuweisungen vorher, 93 nachher,
  weggefallen ist `_ = ss.conn.SetReadDeadline(…)`, jetzt mit `t.Fatal` geprüft), kein
  `//nolint` (Verifikation, Abschnitt 1 Punkt 2). `schliesse` prüft den Fehler von
  `Close` (V5, R3), die Testdaten-Funktionen `bereit`, `beendet` und `ergebnis` liefern
  je Aufruf frische Werte (V6, V8). Die drei `_ = s.Close()` in `upstream_test.go` sind
  Bestand und waren nie ein Befund von `errcheck` (F-446, keine Aktion).

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

- **Was hat funktioniert:** Der Schnitt hielt: drei Liefer-Punkte, drei Testdateien, eine
  Brücke in `postgres` und eine Stelle Produkt-Code, keine Rückführung aus §4. Die beiden
  Regeln aus `de770f7` (seit slice-harness-blackbox-kern) haben gewirkt. Halt vor dem
  Code: Der Implementer gab die Randform „Session ohne Literal, über TCP nicht
  deterministisch“ vor dem ersten Code-Commit zurück, der Architect entschied sie in
  `58e863c` (`newSession`, `SPEC-049` Punkt 7), und `1c30e43` ändert im Plan nur §7; das
  Review meldet zu `AGENTS.md` §3.12 keinen Befund, die Verifikation keine Randform
  außerhalb von §6. Belege in §7: AST-Messung, Testliste, Lint-Zeilen vorher und nachher
  und die Mutationstabelle standen unter *Belege des Implementers*; Review (F-447) und
  Verifikation (Abschnitt 1 und 4) prüften gegen sie und fuhren alle acht Mutationen
  nach, alle rot. Auch der grüne Mutant stand dort („Ohne roten Test“), und von dort
  fanden Review (R1) und Verifikation (Sonde) die Lücke.
- **Was ging anders als geplant:** Zwei Punkte, keine Nacharbeit am Code. Die
  Gesamtzahlen von `make lint` in §7 waren falsch (V-84): 79 → 67 statt gemessen
  72 → 60; die Zeilen unter beiden Pfaden und die Differenz 12 stimmten. Berichtigt mit
  dieser Closure in §7 und in Liefer-Punkt 2 der DoD; kein Folge-Slice übernahm die
  falsche Zahl. Und §6 band die Grenze „Urteil des Review“ an die Kopplung von `Open`
  an `newSession`; darunter fiel im Bericht auch R1, ein Mutant, der das Verhalten
  ändert und über die Schnittstelle fangbar ist (F-444, V-85). Für die Kopplung selbst
  (V9, äquivalent) trägt die Grenze. Den Wortlaut von `SPEC-049` Punkt 7 hat der
  Architect nach F-445 in `f0ea7be` eindeutig gemacht.
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-07-review-slice-harness-blackbox-driven.md`: „0 HIGH · 0 MEDIUM
    · 0 LOW · 4 INFO (F-444 `Open` ruft `newSession` für die zurückgegebene Session,
    Verhalten unverändert, Rückgabe weiter ungeprüft wie vorher; F-445 Bezug in
    `SPEC-049` Punkt 7 nach dem Einschub mehrdeutig; F-446 `_ = s.Close()` im Bestand
    neben `schliesse`; F-447 Mutationen reichen). Wiederkehrende Klasse: keine.“
    Verifikation `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-driven.md`,
    Urteil: drei Liefer-Punkte erfüllt und selbst belegt, Brücke hält Punkt 7 in der
    Fassung `f0ea7be`, acht Mutationen aus §7 rot aus dem richtigen Grund; vor der
    Closure V-84 berichtigen (Beleg mit falscher Zahl), V-85 ein Vorschlag für Register
    und Lerneintrag (Grenze trägt, Lücke adressieren). `make gates` grün an `f0ea7be`.
- **Steering-Loop-Eintrag:** Geschärfte Regel: Ein grüner Mutant wird in §7 eingeordnet
  — lässt er das Verhalten gleich, trägt ihn die Grenze in §6; ändert er es, bekommt er
  einen Test oder, wo §1 neue Fälle ausschließt, eine Test-Idee mit Grenze, der die
  Closure eine Adresse gibt; eine Grenze „nur Urteil“ deckt nur den äquivalenten Mutanten
  — liegt in `.claude/commands/implement-slice.md Schritt 19`.
  Auslöser: F-444 und V-85 (R1 gegen V9). Retirement-Check von `AGENTS.md` §3.12 und der
  Randform-Rückgabe sowie von *Belege in §7* (beide seit slice-harness-blackbox-kern):
  nicht wieder aufgetreten, beide bleiben. §3.9 und §3.11 ohne Befund (Review,
  Negativbefunde), §3.10 ohne neuen Vertrag.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `76bfe09`.
  - `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft/` — **neu**, 1× (F-444,
    V-85). `state.md`: `offen`, mit der Adresse `slice-v1-abschluss-anmeldung` §6. Der
    Ausgang `geplant` ist nach der Form des Registers eine Antwort auf die Schwelle 3×;
    darunter ist `offen` der Stand, die Adresse steht daneben. `slice-v1-abschluss-anmeldung`
    nimmt an: Er liegt in `open/`, baut die Aufbau-Schleife in `Open` für den
    Anmeldeaustausch um (§3: `internal/adapters/driven/postgres`, Weiterleitung der
    Anmeldenachrichten), sein §1 schließt Tests an `Open` nicht aus, und sein §6 trägt
    mit dieser Closure Test-Idee und Grenze als Punkt mit eigenem Ausgang.
  - Ohne Beleg: `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 3×,
    verkörpert; Halt vor dem Code hat gewirkt), `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht`
    (bleibt 3×, verkörpert; die Belege erreichten beide Prüfer),
    `BEO-REPO/plan-folgt-korrektur-nicht` (bleibt 14×; Review Hard Rule 3.9 ohne Befund,
    der Kopf deckt auch `f0ea7be`), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`
    (bleibt 16×; Review Hard Rule 3.11 ohne Befund),
    `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (bleibt 13×; kein neuer Vertrag),
    `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×; Risiko *Prüfung geht im
    Umbau verloren* entfallen), `BEO-REPO/randform-wellenlos-ohne-architect-vor-code`
    (bleibt 1×; §4 nannte den Architect, er entschied vor `1c30e43`),
    `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (bleibt 1×; jede Mutation lief in
    einem frischen Pfad per Bind-Mount).

  Einmalig und nicht eingetragen: V-84 (Gesamtzahl im Beleg falsch, Zeilen richtig). Sie
  passt zu keinem Eintrag: Die Belege erreichten die Prüfer
  (`implementer-bericht-erreicht-pruefer-nicht` nicht wieder aufgetreten), und die Zahl
  war keine Zusage über eine Prüfung hinaus (`zusage-im-kommentar-weiter-als-pruefung`),
  sondern eine falsch übertragene Messung. F-445 (mit `f0ea7be` erledigt), F-446
  (Bestand, keine Aktion), F-447 (Mutationen reichen). Mit diesem Slice erreicht
  **kein** Eintrag die Schwelle 3× neu; über der Schwelle stehen nur verkörperte
  Einträge (`implementer-bericht-erreicht-pruefer-nicht` 3×,
  `randform-im-code-entschieden-dann-zurueckgegeben` 3×,
  `spec-randform-erst-im-review-entschieden` 11×,
  `negativtests-fehlen-bei-neuem-vertrag` 13×, `plan-folgt-korrektur-nicht` 14×,
  `zusage-im-kommentar-weiter-als-pruefung` 16×).
- **Folge-Slices:** `slice-v1-abschluss-anmeldung` (Test-Idee zu R1, §6),
  `slice-lint-bestand-kern-driven` (die vier Befunde im Produkt-Code dieser Pakete,
  Zeilen in §7 *Nachher*), als nächster in der Reihe
  `slice-harness-blackbox-pgwire`, dann `slice-harness-blackbox-einstieg`,
  `slice-harness-coverage` und `slice-harness-lint`.
- **Risiken aus §6:** drei, alle **entfallen**, mit Begründung: *Prüfung geht im Umbau
  verloren* (R1 war vorher ungeprüft, ins Register), *Abdeckung sinkt* (kein Pfad
  abgegeben), *Befund verdeckt statt behoben* (kein neues `_ =`, kein `//nolint`). Die
  Randformen in §6 sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker — `liegt in` nennt
  `.claude/commands/implement-slice.md Schritt 19`; `grep -n "seit slice-harness-blackbox-driven" .claude/commands/implement-slice.md`
  findet ihn in Schritt 19. Folge-Slice — `slice-v1-abschluss-anmeldung`,
  `slice-lint-bestand-kern-driven`, `slice-harness-blackbox-pgwire`,
  `slice-harness-blackbox-einstieg` und `slice-harness-coverage` liegen in `open/`,
  `slice-harness-lint` in `next/`; `slice-v1-abschluss-anmeldung` nennt den Punkt in §6.
  Register — die neue Kennung und die acht ohne Beleg bestehen als Verzeichnis, jedes mit
  nicht leerem `evidence/`. Die nächste Welle-Closure prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-07-review-slice-harness-blackbox-driven.md`
  (bis `1c30e43`; F-444 bis F-447), Verifikation
  `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-driven.md` (bis
  `874b29d`, `make gates` an `f0ea7be`; V-84, V-85), Architect `58e863c` und `f0ea7be`,
  Entscheidung [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) unverändert.
  Validierung: n/a, der Slice ändert Tests, kein End-Nutzer-Verhalten.

**Belege des Implementers** (Arbeitsbaum auf `58e863c` mit dem Diff des Commits, der
diesen Abschnitt anlegt):

*White-Box-Zugriffe, per AST gemessen* (go/types im Image der Stufe `deps`, ohne Netz;
gezählt sind Bezeichner in Testdateien, die auf ein unexportiertes Objekt einer
Produkt-Datei desselben Pakets zeigen, dazu Literale, `var` und `new` unexportierter
Typen in Testdateien):

| Paket | vorher (`aff75ee`) | nachher |
|---|---|---|
| `recording` | 0 Zugriffe, 0 angelegte Werte | 0, 0 |
| `postgres` | `session` (Typ) 5, `conn` 3, `fe` 1, `toFrontendMessage` 2; ein `session`-Literal (`TestCloseMitSchreibfrist`); `lies` ein Fehltreffer der Suche (lokaler Helfer) | nur in `export_test.go`: `toFrontendMessage`, `newSession`, `session` und `conn` in `Verbindung`; 0 angelegte Werte |

Die Brücke `internal/adapters/driven/postgres/export_test.go` hat drei Funktionen:
`ToFrontendMessage` reicht an `toFrontendMessage` weiter, `NeueSession` an `newSession`
(liefert `driven.UpstreamSession`), `Verbindung` liefert `conn` einer übergebenen
Session. `recording` braucht keine Brücke.

*Testliste* (`go test -list .` je Paket, Image der Stufe `deps`): vorher und nachher
gleich, `postgres` 14 Tests, `recording` 13; `diff` leer. Die Abdeckungs-Deklarationen
stehen unverändert (der Diff der Testdateien berührt keine Kommentarzeile `Abdeckung:`).

*`make lint`, Zeilen unter beiden Pfaden.* Vorher (`aff75ee`, Exit 2, 72 Befunde im Repo;
die Gesamtzahl vorher und nachher nach der Verifikation, V-84: golangci-lint `72 issues:`
an `aff75ee`, `60 issues:` an `874b29d` und im Arbeitsbaum):

```text
internal/adapters/driven/postgres/upstream_extended_test.go:1:9: package should be `postgres_test` instead of `postgres` (testpackage)
internal/adapters/driven/postgres/upstream_extended_test.go:98:15: Error return value of `s.Close` is not checked (errcheck)
internal/adapters/driven/postgres/upstream_extended_test.go:178:15: Error return value of `s.Close` is not checked (errcheck)
internal/adapters/driven/postgres/upstream_extended_test.go:206:16: Error return value of `s.Close` is not checked (errcheck)
internal/adapters/driven/postgres/upstream_extended_test.go:220:16: Error return value of `s.Close` is not checked (errcheck)
internal/adapters/driven/postgres/upstream_extended_test.go:314:15: Error return value of `s.Close` is not checked (errcheck)
internal/adapters/driven/postgres/upstream_test.go:1:9: package should be `postgres_test` instead of `postgres` (testpackage)
internal/adapters/driven/postgres/upstream_test.go:65:5: bereit is a global variable (gochecknoglobals)
internal/adapters/driven/postgres/upstream_test.go:184:5: beendet is a global variable (gochecknoglobals)
internal/adapters/driven/postgres/upstream_test.go:186:5: ergebnis is a global variable (gochecknoglobals)
internal/adapters/driven/recording/yaml_test.go:1:9: package should be `recording_test` instead of `recording` (testpackage)
internal/adapters/driven/recording/yaml_test.go:434:7: const syncGruppe is unused (unused)
internal/adapters/driven/postgres/upstream.go:94:25: unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _ (revive)
internal/adapters/driven/postgres/upstream.go:220:1: calculated cyclomatic complexity for function toResponse is 19, max is 15 (cyclop)
internal/adapters/driven/postgres/upstream.go:220:1: cyclomatic complexity 18 of func `toResponse` is high (> 15) (gocyclo)
internal/adapters/driven/recording/yaml.go:538:1: cognitive complexity 22 of func `fromDTO` is high (> 20) (gocognit)
```

Nachher (Exit 2, 60 Befunde im Repo, 12 weniger; in Testdateien unter beiden Pfaden
keiner, auch keiner der eigenen Prüfungen an `export_test.go`; die vier im Produkt-Code
bleiben für `slice-lint-bestand-kern-driven`, in `upstream.go` um sechs Zeilen
verschoben durch `newSession`):

```text
internal/adapters/driven/postgres/upstream.go:100:25: unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _ (revive)
internal/adapters/driven/postgres/upstream.go:226:1: calculated cyclomatic complexity for function toResponse is 19, max is 15 (cyclop)
internal/adapters/driven/postgres/upstream.go:226:1: cyclomatic complexity 18 of func `toResponse` is high (> 15) (gocyclo)
internal/adapters/driven/recording/yaml.go:538:1: cognitive complexity 22 of func `fromDTO` is high (> 20) (gocognit)
```

Behoben ohne `//nolint`, ohne `_ =` vor einem gemeldeten Fehler und ohne Änderung an
`.golangci.yml`: `errcheck` mit dem Helfer `schliesse` (prüft den Fehler von `Close` mit
`t.Error`), `gochecknoglobals` mit den Funktionen `bereit`, `beendet` und `ergebnis`, die
je Aufruf neue Werte liefern, `unused` durch Entfernen von `syncGruppe` (jede Verwendung
war eine lokale Variable gleichen Namens), `testpackage` durch die `_test`-Pakete.

*Mutationen* (`AGENTS.md` §3.10, Risiko *Prüfung geht im Umbau verloren*): je Mutant ein
frischer Pfad außerhalb des Repos (`cp -r` ohne `-p`), Tests per Bind-Mount im Image der
Stufe `deps` mit `go test -count=1 -run`; ohne Mutation ist jeder genannte Test grün, mit
ihr rot.

| Zusage | Mutation | roter Test |
|---|---|---|
| `Close` sendet Terminate mit Schreibfrist (Session über `NeueSession` auf `net.Pipe`) | `SetWriteDeadline` in `Close` entfernt | `TestCloseMitSchreibfrist` |
| scheitert Send nach gelesener ErrorResponse, ist es PGR-E6001 (Schreibhälfte über `Verbindung` geschlossen) | `Send` liefert beim Flush-Fehler PGR-E4003 statt `verbindungsende` | `TestSendNachFehlerantwort/ErrorResponse_gelesen` |
| Describe ohne gültige Zielart ist PGR-E1000 (über `ToFrontendMessage`) | Prüfung `!ok` der Zielart nie wahr | `TestReceiveFehler/Zielart` |
| ein unbekannter Nachrichtentyp ist PGR-E1000 (über `ToFrontendMessage`) | Zweig `default` liefert `Sync` ohne Fehler | `TestReceiveFehler/Zielart` |
| ein Fehler von `Close` ist ein Testfehler (Helfer `schliesse`) | `Close` liefert `net.ErrClosed` | `TestSendUndReceive`, `TestReceiveEndetMitReadyForQuery`, `TestReceiveFehler/COPY`, `TestReceiveFehler/Verbindungsende`, `TestSendUndReceiveGleichzeitig` |
| die letzte ErrorResponse zählt (Testdaten `beendet()` und `ergebnis()` als Funktionen) | `merke` behält die erste ErrorResponse | `TestFehlerantwortVorDemAbbruch/zwei_Fehlerantworten,_die_letzte_zählt` |
| der Aufbau liefert nur ParameterStatus und ReadyForQuery (Testdaten `bereit()` als Funktion) | `BackendKeyData` als Antwort im Aufbau | `TestOpenUndQuery` |
| eine fremde Formatkennung ist PGR-E3003 (`recording_test`) | Prüfung der Formatkennung in `Unmarshal` nie wahr | `TestUnmarshalFehler/fremdes_Format` |

Ohne roten Test, wie §6 sagt: `Open` gibt am ReadyForQuery eine neue Session aus
`newSession(conn)` zurück statt der, über deren Frontend es den Aufbau las; alle Tests
von `postgres` bleiben grün. Ob `Open` die Session über `newSession` anlegt, bleibt
Urteil des Review. Die Lesefrist, die `oeffne` über `Verbindung` setzt, verhindert nur
ein Hängen und hat keine eigene Mutation; dass `Verbindung` die Verbindung der Session
liefert, zeigt die zweite Zeile (das `CloseWrite` daran lässt `Send` scheitern).

*Läufe:* `go vet` und `go test` beider Pakete und `gofmt -l` (leer) im Image der Stufe
`deps`; `make lint` vorher und nachher wie oben; `make gates` grün am Stand dieses
Commits.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `175fd2d` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (12×, verkörpert in `AGENTS.md`
  §3.10) — der rote Fall für diese Pfade gehört zum Gate und liegt bei
  `slice-harness-lint`; dieser Slice liefert keinen neuen Vertrag.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — ein Risiko in §6.
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) betrifft diesen Slice nicht
  mehr, er führt keine Gegenprobe.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 nennt den
  Architect vor dem Code.
- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — die Brücke ist
  genau eine Randform, die ein Umbau still entscheiden würde; sie steht in §6.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
