# Review-Report: slice-harness-abdeckung-gate — 2026-10-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-harness-abdeckung-gate`, Diff `cbeb3381..HEAD` (Commits 01523fa0, cacaf26a (Architect), d7fede16; Pfade `tools/test`, `tools/harness/lint-gegenprobe.sh`, `harness/`, `docs/user`, `.d-check.yml`, Plan). Die Commits zu `slice-harness-devcontainer-austausch-repo` im selben Bereich gehören einem parallelen Lauf und sind nicht Gegenstand.

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

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

- Slice-Plan `slice-harness-abdeckung-gate` (§1 bis §8, Stand HEAD)
- `SPEC-048` Punkte 1 bis 8 (Stand HEAD, mit den Entscheidungen des Architect aus cacaf26a)
- aktive ADRs: [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0033](../plan/adr/0033-gate-nachweise-in-der-abdeckung.md)
- berührte `LH-*`-IDs: `LH-QA-07`
- `AGENTS.md` (Hard Rules §3.7 bis §3.14)

**Eigene Mutationen** (Kopie von `tools/test/` im Scratchpad, genau eine Ersetzung mit geprüfter Trefferzahl, die unveränderte Gegenprobe gegen den Mutanten). 41 Mutanten gefahren, 36 rot, 5 grün (siehe F-1, F-2). Rot gesehen u. a.: i = 0, n zweistellig, `b < 1`, `a > b + 1`, lexikalische Sortierung der Zeile, Elementindex ohne Zahlsortierung (äquivalent bis 9 Elemente), Zeile als Schlüssel entfernt, Symlink gelesen, letzter und erster Teil bei der Vollständigkeit, `exit 0` am Ende, `veraltet=0`, `exit 0` nach Fehlformen, Exit 1 statt 2 an zwei der drei Gründe, Reihenfolge (3) und (4) getauscht, `LH-RB` aus dem Muster, `+`, Leerzeichen und `01` an i und an n, Anker, Tab und Leerzeichen in der eingerückten Gate-Deklaration, `#` ohne Leerzeichen, Ordnung der Tabellen E2E und Gate, E2E-Muster tiefer, `chmod`, erste statt letzte festlegende Deklaration, Ausgang 2 mit Form der Kennung.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Punkt 1 von `SPEC-048` sagt jetzt „eingerückt (Leerzeichen oder Tab)“ auch für Test-Deklarationen zu. Der Fall `eingerueckt` prüft nur Tab. Der Mutant `[ \t]+` → `[\t]+` im `test_awk` (Leerzeichen nicht mehr erkannt) bleibt grün. | `AGENTS.md` §3.10 (Aufzählung je Element) | `tools/test/abdeckung-gegenprobe.sh` · `abgelehnt eingerueckt internal/a_test.go` | ja — Mutant `[\t]+` in `tools/test/abdeckung.sh`, `make abdeckung-gegenprobe` bleibt grün | Aufzählung in der Zusage ohne Fall je Element |
| F-2 | MEDIUM | Die Zahlform-Fälle decken `+`, ein Leerzeichen nach `-`, `01`, `0`, mehrstellig, nicht aber ein Minus als Vorzeichen (`Messung--1-von-2`, `Messung-1-von--2`), Tab als Leerraum und nachgestellten Leerraum. Vier Mutanten (`-?` an i und n, `[\t]?` an i und n, ` ?` hinter n) bleiben grün, obwohl Punkt 8 „Vorzeichen, Leerraum“ als Fehlformen zusagt. | `AGENTS.md` §3.10 (Aufzählung je Element); `SPEC-048` Punkt 8 | `tools/test/abdeckung-gegenprobe.sh` · `for zf in 'Messung-10-von-12'` | ja — die Mutanten am Muster `pform` in `tools/test/abdeckung.sh` | Aufzählung in der Zusage ohne Fall je Element |
| F-3 | LOW | §7 des Plans nennt die Größe 823 (Gegenprobe 474, `SPEC-048` nicht in der Zerlegung). Am Stand HEAD sind es 905 (Gegenprobe 524, `SPEC-048` 32, vom Architect). §4 sagt zugleich, `SPEC-048` zähle nicht, der Pfad steht aber in der Metrik. | `AGENTS.md` §3.9, §3.14 | `docs/plan/planning/in-progress/slice-harness-abdeckung-gate.md` · `Summe aus Einfügungen und Löschungen **823**` | ja — `git diff --numstat cbeb3381..HEAD -- <Pfade der Metrik>` ergibt 905 | Plan folgt dem Diff nicht (Größe) |
| F-4 | LOW | Das Rezept von `abdeckung-check` in `harness/mk/abdeckung.mk` (`--check`) und die Zeile `GATE_CHECKS += …` sind in keinem Fall der Gegenprobe gezeigt; ein Mutant ohne `--check` schreibt die Tabellen im Gate und endet grün. Der Plan führt den Anschluss an `GATE_CHECKS` als Grenze, das Rezept mit `--check` nicht. Bestand, in diesem Diff nur im Kommentar berührt. | `AGENTS.md` §3.10, §3.11 | `harness/mk/abdeckung.mk` · `@bash tools/test/abdeckung.sh --check` | ja — Gegenprobe in einem Temp-Baum wie bei `kopf-check-gegenprobe` | Verdrahtung des Aufrufers ohne Mutation |
| F-5 | INFO | Der Kommentar in `harness/mk/abdeckung.mk` nennt `stat` unter den Programmen des Laufs; `tools/test/abdeckung.sh` ruft es nicht auf. | `AGENTS.md` §3.7 | `harness/mk/abdeckung.mk` · `(bash, awk, sort, find, cmp, stat)` | ja — `grep -w stat tools/test/abdeckung.sh` ohne Treffer | Kommentar nennt mehr als der Code tut |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/test/abdeckung.sh` (Gate-Kopf, Teile, Vollständigkeit erster/mittlerer/letzter, Ordnung bytegenau nach Pfad/Zeile/Element, Reihenfolge der Gründe, Symlink, Ausgang 2 in drei Gründen und Reihenfolge, Exit 1 bei Fehlform und bei veralteter Tabelle) | geprüft, ohne Befund außer F-1, F-2; die Mutanten dazu sind rot |
| `tools/test/abdeckung-gegenprobe.sh` (Kopf, Äquivalenz-Liste, Fallnamen im Kopf von `abdeckung.sh`) | geprüft, ohne Befund; alle genannten Fallnamen kommen vor, äquivalente Mutanten und akzeptierte Negative stimmen in Vertrag, Plan und Kopf überein |
| `harness/vertraege/SPEC-048-abdeckung.md` | geprüft, ohne Befund; Punkte 6 bis 8 stimmen mit dem Code, Ausnahme die Aufzählungen aus F-1, F-2 |
| `harness/sensors/abdeckung-check.md`, `harness/README.md` (Gate-Index) | geprüft, ohne Befund; Exit-Tabelle, fünf Tabellen und Grenze 3 stimmen mit Code und Gegenprobe |
| `docs/user/abdeckung-*.md`, `.d-check.yml` | geprüft, ohne Befund; E2E- und Unit-Tabelle unverändert (byte-gleich), `LH-QA-07` `teilweise` nur in der Gesamtsicht (`1/3 Gate; 2/3 —; 3/3 Gate`), nicht in `abdeckung-vollstaendig.md`; `make abdeckung-check` Exit 0, `make docs-check` 0 Befunde. `make doc-trace` nicht gefahren (Docker) |
| `tools/harness/lint-gegenprobe.sh` (Kopf) | geprüft, ohne Befund; reine Kommentarzeilen, keine Änderung am Ablauf |
| Plan (§1, §3, §6, §8, Nachzählen §3.13, Adresse `slice-harness-integration-volumes`) | geprüft, ohne Befund außer F-3; zwei Schichten, drei Liefer-Punkte |
| `AGENTS.md` §3.2 (Suppression) | geprüft, ohne Befund; kein `nolint`, keine Lockerung eines Gates |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Aufzählung in der Zusage ohne Fall je Element · Plan folgt dem Diff nicht (Größe) · Verdrahtung des Aufrufers ohne Mutation · Kommentar nennt mehr als der Code tut

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 sind MEDIUM (fehlender Fall je Element einer zugesagten Aufzählung, §3.10); behoben, wenn je ein Fall die Mutanten rot färbt. F-3 bis F-5 blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
