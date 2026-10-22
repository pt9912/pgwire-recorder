# Slice slice-version-option: `--version` (und `-v`) als globale Option

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** — (wellenlos). `welle-v1-nachtrag` geprüft und **nicht** gewählt: Ihr Ziel sind die zwei Anforderungsreste von v1 und der Nachweis von Abnahmeszenario 11, ihr Closure-Trigger hängt am Homebrew-Nachweis; dieser Slice liefert keine Anforderung dieser Welle und läge hinter ihrem letzten Slice. Ihn aufzunehmen dehnte die Welle, bis ihr Closure-Trigger später käme (§6 ihres Plans, Out-of-Scope-Disziplin). Eine Welle für die Arbeit nach v1 gibt es noch nicht; entsteht eine, kann sie ihn einsammeln.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung)

**Berührte Spec-Stellen:** `LH-FA-01.a` (CLI-Abschnitt der Spezifikation: *Hilfe vor Prüfung*, *Ende der Optionen* und der Satz zur Option `--version`) · `LH-FA-01` · `LH-FA-14.a` (nur als Abgrenzung genannt: die Ausführlichkeit trägt `--log-level`, der Slice ändert sie nicht) · `SPEC-014` (Exit-Code der ungültigen CLI-Verwendung; gelesen, der Slice ändert ihn nicht)

**Entscheidung:** Entscheidung des Nutzers vom 2026-10-21 (wörtlich): „'--version -v' als slice einplanen“. Entschieden ist damit, **dass** der Slice eingeplant wird und dass er beide Formen zum Gegenstand hat; **wie** die Option sich verhält (Ausgabe, Exit-Code, Formen, Stelle, Vorrang, ob `-v` wirklich belegt wird) ist nicht entschieden und steht in §6 als Rückgabe an den Architect. Keine ADR.

**Verantwortlich:** pt9912
**Autor:** Planner. **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Das gebaute Binary nimmt `--version` als globale Option an (und `-v`, wenn der Architect sie belegt) und verhält sich so, wie die Spezifikation es vorher festlegt; Hilfetext, Benutzerhandbuch und `README.md` nennen die Option im Ist-Zustand. Das ist eine **Änderung der Spezifikation**, kein Fehler im Code: Heute sagt `spec/spezifikation.md` (`LH-FA-01.a`) „eine Option `--version` gibt es nicht (`PGR-E2001`)“, das Handbuch sagt „Eine Option `--version` gibt es nicht; die Version zeigt `pgwire-recorder version`“, und `TestParseVersion` in `internal/adapters/driving/cli/cli_test.go` hält genau das fest. Der Slice dreht diese Aussage um — an allen vier Stellen oder an keiner.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Kommando `pgwire-recorder version` — bleibt unverändert (Klasse 2, Bestand bleibt bewusst stehen). Grund: Es ist in der Spezifikation beschrieben, im Handbuch genannt und von Tests gedeckt; die Option tritt daneben und ersetzt es nicht. Ein Slice, der es entfernt, wäre eine eigene Entscheidung des Nutzers.
- Ein allgemeines Schema für Kurzoptionen (`-l` für `--listen` und dergleichen) — anderer Vorgang (Klasse 3). Grund: Heute kennt der Parser als kurze Formen allein die vier Hilfe-Angaben; ein Schema über alle Kommandos legt Namen, Kollisionen und Hilfetexte für jede Option fest und ist weder klein noch von dieser Entscheidung gedeckt.
- Änderung des Lastenhefts — anderer Vorgang (Klasse 3). Grund: `LH-FA-01` nennt keine Optionsnamen, und die Namen der Optionen legt das Technik-Stratum fest (`LH-FA-01.a`); eine Anforderung kommt nicht hinzu. Hält der Architect es doch für eine neue Anforderung, ist das ein eigener Vorgang mit **Minor-Bump** des Lastenhefts (`harness/conventions.md` §Versionierung des Lastenhefts) — Rückgabe (m) in §6.
- `-v` als „verbose“ — Bestand bleibt bewusst stehen (Klasse 2). Grund: Die Ausführlichkeit der Ausgabe trägt `--log-level` (`LH-FA-14.a`); dieser Slice legt `-v` nicht auf „verbose“, und belegt der Architect `-v` mit der Version, ist der Name dafür dauerhaft weg (Randform (f) und Risiko in §6).
- Eine Versionsangabe für andere Gegenstände (Version des Aufzeichnungsformats, der Konfigurationsdatei, des Protokolls) — anderer Vorgang (Klasse 3). Grund: `--version` meint die Programmversion; was `version: 2` einer Aufzeichnung bedeutet, steht anderswo und bleibt unberührt.
- Produkt-Code außerhalb des CLI-Adapters (Services, Ports, Driven-Adapter) — Schicht-Abgrenzung (Klasse 4). Grund: Die Option wird beim Lesen der Argumente beantwortet; die Version selbst liefert der Bestand, der schon `pgwire-recorder version` beantwortet.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 — die Option** ([`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), `LH-FA-01.a`): Das gebaute Binary beantwortet jede Form, die die Spezifikation als Versions-Angabe festlegt, an jeder Stelle, die sie nennt, mit der dort festgelegten Ausgabe und dem dort festgelegten Exit-Code; jede andere Form bleibt, was sie heute ist. Je Zusage der Spezifikationsstelle ein Fall und eine Mutation, die nur ihn bricht, vom Implementer selbst rot gesehen, mit Zusage · Mutation · rotem Test im Bericht (`AGENTS.md` §3.10); die Aufzählungen der Spezifikation (Formen, Stellen, Vorrang) zählen je Element. `TestParseVersion` in `internal/adapters/driving/cli/cli_test.go` und seine Abdeckungs-Deklaration beschreiben danach das neue Verhalten, und `make abdeckung-check` ist grün.
- [ ] **Liefer-Punkt 2 — Doku im Ist-Zustand** ([`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung)): Hilfetext (globale Hilfe und, wenn die Spezifikation es so entscheidet, die Hilfe je Kommando), `docs/user/benutzerhandbuch.md` (der Absatz, der heute „Eine Option `--version` gibt es nicht“ sagt, und der Kopf, wenn er berührt ist) und `README.md` nennen die Option so, wie das gebaute Binary sie zeigt — jeder Satz mit einem Lauf gegen das Binary belegt, keine Chronik, kein Zielstand, im Handbuch kein Verweis auf Spezifikation, ADR, Slice oder Welle (`AGENTS.md` §3.11). `grep -rn -- '--version' docs/user README.md` am Endstand findet keine Stelle mehr, die die Option verneint.
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
| `internal/adapters/driving/cli/cli_test.go` | update | `TestParseVersion` kehrt sich um: Fälle je Form, je Stelle und je Vorrang, mit Abdeckungs-Deklaration |
| `test/integration/<name>_e2e_test.go` | neu oder update | ein Fall gegen das gebaute Binary für jeden Satz des Handbuchs (Ausgabe, Exit-Code), mit Abdeckungs-Deklaration (`LH-FA-01/...`) |
| `docs/user/benutzerhandbuch.md` | update | der Absatz zur Hilfe und zur Version (heute „Eine Option `--version` gibt es nicht“) |
| `README.md` | update | nur, wenn die Option dort genannt gehört (heute nennt `README.md` sie nicht) |
| `docs/user/abdeckung-*.md` | update (erzeugt) | `make abdeckung` nach den neuen Deklarationen |

**Nicht in dieser Liste, weil vor dem Start:** `spec/spezifikation.md` (`LH-FA-01.a`) — die Randformen aus §6 entscheidet der Architect dort **vor** dem ersten Code-Commit (`AGENTS.md` §3.12); dieser Commit gehört nicht zum Slice, er ist sein Start-Trigger (§4). Ändert der Architect die Spezifikation **während** des Slice (Randform-Rückgabe durch den Implementer), steht sie ab dann in dieser Liste und zählt als dritte Schicht — dann greift der vorgemerkte Schnitt aus §4.

Nicht in dieser Liste: Produkt-Code außerhalb des CLI-Adapters, Gates, Verträge unter `harness/vertraege/`, das Lastenheft (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next`): `slice-erster-release-homebrew-nachweis` liegt in `done/` (Stelle in der Reihenfolge: nach dem Homebrew-Nachweis, vor dem Slice zum vollen Feldvergleich von `play`; Entscheidung des Nutzers vom 2026-10-21 zur Reihenfolge, die Stelle selbst ist Empfehlung des Planners). **Reihenfolge seit dem 2026-10-22:** hinter `slice-harness-regeln-rueckbau` (Entscheidung des Nutzers vom 2026-10-22, „ja“ auf den Bericht des Planners, so vom Koordinator übermittelt: zuerst der Rückbau, dann dieser Slice); der Rückbau ändert Marke, Meldeschwelle, Schätzung und den Geltungsbereich von `AGENTS.md` §3.10 und §3.12. Marke, Meldeschwelle und Schätzung dieses Plans sind Bestand und bleiben stehen; der Planner zieht §4 beim Anspruch nach der neuen Regel nach (Marke ohne Meldeschwelle und Schätzung; Entscheidung des Nutzers vom 2026-10-22, „1. ja“, §6 Randform 7 von `slice-harness-regeln-rueckbau`). **Start** (`next` → `in-progress`): zusätzlich liegt `slice-harness-regeln-rueckbau` in `done/`, und sind die Randformen (a) bis (l) aus §6 in `spec/spezifikation.md` entschieden und vom Architect geprüft (`AGENTS.md` §3.12); eine offene Randform entscheidet der Implementer nicht. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal/adapters/driving/cli test docs/user README.md` (Tests eingeschlossen, `AGENTS.md` §3.14) **über 400** liegt; **Meldeschwelle 320** — erreicht die Summe 320, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet. **Basis** ist der Commit des Übergangs nach `in-progress/`; seinen Hash trägt der Planner beim Übergang hier ein. Schätzung des Planners in dieser Metrik, ohne Messung: Parser und Hilfetexte 40 bis 70, Fälle im Paket-Test mit dem Faktor rund 2 gegenüber dem Produktcode (Tabelle je Form, Stelle und Vorrang, §3.10 zählt je Element) 80 bis 150, E2E-Fall 20 bis 40, Handbuch und `README.md` 20 bis 40, erzeugte Abdeckungstabellen 5 bis 15 — zusammen 165 bis 315, Mitte 240 (60 % der Marke).
- `in-progress` → `next`: **vorgemerkter Schnitt.** Zählt die Spezifikation als dritte Schicht (weil der Architect sie im Slice ändert, §3) oder entscheidet er ein Kurzoptionen-Schema (§1), geht der Handbuch- und `README.md`-Teil in einen eigenen Slice `slice-version-option-doku` direkt dahinter (`AGENTS.md` §3.11 und §3.13, Schichtzählung); dieser Slice behält dann Parser, Tests und Hilfetext.
- `in-progress` → `open`: Der Architect entscheidet eine Randform aus §6 nicht, weil sie eine Entscheidung des Nutzers verlangt (etwa (f): `-v` belegen oder nicht) — der Slice wartet nicht in `in-progress/`, er geht zurück und hält das WIP-Limit frei.

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

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert einen neuen Vertrag (eine Option am Protokollrand der Kommandozeile); **keine dieser Randformen ist entschieden**, alle gehen als Rückgabe an den Architect, Entscheidungsort ist `spec/spezifikation.md` (`LH-FA-01.a`), vor dem ersten Code-Commit. Der Planner entscheidet sie nicht und schlägt keine Lösung vor, wo er keine messen kann; was er gemessen hat, steht bei (f).

| Randform | Frage an den Architect | Stelle |
|---|---|---|
| (a) Ausgabe und Exit-Code | Dieselbe Zeile wie `pgwire-recorder version` (`pgwire-recorder <Version>` auf `stdout`) und Exit-Code 0? Oder eine eigene Form? | `LH-FA-01.a`, Abschnitt zu `version` |
| (b) Formen | Welche Schreibweisen sind eine Versions-Angabe: `--version`, `-version`, `--v`, `-v`, mit `=` und beliebigem Wert (`--version=`, `-v=x`)? Die Liste ist abzuschließen wie die der Hilfe-Angabe („Die Liste ist abgeschlossen“, vier Formen); jede genannte Form ist eine eigene Zusage mit eigenem Fall (§3.10) | `LH-FA-01.a`, *Hilfe vor Prüfung* als Vorbild |
| (c) Stelle | Gilt die Option vor dem Kommando (`--version record`), ohne Kommando, hinter einem bekannten Kommando (`record --version`), hinter einem unbekannten (`bogus --version`) und an der Stelle eines Optionswerts (`--input --version`)? Je Stelle eine Entscheidung | `LH-FA-01.a`, *Hilfe vor Prüfung* |
| (d) Vorrang | Geht die Versions-Angabe jeder Prüfung vor, wie die Hilfe („Hilfe vor Prüfung“: unbekanntes Kommando, unbekannte Option, ungültiger Wert, fehlende Pflichtoption, Umgebungsvariable, Konfigurationsdatei)? Und wer gewinnt, wenn Hilfe-Angabe und Versions-Angabe zusammen dastehen (`--help --version`, `--version --help`) — die erste, die Hilfe immer, oder die Version immer? | `LH-FA-01.a`, *Hilfe vor Prüfung* |
| (e) Nach `--` | Ist eine Versions-Angabe nach dem ersten `--` ein gewöhnliches Argument, wie die Hilfe-Angabe? Dann lehnen `record`, `replay` und `play` sie mit `PGR-E2001` ab, und `--input=--version` ist ein Wert | `LH-FA-01.a`, *Ende der Optionen* |
| (f) `-v` und „verbose“ | **Messung des Planners (2026-10-21, kein Vorschlag):** `-v` ist im Repo nicht vergeben. Der Parser kennt als kurze Formen allein `-h`, `--h`, `-help`, `--help` (`internal/adapters/driving/cli/cli.go`); `grep` nach `-v` als Option in `cmd/`, `internal/`, `test/`, `spec/`, `docs/user/benutzerhandbuch.md` und `README.md` findet keine Stelle (Treffer allein als Testwert `"v"` und in `docker rm -fv`), und die Ausführlichkeit trägt `--log-level`. **Frage:** `-v` jetzt für die Version belegen — dann ist der Name für „verbose“ dauerhaft weg und ein Nutzer, der ihn aus anderen Werkzeugen als „verbose“ kennt, bekommt die Version — oder nur `--version` liefern? Die Entscheidung des Nutzers vom 2026-10-21 nennt beide Formen, legt aber das Verhalten nicht fest | `LH-FA-01.a`; Entscheidung des Nutzers nötig, wenn der Architect sie nicht allein trägt |
| (g) `version` mit der Option | `version` nimmt heute keine Argumente (`PGR-E2001`, „version nimmt keine Argumente“, `TestParseVersion`). Bleibt `version --version` ein Fehler, oder ist es die Version? | `LH-FA-01.a`, Abschnitt zu `version` |
| (h) Meldung bei Ablehnung | Bleibt `PGR-E2001` der Code für jede Stelle, an der eine Versions-Angabe abgelehnt wird, und ändert sich ein Meldungstext? (`make meldungskatalog-check` hält Code, Tabelle und Handbuch zusammen) | `LH-FA-01.a`; `SPEC-034` (Katalog), wenn ein Text hinzukommt |
| (i) Hilfetext | Nennt die globale Hilfe die Option, und nennt sie jede Kommando-Hilfe (`hilfen`)? Unter welcher Überschrift? | `LH-FA-01.a` |
| (j) Umgebung und Konfigurationsdatei | Liest die Option nichts, wie `version` („`version` liest weder Umgebungsvariablen noch Konfigurationsdatei“)? | `LH-FA-01.a` |
| (k) Plattform (`AGENTS.md` §3.12) | Keine Zusage hängt an Betriebssystem oder Dateisystem — Argumentliste, `stdout` und Exit-Code verhalten sich auf allen Plattformen von `LH-QA-03` gleich; geprüft wird wie für `version` unter Linux. Zu bestätigen, nicht zu entscheiden | `LH-FA-01.a`; `LH-QA-03` |
| (l) Aufruf-Voraussetzung, Aufrufarten des Treibers (`AGENTS.md` §3.12) | Entfällt: Der Slice liefert keine Werte über mehrere Aufrufe einer Session und kein Ergebnis, das ein Hook oder ein Gate glaubt | — |
| (m) Stratum der Festlegung | Ist die Option eine Festlegung des Technik-Stratums unter `LH-FA-01` (dann bleibt das Lastenheft unberührt, §1) oder eine neue Anforderung (dann ein eigener Vorgang mit Minor-Bump, `harness/conventions.md` §Versionierung des Lastenhefts)? | `harness/README.md` §Source precedence |

**Risiken:**

- `-v` wird mit der Version belegt, und ein Nutzer, der `-v` als „verbose“ erwartet, bekommt ohne Fehlermeldung die Version; der Name ist danach nicht mehr frei — **Ausgang:** offen bis Closure (entschieden in Randform (f); wird `-v` nicht belegt, entfällt das Risiko mit dieser Begründung).
- Die Liste der Formen in (b) und der Stellen in (c) wächst wie bei der Hilfe-Angabe (vier Formen × fünf Stellen), und §3.10 verlangt je Element einen Fall mit eigener Mutation; der Slice wächst über Marke oder Meldeschwelle — **Ausgang:** offen bis Closure (Marke 400, Meldeschwelle 320, Metrik in §4).
- Die Aussage „eine Option `--version` gibt es nicht“ steht an vier Stellen (Spezifikation, Handbuch, `TestParseVersion`, Abdeckungs-Deklaration dieses Tests); bleibt eine stehen, widerspricht das Repo sich selbst — **Ausgang:** offen bis Closure (`grep -rn -- '--version' spec docs/user README.md internal test` am Endstand, je Treffer eine Zeile im Bericht).
- Die Spezifikation legt den Vorrang (d) anders fest, als die Hilfe ihn kennt, und der Parser braucht dafür mehr als einen Zweig (zwei Durchläufe über die Argumente) — **Ausgang:** offen bis Closure (bleibt der Produktcode über der Schätzung, zählt das als Auftreten von `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`).

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

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (gemergter Stand, 2026-10-21). Treffer, die diesen Slice betreffen: `spec-randform-erst-im-review-entschieden` (47×, verkörpert in `AGENTS.md` §3.12) — der Grund, warum §6 die Randformen als Rückgabe führt und keine selbst entscheidet; `randform-im-code-entschieden-dann-zurueckgegeben` (26×, verkörpert) — derselbe Rand, aus der Sicht des Implementers; `negativtests-fehlen-bei-neuem-vertrag` (64×, verkörpert in `AGENTS.md` §3.10) — je Zusage und je Element der Aufzählungen ein Fall mit Mutation, in DoD 1 genannt; `zusage-im-kommentar-weiter-als-pruefung` (89×, verkörpert in §3.11) — DoD 2 bindet jeden Satz an einen Lauf gegen das Binary; `schaetzung-in-zeilen-unterschaetzt-tests` (18×, geplant: `slice-harness-schaetzung-planner-faktor`) und `groesse-erst-in-der-pruefung-des-architect-sichtbar` (4×, dieselbe Adresse) — §4 nennt deshalb Faktor, Schätzspanne und Meldeschwelle; `verhalten-nur-unter-linux-geprueft` (3×, verkörpert in §3.12) — Randform (k). Keiner dieser Einträge erreicht mit diesem Slice eine Schwelle neu (die vier verkörperten zählen erst wieder, wenn die Regel trotz ihrer Verkörperung versagt).

**Zählung** (Planner, 2026-10-21, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 2** (1: die Option mit ihren Fällen; 2: Hilfetext, Handbuch und `README.md`). **Schichten: zwei** — der CLI-Adapter des Hexagons mit seinen Tests (`internal/adapters/driving/cli/`, dazu der E2E-Fall unter `test/integration/`, der dieselbe Fläche prüft) und die Nutzer-Doku (`docs/user/`, `README.md`); die Spezifikation ist eine **eigene** Schicht nach der Schichtteilung aus `slice-harness-d-check-v0-85` §8, liegt aber vor dem Start und damit außerhalb des Slice (§3, §4). Zählt sie der Architect doch zum Slice — oder zählt er den E2E-Fall als eigene Schicht —, sind es drei, und der vorgemerkte Schnitt aus §4 (`slice-version-option-doku`) greift. Der Slice trägt keine Sendung in einen anderen Plan ein und nimmt keine an; keine Nachzählung bei einem Nehmer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und Entscheidung gehen dem Code voraus; der Architect legt die Randformen vor dem ersten Code-Commit fest).

Nachgezählt beim Eintragen der Reihenfolge (Planner, 2026-10-22, `AGENTS.md` §3.13, Geber `slice-harness-regeln-rueckbau`): Nur der Start in §4 ist geändert (Reihenfolge, Hinweis auf den Bestand von Marke, Meldeschwelle und Schätzung); Gegenstand und DoD unverändert, Liefer-Punkte weiter 2, Schichten weiter zwei.
