# Slice slice-v1-abschluss-zeitangaben-pruefen: Zeitangaben prüfen: absteigende Werte beim Lesen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen)

**Berührte Spec-Stellen:** `LH-FA-21.a` · `SPEC-001` · `SPEC-043` · `SPEC-034`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-11.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Beim Laden einer Aufzeichnung — `replay`, `play`, in YAML wie in SQLite — sind `offset_ms`, die innerhalb einer Session absteigen, `PGR-E3003` (Exit-Code 3): die Aufzeichnung ist beschädigt. Gleiche Werte hintereinander und ein kleinerer Wert in der nächsten Session sind gültig. Das ist beobachtbar an einer Datei von Hand, ohne dass `--record-timing` oder `--keep-timing` existieren.

**Übernommen aus `slice-v1-abschluss-zeitangaben`** (dort bis zum Schnitt DoD-Punkt 1, der letzte Halbsatz „ein ungültiger Wert ist `PGR-E3003` (Test)“): die Prüfung der Werte beim Lesen. Den Typ in SQLite (`offset_ms` weder `NULL` noch eine ganze Zahl) liefert der Leser schon (`slice-v1-abschluss-sqlite-format-lesen`, dort §1: Die Spalte hat einen geprüften Typ); die Form des Werts in YAML, auch das Vorzeichen, und das Absteigen liefert dieser Slice. Das Warten beim Einspielen, der andere Liefer-Punkt jenes Slice, liegt bei `slice-v1-abschluss-zeitangaben`.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).


**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--record-timing` und das Schreiben von `offset_ms` mit Messpunkt, Uhr und Option — `slice-v1-abschluss-zeitangaben-aufzeichnen` (dort §1 und DoD, mit der Kennung dieses Slice): Dieser Slice prüft, was eine Datei trägt, nicht wie sie entsteht; das Schreiben setzt die Prüfung voraus, damit seine Roundtrip-Probe die geprüften Leser benutzt.
- `--keep-timing` und das Warten beim Einspielen — `slice-v1-abschluss-zeitangaben` (dort §1 und DoD): Das Warten setzt geprüfte, nicht absteigende Werte voraus; die Prüfung liegt hier, weil sie auch ohne Option gilt (§6, Randform *Prüfung ohne Wunsch*).
- Der Uhr-Port und die Systemuhr — `slice-v1-abschluss-zeitangaben-uhr`: Die Prüfung vergleicht Werte der Datei und braucht keine Uhr.
- Fehlt `offset_ms` einer Interaktion, ist das hier kein Fehler — Bestand und Zusage der Spezifikation: Nur zeitgetreues Einspielen verlangt die Werte (`PGR-E2003`, `slice-v1-abschluss-zeitangaben`).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Produkt-Code im Kern, im PGWire-Adapter, im CLI-Adapter und im Bootstrap — Schicht-Abgrenzung: Die Prüfung liegt im Recording-Adapter an einer Stelle, `Unmarshal` in `yaml.go`, den auch der Leser von SQLite benutzt (Entscheidung des Architect vom 2026-10-11: kein Helfer im Domain Model, keine zweite Kopie); mehr Schichten berührt der Slice nicht (Rückführung in §4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen): Eine Aufzeichnung, in der `offset_ms` innerhalb einer Session absteigen, lädt in YAML und in SQLite nicht: `PGR-E3003`, Exit-Code 3, die Meldung nennt Session und Interaktion und keinen Wert (Test je Format, Negativ); gleiche Werte hintereinander, ein kleinerer Wert am Anfang der nächsten Session und eine Interaktion ohne Wert zwischen zwei mit Wert laden (Test je Format, Boundary); in YAML ist jede Schreibweise von `offset_ms`, die keine kanonische Dezimalzahl ist, `PGR-E3003` (Test je Form: `1.0`, `1.5`, `1e3`, `0x10`, Text in Anführungszeichen, `true`, `null`, Wert über `int64`; Negativ). Tests: `TestZeitangabenAbsteigendBeschaedigt` (Negativ, je Format), `TestZeitangabenGueltig` (Boundary, je Format), `TestZeitangabenFormYAML` (je Schreibweise) und `TestZeitangabenFormSQLite`; die Rangfolge von Formatkennung, Version und Inhalt trägt `TestZeitangabenVersionVorForm`, die Reihenfolge der Meldungen `TestUnmarshalFehlerReihenfolge`. Ohne den Fix aus dem richtigen Grund rot: Mutant R0 der Verifikation (`yaml.go` auf den Stand `8e430ce`) und die Mutanten in §7.
- [x] Die Prüfung gilt beim Laden, nicht erst mit einer Option: `replay` und `play` auf einer solchen Datei enden beim Start mit `PGR-E3003` (Test am Binary, je ein Aufruf); `make abdeckung` ist gefahren, `make abdeckung-check` grün. Test: `TestE2ENichtAbsteigendeZeitangaben` (beide Formate, beide Befehle); die Mutanten E1 und E2 und die der Verifikation am Binary (R0, A1, B1) machen ihn rot.
- [x] `make gates` grün (Stand `0d8f285`, Lauf des Planners vor der Closure; das Ergebnis nach dem letzten Commit steht im Bericht, ein Lauf belegt nur den Stand, auf dem er lief).
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
| `internal/adapters/driven/recording` (`yaml.go`: `Unmarshal` mit Formatkennung und Version vor dem Inhalt, `formVorpruefung`, `offsetForm`, `kanonischeDezimalzahl`, `sessionFromDTO`, `geprueftFromDTO`) | update | Form des Werts am YAML-Baum, Vergleich der `offset_ms` je Session in `sessionFromDTO`; SQLite läuft über `Unmarshal` durch dieselbe Stelle; Meldung ohne Wert aus der Datei |
| Tests im Recording-Adapter (`zeitangaben_pruefen_test.go`, `yaml_test.go`) und `test/integration` (`zeitangaben_pruefen_e2e_test.go`, die Datei `testdata/zeitangaben-absteigend.sqlite`) | neu/update | Negativ und Boundary je Format nach `LH-FA-21`; `replay` und `play` am Binary; die Reihenfolge der Meldungen in `TestUnmarshalFehlerReihenfolge` folgt der Vorprüfung vor dem Dekodieren |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` erzeugt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-sqlite-format-lesen` liegt in `done/` (Voraussetzung: das Lesen der SQLite-Spalte; Schritt 17 der Reihenfolge in §5 von [welle-v1-abschluss](welle-v1-abschluss.md); WIP-Limit 1). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar, oder die Prüfung braucht eine zweite Schicht (Kern, PGWire-Adapter, CLI oder Bootstrap mit Logik); Schnitt dann: YAML und SQLite getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer widerspricht der Lesart von „ignoriert“ in `LH-FA-21.a` (Randform *Prüfung ohne Wunsch*, §6) und verlangt eine Änderung des Lastenhefts; dann zuerst seine Entscheidung.

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

**Randformen** (`AGENTS.md` §3.12) — vom Architect vor dem Code am 2026-10-11 entschieden, Ort: [`LH-FA-21.a`](../../../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Prüfen* und *Einspielen*. Die Entscheidungen der Randformen, die nicht der Nutzer trägt, sind Entscheidungen des Architect (Empfehlung, bis der Nutzer anderes sagt; der Nutzer hat am 2026-10-11 verfügt, ohne Nachfrage weiterzumachen).

- *Absteigend:* Werte, die innerhalb einer Session absteigen, machen die Aufzeichnung beschädigt (`PGR-E3003`); der Vergleich gilt je Session, nicht über Sessions hinweg — entschieden, `LH-FA-21.a` *Prüfen*.
- *Gleicher Wert:* gültig; absteigend ist nur ein Wert, der kleiner ist als der vorige — entschieden ebenda.
- *Mischform:* eine Interaktion ohne `offset_ms` zwischen zwei mit Wert ist beim Laden gültig; verglichen werden die vorhandenen Werte in der Reihenfolge der Interaktionen — entschieden ebenda.
- *Prüfung ohne Wunsch:* Die Prüfung gilt beim Laden jeder Aufzeichnung, in beiden Formaten, für `replay` und `play`, mit und ohne `--keep-timing`; „ignoriert“ meint das Warten, nicht die Prüfung — entschieden ebenda (*Prüfen* und Ende von *Einspielen*). Das Lastenheft wird nicht geändert: Eine beschädigte Datei ist bei negativen Werten schon heute abgelehnt, und „ignoriert“ (LH-FA-21 Boundary) beschreibt das Verwenden. Widerspricht der Nutzer, ist es eine Änderung am Lastenheft (Rückführung in §4).
- *Form des Werts in YAML:* ein Skalar in kanonischer Dezimalschreibweise (`0` oder Ziffern ohne führende Null, ohne Vorzeichen, im Bereich von `int64`); `1.0`, `1.5`, `1e3`, `0x10`, Text in Anführungszeichen, `true`, `null` und ein Wert über dem Bereich sind beschädigt; `offset_ms: null` ist nicht „fehlt“ — entschieden ebenda. Die Bibliothek entscheidet nichts still; die Prüfung am YAML-Baum liegt vor dem Dekodieren (`formVorpruefung`).
- *Meldung:* nennt Session und Interaktion, nie einen Wert aus der Datei und keine Zeilennummer, in beiden Formaten — entschieden ebenda (`BEO-REPO/yaml-meldung-nennt-wert-und-zeile` bleibt für die übrigen Meldungen von `yaml` offen, dieser Slice berührt ihn nicht). Mehrere Fehler in einer Datei: jeder ist `PGR-E3003`, welcher gemeldet wird, ist nicht zugesagt — akzeptiertes Negativ (Code und Exit-Code sind dieselben; ein Rang der Fehler wäre eine Zusage ohne Nutzen).
- *Rangfolge der Prüfung (F-637, Architect am 2026-10-11, Empfehlung bis zur Antwort des Nutzers):* In `yaml` wie in `sqlite` gilt Lesbarkeit der Datei, Formatkennung, Version, Inhalt, und jede Stufe geht den folgenden voran — entschieden, `SPEC-001` *Rangfolge der Prüfung* (für `sqlite` schon in `SPEC-043`). Zugesagt ist damit: `version: 2` ist `PGR-E3002` auch bei Fehlern im Inhalt; falsche Formatkennung ist `PGR-E3003` auch bei unbekannter Version; `version: abc` (und fehlende, `01`, `1.0`, `"1"`, `null`) ist `PGR-E3003` ohne Wert in der Meldung, nie „Version 0“. Zwischen mehreren Fehlern des Inhalts bleibt nicht zugesagt, welcher gemeldet wird (akzeptiertes Negativ, wie oben). Der Code vor dieser Entscheidung (Befund der Form zurückgehalten, Fehler der Bibliothek hinter dem Befund; der Fehler der Bibliothek beim Dekodieren kam vor Formatkennung und Version; `version: abc` ergab `PGR-E3002 Version 0`) entsprach ihr nicht; der Implementer richtet ihn aus (Auftrag in der Übergabe des Architect).
- *Fixture `testdata/zeitangaben-absteigend.sqlite` (F-639, Architect am 2026-10-11):* akzeptiert mit dokumentierter Herkunft — akzeptiertes Negativ: kein Generator, weil `test/integration` weder Adapter noch Bibliothek von `sqlite` einbinden darf und ein Hilfsprogramm dafür den Import nur verschöbe oder `make a-check` lockerte (`AGENTS.md` §3.6). Das Rezept steht im Kommentar der Variable `absteigendSQLite`; ein Schema-Wechsel, der die Datei ungültig macht, lässt `TestE2ENichtAbsteigendeZeitangaben` rot werden, weil die Meldung „Session 1, Interaktion 2“ dann ausbleibt — nachgestellt: Mit der Spalte `offset_ms` umbenannt in der Datei meldet der Leser „no such column: offset_ms“, und der Test ist rot.
- *Eine Session mit einer oder ohne Interaktion:* nichts zu vergleichen, gültig — entschieden ebenda.

**Risiken:**

- Zwei Leser, zwei Kopien des Vergleichs: eine bricht, die andere nicht (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, verkörpert in `AGENTS.md` §3.10) — **Ausgang: entfallen**, denn es gibt nur eine Kopie: `Unmarshal` in `yaml.go` trägt den Vergleich, `sqlite` liest durch dieselbe Stelle; jede Mutation des Vergleichs ist in beiden Formaten rot (M1 bis M6, `TestZeitangabenAbsteigendBeschaedigt` und `TestZeitangabenGueltig` je Format; am Binary E1 mit der Datei von `sqlite`), die Verifikation bestätigte es mit A1 und R0 an beiden Formaten.
- Die Prüfung entscheidet im Code eine Randform still, die §6 nicht nennt (`BEO-REPO/spec-randform-erst-im-review-entschieden`, verkörpert in `AGENTS.md` §3.12): Der Implementer hält an und gibt sie zurück — **Ausgang: eingetreten**: Die Rangfolge von Formatkennung, Version und Inhalt war nicht genannt und stand im Code (F-637). Eine Folge-Arbeit bleibt nicht: Der Architect entschied sie in diesem Slice in `SPEC-001` (`ae1996d`), der Code folgte (`3cf1f57`), die Verifikation fand beides übereinstimmend; gezählt ist der Fall unter `BEO-REPO/spec-randform-erst-im-review-entschieden`.

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

Geschrieben bei der Closure (Planner, 2026-10-11, Stand `0d8f285`).

- **Belege zur DoD:** die Testverweise stehen in der DoD; *Belege des Implementers* und *Nach dem Review* unten tragen Zusage, Mutation und roten Test. Die Verifikation (`docs/reviews/2026-10-11-verifikation-slice-v1-abschluss-zeitangaben-pruefen.md`) bestätigt beide Liefer-Punkte: Ohne den Fix (Mutant R0) werden die Tests aus dem richtigen Grund rot, 93 Dateien am gebauten Binary beider Formate verhalten sich wie zugesagt, 13 Mutanten des Vergleichs und der Form und sechs der Version rot, drei am Binary rot. Der Namensabgleich des Planners: jeder Testname der DoD hat genau eine Zeile `^func <Name>(` unter `internal/` oder `test/`. `make gates` lief grün am Stand `0d8f285` (Planner, vor den Closure-Commits); nach dem letzten Commit steht das Ergebnis im Bericht.
- **Was hat funktioniert:** Der Schnitt von `slice-v1-abschluss-zeitangaben` in vier Teile trug: zwei Liefer-Punkte, eine Schicht, ein Diff von 90 Zeilen Code, in einer Sitzung prüfbar. Die Prüfung liegt an einer Stelle (`Unmarshal`), `sqlite` liest durch sie, und jede Mutation des Vergleichs war in beiden Formaten rot. Die Verifikation baute ihre Dateien mit einer Gegenstelle, die nicht der eigene Code ist (echter PostgreSQL 17, Python mit rohem `sqlite3`), und fand die Zusage, die der Test nur gegen ein Wort prüfte (V-182).
- **Was ging anders als geplant:**
  1. Die Rangfolge von Formatkennung, Version und Inhalt stand nicht in §6 und entschied der Code (F-637); der Architect entschied sie nachträglich in `SPEC-001`, der Code wurde ausgerichtet.
  2. Der Plan folgte der Korrektur im Code nicht (F-638) und trug einen Satzrest (V-184).
  3. Eine Zusage der Spezifikation („nennt keinen Wert“) war nur gegen ein Wort geprüft; der Mutant mit dem Wert der Datei blieb grün (V-182).
  4. Die Fixture `zeitangaben-absteigend.sqlite` kam ohne Herkunft (F-639); der Architect nahm sie mit dokumentiertem Rezept als akzeptiertes Negativ an (§6).
- **Summary-Zeilen:** Review `docs/reviews/2026-10-11-review-slice-v1-abschluss-zeitangaben-pruefen.md`: F-637 (MEDIUM): Rangfolge von Version, Form und Bibliothek im Code entschieden, §6 nennt sie nicht. F-638 (LOW): der Plan folgt der Korrektur nicht. F-639 (INFO): binäre Fixture ohne Herkunft. Verifikation: „0 HIGH · 1 MEDIUM · 1 LOW · 2 INFO. V-182: Zusage ‚Meldung ohne Wert‘ zur Version nur gegen das Wort geprüft, Mutant C7 grün. V-183: ein Lauf von C3 zusätzlich rot, nicht reproduziert. V-184: Rest in §1. V-185: `version: !!int 1` von der Spezifikation nicht genannt.“ Ausgänge: F-637 durch `ae1996d` (Spezifikation) und `3cf1f57` (Code); F-638 und V-184 durch die Nachzüge von Plan 1 und 7 (`e4bb247`, `0d8f285`, die Berichtigung von §1 in der Closure); F-639 durch das Rezept im Kommentar von `absteigendSQLite` und die Probe in §6; V-182 durch `0d8f285` (Test prüft den Wert, C7 rot); V-183 und V-185 unten.
- **Entscheidungen des Planners zur Verifikation:**
  - *V-183, einmal rot, nicht reproduziert:* Kein Eintrag im Register. Es fehlt der Beleg: Die Ausgabe des ersten Laufs ist nicht aufgehoben, drei weitere Läufe der Kopie, zweimal zwölf Wiederholungen und drei Grundläufe waren grün, eine Ursache ist nicht benannt. Eine Zeile ohne Beleg ließe das Register Behauptung zählen. Tritt es wieder auf, entsteht ein Eintrag mit der Ausgabe des Laufs; die Verifikation nennt als Kandidat die Reihenfolge einer `map` in `ladeBeide`, die Tests lesen sie ohne feste Ordnung. Kein Nachzug am Test, solange keine Ausgabe vorliegt.
  - *V-185, `version: !!int 1`:* Als benannte Spezifikationslücke unten geführt und als `BEO-REPO/version-mit-ausdruecklichem-typ-nicht-genannt` (1×) eingetragen; Empfehlung: Der Architect ergänzt die Form in `SPEC-001`, wenn er die Stelle berührt; kein Folge-Slice, kein Lastenheft berührt.
  - *F-639:* einmalig, nicht eingetragen (akzeptiertes Negativ des Architect, Rezept im Kommentar).
- **Steering-Loop-Eintrag:** Geschärfte Regel: Der Beleg einer Zusage ist die Mutation, die nur sie bricht, im Wortlaut der Spezifikation; der rote Test wird aus diesem Grund rot, und rot vieler Tests oder einer anderen Prüfung belegt die Zusage nicht. Anlass: V-182, nach F-620 und V-175 der dritte Beleg derselben Klasse (`BEO-REPO/beleg-traegt-zusage-aus-anderem-grund`, 3×). Die Zeile in Schritt 19 sagt, was die Verifikation bei „nennt keinen Wert“ tat: den Wert der Datei in die Meldung hängen, nicht ein Wort. Liegt in `.claude/commands/implement-slice.md` Schritt 19 (seit slice-v1-abschluss-zeitangaben-pruefen). Zweitens benannte Spec-Lücke: `SPEC-001` nennt `version: !!int 1` nicht (V-185, `BEO-REPO/version-mit-ausdruecklichem-typ-nicht-genannt`, 1×).

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist wieder aufgetreten (F-638); die Regel bleibt. §3.10 (seit welle-extended-query) ist wieder aufgetreten (V-182); die Regel bleibt, der Sensor ist mit `slice-harness-mutation` geplant. §3.12 (seit slice-harness-randformen-vor-code) ist in der Hälfte *im Review entschieden* wieder aufgetreten (F-637); die Regel bleibt. §3.11 und §3.13 sind nicht wieder aufgetreten: Der Nehmer nimmt an (unten), die Zusagen in Kommentaren und Plan trugen ihre Tests.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `0d8f285` (Zähler = Dateien unter `evidence/`); jeder Beleg trägt `evidence/slice-v1-abschluss-zeitangaben-pruefen.md`.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden`: **Beleg**, 22× → 23× (F-637), Stand verkörpert, bleibt.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 29× → 30× (F-638, V-184), Stand verkörpert, bleibt.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 26× → 27× (V-182), Stand verkörpert, Sensor `slice-harness-mutation` geplant, bleibt.
  - `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund`: **Beleg**, 2× → 3× (V-182), Schwelle erreicht, Ausgang **verkörpert** mit diesem Slice: Zielort `.claude/commands/implement-slice.md` Schritt 19, Herkunfts-Anker `seit slice-v1-abschluss-zeitangaben-pruefen`.
  - Neu: `BEO-REPO/version-mit-ausdruecklichem-typ-nicht-genannt` (1×, V-185), Stand offen.
  - Berichtigt: `BEO-REPO/play-ohne-sitzung-ohne-meldung`: Der Satz, die Spezifikation sage zu `play` ohne Session nichts, ist überholt; `LH-FA-20.a` *Start* entscheidet, `play` verbindet sich nicht und endet mit Exit-Code `0` (Empfehlung des Architect: bewusst schweigen). Stand **gestrichen** mit dieser Begründung, die Beobachtung steht bei 1×.
  - Ohne Beleg: `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` (bleibt 2×; die Meldungen der Prüfung nennen keinen Wert, der Slice berührt die übrigen Meldungen von `yaml` nicht), `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (bleiben 2×; der Schnitt blieb bei zwei Liefer-Punkten und einer Schicht, die Rückführungen aus §4 traten nicht ein), `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt 7×; der Nehmer nimmt an).
  - Einmalig und nicht eingetragen: V-183 (kein Beleg), F-639 (akzeptiertes Negativ).
- **Folge-Slices:** keiner neu. Der nächste Schritt nach §5 von [welle-v1-abschluss](welle-v1-abschluss.md) ist `slice-v1-abschluss-zeitangaben-uhr` (Schritt 20); `slice-v1-abschluss-zeitangaben-doku` trägt Handbuch und README zu den vier Teilen.
- **Risiken aus §6:** zwei. *Zwei Leser, zwei Kopien*: **entfallen**, es gibt eine Kopie des Vergleichs in `Unmarshal`, beide Formate lesen durch sie, die Mutanten sind in beiden rot. *Randform im Code still entschieden*: **eingetreten** (F-637), erledigt in diesem Slice durch `ae1996d` und `3cf1f57`; ein Folge-Slice bleibt nicht, gezählt unter `BEO-REPO/spec-randform-erst-im-review-entschieden`. Die Randformen in §6 sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker: Der Eintrag trägt die Zeile `liegt in` mit `.claude/commands/implement-slice.md` Schritt 19; der Absatz steht dort mit dem Anker `seit slice-v1-abschluss-zeitangaben-pruefen`, und `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund/state.md` nennt Zielort und Anker. Folge-Slice: Der Nehmer `slice-v1-abschluss-zeitangaben` liegt in `next/`; `grep -n "slice-v1-abschluss-zeitangaben-pruefen"` findet die Kennung des Gebers in seinem §1 (Abgabe des Halbsatzes „ein ungültiger Wert ist `PGR-E3003`“ und *Ausdrücklich NICHT*) und in seinem Start, der den Slice in `done/` verlangt; sein §1 schließt die Sendung nicht aus. Ebenso `slice-v1-abschluss-zeitangaben-aufzeichnen` (§1, Start) und `slice-v1-abschluss-zeitangaben-doku` (§1, DoD: Handbuch und README zur Prüfung). Register: die vier Einträge mit Beleg und der neue tragen `evidence/slice-v1-abschluss-zeitangaben-pruefen.md`; die übrigen genannten bestehen als Verzeichnis mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft die Paarungen erneut.

### Belege des Implementers

**Stand.** Geliefert sind beide Liefer-Punkte der DoD mit `78a6b4c`, `341b569` und `3cf1f57`. Geändert: `internal/adapters/driven/recording/yaml.go`
(`Unmarshal`, `formVorpruefung`, neu `offsetForm` und `kanonischeDezimalzahl`, `sessionFromDTO`, `geprueftFromDTO`),
`yaml_test.go` (`TestUnmarshalFehlerReihenfolge`), `docs/user/abdeckung-*.md`; neu `zeitangaben_pruefen_test.go` und
`test/integration/zeitangaben_pruefen_e2e_test.go` mit der Datei `testdata/zeitangaben-absteigend.sqlite` (32 KB, eingebettet; ein Test unter `test/integration` darf weder den Adapter noch die Bibliothek von `sqlite` einbinden, `make a-check`). Handbuch, `README.md`, Kern, PGWire-Adapter, CLI und Bootstrap sind
unverändert. Die Prüfung der Form des Werts liegt in `formVorpruefung` und läuft vor dem Dekodieren; die negative
Prüfung `offset_ms negativ` in `geprueftFromDTO` entfällt, weil `-` keine kanonische Dezimalzahl ist. Der Befund der
Vorprüfung gilt nach Formatkennung und Version (`TestZeitangabenVersionVorForm`). `--keep-timing` gibt es in diesem
Stand nicht (`slice-v1-abschluss-zeitangaben`); die Tests am Binary rufen `replay` und `play` ohne die Option.
**Größe:** `git diff --numstat 8e430ce..3cf1f57` (ohne Plan, Review und Spezifikation) zählt 456 Zeilen hinzu und 20 weg; davon 90 Go-Code (mit Kommentaren, ohne
Kommentare weniger), 359 Go-Tests, 7 Tabellen, dazu die Datei von 32 KB.

**Zusage · Mutation · roter Test.** Jede Mutation lief in einer frischen Kopie (`cp -r` ohne `-p`: `go.mod`, `go.sum`,
`cmd/`, `internal/`, `test/`) mit einer Ersetzung in `yaml.go` und `go test` im Bind-Mount des gepinnten Images (der
Mutant kommt ohne Build-Kontext an); E1 und E2 mit `make test-integration` in einer Kopie des ganzen Baums. Die
Mutanten M1 bis M18, V1 bis V11 und C7 stehen unten.

| Zusage | Mutation | roter Test |
|---|---|---|
| absteigende Werte einer Session sind `PGR-E3003`, in `yaml` und `sqlite` | Vergleich entfernt (M2) | `TestZeitangabenAbsteigendBeschaedigt` (sechs Fälle, je Format); am Binary E1 `TestE2ENichtAbsteigendeZeitangaben` |
| ein gleicher Wert hintereinander ist gültig | `<=` statt `<` (M1) | `TestZeitangabenGueltig`, `TestZeitangabenAbsteigendBeschaedigt/nach_gleichen_Werten` |
| der Vergleich gilt je Session, ein kleinerer Wert in der nächsten ist gültig | Vorgänger über Sessions hinweg behalten (M3) | `TestZeitangabenGueltig/kleinerer_Wert_in_der_nächsten_Session` |
| Mischform: Vergleich der vorhandenen Werte, auch über eine Interaktion ohne Wert | ein fehlender Wert setzt den Vorgänger zurück (M4); er zählt als 0 (M5) | `TestZeitangabenAbsteigendBeschaedigt/über_eine_Interaktion_ohne_Wert` (M4, M5), `TestZeitangabenGueltig` (M5) |
| der Vergleich gilt gegen den vorigen Wert, nicht gegen den ersten | nur gegen den ersten Wert (M6) | `TestZeitangabenAbsteigendBeschaedigt/nicht_benachbart` |
| die Meldung nennt keinen Wert (Vergleich) | Wert in der Meldung (M7) | `TestZeitangabenAbsteigendBeschaedigt` |
| die Meldung nennt Session und Interaktion (Vergleich) | Interaktion nullbasiert (M8) | `TestZeitangabenAbsteigendBeschaedigt` |
| Form: einfacher Skalar, kein Text in Anführungszeichen, kein ausdrücklicher Typ | `Style` nicht geprüft (M9) | `TestZeitangabenFormYAML` (`"…"`, `'…'`, `!!int`); am Binary E2 (Text) |
| Form: keine führende Null | `s[0] == '0'` entfernt (M10) | `TestZeitangabenFormYAML/führende_Null` |
| Form: Bereich von `int64` | Bereich ungeprüft (M11) | `TestZeitangabenFormYAML/über_int64`, `/weit_über_int64` |
| Form: kein Vorzeichen | `-` erlaubt (M12) | `TestZeitangabenFormYAML/Minus`, `TestFelderDerVersion1`, `TestUnmarshalFehlerReihenfolge` |
| Form: `null`, `~` und leerer Wert sind beschädigt, nicht „fehlt“ | `!!null` angenommen (M13) | `TestZeitangabenFormYAML/null`, `/Tilde`, `/leerer_Wert` |
| die Meldung nennt keinen Wert und keine Zeile (Form) | Wert (M14) und Zeile (M15) in der Meldung; Vorprüfung nach dem Dekodieren (die Meldung der Bibliothek nennt Wert und Zeile, E2) | `TestZeitangabenFormYAML` je Schreibweise; E2 am Binary |
| die Meldung nennt Session und Interaktion (Form) | Session und Interaktion vertauscht (M16) | `TestZeitangabenFormYAML` (Session 2, Interaktion 3) |
| die Form wird in jeder Session geprüft | nur die erste Session (M17) | `TestZeitangabenFormYAML` |
| die Meldung zur Form trägt Session und Interaktion auch in `sqlite` | Befund der Vorprüfung als Ursache umhüllt (M18) | `TestZeitangabenFormSQLite` |
| eine kanonische Version außer 1 ist `PGR-E3002`, auch `0`, negativ, über `int64`, auch bei Fehlern im Inhalt | `0` als `PGR-E3003` (V6); negativ als `PGR-E3003` (V6b); über `int64` als `PGR-E3003` (V7) | `TestZeitangabenVersionVorForm` (`Version 0`, `-1`, `über int64`) |
| `version: 2` ist `PGR-E3002` auch bei unbekanntem Schlüssel oder Fehler im Inhalt | Decoder vor der Prüfung der Version (V2); Prüfung des Inhalts vor der Version (V10) | `TestZeitangabenVersionVorForm` (`Version 2`, `Version 2, unbekannter Schlüssel`) |
| fehlende Version oder Wert, der keine kanonische Dezimalzahl ist: `PGR-E3003` ohne Wert | `abc` als Version 0 (V3); fehlend als 0 (V4); führende Null erlaubt (V5); Text in Anführungszeichen angenommen (V8); die Meldung trägt das Wort „Version“ (V11) | `TestZeitangabenVersionVorForm` (`abc`, `fehlt`, `01`, `-0`, `1.0`, `"1"`, `null`, `0x2`) |
| die Meldung zu einer Version, die keine kanonische Dezimalzahl ist, nennt den Wert der Datei nicht | `v.Value` in der Meldung (C7) | `TestZeitangabenVersionVorForm` je Fall (`abc`, `01`, `-0`, `1.0`, `"1"`, `null`, `0x2`; der Text enthält den Wert nicht) |
| eine falsche oder fehlende Formatkennung ist `PGR-E3003`, auch bei unbekannter Version | Version vor Formatkennung (V1); Formatkennung nicht verglichen (V9) | `TestZeitangabenVersionVorForm`, `TestUnmarshalFehler/fremdes_Format` |
| `replay` und `play` enden beim Start mit Exit-Code 3 | Vergleich entfernt (E1); Prüfung des Styles entfernt (E2) | `TestE2ENichtAbsteigendeZeitangaben` (beide Formate, beide Befehle) |

Summe: 30 Mutanten der Unit-Tests (M1 bis M18, V1 bis V11, C7; V1 ist die Umstellung von Formatkennung und Version in einem Mutanten) und 2 am Binary (E1, E2); rot 32, grün 0 (E1 mit der Datei von `sqlite`: der Aufruf von `replay` endet nicht binnen 15 s). Ein Mutant des Vergleichs trifft beide
Formate zugleich, weil `sqlite` durch dieselbe Stelle liest; die Fälle von `TestZeitangabenAbsteigendBeschaedigt`
und `TestZeitangabenGueltig` laufen je Format (Teilnamen `, yaml` und `, sqlite`), beide Teile wurden rot gesehen.
Die Prüfung der Art des Knotens entfiel nach der ersten Fassung: Eine Liste oder Abbildung hat den Wert "", den
`kanonischeDezimalzahl` ablehnt; `TestZeitangabenFormYAML` hält beide als Fälle.

**Gegenprobe zur Fixture (Lehre V-175).** Jeder rote Fall trifft die genannte Prüfung: Die Fälle der Form haben davor
eine Interaktion mit gültigem Wert und sind sonst gültig, die Fälle des Vergleichs sind sonst gültig und tragen
keine Form, die die Vorprüfung ablehnt; die roten Läufe nennen die Prüfung (M7, M14) oder den Ort der Meldung (M8, M16).

**Läufe (Stand `3cf1f57`).** `make test` grün, `make test-integration` grün (mit `TestE2ENichtAbsteigendeZeitangaben`),
`make lint` ohne Befund, `make a-check` ohne Befund, `make abdeckung` gefahren und die Tabellen mit `78a6b4c` eingetragen. `make gates` läuft nach
dem letzten Commit auf sauberem Baum; sein Ergebnis steht im Bericht, nicht hier (ein Gate-Lauf belegt nur den Stand,
auf dem er lief).

**Namensabgleich.** Jeder Testname dieses Abschnitts hat `^func <Name>(` unter `internal/` oder `test/`.

**Nach dem Review (F-637, F-639).** `Unmarshal` liest Formatkennung und Version aus dem YAML-Baum vor dem Decoder und vor
`formVorpruefung` (`kopfVorpruefung`), in der Rangfolge von `SPEC-001`: Lesbarkeit, Formatkennung, Version, Inhalt. Das Zurückhalten
des Befunds der Vorprüfung und `err != nil && vor == nil` entfallen; kein Pfad meldet mehr „Version 0“. M19 und M20 der ersten
Fassung sind durch V1 bis V11 ersetzt; „Prüfung der Version vor dem Dekodieren“ (früher M20) belegt V2 mit
`TestZeitangabenVersionVorForm/Version_2,_unbekannter_Schlüssel`. Die Fixture hat ihr Rezept im Kommentar von `absteigendSQLite`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand des Schnitts, 2026-10-11). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) — die Prüfung ist ein neuer Vertrag, jede Zusage eine Mutation (Risiko in §6).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — Randformen stehen in §6.
- `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` (2×, offen) — die Meldung der Prüfung nennt keinen Wert (Randform *Meldung*); sie berührt den Eintrag nicht neu.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — dieser Schnitt lässt keine Hälfte an der Grenze: zwei Liefer-Punkte, eine Schicht, kleiner Diff; die Größe misst der Architect vor dem Code, und ein Befund dort wäre ein weiterer Beleg.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten: eine (Recording-Adapter, eine Stelle für beide Formate); nachgezählt vom Architect am 2026-10-11 (SQLite liest über `Unmarshal`, ein Helfer im Domain Model entfällt). Nachgezählt beim Schnitt (2026-10-11, `AGENTS.md` §3.13); die Sendung aus `slice-v1-abschluss-zeitangaben` (DoD-Halbsatz) steht in DoD-Punkt 1.

Nachgezählt am Diff (Implementer, `3cf1f57`): Liefer-Punkte 2, Schichten eine (Recording-Adapter); Produkt-Code 90 Zeilen mit Kommentaren in `yaml.go`, ein Bootstrap ohne Änderung.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
