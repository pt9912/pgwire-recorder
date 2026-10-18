# Slice slice-harness-gegenprobe-baumkopie: Gegenproben kopieren den Arbeitsbaum nach `git ls-files` statt mit `tar` und `cp -R`, durch einen gemeinsamen Helfer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt am 2026-10-18 nach dem Fund des Nutzers vom selben Tag (ein voller Lauf von `make gates` brach in `a-check-negativ` mit `tar: ./.harness/state/spans: Datei hat sich beim Lesen geändert` ab, weil die Hooks laufender Agenten dort schreiben; der Vorschlag, die Kopie nach `git ls-files` aufzubauen, ist ein Vorschlag des Nutzers, die Randformen in §6 sind Empfehlungen des Planners). Der Slice liegt in `open/`: Er ist nicht priorisiert, und Vertrag und ADR des Architect liegen nicht vor (§4 *Start*).

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als zwei der drei Gegenproben dessen Messung 1 und 3 tragen (Abdeckungs-Deklaration im Kopf der Gegenproben); kein Anspruch des Produkts, der Slice ändert Werkzeuge des Harness. Bindung an das Register: `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten` (neu, 1×, Beleg `evidence/slice-harness-devcontainer-einrichten-geruest.md`). Keine aktive ADR; ob eine entsteht, entscheidet der Architect (§4 *Start*).

**Berührte Spec-Stellen:** neuer Vertrag des Helfers und seiner Gegenprobe in `harness/vertraege/` (Nummer vergibt der Architect; nächste freie Nummer nach `SPEC-058`, also 059; `SPEC-058` ist durch `container-start` vergeben (Zweig architekt-geruest, noch nicht auf main); der Architect vergibt sie); Ortsangaben „Kopie des Arbeitsbaums“ in [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) (Punkt 13, Absatz zu `lint-gegenprobe`, `coverage-gegenprobe`, `a-check-negativ`) und [`SPEC-056`](../../../../harness/vertraege/SPEC-056-testabdeckung-schwelle.md) (Punkt 11); [`SPEC-049`](../../../../harness/vertraege/SPEC-049-lint-profil.md) (Gegenprobe des Lint-Gates in Kopien des Arbeitsbaums, Grenze des `touch`).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein voller Lauf von `make gates` bricht nicht mehr daran, dass eine Datei im Arbeitsbaum sich während der Kopie ändert: Ein Helfer `tools/harness/baumkopie.sh` kopiert den Arbeitsbaum nach `git ls-files` (versionierte und unversionierte, nicht ignorierte Dateien des Arbeitsbaums, ohne `.git`, ohne `.harness/state/` und andere ignorierte Laufzeitdateien), und die drei Gegenproben, die den Arbeitsbaum kopieren (`a-check-negativ`, `lint-gegenprobe`, `coverage-gegenprobe`), rufen ihn statt `tar` und `cp -R`.

**Betroffen (Zählung, geprüft am Stand `1c49c56f`):** Den Arbeitsbaum als Ganzes kopieren genau drei Skripte: `tools/arch/a-check-negativ.sh` (`tar`, eine Kopie je Fall), `tools/harness/lint-gegenprobe.sh` (Funktion `kopie`, `cp -R` über `*` und `.[!.]*` ohne `.git`) und `tools/harness/coverage-gegenprobe.sh` (Funktion `kopie`, ebenso, dazu ohne `README.md` und `docs/user/`). Nur das erste brach nachweislich (`tar` meldet die Änderung); die beiden anderen brechen nicht, nehmen aber einen halb geschriebenen Stand der Laufzeitdateien mit, derselbe Mangel. **Nicht betroffen:** `stop-hook-gegenprobe` (`cp -a` und `cp` auf ein Temp-Repo, das das Skript selbst anlegt), `gates-schnell-gegenprobe` (ebenso, `cp -r` eines eigenen Temp-Baums), `meldungskatalog-check-gegenprobe`, `kopf-check-gegenprobe`, `container-gegenprobe`, `postgres-versionen-gegenprobe`, `integration-volumes-gegenprobe` (kopieren einzelne benannte Dateien in selbst angelegte Temp-Bäume) und `abdeckung-gegenprobe` (baut ihren Baum selbst).

**Gemeinsamer Helfer, nicht drei Kopien.** Drei Skripte mit derselben Randform-Liste (§6) wären drei Stellen für jede Entscheidung; der Helfer trägt sie einmal, hat einen Vertrag und eine Gegenprobe. Die Aufrufer behalten ihre eigenen Ausschlüsse (`docs/user/`, `README.md`) und ihr `touch` nach dem Schreiben einer Mutation; der Helfer kennt sie nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Klon (`git clone .` oder von `origin`) als Quelle der Kopie — der Klon liefert nur den committeten Stand, die Gegenprobe soll aber den Arbeitsbaum prüfen, wie er ist (unversionierte neue Dateien, Änderungen vor dem Commit); `origin` bräuchte Netz und einen gepushten Stand (Vorschlag des Nutzers vom 2026-10-18).
- Eine Momentaufnahme des Arbeitsbaums (Dateisystem-Snapshot, Sperre gegen Schreiber) — anderer Vorgang und über die Host-Voraussetzungen (`bash`, `git`, `docker`) hinaus; der Helfer zeigt nur, was er zum Zeitpunkt des Lesens fand, und sagt das (§6, *Datei ändert sich beim Lesen*).
- Eine kleinere Kopie, die nur die Pfade enthält, die das jeweilige Gate liest (etwa nur `**/*.go` für `a-check-negativ`, nur den Build-Kontext nach `.dockerignore` für Lint und Abdeckung) — anderer Vorgang (Laufzeit der Gates statt Beständigkeit des Laufs, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`) und eine eigene Aussage je Gate; die Kopie nach `git ls-files` ist die Obermenge, auf der ein späterer Slice einschränken kann. Eine Adresse gibt es nicht; ein Slice legt der Planner an, wenn der Nutzer es entscheidet.
- Die Gegenproben, die selbst angelegte Temp-Bäume kopieren (Liste oben) — Bestand bleibt bewusst stehen: Sie lesen den Arbeitsbaum nicht als Ganzes und können an dieser Ursache nicht brechen; der Helfer wäre dort kein Gewinn.
- Das Verhalten der Hooks laufender Agenten (Spans unter `.harness/state/`) und der Stop-Hook — Schicht-Abgrenzung: Die Ursache liegt im Schreiben der Spans, die Wirkung im Kopieren; behoben wird die Wirkung, und `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` bleibt eigene Beobachtung mit eigener Gegenmaßnahme (`slice-harness-devcontainer-einrichten`).
- Eine Änderung am Build-Kontext (`.dockerignore`) oder an `.a-check.yml` — Schicht-Abgrenzung: Der Slice ändert nur, wie die Kopie entsteht, nicht, was ein Gate daraus liest.
- Handbuch, README, `docs/developer/` — Schicht-Abgrenzung und Größe: kein Verhalten des gebauten Binary; der Vertrag und der Gate-Index in `harness/README.md` tragen die Aussage.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** `tools/harness/baumkopie.sh <quelle> <ziel>` legt in `<ziel>` eine Kopie der Dateien an, die `git ls-files` des Arbeitsbaums von `<quelle>` nennt (versioniert und unversioniert-nicht-ignoriert, Inhalt des Arbeitsbaums, nicht des Index), und endet bei jeder Randform aus §6 so, wie der Vertrag es entscheidet (Exit-Code und Meldung); der Vertrag steht in `harness/vertraege/` vor dem ersten Code-Commit.
- [ ] **Liefer-Punkt 2:** `make baumkopie-gegenprobe` (Fragment unter `harness/mk/`, an `GATE_CHECKS`, Zeile im Gate-Index von `harness/README.md` §Sensors, Eintrag in der Ausnahmeliste und in `gates-schnell-gegenprobe`, soweit der Vertrag es verlangt) trägt je Randform und je Element ihrer Aufzählungen aus §6 einen Fall und eine Mutation des Helfers, die nur ihn rot färbt (`AGENTS.md` §3.10); die unveränderte Kopie ist grün; ein Fall mit einer unter `.harness/state/` laufend schreibenden Datei zeigt, dass die Kopie nicht abbricht und die Datei nicht enthält.
- [ ] **Liefer-Punkt 3:** `a-check-negativ`, `lint-gegenprobe` und `coverage-gegenprobe` rufen den Helfer statt `tar` und `cp -R`, ihr übriges Verhalten bleibt (Ausschlüsse der Aufrufer, `touch`, Aufräumen); die Ortsangaben in `SPEC-052` Punkt 13, `SPEC-056` Punkt 11 und, soweit betroffen, `SPEC-049` und den Sensors-Zeilen von `harness/README.md` sagen „Kopie nach `git ls-files`“ statt „Kopie des Arbeitsbaums“ und tragen die Kennung des Vertrags des Helfers. Dass die Gegenproben den Helfer rufen (und nicht mehr `tar`/`cp -R` über den Baum), prüft ein Fall in `baumkopie-gegenprobe` je Skript, sonst steht der Satz nicht in einem Kommentar (`AGENTS.md` §3.11).
- [ ] `make gates` grün, **ohne** dass ein laufender Hook das Ergebnis verändert: ein voller Lauf, während unter `.harness/state/spans` geschrieben wird (der Implementer stellt das mit einer Schleife her und nennt es im Bericht).
- [ ] Doku-Update: `harness/README.md` §Sensors (Zeile des neuen Gates, Zeilen der drei Gegenproben).
- [ ] Beobachtungs-Register: `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten` mit dem Beleg dieses Slice geführt (liegt an, 1×); kein Zähler wird gesetzt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-<NNN>-baumkopie.md` | neu (Architect, vor dem Code) | Vertrag des Helfers: Aufruf, Umfang der Kopie, Entscheidung je Randform aus §6, Exit-Codes, Grenze (keine Momentaufnahme) |
| `tools/harness/baumkopie.sh` | neu | Helfer: `git -C <quelle> ls-files -z` mit `--cached --others --exclude-standard`, Zusammenführen doppelter Einträge, Kopie ohne Folgen von Symlinks mit erhaltenem Modus, Zielprüfung |
| `tools/harness/baumkopie-gegenprobe.sh`, `harness/mk/baumkopie-gegenprobe.mk` | neu | Gegenprobe in Temp-Repos, die sie selbst anlegt (kein Docker): je Randform ein Fall und eine Mutation (Liefer-Punkt 2) |
| `tools/arch/a-check-negativ.sh`, `tools/harness/lint-gegenprobe.sh`, `tools/harness/coverage-gegenprobe.sh` | update | Aufruf des Helfers statt `tar` und `cp -R`; Kommentarköpfe nennen den Helfer als Quelle der Kopie (Liefer-Punkt 3) |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | update (soweit der Vertrag das Gate in die Ausnahmeliste stellt) | die Liste im Fragment und die Liste der Gegenprobe müssen sich decken |
| `harness/vertraege/SPEC-052-gate-kette.md`, `SPEC-056-testabdeckung-schwelle.md`, ggf. `SPEC-049-lint-profil.md`, `harness/README.md` | update | Ortsangaben nachgezogen (Liefer-Punkt 3) |

## 4. Trigger

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. **Vorbedingung des Starts: Der Architect hat den Vertrag in `harness/vertraege/` geschrieben, mit der Entscheidung je Randform aus §6 (`AGENTS.md` §3.12), und entscheidet, ob eine ADR entsteht** (Wahl der Quelle `git ls-files` gegen Klon und Snapshot; Gate-Name und Stelle in der Ausnahmeliste von `SPEC-052` Punkt 2 und 4). Vertrag und ADR sind Commits des Architect vor dem Move nach `in-progress/`. Reihenfolge: keine Abhängigkeit von einem anderen Slice. `slice-harness-lint-warten-ohne-frist` (`open/`) ändert ebenfalls `tools/harness/lint-gegenprobe.sh`; welcher zuerst startet, ist entschieden (Nutzer, 2026-10-18, Neuschnitt der Wellen 2 und 3): dieser Slice zuerst, `slice-harness-lint-warten-ohne-frist` rebased (dort §4 und §6). Beide stehen in welle-gegenproben-und-mutanten, dieser an erster Stelle (§5 dort); beide ändern `tools/harness/lint-gegenprobe.sh`, und dieser behebt die Flake der Läufe.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Planner, 2026-10-18; Metrik wie oben; Gegenprobe mit mindestens dem Fünffachen des Helfers, weil eine Tabelle von Randformen mit einer Mutation je Element das Gewicht trägt, Beleg `slice-v1-abschluss-protokollrand` rund 8, Entwurf in `slice-harness-schaetzung-planner-faktor`, noch nicht in `AGENTS.md` §3.14):**

| Teil | Zeilen |
|---|---|
| Helfer `baumkopie.sh` (Aufruf, Liste, Kopie, Zielprüfung, Exit-Codes) | 60 bis 90 |
| Umstellung der drei Skripte (Einfügungen und Löschungen, mit Kopfkommentaren) | 30 bis 50 |
| Fragment, Ausnahmeliste, `gates-schnell-gegenprobe` | 20 bis 35 |
| Gegenprobe (14 Randformen mit Mutation, Fünffaches des Helfers) | 300 bis 450 |
| Ortsangaben in den Verträgen und in `harness/README.md` (lange Zeilen zählen je einmal ein und aus) | 15 bis 30 |
| **Zusammen** | **425 bis 655** |

Meldeschwelle 700, **Marke 900**. Die Mitte (540) liegt bei 60 % der Marke, die Obergrenze (655) bei 73 %, unter der Meldeschwelle; ein Schnitt unterbleibt. Die Spanne hängt am Faktor der Gegenprobe: Zählt der Architect mehr als 14 Randformen, rechnet er neu.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Helfer mit Gegenprobe getrennt von der Umstellung der drei Skripte).
- `in-progress` → `open` (blockiert): Der Vertrag entscheidet eine Randform aus §6 so, dass eine der drei Gegenproben in der Kopie nach `git ls-files` rot wird, weil sie eine ignorierte Datei braucht (Befund an den Planner).

## 5. Closure-Trigger

DoD vollständig, ein voller `make gates` grün **während** unter `.harness/state/spans` geschrieben wird (Beleg im Bericht), Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben, die drei Paarungen getragen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des neuen Vertrags (`AGENTS.md` §3.12).** Je Randform steht die Empfehlung des Planners; entschieden ist keine, der Architect legt sie vor dem Code in den Vertrag (`harness/vertraege/`), und der Implementer entscheidet keine still (Randform-Rückgabe).

1. **Unversionierte, nicht ignorierte Datei** — kopiert (Vorschlag des Nutzers: der Stand im Arbeitsbaum). Entschieden: offen (Vertrag).
2. **Ignorierte Datei** (`.harness/state/`, `dist/`) — nicht kopiert. Eine Datei, die ein Gate aus einer ignorierten Quelle liest, fehlt dann in der Kopie; der Slice prüft das für die drei Gegenproben (Risiko 1). Entschieden: offen.
3. **Versionierte Datei im Arbeitsbaum gelöscht, nicht gestaged** — `git ls-files` nennt sie noch; Empfehlung: nicht kopiert, weil der Arbeitsbaum zählt, ohne Fehler. Gestaged gelöscht: nicht im Index, nicht genannt. Entschieden: offen.
4. **Versionierte Datei geändert** — Inhalt des Arbeitsbaums, nicht des Index (auch gestaged-nicht-committet). Entschieden: offen.
5. **Symlink** — als Symlink kopiert, nie gefolgt; ein hängender Symlink und einer, der aus dem Baum zeigt, bleiben unverändert (Ziel nicht geprüft, nicht gelesen). Entschieden: offen.
6. **Dateiname mit Leerzeichen, Tabulator, Zeilenumbruch, Anführungszeichen oder Nicht-UTF-8** — je ein Fall; Empfehlung: Liste mit NUL getrennt (`-z`) gelesen und weitergereicht, nie als Zeilen. Entschieden: offen.
7. **Gitlink (Submodul) und eingebettetes Repository im Baum** — Empfehlung: Exit ungleich 0 mit der Nennung des Pfads, kein stilles Auslassen; dieses Repo hat keines. **Verknüpfter Worktree als Quelle** (`.git` ist eine Datei, wie in diesem Zweig) — arbeitet wie ein gewöhnlicher Baum. Entschieden: offen.
8. **Quelle kein Git-Arbeitsbaum** (kein `.git`, `safe.directory` verweigert `git`) — Empfehlung: Exit 2 mit Meldung, **kein Rückfall auf `tar`**, sonst wechselt die Semantik still. Entschieden: offen.
9. **Datei verschwindet zwischen `ls-files` und Kopie** — Empfehlung: übersprungen und gemeldet, Exit 0, wenn sie dann fehlt; jeder andere Lesefehler (Berechtigung) Exit ungleich 0. Entschieden: offen.
10. **Datei ändert sich beim Lesen** (versionierte oder unversionierte, nicht ignorierte Datei, die ein Agent bearbeitet) — GNU-`tar` endet dafür mit 1 (`--warning=no-file-changed` kennt `bsdtar` nicht), `cp` nicht; die Kopie ist keine Momentaufnahme. Empfehlung: kein Fehler, der Vertrag sagt es als Grenze; Wahl des Kopierwerkzeugs (GNU gegen BSD, die Host-Voraussetzung nennt nur `bash`, `git`, `docker`, `make`) ist Sache des Architect. Entschieden: offen.
11. **Modus-Bits und Zeiten** — ausführbar bleibt ausführbar; die Zeit der Kopie ist die des Kopierens (nicht erhalten), weil der Build-Kontext Mutationen über neuere Zeiten überträgt (`SPEC-049`, Grenze des `touch`). Entschieden: offen.
12. **Ziel** — fehlt: angelegt; existiert und ist nicht leer: Exit ungleich 0; liegt innerhalb der Quelle: Exit ungleich 0 (sonst checkt `git add -A` die Kopie ein, `BEO-REPO/mutanten-skript-checkt-kopien-in-den-arbeitsbaum-ein`); Pfad mit Leerzeichen. Je ein Fall. Entschieden: offen.
13. **Aufruf** — aus einem Unterverzeichnis der Quelle; mit `GIT_DIR`/`GIT_INDEX_FILE` aus der Umgebung (gesetzt in Git-Hooks); mit `sparse-checkout` oder `skip-worktree` (wie Punkt 3); Merge-Konflikt im Index (Pfad mehrfach genannt, Empfehlung: einmal kopiert). Je Aufruf ein Fall. Entschieden: offen.
14. **Leeres Verzeichnis** — Git kennt es nicht, es entsteht nicht in der Kopie; ein Gate, das eines braucht, legt es selbst an. Entschieden: offen.

**Aufrufarten des Treibers (§3.12):** Der Helfer liefert keinen Stand, keinen Nachweis und keinen Exit, den ein Hook oder ein Gate später glaubt; `make -n`, `-q`, `-t`, `-i` und `MAKEFLAGS` betreffen das Rezept des Gates `baumkopie-gegenprobe` wie das jeder Gegenprobe (`SPEC-052` Punkt 8), nicht den Helfer. **Aufruf-Voraussetzung (§3.12):** Der Helfer sagt keine Ordnung über mehrere Aufrufe zu.

- Die Kopie nach `git ls-files` lässt eine Datei weg, die ein Gate aus einer ignorierten Quelle liest, und die Gegenprobe wird dadurch rot oder, schlimmer, grün aus dem falschen Grund — **Ausgang (vorgesehen):** entfallen, wenn die unveränderte Kopie in allen drei Gegenproben grün ist und je Gegenprobe eine Mutation weiter rot wird (Beleg: ein voller `make gates`); sonst eingetreten → Rückführung nach `open` (§4).
- GNU-`tar` bricht auch bei einer versionierten Datei, die ein Agent bearbeitet; die Kopie nach `ls-files` behebt nur die Laufzeitdateien — **Ausgang (vorgesehen):** entfallen, wenn der Vertrag Randform 10 entscheidet und ein Fall sie zeigt; sonst weiter offen → `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten`.
- Die Kopie ist keine Momentaufnahme: Eine Gegenprobe kann einen Stand mit halb geschriebener Datei mutieren und grün oder rot werden aus einem Grund außerhalb der Mutation — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten` (Grenze im Vertrag; ein Lauf in einem Dev-Container, `slice-harness-devcontainer-einrichten`, vermeidet die Schreiber, ohne dass der Slice darauf baut).
- Der Slice liefert drei Gegenproben und einen Helfer samt Gate und kommt mit 14 Randformen nahe an die Obergrenze der Schätzung — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 700 bleibt; darüber meldet der Implementer (Entscheidung „weiter“ oder „schneiden“ als datierte Zeile in §7).
- `slice-harness-lint-warten-ohne-frist` ändert dieselbe Datei `tools/harness/lint-gegenprobe.sh` — die Reihenfolge ist entschieden (dieser Slice zuerst, Nutzer 2026-10-18) und steht in dessen §4 und §6, die Bindung ist keine Sendung (kein Gegenstand wandert), deshalb steht nichts in dessen §1 oder DoD. **Ausgang (vorgesehen):** entfallen, sobald dieser Slice in `done/` liegt und `slice-harness-lint-warten-ohne-frist` auf dessen Stand rebased hat.

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten/` neu angelegt am 2026-10-18 (Planner), Beleg `evidence/slice-harness-devcontainer-einrichten-geruest.md`, Zähler 1×
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (Liefer-Punkt 2 trägt je Randform eine Mutation), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen vor dem Code), `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (die Kopie wird kleiner, nicht größer; eine kleinere Kopie ist ausgeschlossen, §1), `BEO-REPO/mutanten-skript-checkt-kopien-in-den-arbeitsbaum-ein` (Randform 12), `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (abgegrenzt, §1), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (§4 rechnet mit dem Fünffachen), `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 nennt die Metrik). Neu: `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten` (1×, mit diesem Plan angelegt). Keiner der Treffer erreicht mit diesem Slice 3×, den Eintrag zur Schätzung ausgenommen, den die Closure zählt.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Liefer-Punkte 3 (Helfer mit Vertrag, Gegenprobe mit Gate, Umstellung der drei Skripte mit Nachzug der Ortsangaben). Schichten 1 (Harness-Werkzeug mit Gegenprobe und Gate-Index; Planung und Vertrag sind keine Schicht; kein Produkt-Code, keine Nutzer- oder Wartungs-Doku). Schätzung in der Metrik von §4: 425 bis 655, **unter der Marke 900** und unter der Meldeschwelle 700. Der Plan trägt keine Sendung in einen anderen Slice ein; es gibt keinen Nehmer, dessen Zählung nachzuziehen wäre.

**Nachgezählt** (`AGENTS.md` §3.9, §3.13, Planner, 2026-10-18, im Commit des Neuschnitts der Wellen 2 und 3): Die Stelle in der Reihenfolge ist eingetragen (§4, §6); Gegenstand, Liefer-Punkte (3) und Schichten (1) ändern sich nicht, die Schätzung bleibt 425 bis 655, unter der Marke 900.
