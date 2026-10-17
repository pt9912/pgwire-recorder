# Verifikations-Bericht: slice-v1-abschluss-container — 2026-10-17

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gemessen an Plan, DoD und Entscheidungen?

**Gegenstand:** `slice-v1-abschluss-container` (Schritt 44a), Diff `d258c0b7..e365401f` (Code 91a1bd02, 7e5ecb79; Nacharbeit 39c07812, 62a287cb, e365401f; Review 0f958d7e).

**Gegen:** Plan §1 bis §8; [`SPEC-055`](../../harness/vertraege/SPEC-055-container-und-binaries.md) (nach 9ec49476, 1d0d2740, 1c60e745); [`ADR-0043`](../plan/adr/0043-binaries-und-container-gate.md) (Accepted); `SPEC-052` Punkte 2, 4, 13; `SPEC-035`; `LH-QA-01`, `LH-QA-03`, `LH-FA-16`; Review-Report `docs/reviews/2026-10-17-review-slice-v1-abschluss-container.md`.

> **Zitier-Form.** Dieser Bericht friert ein: Kennung, nicht Adresse.

## Ergebnis

**Konform, ohne Blocker.** Kein HIGH, kein MEDIUM. Zwei Befunde (1 LOW, 1 INFO) und ein Hinweis. Die Nacharbeit zu F-1 und F-2 des Reviews trägt: Die Brechprobe am echten `Dockerfile` und an der Gegenprobe wird rot, wie behauptet. Die DoD-Punkte der Lieferung sind erfüllt; die Closure-Punkte (Notiz, Lerneintrag, Register, Paarungen, Risiko-Ausgänge) stehen noch aus, wie vorgesehen.

## 1. DoD gegen den Diff

| DoD-Punkt | Befund |
|---|---|
| `make binaries`: sechs Binaries aus `SPEC-035` mit `SHA256SUMS` | erfüllt; selbst gefahren (29 s), `sha256sum -c` grün, Namen `pgwire-recorder_<VERSION>_<os>_<arch>`, `.exe` unter Windows |
| `make test-container`: zwei Builds, gleiche Summe, Kopf je Binary, Image je Linux-Plattform mit demselben Binary, Smoke (`replay` ohne PostgreSQL) | erfüllt; selbst gefahren, grün in 1 m 08 s: Punkte 4 bis 7 grün, Smoke auf `linux/amd64` |
| Gegenprobe fängt je Zusage eine Mutation (`AGENTS.md` §3.10) | erfüllt für die Stichprobe (Abschnitt 2); die Tabelle in §7 deckt Punkte 1 bis 9, 11, 12 |
| `make gates` grün | Nachweis trägt: Hash des Arbeitsbaums (`working-tree-hash.sh`) gleich dem Inhalt von `.harness/state/gates-passed.diffsha` (`c186527d…`), also Code-Stand gleich dem des grünen Laufs an e365401f; Wiederholung nicht nötig |
| Review, Report liegt vor | erfüllt (`0f958d7e`); die Nacharbeit wurde nicht erneut reviewt, F-1 und F-2 habe ich selbst geprüft (Abschnitt 2) |
| Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen | noch offen (Closure-Schritt); Vorschlag in Abschnitt 6 |

## 2. Bewusstes Brechen (Stichprobe, selbst gefahren)

| Zusage | Mutation | Ergebnis |
|---|---|---|
| F-1: `VERSION` geht in `-X main.version` des echten `Dockerfile` | `-X main.version=$VERSION` aus der Kopie des echten `Dockerfile` entfernt (Kopie außerhalb des Arbeitsbaums), `make test-container` | rot nach Punkt 6 im Smoke `version`: `version gibt 'pgwire-recorder dev' aus statt 'pgwire-recorder 0.0.0-test' und LF`, Exit 1; Punkte 4 bis 6 grün, also erst der Smoke fängt es, wie der Vertrag zusagt |
| F-1 an der Gegenprobe | `p2-vorgabe`, `p2-dev-gate`, `p2-dockerfile-x`, `p2-dev` | je genau `p2-vorgabe-test`, `p2-dev-gate`, `p2-dockerfile-x`, `p2-gueltig` rot |
| F-1 am Gate | `make test-container VERSION=dev` | `VERSION 'dev' ist am Gate unzulässig`, Exit 2, vor `docker` |
| F-2 Kurzflags | `p7-volume-kurz` (`-dv`), `p7-user-kurz` (`-u0`), `p7-user`, `p7-volume` | je nur `p7-kein-volume` bzw. `p7-kein-user` rot |
| Punkt 1/8 | `p8-vor-docker` | alle Fälle der Punkte 1 und 2, dazu `p2-dev-gate`, `p8-1-vor-5` rot (genau die Erwartung aus e365401f) |
| Punkt 3 | `p3-dist`, `p3-dist-leer`, `p3-sort`, `p3-vor-docker`, `p3-dist-fehlt` | je die in §7 genannten Fälle |
| Punkt 4 | `p4-filter-build`, `p4-filter-dist` (`--no-cache-filter`) | `p4-cache-umgehen`, `p4-gruen` |
| Punkt 5 | `p5-elf-magic`, `p5-macho-magic`, `p5-macho-cpu` | je nur Format- bzw. Architektur-Fälle der Plattformen |
| Punkt 7 | `p7-server`, `p7-server-amd64`, `p7-server-arm64` | `p7-server-fremd`, `p7-server-arm64`, `p7-gruen` |
| Punkt 9 | `p9-weiter` (SIGTERM-Weitergabe) | die fünf `-signal`-Fälle |
| Punkt 11 | `p11-i-gate` | `p11-gate-i`, `p11-gate-makeflags-i` |
| Punkt 12 | `p12-ziel`, `p12-gegenprobe` | `p12-gate-checks-ziel`, `p12-gate-checks-gegenprobe` |

Die unveränderte Kopie besteht jeden Fall (jeder Lauf von `container-gegenprobe.sh` oben endet grün). `make gates-schnell-gegenprobe` grün (Ausnahmeliste deckt `GATE_CHECKS`).

## 3. Läufe und Reste

- `make test-container`: grün; danach kein Container, kein Netz, kein Tag `pgr-ct-*`, kein `pgr-container.*` in `/tmp`. Dasselbe nach der Handmutation und nach einem versehentlich nach 8 s abgebrochenen Lauf (siehe Hinweis 1). `pgr-dev` ist der Container des anderen Agenten, nicht von diesem Slice.
- `make binaries`: Hash des Arbeitsbaums davor und danach gleich, `git status` leer (`.gitignore` mit `dist/`); `dist/` danach von mir entfernt.
- `make gates-schnell` Exit 0 (10 s), `make abdeckung-check`, `make meldungskatalog-check`, `make kopf-check` Exit 0. Stop-Hook: Hash des Arbeitsbaums gleich Nachweis des vollen Laufs, Hash gleich Stand der schnellen Kette.

## 4. Entscheidungs-Konformität

`SPEC-055` Punkte 1 bis 13: Der Code folgt dem Vertrag (Plattformlisten nur in `container.mk`, Eingabe vor `docker`, `VERSION` 0.0.0-test am Gate, `DIST` leer als Eingabefehler, Byte-Reihenfolge, Stage-Verzeichnis, Aufräumen, Rezepte ohne `$(MAKE)`, Anschluss an `GATE_CHECKS`). [`ADR-0043`](../plan/adr/0043-binaries-und-container-gate.md): Binaries sind reproduzierbar belegt, das Image nicht (Grenze), Manifestliste und Registries nicht im Slice, wie §1 sagt. `SPEC-052` Punkte 2, 4, 13: Ausnahmeliste enthält `test-container` und `test-container-gegenprobe`, Gegenprobe grün; die Skripte nennen keine Pfade unter `docs/plan/planning/`, `docs/reviews/`, `docs/user/`. Der Produktcode (`cmd/`, `internal/`) und `test/` sind im Diff nicht berührt; keine ADR geändert, keine Moves.

## 5. Plan-vs-Diff

§3 des Plans deckt die Dateien (Dockerfile, `.gitignore`, `container.mk`, vier Skripte, Ausnahmeliste, `SPEC-052`, `harness/README.md`). Darüber hinaus liegt `SPEC-055` im Diff (Architect-Commits, in §1 genannt). §1 (Abgrenzung zu `slice-v1-abschluss-container-doku`) stimmt mit dem Diff: Handbuch, README, `releasing.md` nicht berührt. Die Aufschlüsselung der Messung, die das Review (F-6) als veraltet fand, steht in §7 nicht mehr.

**Messung nachgerechnet** (`git diff --numstat d258c0b7..HEAD -- Dockerfile harness tools test`): 1176 Einfügungen, 15 Löschungen, **1191**, wie behauptet. Marke 1200: 9 Zeilen Rest; Meldeschwelle 1000 überschritten (gemeldet). Die Closure (Plan §7, Register) liegt außerhalb der Pfade und kostet nichts. Jede weitere Änderung unter `harness/` oder `tools/` zählt dagegen: Eine Korrektur zu V-2 als Satzänderung in `SPEC-055` passt (wenige Zeilen), ein neuer Fall in der Gegenprobe nicht sicher.

**Äquivalenz-Aussagen in §7 nachgeprüft** (Register `gruener-mutant-faelschlich-aequivalent` steht bei 5×):
- `DIST` als reguläre Datei: Invariante (`mkdir -p` auf eine Datei endet ungleich 0, vor dem ersten `docker`) und Ort (`binaries.sh`) stimmen. Die Mutation `p3-vor-docker` (`exit 2` entfernt) lässt `p3-dist-datei` grün, `%P3dist` färbt ihn; beides gefahren.
- `LC_ALL=C sort` zu `sort`: Die sechs Namen haben dasselbe Präfix `pgwire-recorder_<VERSION>_` und unterscheiden sich in Betriebssystem und Architektur; ich habe je einmal gemischte Eingabe unter `C`, `C.UTF-8`, `en_US.UTF-8`, `de_DE.UTF-8` sortiert, dieselbe Reihenfolge. Äquivalent, Begründung trägt (kein Zähler-Zuwachs).

## Findings

| ID | Kategorie | Befund | Zuständigkeit |
|---|---|---|---|
| V-1 | LOW | §7 *Belege des Implementers* nennt als Stand „HEAD nach 7e5ecb79“ und für `make gates` 937 s. Der grüne Lauf am Stand e365401f (nach der Nacharbeit; Nachweis `c186527d…` passt zum Arbeitsbaum) steht dort weder mit Stand noch mit Dauer (laut Übergabe 927 s); ebenso fehlen die Läufe von `test-container-gegenprobe`, `abdeckung-check` und `gates-schnell` nach der Nacharbeit. Ein Beleg, den die DoD verlangt (`make gates` grün), ist damit nur durch den Hash belegt, nicht in der Tabelle am richtigen Stand. | Implementer (Closure: Läufe am Endstand eintragen; §7 ist nicht in den Pfaden der Marke) |
| V-2 | INFO | `SPEC-055` Punkt 2 sagt zu: „eine Variable der Umgebung gilt nicht“ (neben dem Vorrang der Befehlszeile). Ich habe es von Hand bestätigt (`VERSION=dev make test-container` fängt an zu bauen, wird nicht abgelehnt), aber kein Fall und keine Mutation der Gegenprobe belegt den Satz (`AGENTS.md` §3.10, §3.11). Der Vorrang der Befehlszeile hat den Fall `p2-dev-gate`. | Architect: Satz enger fassen oder streichen (wenige Zeilen, passt in den Rest der Marke), oder Fall nachziehen (Marke prüfen) |

## Hinweise (kein Befund)

1. Ich habe zur Kontrolle des Satzes in V-2 `VERSION=dev make test-container` per Hintergrundjob nach 8 s abgebrochen. Danach: keine Container, Netze, Tags oder Temp-Verzeichnisse (geprüft). Der Abbruch zeigt nebenbei, dass der Aufräum-Pfad bei einem von außen beendeten Lauf trägt.
2. Der Flake `p7-int-treiber` im zweiten Gate-Lauf (Klasse `gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`, dritte Berührung, geplant auf `slice-harness-gegenprobe-matrix-p7-flanke`) gehört nicht zu diesem Slice; in meinen Läufen trat er nicht auf (die Mutationen von `p7-int-treiber` gehören nicht zu meiner Stichprobe).
3. Nicht von mir wiederholt: der volle `make gates` (Nachweis passt zum Stand), die übrigen Mutationen der Tabelle (Stichprobe der Gruppen 1 bis 5, 7, 9, 11, 12), und die Zusagen mit akzeptierten Negativen (Smoke der anderen Plattform, `SIGKILL`, `-j`).

## 6. Risiko-Ausgänge (Vorschlag zur Closure, nicht eingetragen)

- Risiko 1 (Handbuch beschreibt vor der Software): entfallen, steht bereits im Plan (Schnitt, §6 des Nehmers).
- Risiko 2 ([`ADR-0043`](../plan/adr/0043-binaries-und-container-gate.md) `Proposed`): entfallen, steht bereits im Plan (`Accepted`).
- Risiko 3 (Gate-Laufzeit): **entfallen** mit Begründung: Die Probe ist gelaufen, der Bericht nennt die Dauer (Gate etwa 1 Minute, Gegenprobe etwa 4 Minuten, voller Lauf 927 s am Endstand; V-1 trägt sie ein); kein Folgebedarf, die Zahl ändert sich nur mit weiteren Gates.
- Risiko 4 (Abdeckung): **entfallen** mit Begründung: `make abdeckung-check` ist grün, ohne dass eine Tabelle `LH-FA-16`, `LH-QA-01`, `LH-QA-03` führt, und verlangt keinen Beleg ([`ADR-0033`](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md) erlaubt, verlangt nicht); keine Deklaration eingetragen.
- Lerneintrag: Kandidat ist die Klasse des Reviews F-1 (Zusage am echten Artefakt, das der Test über eine Attrappe nur mittelbar sieht; hier gelöst durch einen Wert, der nicht die Vorgabe des Codes ist, plus Smoke am echten `Dockerfile`). Ob sie eine Verkörperung braucht oder in den Zähler geht, entscheidet der Planner.

## Übergabe

An den Planner. V-1 geht mit der Closure an den Implementer, V-2 an den Architect (INFO, blockiert nicht). Kein Produktcode von mir, keine Änderung an Plan, Vertrag oder Skripten.
