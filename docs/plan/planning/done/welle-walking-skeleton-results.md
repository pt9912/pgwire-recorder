# Welle welle-walking-skeleton — Walking Skeleton — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-walking-skeleton
**Abschluss:** 2026-10-04
**Verantwortlich:** pt9912

## Was wurde geliefert?


Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- Docker-only Build und Test über ein Multistage-Dockerfile, Architektur-Gate mit Gegenprobe, Integrationstests gegen ein gepinntes PostgreSQL (`slice-walking-skeleton-build-gates`, `slice-walking-skeleton-record`).
- `pgwire-recorder record` vermittelt einfache Anfragen zwischen Client und PostgreSQL und schreibt eine geprüfte YAML-Aufzeichnung (`LH-FA-02` durch Tests in allen drei Pfaden belegt).
- `pgwire-recorder replay` beantwortet Anfragen aus der Aufzeichnung ohne PostgreSQL mit strict sequential matching, Diagnose und Exit-Code 5 (`LH-FA-09`, `LH-FA-10` durch Tests belegt).
- Abdeckung je Anforderung und Pfad mit RTM nur für vollständig belegte Anforderungen ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)).

- Rollen-Trennung: Review, Folge-Review und Verifikation in frischem Kontext fanden in jedem Slice echte Mängel, die grüne Gates verdeckten (verworfene Session aufgezeichnet, Exit-Code 0 nach Fehlern, abgeschnittene Datei gültig, zu schmales Orakel im Replay).
- Gegenproben zu jedem Gate (Architektur, Commit-Träger, Abdeckung): Mutationen machen sie nachweisbar rot.
- Black-Box-E2E gegen das gebaute Binary, mit einer zweiten Phase bei gestopptem PostgreSQL.

- Die Slices lieferten mehr, als ihre Pläne sagten, und die Pläne folgten den Korrekturen nicht (Steering-Loop unten).
- Fünf ADRs entstanden während der Umsetzung ([ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) bis [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0024](../../adr/0024-sqlite-bibliothek.md) angenommen): Commit-Träger, Build, YAML-Bibliothek, Abdeckung.
- Die Spezifikation regelt jetzt fremde erste Nachrichten (`PGR-W3003`, `SPEC-045`).
- Die Passwort-Anmeldung beim Upstream ist nicht vermittelt — Folge-Slice `slice-v1-abschluss-anmeldung`.

## Steering-Loop-Einträge


Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Guide** geschärft: Wer im Slice Code, Gates oder Tests ändert, zieht im selben Commit §1, §3 und §6 des Slice-Plans nach
  — liegt in `AGENTS.md §3.9`.
  Auslöser: `BEO-REPO/plan-folgt-korrektur-nicht` (slice-walking-skeleton-build-gates, slice-walking-skeleton-record, slice-walking-skeleton-replay — 3×).


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

- `slice-extended-query-modell` (Extended-Interaktion im Modell und im Leser) — startet `welle-extended-query`.
- `slice-v1-abschluss-anmeldung` (Anmeldung des Clients im Record-Modus vermitteln) — `welle-v1-abschluss`.
- `slice-replay-semantik-mismatch` (Extended-Mismatch, `--fail-on-unconsumed`) — `welle-replay-semantik`.

## Verifikation


Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- `make gates` Exit 0 auf dem Stand der Closure (Commit `988b2dd`): Build, Unit-Tests, Integration in zwei Phasen, Architektur-Gate mit Gegenprobe, Abdeckung mit Gegenprobe, Commit-Träger-Gegenprobe, Doku-Gate (133 Dateien, 0 Befunde), Baseline.
- Walking-Skeleton-Smoke: `TestE2EVorbereitungOhnePostgres` zeichnet `SELECT 1;` über `record` gegen das gepinnte PostgreSQL auf; der Runner stoppt die Instanz; `TestE2EOhnePostgresReplay` erhält im Replay dieselbe Sicht (Spalten, Typ-OIDs, Zeilen, Befehlsabschluss) — ohne interaktive Eingabe.
- Trigger-Audit: Carveouts 0 offen (`docs/plan/carveouts/` leer); keine bootstrap-aware Reifestufe im Repo; Re-Evaluierungs-Trigger der ADRs: die eingetretenen von [ADR-0016](../../adr/0016-einspielen-anmeldung-und-tls.md) und [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md) sind durch [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) und [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md) bedient, sonst keiner eingetreten.
- Drei Paarungen: Anker `AGENTS.md §3.9` trägt `seit welle-walking-skeleton`, die beiden Sensor-Anker der Slices tragen `seit slice-…`; jeder Folge-Slice liegt als Datei in `open/`; jedes genannte `BEO-REPO/…` existiert mit nicht leerem `evidence/`.
- Archivierung: `make archive-welle WELLE=welle-walking-skeleton` hat die drei Slices und den Welle-Plan nach `done/welle-walking-skeleton/` verschoben und durch Stubs mit `archiv.zip` ersetzt; der Commit-Hook nimmt die Welle-Kennung seit [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) an (Register `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung`, behoben).
