# Slice slice-harness-blackbox-kern: Black-Box-Tests im Kern (Domain Model und Services)

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

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (Punkt 7 und 8: Brücke und dauerhafte Ausnahmen, gegen die umgestellt wird. Geändert vom Architect unter der Kennung dieses Slice: Punkt 7 und *Grenze*, kein Wert eines unexportierten Typs aus Brücke oder Test, nach F-441/F-442; Punkt 5, `contextcheck` in Testdateien, nach V-83)

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

**Ziel:** Die Unit-Tests unter `internal/hexagon/model` und `internal/hexagon/services` laufen als Black-Box-Pakete
(`package <name>_test`) und prüfen über die exportierte Schnittstelle des Pakets;
keine Testdatei unter diesen Pfaden hat in `make lint` einen Befund. Dafür behebt der
Slice in den Testdateien, die er umschreibt, die 31 Befunde aus der Messung in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Stand `79f40e1`; Aufteilung je Paketgruppe in der Fassung `92d1b86`
der ADR): `testpackage` 9, `contextcheck` 5, `gochecknoglobals` 17. Wo ein Test einen unexportierten Teil braucht, geht
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

- Die Tests der übrigen Pakete — sie übernehmen `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`.
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
- Produkt-Verhalten, Lastenheft und die Spezifikation außer `SPEC-049` Punkt 5
  (`contextcheck`), Punkt 7 und *Grenze* — Schicht-Abgrenzung: Der Slice ändert
  Testdateien; die Stellen in `SPEC-049` hat der Architect nach F-441/F-442 und V-83
  entschieden, damit die Umstellungs-Slices sie nicht je für sich auslegen.
  `.golangci.yml` und die Gegenprobe des Lint-Gates ändert er nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die Testdateien unter `internal/hexagon/model` und `internal/hexagon/services` gehören zu `_test`-Paketen; die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau gleich, keine Prüfung ist
      entfallen, die Abdeckungs-Deklarationen nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) stehen unverändert und
      `make abdeckung-check` ist grün. Unexportierte Teile erreicht ein Test nur über
      die Brücke nach `SPEC-049` Punkt 7.
- [ ] `make lint` meldet in den Testdateien unter `internal/hexagon/model` und `internal/hexagon/services` keinen Befund: Die
      31 Befunde der Messung sind behoben, ohne `//nolint` und ohne Änderung an
      `.golangci.yml`; die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8 blendet das
      Profil aus. Die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem Umbau
      stehen im Bericht; Befunde im Produkt-Code bleiben für `slice-lint-bestand-kern-driven`.
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
| `internal/hexagon/model/extended_test.go`, `fehler_test.go` | refactor | `package model_test` |
| `internal/hexagon/services/*_test.go` (sieben Dateien) | refactor | `package services_test`; Fakes der Ports und Testhelfer wandern mit |
| `internal/hexagon/model/export_test.go` | neu | Brücke nach `SPEC-049` Punkt 7: `Klasse` reicht an `klassen` weiter |
| `internal/hexagon/services/export_test.go` | neu | Brücke nach `SPEC-049` Punkt 7: `IstLebendpruefung`, `VTLeerraum`, `Abweichung`, `ParameterStelle` reichen weiter; keine Brücke für `letzteNummer` (§6 *Brücke mit Cursor*) |
| `internal/hexagon/services/replay_lebendpruefung_test.go` | update | `TestReplayLetzteNummer` prüft die letzte Nummer über `ReplayService` (§6 *Brücke mit Cursor*) |
| dieselben Testdateien | update | übrige Befunde aus `make lint` in den Dateien, die der Umbau ohnehin umschreibt: `contextcheck` 5 (je Subtest `ctx := t.Context()` in `TestRecordEndeNachErfolgreichemUpstream` statt des geteilten neuen Kontexts), `gochecknoglobals` 17 (Testdaten als Funktionen, die je Aufruf den Wert liefern) |
| `internal/hexagon/model/fehler_test.go` | update | `model.Meldung`-Literale mit Feldnamen: `go vet` (`composites`) lehnt ungeschlüsselte Literale eines fremden Pakets ab |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint-werkzeug` und
`slice-harness-upgrade-v6-16` liegen in `done/` (WIP-Limit 1,
Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05 und 2026-10-06:
`slice-harness-lint-werkzeug`, `slice-harness-upgrade-v6-16`, die vier Umstellungs-Slices,
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
  eines unexportierten Typs legt nur der Produkt-Code an, weder Brücke noch Test
  (entschieden vom Architect am 2026-10-06 nach F-441/F-442, `SPEC-049` Punkt 7).
- **Übrige Befunde der Testdateien** — `contextcheck`: `t.Context()` statt eines neuen Kontexts (`SPEC-049` Punkt 5, entschieden vom Architect am 2026-10-06 nach V-83; die `context.WithoutCancel` aus Entscheidung 5 gilt für den Produkt-Code); `gochecknoglobals`: Testdaten und Fakes in Funktionen oder als Konstanten. Kein `_ =` vor einem Fehler, den
  `errcheck` meldet, und keine Ausnahme, die nur Bestand aussetzt (Entscheidung 5).
- **White-Box-Zugriffe im Bestand** — Namensabgleich per Suche am Stand `ce50a10`
  (ungemessen, kann Fehltreffer enthalten, wo ein Testhelfer gleich heißt): in `model` keine; in `services` u. a. `abweichung`, `cursor`, `istLebendpruefung`, `vtLeerraum`, `laufende`, `letzteNummer`, `mitten`.
  Je Zugriff: über die exportierte Schnittstelle prüfbar, über die Brücke, oder
  Befund. Gemessen beim Umbau (Bezeichner der Testdateien gegen die Paketebene des
  Produkt-Codes, per AST): in `model` `klassen` (Fehltreffer der Suche), in `services`
  `istLebendpruefung`, `vtLeerraum`, `abweichung`, `parameterStelle` über die Brücke,
  `cursor` mit `letzteNummer` über die exportierte Schnittstelle (*Brücke mit Cursor*);
  keiner als Befund. `laufende` und `mitten`
  waren Fehltreffer (Kommentartext). Die Services-Tests prüfen Zwischenstände des Replays (Cursor, laufende Nummer); ob das Verhalten über `ReplayService` allein beobachtbar ist, entscheidet je Fall die Brücke.
- **Fakes der Ports** — die Ports sind exportiert (`internal/hexagon/ports/...`); ein
  Fake in einem `_test`-Paket implementiert sie unverändert. Ein Fake, der auf
  unexportierte Felder eines Produkt-Typs greift, fällt unter die Brücke.
- **Mutationstests auf unexportierte Teile** — eine Mutation, die bisher ein Test auf
  eine unexportierte Funktion fing, muss auch nach dem Umbau fangen; sonst fehlt eine
  Prüfung, obwohl die Testliste gleich ist.
- **Brücke mit Cursor** — entschieden vom Architect am 2026-10-06 nach F-441/F-442
  (`SPEC-049` Punkt 7): Die Brücke legt keinen `cursor` an; die Funktion
  `LetzteNummer` entfällt aus `export_test.go`. Der Produkt-Code legt seine Cursor mit
  `status: "I"` an, ein Cursor nur mit `session` ist ein Zustand, den das Produkt nie
  erzeugt. `TestReplayLetzteNummer` bleibt unter seinem Namen und prüft über
  `ReplayService`: Session mit den Interaktionen 2 und 4, beide Anfragen gesendet, eine
  dritte Anfrage ist `PGR-E5001` mit „nach Interaktion 4 erwartet die Aufzeichnung keine
  weitere“ (wie `TestReplayLebendpruefungNummerNachDemEnde` für den Extended-Weg). Ohne
  Abdeckungs-Deklaration, wie bisher. Rot sein muss danach die Mutation „Zahl der
  Interaktionen statt Nummer“ und „erste statt letzte Interaktion“ in `letzteNummer`.
  **Akzeptiertes Negativ:** Der Fall „ohne Interaktion 0“ entfällt. Das Produkt ordnet
  keine Session ohne erwartete Interaktion zu (`NewReplayService` nimmt sie nicht auf),
  der Zweig ist über die Schnittstelle nicht erreichbar; der Fall prüfte einen Zustand,
  den das Produkt nicht erzeugt, und ist keine entfallene Prüfung im Sinn der DoD.
  Den Zweig selbst lässt dieser Slice stehen (kein Produkt-Code, §1); er hat keinen
  eigenen Folge-Slice, die Abdeckung misst `slice-harness-coverage`.

**Risiken:**

- **Prüfung geht im Umbau verloren** — ein Test bleibt in der Liste, prüft aber
  weniger, weil ein unexportierter Vergleich wegfiel
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Je umgeschriebener (nicht nur
  umgestellter) Test eine Mutation, die er weiter fängt. — **Ausgang:** — (bei Closure)
- **Abdeckung sinkt** — Black-Box-Tests erreichen unexportierte Pfade seltener; die
  Zahl misst erst `slice-harness-coverage` nach allen vier Umstellungs-Slices.
  — **Ausgang:** — (bei Closure)
- **Befund verdeckt statt behoben** — `_ =` vor einem ungeprüften Fehler, eine Globale
  als Funktion mit demselben geteilten Zustand: `make lint` ist grün, der Test nicht
  besser. Review prüft die Form, das Werkzeug nur die Zahl (Grenze von `SPEC-049`).
  — **Ausgang:** — (bei Closure)

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>

**Mutationen des Implementers** (`AGENTS.md` §3.10, Prüfung geht im Umbau verloren):
je Mutant ein frischer Pfad außerhalb des Repos (`cp -r` ohne `-p`), Tests per
Bind-Mount im Image der Stufe `deps` mit `go test -run`; ohne Mutation ist jeder
genannte Test grün, mit ihr rot.

| Zusage | Mutation | roter Test |
|---|---|---|
| jeder Fehlercode ergibt eine Klasse mit Namen (über die Brücke `Klasse`) | Klasse 6 aus `klassen` gestrichen | `TestCodeTabelle` |
| eine Kette trägt den Code des klassifizierten Fehlers (Literale mit Feldnamen) | `Meldungen` setzt immer `CodeInternal` | `TestFehlerKette` |
| ein nicht klassifizierter gleichrangiger Fehler folgt der ersten Meldung | an die letzte statt an die erste Meldung angehängt | `TestFehlerGleichrangig` |
| `;` macht eine Anfrage zu keiner Lebendprüfung (über `IstLebendpruefung`) | `;` zum Leerraum | `TestLebendpruefungErkennung` |
| `\v` ist nur mit `vt` Leerraum | Bedingung `vt` entfernt | `TestLebendpruefungVertikalerTabulator` |
| `\v` ab Hauptversion 17 (über `VTLeerraum`) | `>=` zu `>` | `TestLebendpruefungServerversion` |
| jedes Feld der Client-Nachricht zählt (über `Abweichung`) | Vergleich von `max_rows` entfernt | `TestReplayExtendedFelder` |
| Parameter-Nummer ab 1 (über `ParameterStelle`) | `i+1` zu `i` | `TestReplayExtendedParameterStelle` |
| letzte Nummer ist die aufgezeichnete (über `ReplayService`) | Zahl der Interaktionen statt Nummer | `TestReplayLetzteNummer` |
| letzte Nummer ist die der letzten Interaktion (über `ReplayService`) | erste statt letzte Interaktion | `TestReplayLetzteNummer` |
| Query nach dem Ende liefert `ErrSessionEnded` (Subtest mit eigenem Kontext) | Prüfung `beendet` nach dem Upstream-Aufruf in `Query` entfernt | `TestRecordEndeNachErfolgreichemUpstream/Query` |
| AwaitServer nach dem Ende liefert `ErrSessionEnded` (Subtest mit eigenem Kontext) | Prüfung `beendet` nach `Receive` in `AwaitServer` entfernt | `TestRecordEndeNachErfolgreichemUpstream/AwaitServer` |
| fortlaufende Nummer der Interaktionen (Testdaten als Funktionen) | Nummer der einfachen Anfrage `+2` statt `+1` | `TestRecordExtendedSyncGruppe` |
| abweichender Nachrichtentyp ist PGR-E5001 (Testdaten als Funktionen) | Vergleich von `type` entfernt | `TestReplayExtendedAbweichung` |
| mit `FailOnUnconsumed` ist Unverbrauchtes ein Fehler (Optionen als Funktion) | Zweig `streng` nie genommen | `TestReplayNichtVerbrauchtMeldung` |

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
