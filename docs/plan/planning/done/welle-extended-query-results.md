# Welle welle-extended-query — Extended Query Protocol — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-extended-query
**Abschluss:** 2026-10-05
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- Das Domain-Modell und das Recording-Format tragen Extended-Interaktionen (`type: extended`, Gruppen aus Client- und Server-Nachrichten nach [ADR-0012](../../adr/0012-extended-query-gruppen.md)); die Randformen des Formats sind in der Spezifikation entschieden (`slice-extended-query-modell`).
- `record` vermittelt `Parse`, `Bind`, `Describe`, `Execute`, `Close`, `Flush` und `Sync` je Session in zwei unabhängigen Richtungen ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)) und zeichnet sie geordnet auf ([`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten); `slice-extended-query-record`).
- `replay` beantwortet Extended-Interaktionen strict sequential aus der Aufzeichnung, mit `ErrorResponse` samt Verwerfen bis zum `Sync` und `PGR-E5001` bei Abweichung ([`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers); `slice-extended-query-replay`).
- Lebendprüfungen von `pgxpool` und `database/sql` beantwortet das Replay außerhalb der Reihe und überspringt aufgezeichnete ([ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)); die Zusage „nur Host und Port“ hält damit für nacheinander genutzte Verbindungen eines Pools ([`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung); `slice-extended-query-lebendpruefung`).
- Abnahmeszenario 7 ist erfüllt: pgx im Standardmodus, aufgezeichnet gegen PostgreSQL, wiedergegeben bei gestoppter Instanz, mit derselben Sicht, darunter die Fehlerantwort `22012` und das Verwerfen im Batch (Fehlerfall von Abnahmeszenario 6 für Extended). Meilenstein M2 ist damit erreicht.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Rollen-Trennung und Validierung: Die Validierung urteilte zweimal rot bei grüner Verifikation und fand damit den Bedarf hinter dem Plan — das Herunterfahren ohne Obergrenze (`slice-extended-query-record`, an `slice-v1-abschluss-betrieb`) und die zeitabhängige Lebendprüfung von Pools (`slice-extended-query-replay`). Die zweite wurde in derselben Welle behoben; die Nachvalidierung zu `7bc1605` bestätigt Frage 2 in zwölf von zwölf Läufen der Lagen S5 bis S7′.
- Mutationstabellen als Beleg der Implementer-Rolle (`slice-extended-query-modell`, `slice-extended-query-record`) und Mutationen der Verifikation (24 von 24 rot in `slice-extended-query-lebendpruefung`); bewusstes Brechen der DoD-Tests zeigte, welche Lagen den Fix belegen und welche nur Regressionsschutz sind (V-35).
- Randformen vor dem Code: In `slice-extended-query-lebendpruefung` nannte §6 die Randformen, und [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) entschied sie vor der ersten Code-Zeile; keine davon öffnete ein Review wieder. Das ist der Gegenbeleg, auf dem `slice-harness-randformen-vor-code` aufbaut.
- Zwei Integrationsphasen (mit und ohne PostgreSQL) tragen Abnahmeszenario 7 für Extended so, wie sie im Walking Skeleton Abnahmeszenario 2 trugen.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Vier Folge-Reviews beim Replay.** `slice-extended-query-replay` brauchte ein Review, vier Folge-Reviews, Verifikation mit Nachverifikation und Validierung. Fast alles nach der ersten Runde betraf die Diagnose, nicht das Matching: was „erwartete Query“ bei `bind`, `execute`, `describe` und `close` heißt, sagte die Spezifikation nicht, und jede Runde fand eine weitere Randform. Konsequenz: `BEO-REPO/spec-randform-erst-im-review-entschieden` geplant in `slice-harness-randformen-vor-code`.
- **Ein HIGH zur Serverversion.** Das Review von `slice-extended-query-lebendpruefung` fand F-350: `\v` ist erst ab PostgreSQL 17 Leerraum, die Spezifikation zählte ihn ohne Version, der Kommentar sagte die Zeichen des Scanners zu, und das Gate prüft nur gegen 17. Der Nutzer band `\v` an `server_version` der Session. Konsequenz: benannte Spec-Lücke mit Adresse `slice-v1-abschluss-postgres-versionen` ([`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)), der vor dem Trigger von welle-v1-abschluss startet.
- **Ein HIGH zum Gegendruck.** Das Review von `slice-extended-query-record` fand F-301 trotz Mutationstabelle: Server und Recorder blockierten sich bei großer Gruppe gegen große Ausgabe. Gegendruck war keine Zusage des Plans. Konsequenz: [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) vor dem nächsten Code-Commit.
- **Der Plan folgte den Korrekturen in jedem Slice nicht**, obwohl `AGENTS.md` §3.9 es seit welle-walking-skeleton verlangt: Kopf ohne berührte Spec-Stelle, Folge-Slice nur teilweise nachgezogen, §3 ohne alle geänderten Dateien. Konsequenz: Sensor geplant in `slice-harness-kopf-sensor`.
- **Die Welle wuchs um einen Slice** (`slice-extended-query-lebendpruefung`, Drift-Log der Roadmap 2026-10-05), weil die Validierung den Bedarf aus [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) für Pools rot sah.
- **Offen, nicht in dieser Closure entschieden:** (1) Die Frage aus `slice-extended-query-record`: Der Gleichzeitigkeitsvertrag am Upstream-Port hängt an der PGWire-Bibliothek, der Re-Evaluierungs-Trigger von [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) nennt genau sie, und kein Gate fährt den Race-Detector — ein Trigger ohne Wächter. (2) Die Nachvalidierung zu `7bc1605` meldet, dass das Handbuch `--fail-on-unconsumed` ohne Hinweis beschreibt, obwohl erst `slice-replay-semantik-mismatch` die Option liefert; ob das Handbuch bis zur Lieferung den Zielstand beschreiben darf, ist eine Frage an den Nutzer. (3) Lage S8 der Nachvalidierung (gleichzeitige Verbindungen eines Pools, 4 von 5 rot) ist Information für `slice-v1-abschluss-sessions`, kein Befund dieser Welle.

## Lerneinträge der Slices

Die Steering-Loop-Einträge aus §7 der vier Slices, als Zeiger (das Original steht im Archiv dieser Welle):

- `slice-extended-query-modell` — benannte Spec-Lücke, geschlossen: Randformen des Recording-Formats (`SPEC-001`, `SPEC-041`) entschieden.
- `slice-extended-query-record` — neuer Sensor: Die Stufe `test` prüft die Formatierung aller Go-Dateien (`gofmt -l`), liegt in `Dockerfile §Stufe test`, Anker `· seit slice-extended-query-record`; dazu benannte Spec-Lücke, offen: [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) sagt für das kontrollierte Beenden keine Obergrenze zu → `slice-v1-abschluss-betrieb`.
- `slice-extended-query-replay` — benannte Spec-Lücke: [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) und [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) hielten unter strengem Replay für Connection-Pools nicht → geschlossen in `slice-extended-query-lebendpruefung` mit [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md); zwei weitere mit Adresse (falsche Protokollart → `slice-extended-query-lebendpruefung`, Parameter-Index in der Diagnose → `slice-replay-semantik-meldungscodes`).
- `slice-extended-query-lebendpruefung` — benannte Spec-Lücke: keine unterstützte Serverversion im Vertrag → `slice-v1-abschluss-postgres-versionen` ([`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Guide** geschärft: Jede Zusage eines neuen Vertrags hat einen Test, der ihre Mutation fängt, und der Implementer fährt die Mutation vor der Übergabe selbst
  — liegt in `AGENTS.md §3.10`.
  Auslöser: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (slice-walking-skeleton-record, slice-walking-skeleton-replay, slice-extended-query-modell, slice-extended-query-record, slice-extended-query-replay, slice-extended-query-lebendpruefung — 6×). Zweiter Zielort: `.claude/commands/implement-slice.md` Schritt 19.
- **Guide** geschärft: Ein Kommentar, Hilfetext, eine Abdeckungs-Deklaration oder Plan-Zeile sagt nur zu, was ein Test oder Gate prüft
  — liegt in `AGENTS.md §3.11`.
  Auslöser: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (slice-walking-skeleton-build-gates, slice-walking-skeleton-record, slice-extended-query-modell, slice-extended-query-record, slice-extended-query-replay, slice-extended-query-lebendpruefung — 6×). Zweiter Zielort: `.claude/commands/implement-slice.md` Schritt 20.
- **Guide** geplant: Randformen eines neuen Vertrags vor dem Code entscheiden (Planner, Architect, Implementer) — Kennung `slice-harness-randformen-vor-code`.
  Auslöser: `BEO-REPO/spec-randform-erst-im-review-entschieden` (slice-extended-query-modell, slice-extended-query-record, slice-extended-query-replay, slice-extended-query-lebendpruefung — 4×).
- **Sensor** geplant: Ein Gate prüft, ob jede `LH-`, `SPEC-` und `ARC-`Kennung aus §1 und §2 eines Slice im Kopf unter `Bezug` bzw. `Berührte Spec-Stellen` steht — Kennung `slice-harness-kopf-sensor`. Die Prosa-Regel `AGENTS.md` §3.9 bleibt verkörpert; nach dem vierten Auftreten über der Schwelle gilt sie als ausgeschöpft (v6.13.0 · `regelwerk/modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3), der Sensor deckt die maschinell entscheidbare Hälfte.
  Auslöser: `BEO-REPO/plan-folgt-korrektur-nicht` (slice-walking-skeleton-build-gates, slice-walking-skeleton-record, slice-walking-skeleton-replay, slice-extended-query-record, slice-extended-query-replay, slice-extended-query-lebendpruefung — 6×). Retirement-Check von §3.9: wieder aufgetreten, die Regel bleibt.

Weitere Einträge über der Schwelle: keine. Unter der Schwelle, nicht gelesen, aber als Sensor-Kandidat benannt: `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` (1×, V-30) — ein Test der zweiten Integrationsphase lässt sich mit `make test-integration` nicht einzeln brechen, weil der Runner nach roter erster Phase endet.

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

- `slice-harness-kopf-sensor` (Sensor für den Kopf des Slice-Plans) — neu, ohne Welle.
- `slice-harness-randformen-vor-code` (Randformen vor dem Code entscheiden) — neu, ohne Welle.
- `slice-v1-abschluss-postgres-versionen` (unterstützte PostgreSQL-Versionen, Versionsmatrix) — welle-v1-abschluss, startet vorab.
- `slice-v1-abschluss-sessions` (Lebendprüfungen bei `--session-assignment connection`, Regel vor der Zuordnung) — welle-v1-abschluss.
- `slice-v1-abschluss-betrieb` (Frist für das Herunterfahren in Record und Replay) — welle-v1-abschluss.
- `slice-v1-abschluss-cancel-ohne-schluessel` (F-310) — welle-v1-abschluss.
- `slice-replay-semantik-mismatch` (`--fail-on-unconsumed`) — welle-replay-semantik.
- `slice-replay-semantik-meldungscodes` (Parameter-Index in der Diagnose, `--log-level`) — welle-replay-semantik.

Warum die beiden neuen ohne Welle: Ihre Closure-Bedingung ist ihre DoD; kein Abnahmeszenario und kein Meilenstein hängt an ihnen, und welle-replay-semantik hat ein eigenes *Mehr* (Abnahmeszenarien 4 und 6 für Simple Query), dem Harness-Arbeit nichts hinzufügt. Eingesammelt werden sie von der nächsten Welle-Closure. Sie wirken erst auf Slices, die nach ihnen geplant und gearbeitet werden; vor dem ersten Slice von welle-replay-semantik geliefert, treffen sie deren Pläne.

Bestand in `open/` (Lese-Schritt, Schritt 3): Alle 18 übrigen Slices tragen eine `Welle:` einer offenen Welle (welle-replay-semantik, welle-v1-abschluss, welle-erster-release); keiner ist wellenlos liegen geblieben, keine Gruppe mit gemeinsamem Auslöser zum Konsolidieren. Alle **bestätigt**: Ihre Priorität trägt die Reihenfolge der Wellen.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- **Slices:** alle vier (`slice-extended-query-modell`, `slice-extended-query-record`, `slice-extended-query-replay`, `slice-extended-query-lebendpruefung`) lagen bei Closure-Beginn (`17a7615`) flach in `done/`; `next/` und `in-progress/` ohne Slice.
- **`make gates`** Exit 0 auf `17a7615`: Build, Unit-Tests mit `gofmt`, Integration in zwei Phasen (`run-integration-tests: gruen`), Architektur-Gate mit Gegenprobe, Abdeckung mit Gegenprobe, Commit-Träger-Gegenprobe, Doku-Gate (186 Dateien, 0 Befunde), Baseline (`v6.13.0`, 54 Dateien).
- **Das *Mehr* aus §3 der Welle**, aus demselben Lauf: Abnahmeszenario 7 — `TestE2EVorbereitungExtendedOhnePostgres` PASS (pgx im Standardmodus über `record` gegen PostgreSQL), danach bei gestoppter Instanz `TestE2EOhnePostgresExtendedReplay` PASS (dieselbe Sicht, Exit-Code 0). Fehlerfall von Abnahmeszenario 6 für Extended: Die Sicht beider Tests enthält `fehler: 22012 division by zero` und das Verwerfen im Batch; `TestE2EReplayExtendedPgx` PASS (dieselbe Sicht in drei Läufen, Abdeckung `LH-FA-11/Happy`). Abweichung: `TestE2EReplayExtendedAbweichung` PASS (`PGR-E5001`, Exit-Code 5). Walking-Skeleton-Smoke: `TestE2EVorbereitungOhnePostgres` PASS, `TestE2EOhnePostgresReplay` PASS.
- **Validierung:** Validierung zu `slice-extended-query-replay` (2026-10-05) samt Abschnitt „Nachvalidierung zu `7bc1605`“: Frage 2 trägt; S8 ist Information für `slice-v1-abschluss-sessions`.
- **Trigger-Audit:** Carveouts 0 offen (`docs/plan/carveouts/` trägt nur `.gitkeep`). Bootstrap-aware Gates: keine Reifestufe im Repo, 0 offen. ADRs: [ADR-0007](../../adr/0007-strict-replay.md) — Trigger für eine Klasse von Anfragen eingetreten, bedient durch [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) (ergänzt, ersetzt nicht), für alle anderen bestätigt; [ADR-0012](../../adr/0012-extended-query-gruppen.md) — bestätigt, pgx im Standardmodus erzeugt dieselben Namen in jedem Lauf (`TestE2EReplayExtendedPgx`, drei Läufe), und kein toleranteres Matching ist als Anforderung aufgenommen; [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) — bestätigt, `replay` bleibt synchron (Validierung, Frage 7), die Bibliothek sagt gleichzeitiges Lesen und Schreiben weiter zu, ohne Wächter (oben, offen); [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) — bestätigt, kein beobachteter Treiber sendet Lebendprüfungen anders. Hard Rules: §3.9 wieder aufgetreten, bleibt; keine Regel mit eingetretenem Auflösungs-Trigger.
- **Drei Paarungen:** Anker — `AGENTS.md` §3.10 und §3.11 tragen `seit welle-extended-query` in der Überschrift, `.claude/commands/implement-slice.md` Schritt 19 und 20 ebenso; §3.9 trägt `seit welle-walking-skeleton`. Folge-Slice — alle acht genannten liegen als Datei in `open/`. Register — jede genannte `BEO-REPO/…`-Kennung besteht als Verzeichnis, und jedes der zwölf Verzeichnisse trägt ein nicht leeres `evidence/`.
- **Geltungsbereich der Sensoren vor der Archivierung:** Das Doku-Gate klassifiziert Slices und Wellen über `docs/plan/planning/**/…` und prüft Links in jeder Datei, also auch in den Stubs; Commit-Träger und `make hook-gegenprobe` lesen den Archiv-Ordner seit [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md). Nicht mit wandert `make doc-trace` (Werkzeug, kein Gate): Es zählt nur Slices flach in `done/`; nach der Archivierung führt die Matrix die vier Slices nicht mehr. Benannt in `harness/README.md` (Werkzeug-Tabelle).
- **Archivierung:** `make archive-welle WELLE=welle-extended-query` nach dem Move des Welle-Plans; Vorschau: 4 Mitglieder, 0 wellenlos, 0 fremd, 0 Review-Reports.
