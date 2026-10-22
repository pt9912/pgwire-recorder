# Slice slice-harness-container-start-u-boot-0-13: Die gemessene Umgebung von `SPEC-058` nennt `u-boot` 0.13.0, belegt durch einen Lauf der Gruppe Docker am Host

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-21). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist. [welle-dev-container](../welle-dev-container.md) ist geprüft und **nicht** gewählt: Ihr Closure-Trigger verlangt alle ihre Slices in `done/` und den Aufbau eines frischen Scratch-Containers; ein sechster Slice, der auf einen Host mit `u-boot` 0.13.0 und einen 19-minütigen Lauf wartet, dehnte sie. Der Nutzer hat am 2026-10-21 die **Stelle** entschieden (hinter `slice-harness-devcontainer-einrichten`), nicht eine Welle.

**Bezug:** kein Anspruch des Produkts — der Slice ändert die gemessene Umgebung im Kopf eines Vertrags eines Harness-Werkzeugs und belegt sie mit einem Lauf. Entscheidungen: [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) (`container-start` ist ein Werkzeug im Repo), [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (die Gegenprobe ist ein Host-Ziel), [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Obergrenze von 600 s je Gegenprobe; Konflikt, §1 und §6), [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Metrik der Marke in §4; Status `Accepted`, `slice-harness-groessen-metrik` trägt sie in `AGENTS.md` §3.14 ein, bis dahin gilt sie aus der ADR). Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Laufzeit der Gegenprobe), `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (die Version am Host), `BEO-REPO/lastabhaengige-zaehlung` (Zeit und Zählung unter Last), `BEO-REPO/gegenprobe-temp-praefix-geteilt` (ein Lauf zur Zeit auf einem Daemon). Herkunft: Meldung in §6 von `slice-harness-devcontainer-einrichten` („die gemessene Umgebung dort nachzuziehen ist ein eigener Vorgang, Kennung vom Planner“); genannt, nicht geändert: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 15 (7) und Punkt 16 (Zustand `nicht gemessen`).

**Berührte Spec-Stellen:** [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) — Kopf (*Plattform*: die gemessene Umgebung) und, soweit der Architect es dort führt, Punkt 15 (Gruppe Docker). Punkt 3 (Schwelle `u-boot` ab 0.10.0 und seine Fälle `p3-*`) bleibt unverändert (§1).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Kopf von [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) nennt die Umgebung, in der `make container-start` und seine Gegenprobe **wirklich gemessen** sind: `u-boot` **0.13.0** statt 0.10.0, mit den Versionen von Docker und `devcontainer` desselben Laufs. Belegt ist der Satz durch einen grünen Lauf der **Gruppe Docker** von `make container-start-gegenprobe` am Host unter 0.13.0 — das ist die Gruppe, die das echte `u-boot` und das echte `devcontainer` ruft (`SPEC-058` Punkt 7, Punkt 15); ohne ihn wäre der Satz eine Behauptung (`AGENTS.md` §3.11). Die Schwelle aus Punkt 3 (`u-boot` ab 0.10.0) bleibt, und mit ihr jeder Fall `p3-*`.

**Anlass.** `u-boot` liegt am Host seit dem 2026-10-21 in 0.13.0 (Angabe des Nutzers), gemessen war 0.10.0, und ein Lauf dazwischen lief unter 0.11.0 (`slice-harness-devcontainer-einrichten-start`, `done/`, §7 Punkt 8). Die Meldung steht in §6 von `slice-harness-devcontainer-einrichten`: Für dessen Ziel folgt daraus nichts, der Nachzug von `SPEC-058` ist ein eigener Vorgang. Dieser Slice ist er.

**Konflikt mit [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zehn-Minuten-Grenze), ausdrücklich benannt:** `make container-start-gegenprobe` braucht für den vollen Lauf **19 min 38 s** (264 Mutationen, Gruppe Docker eingeschlossen; gemessen am 2026-10-18 am Host, Beleg `../observations/BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/evidence/slice-harness-devcontainer-einrichten-start.md`) und reißt damit die Obergrenze von 600 s um mehr als das Dreifache. Dieser Slice **hebt den Konflikt nicht auf** und senkt die Zeit nicht: Nach `SPEC-060` Punkt 16 (*Aufnahme*) und Punkt 15 (7) gilt die Grenze für eine Gegenprobe erst **ab ihrer Aufnahme**; `container-start-gegenprobe` steht dort im benannten Zustand `nicht gemessen` mit `slice-harness-gegenproben-grenze-rest` als Nehmer, und für diesen Zustand ist die Zeit nicht zugesagt, sondern wird gemessen und genannt. Dieser Slice fährt deshalb **nur die Gruppe Docker**, nicht den vollen Lauf — dieselbe Lesart, die [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) unter *Übergang* für eine nicht aufgenommene Gegenprobe gibt (engster Lauf beim Implementer, voller Lauf beim Verifier).

**Sendung an `slice-harness-gegenproben-grenze-rest`** (eingetragen vom Planner am 2026-10-21 im selben Commit, mit Nachzählung in dessen §8; `AGENTS.md` §3.13): Jener Slice misst und schneidet `container-start-gegenprobe` (dessen §1 nennt sie in der Liste der 18, seine DoD misst sie am Host). Er **nimmt an**: (a) die Zeit der Gruppe Docker unter 0.13.0, die dieser Slice misst und in §7 belegt, als Vormessung — sie ersetzt seinen vollen Lauf nicht; (b) die Zahl 19 min 38 s als gemessenen Stand über 600 s, aus dem folgt, dass er `container-start-gegenprobe` **schneidet** statt sie aufzunehmen; (c) dass seine Messung unter der Version läuft, die dieser Slice in den Kopf einträgt. `grep` nach `slice-harness-container-start-u-boot-0-13` findet die Sendung in dessen §1, §6 und §8.

**Der Vertragstext ist nicht Sache des Implementers.** Den Satz im Kopf von `SPEC-058` schreibt der **Architect** (`AGENTS.md` §3.12, Entscheidung des Nutzers: Vertragstext geht an den Architect zurück); der Implementer fährt den Lauf, legt die gemessenen Versionen und Zeiten vor und meldet jede Abweichung, statt den Kopf selbst umzuschreiben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Schwelle `u-boot` ab 0.10.0 in Punkt 3 und ihre Fälle** (`p3-*`: 0.9.0 abgelehnt, 0.10.0, 0.100.0 und 1.0.0 angenommen) — Bestand bleibt bewusst stehen (Klasse 2). Grund: 0.13.0 liegt über der Schwelle und wird von den vorhandenen Fällen gedeckt; eine Anhebung wäre eine Verschärfung der Voraussetzungen am Host mit eigener Entscheidung, und sie schlösse Hosts aus, auf denen das Werkzeug läuft.
- **Die Zeit von `container-start-gegenprobe` senken oder sie in die Tabelle von `SPEC-060` Punkt 16 aufnehmen** — `slice-harness-gegenproben-grenze-rest` (Klasse 1, `open/`; nimmt die Sendung oben in §1, §6 und §8 an). Grund: Aufnahme heißt gleiches Ergebnis in 600 s, und das ist Arbeit am Skript, nicht an der gemessenen Umgebung.
- **Der volle Lauf der Gegenprobe** (alle Gruppen, 264 Mutationen) — anderer Vorgang (Klasse 3): Er gehört dem Nehmer der Aufnahme und dem Verifier am Endstand ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) *Übergang*). Grund: Die Version berührt die Gruppe Docker (echtes `u-boot`, echtes `devcontainer`); die Gruppen mit Attrappen prüfen Formen, die von der Version nicht abhängen, und 19 Minuten für dieselbe Aussage sind kein Gewinn.
- **`docs/developer/devcontainer.md`** (§1 und §2 nennen 0.10.0 als Stand der Proben, §13 die Frage `podman`/`docker`) — Bestand bleibt bewusst stehen (Klasse 2). Grund: Die Sätze dort tragen mehr als dieser Lauf belegt (Features `sshd` und `claude-code`, Profil `sandbox`, Claude Code im Container, `make gates` im Klon); ein Satz mit 0.13.0 reichte weiter als die Prüfung (`AGENTS.md` §3.11). Das Dokument sagt „Stand der Proben“ und bleibt damit wahr; wer es nachzieht, fährt dessen Proben neu, und das ist ein eigener Vorgang (§6 nennt die Folge).
- **§6 und §9 von `docs/developer/devcontainer.md`** — `slice-harness-devcontainer-einrichten-doku` (`open/`, dessen §1 führt genau diese Abschnitte als eigenen Gegenstand; keine Sendung, nur eine Grenze). Grund: Zwei Slices an denselben Abschnitten erzeugen einen Konflikt statt eines Dokuments.
- **Produkt-Code, Gates und die übrigen Gegenproben** (`cmd/`, `internal/`, `tools/harness/*` außer dem Lauf, `harness/mk/`, `harness/README.md`) — Schicht-Abgrenzung (Klasse 4). Grund: Der Slice ändert einen Satz im Kopf eines Vertrags und fährt einen Lauf; kein Werkzeug verhält sich danach anders, und beim Review ist sofort prüfbar, was er nicht angefasst hat. Die Zeile in `harness/README.md` nennt die Schwelle (`u-boot` ab 0.10.0), nicht die gemessene Version, und bleibt deshalb unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkte (zwei):

- [ ] **Liefer-Punkt 1 — der Kopf von `SPEC-058` nennt die gemessene Umgebung des Laufs:** `u-boot` 0.13.0 und die Versionen von Docker und `devcontainer`, unter denen Liefer-Punkt 2 lief, mit dem Datum des Laufs; die Grenze darunter (macOS, Windows, Docker Desktop, andere Versionen von `devcontainer` nicht zugesagt) bleibt, ebenso Punkt 3 mit seinen Fällen. Geschrieben vom **Architect** (§1); der Satz nennt keine Version, die nicht gelaufen ist.
- [ ] **Liefer-Punkt 2 — der Lauf, der den Satz belegt:** `make container-start-gegenprobe` mit der **Gruppe Docker** am Host unter `u-boot` 0.13.0 ist grün; §7 nennt Wanduhr in Sekunden, Last (`/proc/loadavg`, erster Wert), die drei Werkzeug-Versionen (`u-boot --version`, `docker version`, `devcontainer --version`), den Hash des Stands und das Datum. Kein anderer Lauf eines Slice fährt gleichzeitig auf demselben Docker-Daemon (`SPEC-058` Punkt 15). Ein roter Fall ist ein Befund an den Architect (§4, §6), kein Haken.

Gates und Closure:

- [ ] `make docs-check` und `make kopf-check` grün (der Slice ändert ein Dokument unter `harness/vertraege/` und diesen Plan).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: Die gemessene Zeit der Gruppe Docker unter 0.13.0 ist ein Beleg in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (`evidence/slice-harness-container-start-u-boot-0-13.md`); weicht das Verhalten unter 0.13.0 von dem unter 0.10.0 ab, dazu ein Beleg in `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`. Keine Beobachtung ist ebenfalls eine Antwort und steht in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md), Kopf *Plattform* | update, **durch den Architect** | Liefer-Punkt 1: die gemessene Umgebung des Laufs (drei Versionen, Datum); Punkt 3 und seine Fälle unverändert |
| [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) Punkt 15, Gruppe Docker | update nur, wenn der Architect die Version dort führt | Liefer-Punkt 1: eine zweite Nennung derselben Umgebung driftet gegen den Kopf; der Architect entscheidet, ob sie dort steht oder auf den Kopf zeigt |
| Lauf `make container-start-gegenprobe` (Gruppe Docker) am Host | Lauf, keine Datei | Liefer-Punkt 2: der Beleg; die Zahlen stehen in §7 dieses Plans, nicht im Vertrag |
| `docs/plan/planning/open/slice-harness-gegenproben-grenze-rest.md` | update, **durch den Planner am 2026-10-21** | Sendung dieses Slice (`AGENTS.md` §3.13): Eintrag in §1 und §6, Nachzählung in dessen §8 — im Commit, der diesen Plan anlegt |
| `docs/plan/planning/observations/BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` | update bei der Closure | Beleg mit der gemessenen Zeit (kein Liefer-Punkt, DoD) |

- **Reihenfolge:** (1) der Architect schreibt den Satz **nicht** vor dem Lauf — gemessen wird zuerst, der Satz trägt die Zahlen des Laufs; (2) Lauf der Gruppe Docker am Host; (3) Kopf von `SPEC-058` mit den gemessenen Werten; (4) Closure mit Beleg im Register. Jeder Commit nennt die Kennung dieses Slice.
- **Werkzeuge am Host:** `docker`, `u-boot` 0.13.0, `devcontainer`, `node`, `ssh`, `ssh-keygen` (`SPEC-058` Punkt 15, *Voraussetzung am Host*). Fehlt eines, endet die Gegenprobe mit Exit ungleich 0 und nennt es; dann greift die Rückführung aus §4.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): `slice-harness-devcontainer-einrichten` liegt in `done/` (**Entscheidung des Nutzers vom 2026-10-21**: dieser Slice steht dahinter) — der Grund ist nicht nur die Reihenfolge: Jener Slice arbeitet am Host und im Container, und die Gruppe Docker darf nicht gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon fahren (`SPEC-058` Punkt 15, `BEO-REPO/gegenprobe-temp-praefix-geteilt`). Dazu: ein Host mit `u-boot` 0.13.0 und den Werkzeugen aus §3, WIP-Limit frei, und der Nutzer priorisiert den Slice.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) Punkt 1 (bis `slice-harness-groessen-metrik` sie in `AGENTS.md` §3.14 trägt, gilt sie aus der ADR):

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- harness/vertraege \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt (sein Hash steht beim Übergang hier). Die Pfade sind `harness/vertraege`; Plan, Register-Beleg und Berichte liegen unter den ausgeschlossenen Pfaden und zählen nicht.

**Marke 60**, **Meldeschwelle 40**. **Schätzung des Planners (2026-10-21): 4 bis 20** — der Satz im Kopf (zwei bis vier Zeilen, Löschungen gezählt) und höchstens eine Zeile in Punkt 15. Faktor Tests zu Produktcode entfällt: kein Code, kein Test, keine neue Mutation (§6). **Die Metrik misst den Aufwand dieses Slice nicht** — der liegt im Lauf am Host (19 Minuten voll, die Gruppe Docker gemessen in Liefer-Punkt 2) und in Zeilen nicht sichtbar; das ist das akzeptierte Negativ von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (*Konsequenzen*), und begrenzt wird der Slice hier durch Liefer-Punkte und Schichten (§8), nicht durch die Zahl.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt **über 60** — das heißt hier, dass der Architect mehr als die gemessene Umgebung ändert (etwa die Schwelle in Punkt 3, Fälle oder Zusagen der Gruppe Docker). Dann schneidet der Planner: die gemessene Umgebung bleibt hier, die Änderung der Zusagen wird ein eigener Slice.
- `in-progress` → `open` (blockiert): Kein Host mit `u-boot` 0.13.0 und den Werkzeugen aus §3 steht zur Verfügung, oder die Gruppe Docker ist unter 0.13.0 **rot** (dann bleibt der Kopf bei der alten Version, und der Befund geht an den Architect: Er entscheidet, ob `SPEC-058` eine Obergrenze oder eine Grenze je Version braucht, `AGENTS.md` §3.12), oder ein anderer Slice fährt auf demselben Daemon.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. Zwei beobachtbare Kriterien: (1) `grep -n "0\.10\.0" harness/vertraege/SPEC-058-container-start.md` findet die Zeichenfolge nur noch in Punkt 3 und seinen Fällen (Schwelle), nicht mehr im Kopf; (2) §7 nennt für den Lauf der Gruppe Docker Ausgang, Zeit, Last, die drei Werkzeug-Versionen und den Hash, und die Zeit steht als Beleg im Register.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Der Slice liefert **keinen neuen Vertrag**: keine neue Zusage, kein neues Werkzeug, keine neue Option, also kein Fall nach `AGENTS.md` §3.10 und keine neue Mutation. Er ändert einen Satz über die **Plattform**, und genau dafür nennt `AGENTS.md` §3.12 (*Die Plattform gehört dazu*) die Randformen; sie stehen hier mit Entscheidung:

| Randform | Frage | Entscheidung |
|---|---|---|
| (a) Welche Versionen nennt der Kopf | nur `u-boot`, oder alle drei Werkzeuge des Laufs | **entschieden:** alle drei des Laufs mit Datum (Liefer-Punkt 1); eine Version ohne Lauf steht dort nicht (`AGENTS.md` §3.11) |
| (b) Zwei gemessene Versionen | frühere Läufe liefen unter 0.10.0 und 0.11.0 | **entschieden:** der Kopf nennt die Umgebung des jüngsten Laufs; eine Liste aller je gelaufenen Versionen trägt er nicht (die Läufe stehen in den Plänen in `done/`) |
| (c) Schwelle gegen gemessene Version | Punkt 3 verlangt `ab 0.10.0`, gemessen ist 0.13.0 | **entschieden:** die Schwelle bleibt (§1, Klasse 2); gemessene Umgebung und Mindestversion sind zwei Aussagen |
| (d) Umfang des Laufs | voller Lauf (19 min 38 s) oder Gruppe Docker | **entschieden:** Gruppe Docker (§1), weil nur sie das echte `u-boot` und `devcontainer` ruft; der volle Lauf gehört dem Nehmer der Aufnahme |
| (e) Plattform des Laufs | Linux-Host; macOS und Windows | **entschieden:** Grenze bleibt wie im Kopf von `SPEC-058` (nicht zugesagt, akzeptiertes Negativ); dieser Slice misst nichts darüber |
| (f) Ausgang bei rotem Fall | Fall der Gruppe Docker unter 0.13.0 rot | **entschieden:** kein Haken, Befund an den Architect, Rückführung nach `open` (§4) |

**Risiken:**

- **Die Gruppe Docker ist unter 0.13.0 rot** (geändertes Verhalten von `u-boot generate`, der Features, des Egress in einer `nft -f`-Transaktion, DNS nur für den lokalen Resolver-Nutzer — die Meldung dazu steht in §6 von `slice-harness-devcontainer-einrichten`). — **Ausgang:** *entfallen*, wenn der Lauf grün ist (Beleg in §7); sonst *eingetreten* → Befund an den Architect und Rückführung nach `open` (§4), dazu ein Beleg in `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`.
- **Der Lauf der Gruppe Docker reißt für sich die 600 s** von [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (der volle Lauf tut es mit 19 min 38 s sicher). — **Ausgang:** *entfallen* mit der Begründung aus §1: `SPEC-060` Punkt 16 (*Aufnahme*) und Punkt 15 (7) binden die Grenze an die Aufnahme, und `container-start-gegenprobe` steht als `nicht gemessen` mit `slice-harness-gegenproben-grenze-rest` als Nehmer; die gemessene Zeit geht als Sendung dorthin und als Beleg ins Register (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`).
- **Die gemessene Zeit hängt an der Last des Hosts** (`BEO-REPO/lastabhaengige-zaehlung`) und ein zweiter Lauf auf demselben Daemon stört (`BEO-REPO/gegenprobe-temp-praefix-geteilt`). — **Ausgang:** *weiter offen* → `BEO-REPO/lastabhaengige-zaehlung`, wenn die Zahl in §7 an der Last hängt; *entfallen*, wenn §7 Last und Alleinlauf belegt.
- **`docs/developer/devcontainer.md` nennt weiter 0.10.0, der Vertrag 0.13.0** (§1, Klasse 2): Ein Leser kann die beiden Zahlen für einen Widerspruch halten. — **Ausgang:** *entfallen*, wenn der Kopf von `SPEC-058` sagt, **was** unter 0.13.0 gemessen ist (Werkzeug und Gegenprobe), und das Dokument weiter „Stand der Proben“ sagt; sonst *eingetreten* → Befund an den Planner, der einen Slice für die Proben des Dokuments anlegt.
- **`u-boot` steht beim Lauf schon auf einer höheren Version als 0.13.0** (das Werkzeug zieht am Host weiter, `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`). — **Ausgang:** *entfallen*, wenn der Kopf die Version des Laufs nennt (Liefer-Punkt 1, Randform (a)) — der Slice trägt die gemessene, nicht die geplante Zahl; die Kennung dieses Slice bleibt dann richtig, auch wenn sie 0.13.0 nennt.
- **Der Nehmer `slice-harness-gegenproben-grenze-rest` liegt bei der Lieferung nicht mehr in `open/` oder `next/`** (schon in `in-progress/` oder `done/`): Dann hat die Sendung keinen Ort (`AGENTS.md` §3.13). — **Ausgang:** *entfallen*, wenn er bei der Lieferung offen liegt und die Sendung in §1, §6 und §8 trägt; sonst *eingetreten* → Befund an den Planner (die Zahl steht dann in der Closure-Notiz dieses Slice, nicht in jenem Plan).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgegangen (Stand 2026-10-21, gemergter Stand). Treffer: `gate-laufzeit-waechst-mit-gegenprobe` (Stand verkörpert; die 19 min 38 s dieser Gegenprobe sind dort belegt — dieser Slice legt einen weiteren Beleg mit der Zeit der Gruppe Docker, DoD), `gegenprobe-haengt-an-umgebung-des-hosts` (die Version am Host ist hier der Gegenstand; ein weiterer Beleg nur, wenn das Verhalten abweicht, §6), `lastabhaengige-zaehlung` (§6), `gegenprobe-temp-praefix-geteilt` (ein Lauf zur Zeit, §4 *Start*), `empfehlung-als-nutzerentscheidung-bezeichnet` (deshalb trägt §4 *Start* die Entscheidung des Nutzers vom 2026-10-21 mit Datum und der Kopf die Welle als Entscheidung des **Planners**). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-21, im Commit, der diesen Plan anlegt):** **Liefer-Punkte 2** — (1) die gemessene Umgebung im Kopf von `SPEC-058`, (2) der Lauf der Gruppe Docker, der sie belegt. Der Lauf ist ein eigener Punkt, weil er mit dem Umfang wächst (Host, Zeit, Versionen) und ohne ihn der Satz eine Behauptung wäre; die Sendung in `slice-harness-gegenproben-grenze-rest` ist **kein** dritter Punkt (sie trägt keinen neuen Gegenstand, nur die Kennung und die Zahl in einen fremden Plan, und die Planung ist keine Schicht). Nicht gezählt: Gate-Läufe, Review-Report, Closure-Notiz, Register-Beleg, Risiko-Ausgänge und die drei Paarungen. **Schichten 1** — der Vertrag eines Harness-Werkzeugs (`harness/vertraege/SPEC-058-container-start.md`); kein Produkt-Code, kein Skript, kein Fragment, keine Nutzer- oder Wartungs-Doku (§1 schließt `docs/developer/` aus, sonst wäre es eine zweite Schicht), die Planung zählt nicht. **Schätzung 4 bis 20** in der Metrik von §4 bei **Marke 60** und Meldeschwelle 40. Damit unter drei Liefer-Punkten und unter zwei Schichten.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (`REPO`: Vertrag und Entscheidung gehen dem Werkzeug voraus; hier folgt der Vertrag einer Messung am Host, und die Messung ist der Beleg, nicht eine Inventur des Bestands).
