# Slice slice-version-option: `--version` (und `-v`) als globale Option

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** — (wellenlos). `welle-v1-nachtrag` geprüft und **nicht** gewählt: Der Slice liefert keine Anforderung dieser Welle, und ihn aufzunehmen dehnte ihren Closure-Trigger (Out-of-Scope-Disziplin, §6 ihres Plans). Eine Welle für die Arbeit nach v1 gibt es noch nicht; entsteht eine, kann sie ihn einsammeln.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung)

**Berührte Spec-Stellen:** `LH-FA-01.a` (CLI-Abschnitt der Spezifikation: *Hilfe vor Prüfung*, *Ende der Optionen* und der Satz zur Option `--version`) · `LH-FA-01` · `LH-FA-14.a` (nur als Abgrenzung genannt: die Ausführlichkeit trägt `--log-level`, der Slice ändert sie nicht) · `SPEC-014` (Exit-Code der ungültigen CLI-Verwendung; gelesen, der Slice ändert ihn nicht)

**Entscheidung:** Entscheidung des Nutzers vom 2026-10-21 (wörtlich): „'--version -v' als slice einplanen“. Entschieden ist damit, **dass** der Slice eingeplant wird und dass er beide Formen zum Gegenstand hat; **wie** die Option sich verhält (Stelle, Ausgabe, ob `-v` wirklich belegt wird) ist nicht entschieden und steht in §6 als zwei Fragen an den Nutzer und eine Festlegung durch den Architect. Keine ADR.

**Verantwortlich:** pt9912
**Autor:** Planner. **Datum:** 2026-10-21; nach `slice-harness-regeln-rueckbau` nachgezogen am 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Das gebaute Binary nimmt `--version` als globale Option an (und `-v`, wenn der Nutzer es entscheidet, §6 Frage 1) und verhält sich so, wie die Spezifikation es vorher festlegt; Hilfetext und Benutzerhandbuch nennen die Option im Ist-Zustand. Das ist eine **Änderung der Spezifikation**, kein Fehler im Code: Heute sagt `spec/spezifikation.md` (`LH-FA-01.a`) „eine Option `--version` gibt es nicht (`PGR-E2001`)“, das Handbuch sagt „Eine Option `--version` gibt es nicht; die Version zeigt `pgwire-recorder version`“, und `TestParseVersion` in `internal/adapters/driving/cli/cli_test.go` hält mit seiner Abdeckungs-Deklaration genau das fest. Der Slice dreht diese Aussage um — an allen vier Stellen oder an keiner.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Kommando `pgwire-recorder version` — bleibt unverändert (Klasse 2, Bestand bleibt bewusst stehen). Grund: Es ist in der Spezifikation beschrieben, im Handbuch genannt und von Tests gedeckt; die Option tritt daneben und ersetzt es nicht. Ein Slice, der es entfernt, wäre eine eigene Entscheidung des Nutzers.
- `README.md` — Bestand bleibt bewusst stehen (Klasse 2). Grund: Es nennt weder das Kommando `version` noch die Option (geprüft 2026-10-22 mit `grep`); es gibt nichts umzudrehen.
- Ein allgemeines Schema für Kurzoptionen (`-l` für `--listen` und dergleichen) — anderer Vorgang (Klasse 3). Grund: Heute kennt der Parser als kurze Formen allein die vier Hilfe-Angaben; ein Schema über alle Kommandos legt Namen, Kollisionen und Hilfetexte für jede Option fest und ist weder klein noch von dieser Entscheidung gedeckt.
- Änderung des Lastenhefts — anderer Vorgang (Klasse 3). Grund: `LH-FA-01` nennt keine Optionsnamen, und die Namen der Optionen legt das Technik-Stratum fest (`LH-FA-01.a`); eine Anforderung kommt nicht hinzu. Hält der Architect es doch für eine neue Anforderung, ist das ein eigener Vorgang mit **Minor-Bump** des Lastenhefts (`harness/conventions.md` §Versionierung des Lastenhefts), nicht dieser Slice.
- `-v` als „verbose“ — Bestand bleibt bewusst stehen (Klasse 2). Grund: Die Ausführlichkeit der Ausgabe trägt `--log-level` (`LH-FA-14.a`); dieser Slice legt `-v` nicht auf „verbose“, und belegt der Nutzer `-v` mit der Version (§6 Frage 1), ist der Name dafür an dieser Stelle weg.
- Eine Versionsangabe für andere Gegenstände (Version des Aufzeichnungsformats, der Konfigurationsdatei, des Protokolls) — anderer Vorgang (Klasse 3). Grund: `--version` meint die Programmversion; was `version: 2` einer Aufzeichnung bedeutet, steht anderswo und bleibt unberührt.
- Produkt-Code außerhalb des CLI-Adapters (Services, Ports, Driven-Adapter) — Schicht-Abgrenzung (Klasse 4). Grund: Die Option wird beim Lesen der Argumente beantwortet; die Version selbst liefert der Bestand, der schon `pgwire-recorder version` beantwortet.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 — die Option** ([`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), `LH-FA-01.a`): Das gebaute Binary beantwortet jede Form, die die Spezifikation als Versions-Angabe festlegt, an jeder Stelle, die sie nennt, mit der dort festgelegten Ausgabe und dem dort festgelegten Exit-Code; jede andere Form bleibt, was sie heute ist. Ein Test je Zusage der Spezifikationsstelle, bei einer Aufzählung (Formen, Stellen) ein Fall je Element, **ohne Mutation** (`AGENTS.md` §3.10: eine Option ohne Wirkung auf den Record-/Replay-Pfad). `TestParseVersion` und seine Abdeckungs-Deklaration beschreiben danach das neue Verhalten, und `make abdeckung-check` ist grün.
- [ ] **Liefer-Punkt 2 — Doku im Ist-Zustand** ([`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung)): Hilfetext (globale Hilfe und, wenn die Spezifikation es so entscheidet, die Hilfe je Kommando) und `docs/user/benutzerhandbuch.md` (der Absatz, der heute „Eine Option `--version` gibt es nicht“ sagt, und jede Stelle, die der neuen Ausgabe oder dem Satz zu `stdout` widerspräche) nennen die Option so, wie das gebaute Binary sie zeigt — jeder Satz mit einem Lauf gegen das Binary belegt, keine Chronik, kein Zielstand, im Handbuch kein Verweis auf Spezifikation, ADR, Slice oder Welle (`AGENTS.md` §3.11). `grep -rn -- '--version' docs/user README.md` am Endstand findet keine Stelle mehr, die die Option verneint.
- [ ] `make gates` grün; die betroffenen Gegenproben nach `harness/README.md` §Sensors gefahren, falls ein Gate oder ein Vertrag berührt ist (hier nicht erwartet, §3).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/cli/cli.go` | update | Erkennung der Versions-Angabe beim Lesen der Argumente (neben der Hilfe-Angabe), Ausgabe und Exit-Code nach der Spezifikation; Zeile in `usage` und, wenn so entschieden, in `hilfen` |
| `internal/adapters/driving/cli/cli_test.go` | update | `TestParseVersion` kehrt sich um: Fälle je Form und je Stelle, mit Abdeckungs-Deklaration |
| `test/integration/<name>_e2e_test.go` | neu oder update | ein Fall gegen das gebaute Binary für jeden Satz des Handbuchs (Ausgabe, Exit-Code), mit Abdeckungs-Deklaration (`LH-FA-01/...`) |
| `docs/user/benutzerhandbuch.md` | update | der Absatz zur Hilfe und zur Version (heute „Eine Option `--version` gibt es nicht“) |
| `docs/user/abdeckung-*.md` | update (erzeugt) | `make abdeckung` nach den neuen Deklarationen |

**Nicht in dieser Liste, weil vor dem Start:** `spec/spezifikation.md` (`LH-FA-01.a`, mit der Zeile der Historie) — die Festlegung aus §6 trifft der Architect dort **vor** dem ersten Code-Commit (`AGENTS.md` §3.12); dieser Commit gehört nicht zum Slice, er ist sein Start-Trigger (§4). Ändert der Architect die Spezifikation **während** des Slice (Rückgabe durch den Implementer), steht sie ab dann in dieser Liste und zählt als dritte Schicht — dann greift der vorgemerkte Schnitt aus §4.

Nicht in dieser Liste: `README.md` (§1), Produkt-Code außerhalb des CLI-Adapters, Gates, Verträge unter `harness/vertraege/`, das Lastenheft (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-erster-release-homebrew-nachweis` und `slice-harness-regeln-rueckbau` liegen in `done/`. Der Nutzer hat den Start der Produkt-Slices am 2026-10-22 mit „leg los“ freigegeben (wörtlich, so vom Koordinator übermittelt). Offen sind (1) die zwei Fragen an den Nutzer aus §6, (2) danach die Festlegung durch den Architect in `spec/spezifikation.md` (`LH-FA-01.a`), von ihm geprüft (`AGENTS.md` §3.12); der Implementer beginnt erst, wenn beides vorliegt. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik und Messbefehl: `AGENTS.md` §3.14; Pfade des Plans: `internal/adapters/driving/cli test docs/user`); die Basis ist der Commit des Übergangs nach `in-progress/`, sein Hash: `49da54b6993287b9d1902b30aacef66e8afc4703`, Zweig `impl-version-option`. Dann hält der Implementer an und meldet dem Planner.
- `in-progress` → `next`: **vorgemerkter Schnitt.** Zählt die Spezifikation als dritte Schicht (weil der Architect sie im Slice ändert, §3) oder entscheidet er ein Kurzoptionen-Schema (§1), geht der Handbuch-Teil in einen eigenen Slice `slice-version-option-doku` direkt dahinter (`AGENTS.md` §3.11 und §3.13, Schichtzählung); dieser Slice behält dann Parser, Tests und Hilfetext.
- `in-progress` → `open`: Der Architect gibt eine Entscheidung zurück, die der Nutzer treffen muss und §6 nicht nennt — der Slice wartet nicht in `in-progress/`, er geht zurück und hält das WIP-Limit frei.

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

**Randformen** (`AGENTS.md` §3.12). Die Option hat keine Wirkung auf den Record-/Replay-Pfad und ist keine neue Wire-, Datei- oder Formatgrenze; §6 nennt deshalb die Spezifikationsstelle der Geschwister-Option und was von ihr abweicht. **Geschwister-Option:** die Hilfe-Angabe, `spec/spezifikation.md` `LH-FA-01.a`, Absätze *Hilfe vor Prüfung* und *Ende der Optionen*. Die Versions-Angabe übernimmt: ganzes Argument; nach dem ersten `--` ein gewöhnliches Argument; liest weder Umgebungsvariablen noch Konfigurationsdatei; `PGR-E2001` bleibt der Code jeder Ablehnung, kein neuer Meldungstext; nichts hängt am Betriebssystem (Argumente, `stdout`, Exit-Code). **Abweichungen** — der Architect legt sie in `LH-FA-01.a` vor dem ersten Code-Commit fest (Empfehlungen des Planners in Klammern; Entscheidungen des Nutzers sind es nicht):

1. *Ausgabe und Exit-Code* (dieselbe Zeile wie `pgwire-recorder version`, auf `stdout`, Exit-Code 0; dann bleiben die Sätze zu `stdout` in Spezifikation und Handbuch wahr).
2. *Formen* (genau `--version` und, je Frage 1, `-v`, jeweils ganze Argumente und **ohne** `=`-Wert: `--version=x` bleibt `PGR-E2001`, anders als die Hilfe-Angabe, die `=Wert` zulässt).
3. *Stelle und Vorrang* (je Frage 2. Bei B nur als einziges Argument; dann ist kein Vorrang gegenüber der Hilfe zu regeln: `--version --help` bleibt die Hilfe, `version --version` bleibt `PGR-E2001`. Bei A wie die Hilfe-Angabe an jeder Stelle vor dem ersten `--`, die Hilfe geht der Version vor, `version --version` ist die Version).

**Fragen an den Nutzer** — Empfehlung des Planners, Entscheidung des Nutzers steht aus:

1. *`-v` belegen?* Messung des Planners (2026-10-21): `-v` ist im Repo nicht vergeben (der Parser kennt als kurze Formen allein die vier Hilfe-Angaben; Treffer allein als Testwert `"v"` und in `docker rm -fv`). Optionen: **A** `--version` und `-v` (wie am 2026-10-21 genannt; Empfehlung) · **B** nur `--version` · **C** `--version` und `-V` (in den PostgreSQL-Werkzeugen wie `pg_dump` ist `-V` die Version und `-v` „verbose“; nicht im Repo gemessen).
2. *Wo gilt die Option?* **B** nur als einziges Argument, `pgwire-recorder --version`, sonst wie heute `PGR-E2001`; `-v` nimmt den Namen dann nur an dieser Stelle (Empfehlung: kleinste Änderung, kleinste Tabelle, hält `record -v` frei) · **A** an jeder Stelle wie die Hilfe, vor jeder Prüfung (`record --version` gibt die Version, auch ohne Pflichtoptionen).

**Risiken:**

- `-v` wird mit der Version belegt, und ein Nutzer, der `-v` als „verbose“ erwartet, bekommt ohne Fehlermeldung die Version; der Name ist danach an dieser Stelle nicht mehr frei — **Ausgang:** offen bis Closure (entschieden mit Frage 1 und 2; wird `-v` nicht belegt, entfällt das Risiko mit dieser Begründung).
- Die Aussage „eine Option `--version` gibt es nicht“ steht an vier Stellen (Spezifikation, Handbuch, `TestParseVersion`, Abdeckungs-Deklaration dieses Tests); bleibt eine stehen, widerspricht das Repo sich selbst — **Ausgang:** offen bis Closure (`grep -rn -- '--version' spec docs/user README.md internal test` am Endstand, je Treffer eine Zeile im Bericht).

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

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (gemergter Stand, 2026-10-22). Treffer, die diesen Slice betreffen: `zusage-im-kommentar-weiter-als-pruefung` (verkörpert in `AGENTS.md` §3.11; DoD 2 bindet jeden Satz an einen Lauf gegen das Binary), `plan-folgt-korrektur-nicht` (verkörpert in §3.9; die Spezifikation steht vor dem Start, §3) und `spec-randform-erst-im-review-entschieden` (verkörpert in §3.12; §6 führt die Festlegung als Aufgabe des Architect). Keiner dieser Einträge erreicht mit diesem Slice eine Schwelle neu (die verkörperten zählen erst wieder, wenn die Regel trotz ihrer Verkörperung versagt).

**Zählung** (Planner, 2026-10-22, `AGENTS.md` §3.13; Geber `slice-harness-regeln-rueckbau`, dessen Reihenfolge den Start in §4 und die Marke 1500 setzt; Randformen-Tabelle, Mutationspflicht, Schätzung und Meldeschwelle entfallen): **Liefer-Punkte: 2** (1: die Option mit ihren Fällen; 2: Hilfetext und Handbuch). **Schichten: zwei** — der CLI-Adapter des Hexagons mit seinen Tests (`internal/adapters/driving/cli/`, dazu der E2E-Fall unter `test/integration/`, der dieselbe Fläche prüft) und die Nutzer-Doku (`docs/user/`; `README.md` ist unberührt, §1); die Spezifikation ist eine **eigene** Schicht, liegt aber vor dem Start und damit außerhalb des Slice (§3, §4). Zählt sie der Architect doch zum Slice — oder zählt er den E2E-Fall als eigene Schicht —, sind es drei, und der vorgemerkte Schnitt aus §4 (`slice-version-option-doku`) greift. Der Slice trägt keine Sendung in einen anderen Plan ein und nimmt keine an.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und Entscheidung gehen dem Code voraus; der Architect legt die Festlegung vor dem ersten Code-Commit fest).
