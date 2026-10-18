# Review-Report: slice-harness-devcontainer-einrichten-start — 2026-10-18

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-harness-devcontainer-einrichten-start`, Diff `aad621a8..1f93ed00` (Pfade `tools`, `harness`, `docs`; Zweig `rev-start`).

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

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

- Slice-Plan `slice-harness-devcontainer-einrichten-start` im Stand des Zweigs (§1 bis §4 mit den Belegen des Implementers, §8); Vergleich mit dem Plan auf `main` (Commit `0a59a42f`, Sendung `arch-limits`)
- [ADR-0045](../plan/adr/0045-container-start-als-werkzeug-im-repo.md) (Accepted), [ADR-0047](../plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md); [SPEC-058](../../harness/vertraege/SPEC-058-container-start.md) Teil II (Punkte 9 bis 12, 13 bis 17 soweit Container, Punkt 15 (g), Punkt 2 Grund der Vorgaben)
- [LH-QA-03](../../spec/lastenheft.md#lh-qa-03--portabilität) (nur Host-Voraussetzungen)
- `AGENTS.md` (Hard Rules, §3.7, §3.9, §3.10, §3.11, §3.12, §3.13, §3.14)
- Eigene Läufe im Worktree: `bash -n` an beiden Skripten (ohne Befund); `make docs-check` (0 Befunde); `make kopf-check` (ohne Ausgabe, Exit 0); die Mutationen `p9-update`, `p9-kollision-zerstoert`, `p9-outcome`, `p9-up-exit`, `p9-up-bei-laeuft`, `p11-known-hosts-geschrieben`, `p11-meldung-stderr`, `p12-rueckbau-docker` und `p15-raeumen-name` per `GP_MUT` mit `GP_NUR=p9-laeuft` (ohne Gruppe Docker; rote Fälle wie im Mutationsprotokoll, `p15-raeumen-name` grün); Metrik nachgerechnet (`git diff --numstat aad621a8..HEAD -- tools harness docs/developer`: 515 und 127, Summe 642; davon `harness/vertraege/` 29, ohne sie 613); Maschinen-Statik: alle `mut`-Zeilen ohne roten Fall gezählt (sechs). Nicht gefahren: volle Gegenprobe, `make gates`, Gruppe Docker (Auftrag des Nutzers: engster Sensor).

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Commit `1f93ed00` hat §5 (Closure-Trigger), §6 (Risiken mit Ausgängen, Randformen A4 bis A25, Entscheidungen Z1 bis Z4), die Überschrift von §7, die Schätzung und die vorab benannten Rückführungs-Bedingungen aus §4 entfernt und die Belege des Implementers (Mutationsprotokoll, Rückgaben, Dokument-Sätze) unter §4 gesetzt. §1 und die DoD verweisen weiter auf §6 („Z2, §6“, „je Randform aus §6“), die Risiken der Closure haben keinen Ort mehr, und der Plan der Sendung `arch-limits` (`0a59a42f`, §1, DoD, §3, §6) trifft auf entfernten Text. | `AGENTS.md` §3.9; `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst und §Lifecycle als State Machine (Rückführungen: Bedingung vorab in §4) | `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten-start.md` · `## 4. Trigger` (danach folgt unmittelbar `## 8.`; es gibt kein `## 5.`, `## 6.`, `## 7.`) | ja — `grep -n '^## '` auf dem Plan; `git diff 8b4564e6..HEAD` am Plan | Plan folgt dem Diff nicht (Abschnitte beim Umbau entfernt) |
| F-2 | MEDIUM | Die Belege, auf die der Plan sich beruft, stehen nicht im Plan: die Messung des Labels `devcontainer.local_folder` (Pflicht vor dem Code, §3), die Laufzeit der Gegenprobe (DoD), der Lauf der Gruppe Docker gegen einen echten Container (alt §5), die Läufe (nur der Name `full4` fällt) und die „Messung `docker update` oben“ sowie „gemessen siehe oben“ (Rückgabe an den Planner 2) verweisen auf Text, den es nicht gibt; die Commit-Message nennt „Messung des Labels, Läufe“, der Diff enthält sie nicht. Das Dokument sagt dazu als Tatsache, ein laufender Container nehme `docker update --memory <x> --memory-swap <x>` an (beide Werte gleich, „gilt bis zum Neubau“). | `AGENTS.md` §3.11; Plan DoD (Dokument: „Jede Aussage ist durch einen Lauf belegt oder steht in §13“) | `docs/developer/devcontainer.md` · „Ein laufender Container nimmt neue Werte an mit `docker update --memory <x> --memory-swap <x> <Name>`“; Plan · „Messung `docker update` oben“ | ja — der Verifier fragt die Läufe ab; `grep -n -i 'Laufzeit\|devcontainer.local_folder' ` am Plan findet nur Aufträge | Beleg im Plan fehlt, Verweis ins Leere |
| F-3 | MEDIUM | `p15-volume-rest-entfernen` wird als äquivalent geführt („nach `docker rm -fv` bleibt kein anonymes Volume“), die Sonde `p15-docker-ende` durchläuft den Zweig nie. Der Mutant ist nicht äquivalent: In dem Fall, in dem der Zweig feuert (ein Volume bleibt), meldet die Gegenprobe es, entfernt es aber nicht mehr, und Punkt 15 (d) sagt „per Name entfernt“. `p15-rm-ohne-v` lässt den Zweig feuern, prüft aber nur die Meldung, nicht das Entfernen. | `AGENTS.md` §3.10; Skill (MEDIUM: Äquivalenter Mutant ohne Invariante, Ort oder Sonde) | `tools/harness/container-start-gegenprobe.sh` · `mut p15-volume-rest-entfernen p 'rest=1; docker volume rm "$v" > /dev/null 2>&1' 'rest=1' 'p15-docker-.*'` | ja — Lauf des Mutanten bleibt grün; kein Fall, der das Entfernen im Zweig beobachtet | Äquivalenz-Behauptung ohne tragende Sonde |
| F-4 | MEDIUM | Die Randformen 1 bis 6 der Rückgaben (zwei Exit-4-Gründe zugleich, Zustände `paused`/`restarting`/`removing`/`dead` wie „besteht, läuft nicht“, fehlender Port 2222 als abweichend, laufender Container bei fehlendem Gerüst, Signal gegen Fehler eines Schrittes, `docker ps`/`inspect` mit Exit ungleich 0 als „kein Container“) sind im Code entschieden (`zust_container`, `s_override`, `s_start`) und erscheinen erst im Commit `1f93ed00` nach dem Code in §4, nicht in §6 vor dem ersten Code-Commit. Konkret bei Rückgabe 4: ein laufender Container zu einem fehlenden `$PGR_DEV` legt ein neues Schlüsselpaar an, überspringt `override` und `start`, und die Probe scheitert nach `PGR_WARTEN` Sekunden mit Exit 5. Bei Rückgabe 6 fährt `start` weiter, wenn `docker ps` scheitert. | `AGENTS.md` §3.12 (Randform vor dem Code; Implementer hält an und gibt zurück); Skill (MEDIUM: Randform im Code-Commit) | `tools/harness/container-start.sh` · `zust_container() {` und `[ "$zc" != laeuft ] || return 0`; Plan · „Rückgaben an den Architect (Randform · Frage)“ | nein — Urteil des Architect; die Rückgaben sind an den Architect gerichtet (kein Konflikt-Pfad nötig, solange er sie beantwortet) | Randform erst nach dem Code genannt |
| F-5 | LOW | Der Filter `name=^/$PGR_NAME$` behandelt `.` aus `PGR_NAME` (laut Punkt 2 zulässig: `[a-z0-9][a-z0-9_.-]*`) als Regex-Zeichen. Bei `PGR_NAME=a.b` meldet das Ziel einen fremden Container `axb` ohne Label als Kollision mit Exit 4. Kein Fall und keine Mutation hat einen Namen mit Punkt. | `AGENTS.md` §3.12 (Randform „unbekannter Wert / Alias“), §3.10 | `tools/harness/container-start.sh` · `docker ps -a --filter "name=^/$PGR_NAME\$"` | ja — Fall mit `PGR_NAME=a.b` und einer Attrappe mit Container `axb` | Metazeichen einer Eingabe im Filter |
| F-6 | LOW | Die Zeile „Der Docker-Daemon des Hosts stürzte ab“ kennzeichnet nur die Aussage zu `docker start` und `--restart` als nicht gemessen. Ungekennzeichnet bleiben „die Dev-Container sind gestoppt (Exit 0)“ und „nur `postStartCommand` stellt die Netzbeschränkung und den Daemon im Container wieder her“; gemessen sind laut Plan nur Klon unverändert und `git fetch`. Dass `make container-start` nach einem Absturz fortsetzt, ist als Lauf nicht belegt (belegt ist `docker stop` und Aufruf). | `AGENTS.md` §3.11; Plan · „die Aussage zu `docker start` und `--restart` ist als Folgerung gekennzeichnet, nicht gemessen“ | `docs/developer/devcontainer.md` · „nur `postStartCommand` stellt die Netzbeschränkung (`egress-init.sh`) und den Daemon im Container wieder her“ | nein — Vergleich mit der Messung des Koordinators (nicht im Plan, siehe F-2) | Messgrenze im Dokument zu eng gezogen |
| F-7 | LOW | Die Zusage „kein `docker update`“ steht ohne Einschränkung im Dokument („ändert sie nie“), in `harness/README.md` („ändert die Limits eines vorhandenen Containers nie“) und im Vertrag; die Mutation `p9-update` ruft es nur im Zweig „läuft“ auf. Ein Aufruf in den Zweigen `gestoppt` oder `keiner` bleibt für die Gegenprobe unsichtbar (der Vertrag verlangt nur den laufenden Fall). | `AGENTS.md` §3.10 (Aufzählung und Ortsangabe zählen je Element), §3.11 | `tools/harness/container-start-gegenprobe.sh` · `mut p9-update s '[ "${i%% *}" != running ] || zc=laeuft'` | ja — Mutant mit `docker update` im Zweig `gestoppt` bleibt grün | Quantor „nie“ trägt eine Mutation |
| F-8 | LOW | Die Gegenprobe meldet „je Zusage eine Mutation“ am Ende und im Kopf „Punkte 1 bis 16“; sechs `mut`-Zeilen haben keinen roten Fall (`p2-vorgabe-warten`, `p15-int`, `p15-term`, `p15-schutz`, `p15-raeumen-name`, `p15-volume-rest-entfernen`). §4 des Plans führt fünf davon; `p2-vorgabe-warten` (Mutant der Vorgabe 60, im Nehmer durch `p11-vorgabe-warten` abgedeckt) fehlt in der Liste der grünen Mutanten. Die Abschlusszeile sagt mehr, als geprüft ist. | `AGENTS.md` §3.11 | `tools/harness/container-start-gegenprobe.sh` · `echo "container-start-gegenprobe: gruen — SPEC-058: Fälle der Punkte 1 bis 16 und je Zusage eine Mutation"` | ja — Zählung der `mut`-Zeilen ohne roten Fall | Abschlusszeile weiter als die Prüfung |
| F-9 | INFO | `p11-opt-benutzer-adresse` deckt zwei Elemente der Aufzählung (`vscode`, `127.0.0.1`) mit einem Fall; die Mutationen `p11-benutzer` und `p11-adresse` färben ihn je einzeln rot. Die übrigen neun Elemente haben je einen Fall und eine Mutation (nur Entfernen oder ein Wert; kein Mutant ändert den Wert von `StrictHostKeyChecking`, `BatchMode` oder `IdentitiesOnly`). | `AGENTS.md` §3.10 | `tools/harness/container-start-gegenprobe.sh` · `benutzer-adresse) t=' vscode@127.0.0.1 ' ;;` | nein — Hinweis | Fall deckt mehrere Elemente |
| F-10 | INFO | Das Label `devcontainer.local_folder` wird mit `$PGR_DEV` wörtlich verglichen. Der Pfad darf ein Symlink-Segment haben (Punkt 2 prüft nur die Form); ob die CLI das Label mit aufgelöstem Pfad setzt, ist weder gemessen noch in einer Randform genannt. Bei Abweichung wäre ein vorhandener Container „ohne Label“ und der Lauf endete mit Exit 4 (Kollision). | `AGENTS.md` §3.12 | `tools/harness/container-start.sh` · `lab=$(docker ps -a --filter "label=devcontainer.local_folder=$PGR_DEV"` | nein — Messung nötig (Zuständigkeit: Architect) | Randform Pfad-Auflösung nicht genannt |
| F-11 | INFO | Der Zweig trägt den Plan ohne die Sendung `arch-limits` aus `0a59a42f` (§1, DoD, §3, §6, §8: `PGR_MEM` `8g`, `p9-kein-update`, Abschnitt „Limits“); der Diff liefert sie. Beim Zusammenführen kollidiert `0a59a42f` mit dem Umbau aus F-1. §3 der Tabelle nennt weiter `pgr-dev2` und Port 2223 (Rückgabe an den Planner 1 des Implementers, nicht Gegenstand des Reviews). | `AGENTS.md` §3.9 | Plan · Kopf von `## 1. Ziel und Abgrenzung` (ohne Limits-Absatz) | ja — `git merge-base --is-ancestor 0a59a42f rev-start` (Exit 1) | Plan-Commit fehlt auf dem Zweig |
| F-12 | INFO | Metrik und Marke stimmen: 515 Einfügungen und 127 Löschungen, Summe 642 wörtlich, 29 in `harness/vertraege/`, ohne sie 613; Marke 900 und Meldeschwelle 700 nicht erreicht. Die drei Rückgaben an den Planner (Pfade `pgr-dev2`/2223, Schätzung in der Metrik, Zeilen der Zweige `arch-limits` und `arch-spec058-punkt15`) gehören dem Planner und werden hier nicht gewertet. | `AGENTS.md` §3.14 | Plan · „Summe 642“ | ja — `git diff --numstat aad621a8..HEAD -- tools harness docs/developer` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/container-start.sh` (Schritte `zustand`, `start`, `warten`, `kind`, Falle für `SIGINT` und `SIGTERM`) | geprüft, Befunde nur F-4, F-5, F-10; `bash -n` ohne Befund; kein Wert der drei Zugangsvariablen gelangt in eine Ausgabe (die Variablen werden am Anfang entfernt, `p9-umgebung` mit drei Sentinels, Mutationen je Variable `p8-umgebung-<Variable>` im Plan; die Ausgabe von `devcontainer up` und `ssh` wird nicht umgeformt und enthält keinen Wert) |
| `tools/harness/container-start-gegenprobe.sh` (Fälle `p9-*`, `p10-*`, `p11-*`, `p12-*`, `p13-*`, `p14-*`, Aufzählungen: zehn Optionen der Probe in elf Elementen mit Mutation, Signale `INT`/`TERM` je Schritt `geruest`, `start`, `warten`, Zustände `running`/`exited`/`created`, Port und Adresse je einzeln) | geprüft, Befunde nur F-3, F-7, F-8, F-9; neun Mutationen per `GP_MUT` gefahren, rote Fälle stimmen mit dem Mutationsprotokoll |
| Grüne Mutanten `p15-int`, `p15-term`, `p15-schutz`, `p15-raeumen-name` | geprüft, ohne Befund (Invariante, Ort und Sonde tragen: Standardverhalten von bash bei Signal mit `EXIT`-Falle, `mktemp -t pgr-cs-`, Namensschutz `p15-name` rot, Sonde `p15-name-ohne-praefix`); `p15-raeumen-name` im Lauf grün bestätigt |
| `harness/README.md` (Zeilen `container-start` und `container-start-gegenprobe`) | geprüft, ohne Befund außer F-7 (Quantor „nie“); Aussagen zur Gruppe Docker stimmen mit dem Skript überein |
| `harness/mk/container-start.mk`, `harness/mk/container-start-gegenprobe.mk` | geprüft, ohne Befund (Rezept unverändert, nur Kopf und Hilfetext; `p14-rezept` und `p14-phony`) |
| `harness/vertraege/SPEC-058-container-start.md` | geprüft, ohne Befund (Änderungen sind die des Architect; Punkt 2 „Grund der Vorgaben“ und Punkt 15 (g) stimmen mit Skript und Gegenprobe überein; `12g` kommt in `tools/` und `harness/` nur als Mutant `p2-vorgabe-speicher` und als Messung vor) |
| `docs/developer/devcontainer.md` | geprüft, Befunde F-2 und F-6; kein Name und keine Mail des Nutzers (der Treffer `pt9912` ist die URL des Werkzeugs `u-boot`, unverändert gegenüber der Basis), keine Zugangswerte; Abschnitte §8.1 „Zwei Arten zu arbeiten“ und §10.1 „Limits“ nennen ihre Messwerte (12g gegen 8g: Spitzen 5,42 und 4,24 GiB, 1701 s in §8) und die Grenze „nicht gemessen“ (auch in §13); die Stellen mit `12g` stehen als Messung gekennzeichnet |
| `make docs-check`, `make kopf-check` | geprüft, ohne Befund (0 Befunde; Exit 0) |
| `AGENTS.md` §3.7 / Kommentare in `container-start.sh` | geprüft, ohne Befund (jeder Kommentar nennt eine Zusage mit Fall oder eine Kopplung; kein Konjunktiv über verworfene Alternativen) |
| Rückgaben des Implementers (sechs an den Architect, drei an den Planner) | eingeordnet: die sechs an den Architect gehören ihm (F-4); die drei an den Planner sind nicht Gegenstand des Reviews (F-12) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Plan folgt dem Diff nicht (Abschnitte beim Umbau entfernt) · Beleg im Plan fehlt, Verweis ins Leere · Äquivalenz-Behauptung ohne tragende Sonde · Randform erst nach dem Code genannt · Metazeichen einer Eingabe im Filter · Messgrenze im Dokument zu eng gezogen · Quantor „nie“ trägt eine Mutation · Abschlusszeile weiter als die Prüfung

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-4 (MEDIUM). F-1 und F-2 sind Plan und Belege (kein Produktcode); F-3 ist ein Mutant, dessen Äquivalenz nicht trägt; F-4 wartet auf die Antwort des Architect auf die Rückgaben 1 bis 6. Ein HIGH liegt nicht vor: kein Wert eines Zugangs erscheint in einer Ausgabe, die Gegenprobe räumt nur Selbstangelegtes (Punkt 15 (g)), und die geprüften Mutationen färben die Fälle ihrer Gruppe.

**Übergabe:** Findings gehen an den Implementer (F-1 und F-2 teilweise Rückkante Review → Plan, mit Merge-Reihenfolge zu `0a59a42f`); F-4 und F-10 zusätzlich an den Architect; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
