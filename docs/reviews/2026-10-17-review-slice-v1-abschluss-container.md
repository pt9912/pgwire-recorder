# Review-Report: slice-v1-abschluss-container — 2026-10-17

**Review-Art:** Code — geprüft gegen Plan, `SPEC-055`, [`ADR-0043`](../plan/adr/0043-binaries-und-container-gate.md), [`ADR-0026`](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md), `SPEC-052` (Punkte 2, 4, 13), `SPEC-035` und die Hard Rules.

**Gegenstand:** `slice-v1-abschluss-container`, Diff `d258c0b7..1bac5050` über `Dockerfile harness tools test .gitignore`.

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form.** Dieser Report friert ein: Kennung, nicht Adresse; Baseline-Stellen als Tag und Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-v1-abschluss-container` (§1, §2, §3, §4, §6, §7, §8)
- [`ADR-0043`](../plan/adr/0043-binaries-und-container-gate.md) (Accepted 2026-10-17), [`ADR-0026`](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md)
- Verträge `SPEC-055` (13 Punkte, Stand nach 1d0d2740), `SPEC-052` (Punkte 2, 4, 13), `SPEC-035`
- `LH-QA-01`, `LH-QA-03`, `LH-FA-16`
- `AGENTS.md` (Hard Rules, §3.6, §3.7, §3.9 bis §3.12, §3.14)

**Eigene Läufe (nacheinander, Docker-Läufe nie parallel):**

- `make binaries` (leeres `dist/`): sechs Binaries, `SHA256SUMS` in Byte-Reihenfolge, `sha256sum -c` grün; der Hash des Arbeitsbaums (`tools/harness/working-tree-hash.sh`) davor und danach gleich, `git status` leer; ein zweiter Aufruf endet mit `DIST 'dist' ist nicht leer` (Exit 2), kein `docker`.
- `make test-container`: grün, 1 m 04 s; danach keine Container, kein Netz, keine Tags `pgr-ct-*`, kein `pgr-container.*` in `/tmp`.
- `make test-container-gegenprobe` über `make`: grün, 3 m 46 s. `make gates-schnell-gegenprobe`: Exit 0.
- Mutationen der Tabelle in §7 (Stichprobe, 26 Namen aus allen Gruppen): Punkt 1 (`p1d`, `p1e-bin-darwin/arm64`, `p8-vor-docker`), 2 (`p2-dollar`, `p2-minus`), 3 (`p3-sort`, `p3-dist-leer`, `p3-teilmenge`, `p3-exe`), 4 (`p4-filter-dist`, `p4-abweichung-windows/arm64`), 5 (`p5-macho-cpu`, `p5-pe-offset-3c`), 6 (`p6-summe-linux/arm64`, `p6-runtime-run`), 7 (`p7-intern`, `p7-user`, `p7-mount`, `p7c-pg`), 8 (`p8-4-weiter`), 9 (`p9-container`, `p9-weiter`, `p9-stage-exit`), 11 (`p11-plus-binaries`, `p11-i-gate`), 12 (`p12-gegenprobe`): in jedem Fall färben genau die in §7 benannten Fälle rot.
- Punkt 13 von `SPEC-052`: Die Gegenprobe lief in einer Kopie des Baums ohne `docs/user/`, `docs/plan/planning/`, `docs/reviews/` und `README.md` (`git archive` ohne diese Pfade); die unveränderte Kopie bestand jeden Fall. Die Skripte und `container.mk` nennen keinen dieser Pfade.
- Handlauf Zeitstempel (in einer Kopie außerhalb des Arbeitsbaums, `-X main.version=${VERSION}-$(date +%s%N)` im `Dockerfile`): `make test-container` rot bei Punkt 4 mit den Namen aller sechs Binaries, kein Image gebaut, danach nichts zurückgeblieben; die Angaben in §7 stimmen (der "Exit 2" ist der Status von `make`, das Skript endet mit 1).
- Eigene Mutationen außerhalb der Tabelle (Kopie ohne Docs, Zusatzzeilen in der Kopie der Gegenprobe): siehe F-1 bis F-3 und F-8.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Entfernt man `-X main.version=$VERSION` aus der Stufe `build` des `Dockerfile`, bleibt `make test-container` mit der Vorgabe `VERSION=dev` grün (gefahren, 1 m 56 s, Kopie außerhalb des Arbeitsbaums): `main.version` ist in `cmd/` `"dev"`, und der Smoke `version` vergleicht gegen `VERSION`. Die Zusage von `SPEC-055` Punkt 2, `VERSION` gehe in `-X main.version` des Builds ein, hat für das echte `Dockerfile` weder einen Fall noch eine Mutation; die Attrappe belegt nur, dass `--build-arg VERSION=` übergeben wird (`p3-version-im-build`). | `AGENTS.md` §3.10; `SPEC-055` Punkt 2 | `Dockerfile` · `-ldflags="-s -w -X main.version=$VERSION"` | ja — `make test-container` mit der Mutation in einer Kopie bleibt grün | Zusage des Vertrags am echten Artefakt ohne roten Fall (Vorgabewert maskiert) |
| F-2 | LOW | Die Fälle `p7-kein-volume` und `p7-kein-user` erkennen nur ` -v `, ` --volume`, ` --user`, ` -u ` als eigenes Wort; `run -dv /tmp:/tmp` und `run -d -u0` an der `record`-Zeile bleiben grün (beide gefahren), während `p7-kein-terminal` zusammengesetzte Kurzflags erkennt. Die Zusage "kein `--user`, `-v`, `--mount`" ist damit je Form nur für eine Schreibweise gefangen. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/harness/container-gegenprobe.sh` · `ist p7-kein-volume '! grep -E "^(run\|create) " "$log" \| grep -qE " -v \| --volume"'` | ja — `GP_MUT` mit zwei Zusatzmutationen | Zusage-Prüfung deckt nur eine Schreibweise der Form |
| F-3 | LOW | Die Köpfe von `run-container.sh` und `binaries.sh` sagen "Ausgang: 2 bei einem Fehler der Eingabe, 1 bei einem roten Punkt"; kein Fall prüft den Wert (`exit 3` an beiden Stellen bleibt grün, gefahren), die Fälle prüfen nur "ungleich 0". `SPEC-055` Punkt 8 verlangt nur "ungleich 0". | `AGENTS.md` §3.11 | `tools/harness/run-container.sh` · `# Ausgang: 2 bei einem Fehler der Eingabe (vor dem ersten Aufruf von docker), 1 bei einem roten` | ja — Mutation `exit 3` bleibt grün | Kommentar sagt weiter zu als geprüft |
| F-4 | LOW | Der Kopf der Gegenprobe sagt über einen Zeitstempel im Binary: "Er färbt Punkt 4 rot, und der Fall der abweichenden Summe belegt das". Der Fall der Attrappe belegt nur den Vergleich; die Aussage über das echte `Dockerfile` stützt allein der einmalige Handlauf in §7 des Plans, den kein Gate wiederholt (der Plan selbst nennt das so). | `AGENTS.md` §3.11 | `tools/harness/container-gegenprobe.sh` · `Ein Zeitstempel im Binary (Punkt 4) ist kein Fall: Er färbt Punkt 4 rot, und der Fall` | nein — der Handlauf ist einmalig; ich habe ihn wiederholt und bestätigt | Kommentar sagt weiter zu als geprüft |
| F-5 | LOW | `harness/README.md` §Sensors sagt zur Gegenprobe, sie prüfe, "dass das Gate nichts in den Arbeitsbaum schreibt, Punkt 10". Der Fall `p10-arbeitsbaum` vergleicht die Liste der Dateinamen vor und nach einem Lauf in der Kopie; das Überschreiben einer vorhandenen Datei (Inhalt) färbt ihn nicht. "Liest nur, was der Build-Kontext zulässt" hat keinen Fall (die Kopie enthält die Pfade nicht, das Lesen scheitert allenfalls zufällig). | `AGENTS.md` §3.11 | `harness/README.md` · `dazu, dass das Gate nichts in den Arbeitsbaum schreibt, Punkt 10` | ja — Mutation, die eine vorhandene Datei überschreibt | Gate-Beschreibung sagt weiter zu als geprüft |
| F-6 | LOW | Die Aufschlüsselung der Messung in §7 des Plans ist veraltet: genannt sind `container-gegenprobe.sh` 764 und `binaries.sh` 68, gemessen `git diff --numstat d258c0b7..HEAD` 767 und 69; die genannten Posten ergeben 1147, die genannte Summe 1157 (`SPEC-055` mit 6 steht nur im Nebensatz). Die Summe 1157 selbst stimmt (1147 Einfügungen, 10 Löschungen, Marke 1200). Das 764 : 361 gilt damit nicht mehr. | `AGENTS.md` §3.9, §3.14 | `docs/plan/planning/in-progress/slice-v1-abschluss-container.md` · `Aufschlüsselung: `container-gegenprobe.sh` 764` | ja — `git diff --numstat` | Plan folgt der Korrektur nicht |
| F-7 | LOW | `binaries.sh` baut in ein Stage-Verzeichnis (`mktemp -d`, außerhalb von `DIST`), verschiebt erst nach dem letzten Build und räumt es bei `EXIT`, `SIGINT` und `SIGTERM`; die Gegenprobe hat dafür Fälle und Mutationen (`p9-*-stage`, `p9-binINT-*`, `res_stage`). `SPEC-055` Punkt 3 und 9 und §6 des Plans nennen dieses Verzeichnis, seinen Ort und das Aufräumen von `make binaries` allein nicht (Punkt 9 nennt die Verzeichnisse der Punkte 4 und 7). | `AGENTS.md` §3.12 | `tools/harness/binaries.sh` · `stage="$(mktemp -d)"` | nein — Vertragslücke, kein Lauf | Randform vom Code entschieden, im Vertrag nicht genannt |
| F-8 | INFO | Äquivalente Mutation, in §7 nicht geführt: `LC_ALL=C sort` zu `sort` in `binaries.sh` bleibt grün (gefahren). Invariante: Die sechs Namen unterscheiden sich nur in Betriebssystem und Architektur (`darwin_amd64` bis `windows_arm64.exe`) bei gleichem Präfix; die Reihenfolge ist in jeder Locale gleich. Die Festlegung "Byte-Reihenfolge (C-Locale)" in `SPEC-055` Punkt 3 hat damit keine beobachtbare Wirkung; `P3sort` entfernt das ganze `sort`. Zuständig: Implementer, für die Liste der Mutanten in §7. | `AGENTS.md` §3.10 | `tools/harness/binaries.sh` · `LC_ALL=C sort) # P3sort` | ja — Mutation | grüner Mutant nicht geführt |
| F-9 | INFO | Bestätigt, ohne Befund: der äquivalente Mutant "`DIST` als reguläre Datei" trägt Invariante (`mkdir -p` auf eine Datei endet mit Exit ungleich 0, vor `docker`) und Ort (§7); die Maskierung bei Punkt 5 (Lauf endet beim ersten roten Kopf, Mutationen mit `nur` begrenzt) ist als Eigenschaft in §7 und im Kopf der Gruppe genannt und von jeder Mutation der Gruppe (zwölf Architektur-/Format-Mutanten) eingehalten. Zählt nicht zu `gruener-mutant-faelschlich-aequivalent`. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-container.md` · `Äquivalente und Grenzen` | ja | — |
| F-10 | INFO | `bereit` in `run-container.sh` liest `docker logs "$1" 2>&1 \| grep -q 'gestartet'` unter `set -o pipefail`; endet `grep -q` vor dem letzten Schreiben von `docker logs`, kann die Pipeline mit 141 als nicht gefunden gelten. Bei der Länge der Ausgabe (eine Meldung) trat es in den Läufen nie auf; kein Fall. Zuständig: Implementer, falls das Gate flackert. | Maintainability | `tools/harness/run-container.sh` · `grep -q 'gestartet'` | nein | Pipeline mit frühem Ende unter pipefail |
| F-11 | INFO | Wartezeiten (je 60 s auf PostgreSQL und auf die Startmeldung) sind Wahl des Codes; die Gegenprobe prüft, dass die Schleife endet (`p7b-pg-nicht-bereit`, `p7b-kein-start` mit `sleep`-Attrappe), nicht die Dauer; Plan §7 und Skript sagen nichts anderes zu. Akzeptierte Negative (Smoke nur auf der Plattform des Servers, `SIGKILL`, `-j`) stehen in `SPEC-055` Punkt 13 und im Plan; kein Kommentar sagt mehr. `p7b-anfrage-rot`, `p7b-kein-start` und `p6-tag` haben keine Mutation in der Tabelle; die Mutationen dazu sind rot (gefahren), die Tabelle führt sie nicht. | `AGENTS.md` §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-container.md` · `Wartezeit auf PostgreSQL und auf die Startmeldung` | ja | — |
| F-12 | INFO | Größe: Die Messung 1157 ist nachgerechnet (Marke 1200, Metrik §3.14); es bleiben 43 Zeilen für Änderungen aus diesem Review und die Closure. Die Gegenprobe (767 Zeilen, rund 190 Fälle, 220 Mutationen) ist als eine Datei mit Gruppen je Punkt und Kopfkommentar zu lesen; Fälle teilen sich je Kopie einen Baum, setzen aber Zustand (`$st`, `$tmp`, `$log`) je `lauf` zurück und legen ihre Verzeichnisse unter eigenen Namen an, kein gemeinsamer Zustand färbt mehrere. Leerzeilen-Blöcke und die Zeile `#MUTATIONEN` ohne Funktion sind ohne Konventions-Anker. | Maintainability | `tools/harness/container-gegenprobe.sh` · `#MUTATIONEN` | ja | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `Dockerfile` (Stufe `dist`, Stufe `runtime` ohne `RUN`) | geprüft; Befund nur F-1 (Wiring von `VERSION`) |
| `harness/mk/container.mk` (Rezepte ohne `$(MAKE)`/`+`, `GATE_CHECKS`, Plattformlisten) | geprüft, ohne Befund |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh`, `harness/vertraege/SPEC-052-gate-kette.md` | geprüft, ohne Befund (`make gates-schnell-gegenprobe` grün, Ausnahmeliste und Deckung) |
| `harness/vertraege/SPEC-055-container-und-binaries.md` (Änderungen 1d0d2740) | geprüft, ohne Befund; Lücke nur F-7 |
| `tools/harness/container-eingabe.sh` | geprüft, ohne Befund |
| `tools/harness/run-container.sh` (Reihenfolge der Gründe, Aufräumen, Signale, Smoke, kein `--user`/`-v`/`--mount`) | geprüft; F-3, F-10 |
| `tools/harness/binaries.sh` | geprüft; F-3, F-7, F-8 |
| `tools/harness/container-gegenprobe.sh` | geprüft; F-2, F-4, F-12 |
| `.gitignore` (`dist/`), `harness/README.md` §Sensors | geprüft; `dist/` hält den Hash und `git status` stabil; F-5 |
| `test/` | geprüft, ohne Befund (Diff leer) |
| Core-Reinheit, Suppression (`nolint`), Spec-Straten, ADR-Immutabilität, Move-Commits | geprüft, ohne Befund (kein Produktcode, keine ADR geändert, keine Moves im Diff) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 6 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Zusage des Vertrags am echten Artefakt ohne roten Fall · Zusage-Prüfung deckt nur eine Schreibweise der Form · Kommentar sagt weiter zu als geprüft · Plan folgt der Korrektur nicht · Randform vom Code entschieden, im Vertrag nicht genannt · grüner Mutant nicht geführt · Pipeline mit frühem Ende unter pipefail

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM); F-2 bis F-7 sind LOW und blockieren nicht.

**Übergabe:** Die Findings gehen an den Implementer (F-7 mit Rückgabe an den Architect, weil der Vertrag die Randform nicht nennt); die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report ersetzt keine Verifikation: DoD und Spec-Konformität prüft der Verifier.
