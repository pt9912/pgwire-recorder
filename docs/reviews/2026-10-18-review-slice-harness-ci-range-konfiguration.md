# Review-Report: slice-harness-ci-range-konfiguration — 2026-10-18

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; das ist die Frage des Verifiers).

**Gegenstand:** `slice-harness-ci-range-konfiguration`, Zweig `rev-rk` (Implementer-Stand dfeb2d23), Diff `origin/main..HEAD` (Basis 368e3d43) über `.d-check.yml`, `harness`, `tools` und den Plan; Basis der Größenmetrik eac90f75.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

> **Zitier-Form.** Dieser Report friert ein; Slices werden mit Kennung zitiert. Die ADR- und SPEC-Kennungen sind auflösende Links.

**Eingangs-Kontext:**

- `slice-harness-ci-range-konfiguration` (Plan in `in-progress/`, §7 *Belege des Implementers*: Befundzeilen, 20 Fälle, 29 Mutationen, Messungen, Läufe, Metrik)
- [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md), [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md)
- [SPEC-064](../../harness/vertraege/SPEC-064-ci.md) Kopf (Teilung), Punkte 10, 11, 12, 15, 16, 17; [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) (Namensregel, Listen)
- `AGENTS.md` (Hard Rules §3.2, §3.7 bis §3.12, §3.14)

**Gefahrene Sensoren (Auftrag: engster Sensor, keine volle Gegenprobe, kein `make gates`):** `make range-konfiguration-gegenprobe` mit `GP_MUT`/`GP_ENTDECKE` am Host als Benutzer uid 1000 (Läufe: `ok` und alle 29 Mutationen, in zwei Läufen `grün (9 Läufe)` und `grün (21 Läufe)`; die gesehenen roten Mengen stimmen mit der Tabelle in §7 überein), eigene Sonden (Muster-Weitung, Signale), `bash tools/harness/gegenproben-regel.sh` über `namensregel` (keine Verletzung), `make docs-check` (1174 Dateien, 0 Befunde, auch mit diesem Report; siehe Negativbefunde). Als root nicht gefahren (`sudo` verlangt ein Passwort); die Aussage des Implementers zu root ist hier nicht bestätigt, aber aus dem Code nicht widerlegt (siehe Negativbefunde).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Muster von `commits.id-patterns` werden nur auf Anwesenheit geprüft, nicht auf ihre Form: Setze ich in der Kopie `ADR-[0-9]{4}` auf `ADR-[0-9]`, `MR-[0-9]{3}` auf `MR-[0-9]` und `(slice\|welle)-[a-z0-9]+(-[a-z0-9]+)*` auf `(slice\|welle)-`, bleibt `GP_MUT=ok` bis auf `p11-lh` grün (gesehen: nur `LH` wurde gefärbt, weil `LH-FA-01` zwei Buchstaben trägt). Kein Fall bringt eine Message mit einer fast passenden Kennung (`ADR-1 x`, `MR-1 x`, `slice- x`), und die Mutation `m-p11-muster-verankert` verankert alle vier Muster zugleich, der Fall `p11-kennung-im-rumpf` trägt nur `MR-001`; ein Anker allein vor `ADR`, `LH` oder `(slice\|welle)` bliebe grün. Die Zusage „genau diesen vier Mustern“ (SPEC-064 Punkt 11) und der Satz im Kopf des Fragments „je Zusage und Element eine Mutation“ reichen weiter als die Prüfung. | `AGENTS.md` §3.10 (Aufzählung je Element), §3.11; [SPEC-064](../../harness/vertraege/SPEC-064-ci.md) Punkt 11 | `tools/harness/range-konfiguration-gegenprobe.sh` · `mut m-p11-muster-verankert` und die Fälle `f_p11_*`; `harness/mk/range-konfiguration.mk` · `je Zusage und Element eine Mutation` | ja — `make range-konfiguration-gegenprobe` mit den Mutanten oben (Muster-Weitung) bleibt grün | Aufzählungselement ohne eigene Mutation (Form des Musters, Anker je Muster) |
| F-2 | MEDIUM | Bricht man die Gegenprobe mit `SIGINT` oder `SIGTERM` ab, bleibt das Temp-Verzeichnis liegen und die Docker-Läufe der Hintergrund-Unterprozesse laufen weiter: Der `EXIT`-Trap löscht `$arbeit`, aber die `mut`-Unterprozesse rufen `mkdir -p "$run"` und legen den Baum neu an (gesehen: nach `SIGTERM` war das Verzeichnis kurz leer, nach 25 s stand `tmp.<zufall>/m-p11-muster-lh-fehlt/` wieder da; nach `SIGINT` blieben `fehler`, `gefahren` und `m-p11-block-fehlt`, zusammen 3,4 MB). SPEC-064 Punkt 16 sagt „räumt ihn weg“, Abbruch und Herunterfahren nennt weder §6 des Plans noch der Vertrag, und die Gegenprobe hat dazu keinen Fall. | `AGENTS.md` §3.12 (Abbruch und Herunterfahren als Randform), §3.10; [SPEC-064](../../harness/vertraege/SPEC-064-ci.md) Punkt 16 | `tools/harness/range-konfiguration-gegenprobe.sh` · `trap 'chmod -R u+rwX "${arbeit:?}" 2> /dev/null; rm -rf "${arbeit:?}"' EXIT` | ja — Skript starten, nach 8 s `SIGINT`/`SIGTERM` an den Prozess, danach `ls $TMPDIR` | Abbruch hinterlässt Temp-Baum und laufende Unterprozesse |
| F-3 | LOW | Die Mutation `m-p12-paths-index` färbt alle acht Fälle von Gruppe p12, ebenso `m-p12-immutable-when-fehlt`; der Fall `p12-index` ist damit nicht von den übrigen unterscheidbar, und die Mutation bricht zwei Einstellungen zugleich (Glob und `immutable-when`). Der Kommentar im Skript erklärt, warum (`head-allow` lehnt jede unveränderliche Datei ohne Status-Zeile ab); „nur ihn bricht“ (§3.10) ist für `p12-index` nicht erfüllt. | `AGENTS.md` §3.10 | `tools/harness/range-konfiguration-gegenprobe.sh` · `mut m-p12-paths-index` | ja — `GP_ENTDECKE=1 GP_MUT=m-p12-paths-index` zeigt acht rote Fälle | Fall ohne trennende Mutation |
| F-4 | LOW | Die Fallnamen aus `FAELLE_P11`/`FAELLE_P12` werden als `f_<name>` gerufen; ein Name ohne Funktion (Tippfehler) endet in `( bad=0; "f_…"; exit "$bad" )` mit Exit 0, der Fall zählt als grün, ohne dass etwas lief. Es gibt keine Prüfung, dass jeder Name eine Funktion hat. | Maintainability | `tools/harness/range-konfiguration-gegenprobe.sh` · `( bad=0; "f_${n//-/_}"; exit "$bad" )` | ja — einen Namen in `FAELLE_P11` umbenennen, `GP_MUT=ok` bleibt grün | Gegenprobe zählt einen nicht gelaufenen Fall als grün |
| F-5 | LOW | Der Kopf des Skripts nennt „bash ≥ 4“, „git ≥ 2.32“ und „GNU sed“; `wait -n` braucht bash 4.3, und keine Stelle prüft eine Version. Ebenso steht im Kopf „Braucht … coreutils“, die Prüfschleife nennt aber `chmod`, `env`, `cat`, `cp` nur als Werkzeugnamen ohne Versionsprüfung. | `AGENTS.md` §3.11 | `tools/harness/range-konfiguration-gegenprobe.sh` · `Braucht bash ≥ 4, make, git ≥ 2.32` | nein — kein Lauf prüft Versionen | Zusage im Kommentar weiter als die Prüfung |
| F-6 | INFO | An den Architect: Drei Fälle gehen über die Fallliste der Punkte 11 und 12 hinaus. `p11-exempt-nicht-am-anfang` (Anker `^` der Ausnahme) und `p11-struktur-allein` („`SPEC-…` zählt weder ab noch als Kennung“, Satz in Punkt 11) und `p12-index` („der Index … ist nie unveränderlich“, Satz in Punkt 12) belegen Aussagen, die der Vertrag trägt, aber nicht in seiner Fall- und Mutationsliste führt. Das ist eine Lücke der Liste, kein Überschuss: Jede der drei Aussagen ist eine Zusage (§3.10). Die Mutationen `m-p11-ausnahme-ohne-anker`, `m-p11-muster-struktur` und `m-p12-paths-index` färben je den Fall (zu `p12-index` siehe F-3). | `AGENTS.md` §3.10, §3.12 | `harness/vertraege/SPEC-064-ci.md` · `Fälle (Punkt 16): p11-ohne-kennung` | nein | Vertragslücke Fallliste |
| F-7 | INFO | An den Architect: `ignore-link-targets` lässt nach der Messung des Implementers nur den Umzug in ein anderes Verzeichnis bei gleichem Dateinamen durch (Fall `p12-link-ziel`); eine Umbenennung bleibt `core-drift-vcs`. SPEC-064 Punkt 12 nennt nur den Umzug („Link-Ziel einer umgezogenen Datei nachgezogen“) und lässt die Umbenennung offen; die Messung hier wurde nicht wiederholt, der Fall `p12-link-ziel` und die zwei Mutationen dazu laufen. Ebenso offen: Der verworfene Versuch „Merge nach Leerzeile“ (Message `x`, Leerzeile, `Merge y`) bleibt nach §7 rot und ist durch keine Einstellung zu trennen; der Vertrag nennt die Stelle der Ausnahme nur als Betreffanfang, das Verhalten für eine spätere Zeile steht nicht in der Zusage. Hier nicht nachgemessen. | `AGENTS.md` §3.12 (Stelle einer Nachricht) | `harness/vertraege/SPEC-064-ci.md` · `Link-Ziel einer umgezogenen Datei nachgezogen` | nein | Vertragslücke Randform |
| F-8 | INFO | Die Messung von CRLF und Kodierung steht in §7 als Messung („keine Zusage, kein Fall“), mit den getesteten Formen einzeln; das trägt Punkt 17. Der letzte Satz („Das Zeilenende allein ist für `vcs` also keine Änderung des Core“) verallgemeinert aus einer Messung an einer Datei und liest sich als Aussage über das Werkzeug; Punkt 17 sagt, ob `core-drift-vcs` entsteht, ist nicht zugesagt. Kein Mangel der Konfiguration, nur die Lesart des Satzes für die Closure. Nicht nachgemessen. | `AGENTS.md` §3.11; [SPEC-064](../../harness/vertraege/SPEC-064-ci.md) Punkt 17 | `docs/plan/planning/done/slice-harness-ci-range-konfiguration.md` · `Das Zeilenende allein ist für` | nein | Messung liest sich als Zusage |
| F-9 | INFO | §3 des Plans hat der Implementer geändert (Fragment-Zeile: 20 Fälle, 29 Mutationen, drei Fälle über den Vertrag hinaus). Das ist nach `AGENTS.md` §3.9 vorgeschrieben (der Plan folgt jeder Korrektur) und kein Verstoß gegen „nur §7“. §1 und §6 bleiben ohne Nachzug: §1 nennt die Zahl der Fälle nicht (kein Widerspruch), §6 nennt den Abbruch nicht als Randform (siehe F-2). Der Kopf (`Berührte Spec-Stellen`) trägt SPEC-064, SPEC-060, SPEC-050, SPEC-063; `make kopf-check` wurde nicht gefahren. | `AGENTS.md` §3.9 | `docs/plan/planning/done/slice-harness-ci-range-konfiguration.md` · `20 Fälle (12 `p11-*`, 8 `p12-*`` | ja — `make kopf-check` | Plan-Nachzug |
| F-10 | INFO | `make doc-commits` und `make doc-immutable` laufen im Worktree `/Development/pgwire-recorder-revrk` nicht (`kein lesbares git-Repository unter /repo: lstat /repo/.git/config: not a directory`), weil `.git` dort eine Datei ist und nur das Arbeitsverzeichnis eingebunden wird. Das ist eine Eigenschaft der Treiber (`d-check.mk`), nicht des Diffs; die Belege des Implementers (Range dieses Slice, Bestand) wurden deshalb hier nicht wiederholt. | Maintainability | `d-check.mk` · `docker run --rm --network none -v` | ja — `make doc-commits RANGE=368e3d43..HEAD` im Worktree | Treiber setzt Klon mit `.git`-Verzeichnis voraus |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.d-check.yml` gegen SPEC-064 Punkt 11 und 12 wörtlich (a) | geprüft, ohne Befund: vier Muster, `exempt-pattern` `^(Merge \|Revert )`, `paths` `docs/plan/adr/[0-9]*.md`, `immutable-when`, `exclude-sections: [Geschichte]`, `status-line`, `head-allow`, `ignore-link-targets: true` stimmen Zeichen für Zeichen mit dem Vertrag; `workflows` steht nicht in `modules:`; Kommentar sagt nur zu, was `range-konfiguration-gegenprobe` prüft (mit F-1 als Grenze der Aussage „je Eintrag“) |
| Mutationen am Host (b) | 29 Mutationen und die unveränderte Kopie gefahren, alle roten Mengen gleich der Tabelle in §7 und dem Skript; mehr als zehn einzeln gelesen (u. a. `m-p11-ausnahme-ohne-anker`, `m-p12-head-allow-eng`, `m-p12-link-ziel-aus`, `m-p12-status-zeile-eng`); Abweichungen nur F-1 bis F-3 |
| Gegenprobe-Mechanik: Vergleich (c) | geprüft, ohne Befund: Exit-Code, stdout und stderr werden gegen erwarteten Text verglichen; die Normierung (40-stellige SHAs, 7-stellige SHAs, Dateizahl, `d-check.mk:N`) verdeckt keine Aussage der Fälle; die Zeile `history-range-guard:` wird aus stdout gefiltert (jede Zeile dieses Präfixes, nicht nur die des Vorlaufs: unkritisch) |
| Gegenprobe-Mechanik: Rechte, Docker-Hülle, root (c) | `chmod -R a+rX` auf dem Fall-Baum und `u+rwX` im Trap vorhanden; die Hülle `docker` führt `$DOCKER` mit dem PATH des Aufrufers aus und ruft sich nicht selbst auf (gesehen: Läufe grün mit `DOCKER=docker`); als root nicht gefahren, im Code hängt nichts an Rechten außer dem Lesen durch das Image; Temp-Aufräumen bei normalem Ende in Ordnung, bei Abbruch siehe F-2 |
| Kommentare, README-Zeile, Fragment (d) | geprüft, ohne Befund außer F-1 und F-5: Die Zeile in `harness/README.md` nennt die Regel an den Schreiber („wer einen der Blöcke ändert, fährt sie“) und zeigt auf `.d-check.yml` als Quelle; kein Kommentar beschreibt eine verworfene Alternative oder abwesenden Text |
| Kette (f) | geprüft, ohne Befund: das Ziel steht per `GEGENPROBEN_CONTAINER += …` im Fragment, nicht in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen Kette; `namensregel` meldet keine Verletzung (Name endet auf `-gegenprobe`, Rezept stammt aus `harness/mk/`); `make gates-schnell-gegenprobe` hier nicht gefahren |
| Größe (§3.14) | Metrik gemessen: `git diff --numstat eac90f75..dfeb2d23 -- tools harness .d-check.yml Makefile` = 285 (19, 1, 11, 254, keine Löschungen), unter der Marke 500; entspricht §7 |
| Messung CRLF/Kodierung (h) | steht als Messung, nicht als Zusage und nicht als Fall; Lesart des Schlusssatzes siehe F-8 |
| `docs/plan/planning/in-progress/` (Plan) | geprüft, Findings F-9; `docs-check` mit diesem Report: 0 Befunde |
| `tools/harness/` außerhalb des neuen Skripts, `harness/mk/` außerhalb des neuen Fragments | nicht berührt, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Aufzählungselement ohne eigene Mutation · Abbruch hinterlässt Temp-Baum und laufende Unterprozesse · Fall ohne trennende Mutation · Gegenprobe zählt einen nicht gelaufenen Fall als grün · Zusage im Kommentar weiter als die Prüfung · Vertragslücke Fallliste · Vertragslücke Randform

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 sind MEDIUM (fehlende Negativtests bei neuem Vertrag; Randform Abbruch weder genannt noch gedeckt).

**Übergabe:** F-1 bis F-5 gehen an den Implementer; F-6 und F-7 an den Architect (Vertragslücken, keine Rückgabe des Implementers); F-8 bis F-10 sind Hinweise. Die Finding-Klassen gehen in die Slice-Closure §7. Der Report ersetzt keine Verifikation.
