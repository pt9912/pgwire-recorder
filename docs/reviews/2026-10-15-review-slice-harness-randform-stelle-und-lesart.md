# Review-Report: slice-harness-randform-stelle-und-lesart — 2026-10-15

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-harness-randform-stelle-und-lesart`, Diff `88f9c5e..95a5bc2` (Absatz in `AGENTS.md` §3.12, Satz in `.claude/agents/architect.md`, Satz in `.claude/commands/implement-slice.md`, Plan und Roadmap).

**Skill:** `.harness/skills/reviewer.md` @ 95a5bc2 ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-randform-stelle-und-lesart` (§1 Ziel und Abgrenzung, §3, §4, §6, §7 Belege des Implementers, §8)
- Entscheidung des Nutzers vom 2026-10-15 (Punkte 1 und 2: Absatz zur Stelle einer Nachricht samt Reihenfolge der Ablehnungsgründe und Zeiger im Architect; ein Satz zu „Lesart“ in der Randform-Rückgabe; kein Sensor)
- Register: `BEO-REPO/spec-randform-erst-im-review-entschieden`, `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`, `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`
- Keine aktive ADR, keine `LH-*`-Kennung (der Slice ändert Regeln des Harness)
- `AGENTS.md` (Hard Rules) §3.7 bis §3.14

Eigene Läufe: `git diff --numstat 88f9c5e..HEAD -- AGENTS.md .claude` (Summe 20, deckt sich mit §7 und §8); `make docs-check` (Ergebnis unter Negativbefunde).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | §6 des Plans nennt als Randform, dass nicht berührte Stellen mit dem Satz „nicht berührt“ in §6 entfallen, und setzt „den Wortlaut in `AGENTS.md` §3.12“ als Entscheidungsort. Der gelieferte Absatz enthält diesen Satz nicht; die Entscheidung steht nur im Plan. | `AGENTS.md` §3.12 (Entscheidungsort der Randform), §3.9 | `docs/plan/planning/in-progress/slice-harness-randform-stelle-und-lesart.md` · „entfallen mit dem Satz „nicht berührt“ in §6, nicht leer“ gegen `AGENTS.md` · „nennt §6 je Nachrichtenart die zulässigen Stellen in der Folge“ | nein — kein Gate vergleicht §6 eines Plans mit dem Wortlaut der Regel | Randform des Plans ohne Träger im genannten Entscheidungsort |
| F-2 | LOW | §3 des Plans trägt weiter „Entwurf des Wortlauts“ und die Entwurfsfassungen; der gelieferte Wortlaut weicht ab (Plural „Nachrichten“ im Falsch-Beispiel, „und die Spezifikation bindet sie“, im Satz für `architect.md` fehlt „fehlt sie, meldest du sie dem Planner“, im Satz für `implement-slice.md` steht „als eigene Feststellung“). §7 zählt die Abweichungen auf, §3 folgt dem Diff nicht. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-harness-randform-stelle-und-lesart.md` · „**Entwurf für `.claude/agents/architect.md`:**“ gegen `.claude/agents/architect.md` · „Zur Liste gehört die Stelle jeder Nachricht“ | nein — `make kopf-check` prüft nur den Kopf, nicht §3 gegen den Diff | Plan §3 folgt dem gelieferten Wortlaut nicht |
| F-3 | INFO | Die Läufe (`make gates`, `make docs-check`) und der Befund zur Flanke stehen nur „im Bericht“ des Implementers; §7 verweist mit „Läufe: siehe Bericht“ auf ein Artefakt, das dem Review nicht vorliegt. Eine Gate-Aussage zu diesem Stand ließ sich hier nicht aus §7 belegen. | Maintainability (Reviewer-Skill, Kontext-Eingang: Belege des Implementers in §7) | `docs/plan/planning/in-progress/slice-harness-randform-stelle-und-lesart.md` · „Läufe: siehe Bericht.“ | nein — Gate-Läufe des Verifiers würden den Stand bestätigen | Beleg der Läufe nur im Bericht |
| F-4 | INFO | Flanke `p7-int-treiber`, Beschreibung unter „Flanke der Versionsmatrix“ unten: Der Diff berührt weder Gate noch Gegenprobe; wiederholt sich das Rot, steht das Register bei 3 und die Flanke ist nach dem Vermerk dort eine Lücke mit eigenem Slice. Weder §7 noch §6 des Plans führen das Auftreten oder einen Ausgang. | `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig` (Stand: „ein drittes Rot macht die Flanke zur Lücke“) | `docs/plan/planning/observations/BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig/state.md` · „ein drittes Rot macht die Flanke zur Lücke und braucht einen eigenen Slice“ | ja — `make test-postgres-versionen-gegenprobe` und `make gates`, wiederholt; ein Rot ohne Änderung an Gate oder Gegenprobe bestätigt es | Flanke der Gegenprobe färbt zusätzlichen Fall |
| F-5 | INFO | Der Absatz in `AGENTS.md` hat ein Falsch/Richtig-Paar nur für die Stelle einer Nachricht (deckt sich mit F-1 des Verbindungsende-Reviews); für die Reihenfolge der Ablehnungsgründe gibt es keins. Zudem nennt der Absatz „die zulässigen Stellen“ mit „je Stelle eine Entscheidung“; ob eine unzulässige Stelle ebenfalls eine Entscheidung braucht, sagt der Wortlaut nicht. | `AGENTS.md` §3.12 | `AGENTS.md` · „nennt §6 je Nachrichtenart die zulässigen Stellen in der Folge“ | nein | Regelabsatz mit Mehrdeutigkeit ohne Beispiel für die zweite Hälfte |
| F-6 | INFO | Zwei Randbefunde ohne erwartete Aktion: (a) Der Absatz trägt „Entscheidung des Nutzers vom 2026-10-15“, der Wortlaut der Antwort („ja deine Vorschläge“) steht im Plan und im Register, nicht im Absatz; `architect.md` verlangt „mit Datum und Wortlaut im Absatz, der sie trägt“. Die vorhandenen Absätze in §3.11 bis §3.14 haben dieselbe Form. (b) Der Satz in `implement-slice.md` setzt bei §7 an, also nach dem Code; was mit schon committetem Code geschieht, regelt der Vorsatz „Ein Absatz im Bericht … ersetzt die Rückgabe nicht, wenn der Commit sie schon trägt“, der Satz widerspricht ihm nicht. | `.claude/agents/architect.md` (Entscheidung des Nutzers), Maintainability | `AGENTS.md` · „Die Stelle einer Nachricht gehört dazu“ ; `.claude/commands/implement-slice.md` · „Schreibst du in §7 „Lesart“, „Auslegung“ oder „§6 nennt nicht““ | nein | Wortlaut der Nutzerentscheidung nicht im Absatz |

### Flanke der Versionsmatrix (`p7-int-treiber`), nur beschrieben

- Stand im Register: Die Mutation `p7-int-treiber` (ersetzt im Treiber `trap 'beenden 130' INT` durch `trap 'kill -KILL $$' INT`) soll nur die Fälle ihrer Gruppe `p7-sig` färben (`p7-sigINT-gruppe`, `p7-sigINT-skript`; Zeile in `tools/test/postgres-versionen-gegenprobe.sh`). Zweimal färbte sie im ersten vollen `make gates` zusätzlich einen TERM-Fall rot (`p7-sigTERM-gruppe` bei `slice-v1-abschluss-postgres-versionen-zusage`, `p7-sigTERM-skript` bei `slice-v1-abschluss-protokollrand`); Einzellauf und zweiter voller Lauf waren jeweils grün, beide Slices berührten Gate und Gegenprobe nicht, die Ursache ist nicht untersucht. Das Register zählt 2, ohne Adresse.
- Befund des Implementers laut Auftrag: dritte Beobachtung derselben Art, zusätzlich rot `p7-sigTERM-skript`, an einem Slice, der nur Prosa unter `AGENTS.md` und `.claude/` ändert. Das konnte ich nicht gegen einen Beleg lesen (nicht in §7, Bericht liegt nicht vor) und nicht selbst nachfahren, weil kein Gate-Lauf Teil dieses Reviews war; hier gilt die Aussage als unbestätigt übernommen.
- Lesart, wenn sie zutrifft: Weil der Slice Gate und Gegenprobe nicht berührt, spricht das Muster (erster voller Lauf rot, Einzellauf und Wiederholung grün, immer ein TERM-Fall neben einer INT-Mutation) für eine zeitabhängige Kopplung zwischen dem Abbruch des Treibers per `SIGKILL` unter INT und einem nachlaufenden TERM-Fall der Gruppe `p7-sig`, nicht für einen Fehler dieses Slice. Das ist eine Beschreibung des Musters, keine Ursache; belegt ist sie nicht.
- Folge für das Register nach dessen eigenem Wortlaut: Mit einem dritten Auftreten erreicht die Beobachtung die Schwelle; sie ist dann keine Notiz mehr, sondern braucht einen eigenen Slice. Das trägt die Closure dieses Slice als Register-Eintrag (Beleg `evidence/…`), nicht der Diff. Nicht Gegenstand dieses Diffs, kein Fix hier.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `AGENTS.md` §3.12 (Absatz *Die Stelle einer Nachricht gehört dazu*) | geprüft; Wortlaut deckt Beschluss Punkt 1 (Stellen, Entscheidung je Stelle, Spezifikation trägt sie, Reihenfolge der Ablehnungsgründe) ohne Zusatz und ohne Lockerung; Stellung nach *Aufrufarten des Treibers*, vor der Begründung; Befunde F-1, F-5, F-6 |
| `AGENTS.md` §3.12 (Retirement-Check, die eine gelöschte Zeile) | geprüft, ohne Befund; die Zeile wird nur verlängert („…gate-kette-schnell, für die Stelle einer Nachricht seit slice-harness-randform-stelle-und-lesart?“), keine Zusage entfällt, keine Bedingung wird enger |
| Falsch/Richtig-Paar | geprüft, ohne Befund bis auf F-5; das Falsch-Beispiel entspricht F-1 aus dem Review zu `slice-v1-abschluss-verbindungsende-format` (`ready_for_query` und `connection_closed`, Stelle im Code entschieden), das Richtig-Beispiel ist für den Vertrag nachprüfbar |
| §3.7 (indikativ, keine Chronik) für alle drei neuen Stellen | geprüft, ohne Befund; Herkunft steht als `seit …`-Anker, keine Konjunktive über verworfene Alternativen |
| §3.8 (ADR ohne Regeldetails) | geprüft, ohne Befund; keine ADR berührt |
| `.claude/agents/architect.md` | geprüft, ohne Befund; der Satz „Zur Liste gehört …“ folgt unmittelbar auf „Du prüfst die Liste in §6 … Fehlt in §6 eine … meldest du sie dem Planner“ und ist so Teil dieser Prüfung der Liste aus §6; der gestrichene Halbsatz „fehlt sie, meldest du sie dem Planner“ ist vom Satz davor getragen, kein Verlust der Zusage |
| `.claude/commands/implement-slice.md` Randform-Rückgabe | geprüft, ohne Befund bis auf F-6 (b); Zeiger auf `AGENTS.md` §3.12 löst auf, kein Widerspruch zu „Kein Code-Commit trägt eine neue Randform in §6“, „Auch eine Lesart ist eine Entscheidung“ und der Rückgabe an den Architect; „als eigene Feststellung“ setzt die Randform „Lesart im Zitat“ aus §6 um |
| Kein Sensor, keine Dateien außerhalb des Beschlusses | geprüft, ohne Befund; Diff berührt nur `AGENTS.md`, `architect.md`, `implement-slice.md`, Plan, Roadmap; `.harness/skills/`, `plan-welle.md`, `planner.md`, `implementer.md` unberührt (Abgrenzung §1 eingehalten) |
| Größe und Zählung (`AGENTS.md` §3.14, §3.13, §8) | geprüft, ohne Befund; Summe aus Einfügungen und Löschungen 20 (14+1, 2+0, 3+0), unter Meldeschwelle 80 und Marke 100; zwei Liefer-Punkte, zwei Schichten, keine Sendung an einen anderen Nehmer, daher keine Nachzählung bei Dritten |
| Register-Zuweisung (`AGENTS.md` §3.13, Adresse nimmt an) | geprüft, ohne Befund; beide Beobachtungen stehen als `geplant` mit der Kennung dieses Slice, §1 schließt sie nicht aus, die DoD führt sie; die Risiko-Ausgänge in §6 sind *weiter offen* mit Zielbeobachtung oder *entfallen*, ohne Kennung eines anderen Slice |
| Roadmap-Zeile „In Arbeit“ | geprüft, ohne Befund; nennt Zustand und Fundort, keine Chronik |
| Move-Commits (`AGENTS.md` §3.3) | geprüft, ohne Befund; `next/` nach `in-progress/` steht in einem eigenen Commit (`1d3e6bf`), der Inhalt der Datei ändert sich danach |
| Suppression (`//nolint`) und Gate-Lockerung | geprüft, ohne Befund; kein Code, kein Gate, kein Profil im Diff |
| `make docs-check` | geprüft, grün mit diesem Report (siehe Übergabe) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Randform des Plans ohne Träger im genannten Entscheidungsort · Plan §3 folgt dem gelieferten Wortlaut nicht · Beleg der Läufe nur im Bericht · Flanke der Gegenprobe färbt zusätzlichen Fall · Regelabsatz mit Mehrdeutigkeit ohne Beispiel für die zweite Hälfte · Wortlaut der Nutzerentscheidung nicht im Absatz

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Der Wortlaut bleibt im Rahmen des Beschlusses vom 2026-10-15, nichts ist darüber hinaus ergänzt oder gelockert; F-1 und F-2 sind Nacharbeit am Plan (§3.9), keine Änderung an Regeln. Zu F-2 hält das Register `BEO-REPO/plan-folgt-korrektur-nicht` die Klasse als wiederkehrend; die Einstufung LOW beruht darauf, dass §7 die Abweichungen vollständig nennt und der Wortlaut inhaltlich gleich bleibt.

**Übergabe:** Findings gehen an den Implementer (F-1, F-2: Plan; F-3: §7 vor der Closure); F-4 geht in die Closure und ans Register (die Entscheidung über einen eigenen Slice gehört dem Planner, nicht diesem Slice); F-5 und F-6 ohne erwartete Aktion. Die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
