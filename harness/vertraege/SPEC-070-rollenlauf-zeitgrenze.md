# SPEC-070 — Zeitgrenze eines Rollenlaufs (`pretooluse-zeitgrenze`, `zeitgrenze-gegenprobe`)

Ein Lauf eines Subagenten (eine Rolle unter `.claude/agents/`) hat höchstens **eine Stunde**. Der Hook
`.claude/hooks/pretooluse-zeitgrenze.sh` setzt das im Werkzeug-Aufruf durch: Er liest den Strom der Erfassung
des laufenden Agenten, rechnet die Wanduhr seit dem Beginn des Laufs und sperrt ab der Grenze die Aufrufe, die
den Lauf verlängern, und lässt die Abwicklung durch (committen, in das Austausch-Repo pushen, lesen, berichten).
`make zeitgrenze-gegenprobe` prüft den Hook. Der Hook ist **kein Gate**: Er steht in `harness/README.md` im
Abschnitt *Werkzeuge*, nicht in `GATE_CHECKS`, und prüft keinen Stand des Repos. Diese Kennung legt fest, was
er zusagt, wann er schweigt, woran er sperrt und was die Gegenprobe belegt. Der Grund der Wahl steht in
[ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md); diese
Kennung wiederholt ihn nicht.

**Stand.** Entscheidung des Nutzers vom 2026-10-21 (wörtlich, vom Koordinator übermittelt): „Wir haben rollen
(siehe .claude/agents) und commands -- max. 1 Stunde“. Dazu gehören aus derselben Entscheidung: Bei der Grenze hält der Lauf an,
committet, pusht in das Austausch-Repo und berichtet, ohne Abbruch und ohne Rücknahme von Arbeit; und sein „ok“ vom selben Tag auf den
Entwurf des Mechanismus: ein deterministischer `PreToolUse`-Hook (bash und awk, wie `pretooluse-command-guard.sh`), kein
Hintergrundprozess, Datenquelle die vorhandene Erfassung. [ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) ist **angenommen**: Entscheidung des
Nutzers vom 2026-10-22, Wortlaut „adr accept“. Damit gilt die Wahl der Mittel darunter: Sperrmenge, Fail-Richtung, Anker, Wanduhr, Wege
über die Grenze, Umgang mit den pflichtigen Läufen (Punkte 2 bis 10). **Geltungsbereich**, bestätigt in einer eigenen Entscheidung
des Nutzers vom 2026-10-22 (Wortlaut „2. ja“ auf die Frage der Koordination, ob der Wächter nur für Rollenläufe gilt und Commands und Hauptsitzung im
Regeltext als nicht erzwungen stehen): Der Hook deckt Rollenläufe; Commands und Hauptsitzung sind nicht erzwungen (Punkt 2, Grenze 12 (1)).
Ein Anker für sie wäre eine eigene Entscheidung (Folgepflicht der ADR). **Grenze bleibt, was noch kein Lauf belegt und kein Fall der
Gegenprobe zeigt** (Grenze 12 (7), Nachweis Punkt 13): dass die Eingabe von `PreToolUse` `agent_id` trägt (im Quelltext von Claude Code 2.1.306 gelesen,
nicht beobachtet), dass `decision: block` an dieser Stelle greift, ob und wo Claude Code stderr bei Exit 0 anzeigt und was bei überschrittenem
`timeout` geschieht. Die Zusagen unten gelten für den Hook, wie ihn die Gegenprobe zeigt, nicht für das Verhalten von Claude Code.

**Was er nicht ist.** Er beendet keinen Lauf, nimmt keine Arbeit zurück, schreibt keine Zeile in den Strom
und keine Datei (Leser, Punkt 1) und gilt nicht für die Hauptsitzung (Punkt 2, Grenze 12 (1)). Er ändert weder
`.claude/hooks/span-emit.sh` noch die Feldliste der Erfassung (`harness/erfassung-feldliste.md`), noch den
Stop-Hook, noch die Zeitgrenze je Gegenprobe ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16); die
Stunde ist eine **zusätzliche** Grenze über dem Lauf.

**Plattform** (`AGENTS.md` §3.12). Zugesagt und geprüft ist **Linux** (die Gegenprobe im Gegenproben-Container
und im Dev-Container, der Nachweis nach Punkt 13 an dem Host, auf dem der Agent läuft; bash ≥ 4, `awk` nach POSIX, GNU `date`, `env` und `timeout`). Der Hook läuft, wo Claude Code
läuft; die drei Betriebssysteme von [`LH-QA-03`](../../spec/lastenheft.md#lh-qa-03--portabilität) gelten dem
Produkt, nicht diesem Werkzeug. **macOS und Windows sind nicht zugesagt** (Grenze; Grund: kein solcher Host
steht zur Messung bereit). Plattformabhängig sind genau drei Dinge, und der Hook nutzt von jedem nur den
gemeinsamen Teil: die Gegenwart (`date -u +%s`, auf Linux, macOS und in Git-Bash vorhanden), die Umrechnung des
Ankers (Rechnung in `awk`, ohne `mktime` und ohne `date -d`, die BSD und mawk nicht kennen) und der awk-Dialekt
(kein gawk-Spezifikum, wie `tools/harness/extract-command.awk`). Wer den Hook dort laufen lässt, tut es auf
eigene Gefahr; was dort geschieht, ist nicht gemessen. Wo die Erfassung nichts schreibt (ein Host ohne den
Träger unter `.harness/state/bin/`), schweigt der Hook (Punkt 7).

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element einer
Aufzählung gibt es einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe
verlangt genau diese Menge roter Fälle je Mutation (mehr und weniger machen sie rot).

**Randform → Punkt.** Die Nummern der Fragen sind die der sechs Fragen in §6 des Slice-Plans
`slice-harness-rollenlauf-zeitgrenze`.

| Randform des Plans | Punkt |
|---|---|
| Frage 1 Matcher, Sperrmenge, Abwicklung, Erkennung am Befehl | 1, 5 |
| Frage 2 Fail-Richtung je Lage, Selbstauskunft | 7 |
| Frage 3 Befehl aus `tool_input.command`, Anker, Strom finden, Wanduhr, Zeitform, Zeitzone, Rückwärtssprung | 3, 5 |
| Frage 4 (a) Hauptsitzung, (b) langer Vordergrund-Befehl, (c) Erfassung nötig, (d) Plattform | 2, 12, 7, Kopf |
| Frage 5 Zeitgrenze des Hooks, leere Umgebung, Unterverzeichnis, Container und Host, Worktree | 8 |
| Frage 6 Weg über die Grenze | 9 |
| Reihenfolge mehrerer Gründe | 4 |
| Stelle der Aufrufe im Ablauf (Anfang, Mitte, ab der Grenze) | 3 (d), 4 |
| Aufruf-Voraussetzung (nacheinander, gleichzeitig) | 8 (5) |
| Pflichtige Läufe in der Stunde | 10 |
| Nachweis im Lauf | 13 |
| Gegenprobe, Zeit, Aufnahme | 14 |
| Lieferung in zwei Slices: Zuordnung der Punkte, allgemeine Freigabe, was nach der ersten Stufe allein zugesagt ist | Stufen (nach dieser Tabelle) |

**Stufen** (Entscheidung des Architect vom 2026-10-22, innerhalb von [ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md), auf den Vorschlag des Planners in §6 Randform 7 des Plans der Stufe 1,
mit den Änderungen unten; die ADR bleibt unverändert, die Stufung ordnet nur die Lieferung einer Entscheidung). Hook und Gegenprobe liegen zusammen über der Größenmarke eines Slice; geliefert wird in zwei Slices:
**Stufe 1** ist `slice-harness-rollenlauf-zeitgrenze` (der Wächter wirkt), **Stufe 2** ist `slice-harness-rollenlauf-zeitgrenze-haertung` (Selbstauskunft und Härtung), und Stufe 2 folgt auf Stufe 1. Jede Zusage steht in
genau einer Stufe, mit ihren Fällen und Mutationen; ein Punkt, der in beiden steht, ist geteilt, und die Tabelle nennt die Fälle je Stufe. Die Namen der Fälle tragen die Nummer des Punktes, den sie prüfen, auch wenn die Stufe
eine andere ist.

| Punkt | Stufe 1 (`slice-harness-rollenlauf-zeitgrenze`) | Stufe 2 (`slice-harness-rollenlauf-zeitgrenze-haertung`) |
|---|---|---|
| 1 Eintrag, Leser, Befehl | alle Fälle | — |
| 2 wer darunter fällt | alle Fälle, auch `p2-agent-in-tool-input` | — |
| 3 Anker, Strom, Zeit | (a) bis (d), (f), (g) mit allen Fällen `p3-*` außer denen rechts; von (e) die **Gestalt**: `p3-form-versatz`, `-bruchteil`, `-klein-z`, `-leerzeichen`, `-leer`, `-fehlt`, `-text-davor` | von (e) der **Bereich**: `p3-form-monat-13`, `-tag-32`, `-februar-30`, `-stunde-24`, `-minute-60`, `-sekunde-60` |
| 4 Reihenfolge der Gründe | — | alle Fälle |
| 5 Sperrmenge, Durchgelassenes | alle Fälle, auch `p5-unicode-in-tool-input`; `p5-bash-unlesbar` prüft die Freigabe | die Zeile `befehl-unlesbar` bei `p5-bash-unlesbar` |
| 6 Ausgabe und Exit | `p6-block-text`, `p6-frei-leer`, `p6-exit-0-block`, `p6-exit-0-frei`, `p6-exit-0-fehler` | `p6-exit-0-nicht-geprueft`, `p6-stderr-eine-zeile` |
| 7 Freigabe und Selbstauskunft | die **Freigabe** in sieben Lagen: `p7-kein-verzeichnis`, `p7-kein-strom`, `p7-strom-leer`, `p7-kein-start`, `p7-eingabe-leer`, `p7-eingabe-kein-json`, `p7-kein-awk` (je Exit 0, stdout leer, eine Mutation „Block statt Freigabe“) | die übrigen Fälle (`p7-eingabe-abgeschnitten`, `-eingabe-doppelt`, `-eingabe-unicode`, `-kennung-session-leer`, `-kennung-session-fehlt`, `-kennung-lang`, `-kennung-pfad`, `-anker-unlesbar`, `-kein-date`, `-befehl-unlesbar`) und die **Zeile auf stderr** in jeder Lage von Punkt 3, 5 und 7 |
| 8 Umgebung und Aufrufe | (1) Wurzel mit `p8-wurzel-projektdir`, `p8-wurzel-ort-der-datei`, `p8-unterverzeichnis`; (4) Worktree und Container (Prosa) | (2) `p8-leere-umgebung`, (3) `p8-gross`, (5) `p8-parallel` |
| 9 kein Weg über die Grenze | `p9-fortsetzung` | — |
| 10, 11, 12 pflichtige Läufe, Plattform, Grenzen | Prosa, kein Fall; gilt ab Stufe 1 | — |
| 13 Nachweis im Agentenlauf | die Handlung | nur, wenn Stufe 2 den Eintrag, den Pfad von `agent_id` oder die Ausgabe der Sperre ändert |
| 14 Gegenprobe, Zeit, Aufnahme | Gerüst, die Fälle von [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16 (9), `p14-aufnahme`, Aufnahme | Erweiterung um die Fälle von Stufe 2; die Zeit bleibt unter 600 s |

**Nach Stufe 1 allein ist zugesagt** (und nichts darüber hinaus):
(1) wer begrenzt ist und dass die Hauptsitzung still bleibt (Punkt 2); (2) Strom, Anker als letzter Start, Wanduhr und die Grenze von 3600 s mit `≥`, der Anker nur in der **Gestalt** `YYYY-MM-DDTHH:MM:SSZ`
(zwanzig Zeichen: vier, zwei und zwei Ziffern mit `-`, `T`, zwei, zwei und zwei Ziffern mit `:`, `Z`; Punkt 3 ohne den Bereich); (3) die Sperrmenge mit den Wörtern, der Wortgrenze und dem Durchgelassenen und Text und Form der Sperre
(Punkte 5 und 6, `p6-block-text`, `p6-frei-leer`); (4) die **Freigabe, mit Exit 0 und ohne Ausgabe auf stdout**, in genau diesen Lagen: Hauptsitzung; kein Verzeichnis; kein Strom (die Datei fehlt oder ist nicht lesbar);
Strom ohne Start oder leer; Anker in falscher Gestalt; Eingabe leer; Eingabe kein JSON; `awk` fehlt; ein Befehl, den `extract-command.awk` nicht liest (Unicode-Escape); ein Schlüssel `agent_id` unterhalb der obersten
Ebene; ein Unicode-Escape in einem anderen Schlüssel als `session_id`, `agent_id` und `tool_name`; über stderr ist in diesen Lagen nichts zugesagt, außer dass die Hauptsitzung dort leer bleibt; (5) die Wurzel nach Punkt 8 (1)
und aus einem Unterverzeichnis; (6) kein Weg über die Grenze (Punkt 9); (7) der Eintrag (Punkt 1), der Nachweis im Agentenlauf (Punkt 13) und die Gegenprobe samt Aufnahme (Punkt 14).
Die Freigabe in den Lagen von (4) ist der Fall, in dem der Hook den Wächter nach §3.1 **nicht** nachahmt: Jener sperrt bei Parse-Zweifel und ohne `awk`, dieser gibt frei (Punkt 7); deshalb gehören diese Lagen in Stufe 1 und nicht zur Härtung.
**Nicht zugesagt nach Stufe 1:** die Freigabe bei jeder anderen Unsicherheit (eine abgeschnittene Eingabe, ein doppelter Schlüssel, ein Unicode-Escape in `session_id`, `agent_id` oder `tool_name`, eine unzulässige Kennung,
`date` fehlt); ein Anker mit einem Wert außerhalb des Bereichs (weder Freigabe noch Sperre); die Reihenfolge der Gründe; jede Zeile auf stderr; das Verhalten bei leerer Umgebung, bei einer Eingabe bis 1 MB innerhalb des
`timeout` und bei gleichzeitigen Aufrufen; Exit 0 in einem Zweig, den kein Fall der Stufe 1 nennt. Die Sätze „Jede Unsicherheit des Hooks ist eine Freigabe“ (Punkt 7) und „Exit-Code 0 in jedem Zweig“ (Punkt 6) gelten
**uneingeschränkt erst nach Stufe 2**; bis dahin gelten sie für die Lagen in (4).
**Zwischen den Stufen:** (a) Der Hook ist nach Stufe 1 in `.claude/settings.json` eingetragen und wirkt; Stufe 2 ändert diese Zusagen nicht. (b) Ein Fall der Stufe 1 wird in Stufe 2 nicht gelockert; er darf eine
Erwartung an stderr dazubekommen (die Zeile der Stufe 2), und Stufe 2 fügt Fälle hinzu, ändert keinen Namen und keine erwartete Menge roter Fälle der Stufe 1. (c) Die Gegenprobe ist nach Stufe 1 aufgenommen (`SPEC-060`
Punkt 16) und bleibt es nach Stufe 2 unter ihrer Obergrenze. (d) Jeder Satz über den Hook (`harness/README.md`, Kommentar im Kopf, Handbuch, Regeltext) nennt nach Stufe 1 nur die Zusagen nach Stufe 1 (`AGENTS.md` §3.11): keine
Selbstauskunft auf stderr, keine allgemeine Freigabe. (e) Bliebe Stufe 2 aus, bliebe der Hook nach Stufe 1 mit diesen Zusagen und der Liste unter „Nicht zugesagt“ stehen; die Spalte Stufe 2 wäre dann nicht geliefert, und
der Vertrag sagt es so. (f) Randformen, die ein Fall der Stufe 2 aufwirft und die Stufe 1 berühren, geben beide Stufen an den Architect zurück; kein Implementer entscheidet sie im Code.

**Gemessen** am Bestand der Erfassung des Hauptbaums am 2026-10-22 (`.harness/state/spans/`, 854 Ströme,
39 569 Zeilen), damit die Punkte unten an der Wirklichkeit hängen und nicht an der Feldliste allein: Jede Zeile
beginnt `{"seq":<n>,"ts":"<ts>","event":"<Name>"`; die Ereignisse sind `PostToolUse` (38 037), `PostToolUseFailure`
(376) und `SubagentStart` (1156), ein Ereignis für das Ende eines Laufs gibt es nicht. 850 Ströme gehören einem
Agenten (Name `<Sitzung mit Unterstrichen>-<agent>.jsonl`, `agent` ist `a` und sechzehn Hexziffern in allen 850),
vier der Hauptsitzung (`<Sitzung>.jsonl`, `agent` leer). Alle 1156 Startzeilen sind streng geformt (`ts` als
`YYYY-MM-DDTHH:MM:SSZ`). 849 Agenten-Ströme beginnen mit `SubagentStart`; einer beginnt mit `PostToolUseFailure`
und hat nie einen Start. 132 Ströme tragen mehr als einen Start (307 spätere Starts, 285 davon unmittelbar nach
einer Zeile des Werkzeugs `SubagentHandback`): ein fortgesetzter Agent schreibt einen neuen Start in denselben
Strom. Vom **letzten** Start bis zur letzten Zeile des Stroms dauern 10 von 849 Läufen länger als eine Stunde,
vom **ersten** Start an 100 Ströme (längster: vier Tage). Bash-Aufrufe: 32 472, die längste Dauer 601 s, 88
Aufrufe mit mindestens 600 s (vermutlich eine Obergrenze des Werkzeugs; nicht zugesagt).
**Gelesen** im Quelltext der installierten Fassung von Claude Code (2.1.306, Zeichenketten des Binaries): Die
Eingabe von `PreToolUse` ist die Basis aller Hook-Ereignisse (`session_id`, `transcript_path`, `cwd`,
`permission_mode`, `agent_id`, `agent_type` …) mit `hook_event_name`, `tool_name`, `tool_input`, `tool_use_id`;
`agent_id` fehlt in der Hauptsitzung. Den Beleg im Agentenlauf trägt Punkt 13.

**Teil 1: der Hook.**

1. **Gegenstand, Eintrag, Leser.** *(Stufe 1)* Der Hook ist `.claude/hooks/pretooluse-zeitgrenze.sh`, aufgerufen wie die
   übrigen Hooks (`bash "$CLAUDE_PROJECT_DIR"/.claude/hooks/pretooluse-zeitgrenze.sh`). Er steht in
   `.claude/settings.json` unter `PreToolUse` in einem **eigenen** Eintrag mit `matcher` `""` (jedes Werkzeug)
   und `timeout` 5; der Eintrag des Wächters nach `AGENTS.md` §3.1 (`matcher` `Bash`) bleibt unberührt. Der Hook
   liest die Eingabe auf stdin und die Datei des Stroms und schreibt nichts außer stdout und stderr (Punkt 6);
   er legt weder Datei noch Marke an und ändert keine Zeile des Stroms. Er ist ein Skript mit genau einem
   Leser für die Felder der Eingabe (`tools/harness/hook-felder.awk`, ein Skript in POSIX-`awk`, das `session_id`,
   `agent_id` und `tool_name` von der obersten Ebene der Eingabe liest) und nimmt den Befehl aus `tool_input.command`
   über `tools/harness/extract-command.awk`, den Extraktor des Wächters nach §3.1, **nicht** aus dem Feld `program`
   des Stroms (das trägt nur das erste Wort: `cd … && make …` steht dort als `cd`). Fälle: `p1-eintrag` (der
   Arbeitsbaum trägt den Eintrag mit `matcher` `""`, `timeout` 5 und dem Befehl; der Eintrag für `Bash` ist
   unverändert vorhanden), `p1-schreibt-nichts` (die Dateiliste des Temp-Baums mit Inhalten ist nach jedem Fall von Punkt 3 bis 7 unverändert),
   `p1-befehl-nicht-aus-program` (ein Strom, dessen `program` `cd` ist, bei einem Befehl `cd x && make y` über der
   Grenze: gesperrt). Mutationen: das Matcher-Feld `Bash` statt `""`; `timeout` 60; der Hook legt eine Marke an (je eine für den Fall mit Block und den ohne); der Befehl wird aus dem
   Feld `program` gelesen.
2. **Wer darunter fällt.** *(Stufe 1)* Jeder Aufruf, dessen Eingabe `agent_id` trägt, also jeder Lauf eines Subagenten, **unabhängig
   von `agent_type`** und `agent_role` (auch ein Typ außerhalb der sechs Rollen, auch ein falsch geschriebener: leer
   heißt unbekannt, nie rollenlos, und die Grenze hängt nicht an der Schreibweise). Ein Aufruf **ohne** `agent_id` (die
   Hauptsitzung) wird **still** durchgelassen: keine Ausgabe auf stdout, keine auf stderr, auch wenn kein Strom
   liegt; die Hauptsitzung hat keinen `SubagentStart` und damit keinen Beginn (Grenze 12 (1), benannt, nicht
   verschwiegen). Ein `agent_id`, das kein String ist (`null`, Zahl, Objekt) oder leer, gilt als fehlend. Fälle:
   `p2-rolle-leer` (`agent_type` leer, Strom über der Grenze: gesperrt), `p2-typ-fremd` (`agent_type` `general-purpose`:
   gesperrt), `p2-hauptsitzung` (keine `agent_id`, Strom der Sitzung über der Grenze vorhanden: stdout und stderr
   leer), `p2-agent-null` (`"agent_id":null`), `p2-agent-leer` (`"agent_id":""`), `p2-agent-zahl` (`"agent_id":5`), `p2-agent-in-tool-input`
   (`tool_input` trägt den Text `"agent_id":"x"`, die oberste Ebene keinen: still wie die Hauptsitzung): je
   stdout und stderr leer. Mutationen: nur die sechs Rollen werden begrenzt (`agent_type` wird gelesen); die
   Hauptsitzung bekommt eine Zeile auf stderr; `null` gilt als Wert `null`; leer gilt als Wert; ein Schlüssel der zweiten Ebene zählt.
3. **Anker, Strom, Zeit.** *(Stufe 1; der Bereich von (e) Stufe 2)* (a) **Der Strom** des Aufrufs ist `<session_id mit _ statt ->-<agent_id>.jsonl` im Verzeichnis
   `<Wurzel>/.harness/state/spans/` (Wurzel: Punkt 8); ein Strom eines **anderen** Agenten derselben Sitzung oder
   derselben Kennung einer anderen Sitzung wird nie gelesen. (b) **Der Anker** ist die **letzte** Zeile des Stroms, die
   mit `{"seq":<n>,"ts":"<Text ohne Anführungszeichen>","event":"SubagentStart"` beginnt und mit `}` endet; eine Zeile, die ein Wert
   so enthält (ein `path` mit dem Text), zählt nicht, denn nur der Anfang der Zeile zählt. Der **letzte** Start, nicht der
   erste, weil ein fortgesetzter Agent in denselben Strom einen neuen Start schreibt und mit ihm eine neue Stunde
   beginnt (Punkt 9); die erste Zeile des Stroms ist kein Anker (in einem gemessenen Strom ist sie `PostToolUseFailure`, und er hat nie einen
   Start), die Änderungszeit der Datei auch nicht (sie ändert sich mit jedem Aufruf). (c) **Die Zeit** ist die Wanduhr: `jetzt − Anker`, `jetzt` aus
   `date -u +%s` (ganze Sekunden seit der Epoche, zeitzonenfrei), `Anker` aus `ts` gerechnet (`ts` ist UTC mit `Z`). Die
   Werkzeugzeit (Summe der `duration_ms`) zählt nicht: Sie ist ein Anteil der Wanduhr (48 % in dem Lauf von 4,34 h, der den
   Anlass gab), und die Stunde des Nutzers ist die, die der Koordinator wartet. (d) **Die Grenze** ist **3600 s, `≥`**:
   3599 s ist frei, 3600 s und mehr ist ab der Grenze (Punkt 5). Gemessen wird bei **jedem** Aufruf, gleich welchen
   Werkzeugs und an welcher Stelle des Laufs: Anfang, Mitte und Ende sind gleich; vor der Grenze ist jeder Aufruf frei.
   (e) **Zeitform des Ankers:** genau `YYYY-MM-DDTHH:MM:SSZ` mit Ziffern, Monat 01 bis 12, Tag 01 bis zum Monatsende des
   Jahres (Schaltjahr nach dem gregorianischen Kalender), Stunde 00 bis 23, Minute 00 bis 59, Sekunde 00 bis 59, Großbuchstaben
   `T` und `Z`; jede andere Form (Versatz statt `Z`, Bruchteil, `z`, Leerzeichen statt `T`, ein Wert außerhalb des Bereichs, leer,
   Text davor oder danach) ist `anker-unlesbar` (Punkt 7); trägt eine Startzeile kein `ts`, ist sie keine Startzeile nach (b) und gilt nicht (`kein-anker`, wenn keine andere da ist). **Gestalt und Bereich** sind getrennt zugesagt: Die Gestalt (zwanzig Zeichen, Ziffern, `-`, `T`, `:`, `Z` an ihren Stellen) prüft Stufe 1, den Bereich (Monat, Tag, Stunde, Minute, Sekunde) Stufe 2; nach Stufe 1 ist für einen Anker in richtiger Gestalt und falschem Bereich weder Freigabe noch Sperre zugesagt. (f) **Rückwärtssprung der Uhr:** liegt der Anker nach
   `jetzt` (Zeit negativ), ist der Aufruf frei. Ein Sprung vorwärts sperrt zu früh (Grenze 12 (6)). (g) **Mehrere Starts,
   unlesbarer letzter Start:** Ist der letzte passende Start unlesbar (Form nach (e)), gilt `anker-unlesbar`; ein früherer Start
   wird dann **nicht** ersatzweise genommen. Eine abgeschnittene letzte Zeile (ohne `}`) ist kein Start, und der vorige gilt.
   Fälle: `p3-unter` (Anker 3599 s vor `jetzt`: frei), `p3-grenze` (3600 s: gesperrt), `p3-ueber` (3601 s: gesperrt),
   `p3-negativ` (Anker 100 s nach `jetzt`: frei), `p3-letzter-start` (zwei Starts, der erste 2 h, der letzte 60 s vor `jetzt`: frei),
   `p3-start-nicht-erste-zeile` (erste Zeile `PostToolUseFailure`, danach ein Start 2 h vor `jetzt`: gesperrt),
   `p3-wert-mit-text` (eine Zeile `PostToolUse`, deren `path` den Text `\"event\":\"SubagentStart\"` mit altem `ts` trägt, danach ein junger Start:
   frei), `p3-start-abgeschnitten` (die letzte Zeile ist ein abgeschnittener Start, davor ein junger: frei; davor ein alter: gesperrt),
   `p3-anderer-agent` (zwei Ströme derselben Sitzung, der eine über der Grenze: der Aufruf des jungen ist frei, der des alten gesperrt),
   `p3-andere-sitzung` (dieselbe `agent_id` in zwei Sitzungen: je nach Sitzung), `p3-unterstrich` (Sitzung mit `-` findet den Strom mit `_`),
   `p3-zeitzone-kalkutta` und `p3-zeitzone-los-angeles` (echte Uhr, `TZ=Asia/Kolkata` bzw. `TZ=America/Los_Angeles`, Anker 60 s bzw. 7200 s
   vor `jetzt`: frei bzw. gesperrt, in beiden Zonen), `p3-schaltjahr` (Anker `2028-02-29T23:59:59Z`, `jetzt` 3601 s später: gesperrt),
   `p3-monatsende` (Anker `2026-12-31T23:30:00Z`, `jetzt` 3600 s später: gesperrt), `p3-start-vorletzter-lesbar` (letzter Start unlesbar, davor ein junger:
   `anker-unlesbar`), dazu die Formen von (e), je eine: `p3-form-versatz`, `p3-form-bruchteil`, `p3-form-klein-z`, `p3-form-leerzeichen`,
   `p3-form-monat-13`, `p3-form-tag-32`, `p3-form-februar-30`, `p3-form-stunde-24`, `p3-form-minute-60`, `p3-form-sekunde-60`, `p3-form-leer`,
   `p3-form-fehlt` (Startzeile ohne `ts`: `kein-anker`), `p3-form-text-davor`. Die Gegenwart liefert in diesen Fällen eine Attrappe von `date` (ein Verzeichnis vor `PATH`, das für
   jeden Aufruf, der mit `+%s` endet, eine feste Zahl ausgibt; die Fälle der Zeitzonen laufen mit dem echten `date`). Mutationen: der erste statt des letzten Starts; nur die erste Zeile
   wird gelesen; eine Teilzeichenkette `"event":"SubagentStart"` irgendwo in der Zeile zählt; die Änderungszeit der Datei als Anker; `>` statt `≥`
   (färbt `p3-grenze`); `ts` als Ortszeit gelesen (färbt beide Fälle der Zone); der Betrag der Zeit statt der Zeit (`p3-negativ`); der Strom ohne `_` gesucht; alle Ströme der Sitzung
   gelesen; der erste lesbare Start vor dem unlesbaren genommen; eine abgeschnittene Zeile zählt; je Form von (e) eine Mutation, die genau sie annimmt;
   Schaltjahr und Monatsende ignoriert (je eine).
4. **Reihenfolge der Gründe.** *(Stufe 2)* Der Hook prüft in dieser Reihenfolge und endet beim ersten Treffer: (1) `awk` fehlt: frei, Grund `kein-awk`;
   (2) Eingabe nicht lesbar: frei, `eingabe-unlesbar`; (3) kein `agent_id`: frei, still (Punkt 2); (4) Kennung unzulässig (Punkt 7): frei,
   `kennung-unzulaessig`; (5) Verzeichnis fehlt: `kein-verzeichnis`; (6) Strom fehlt: `kein-strom`; (7) kein lesbarer Anker: `kein-anker` (kein Start
   in der Datei, auch eine leere Datei) oder `anker-unlesbar`; (8) `date` fehlt: `kein-date`; (9) Zeit unter der Grenze: frei, still; (10) ab der Grenze:
   Punkt 5. Fälle (Paare benachbarter Gründe): `p4-awk-vor-eingabe`, `p4-eingabe-vor-agent`, `p4-agent-vor-kennung`, `p4-kennung-vor-verzeichnis`,
   `p4-verzeichnis-vor-strom`, `p4-strom-vor-anker`, `p4-anker-vor-date`, `p4-date-vor-zeit` (ohne `date` und mit Anker über der Grenze: frei mit `kein-date`).
   Mutationen: je Paar vertauscht (acht).

**Sperrmenge, Ausgabe, Lagen ohne Prüfung.**

5. **Was ab der Grenze gesperrt wird und was durchgelassen.** *(Stufe 1; die Zeile `befehl-unlesbar` Stufe 2)* Gesperrt sind (a) die Werkzeuge `Write`, `Edit`, `NotebookEdit` und `Agent` (je Name
   ein Element) und (b) Aufrufe des Werkzeugs `Bash`, deren `tool_input.command` eines der **Wörter** `make`, `docker` und `sleep` enthält (je Wort ein
   Element). Ein Wort ist die Teilkette, vor der Anfang oder ein Zeichen **außerhalb** von `A-Za-z0-9_.-` steht und nach der Ende oder ein solches Zeichen folgt;
   die Prüfung ist **blind für Anführungszeichen**: Ein Wort in einer Commit-Message, in einem Argument von `ssh`, in einem Heredoc sperrt den Aufruf ebenso (der
   Text der Sperre sagt es, Punkt 6). **Durchgelassen** ist alles andere, namentlich: `Read`, `Grep`, `Glob`, `SubagentHandback` und jedes Werkzeug, das die
   Sperrmenge nicht nennt (auch eines, das es heute nicht gibt); `Bash` ohne die drei Wörter, darunter `git add`, `git commit`, `git push`, `git status`,
   `git diff`, `git log` und `git mv`, auch über `ssh` in den Arbeitsklon eines Dev-Containers (der Hook liest den Befehl, nicht, was er ruft). Eine Rolle muss bei
   der Grenze committen, in das Austausch-Repo pushen und berichten können; eine zu enge Freigabe nähme ihr den Bericht, eine zu weite ließe die Arbeit weiterlaufen,
   und der erste Fehler wiegt schwerer (Grund der kleinen Sperrmenge, [ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md)). Ein `Bash`-Aufruf, dessen Befehl der Hook nicht lesen kann (Parse-Zweifel von
   `extract-command.awk`, auch ein `\u`-Escape im Befehl), ist **frei** (Punkt 7) und wird nicht gesperrt. Stolperdraht, keine Sandbox (Grenze 12 (5)).
   Fälle, je Element: `p5-sperre-write`, `p5-sperre-edit`, `p5-sperre-notebookedit`, `p5-sperre-agent`; `p5-sperre-make`, `p5-sperre-docker`, `p5-sperre-sleep`;
   Wortgrenze: `p5-pfad-make` (`/usr/bin/make x`: gesperrt), `p5-wort-in-message` (`git commit -m "make gates gruen"`: gesperrt), `p5-wort-in-ssh` (`ssh h 'cd x && make y'`
   und `ssh h 'make y'`: je gesperrt), `p5-zeilenumbruch` (`git status` und in der nächsten Zeile `make y`: gesperrt), `p5-strichwort` (`make-x`, `docker-compose.yml`, `sleep.sh`,
   `gnumake`, `usleep`, `Makefile`: je frei), `p5-punkt-wort` (`git add make.mk` und `git add docker.mk`: je frei); durchgelassen: `p5-frei-read`, `p5-frei-grep`, `p5-frei-glob`,
   `p5-frei-handback`, `p5-frei-unbekannt` (`mcp__x__y`), `p5-frei-git-commit`, `p5-frei-git-push`, `p5-frei-git-status`, `p5-frei-git-diff`, `p5-frei-git-log`,
   `p5-frei-git-mv`, `p5-frei-git-add`, `p5-frei-ssh-git` (`ssh h 'git commit -m x'`); `p5-bash-unlesbar` (ein Unicode-Escape im Befehl: frei; ab Stufe 2 mit der Zeile `befehl-unlesbar`, nicht `eingabe-unlesbar`), `p5-unicode-in-tool-input` (`Write` über der Grenze, ein Unicode-Escape in `content`: gesperrt). Mutationen: je Werkzeug und je Wort ein Element der Sperrmenge gestrichen (sieben); je durchgelassenem Namen ein Element, das es sperrt (`Read`,
   `Grep`, `Glob`, `SubagentHandback`, `git commit`, `git push`, `git status`, `git diff`, `git log`, `git mv`, `git add`, der unbekannte Name); die Wortgrenze ohne `-`
   (färbt `p5-strichwort`), ohne `.` (färbt `p5-strichwort` und `p5-punkt-wort`); die Wortgrenze verlangt Leerraum davor (färbt `p5-pfad-make`, `p5-wort-in-message`, `p5-wort-in-ssh`); die Prüfung
   liest nur die erste Zeile des Befehls (färbt `p5-zeilenumbruch`); ein Wort in Anführungszeichen wird ausgenommen (färbt `p5-wort-in-message`, `p5-wort-in-ssh`); ein Unicode-Escape irgendwo in der Eingabe gilt als Zweifel (färbt `p5-unicode-in-tool-input`); ein unlesbarer Befehl sperrt (färbt `p5-bash-unlesbar`).
6. **Ausgabe und Exit-Code.** *(Stufe 1; die Zeile auf stderr Stufe 2)* Im Fall „frei“ gibt der Hook **nichts** auf stdout aus (wie beim Wächter nach §3.1; was Claude Code mit einer leeren Ausgabe tut, ist nicht zugesagt, Grenze 12 (7)). Im Fall „gesperrt“ gibt er auf stdout **genau eine Zeile** aus, ein JSON-Objekt in der Form des Wächters nach §3.1 (`{"decision":"block","reason":"<Text>"}`), und `<Text>` ist
   dieser, Zeichen für Zeichen, auf einer Zeile:

   ```
   Zeitgrenze von einer Stunde erreicht. Halte an: committe den Stand, pushe ihn in das Austausch-Repo und berichte Stand und Offenes mit SubagentHandback. Nimm keine Arbeit zurück. Frei bleiben Read, Grep, Glob, SubagentHandback und Bash ohne die Wörter make, docker und sleep (auch git commit und git push, auch über ssh). Gesperrt sind Write, Edit, NotebookEdit, Agent und Bash, das make, docker oder sleep nennt, auch in einer Commit-Message.
   ```

   Im Fall „nicht geprüft“ (Punkt 7) steht auf stderr **genau eine Zeile** `zeitgrenze: nicht geprüft (<Grund>)` und nichts auf stdout; sonst schreibt der Hook nichts auf stderr. Der Exit-Code
   ist **0 in jedem Zweig** (uneingeschränkt erst nach Stufe 2, siehe *Stufen*), auch bei einem Fehler im Skript (ein Skript, das mit 2 endete, sperrte; ein anderer Code zeigte eine Fehlermeldung an; beides widerspricht der Zusage, dass jede
   Unsicherheit eine Freigabe ist). Der Hook hat kein `set -e`. Fälle: `p6-block-text` (stdout ist die Zeile mit dem Text, byte-gleich, danach ein Zeilenende; stderr leer), `p6-frei-leer`
   (stdout und stderr leer, Exit 0), `p6-exit-0-block`, `p6-exit-0-frei`, `p6-exit-0-nicht-geprueft` (Exit 0), `p6-exit-0-fehler` (eine Datei des Stroms ohne Leserecht: Exit 0 und `kein-strom`; entfällt
   als root, die Gegenprobe nennt es), `p6-stderr-eine-zeile` (je Grund von Punkt 7 genau eine Zeile in der Form). Mutationen: stdout trägt ein Zeichen mehr oder weniger im Text (je Satz
   der drei Sätze des Textes einer); das Objekt trägt `approve`; stderr nennt den Block; Exit 2 bei Block; Exit 1 bei „nicht geprüft“; `set -e`; zwei Zeilen auf stderr; stdout trägt etwas im
   Fall „nicht geprüft“.
7. **Fail-open und Selbstauskunft.** *(Stufe 1 für sieben Lagen der Freigabe, sonst Stufe 2)* **Jede Unsicherheit des Hooks ist eine Freigabe** (uneingeschränkt erst nach Stufe 2; nach Stufe 1 für die Lagen unter *Stufen*), mit einer Zeile auf stderr, die den Grund nennt; ausgenommen ist die Hauptsitzung (Punkt 2, still). Die Gründe
   sind abschließend: `kein-awk` (Programm nicht auffindbar), `eingabe-unlesbar` (leer, kein JSON, abgeschnitten, doppelter Schlüssel `session_id`, `agent_id` oder `tool_name` auf der obersten
   Ebene, ein Unicode-Escape in einem der drei Werte), `kennung-unzulaessig` (`session_id` leer oder fehlend, oder `session_id` oder `agent_id` länger als 128 Zeichen oder mit einem Zeichen außerhalb von
   `A-Za-z0-9_-`: Schutz vor einem Pfad aus der Eingabe), `kein-verzeichnis`, `kein-strom` (die Datei fehlt oder ist nicht lesbar), `kein-anker` (kein Start in der Datei, auch eine leere Datei),
   `anker-unlesbar`, `kein-date`, `befehl-unlesbar` (nur bei `Bash` ab der Grenze, Punkt 5). Ein Unicode-Escape im Wert eines **anderen** Schlüssels (ein `Write` mit einem Steuerzeichen in `content`)
   ist kein Zweifel (`p5-unicode-in-tool-input`), und Schlüssel **unterhalb** der obersten Ebene (`tool_input.agent_id`) zählen nicht (`p2-agent-in-tool-input`). **Die Richtung ist die umgekehrte des Wächters nach §3.1**, mit Grund: Der Wächter nach
   §3.1 braucht nur den Text des Befehls, der immer da ist, und ein verfehlter Block ist ein installierter Paketmanager auf dem Host, der nicht zurückzunehmen ist; dieser Hook braucht
   gitignorierten Zustand, der in einem frischen Klon, nach einem Aufräumen der Erfassung, auf einem Host ohne Träger und in einem Worktree fehlen kann, und ein Block aus Zweifel hielte dort **jede
   Arbeit jedes Subagenten** an, still und ohne dass jemand die Ursache sähe; ein verfehlter Block heißt nur, dass diese Grenze für diesen Lauf nicht greift, und die Regel an den Leser bleibt.
   Dass der Hook seine Leere nennt, ist die Gegenmaßnahme dafür, dass er dort schweigt, wo er gebraucht wird; ob und wo Claude Code stderr bei Exit 0 anzeigt, ist nicht zugesagt (Grenze 12 (7)).
   Fälle, je Grund einer oder mehrere: `p7-kein-awk` (`PATH` ohne `awk`), `p7-eingabe-leer`, `p7-eingabe-kein-json`, `p7-eingabe-abgeschnitten`, `p7-eingabe-doppelt`, `p7-eingabe-unicode`
   (ein Unicode-Escape in `agent_id`), `p7-kennung-session-leer`, `p7-kennung-session-fehlt`, `p7-kennung-lang` (129 Zeichen, je für `session_id` und `agent_id`), `p7-kennung-pfad` (`agent_id` `../x`),
   `p7-kein-verzeichnis`, `p7-kein-strom`, `p7-strom-leer` und `p7-kein-start` (Ströme ohne Startzeile: nur `PostToolUse`-Zeilen; beide `kein-anker`), `p7-anker-unlesbar` (die Fälle von Punkt 3 (e)
   tragen den Grund), `p7-kein-date` (`PATH` ohne `date`), `p7-befehl-unlesbar` (`Bash` über der Grenze, ein Unicode-Escape im Befehl: frei, Grund `befehl-unlesbar`, nicht `eingabe-unlesbar`),
   Je Fall wird geprüft: Exit 0, stdout leer, **ab Stufe 2** stderr genau die Zeile (in Stufe 1 nur Exit 0 und stdout). Mutationen: je Grund ein Block statt der Freigabe (neun; für `befehl-unlesbar` gesperrt statt frei); die Zeile entfällt; der Grund wird vertauscht; ein doppelter Schlüssel nimmt den ersten Wert; die Kennung wird
   ungeprüft in den Pfad gesetzt (färbt `p7-kennung-pfad` mit einem Strom im übergeordneten Verzeichnis).
8. **Umgebung und Aufrufe des Treibers.** *(Stufe 1: (1) und (4); Stufe 2: (2), (3), (5))* (1) **Wurzel** ist `$CLAUDE_PROJECT_DIR`; ist sie leer oder nicht gesetzt, die Wurzel aus dem Ort der Datei (`<wurzel>/.claude/hooks/…`), wie
   `.claude/hooks/span-emit.sh`: der Hook liest dort, wo der Emitter schreibt. Das Arbeitsverzeichnis des Aufrufs ändert nichts (Aufruf aus einem Unterverzeichnis, in einem Worktree, in einem anderen Verzeichnis).
   (2) **Leere Umgebung:** mit `env -i`, einem gesetzten `CLAUDE_PROJECT_DIR` und ohne `PATH` findet der Hook seine Werkzeuge über den Standardpfad von bash (Linux: `/usr/bin`) und urteilt wie sonst; der Fall läuft mit der echten Uhr. (3) **Zeitgrenze des Hooks:**
   Der Eintrag nennt `timeout` 5; das Skript endet bei einer Eingabe bis 1 MB (ein `Write` mit großem `content`) in unter 5 s (gemessen mit `extract-command.awk`: 0,15 s für 1 MB, gawk 5.2.1); was Claude Code bei
   Überschreitung tut, ist nicht zugesagt (Grenze 12 (7)). Die Zeit wächst mit dem Strom, ohne dass über seine Größe etwas zugesagt ist; gemessen: Das Suchen des Ankers in einem Strom von 5 MB (11 166 Zeilen, das
   Sechsfache des größten Stroms des Bestands) dauert 21 ms. (4) **Worktree und Container:** Die Erfassung ordnet `slice`, `branch` und `commit` dem Zustand des Hauptbaums zu
   (`BEO-REPO/span-zuordnung-folgt-dem-hauptbaum-nicht-dem-klon`); der Hook braucht nur den **Ort** des Stroms und den **Beginn**, beides hängt nicht an diesen Feldern, und er ist von dieser Beobachtung nicht
   betroffen. Läuft Claude Code im Hauptbaum und arbeitet der Agent in einem Worktree oder im Klon eines Dev-Containers (Aufrufe über `ssh`), lesen Hook und Emitter dieselbe Wurzel nach (1); wo Claude Code `CLAUDE_PROJECT_DIR` setzt, ist nicht zugesagt. Läuft eine Sitzung **in** einem
   Worktree ohne `.harness/state/bin/` (gitignored), schreibt der Emitter nichts, und der Hook schweigt mit `kein-verzeichnis` oder `kein-strom`. Läuft Claude Code im Dev-Container, gilt dieselbe Regel (1) für den
   Container; geprüft ist das nicht. Geprüft ist keiner dieser Orte durch die Gegenprobe (sie läuft ohne Claude Code); den Beleg im Agentenlauf trägt Punkt 13, und nur für den Ort, an dem er läuft. (5) **Aufruf-Voraussetzung:** Der Hook hat keinen Zustand und sagt keine Ordnung über
   mehrere Aufrufe zu; Aufrufe einer Session nacheinander und gleichzeitig (parallele Werkzeug-Aufrufe eines Agenten) sind zugesagt, weil er nur liest; ein Start, den der Emitter in dem Augenblick schreibt, in dem der Hook liest, ist
   eine abgeschnittene Zeile (Punkt 3 (g)), und eine Entscheidung auf dem vorigen Start ist die Folge (akzeptiertes Negativ; Grund: ein Fenster von Mikrosekunden, und der Fall trifft nur den ersten Aufruf nach einer Fortsetzung).
   Fälle: `p8-wurzel-projektdir` (Strom nur unter `$CLAUDE_PROJECT_DIR`, anderes Arbeitsverzeichnis: wie Punkt 3), `p8-wurzel-ort-der-datei` (`CLAUDE_PROJECT_DIR` leer, Hook als `<wurzel>/.claude/hooks/…` im Temp-Baum: gesperrt),
   `p8-unterverzeichnis`, `p8-leere-umgebung` (`env -i CLAUDE_PROJECT_DIR=… bash …`), `p8-gross` (1 MB `content` bei `Write` über der Grenze unter `timeout 5`: gesperrt, Exit 0; unter der Grenze: frei), `p8-parallel` (vier Aufrufe
   gleichzeitig, je ihr Ergebnis, nach ihnen keine Datei). Mutationen: die Wurzel aus dem Arbeitsverzeichnis; der Ort der Datei wird nicht als Rückfall genommen; der Pfad relativ zum Arbeitsverzeichnis; das Skript wartet 6 s;
   der Hook legt eine Sperrdatei an (färbt `p8-parallel` und `p1-schreibt-nichts`).
9. **Kein Weg über die Grenze im Hook.** *(Stufe 1)* Es gibt keine Umgebungsvariable, keine Datei und keinen Satz im Auftrag, die die Grenze heben; ein Weg wäre aus der Grenze eine Bitte (Entscheidung des Nutzers vom 2026-10-21:
   „max. 1 Stunde“). Der Weg ist, was die Entscheidung schon vorsieht: **Der Koordinator entscheidet nach dem Bericht über einen neuen Lauf.** Setzt er den angehaltenen Agenten fort, schreibt die Erfassung einen neuen
   `SubagentStart` in denselben Strom (gemessen, Kopf: 307 spätere Starts), und mit ihm beginnt nach Punkt 3 (b) eine neue Stunde; die Rolle selbst hebt ihre Grenze nicht. Dass die Datei des Stroms nicht
   zugriffsbeschränkt ist, die Rolle sie also verändern könnte, ist Grenze 12 (5). Fall `p9-fortsetzung` (ein Strom mit zwei Starts, der zweite 60 s vor `jetzt`: frei, auch wenn der erste 5 h zurückliegt; derselbe Strom
   mit dem zweiten Start 2 h zurück: gesperrt); Mutation: der erste Start zählt (färbt beide).
10. **Die Stunde und die pflichtigen Läufe.** *(Prosa, gilt ab Stufe 1)* Die Stunde trägt, was jede Rolle **muss**, sofern die Arbeit selbst darunter bleibt; für zwei Fälle trägt sie es nicht, (c) und (d), und dafür gilt je eine Ausnahme.
    (a) `make gates` (einmal am Endstand): 265 bis 303 s auf ruhigem Host, einmal 9 min 53 s (Messwerte der Pläne in `done/` und `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`). (b) Der volle Lauf einer
    **aufgenommenen** Gegenprobe ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16, `container-einrichten-gegenprobe`, diese Gegenprobe): höchstens 1200 s ([ADR-0059](../../docs/plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md), vorher 600 s nach [ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md)).
    Beide zusammen belegen im Regelfall bis zu 25 min (rund 300 s und bis zu 1200 s), im ungünstigen Fall rund 30 min (600 s und 1200 s); es bleiben mindestens 30 min für die Arbeit. **Das trägt.** (c) Der volle Lauf einer Gegenprobe im Zustand `nicht gemessen` (heute `release-gegenprobe` mit 25 bis 30 min auf ruhigem Host
    und Stunden unter Last, `test-container-gegenprobe`, `container-start-gegenprobe` mit 19 min 38 s am Host): **das trägt nicht neben der Arbeit, und kein Rollenlauf fährt ihn** (Entscheidung des Nutzers vom 2026-10-22, Wortlaut
    „agents.md anpassen“, vom Koordinator übermittelt; sie ändert den Übergang in [ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md) und die Ausnahme von `AGENTS.md` §6 Punkt 6). Entscheidung des Nutzers vom 2026-10-25 (ändert die Fassung vom 2026-10-22): Implementer und Verifier fahren den vollen Lauf einer aufgenommenen Gegenprobe im Hintergrund nach [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 17, nie im Vordergrund; den engsten Fall (`GP_MUT=<Muster>`) fahren sie während der Arbeit. Aufgenommen sind `release-gegenprobe` und `test-container-gegenprobe` mit `slice-harness-gegenproben-schneller`; `make gegenproben` am Ende einer Welle (d). Ein Bruch außerhalb des engsten Falls
    bleibt bis dahin unbemerkt (akzeptiertes Negativ; Grund: der volle Lauf trägt nicht neben der Arbeit; `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`). Für eine nicht aufgenommene Gegenprobe ist der engste Fall die Fahrt nach [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 14 (Entscheidung des Nutzers vom 2026-10-22, Wortlaut „ok“ auf den Vorschlag der Koordination; der Punkt überlässt Rolle und Zeitpunkt den Anweisungen der Rollen), und die Frist trägt die Tabelle mit ihrem Nehmer.
    (d) `make gegenproben` am Ende einer Welle (18 min 15 s für die Liste des Containers, Messwert in `done/`) fährt die Einlieferung in das Austausch-Repo auf dem eingelieferten Hash
    ([`SPEC-061`](SPEC-061-gegenproben-nach-einlieferung.md)), `make gegenproben-host` der Koordinator am Host in der Hauptsitzung; **keine Rolle wartet darauf innerhalb ihrer Stunde.** Der Hook sperrt ein `make`, das **vor** der Grenze begann, nicht
    (Grenze 12 (2)); ein Lauf im Hintergrund läuft über die Grenze weiter. Kein Fall (Grund: der Punkt sagt eine Zuordnung von Läufen zu Rollen und eine Rechnung zu, keine Eigenschaft des Hooks, und ein Fall dafür wäre ein
    Gate, das diese Entscheidung verlässt); der Beleg sind die Messwerte, die die Anweisungen der Rollen tragen.

**Teil 2: Nachweis, Grenzen, Gegenprobe.**

11. **Plattform.** *(Prosa, gilt ab Stufe 1)* Siehe den Kopf dieser Kennung. Je Plattform von `LH-QA-03`: Linux geprüft, macOS und Windows Grenze (keine Zusage, kein Fall). Ein Satz in einem Handbuch über die Zeitgrenze nennt nur Linux (§3.11).
12. **Grenzen.** *(Prosa, gilt ab Stufe 1)* (1) **Die Hauptsitzung** hat keinen `SubagentStart` und damit keinen Beginn: Ein Command, das in ihr läuft (`.claude/commands/`), ist nicht begrenzt. Das Budget des Nutzers nennt „rollen und commands“;
    diese Kennung deckt die Rollen. Ein Anker für die Hauptsitzung ist eine eigene Entscheidung ([ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md), Folgepflicht), nicht ein Zusatz dieses Hooks. (2) **Innerhalb eines laufenden Vordergrund-Befehls** feuert kein
    Hook; die Zusage lautet darum „beim ersten Aufruf ab der Grenze“, nicht „nach einer Stunde“. Wie lange die Überschreitung dauern kann, ist nicht zugesagt (gemessen: längster Bash-Aufruf 601 s von 32 472; ein Zeitlimit des Werkzeugs ist dessen Sache), ein Befehl im
    Hintergrund läuft weiter; der Hook beendet keinen Prozess (Abbruch ist nicht gewollt: ein Wächter, der einen Lauf tötet, verlöre den Bericht). (3) **Der Hook wirkt nur, wo die Erfassung liegt** (Punkt 7, 8). (4) **Denkzeit:**
    Zeit, in der ein Agent kein Werkzeug aufruft, feuert keinen Hook; ein Lauf, der lange nur denkt, bemerkt die Grenze beim nächsten Aufruf. (5) **Stolperdraht, keine Sandbox:** Ein Befehl, der das Wort umgeht (`m""ake`, ein Skript,
    `bash -c`, `ssh`-Skript, das `make` erst darin nennt), und das Werkzeug `Monitor` werden nicht erkannt; der Strom liegt gitignored und nicht zugriffsbeschränkt (`harness/erfassung-feldliste.md`), und wer ihn
    ändert, ändert seinen Anker. Die Grenze ist eine Bremse für einen Agenten, der sie lesen und befolgen soll, nicht eine für einen, der sie umgehen will. (6) **Ein Sprung der Uhr nach vorn** sperrt zu früh; die Zeit ist die
    Zeit des Hosts. (7) **Verhalten von Claude Code** (gemessen ist allein die Fassung 2.1.306, jede andere ist nicht zugesagt und der Beleg nach Punkt 13 nennt die Fassung): die Form `decision: block` an dieser Stelle, die Anzeige von stderr bei Exit 0 und das Verhalten bei überschrittenem `timeout` sind nicht zugesagt; das Erste belegt
    Punkt 13 an der installierten Fassung, die beiden anderen sind aus dem Klon nicht prüfbar. (8) **Der Hook ist im Strom nicht sichtbar:** Ein gesperrter Aufruf erzeugt keine Zeile, und der Hook schreibt keine; dass er
    sperrte, steht im Bericht der Rolle. (9) **Zeit des Modells:** Die Stunde misst die Wanduhr, nicht Kosten oder Token.
13. **Nachweis im Agentenlauf — Handlung, kein Sensor.** *(Stufe 1)* Was die Gegenprobe nicht belegt (Claude Code ruft den Hook mit `agent_id` auf, die Form des Blocks greift, der Eintrag in `.claude/settings.json` greift),
    belegt **einmal je Lieferung** der Implementer in seinem eigenen Lauf. Schritte: (1) Er legt eine Datei `zeitgrenze-marke-<Zahl>` an und liest sie mit `Read`; seinen Strom erkennt er an der Zeile `Read` mit diesem Pfad
    (die Datei `<Sitzung mit _>-<agent>.jsonl` im Verzeichnis der Ströme). (2) Er sichert den Strom mit `cp` und nennt die Summe nach `sha256sum`. (3) Er setzt in der letzten Startzeile `ts` mit `sed` um zwei Stunden zurück.
    (4) Er ruft `Bash` mit dem Text `echo make` auf (gesperrt) und liest die Meldung im Ergebnis des Werkzeugs, dann `Write` auf eine Datei im Temp-Ordner (gesperrt), dann `Bash` mit `git status` (durchgelassen).
    (5) Er stellt den Strom aus der Sicherung wieder her (`cp`, durchgelassen); die Summe nach `sha256sum` ist die aus (2). **Beleg:** eine eigene Datei unter `docs/reviews/` (Lauf-Beleg, von der Größenzählung ausgenommen,
    [ADR-0056](../../docs/plan/adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)), die den Kopf (Datum, Stand, Plattform, Fassung von Claude Code), die Befehle und ihre Ausgabe **ungekürzt** trägt (gekürzt wird nur ein
    Geheimnis oder Token, mit Vermerk); §7 des Plans nennt den Pfad der Datei, ihre Summe nach `sha256sum`, die Zeile mit dem Exit-Code und den Kopf. Zeigt der Beleg, dass die Eingabe von `PreToolUse` **kein** `agent_id` trägt
    oder die Form des Blocks nicht greift, hält der Implementer an und gibt es dem Architect zurück: Die Wahl des Ankers in Punkt 3 (a) und die Form in Punkt 6 fallen, es folgt eine neue Entscheidung. Der Beleg ist ein Lauf der
    Umgebung des Agenten; ein grüner Lauf der Gegenprobe belegt ihn nicht (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`).
14. **Gegenprobe: Aufbau, Zeit, Aufnahme.** *(Stufe 1, in Stufe 2 erweitert)* `zeitgrenze-gegenprobe` ist `.PHONY` und steht in `GEGENPROBEN_CONTAINER` (bash, awk, `env`, `timeout`; ohne Docker und ohne Claude Code), nicht in `GATE_CHECKS` und nicht in der
    Ausnahmeliste der schnellen Kette ([ADR-0048](../../docs/plan/adr/0048-gegenproben-verlassen-die-gate-kette.md); die Zuordnung prüft [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 2, ein eigener Fall dafür entfällt); ihr Fragment
    ist `harness/mk/zeitgrenze.mk`, ihr Skript `tools/dev/harness/zeitgrenze-gegenprobe.sh`. Sie fährt **Attrappen-Ströme** (JSONL in einem Temp-Ordner je Fall, kein Agentenlauf, keine echte Erfassung) und eine Attrappe von `date`;
    **kein Fall wartet** (die Zeit kommt aus dem Anker und der Attrappe, nicht aus `sleep`; `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`). Mutationen sind Textersatz in einer Kopie des Hooks, des Lesers der Felder und der
    Einträge. Die unveränderte Kopie ist grün. Die Fälle von [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16 (9) (Zeitbudget, Zeile mit Zeit und Last, `GP_BUDGET`, Arbeitsordner je Fall und je Mutation entfernt, `GP_MUT`,
    `GP_PARALLEL`) gelten ohne Abweichung; der Code von `container-einrichten-gegenprobe.sh` darf übernommen werden; **Punkt 17 von `SPEC-060`** (Start mit ignoriertem `SIGINT`) gilt für deren Fälle mit Signal. Erwartet sind 30 bis
    120 s für den vollen Lauf (Schätzung des Architect: nach Stufe 1 rund 80 Fälle, nach Stufe 2 zusammen rund 110, dazu die Fälle von `SPEC-060` Punkt 16 (9); rund 180 Mutationen, jede fährt nur die Fälle ihrer Gruppe, je Fall wenige Prozess-Starts), Obergrenze nach `SPEC-060` Punkt 16 (1), auch nach Stufe 2; die gemessene Zeit steht in §7.
    Die Zeile der Gegenprobe in der Tabelle von [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 16 trägt der Slice, der sie anlegt, mit der Lieferung als `aufgenommen`; vorher steht sie dort nicht, weil sie nicht besteht. Fall
    `p14-aufnahme` (die Tabelle in `SPEC-060` nennt `zeitgrenze-gegenprobe` mit dem Zustand `aufgenommen`); Mutation: die Zeile fehlt.
