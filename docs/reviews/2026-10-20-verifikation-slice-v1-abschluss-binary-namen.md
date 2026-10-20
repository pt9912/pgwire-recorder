# Verifikationsbericht: slice-v1-abschluss-binary-namen — 2026-10-20

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gegen Plan, DoD und Entscheidung?

**Gegenstand:** Diff `de41645c..int-namen` (Kopf `b70e9c95`), Zweig `ver-namen`. Der Review-Report und die Nacharbeit zu F-1 (`c1326b0c`) sind im Diff enthalten.
**Maßstab:** Slice-Plan (§1, §2 DoD, §3, §6, §7), [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md) (Accepted), [SPEC-055](../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 2 und 3, [SPEC-067](../../harness/vertraege/SPEC-067-release.md) Punkt 4, 5 und 9 (Verträge im Diff unverändert), Review-Report vom selben Tag.

## Ergebnis in einem Satz

Beide Liefer-Punkte der DoD sind am Endstand belegt, die Mutationen aus §7 sind in der Stichprobe rot gesehen, die unveränderte Kopie ist grün, und das Handbuch stimmt mit einem echten Bau überein. Befunde: kein HIGH, kein MEDIUM, ein LOW (`make release-gegenprobe` in einem Zug nicht belegt, Vorbefund `m-p16-*` nicht von mir bestätigt), drei INFO.

## Selbst gefahrene Sensoren (nacheinander)

| Sensor | Stand | Ergebnis |
|---|---|---|
| `GP_ENTDECKE=1 GP_MUT=^(p5-name-altes-schema\|p3-name\|p3-exe\|p3-sums-zeile-linux/amd64\|p3-dist-andere-version\|p3-dist-altes-schema\|p3-versionsdatei)$ bash tools/harness/container-gegenprobe.sh` | `b70e9c95` | Exit 0, „unveränderte Kopie grün, je Mutation genau die benannten Fälle rot“; die roten Mengen entsprechen der Tabelle in §7 (`p3-name` 12 Fälle, `p3-exe` 4, `p3-versionsdatei` nur `p3-namen`, `p3-dist-altes-schema` allein, `p5-name-altes-schema` 15 Fälle) |
| `GP_ENTDECKE=1 GP_MUT=^(m-p4-altes-schema-angenommen\|m-p9-achter-anhang\|m-p4-liste-nicht-gezaehlt\|m-p9-name-fehlt-windows-arm64\|m-p5-summe-name-os)$ bash tools/harness/release-gegenprobe.sh` | `b70e9c95` | „grün (5 Läufe)“; `m-p4-altes-schema-angenommen` → `p4-altes-schema` allein; `m-p4-liste-nicht-gezaehlt` → `p4-datei-zu-viel` allein; `m-p9-achter-anhang` → `p9-anhaenge`, `p9-release-existiert-gleich`, `p11-wiederanlauf-release-teilweise`, `p11-wiederanlauf-release-vollstaendig`; `m-p9-name-fehlt-windows-arm64` mit `p4-namen-je-plattform-windows-arm64` |
| `GP_MUT=^ok$ bash tools/harness/release-gegenprobe.sh` | `b70e9c95` | „grün (1 Läufe)“ (unveränderte Kopie) |
| `make binaries DIST=/tmp/vd VERSION=1.2.3` in leerem `DIST` (echtes Docker) | `b70e9c95` | „6 Binaries und SHA256SUMS“; `ls -A` = sieben Einträge, `sha256sum -c` sechs mal `OK`, `version` gibt `pgwire-recorder 1.2.3`; zweiter Aufruf mit `VERSION=2.0.0` im selben `DIST`: `binaries: Eingabe: DIST ... ist nicht leer`, Exit 2 von `make`, sieben Einträge unverändert |
| `make docs-check` | `b70e9c95` | 1293 Dateien, 0 Befunde |

`make gates` nicht gefahren (am Host grün, Auftrag).

## Von Hand gebrochen (echte Datei im Klon, danach `git checkout`)

| Bruch | Lauf | Ergebnis |
|---|---|---|
| `binaries.sh`: Zeile `# P3name` zurück auf `pgwire-recorder_${VERSION}_${os}_${arch}` | `GP_MUT=^nomatch$ container-gegenprobe.sh` | rot: „die unveränderte Kopie besteht nicht“ mit `p2-gueltig`, `p2-vorgabe-test`, `p3-exe`, `p3-name-ohne-version`, `p3-name-trenner`, `p3-namen`, `p3-sums-*` u. a. Der Test ist ohne den Fix aus dem richtigen Grund rot. |
| `release.sh`: `erwartete_namen` gibt `pgwire-recorder_${os}_${a}$e` aus | `GP_MUT=^ok$ release-gegenprobe.sh` | rot: „unveränderte Kopie ist nicht grün“ (`p10-*`, `p11-*`, `p16-*` u. a.) |

Beide Dateien sind danach wiederhergestellt (`git status` leer).

## DoD-Stand

| Punkt | Verdikt |
|---|---|
| 1: `binaries.sh`, `run-container.sh`, `release.sh` schreiben und erwarten `pgwire-recorder-<os>-<arch>`, `.exe` nur bei `windows`; Fälle und Mutationen aus `SPEC-055` Punkt 3 und `SPEC-067` Punkt 4 samt alter Namen nachgeführt; Mutation und roter Fall je Zusage im Bericht; Zeile in `harness/README.md` | erfüllt. Diff: `binaries.sh` Zeile `# P3name`, `run-container.sh` `binaer()`, `release.sh` `erwartete_namen` und Summe im Schritt `image` (`pgwire-recorder-linux-$arch`). Die Fälle `p3-namen`, `p3-exe`, `p3-name-ohne-version`, `p3-name-trenner`, `p3-sums-zeile-*`, `p3-dist-andere-version`, `p3-dist-altes-schema`, `p4-altes-schema`, `p4-namen-je-plattform-*` stehen in den Gegenproben; Stichproben rot gesehen (oben). Die Zeile in `harness/README.md` §Sensors nennt Schema, `DIST` bei anderer `VERSION` und [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md). §7 trägt die Tabelle Zusage · Mutation · roter Fall. |
| 2: Handbuch (`make binaries`) und `releasing.md` im Ist-Zustand ohne Chronik; Abschnitt „Binary vom GitHub-Release“ und `README.md` bleiben; Beispiele laufen gegen das gebaute Binary | erfüllt. Namensblock des Handbuchs = `ls -A` des echten Baus (Diff leer); Aussagen zu `VERSION=1.2.3`, `sha256sum -c`, nicht leerem `DIST` mit anderer `VERSION` von mir bestätigt. Der Abschnitt zu `v0.1.0` im Handbuch (zehn Treffer des alten Schemas) ist unverändert; `README.md` nicht im Diff. Keine Verweise auf Spezifikation oder Slices im Handbuch-Diff. `releasing.md` §2, §3, Version 0.11 und Historienzeile nachgezogen. |
| 3: `make gates` grün; `make test-container-gegenprobe`, `make release-gegenprobe` am Endstand | teilweise belegt: `make gates` am Host grün (Auftrag). `test-container-gegenprobe` im Auswahllauf grün, die unveränderte Kopie dort grün. `release-gegenprobe` nur in Teilen belegt, siehe F-V1. |
| 4: Review-Report | erfüllt (`2026-10-20-review-slice-v1-abschluss-binary-namen.md`, 0 HIGH, 0 MEDIUM, F-1 LOW behoben in `c1326b0c`, F-2 und F-3 INFO). Das Nachziehen von F-1 (Mutation `p5-name-altes-schema`) habe ich rot gesehen. |
| 5 bis 8: Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen, gehören zur Closure, siehe unten |

## Plan-vs-Code, Entscheidung, Sendungen

- Der Diff berührt genau die Dateien aus §3 (sieben Zeilen der Tabelle) plus Plan und Review-Report. `Dockerfile`, `harness/mk/container.mk`, `cmd/`, `internal/`, die Verträge, `README.md`, ADR sind nicht im Diff. Die Abgrenzung §1 hält.
- [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md): Schema `pgwire-recorder-<os>-<arch>`, keine Versionsdatei in `DIST` (Mutation `p3-versionsdatei` rot), Version nur im Binary (`-X main.version`, `build-arg` unverändert, Mutation `p3-build-arg`). Image-Tag und Registries unverändert.
- Metrik (`AGENTS.md` §3.14): `git diff --numstat de41645c..b70e9c95 -- tools harness/README.md docs/user docs/maintainer README.md` = 140 Einfügungen + 77 Löschungen = **217** (§7 nennt 216 am Stand `9492d6c4`; die eine Zeile mehr ist die Mutation aus `c1326b0c`), Marke 800 nicht erreicht, Schätzung 360 bis 510 unterschritten. Zwei Liefer-Punkte, zwei Schichten (Harness-Werkzeug, Nutzer- und Wartungs-Doku) wie in §8.
- Sendungen (§3.13): `slice-v1-abschluss-homebrew` und `slice-erster-release-homebrew-nachweis` nennen die Kennung des Gebers je dreimal (Bindung, Start-Trigger, Zählung); der Diff trägt in diese Pläne nichts neu ein, der Commit `64a6a977` vor dem Slice hat es getan. Die Adresse für den Abschnitt zu `v0.1.0` ist `slice-v1-abschluss-release-v0-1-1` (nennt den Geber).
- Fundstellen des alten Schemas (`git grep -nI pgwire-recorder_` ohne Baseline, `docs/reviews/`, `done/`): `container-gegenprobe.sh` 7, `release-gegenprobe.sh` 2, `benutzerhandbuch.md` 10 (Abschnitt `v0.1.0`), `SPEC-055` 1, `SPEC-067` 2, [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md) 3, `slice-v1-abschluss-release-v0-1-1` 1, dieser Plan 7. In `binaries.sh`, `run-container.sh`, `release.sh`, `harness/README.md`, `releasing.md`, `README.md` kein Treffer. Alle Treffer nennen das alte Schema als alt oder `v0.1.0`.

## Befunde

| ID | Kategorie | Befund | Ort | Zuständig |
|---|---|---|---|---|
| F-V1 | LOW | Die DoD verlangt `make release-gegenprobe` am Endstand. Der Implementer hat sie in fünf Teilen gefahren, im Teil `^m-p16-` sind sechs Mutationen rot (allein grün; am Stand `de41645c` ebenso rot, „Vorbefund der Umgebung, Ursache nicht untersucht“). Ich habe `m-p16-*` nicht bestätigt: Mein Lauf mit `GP_MUT=^m-p16-` (130 Mutationen) füllte den Datenträger bis auf 207 MB und wurde von mir abgebrochen, das Temp-Verzeichnis (3,6 GB) entfernt; der Platz ist wiederhergestellt. Der volle Lauf in einem Zug ist damit weder vom Implementer noch von mir belegt. Der Slice ändert den Schritt `scan` nicht, ein Zusammenhang mit den Namen ist unwahrscheinlich, aber unbelegt. | Plan §7 *Läufe* (`release-gegenprobe`) | Planner (Urteil: Lauf auf Maschine mit Platz verlangen oder die Umgebung als Carveout/Beobachtung führen) |
| F-V2 | INFO | Ein Auslöser des Plattenüberlaufs ist der Arbeitsordner der Gegenprobe (rund 40 MB je Mutation, 130 Mutationen in `m-p16-*`); der Hinweis steht schon in §7, hier nur als Beobachtung für das Register. | `release-gegenprobe.sh` | Register |
| F-V3 | INFO | §7 *Befundzeilen* nennt für `container-gegenprobe.sh` sechs Treffer des alten Schemas am Stand `9492d6c4`; am Endstand sind es sieben (die Mutation `p5-name-altes-schema` aus `c1326b0c`). Die Zusage „nur das alte Schema als alt genannt“ bleibt wahr. | Plan §7 | Implementer / Closure (Zahl nachziehen) |
| F-V4 | INFO | §7 nennt die Metrik 216 am Stand `9492d6c4`; am Endstand 217. Die Marke ist nicht berührt. | Plan §7 | Closure |

Keine DoD-Verletzung unter den Liefer-Punkten, keine Entscheidung verletzt, keine Randform vom Implementer entschieden. Das Fehlen eines Belegs, den die DoD verlangt, ist nur F-V1.

## Was für die Closure (Planner) noch fehlt

1. DoD-Häkchen setzen (Liefer-Punkte 1 und 2, Review; Punkt 3 erst nach Entscheidung zu F-V1).
2. §7: Closure-Notiz mit Lerneintrag (geschärfte Regel, neuer Sensor oder benannte Spec-Lücke; „Wird bei Closure gefüllt“ steht noch); die Zeile wird mit Herkunfts-Anker `liegt in` nur geschrieben, wenn etwas verkörpert wurde.
3. Beobachtungs-Register (`docs/plan/planning/observations/`): fortschreiben (Zähler `BEO-REPO/verhalten-nur-unter-linux-geprueft`, `folge-slice-adresse-nimmt-nicht-an`, `folge-slice-liefert-handbuch-teil-nicht`, ggf. Plattenbedarf der Gegenprobe) oder „keine Beobachtung angefallen“ notieren.
4. Risiko-Ausgänge aus §6: (a) Fundstellen des alten Schemas: Ausgang *entfallen* mit Begründung (grep am Endstand nur Treffer, die das alte Schema als alt nennen, Zahl 7 statt 6 nachziehen); (b) Abschnitt `v0.1.0` in Handbuch und `README.md`: *weiter offen*, übergeben an `slice-v1-abschluss-release-v0-1-1` (dort bereits §1, DoD 2, §8).
5. Drei Paarungen (Anker · Folge-Slice · Register) bestätigen; die Folge-Slices tragen die Kennung des Gebers bereits.
6. Danach der reine `git mv` nach `done/` in eigenem Commit (`AGENTS.md` §3.3), Welle-Plan `welle-v1-nachtrag` nachziehen.

Dieser Bericht ist ein Lauf-Beleg; er repariert nichts.
