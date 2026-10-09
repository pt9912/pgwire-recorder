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

- [x] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings): `record` schreibt die Aufzeichnung in eine temporäre Datei und verschiebt
      sie danach an `--output`; ein Abbruch während des Schreibens, auch nach dem
      Zwangsende des Herunterfahrens, lässt unter `--output` keine teilweise Datei zurück
      (Test).
- [x] [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Ein vorhandenes `--output` wird ohne `--force` abgelehnt
      (Exit-Code 2, `SPEC-014`), gleich aus welcher Quelle `--force` kommt (Test). Beleg
      in §7 für beide Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
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
| `spec/spezifikation.md` (`LH-FA-07.a`) | keine Änderung durch den Implementer | Die Randformen in §6 sind entschieden (Architect, 2026-10-09, vor dem Code); eine Randform, die §6 nicht nennt, geht an den Architect zurück (`AGENTS.md` §3.12) |
| `internal/adapters/driven/recording` | update | temporäre Datei, atomares Verschieben (vorhanden); neu nach §6: ein vorhandener Pfad, der keine reguläre Datei ist, beim Start `PGR-E3001`, ohne und mit `--force` (vorher `PGR-E2002` ohne, ein Fehler erst beim ersten Schreiben mit `--force`); temporäre Datei nach einem Fehlschlag entfernt, ein Fehler beim Entfernen als Ursache derselben Meldung (vorher blieb sie liegen). Die Dateioperationen von `Prepare` (Probedatei anlegen und entfernen) und `Write` (Zufall, Schreiben, Synchronisieren, Verschieben, Entfernen) stehen in einer Tabelle von Funktionen, die die Unit-Tests über `export_test.go` ersetzen. Nach den zwei Rückgaben aus §7 (Architect, 2026-10-09) ändert sich das Verhalten nicht; die Tabelle nimmt das Setzen der Rechte und das Schließen auf, damit ein Test sie scheitern lassen kann. Nach den Übergaben des Reviews an den Architect (F-543, F-544; Architect, 2026-10-09): Das Schließen der Probedatei läuft über die Tabelle (`schliessen`), damit ein Test es scheitern lassen kann; sonst ändert sich das Verhalten nicht, der Code folgt beiden Randformen schon. Nach der Randform-Rückgabe zur Probedatei (Architect, 2026-10-09, §6): Scheitert das Schließen, hängt der Fehler des Entfernens als Ursache an **dieselbe** Meldung (wie `schreibe`), statt als zweite gleichrangige; sonst entfernt der Code schon |
| `internal/adapters/driving/cli` | keine Änderung | `--output` und `--force` liefert der allgemeine Leser (§1); abgelehnt wird im Recording-Adapter |
| `internal/adapters/driven/recording` (Unit-Tests `schreiben_test.go`), `test/integration` (`schreiben_e2e_test.go`, Helfer `startProzessIn` in `record_e2e_test.go`) | update | Boundary/Negative nach `LH-FA-07.a`; `--force` aus Kommandozeile, Umgebung und Datei nach `LH-FA-17.a`. Neu nach den zwei Rückgaben aus §7: Fehlschlag beim Setzen der Rechte und beim Schließen (`PGR-E3001`, Zieldatei unverändert, temporäre Datei entfernt) und Rechte des Ziels bei ersetzter Verknüpfung, je mit Mutation (`AGENTS.md` §3.10). Neu nach den Übergaben des Reviews an den Architect: (a) Schließen und Entfernen der Probedatei scheitern je für sich → `PGR-E3001` beim Start, mit dem Text des Fehlers als Ursache; Mutation: den Fehler verwerfen (Mutant P des Reviews), muss rot werden. (b) Ein anderer Fehler als ein belegter Name beim Anlegen der temporären Datei (etwa ein zu langer Name) → sofort `PGR-E3001` mit der Ursache des Betriebssystems, nicht *kein freier Name*; Mutation: jeder Fehler führt zu einem neuen Versuch (Mutant C des Reviews), muss rot werden. Nach F-537 des Reviews: nicht beschreibbares Verzeichnis auch für eine vorhandene Datei und eine Verknüpfung auf sie mit `--force` → `PGR-E3001` beim Start, mit Ursache; Mutation: Probedatei für einen vorhandenen Pfad ausgelassen (Mutant D des Reviews), muss rot werden. Nach der Randform-Rückgabe zur Probedatei (Architect, 2026-10-09, §6): (c) Schließen scheitert, Entfernen gelingt → eine Meldung `PGR-E3001` mit der Ursache des Schließens, keine `.probe` im Verzeichnis; Mutation: nach gescheitertem Schließen nicht entfernen, muss rot werden. (d) Schließen und Entfernen scheitern → genau **eine** Meldung `PGR-E3001`, die beide Texte trägt, die Probedatei bleibt; Mutationen: zwei gleichrangige Meldungen (heutiger Code), Ursache des Entfernens verworfen — je rot. Grenze 6 in §7 wird damit zum Beleg |
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
  des Reviews). Der Code folgt (`pruefe`: `schliessen` und `entfernen` aus der Tabelle, der
  Fehler als Ursache der Meldung), der Test ist `TestPrepareProbedateiScheitert` (§3).
- **Schließen der Probedatei scheitert, Entfernen danach** — `record` entfernt sie dennoch;
  scheitert auch das Entfernen, bleibt sie liegen, und der Fehler des Entfernens folgt als
  Ursache **derselben** Meldung `PGR-E3001` (wie bei der temporären Datei, *Fehlschlag*);
  entschieden in `LH-FA-07.a` *Zielpfad beim Start*, *Verzeichnis* (Architect, 2026-10-09,
  Randform-Rückgabe des Implementers, §7 Grenze 6). Der Code folgt (`pruefe`: der Fehler
  des Entfernens als Ursache der Meldung des Schließens), der Test ist
  `TestPrepareProbedateiScheitert` (Fälle `Schliessen`, `Schliessen und Entfernen`; §3).
- **Weitere Randformen der Probedatei** — akzeptierte Negative, keine eigene Zusage
  (Architect, 2026-10-09): *Belegter Name* — eine Probedatei überschreibt nichts
  (*Übrig gebliebene Datei*); wie oft ein neuer Zufallsteil gezogen wird, sagt
  `LH-FA-07.a` nicht zu (der Code übernimmt die Wiederholung von `os.CreateTemp`), und
  scheitert das Anlegen am Ende, ist es *Anlegen gescheitert*. *Rechte der Probedatei* —
  ohne Belang, sie wird sofort entfernt und nie gelesen. *Abbruch zwischen Anlegen und
  Entfernen* — die liegen gebliebene Probedatei ist eine übrig gebliebene Datei eines
  früheren Laufs (*Übrig gebliebene Datei*).
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
  `slice-v1-abschluss-betrieb`) — **Ausgang:** weiter offen, ins Beobachtungs-Register als
  `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×). Nicht eingetreten: Unter Linux im selben
  Dateisystem ersetzt das Verschieben in einem Schritt (§7 *Grenzen* 2, Verifikation zu
  Punkt 1), die Rückführungen aus §4 traten nicht ein. Nicht entfallen: Geprüft ist keine
  andere Plattform, und was *bestmöglich* unter macOS und Windows heißt, sagt `LH-FA-07.a`
  nicht (§7, Steering-Loop-Eintrag).

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
- **Was hat funktioniert:** Die Prüfung des Architect vor dem Code (`ad1b7ba`) entschied die
  Randformen aus §6 in `LH-FA-07.a`, bevor der Code entstand; Review und Verifikation fanden
  am Binary keine Abweichung von ihnen. Die Belege in §7 stehen je Zusage als Zusage ·
  Mutation · roter Test: Der Implementer fuhr 31 und 11 Mutanten vor der ersten Übergabe,
  danach 5, 14 und 5, alle formatgerecht rot; drei, die im ersten Lauf nur an `gofmt` rot
  waren, erkannte er selbst als aus dem falschen Grund rot und wiederholte sie. Die
  Verifikation fuhr 11 Unit- und 4 E2E-Mutanten nach, alle aus dem genannten Grund rot. Die
  dritte Randform-Rückgabe (`9086fd7`) kam ohne neuen Code: Der Implementer hielt an, sagte
  nichts zu und testete nichts, bis der Architect entschied (`9996a8c`). Der Diff passte in
  eine Review-Sitzung (F-546).
- **Was ging anders als geplant:**
  1. Drei Runden Randform-Rückgabe statt keiner: Setzen der Rechte und Schließen der
     temporären Datei, im Code entschieden und danach zurückgegeben (`0311f72` → `b74e16e`,
     F-539); Schließen und Entfernen der Probedatei und ein anderer Fehler als ein belegter
     Name, erst im Review gefunden (F-543, F-544 → `278d929`); die Kombination *Schließen
     scheitert, dann Entfernen*, von der zweiten Entscheidung offen gelassen (`9086fd7` →
     `9996a8c`). Gemeinsam ist den drei Runden, dass §6 die Fehlschläge nach dem Wortlaut von
     `LH-FA-07.a` aufzählte, nicht nach den Dateioperationen des Codes.
  2. Drei grüne Mutanten fand erst das Review (D, C, P zu F-537, F-544, F-543), trotz 47
     Mutanten des Implementers vor dem Review; jede Zusage hielt nur für einen Teil ihrer Fälle oder stand
     nur im Kommentar.
  3. Der Plan folgte dem Code zweimal nicht (F-540, V-127), und die Abgrenzung zu `sqlite`
     nannte einen Nehmer, der die Sendung nicht führte (F-538) — das erste Auftreten nach der
     Verkörperung in `AGENTS.md` §3.13.
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-09-review-slice-v1-abschluss-schreiben.md`: „0 HIGH · 4 MEDIUM · 3
    LOW · 3 INFO. F-537: Ein nicht beschreibbares Verzeichnis bei vorhandener Datei mit
    `--force` hält keine Mutation (Mutant D grün, Sonde am Binary). F-538: §6 weist die
    Entscheidung für `sqlite` einem Slice zu, der sie nicht führt; §7 sagt das Gegenteil.
    F-539: Entfernen nach Fehlschlag von Setzen der Rechte und Schließen ist im Code
    entschieden und erst danach zurückgegeben. F-540: §1 nennt den CLI-Adapter als geändert,
    §3 und der Diff nicht. F-541: §7 *„Kein Code-Commit ändert §6“* trifft für `24fc87c` nicht
    zu. F-542: Rang-Zeiger `SPEC-033` an der Rechte-Regel. F-543: Der Kommentar an `Prepare`
    sagt `PGR-E3001` für das Entfernen der Probedatei zu; das steht weder in der Spezifikation
    noch in einem Test. F-544: Nur ein belegter Name führt zu einem neuen Zug; das hält keine
    Mutation, und keine Zusage verlangt es. F-545: Die Zusammenfassung *Fehlermodi* ist enger
    als *Fehlschlag*. F-546: Der Diff war in einer Sitzung prüfbar. Wiederkehrende Klassen:
    `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-537),
    `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (F-538),
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (F-539),
    `BEO-REPO/plan-folgt-korrektur-nicht` (F-540),
    `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-541, F-543).“ Verifikation
    `docs/reviews/2026-10-09-verifikation-slice-v1-abschluss-schreiben.md`: „0 HIGH · 0
    MEDIUM · 1 LOW · 1 INFO (V-127: §6 und §7 nennen `closeErr` und `removeErr`, die `3a7018b`
    entfernt hat; V-128: das nicht beschreibbare Verzeichnis ist am Binary bestätigt, ein
    E2E-Test fehlt weiter als benannte Grenze).“ Ausgänge: F-537, F-542, F-543 und F-544 in
    `ca5124f` (F-543 und F-544 nach `278d929`), F-540 in `ca5124f`, F-538 in `278d929`, F-541
    in `9086fd7`, F-545 in `278d929`; F-539 ins Register; F-546 von der Verifikation
    eingeordnet (Risiko unten); V-127 in `b15815f`; V-128 bleibt Grenze 5, kein Auftrag.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke: `LH-FA-07.a` *Schritte je Schreibvorgang*
  sagt „atomar beziehungsweise bestmöglich atomar“ zu, bestimmt aber nicht, was *bestmöglich*
  auf einer Plattform heißt, die das Ersetzen einer Datei nicht atomar bietet, und wie es
  geprüft wird. [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität) sieht
  ausführbare Dateien auch für macOS und Windows vor; Tests, Gate und Verifikation laufen nur
  unter Linux in einem Dateisystem (§7 *Grenzen* 2). Die Rückführung `in-progress` → `open`
  aus §4 nannte genau diese Lücke als Bedingung; sie trat nicht ein, weil keine andere
  Plattform gefahren wurde, nicht weil die Lücke geschlossen ist. Ob die Spezifikation die
  Zusage je Plattform schärft oder als Grenze nennt, entscheidet der Architect mit dem
  Nutzer; gezählt ist sie unter `BEO-REPO/verhalten-nur-unter-linux-geprueft`. Kein Feld
  `liegt in`: Mit diesem Slice ist nichts verkörpert.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist wieder aufgetreten
  (F-540, V-127); die Regel bleibt. §3.10 (seit welle-extended-query) ist wieder aufgetreten
  (F-537, F-544); die Regel bleibt, der Sensor ist mit `slice-harness-mutation` geplant. §3.11
  (seit welle-extended-query) ist wieder aufgetreten (F-541, F-543); die Regel bleibt. §3.12
  (seit slice-harness-randformen-vor-code) ist in beiden Hälften wieder aufgetreten: ein
  Code-Commit entschied, was danach zurückging (F-539), und das Review fand Randformen
  außerhalb von §6 (F-543, F-544); die Regel bleibt. §3.13 (seit
  slice-lint-bestand-kern-driven) ist zum ersten Mal wieder aufgetreten (F-538, gesetzt vom
  Architect in `ad1b7ba`); die Regel bleibt. Die Regel *Beleg im Plan, nicht im Bericht*
  (`.claude/agents/implementer.md`, seit slice-harness-blackbox-kern) ist nicht wieder
  aufgetreten: §7 trägt jeden Gate-Lauf und jede Mutanten-Reihe. `implement-slice` Schritt 19
  (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`): nicht wieder aufgetreten, jeder Mutant
  in einer frischen Kopie.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `b15815f` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`: **Beleg**, 5× → 6× (F-539),
    Stand verkörpert, bleibt.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden`: **Beleg**, 17× → 18× (F-543, F-544;
    die drei Runden der Rückgabe stehen im Beleg), Stand verkörpert, bleibt.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 20× → 21× (F-540, V-127), Stand
    verkörpert, bleibt.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 19× → 20× (F-537, F-544),
    Stand verkörpert, bleibt.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 25× → 26× (F-541, F-543),
    Stand verkörpert, bleibt.
  - `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`: **Beleg**, 3× → 4× (F-538), Stand
    verkörpert, bleibt.
  - `BEO-REPO/verhalten-nur-unter-linux-geprueft`: **neu**, 1× (Risiko aus §6, Ausgang
    *weiter offen*), Stand offen.
  - Ohne Beleg: `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (bleibt 1×, offen;
    anderer Gegenstand: Serverversion, nicht Plattform),
    `BEO-REPO/slice-waechst-durch-uebernahmen` (bleibt 2×, offen; der Slice nahm nach der
    Anlage nichts auf), `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (bleibt 1×, offen;
    die grünen Mutanten D, C und P waren nicht als äquivalent eingestuft, sondern nicht
    gefahren), `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (bleibt 1×, offen),
    `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (bleibt 4×, verkörpert).
  - Einmalig und nicht eingetragen: F-542 (Rang-Zeiger), F-545 (Zusammenfassung enger als die
    Regel), F-546 und V-128 (Hinweise). Mit diesem Slice erreicht kein Eintrag die Schwelle 3×
    neu; über der Schwelle stehen nur Einträge mit Ausgang (`commit-nennt-struktur-kennung`
    3× und `roter-lauf-haengt-bis-zum-zeitlimit` 3× geplant;
    `folge-slice-adresse-nimmt-nicht-an` 4×, `implementer-bericht-erreicht-pruefer-nicht` 4×,
    `randform-im-code-entschieden-dann-zurueckgegeben` 6×,
    `spec-randform-erst-im-review-entschieden` 18×, `negativtests-fehlen-bei-neuem-vertrag`
    20×, `plan-folgt-korrektur-nicht` 21×, `zusage-im-kommentar-weiter-als-pruefung` 26×
    verkörpert; `white-box-liste-vor-code-nur-namenssuche` 3× gestrichen).
- **Folge-Slices:** keiner neu. `slice-v1-abschluss-sqlite-format` (`next/`) ist Adresse für die
  *Abgrenzung zu `sqlite`* aus §6 und führt sie seit `278d929` in §1 unter *Übernommen von
  `slice-v1-abschluss-schreiben`*. Der nächste Schritt nach §5 von
  [welle-v1-abschluss](../welle-v1-abschluss.md) ist `slice-v1-abschluss-einspielen`.
- **Risiken aus §6:** eines. *Atomarität plattformabhängig*: **weiter offen**, ins Register
  als `BEO-REPO/verhalten-nur-unter-linux-geprueft`; Begründung in §6. Die Randformen in §6
  sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt kein Feld `liegt in`; nichts zu
  prüfen. Folge-Slice: `slice-v1-abschluss-sqlite-format` liegt in `next/`; `grep -n
  "slice-v1-abschluss-schreiben"` findet die Kennung in seinem §1 unter *Übernommen von*
  (Temporäre Datei und Probedatei für `sqlite`, Entscheidung in `LH-FA-22.a` vor seinem Code);
  sein §1 *Ausdrücklich NICHT* schließt nur Formatumwandlung, weitere Formate und das
  Standardformat YAML aus. Register: die sechs Einträge mit Beleg und der neue tragen
  `evidence/slice-v1-abschluss-schreiben.md`; die übrigen genannten Einträge bestehen als
  Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft die
  Paarungen erneut.

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
Verhalten des Codes ist unverändert. Nach dem Review (`5bff9bd`) und der Entscheidung des
Architect zu seinen Übergaben (`278d929`) schließt `ca5124f` die Probedatei über `schliessen`
aus `dateiOps` (sonst gleiches Verhalten), erweitert `TestPrepareVerzeichnisNichtBeschreibbar`
um eine vorhandene Datei und eine Verknüpfung auf sie mit `--force` und um die Ursache
(F-537), fügt `TestPrepareProbedateiScheitert` (F-543) und
`TestWriteAndererFehlerBeimAnlegen` (F-544) hinzu, setzt im Kommentar an `Write` den Verweis
für die Rechte auf `LH-FA-07.a` *Rechte* (F-542), zieht §1 (Schicht-Abgrenzung: CLI-Adapter
unverändert, F-540), §3, §6, Handbuch und Abdeckung nach. Nach der Entscheidung des
Architect zur Randform-Rückgabe (`9996a8c`) meldet `pruefe` ein gescheitertes Schließen der
Probedatei als eine Meldung `PGR-E3001`, an die der Fehler des Entfernens als Ursache
hängt (vorher zwei gleichrangige Meldungen über `errors.Join`; `closeErr` und `removeErr`
entfallen); `TestPrepareProbedateiScheitert` prüft je Fall genau eine Meldung, die
Ursachen und die Zahl der liegenden Probedateien und erhält den Fall
`Schliessen und Entfernen`. Kommentar an `Prepare` und Abdeckung nachgezogen; das Handbuch
nennt die Probedatei nicht und bleibt.

**Randformen.** Kein Code-Commit fügt §6 eine Randform hinzu. Zwei Code-Commits ändern in §6
nur den Verweis auf den Test einer schon entschiedenen Randform: `24fc87c` (Rechte bei
ersetzter Verknüpfung → `TestWriteRechteVerknuepfung`) und `ca5124f` (die zwei Randformen aus
`278d929` → `TestWriteAndererFehlerBeimAnlegen`, `TestPrepareProbedateiScheitert`). Der
Code-Commit zur Randform-Rückgabe ändert in §6 nur den Stand der schon entschiedenen
Randform *Schließen der Probedatei scheitert, Entfernen danach* (`9996a8c`) auf ihren Test.
Entschieden ist nur, was §6 nennt. Zwei Punkte
gingen als Frage an den Architect zurück; beide sind entschieden (`b74e16e`, `LH-FA-07.a`),
das Verhalten des Codes blieb, und je ein Test sagt sie zu (DoD-Punkt 1 unten):

| Randform | Frage | Entschieden |
|---|---|---|
| Fehlschlag beim Setzen der Rechte (`Chmod`) oder beim Schließen der temporären Datei | `LH-FA-07.a` *Fehlschlag* nennt Anlegen, Schreiben, Synchronisieren und Verschieben. Der Code behandelt beide wie einen Fehlschlag des Schreibens (`PGR-E3001`, Zieldatei unverändert, temporäre Datei entfernt). Gehören sie zum Schreiben? | ja, beide gehören zum *Fehlschlag*; Test `TestWriteFehlschlag` (Fälle `Rechte`, `Schliessen`) |
| Rechte, wenn `--force` eine symbolische Verknüpfung auf eine Datei ersetzt | Der Code übernimmt die Zugriffsrechte des Ziels (`os.Stat` folgt der Verknüpfung), wie vor dem Slice. *Rechte* sagt „eine ersetzte behält ihre Zugriffsrechte“, *Vorhanden* zählt nach dem Ziel; ist das Ziel gemeint? | ja, die Zieldatei erhält die Zugriffsrechte des Ziels; Test `TestWriteRechteVerknuepfung` |

Ein Fehler von `os.Stat` beim Schreiben ist ein akzeptiertes Negativ (§6); dazu keine Zusage
und kein Test.

`slice-v1-abschluss-sqlite-format` ist Adresse für die *Abgrenzung zu `sqlite`* in §6; die
Sendung steht dort in §1 unter *Übernommen von `slice-v1-abschluss-schreiben`* (Architect,
`278d929`; §3.13). Die Nacharbeit `ca5124f` nennt keine weitere Adresse.

**Größe.** Der Diff `ad1b7ba..65968e6` umfasst 753 hinzugefügte und 29 entfernte Zeilen: Code
71 (`yaml.go`), Unit-Tests 465, Integration 167, Handbuch 23, Abdeckung 24, Plan 3. Die
Nacharbeit `24fc87c` fügt Code 15 (`yaml.go`), Unit-Tests 46 (`schreiben_test.go`,
`export_test.go`), Handbuch 4 und Abdeckung 3 Zeilen (geändert und hinzugefügt) hinzu. Die
Nacharbeit `ca5124f` umfasst Code 8/7 (`yaml.go`, hinzugefügt/entfernt), Unit-Tests 107/7,
Handbuch 2/2, Abdeckung 3/1 und Plan 9/6. Die Nacharbeit zur Randform-Rückgabe ändert Code
(`yaml.go`, 7 hinzugefügt, 13 entfernt), Unit-Tests (`schreiben_test.go`, 26/16) und eine
Zeile der Abdeckung. Ob er
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
| `make gates` | Commit dieser Belege (`9bd65c5`) | Exit 0, vor der Übergabe |
| `make test`, `make lint`, `make docs-check`, `make kopf-check`, `make abdeckung-check` | vor `ca5124f` | Exit 0, 0 Befunde |
| `make abdeckung` | vor `ca5124f` | `abdeckung-unit.md` nachgezogen |
| `make gates` | Commit dieser Belege zu `ca5124f` | Exit 0, vor der Übergabe |
| `make abdeckung` | vor dem Commit zur Randform-Rückgabe | `abdeckung-unit.md` nachgezogen |
| `make gates` | Commit zur Randform-Rückgabe | Exit 0, vor der Übergabe |

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

Zur Nacharbeit `ca5124f` (Übergaben des Reviews): je Mutant eine frische Kopie von `go.mod`,
`go.sum`, `cmd`, `internal` und `test` (`cp -r` ohne `-p`), per Bind-Mount in einem Container
der Stufe `deps` mit `--network=none`; zuerst `gofmt -l ./internal` (leer, sonst Abbruch),
dann `go vet` und `go test -count=1` des Pakets `internal/adapters/driven/recording`. Das
Skript ersetzt genau ein Vorkommen und bricht bei einer anderen Trefferzahl ab. Gefahren:
eine unveränderte Kopie (grün) und 14 Mutanten (Tabelle unten, Zeilen mit F-537, F-543,
F-544), alle formatgerecht und rot im genannten Test. Ein erster Ansatz für *nicht entfernte
Probedatei bleibt* (zweites `os.Remove` nach jedem Entfernen) war in acht Tests rot, aus dem
falschen Grund (das zweite Entfernen scheitert immer); ersetzt durch Mutant *Entfernen nach
gescheitertem Entfernen nachgeholt*, rot nur im genannten Fall.

Zur Nacharbeit nach der Randform-Rückgabe (`9996a8c`): derselbe Weg (frische Kopie je Mutant,
Bind-Mount, Stufe `deps`, `--network=none`, `gofmt -l ./internal` zuerst). Gefahren: eine
unveränderte Kopie (grün) und fünf Mutanten an `pruefe` (Tabelle unten, Zeilen mit
*Randform-Rückgabe*), alle formatgerecht und rot in `TestPrepareProbedateiScheitert`, je
nur in den genannten Fällen.

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
| nicht beschreibbares Verzeichnis `PGR-E3001` auch für eine vorhandene Datei mit `--force` (F-537) | Probedatei für einen vorhandenen Pfad ausgelassen (`case err == nil: return nil`, Mutant D des Reviews) | `TestPrepareVerzeichnisNichtBeschreibbar` (`vorhanden mit --force`, `Verknüpfung mit --force`), `TestPrepareProbedateiScheitert` (`vorhanden=true`) |
| … auch für eine Verknüpfung auf eine vorhandene Datei mit `--force` (F-537, weitere Ausprägung von *vorhanden*) | Probedatei ausgelassen, wenn `os.Lstat` eine Verknüpfung zeigt | `TestPrepareVerzeichnisNichtBeschreibbar` (nur `Verknüpfung mit --force`) |
| Anlegen der Probedatei gescheitert: Fehler als Ursache | Ursache `nil` statt des Fehlers | `TestPrepareVerzeichnisNichtBeschreibbar` (alle drei Fälle) |
| Schließen der Probedatei gescheitert: `PGR-E3001` (F-543) | Fehler von `schliessen` verworfen | `TestPrepareProbedateiScheitert` (`Schliessen`, beide Fälle) |
| … das Schließen läuft über die Tabelle | `probe.Close()` statt `ops.schliessen(probe)` | `TestPrepareProbedateiScheitert` (`Schliessen`, beide Fälle) |
| Schließen der Probedatei gescheitert: Fehler als Ursache | Fehler des Schließens verworfen (`errors.Join(ops.entfernen(name))` als Ursache) | `TestPrepareProbedateiScheitert` (`Schliessen`, `Schliessen und Entfernen`) |
| Entfernen der Probedatei gescheitert: `PGR-E3001` (F-543) | Fehler von `entfernen` verworfen (Mutant P des Reviews) | `TestPrepareProbedateiScheitert` (`Entfernen`, beide Fälle) |
| Entfernen der Probedatei gescheitert: Fehler als Ursache | Ursache `nil` an „Probedatei nicht zu entfernen“ | `TestPrepareProbedateiScheitert` (`Entfernen`) |
| nicht entfernte Probedatei bleibt liegen | nach gescheitertem Entfernen `os.Remove` nachgeholt | `TestPrepareProbedateiScheitert` (`Entfernen`: keine `.probe` im Verzeichnis) |
| Schließen gescheitert: Probedatei dennoch entfernt (Randform-Rückgabe) | nach gescheitertem Schließen nicht entfernt | `TestPrepareProbedateiScheitert` (`Schliessen`: `.probe` liegt; `Schliessen und Entfernen`: Ursache des Entfernens fehlt) |
| Schließen und Entfernen gescheitert: genau eine Meldung (Randform-Rückgabe) | zwei gleichrangige Meldungen über `errors.Join` (Code vor `3a7018b`) | `TestPrepareProbedateiScheitert` (nur `Schliessen und Entfernen`: zwei Meldungen) |
| … sie trägt den Fehler des Entfernens als Ursache | Fehler des Entfernens verworfen (`_ = ops.entfernen(name)`) | `TestPrepareProbedateiScheitert` (nur `Schliessen und Entfernen`) |
| … und den Fehler des Schließens als Ursache | Fehler des Schließens verworfen (`errors.Join(ops.entfernen(name))`) | `TestPrepareProbedateiScheitert` (`Schliessen`, `Schliessen und Entfernen`) |
| … und die Probedatei bleibt liegen | nach gescheitertem Entfernen `os.Remove` nachgeholt, im Zweig des Schließens | `TestPrepareProbedateiScheitert` (nur `Schliessen und Entfernen`: keine `.probe`) |
| Probedatei sofort entfernt | Entfernen weggelassen | `TestPrepareProbedatei`, `TestPrepareVerzeichnis`, `TestRoundtrip`, `TestUebrigGebliebeneDateien` |
| Probedatei `.<Name>.<Zufallsteil>.probe` | Name ohne Punkt, Endung `.pruef` | `TestPrepareProbedatei` |
| Probedatei im Verzeichnis der Zieldatei | Probedatei in `os.TempDir()` | `TestPrepareProbedatei`, `TestPrepareVerzeichnis`, `TestPrepareVorhandeneDatei` |
| Verschieben ersetzt die Verknüpfung, nicht ihr Ziel | Pfad vor dem Schreiben über `filepath.EvalSymlinks` aufgelöst | `TestWriteVerknuepfungErsetzt` |
| temporäre Datei im Verzeichnis der Zieldatei | temporäre Datei in `os.TempDir()` | `TestWriteNameDerTempDatei`, `TestWriteBelegterName`, `TestWriteZehnBelegteNamen`, `TestWriteEntfernenScheitert` |
| Name mit 16 Hexziffern | 7 statt 8 Zufallsbytes | `TestWriteNameDerTempDatei`, `TestWriteBelegterName`, `TestWriteZehnBelegteNamen` |
| exklusiv angelegt, belegter Name nicht überschrieben | `O_TRUNC` statt `O_EXCL` | `TestWriteBelegterName`, `TestWriteZehnBelegteNamen` |
| nach zehn belegten Namen `PGR-E3001` | elf Versuche | `TestWriteZehnBelegteNamen` |
| nur ein belegter Name führt zu einem neuen Versuch, ein anderer Fehler beim Anlegen sofort (F-544) | jeder Fehler führt zu einem neuen Versuch (`err != nil` statt `fs.ErrExist`, Mutant C des Reviews) | `TestWriteAndererFehlerBeimAnlegen` (`Name zu lang`, `Verzeichnis fehlt`: zehn Züge, „kein freier Name“) |
| … auch ein fehlendes Verzeichnis (weitere Ausprägung) | neuer Versuch auch bei `fs.ErrNotExist` | `TestWriteAndererFehlerBeimAnlegen` (nur `Verzeichnis fehlt`) |
| … auch ein zu langer Name (weitere Ausprägung) | neuer Versuch bei jedem Fehler außer `fs.ErrNotExist` | `TestWriteAndererFehlerBeimAnlegen` (nur `Name zu lang`) |
| anderer Fehler beim Anlegen: Fehler des Betriebssystems als Ursache | Ursache `nil` an „nicht anzulegen“ | `TestWriteAndererFehlerBeimAnlegen` (beide Fälle) |
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
5. *Nicht beschreibbares Verzeichnis* ist in den Unit-Tests über `Eingriffe` (`Probe` mit
   `fs.ErrPermission`) gefahren, nicht über die Rechte eines echten Verzeichnisses; die
   Tests laufen im Container als `root`, dem `0555` das Anlegen nicht verwehrt. Am Binary
   hat das Review es mit `--user 65534` gezeigt (Sonde zu F-537); ein E2E-Test dazu fehlt.
6. *Schließen der Probedatei gescheitert* — keine Grenze mehr: Der Architect hat die Frage
   entschieden (`9996a8c`, `LH-FA-07.a` *Verzeichnis*); `record` entfernt die Probedatei
   dennoch, ein Fehler des Entfernens folgt als Ursache derselben Meldung. Beleg sind die
   fünf Zeilen *Randform-Rückgabe* in der Tabelle zu DoD-Punkt 1. Die drei akzeptierten
   Negative der Probedatei in §6 (*Belegter Name*, *Rechte der Probedatei*, *Abbruch
   zwischen Anlegen und Entfernen*) haben keine Zusage und keinen Test.

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
