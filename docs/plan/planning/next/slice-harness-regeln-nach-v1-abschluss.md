# Slice slice-harness-regeln-nach-v1-abschluss: Regeln aus dem Verbindungsende-Block — Abdeckungstabellen zählen nicht als Schicht, Wortlaut-Suche und Arbeitsbaum vor E2E-Mutanten, offene Fragen zu §3.10

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Angelegt am 2026-10-16 nach Entscheidung des Nutzers vom selben Tag („ja“ zu den Vorschlägen des Planners: der Regel-Slice wird angelegt, eingeordnet hinter Schritt 42 von welle-v1-abschluss und vor `slice-harness-schaetzung-planner-faktor`; seit dem 2026-10-16 hinter dem Release, siehe §4 *Start*); eingeordnet in `next/`, weil er mit einer Stelle in der Reihenfolge eingeplant ist (§4 *Start*, Drift-Log der Roadmap). Den **Wortlaut** hat der Nutzer nicht entschieden (§1).

**Bezug:** kein Anspruch des Produkts — der Slice ändert Regeln des Harness. Bindung an das Register: `BEO-REPO/schichtteilung-je-plan-verschieden`, `BEO-REPO/zahlwort-in-zusage-ohne-fall-je-seite` (1×), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Regelstellen: `AGENTS.md` §3.13 (seit slice-lint-bestand-kern-driven), §3.10 (seit welle-extended-query) und `.claude/commands/implement-slice.md`. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Regeln, die der Block der Slices zum Verbindungsende (Schritte 35a bis 42 von welle-v1-abschluss) als Lücken benannt hat, stehen im Briefing und in den Commands: (1) `AGENTS.md` §3.13 sagt, dass von `make abdeckung` erzeugte Tabellen in `docs/user/abdeckung-*.md` nicht als Schicht zählen, (2) `.claude/commands/implement-slice.md` verlangt in Schritt 17 die Suche des alten Wortes über das gesamte `internal/` nach einer Wortlaut-Entscheidung und in Schritt 19 den committeten Arbeitsbaum vor einem E2E-Mutanten, (3) `AGENTS.md` §3.10 beantwortet die vier offenen Fragen aus dem Register oder streicht sie.

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Regeländerungen in `AGENTS.md` und in `.claude/commands/` sind eine Entscheidung des Nutzers. Der Nutzer hat am 2026-10-16 („ja“) bestätigt, dass dieser Slice mit diesem Gegenstand angelegt und an der genannten Stelle eingeordnet wird; **den Wortlaut hat er nicht entschieden**. Der Implementer legt den Wortlaut aus §3 (oder seine Überarbeitung) dem Nutzer vor und committet beide Dateien erst nach der Bestätigung, **Nutzerentscheidung über den Wortlaut vor der Umsetzung**; bis dahin ist er eine Empfehlung. Was der Nutzer nicht bestätigt, nimmt der Implementer nicht auf und nennt es in §7.

**Beleg-Stand je Regel** (Stand 2026-10-16, gemergter Stand des Registers; die Belege stehen dort und in den Closure-Notizen der Slices, nicht hier):

1. *Abdeckungstabellen zählen nicht als Schicht* (`AGENTS.md` §3.13, Satz analog zu den Meldungscode-Konstanten, mit Begründung und Retirement-Check `seit slice-harness-regeln-nach-v1-abschluss`; „Handbuch und README zählen weiter“): `slice-v1-abschluss-verbindungsende-replay-extended-flush` §7, Frage (b): Die Zählung von `docs/user/` als dritte Schicht war offen; der Slice zählte konservativ mit Markierung, zwei ohne; `slice-v1-abschluss-verbindungsende-upstream-extended` §8 zählt ebenso. Register: `BEO-REPO/schichtteilung-je-plan-verschieden`.
2. *Wortlaut-Suche über das gesamte `internal/`* (`implement-slice.md` Schritt 17): Belege `slice-v1-abschluss-verbindungsende-replay-extended-flush` §7, Lerneintrag (1): Der alte Wortlaut „Sync“ stand nach der Entscheidung des Architect in Kommentaren dreier Pfade; Review F-1 und F-2 und Verifikation V-1 und V-2 fanden ihn je in einem anderen Schritt (ein Auftreten, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). **Zweiter Beleg** (Beschluss des Nutzers vom 2026-10-16, „ja deine Empfehlungen“, Vorschlag aus der Closure von `slice-v1-abschluss-verbindungsende-einspielen-extended`, Lerneintrag (1)): Dort stand „bis zum ReadyForQuery“ in Kommentaren dreier Pfade derselben Funktion; die Nacharbeit zu F-3 zog den Pfad des Slice nach, die Verifikation (V-1) fand `Play` und drei Stellen im Pfad mit Vergleich, indem sie das Wort über den Einspiel-Pfad suchte. *Arbeitsbaum vor E2E-Mutanten committen* (Schritt 19, Weg der Mutanten): dort Lerneintrag (4), ein Auftreten (`git checkout -- internal` verwarf eine nicht eingecheckte Änderung am Adapter-Test), kein Register-Eintrag.
3. *Offene Fragen zu `AGENTS.md` §3.10*, alle vier **Vorschläge** und Entscheidung des Nutzers vor der Umsetzung:
   - *Zahlwort in einer Zusage* („genau eine“, „jede andere Zahl“, „alle Zeilen“): je Seite ein Fall und eine Mutation, die nur sie bricht. Register: `BEO-REPO/zahlwort-in-zusage-ohne-fall-je-seite` (1×, `slice-v1-abschluss-verbindungsende-sqlite-schreiber` V-1); verwandt, dort nicht mitgezählt: `slice-v1-abschluss-sessions` (V-305, F-685) und `slice-v1-abschluss-protokollrand` (F-5, V-3) bei `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; zählte der Nutzer sie mit, wären es 3×.
   - *Das Feld einer Aufzählung trägt keine Fixture, die ein anderes Element prüft*: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, Beleg `slice-v1-abschluss-verbindungsende-format-schreiber` (F-1, F-5), ein Beleg dieser Form; die Frage steht dort offen.
   - *Das Glied einer Konjunktion in einer Bedingung, die eine Zusage trägt, als Element der Aufzählung*: `slice-v1-abschluss-verbindungsende-replay` §7, Lerneintrag (1) (F-1: der Mutant, der das Glied entfernt, blieb grün, bis das Review ihn baute); ein Beleg, Empfehlung des Planners damals: abwarten bis zum zweiten.
   - *Die Mutationstabelle in §7 wird vor der Übergabe gegen §6 gelesen* (Zeile je Element, jede Zeile nennt einen Eingriff, den der Code hat, und einen Test, den der Nachbau färbt): Belege bei `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: `slice-v1-abschluss-verbindungsende-sqlite` (F-2, V-1, V-2: Zeile nannte eine Umwandlung, die es nicht gibt), `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (F-2: Paare als eine Zeile; eine Zeile nannte einen roten Test, den der Nachbau nicht färbt) und `slice-v1-abschluss-verbindungsende-replay` (F-4, V-2: zwei Zeilen nannten einen Eingriff, den der Produktcode nicht hat); mindestens drei Auftreten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Sensor, der Schichten, Mutationstabellen oder Wortlaut prüft (etwa eine Erweiterung von `kopf-check` oder ein Gate gegen Kommentarprosa) — anderer Vorgang (Werkzeug statt Regel): `slice-harness-mutation` trägt den Sensor für Kommentarprosa; ein Sensor hier bekäme einen Vertrag in `harness/vertraege/` und eine Gegenprobe.
- Die Regel zur Schätzung, zur Meldeschwelle und zur Zählung nach der Obergrenze (`AGENTS.md` §3.14) — `slice-harness-schaetzung-planner-faktor` (Adresse; sein Wortlaut steht dort): ein anderer Absatz, ein anderer Nutzer-Entscheid.
- Änderung an `.claude/agents/`, `.harness/skills/`, `plan-welle.md`, `close-welle.md` und an den Hooks — Schicht-Abgrenzung und Größe: Die vier Regeln gehören in das Briefing und in den Command des Implementers; eine dritte Datei brächte einen vierten Liefer-Punkt.
- Der Abgleich der Sätze einer Stelle der Spezifikation durch den Architect (`.claude/agents/architect.md`, `AGENTS.md` §3.12; Beschluss des Nutzers vom 2026-10-16, Register `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut`) — `slice-harness-regeln-spezifikation-satzabgleich` (Klasse 1, Schnitt des Planners, 2026-10-16): Mit ihm wären es vier Liefer-Punkte in drei Schichten (Briefing, Commands, Agenten-Rollen), über beiden Grenzen; der Nehmer nimmt die Sendung in §1 (*Geschnitten aus*), DoD und §3 an.
- Rückwirkendes Nachzählen der Schichten in bestehenden Plänen — Bestand bleibt bewusst stehen: Deren Zählungen stehen in ihrem §8 mit Markierung; gezählt wird neu, wenn ein Plan angefasst wird.
- Produkt-Code, Gates, Spezifikation, Benutzerhandbuch — Schicht-Abgrenzung: Der Slice ändert zwei Dokumente des Harness und kein Verhalten des gebauten Binary; **kein Produktcode**.
- Eine Aussage über Handbuch und README als Schicht — bleibt: Beide zählen weiter als Nutzer- und Wartungs-Doku (`AGENTS.md` §3.13, Entscheidung des Nutzers vom 2026-10-10); der Satz grenzt nur die erzeugten Tabellen aus.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** `AGENTS.md` §3.13 trägt den Satz: Von `make abdeckung` erzeugte Tabellen in `docs/user/abdeckung-*.md` zählen nicht als Schicht (Handbuch und README zählen weiter), analog zu den Meldungscode-Konstanten, mit Begründung und Retirement-Check `seit slice-harness-regeln-nach-v1-abschluss`. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] **Liefer-Punkt 2:** `.claude/commands/implement-slice.md` Schritt 17: Wird eine Wortlaut-Entscheidung nachgezogen, wird das alte Wort über das gesamte `internal/` gesucht (auch `ports/`), jeder Treffer mit Entscheidung in §7; Schritt 19: Der Arbeitsbaum ist vor E2E-Mutanten committet. Derselbe Vorbehalt der Bestätigung.
- [ ] **Liefer-Punkt 3:** `AGENTS.md` §3.10 trägt, was der Nutzer von den vier Vorschlägen aus §1 bestätigt (Zahlwort, Fixture-Träger, Konjunktionsglied, Tabelle gegen §6 lesen), je mit Paar Falsch/Richtig und Retirement-Check `seit slice-harness-regeln-nach-v1-abschluss`; ein nicht bestätigter Vorschlag steht nicht in `AGENTS.md` und ist in §7 mit „nicht aufgenommen“ und dem Beleg-Stand genannt. Dasselbe Verfahren für den Wortlaut.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage von §3.10 oder §3.13; der Diff von `AGENTS.md` ergänzt (Beleg: `git diff` ohne Löschzeile in den Regeln selbst, in §7).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-regeln-nach-v1-abschluss` in `AGENTS.md` §3.13 und §3.10 (`liegt in`), soweit aufgenommen.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: für jede aufgenommene Regel der Eintrag von offen/`geplant` auf `verkörpert` (Zielort und Anker); für jeden nicht aufgenommenen Vorschlag der Zähler unverändert und der Beleg-Stand genannt; `BEO-REPO/schichtteilung-je-plan-verschieden` mit Vermerk.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.13 | update (ergänzen) | Satz zu den erzeugten Abdeckungstabellen nach dem Absatz zu den Meldungscode-Konstanten; Paar Falsch/Richtig; Begründung und Retirement-Check |
| `AGENTS.md` §3.10 | update (ergänzen, soweit bestätigt) | vier Vorschläge aus §1, je ein Satz mit Paar Falsch/Richtig |
| `.claude/commands/implement-slice.md` | update | Schritt 17 (Wortlaut-Suche) und Schritt 19 (Arbeitsbaum vor E2E-Mutanten) |
| `docs/plan/planning/observations/BEO-REPO/` (`state.md` der betroffenen Einträge) | update | Ausgang und Vermerke bei Closure |

Tests und Gates: keine neuen. Der Slice liefert Prosa (§3.10 verlangt Mutationen für Verträge von Code und Gates); die Wirkung zeigt das Register (Auftreten trotz der Regel).

**Entwurf des Wortlauts** (Empfehlung des Planners, keine Entscheidung des Nutzers; Nutzerentscheidung über den Wortlaut vor der Umsetzung):

> *§3.13, nach den Meldungscode-Konstanten:* Von `make abdeckung` erzeugte Tabellen in `docs/user/abdeckung-*.md` zählen nicht als Schicht; Handbuch und README zählen weiter.
>
> *`implement-slice.md`, Schritt 17, Ergänzung:* Wird eine Wortlaut-Entscheidung der Spezifikation nachgezogen, sucht der Implementer das alte Wort über das gesamte `internal/` (auch `ports/`) und nennt jeden Treffer mit Entscheidung in §7.
>
> *Schritt 19, Ergänzung:* Vor einem E2E-Mutanten ist der Arbeitsbaum committet; `git checkout -- internal` verwirft sonst nicht eingecheckte Änderungen.
>
> *§3.10, vier Sätze (je einzeln zu bestätigen):* ein Zahlwort in einer Zusage („genau eine“, „jede andere Zahl“, „alle Zeilen“) bekommt je Seite einen Fall und eine Mutation, die nur sie bricht; das Feld einer Aufzählung trägt keine Fixture, die ein anderes Element prüft; das Glied einer Konjunktion in einer Bedingung, die eine Zusage trägt, zählt als Element; die Mutationstabelle in §7 wird vor der Übergabe gegen §6 gelesen.

**Reihenfolge:** `AGENTS.md` zuerst (die Quelle), dann der Command, beide in einem Commit nach der Bestätigung; der Plan des Slice folgt in demselben Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [welle-v1-nachtrag](../welle-v1-nachtrag.md) liegt in `done/` und damit `slice-v1-abschluss-verbindungsende-doku` (Schritt 42) und `slice-v1-abschluss-verbindungsende-upstream-extended`, sodass die Belege des ganzen Blocks vorliegen (Voraussetzung der Stelle, nicht der Technik); der Slice steht in der wellenlosen Reihe hinter `slice-harness-slice-mv-exakter-treffer` und vor `slice-harness-schaetzung-planner-faktor` (Auftrag des Nutzers vom 2026-10-16) und ist seit demselben Tag hinter das Release verschoben (Entscheidung des Nutzers vom 2026-10-16, Drift-Log der Roadmap): wellenlos hinter [welle-v1-nachtrag](../welle-v1-nachtrag.md) und vor der übrigen Harness-Reihe (`slice-harness-commit-struktur-id` zuerst). Die Reihe der fünf: `slice-harness-slice-mv-exakter-treffer`, `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`, `slice-harness-gegenprobe-matrix-p7-flanke`. WIP-Limit 1. Der Implementer legt den Wortlaut dem Nutzer vor und committet `AGENTS.md` und `implement-slice.md` erst nach der Bestätigung; bis dahin hält er an und gibt an den Aufrufer zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- AGENTS.md .claude`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 200** geht der Slice zurück zur Zerlegung. **Meldeschwelle 160:** Erreicht die Summe 160, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); die Entscheidung des Planners steht als datierte Zeile in §7. Schätzung des Planners vor dem Anspruch: **60 bis 120, Mitte 90, 45 % der Marke**; Faktor Tests zu Produktcode entfällt (kein Code, kein Test, §3); Basis der Schätzung: neun Sätze mit vier Paaren Falsch/Richtig in `AGENTS.md` (rund 15 Zeilen je Regel mit Paar und Retirement-Check) und zwei Sätze im Command.
- `in-progress` → `next`: Die Regeln verlangen eine Änderung an einer dritten Datei (`architect.md`, `planner.md`, `plan-welle.md`, ein Skill) oder einen Sensor — zurück zur Zerlegung (§1).
- `in-progress` → `open`: Der Nutzer bestätigt keinen Wortlaut und hat keine Fassung, die er bestätigen würde — der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut in `AGENTS.md` steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Bestätigung des Wortlauts durch den Nutzer in §7 mit Datum, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md` §3.13 und §3.10, `liegt in`, soweit aufgenommen).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regeln** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber Regeln, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist der Wortlaut, vor dem Commit bestätigt:

- *Abdeckungstabellen:* die Regel gilt für die von `make abdeckung` geschriebenen Tabellen `docs/user/abdeckung-*.md`; eine von Hand gepflegte Datei unter `docs/user/` ist Doku und zählt. Handbuch und README zählen weiter (Entscheidung vom 2026-10-10).
- *Wortlaut-Suche:* gesucht wird das alte Wort über `internal/` samt `ports/`; Treffer in Tests und in `docs/` nennt der Implementer, wenn er sie berührt; ein Treffer, der das Wort in anderer Bedeutung trägt (`Sync` als Nachricht des Protokolls), steht mit „richtig“ und Grund in §7.
- *Arbeitsbaum vor E2E-Mutanten:* gilt für Mutanten, die der Implementer im Arbeitsbaum fährt; Mutanten in einer Kopie außerhalb des Repos (Wege der Mutanten in den Closure-Notizen) sind ausgenommen.
- *Vorschläge aus §3.10:* jeder einzeln bestätigt oder nicht; eine Teilmenge ist erlaubt.
- *Aufruf-Voraussetzung, Aufrufarten des Treibers:* nicht berührt (`AGENTS.md` §3.12).

**Risiken:**

- Der Nutzer ändert den Wortlaut, streicht einen Vorschlag oder bestätigt keinen — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht in `AGENTS.md`; bei Ablehnung geht der Slice zurück nach `open/`, §4).
- Die Regeln sind Prosa ohne Sensor: ein Plan, der erzeugte Tabellen als Schicht zählt, oder ein Mutant im ungesicherten Arbeitsbaum bleibt unbemerkt, bis das Review ihn findet — **Ausgang:** weiter offen → `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Auftreten trotz der Regel zählt weiter; ein Sensor ist abgegrenzt, §1).
- Der Satz zu den Abdeckungstabellen widerspricht der Teilung in `slice-doku-ist-stand` (Handbuch und README zählen als Nutzer-Doku) — **Ausgang:** offen bis Closure (erwartet: entfallen, der Wortlaut grenzt nur die erzeugten Tabellen aus und nennt Handbuch und README ausdrücklich).
- Die Änderungen in `AGENTS.md` treffen auf `slice-harness-schaetzung-planner-faktor` (§3.14) und `slice-harness-regeln-spezifikation-satzabgleich` (§3.12), derselbe Rang, andere Abschnitte — **Ausgang:** offen bis Closure (erwartet: entfallen, die Abschnitte überlappen nicht; der Slice, der später schließt, rechnet den Wortlaut des ersten ein).
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 60 bis 120 gegen Marke 200) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Die Bestätigung des Wortlauts durch den Nutzer (Datum, Fundstelle, je Vorschlag), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) und das Ergebnis der Schätzung gegen die Messung stehen hier.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-16). Treffer in der Sub-Area `REPO`: `BEO-REPO/schichtteilung-je-plan-verschieden`, `BEO-REPO/zahlwort-in-zusage-ohne-fall-je-seite` (1×), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — sie sind der Auslöser (§1, Beleg-Stand); `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (9×) und `BEO-REPO/meldeschwelle-ohne-artefakt` (2×) gehören `slice-harness-schaetzung-planner-faktor`. Eine Schwelle wird mit diesem Slice nicht neu erreicht (keine Evidence-Datei entsteht).

Nachgezählt (Planner, 2026-10-16, `AGENTS.md` §3.13): Liefer-Punkte: 3 (§3.13-Satz, `implement-slice.md` Schritte 17 und 19, §3.10-Sätze). Schichten: zwei, das Briefing (`AGENTS.md`) und die Commands (`.claude/commands/`); Planung (Register, Plan) zählt keine Schicht; kein Produkt-Code, keine Nutzer-Doku. Der Slice nimmt `slice-harness-schaetzung-planner-faktor` nichts ab und trägt keine Sendung in einen Nehmer ein: keine Nachzählung bei einem anderen Slice. Schätzung in der Metrik von §4: 60 bis 120 (Mitte 90, 45 % der Marke 200); Faktor Tests zu Produktcode entfällt; der Plan liegt unter 70 % der Marke und bleibt ungeteilt.

Nachgezählt beim Prüfen der Sendung „Abgleich der Sätze einer Stelle“ (Planner, 2026-10-16, `AGENTS.md` §3.13, Beschluss des Nutzers vom selben Tag): Die Sendung würde `.claude/agents/architect.md` und `AGENTS.md` §3.12 berühren, einen vierten Liefer-Punkt und eine dritte Schicht (Agenten-Rollen). Der Slice steht bei drei Liefer-Punkten in zwei Schichten und läge damit über beiden Grenzen; **nicht eingetragen, geschnitten** (nach Lieferwert, nicht nach Schichten): Die Sendung liegt bei `slice-harness-regeln-spezifikation-satzabgleich` (zwei Liefer-Punkte, zwei Schichten, Schätzung 20 bis 45). Dieser Slice bleibt bei drei Liefer-Punkten, zwei Schichten, Schätzung 60 bis 120 (Marke 200) unverändert; der Zweitbeleg zu Punkt 2 (§1) ist ein Beleg, kein Liefer-Punkt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

Nachgezählt beim Verschieben hinter das Release (Planner, 2026-10-16, `AGENTS.md` §3.13): Start und Stelle geändert, Gegenstand und DoD unverändert; der Slice trägt keine Sendung ein und nimmt keine ab: Liefer-Punkte weiter 3, Schichten weiter zwei. Kein Slice vor dem Release nennt ihn als Voraussetzung; `slice-v1-abschluss-verbindungsende-hilfetext-binary` (`in-progress/`) nennt ihn als Adresse der offenen Frage zur Zählung der erzeugten Abdeckungstabellen und zählt bis dahin konservativ mit Markierung; dasselbe gilt für die Slices bis zum Release, die erzeugte Tabellen ändern.
