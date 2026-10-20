# Slice slice-harness-gegenproben-schneller: Die Gegenproben laufen deutlich schneller und räumen ihren Arbeitsordner je Mutation auf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-20). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist, und keine Anforderung der Welle `welle-v1-nachtrag` hängt an ihm; würde er in sie aufgenommen, hinge deren Closure-Trigger an einem Werkzeug-Slice, den kein Slice der Welle voraussetzt. Angelegt bei der Closure von `slice-v1-abschluss-binary-namen` (Lerneintrag, benannte Lücke; Beobachtung des Verifiers F-V2 und des Reviews F-2) und nach dem Hinweis des Nutzers vom 2026-10-20, die Ursache der langen Läufe seien verschachtelte Schleifen mit vielen `sed`- und `awk`-Aufrufen; das ist die Vermutung des Nutzers, belegt ist sie nicht (§6). Der Slice liegt in `open/`: Er ist nicht priorisiert, und was der Architect nach dem Lesen der Skripte entscheidet, liegt nicht vor (§4 *Start*).

**Bezug:** kein Lastenheft-Bezug — der Slice ändert Werkzeuge des Harness, keinen Anspruch des Produkts. Bindung an das Register: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (die Größe „freier Platz des Datenträgers“, Beleg `evidence/slice-v1-abschluss-binary-namen.md`) und `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Laufzeit der Gegenproben). Entscheidung: [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenproben stehen nicht in `GATE_CHECKS`; hier unverändert, §1).

**Berührte Spec-Stellen:** `SPEC-067` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md); `make release-gegenprobe`) · `SPEC-055` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md); `make test-container-gegenprobe`) · `SPEC-060` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-060-gegenproben.md); Regeln für Gegenproben). Ob und wo der Vertrag die Zusage „der Arbeitsordner einer Mutation ist nach ihr entfernt“ trägt, entscheidet der Architect vor dem Code (§6, Rückgabe 1).

**Verantwortlich:** —
**Autor:** Planner. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Gegenproben `tools/harness/release-gegenprobe.sh` (2045 Zeilen, 14 Zeilen mit `sed`, 10 mit `awk`, 59 Schleifen) und `tools/harness/container-gegenprobe.sh` (790 Zeilen) laufen deutlich schneller, bei gleichem Ergebnis, und der Arbeitsordner einer Mutation wird nach ihr entfernt statt am Ende des Laufs. Bisher wächst der Arbeitsordner von `make release-gegenprobe` um rund 40 MB je Mutation, über 5 GB bei 535 Mutationen, und wird erst am Ende freigegeben; ein voller Lauf braucht rund 6 GB und scheitert auf einem Host mit weniger freiem Platz als scheinbarer Befund (`slice-v1-abschluss-binary-namen` §7, Verifikation F-V1 und F-V2, Review F-2: sechs Mutationen `m-p16-*` rot im Lauf, allein gefahren grün). **„Deutlich schneller“ ist keine Zahl:** Der Architect legt nach dem Lesen der Skripte fest, was gemessen wird und um wieviel die Zeit sinken muss (§6, Rückgabe 2); ohne diese Zahl ist die DoD nicht prüfbar.

**Der Slice beginnt mit dem Lesen der Skripte durch den Architect, nicht mit einer Messreihe.** Die Ursache der Laufzeit ist nicht belegt: Der Nutzer vermutet verschachtelte Schleifen mit vielen Aufrufen von `sed` und `awk`; ebenso möglich sind Kopien des Arbeitsbaums, `docker`-Aufrufe der Attrappen, das Entpacken von Archiven je Fall oder `git`-Repositories je Fall. Der Architect liest, schlägt den Schnitt vor und misst nur, was er zur Entscheidung braucht (eine Auswahl `GP_MUT`, kein Gesamtlauf).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung dessen, was die Gegenproben prüfen: Ergebnis, Fälle, Mutationen und die Mengen der roten Fälle je Mutation bleiben gleich (Klasse 3, anderer Vorgang). Grund: Der Slice ändert, wie schnell geprüft wird, nicht, was; eine gestrichene oder zusammengelegte Mutation senkte die Prüfung, nicht die Laufzeit (`AGENTS.md` §3.6, Gates nicht ohne ADR lockern).
- Ein Nachtlauf der Gegenproben in der CI — ein anderer Vorgang (Klasse 3). Grund: Er wäre eine neue Außenwirkung des Harness und braucht eine ADR; die Slices von `welle-ci-nachtlaeufe` und `slice-harness-ci-gegenproben` (`open/`) tragen ihn, dieser Slice nimmt ihn nicht an.
- Die übrigen Gegenproben (`make hook-gegenprobe`, `make lint-gegenprobe`, `make workflows-check-gegenprobe` und die weiteren der Liste `GEGENPROBEN_CONTAINER`) — Bestand bleibt bewusst stehen (Klasse 2). Grund: Gemessen und gelesen sind die zwei größten; ob eine andere dieselbe Form trägt, zeigt das Lesen des Architect, und was er dort findet, geht in einen Folge-Slice, nicht in diesen.
- Die Aufnahme einer Gegenprobe in `GATE_CHECKS` oder `make gates` — Klasse 3 (`ADR-0048`). Grund: Schneller macht sie nicht zum Gate der Kette.
- Produkt-Code (`cmd/`, `internal/`) und `Dockerfile` — Schicht-Abgrenzung (Klasse 4): Der Slice berührt allein Skripte unter `tools/harness/`, ihre Gegenproben und die Zeile in `harness/README.md`.
- Der Gegenstand `slice-harness-gegenprobe-baumkopie` (Kopie des Arbeitsbaums nach `git ls-files` durch einen gemeinsamen Helfer, `open/`) — Klasse 1 des anderen Slice, nicht dieses. Grund: Er trägt die Kopie des Baums für `a-check-negativ` und weitere; berührt der Architect hier dieselbe Stelle, nennt er das und gibt es an den Planner.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] `make release-gegenprobe` läuft in der Zeit, die der Architect nach dem Lesen festlegt (Messung und Zahl in §3 und im Vertrag, §6 Rückgabe 2), mit unverändertem Ergebnis: die unveränderte Kopie grün, je Mutation genau die benannten Fälle rot (`GP_ENTDECKE=1`, Vorher-Nachher-Vergleich der roten Mengen je Mutation).
- [ ] `make test-container-gegenprobe` ebenso: Zeit nach Festlegung, Ergebnis unverändert.
- [ ] Der Arbeitsordner je Mutation wird nach der Mutation entfernt (nach einer grünen und nach einer roten, bei Abbruch durch `SIGINT` und `SIGTERM` nach den Randformen des Architect); der Platzbedarf eines vollen Laufs von `make release-gegenprobe` liegt unter einer Grenze, die der Architect festlegt, und der volle Lauf in einem Zug ist am Endstand belegt (`slice-v1-abschluss-binary-namen`, F-V2).
- [ ] `make gates` grün; die Gegenproben, die der Slice ändert, am Endstand in einem Zug (`AGENTS.md` §6); die übrigen Gegenproben der Listen bleiben unberührt (§1).
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

Der Architect füllt diese Tabelle nach dem Lesen der Skripte (Start, §4); der Planner schreibt hier nur die Kandidaten.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/release-gegenprobe.sh` | update | Laufzeit; Arbeitsordner je Mutation entfernen |
| `tools/harness/container-gegenprobe.sh` | update | Laufzeit; Arbeitsordner je Mutation entfernen |
| `harness/README.md` (§Sensors, Zeilen der beiden Gegenproben) | update, falls der Aufruf oder die Zusage sich ändert | nennt, was die Gegenprobe zusagt |
| `harness/vertraege/` | durch den Architect vor dem Code | Zusage zum Arbeitsordner und Zielzahl (§6, Rückgaben) |

Nicht in dieser Liste: Produkt-Code, `Dockerfile`, andere Gegenproben, CI (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next`): Der Nutzer priorisiert den Slice; Reihenfolge laut Auftrag vom 2026-10-20: nach `slice-v1-abschluss-homebrew` (der Slice steht beim Schreiben in `next/`); eine Stellung zu `slice-v1-abschluss-release-v0-1-1` und `slice-erster-release-homebrew-nachweis` ist nicht entschieden, keiner der beiden setzt diesen Slice voraus. Der erste Schritt nach `next` → `in-progress` ist das Lesen der zwei Skripte durch den Architect (kein langes Messen) und das Schreiben von Zielzahl, Randformen und Schnitt vor dem ersten Code-Commit (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness/README.md harness/vertraege` (Tests und Gegenproben eingeschlossen, `AGENTS.md` §3.14) **über 1000** liegt; `<Basis>` ist der Hash von `main` beim Übergang nach `in-progress/`, den der Architect in §3 einträgt. Die Schätzung des Architect in dieser Metrik steht in §3 (Faktor von Tabellen- und Gegenprobe-Code zum Produktcode rund zwei und mehr; hier ist der Code selbst Gegenprobe-Code, deshalb ohne Faktor, aber mit Zeilen, die ein Umbau verschiebt). Ist die Marke vor dem Code absehbar überschritten, schneidet der Architect in §3 nach den zwei Skripten (je ein Liefer-Punkt ist bereits getrennt, DoD 1 bis 3).
- `in-progress` → `next`: Das Lesen zeigt, dass die Laufzeit nicht an den Skripten hängt (etwa an `docker` oder dem Datenträger des Hosts); dann ist „schneller“ kein Gegenstand dieses Slice, und der Architect gibt die Ursache an den Planner zurück.
- `in-progress` → `open`: Die Zielzahl oder die Zusage zum Arbeitsordner verlangt eine ADR, die der Nutzer nicht annimmt.

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

**Randformen** (`AGENTS.md` §3.12; vom Architect vor dem ersten Code-Commit zu entscheiden und im Vertrag festzuhalten; der Planner nennt sie, entscheidet sie nicht):

| Randform | Frage | Entscheidung |
|---|---|---|
| (a) Arbeitsordner einer roten Mutation | Bleibt er für die Diagnose stehen (Pfad im Log) oder wird er auch dann entfernt? | offen, Architect |
| (b) Abbruch mitten in einer Mutation | `SIGINT` und `SIGTERM`: der Ordner der laufenden Mutation wird entfernt; die Fälle `p13-*` tragen das heute für den ganzen Lauf | offen, Architect |
| (c) Parallele Läufe (`GP_PARALLEL`) | Mehrere Mutationen zugleich: Platzbedarf und Entfernen je Ordner, kein gemeinsames Verzeichnis | offen, Architect |
| (d) Aufrufarten außerhalb des Normalfalls | `GP_MUT`, `GP_ENTDECKE`, `GP_FALL` und die Aufrufe unter `make` mit `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` (`AGENTS.md` §3.12, Aufrufarten des Treibers): unverändert, mit je einem Fall | offen, Architect |
| (e) Plattform | Gegenproben laufen im Gegenproben-Container unter Linux; macOS und Windows nicht zugesagt | offen, Architect (Grenze im Vertrag) |
| (f) Umgebung der Fahrt | freier Platz, Last, root, Locale und Zeitzone (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`): der Slice fährt die Gegenprobe vor der Übergabe einmal unter fremdem Benutzer, fremder Locale und unter Last | offen, Architect |

**Rückgaben an den Architect** (nach dem Lesen, vor dem Code):

1. Ort und Wortlaut der Zusage „Arbeitsordner je Mutation entfernt“ (`SPEC-060` oder die Verträge `SPEC-055` und `SPEC-067`) und die Randformen (a) bis (f).
2. Messgröße und Zielzahl für „deutlich schneller“ (etwa Zeit der unveränderten Kopie plus einer festen Auswahl `GP_MUT`, vorher und nachher am selben Stand, bei ruhigem Host; Last macht Zeiten unbrauchbar, `BEO-REPO/lastabhaengige-zaehlung`), und ob eine ADR nötig ist (kein Gate wird gelockert, `AGENTS.md` §3.6).
3. Der Schnitt, falls beide Skripte zusammen über der Marke liegen.

**Risiken:**

- Die Ursache der Laufzeit ist nicht belegt (Vermutung des Nutzers: verschachtelte Schleifen mit `sed` und `awk`) — **Ausgang:** offen bis Closure (entfallen: Ursache am Endstand gemessen und behoben; eingetreten: Ursache liegt außerhalb der Skripte, Folge-Slice mit Kennung).
- Das Ergebnis ändert sich unbemerkt: Ein schnellerer Fall färbt eine andere Menge roter Fälle oder bleibt grün, wo er rot sein müsste — **Ausgang:** offen bis Closure (Gegenmittel: Vorher-Nachher-Vergleich der roten Mengen aller Mutationen mit `GP_ENTDECKE=1`, einmal am Endstand in einem Zug; entfallen, wenn der Vergleich leer ist).
- Zeitmessungen schwanken mit der Last des Hosts — **Ausgang:** offen bis Closure (`BEO-REPO/lastabhaengige-zaehlung`: weiter offen, wenn eine Zahl der DoD an der Last hängt).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `gegenprobe-haengt-an-umgebung-des-hosts` (9×, Größe „freier Platz“ neu; der Ausgang liegt beim Nutzer, dieser Slice ist ein Teil der Antwort auf diese Größe, nicht die Zuweisung), `gate-laufzeit-waechst-mit-gegenprobe` (verkörpert; die Laufzeit der Gegenproben außerhalb von `make gates` ist hier der Gegenstand), `lastabhaengige-zaehlung` (Zeitmessung, §6), `gegenprobe-temp-praefix-geteilt` (Arbeitsordner, Präfix; der Architect liest den Eintrag vor dem Entfernen je Mutation). Der Stand der Zähler ist beim Start zu lesen; erreicht einer mit diesem Slice drei, braucht er einen Folge-Slice.

**Zählung** (Planner, 2026-10-20, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 3** (1: `release-gegenprobe.sh` schneller; 2: `container-gegenprobe.sh` schneller; 3: Arbeitsordner je Mutation entfernt, in beiden Skripten). **Schichten: eine** — Harness-Werkzeug (Skripte, Zeilen in `harness/README.md`); Verträge liegen vor dem Code beim Architect und zählen nicht als Liefer-Punkt. Der Slice steht an der Grenze von drei Punkten; der Architect zählt nach dem Lesen neu und schneidet in §3, wenn die Skripte nicht in einer Review-Sitzung prüfbar sind. Es sind noch keine Sendungen anderer Slices eingetragen. Größe: Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Verträge und Konventionen des Harness gehen den Skripten voraus).
