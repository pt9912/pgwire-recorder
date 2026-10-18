# Review-Report: slice-harness-gegenproben-aus-gates-fahrstelle-umgebung — 2026-10-18

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung`, Zweig `rev-fu`, `git diff 02a606eb..75a3072f` (Kopf `75a3072f`; Basis der Größenmetrik `4219bab6`). Der Diff enthält per Cherry-Pick `597f2ed0` die Entscheidung des Architect (Commit `d12f5ab1`) zu [SPEC-061](../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) Punkte 5 und 11 sowie die §6-Zeilen des Plans.

**Skill:** `.harness/skills/reviewer.md` @ `229b5efb` · <!-- d-check:ignore (Skill-Pfad des Repos, kein Link-Ziel dieses Reports) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

> **Zitier-Form.** Dieser Report friert ein; er nennt Kennungen (`slice-<Kennung>`, `make <target>`), keine Lifecycle-Pfade. Ein `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Plan `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` (§1, §3, §4, §6, §7 *Belege des Implementers*)
- [SPEC-061](../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) Punkte 5, 9, 11, 13, 16 (Fassung des Architect)
- `AGENTS.md` §3.2, §3.7 bis §3.12, §3.14; `harness/README.md` §Sensors
- keine `LH-*`-Kennung im Plan-Kopf; keine ADR im Diff genannt
- Reviewer-Eigenläufe (Host, bash 5.2.21, Benutzer ohne root, Last gering): 12 Mutationen des Bestands mit `GP_MUT` und `GP_ENTDECKE=1` (`m-p5-keine-fd-schliessen`, `m-p5-nur-fd5`, `m-p5-nur-3-bis-9`, `m-p5-nur-ueber-3`, `m-p5-schliesst-fd2`, `m-p5-lauf-ohne-exec`, `m-p5-lauf-als-kind`, `m-p5-make-mit-sperre`, `m-p11-hook-ohne-reset`, `m-p11-make-mit-und`, `m-p11-lauf-trap-int`, `m-p11-lauf-trap-quit`) plus die unveränderte Kopie: alle rot genau wie in der Tabelle von §7, Abschlusszeile `gruen`. Danach in einer Kopie unter dem Scratchpad (nichts im Repo geändert) sieben eigene Mutationen und die Sonde der Pipe (siehe F-1, F-2). Keine volle Gegenprobe, kein `make gates`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Mutationstabelle hat keine Mutation, die nur ein einzelnes Element einer Aufzählung bricht: `env --default-signal=INT` allein und `=QUIT` allein fehlen (nur beide zugleich entfallen), ebenso „alle Deskriptoren außer 5“ und „außer 7“ (vorhanden sind „nur 5“, „nur 3 bis 9“, „über 3“). Die Fälle fangen sie: meine Kopie mit `=INT` färbt genau `p11-erbt-quit`, mit `=QUIT` genau `p11-erbt-int`, „außer 7“ genau `p5-eof-fd7`, „außer 5“ `p5-eof-fd5` und `p5-fd-lauf`, `p5-fd-make`, `p5-fd-kind`; im Bestand fährt sie niemand. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` · `mut m-p11-hook-ohne-reset "p11-erbt-int p11-erbt-quit"` | ja — `GP_MUT` mit den vier Mutationen in der Kopie | Mutation je Aufzählungselement fehlt (Fall vorhanden) |
| F-2 | MEDIUM | Der grüne Mutant `exec "$BASH" "$@"` → `"$BASH" "$@"` (`m-p5-lauf-ohne-exec`) gilt in §7 als äquivalent; der Ort ist das Verhalten von `bash -c`, nicht eine Stelle im Repo, und es hängt an der Version: das Auslassen des `fork` für den letzten Befehl einer Liste im `-c`-String kenne ich aus den Änderungen von bash 5.1 (am Host nicht gemessen, nur 5.2.21 steht bereit; dort bleibt der Mutant grün, die Sonde `…; :` rot an `p5-ein-prozess`). Die Gegenprobe nennt „bash ≥ 5“, [SPEC-061](../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) Kopf „bash ≥ 4“: Unter einer bash ohne diese Optimierung wäre der Mutant nicht äquivalent, `m-p5-lauf-ohne-exec` meldete „erwartet nichts“ mit rotem `p5-ein-prozess` und die Gegenprobe wäre rot. | `AGENTS.md` §3.10 „Ein grüner Mutant gilt als äquivalent erst mit Invariante, Ort und Sonde“ | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates-fahrstelle-umgebung.md` · „(b) Ort: die Ausführung des `-c`-Strings in bash“ | nein — erst ein Lauf mit bash 5.0 würde es zeigen (nicht verfügbar) | Äquivalenz stützt sich auf Interpreter-Version, Ort nicht im Repo |
| F-3 | LOW | Titel, §1 *Ziel* (2) und DoD Liefer-Punkt 1 des Plans sagen, der Lauf ignoriere „weder `SIGINT` noch `SIGQUIT`“ und `SIGINT`/`SIGQUIT` seien „in Lauf, `make` und den Kindern“ nicht ignoriert. Punkt 11 in der Fassung des Architect nimmt `SIGQUIT` im Lauf selbst aus (bash ignoriert es selbst), und der Code liefert das; §1 und DoD sind nach der Entscheidung nicht nachgezogen (die Änderung am Plan beschränkt sich auf §6-Zeilen und §7). | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates-fahrstelle-umgebung.md` · „weder `SIGINT` noch `SIGQUIT` auf „ignoriert“; die Gegenprobe des Hooks zeigt beides“ | nein — `make kopf-check` prüft nur den Kopf | Plan folgt der Entscheidung des Architect nicht (§1, DoD) |
| F-4 | LOW | Der Verweis „Punkt 16 (10)“ für das Negativ „echter `git push` über `ssh` ist nicht geprüft“ löst nicht auf: (10) beschreibt den Start mit `setsid -f` und die Plattform und nennt weder `ssh` noch das Nicht-Prüfen des Pushs. Das Negativ steht in Punkt 5; README-Zeile und §7 zitieren den Verweis weiter. | `AGENTS.md` §3.11 (Rang-Zeiger; Zusage nur, was geprüft ist) | `harness/README.md` · „nicht geprüft sind ein echter `git push` über `ssh` (Punkt 16 (10))“ | ja — Lesen von Punkt 16 (10) gegen den Verweis | Zeiger auf Vertragsstelle trägt die Aussage nicht |
| F-5 | INFO | Kommentar im Hook: „`setsid -f` statt `cmd &`: ein einfacher Hintergrundbefehl … startet mit ignoriertem SIGINT und SIGQUIT“ nennt die verworfene Alternative. Er schreibt an den, der die Zeile ändert, und trägt eine Kopplung (Grund, warum `&` hier falsch ist; Test: `m-p11-hook-mit-und`); ich werte ihn als Abgrenzung, nicht als Befund der Klasse „verworfene Alternative“. Grenzfall, Urteil beim Architect. | `AGENTS.md` §3.7 | `tools/harness/gegenproben-nach-einlieferung.sh` · „`setsid -f` statt `cmd &`“ | nein | Kommentar nennt verworfene Alternative (Grenzfall Abgrenzung) |
| F-6 | INFO (an den Architect) | Der Hook prüft `env` mit `--default-signal` nicht (Punkt 13 prüft `git flock setsid make pgrep`; Punkt 16 (10) nimmt nur `setsid -f` als akzeptiertes Negativ aus). Fehlt die Option (coreutils < 8.31), endet `env` im verworfenen stderr, kein Lauf entsteht, die Marke bleibt verwaist, der Hook endet mit Exit 0. Der Kopf des Vertrags nennt „GNU-Coreutils“ ohne Version. Ebenso ohne Fall und ohne Meldung: kein `/proc` (die Schleife läuft leer, der Push könnte hängen); der Kopf bindet die Zusage an Linux. | `AGENTS.md` §3.12 (Plattform gehört dazu) | `tools/harness/gegenproben-nach-einlieferung.sh` · `setsid -f env --default-signal=INT,QUIT` | nein | Voraussetzung des Aufrufs im Vertrag nicht als Negativ geführt |
| F-7 | INFO (an den Architect, Anmerkung des Implementers) | Die Listen im Vertrag stimmen nicht mit den gemessenen Mengen überein: Punkt 11 nennt für „Hook mit `cmd &`“, `trap '' INT`, `trap '' QUIT` und „`make` mit `cmd &`“ Mengen ohne `p11-erbt-*`, gemessen färben diese Mutationen sie mit (bestätigt: `m-p11-lauf-trap-int`, `m-p11-lauf-trap-quit`, `m-p11-make-mit-und`). Zusätzlich nennt Punkt 5 weder den Fall `p5-eof-fd3` noch `p5-ein-prozess` und nicht die Mutationen „über 3“, „ab 2“ und „Lauf bleibt Elternprozess“, die Gegenprobe führt sie. Kein Befund gegen den Implementer. | Maintainability | `harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md` · „Mutationen: der Hook startet den Lauf mit `cmd &` statt `setsid -f` (färbt die sieben Fälle ohne `erbt`)“ | ja — `GP_ENTDECKE=1` gegen die Listen | Vertragsliste und gemessene Menge weichen ab |
| F-8 | INFO | Zu `p2-kein-verzeichnis` (zweimal bei Last über 30 zusätzlich rot, Ursache ungeklärt): §7 benennt beide Vorkommen mit Mutation, Last und den grünen Wiederholungen, und nennt den Lauf mit Rot und den Wiederholungslauf getrennt. Der Beleg ist ehrlich; geklärt ist die Flackerigkeit damit nicht, und ein Gate-naher Fall des Bestands bleibt unter Last unzuverlässig. Eigener Lauf (Last unter 4): `ok` und 12 Mutationen grün. | `AGENTS.md` §3.11 | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates-fahrstelle-umgebung.md` · „die Ursache ist ungeklärt“ | nein | Flackernder Fall des Bestands, Ursache offen |

### Zu den Prüfpunkten (a) bis (g)

- **(a)** SPEC-061 Punkt 5: Deskriptoren 3, 5, 7 und 12 haben je einen Fall (`p5-eof-fd3|5|7|12`); die Sonde „Deskriptor ist eine offene Pipe“ bestätigt: mit `/dev/null` statt Pipe in `eof_gruppe` sind alle vier rot mit der Meldung des Vertrags. 0 bis 2 bleiben offen: `p5-fd-0|1|2`, Mutationen `m-p5-erbt-stdin|stdout`, `m-p5-schliesst-fd2` (rot genau `p5-fd-2`, selbst gefahren). Punkt 11: `SIGINT` im Lauf, `SIGINT`/`SIGQUIT` in `make`, Kind, Enkel je ein Fall (`g11s`); Reset bei geerbtem ignorierten Signal `p11-erbt-int|quit` (`g11e`, ohne `timeout`, begründet durch M2). Die Mutationen für Kind und Enkel färben immer `make` mit, weil ein ignoriertes Signal nach unten erbt; eine Mutation, die nur Kind oder Enkel bricht, gibt es im Hook-Code nicht (keine Lücke, aber in §7 nicht begründet). Einzelelement-Lücken: F-1. Selbst gefahren: 12 Mutationen des Bestands, 7 eigene, alle wie erwartet.
- **(b)** Hook-Code: Pfad mit Leerzeichen und Argumente mit `$(…)`/Apostroph kommen unverändert in `"$@"` an (Probe im Verzeichnis `a b`); Aufruf ohne Argumente startet nichts (`args` leer, `p4-stdin-leer`); die Schleife über `/proc/self/fd/*` schließt jede Nummer über 2 (Eigenfall, auch 255 entfällt in `bash -c` nicht, weil dort keiner steht); `haupt` endet mit `return 0`, `exit 0` am Ende, `setsid -f` ohne `set -e`: Exit 0 in jedem Fall. Fehlendes `/proc` und fehlendes `env --default-signal`: F-6. Anker im Vertrag: F-4.
- **(c)** README-Zeile: 111 Fälle (Lauf: „gefahren: 111 Fälle“) und 144 Mutationen (`grep -c '^mut '` = 144) stimmen; die Sätze „nicht ignoriert, auch wenn der Aufrufer des Hooks sie ignoriert“ decken `p11-erbt-*` (Prüfung an `make` und Kind, Lauf und Enkel folgen aus der Vererbung); Ausnahmen der Zeile (Sperre 9, Skript-Deskriptor 255, `SIGQUIT` der Shell des Laufs) stehen in §7 und in der Zeile unter „nicht geprüft“. Befund nur F-4.
- **(d)** Der Diff des Plans gegen die Basis trägt zwei Blöcke: §6-Zeilen (aus `597f2ed0`, Architect) und §7 (Implementer); §1 bis §5 sind unverändert. Die Aussage in §7 stimmt. §1/DoD folgen der Entscheidung nicht: F-3.
- **(e)** `m-p5-lauf-ohne-exec`: Invariante, Ort und Sonde stehen in §7, die Sonde (`…; :` rot, ohne `exec` grün an `p5-ein-prozess`) ist am Host reproduziert; Ort und Version: F-2.
- **(f)** Siehe F-8.
- **(g)** Größe nach `AGENTS.md` §3.14: Summe der `--numstat`-Spalten über `tools`, `harness/mk` und `harness/README.md` von `4219bab6..HEAD` = 1+1 (README) + 141+19 (Gegenprobe) + 5+1 (Hook) = 168; unter Meldeschwelle 700 und Marke 900; passt zu §7 (147 Einfügungen, 21 Löschungen).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/gegenproben-nach-einlieferung.sh` (Start des Laufs, Fortbestand der Zusagen aus Punkt 5 und 9) | geprüft, ohne Befund außer F-5 und F-6 |
| `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` (Attrappe, Hilfen `rec`/`kind`, Fälle `g5*`, `g11s`, `g11e`, Mutationsteil) | geprüft, ohne Befund außer F-1 |
| `harness/README.md` §Sensors (Zeile der Gegenprobe) | geprüft; Zahlen und Quantoren belegt, Befund F-4 |
| `harness/vertraege/` (Fassung des Architect) | geprüft gegen Code; Anmerkungen F-6, F-7 an den Architect |
| `docs/plan/planning/` (Plan §6 und §7) | geprüft; F-3, F-8 |
| `harness/mk/`, `.golangci.yml`, `cmd/`, `internal/`, `spec/` | nicht im Diff, nicht geprüft |
| Suppression (`nolint`, Ausnahmen im Lint-Profil) | geprüft, ohne Befund (keine im Diff) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Mutation je Aufzählungselement fehlt (Fall vorhanden) · Äquivalenz stützt sich auf Interpreter-Version, Ort nicht im Repo · Plan folgt der Entscheidung des Architect nicht (§1, DoD) · Zeiger auf Vertragsstelle trägt die Aussage nicht · Voraussetzung des Aufrufs im Vertrag nicht als Negativ geführt · Vertragsliste und gemessene Menge weichen ab · Flackernder Fall des Bestands, Ursache offen

## Verdikt

**Merge-blockierend:** ja, wegen F-2 (MEDIUM, nach Regel des Skills); der Befund ist klein: §7 ergänzt, welche Eigenschaft den Mutanten an die Version bindet, oder der Mutant bekommt einen Fall. Kein Befund im Hook-Code selbst: Alle Zusagen von Punkt 5 und 11 sind mit Fall und Mutation belegt, die Eigenläufe bestätigen sie.

**Übergabe:** F-1, F-2, F-4 an den Implementer; F-3, F-6, F-7 an den Architect (Plan §1/DoD und Vertrag); F-5, F-8 zur Kenntnis. Die Finding-Klassen gehen in die Slice-Closure §7.
