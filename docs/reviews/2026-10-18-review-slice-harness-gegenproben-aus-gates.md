# Review-Report: slice-harness-gegenproben-aus-gates — 2026-10-18

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; das ist die Frage des Verifiers).

**Gegenstand:** `slice-harness-gegenproben-aus-gates` (Kettenumbau), Zweig `rev-kette`, Diff `b4ce87d3..2f9717fc` über `tools harness AGENTS.md docs`.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

> **Zitier-Form.** Dieser Report friert ein; Slices werden mit Kennung zitiert. Die ADR- und SPEC-Kennungen sind auf Wunsch des Auftrags auflösende Links.

**Eingangs-Kontext:**

- `slice-harness-gegenproben-aus-gates` (Plan in `in-progress/`, §7 *Belege des Implementers*: Mutationstabelle, Läufe, Metrik, Rückgaben 1 bis 4)
- [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md) (Accepted), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md), [ADR-0047](../plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)
- [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Teil 1 (Punkte 1 bis 3; Teil 2 gehört `slice-harness-gegenproben-aus-gates-sammelziele`), [SPEC-052](../../harness/vertraege/SPEC-052-gate-kette.md) (Punkte 1, 2, 4, 9, 13)
- keine `LH-*`-Anforderung des Produkts berührt (Plan: nur Abdeckungs-Deklaration der Gegenproben)
- `AGENTS.md` (Hard Rules §3.7, §3.9, §3.10, §3.11, §3.12, §3.14; §6 Schritt 6)

Gefahren, engster Sensor (Auftrag des Nutzers): `bash -n` über alle Skripte unter `tools/`, `make -pq help` im Arbeitsbaum des Zweigs mit dem Code von `namensregel` in zwei Locales, `make docs-check`, `make kopf-check`. Keine Gegenprobe, kein `make gates`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Namensregel erkennt ein Ziel mit Rezept am englischen Text `#  recipe to execute (from '…` der Ausgabe von `make -p`. Unter `LANG=de_DE.UTF-8` (GNU make 4.3) heißt die Zeile `Auszuführendes Rezept (aus „…“`; `gp_rezeptziele` liefert nichts, und `namensregel` meldet für den unveränderten Arbeitsbaum dreizehnmal `p2-v6-container`/`p2-v6-host`. Damit ist `make gates-schnell-gegenprobe` am Host mit deutscher Locale rot, im Dev-Container (Locale `C`) grün; `LC_ALL=C` setzt weder die Gegenprobe noch die Regel, und der §7-Beleg „Sonde: GNU make 4.4.1 im Container“ deckt nur eine Locale. | Maintainability; [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Plattform (Linux, GNU make) und Punkt 2 (Regel gilt „im Arbeitsbaum selbst wie in der Kopie“); `AGENTS.md` §3.10 | `tools/harness/gegenproben-regel.sh` · `/^#  recipe to execute \(from .harness\/mk\//`; `tools/harness/gates-schnell-gegenprobe.sh` · `make -s --no-print-directory -C "$1" -pq help 2> /dev/null \|\| true` | ja — `. tools/harness/gegenproben-regel.sh; make -s -pq help > db; namensregel db` unter `LANG=de_DE.UTF-8` (dreizehn Zeilen) und unter `LC_ALL=C` (leer) | Ausgabe von make locale-abhängig gelesen |
| F-2 | MEDIUM | Die Zusage, dass die Regel im Arbeitsbaum selbst gilt (SPEC-060 Punkt 2, „im Arbeitsbaum selbst wie in der Kopie“; Gate-Index: „und die Regel gegen den Arbeitsbaum selbst“), hat keinen Fall und keine Mutation in der Gegenprobe. §7 führt die Zeile als „n/a … von Hand gesehen“; streicht man die zwei Zeilen `ist="$(deckung "$wurzel" …)"` und `[ -z "$ist" ] \|\| melde "Arbeitsbaum: …"`, bleiben alle Fälle und alle Mutationen der Regel grün (die Kopien rufen `deckung` über `pruefe`). | `AGENTS.md` §3.10 (Ortsangabe zählt je Schritt), §3.11 | `tools/harness/gates-schnell-gegenprobe.sh` · `[ -z "$ist" ] \|\| melde "Arbeitsbaum: Ausnahmeliste und GATE_CHECKS decken sich nicht: $ist"`; `harness/README.md` · `und die Regel gegen den Arbeitsbaum selbst` | ja — die zwei Zeilen streichen, `make gates-schnell-gegenprobe` bleibt grün | grüner Mutant: Prüfung im Arbeitsbaum ohne Mutation |
| F-3 | LOW | Dass `make gates-schnell` keine Gegenprobe fährt (Gate-Index: „gehören nicht zu … `make gates-schnell`“; [SPEC-052](../../harness/vertraege/SPEC-052-gate-kette.md) Punkt 1 und 2), prüft die Gegenprobe nur über `GATE_CHECKS` (V1). Eine Änderung des Rezepts von `gates-schnell`, die `$(GEGENPROBEN_CONTAINER)` mitfährt, färbt keinen Fall: der Lauf der Kopie prüft nur, dass die Gates aus `LANGSAM` fehlen und die aus `SCHNELL` laufen; die Gegenproben der Fixture stehen nicht im Log-Vergleich. | `AGENTS.md` §3.10, §3.11 | `tools/harness/gates-schnell-gegenprobe.sh` · `for g in $LANGSAM; do grep -qxF "$g" "$arbeit/log" && echo "p2-langsam"; done` | ja — Rezept von `gates-schnell` in der Kopie um die Liste erweitern, Gegenprobe bleibt grün | Zusage zur schnellen Kette ohne Mutation |
| F-4 | LOW | Zu `p2-andere-namen` listet §7 „vier Mutationen, je eine für `x-gegenprobe-alt`, `gegenprobe`, `gegenproben` und `gegenproben-host`, `x-negativ`“. Im Code sind es vier Mutationen für fünf Namen; `gegenproben` und `gegenproben-host` fallen in dieselbe (Präfix-)Mutation, und kein Mutant bricht nur einen von beiden. Der Fall ist ein einziger Fall mit fünf Namen, so dass die Aufzählung nicht je Element getrennt belegt ist. | `AGENTS.md` §3.10 (Aufzählung je Element), §3.11 | `tools/harness/gates-schnell-gegenprobe.sh` · `regelmut p2-andere-praefix '*-gegenprobe) echo n1' 'gegenprobe*\|*-gegenprobe) echo n1' "p2-andere-namen"` | ja | Aufzählung ohne Mutation je Element |
| F-5 | LOW | Der Kopf von `gates-schnell-gegenprobe.sh` sagt „Ihre elf Fälle (p2-*)“; das Skript führt dreizehn (`p2-v4-host`, `p2-rezept-ausserhalb-fragment` dazu), und Meldung und Gate-Index sagen dreizehn. | `AGENTS.md` §3.7, §3.11 | `tools/harness/gates-schnell-gegenprobe.sh` · `Ihre elf Fälle (p2-*) legt je eine` | ja — Zählung der `regelfall`-Aufrufe | Zahl im Kommentar weicht vom Code ab |
| F-6 | LOW | Die Zeile von `make container-start-gegenprobe` im Gate-Index trägt den eingefügten Halbsatz mitten in einer Aufzählung („… der Fall `p16-keine-kette` (…); sie steht in `GEGENPROBEN_HOST` und das Aufräumen der Temp-Verzeichnisse bei Ende …“), so dass die Aufzählung „dazu der Fall …, der Fall … und das Aufräumen“ zerreißt. Dass sie in der Liste für den Host steht und dass sie dort „heute nur“ allein steht, prüft kein Fall: die Regel verlangt nur genau eine Liste, `p16-keine-kette` prüft die Abwesenheit in `GATE_CHECKS` und Ausnahmeliste. | `AGENTS.md` §3.7, §3.11 | `harness/README.md` · `sie steht in `GEGENPROBEN_HOST` und das Aufräumen der Temp-Verzeichnisse` | ja — Lesen; Verschieben von `container-start-gegenprobe` in die Liste für den Container lässt `make gates-schnell-gegenprobe` grün ([SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 3 nennt das als akzeptiertes Negativ, nicht aber die Prosa im Gate-Index) | Prosa im Gate-Index weiter als Prüfung / zerrissener Satz |
| F-7 | LOW | Der Plan folgt dem Diff nicht an vier Stellen: (a) §3 und §2 führen `gegenproben-regel.sh` nicht (neue Datei, nur in §7 „Neu“) und `harness/mk/container-start-gegenprobe.mk` und `tools/harness/container-start-gegenprobe.sh` nicht (beide geändert); (b) §2 und §3 sagen „11 Fälle“, die Gegenprobe führt dreizehn; (c) §3 verlangt den Nachzug des Kommentarkopfs von `container-start.mk`, der Diff ändert ihn nicht (sein Text bleibt wahr); (d) §4 Start (a) und der Kopf sagen „Commits des Architect vor dem Move“, in der Historie stehen `66ad7abc` (SPEC-060) und `93434376` (ADR) nach dem Move `350139c5`. Zusätzlich sagen §6 (K1 bis K7) und der Kopf von [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) („bis zur Annahme der ADR eine Empfehlung“) noch Empfehlung, obwohl [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md) Accepted ist. | `AGENTS.md` §3.9, §3.7 | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates.md` · `die Fälle von Teil 1 des Vertrags (nach dem Architect 11)`; `harness/vertraege/SPEC-060-gegenproben.md` · `bis zur Annahme der ADR eine Empfehlung` | nein — Lesen | Plan folgt dem Diff nicht (`BEO-REPO/plan-folgt-korrektur-nicht`) |
| F-8 | LOW | Für `lint-gegenprobe` und `coverage-gegenprobe` hat der Implementer die neuen Fälle `p10-gegenprobe-liste`/`p11-gegenprobe-liste` nur als „Code des Falls an der Datenbank einer Kopie“ rot gesehen (Lauf der Gegenprobe 283 s bzw. 105 s nicht gefahren); `kopf-check` und `meldungskatalog-check` von Hand an Kopien. Nur `container-`, `postgres-versionen-` und `stop-hook-gegenprobe` tragen die Mutation im Skript. Die Fälle prüfen die Datenbank direkt, die Lücke ist klein; dass die Mutation vom Skript gefahren wird, behauptet §7 für diese vier nicht. | `AGENTS.md` §3.10 | `tools/harness/lint-gegenprobe.sh` · `fall=p10-gegenprobe-liste`; `tools/harness/coverage-gegenprobe.sh` · `p11-gegenprobe-liste` | ja — Mutation `GEGENPROBEN_CONTAINER += lint-gegenprobe` streichen, Lauf der Gegenprobe | Anschluss-Fall ohne Mutation im Skript |
| F-9 | INFO | Die Gegenprobe baut ihre Fixture mit `GEGENPROBEN_CONTAINER := …`/`GEGENPROBEN_HOST := …`; die Fragmente des Arbeitsbaums hängen mit `+=` an eine nicht gesetzte Variable (rekursive Variable, `make -p` zeigt den Wert unausgewertet). Die Regel liest beide Formen gleich, solange die Fragmente Namen wörtlich nennen; bildete ein Fragment einen Eintrag aus einer Variablen, sähe `namensregel` `$(X)` als Eintrag ohne Ziel (V6). Das ist Rückgabe 3 des Implementers an den Architect (Initialisierung in `Makefile` oder Vertrag), nicht Gegenstand dieses Diffs. | Maintainability | `harness/mk/gates-schnell.mk` · `GEGENPROBEN_CONTAINER += gates-schnell-gegenprobe` | nein | Fixture-Form weicht von der Form der Fragmente ab |
| F-10 | INFO | Rückgaben des Implementers: (1) `container-start-gegenprobe` am Host einmal ganz fahren (`u-boot`, `devcontainer` fehlen im Dev-Container): berechtigt, Verifier/Nutzer am Host; der geänderte Anker `GATE_SCHNELL_AUSNAHME := test-container ` trifft die Zeile in `harness/mk/gates-schnell.mk` (Leerzeichen folgt). (2) „§4-Satz laut Plan“ in `AGENTS.md`: der Diff ändert nur die erste Hälfte von Schritt 6 (eine Zeile), Wortlaut bis „Zuordnung in `harness/README.md` §Sensors“, Punkt statt Semikolon; die zweite Hälfte (Ende der Welle, Release, Sammelziele) steht nicht da, wie `AGENTS.md` §3.11 es verlangt. (3) Zwei Fälle über den Vertrag hinaus: kein Verstoß, Aufnahme oder Streichung liegt beim Architect. (4) `Makefile` Z. 2 („jedes [Fragment] haengt seine Checks an GATE_CHECKS“) ist für die Gegenproben nicht mehr wahr; Datei ist generiert und nicht in der Dateiliste des Plans, Rückgabe berechtigt. | Maintainability | Plan §7 · `Rückgaben` | nein | — |
| F-11 | INFO | Metrik nachgerechnet: `git diff --numstat b4ce87d3..HEAD -- tools harness AGENTS.md` = 527 + 180 = **707**; davon `harness/vertraege/` 225 + 44 = 269; ohne sie **438** (302 + 136), Basis-Wechsel auf `9ce157a4` ändert die Zahl nicht, weil `9ce157a4` nur Pläne trägt. Marke 900 nicht erreicht, Meldeschwelle 700 in der wörtlichen Lesart um 7 überschritten, wie §7 sagt; §7 nennt beide Zahlen. | `AGENTS.md` §3.14 | Plan §7 · `Metrik` | ja — Befehl wie oben | — |

Geprüft, ohne Befund, zu den Aufträgen des Auftrags:

- **Je Element der Zusagen von SPEC-060 Teil 1** (Punkt 2: N1, N2, V1 bis V4 je Klasse, V4 je Liste, V6 je Liste, `p2-andere-namen`, Rezept aus Fragment): Fall und Mutation, die nur ihn färbt, im Code vorhanden (`regelfall`/`regelmut`); die Mutationen ersetzen genau eine Stelle und melden „ändert die Regel nicht“, wenn der Anker fehlt. Punkt 1 (Reihenfolge der Liste, Sammelziele in keiner Liste) und Punkt 3 (Wahl der Liste) tragen keinen Fall: Punkt 3 ist akzeptiertes Negativ, Reihenfolge und Sammelziele gehören dem Folge-Slice.
- **AGENTS.md §6 Schritt 6:** nur die erlaubte Änderung, eine Zeile, erste Hälfte des vom Nutzer entschiedenen Satzes; keine Ziele genannt, die es noch nicht gibt.
- **Kommentare der Fragmente** (`harness/mk/*.mk`): Indikativ, nennen die Liste; `abdeckung.mk` behält den Satz zum Anschluss von `abdeckung-check` (wahr). Die Gegenprobe-Köpfe von `kopf-check`, `lint`, `meldungskatalog-check`, `coverage`, `stop-hook` stimmen mit den Fällen überein.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `harness/mk/` (zwölf geänderte Fragmente + `container-start-gegenprobe.mk`) | geprüft, ohne Befund (außer F-9 als INFO); je Gegenprobe eine Zeile in genau einer Liste, `GATE_CHECKS` ohne Gegenprobe |
| `harness/sensors/` (vier Texte) | geprüft, ohne Befund; Aussagen „hängt an der Gate-Kette“ für Gegenproben ersetzt |
| `harness/conventions.md` | geprüft, ohne Befund (Zeile `modul-13-quality-gates.md` nennt beide Listen) |
| `harness/README.md` §Sensors | geprüft; Befund F-2 (Prosa über Arbeitsbaum), F-6 (Zeile von `container-start-gegenprobe`); jede andere Gegenprobe-Zeile nennt ihre Liste |
| `harness/vertraege/` (SPEC-052, -053, -054, -055, -056, -057, -058, README) | geprüft, ohne Befund außer F-7 (SPEC-060 Kopf); Nachzug stimmt mit der Ausnahmeliste auf drei Gates überein |
| `tools/harness/gegenproben-regel.sh` | geprüft; Befund F-1; `bash -n` sauber |
| `tools/harness/gates-schnell-gegenprobe.sh` | geprüft; Befunde F-2, F-3, F-4, F-5; `bash -n` sauber |
| `tools/harness/{container,coverage,kopf-check,lint,meldungskatalog-check,stop-hook,container-start}-gegenprobe.sh`, `tools/test/{postgres-versionen-gegenprobe,abdeckung}.sh` | geprüft; Befund F-8; die Umbenennung `liste` zu `gp_liste` in `container-` und `postgres-versionen-gegenprobe.sh` löst den gemeldeten Namenskonflikt, `bash -n` sauber |
| `AGENTS.md` | geprüft, ohne Befund |
| `docs/plan/` (Plan, ADR-Index, die ADR zu den Gegenproben) | geprüft; Befund F-7; keine Accepted-ADR inhaltlich geändert ([ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md) neu, Status-Wechsel in `0ba4b696`) |
| Commit-Reihenfolge Move/Inhalt | geprüft, ohne Befund: `350139c5` ist der reine Move, Inhalt danach |

`make docs-check` (nach dem Link in der Negativbefund-Zeile zu den ADRs) und `make kopf-check` im Worktree: grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Ausgabe von make locale-abhängig gelesen · grüner Mutant: Prüfung im Arbeitsbaum ohne Mutation · Zusage zur schnellen Kette ohne Mutation · Aufzählung ohne Mutation je Element · Zahl im Kommentar weicht vom Code ab · Prosa im Gate-Index weiter als Prüfung · Plan folgt dem Diff nicht · Anschluss-Fall ohne Mutation im Skript

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM). F-1 macht die Gegenprobe der Kette am Host mit deutscher Locale rot, F-2 lässt die Prüfung der Regel im Arbeitsbaum ohne Mutation. LOW und INFO blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; Rückgaben 3 (Fälle über den Vertrag hinaus, Initialisierung der Listen) und 4 (`Makefile` Z. 2) an Architect/Planner, der Rest von F-7 (Plan, SPEC-060 Kopf) an Planner/Architect. Die Finding-Klassen gehen in die Closure §7 und von dort in den Zähler. Dieser Report ersetzt keine Verifikation.
