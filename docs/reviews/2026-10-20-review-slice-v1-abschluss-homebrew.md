# Review-Report: slice-v1-abschluss-homebrew — 2026-10-20

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (`SPEC-067` Punkt 1c und 19, `AGENTS.md` §3.7, §3.9 bis §3.14).

**Gegenstand:** `slice-v1-abschluss-homebrew`, Diff `8dbe6a14..impl-brew` (Kopf `01850212`), vier Commits (`43364350`, `e964ff19`, `2dfcd669`, `01850212`), sieben Dateien.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-20

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, keine Lifecycle-Pfade, und Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill). Das `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-v1-abschluss-homebrew` (§1 bis §4, §6 *Randformen*, §8); §7 trägt keine Belege des Implementers (der Implementer wurde vor dem Bericht abgebrochen), die Zusagen und die Mutation je Zusage sind aus dem Diff und `SPEC-067` Punkt 19 (9) abgeleitet
- [ADR-0053](../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md), [ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md), `SPEC-067` Punkt 1c und 19, `SPEC-055` Punkt 3
- berührte Anforderungen `LH-FA-19`, `LH-FA-01`
- `AGENTS.md` (Hard Rules §3.7, §3.9 bis §3.14)

**Läufe des Reviewers** (nur engste Sensoren, im Dev-Container `pgr-dev2`, je Lauf unter 3 Minuten, kein voller Lauf von `make release-gegenprobe`, kein `make gates`): `GP_MUT` auf `m-p19-vorabversion-nur-rc` und `m-p19-push-force` (grün); auf `m-p1-homebrew-*`, `m-p19-aktuell-immer-neu`, `m-p19-hoeher-ueberschrieben` und `m-p19-token-in-datei` (grün, 7 Läufe); zwei Sonde-Mutationen in einer Wegwerfkopie der Gegenprobe (`tools/harness/gp-tmp.sh`, danach gelöscht): Prüfung der Vorlage in `release.sh` entfernt (Ergebnis siehe F-1) und Token als Argument beim Klon (rot in `p19-token-nicht-im-argument` und `p19-token-nicht-in-url-und-config`). Von Hand: Vorlage gegen den Text in `SPEC-067` Punkt 19 (4) mit `cmp` (gleich); Job `homebrew` gegen Punkt 1c Zeile für Zeile (gleich).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Schritt hat einen Zweig für eine fehlende oder leere Vorlage (Exit 1, Stichwort `vorlage`), den weder `SPEC-067` Punkt 19 noch §6 des Plans noch der Kopfkommentar von `release.sh` (Liste der Stichworte) nennen und den kein Fall und keine Mutation prüft: Die Sonde-Mutation „Prüfung entfernt“ färbt keinen Fall (`m-p19-xx-vorlage-leer => nichts`). Ohne die Prüfung schriebe ein Nachtrag aus einem Klon ohne Vorlage eine leere Formel in den Tap; der nächste Lauf endete dort mit `tap-formel` (keine Zeile `version`), bis jemand den Tap von Hand repariert. | `AGENTS.md` §3.12 (Randform im Code-Commit, nicht in §6 und nicht in der Spezifikation) · §3.10 (Zusage ohne Mutation) | `tools/harness/release.sh` · `homebrew_fehler vorlage "die Vorlage $VORLAGE fehlt oder ist leer"` | ja — ein Fall mit fehlender und einer mit leerer Vorlage samt Mutation, die die Prüfung entfernt | Randform im Code entschieden, ohne Fall und Mutation |
| F-2 | LOW | `release.sh` trägt sechs Anker-Kommentare des Schrittes, auf die weder eine Mutation noch ein Fall der Gegenprobe zeigt: `P19aktuell`, `P19liste`, `P19summen-aufruf`, `P19summen-zahl`, `P19url`, `P19version` (die Mutationen dazu ersetzen den Text der Zeile). Ein Anker, den niemand liest, ist keine Kopplung. | `AGENTS.md` §3.7 (Klassen eines Kommentars: Kopplung) | `tools/harness/release.sh` · `# P19liste` | ja — `grep` je Anker in `tools/harness/release-gegenprobe.sh` | Anker-Kommentar ohne Leser |
| F-3 | LOW | `SPEC-067` Punkt 19 (9) nennt als Beispiel-Mutation „der Token steht in der URL“. Die Mutationen zum Token (`m-p19-token-im-argument-git`, `-gh`, `-in-config`, `-in-datei`, `-im-log-*`) tragen keine solche Mutation der Klon-URL. Die Sonde des Reviewers (Token als Argument des Klons) wird von den Fällen rot gefärbt, die Zusage ist also gefangen, aber nicht als benannte Mutation geführt. | `AGENTS.md` §3.10 (je Element eine Mutation) · `SPEC-067` Punkt 19 (9) | `tools/harness/release-gegenprobe.sh` · `mh m-p19-token-im-argument-git` | ja — `GP_ENTDECKE=1 GP_MUT='^m-p19-token-in-url$'` nach Ergänzung | Beispiel-Mutation der Spezifikation fehlt |
| F-4 | INFO | `SPEC-067` Kopf nennt [ADR-0053](../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) noch „Proposed“ („Stand der Entscheidung zu H“) und Punkt 19 (6) „Entscheidung des Nutzers steht aus“; die ADR ist `Accepted` (2026-10-20) und der Plan nennt das so. Vertragstext: Rückgabe an den Architect, der Implementer ändert ihn nicht. | `AGENTS.md` §3.9 · §3.12 (Vertragstext zurückgeben) | `harness/vertraege/SPEC-067-release.md` · „ist Proposed; was dort die Empfehlung des Architect ist“ | nein | Vertragstext folgt der Annahme der ADR nicht |
| F-5 | INFO | §7 des Plans trägt keine Belege des Implementers (Mutationstabelle, Läufe); der Eingang des Reviews fehlt. Die Ableitung steht in der Kopfzeile *Eingangs-Kontext*; die DoD-Frage (Zusage · Mutation · roter Fall) gehört dem Verifier. Hinweis dorthin: die Fälle `p19-*` und `m-p19-*` sind nicht vollständig gefahren worden. | Maintainability · Skill §Kontext-Eingang | `docs/plan/planning/in-progress/slice-v1-abschluss-homebrew.md` · „Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).“ | nein | Belege des Implementers fehlen im Eingang |
| F-6 | INFO | `releasing.md` sagt Verhalten von GitHub zu, das keine Gegenprobe zeigt („Re-run failed jobs“ wiederholt nur den Job `homebrew`, baut das Image nicht neu) und beschreibt den Zustand des echten Taps („im Tap liegt keine Formel“). `SPEC-067` Punkt 15 führt das Verhalten von GitHub als Grenze; die Tap-Aussage ist aus dem Klon nicht prüfbar. | `AGENTS.md` §3.11 | `docs/maintainer/releasing.md` · „im Tap liegt keine Formel“ | nein | Zusage über Verhalten außerhalb des Klons |
| F-7 | INFO | Die Marke aus §4 (Summe aus Einfügungen und Löschungen, `git diff --numstat 0256ddff..<letzter Commit> -- tools .github harness/README.md docs/maintainer`, Grenze 1300) misst 1078; die Basis `0256ddff` liegt vor `slice-v1-abschluss-binary-namen`, die Zahl enthält dessen Dateien (`binaries.sh`, `container-gegenprobe.sh`, `run-container.sh`). Die Marke ist nicht eingetreten. | `AGENTS.md` §3.14 | `docs/plan/planning/in-progress/slice-v1-abschluss-homebrew.md` · „über 1300“ | ja — `git diff --numstat` wie in §4 | Basis der Größenmarke enthält fremde Commits |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.github/workflows/release.yml` (Job `homebrew`) | geprüft, ohne Befund: Zeilen gleich Punkt 1c; `needs: release` und `if: github.event_name == 'push'` (nur nach grünem Release, nie im Trockenlauf); `permissions: contents: read` im Job, `permissions: {}` oben unverändert; Checkout per SHA mit Tag-Kommentar, gleicher Pin wie im Job `release`, `persist-credentials: false`; kein weiteres `uses:`; `HOMEBREW_TAP_GITHUB_TOKEN` nur im Schritt `Homebrew`, kein `packages: write`, kein Registry-Geheimnis |
| `tools/harness/homebrew-formula.rb.tmpl` | geprüft, ohne Befund: byte-gleich dem Text aus Punkt 19 (4); Namen `pgwire-recorder-<os>-<arch>` ohne Version ([ADR-0054](../plan/adr/0054-binaries-ohne-version-im-dateinamen.md)); `bin.install Dir["pgwire-recorder-*"].fetch(0) => "pgwire-recorder"`; vier Plattformen, kein Windows |
| `tools/harness/release.sh` Schritt `homebrew`: Eingaben und Reihenfolge der Gründe | geprüft, ohne Befund: (a) Ereignis/Ref, (b) Tag-Form, (c) Vorabversion vor den Token ohne `gh` und `git`, (d) `GH_TOKEN`, (e) `RELEASE_TAP_TOKEN`; Exit 2 nennt den Namen, nie einen Wert |
| `release.sh` Summen je Plattform | geprüft, ohne Befund: genau eine Zeile je Name, 64 Hexziffern klein, Summe je Plattform über `summe[$k]` aus der Zeile des eigenen Namens; altes Schema wird nicht angenommen; Prüfung vor dem Klon |
| `release.sh` Token | geprüft, ohne Befund: nur über `GIT_ASKPASS` (Skript liest die Umgebungsvariable), `GIT_TERMINAL_PROMPT=0`, `-c credential.helper=`; kein `set -x`; Fehlerausgabe von `git` nur mit Präfix `release: git: `, ohne Token |
| `release.sh` Tap und Fehler des Pushs | geprüft, ohne Befund: leerer Tap mit Zweig `main`, Zweig-, Pfad- und Formelprüfung, `sort -V` für höher/gleich/niedriger, Commit mit Identität und Nachricht nach Punkt 19 (5), `git push origin main` ohne `--force`; abgewiesener Push endet mit Exit 1 `tap-push`, ohne Wiederholung; ein roter Job lässt Release und Registries stehen |
| `tools/harness/release-gegenprobe.sh` | geprüft: alle Fälle und Beispiel-Mutationen aus Punkt 19 (9) haben einen Namen (bis auf F-3); Mutationen je Plattform für Adresse, Summe und Zeile in `SHA256SUMS`; stichprobenhaft sieben Läufe grün (siehe Läufe), kein voller Lauf |
| `harness/README.md` (Zeile `make release-gegenprobe`) | geprüft, ohne Befund: nennt `p1-homebrew-*`, `p19-*` und die Grenzen; keine Zusage über Verhalten, das die Gegenprobe nicht zeigt |
| `docs/maintainer/releasing.md` (§1 bis §4, §6, §7, §11, Änderungshistorie) | geprüft, bis auf F-6 ohne Befund: Reihenfolge der Gründe, Stichworte und Nachtrag stimmen mit Punkt 19 überein |
| Plan (§1, §3, §6, §8) | geprüft, ohne Befund gegen §3.9: §3 nennt die Dateien des Diffs; §8 (drei Liefer-Punkte, zwei Schichten) trägt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Randform im Code entschieden, ohne Fall und Mutation · Anker-Kommentar ohne Leser · Beispiel-Mutation der Spezifikation fehlt · Vertragstext folgt der Annahme der ADR nicht · Belege des Implementers fehlen im Eingang · Zusage über Verhalten außerhalb des Klons · Basis der Größenmarke enthält fremde Commits

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) bis zur Entscheidung der Randform (Architect: Stichwort und Exit-Code der fehlenden oder leeren Vorlage in `SPEC-067` Punkt 19) und zum Fall mit Mutation. Kein HIGH.

**Übergabe:** Findings gehen an den Implementer; F-4 geht an den Architect; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation: DoD-Konformität prüft der Verifier.
