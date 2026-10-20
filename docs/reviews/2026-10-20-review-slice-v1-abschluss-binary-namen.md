# Review-Report: slice-v1-abschluss-binary-namen — 2026-10-20

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Entscheidung ADR-0054, Verträge `SPEC-055` Punkt 2 und 3, `SPEC-067` Punkt 4, 5 und 9, `AGENTS.md` §3.7, §3.9 bis §3.14).

**Gegenstand:** `slice-v1-abschluss-binary-namen`, Diff `de41645c..impl-namen` (Kopf `02a04691`), sechs Commits, neun Dateien.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-20

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, keine Lifecycle-Pfade, und Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill). Das `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-v1-abschluss-binary-namen` (§1 Ziel und Abgrenzung, §3 Plan, §6 Randformen, §7 Belege des Implementers)
- ADR-0054 (Accepted), `SPEC-055` Punkt 2 und 3, `SPEC-067` Punkt 4, 5, 9 und 19 (3)
- berührte Anforderungen `LH-FA-01`, `LH-FA-19`, `LH-QA-03`
- `AGENTS.md` (Hard Rules §3.7, §3.9 bis §3.14)

**Läufe des Reviewers** (engste Sensoren, nacheinander, im Dev-Container `pgr-dev2`, Datenträger vor jedem Lauf geprüft, nie unter 3 GB frei): `tools/harness/container-gegenprobe.sh` mit `GP_MUT` auf acht Mutationen von Punkt 3 (`p3-exe-nur-arm64`, `p3-name-version`, `p3-name-trenner-hinten`, `p3-name-gross`, `p3-name-ziffern`, `p3-versionsdatei`, `p3-sums-zeile-linux/arm64`, `p3-dist-andere-version`): grün, rote Mengen wie in §7 genannt; `tools/harness/release-gegenprobe.sh` mit `GP_MUT=^(m-p4-|m-p9-name-fehlt-)` (15 Läufe): grün; die sechs `m-p16-*` aus §7 allein: grün; `^m-p16-` im Ganzen nach 40 von 130 Mutationen abgebrochen (4,8 GB auf 3,2 GB frei); eine Handmutation in einer Kopie des Arbeitsbaums (`binaer()` in `run-container.sh` mit altem Namen, danach zurückgesetzt) gegen `GP_MUT=^ok$`. Kein `make gates`, kein Integrationslauf.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Zusage „`run-container.sh` erwartet `pgwire-recorder-<os>-<arch>`“ (DoD 1, §3) hat keine benannte Mutation in §7 und keine in `container-gegenprobe.sh`. Eine Handmutation (`binaer()` mit altem Namen) färbt die unveränderte Kopie rot (`p2-gueltig`, `p4-gruen`, `p5-*`), die Zusage ist also geprüft, aber nicht als Mutation geführt und nicht vom Implementer gefahren. | `AGENTS.md` §3.10 (Mutation je Zusage) | `tools/harness/run-container.sh` · `local n="pgwire-recorder-$1-$2"`; `docs/plan/planning/in-progress/slice-v1-abschluss-binary-namen.md` · „Zusagen und Mutationen“ | ja — `make test-container-gegenprobe` mit der Handmutation | Zusage einer geänderten Datei ohne benannte Mutation |
| F-2 | INFO | §7 nennt die sechs roten `m-p16-*` einen „Vorbefund der Umgebung“ und sagt „Ursache nicht untersucht“. Beobachtet: Der Arbeitsordner der Gegenprobe wächst um rund 40 MB je Mutation (4,8 GB auf 3,2 GB frei nach 40 Mutationen, bei 130 Mutationen des Teils und 535 insgesamt); die sechs Mutationen sind allein gefahren grün. Die Datenträger-Erschöpfung gegen Ende des Laufs ist damit die naheliegende Ursache, am Endstand aber nicht in einem Zug belegt. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-binary-namen.md` · „ein Vorbefund der Umgebung, nicht dieses Slice“ | nein — ein voller Lauf braucht rund 6 GB, auf diesem Host nicht fahrbar | Gegenprobe-Lauf in einem Zug auf dem Host nicht fahrbar; Ursache eines Vorbefunds nicht belegt |
| F-3 | INFO | `SPEC-055` Punkt 3 verlangt „je Element ein Fall und eine Mutation, die nur ihn bricht“. Die Mutationen der Zeichenmenge (`p3-name-gross`, `p3-name-ziffern`, `p3-exe-ohne-punkt`) und des Namens färben jeweils mehrere Fälle (`p3-namen` und die Folgefälle `p3-exe`, `p3-sums-*`); einen eigenen Fall je Zeichenklasse gibt es nicht, `p3-namen` trägt sie gemeinsam. §7 führt die Mitläufer offen mit „+“; die Invariante (Name gleich einer festen Liste) macht die Klassen redundant. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `harness/vertraege/SPEC-055-container-und-binaries.md` · „je Element ein Fall und eine Mutation, die nur ihn bricht“; `tools/harness/container-gegenprobe.sh` · `mut p3-name-ziffern` | ja — `GP_ENTDECKE=1` zeigt die roten Mengen | Wortlaut des Vertrags weiter als die Fälle |

### Zu den Schwerpunkten des Auftrags

1. **Altes Schema.** `git grep pgwire-recorder_` (ohne vendored Baseline, `done/`, `docs/reviews/`) am Kopf: `binaries.sh`, `run-container.sh`, `release.sh`, `harness/README.md`, `releasing.md` ohne Treffer. Die übrigen Treffer einzeln gelesen: Handbuch Zeilen 101 bis 121 (`v0.1.0`, Bestand nach §1), ADR-0054 (Beschreibung des Alten), `SPEC-055` Punkt 3 und `SPEC-067` Punkt 4 und 19 (3) (nennen das alte Schema als nicht angenommen), `container-gegenprobe.sh` (Fall `p3-dist-altes-schema`, Mutationen `p3-name*`), `release-gegenprobe.sh` (Modus `altes-schema` der Attrappe, Mutation `m-p4-altes-schema-angenommen`), Plan und `slice-v1-abschluss-release-v0-1-1` (übergebener Bezug). Kein Treffer ist ein gültiger Name. Auch `${VERSION}_`-Muster außerhalb davon: keine.
2. **Mutationen.** Stichproben (siehe Läufe): jede färbt genau die in §7 genannten Fälle, die unveränderten Kopien sind grün. `p3-name-version` lässt `p3-name-trenner` grün (wie behauptet). `m-p4-liste-nicht-gezaehlt` und `m-p4-altes-schema-angenommen` färben je genau einen Fall. `m-p9-name-fehlt-<os>-<arch>` färbt den Fall `p4-namen-je-plattform-<os>-<arch>` der Plattform. Kein grüner Mutant, keiner als äquivalent geführt.
3. **Handbuch und `releasing.md`.** Abschnitt `make binaries`: Namen, `VERSION` im Binary, „muss leer sein“ auch bei Dateien anderer `VERSION` stimmen mit `binaries.sh` und mit den Fällen `p3-dist-andere-version`, `p3-name-ohne-version` überein; Ist-Zustand, keine Chronik, kein Verweis auf Spezifikation oder Slices. `releasing.md` §2 und §3 stimmen, Version 0.11 und Historienzeile vorhanden.
4. **Plan.** §3 trägt die geänderten Dateien und Fälle des Diffs; §1, §6 und der Kopf sind vom Diff unberührt und treffen weiter zu (§6 ohne neue Randform im Code-Commit); die Verträge und ADR-0054 sind unverändert. §7 steht allein im letzten Commit. Metrik nachgerechnet: `git diff --numstat de41645c..9492d6c4 -- tools harness/README.md docs/user docs/maintainer README.md` = 139 + 77 = 216, Marke über 800 nicht erreicht.
5. **`m-p16-*`.** Die sechs Mutationen sind allein grün (nachgefahren). Der Diff berührt weder den Schritt `scan` noch dessen Fälle. Mit der Wachstumsrate des Arbeitsordners ist eine Erschöpfung des Datenträgers die plausible Ursache (F-2); ein echter Fehler, den dieser Slice berührt, ist nicht erkennbar.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/binaries.sh`, `tools/harness/release.sh` (`erwartete_namen`, Liste `P4liste`, Anhänge `P9anhaenge`, Summe im Schritt `image`), Syntax | geprüft, ohne Befund |
| `tools/harness/run-container.sh` (`binaer()`) | geprüft, ohne Befund außer F-1 |
| `tools/harness/container-gegenprobe.sh` (`bname`, `namen_soll`, `namen_trenner`, `namen_gleich_summen_verschieden`, Fälle und Mutationen von Punkt 3 und 4) | geprüft, ohne Befund außer F-3 |
| `tools/harness/release-gegenprobe.sh` (Attrappe von `make`, `NAMEN`, `ALLE7`, Fälle `p4-altes-schema`, `p4-namen-je-plattform-*`, Mutationen zu Punkt 4, 5, 9) | geprüft, ohne Befund |
| `harness/README.md` (Zeile `make binaries`), `docs/maintainer/releasing.md`, `docs/user/benutzerhandbuch.md` (Abschnitt `make binaries`) | geprüft, ohne Befund |
| `harness/vertraege/`, `docs/plan/adr/` (unverändert im Diff; keine `Accepted`-ADR überschrieben) | geprüft, ohne Befund |
| Sendungen und Zählung (§3.13): der Diff trägt keine Adresse neu oder geändert ein | geprüft, ohne Befund |
| Kommentare in den geänderten Skripten (`AGENTS.md` §3.7, §3.11) | geprüft, ohne Befund |
| Suppressionen (`nolint`), Core-Reinheit (`internal/hexagon`) | geprüft, ohne Befund; der Diff berührt weder Go-Code noch `Dockerfile` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage einer geänderten Datei ohne benannte Mutation · Gegenprobe-Lauf in einem Zug auf dem Host nicht fahrbar; Ursache eines Vorbefunds nicht belegt · Wortlaut des Vertrags weiter als die Fälle

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; F-1 ist LOW, F-2 und F-3 sind INFO.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation; DoD-Konformität prüft der Verifier separat (F-2 betrifft dessen Lauf von `make release-gegenprobe`).
