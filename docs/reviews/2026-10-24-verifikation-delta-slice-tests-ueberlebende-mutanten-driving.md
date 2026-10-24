# Verifikation (Delta) — slice-tests-ueberlebende-mutanten-driving: Frist 30 s an `<-danach`

Stand: main bdb53346, Zweig `ver-drv2`. Rolle: Verifier. Gegenstand: die eine offene Mutation
(Frist als Literal an `<-danach` in `internal/bootstrap/play_anmeldung_test.go`, Commit 9d4ddda9;
Entscheidung SPEC-038, Commit 0950edcc). Gefahren wurde nur dies; kein `make gates`, keine Gegenproben.

## Aufbau

- Drei frische Kopien per `git archive bdb53346` unter dem Scratchpad (`orig`, `mut`, `nofrist`).
- Image `golang:1.27.3-alpine@sha256:b8c9b82a...5e9c`, Bind-Mount, Tests
  `^(TestRunPlayPasswort|TestRunPlayErstesSignalInAnmeldung)$`, `-count=1` (Lauf 2: `-count=20`).
- Mutation: je eine zeilengenaue Löschung der beiden `typen <- string(gelesen)` in der Lesestufe des
  Fake-Servers (`sed`, Trefferzahl vorher 2, danach 0; `typen <- ""` in der Anmeldestufe bleibt).
- Gegenprobe-Kopie `nofrist`: derselbe Mutant, dazu die zwei Zeilen
  `case <-time.After(30 * time.Second): t.Fatal("der Fake-Server sieht binnen 30 s kein Ende ...")`
  an `<-danach` in `TestRunPlayPasswort` (je Fall) und `TestRunPlayErstesSignalInAnmeldung` entfernt
  (diff gegen `orig`: genau 2 + 2 + 2 Zeilen).

## Mutationstabelle

| # | Kopie | Lauf | Erwartung | Ergebnis |
|---|---|---|---|---|
| 1 | mut (Fake-Server füllt `danach` nie) | `-count=1 -timeout 250s` | rot nach der Frist mit der Meldung | rot: `Variable_bei_host:port` 30,03 s, `Platzhalter_vor_der_Variable` 30,03 s, `TestRunPlayErstesSignalInAnmeldung` 30,51 s; jeweils „der Fake-Server sieht binnen 30 s kein Ende der Verbindung“; Paket 90,6 s, weit unter 250 s; `falsche_Variable`, `falscher_Platzhalter` grün (kein `danach` gelesen) |
| 2 | orig (ohne Mutant) | `-count=20` | grün | `ok ... 10.552s` |
| 3 | nofrist (Mutant, Frist entfernt) | `-count=1 -timeout 90s` | hängt bis zum Paket-Timeout | `panic: test timed out after 1m30s`, laufend `TestRunPlayPasswort/Variable_bei_host:port (1m30s)` |

## Ausgaben (Auszug)

Lauf 1:
```
--- FAIL: TestRunPlayPasswort (60.08s)
    --- FAIL: TestRunPlayPasswort/Variable_bei_host:port (30.03s)
    --- FAIL: TestRunPlayPasswort/Platzhalter_vor_der_Variable (30.03s)
    --- PASS: TestRunPlayPasswort/falsche_Variable (0.00s)
    --- PASS: TestRunPlayPasswort/falscher_Platzhalter (0.00s)
--- FAIL: TestRunPlayErstesSignalInAnmeldung (30.51s)
FAIL	github.com/pt9912/pgwire-recorder/internal/bootstrap	90.599s
```
Lauf 2: `ok  github.com/pt9912/pgwire-recorder/internal/bootstrap  10.552s` (20 Läufe).
Lauf 3: `panic: test timed out after 1m30s` / `running tests: TestRunPlayPasswort (1m30s)`.

## Urteil

Behauptung des Implementers **bestätigt**:
- Der Mutant wird rot nach der Frist (30,03 s / 30,03 s / 30,51 s; Implementer: 30,03 s bzw. 30,54 s),
  mit der benannten Meldung, nicht über den Paket-Timeout.
- Ohne Mutant 20 Läufe grün (Implementer meldete 50; hier nur 20 verlangt und gefahren).
- Die Frist trägt die Wirkung: ohne sie hängt derselben Mutant bis zum Paket-Timeout (Gegenprobe 3).
  Die Tests wären also ohne die Frist gegen diesen Mutanten nicht rot, sondern ein Paket-Timeout.

Grenzen: Gefahren ist nur die beschriebene Mutation (beide Stellen zugleich); eine Mutation je einzelner
Stelle und `TestRunPlayOhnePasswort`/`...Schweigen` (ebenfalls mit Frist) sind nicht Gegenstand dieses Laufs.
Im Lauf 3 brach der Timeout beim ersten hängenden Test ab; die zwei übrigen Tests wurden dort nicht erreicht,
ihre Frist zeigt Lauf 1. Module wurden zur Laufzeit geladen (Netz), nicht Teil der Aussage.
