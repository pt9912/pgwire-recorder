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
2. [`spec/spezifikation.md`](spec/spezifikation.md) für das Produkt und [`harness/vertraege/`](harness/vertraege/README.md) für die Werkzeuge — technisch verbindlich, fortschreibbar.
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

Wer im Slice Code, Gates, Tests, die Spezifikation oder Harness-Artefakte
(`.claude/commands/`, `.harness/skills/`, Hooks) ändert, zieht im selben
Commit §1 (Ziel und Abgrenzung), §3 (Plan) und §6 (Risiken) des Slice-Plans nach,
dazu den Kopf (`Bezug`, `Berührte Spec-Stellen`) und betroffene Folge-Slices. Die
Harness-Artefakte zählen seit slice-v1-abschluss-zeitangaben-doku (Entscheidung des
Nutzers vom 2026-10-12): Ein Satz in `implement-slice.md`, den der Plan weder in §1 noch
in §3 nennt, ist derselbe Mismatch (`BEO-REPO/plan-folgt-korrektur-nicht`). Entfernen
des Zusatzes setzt den Retirement-Check voraus: Ist die Beobachtung seit
slice-v1-abschluss-zeitangaben-doku wieder aufgetreten?

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

**Aufzählung und Ortsangabe zählen je Element** (seit slice-harness-gate-kette-stop-hook-liste,
Entscheidung des Nutzers vom 2026-10-13): Je Element einer Aufzählung in einer Zusage und je
Schritt einer Ortsangabe (zum Beispiel „steht vor allen Prüfungen“) ein Fall und eine Mutation,
die nur ihn bricht.

**Falsch:** Die Zusage nennt „Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8“ im Pfad und
„davor“ für den Schleifenschutz; die Gegenprobe hat je Zusage einen Fall, Mutationen einzelner
Pfadformen (F-703) und einzelner Schritte vor dem Schleifenschutz (F-704) bleiben grün, und das
Review findet sie.
**Richtig:** Ein Fall je Pfadform und einer je Schritt vor dem Schleifenschutz, jeder mit einer
Mutation, die nur ihn bricht.

**Begründung:** Eine Zusage mit Aufzählung oder Ortsangabe trägt mehrere Zusagen, und eine
Mutation je Zusage zählte sie als eine (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`).
Entfernen oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit
slice-harness-gate-kette-stop-hook-liste wieder aufgetreten?

**Ein grüner Mutant gilt als äquivalent erst mit Invariante, Ort und Sonde** (seit welle-v1-abschluss,
Entscheidung des Nutzers vom 2026-10-18): Wer in §7 einen grünen Mutanten äquivalent nennt, nennt (a) die
Invariante, aus der gleiches Verhalten folgt, (b) den Ort, der sie erzwingt (Typ, `Validate`, Leser, eine
vorgelagerte Prüfung im selben Pfad; ein Kommentar oder die Dokumentation des Modells ist kein Ort) und (c)
eine Sonde, die Original und Mutant trennen würde, wenn die Invariante nicht gälte, mit ihrem Lauf. Ist die
Invariante die Zusage, die der Slice liefert, oder gilt sie nur laut Dokumentation, ist der Mutant nicht
äquivalent und bekommt einen Fall. Eine Prüfung wird nicht als überflüssig entfernt, solange (a) bis (c) in
§7 fehlen.

**Falsch:** §7 führt „`Send` liest nach dem Ende“ als äquivalent, Invariante „`Send` enthält keinen
Lesezugriff“; das ist die Zusage selbst, und der Mutant bleibt grün (F-1 in
`slice-v1-abschluss-verbindungsende-upstream`).
**Richtig:** §7 nennt die Invariante, den Ort, der sie erzwingt, und die Sonde mit Ergebnis „gleich“; fehlt
eines davon, bekommt der Mutant einen Test.

**Begründung:** Die Begründung eines äquivalenten Mutanten stützt sich auf die Eingaben, an die der
Implementer gedacht hat; eine andere unterscheidet ihn, und erst Review oder Verifikation finden sie
(`BEO-REPO/gruener-mutant-faelschlich-aequivalent`, fünf Slices). Entfernen oder Lockern setzt den
Retirement-Check voraus: Ist die Beobachtung seit welle-v1-abschluss wieder aufgetreten?

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
Harness-Werkzeuge (seit slice-lastenheft-pruefbarkeit): in `harness/vertraege/` (Ort seit
slice-harness-vertraege-umzug-gate-kette); auch eine Entscheidung des
Nutzers wird dort festgehalten. Eine ADR trägt nur Entscheidung und Gründe und zeigt mit
`Schärft:` auf diese Stelle (§3.8). Der Architect prüft die Liste vor dem Code. Eine nicht
genannte oder offene Randform entscheidet der Implementer nicht, er hält an und gibt
sie dem Architect zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe). Eine
Randform, die erst ein Commit mit Code, Tests oder Gates in §6 einträgt, ist nicht genannt
(seit slice-harness-blackbox-kern).

**Aufruf-Voraussetzung gehört dazu** (seit slice-v1-abschluss-zeitangaben-aufzeichnen): Sagt ein Vertrag eine Ordnung
oder Eindeutigkeit der Werte zu, die er über mehrere Aufrufe einer Session liefert (nicht absteigend, fortlaufend), nennt §6
als Randform, ob diese Aufrufe nacheinander oder gleichzeitig kommen, und die Spezifikation bindet die Zusage an diese
Voraussetzung. Bemerkt der Implementer ein Gegenbeispiel, hält er an und gibt es zurück; ein Satz in seinem Bericht ist keine Rückgabe.

**Falsch:** Die Zusage „Werte steigen nicht ab“ steht ohne Bedingung, das Gegenbeispiel (zwei gleichzeitige Aufrufe) steht als
Satz im Bericht, und das Review findet es.
**Richtig:** §6 nennt „Aufrufe einer Session nacheinander“, die Spezifikation bindet die Zusage daran, der Port nennt es im Kommentar.

**Falsch:** Der Leser behandelt `null` wie einen fehlenden Schlüssel, die
Spezifikation sagt nichts dazu, und das Review findet es.
**Richtig:** §6 nennt `null` und fehlenden Schlüssel und zeigt auf die Stelle der
Spezifikation, die beide entscheidet; dann folgt der Code.

**Aufrufarten des Treibers gehören dazu** (seit slice-harness-gate-kette-schnell, Vorschlag des Planners; der Nutzer kann ihn
streichen): Liefert ein Werkzeug ein Ergebnis, das ein Hook oder ein Gate später glaubt (ein Stand, ein Nachweis, ein Exit-Code),
nennt §6 als Randform die Aufrufarten außerhalb des Normalfalls, mit denen der Treiber (`make`, `bash`, `git`) das Rezept leerlaufen
lässt oder Fehler überliest, bei `make` `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` aus der Umgebung, auch in Kombination, und je eine
Entscheidung; die Spezifikation trägt sie.

**Falsch:** `make -n gates-schnell` führt das Rezept aus, weil es `$(MAKE)` enthält, fährt kein Gate und schreibt den Stand, den der
Hook als „grün“ liest; §6 nannte die Flags nicht, das Review fand es.
**Richtig:** §6 nennt `-n`, `-q`, `-t` (kein Stand, ein vorhandener bleibt) und `-i` (kein Stand, ein älterer entfällt) samt dem
Vorrang bei Kombination, und je Flag steht ein Fall in der Gegenprobe.

**Die Stelle einer Nachricht gehört dazu** (seit slice-harness-randform-stelle-und-lesart, Entscheidung des Nutzers vom
2026-10-15): Enthält ein Vertrag eine Nachricht, die eine Folge beendet oder ersetzt (Ende, `ready_for_query`, `sync`), nennt §6
je Nachrichtenart die zulässigen Stellen in der Folge (Anfang, Mitte, direkt vor dem Ende, letzte Gruppe, nicht letzte Gruppe)
und je Stelle eine Entscheidung; die Spezifikation trägt sie. Kann eine Folge aus mehreren Gründen abgelehnt werden, nennt §6 die
Reihenfolge der Gründe, und die Spezifikation bindet sie.

**Falsch:** §6 nennt die Nachrichten `ready_for_query` und `connection_closed`, aber nicht, ob sie vor dem Ende oder mitten in der
Folge stehen dürfen; der Implementer entscheidet beide im Code, und das Review findet sie (F-1 in
`slice-v1-abschluss-verbindungsende-format`).
**Richtig:** §6 führt je Nachricht die Stellen auf und zeigt auf die Stelle der Spezifikation, die jede entscheidet; dann folgt der
Code.

**Die Plattform gehört dazu** (seit welle-v1-abschluss, Entscheidung des Nutzers vom 2026-10-18): Sagt ein
Vertrag eine Eigenschaft zu, die vom Betriebssystem oder Dateisystem abhängt (atomares Verschieben, Sperren,
Transaktionen auf einer Datei, Signale, Pfadform, Zeilenende), nennt §6 als Randform die Plattformen von
`LH-QA-03`, auf denen sie geprüft ist, und je übrige Plattform eine Entscheidung: Grenze in der Spezifikation
an der Zusage („geprüft nur unter Linux“) oder eine Zusage mit Test. Die Spezifikation trägt die Grenze; ein
Handbuchsatz über die Eigenschaft nennt nur die geprüfte Plattform (§3.11).

**Falsch:** Die Zusage „atomar“ steht ohne Plattform, die Tests laufen unter Linux, und die Grenze für macOS
und Windows steht als Risiko in §6 mit dem Ausgang „weiter offen“ (drei Slices in Folge).
**Richtig:** §6 nennt Linux (geprüft), macOS und Windows (Grenze) und zeigt auf die Stelle der Spezifikation,
die sie als Grenze festhält; dann folgt der Code.

**Begründung:** Die Plattformgrenze trat als Risiko in drei Slices auf, jedes Mal erst im Review oder als
§6-Risiko mit dem Ausgang „weiter offen“ (`BEO-REPO/verhalten-nur-unter-linux-geprueft`). Entfernen oder
Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit welle-v1-abschluss wieder aufgetreten?

**Begründung:** Was der Code still entscheidet, findet erst das Review, Runde um Runde
eine weitere Randform (`BEO-REPO/spec-randform-erst-im-review-entschieden`). Ein Diff,
der eine Randform entscheidet, die §6 nicht nennt, ist ein Review-Befund gegen diese
Regel. Entfernen oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung
seit slice-harness-randformen-vor-code wieder aufgetreten, für die Aufruf-Voraussetzung seit
slice-v1-abschluss-zeitangaben-aufzeichnen, für die Aufrufarten des Treibers seit
slice-harness-gate-kette-schnell, für die Stelle einer Nachricht seit
slice-harness-randform-stelle-und-lesart?

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

**Auch eine Ortsangabe wird nachgezählt** (seit slice-harness-vertraege-umzug-gate-kette,
Entscheidung des Nutzers vom 2026-10-13): Die Zählung steht je Nehmer in dessen §8, im Commit,
der die Sendung einträgt; ein Hinweis an den Planner ersetzt sie nicht. Liegt der Nehmer schon
darüber, nennt die Zählung das und den vorgemerkten Schnitt.

**Falsch:** Der Umzug einer Vertragsstelle trägt die neue Ortsangabe in vier Nehmer ein; der
Bericht meldet die Zählung als Hinweis, §8 der Nehmer bleibt leer, und das Review findet es.
**Richtig:** Jeder der vier Nehmer führt im selben Commit eine Zählung mit der Kennung des Gebers.

**Begründung:** Die Regel galt dem Wortlaut nach für eine Sendung mit Inhalt, und eine Ortsangabe
las sich als keine (`BEO-REPO/slice-waechst-durch-uebernahmen`, zwei Slices in Folge). Entfernen
oder Lockern setzt den Retirement-Check voraus: Ist die Beobachtung seit
slice-harness-vertraege-umzug-gate-kette wieder aufgetreten?

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

**Meldungscode-Konstanten zählen nicht** (seit slice-v1-abschluss-tls-client, Entscheidung des
Nutzers vom 2026-10-13): Eine neue Konstante des Meldungskatalogs in `internal/hexagon/model`
(`PGR-…`) ist keine Schicht; Logik im Model zählt.

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

### 3.14 Eine Rückführungs-Marke nennt ihre Metrik, und die Schätzung rechnet in ihr (seit slice-harness-gate-kette-stop-hook-liste)

Entscheidung des Nutzers vom 2026-10-13. Jede Rückführungs-Marke in §4 eines Plans („zu groß, über
N Zeilen mit Tests“) nennt ihre Metrik: die Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice>` über die Pfade, die der Plan nennt, Tests
eingeschlossen. Der Architect gibt Schätzungen in dieser Metrik an und rechnet für Tabellen- und
Gegenprobe-Code mit rund dem Doppelten des Produktcodes.

**Falsch:** §4 nennt „über 700 Zeilen mit Tests“; gemessen sind 630 (netto), 704 (Einfügungen) und
778 (Summe), ob die Marke eintrat, hängt an der Lesart, und die Schätzung von 440 bis 620 Zeilen
endet bei 1065, weil die Gegenprobe zwei Drittel des Diffs trägt.
**Richtig:** §4 nennt „Summe aus Einfügungen und Löschungen von `git diff --numstat 9662e13..<letzter
Commit> -- tools harness .claude`, über 700“; der Architect schätzt in dieser Zahl und setzt die
Gegenprobe mit dem Doppelten des Produktcodes an.

**Begründung:** Eine Marke ohne Metrik lässt die Probe an der Lesart hängen
(`BEO-REPO/groessenprobe-ohne-messbaren-stand`), und eine Schätzung ohne den Faktor der Tests liegt
unter der Messung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`); der ungeteilte Plan über der Grenze zeigte
sich erst in der Prüfung des Architect (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`).
Entfernen oder Lockern setzt den Retirement-Check voraus: Ist eine der drei Beobachtungen seit
slice-harness-gate-kette-stop-hook-liste wieder aufgetreten?

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
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`). Wer ein Gate, eine Gegenprobe oder einen Vertrag ändert, fährt die betroffene Gegenprobe (`make <name>`, Zuordnung in `harness/README.md` §Sensors); am Ende jeder Welle und vor dem Release fahren `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host.
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe `harness/README.md` §Guides) →
Verifier. Kein Self-Review — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).
