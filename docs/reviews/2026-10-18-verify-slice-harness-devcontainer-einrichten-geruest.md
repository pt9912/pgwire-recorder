# Verifikation slice-harness-devcontainer-einrichten-geruest

**Rolle:** Verifier. **Stand:** Zweig `rev-geruest`, Commit `4f72a7a5`, Basis `8b19dc66`. **Datum:** 2026-10-18.
**Gegenstand:** DoD (§2) und Belege (§7) des Slice-Plans `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten-geruest.md` gegen [`SPEC-058`](../../harness/vertraege/SPEC-058-container-start.md) Teil I und [ADR-0045](../plan/adr/0045-container-start-als-werkzeug-im-repo.md) (Accepted).

## Urteil

DoD-Liefer-Punkte erfüllt und belegt, mit **einem eigenen roten Lauf am Host** (V-1) und vier kleineren Befunden. Closure-Pflichten (Notiz, Register, Risiko-Ausgänge, Paarungen) sind erwartungsgemäß noch offen (der Planner schreibt sie).

## Eigene Läufe (Host, Docker-Daemon des Hosts, `pgr-dev2` unberührt)

| Lauf | Ergebnis |
|---|---|
| `make container-start-gegenprobe` (volle Gegenprobe, 153 Mutationszeilen) | **rot**, Exit 2 von `make`, 3 min 06 s: eine Mutation, `p7-doctor` (V-1); alle übrigen Fälle und Mutationen wie erwartet |
| `GP_MUT='^$'` (nur unveränderte Kopie) | grün, 11 s; hinterlässt kein `pgr-cs-*` in `/tmp` |
| 14 Mutanten einzeln (`GP_ENTDECKE=1`): `p1-repo`, `p3-u-boot-lexikalisch`, `p5-praefix`, `p6-passphrase-anlegen`, `p8-behalte-oauth`, `p13-exit-null`, `p14-plus`, `p15-ende`, `p15-exit-trap`, `p16-scharf` und die vier grünen | rot genau wie in der Tabelle von §7 (z. B. `p14-plus` färbt `p14-n`, `-q`, `-t`, `-n-i`, `-rezept`, vier `makeflags`-Fälle; `p8-behalte-oauth` nur `p8-umgebung-CLAUDE_CODE_OAUTH_TOKEN`) |
| `make gates-schnell` | Exit 0 (1 min 60 s) |
| `make docs-check` | Exit 0, 1056 Dateien, 0 Befunde |
| `make kopf-check` | Exit 0 |
| `make gates` | nicht gefahren (Auftrag); Beleg des Implementers: Exit 0, 1630 s auf `b311d82a` (Log nicht gesehen; Endstand `4f72a7a5` ändert nur den Plan) |

## DoD-Punkte

- **Liefer-Punkt 1, Ziel `container-start` Teil I** (`harness/mk/container-start.mk`, `tools/harness/container-start.sh`): erfüllt. Schritte `eingaben`, `werkzeuge`, `zustand`, `geruest`, `schluessel`, `override` in der Reihenfolge von Punkt 4; Exit 2 bis 5; Zwischenstand Exit 6 mit der Zeile aus Punkt 13; Rezept eine Zeile ohne `$(MAKE)` und `+`, `.PHONY` ohne Voraussetzung; Zugangsvariablen zu Beginn entfernt; in `harness/README.md` §Werkzeuge genannt. Code gegen Vertrag gelesen: keine Abweichung gefunden.
- **Liefer-Punkt 2, Gegenprobe**: Fälle und Mutationen je Punkt vorhanden und am Host bis auf V-1 wie behauptet; Fall `p16-gate-kette` und `p7-echt-u-boot` (echtes `u-boot` 0.10.0) grün; hängt an `GATE_CHECKS`, steht in `GATE_SCHNELL_AUSNAHME`; Index in `harness/README.md`; Fixture von `gates-schnell-gegenprobe` nachgezogen (eine Zeile).
- **Gates grün**: Beleg des Implementers vorhanden (§7), Teilnachweis durch meine Läufe (V-1 ausgenommen).
- **Review**: Report liegt (untracked im Worktree: `2026-10-18-review-slice-harness-devcontainer-einrichten-geruest.md`).
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen**: offen (Planner). §7 enthält weder „keine Beobachtung“ noch Register-Zeilen; die Ausgänge der §6-Risiken stehen als Vorschlag im Plan, nicht als Ergebnis.

## Befunde

**V-1 (hoch, DoD „Gegenprobe grün“ am Host nicht reproduziert).** Mutation `p7-doctor` (`ub generate devcontainer; ub doctor`) färbt am Host zusätzlich `p7-echt-u-boot` (erwartet: nur `p7-kein-up-doctor`, `p7-reihenfolge`). Reproduziert dreimal (Vollauf, zweimal einzeln). Ursache gemessen: Im Pfad der Gegenprobe (`$w/bin:$sys`, ohne `git`, `docker` als Attrappe) meldet das echte `u-boot doctor` `git.installed`, `docker.installed`, `docker.compose.installed`, `docker.reachable` als Fehler und endet mit Exit 11; `ub` meldet Exit 5, `gruen` (Exit 6 erwartet) wird rot. Ob die Mutationstabelle im Container des Implementers grün war, weil dort `doctor` anders endete, ist nicht belegt; der Beleg „Exit 0“ gilt für dessen Umgebung. Folge: `make gates` am Host mit `u-boot` 0.10.0 und `/usr/bin/git` ist mit großer Wahrscheinlichkeit rot an diesem Gate. An den Implementer: Erwartung für `p7-doctor` umgebungsunabhängig machen (der Fall `p7-echt-u-boot` darf nicht vom Exit von `doctor` abhängen, oder die Erwartung nennt ihn) und in beiden Umgebungen belegen. (`p7-up` bleibt am Host grün, weil `u-boot up` ohne Dienste Exit 0 liefert; kein Befund.)

**V-2 (mittel, Risiko-Ausgang „entfallen mit Beleg: nach Ende kein Temp-Verzeichnis“).** Die unveränderte Kopie räumt auf (gemessen). Die Mutanten `p15-ende` (`kill -KILL`) und `p15-exit-trap` hinterlassen absichtlich leere `/tmp/pgr-cs-*`-Verzeichnisse: zwei nach `p15-ende` allein, zwölf nach zwei Vollläufen; jeder `make gates` hinterlässt rund sechs. Der Beleg des Implementers sagt das nicht; die Zusage „Aufräumen bei Ende“ (Punkt 15 (f)) gilt für die Gegenprobe selbst, nicht für ihre Mutanten, und `SPEC-058` Punkt 15 nimmt das nicht aus. Ausgang des Risikos „hinterlässt Verzeichnisse“ trägt damit nur mit dieser Einschränkung; Architect/Planner entscheiden (Mutant räumt selbst auf, oder akzeptiertes Negativ im Vertrag). Ich habe meine Reste per `rmdir` nach Namen entfernt.

**V-3 (niedrig, Adresse nimmt nicht an, `AGENTS.md` §3.13).** Grüner Mutant `p2-vorgabe-warten` (Vorgabe 61 statt 60) ist in Teil I unbeobachtbar; `PGR_WARTEN` wird nur geprüft. Der Implementer nennt als Adresse `slice-harness-devcontainer-einrichten-start`. Dessen §1 nennt die Vorgabe 60 nicht (kein Treffer für „Vorgabe“/„60“ außer in A8 bzw. der Sendung ohne Vorgabe), `grep` nach der Kennung dieses Slice findet die Bindung, nicht die Vorgabe. Der Punkt ist als „an den Planner“ gemeldet; er ist noch nicht eingetragen. Die Sonde `grep -n PGR_WARTEN tools/harness/container-start.sh` (Zeilen 6, 36, 59, 60, 63) stimmt: einziger Leser ist `pruefe_PGR_WARTEN`.

**V-4 (niedrig, grüne Mutanten `p15-int`, `p15-term`, `p15-schutz`).** `p15-int` und `p15-term`: die Fälle `p15-aufraeumen-int`/`-term` bleiben grün, weil bash bei Tod durch das Signal die `EXIT`-Falle ausführt und 130/143 liefert (Invariante gemessen: leere Menge roter Fälle bei `GP_ENTDECKE=1`, Ort bash, Sonde der Lauf mit `GP_MUT='^p15-(int|term)$'`); Vorgang wie `SPEC-057` Punkt 7, vom Vertrag Punkt 12 als äquivalent vorgesehen. Akzeptiert. `p15-schutz`: `$arbeit` stammt immer aus `mktemp -d -t pgr-cs-XXXXXX`, das Präfix ist über `p15-praefix` gedeckt; der Fall (e) „Name aus der Umgebung“ aus Punkt 15 hat in Teil I keine Eingabe. Akzeptiert, Teil II trägt (e) für Container.

**V-5 (niedrig, Größe).** §4 misst `git diff --numstat 8b19dc66..<letzter Commit> -- tools harness docs/developer`. Wörtlich gemessen: **1017** (Marke 900 überschritten, Meldeschwelle 700 überschritten); der Implementer nennt 734 „ohne Vertragsdateien“, weil §4 Vertrag und ADR ausnimmt (SPEC-058 279, SPEC-052 2, `vertraege/README.md` 2). Die Basis liegt vor den Commits des Architect, daher zählt die Metrik sie mit; die Ausnahme steht im Fließtext, nicht in der Metrik. Der Planner entscheidet die Lesart; Rückführung nach §4 trat bei 734 nicht ein. Meine Zahl pro Datei: Gegenprobe 512, Skript 197, Fragment 19, Rest 7 plus Vertrag 283.

**V-6 (Hinweis, Rückgaben R1 bis R6).** Sechs Fragen an den Architect stehen offen; der Code tut das in §7 Genannte. R1 (fehlendes `u-boot` macht `make gates` rot, kein Überspringen) ist wegen V-1 und der Host-Voraussetzung `u-boot` für das Gate relevant: das Gate hängt nun an `u-boot`, `node`, `ssh-keygen` im Pfad, `harness/README.md` nennt das nur bei der Gegenprobe.

## Plan-vs-Code-Diff

§3 Tabelle deckt alle sieben geänderten/neuen Dateien (Skript, Gegenprobe, zwei Fragmente, Ausnahmeliste, Fixture-Zeile, Index). Zusätzlich im Diff und nicht in §3: `harness/vertraege/SPEC-052-gate-kette.md` (1 Zeile) und `harness/vertraege/README.md` (2 Zeilen), Architect-Commits; [ADR-0045](../plan/adr/0045-container-start-als-werkzeug-im-repo.md) und der Index in `docs/plan/adr/README.md`. Kein Produktcode, kein `//nolint`, kein Eintrag in `docs/developer/` (wie §1 zusagt). §1 (Z1: Exit 6) und §6 stimmen mit Code und Vertrag überein.

## Nicht geprüft

Vollständiger `make gates` (Auftrag), Log des Implementer-Laufs, Läufe gegen echtes `devcontainer` und Container (Teil II).
