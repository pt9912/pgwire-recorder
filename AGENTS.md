# AGENTS.md — Briefing für AI-Coding-Agenten

---

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementer-Agent immer
einhalten muss.

Regeln dieser Datei: Baseline-Regelwerk `modul-09-implementierung.md`
§Ziel-Form: AGENTS.md — sie trägt Hard Rules und Pointer auf kanonische
Quellen, sie dupliziert deren Inhalt nicht; sonst entsteht Drift.

**Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle** (Source Precedence — siehe
`harness/README.md`).

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md).

Das **Regelwerk der adoptierten Baseline** ist die **präsente,
nachschlagbare Vertiefung** zu diesem Briefing: ein self-navigierbares
**Modul-Bundle** (`README.md` = Index). Beim Bootstrap wird das
self-contained Release-ZIP
(<https://github.com/pt9912/ai-harness-course/releases/download/v6.18.0/lab-regelwerk.zip>)
**committet vendored** unter `.harness/baseline/<tag>/{regelwerk,templates}/`
(Regelwerk *und* Templates parallel, netzlos materialisiert samt `SHA256SUMS`
— Vorgehen siehe
Baseline-Regelwerk `modul-02-harness-bootstrap.md`;
Quelle/Stand in [`harness/conventions.md`](harness/conventions.md) §Baseline).

Die verkörperte Form (dieses Briefing, die Konventionen, deine
ausgefüllten Artefakte) **führt**; das Regelwerk wird **pro Entscheidung
nachgeschlagen, deren
operative Detailtiefe das Briefing nicht trägt** — Trigger-Klassen,
Sub-Area-Qualifikation, Carveout-vs-Reconciliation, Modus-Diagnose. Dabei
**nur den benötigten Abschnitt** laden (README ist der Index), **nicht das
ganze Regelwerk im Kontext halten**. Breiterer Pflicht-Blick bleibt bei:
Bootstrap, Änderung an [`harness/conventions.md`](harness/conventions.md)
(Adaptionen `MR-<NNN>`, Source-Precedence, ID-Schema), Drift-Audit gegen die
Baseline (Baseline-Regelwerk `modul-02-harness-bootstrap.md`
§Freshness-Audit der vendored Baseline — darunter die Stichprobe gegen
den Bestand, die auch bei aktuellem Pin läuft). Derivativ: bei Konflikt gelten die kanonischen Quellen.

Die **Skelett-Vorlagen** der Baseline liegen **vendored** unter
`.harness/baseline/<tag>/templates/` (aus demselben Baseline-Bundle) und
tragen zwei Rollen: als **Referenz-Form**, auf die das Regelwerk mit
`../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu
`regelwerk/` vendored), und als Vorlage, die beim Anlegen neuer Artefakte
(ADR, Slice, Welle, Carveout, Review-Report) **kopiert und ausgefüllt** wird
statt frei zu formulieren.

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — vertraglich abnahmebindend.
2. [`spec/spezifikation.md`](spec/spezifikation.md) — technisch verbindlich, fortschreibbar.
3. [`spec/architecture.md`](spec/architecture.md) — Komponenten- und Sequenzsicht.
4. [`docs/plan/adr/`](docs/plan/adr/) — ADR-Verzeichnis und -Index.
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — Wellen-Sequenz.
6. `docs/user/*` *(falls vorhanden)* — Operations, Quality, Releasing. <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
7. [`README.md`](README.md) — Projekt-Überblick.
8. **AGENTS.md (diese Datei).**
9. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln



### 3.1 Docker-only



Kein lokales <venv/SDK/Toolchain-Install>. Alles läuft über `make`
(das Docker nutzt). Host braucht nur Docker und GNU `make`.

**Falsch:** <z.B. `pip install ...`>
**Richtig:** <z.B. `<make-target>`>

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

### 3.2 Suppression-Verbot

Eine Direktive `nolint` im Go-Code bricht das Gate `make lint`. Ausnahmen leben nur
in `.golangci.yml`, jede Regel unter `exclusions.rules` mit einem Kommentarblock
`# Why:` unmittelbar darüber; eine Regel, die nichts ausblendet, bricht das Gate
ebenso. Vertrag und Lesart: [`harness/sensors/lint.md`](harness/sensors/lint.md).

**Falsch:** `os.Remove(pfad) //nolint:errcheck` im Code.
**Richtig:** den Befund beheben; ist eine Ausnahme dauerhaft richtig, eine Regel in
`.golangci.yml` mit `# Why:` und einem Grund, der auch für neuen Code gilt.

**Begründung:** Eine Ausnahme im Code sieht nur, wer die Zeile liest; im Profil
steht sie an einer Stelle mit ihrem Grund, und die Gegenprobe des Gates prüft sie
([ADR-0034](docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md)).

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird, sind das
zwei Commits — der Move-Commit bleibt rein (Git erkennt R-Rename). Welcher
zuerst kommt, sagt der Vorgang:

1. Regelfall: `git mv source target` → eigener Commit, dann Inhalt umschreiben.
2. Lifecycle-Übergang nach `done/`: erst der Inhalt (DoD-Häkchen,
   Closure-Notiz), dann der reine `git mv` — die Notiz ist die Bedingung für
   `done/`, nicht ihre Folge.

**Begründung:** Sonst fällt die Rename-Detection unter die 50%-
Similarity-Schwelle und `git log --follow` wird unzuverlässig.

### 3.4 Architektur ist sprach- und meilensteinfrei

`spec/architecture.md` darf Pfade zu **Code-Modulen** referenzieren
(`src/service/`), aber **keine** Wellen, Slices, Commit-Hashes oder
Closure-Daten. Die zeitliche Schicht lebt in
`docs/plan/planning/` und den späteren Closure-Notizen. Auch **keine
ADR-Bezüge**: Die Sicht steht im Stabilitäts-Rang über der ADR; welche ADR
eine Aussage verbindlich macht, deklariert die ADR in ihrem `Schärft:`-Feld.

Diese Regel ist *verkörpert*, nicht hier entschieden — sie folgt aus dem
Sicht-Stratum (Baseline-Regelwerk `modul-03-spec.md`
§Ziel-Form: Architektur-Sicht).

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NNNN`.

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Architekturregel)
ist ein ADR, kein PR-Kommentar. Eine befristete Ausnahme für einen Teil (einen
Layer, einen Pfad) ist keine Senkung, sondern ein Carveout mit Trigger und
Folge-Slice; die Schwelle selbst bleibt.

### 3.7 Ein Kommentar beschreibt, was da ist

Gilt für Code, Konfiguration und Skripte — und für Zustandsfelder (unten).
Ein Kommentar trägt eine dieser Klassen — **Zusage · Kopplung · Abgrenzung ·
Rang-Zeiger · Grenze** — und schreibt an den, der die Stelle *ändert*, nicht an den, der die
Entscheidung *trifft*. Regeln dieser Sektion: Baseline-Regelwerk
`grundlagen-harness-dateien.md` §Was ein Kommentar trägt.

**Falsch:** <z.B. „Ohne dieses Feld behauptete die Ausgabe eine Verteilung,
die nicht stattgefunden hat"> — Konjunktiv über die verworfene Alternative.
**Richtig:** <z.B. „Verteilt ist wahr, wenn die Splitting-Regel angewendet
werden konnte"> — Indikativ über den Zustand.

**Falsch:** <z.B. „die frühere Fassung prüfte nur die Länge"> — beschreibt
abwesenden Text.
**Richtig:** die geltende Zusage nennen; die vorige hält `git`.

**Zustandsfelder ebenso:** Eine `Stand`-/`Status`-Zelle in Roadmap,
Beobachtungs-Register oder Meilenstein-Tabelle nennt den Zustand und den Beleg
als auflösbaren Anker, nicht die Chronik; das Drift-Log der Roadmap trägt nur
Umplanungen, keine Schließungen und keine erreichten Meilensteine.

**Begründung:** Die Abwägung gehört in die ADR, die Historie in `git`, die
Herkunft in **ein** auflösbares Feld (`LH-*`, `ADR-*`, `· seit welle-<Kennung>`).
Was daneben steht, liest jeder Lauf mit und bezahlt es mit Kontext.





### 3.8 Eine ADR trägt Entscheidung und Gründe, keine Regeldetails

Eine ADR hält fest, was entschieden wurde und warum. Die Regeln im Einzelnen
(Fälle, Ausnahmen, Codes, Rangfolgen) stehen in der Spezifikation; die ADR zeigt
über `Schärft:` auf sie und wiederholt sie nicht.

**Falsch:** die Fallunterscheidungen einer Spezifikationsstelle in der ADR
nacherzählen.
**Richtig:** die Entscheidung, die Alternativen und die Gründe nennen und für die
Einzelregeln auf die Spezifikationsstelle zeigen.

**Begründung:** Angenommene ADRs sind unveränderlich (3.5). Eine wiederholte Regel
veraltet mit der Spezifikation und zwingt zu einer Ersetzungs-ADR, die nur den
Text berichtigt.

### 3.9 Der Slice-Plan folgt jeder Korrektur (seit welle-walking-skeleton)

Wer im Slice Code, Gates, Tests oder die Spezifikation ändert, zieht im selben
Commit §1 (Ziel und Abgrenzung), §3 (Plan) und §6 (Risiken) des Slice-Plans nach,
dazu den Kopf (`Bezug`, `Berührte Spec-Stellen`) und betroffene Folge-Slices.

**Sensor für die Kopf-Hälfte** (seit slice-harness-kopf-sensor): `make kopf-check`
meldet rot, wenn §1 oder §2 eines Slice-Plans in `open/`, `next/` oder
`in-progress/` eine Lastenheft-, Spezifikations- oder Sicht-Kennung nennt, die der
Kopf nicht führt ([ADR-0032](docs/plan/adr/0032-kopf-sensor-fuer-slice-plaene.md)).
Ob §1, §3 und §6 inhaltlich dem Diff folgen, bleibt Urteil von Review und Verifikation.

**Falsch:** den Mismatch liefern, während §1 ihn noch „welle-replay-semantik“
zuweist.
**Richtig:** §1 nennt, was der Slice liefert und was er abgibt; der Folge-Slice
nennt, was schon geliefert ist.

**Begründung:** Review und Verifikation messen gegen den Plan; ein Plan, der dem
Code nicht folgt, erzeugt in jeder Runde dieselben Findings
(`BEO-REPO/plan-folgt-korrektur-nicht`). Entfernen oder Lockern dieser Regel
setzt den Retirement-Check voraus: Ist die Beobachtung seit welle-walking-skeleton
wieder aufgetreten?

### 3.10 Jede Zusage eines neuen Vertrags fängt eine Mutation (seit welle-extended-query)

Wer einen neuen Vertrag liefert (Protokollrand, Format, Leser, Diagnose, Option,
Gate), schreibt zu jeder seiner Zusagen einen Test, der rot wird, wenn der Code die
Zusage bricht, und fährt diese Mutation vor der Übergabe selbst. Der Bericht nennt je
Zusage: Zusage · Mutation · roter Test.

**Falsch:** Die Tests decken den Gutfall; dass eine Mutation grün bleibt, findet erst
das Review.
**Richtig:** Je Zusage eine Mutation, vom Implementer selbst rot gesehen.

**Begründung:** Ein grünes Gate belegt, dass nichts bricht, nicht, dass eine Zusage
geprüft ist (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`). Entfernen oder
Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit
welle-extended-query wieder aufgetreten?

### 3.11 Ein Kommentar sagt nur zu, was ein Test prüft (seit welle-extended-query)

Kommentar, Hilfetext, Abdeckungs-Deklaration, Plan-Zeile und Benutzerhandbuch sagen
nur zu, was ein Test oder Gate prüft. Reicht der Satz weiter, wird er enger gefasst oder
bekommt seinen Test (§3.10). Ein Beispiel im Benutzerhandbuch gilt erst als belegt, wenn
es als Datei gegen das gebaute Binary läuft; gelesen gegen den Code ist es nicht geprüft
(seit slice-v1-abschluss-upstream-verbinden). Prosa über ein Gate (Command, Plan, Kommentar)
beschreibt sein Verhalten nicht nach: Sie nennt die Regel an den Schreiber und zeigt auf die
Konfiguration als Quelle (`.d-check.yml`, `.a-check.yml`, `.golangci.yml`); ein Verhalten des
Gates gilt erst als zugesagt, wenn eine Gegenprobe es zeigt (seit slice-harness-d-check-v0-85).

**Falsch:** „prüft beide Adapter“ über einer Gegenprobe, die einen prüft.
**Richtig:** „prüft den Recording-Adapter“, oder die Gegenprobe prüft beide.

**Falsch:** Das Handbuch zeigt eine Konfigurationsdatei mit einem Abschnitt, den das
gebaute Binary mit `PGR-E2004` ablehnt; Implementer und Review hielten ihn gegen den Code
gelesen für richtig.
**Richtig:** Das Beispiel liegt als Datei vor und läuft gegen das gebaute Binary, etwa mit
`config show`, ohne Meldung.

**Handbuch und README beschreiben den Ist-Zustand** (seit
slice-v1-abschluss-einspielen-laufsteuerung, Entscheidung des Nutzers vom 2026-10-10):
`docs/user/benutzerhandbuch.md` und `README.md` beschreiben nur das Verhalten, das das gebaute
Binary heute zeigt. Keine Chronik („noch nicht“, „kommt mit“, „bis zu“), kein Zielstand; das
Zielbild steht in `spec/`. Das Handbuch verweist nicht auf Spezifikation, Lastenheft, ADRs,
Slices, Wellen oder Reviews. `README.md` darf auf `spec/` zeigen (Verweis nach oben,
`harness/README.md` §Source precedence), ohne Chronik und ohne zu sagen, was davon noch kommt.
Ein Slice, der Verhalten liefert, liefert seinen Teil von Handbuch und README im selben Slice
oder, wenn die Schichtzählung (§3.13) das nicht zulässt, im Doku-Folge-Slice `<Kennung>-doku`
direkt dahinter, in derselben Welle. Die Welle schließt erst mit beiden, und ein Release
entsteht erst danach; zwischen Code-Slice und Doku-Slice hinkt das Handbuch dem Binary hinterher
(seit slice-doku-ist-stand, Entscheidung des Nutzers vom 2026-10-10).

**Falsch:** Das Handbuch beschreibt `--compare-responses` bei `play`, das das gebaute Binary mit
`PGR-E2001` ablehnt, und eine Closure nimmt den Zielstand bis zum Ende der Welle hin.
**Richtig:** Das Handbuch nennt nur, was das Binary kann; der Slice, der `--compare-responses`
liefert, schreibt dessen Abschnitt.

**Falsch:** Ein Command zählt auf, wann `make docs-check` rot ist; jede Runde der Prüfung
findet eine weitere Form, die grün bleibt.
**Richtig:** Der Command nennt die Regel, Kennungen als Links zu schreiben, und sagt, dass
`.d-check.yml` festlegt, was das Gate prüft.

**Begründung:** Eine Zusage ohne Prüfung liest jeder Lauf als geprüft
(`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Entfernen oder Lockern setzt
den Retirement-Check voraus: Ist die Beobachtung seit welle-extended-query wieder
aufgetreten, für das Handbuch seit slice-v1-abschluss-upstream-verbinden, für Prosa über ein
Gate seit slice-harness-d-check-v0-85, für den Ist-Zustand von Handbuch und README seit
slice-v1-abschluss-einspielen-laufsteuerung?

### 3.12 Randformen eines neuen Vertrags sind vor dem Code entschieden (seit slice-harness-randformen-vor-code)

Ein Slice, der einen neuen Vertrag liefert (§3.10), nennt in §6 des Slice-Plans dessen
Randformen, etwa fehlender, leerer oder unbekannter Wert, Alias, Abbruch und
Herunterfahren, was ein Diagnose-Feld je Nachrichtenart heißt, Verhalten je
Serverversion, und je Randform, wo sie entschieden ist. Entschieden ist jede vor dem
ersten Code-Commit, in der Spezifikation (Technik-Stratum), auch für
Harness-Werkzeuge (seit slice-lastenheft-pruefbarkeit); auch eine Entscheidung des
Nutzers wird dort festgehalten. Eine ADR trägt nur Entscheidung und Gründe und zeigt mit
`Schärft:` auf diese Stelle (§3.8). Der Architect prüft die Liste vor dem Code. Eine nicht
genannte oder offene Randform entscheidet der Implementer nicht, er hält an und gibt
sie dem Architect zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe). Eine
Randform, die erst ein Commit mit Code, Tests oder Gates in §6 einträgt, ist nicht genannt
(seit slice-harness-blackbox-kern).

**Falsch:** Der Leser behandelt `null` wie einen fehlenden Schlüssel, die
Spezifikation sagt nichts dazu, und das Review findet es.
**Richtig:** §6 nennt `null` und fehlenden Schlüssel und zeigt auf die Stelle der
Spezifikation, die beide entscheidet; dann folgt der Code.

**Begründung:** Was der Code still entscheidet, findet erst das Review, Runde um Runde
eine weitere Randform (`BEO-REPO/spec-randform-erst-im-review-entschieden`). Ein Diff,
der eine Randform entscheidet, die §6 nicht nennt, ist ein Review-Befund gegen diese
Regel. Entfernen oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung
seit slice-harness-randformen-vor-code wieder aufgetreten?

### 3.13 Eine Adresse nimmt an, bevor sie genannt wird (seit slice-lint-bestand-kern-driven)

Wer einem anderen Slice etwas zuweist (Folge-Slice, Risiko-Ausgang *eingetreten*,
Register-Ausgang *geplant*, Abgrenzung der Klasse 1), liest vorher dessen §1
*Ausdrücklich NICHT* und dessen DoD. Im selben Commit trägt er die Sendung in §1 oder
die DoD des Nehmers ein, mit der Kennung des Gebers. Trifft ein Ausschluss des Nehmers
die Sendung oder liegt der Nehmer in `done/`, ist er keine Adresse; der Punkt geht an
den Planner.

**Falsch:** §6 weist grüne Mutanten `slice-harness-coverage` zu, dessen §1 neue Tests
über dem gemessenen Stand ausschließt.
**Richtig:** Der Nehmer nennt die Mutanten in §1 und DoD mit der Kennung des Gebers;
`grep` nach ihr findet sie im Nehmer.

**Begründung:** Beim Zuweisen geprüft war nur, dass es den Nehmer gibt; dass er
annimmt, fand erst das Review (`BEO-REPO/folge-slice-adresse-nimmt-nicht-an`). Entfernen
oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit
slice-lint-bestand-kern-driven wieder aufgetreten?

**Auch eine Bindung ist eine Sendung** (seit slice-v1-abschluss-einspielen-laufsteuerung): Ruht
ein akzeptiertes Negativ, eine Grenze oder eine Randform darauf, dass ein anderer Slice etwas
einhält, ist das eine Zuweisung wie die vier oben. Derselbe Commit trägt die Bindung in §1 des
Nehmers ein, mit der Kennung des Gebers, und nennt sie dort in §6.

**Falsch:** Ein akzeptiertes Negativ schließt einen Fall aus, weil zwei spätere Slices an den
Vertrag eines Ports gebunden seien; keiner der beiden nennt die Bindung.
**Richtig:** Beide nennen sie in §1 mit der Kennung des Gebers und in §6 als Randform.

**Begründung:** Die Aufzählung der vier Formen las sich als abschließend, und die Bindung ging
ohne Eintrag hinaus; das Review fand es (`BEO-REPO/folge-slice-adresse-nimmt-nicht-an`).
Entfernen oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit
slice-v1-abschluss-einspielen-laufsteuerung wieder aufgetreten?

**Nachzählen beim Eintragen** (seit slice-v1-abschluss-einspielen): Wer eine Sendung in
einen Nehmer einträgt, zählt im selben Commit dessen Liefer-Punkte und Schichten nach,
mit der Sendung (Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice:
höchstens drei Liefer-Punkte, höchstens zwei Schichten). Läge der Nehmer darüber, geht der
Punkt an den Planner statt in den Nehmer.

**Schichten zählen alle Pläne gleich** (seit slice-doku-ist-stand, Entscheidung des Nutzers
vom 2026-10-10): nach der Teilung in `slice-harness-d-check-v0-85` §8. Nutzer- und
Wartungs-Doku (`docs/user/`, `docs/maintainer/`, `README.md`) ist eine Schicht, die Planung
(Pläne, Roadmap, Register) keine; im Produkt-Code zählt jede Schicht des Hexagons, ebenso ihre
Tests. Liegt ein Slice mit seinem Handbuch- und README-Teil über zwei Schichten, geht dieser
Teil in einen eigenen Slice `<Kennung>-doku` direkt hinter ihm (§3.11).
Der Bootstrap (`internal/bootstrap`, `composition_root`) zählt nur dann als Schicht, wenn er
Logik enthält; reine Verdrahtung (ein Feld durchreichen, ein Kommentar) ist keine Schicht
(seit slice-v1-abschluss-einspielen-anmeldung, Entscheidung des Nutzers vom 2026-10-10).
Kommt Logik hinein (Datei lesen, Fehler einstufen, verzweigen), zählt er, und der Plan prüft
die Grenze neu.

**Falsch:** Ein Plan zählt den Bootstrap für das Reichen eines Feldes als dritte Schicht, ein
anderer für dieselbe Änderung nicht.
**Richtig:** Verdrahtung zählt nicht; liest der Bootstrap eine Datei oder stuft er einen
Fehler ein, zählt er.

**Falsch:** Ein Plan schreibt „Handbuch und README zählen als Dokumentation, nicht als Schicht“
und bleibt mit zwei Code-Schichten bei zwei.
**Richtig:** Der Plan zählt die Dokumentation als dritte Schicht und gibt ihren Teil an den
Doku-Folge-Slice ab.

**Begründung:** Geber und Nehmer zählten im selben Commit mit verschiedenen Teilungen
(`BEO-REPO/schichtteilung-je-plan-verschieden`). Entfernen oder Lockern setzt den
Retirement-Check voraus: Ist die Beobachtung seit slice-doku-ist-stand wieder aufgetreten?
Für die Zählung des Bootstrap gilt derselbe Check seit slice-v1-abschluss-einspielen-anmeldung.

**Falsch:** Die siebte Übernahme wandert in einen vorhandenen Liefer-Punkt; der Nehmer
bleibt bei drei Punkten und berührt vier Schichten, die Größe findet erst die Prüfung des
Architect vor dem Code.
**Richtig:** Der Commit, der die Sendung einträgt, zählt den Nehmer mit ihr nach und nennt
die Zählung in dessen §8; läge er darüber, entscheidet der Planner, ob geschnitten wird oder
ein anderer Slice annimmt.

**Begründung:** Die Größenregel zählt Punkte, und ein Punkt nimmt Übernahme um Übernahme
auf, ohne dass die Zahl steigt (`BEO-REPO/slice-waechst-durch-uebernahmen`). Entfernen
oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit
slice-v1-abschluss-einspielen wieder aufgetreten?

## 4. Quality Gates

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt. Der Gate-Index steht **einmal**, in
[`harness/README.md`](harness/README.md) §Sensors — dort steht auch die
*Bindung* jedes Targets. Diese Datei führt die Liste nicht.

Kein Target nennen, das im Makefile nicht existiert — auch nicht in Prosa.

## 5. Dokumentations-Regeln



| # | Regel | Datei |
|---|---|---|
| 1 | **Anforderungs-IDs und ADR-Nummern** müssen in PRs/Commits referenziert sein — sie sagen, welche Zusage oder Entscheidung berührt ist. Struktur-IDs (`SPEC-<NNN>`, `ARC-<NNN>`) adressieren *innerhalb* der Spec und gehören nicht in die Commit-Message. | — |
| 2 | Vergeben werden IDs beim Spec-/ADR-Schreiben nach dem in `harness/conventions.md` deklarierten ID-Schema (Default: `<PREFIX>-FA-<NN>` / `<PREFIX>-QA-<NN>` / `<PREFIX>-RB-<NN>` aus dem Lastenheft, `SPEC-<NNN>` in der Spezifikation, `ARC-<NNN>` in der Sicht, ADR-Nummern über den ADR-Index) — nie ad hoc im PR. | — |
| 3 | Neue ADRs müssen den ADR-Index aktualisieren. | — |
| 4 | Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in `spec/architecture.md`. | — |

## 6. Minimal Agent Workflow

Pro Slice:

1. `harness/README.md` lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe `harness/README.md` §Guides) →
Verifier. Kein Self-Review — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).
