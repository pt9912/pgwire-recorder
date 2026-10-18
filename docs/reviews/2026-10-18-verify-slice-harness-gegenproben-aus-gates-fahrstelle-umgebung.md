# Verifikation slice-harness-gegenproben-aus-gates-fahrstelle-umgebung

Stand: Zweig `rev-fu`, Kopf `205ffcc0`, Basis der Implementierung `02a606eb`, Basis der Metrik `4219bab6`.
Maßstab: [Plan](../plan/planning/done/slice-harness-gegenproben-aus-gates-fahrstelle-umgebung.md) §2 und §7,
[SPEC-061](../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) (Plattform bash ≥ 5.1: Entscheidung des Nutzers, nicht gemessen),
[ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md), `AGENTS.md`.
Host-Läufe als Benutzer `db` (kein root), bash 5.2.21, coreutils 9.4. Keine volle Gegenprobe, kein `make gates` (Auftrag; die volle Gegenprobe lief laut §7 im Container grün).
Eigene Skripte und Kopien lagen außerhalb des Repos.

## Urteil: bestanden

0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO. Die Liefer-Punkte sind erfüllt und belegt; Abweichungen von Plan §1 bis §5 durch den Implementer gibt es nicht. Die Closure-Pflichten (Häkchen, Notiz, Risiko-Ausgänge) stehen noch aus und sind der Closure vorbehalten (Befund V-2).

## Ergebnis je Prüfpunkt

| # | Punkt | Ergebnis |
|---|---|---|
| 1 | DoD 1 und 2, Fälle rot ohne den Fix, Stichproben-Mutationen, Vergleich Hook alt/neu | bestanden |
| 2 | Plan-vs-Code-Diff (`AGENTS.md` §3.9) | bestanden |
| 3 | `AGENTS.md` §3.11: README-Zeile, Kommentare, Hilfetext | bestanden, 1 LOW (V-1) |
| 4 | Risiken aus §6 und ihre Ausgänge | vorgesehen, für die Closure vollständig (V-2) |
| 5 | Größenmetrik | 224, bestätigt |
| 6 | `make docs-check` am Kopf | Exit 0, 1156 Dateien, 0 Befunde |

## 1. DoD und Belege

**Fall ohne den Fix aus dem richtigen Grund rot (Modul 11).** Ich habe die Gegenprobe des Kopfes gegen den Hook von `02a606eb` gefahren (Kopie in Temp, nur `GP_MUT='^ok$'`, also die unveränderte Kopie).
Rot sind genau: `p5-eof-fd3 p5-eof-fd5 p5-eof-fd7 p5-eof-fd12 p5-fd-lauf p5-fd-make p5-fd-kind` (Deskriptoren), `p11-sigint-lauf p11-sigint-make p11-sigint-kind p11-sigint-enkel p11-sigquit-make p11-sigquit-kind p11-sigquit-enkel` (`SigIgn` aus `/proc`), `p11-erbt-int p11-erbt-quit` (Aufrufer mit `trap '' INT` bzw. `QUIT`) und `p13-werkzeug-env` (die Probe `env --default-signal` fehlt im alten Hook). Das deckt sich mit dem Beleg in §7 (dort ohne `p13-werkzeug-env`, weil dieser Fall später dazukam; das Rot ist die Folge der Probe, kein anderer Grund).
Der Hook am Kopf (unveränderte Kopie) ist grün.

**Aufgabe „ssh-ähnlicher Aufrufer“ (Deskriptor 5/7/12 offen vor dem Hook-Aufruf, Aufrufer mit `trap '' INT`).** Die Fälle `p5-eof-fd3/5/7/12` halten eine Pipe an dem Deskriptor, deren Leser bis EOF wartet, und prüfen vorab, dass der Deskriptor eine offene Pipe ist; `p11-erbt-*` rufen den Hook aus `sh -c 'trap "" INT; exec …'` ohne `timeout`. Alter Hook: rot (siehe oben), neuer Hook: grün. Das ist das geforderte Vorher/Nachher; ein echter `git push` über `ssh` ist akzeptiertes Negativ (`SPEC-061` Punkt 16 (11)), den Beleg trägt `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`.

**Stichproben-Mutationen, selbst gefahren (`GP_MUT` mit `GP_ENTDECKE=1`).** Gefahren sind 51 Mutationen der Gegenprobe (Muster `^(ok|m-p5-(keine-fd|nur|ausser|schliesst|wartet|gleiche|lauf)|m-p11-(hook|make-mit|lauf)|m-p13-ohne-env|m-bash)` und `^(m-bash|m-p5-erbt|m-p5-stdout|m-p5-gestartet|m-p9-marke|m-p1-liest|m-p12-nicht|m-p3-gleichzeitig)`), dazu die unveränderte Kopie. In allen Läufen: je Mutation genau die erwartete Menge rot, Schlusszeile „grün“. Namentlich geprüft, mit der gemessenen Menge:

| Mutation | rote Fälle (gemessen = erwartet) |
|---|---|
| `m-p5-keine-fd-schliessen` | `p5-eof-fd3/5/7/12 p5-fd-kind p5-fd-lauf p5-fd-make` |
| `m-p5-nur-ueber-3` | `p5-eof-fd3` |
| `m-p5-nur-fd5` | `p5-eof-fd3 p5-eof-fd7 p5-eof-fd12` |
| `m-p5-nur-3-bis-9` | `p5-eof-fd12` |
| `m-p5-schliesst-fd2` | `p5-fd-2` |
| `m-p5-ausser-5` / `m-p5-ausser-7` | `p5-eof-fd5 p5-fd-kind p5-fd-lauf p5-fd-make` / `p5-eof-fd7` |
| `m-p5-wartet` (`setsid -w`) | 15 Fälle, darunter `p5-ein-prozess`, `p5-vor-ende`, `p5-eof-*` |
| `m-p5-gleiche-sitzung` | `p5-gruppe` |
| `m-p5-lauf-als-kind` | `p5-ein-prozess` |
| `m-p5-lauf-ohne-exec` | nichts (äquivalent, siehe V-3) |
| `m-p11-hook-nur-int` / `-nur-quit` | `p11-erbt-quit` / `p11-erbt-int` |
| `m-p11-hook-ohne-reset` | `p11-erbt-int p11-erbt-quit` |
| `m-p11-hook-mit-und`, `m-p11-make-mit-und` | wie §7 |
| `m-p11-lauf-trap-int` / `-quit` | 5 bzw. 4 Fälle wie §7 |
| `m-p13-ohne-env` | `p13-werkzeug-env` |
| `m-bash-*` (7 Stück, `metag`) | je genau `p5-bash-version` |

Eigene Mutanten des Hooks, unabhängig von der Mutationsliste (Kopien in Temp, nur `^ok$`): (a) `-gt 2` → `-gt 4` im Schließen ergibt rot an genau `p5-eof-fd3`; (b) `--default-signal=INT,QUIT` → `--default-signal=TERM,HUP` ergibt rot an genau `p11-erbt-int p11-erbt-quit`. Beide Aussagen sind von der Gegenprobe gedeckt, ohne dass die Mutation in ihrer Liste stünde. Die Versionsprüfung: `GP_BASH_VERSION=5.0.19` endet mit Exit 2 und `bash >= 5.1 nötig (hier 5.0.19)`.

**Belege des Implementers in §7.** Vorhanden: Mutationstabelle, Befunde vorher/nachher, grüne Mutanten mit Invariante/Ort/Sonde, Sonde `p5-eof-*`, akzeptierte Negative, Größe, Läufe (Gegenprobe im Container 113 Fälle und 156 Mutationen, `make gates`, `gates-schnell-gegenprobe`, `docs-check`), Findings des Reviews mit Ausgang. Gezählt am Kopf: 113 Fälle (`FAELLE`), 149 `mut` plus 7 `metag` gleich 156. Fehlt kein Beleg, den die DoD verlangt. `make gates` und die volle Gegenprobe habe ich nicht neu gefahren (Auftrag); die Behauptung stützt sich auf die Zeitangaben in §7 und meine Stichproben.

**Liefer-Punkte.** L1 (Skript): `setsid -f env --default-signal=INT,QUIT "$BASH" -c '<Schleife über /proc/self/fd/*, alles über 2 schließen>; exec "$BASH" "$@"'`; ohne Nummernliste, `SIGINT`/`SIGQUIT` für `make` und Nachkommen nicht ignoriert (Fälle grün), `SIGINT` im Lauf nicht ignoriert. L2 (Gegenprobe): Fälle für Deskriptor 3, 5, 7, 12 mit Pipe, `SigIgn` für Lauf/make/Kind/Enkel, Aufrufer mit `trap`, je Mutation nur ihr Fall (Tabelle oben). Beide erfüllt.

## 2. Plan-vs-Code-Diff (`AGENTS.md` §3.9)

`git diff 02a606eb..HEAD -- docs/plan/planning/in-progress/…umgebung.md`: geändert sind die Tabelle in §6 (Spalte „Wo entschieden“, vom Architect, Cherry-picks `d12f5ab1`, `a25e068c`, `0fa80fce` aus Zweig `arch-fu`, dazu zwei neue Zeilen) und der neue Abschnitt „Belege des Implementers“ in §7. §1 bis §5, Kopf (`Bezug`, `Berührte Spec-Stellen`, `Verantwortlich`) und DoD sind unverändert. Der Implementer hat nur §7 geändert, wie er sagt.
Geändert außerhalb des Plans: `SPEC-061` (Architect), `tools/harness/gegenproben-nach-einlieferung.sh`, dessen Gegenprobe, `harness/README.md` (eine Zeile). Das entspricht §3 des Plans (drei Dateien); `SPEC-061` liegt außerhalb der Metrik-Pfade und ist Sache des Architect.
Die Zusage „setsid -f ändert keine Zusage von Punkt 5/9“ (Risiko 3): Von den Fällen des Bestands ist keiner entfernt oder geändert (`git diff` zeigt als Löschungen nur die Zeile der Fallliste, die Attrappe/Steuer-Zeilen, die Aufrufzeile mit `timeout` und die Mutationszeilen). Vier Mutationszeilen haben zusätzliche rote Fälle, in §7 je erklärt; `m-p5-wartet` und `m-p5-gleiche-sitzung` haben einen neuen Mutationstext, weil der Startbefehl sich ändert.

## 3. `AGENTS.md` §3.11 (Zusage nur, wo ein Test sie prüft)

- Zeile `make gegenproben-nach-einlieferung-gegenprobe` in `harness/README.md`: „113“ und „156“ stimmen mit dem Zählen am Kopf überein; „bash ≥ 5.1, sonst Exit 2“ (Fall `p5-bash-version`, Mutationen `m-bash-*`); „coreutils ≥ 8.31“ und „`env` mit `--default-signal`“ (Fall `p13-werkzeug-env`); Verweis „Punkt 16 (11)“ steht im Vertrag; „`SIGINT` im Lauf und `SIGINT` und `SIGQUIT` in `make` und seinen Nachkommen“ stimmt mit `g11s` überein (der Fall `p11-sigquit-lauf` gibt es bewusst nicht); „nicht geprüft: `SIGQUIT` in der Shell des Laufs“ ist genannt. Die Quantoren („ohne geerbten Deskriptor über 2“) tragen die Ausnahmen von §7 (Sperre 9, Skript-Deskriptor 255 im Lauf); die Zeile selbst nennt diese beiden Ausnahmen nicht (siehe V-1, nicht gesondert befundet, die Fälle `p5-fd-lauf` etc. prüfen Deskriptoren > 2 ohne diese beiden).
- Kommentare im Hook: der neue Kommentar nennt `SIGINT`/`SIGQUIT` nur für den Start durch `setsid -f` und `env`; Mutation `m-p11-hook-mit-und` existiert. Bestanden.
- Suche nach der Form „weder `SIGINT` noch `SIGQUIT`“ für den Lauf (F-3, beim Planner): im Plan §1/Titel/DoD wie bekannt; in README, Hook und Vertrag steht sie nicht. Einzig der Kommentar der Gegenprobe (V-1).

## Befunde

**V-1 (LOW) Kommentar der Gegenprobe sagt mehr, als `g11s` prüft.**
`tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` Zeile 445: „g11s — SIGINT und SIGQUIT in Lauf, make und Kind“. Die Funktion prüft `SIGQUIT` für `lauf` bewusst nicht (`[ "$s" = lauf ] || …`), der Vertrag nimmt `SIGQUIT` im Lauf selbst aus (Punkt 11, M1). Das ist dieselbe Form wie F-3 des Reviews, nur in einem Kommentar statt im Plan; `AGENTS.md` §3.11 verlangt, dass ein Kommentar nur zusagt, was ein Test prüft. Wirkung: gering (die Zeile hält niemand für eine Zusage des Hooks). Korrektur: „SIGINT in Lauf, SIGINT und SIGQUIT in make, Kind und Enkel“. Der Plan (§1/Titel/DoD) bleibt wie vereinbart beim Planner.

**V-2 (INFO) Closure-Pflichten stehen aus; Ausgänge der Risiken aus §6.**
Alle DoD-Häkchen sind leer, §7 trägt „(Wird bei der Closure geschrieben.)“; das ist vor der Closure zu erwarten. Risiken aus §6 und ihre vorgesehenen Ausgänge (in §6 selbst als „Ausgang (vorgesehen)“ eingetragen, in der Closure-Notiz von §7 zu übernehmen):
1. Weitere vererbte Zustände (Umgebung von `git`, `umask`, Arbeitsverzeichnis, Limits, `SigBlk`): weiter offen → `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale`. Der Eintrag existiert im Register (1×).
2. Die Fälle ersetzen die Fahrt nicht: weiter offen → `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (Eintrag und Beleg dieses Slice liegen im Register).
3. `setsid -f` ändert eine Zusage von Punkt 5 oder 9: **entfallen**, Begründung ist am Endstand belegt (kein Fall der Punkte 5 und 9 entfernt oder geändert; die Gegenprobe am Kopf grün). Die Closure trägt die Begründung ein.
4. Der Slice wächst durch Randformen: **entfallen** mit Messung (224 gegen Marke 900; die Randformen des Architect liegen in §6 vor dem Code; die Erweiterungen nach dem Review, bash-Prüfung und `env`-Probe, sind darin gemessen).
Nichts fehlt außer dem Eintrag in §7 selbst. Ausgang 3 und 4 sind Urteil: Beide tragen nach meiner Messung.
Die Folge-Slice-Paarung: `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (in `open/`) nennt die Kennung dieses Slice in §1 (Sendung), §4 (Start) und §6 (Risiko) — gefunden, Annahme gegeben.

**V-3 (INFO) Äquivalenz von `m-p5-lauf-ohne-exec` hängt an der Version, ihr Ort liegt in der Gegenprobe, nicht im Hook.**
§7 nennt Invariante (ab bash 5.1 kein `fork` für den letzten Befehl eines `-c`-Strings), Ort (Versionsprüfung der Gegenprobe) und Sonde (Messung des Architect mit `bash:4.4` bis `bash:5.2`; im Container `m-p5-lauf-als-kind` rot). Ich habe auf dem Host (bash 5.2.21) bestätigt: `m-p5-lauf-ohne-exec` bleibt grün, `m-p5-lauf-als-kind` ist rot an `p5-ein-prozess`; eine ältere bash konnte ich nicht fahren (Auftrag). Der Hook selbst prüft seine bash-Version nicht (Vertrag, Entscheidung des Nutzers). Das erfüllt `AGENTS.md` §3.10 (a) bis (c), mit dem Hinweis, dass der Ort die Gegenprobe ist, nicht der Hook. Kein Handlungsbedarf, es sei denn, der Nutzer will die Prüfung in den Hook ziehen.

**V-4 (INFO) Läufe unter Last.**
Laut §7 war ein voller Lauf (Last 30 bis 46, parallel `make gates`) rot an einer zusätzlichen Zeile `p13-werkzeug-git` bei `m-p5-nur-ueber-3`; zusammen mit `p2-kein-verzeichnis` sind das drei einmalige Rots des Bestands, Ursache ungeklärt (F-8 des Reviews). Meine 51 Mutationen mit `GP_PARALLEL` 4 am Host liefen ohne Rot. Das Muster gehört in das Register (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` Klasse Last); kein Befund gegen den Slice.

## 5. Größenmetrik

`git diff --numstat 4219bab6..HEAD -- tools harness/mk harness/README.md`: `harness/README.md` 1/1, Gegenprobe 192/23, Hook 6/1; Einfügungen 199, Löschungen 25, **Summe 224**. Der Implementer meldet 224: bestätigt. Schätzung des Plans 90 bis 180 (Faktor knapp 1,3 über der Obergrenze, ein Review-Zyklus mit bash-Prüfung und `env`-Probe kam hinzu); Meldeschwelle 700 und Marke 900 weit unterschritten.

## 6. `make docs-check`

Am Kopf (vor dem Anlegen dieses Reports): `d-check: 1156 Datei(en) geprüft, 0 Befund(e)`. Der Report selbst wird vor dem Commit erneut geprüft (siehe Commit-Nachricht, wenn rot).

## Gefahren

Gegenprobe am Host: 51 Mutationen einzeln plus unveränderte Kopie, 3 Kopien mit Hook-Mutanten (alt, Schwelle 4, falsche Signalliste); `make docs-check`; Zählen von Fällen und Mutationen; `git diff`-Prüfungen von Plan, Hook, Gegenprobe, README. Nicht gefahren: volle Gegenprobe, `make gates`, ältere bash-Versionen, echter `git push` über `ssh`.
