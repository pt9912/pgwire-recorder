# Slice slice-v1-abschluss-schreiben: Atomares Schreiben, --output und --force

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0005](../../adr/0005-recording-store-ist-driven-adapter.md)

**Berührte Spec-Stellen:** `LH-FA-07.a` · `LH-FA-13.a` · `LH-FA-17.a` · `SPEC-014` · `SPEC-016` · `SPEC-023` · `SPEC-034`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record` schreibt die Aufzeichnung atomar (temporäre Datei, dann Verschieben), und lehnt ein vorhandenes `--output` ohne `--force` ab.

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Schreiben* (atomares Schreiben aus
DoD-Punkt 1, `--output` und `--force` aus DoD-Punkt 2; Entscheidung des Nutzers vom
2026-10-08, Schnitt nach F-345; dort §7 `Gegenstand:`). Im Einzelnen, je mit dem
ursprünglichen Geber:

- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): atomares Schreiben,
  `--output` und `--force`; das atomare Schreiben gilt auch für die Aufzeichnung nach
  dessen Zwangsende (DoD-Punkt 1).

**Abgegeben** an `slice-v1-abschluss-konfiguration` (dort §1, *Übernimmt*, mit der Kennung
dieses Slice; Entscheidung des Architect vom 2026-10-08 zu F-496): die Anmeldung von
`--output` und `--force` am allgemeinen Leser, damit ihre Umgebungsvariablen gelesen und
geprüft werden, und die strengen Werte für `--force` (aus `slice-replay-semantik-mismatch`).
Die Schlüssel `output` und `force` der Konfigurationsdatei folgen aus dieser Anmeldung und
kommen mit `slice-v1-abschluss-konfigurationsdatei` (dort §6, Rückgabe 12). Hier bleiben das
atomare Schreiben und das Verhalten bei vorhandenem `--output`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Herunterfahren und die Frist — `slice-v1-abschluss-herunterfahren`; er liegt vor diesem
  Slice, und das atomare Schreiben gilt auch für die Aufzeichnung nach seinem Zwangsende.
- Der allgemeine Leser, die Anmeldung von `--output` und `--force` und ihre Werte —
  `slice-v1-abschluss-konfiguration` (oben, *Abgegeben*).
- Die Konfigurationsdatei und ihre Schlüssel `output` und `force` —
  `slice-v1-abschluss-konfigurationsdatei`; dieser Slice hat keinen Code für die Datei.
- Das Schreiben im SQLite-Format — anderer Vorgang: `slice-v1-abschluss-sqlite-format`
  schreibt seine Datei je Session in einer Transaktion; dieser Slice schreibt die Datei
  des Standardformats YAML.
- Code im Kern, im PGWire-Adapter, im CLI-Adapter und im Bootstrap — Schicht-Abgrenzung:
  Der Slice ändert nur den Recording-Adapter; `--output` und `--force` liefert der allgemeine
  Leser (oben, *Abgegeben*), abgelehnt wird im Recording-Adapter.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings): `record` schreibt die Aufzeichnung in eine temporäre Datei und verschiebt
      sie danach an `--output`; ein Abbruch während des Schreibens, auch nach dem
      Zwangsende des Herunterfahrens, lässt unter `--output` keine teilweise Datei zurück
      (Test).
- [ ] [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Ein vorhandenes `--output` wird ohne `--force` abgelehnt
      (Exit-Code 2, `SPEC-014`), gleich aus welcher Quelle `--force` kommt (Test). Beleg
      in §7 für beide Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `spec/spezifikation.md` (`LH-FA-07.a`) | keine Änderung durch den Implementer | Die Randformen in §6 sind entschieden (Architect, 2026-10-09, vor dem Code); eine Randform, die §6 nicht nennt, geht an den Architect zurück (`AGENTS.md` §3.12) |
| `internal/adapters/driven/recording` | update | temporäre Datei, atomares Verschieben (vorhanden); neu nach §6: ein vorhandener Pfad, der keine reguläre Datei ist, beim Start `PGR-E3001`, ohne und mit `--force` (vorher `PGR-E2002` ohne, ein Fehler erst beim ersten Schreiben mit `--force`); temporäre Datei nach einem Fehlschlag entfernt, ein Fehler beim Entfernen als Ursache derselben Meldung (vorher blieb sie liegen). Die Dateioperationen von `Prepare` (Probedatei anlegen und entfernen) und `Write` (Zufall, Schreiben, Synchronisieren, Verschieben, Entfernen) stehen in einer Tabelle von Funktionen, die die Unit-Tests über `export_test.go` ersetzen. Nach den zwei Rückgaben aus §7 (Architect, 2026-10-09) ändert sich das Verhalten nicht; die Tabelle nimmt das Setzen der Rechte und das Schließen auf, damit ein Test sie scheitern lassen kann. Nach den Übergaben des Reviews an den Architect (F-543, F-544; Architect, 2026-10-09): Das Schließen der Probedatei läuft über die Tabelle (`schliessen`), damit ein Test es scheitern lassen kann; sonst ändert sich das Verhalten nicht, der Code folgt beiden Randformen schon |
| `internal/adapters/driving/cli` | keine Änderung | `--output` und `--force` liefert der allgemeine Leser (§1); abgelehnt wird im Recording-Adapter |
| `internal/adapters/driven/recording` (Unit-Tests `schreiben_test.go`), `test/integration` (`schreiben_e2e_test.go`, Helfer `startProzessIn` in `record_e2e_test.go`) | update | Boundary/Negative nach `LH-FA-07.a`; `--force` aus Kommandozeile, Umgebung und Datei nach `LH-FA-17.a`. Neu nach den zwei Rückgaben aus §7: Fehlschlag beim Setzen der Rechte und beim Schließen (`PGR-E3001`, Zieldatei unverändert, temporäre Datei entfernt) und Rechte des Ziels bei ersetzter Verknüpfung, je mit Mutation (`AGENTS.md` §3.10). Neu nach den Übergaben des Reviews an den Architect: (a) Schließen und Entfernen der Probedatei scheitern je für sich → `PGR-E3001` beim Start, mit dem Text des Fehlers als Ursache; Mutation: den Fehler verwerfen (Mutant P des Reviews), muss rot werden. (b) Ein anderer Fehler als ein belegter Name beim Anlegen der temporären Datei (etwa ein zu langer Name) → sofort `PGR-E3001` mit der Ursache des Betriebssystems, nicht *kein freier Name*; Mutation: jeder Fehler führt zu einem neuen Versuch (Mutant C des Reviews), muss rot werden. Nach F-537 des Reviews: nicht beschreibbares Verzeichnis auch für eine vorhandene Datei und eine Verknüpfung auf sie mit `--force` → `PGR-E3001` beim Start, mit Ursache; Mutation: Probedatei für einen vorhandenen Pfad ausgelassen (Mutant D des Reviews), muss rot werden |
| `docs/user/benutzerhandbuch.md` | update | atomares Schreiben, `--output` und `--force`, übrig gebliebene temporäre Dateien entfernt der Anwender; das Verzeichnis muss beschreibbar sein, auch mit `--force` |
| `docs/user/abdeckung-*.md` | update | über `make abdeckung` aus den Deklarationen der neuen Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-konfiguration` und
`slice-v1-abschluss-konfigurationsdatei` liegen in `done/` (`--force` wird über den allgemeinen
Leser gelesen, auch aus der Datei). Schritt 7 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08, nach den beiden Schnitten vom 2026-10-09 zwei Schritte später). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Atomares Schreiben verlangt
  plattformspezifische Lösungen — zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Das Verschieben ist auf einer
  Zielplattform nicht atomar zu haben, und die Spezifikation sagt nicht, was dann gilt;
  dann zuerst die Entscheidung, gegebenenfalls ein Carveout.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist:

- **Ort und Name der temporären Datei** — im Verzeichnis der Zieldatei,
  `.<Name der Zieldatei>.<16 Hexziffern>.tmp`, exklusiv angelegt, nach zehn belegten Namen
  `PGR-E3001`; entschieden in `LH-FA-07.a` *Temporäre Datei* (Architect, 2026-10-09).
- **Anderer Fehler als ein belegter Name beim Anlegen der temporären Datei** — kein neuer
  Versuch, sofort `PGR-E3001` mit dem Text des Fehlers als Ursache; entschieden in
  `LH-FA-07.a` *Temporäre Datei*, *Name* (Architect, 2026-10-09, Übergabe F-544 des
  Reviews). Der Code folgt (`fs.ErrExist`), der Test ist `TestWriteAndererFehlerBeimAnlegen`
  (§3).
- **Schließen oder Entfernen der Probedatei scheitert** — `PGR-E3001` beim Start, der Text
  des Fehlers als Ursache, eine nicht entfernte Probedatei bleibt liegen; entschieden in
  `LH-FA-07.a` *Zielpfad beim Start*, *Verzeichnis* (Architect, 2026-10-09, Übergabe F-543
  des Reviews). Der Code folgt (`schliessen` aus der Tabelle, `closeErr`, `removeErr`), der
  Test ist `TestPrepareProbedateiScheitert` (§3).
- **Übrig gebliebene temporäre Datei oder Probedatei eines früheren oder gleichzeitigen
  Laufs** — weder gelesen noch überschrieben, entfernt oder gemeldet; entschieden in
  `LH-FA-07.a` *Temporäre Datei* (Architect, 2026-10-09).
- **`--output` ist ein Verzeichnis oder sonst keine reguläre Datei** — `PGR-E3001` beim
  Start, auch mit `--force`; entschieden in `LH-FA-07.a` *Zielpfad beim Start* (Architect,
  2026-10-09).
- **`--output` liegt in einem fehlenden oder nicht beschreibbaren Verzeichnis** — `PGR-E3001`
  beim Start, kein Verzeichnis wird angelegt (bestätigt: `PGR-E3…` nach `SPEC-016`, nicht
  `PGR-E2…`); entschieden in `LH-FA-07.a` *Zielpfad beim Start* (Architect, 2026-10-09).
- **Nicht prüfbarer Pfad** (der Zustand lässt sich nicht feststellen) — `PGR-E3001`;
  entschieden in `LH-FA-07.a` *Zielpfad beim Start* (Architect, 2026-10-09).
- **Symbolische Verknüpfung als `--output`** — vorhanden nach ihrem Ziel, ins Leere nicht
  vorhanden; das Verschieben ersetzt die Verknüpfung, nicht ihr Ziel; entschieden in
  `LH-FA-07.a` *Zielpfad beim Start* (Architect, 2026-10-09).
- **Pfad, der erst nach dem Start entsteht** — vom nächsten Schreibvorgang ohne Prüfung
  ersetzt (Grenze, `LH-FA-08` prüft beim Start); entschieden in `LH-FA-07.a` *Zielpfad beim
  Start* (Architect, 2026-10-09).
- **Rechte der Zieldatei** — neu `0666` nach der umask, ersetzt mit ihren Zugriffsrechten;
  Eigentümer, Gruppe und weitere Attribute nicht erhalten (Grenze); entschieden in
  `LH-FA-07.a` *Temporäre Datei* (Architect, 2026-10-09).
- **Rechte, wenn `--force` eine symbolische Verknüpfung auf eine Datei ersetzt** — die
  Zieldatei erhält die Zugriffsrechte des Ziels der Verknüpfung (wie *Vorhanden* nach dem
  Ziel zählt); entschieden in `LH-FA-07.a` *Temporäre Datei*, *Rechte* (Architect,
  2026-10-09, Rückgabe des Implementers aus §7). Der Code folgt (`os.Stat`), der Test ist
  `TestWriteRechteVerknuepfung` (§3).
- **Zieldatei beim Schreiben nicht prüfbar** (`os.Stat` vor dem Setzen der Rechte scheitert
  anders als mit „nicht vorhanden“) — akzeptiertes Negativ, keine Zusage: Der Pfad ist beim
  Start geprüft (*Zielpfad beim Start*), und ein Pfad, der sich danach ändert, liegt in der
  dort genannten Grenze; der Code behandelt ihn wie eine neue Datei (Architect, 2026-10-09).
- **Vorhandenes `--output` ohne `--force`** — Exit-Code 2 (`PGR-E2002`); entschieden in
  `SPEC-014` und `LH-FA-07.a`.
- **Werte von `--force`** — nur `true` oder `false`; entschieden in `LH-FA-17.a`, geliefert von
  `slice-v1-abschluss-konfiguration` (§1, *Abgegeben*).
- **Fehlschlag von Anlegen, Setzen der Rechte, Schreiben, Synchronisieren, Schließen oder
  Verschieben** — `PGR-E3001`,
  Zieldatei unverändert, temporäre Datei entfernt; misslingt das Entfernen, bleibt sie, und
  der Fehler folgt als Ursache derselben Meldung; entschieden in `LH-FA-07.a` *Temporäre
  Datei* (Architect, 2026-10-09; Setzen der Rechte und Schließen ergänzt am selben Tag auf
  Rückgabe des Implementers aus §7). Exit-Code und Fortgang danach: `LH-FA-13.b` (bestehend).
- **Abgrenzung zu `sqlite`** — *Zielpfad beim Start* gilt über `LH-FA-22.a` auch für
  `sqlite`; *Temporäre Datei* beschreibt `yaml`. Was davon für `sqlite` gilt, entscheidet
  `slice-v1-abschluss-sqlite-format` vor seinem Code (akzeptiertes Negativ hier: Der Adapter
  für `sqlite` existiert noch nicht, eine Entscheidung ohne seinen Plan wäre geraten). Der
  Nehmer führt die Sendung in seinem §1 mit dieser Kennung (`AGENTS.md` §3.13; Architect,
  2026-10-09, Übergabe F-538 des Reviews).

**Risiken:**

- Atomarität des Verschiebens ist plattformabhängig („bestmöglich atomar“; übernommen aus
  `slice-v1-abschluss-betrieb`) — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten.
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

### Belege des Implementers

**Stand.** Geliefert sind DoD-Punkt 1 und 2 in drei Commits: `aad7085` (Recording-Adapter:
keine reguläre Datei `PGR-E3001` beim Start, ohne und mit `--force`; temporäre Datei nach
einem Fehlschlag entfernt, Fehler beim Entfernen als Ursache derselben Meldung; Unit- und
Integrationstests; Abdeckung), `df4f591` (Handbuch, keine Beispiele geändert), `65968e6`
(Probedatei über dieselbe Tabelle der Dateioperationen; Tests zu ihrem Namen und zu einem
nicht beschreibbaren Verzeichnis). Neu sind
`internal/adapters/driven/recording/schreiben_test.go`,
`internal/adapters/driven/recording/export_test.go` (`PruefeMit`, `SchreibeMit`, `Eingriffe`)
und `test/integration/schreiben_e2e_test.go`; geändert `yaml.go` (`pruefe`, `schreibe`,
`dateiOps`, `betriebssystem`, `neueTempDatei` mit Zufall als Parameter) und
`record_e2e_test.go` (Helfer `startProzessIn` mit Arbeitsverzeichnis). CLI-Adapter, Kern,
PGWire-Adapter und Bootstrap sind unverändert (§1). Nach den zwei Entscheidungen des
Architect (`b74e16e`) nimmt `24fc87c` das Setzen der Rechte und das Schließen in `dateiOps`
und `Eingriffe` auf (`Rechte`, `Schliessen`), ergänzt `TestWriteFehlschlag` um beide Fälle,
fügt `TestWriteRechteVerknuepfung` hinzu und zieht Handbuch und Abdeckung nach; das
Verhalten des Codes ist unverändert.

**Randformen.** Kein Code-Commit ändert §6. Entschieden ist nur, was §6 nennt. Zwei Punkte
gingen als Frage an den Architect zurück; beide sind entschieden (`b74e16e`, `LH-FA-07.a`),
das Verhalten des Codes blieb, und je ein Test sagt sie zu (DoD-Punkt 1 unten):

| Randform | Frage | Entschieden |
|---|---|---|
| Fehlschlag beim Setzen der Rechte (`Chmod`) oder beim Schließen der temporären Datei | `LH-FA-07.a` *Fehlschlag* nennt Anlegen, Schreiben, Synchronisieren und Verschieben. Der Code behandelt beide wie einen Fehlschlag des Schreibens (`PGR-E3001`, Zieldatei unverändert, temporäre Datei entfernt). Gehören sie zum Schreiben? | ja, beide gehören zum *Fehlschlag*; Test `TestWriteFehlschlag` (Fälle `Rechte`, `Schliessen`) |
| Rechte, wenn `--force` eine symbolische Verknüpfung auf eine Datei ersetzt | Der Code übernimmt die Zugriffsrechte des Ziels (`os.Stat` folgt der Verknüpfung), wie vor dem Slice. *Rechte* sagt „eine ersetzte behält ihre Zugriffsrechte“, *Vorhanden* zählt nach dem Ziel; ist das Ziel gemeint? | ja, die Zieldatei erhält die Zugriffsrechte des Ziels; Test `TestWriteRechteVerknuepfung` |

Ein Fehler von `os.Stat` beim Schreiben ist ein akzeptiertes Negativ (§6); dazu keine Zusage
und kein Test.

Kein anderer Slice ist als neue Adresse genannt (§3.13).

**Größe.** Der Diff `ad1b7ba..65968e6` umfasst 753 hinzugefügte und 29 entfernte Zeilen: Code
71 (`yaml.go`), Unit-Tests 465, Integration 167, Handbuch 23, Abdeckung 24, Plan 3. Die
Nacharbeit `24fc87c` fügt Code 15 (`yaml.go`), Unit-Tests 46 (`schreiben_test.go`,
`export_test.go`), Handbuch 4 und Abdeckung 3 Zeilen (geändert und hinzugefügt) hinzu. Ob er
in eine Review-Sitzung passt, urteilt das Review.

**Läufe.**

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make test` | vor `aad7085`, vor `65968e6` | Exit 0 |
| `make lint` | vor `aad7085`, vor `65968e6` | Exit 0, `0 issues.` |
| `make test-integration` | vor `aad7085` | Exit 0, `run-integration-tests: gruen`, die vier neuen `TestE2ERecord…` PASS |
| `make docs-check` | vor `aad7085`, vor `65968e6` | 0 Befunde |
| `make kopf-check`, `make abdeckung-check` | vor `aad7085` | Exit 0 |
| `make abdeckung` | vor `aad7085`, vor `65968e6` | Tabellen nachgezogen |
| `make test`, `make lint`, `make docs-check`, `make kopf-check`, `make abdeckung-check` | vor `24fc87c` | Exit 0, 0 Befunde |
| `make abdeckung` | vor `24fc87c` | `abdeckung-unit.md` nachgezogen |
| `make gates` | Commit dieser Belege | Exit 0, vor der Übergabe |

**Weg der Mutanten.** Je Mutant eine frische Kopie des Arbeitsbaums unter dem Scratch-Pfad der
Sitzung (`cp -r` ohne `-p`, neue mtime); ein Skript ersetzt genau ein Vorkommen (oder das
genannte n-te) und bricht bei einer anderen Trefferzahl ab. Unit-Mutanten: `docker build
--target test` in der Kopie (gofmt, vet, `go test ./...`). Integrations-Mutanten: Stufe
`integration` der Kopie mit eigenem Tag, eigenem internem Netz und eigenem
PostgreSQL-Container, gefahren nur der genannte Test (`-test.run '^Name$'`); Container, Netz
und Image entfernt danach ein `trap`. Die Kopie wird danach gelöscht; der Arbeitsbaum blieb
unberührt. Gefahren: 31 Unit-Mutanten (die Mutanten an `Prepare` nach `65968e6` neu) und 11
Integrations-Mutanten, alle rot. Ein Mutant (Form des Namens der Probedatei) war vor
`65968e6` grün, weil die Probedatei sofort entfernt ist und ihr Name von außen nicht zu sehen
war; mit `TestPrepareProbedatei` ist er rot. Nach `24fc87c` gefahren: eine unveränderte Kopie
(grün) und fünf Unit-Mutanten, alle rot im genannten Test. Drei davon (Schließen übergangen,
Rechte und Schließen ohne Entfernen) waren im ersten Lauf nur an `gofmt` rot, also aus dem
falschen Grund; formatgerecht wiederholt, sind sie im Test rot.

**DoD-Punkt 1 — atomares Schreiben, Zielpfad beim Start (`LH-FA-07.a`).**

| Zusage | Mutation | rote Tests |
|---|---|---|
| keine reguläre Datei ist `PGR-E3001` | Prüfung `IsRegular` entfernt | `TestPrepareKeineRegulaereDatei` (alle sechs Fälle); E2E: `TestE2ERecordZielIstVerzeichnis` (Exit 2 statt 3) |
| auch mit `--force` | Prüfung nur ohne `replace` | `TestPrepareKeineRegulaereDatei` (drei Fälle mit `replace`) |
| auch eine andere Art als ein Verzeichnis | `IsDir` statt `IsRegular` | `TestPrepareKeineRegulaereDatei` (benannte Pipe) |
| vorhanden nach dem Ziel einer Verknüpfung, ins Leere nicht vorhanden | `os.Lstat` statt `os.Stat` | `TestPrepareVerknuepfung` (beide Fälle), `TestPrepareNichtPruefbar` |
| nicht prüfbar ist `PGR-E3001` | jeder Fehler von `Stat` als „nicht vorhanden“ | `TestPrepareNichtPruefbar` (Verknüpfungsschleife: kein Fehler; Pfad unter einer Datei: andere Meldung) |
| fehlendes Verzeichnis `PGR-E3001`, kein Verzeichnis angelegt | `os.MkdirAll` vor der Probedatei | `TestPrepareVerzeichnis`, `TestPrepareVorhandeneDatei` |
| nicht beschreibbares Verzeichnis `PGR-E3001` | Fehler mit `fs.ErrPermission` beim Anlegen der Probedatei übergangen | `TestPrepareVerzeichnisNichtBeschreibbar` |
| Probedatei sofort entfernt | Entfernen weggelassen | `TestPrepareProbedatei`, `TestPrepareVerzeichnis`, `TestRoundtrip`, `TestUebrigGebliebeneDateien` |
| Probedatei `.<Name>.<Zufallsteil>.probe` | Name ohne Punkt, Endung `.pruef` | `TestPrepareProbedatei` |
| Probedatei im Verzeichnis der Zieldatei | Probedatei in `os.TempDir()` | `TestPrepareProbedatei`, `TestPrepareVerzeichnis`, `TestPrepareVorhandeneDatei` |
| Verschieben ersetzt die Verknüpfung, nicht ihr Ziel | Pfad vor dem Schreiben über `filepath.EvalSymlinks` aufgelöst | `TestWriteVerknuepfungErsetzt` |
| temporäre Datei im Verzeichnis der Zieldatei | temporäre Datei in `os.TempDir()` | `TestWriteNameDerTempDatei`, `TestWriteBelegterName`, `TestWriteZehnBelegteNamen`, `TestWriteEntfernenScheitert` |
| Name mit 16 Hexziffern | 7 statt 8 Zufallsbytes | `TestWriteNameDerTempDatei`, `TestWriteBelegterName`, `TestWriteZehnBelegteNamen` |
| exklusiv angelegt, belegter Name nicht überschrieben | `O_TRUNC` statt `O_EXCL` | `TestWriteBelegterName`, `TestWriteZehnBelegteNamen` |
| nach zehn belegten Namen `PGR-E3001` | elf Versuche | `TestWriteZehnBelegteNamen` |
| nicht vor dem zehnten Versuch aufgegeben | neun Versuche | `TestWriteZehnBelegteNamen` |
| Fehlschlag Anlegen `PGR-E3001` | Code `PGR-E1000` am Anlegen | `TestWriteFehlschlag` (Anlegen), `TestWriteZehnBelegteNamen`, dazu `TestRunRecordSchreibfehlerJeStufe` und `TestRunRecordSchreibfehlerVorGemerkterKlasse` |
| Fehlschlag Schreiben (nach einem Teil der Daten): `PGR-E3001`, Zieldatei unverändert | Fehler des Schreibens übergangen | `TestWriteFehlschlag` (Schreiben: kein Fehler, Zieldatei mit halber Aufzeichnung), `TestWriteEntfernenScheitert` |
| Fehlschlag Synchronisieren: `PGR-E3001`, Zieldatei unverändert | Fehler des Synchronisierens übergangen | `TestWriteFehlschlag` (Synchronisieren) |
| Fehlschlag Setzen der Rechte: `PGR-E3001`, Zieldatei unverändert | Fehler von `rechte` verworfen (`_ =` statt `merr =`) | `TestWriteFehlschlag` (Rechte: kein Fehler, Zieldatei ersetzt) |
| Fehlschlag Schließen: `PGR-E3001`, Zieldatei unverändert | Fehler von `schliessen` verworfen | `TestWriteFehlschlag` (Schliessen: kein Fehler, Zieldatei ersetzt) |
| temporäre Datei nach Fehlschlag des Setzens der Rechte entfernt | nach dem Fehler von `rechte` sofort `PGR-E3001` ohne Entfernen | `TestWriteFehlschlag` (Rechte: `.tmp` liegt im Verzeichnis) |
| temporäre Datei nach Fehlschlag des Schließens entfernt | nach dem Fehler von `schliessen` sofort `PGR-E3001` ohne Entfernen | `TestWriteFehlschlag` (Schliessen: `.tmp` liegt im Verzeichnis) |
| Fehlschlag Verschieben: `PGR-E3001` | Fehler des Verschiebens übergangen | `TestWriteFehlschlag` (Verschieben), `TestWriteEntfernenScheitert` |
| temporäre Datei nach Fehlschlag des Schreibens oder Synchronisierens entfernt | Entfernen im ersten Zweig weggelassen | `TestWriteFehlschlag` (Schreiben, Synchronisieren), `TestWriteEntfernenScheitert` |
| temporäre Datei nach Fehlschlag des Verschiebens entfernt | Entfernen im zweiten Zweig weggelassen | `TestWriteFehlschlag` (Verschieben), `TestWriteEntfernenScheitert`; E2E: `TestE2ERecordSchreibfehlerNachZwangsende` (zwei `.tmp` im Verzeichnis) |
| Fehler beim Entfernen folgt als Ursache (Schreiben) | Fehler des Entfernens verworfen, erster Zweig | `TestWriteEntfernenScheitert` (Schreiben) |
| Fehler beim Entfernen folgt als Ursache (Verschieben) | Fehler des Entfernens verworfen, zweiter Zweig | `TestWriteEntfernenScheitert` (Verschieben) |
| in derselben Meldung | Fehler des Entfernens als eigene Meldung `PGR-E3001` | `TestWriteEntfernenScheitert` (zwei Meldungen) |
| übrig gebliebene Probedatei nicht entfernt | `Prepare` entfernt `.<Name>.*.probe` | `TestUebrigGebliebeneDateien` |
| übrig gebliebene temporäre Datei nicht entfernt | `Write` entfernt `.<Name>.*.tmp` | `TestUebrigGebliebeneDateien`, `TestWriteBelegterName`, `TestWriteZehnBelegteNamen` |
| übrig gebliebene Dateien nicht gemeldet | `Prepare` meldet `PGR-E3001`, wenn `.<Name>.*` liegt | `TestUebrigGebliebeneDateien` |
| Pfad, der erst nach dem Start entsteht, ohne Prüfung ersetzt | `Write` lehnt einen vorhandenen Pfad ab | `TestWritePfadNachDemStart` und 4 weitere Unit-Tests |
| neue Datei `0666` nach der umask | `0600` statt `0666` | `TestWriteRechteUnterUmask022` |
| ersetzte Datei behält ihre Zugriffsrechte | `Chmod` weggelassen | `TestWriteRechteUndAtomar` |
| ersetzte Verknüpfung auf eine Datei: Zieldatei mit den Zugriffsrechten des Ziels | `os.Lstat` statt `os.Stat` beim Schreiben | `TestWriteRechteVerknuepfung` (`0777` statt `0640`) |
| nach dem Zwangsende vollständig ersetzt, keine temporäre Datei | Verschieben weggelassen | E2E: `TestE2ERecordZwangsendeSchreibtVollstaendig` |

**DoD-Punkt 2 — vorhandenes `--output` und `--force` aus jeder Quelle (`LH-FA-08`,
`LH-FA-17`).** Je Quelle und Wert ein eigener Fall in `TestE2ERecordVorhandeneZieldatei`; je
Mutant ist nur der genannte Fall rot.

| Zusage | Mutation | roter Fall |
|---|---|---|
| vorhandene Datei ohne `--force` `PGR-E2002`, mit `--force` ersetzt | `Prepare` lehnt jede vorhandene reguläre Datei ab | `Kommandozeile/true`, `Umgebung/true`, `Datei/true` |
| `true` aus der Kommandozeile ersetzt | Kommandozeile für `force` übergangen | `Kommandozeile/true` |
| `true` aus der Umgebung ersetzt | Umgebung für `force` übergangen | `Umgebung/true` |
| `true` aus der Datei ersetzt | Datei für `force` übergangen | `Datei/true` |
| `false` aus der Kommandozeile lehnt ab | Wert der Kommandozeile für `force` als `true` | `Kommandozeile/false` |
| `false` aus der Umgebung lehnt ab | Wert der Umgebung für `force` als `true` | `Umgebung/false` |
| `false` aus der Datei lehnt ab | Wert der Datei für `force` als `true` | `Datei/false` |
| ohne Quelle lehnt ab | Standardwert `true` | `keine/false` |

Die Unit-Zeile `PGR-E2002` ohne `--force` hält `TestPrepareVorhandeneDatei` (vor dem Slice),
die Zeile mit Verknüpfung `TestPrepareVerknuepfung` (oben, `os.Lstat`).

**Grenzen.**

1. *Liest nicht* (übrig gebliebene Datei): Dass `record` eine übrig gebliebene Datei nicht
   liest, ist über die Schnittstelle nicht zu sehen; kein Mutant.
2. *Bestmöglich atomar*: Die Tests laufen nur unter Linux im selben Dateisystem; ob das
   Verschieben auf einer anderen Plattform atomar ist, prüft kein Test (Risiko in §6).
3. *Harter Abbruch* (`SIGKILL`) während des Schreibens ist nicht gefahren; die
   Spezifikation erlaubt dann eine liegen gebliebene temporäre Datei, die Zieldatei hält das
   Verschieben in einem Schritt (`TestWriteFehlschlag`, Schreiben nach einem Teil der Daten).
4. Im Integrationstest scheitert das Schreiben nach dem Zwangsende nur über ein Verzeichnis an
   `--output` (Verschieben); ein Fehlschlag von Anlegen, Schreiben oder Synchronisieren im
   gebauten Binary ist ohne Eingriff nicht herzustellen und nur in den Unit-Tests über
   `Eingriffe` gefahren.

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `3864e44` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/spec-randform-erst-im-review-entschieden` (12×, verkörpert in `AGENTS.md`
  §3.12) — darum nennt §6 die Randformen, und der Architect prüft sie vor dem Code.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (14×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (19×, §3.11) — je Zusage eine
  Mutation, Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (17×, §3.9) — dieser Plan entsteht aus einem
  Schnitt; jede Korrektur zieht §1, §3 und §6 im selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — jede Übernahme oben steht
  mit der Kennung des Gebers; die Geber zeigen im selben Commit hierher.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
