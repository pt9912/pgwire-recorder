# Review slice-harness-vertraege-umzug-gate-kette

**Rolle:** Reviewer. **Datum:** 2026-10-13.
**Gegenstand:** Diff `8b1ede6..7249e6e` (Move-Commit bis HEAD) gegen den Plan `docs/plan/planning/in-progress/slice-harness-vertraege-umzug-gate-kette.md`, [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), [MR-002](../../harness/conventions.md#mr-002), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md), `AGENTS.md` §3.3, §3.7 bis §3.14.
**Summary:** 0 HIGH · 0 MEDIUM · 3 LOW · 3 INFO.

## Findings

### F-1 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.11 (Kommentar sagt nur zu, was ein Test prüft; Kommentar beschreibt, was da ist)
- `pfad`: `tools/harness/stop-hook-gegenprobe.sh` · "Art X, Hook und Spezifikation nicht inert" (Kopfkommentar) und "Hook, Skripte der Liste und Spezifikation sind nicht inert." (über `f_p12_selbst_hook`)
- `befund`: Der Diff ergänzt den Fall `p12-selbst-vertrag` und die Mutation, zieht aber die beiden Kommentare nicht nach, die aufzählen, welche Selbständerungen die Gegenprobe prüft. Beide nennen den Vertrag nicht, obwohl Punkt 12 ihn nun nennt und ein Fall ihn prüft; der Kommentar ist enger als die Prüfung, kein Zusage-Überschuss.
- `verifizierbar`: nein
- `klasse`: Aufzählung im Kommentar nach neuem Fall nicht nachgezogen

### F-2 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.9 (Der Slice-Plan folgt jeder Korrektur)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-vertraege-umzug-gate-kette.md` · "weil der Vertrag 261 statt 262 Zeilen trägt"
- `befund`: §7 nennt den Vertrag mit 261 Zeilen (zweimal), die Datei trägt 262 (`wc -l`; §1 nennt 262 richtig). Die Vergleichsbefehle in §7 ("Zeilen 3216 bis 3476") und in der DoD/§3 ("3216 bis 3477") weichen voneinander ab; mit 3477 meldet der Vergleich zusätzlich `262d262` (Leerzeile).
- `verifizierbar`: ja (`wc -l harness/vertraege/SPEC-052-gate-kette.md`)
- `klasse`: Zahl im Plan nach Messung nicht nachgezogen

### F-3 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Nachzählen beim Eintragen), §3.9
- `pfad`: `docs/plan/planning/open/slice-harness-lh-links-pflicht.md` · "zwei Schichten (Vertrag: neue Stelle"
- `befund`: Das Etikett der Schicht wechselt in einer von drei Zählungen von "Spezifikation" auf "Vertrag"; die jüngste Zählung ("zwei Schichten (Spezifikation; Harness: `.d-check.yml`, Commands, Links in `harness/`)") behält "Spezifikation". Nach §3 des Gebers liegt die Stelle jetzt in der Schicht Harness; der Nehmer zählt also uneinheitlich, ohne dass §8 die Zählung neu nennt. Der Implementer hat das als Hinweis an den Planner gemeldet (§7); die Zahl bleibt unter der Grenze, die Entscheidung liegt beim Planner.
- `verifizierbar`: nein
- `klasse`: Ortsänderung im Nehmer ohne einheitliche Nachzählung in §8

### F-4 INFO

- `quelle`: Maintainability · `pfad`: `harness/vertraege/SPEC-052-gate-kette.md` · "(jede geänderte, hinzugekommene oder gelöschte Datei auf der Positivliste und grüner Stand der schnellen Kette)"
- `befund`: Punkt 22 (b) ist für sich allein eindeutig (Freigabe in beiden Aufrufarten). Die Klammer fasst "freigibt nach Punkt 9" als "Datei auf der Positivliste" und lässt Symlinks (Punkt 15) und die Art `X` (Punkt 12) unerwähnt; der Verweis "nach Punkt 9" trägt sie. Wortlaut vom Architect, Entscheidung des Nutzers; keine Aktion für diesen Diff.
- `klasse`: Klammer enger als verwiesener Punkt (benannt)

### F-5 INFO

- `quelle`: `AGENTS.md` §3.5, [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) (akzeptiertes Negativ) · `pfad`: `docs/plan/adr/0038-gate-kette-voll-und-schnell.md` · "§11 Harness-Werkzeuge der Spezifikation"
- `befund`: Die angenommene [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) nennt weiter §11 der Spezifikation als Ort von `SPEC-052`; sie ist unveränderlich und der Fließtext ist kein Link. Ebenso `MR-002` (Zustandsabsatz "Der Umzug geschieht erst im Umzug-Slice ..."): Closure-Arbeit nach §5 des Plans, vor `git mv` nach `done/` zu klären.
- `klasse`: Fließtext in Entscheidung bleibt (akzeptiert) / Zustand von MR-002 offen bis Closure

### F-6 INFO

- `quelle`: `AGENTS.md` §3.10 · `pfad`: `tools/harness/stop-hook-gegenprobe.sh` · "aufnahme 'spec/*' p12-selbst-spec"
- `befund`: Punkt 12 nennt `spec/spezifikation.md` nicht mehr; `p12-selbst-spec` bleibt als Fall der Pfadklasse `spec/` (Entscheidung des Nutzers, Randform 2) und prüft damit nur noch die allgemeine Regel "nicht auf der Positivliste". Kein Verstoß.
- `klasse`: Fall ohne namentliche Zusage (entschieden)

## Eigene Prüfung (nicht auf Zusage)

- Wortgleichheit `SPEC-052`: `git show 231089e:spec/spezifikation.md | sed -n '3216,3477p' | sed '1s/^### /# /' | diff - harness/vertraege/SPEC-052-gate-kette.md`: genau zwei Hunks, `120,122c120,122` (Punkt 12) und `247,249c247,250` (Punkt 22 (b)), dazu `262d262` (Leerzeile am Ende); sonst leer, `#` statt `###` nur in der Überschrift. Mehr Abweichung gibt es nicht.
- Punkt 12: beide Stellen im Wortlaut der Randform 2 ("und dieser Vertrag nennt sie"; "oder diesem Vertrag (`harness/vertraege/SPEC-052-gate-kette.md`) ist nicht inert"). Die Aufzählung (Hook, `dateiliste.sh`, `working-tree-hash.sh`, Vertrag) hat je Element einen Fall (`p12-selbst-hook`, `-dateiliste`, `-hash`, `-vertrag`). Der Vertragstext enthält kein "Spezifikation" mehr (`grep` leer).
- Punkt 22 (b): Wortlaut wie von Randform 1 vorgegeben; der Satz "beide Aufrufarten gelten auch ..." ist entfallen; Rest von Punkt 22 unverändert (Diff).
- Historie: Die fünf Zeilen 3576 bis 3580 am Stand `231089e` stehen wortgleich am Ende der Tabelle in `harness/vertraege/README.md`, nach der Zeile vom 2026-10-08, zwei mit Datum 2026-10-14. Die Zeile von `SPEC-052` steht in der Tabelle der Verträge.
- Fall `p12-selbst-vertrag`: in `FAELLE`, `DATEIEN_FIXTURE` (Pfad ergänzt), `f_p12_selbst_vertrag` (`pfad_fall ... block harness/vertraege/SPEC-052-gate-kette.md`) und `aufnahme 'harness/vertraege/*' p12-selbst-vertrag`. Fallzahl: `FAELLE` zählt 139 Namen mit `p`-Nummer, wie §7 sagt; in Doku, `harness/README.md`, `harness/sensors/` und Kommentaren steht keine Fallzahl, die nachzuziehen wäre.
- Mutation 2 nachgefahren: `harness/vertraege/*) return 0 ;;` in die Positivliste von `.claude/hooks/stop-require-gates.sh`, `touch`, `make stop-hook-gegenprobe`: `ROT — Kopie 'ok' ist nicht grün: p12-selbst-vertrag`, nur dieser Fall. Danach `git checkout` des Hooks, Baum sauber.
- Spezifikation: `## 11. Historie` ist die letzte Überschrift; `grep -n "harness/\|vertraege\|SPEC-052" spec/*.md` ohne Treffer; kein Link auf `#12-historie` im Repo.
- Verweise: Repo-weiter `grep` (ohne `.git`, `.harness`, `done/`, `reviews/`, `adr/`) auf `spezifikation.md` neben `§11`, `SPEC-052`, `spec-052`: Treffer nur im Plan selbst (Zitate), `MR-002` (F-5) und `slice-harness-lint-warten-ohne-frist` (dort ist §11 die umnummerierte *Historie*, richtig). Kein `§12` mehr in `open/`, `next/`, `.claude`, `AGENTS.md`, `harness/`. Kein Link `spezifikation.md#spec-052` außerhalb des Plans.
- Kommentare in `tools/`, `.claude/`, `harness/mk/`: alle nennen `harness/vertraege/SPEC-052-gate-kette.md`; die übrigen Nennungen von `SPEC-052 Punkt N` sind Kennung ohne Pfad. Außer F-1 kein veralteter Ort.
- `AGENTS.md` §3.12: nennt `harness/vertraege/` ohne Ausnahme und ohne `§11` (Herkunfts-Anker in der Form von §3.7). Sendung aus dem Ordner-Slice damit getragen.
- Vier Pläne in `open/`/`next/`: Ort `harness/vertraege/`, nächste Nummer 053, Fließtext "Abschnitt für Harness-Werkzeuge der Spezifikation" berichtigt; Nehmer liegen in `open/`/`next/`, keiner schließt die Sendung aus (§3.13). Zählung der Schichten: siehe F-3; die anderen drei Pläne zählen keine Schichten neu und gewinnen keinen Liefer-Punkt.
- §3.3: Der Move `8b1ede6` ist rein; die drei Folge-Commits enthalten keinen Rename oder Löschung einer Datei (`--name-status` ohne `R`/`D`).
- §3.5: Kein ADR im Diff geändert.
- Größe (§3.14, Pfade des Plans): `git diff --numstat 8b1ede6..HEAD -- spec harness AGENTS.md tools .claude docs/plan docs/reviews` Summe 677 (+355, −322); bis `ebccbeb` 636; Marke 900 und Meldeschwelle 800 nicht erreicht; Schätzung des Architect 680 bis 780 mit Fall, 620 bis 700 ohne.
- Zusagen ohne Gate (§3.11): "Spezifikation ohne Verweis auf `harness/`" (Links fängt `matrix-forbidden`, Prosa nur `grep`), Kommentar-Pfade (nur `grep`), `AGENTS.md` §3.12 (nur `grep`); der Plan nennt das in §7 selbst und sagt keine Sperre zu. Kein Verstoß.
- §3.9: Kopf, §1, §3, §6, §7 folgen dem Diff bis auf F-2; `Berührte Spec-Stellen` nennt `SPEC-049` als nicht berührt, wie im Diff.
- §3.12: Keine Randform im Code-Commit neu in §6; Randform 2 und der Fall stehen vor dem ersten Code-Commit im Plan (Entscheidung des Nutzers).
- `gofmt`/Lint: kein Go-Code im Diff, keine `//nolint`-Direktive.

## Geprüft, ohne Befund

- `harness/vertraege/` (README und `SPEC-052-gate-kette.md`): ohne Befund (F-4 nur Wortlaut des Architect).
- `spec/`: ohne Befund.
- `harness/README.md`, `harness/mk/`: ohne Befund.
- `.claude/hooks/stop-require-gates.sh`: ohne Befund (nur Kommentarzeile).
- `tools/harness/` (`dateiliste.sh`, `record-gates.sh`, `gates-schnell.sh`, `gates-schnell-gegenprobe.sh`): ohne Befund; `stop-hook-gegenprobe.sh`: F-1, F-6.
- `AGENTS.md`: ohne Befund.
- `docs/plan/planning/open/`, `next/`: F-3; sonst ohne Befund.
- `docs/plan/adr/`, `docs/reviews/`: unverändert, ohne Befund (F-5).
