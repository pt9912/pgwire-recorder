# Slice slice-v1-abschluss-zeitangaben-aufzeichnen: Zeitangaben aufzeichnen (`--record-timing`)

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

**Berührte Spec-Stellen:** `LH-FA-21.a` · `LH-FA-17.a` · `SPEC-001` · `SPEC-002` · `SPEC-004` · `SPEC-041` · `SPEC-043` · `SPEC-034` · `ARC-001` · `ARC-002`

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

**Ziel:** `record --record-timing` trägt an jeder Interaktion `offset_ms` ein — ganze Millisekunden zwischen der Annahme der Verbindung und dem Eintreffen der ersten Client-Nachricht der Interaktion, mit der monotonen Uhr des Ports —, in YAML und in SQLite; ohne die Option fehlt das Feld und die Ausgabe ist unabhängig vom Zeitverhalten. Die Aufzeichnung lässt sich in beiden Formaten laden (Roundtrip).

**Übernommen aus `slice-v1-abschluss-zeitangaben`** (dort bis zum Schnitt DoD-Punkt 1 und §3, Zeilen *PGWire-Adapter*, *Recording-Adapter* und *CLI-Adapter*, Option `--record-timing`): das Aufzeichnen der Zeitangaben.

**Übernommen aus `slice-v1-abschluss-sqlite-format-lesen`** (dort §1, Abgrenzung): die Zeitangaben in der SQLite-Datei. Das Lesen der Spalte `interaction.offset_ms` liefert jener Slice (nur Spalte mit geprüftem Typ: negative und nicht ganze Werte sind `PGR-E3003`); dieser Slice schreibt sie mit `--record-timing` und belegt den Roundtrip in beiden Formaten.

**Übernommen aus `slice-v1-abschluss-einspielen`** (dort §1, Abgrenzung): Zeitangaben; hier die Hälfte beim Aufzeichnen, die Hälfte beim Einspielen liegt in `slice-v1-abschluss-zeitangaben`.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).


**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Prüfung absteigender Werte beim Lesen — `slice-v1-abschluss-zeitangaben-pruefen` (dort §1 und DoD, mit der Kennung dieses Slice): Dieser Slice schreibt aufsteigende Werte und lädt sie über die geprüften Leser; die Prüfung selbst liegt dort und steht in der Reihenfolge vor diesem Slice.
- Der Uhr-Port, die Fake-Uhr und die Systemuhr — `slice-v1-abschluss-zeitangaben-uhr` (dort §1 und DoD): Dieser Slice benutzt den Port, wie jener ihn liefert; ändert er die Signatur, ist das ein Plan-Wechsel, nicht eine Ergänzung.
- `--keep-timing`, `--timing-mode`, `--timing-reference` und das Warten beim Einspielen — `slice-v1-abschluss-zeitangaben` (dort §1 und DoD): Abnahmeszenario 13 läuft erst, wenn beide Hälften da sind; es schließt in jenem Slice, der hinter diesem steht.
- Eine neue Version des Recording-Formats — jedes neue Feld, auch ein optionales, erhöht `version` (`SPEC-001`); `offset_ms` steht schon in Version 1 (`SPEC-002`, `SPEC-043`), dieser Slice fügt kein Feld hinzu (Rückführung in §4).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Änderungen am PGWire-Adapter — Schicht-Abgrenzung, solange der Messpunkt im Kern liegen kann (Randform *Messpunkt* in §6); liegt er im Adapter, ist das die dritte Schicht und die Rückführung in §4.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen): `record --record-timing` trägt `offset_ms` an jeder Interaktion ein, einfach wie Extended, gemessen mit der Uhr des Ports; die Aufzeichnung lädt in YAML und in SQLite per Roundtrip mit den Werten (E2E am Binary mit beiden Formaten und Test mit Fake-Uhr im Service); ohne die Option fehlt das Feld und die Ausgabe ist byte-gleich der heutigen (Test, Boundary).
- [ ] Option `--record-timing` für `record` mit Umgebungsvariable und Schlüssel der Konfigurationsdatei nach den Regeln der Optionstabelle (`LH-FA-17.a`, Default `false`, Priorität Option vor Variable vor Datei, ungültiger Wert `PGR-E2001` bzw. `PGR-E2004`); bei `replay` und `play` unbekannt (`PGR-E2001`); `config show` und `--help` nennen sie (Test am Binary).
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
| `internal/hexagon/services` (Record-Service) | update | `offset_ms` je Interaktion mit der Uhr des Ports, Basis je Session = Annahme der Verbindung |
| `internal/hexagon/ports/driving` (Recorder) | update (nur wenn der Messpunkt Zeit mitgeben muss) | Entscheidung des Architect (§6) |
| `internal/adapters/driving/cli` | update | Option `--record-timing`, Umgebungsvariable, Schlüssel im Abschnitt `record:`, `config show`, Hilfe |
| `internal/bootstrap` | update | Uhr an den Record-Service reichen (reine Verdrahtung, keine Schicht, `AGENTS.md` §3.13) |
| Tests, `test/integration` | neu/update | Service mit Fake-Uhr, Roundtrip beider Formate am Binary, Optionstabelle; `docs/user/abdeckung-*.md` von `make abdeckung` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-zeitangaben-pruefen` und `slice-v1-abschluss-zeitangaben-uhr` liegen in `done/`, dazu `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` (der Roundtrip in SQLite; Schritt 21 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md); WIP-Limit 1). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Messpunkt liegt im PGWire-Adapter — das wären Kern, PGWire-Adapter und CLI-Adapter, drei Schichten —, oder der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Messpunkt im Adapter mit der Signatur des Ports als eigener Slice vor der Option.
- `in-progress` → `next`: Das Zeitmodell verlangt eine Änderung am Recording-Format über das in Version 1 enthaltene Feld `offset_ms` hinaus; jedes neue Feld, auch ein optionales, erhöht `version` (`SPEC-001`) — zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Die Fake-Uhr bildet die monotone Uhr nicht ab oder der Port aus `slice-v1-abschluss-zeitangaben-uhr` trägt den Messpunkt nicht — Carveout, dann zuerst der Plan jenes Slice.

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

**Randformen** (`AGENTS.md` §3.12) — Liste für den Architect, geschrieben beim Schnitt vom 2026-10-11. Entschieden ist eine Randform erst, wenn die Spezifikation (Technik-Stratum) sie vor dem ersten Code-Commit trägt; die Planung hat die Spezifikation nicht geändert. Je Punkt steht, wo er entschieden ist, sonst *offen*: Der Architect entscheidet vor dem Code oder gibt die Frage an den Nutzer (Lastenheft-Bezug). Der Implementer entscheidet keine offene Randform (`.claude/commands/implement-slice.md`, Randform-Rückgabe).

- *Messpunkte:* Basis ist die Annahme der Verbindung, Ende die Ankunft der ersten Client-Nachricht der Interaktion — entschieden in `LH-FA-21.a` *Aufzeichnen*; ob „Ankunft“ der Moment ist, in dem der Adapter die Nachricht liest, oder der, in dem der Use Case sie übernimmt (bei hintereinander geschickten einfachen Anfragen liegen beide auseinander) — offen, Entscheidung des Architect; sie legt fest, ob der Slice den PGWire-Adapter berührt (§4).
- *Rundung:* `offset_ms` ist eine ganze Zahl ≥ 0 Millisekunden — entschieden; ob abgeschnitten oder gerundet wird — offen (Empfehlung des Planners: abschneiden, damit der Wert nie größer als die gemessene Zeit ist).
- *Erste Interaktion:* kann 0 sein; `offset_ms` ist nie negativ — entschieden (ganze Zahl ≥ 0).
- *Nicht absteigend:* Die erzeugte Aufzeichnung darf die Prüfung aus `slice-v1-abschluss-zeitangaben-pruefen` nicht verletzen; Messung in Ankunftsreihenfolge je Session, auch wenn beide Richtungen einer Extended-Session gleichzeitig laufen (`SPEC-041`) — offen, wie die Reihenfolge der Interaktionen und der Werte zusammenfällt; der Roundtrip lädt über die geprüften Leser.
- *Extended-Interaktion:* das Feld steht auf der Ebene der Interaktion (`SPEC-002`, `SPEC-041`), Messpunkt ist die erste Client-Nachricht der Interaktion (`Parse`, `Bind`, … der Gruppe) — entschieden im Wortlaut; ihre Bestimmung bei überlappenden Gruppen — offen.
- *Lebendprüfung, Fehlerantwort, Abbruch und Herunterfahren:* `record` behandelt eine Lebendprüfung wie jede andere Anfrage (`LH-FA-09.a`), sie trägt das Feld; ob eine Interaktion, die beim Herunterfahren nicht abgeschlossen wird, den Wert behält — offen.
- *Mehrere Sessions und `--record-empty-sessions`:* Basis je Session ihre eigene Annahme (`LH-FA-21.a`: die Zeit zählt je Session); eine Session ohne Interaktion trägt kein Feld — entschieden.
- *Ohne die Option:* Das Feld fehlt, und die Ausgabe ist unabhängig vom Zeitverhalten des Laufs, ausgenommen `Flush`-Gruppen ohne wartenden Client (`SPEC-004`) — entschieden.
- *Option:* Name, Umgebungsvariable `PGWIRE_RECORDER_RECORD_TIMING`, Schlüssel, Default und Wertemenge stehen in der Optionstabelle und in `LH-FA-17.a` — entschieden; bei `replay` und `play` ist die Option unbekannt (`PGR-E2001`).

**Risiken:**

- Zeitabstände auf einem ausgelasteten Rechner schwanken; der E2E-Test prüft nur Ordnung und Untergrenze (nicht absteigend, nicht negativ), nicht eine bestimmte Dauer — **Ausgang:** offen bis Closure.
- Der Messpunkt liegt entgegen der Annahme im PGWire-Adapter, und der Slice berührt drei Schichten (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, 2×, offen) — **Ausgang:** offen bis Closure; die Rückführung steht in §4.

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

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand des Schnitts, 2026-10-11). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — neuer Vertrag, Randformen in §6.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) — Der Messpunkt ist der Grund, dass die Größe erst die Prüfung des Architect zeigen kann; die Rückführung in §4 nennt die Bedingung vorab. Ein Eintreten wäre ein weiterer Beleg.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — Hilfetext und Abdeckungs-Deklaration sagen nur zu, was ein Test prüft.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten: zwei (Kern mit dem Record-Service, CLI-Adapter); der Bootstrap reicht die Uhr weiter und zählt als reine Verdrahtung nicht; Messpunkt im PGWire-Adapter wäre die dritte Schicht (§4). Nachgezählt beim Schnitt (2026-10-11, `AGENTS.md` §3.13); die Sendungen aus `slice-v1-abschluss-sqlite-format-lesen` und `slice-v1-abschluss-einspielen` stehen in §1 mit den Kennungen der Geber, in DoD-Punkt 1.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
