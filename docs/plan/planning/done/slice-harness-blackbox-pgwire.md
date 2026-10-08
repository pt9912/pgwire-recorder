# Slice slice-harness-blackbox-pgwire: Black-Box-Tests des PGWire-Servers

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

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (Punkt 7 und 8: Brücke und dauerhafte Ausnahmen, gegen die umgestellt wird; der Slice ändert die Stelle nicht)

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

**Ziel:** Die Unit-Tests unter `internal/adapters/driving/pgwire` laufen als Black-Box-Pakete
(`package <name>_test`) und prüfen über die exportierte Schnittstelle des Pakets;
keine Testdatei unter diesen Pfaden hat in `make lint` einen Befund. Dafür behebt der
Slice in den Testdateien, die er umschreibt, die 6 Befunde aus der Messung in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Stand `79f40e1`; Aufteilung je Paketgruppe in der Fassung `92d1b86`
der ADR): `testpackage` 3, dazu `gochecknoglobals` und `revive`, zusammen 3. Wo ein Test einen unexportierten Teil braucht, geht
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

- Die Tests der übrigen Pakete — sie übernehmen `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`, `slice-harness-blackbox-einstieg`.
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
- Befunde im Produkt-Code dieser Pakete — übernimmt `slice-lint-bestand-driving` nach den
  Umstellungs-Slices.
- Produkt-Verhalten, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Testdateien; `.golangci.yml` und die Gegenprobe des Lint-Gates ändert er nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Die Testdateien unter `internal/adapters/driving/pgwire` gehören zu `_test`-Paketen; die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau gleich, keine Prüfung ist
      entfallen, die Abdeckungs-Deklarationen nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) stehen unverändert und
      `make abdeckung-check` ist grün. Unexportierte Teile erreicht ein Test nur über
      die Brücke nach `SPEC-049` Punkt 7. — Verifikation, Abschnitt 1 Punkt 1 und
      Abschnitt 2: die drei Testdateien sind `package pgwire_test`, im Paket des Codes
      liegt nur `export_test.go`; `go test -list .` an `6ff1b3b` und `1b6fa29` je 40
      Tests, `diff` leer, 53 Zeilen `--- PASS` an beiden Ständen; 25 Deklarationen Wort
      für Wort und in der Zuordnung gleich, `make abdeckung-check` Exit 0. Keine
      Prüfung entfallen bis auf den Rückgabewert von `weiterlesen`, den das Produkt nicht
      liest (G1, G2 äquivalent; zulässig nach `SPEC-049` Punkt 7, letzter Fall; §6
      Risiko *Prüfung geht im Umbau verloren*). Die Brücke reicht nur weiter; die
      White-Box-Zugriffe aus §7 sind auf anderem Weg nachgemessen.
- [x] `make lint` meldet in den Testdateien unter `internal/adapters/driving/pgwire` keinen Befund: Die
      6 Befunde der Messung sind behoben, ohne `//nolint` und ohne Änderung an
      `.golangci.yml`; die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8 blendet das
      Profil aus. Die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem Umbau
      stehen im Bericht; Befunde im Produkt-Code bleiben für `slice-lint-bestand-driving`.
      — Verifikation, Abschnitt 1 Punkt 2: die Zeilen vorher und nachher in §7
      *Belege des Implementers*, Wort für Wort bestätigt; 6 → 0 in `*_test.go`, 14 im
      Produkt-Code an beiden Ständen gleich; modulweit 60 → 54 ohne neuen Befund
      (`comm -13` leer); kein `//nolint`, `.golangci.yml` unverändert, kein neues `_ =`
      (je 44 Blank-Zuweisungen).
- [x] `make gates` grün. — an `1b6fa29` (Go-Code gleich `b37d6e7`; Verifikation,
      Abschnitt 1 Punkt 3 und Abschnitt 6); `edcb4a1` und diese Closure ändern nur
      Berichte, Pläne, Roadmap und das Register, `make docs-check` und `make kopf-check`
      grün am Stand dieser Closure.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      — `docs/reviews/2026-10-07-review-slice-harness-blackbox-pgwire.md` (bis
      `b37d6e7`; F-448 MEDIUM, F-449 LOW, F-450 bis F-453 INFO), Nacharbeit `1b6fa29`.
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
| `internal/adapters/driving/pgwire/server_test.go`, `server_extended_test.go`, `server_meldung_test.go` | refactor | `package pgwire_test`; die ersten beiden sind die größten Testdateien des Repos |
| `export_test.go` je Paket, nur wo nötig | neu | Brücke zu unexportierten Teilen nach `SPEC-049` Punkt 7 |
| dieselben Testdateien | update | übrige Befunde aus `make lint` (dazu `gochecknoglobals` und `revive`, zusammen 3) in den Dateien, die der Umbau ohnehin umschreibt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-blackbox-driven` liegt in `done/` (WIP-Limit 1,
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
- **Übrige Befunde der Testdateien** — `contextcheck`: Kontext aus `t.Context()`, in einer Funktion für `t.Cleanup` aus `context.WithoutCancel(t.Context())` (`SPEC-049` Punkt 5); `gochecknoglobals`: Testdaten in Funktionen oder als Konstanten; `revive`: nach der Meldung. Kein `_ =` vor einem Fehler, den
  `errcheck` meldet, und keine Ausnahme, die nur Bestand aussetzt (Entscheidung 5).
- **White-Box-Zugriffe im Bestand** — per AST gemessen am Stand `6ff1b3b`, Zuordnung je
  Zugriff in §7 (*Belege des Implementers*): über die Brücke oder über die exportierte
  Schnittstelle, kein Befund. Fristen und Weckmechanik folgen dem Punkt
  *Export-Test-Brücke*: `meldeFrist` reicht die Brücke als Konstante weiter, gelesen,
  nicht gesetzt; `richtungen` legt nur das Produkt an, `wecke` und `weiterlesen` prüft der
  Test über `Handle` an einer Verbindung, die er stellt. Grenze: Den Rückgabewert von
  `weiterlesen` verwirft `clientRichtung`; über die exportierte Schnittstelle ist er nicht
  zu sehen, ein Mutant nur an ihm ist dort äquivalent (§7, *Ohne roten Test*).
- **Fakes der Ports** — die Ports sind exportiert (`internal/hexagon/ports/...`); ein
  Fake in einem `_test`-Paket implementiert sie unverändert. Ein Fake, der auf
  unexportierte Felder eines Produkt-Typs greift, fällt unter die Brücke.
- **Mutationstests auf unexportierte Teile** — eine Mutation, die bisher ein Test auf
  eine unexportierte Funktion fing, muss auch nach dem Umbau fangen; sonst fehlt eine
  Prüfung, obwohl die Testliste gleich ist.
- **Eigener Schnitt** — pgwire trägt die meisten Zugriffe; darum ein Slice für sich.
  Reicht auch das nicht für eine Review-Sitzung, gilt die Rückführung aus §4 (je
  Testdatei).

**Risiken:**

- **Prüfung geht im Umbau verloren** — ein Test bleibt in der Liste, prüft aber
  weniger, weil ein unexportierter Vergleich wegfiel
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Je umgeschriebener (nicht nur
  umgestellter) Test eine Mutation, die er weiter fängt. — **Ausgang:** entfallen: Die
  Testliste ist gleich (40 Tests, 53 Zeilen `--- PASS` an beiden Ständen), die 25
  Deklarationen sind gleich zugeordnet, die 15 Mutationen der Tabelle in §7 sind in der
  Verifikation rot aus dem genannten Grund (Abschnitt 4), ebenso R2 bis R4 des Reviews.
  Weggefallen ist genau eine Prüfung: der **Rückgabewert von `weiterlesen`**, den der
  alte `TestWeiterlesenNachWecken` am selbst angelegten `richtungen`-Wert las. Das
  Produkt liest ihn nicht (`clientRichtung` ruft `weiterlesen` als Anweisung), die
  Mutanten G1 und G2 der Verifikation (Review R1) sind über die Schnittstelle
  äquivalent, und nach `SPEC-049` Punkt 7, letzter Fall, ist die Prüfung damit
  zulässig weggefallen (Verifikation, Abschnitt 3; F-449, V-86). Kein Verlust einer
  Prüfung, die Verhalten fängt; die Grenze steht in §6 *White-Box-Zugriffe im Bestand*,
  und ihr Auslöser ist an `slice-lint-bestand-driving` §6 gegeben (§7, *Folge-Slices*).
- **Abdeckung sinkt** — Black-Box-Tests erreichen unexportierte Pfade seltener; die
  Zahl misst erst `slice-harness-coverage` nach allen vier Umstellungs-Slices.
  — **Ausgang:** entfallen: Der Umbau gibt keinen Pfad des Produkts ab. Jeder frühere
  Zugriff auf Unexportiertes geht über die Brücke an dieselbe Funktion oder Konstante
  (Verifikation, Abschnitt 2), die vier `&Server{log: …}` werden `NewRecordServer(nil,
  …)` mit demselben Zustand, und den einzigen umgeschriebenen Test fährt jetzt das
  Produkt über `Handle` durch `recordSitzung`, `wecke` und `weiterlesen`. Nicht mehr
  erreicht ist nur die Rückgabe von `weiterlesen`, die kein Produkt-Pfad liest. Die
  Zahl für alle vier Umstellungs-Slices misst `slice-harness-coverage` nach seinem
  eigenen §1, nicht als Ausgang dieses Risikos.
- **Befund verdeckt statt behoben** — `_ =` vor einem ungeprüften Fehler, eine Globale
  als Funktion mit demselben geteilten Zustand: `make lint` ist grün, der Test nicht
  besser. Review prüft die Form, das Werkzeug nur die Zahl (Grenze von `SPEC-049`).
  — **Ausgang:** entfallen: Kein neues `_ =` (je 44 Blank-Zuweisungen an beiden Ständen;
  `_ = fe.Flush()` ist aus `sende` übernommen), kein `//nolint`, `.golangci.yml`
  unverändert (Verifikation, Abschnitt 1 Punkt 2). `nebenher(fe)` legt je Aufruf eine
  eigene Sperre an, kein geteilter Zustand über Verbindungen oder Tests (F-452;
  Verifikation, Abschnitt 7: 21 Aufrufe, kein `t.Parallel`); die Sperre war vorher
  ebenso ungeprüft (§7, *Ohne roten Test*).

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
  Brücke, kein Produkt-Code, ein umgeschriebener Test, keine Rückführung aus §4 und keine
  Entscheidung des Architect während der Arbeit. Die Brücken-Regel aus `SPEC-049` Punkt 7
  in der Fassung `f0ea7be` trug ohne Nachfrage: Der einzige Wert eines unexportierten Typs
  (`richtungen`) ging an die Schnittstelle, `TestWeiterlesenNachWecken` fährt die
  Weckmechanik jetzt über `Handle` an einer vom Test gestellten Verbindung. Belege in §7:
  AST-Messung, Testliste, Lint-Zeilen vorher und nachher, 15 Mutationen und die grünen
  Mutanten unter *Ohne roten Test* standen dort; Review und Verifikation prüften gegen sie,
  die Verifikation maß die White-Box-Zugriffe auf anderem Weg nach (Compiler-Fehler unter
  Punkt-Import) und fand dieselben Zahlen. Die Einordnung grüner Mutanten nach
  `.claude/commands/implement-slice.md` Schritt 19 (seit slice-harness-blackbox-driven) hat
  für zwei von drei gewirkt: `meldeFrist = 3 * time.Second` und `nebenher` ohne Sperre
  standen mit Test-Idee und Grenze in §7, beide hat das Review bestätigt (F-451, F-452).
- **Was ging anders als geplant:** Zwei Punkte, keine Nacharbeit am Code. §6 nannte die
  Frage, ob die Brücke eine Frist schreiben darf, weiter offen, obwohl der Punkt
  *Export-Test-Brücke* (seit `26c13aa`) und der Code sie beantworteten; `b37d6e7` zog im
  Plan nur §7 nach (F-448 MEDIUM, drittes Auftreten der Klasse in Reviews). Und ein dritter
  grüner Mutant, der Rückgabewert von `weiterlesen`, war in §7 nicht eingeordnet, weil der
  Implementer ihn nicht mutiert hatte; das Review fand ihn (F-449, R1). Beides behob
  `1b6fa29` nur im Plan; die Verifikation bestätigt §6 widerspruchsfrei und die Einordnung
  als äquivalent (Abschnitt 3).
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-07-review-slice-harness-blackbox-pgwire.md`: „0 HIGH · 1 MEDIUM
    (F-448: §6 nennt die Frage, ob die Brücke schreiben darf, weiter offen, obwohl Punkt
    *Export-Test-Brücke* und Diff sie beantworten) · 1 LOW (F-449: Rückgabewert von
    `weiterlesen` im Umbau ungeprüft, grüner Mutant nicht in §7 eingeordnet) · 4 INFO
    (F-450 Folge in `TestWeiterlesenNachWecken` praktisch deterministisch, eine verdeckte
    Abhängigkeit vom Startphasen-Wächter; F-451 Wert von `meldeFrist` ist eine Spec-Lücke
    aus dem Bestand; F-452 `nebenher` gleichwertig; F-453 Mutationen reichen bis auf
    F-449). Wiederkehrende Klasse: „Plan folgt Korrektur nicht“
    (`BEO-REPO/plan-folgt-korrektur-nicht`, drittes Auftreten in Reviews).“
    Verifikation `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-pgwire.md`,
    Urteil: drei Liefer-Punkte erfüllt und selbst belegt, Brücke hält Punkt 7, die 15
    Mutationen aus §7 rot aus dem genannten Grund, die grünen Mutanten tragfähig
    eingeordnet, F-448 behoben, F-449 trägt, F-450 ohne Befund (2000 grüne Läufe unter
    Drossel), F-451 bestätigt als Lücke aus dem Bestand; kein Befund blockiert die Closure,
    V-86 und V-87 sind Vorschläge für Risiko-Ausgang, Folge-Slices und Lerneintrag.
    `make gates` grün an `1b6fa29`.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke: Die Spezifikation nennt keinen Wert für
  die Frist, mit der der Recorder beim Ende einer Session die Fehlerantwort an einen nicht
  lesenden Client schreibt (`meldeFrist` in `internal/adapters/driving/pgwire/server.go`,
  im Code eine Sekunde). Der Absatz *Abbruch* im Record-Teil sagt nur, dass das Schließen
  der Client-Verbindung ein blockiertes Schreiben beendet; das Lastenheft nennt die Frist
  ebenfalls nicht. Den Wert entscheidet bisher allein der Code, seit
  `slice-extended-query-record` (`09cc0e9`, dort Verifikation V-22); der Mutant
  `meldeFrist = 3 * time.Second` ist an `6ff1b3b` und an `1b6fa29` grün (Verifikation G3).
  Ein Test mit der Schranke als Literal schriebe einen nicht spezifizierten Wert fest; er
  braucht zuerst die Stelle in der Spezifikation (F-451, V-87). Adresse:
  `slice-v1-abschluss-betrieb` §6, als Randform mit Herkunfts-Anker
  `seit slice-harness-blackbox-pgwire`: Der Architect entscheidet den Wert vor dem Code in
  der Spezifikation, danach folgt der Test mit der Schranke als Literal. Kein neuer Beleg
  im Register: Der Fund steht als V-22 schon in
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und zählte sonst doppelt.
  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) — wieder aufgetreten
  (F-448), die Regel bleibt; `make kopf-check` liest nur den Kopf gegen §1 und §2, einen
  Widerspruch innerhalb von §6 fängt er nicht und soll er nach seinem Vertrag nicht.
  `.claude/commands/implement-slice.md` Schritt 19 (seit slice-harness-blackbox-driven) —
  die Einordnung trug für die gefundenen grünen Mutanten; F-449 ist ein Mutant, den der
  Implementer nicht gefahren hat, kein falsch eingeordneter; die Regel bleibt. §3.12 und
  die Randform-Rückgabe: ohne Befund im Code-Commit (Review, Negativbefund), bleiben.
  §3.10 ohne neuen Vertrag, §3.11 ohne Befund.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `edcb4a1` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/plan-folgt-korrektur-nicht` — **Beleg**, 14× → 15× (F-448), verkörpert in
    `AGENTS.md` §3.9 und `make kopf-check`; die Schwelle lag schon vorher darüber.
  - Ohne Beleg: `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×; die
    weggefallene Prüfung ist ein Test-Vergleich ohne Produkt-Leser und zulässig nach
    `SPEC-049` Punkt 7, keine Gate-Regel), `BEO-REPO/randform-wellenlos-ohne-architect-vor-code`
    (bleibt 1×; die Randformen standen vor dem Code entschieden in §6, Review Hard Rule
    3.12 ohne Befund im Code-Commit), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`
    (bleibt 3×, verkörpert; keine Randform im Code entschieden),
    `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (bleibt 3×, verkörpert; die
    Belege in §7 erreichten beide Prüfer), `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`
    (bleibt 1×; jede Mutation lief in einem frischen Pfad per Bind-Mount),
    `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (bleibt 16×; Review Hard Rule 3.11
    ohne Befund; V-22 dort ist der Fund zur Frist, kein neuer Beleg),
    `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (bleibt 13×; kein neuer Vertrag),
    `BEO-REPO/spec-randform-erst-im-review-entschieden` (bleibt 11×; F-451 entscheidet
    keine Randform, sondern benennt eine aus dem Bestand).

  Einmalig und nicht eingetragen: F-449 als Klasse (ein grüner Mutant, den der Implementer
  nicht gefahren hat; Schritt 19 verlangt die Einordnung gefundener Mutanten, nicht eine
  vollständige Suche), F-450 (verdeckte Abhängigkeit vom Startphasen-Wächter, 2000 grüne
  Läufe), F-452 (keine Aktion), F-453 (Mutationen reichen). Mit diesem Slice erreicht
  **kein** Eintrag die Schwelle 3× neu; über der Schwelle stehen nur verkörperte Einträge
  (`implementer-bericht-erreicht-pruefer-nicht` 3×,
  `randform-im-code-entschieden-dann-zurueckgegeben` 3×,
  `spec-randform-erst-im-review-entschieden` 11×,
  `negativtests-fehlen-bei-neuem-vertrag` 13×, `plan-folgt-korrektur-nicht` 15×,
  `zusage-im-kommentar-weiter-als-pruefung` 16×).
- **Folge-Slices:** `slice-v1-abschluss-betrieb` (Wert von `meldeFrist`, Randform in §6:
  Architect entscheidet den Wert in der Spezifikation, dann Test mit der Schranke als
  Literal; nimmt an: liegt in `open/`, setzt die Frist in `recordSitzung` und das Ende
  von Sessions unter `--shutdown-timeout` um, §1 schließt weder Spezifikation noch Tests
  an `pgwire` aus), `slice-lint-bestand-driving` (die 14 Befunde im Produkt-Code dieses
  Pakets, Zeilen in §7 *Nachher*; dazu in §6 die Grenze zum Rückgabewert von
  `weiterlesen` beim Zerlegen von `clientRichtung`, V-86), als nächster in der Reihe
  `slice-harness-blackbox-einstieg`, dann `slice-harness-coverage` und
  `slice-harness-lint`.
- **Risiken aus §6:** drei, alle **entfallen**, mit Begründung: *Prüfung geht im Umbau
  verloren* (einzige weggefallene Prüfung ist der Rückgabewert von `weiterlesen`, G1/G2
  äquivalent, zulässig nach `SPEC-049` Punkt 7, letzter Fall), *Abdeckung sinkt* (kein
  Produkt-Pfad abgegeben), *Befund verdeckt statt behoben* (kein neues `_ =`, kein
  `//nolint`, Sperre je Verbindung). Die Randformen in §6 sind Entscheidungen, keine
  Risiken.
- **Drei Paarungen:** Anker — der Lerneintrag ist eine benannte Spec-Lücke, keine
  verkörperte Regel, daher kein `liegt in`; der Herkunfts-Anker
  `seit slice-harness-blackbox-pgwire` steht an der Adresse,
  `grep -n "seit slice-harness-blackbox-pgwire" docs/plan/planning/done/slice-v1-abschluss-betrieb.md`
  findet ihn in §6. Folge-Slice — `slice-v1-abschluss-betrieb`,
  `slice-lint-bestand-driving`, `slice-harness-blackbox-einstieg` und
  `slice-harness-coverage` liegen in `open/`, `slice-harness-lint` in `next/`;
  `slice-v1-abschluss-betrieb` und `slice-lint-bestand-driving` nennen ihren Punkt in §6.
  Register — `BEO-REPO/plan-folgt-korrektur-nicht` trägt
  `evidence/slice-harness-blackbox-pgwire.md`; die acht ohne Beleg bestehen als
  Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-07-review-slice-harness-blackbox-pgwire.md`
  (bis `b37d6e7`; F-448 bis F-453), Nacharbeit `1b6fa29`, Verifikation
  `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-pgwire.md` (bis `1b6fa29`,
  `make gates` an `1b6fa29`; V-86, V-87), Entscheidung
  [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) unverändert.
  Validierung: n/a, der Slice ändert Tests, kein End-Nutzer-Verhalten.

**Belege des Implementers** (Arbeitsbaum auf `6ff1b3b` mit dem Diff des Commits, der
diesen Abschnitt anlegt):

*White-Box-Zugriffe, per AST gemessen* (go/types im Image der Stufe `deps`, ohne Netz;
gezählt sind Bezeichner in Testdateien des Pakets `pgwire`, die auf ein unexportiertes
Objekt einer Produkt-Datei desselben Pakets zeigen, dazu Literale, `var`, `new` und
`:=` unexportierter Typen in Testdateien). Vorher (`6ff1b3b`), je Zugriff mit Zuordnung:

| Zugriff | Stellen | Zuordnung |
|---|---|---|
| `(*Server).handle` | 10 | Brücke `Handle` (reicht an `handle` des übergebenen, von `NewRecordServer` oder `NewReplayServer` erzeugten `Server` weiter; die Verbindung stellt der Test) |
| `(*Server).fail`, `(*Server).note` | 2, 3 | Brücke `Fail`, `Note` (an einem übergebenen `Server`) |
| Feld `log` in `&Server{log: …}` | 4 | exportierte Schnittstelle: `NewRecordServer(nil, …)`, derselbe Zustand ohne Use Case |
| `toClientMessage`, `toMessage` | 4, 4 | Brücke `ToClientMessage`, `ToMessage` |
| `codeSSLRequest`, `codeGSSEncRequest`, `codeCancelRequest`, `majorSpezial` | je 1 | Brücke, Konstanten |
| `meldeFrist` | 2 | Brücke, Konstante `MeldeFrist`; gelesen, nicht gesetzt |
| Typ `richtungen`, Felder `conn` und `weck`, Methoden `wecke` und `weiterlesen` | 1, 1, 1, 1, 2 | exportierte Schnittstelle (`SPEC-049` Punkt 7, letzter Fall): Den Wert legte `TestWeiterlesenNachWecken` als Literal an; das Produkt erzeugt ihn nur in `recordSitzung`, die Brücke bekommt ihn weder übergeben noch lässt sie ihn erzeugen. Der Test ist umgeschrieben (unten) |

Werte unexportierter Typen in Testdateien vorher 2 (`richtungen`-Literal und die
Variable `r`, beide in `TestWeiterlesenNachWecken`). Nachher: die zehn Objekte aus den
Zeilen der Brücke stehen je einmal in `export_test.go` und sonst nirgends; 0 Werte
unexportierter Typen. Kein Zugriff blieb Befund; kein Test setzt eine Frist oder weckt an
einem Wert, den er selbst angelegt hat: `TestWeiterlesenNachWecken` weckt über das Ende
von `ctx` an der Session, die `Handle` aus einer vom Test gestellten Verbindung
(`fristSpion`, ein `net.Conn`) erzeugt. Fehltreffer der Suche in §6: `fehler` (lokale
Variable), `startup` (Testhelfer).

Die Brücke `internal/adapters/driving/pgwire/export_test.go` hat fünf Konstanten
(`CodeCancelRequest`, `CodeSSLRequest`, `CodeGSSEncRequest`, `MajorSpezial`,
`MeldeFrist`) und fünf Funktionen (`Handle`, `Fail`, `Note`, `ToClientMessage`,
`ToMessage`), die nur weiterreichen. Kein Produkt-Code geändert.

*Testliste* (`go test -list .`, Image der Stufe `deps`): vorher und nachher 40 Tests,
`diff` leer. Die Abdeckungs-Deklarationen stehen unverändert (der Diff berührt keine
Kommentarzeile `Abdeckung:` und keinen Kommentar darunter); `make abdeckung-check` grün.

*`make lint`.* Vorher (`6ff1b3b`, Exit 2): golangci-lint `60 issues:` im Modul, davon
20 unter `internal/adapters/driving/pgwire`, 6 in dessen Testdateien:

```text
internal/adapters/driving/pgwire/server_extended_test.go:42:5: sendeMu is a global variable (gochecknoglobals)
internal/adapters/driving/pgwire/server_extended_test.go:25:37: context-as-argument: context.Context should be the first parameter of a function (revive)
internal/adapters/driving/pgwire/server_extended_test.go:313:35: context-as-argument: context.Context should be the first parameter of a function (revive)
internal/adapters/driving/pgwire/server_extended_test.go:1:9: package should be `pgwire_test` instead of `pgwire` (testpackage)
internal/adapters/driving/pgwire/server_meldung_test.go:1:9: package should be `pgwire_test` instead of `pgwire` (testpackage)
internal/adapters/driving/pgwire/server_test.go:1:9: package should be `pgwire_test` instead of `pgwire` (testpackage)
```

Nachher (Arbeitsbaum, Exit 2): `54 issues:` im Modul, 6 weniger; unter
`internal/adapters/driving/pgwire` 14, alle im Produkt-Code und dieselben Zeilen wie
vorher (`diff` der übrigen Befunde leer), keiner in einer Testdatei und keine
`lint:`-Zeile der eigenen Prüfungen (auch nicht an `export_test.go`). Die 14 bleiben für
`slice-lint-bestand-driving`:

```text
internal/adapters/driving/pgwire/server.go:302:2: found a struct that contains a context.Context field (containedctx)
internal/adapters/driving/pgwire/server.go:440:47: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/pgwire/server.go:462:34: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/pgwire/server.go:474:40: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/pgwire/server.go:182:1: calculated cyclomatic complexity for function replaySitzung is 21, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:435:1: calculated cyclomatic complexity for function clientRichtung is 17, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:569:1: calculated cyclomatic complexity for function startup is 17, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:739:1: calculated cyclomatic complexity for function toMessage is 20, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:182:1: cognitive complexity 37 of func `(*Server).replaySitzung` is high (> 20) (gocognit)
internal/adapters/driving/pgwire/server.go:435:1: cognitive complexity 22 of func `(*richtungen).clientRichtung` is high (> 20) (gocognit)
internal/adapters/driving/pgwire/server.go:182:1: cyclomatic complexity 20 of func `(*Server).replaySitzung` is high (> 15) (gocyclo)
internal/adapters/driving/pgwire/server.go:739:1: cyclomatic complexity 19 of func `toMessage` is high (> 15) (gocyclo)
internal/adapters/driving/pgwire/server.go:62:22: net.Listen must not be called. use (*net.ListenConfig).Listen (noctx)
internal/adapters/driving/pgwire/server.go:728:7: unused-receiver: method receiver 's' is not referenced in method's body, consider removing or renaming it as _ (revive)
```

Behoben ohne `//nolint`, ohne neues `_ =` und ohne Änderung an `.golangci.yml`:
`testpackage` durch `package pgwire_test`; `gochecknoglobals` (`sendeMu`) durch die
Funktion `nebenher(fe)`, die je Aufruf eine eigene Sperre anlegt und eine Funktion
`sende` für diese eine Verbindung liefert, statt einer Sperre für alle Verbindungen des
Pakets; `revive` `context-as-argument` durch `ctx` als ersten Parameter von
`verbindeExtended` und `verbindeReplay`. Das `_ = fe.Flush()` in `nebenher` ist die Zeile
aus `sende`, unverändert übernommen; im neuen Test ist der Fehler von `SetDeadline`
geprüft.

*Umgeschrieben:* nur `TestWeiterlesenNachWecken`. Er startet eine Record-Session über
`Handle` auf `fristSpion`, der jedes Setzen der Lesefrist (`frist`, `zurück`) und jedes an
ihr gescheiterte Lesen (`abgelaufen`) protokolliert, beendet `ctx` bei `Shutdown` ohne
Freigabe und erwartet genau `frist abgelaufen abgelaufen zurück`, danach das Lesen der
nächsten Nachricht (`c:parse`); 20 Läufe mit `-count=20` grün. Alle übrigen Tests sind umgestellt (Paketname,
Brücke, Präfix `pgwire.`), ihre Prüfungen unverändert.

*Mutationen* (`AGENTS.md` §3.10, Risiko *Prüfung geht im Umbau verloren*): je Mutant ein
frischer Pfad außerhalb des Repos (`cp -r` ohne `-p` aus dem Arbeitsbaum), Tests per
Bind-Mount im Image der Stufe `deps` mit `go test -count=1 -run`; ohne Mutation ist jeder
genannte Test grün, mit ihr rot.

| Zusage | Mutation | roter Test |
|---|---|---|
| nach dem Wecken lässt `weiterlesen` die Lesefrist stehen (umgeschrieben) | `weiterlesen` setzt die Frist auch nach dem Wecken zurück | `TestWeiterlesenNachWecken` (`frist abgelaufen zurück`) |
| dasselbe, das Signal wird nicht verbraucht | `weiterlesen` ohne `select`, setzt immer zurück | `TestWeiterlesenNachWecken` |
| das Wecken hinterlegt ein Signal | `wecke` ohne Senden in `weck` | `TestWeiterlesenNachWecken` |
| ohne Wecken setzt `weiterlesen` die Frist zurück und die Richtung liest weiter | `weiterlesen` setzt ohne Signal nicht zurück | `TestWeiterlesenNachWecken` (kein `zurück` binnen 2 s) |
| Parameterwerte werden kopiert (über `ToClientMessage`) | `Bind` übernimmt den Puffer der Bibliothek | `TestToClientMessage` |
| eine fremde Zielart ist ein Fehler (über `ToClientMessage` und `Handle`) | `zielart` liefert im `default` `TargetStatement` | `TestToClientMessage`, `TestExtendedNichtUnterstuetzt` |
| ein unbekannter Antworttyp wird nicht erfunden (über `ToMessage`) | `default` liefert `NoData` | `TestToMessageErfindetNichts` |
| ParameterDescription trägt die Typen (über `ToMessage`) | `ParameterOIDs: nil` | `TestToMessageExtended` |
| SQLSTATE nach Klasse (über `Fail`) | Klasse 6 liefert `XX000` | `TestFehlerantwortJeKlasse`, `TestFailGleichrangig` |
| genau eine ErrorResponse (über `Fail`) | `fail` sendet sie zweimal | `TestFailGleichrangig` |
| der erste Fehler zählt (über `Note`) | `note` überschreibt `firstCode` immer | `TestErsterFehlerZaehlt` |
| je Meldung eine Log-Zeile (über `Note`) | `note` loggt nur die erste | `TestNoteGleichrangig` |
| die Fehlerantwort beim Ende schreibt höchstens `MeldeFrist` lang | `SetWriteDeadline` in `beende` entfernt | `TestFehlerantwortMitFrist` |
| SSLRequest wird erkannt (Konstante über die Brücke) | `codeSSLRequest = 80877199` | `TestSSLUndGSSMitN` (die Bibliothek dekodiert den echten Code) |
| Client-Nachrichten gehen mit ihren Feldern an den Use Case (über `Handle` und `nebenher`) | `Parse` trägt als `Statement` die Anfrage | `TestExtendedSyncGruppe` |

*Ohne roten Test*, eingeordnet nach `.claude/commands/implement-slice.md` Schritt 19:

- Der Rückgabewert von `weiterlesen`, beide Richtungen: `return true` → `return false`
  nach dem Wecken und `return false` → `return true` ohne Signal. Beide sind an
  `b37d6e7^` rot (`TestWeiterlesenNachWecken`, „Wecken nicht gemeldet“ und „Wecken ohne
  Signal gemeldet“) und an `b37d6e7` in allen Tests des Pakets grün (Review F-449, R1;
  hier nachgefahren, am Altstand aus `git archive` in einem frischen Pfad). Äquivalent
  über die Schnittstelle: `clientRichtung` ruft `weiterlesen` als Anweisung und verwirft
  den Wert, Lesefrist und Lesen bleiben gleich. Nach `SPEC-049` Punkt 7, letzter Fall,
  ist die Prüfung damit zulässig weggefallen. Grenze in §6 (*White-Box-Zugriffe im
  Bestand*): Liest das Produkt den Wert künftig, braucht er einen Test über die
  Schnittstelle.
- `meldeFrist = 3 * time.Second` bleibt grün: `TestFehlerantwortMitFrist` wartet
  `MeldeFrist` plus eine Sekunde, vorher wie nachher über dieselbe Konstante. Der Mutant
  ändert das Verhalten (die Dauer), über die Schnittstelle fangbar nur gegen einen festen
  Wert; weder Spezifikation noch Lastenheft nennen ihn. Test-Idee: die Schranke als
  Literal. Grenze: Der Wert ist nicht spezifiziert, §1 schließt geänderte Erwartungen
  aus; die Closure gibt der Idee eine Adresse.
- `nebenher` ohne Sperre (Lock und Unlock entfernt) bleibt grün, 10 Läufe: Ohne Sperre
  schreibt ein zweiter Aufruf in den Puffer von `fe`, während der vorige Flush läuft. Das
  ändert das Verhalten des Testhelfers nur bei gleichzeitigem Zugriff; die Tests warten
  zwischen zwei Aufrufen auf den Use Case. Test-Idee: `go test -race`. Grenze: Das Image
  baut mit `CGO_ENABLED=0` und hat keinen C-Compiler, der Race-Detektor läuft dort
  nicht; die Sperre stand vorher ebenso ungeprüft in `sendeMu`.

*Läufe:* `gofmt -l` (leer), `go vet` und `go test` des Pakets im Image der Stufe `deps`;
`make lint` vorher und nachher wie oben; `make abdeckung-check` grün; `make gates` grün am
Stand dieses Commits.

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
