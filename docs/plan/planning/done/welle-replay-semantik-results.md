# Welle welle-replay-semantik — Replay-Semantik: Mismatch, Fehlerreplay, Meldungscodes — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-replay-semantik
**Abschluss:** 2026-10-06
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- `replay --fail-on-unconsumed` meldet nicht verbrauchte Interaktionen und Sessions beim Beenden als Fehler `PGR-E5002` mit Exit-Code 5 ([`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)); dazu die Regeln „Hilfe vor jeder Prüfung“ und „`--` beendet die Optionen“ für `record`, `replay` und den globalen Aufruf ([`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung); `slice-replay-semantik-mismatch`).
- PostgreSQL-Fehlerantworten werden für Simple Query aufgezeichnet und wiedergegeben, mit Resultsets und Transaktionen; eine Fehlerantwort vor dem Abbruch einer Verbindung ist `PGR-E6001`. Ein Dreifachvergleich (direkt, über `record`, über `replay`) hält die Sicht des Clients gleich ([`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol); `slice-replay-semantik-fehlerreplay`).
- Fehler tragen einen Meldungscode im Fehlertext, eine Fehlerkette wird per Tiefensuche eingeordnet, `--log-level` steuert die Diagnose, und die Abweichungs-Diagnose bei Extended Query nennt die Parameter-Nummer ohne Wert ([`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md); `slice-replay-semantik-meldungscodes`).
- Abnahmeszenario 4 (Replay-Abweichung, [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage)) und 6 (Fehlerreplay, [`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers)) sind für Simple Query end-to-end nachgewiesen, jede Zusage mit einer roten Mutation (Verifikation unten). Für Extended trug sie welle-extended-query; der Exit-Code beim Herunterfahren folgt in welle-v1-abschluss. M3 bleibt offen.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **`AGENTS.md` §3.12 griff im letzten Slice jedes Mal vor dem Code.** In `slice-replay-semantik-meldungscodes` standen die Randformen vor dem ersten Code-Commit in der Spezifikation (`ac12870`); die drei Rückgaben des Implementers entschied `f2c428a` vor dem Code `1348a9a`, F-401 und F-402 entschied `ba2de96` vor `7dc26b7`, V-57 entschied `a749370` vor dem Test `596aa7f`. Es ist der erste Slice unter §3.12, in dem keine Rückgabe einer Lesart folgte, die schon im Code stand (`BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` bleibt 2×). In den beiden Slices davor hielt die Reihenfolge ebenso (`001403f` vor `6e8d2dd`; `a8c490e`, `186f476` vor `212a3d4`).
- **Mutationen der Verifikation fanden, was die Deklarationen verdeckten.** V-52 (verneinende Zusage von Szenario 4 ohne tragenden Test) fand die Verifikation von `slice-replay-semantik-fehlerreplay`, obwohl `make abdeckung-check` grün war; daraus wurde Schritt 1 der Welle-Closure geschärft, und diese Closure ist die erste, die ihn fährt.
- **Der Dreifachvergleich** `dreiSichten` (direkt, record, replay) färbte jede Mutation nur im Replay oder nur beim Aufzeichnen an der richtigen Sicht rot (VF12 bis VF20).
- **Mutation über eine frische Kopie** (Bind-Mount oder neuer Pfad, `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`): Seit `slice-replay-semantik-mismatch` nennen die Berichte den Weg; auch EV3 dieser Closure lief so.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Viele Randform-Runden.** Die Reihenfolge von §3.12 hielt, die Vollständigkeit der Liste in §6 nicht: drei Runden an der Hilfe-Regel in `slice-replay-semantik-mismatch` (F-382, Rückgabe, V-45), vier Runden an Zeile 11 in `slice-replay-semantik-fehlerreplay` (11a bis 11h, darunter F-390 als echter Fehler im Extended-Pfad), zwei still im Code entschiedene Formen und eine offene in `slice-replay-semantik-meldungscodes` (F-401, F-402, V-57). Konsequenz: `BEO-REPO/spec-randform-erst-im-review-entschieden` bleibt Prosa (Steering-Loop unten, mit Begründung).
- **Ein roter Stand wurde gepusht.** Nach dem Move von `slice-replay-semantik-mismatch` nach `done/` verlinkte der Verifikationsbericht den Plan unter `in-progress/`; das Doku-Gate war rot. Behoben in `d5dbd55` (Pfad als Code-Span statt Link). Die Slice-Closure lag davor und hat den Fall nicht ins Register eingetragen; er steht dort bisher nicht (an den Nutzer gemeldet).
- **`slice-replay-semantik-mismatch` wuchs über seine Option hinaus** um die Regeln aus `LH-FA-01.a`, benannt erst durch V-46; die Nachbarregel desselben Eingangs stand nicht in §6.
- **Nacharbeit ohne eigenes Review.** In allen drei Slices lief die Nacharbeit nach Review oder Verifikation (`29d02a1`, `ce50a10`, `596aa7f`) ohne eigenes Review; die Mutationen zu `596aa7f` fuhr nur der Implementer. EV3 hat diese Closure selbst nachgefahren (Verifikation unten).
- **Nebenbefunde, in dieser Closure behandelt oder weitergegeben:**
  1. `harness/conventions.md` §Modus-Deklaration pro Sub-Area begründete Greenfield mit „ein Go-Gerüst ohne Funktion“. Die Zeile ist Harness-Doku und mit dieser Closure nachgezogen: Sie nennt den Go-Code von Record und Replay mit Tests; Greenfield bleibt, weil Spezifikation und Entscheidungen dem Code vorausgehen.
  2. Der Kommentar an `Meldungen` (`internal/hexagon/model/fehler.go`, Zeilen 171 bis 173) sagt „außen nach innen und in der Reihenfolge seiner Ursachen“; die Spezifikation sagt seit `a749370` „Tiefensuche“. Der Code ist richtig, der Kommentar folgt der Spezifikation nicht (`AGENTS.md` §3.11). Produkt-Code, in dieser Closure nicht geändert; übernommen von `slice-harness-lint` (§3, Zeile `internal/hexagon/model/fehler.go`), der den Kommentar-Bestand ohnehin anfasst.
  3. V-61 (`BEO-REPO/commit-nennt-struktur-kennung`, 1×): Commit-Messages nennen `SPEC-*`, obwohl `AGENTS.md` §5 Regel 1 es ausschließt. Sensor-Kandidat für den Commit-Hook (ausgeschlossene Kennung ablehnen); unter der Schwelle, deshalb weder geplant noch verkörpert.
  4. Die Folgepflicht aus [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) (ein Gate, das Code-Tabelle und Katalog abgleicht) hat keine Slice-Kennung: `slice-v1-abschluss-container` nennt das Gate in §1 als „Folge-Slice“ ohne Kennung. Offen an den Nutzer (Trigger-Audit unten).

## Lerneinträge der Slices

Die Steering-Loop-Einträge aus §7 der drei Slices, als Zeiger (das Original steht im Archiv dieser Welle):

- `slice-replay-semantik-mismatch` — geschärfte Regel: Ein Mutationsnachweis zählt nur, wenn der Mutant im Build angekommen ist; Mutation und Zurücksetzen laufen mit `touch`, frischem Pfad oder Bind-Mount, und der Bericht nennt den Weg — liegt in `.claude/commands/implement-slice.md Schritt 19`, Anker `seit slice-replay-semantik-mismatch`. Auslöser V-50, `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×).
- `slice-replay-semantik-fehlerreplay` — geschärfte Regel: Ein Abnahmeszenario im Closure-Trigger einer Welle gilt erst, wenn jede seiner Zusagen, auch die verneinende, ein Test hält, der unter ihrer Mutation rot wird — liegt in `.claude/commands/close-welle.md Schritt 1`, Anker `seit slice-replay-semantik-fehlerreplay`. Auslöser V-52, Klasse `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- `slice-replay-semantik-meldungscodes` — geschärfte Regel: Eine verneinende Zusage ist im Mutationsbericht eine eigene Zeile, und ihr Test liest bis zum Ende des Stroms und zählt — liegt in `.claude/commands/implement-slice.md Schritt 19`, Anker `seit slice-replay-semantik-meldungscodes`. Auslöser V-56, Klasse `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.

Eingesammelt werden mit dieser Welle auch die beiden wellenlosen Slices seit der letzten Closure; ihre Lerneinträge: `slice-harness-randformen-vor-code` — Guide geschärft, `AGENTS.md` §3.12; `slice-harness-kopf-sensor` — neuer Sensor `make kopf-check` für die Kopf-Hälfte von §3.9.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Lese-Schritt am Register `../observations/BEO-REPO/` (Stand `e05e4bd`): vier Einträge stehen über der Schwelle, alle schon als Prosa verkörpert und in jedem Slice dieser Welle wieder aufgetreten. Nach v6.13.0 · `regelwerk/modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3, gilt die Prosa-Form damit als ausgeschöpft: Je Eintrag folgt ein Sensor oder die Begründung, warum keiner möglich ist. Entschieden vom Nutzer am 2026-10-06.

- **Sensor** geplant: Ein Mutations-Gate (Werkzeug für Mutationstests, Docker-only, per Digest gepinnt, mit Schwelle für den Bestand oder die geänderten Pakete) fängt Zusagen, die kein Test hält — Kennung `slice-harness-mutation`. Die Prosa-Regeln `AGENTS.md` §3.10 und §3.11 bleiben verkörpert.
  Auslöser: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (11×; in dieser Welle slice-replay-semantik-mismatch, slice-replay-semantik-fehlerreplay, slice-replay-semantik-meldungscodes, dazu slice-harness-kopf-sensor und slice-harness-randformen-vor-code) und `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (10×; in dieser Welle die drei Slices und slice-harness-kopf-sensor). Retirement-Check von §3.10 und §3.11: wieder aufgetreten, beide bleiben.
- **Guide** bestätigt, kein weiterer Sensor: Der Slice-Plan folgt jeder Korrektur — liegt in `AGENTS.md §3.9`.
  Auslöser: `BEO-REPO/plan-folgt-korrektur-nicht` (10×; in dieser Welle die drei Slices und slice-harness-kopf-sensor). Die maschinell entscheidbare Hälfte hält `make kopf-check` (seit slice-harness-kopf-sensor). **Warum kein weiterer Sensor möglich ist:** Was danach auftrat, lag in Folge-Slice-DoDs, in der Vollständigkeit der Sicht-Stellen, in §3-Begründungen und im Welle-Plan; ob ein Plan dem Diff inhaltlich folgt, ist Urteil über Bedeutung, kein Abgleich von Kennungen. Review und Verifikation fanden jeden Fall vor der Closure (F-384, F-386, F-388, V-46 bis V-48; F-395, F-396, V-53, V-54; F-400, F-403, F-405). Retirement-Check von §3.9: wieder aufgetreten, die Regel bleibt.
- **Guide** bestätigt, kein Sensor: Randformen eines neuen Vertrags sind vor dem Code entschieden — liegt in `AGENTS.md §3.12`.
  Auslöser: `BEO-REPO/spec-randform-erst-im-review-entschieden` (8×; in dieser Welle die drei Slices und slice-harness-kopf-sensor). **Warum kein Sensor möglich ist:** Die Reihenfolge hielt in jeder Runde (oben, *Was hat funktioniert*); was wiederkehrt, ist eine Randform, die §6 nicht nennt. Ob eine Liste von Randformen vollständig ist, kann keine Maschine entscheiden — eine nicht genannte Form ist für sie unsichtbar. Review und Verifikation fanden jede vor der Closure (F-382, V-45; F-390, F-391; F-401, F-402, V-57). Retirement-Check von §3.12: wieder aufgetreten, die Regel bleibt.

Weitere Einträge über der Schwelle: keine. Alle übrigen stehen bei 1× oder 2× (`randform-im-code-entschieden-dann-zurueckgegeben` und `werkzeug-commit-ohne-zugelassene-kennung` je 2×) und werden in dieser Closure nicht gelesen.

**Lücke der Ablage:** Die Begründung „kein Sensor möglich“ hat in `state.md` keine Form. Die geschlossene Menge der Ausgänge (`observations/README.md`, *Die drei Ausgänge*) kennt `verkörpert`, `geplant` und `gestrichen`; die Zeile `**Sensor:**`, die `plan-folgt-korrektur-nicht` führt, trägt einen dieser Werte. Die `state.md` von `plan-folgt-korrektur-nicht` und `spec-randform-erst-im-review-entschieden` bleibt deshalb unverändert (`verkörpert`); die Begründung steht nur hier. Ob die Ablage dafür eine Form bekommt, ist an den Nutzer gemeldet.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht im Register `../observations/` (eine Datei je Auftreten unter `evidence/`).
Was in dieser Welle **3×** erreicht hat, steht oben unter
*Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- `slice-harness-mutation` (Mutationstests als Gate) — neu, ohne Welle, nach `slice-harness-coverage`.
- `slice-harness-lint` — übernimmt den Kommentar an `Meldungen` (Nebenbefund 2).
- `slice-v1-abschluss-betrieb` (Frist und weiteres Signal, Schlüssel `fail_on_unconsumed` und `log_level`, Hilfe vor Prüfung für `config show`, allgemeiner Leser) — welle-v1-abschluss.
- `slice-v1-abschluss-einspielen` (`play`, `--log-level` bei `play`, unerreichbare Teile von `LH-FA-20.a`) — welle-v1-abschluss.
- `slice-v1-abschluss-sessions` (`connection` bei Sessions ohne Interaktion) — welle-v1-abschluss.
- `slice-v1-abschluss-antwortvergleich` (unerreichbare Teile von `LH-FA-24.a`) — welle-v1-abschluss.

Warum der neue ohne Welle: Seine Closure-Bedingung ist seine DoD; kein Abnahmeszenario und kein Meilenstein hängt an ihm.

Bestand in `open/` (Lese-Schritt, Schritt 3), erst gruppiert: Die acht wellenlosen Harness-Slices (`slice-lastenheft-pruefbarkeit`, `slice-harness-lint`, die vier `slice-harness-blackbox-*`, `slice-harness-coverage`, `slice-harness-mutation`) teilen den Auslöser (Entscheidungen des Nutzers vom 2026-10-05 und 2026-10-06) und sind schon eine geordnete Reihe; jeder trägt einen eigenen Gegenstand, kein Konsolidieren. Die übrigen 15 tragen eine `Welle:` einer offenen Welle (welle-v1-abschluss, welle-erster-release). Alle **bestätigt**: Ihre Priorität trägt die Reihenfolge der Wellen und der Reihe.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- **Slices:** alle drei (`slice-replay-semantik-mismatch`, `slice-replay-semantik-fehlerreplay`, `slice-replay-semantik-meldungscodes`) lagen bei Closure-Beginn (`e05e4bd`) flach in `done/`; `next/` und `in-progress/` ohne Slice.
- **`make gates`** Exit 0 auf `e05e4bd`: Build, Stufe `test` (`gofmt`, `go vet`, Unit-Tests; aus dem Build-Cache desselben Quellstands), Integration in zwei Phasen (`run-integration-tests: gruen`, 47 E2E-Tests `PASS`, 0 `FAIL`), Architektur-Gate mit Gegenprobe, Abdeckung mit Gegenprobe, Commit-Träger-Gegenprobe, Kopf-Gegenprobe, Doku-Gate (243 Dateien, 0 Befunde), Baseline (`v6.13.0`, 54 Dateien).
- **Das *Mehr* aus §3 der Welle**, aus demselben Lauf: Szenario 4 — `TestE2EReplayAbweichung` PASS, `TestE2EReplayNachDemEnde` PASS; Szenario 6 — `TestE2EFehlerreplayEinfach` PASS (alle Simple Query über `pgconn.Exec`).
- **Rote Mutation je Zusage** (`.claude/commands/close-welle.md` Schritt 1, seit slice-replay-semantik-fehlerreplay):
  - Szenario 4, abweichende Anfrage, bejahend und verneinend: EV1 und EV2 rot in `TestE2EReplayAbweichung`, Gegenprobe gegen den Test aus `852c6d0` grün (Verifikation von `slice-replay-semantik-meldungscodes`, `docs/reviews/2026-10-06-verifikation-slice-replay-semantik-meldungscodes.md`, Abschnitt 2 und 3; V-52 geschlossen).
  - Szenario 4, Anfrage nach dem Ende der Aufzeichnung (EV3, verneinende Zusage „kein Ergebnis“): von dieser Closure nachgefahren, weil ihn bisher nur der Implementer gefahren hatte. Frische Kopie per `git archive e05e4bd` unter neuem Pfad; Mutation: `ReplayService.Query` liefert nach der letzten Interaktion deren Antworten ohne `ReadyForQuery` neben `PGR-E5001`, der PGWire-Adapter sendet sie vor dem Fehler. `make test-integration` in der Kopie: Exit 2, genau ein Test rot — `TestE2EReplayNachDemEnde`: „erwartet PGR-E5001 ohne Ergebnis, erhalten 1 Ergebnisse, FATAL: Replay [PGR-E5001]: Session 1: nach Interaktion 1 erwartet die Aufzeichnung keine weitere“. Der richtige Grund; ohne Mutation grün im Gate-Lauf oben. V-60 ist damit E2E geschlossen.
  - Szenario 6: ES1 (SQLSTATE vertauscht), ES2 (Meldungscode an wiedergegebener Fehlerantwort, verneinende Zusage aus `LH-FA-11.a`) und ES3 (nur im Replay) rot in `TestE2EFehlerreplayEinfach` an `7dc26b7` (dieselbe Verifikation, Abschnitt 3); VF12 bis VF20 rot an der richtigen Sicht (`docs/reviews/2026-10-06-verifikation-slice-replay-semantik-fehlerreplay.md`).
- **Verifikation der Slices:** `slice-replay-semantik-mismatch` — 23 von 23 Mutationen rot (VM01 bis VM23), `make gates` grün an `29d02a1`; `slice-replay-semantik-fehlerreplay` — 22 von 23 rot (VF24 grün außerhalb des Slice), Liefer-Punkte 1 bis 3 bestätigt; `slice-replay-semantik-meldungscodes` — 56 von 64 rot, Liefer-Punkte 1 und 3 bestätigt, 2 mit Einschränkung V-58, behoben in `596aa7f`. Validierung: keine in dieser Welle; die Rollen-Sequenz sieht sie vor größeren Wellen vor, welle-v1-abschluss ist die nächste.
- **Trigger-Audit:**
  - Carveouts: 0 offen (`docs/plan/carveouts/` trägt nur `.gitkeep`).
  - Bootstrap-aware Gates: keine Reifestufe im Repo (`harness/mk/`, `harness/README.md` §Sensors ohne Stufe), 0 offen. Die Stufe `testpackage` entsteht erst mit `slice-harness-lint`.
  - [ADR-0007](../../adr/0007-strict-replay.md) — bestätigt: Diese Welle hat die Strenge verschärft (`--fail-on-unconsumed`), kein toleranter Modus ist als Anforderung aufgenommen; die eine Klasse, für die der Trigger eintrat, bedient [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md).
  - [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md) — bestätigt: Die Bibliothek deckt alle Nachrichtentypen dieser Welle ab. Dass sie leere Felder nicht sendet (leere Meldung `M`), ist ein akzeptiertes Negativ in `slice-replay-semantik-fehlerreplay`, kein fehlender Nachrichtentyp; der E2E-Test liest PGWire ohne die Bibliothek, damit er ihr Verhalten nicht voraussetzt.
  - [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) — erste Hälfte des Triggers eingetreten (Codes mehrerer Klassen in einer Fehlerkette), geprüft in `slice-replay-semantik-meldungscodes` §6 und im Review: Die Entscheidung bleibt, die Rangregel (Tiefensuche) ist Spezifikationsdetail nach `AGENTS.md` §3.8, keine Folge-ADR. Zweite Hälfte nicht eingetreten: 19 Codes im Quelltext (`internal/hexagon/model/fehler.go`), höchstens vier je Klasse. Die Folgepflicht der ADR (Gate für Code-Tabelle und Katalog) ist kein Trigger, hat aber keine Slice-Kennung (Nebenbefund 4).
  - [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) — bestätigt: Kein beobachteter Treiber sendet Lebendprüfungen anders, weder eine zeitunabhängige Aufzeichnung noch ein toleranter Modus ist als Anforderung aufgenommen.
  - Hard Rules: §3.9 bis §3.12 wieder aufgetreten, alle bleiben; keine Regel mit eingetretenem Auflösungs-Trigger.
- **Drei Paarungen:** Anker — `AGENTS.md` §3.9 trägt `seit welle-walking-skeleton` in der Überschrift, §3.12 `seit slice-harness-randformen-vor-code`; die Anker der Lerneinträge (`seit slice-replay-semantik-mismatch`, `-fehlerreplay`, `-meldungscodes`) stehen in `.claude/commands/implement-slice.md` Schritt 19 bzw. `.claude/commands/close-welle.md` Schritt 1. Folge-Slice — alle sechs genannten liegen als Datei in `open/`. Register — jede genannte `BEO-REPO/…`-Kennung besteht als Verzeichnis, und jedes der 18 Verzeichnisse trägt ein nicht leeres `evidence/`.
- **Geltungsbereich der Sensoren vor der Archivierung:** Das Doku-Gate prüft Links in jeder Datei unter `docs/`, also auch in Stubs eine Ebene tiefer; Commit-Träger und `make hook-gegenprobe` lesen den Archiv-Ordner ([ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md)). `make kopf-check` prüft nur `open/`, `next/` und `in-progress/` und ist von der Archivierung nicht berührt. `make doc-trace` (Werkzeug, kein Gate) zählt nur Slices flach in `done/` und führt die fünf eingesammelten danach nicht mehr, wie in `harness/README.md` benannt.
- **Archivierung:** `make archive-welle WELLE=welle-replay-semantik` nach dem Move des Welle-Plans. Die Vorschau sammelt 3 Mitglieder, 2 wellenlose (`slice-harness-randformen-vor-code`, `slice-harness-kopf-sensor`, beide seit der Closure von welle-extended-query in `done/`), 0 fremde Slices und 0 Review-Reports ein. Gewollt: Die Closure von welle-extended-query hat beide ohne Welle angelegt mit „eingesammelt werden sie von der nächsten Welle-Closure“, und `.claude/commands/close-welle.md` Schritt 4 sammelt die wellenlosen seit der letzten Closure ein. Die übrigen wellenlosen Slices liegen in `open/` und bleiben liegen.
