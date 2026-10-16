# Slice slice-harness-meldungskatalog-gate: Gate für Code-Tabelle und Meldungskatalog

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein
Abnahmeszenario und kein Meilenstein hängt an ihm. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom 2026-10-07:
in der wellenlosen Reihe nach `slice-harness-integration-wait` und vor
`slice-harness-abdeckung-gate` (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler) (stabile Meldungscodes je Ursache). Bindung an Entscheidungen: [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) (Folgepflicht: ein Gate, das Code-Tabelle und Katalog abgleicht).

**Berührte Spec-Stellen:** `SPEC-034` (gelesen, Tabelle der Codes; geändert nur, wenn der Architect einen Bestands-Befund dort berichtigt) · [`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (Vertrag des Gates, vom Architect vor dem Code geschrieben; die nächste freie Nummer ist danach 055) · `SPEC-052` Punkt 13 (Beleg des Handbuchs in der schnellen Kette, vom Architect nachgezogen)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Gate meldet rot, wenn ein Meldungscode `PGR-…` im Go-Code nicht im
Meldungskatalog der Spezifikation (Tabelle in `SPEC-034`) steht, und umgekehrt, wenn
ein Code des Katalogs im Go-Code fehlt, es sei denn, er ist ein Rückfall `…000`; und
wenn die Codes der Fehlercode- und Warnungstabellen des Handbuchs nicht genau die der
Code-Tabelle sind. Was als Quelle, Code, Befund und Sperre zählt, entscheidet
[`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (Architect,
vor dem ersten Code-Commit); die Randformen dazu stehen in §6 mit ihrem Punkt. Geliefert wie
`make kopf-check` und `make abdeckung-check`: ein Prüfskript unter `tools/harness/`
mit Gate-Ziel, eine Gegenprobe als eigenes Gate-Ziel, beide über ein Fragment unter
`harness/mk/` an `GATE_CHECKS`, eine Sensor-Datei unter `harness/sensors/` und je eine
Zeile in `harness/README.md` §Sensors.

**Herkunft:** [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) nennt als
Folgepflicht ein Gate, das Code-Tabelle und Katalog abgleicht, sobald beide
existieren; seit `slice-replay-semantik-meldungscodes` liegt die Code-Tabelle in
`internal/hexagon/model/fehler.go`. Die Closure von welle-replay-semantik fand, dass die
Folgepflicht keine Slice-Kennung hat (Nebenbefund 4 in
`welle-replay-semantik-results.md`): welle-v1-abschluss §6 und
`slice-v1-abschluss-container` §1 nannten das Gate als „Folge-Slice“ ohne Kennung.
Entscheidung des Nutzers vom 2026-10-07 (Audit nach [`AGENTS.md`](../../../../AGENTS.md)
§3.13, Fund 1): eigener wellenloser Slice, verglichen wird der Go-Code mit dem Katalog
der Spezifikation. Das Handbuch kommt als dritte Quelle dazu (Empfehlung des Architect
vom 2026-10-17, §6 Punkt 8; der Nutzer hat es nicht ausgeschlossen): Es kostet eine
Funktion im Skript und schützt `slice-v1-abschluss-container`, bevor der den Katalog der
Betriebsdokumentation fortschreibt.

**Übernommen aus welle-v1-abschluss** (dort §6, Out-of-Scope) **und aus
`slice-v1-abschluss-container`** (dort §1, Abgrenzung): das Gate für Code-Tabelle und
Katalog; das ist das Ziel oben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Codes erzeugen, die der Katalog führt und der Code noch nicht — Schicht-Abgrenzung:
  kein Verhalten im Produkt-Code und keine Änderung an `internal/hexagon/model/fehler.go`
  oder `SPEC-034`; das Gate ist am Stand des Starts grün (§6 *Bestand*), die vier Codes des
  Katalogs ohne Konstante sind Rückfälle, die `SPEC-054` Punkt 4b zulässt.
- Klasse, Exit Code und Form des Fehlertexts eines Codes prüfen — Bestand bleibt
  stehen: Das prüfen die Tests aus `slice-replay-semantik-meldungscodes`
  (`fehler_test.go`, Kopf und Exit Code in `cli_test.go`, `bootstrap_test.go`); das
  Gate gleicht Mengen von Codes ab, kein Verhalten.
- Codes, die andere Abschnitte der Spezifikation, das Lastenheft oder Slice-Pläne in
  Prosa nennen — anderer Vorgang: Der Abgleich gilt der Code-Tabelle und dem Katalog,
  nicht jeder Nennung; eine Prosa-Nennung eines unbekannten Codes ist Urteil des
  Reviews.
- Den Katalog neu ordnen oder Codes neu vergeben — anderer Vorgang: Vergabe ist
  Spezifikationsarbeit eines Slice, der eine Ursache einführt (`SPEC-034`
  §Stabilität); dieser Slice berichtigt nur einen Bestands-Befund, den sein Vertrag
  rot macht, und nur nach Entscheidung des Architects (§6, *Bestand*).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler):
      Der Vertrag `SPEC-054` steht vor dem ersten Code-Commit in `harness/vertraege/`
      und entscheidet jede Randform aus §6; das Prüfskript (bash und awk, ohne Docker)
      hängt mit seinem Ziel über ein Fragment unter `harness/mk/` an `GATE_CHECKS`, ist
      am Stand des Starts grün, ohne dass `fehler.go` oder `SPEC-034` sich ändern, und
      steht in der schnellen Kette; die Gegenprobe steht in `GATE_SCHNELL_AUSNAHME`
      (`harness/mk/gates-schnell.mk`) und in `LANGSAM` von `tools/harness/gates-schnell-gegenprobe.sh`.
- [ ] Gegenprobe als eigenes Gate-Ziel in Temp-Bäumen: je Punkt von `SPEC-054` und je
      Element einer Aufzählung darin (`AGENTS.md` §3.10) ein Fall, der abgelehnt bzw.
      angenommen werden muss — darunter ein Code nur im Go-Code, nur im Katalog, nur im
      Handbuch, in falscher Form; dazu in einem Temp-Baum mit
      Makefile und Fragment, dass beide Ziele an `GATE_CHECKS` hängen und über `make`
      bei rotem Skript mit Exit ungleich 0 enden. Je Zusage des Skripts ist die Mutation
      vom Implementer rot gesehen ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die
      Tabelle Zusage · Mutation · roter Test steht in §7.
- [ ] Sensor-Datei `harness/sensors/meldungskatalog-check.md` (Vertrag, Grenze, Ausgabe
      und Ausgänge, wie `harness/sensors/kopf-check.md`) und je eine Zeile für beide
      Ziele in `harness/README.md` §Sensors mit Bindung an
      [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) und den neuen Vertrag.
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
| `harness/vertraege/SPEC-054-meldungskatalog.md`, `README.md` des Ordners, `SPEC-052-gate-kette.md` Punkt 13 | neu, update (Architect, vor Code, geliefert) | Vertrag des Gates; Beleg des Handbuchs in der schnellen Kette |
| `docs/plan/adr/` | keine ADR | [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) trägt die Folgepflicht, die Form folgt dem Muster der Sensoren; die Gründe (Handbuch als dritte Quelle, Rückfälle zugelassen, Go-Test unmöglich wegen des Build-Kontexts) stehen im Vertrag |
| `tools/harness/meldungskatalog-check.sh` | neu | liest die Codes aus `fehler.go`, aus der Tabelle in `SPEC-034` und aus den Tabellen des Handbuchs und vergleicht die Mengen (`SPEC-054`) |
| `tools/harness/meldungskatalog-check-gegenprobe.sh` | neu | Fehlformen in Temp-Bäumen je Punkt des Vertrags, wie `tools/harness/kopf-check-gegenprobe.sh`; dazu der Fall mit Makefile und Fragment |
| `harness/mk/meldungskatalog-check.mk` | neu | zwei Ziele, Prüfung und Gegenprobe, beide an `GATE_CHECKS` |
| `harness/mk/gates-schnell.mk` | update | `meldungskatalog-check-gegenprobe` in `GATE_SCHNELL_AUSNAHME` |
| `tools/harness/gates-schnell-gegenprobe.sh` | update | `meldungskatalog-check-gegenprobe` in `LANGSAM`, der Liste der Attrappen; ohne den Eintrag meldet die Deckung der Ausnahmeliste (`SPEC-052` Punkt 4) in jeder Kopie `p4-eintrag` |
| `harness/sensors/meldungskatalog-check.md` | neu | Vertrag-Zeiger, Grenze, Ausgabe und Ausgänge, Sperren |
| `harness/README.md` §Sensors | update | zwei Zeilen mit Bindung |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Slices von welle-v1-abschluss außer
`slice-v1-abschluss-container` liegen in `done/`
(WIP-Limit 1); dieser Slice ist Schritt 43 der Reihenfolge. Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die drei Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-13, vor die Release-Slices gezogen), dann die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-abdeckung-gate` und
`slice-harness-coverage`, dann die Slices von welle-erster-release und welle-v1-nachtrag, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-lint-warten-ohne-frist`,
`slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md). Technisch hängt
dieser Slice an keinem davon; die Code-Tabelle und der Katalog liegen seit
welle-replay-semantik vor. Er steht direkt vor `slice-v1-abschluss-container`, dessen §1
das Gate als vorher geliefert nennt: Es misst die Codes, die die Slices der Welle
eingeführt haben, bevor die Betriebsdokumentation den Katalog schreibt, und die Frage, ob
es diesen Katalog mitprüft (§6, Punkt 8), ist vor dessen Entstehen entschieden. Dass die
Slices der Welle ihre Codes nicht unter dem Gate liefern, ist die Folge der Reihenfolge;
was sie am Katalog vorbei eingeführt haben, berichtigt dieser Slice nach §6 *Bestand*. Vor dem ersten
Code-Commit entscheidet der Architect die Randformen aus §6 im Vertrag und ob eine
ADR nötig ist (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und
  Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness`
  liegt über 900; Basis ist `8d14a918` (HEAD vor dem Anspruch; er enthält Vertrag und Plan, die Move-Commits berühren `tools` und `harness` nicht). Die Schätzung des
  Architect (2026-10-17) in dieser Metrik: Spanne 450 bis 750, Mitte 580, Obergrenze 750
  (Skript 110 bis 150; Gegenprobe 300 bis 450, das Zwei- bis Dreifache des Skripts nach dem
  Faktor aus `AGENTS.md` §3.14; Fragment, Sensor-Datei und Zeilen in `harness/README.md` und
  `gates-schnell.mk` rund 80); die Schätzung liegt unter der Marke, ein Schnitt vor dem Code
  ist nicht geboten. Meldeschwelle: Stehen bei 750 noch Fälle der Gegenprobe aus, meldet der
  Implementer es dem Architect, bevor er weiterschreibt.
- `in-progress` → `open` (blockiert — Carveout?): Die Richtung „umgekehrt“ (§6) ist
  nur als Gleichheit entscheidbar und färbt das Gate am Stand des Starts rot, ohne dass
  der Nutzer eine Stufung oder einen Carveout entschieden hat.

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

**Randformen des neuen Vertrags** ([`AGENTS.md`](../../../../AGENTS.md) §3.12). Der Architect
hat jede am 2026-10-17 im Vertrag
[`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) entschieden; die Nummer
nennt den Punkt. Der Implementer entscheidet keine still, sondern gibt eine nicht genannte
zurück. Stand am 2026-10-17: C (Konstanten) 30 Codes, K (Katalog) 34, H (Handbuch) 30; die
Zahlen „34 gegen 19“ des ersten Plans sind durch die Slices von welle-v1-abschluss überholt.

1. **Richtung „umgekehrt“, nicht erzeugte Codes** — Punkt 4b: K ohne C ist ein Befund, außer
   bei Endung `000`; eine Markierung im Katalog gibt es nicht, und am Stand des Starts ist
   nichts mehr offen.
2. **Rückfälle** — Punkt 4b und 5: Ein nicht erzeugter Rückfall muss nicht im Code stehen; je
   Klasse 1 bis 6 steht `PGR-E<k>000` in K.
3. **Reserviert, zurückgezogen** — Punkt 7: keine Form; der Slice, der zurückzieht, ändert den
   Vertrag.
4. **Fundort im Code** — Punkt 1 (nur die Konstanten `Code…` in `fehler.go`) und 6 (ein
   Literal `"PGR-` in einer anderen Nicht-Test-Datei unter `cmd/` und `internal/` ist ein
   Befund).
5. **Tests und Kommentare** — Grenze: nicht gelesen.
6. **Form** — Punkte 2 und 3: ERE `PGR-[EWI][0-9]{4}` ganz; doppelte Konstante, Katalog- und
   Handbuch-Zeile sind `doppelt`.
7. **Grenze des Katalogs; fehlende Tabelle** — Punkte 1 und 8: nur die erste Zelle der
   Tabellenzeilen im Abschnitt `### SPEC-034 `; fehlende Datei, Abschnitt oder Zeilen sind
   eine Sperre (Exit 2).
8. **Katalog der Betriebsdokumentation** — Punkt 4c und 4d: Das Handbuch ist dritte Quelle und
   führt genau die Codes der Code-Tabelle. Empfehlung des Architect, nicht Entscheidung des
   Nutzers; sie gilt, bis er widerspricht. Die Rückführung in §4 bindet nicht an die Quelle,
   sondern an die Größe: Die Schätzung liegt unter der Marke.
9. **Klasse** — Grenze: Die Spalten werden nicht gelesen.
10. **Form der Meldung und Ausgänge** — Punkte 9 und 8: Zeile, Sortierung, Exit 0, 1 und 2.
11. **Gegenprobe** — Punkt 11 und DoD: je Punkt ein Fall; der Anschluss-Fall mit Makefile und
    Fragment.
12. **Aufrufarten** — Punkt 10: kein Stand und kein Ergebnis außer dem Ausgang; keine
    Sonderbehandlung von `-n`, `-q`, `-t`, `-i`.
13. **Wahl der Kette** — Punkt 11: `meldungskatalog-check` in der schnellen Kette (`SPEC-052`
    Punkt 13), die Gegenprobe in der Ausnahmeliste.

**Risiken:**

- **Bestand.** Am Stand des Starts (2026-10-17) ist das Gate nach `SPEC-054` grün: C ist
  Teilmenge von K, K ohne C sind vier Rückfälle, H gleich C, kein Literal außerhalb von
  `fehler.go`. Ändert sich der Stand bis zum Start, ist ein Rot ein Befund für den
  Architect, kein stilles Rot und keine Stufung ([`AGENTS.md`](../../../../AGENTS.md) §3.6)
  — **Ausgang:** offen bis Closure.
- **Folgelast für welle-v1-abschluss.** Jeder Slice, der einen Code einführt, muss ab
  diesem Slice Katalog und Code-Tabelle zusammen ändern; das Gate macht das sichtbar,
  die Pläne nennen es heute nicht — **Ausgang:** offen bis Closure.
- **Adresse aus `slice-v1-abschluss-container`.** Dieser Slice schließt vor dem
  Container-Slice; der verweist auf ihn als Liefernden, nicht als späteren Folge-Slice.
  Schreibt der Container-Slice den Katalog der Betriebsdokumentation fort, hält dieses
  Gate die Tabellen `### Fehlercodes` und `### Warnungen` des Handbuchs (`SPEC-054` Punkt 1),
  nicht andere Betriebsdokumente — **Ausgang:** offen bis Closure.
- **Lesart ohne Entscheidungsort.** Ein Werkzeug-Vertrag, dessen Lesarten erst im Code
  festgelegt werden, wiederholt `BEO-REPO/harness-lesart-ohne-entscheidungsort` (2×) —
  **Ausgang:** offen bis Closure.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Basis der Marke `8d14a918`; Code-Commit `9b7ff671` (Prüfskript, Gegenprobe, Fragment, Sensor-Datei,
Zeilen in `harness/README.md`, `GATE_SCHNELL_AUSNAHME`, `LANGSAM`). Das Gate startet grün am Stand des
Starts: `bash tools/harness/meldungskatalog-check.sh` endet mit 0 ohne Ausgabe (C 30 Codes, K 34, H 30;
`fehler.go`, `SPEC-034`, Handbuch unverändert).

**Weg der Mutationen.** Je Mutant eine frische Kopie von Skript, Gegenprobe, Fragmenten und `Makefile` in ein Temp-Verzeichnis
(kein Docker, kein Build-Kontext; die mtime-Regel gilt nicht), im Skript, im Fragment `meldungskatalog-check.mk` oder in
`gates-schnell.mk` genau eine Zeichenkette ersetzt (der Lauf bricht ab, wenn sie nicht genau einmal vorkommt), dann die
Gegenprobe aus der Kopie gefahren; die roten Fälle sind die Namen `Fall '…'` auf stderr. Der Läufer steht am Ende dieses
Abschnitts im Wortlaut; die Spalte *Mutation* nennt die Änderung, der Läufer die Zeichenketten. 109 Mutanten: 107 färben mindestens
den gezielten Fall ihrer Zusage rot, 2 sind äquivalent (unten).

**Mutationstabelle** (Zusage · Mutation · roter Fall). Bei Mutanten, die jeden Gutfall treffen (Anzahl in Klammern),
steht der gezielte Fall voran; die Zahl nennt alle roten Fälle des Laufs. Mutanten, die einen Aufbau aller Fälle
brechen (`k-backtick-pflicht`, `h-…-nicht-gelesen`, `4a-weg`, `9-stdout`), färben zwangsläufig auch jeden Fall rot, dessen
Baum die Quelle benutzt; das ist ihre Invariante, kein Mangel des gezielten Falls.

| Zusage (Punkt von `SPEC-054`) | Mutation | roter Fall |
|---|---|---|
| 1 C: beliebiger Einzug | Einzug nur Tabulator | `nr1-c-einzug-leerzeichen` |
| 1 C: beliebiger Einzug (auch keiner) | Einzug Pflicht | `nr1-c-ohne-einzug` |
| 1 C: Leerraum vor `=` | nur Leerzeichen vor `=` | `nr1-c-leerraum-tab` |
| 1 C: Leerraum nach `=` | nur Leerzeichen nach `=` | `nr1-c-leerraum-tab` |
| 1 C: Name `Code` allein ist Konstante | Name braucht ein Zeichen nach `Code` | `nr1-c-nur-code` |
| 1 C: Wortzeichen im Namen | Unterstrich kein Wortzeichen | `nr1-c-wortzeichen` |
| 1 C: Zeile beginnt mit `Code` (Kommentar, anderer Anfang zählt nicht) | Zeilenanfang nicht verankert | `nr1-c-anderer-name` `nr1-c-kommentar` |
| 1 C: Name beginnt mit `Code` | Namenspräfix nicht verlangt | `nr1-c-anderer-name` `nr6-fehler-go-anderer-name` |
| 1 C: Leerraum, dann `=` | Typname vor `=` zugelassen | `nr1-c-typisiert` |
| 1 C: Zeichenkette in doppelten Anführungszeichen | Backtick-Wert zugelassen | `nr1-c-backtick` |
| 1: erste Zelle in Backticks; Kopf- und Trennzeile keine Zeile | Backticks nicht verlangt | `nr1-k-keine-zeile` (68 Fälle) |
| 1 K: nur der Abschnitt | Lesen beginnt vor der Überschrift | `nr1-k-davor` |
| 1 H: nur die beiden Abschnitte | Lesen beginnt vor der Überschrift | `nr1-h-davor` |
| 1 K: Abschnitt endet an der nächsten `#`-Zeile | Ende nicht erkannt | `nr1-k-danach` |
| 1 H: Abschnitt endet an der nächsten `#`-Zeile | Ende nicht erkannt | `nr1-h-danach` `nr1-h-zwischen` |
| 1 H: `### Warnungen` gehört zu H | Abschnitt nicht gelesen | `nr4d-warnung` `nr8-h-nur-warnungen` (63 Fälle) |
| 1 H: `### Fehlercodes` gehört zu H | Abschnitt nicht gelesen | `nr4c-fehler` `nr8-h-nur-fehlercodes` (64 Fälle) |
| 8: H zählt beide Abschnitte zusammen | Sperre, wenn `### Fehlercodes` leer | `nr8-h-nur-warnungen` (66 Fälle) |
| 4b: `…000` in K ohne C zugelassen | Ausnahme entfernt | `nr4b-rueckfall-w` `nr4b-endung-00` (65 Fälle) |
| 4b: Ausnahme gilt für jede Schwere | nur `E…000` | `nr4b-rueckfall-w` |
| 4b: Ausnahme gilt nur für Endung `000` | Endung `00` genügt | `nr4b-endung-00` |
| 2: Form ganz (Ende) | Endanker entfernt | `nr2-form-zu-lang` `nr2-form-zwei-codes` |
| 2: Form ganz (Anfang) | Anfangsanker entfernt | `nr2-form-vorsatz` `nr2-form-zwei-codes` |
| 2: Schwere `I` gehört zur Form | `I` nicht zugelassen | `nr2-schwere-i` |
| 2: Form-Code nimmt an 3 bis 5 nicht teil | Form-Code läuft weiter | `nr2-form-c` `nr9-zeile-numerisch` (8 Fälle) |
| 3: je Code und Quelle einmal | jede Wiederholung gemeldet | `nr3-doppelt-dreifach` |
| 3: an der zweiten Fundstelle | Zeile der ersten | `nr3-doppelt-c` `nr3-doppelt-dreifach` `nr3-doppelt-h` `nr3-doppelt-h-abschnitte` `nr3-doppelt-k` `nr9-befund-vor-zeile` |
| 3: der Code zählt einmal | Wiederholung trägt in die Liste ein | `nr9-befund-vor-zeile` |
| 3: in C | nicht erkannt | `nr3-doppelt-c` `nr3-doppelt-dreifach` `nr9-befund-vor-zeile` |
| 3: in K | nicht erkannt | `nr3-doppelt-k` |
| 3: in H (auch über beide Abschnitte) | nicht erkannt | `nr3-doppelt-h` `nr3-doppelt-h-abschnitte` |
| 4a | Regel entfernt | `nr4a-nicht-im-katalog` (14 Fälle) |
| 4b | Regel entfernt | `nr2-schwere-i` `nr4b-endung-00` `nr4b-nicht-in-c` `nr9-bytefolge` `nr9-sortiert` |
| 4c | Regel entfernt | `nr4c-fehler` `nr4c-warnung` (13 Fälle) |
| 4d | Regel entfernt | `nr4d-in-k` `nr4d-nicht-in-k` `nr4d-warnung` `nr9-sortiert` |
| 4d: auch wenn K den Code führt | nur wenn K ihn nicht führt | `nr4d-in-k` |
| 5: Klasse 6 | Schleife bis 5 | `nr5-klasse-6` |
| 5: Klasse 1 | Schleife ab 2 | `nr5-klasse-1` |
| 5: Klasse 2 | Klasse übersprungen | `nr5-klasse-2` |
| 5: Klasse 3 | Klasse übersprungen | `nr5-klasse-3` |
| 5: Klasse 4 | Klasse übersprungen | `nr5-klasse-4` |
| 5: Klasse 5 | Klasse übersprungen | `nr5-klasse-5` |
| 5: Rückfall-Befund ohne Zeile | Zeile `0` | `nr5-klasse-1` `nr5-klasse-2` `nr5-klasse-3` `nr5-klasse-4` `nr5-klasse-5` `nr5-klasse-6` |
| 6: `cmd/` | nicht durchsucht | `nr6-cmd` `nr6-kommentar-mit-anfuehrung` `nr9-sortiert` |
| 6: `internal/` | nicht durchsucht | `nr6-anderes-fehler-go` `nr6-internal` |
| 6: Test-Dateien ausgenommen | `_test.go` gelesen | `nr6-test` |
| 6/Grenze: `test/` nicht gelesen | `test/` durchsucht | `nr6-test-verzeichnis` (76 Fälle) |
| 6: nur Go-Dateien | jede Datei | `nr6-nicht-go` |
| 6: nur die Code-Tabelle ausgenommen | jede `fehler.go` ausgenommen | `nr6-anderes-fehler-go` |
| 6: Code-Tabelle ausgenommen | Code-Tabelle durchsucht | `nr6-fehler-go-anderer-name` (65 Fälle) |
| 6: Anführungszeichen unmittelbar davor | `PGR-` ohne Anführungszeichen | `nr6-backtick` `nr6-kommentar-ohne-anfuehrung` (7 Fälle) |
| 6: Code aus Wortzeichen | Code bis Leerzeichen | `nr6-anderes-fehler-go` `nr6-cmd` `nr6-internal` `nr6-kommentar-mit-anfuehrung` `nr9-sortiert` |
| 8: Reihenfolge C vor K | getauscht | `nr8-reihenfolge-c-vor-k` |
| 8: Reihenfolge K vor H | getauscht | `nr8-reihenfolge-k-vor-h` |
| 8: Dateien vor Verzeichnissen | getauscht | `nr8-reihenfolge-datei-vor-verzeichnis` |
| 8: `<pfad> fehlt` für C | Prüfung für C entfernt | `nr8-c-fehlt` `nr8-reihenfolge-c-vor-k` |
| 8: `nicht lesbar` für C | Prüfung für C entfernt | `nr8-unlesbar-F` |
| 8: `<pfad> fehlt` für K | Prüfung für K entfernt | `nr8-k-fehlt` `nr8-reihenfolge-k-vor-h` |
| 8: `nicht lesbar` für K | Prüfung für K entfernt | `nr8-unlesbar-S` |
| 8: `<pfad> fehlt` für H | Prüfung für H entfernt | `nr8-h-fehlt` `nr8-reihenfolge-datei-vor-verzeichnis` |
| 8: `nicht lesbar` für H | Prüfung für H entfernt | `nr8-unlesbar-B` |
| 8: `cmd/ fehlt` | Prüfung entfernt | `nr8-cmd-fehlt` `nr8-reihenfolge-verzeichnis-vor-abschnitt` |
| 8: Verzeichnisse vor Abschnitten | Prüfung entfernt (awk-Sperre gilt zuerst) | `nr8-cmd-fehlt` `nr8-reihenfolge-verzeichnis-vor-abschnitt` |
| 8: Abschnitt K vor `### Fehlercodes` | getauscht | `nr8-reihenfolge-abschnitt-k-vor-f` |
| 8: `### Fehlercodes` vor `### Warnungen` | getauscht | `nr8-reihenfolge-abschnitt-f-vor-w` |
| 8: Abschnitte vor `ohne Code` | getauscht | `nr8-reihenfolge-abschnitt-vor-quelle` |
| 8: `ohne Code` C vor K | getauscht | `nr8-reihenfolge-quelle-c-vor-k` |
| 8: `ohne Codes` K vor H | getauscht | `nr8-reihenfolge-quelle-k-vor-h` |
| 8: Überschrift `### SPEC-034 ` mit Leerzeichen | Präfix ohne Leerzeichen | `nr8-abschnitt-k` `nr8-reihenfolge-abschnitt-k-vor-f` `nr8-reihenfolge-abschnitt-vor-quelle` `nr8-sperre-vor-befund` |
| 8: Zeile `### Fehlercodes` ganz | Präfix genügt | `nr8-abschnitt-fehlercodes` |
| 8: Zeile `### Warnungen` ganz | enthalten genügt | `nr8-abschnitt-warnungen` |
| 8: keine Befund-Zeile bei Sperre | Befunde vor der Sperre ausgegeben | `nr8-sperre-vor-befund` |
| 9: Ausgang 2 (Sperre in der Shell) | Exit 1 | `nr8-c-fehlt` `nr8-cmd-fehlt` (11 Fälle) |
| 9: Ausgang 2 (Sperre im awk) | awk endet mit 1 | keiner (äquivalent, siehe unten) |
| 9: Ausgang 2 (Sperre aus awk) | Exit 1 | `nr8-abschnitt-k` `nr8-c-ohne-code` (12 Fälle) |
| 9: sortiert nach Pfad | Schlüssel entfernt | `nr5-klasse-1` `nr9-sortiert` |
| 9: sortiert nach Code | Schlüssel entfernt | `nr9-bytefolge` `nr9-code-vor-zeile` |
| 9: sortiert nach Befund | Schlüssel entfernt | `nr9-befund-vor-zeile` |
| 9: Zeile als Zahl | lexikalisch | `nr9-zeile-numerisch` |
| 9: Bytefolge | `LC_ALL=C` entfernt | `nr9-bytefolge` |
| 9: stdout leer | Ausgabe auf stdout | `nr2-form-c` `nr3-doppelt-c` (41 Fälle) |
| 9: Ausgang 1 bei Befund | Exit 0 | `nr2-form-c` `nr3-doppelt-c` (42 Fälle) |
| 9: `<pfad>[:<zeile>]`, Zeile fehlt bei Rückfall | Zeile immer ausgegeben | `nr5-klasse-1` `nr5-klasse-2` `nr5-klasse-3` `nr5-klasse-4` `nr5-klasse-5` `nr5-klasse-6` |
| 9: Zeilenform | Präfix geändert | `nr5-klasse-1` `nr5-klasse-2` `nr5-klasse-3` `nr5-klasse-4` `nr5-klasse-5` `nr5-klasse-6` |
| 2: Befund `Form` | Text geändert | `nr2-form-c` `nr2-form-k` (8 Fälle) |
| 3: Befund `doppelt` | Text geändert | `nr3-doppelt-c` `nr3-doppelt-dreifach` `nr3-doppelt-h` `nr3-doppelt-h-abschnitte` `nr3-doppelt-k` `nr9-befund-vor-zeile` |
| 4a: Befund `nicht im Katalog` | Text geändert | `nr4a-nicht-im-katalog` (14 Fälle) |
| 4c: Befund `nicht im Handbuch` | Text geändert | `nr4c-fehler` `nr4c-warnung` (13 Fälle) |
| 5: Befund `Rückfall fehlt` | Text geändert | `nr5-klasse-1` `nr5-klasse-2` `nr5-klasse-3` `nr5-klasse-4` `nr5-klasse-5` `nr5-klasse-6` |
| 6: Befund `Code außerhalb der Code-Tabelle` | Text geändert | `nr6-anderes-fehler-go` `nr6-cmd` `nr6-internal` `nr6-kommentar-mit-anfuehrung` `nr9-sortiert` |
| 9: Befund `nicht in der Code-Tabelle` (K) | Text geändert | `nr2-schwere-i` `nr4b-endung-00` `nr4b-nicht-in-c` `nr9-bytefolge` `nr9-sortiert` |
| 9: Befund `nicht in der Code-Tabelle` (H) | Text geändert | `nr4d-in-k` `nr4d-nicht-in-k` `nr4d-warnung` `nr9-sortiert` |
| 10: `MAKEFLAGS=-n` ändert die Aussage nicht | Skript endet 0 bei `-n` | `nr10-makeflags-n` |
| 10: `MAKEFLAGS=-q` ändert die Aussage nicht | Skript endet 0 bei `-q` | `nr10-makeflags-q` |
| 10: `MAKEFLAGS=-t` ändert die Aussage nicht | Skript endet 0 bei `-t` | `nr10-makeflags-t` |
| 10: `MAKEFLAGS=-i` ändert die Aussage nicht | Skript endet 0 bei `-i` | `nr10-makeflags-i` |
| 10: kein Stand, grüner Lauf | legt Datei an | `nr10-schreibt-nichts-gruen` |
| 10: kein Stand, roter Lauf | legt Datei an | `nr10-schreibt-nichts-rot` |
| 10: ohne Argument das aktuelle Verzeichnis | anderer Vorgabewert | `nr10-wurzel-aktuell` `nr11-ziel-gruen` |
| 11: Prüfung an `GATE_CHECKS` | entfernt | `nr11-gate-checks` |
| 11: Gegenprobe an `GATE_CHECKS` | entfernt | `nr11-gate-checks` |
| 11: Exit ungleich 0 bei rotem Skript (Prüfung) | `|| true` | `nr11-scharf` |
| 11: Exit ungleich 0 bei rotem Skript (Gegenprobe) | `|| true` | `nr11-scharf` |
| 11: Prüfung ruft das Skript | `true` | `nr11-scharf` |
| 11: Prüfung ruft das Skript | falscher Pfad | `nr11-ziel-gruen` |
| 10: `make -t` hinterlässt nichts | `.PHONY` entfernt | `nr11-make-i` `nr11-make-t` `nr11-makeflags-umgebung` |
| 11: Gegenprobe in der Ausnahmeliste | entfernt | `nr11-ausnahme-gegenprobe` |
| 11: Prüfung in der schnellen Kette | Prüfung in der Ausnahmeliste | `nr11-schnelle-kette` |
| 8: `internal/ fehlt` | Prüfung entfernt | keiner (äquivalent, siehe unten) |

**Äquivalente Mutanten.**

- `8-exit-awk-sperre` (awk endet bei einer Sperre mit 1 statt 2): Invariante: das Skript bildet jeden Ausgang ungleich 0 von
  `awk` auf Exit 2 ab; Ort: `if ! befunde="$(awk …)"; then exit 2` in `tools/harness/meldungskatalog-check.sh`. Der Mutant
  an diesem Ort (`8-awk-fehler-exit-2`, Exit 1) ist rot; der Ausgang des awk-Prozesses ist nur ein Zwischenwert.
- `8-ohne-internal` (die Sperre `internal/ fehlt` entfällt): Invariante: die Code-Tabelle liegt unter `internal/`, ihr Fehlen
  gilt in der Reihenfolge von Punkt 8 zuerst; fehlt `internal/`, endet der Lauf vorher mit `<C> fehlt`. Die Zeile
  `internal/ fehlt` ist nicht erreichbar; ein Fall dafür ist nicht herstellbar.

**Messung der Größe** (Metrik §4: Summe Einfügungen und Löschungen von
`git diff --numstat 8d14a918..<letzter Commit> -- tools harness`).

Summe 488 (Marke 900, Meldeschwelle 750, nicht erreicht), gemessen am Code-Commit `9b7ff671` mit `git diff --numstat 8d14a918..9b7ff671 -- tools harness`;
das Plan-Commit danach berührt `docs/` und `tools`/`harness` nicht. Aufschlüsselung (Einfügungen + Löschungen): `tools/harness/meldungskatalog-check.sh` 132;
`tools/harness/meldungskatalog-check-gegenprobe.sh` 261 (das Zweifache des Skripts; Schätzung des Architect 300 bis 450 für die Gegenprobe);
`harness/sensors/meldungskatalog-check.md` 71; `harness/mk/meldungskatalog-check.mk` 18; `harness/README.md` 2; `harness/mk/gates-schnell.mk` 2;
`tools/harness/gates-schnell-gegenprobe.sh` 2. Die Schätzung 450 bis 750 (Mitte 580) trifft mit 488 den unteren Teil der Spanne.

**Läufe.**

- `bash tools/harness/meldungskatalog-check.sh` am Stand des Starts: Exit 0, keine Ausgabe.
- `make meldungskatalog-check-gegenprobe`: grün (Ausgabe `gruen — Fehlformen je Punkt von SPEC-054 abgelehnt, Gutfälle angenommen`); `make gates-schnell-gegenprobe`: grün.
- Mutanten: 109 gefahren, 107 rot, 2 äquivalent (oben); die Mutanten nach der Korrektur der Fälle `nr1-k-davor`, `nr1-h-davor` und `nr8-sperre-vor-befund`
  (anfangs grün: der Anfangszustand im Abschnitt und die Befunde vor der Sperre hatten keinen Fall, der sie erreichte) sind erneut gefahren.
- `make gates-schnell`: Exit 0. `make kopf-check`: Exit 0, keine Ausgabe. `make docs-check`: `938 Datei(en) geprüft, 0 Befund(e)`.
- `make lint`: nicht gesondert gefahren; kein Go-Code berührt (`cmd/`, `internal/`, `test/`, `.golangci.yml` unverändert), die Stufe lief in `make gates-schnell` und `make gates` aus dem Zwischenspeicher grün.
- `make gates` (ein Lauf, im Hintergrund über das Werkzeug gestartet, Stand: HEAD `9b7ff671` mit dem Plan, wie er vor dieser Zeile stand): Exit 0, darin `meldungskatalog-check-gegenprobe: gruen` und `stop-hook-gegenprobe: gruen — 139 Fälle`.

**Rückgabe an den Architect: Randform · Frage.** Der Vertrag nennt die folgenden Punkte nicht oder nur durch einen Platzhalter;
der Code trägt jeweils die in der Spalte *vorläufig* genannte Form, und die Gegenprobe zeigt sie. Der Slice bleibt in
`in-progress/`, bis der Architect entschieden hat; jede Entscheidung ändert eine Zeile und einen Fall.

| Randform | Frage | vorläufig im Code (Fall) |
|---|---|---|
| Name einer Konstante in C | Ist ein Name genau `Code` eine Konstante („`Code` und Wortzeichen“)? | ja (`nr1-c-nur-code`) |
| Konstante mit Typangabe | Ist `CodeT Typ = "…"` eine Konstante (Punkt 1: „Leerraum, `=`, Leerraum“)? | nein (`nr1-c-typisiert`) |
| Wortlaut der Sperre | Was setzt `<überschrift>` und `<verzeichnis>` ein? | `Abschnitt ### SPEC-034 fehlt`, `Abschnitt ### Fehlercodes fehlt`, `Abschnitt ### Warnungen fehlt`, `cmd/ fehlt` (`nr8-abschnitt-*`, `nr8-cmd-fehlt`) |
| Reihenfolge innerhalb einer Sperrklasse | Punkt 8 ordnet die Klassen, nicht die Glieder: Dateien, Abschnitte, Quellen | C, K, H; K, Fehlercodes, Warnungen; C, K, H (`nr8-reihenfolge-*`) |
| Überschrift der Handbuch-Abschnitte | Ist `### Fehlercodes` die ganze Zeile oder ein Präfix? Punkt 1 schreibt bei K das Leerzeichen mit, bei H nicht | ganze Zeile (`nr8-abschnitt-fehlercodes`, `nr8-abschnitt-warnungen`); K: Präfix mit Leerzeichen (`nr8-abschnitt-k`) |
| Pfad ist ein Verzeichnis, Wurzel fehlt | „fehlt oder ist nicht lesbar“ | `<pfad> fehlt` für alles, was keine reguläre Datei ist; fehlende Wurzel: `<C> fehlt` (nicht durch einen Fall belegt) |
| Nicht lesbares Verzeichnis oder nicht lesbare Datei im Baum für Punkt 6 | Punkt 8 nennt nur die drei Dateien. Gemessen: ein Verzeichnis mit Rechten 000 unter `cmd/` endet mit Exit 0 und einer Meldung von `find` (die Datei darin ist ungeprüft, grün); eine Datei mit Rechten 000 endet mit Exit 2 und der Meldung von `awk`. | beides ungeregelt, kein Fall; das stille Grün entspricht `BEO-REPO/gate-uebergeht-ablage-eintrag-still` |
| Mehrere `"PGR-` in einer Zeile; Tabellenzeile mit Einzug oder ohne schließendes `|` | Punkt 6 und Punkt 1 sagen „eine Zeile“ bzw. „erste Zelle“ | eine Zeile ein Befund, mit dem ersten Literal; Tabellenzeile beginnt mit `|` in Spalte 1 und braucht kein schließendes `|` (kein Fall) |

**Abweichung von der Vorgabe des Architect.** Die Vorgabe „LANGSAM bleibt unverändert“ in `tools/harness/gates-schnell-gegenprobe.sh`
trägt nicht: Die Attrappen dieser Gegenprobe kopieren `harness/mk/gates-schnell.mk` des Arbeitsbaums, und jeder Eintrag der Ausnahmeliste, den
`LANGSAM` nicht führt, ist in jeder Kopie ein `p4-eintrag`. Gemessen: mit dem neuen Eintrag in `GATE_SCHNELL_AUSNAHME` und ohne ihn
in `LANGSAM` endet `make gates-schnell-gegenprobe` rot und meldet Mutanten wie `m-hash-vor-dem-lauf` neben ihrem Punkt auch mit
`p4-eintrag`; mit dem Eintrag in beiden grün. Geändert ist ein Wort in einer Zeile (Plan §3, DoD). Die Gegenprobe der schnellen Kette
selbst, ihre Fälle und ihr Vertrag sind unberührt.

**Die Gegenprobe liest den Arbeitsbaum nicht** (Regel des Schreibers, `SPEC-052` Punkt 13): Im Lauf aus einer Kopie, die nur
`Makefile`, `harness/mk/` und die beiden Skripte trägt (kein `spec/`, `docs/`, `cmd/`, `internal/`), endet die Gegenprobe grün.

**Beobachtungen für die Closure** (Register nicht fortgeschrieben, das ist Sache der Closure).

- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (2×): mit diesem Slice zum dritten Mal berührt; die Tabelle oben ist die Liste, die der
  Vertrag nicht entschied.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still`: ein nicht lesbares Verzeichnis im Baum für Punkt 6 bleibt still (Tabelle oben).
- Eine Vorgabe des Architect über eine Datei, die der Slice nicht berührt, ohne Messung (`LANGSAM`, oben): Der Lauf von
  `make gates-schnell-gegenprobe` vor der Vorgabe hätte sie widerlegt.
- Mutanten, die einen Aufbau aller Gutfälle brechen (Kopf- und Trennzeile von K und H, die Ausnahme `…000`, `4a`), sind nur über Fälle zu
  trennen, die den Aufbau mitbenutzen; die gezielten Fälle stehen in der Tabelle voran.

**Läufer der Mutationen** (Python 3, aufgerufen mit `python3 -I <datei>` aus der Wurzel; ohne Argument alle Mutanten, mit Argumenten die
Mutanten, deren Name ein Argument enthält):

```text
import os, re, shutil, subprocess, sys, tempfile
R = "/Development/pgwire-recorder"
SK = "tools/harness/meldungskatalog-check.sh"
GP = "tools/harness/meldungskatalog-check-gegenprobe.sh"
MK = "harness/mk/meldungskatalog-check.mk"
GS = "harness/mk/gates-schnell.mk"
FILES = {"sk": SK, "mk": MK, "gs": GS}
M = []
def m(name, key, old, new, n=1):
    M.append((name, key, old, new, n))
CR = '/^[ \\t]*Code[A-Za-z0-9_]*[ \\t]+=[ \\t]+"[^"]*"/'
m("c-einzug-nur-tab", "sk", '/^[ \\t]*Code', '/^[\\t]*Code')
m("c-einzug-pflicht", "sk", '/^[ \\t]*Code', '/^[ \\t][ \\t]*Code')
m("c-leerraum-vor-gleich", "sk", 'Code[A-Za-z0-9_]*[ \\t]+=', 'Code[A-Za-z0-9_]* +=')
m("c-leerraum-nach-gleich", "sk", '=[ \\t]+"[^"]*"/', '= +"[^"]*"/')
m("c-nur-code", "sk", 'Code[A-Za-z0-9_]*[ \\t]+=', 'Code[A-Za-z0-9_]+[ \\t]+=')
m("c-wortzeichen-unterstrich", "sk", 'Code[A-Za-z0-9_]*[ \\t]+=', 'Code[A-Za-z0-9]*[ \\t]+=')
m("c-anker-weg", "sk", '/^[ \\t]*Code', '/[ \\t]*Code')
m("c-name-code-weg", "sk", 'Code[A-Za-z0-9_]*[ \\t]+=', '[A-Za-z0-9_]*[ \\t]+=')
m("c-typisiert", "sk", 'Code[A-Za-z0-9_]*[ \\t]+=', 'Code[A-Za-z0-9_]*[ \\t]+([A-Za-z0-9_]+[ \\t]+)?=')
m("c-backtick", "sk", '=[ \\t]+"[^"]*"/', '=[ \\t]+["`][^"`]*["`]/')
m("k-backtick-pflicht", "sk", ' || substr(s, 1, 1) != "`" || substr(s, length(s), 1) != "`") return 0', ') return 0')
m("k-davor", "sk", 'imabschnitt = 0\n  }\n  art == "C"', 'imabschnitt = (art == "K")\n  }\n  art == "C"')
m("h-davor", "sk", 'imabschnitt = 0\n  }\n  art == "C"', 'imabschnitt = (art == "H")\n  }\n  art == "C"')
m("k-danach", "sk", '    else if (substr($0, 1, 1) == "#") imabschnitt = 0\n    else if (imabschnitt && tabellenzeile($0)) eintrag("K"', '    else if (imabschnitt && tabellenzeile($0)) eintrag("K"')
m("h-ende", "sk", '    else if (substr($0, 1, 1) == "#") imabschnitt = 0\n    else if (imabschnitt && tabellenzeile($0)) eintrag("H"', '    else if (imabschnitt && tabellenzeile($0)) eintrag("H"')
m("h-warnungen-nicht-gelesen", "sk", '{ imabschnitt = 1; hatw = 1 }', '{ imabschnitt = 0; hatw = 1 }')
m("h-fehlercodes-nicht-gelesen", "sk", '{ imabschnitt = 1; hatf = 1 }', '{ imabschnitt = 0; hatf = 1 }')
m("h-ohne-codes-fehlercodes-leer", "sk", 'if (!anzahl["H"]) sperre', 'if (!anzahl["HF"]) sperre')
m("rueckfall-ausnahme-weg", "sk", ' && code !~ /000$/', '')
m("rueckfall-ausnahme-nur-e", "sk", 'code !~ /000$/', 'code !~ /E[0-9]000$/')
m("rueckfall-ausnahme-00", "sk", 'code !~ /000$/', 'code !~ /00$/')
m("form-dollar", "sk", '[0-9][0-9][0-9][0-9]$/) { ausgabe', '[0-9][0-9][0-9][0-9]/) { ausgabe')
m("form-caret", "sk", '/^PGR-[EWI]', '/PGR-[EWI]')
m("form-ohne-i", "sk", '[EWI]', '[EW]')
m("form-ohne-return", "sk", '"Form", zeile); return }', '"Form", zeile) }')
m("doppelt-ohne-gemeldet", "sk", 'if (!((q, code) in gemeldet)) { gemeldet[q, code] = 1; ausgabe(pfad, code, "doppelt", zeile) }', 'ausgabe(pfad, code, "doppelt", zeile)')
m("doppelt-erste-zeile", "sk", 'ausgabe(pfad, code, "doppelt", zeile)', 'ausgabe(pfad, code, "doppelt", erste[q, code])')
m("doppelt-ohne-return", "sk", '"doppelt", zeile) }\n      return\n    }', '"doppelt", zeile) }\n    }')
m("doppelt-nicht-c", "sk", 'if ((q, code) in erste) {', 'if (q != "C" && (q, code) in erste) {')
m("doppelt-nicht-k", "sk", 'if ((q, code) in erste) {', 'if (q != "K" && (q, code) in erste) {')
m("doppelt-nicht-h", "sk", 'if ((q, code) in erste) {', 'if (q != "H" && (q, code) in erste) {')
m("4a-weg", "sk", 'if (!(("K", code) in erste)) ausgabe(c, code, "nicht im Katalog", erste["C", code])', '')
m("4b-weg", "sk", 'if (!(("C", code) in erste) && code !~ /000$/) ausgabe(k', 'if (0) ausgabe(k')
m("4c-weg", "sk", 'if (!(("H", code) in erste)) ausgabe(c, code, "nicht im Handbuch", erste["C", code])', '')
m("4d-weg", "sk", 'if (!(("C", code) in erste)) ausgabe(h', 'if (0) ausgabe(h')
m("4d-nicht-wenn-in-k", "sk", 'if (!(("C", code) in erste)) ausgabe(h', 'if (!(("C", code) in erste) && !(("K", code) in erste)) ausgabe(h')
m("5-bis-5", "sk", 'i <= 6; i++) {\n      code = "PGR-E"', 'i <= 5; i++) {\n      code = "PGR-E"')
m("5-ab-2", "sk", 'for (i = 1; i <= 6; i++)', 'for (i = 2; i <= 6; i++)')
for k in (2,3,4,5):
    m("5-ohne-klasse-%d"%k, "sk", 'for (i = 1; i <= 6; i++) {\n      code = "PGR-E"', 'for (i = 1; i <= 6; i++) {\n      if (i == %d) continue\n      code = "PGR-E"'%k)
m("5-zeile-null", "sk", '"Rückfall fehlt", "")', '"Rückfall fehlt", 0)')
m("6-ohne-cmd", "sk", 'find cmd internal', 'find internal')
m("6-ohne-internal", "sk", 'find cmd internal', 'find cmd')
m("6-mit-test-dateien", "sk", " ! -name '*_test.go'", '')
m("6-mit-test-verzeichnis", "sk", 'find cmd internal', 'find cmd internal test')
m("6-nicht-nur-go", "sk", "-name '*.go'", "-name '*'")
m("6-fehler-go-ausgenommen-nach-name", "sk", '! -path "$C"', '! -name fehler.go')
m("6-fehler-go-nicht-ausgenommen", "sk", ' ! -path "$C"', '')
m("6-ohne-anfuehrungszeichen", "sk", 'index($0, "\\"PGR-")', 'index($0, "PGR-")')
m("6-code-bis-leerzeichen", "sk", 'match(rest, /^[A-Za-z0-9_]*/)', 'match(rest, /^[^ ]*/)')
m("8-datei-s-vor-c", "sk", 'for datei in "$C" "$K" "$H"', 'for datei in "$K" "$C" "$H"')
m("8-datei-h-vor-k", "sk", 'for datei in "$C" "$K" "$H"', 'for datei in "$C" "$H" "$K"')
m("8-verzeichnis-vor-datei", "sk", 'for datei in "$C" "$K" "$H"; do', '[ -d cmd ] || sperre "cmd/ fehlt"\nfor datei in "$C" "$K" "$H"; do')
for q in ("C","K","H"):
    m("8-fehlt-nicht-%s"%q, "sk", '[ -f "$datei" ] || sperre', '[ -f "$datei" ] || [ "$datei" = "$%s" ] || sperre'%q)
    m("8-lesbar-nicht-%s"%q, "sk", '[ -r "$datei" ] || sperre', '[ -r "$datei" ] || [ "$datei" = "$%s" ] || sperre'%q)
m("8-ohne-internal", "sk", "for verzeichnis in cmd internal", "for verzeichnis in cmd")
m("8-ohne-cmd", "sk", 'for verzeichnis in cmd internal', 'for verzeichnis in internal')
m("8-verzeichnis-nach-abschnitt", "sk", '[ -d "$verzeichnis" ] || sperre "$verzeichnis/ fehlt"', 'true')
m("8-abschnitt-k-f-getauscht", "sk", '    if (!hatk) sperre("Abschnitt ### SPEC-034 fehlt")\n    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")\n', '    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")\n    if (!hatk) sperre("Abschnitt ### SPEC-034 fehlt")\n')
m("8-abschnitt-f-w-getauscht", "sk", '    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")\n    if (!hatw) sperre("Abschnitt ### Warnungen fehlt")\n', '    if (!hatw) sperre("Abschnitt ### Warnungen fehlt")\n    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")\n')
m("8-quelle-vor-abschnitt", "sk", '    if (!hatk) sperre("Abschnitt ### SPEC-034 fehlt")\n    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")\n    if (!hatw) sperre("Abschnitt ### Warnungen fehlt")\n    if (!anzahl["C"]) sperre(c " ohne Code")\n', '    if (!anzahl["C"]) sperre(c " ohne Code")\n    if (!hatk) sperre("Abschnitt ### SPEC-034 fehlt")\n    if (!hatf) sperre("Abschnitt ### Fehlercodes fehlt")\n    if (!hatw) sperre("Abschnitt ### Warnungen fehlt")\n')
m("8-quelle-c-k-getauscht", "sk", '    if (!anzahl["C"]) sperre(c " ohne Code")\n    if (!anzahl["K"]) sperre(k " ohne Codes")\n', '    if (!anzahl["K"]) sperre(k " ohne Codes")\n    if (!anzahl["C"]) sperre(c " ohne Code")\n')
m("8-quelle-k-h-getauscht", "sk", '    if (!anzahl["K"]) sperre(k " ohne Codes")\n    if (!anzahl["H"]) sperre(h " ohne Codes")\n', '    if (!anzahl["H"]) sperre(h " ohne Codes")\n    if (!anzahl["K"]) sperre(k " ohne Codes")\n')
m("8-abschnitt-k-praefix", "sk", 'index($0, "### SPEC-034 ") == 1', 'index($0, "### SPEC-034") == 1')
m("8-abschnitt-fehlercodes-praefix", "sk", '$0 == "### Fehlercodes"', 'index($0, "### Fehlercodes") == 1')
m("8-abschnitt-warnungen-enthaelt", "sk", '$0 == "### Warnungen"', 'index($0, "### Warnungen") > 0')
m("8-befunde-vor-sperre", "sk", "if ! befunde=", "if false; then :; fi\nif ! befunde=") if False else None
m("8-sperre-gibt-befunde-aus", "sk", "then\n  exit 2\nfi", "then\n  printf '%s\\n' \"$befunde\" >&2\n  exit 2\nfi")
m("8-exit-shell-sperre", "sk", 'sperre() { echo "meldungskatalog-check: $*" >&2; exit 2; }', 'sperre() { echo "meldungskatalog-check: $*" >&2; exit 1; }')
m("8-exit-awk-sperre", "sk", 'print "meldungskatalog-check: " grund > "/dev/stderr"; exit 2', 'print "meldungskatalog-check: " grund > "/dev/stderr"; exit 1')
m("8-awk-fehler-exit-2", "sk", "then\n  exit 2\nfi", "then\n  exit 1\nfi")
m("9-sort-ohne-pfad", "sk", "-k1,1 -k2,2", "-k2,2")
m("9-sort-ohne-code", "sk", "-k2,2 -k3,3", "-k3,3")
m("9-sort-ohne-befund", "sk", "-k3,3 -k4,4n", "-k4,4n")
m("9-sort-zeile-lexikalisch", "sk", "-k4,4n", "-k4,4")
m("9-locale", "sk", "export LC_ALL=C\n", "")
m("9-stdout", "sk", "    }' >&2\nexit 1", "    }'\nexit 1")
m("9-exit-befund", "sk", "    }' >&2\nexit 1", "    }' >&2\nexit 0")
m("9-ohne-zeile-format", "sk", 'if ($4 == "") printf', 'if (0) printf')
m("9-praefix", "sk", 'printf "meldungskatalog-check: %s: %s: %s\\n", $1, $2, $3', 'printf "meldungskatalog-pruefung: %s: %s: %s\\n", $1, $2, $3')
for t in ("Form","doppelt","nicht im Katalog","nicht im Handbuch","Rückfall fehlt","Code außerhalb der Code-Tabelle"):
    m("9-text-"+t, "sk", '"%s"'%t, '"%sX"'%t)
m("9-text-nicht-in-ct-k", "sk", 'ausgabe(k, code, "nicht in der Code-Tabelle"', 'ausgabe(k, code, "nicht in der CodeTabelle"')
m("9-text-nicht-in-ct-h", "sk", 'ausgabe(h, code, "nicht in der Code-Tabelle"', 'ausgabe(h, code, "nicht in der CodeTabelle"')
for f in "nqti":
    m("10-makeflags-"+f, "sk", "export LC_ALL=C\n", 'export LC_ALL=C\ncase "${MAKEFLAGS:-}" in *%s*) exit 0;; esac\n'%f)
m("10-stand-gruen", "sk", 'if [ -z "$befunde" ]; then\n  exit 0', 'if [ -z "$befunde" ]; then\n  touch .stand; exit 0')
m("10-stand-rot", "sk", "    }' >&2\nexit 1", "    }' >&2\ntouch .stand\nexit 1")
m("10-wurzel-aktuell", "sk", 'wurzel="${1:-.}"', 'wurzel="${1:-/nirgends}"')
m("11-check-nicht-in-gate-checks", "mk", "GATE_CHECKS += meldungskatalog-check meldungskatalog-check-gegenprobe", "GATE_CHECKS += meldungskatalog-check-gegenprobe")
m("11-gegenprobe-nicht-in-gate-checks", "mk", "GATE_CHECKS += meldungskatalog-check meldungskatalog-check-gegenprobe", "GATE_CHECKS += meldungskatalog-check")
m("11-check-ueberliest-rot", "mk", "@bash tools/harness/meldungskatalog-check.sh\n", "@bash tools/harness/meldungskatalog-check.sh || true\n")
m("11-gegenprobe-ueberliest-rot", "mk", "@bash tools/harness/meldungskatalog-check-gegenprobe.sh\n", "@bash tools/harness/meldungskatalog-check-gegenprobe.sh || true\n")
m("11-check-ruft-nichts", "mk", "@bash tools/harness/meldungskatalog-check.sh\n", "@true\n")
m("11-check-falscher-pfad", "mk", "@bash tools/harness/meldungskatalog-check.sh\n", "@bash tools/harness/meldungskatalog-check-x.sh\n")
m("11-ohne-phony", "mk", ".PHONY: meldungskatalog-check meldungskatalog-check-gegenprobe\n", "")
m("11-ausnahme-ohne-gegenprobe", "gs", " meldungskatalog-check-gegenprobe", "")
m("11-ausnahme-mit-check", "gs", " meldungskatalog-check-gegenprobe", " meldungskatalog-check-gegenprobe meldungskatalog-check")
M = [x for x in M if x]

only = sys.argv[1:]
res = []
for name, key, old, new, n in M:
    if only and not any(o in name for o in only): continue
    t = tempfile.mkdtemp(dir=os.environ.get("TMPDIR", "/tmp"))
    for rel in (SK, GP, MK, GS, "Makefile"):
        os.makedirs(os.path.dirname(os.path.join(t, rel)) or t, exist_ok=True)
        shutil.copy(os.path.join(R, rel), os.path.join(t, rel))
    p = os.path.join(t, FILES[key]); s = open(p).read()
    if s.count(old) != n:
        print("%s: ALT-STRING %d x gefunden" % (name, s.count(old))); shutil.rmtree(t); continue
    open(p, "w").write(s.replace(old, new))
    r = subprocess.run(["bash", os.path.join(t, GP)], capture_output=True, text=True)
    red = sorted(set(re.findall(r"Fall '([^']*)'", r.stderr)))
    print("%s [exit %d]: %s" % (name, r.returncode, " ".join(red) if red else "GRÜN (Mutant überlebt)"))
    shutil.rmtree(t)
```

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-07); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — 13×, über der Schwelle, verkörpert ([`AGENTS.md`](../../../../AGENTS.md) §3.10); die Gegenprobe ist eigener Liefer-Punkt (§2).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 17×, verkörpert (§3.11); Sensor-Datei und Zeile in `harness/README.md` sagen nur zu, was Skript und Gegenprobe prüfen.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 11×, verkörpert (§3.12); die Randformen stehen in §6 und sind vor dem Code zu entscheiden.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` — 1×; der Slice ist wellenlos, §4 nennt die Entscheidung des Architects vor dem Code.
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` — 2×; mit diesem Slice das dritte Mal berührt, falls eine Lesart erst im Code entschieden wird (§6).
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` — 3×, verkörpert (§3.13); die Geber nennen diesen Slice, dieser nennt sie in §1.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` — 1×; ein Gate, das Dateien sucht, darf eine nicht lesbare Quelle nicht still übergehen (§6, Punkt 7).
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` — 2×; das Fragment unter `harness/mk/` muss so an `GATE_CHECKS` hängen, wie es sich liest (Gegenprobe mit Makefile, §2).

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure ihrer Welle, nicht hier.

**Nachgezählt** mit der Sendung aus `slice-harness-vertraege-umzug-gate-kette` (Ort neuer
Kennungen: `harness/vertraege/`, nächste freie Nummer 053; `AGENTS.md` §3.13): drei
Liefer-Punkte (§2: Prüfung, Gegenprobe, Sensor-Datei, unverändert). Schichten ohne die bedingte
Zeile von §3 (Bestands-Befund in `internal/hexagon/model/fehler.go` und `SPEC-034`): eine,
Harness (Vertragsdatei, `tools/`, `harness/mk/`, `harness/sensors/`, `harness/README.md`). Mit
der bedingten Zeile kämen Produkt-Code und Spezifikation dazu, also drei; das war vor der
Sendung ebenso (die Vertragsdatei lag in der Spezifikation) und ist Sache der Entscheidung des
Architects zu *Bestand* in §6, die Sendung ändert daran nichts.

**Nachgezählt** mit der Vorprüfung des Architect (2026-10-17, `AGENTS.md` §3.13): drei
Liefer-Punkte (§2: Prüfung mit Anschluss, Gegenprobe, Sensor-Datei mit Zeilen in
`harness/README.md`). Schichten: eine, Harness; die bedingte Zeile für `fehler.go` und
`SPEC-034` ist entfallen (Bestand grün). Das Handbuch wird gelesen, nicht geändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
