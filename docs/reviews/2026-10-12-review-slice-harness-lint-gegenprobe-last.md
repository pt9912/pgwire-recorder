# Review-Report: slice-harness-lint-gegenprobe-last — 2026-10-12

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3 und §6, gegen [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und gegen `AGENTS.md` §3.3, §3.6, §3.7, §3.9, §3.11 und §3.12. Die DoD prüft der Verifier.

**Gegenstand:** `2d1f0ca..b4ee633` (zwei Commits: `10dc4c8` ändert `tools/harness/lint-gegenprobe.sh`, `b4ee633` trägt §7 *Belege des Implementers* im Plan nach). `make gates` an `b4ee633` grün laut Auftrag (nicht von diesem Lauf gefahren).

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

## Summary

HIGH 0 · MEDIUM 0 · LOW 0 · INFO 3.

Der Diff setzt `parallel=6` auf `parallel=3`, ersetzt die 6 im Kopf, streicht den Satz „den Default fährt jeder Lauf von `make gates`" und liest 3 als Höchstzahl je Lauf des Skripts. Die Zweige für gesetzte und ungültige Werte (`^[1-9][0-9]*$`, Ausgang 2, feste Zeile auf stderr) stehen unverändert im Diff-Kontext. Kein Fall prüft den Standard; das ist das akzeptierte Negativ aus §6 (Weg 2, Entscheidung des Nutzers).

## Findings

### INFO-1 — §7 beschreibt Load im Maximum mit „etwas höher"

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.11, Zahlen und Aussage über den Beleg)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-lint-gegenprobe-last.md` · „Load und RAM im Maximum etwas höher"
- `befund`: Die Tabelle nennt als Load-Maximum 38,0 und 37,7 gegen 31,9 aus §1 (rund 19 % mehr); „etwas" untertreibt das. Die Zahlen selbst stimmen mit der Messung überein (Parallelität 6/3/2 wie in §1, Laufzeit 166 s und 167 s gegen 167 s, Mittel 27,8 und 31,0 gegen 27,5, RAM-Maximum 19,3 GB gegen 18,7 GB bei höherer Basis 15,6 gegen 14,5 GB); der Satz ist die einzige Abweichung, und der Plan nennt den belegten Grund (Host bei Start mit Load 22).
- `verifizierbar`: nein
- `klasse`: Wertung im Beleg enger als die Zahl daneben

### INFO-2 — Suchumfang in §7 enger als in §3

- `kategorie`: INFO
- `quelle`: Maintainability (Plan §3, `AGENTS.md` §3.9)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-lint-gegenprobe-last.md` · „`PARALLEL` über `tools/`, `harness/`, `spec/`, `docs/user`"
- `befund`: §3 verlangt die Suche über `docs/`, §7 nennt `docs/user`. Eine eigene Suche des Reviews über das ganze Repo (ohne `.git` und Baseline) findet `LINT_GEGENPROBE` nur im Skript, im abgeschlossenen `done/slice-harness-lint.md`, in dessen Review- und Verifikationsberichten, in zwei Evidence-Dateien des Registers, im eigenen Plan und in der Roadmap; keine dieser Stellen nennt 6 als gültigen Standard.
- `verifizierbar`: ja (`grep -rn LINT_GEGENPROBE`)
- `klasse`: Suchumfang im Beleg weicht vom Plan ab

### INFO-3 — Zeilenangaben im Plan nennen den Stand vor dem Diff

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.9)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-lint-gegenprobe-last.md` · „`parallel=6` (Zeile 126) wird `parallel=3`"
- `befund`: §3 und die DoD nennen Zeilen 14, 18 bis 19 und 126 des Standes vor `10dc4c8`; nach dem Diff steht `parallel=3` in Zeile 125 und der Kopf ist umgebrochen. Die Angaben beschreiben die Änderung vor Code und führen zum Anker (`parallel=`), es entsteht kein Fehlverweis. Im Skript selbst bleibt Zeile 14 des Kopfs mit 136 Zeichen deutlich länger als die Nachbarzeilen; ohne Konventions-Anker kein Finding.
- `verifizierbar`: ja
- `klasse`: Zeilenverweis im Plan altert mit dem Diff

## Gezielte Prüfungen laut Auftrag

- **Kopf des Skripts:** „Höchstzahl gleichzeitiger Läufe der Stufe je Lauf dieses Skripts, keine Obergrenze der Last von `make gates`. Nicht gesetzt: 3." Keine Zusage, ein Lauf prüfe den Standard (der Satz ist gestrichen, `GEPRÜFT DURCH` nennt nur `einstellung-<wert>`); keine Aussage zu OOM, Last, Speicher oder Zeit. Entspricht §6 *Lesart des Standards 3* und §3.11.
- **Die 6 als Standard:** nirgends mehr. `harness/sensors/lint.md`, `harness/mk/lint.mk` und `Makefile` nennen weder Zahl noch Einstellung (Suche nach `parallel`, `gleichzeitig`, `sechs`); die Spezifikation und `docs/user` ebenso nicht. Die 6 steht noch in historischen Berichten und Evidence-Dateien (Befund von damals, bleibt laut Plan §3).
- **Gesetzte und ungültige Werte:** `elif [[ ... =~ ^[1-9][0-9]*$ ]]` und der `else`-Zweig mit Ausgang 2 liegen außerhalb der geänderten Zeilen; die Fälle `einstellung-<wert>` (Zeile 333 ff.) sind nicht berührt. Verhalten unverändert.
- **§3.6:** Der Standard ist die Höchstzahl eines Werkzeug-Parameters; keine Schwelle, kein Linterprofil, keine Architekturregel ändert sich, kein Fall entfällt. Keine Lockerung im Sinn von §3.6; die Entscheidung des Architect (keine Zeile in `SPEC-049`) liegt in §6 und trägt.
- **§3.12:** Alle Randformen stehen in §6 und sind vor dem Code als vom Architect geprüft vermerkt; der Code-Commit `10dc4c8` ändert §6 nicht, `b4ee633` nur §7. Keine neue Randform im Code-Commit.
- **§3.3:** Beide Commits ohne Move; der Move `2d1f0ca` liegt vor dem Diff und ist rein.
- **§3.10 / akzeptiertes Negativ:** `parallel=3` zu `parallel=6` bleibt grün; §6 und §7 nennen das offen als Weg (2), ohne Zusage im Kopf. Kein Widerspruch zwischen Kopf und Prüfung.
- **§3.13:** Der Diff nennt keine neue Adresse (Folge-Slice, Register-*geplant*); der Ausgang `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` ist *weiter offen* und wird bei Closure angelegt, kein Nehmer-Slice betroffen.
- **Zahlen in §7:** siehe INFO-1; Lauf 1 hat keine RAM-Zahl und sagt das offen. §7 behauptet nicht, die Senkung behebe den OOM, und nennt das als „Nicht belegt".

## Negativbefund (geprüft, ohne Befund)

- `tools/harness/` (`lint-gegenprobe.sh`, Diff und Kontext): ohne Befund.
- `harness/mk/` (`lint.mk`), `Makefile`: ohne Befund.
- `harness/sensors/lint.md`, `harness/README.md` §Sensors: ohne Befund.
- `spec/`: ohne Befund (keine Zeile berührt, keine nötig).
- `docs/plan/adr/` ([ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md)): ohne Befund (nennt die Zahl nicht).
- `docs/user/`, `README.md`, `AGENTS.md`, `.claude/`: ohne Befund (keine Fundstelle).
- `docs/plan/planning/in-progress/slice-harness-lint-gegenprobe-last.md` (§7-Diff): ohne Befund über INFO hinaus.
