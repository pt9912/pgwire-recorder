# Slice slice-harness-gate-kette-stop-hook-freigabe: Der Stop-Hook gibt bei geänderten inerten Pfaden frei, wenn die schnelle Kette grün ist

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeordnet am 2026-10-12 nach Entscheidung des
Nutzers vom 2026-10-13 (die Gate-Kette steht vor den übrigen Release-Slices, vor
`slice-v1-abschluss-tls-client`): direkt hinter `slice-harness-gate-kette-stop-hook-liste`
(WIP-Limit 1); Reihenfolge in §4 *Start* und im Drift-Log der Roadmap. Der Schnitt der Hook-Arbeit in
`slice-harness-gate-kette-stop-hook-liste` und diesen Slice (vor dem Start, aus dem Plan
`slice-harness-gate-kette-stop-hook`) ist **Entscheidung des Nutzers vom 2026-10-12**
(„Deine Empfehlungen“).

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Wahl der
Kette am geänderten Pfad; Accepted 2026-10-12), das sich auf
[`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) stützt.

**Berührte Spec-Stellen:** `SPEC-052` (von `slice-harness-gate-kette-schnell` in
`spezifikation.md` §11 angelegt, mit den Randformen vom Architect entschieden; dieser Slice ändert
sie nicht; der Architect hat Punkte 9, 10, 12, 15 und 22 am 2026-10-14 vor dem Code geschärft, §6
*Ergänzung*, und nach dem Review (b) die Punkte 9, 12, 13, 20 und 21 geändert, §6 *Entscheidungen des
Architect nach dem Review*; Punkt 9 zuletzt in `278140c`). Der Umzug der Verträge der Harness-Werkzeuge nach `harness/vertraege/`
([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), Accepted 2026-10-13) folgt
direkt auf diesen Slice und zieht den Verweis auf `SPEC-052` nach; dieser Slice liegt vor ihm und ändert
den Ort nicht

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Stop-Hook `.claude/hooks/stop-require-gates.sh` wählt die Kette nach `SPEC-052`: Er gibt
frei, wenn jede seit dem Stand des letzten **vollen** grünen Laufs geänderte, hinzugekommene oder
gelöschte Datei auf der Positivliste liegt (zu Beginn `docs/plan/planning/`, `docs/reviews/`, `docs/user/` und die Wurzel-`README.md`; Entscheidung des Nutzers vom 2026-10-13, Beleg in `slice-harness-gate-kette-schnell` §6) und der
Stand der schnellen Kette (`gates-schnell`) grün ist; jeder andere Fall verlangt `make gates`
(fail-closed). Die Gegenprobe `stop-hook-gegenprobe` bekommt je Randform dieses Slice einen Fall, der nur
im inerten Fall freigibt. Der Hook ist das Restrisiko der Entscheidung
([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*): Ohne seine Gegenprobe
schließt der Slice nicht. Der Slice liefert die Ersparnis der Entscheidung: Ein Stop nach einem
Doku-Commit verlangt den vollen Lauf nicht mehr.

**Voraussetzung (Geber):** `slice-harness-gate-kette-stop-hook-liste` liegt in `done/`: die Dateiliste des
vollen Laufs mit ihrem Leser (dieser Slice liest sie über dieselbe Funktion, keine zweite Logik), der Hook,
der nur bei identischer Liste freigibt und sonst `make gates` verlangt, die Prüfung nicht aufgelöster
Index-Einträge **vor** jeder Entscheidung nach Pfad (dieser Slice stellt die Positivliste nach ihr, nie
davor), und das Gate `stop-hook-gegenprobe` mit dem Gerüst der Temp-Repos (dieser Slice fügt Fälle hinzu,
er legt das Gate nicht an). `slice-harness-gate-kette-schnell` liegt davor in `done/`: `SPEC-052`, das Ziel
`gates-schnell`, der Stand der schnellen Kette `.harness/state/gates-schnell-passed.diffsha` (Hash des
Arbeitsbaums aus `tools/harness/working-tree-hash.sh`). **Bindung (Geber `slice-harness-gate-kette-schnell`,
`SPEC-052` Teil 1, Punkte 5 bis 7):** Der Hook stützt die Freigabe auf den Stand der schnellen Kette und nur auf ihn: Er
vergleicht die Datei mit dem Hash des aktuellen Arbeitsbaums und liest sie unverändert. Zugesagt ist, dass der Stand nur
nach einem Lauf entsteht, in dem die Gates tatsächlich liefen (kein Stand bei `-n`, `-q`, `-t` und `-i`, ein älterer entfällt
nach Rot, Abbruch und `-i`), und dass der volle Nachweis davon getrennt bleibt. **Nicht zugesagt** und nicht Grundlage des Hooks
sind: dass nur die schnelle Kette die Datei schreibt und dass der volle Lauf sie nicht löscht (akzeptierte Negative des Gebers,
§6). **Name der Dateiliste** (gewählt vom Implementer von `slice-harness-gate-kette-stop-hook-liste`, eingetragen
vom Planner bei dessen Closure): `.harness/state/gates-passed.dateiliste`, neben `gates-passed.diffsha` (dem Hash dieser Liste); die
Funktionen `dateiliste_erzeugen` und `dateiliste_pruefen` und die Variable `DATEILISTE_STAND` stehen in `tools/harness/dateiliste.sh`, das der Hook mit `source` lädt.
**Gerüst des Gebers in `tools/harness/stop-hook-gegenprobe.sh`**, das dieser Slice erweitert: die Tabelle `FAELLE`, die Hilfen `fixture`, `kopie`,
`stand` (schreibt den Stand mit `record-gates.sh` im Normalaufruf), `hook`, `ist`, `stub` und `sedmut`; ein Fall trägt die Punktnummer
im Namen (`p<Punkt>-…`). **Zusagen des Gebers, auf die dieser Slice sich stützt** (`SPEC-052` Teil 2, vom Geber mit Fall geprüft; dieser Slice ändert sie nicht): Punkt 11 (Form
der Liste, Arten `F`, `L`, `G`, `X`), 16 (identische Liste gibt frei), 17 (Konflikt vor jeder Entscheidung), 19 (fehlender Stand sperrt),
20 (beschädigte Liste, Abbruch als `block`; der Schleifenschutz ist akzeptierte Grenze und überstimmt auch diesen Slice), 21 (`make -i gates` schreibt
keinen Nachweis, ein älterer entfällt) und 22 (Unterverzeichnis, leere Umgebung). Die Positivliste steht nach der Prüfung des Konflikts und ersetzt die Prüfung der Liste des Standes nicht.
`SPEC-052` Teil *Hook* gilt nach seinem Eingangssatz (Wortlaut: `slice-harness-gate-kette-schnell` §6, Zeile *Zwischenstand*), soweit ein Fall von `stop-hook-gegenprobe` die Zusage nennt;
dieser Slice erweitert das Gate um die Fälle der Positivliste, jeder mit der Punktnummer aus `SPEC-052` in seinem Namen (wie im Geber), und ändert die Spezifikation nicht: Der Eingangssatz bleibt mit diesen Fällen wahr, ohne dass jemand ihn nachzieht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Dateiliste, ihre Form, ihr Schreiber und der Vergleich „identisch oder voll“, die Randformen
  fehlender oder beschädigter Stand, nicht lesbarer Vergleich, Abbruch, nicht aufgelöster Index-Eintrag
  und Schleifenschutz — `slice-harness-gate-kette-stop-hook-liste` (liegt vor diesem Slice in `done/`):
  Dieser Slice setzt sie voraus und ändert sie nicht; ändert er sie, hat er den Plan geändert.
- `SPEC-052`, die schnelle Kette, ihr Stand und deren Gegenprobe — `slice-harness-gate-kette-schnell`
  (liegt in `done/`): Dieser Slice ändert die Spezifikation nicht; gibt der Implementer eine Randform
  zurück, entscheidet der Architect in `SPEC-052`.
- Aufnahme weiterer Pfade in die Positivliste (`docs/plan/adr/`, `spec/`, `harness/README.md`) —
  Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Setzungen des
  Nutzers*; eine Aufnahme ist eine Änderung von `SPEC-052` mit dem Beleg, dass kein langsames Gate den Pfad
  liest. `docs/user/` (einschließlich der Abdeckungstabellen) und die Wurzel-`README.md` stehen auf der
  Positivliste, belegt vom Architect am 2026-10-13 (`slice-harness-gate-kette-schnell` §6). Es gibt dafür
  keinen Slice und keine Kennung.
- Änderung von `make gates` — [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung*
  Punkt 1: Die volle Kette bleibt der Nachweis für alles, was nicht auf der Positivliste steht.
- Die Hooks `.claude/hooks/pretooluse-command-guard.sh` und `.claude/hooks/span-emit.sh` und den
  Eintrag in `.claude/settings.json` — anderer Vorgang: Der Slice ändert die Wahl der Kette im
  Stop-Hook, nicht den Befehls-Wächter oder die Erfassung.
- Eine Prüfung gegen die echte Hook-Umgebung von Claude Code — die Gegenprobe fährt den Hook in
  Temp-Repos mit dessen Eingabe und liest dessen Ausgabe; ob die Umgebung die Antwort so liest, prüft
  sie nicht (Risiko in §6).
- Produkt-Code, Benutzerhandbuch, `README.md` — Schicht-Abgrenzung.
## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1:** `.claude/hooks/stop-require-gates.sh` gibt nach `SPEC-052` frei, wenn alle
      geänderten Pfade auf der Positivliste liegen (Grenze mit Schrägstrich, bytegenau und case-sensitiv;
      gelöscht, umbenannt mit beiden Pfaden, untracked wie getrackt; Pfad mit `.`/`..`/leerem Segment
      oder absolut gilt als beschädigte Liste, geprüft an den Pfaden der Differenz; Art `X` ist nie inert;
      ein NUL im Pfad endet als Feldzahl-Fehler des Gebers, ohne eigenen Fall) **und** der Stand der schnellen Kette grün ist
      (Hash des aktuellen Arbeitsbaums, aus `tools/harness/working-tree-hash.sh` gerufen, nicht nachgebaut,
      gleich der Datei; fehlend, ungleich, leer oder nicht lesbar heißt kein Stand; fehlt das Skript oder
      scheitert es, gilt `block`, beim Schleifenschutz `approve`); bei identischer Liste gibt er ohne Stand der schnellen Kette frei; jeder andere Fall
      verlangt `make gates`. Der Kopf-Kommentar sagt die Positivliste zu, soweit die Gegenprobe sie prüft
      (`AGENTS.md` §3.7, §3.11).
- [x] **Liefer-Punkt 2:** Ein Symlink ist inert, wenn Link-Pfad und das lexikalisch gegen das Verzeichnis
      des Links aufgelöste Ziel (Eintrag `L`) auf der Positivliste liegen; Ziel außerhalb des Repos,
      absolut oder nicht auflösbar ist nicht inert; ein Ziel, das ein Verzeichnis der Liste selbst benennt, ebenso. Liegt irgendwo im
      Arbeitsbaum ein Link außerhalb der Positivliste, dessen Ziel auf sie zeigt, sie benennt oder ihr Vorfahr
      ist (`docs`, `docs/plan`, die Wurzel), ist jede Änderung nicht inert (Spiegelregel; `SPEC-052` Punkt 15).
- [x] **Liefer-Punkt 3:** Die Gegenprobe `stop-hook-gegenprobe` macht je Randform dieses Slice (§6) einen
      Fall in einem Temp-Repo, der nur im inerten Fall freigibt; ihre Zeile in `harness/README.md`
      §Sensors nennt die Positivliste und die Bindung an
      [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), soweit die Zeile aus
      `slice-harness-gate-kette-stop-hook-liste` sie bisher nicht trägt; `harness/conventions.md` (Träger-Tabelle,
      Zeile `grundlagen-durchsetzungsschicht.md`) bleibt richtig oder wird nachgezogen.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Der Reviewer prüft
      fail-closed ausdrücklich: Jede Abweichung von „inert und Stand grün“ verlangt `make gates`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Die Mutationen und ihre roten Fälle stehen als Tabelle in §7 (Zusage · Mutation · roter
      Fall), vom Implementer selbst gefahren (`AGENTS.md` §3.10).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben: eine weitere Datei
      `evidence/slice-harness-gate-kette-stop-hook-freigabe.md` in
      `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` mit der Messung dieses Slice (Dauer von
      `make gates` und der schnellen Kette; Freigabe eines Stops nach einem Doku-Commit; Last und
      Laufzeit dessen, was der Hook spart); der Zähler folgt aus den Dateien (4×; den dritten Beleg lieferte
      `slice-harness-gate-kette-schnell`) und der Eintrag, dort mit dem Ausgang *geplant* (Kennung dieses Slice)
      eingetragen, bekommt seinen Ausgang *verkörpert* (Zielort und Herkunfts-Anker `seit
      slice-harness-gate-kette-stop-hook-freigabe`), sofern die Messung die Senkung zeigt, sonst *offen* mit
      Begründung. Die Messung ist Beleg, kein Liefer-Punkt.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/hooks/stop-require-gates.sh` | update | Positivliste, Symlink-Regel, Stand der schnellen Kette; Rückfall auf voll bei jedem anderen Fall; Kopf-Kommentar |
| `tools/harness/` (Gegenprobe des Hooks aus `slice-harness-gate-kette-stop-hook-liste`) | update | je Randform dieses Slice ein Fall in einem Temp-Repo; das Fragment unter `harness/mk/` und die Anbindung an `GATE_CHECKS` bestehen |
| `harness/README.md` §Sensors | update, falls nötig | Zeile der Gegenprobe nennt die Positivliste |
| `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) | update, falls nötig | nennt die Träger des geänderten Hooks |
| `harness/mk/stop-hook-gegenprobe.mk` (Kopf) | update | nennt die Punkte 9 bis 22 |
| `.claude/commands/implement-slice.md`, `plan-welle.md`, `close-welle.md` | update | nennen die Positivliste (`inert` im Hook als Quelle) und `make gates-schnell` |
| `tools/harness/record-gates.sh` (Kopf) | update | nennt die Freigabe über die Positivliste statt „nur bei identischer Liste“ |
| `tools/hook/commit-msg-gegenprobe.sh` | update | legt Slices und Wellen ihrer Fälle in einem eigenen Temp-Repo an und liest den Planungsbaum des Arbeitsbaums nicht (Beleg von `SPEC-052` Punkt 13) |

- Der Hook liest die Liste über dieselbe Funktion wie in `slice-harness-gate-kette-stop-hook-liste`; die
  Positivliste ist die einzige neue Logik und steht nach der Prüfung nicht aufgelöster Index-Einträge.
- Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`.
## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gate-kette-stop-hook-liste` liegt in `done/`
(WIP-Limit 1); dieser Slice steht direkt dahinter und vor `slice-v1-abschluss-tls-client` (Schritt 28) und den übrigen Release-Slices (Drift-Log der Roadmap vom
2026-10-13). Die Reihenfolge ist Entscheidung des Nutzers vom 2026-10-13. Die Rückfragen des Planners hat der Architect am 2026-10-13 entschieden (§4 Größe, §6 Mutationen; Merge-Konflikt und Eingangssatz im Geber). Der Name der
Dateiliste steht in §1 (Eintrag des Planners bei der Closure des Gebers; `grep` nach ihm in diesem Plan
ist nicht leer). Die Randformen aus §6 sind in `SPEC-052` festgelegt; dieser Slice ändert sie nicht.

**Größe — Schätzung des Architect vom 2026-10-14 nach `AGENTS.md` §3.14 (ersetzt die Zahl 280 bis 420, die vor
§3.14 in anderer Lesart stand):** Metrik ist die Summe aus Einfügungen und Löschungen von `git diff --numstat
<Basis>..<letzter Commit des Slice> -- tools harness .claude`, Tests eingeschlossen; Basis ist der Commit, der
diesen Slice nach `in-progress/` bewegt (Plan, Register und Spezifikation liegen außerhalb dieser Pfade). Geschätzt
**550 bis 750 Zeilen**: Hook 130 bis 170 (Pfadregel mit Zulässigkeit, Differenz beider Listen, Symlink-Auflösung
und Spiegelregel, Stand der schnellen Kette, Kopf-Kommentar); Gegenprobe 400 bis 570 (Gerüst-Erweiterung 60: Dateien
der Positivliste in `fixture`, `stand_schnell`, Helfer je Pfadklasse; gut fünfzig Fälle mit je einer Mutation, weil
`AGENTS.md` §3.10 je Element einer Aufzählung einen Fall und eine Mutation verlangt; der Faktor 2 gegenüber dem
Produktcode ergäbe 300 und gilt hier als Untergrenze, weil die Aufzählungen der Positivliste und des Symlinks mehr
Fälle tragen als der Satz des Gebers); `harness/README.md` und `harness/conventions.md` 5 bis 15. Zum Vergleich: der
Geber lag bei Schätzung 440 bis 620 und Messung 1065, davon Gegenprobe 769. In *einer* Review-Sitzung prüfbar: ja, mit
deutlichem Abstand unter der Messung des Gebers, die der Verifier „am Rand“ nannte. Ein Schnitt vorab ist nicht
nötig; der Schnitt, falls die Marke eintritt, steht unten.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff liegt über einer Review-Sitzung
  (Marke: die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> --
  tools harness .claude` liegt über 900; die alte Marke war 500 in unbenannter Lesart, die Schätzung des Architect liegt bei 550 bis
  750, der Geber maß 1065 und war „am Rand“). Schnitt dann nach Lieferwert: erst die Positivliste mit Pfaden, Stand der
  schnellen Kette und Freigabe (Liefer-Punkt 1 und ihre Fälle); danach die Symlink-Regel mit Spiegelregel und
  ihren Fällen (Liefer-Punkt 2) als eigener Slice. In der Zwischenstufe ist ein geänderter Link nie inert, und die
  Spiegelregel steht schon im ersten Slice (sonst gäbe es eine Lücke); beide mit der Gegenprobe ihrer Zusagen,
  denn ohne Gegenprobe schließt keiner. Kennungen legt der Planner an, nicht dieser Plan.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt
  [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) zurück, oder das Format der Antwort des
  Stop-Hooks ändert sich so, dass die Gegenprobe nicht mehr gegen die echte Eingabe und Ausgabe des
  Hooks läuft.
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (vom Auftraggeber gefahren), die Gegenprobe des Hooks rot unter
jeder Mutation der Tabelle in §6, die Messung als Evidence im Register, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag (Wahl der Kette am Pfad).
Die Randformen sind vom Architect am 2026-10-12 **entschieden**, und zwar in der Tabelle von §6 des
Gebers `slice-harness-gate-kette-schnell` (Fundstelle `SPEC-052`; dessen erster Commit überträgt sie in
die Spezifikation, vor jedem Code). Dieser Slice trägt davon, ohne Doppelung mit
`slice-harness-gate-kette-stop-hook-liste`:

- Pfad auf der Positivliste (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`, `README.md`); Pfad auf der
  Grenze (`docs/plan/planning-x/`, `docs/reviews-alt/`, `docs/users/`, eine Datei `docs/reviews` oder
  `docs/user`, `README.md.bak`, `harness/README.md`); Groß- und Kleinschreibung (`docs/Plan/planning/`,
  `readme.md`);
- Pfad mit Segment `.`, `..` oder leerem Segment, absoluter Pfad, Pfad mit NUL: beschädigte Liste, nie
  inert;
- gelöschte oder umbenannte Datei (beide Pfade gelten als geändert und müssen inert sein);
- untracked Datei (zählt wie eine getrackte; eine durch `.gitignore` ausgeblendete ist nicht Teil der
  Liste, Bestand);
- Symlink (Link-Pfad und Ziel inert; Ziel lexikalisch gegen das Verzeichnis des Links aufgelöst; Ziel
  außerhalb des Repos, absolut oder nicht auflösbar nicht inert) und gespiegelter Link (Link außerhalb der
  Positivliste mit Ziel darin: jede Änderung nicht inert; der Bestand hat Links nur unter `.claude/rules/`
  mit Zielen außerhalb);
- leerer Diff gegen den vollen Lauf ohne Stand der schnellen Kette (frei, Teil des Gutfalls der
  Identität aus `slice-harness-gate-kette-stop-hook-liste`; hier der Fall, dass ein vorhandener roter
  Stand ihn nicht sperrt);
- Stand der schnellen Kette fehlt, ist rot oder gehört zu einem anderen Arbeitsbaum (der Hash des
  aktuellen Arbeitsbaums weicht von der Datei ab oder die Datei fehlt): voll;
- Bindung an den Geber `slice-harness-gate-kette-stop-hook-liste` (§1, eingetragen bei dessen Closure): Der Hook liest den Stand des
  vollen Laufs nur über `dateiliste_pruefen` und `dateiliste_erzeugen` aus `tools/harness/dateiliste.sh`, ein zweiter Leser entfällt; die Fälle der
  Positivliste legen den Stand mit `record-gates.sh` im Normalaufruf an, `make -i gates` ist dafür kein Weg (Punkt 21); der Schleifenschutz
  (`stop_hook_active`, Punkt 20) steht vor allen Prüfungen dieses Slice und bleibt dort;
- Bindung an den Geber `slice-harness-gate-kette-schnell` (`SPEC-052` Punkte 5 bis 7, §1): Der Hook liest den Stand
  unverändert und vergleicht nur den Hash; er stützt sich nicht darauf, dass nur die schnelle Kette die Datei schreibt
  oder dass der volle Lauf sie stehen lässt (nicht zugesagt; der Geber hat beides als akzeptiertes Negativ entschieden);
- gleichzeitige Läufe (**Aufruf-Voraussetzung**: nacheinander, `AGENTS.md` §3.12): `SPEC-052` bindet jede
  Zusage daran, der Hook sagt gleichzeitige Läufe nicht zu.

**Ergänzung — Entscheidung des Architect, 2026-10-14 (Fundstelle `SPEC-052`, vor dem Code in die Spezifikation
eingetragen).** Alle Randformen dieses Plans sind entschieden; keine ist „offen“. Punkte 9, 12, 13 und 15, die der Geber
nicht prüft und dieser Slice liefert, sind es mit den Zusätzen unten; Punkt 13 (Beleg für `docs/user/` und `README.md`)
liefert keinen eigenen Fall, er begründet die Aufnahme (die Fälle sind die der Pfadklassen in Punkt 12):

- **Reihenfolge** (Punkt 9): Schleifenschutz, Konflikt, fehlender oder beschädigter Stand des vollen Laufs, Identität, Pfadregel,
  Stand der schnellen Kette. Die Pfadregel ersetzt keine frühere Prüfung: eigener Fall je Schritt, der vor ihr steht, mit einer
  inerten Änderung und grünem Stand der schnellen Kette (Konflikt, fehlender Stand, beschädigter Stand: `block`).
- **Wahl der Kette und Aufruf**: kein Wahl-Skript. Die Pfadregel steht im Hook selbst und nirgends sonst, in reinem `bash`; der Hook ruft nur
  `tools/harness/working-tree-hash.sh` zusätzlich (Punkt 10). Fehlt das Skript: `block`, beim Schleifenschutz `approve` (je ein Fall). Aufrufarten
  des Treibers für den Hook sind in `SPEC-052` Punkt 22 entschieden und gelten für die Freigabe mit (Fall in Unterverzeichnis und in
  leerer Umgebung mit inertem Unterschied und grünem Stand der schnellen Kette). Der Hook ruft kein `make`: `-n`, `-q`, `-t`, `-i`
  und `MAKEFLAGS` sind für den Stand der schnellen Kette in `slice-harness-gate-kette-schnell` (Punkt 6) entschieden und dort
  gefahren; hier gilt: Der Hook vergleicht den Hash und prüft den Ursprung der Datei nicht, weil ein Lauf mit `-i` keinen
  Stand schreibt und einen älteren löscht (kein Fall im Hook, akzeptiertes Negativ mit diesem Grund; Punkt 10).
- **Hash-Vergleich** (Punkt 10): der Hash, nicht die Bytes der Datei; Datei leer oder fehlend: kein Stand (je ein Fall).
- **Pfade der Differenz** (Punkt 12): Zulässigkeit wird an den Pfaden der Differenz beider Listen geprüft; die Fälle für `.`, `..`,
  leeres Segment und absoluten Pfad legen den Pfad in den Stand des vollen Laufs (Eintrag der Form F mit gültigem Hash), der
  Arbeitsbaum hat ihn nicht; je Form ein Fall. NUL: kein eigener Fall (Feldzahl-Fehler des Gebers). Die Seite des
  Arbeitsbaums hat keinen herstellbaren Fall (akzeptiertes Negativ).
- **Art `X`** (Gitlink im Pfad der Positivliste geändert): nie inert, ein Fall.
- **Symlink** (Punkt 15), je Element ein Fall: Link und Ziel inert (frei); Ziel außerhalb des Repos; absolutes Ziel; Ziel über die
  Wurzel hinaus (nicht auflösbar); Ziel in einem Pfad außerhalb der Liste innerhalb des Repos; Ziel ist ein Verzeichnis der Liste
  selbst; Wechsel der Art in beide Richtungen (Datei zu Link, Link zu Datei, jeweils mit dem nicht inerten Ziel auf einer Seite). Spiegelregel,
  je Zieltyp ein Fall: Ziel auf der Liste; Ziel ist ein Verzeichnis der Liste; Ziel ist Vorfahr (`docs`); Ziel ist die Wurzel (`..` vom
  Verzeichnis eines Links aus); absolutes Ziel unter der Wurzel. Gegenrichtungen (absolutes Ziel außerhalb, Ziel über die Wurzel
  hinaus spiegeln nicht) sind ein Ärgernis, keine Lücke, ohne Mutation.
- **Positivliste ändert sich selbst** (Punkt 12): Änderung am Hook, an `tools/harness/dateiliste.sh`, an
  `tools/harness/working-tree-hash.sh` und an `spec/spezifikation.md` ist nicht inert und verlangt `make gates` (je ein Fall mit grünem Stand
  der schnellen Kette). Akzeptiertes Negativ: Der geänderte Hook entscheidet über seine eigene Änderung.
- **Nicht inert, je Element der Aufzählung in Punkt 9** (Code, Gates, Profile, Skripte, `Dockerfile`, Fragmente, `go.mod`) und die
  Bedingung (b) von Punkt 13 (`.dockerignore`) sowie `.d-check.yml`: je ein Fall mit inertem Unterschied daneben nicht nötig, der
  Pfad allein genügt (`internal/`, `Makefile`, `harness/mk/`, `.golangci.yml`, `tools/`, `Dockerfile`, `go.mod`, `.dockerignore`,
  `.d-check.yml`).
- **Gegenprobe-Gerüst** (Hinweis für den Implementer, keine Randform): Die Zeilenkonstanten des Hooks in `stop-hook-gegenprobe.sh`
  (`L_CMP`, `L_KONFLIKT`, `L_SCHLUSS` und die übrigen) bleiben erhalten oder werden im selben Commit nachgezogen; `fixture` bekommt
  Dateien auf der Positivliste, ohne dass ein Fall des Gebers sich ändert (deren Pfade liegen außerhalb).

**Entscheidungen des Architect nach dem Review, 2026-10-14 (b)** (F-713 bis F-722 des Reviews vom 2026-10-13; `SPEC-052`
Punkte 9, 12, 13, 20, 21 sind geändert, vor dem Code der Nacharbeit):

- **F-713 (MEDIUM): Entscheidung (c), die Vorbedingung beseitigen.** Das Skript `tools/hook/commit-msg-gegenprobe.sh` legt die
  Slices und Wellen seiner Positivfälle in einem eigenen Temp-Repo an (`git init`, Dateien unter `docs/plan/planning/{open,done}/`
  und flach in `docs/plan/planning/` und `done/` samt Closure-Notiz, `git add`; der Träger `.githooks/commit-msg` und
  `tools/harness/commit-msg-traceability.sh` werden dorthin kopiert, soweit der Träger sie braucht) und fährt die bisherigen
  fünfzehn Fälle dort. Es liest den Planungsbaum des Arbeitsbaums nicht mehr; das Abbrechen „kein Slice … im Index gefunden“
  entfällt mit der Vorbedingung. Die Namen im Fall dürfen feste Namen sein. Damit liest kein Gate der Ausnahmeliste
  `docs/plan/planning/`, `SPEC-052` Punkt 13 trägt den Beleg, die Positivliste bleibt, und die Ausnahmeliste bleibt unverändert.
  Verworfen: (a) `hook-gegenprobe` in die schnelle Kette (0,17 s gemessen, trüge, aber die schnelle Kette ist nach
  [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) Entscheidung 2 „ohne Gegenproben“ gefasst, und `SPEC-052` Punkte 2 und 4 bekämen
  eine Ausnahme von der Namensregel, §1 dieses Plans schließt die Änderung der schnellen Kette aus); (b) den Pfad von der Liste nehmen
  (verlöre die Hauptersparnis des Slice); (d) akzeptiertes Negativ (das Negativ der ADR deckt einen künftigen Leser, nicht
  einen vorhandenen). **Verdikt zu [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) Entscheidung 3:** Die Entscheidung gilt, und der Stand
  dieses Slice hat sie verletzt (ein vorhandenes langsames Gate las den Pfad, der Beleg in Punkt 13 fehlte für
  `docs/plan/planning/`); kein Widerspruch der ADR, keine Folge-ADR, weil die Ausnahmeliste und der Beleg in der Spezifikation stehen.
  Der Mangel ist **vor der Closure zu heilen** (die Closure ist ohne ihn blockiert, HIGH, wenn die Heilung unterbleibt). **Akzeptiertes
  Negativ** (Grund: gleiche Folge wie in der ADR): Führt ein künftiges Gate der Ausnahmeliste den Planungsbaum wieder als Eingabe, ändert es
  Skript oder Gate und löst den vollen Lauf aus; das Review dieser Änderung führt den Beleg neu. Ein Dauerfall dafür (das Skript in
  einer Kopie ohne offene Slices fahren) entfällt: Er wäre ein weiteres langsames Gate für einen Fehler, der beim Ändern des Skripts
  entsteht und dort auffällt. **Kosten:** 25 bis 40 Zeilen in `tools/hook/commit-msg-gegenprobe.sh`; kein Fragment, keine Ausnahmeliste, keine
  Spezifikation. **Zuordnung:** Liefer-Punkt 1 (Bedingung der Freigabe: kein langsames Gate liest die Positivliste); §3 nimmt die Datei auf.
- **F-717:** Der letzte Satz von `SPEC-052` Punkt 21 ist geändert (die Freigabe über die Positivliste ist genannt). `tools/harness/record-gates.sh`
  (F-714) zieht der Implementer nach.
- **F-718:** Die Ersetzung der drei äquivalenten Mutationen in §7 ist **bestätigt**; die Tabelle oben trägt sie (Zeilen *leere Stand-Datei*,
  *Hash-Skript*, *Pfad mit `..`*). Die Zusage „absolut nie inert“ ist zweifach gesichert, die Mutation lockert beide; das Minimum ist eine Mutation mit
  zwei Änderungen.
- **F-720:** Der `reason` bei fehlendem oder ungleichem Stand der schnellen Kette (jede Änderung inert) lautet
  ``The fast gates chain has no passed run for this working tree. Run `make gates-schnell` (or `make gates`).`` Der Implementer
  ändert `schnell_gruen` und `voll` (Hinweis als Parameter, Standard `make gates`); alle übrigen Antworten behalten
  ``Run `make gates`.`` Die Gegenprobe **prüft den `reason`**, aber nur in diesem Punkt: `p10-fehlt` (voller Stand da, schneller fehlt,
  Änderung inert; nicht `p9-stand-fehlt`, das den Stand des vollen Laufs prüft und `make gates` nennt) verlangt `make gates-schnell`
  im `reason`, ein neuer Fall `p9-stand-fehlt-nicht-inert` (Stand fehlt, Änderung in `internal/`) verlangt, dass er es **nicht** nennt. Der
  Treiber der Fälle bekommt dafür einen Helfer (`grund <Name> <Teilstring>`/`kein_grund`), der den `reason` aus der Antwort liest;
  andere Fälle prüfen weiter nur `decision`. Das löst F-719 mit: die Ordnung „Pfadregel, zuletzt Stand“ ist beobachtbar.
- **F-716:** Fall und Mutation, Punkt 12 bleibt (siehe Tabelle: `p12-unzulaessig-leer-ende`, Pfad `docs/user/x/` mit gültigem
  Hash im Stand; Mutation: `| */` aus `zulaessig` entfernen). Der Fall legt ihn in den Stand, nicht in den Arbeitsbaum (wie die übrigen
  Formen).
- **F-721:** Keine Randform, ein akzeptiertes Negativ, in `SPEC-052` Punkt 20 gefasst (keine Dauer zugesagt, kein `timeout` im Eintrag; was
  der Treiber bei Überschreiten tut, prüft die Gegenprobe nicht). Es gehört zum Risiko „Umgebung von Claude Code“ in der Liste unten
  (`BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`), kein neuer Eintrag. Eine Änderung an `.claude/settings.json` bleibt
  außerhalb dieses Slice (§1).
- **F-722:** `export LC_ALL=C` bleibt als Implementierung ohne Zusage und ohne Fall (kein erreichbares Fehlverhalten, aber die Zeile
  hält die bytegenaue Mustererkennung in jeder Locale; sie zu streichen spart eine Zeile und öffnet eine Frage). Der Kopf des Hooks
  sagt dazu nichts zu (so schon im Review). Kein Eintrag in der Spezifikation.
- **Größe und Zählung (§3.13, §3.14):** 3 Liefer-Punkte (die Heilung von F-713 steht unter Punkt 1), 1 Schicht (Harness-Werkzeug); die
  Nacharbeit schätze ich auf 60 bis 100 Zeilen in `tools harness .claude` (Skript des Commit-Trägers 25 bis 40, Hook 5, Gegenprobe des
  Hooks 25 bis 40 mit Helfer, zwei Fällen, drei Mutationen). Stand 639; erwartet 700 bis 740, unter der Marke 900.

Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie zurück.

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Die Zeilen mit „Planner“ in der letzten Spalte sind vom Planner abgeleitet; der
Architect hat sie am 2026-10-13 bestätigt oder geschärft (Spalte *Herkunft*). Je Zeile bricht die Mutation nur ihre Zusage:

| Zusage | Mutation | roter Fall | Herkunft |
|---|---|---|---|
| Nur ein Pfad der Positivliste gibt frei (`spec/`) | `spec/` in die Positivliste aufnehmen | Änderung in `spec/` gibt frei → rot | Architect |
| Ein Pfad auf der Grenze zählt nicht | Präfix ohne Verzeichnisgrenze vergleichen | `docs/plan/planning-x/` gibt frei → rot | Architect |
| Gelöschte und umbenannte Dateien zählen mit beiden Pfaden | nur den neuen Pfad prüfen | Umbenennung aus einem nicht inerten Pfad in einen inerten gibt frei → rot | Architect |
| Ein Symlink zählt mit seinem Ziel | Link ohne Ziel vergleichen | Link in inertem Pfad mit Ziel außerhalb gibt frei → rot | Architect |
| Ohne grünen Stand der schnellen Kette gibt der Hook nicht frei | Prüfung des Stands der schnellen Kette entfernen | rote oder fehlende Kette gibt frei → rot | Architect |
| Ein Link außerhalb der Positivliste mit Ziel darin macht jede Änderung nicht inert | die Spiegelregel entfernen | Änderung in `docs/plan/planning/` gibt frei → rot | Architect |
| Spiegelregel je Zieltyp: Ziel auf der Liste · Verzeichnis der Liste · Vorfahr (`docs`) · Wurzel · absolutes Ziel unter der Wurzel | je Zieltyp nur diesen Typ aus der Spiegelregel nehmen | Änderung in `docs/plan/planning/` gibt frei → rot, der Fall des Typs allein | Architect, 2026-10-14 |
| Ein Pfad mit `..`, absolut oder leerem Segment gilt als beschädigte Liste (je Form) | je Form nur diese nicht prüfen (fünf Mutationen: `..`, `.`, leeres Segment in der Mitte, leeres Segment am Ende, absolut). Absolut ist zweifach gesichert (Prüfung auf `/` in `zulaessig` und der Pfadanfang in `inert`): die Mutation lockert **beide** (`m-absolut-relativ`); die Prüfung allein zu entfernen (`m-ohne-absolut`) ist äquivalent und keine Mutation | `docs/plan/planning/../../internal/x.go`, `docs/plan/planning/./x`, `docs/plan/planning//x`, `docs/user/x/` (Ende auf `/`), absolut mit inerter Endung im Stand gibt frei → rot, der Fall der Form allein | Architect, geschärft 2026-10-14; Architect, 2026-10-14 (b) |
| Ein leerer Diff gibt ohne Stand der schnellen Kette frei | Stand der schnellen Kette auch dann verlangen | Hook blockiert bei identischem Baum, auch mit fehlendem, leerem und ungleichem Stand → rot | Architect |
| Der Stand der schnellen Kette gehört zum aktuellen Arbeitsbaum | nur das Vorhandensein der Stand-Datei prüfen, nicht ihren Hash | Stand eines anderen Arbeitsbaums gibt frei → rot | Planner, geschärft (Architect, 2026-10-13) |
| Eine leere Stand-Datei der schnellen Kette ist kein Stand | eine vorhandene leere Datei als Stand gelten lassen (`m-schnell-leer-gilt`; „leere Datei verlangen“ wäre äquivalent, weil ein leerer Hash ohnehin ungleich ist) | `p10-leer` gibt frei → rot (`p10-fehlt`, `p10-fremd` bleiben grün) | Architect, 2026-10-14; bestätigt 2026-10-14 (b) |
| Der Hook ruft `working-tree-hash.sh`; fehlt es oder scheitert es, gilt `block`, beim Schleifenschutz `approve` | je Zusage eine Mutation, die nur sie bricht: Scheitern des Skripts überlesen (`|| true`, `m-skript-ueberlesen`) · bei fehlendem Skript den gespeicherten Stand nehmen (`m-skript-fehlt-ersatz`; `|| true` allein wäre hier äquivalent, die leere Ausgabe wird abgelehnt) · leere Ausgabe gelten lassen (`m-hash-leer-gilt`) · Schleifenschutz hinter den Aufruf stellen | `p10-skript-scheitert` · `p10-skript-fehlt` · `p10-skript-leer` gibt frei → rot, je der Fall allein; mit `stop_hook_active: true` → `approve` → rot | Architect, 2026-10-14; bestätigt 2026-10-14 (b) |
| Die Pfadregel ersetzt keine frühere Prüfung (Reihenfolge, je Schritt ein Fall) | Pfadregel vor die Konfliktprüfung · vor die Prüfung des Standes des vollen Laufs auf Vorhandensein · vor `dateiliste_pruefen` stellen (drei Mutationen) | inerter Unterschied bei Konflikt · fehlendem Stand · beschädigtem Stand gibt frei → rot, jeweils der Fall des Schritts allein | Architect, 2026-10-14 |
| Gelöschte Datei außerhalb der Positivliste zählt | Pfade nur aus der Liste des Arbeitsbaums prüfen | gelöschte Datei in `internal/` neben grünem Stand der schnellen Kette gibt frei → rot | Architect, 2026-10-14 |
| Eine Umbenennung prüft den alten Pfad | nur den alten Pfad ignorieren (neuen prüfen) | Umbenennung aus `docs/user/` nach `internal/` gibt frei → rot (Gegenstück zur Zeile oben) | Architect, 2026-10-14 |
| Ein geänderter Eintrag der Art `X` ist nie inert | `X` wie `F` behandeln | geändertes Gitlink unter `docs/user/` gibt frei → rot | Architect, 2026-10-14 |
| Jeder Eintrag der Positivliste gibt frei: `docs/plan/planning/`, `docs/reviews/`, `README.md`, in jeder Tiefe | je Eintrag entfernen (drei Mutationen) · nur die erste Ebene unter dem Eintrag zulassen | Änderung in `docs/plan/planning/a/b/c.md`, `docs/reviews/…`, `README.md` gibt nicht frei → rot, je Eintrag allein | Architect, 2026-10-14 |
| Die Grenze gilt je Eintrag: `docs/plan/planning-x/`, `docs/reviews-alt/`, `docs/users/`, Datei `docs/reviews`, Datei `docs/user`, `README.md.bak` | je Eintrag den Schrägstrich bzw. den Gleichheitsvergleich weglassen (sechs Mutationen) | der Pfad des Eintrags gibt frei → rot, je Eintrag allein | Architect, 2026-10-14 |
| Der Hook kennt die Positivliste nur an einer Stelle und nicht inert sind Hook, `dateiliste.sh`, `working-tree-hash.sh`, die Spezifikation und die Aufzählung in Punkt 9 | je Pfad in die Positivliste aufnehmen (zehn Mutationen: Hook, `dateiliste.sh`, `working-tree-hash.sh`, `internal/`, `Makefile`, `harness/mk/`, `.golangci.yml`, `Dockerfile`, `go.mod`, `.dockerignore`, `.d-check.yml` → je eine) | Änderung des Pfades gibt frei → rot, je Pfad allein | Architect, 2026-10-14 |
| Freigabe in Unterverzeichnis und leerer Umgebung | Aufruf des Hash-Skripts relativ zum Aufrufer · Programm außerhalb der bisherigen verlangen | inerter Unterschied und grüner Stand: `approve` aus Unterverzeichnis und aus `env -i` → rot | Architect, 2026-10-14 |
| Symlink je Element (Link und Ziel inert · Ziel außerhalb · absolut · über die Wurzel · Ziel in Pfad außerhalb der Liste · Ziel ist Verzeichnis der Liste · Wechsel der Art in beide Richtungen) | je Element nur dieses Element lockern | der Fall des Elements allein (inert: frei; sonst: gibt frei → rot) | Architect, 2026-10-14 |
| `docs/user/` ist inert, auch die Abdeckungstabellen | `docs/user/` aus der Positivliste nehmen | Änderung in `docs/user/benutzerhandbuch.md` und in `docs/user/abdeckung-e2e.md` gibt nicht frei → rot | Architect, 2026-10-13 |
| Nur die Wurzel-`README.md` ist inert | den Namen `README.md` in jeder Tiefe vergleichen | `harness/README.md` gibt frei → rot | Architect, 2026-10-13 |
| Die Positivliste vergleicht case-sensitiv | Groß-/Kleinschreibung nicht unterscheiden | `docs/Plan/planning/` und `readme.md` geben frei → rot | Planner, bestätigt (Architect, 2026-10-13) |
| Eine untracked Datei zählt wie eine getrackte, nach ihrem Pfad | untracked Pfade pauschal als inert behandeln | untracked Datei in `internal/` neben geändertem inertem Pfad gibt frei → rot | Planner, geschärft (Architect, 2026-10-13; die Gegenrichtung, untracked inerter Pfad blockiert, ist ein Ärgernis, keine Lücke, und hat keine Mutation) |
| Fehlt der Stand der schnellen Kette bei inerter Änderung, nennt der `reason` `make gates-schnell`; fehlt er bei nicht inerter Änderung, nur `make gates` (`SPEC-052` Punkt 9; der Wortlaut für Stand eines anderen Arbeitsbaums, leer und nicht lesbar ist nicht zugesagt: akzeptiertes Negativ, die Zeile des Hooks ist dieselbe, ein Kunde des Wortlauts fehlt) | im Hook den Hinweis `make gates-schnell` auch in den anderen Fällen nennen · ihn im Fall des fehlenden Standes weglassen · `schnell_gruen` vor die Pfadregel stellen | `p10-fehlt` (`reason` nennt `make gates-schnell`), `p9-stand-fehlt-nicht-inert` („Stand fehlt, Änderung in `internal/`“; `reason` nennt es nicht) → rot, je der Fall allein; die dritte Mutation färbt den zweiten Fall (die Ordnung ist damit beobachtbar) | Architect, 2026-10-14 (b) |
| Kein Fall von `hook-gegenprobe` bezieht Slice oder Welle aus dem Index des Arbeitsbaums | das Skript liest wieder aus dem Index des Arbeitsbaums (die alte Fassung von `erster` und der Welle-Zeilen) | in einer Kopie des Arbeitsbaums, in der alle Slices aus `open/` nach `next/` bewegt sind, endet die Mutante mit Exit 1, das Skript der Lieferung mit 0; Beleg einmal von Hand in §7, kein Dauerfall (akzeptiertes Negativ unten) | Architect, 2026-10-14 (b) |

Die Mutation „nicht aufgelöster Index-Eintrag“ der Tabelle des Architect steht in
`slice-harness-gate-kette-stop-hook-liste` (dort als Konflikt bei sonst identischer Liste; Verteilung des
Planners, damit die Prüfung vor der Pfad-Entscheidung liegt und in der Zwischenstufe nicht fehlt).

**Risiken:**

- Der Hook ist das Restrisiko der Entscheidung: Ein Fehler dort bleibt ungeprüft, wo die Gegenprobe
  ihn nicht fängt ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*) — **Ausgang:** entfallen: Die Gegenprobe macht je
  Randform einen Fall, jede Mutation der Tabelle ist rot gesehen, und ohne sie schließt der Slice
  nicht (DoD-Punkt 3).
- Die Gegenprobe fährt den Hook in Temp-Repos mit dessen Eingabe, nicht in der Umgebung von Claude
  Code; ob die Umgebung die Antwort (`decision`) so liest, prüft sie nicht — **Ausgang:** weiter offen:
  → `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (von `slice-harness-gate-kette-stop-hook-liste`
  angelegt, 1×; mit diesem Slice 2×, Evidence bei der Closure).
- Ein Pfad der Positivliste wird später von einem langsamen Gate gelesen ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *akzeptiertes
  Negativ*) — **Ausgang:** entfallen: Ein späterer Leser ändert Code oder Gates und löst damit den
  vollen Lauf aus; die ADR hat das Negativ akzeptiert.
- Zwei Läufe gleichzeitig (der Hook liest, während die Kette ihren Stand schreibt) — **Ausgang:**
  entfallen: `SPEC-052` bindet die Zusagen an Läufe nacheinander (Aufruf-Voraussetzung, entschieden
  am 2026-10-12, `AGENTS.md` §3.12), und der Hook sagt gleichzeitige Läufe nicht zu; ein halb
  geschriebener Stand endet bei voll („kein Stand“).
- Die Positivliste liest sich als Zusage, die die Gegenprobe nur an Pfaden auf und neben der Grenze
  prüft (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 2×) — **Ausgang:** entfallen: Der
  Kopf-Kommentar sagt nur zu, was ein Fall der Gegenprobe prüft (`AGENTS.md` §3.11), und jede Mutation
  der Tabelle ist rot gesehen; tritt es dennoch auf, legt das Review den dritten Beleg an.
- Der geänderte Hook entscheidet über seine eigene Änderung (Positivliste ändert sich selbst): Ein Fehler im neuen Hook, der ihn
  immer freigeben ließe, wäre sein eigener Freigabegrund — **Ausgang:** entfallen: Die Änderung am Hook ist nicht inert, der Hook gibt sie nur
  frei, wenn er sich selbst falsch liest, und das fangen die Fälle der Gegenprobe und das Review; ein akzeptiertes Negativ von `SPEC-052`
  Punkt 12 mit diesem Grund.
- Der Hook gibt öfter frei als der volle Lauf für dieselbe Menge — **Ausgang:** entfallen: Er gibt nur frei, wenn jeder Unterschied zum
  Stand des vollen Laufs auf der Positivliste liegt, kein Gate der Ausnahmeliste diese Pfade liest (`SPEC-052` Punkt 13, seit 2026-10-14 (b) auch für `docs/plan/planning/` und
  `docs/reviews/`; `hook-gegenprobe` bezieht ihre Fälle aus einem eigenen Temp-Repo) und die schnelle
  Kette für genau diesen Arbeitsbaum grün war; jede Abweichung ist `block`, und jeder dieser Schritte hat einen Fall mit einer Mutation (Tabelle).
- Der Umzug der Verträge nach `harness/vertraege/` ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)) berührt
  diesen Slice nicht: `harness/vertraege/` steht nicht auf der Positivliste, ein Umzug ändert nur Pfade, die den vollen Lauf verlangen —
  **Ausgang:** entfallen: Der Umzug-Slice zieht Verweise auf `SPEC-052` (Kopf, §1, Kommentare in Hook und Gegenprobe) selbst nach; ändert er eine
  Zusage, hat er den Vertrag geändert.
- Der Nutzen schrumpft, wenn die Gegenproben schneller werden ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Re-Evaluierungs-Trigger*)
  — **Ausgang:** entfallen: Der Slice sagt keine Zeitersparnis zu; die Messung steht als Beleg im
  Register.
## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Code-Commits:** `6c6182b` (Hook und Gegenprobe), `a1fca6c` (Gate-Zeile, Träger-Tabelle, Fragment-Kopf, drei Commands),
`34a268f` und `b512b23` (Gegenprobe: zwei Rotfärbungen nachgezogen). Der Hook ist in einer Kopie im Scratchpad entwickelt und erst nach grüner
Gegenprobe eingespielt worden; der alte Stand: `git show 5822f9f:.claude/hooks/stop-require-gates.sh` und `git show 5822f9f:tools/harness/stop-hook-gegenprobe.sh`
(zurück: `git checkout 5822f9f -- .claude/hooks/stop-require-gates.sh tools/harness/`).

**Aufbau des Hooks:** Die Positivliste steht einmal, in `inert()`; `zulaessig()` lehnt absolute Pfade und Segmente `.`, `..`, leer ab; `ziel_lage()` löst das Ziel eines
Links lexikalisch auf; `spiegel()` ist die Spiegelregel; `schnell_gruen()` ruft `bash tools/harness/working-tree-hash.sh` (kein Nachbau) und vergleicht den Hash mit dem Inhalt
von `.harness/state/gates-schnell-passed.diffsha`; `freigabe()` fährt Pfadregel, dann Stand der schnellen Kette und endet in `frei` oder `voll`. Der Hook ruft kein `make`. Die Gegenprobe
setzt `export LC_ALL=C` im Hook voraus (bytegenaue Muster, auch für Nicht-UTF-8-Pfade); das ist Implementierung, keine Randform.

**Größe** (Metrik aus §4, `git diff --numstat 8349a65..b512b23 -- tools harness .claude`): 610 Einfügungen, 29 Löschungen, Summe **639** (Schätzung 550 bis 750, Marke 900).
Hook 200/8 (Schätzung 130 bis 170: die Zulässigkeit, die Auflösung des Ziels und die Spiegelregel brauchen je eine Funktion), Gegenprobe 392/9 (Schätzung 400 bis 570),
`harness/README.md`, `harness/conventions.md`, Fragment-Kopf je 1/1, drei Commands je 5/3. 70 neue Fälle (Schätzung: gut fünfzig), 136 insgesamt.

**Befund der Gegenprobe (vorher/nachher):** Der erste Lauf der neuen Mutanten meldete fünf Rotfärbungen, die nicht stimmten: (1) `m-pfad-vor-konflikt` färbte auch `p9-stand-beschaedigt`
(die Freigabe stand dort vor beiden späteren Prüfungen: Erwartung berichtigt); (2) derselbe Fall unter `m-pfad-vor-stand-fehlt`, dessen Mutant sich nicht auf das Vorhandensein beschränkte
(jetzt `[ -f "$stand" ] || freigabe`); (3) `m-absolut-relativ` färbte `p12-unzulaessig-absolut` nicht, weil der Pfad `/abs/docs/user/x.md` auch ohne `/` nicht mit `docs/user/` beginnt
(Pfad im Fall jetzt `/docs/user/x.md`); (4) `m-nur-arbeitsbaum` färbte `p14-geloescht-nicht-inert` nicht, weil eine im Arbeitsbaum gelöschte, im Index gebliebene Datei als `G`-Eintrag
in der Liste des Arbeitsbaums steht (neuer Fall `p14-geloescht-index-nicht-inert` mit `git rm`); (5) `m-spiegel-ohne-verzeichnis` färbte auch `p15-spiegel-absolut`, dessen Ziel ein Verzeichnis der Liste ist
(Erwartung berichtigt). Danach grün.

**Läufe** (Code-Commit `b512b23`, Belege-Commit `2a95992`, Baum sauber): `make gates` Exit 0, 7 min 18 s (438 s); darin `stop-hook-gegenprobe: gruen — 136 Fälle der Punkte 1, 9 bis 22 …`
(gemessen für sich: 2 min 3 s; der Geber lag bei 12 s mit 40 Fällen, der Stand davor bei 31 s mit 66), `gates-schnell-gegenprobe: gruen`, `d-check: 693 Datei(en) geprüft, 0 Befund(e)`,
`baseline-verify: v6.18.0 OK`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`. Der erste Lauf von `make gates` am Stand `34a268f` dauerte 7 min 16 s (436 s). Danach der Hook im
Arbeitsbaum: `{"decision":"approve"}`; ein Doku-Commit (dieser Eintrag, unter `docs/plan/planning/`) und `make gates-schnell` (Exit 0, 8 s, Docker-Cache warm) ändern daran nichts:
`{"decision":"approve"}` ohne erneuten vollen Lauf. Die Last dessen, was der Hook spart, ist der volle Lauf (7 min 18 s, Docker-Builds, Integrationstests, die Gegenproben); die schnelle Kette
fährt keine von ihnen. Der Hook ruft das Hash-Skript einmal (die Liste ein zweites Mal); gemessen im Repo mit einem inerten Unterschied: 4,5 s je Aufruf des Hooks.

**Mutationen** (`bash tools/harness/stop-hook-gegenprobe.sh`; Mutant = `sed`/`awk` auf der Kopie von Hook oder Skript im Temp-Repo, vom Skript selbst gefahren, frische Kopie je Mutant: kein
Docker-Build-Kontext. „rot“: der Fall meldet sich, „grün bleibt“: der Fall bleibt mit dem Mutanten grün. Die Mutanten in Klammern stehen mit ihrem Namen in der Gegenprobe):

| Zusage | Mutation | roter Fall (grün bleibt) |
|---|---|---|
| Der Konflikt steht vor der Pfadregel (Punkt 9) | `m-pfad-vor-konflikt`: Freigabe vor `git ls-files -u` | `p9-konflikt` (grün: `p9-stand-fehlt`) |
| Der fehlende Stand steht vor der Pfadregel | `m-stand-fehlt-als-unveraendert`: `[ -f "$stand" ] \|\| frei`; `m-pfad-vor-stand-fehlt`: `\|\| freigabe` | `p9-stand-fehlt` (nur unter dem ersten); der zweite bleibt grün, äquivalent: ohne Datei bricht das Lesen in der Freigabe ab, Antwort `block` |
| Der beschädigte Stand steht vor der Pfadregel | `m-pfad-vor-pruefung`: Freigabe vor `dateiliste_pruefen` | `p9-stand-beschaedigt` (grün: `p9-konflikt`, `p9-stand-fehlt`) |
| Jeder Pfad der Positivliste gibt frei, in jeder Tiefe | `m-ohne-planning`, `m-ohne-reviews`, `m-ohne-user`, `m-ohne-readme` (je ein Eintrag entfernt) | `p12-planning-tief`; `p12-reviews-tief`; `p12-user-tief`, `p13-docs-user`, `p9-freigabe`; `p12-readme` |
| … nicht nur die erste Ebene | `m-planning-erste-ebene`, `m-reviews-erste-ebene`, `m-user-erste-ebene` | `p12-planning-tief`; `p12-reviews-tief`; `p12-user-tief` (grün: `p9-freigabe`) |
| `docs/user/` ist inert, auch die Abdeckungstabellen | `m-ohne-user` | `p13-docs-user` (Handbuch und `abdeckung-e2e.md`) |
| Nur die Wurzel-`README.md` | `m-readme-jede-tiefe`: `README.md \| */README.md` | `p12-harness-readme` (grün: `p12-readme`) |
| Case-sensitiv | `m-case-egal`: `shopt -s nocasematch` | `p12-case-planning`, `p12-case-readme` (grün: `p9-freigabe`) |
| Die Grenze gilt je Eintrag | sechs Mutanten `m-aufnahme-14` bis `-19`: `docs/plan/planning-*`, `docs/reviews-*`, `docs/users/*`, `docs/reviews`, `docs/user`, `README.md.bak` je als zusätzlicher Eintrag | `p12-grenze-planning-x`; `-reviews-alt`; `-users`; `-datei-reviews`; `-datei-user`; `-readme-bak`, je der Fall allein |
| `harness/README.md` ist nicht inert | `m-aufnahme-20`: `harness/README.md` aufgenommen | `p12-harness-readme` |
| Die Positivliste kennt nur diese Pfade; nicht inert sind Hook, `dateiliste.sh`, `working-tree-hash.sh`, Spezifikation und die Aufzählung in Punkt 9 | dreizehn Mutanten `m-aufnahme-1` bis `-13` (`-5` ist `tools/*`, mit drei Fällen): `spec/*`, Hook, `dateiliste.sh`, `working-tree-hash.sh`, `internal/*`, `Makefile`, `harness/mk/*`, `.golangci.yml`, `Dockerfile`, `go.mod`, `.dockerignore`, `.d-check.yml` je aufgenommen | `p12-selbst-spec`; `p12-selbst-hook`; `p12-selbst-dateiliste`; `p12-selbst-hash`; `p9-nicht-inert-internal`, `-makefile`, `-harness-mk`, `-golangci`, `-dockerfile`, `-go-mod`, `-dockerignore`, `-d-check`; `tools/*`: `p9-nicht-inert-tools`, `p12-selbst-dateiliste`, `p12-selbst-hash` |
| Pfad mit `..` im Unterschied ist nie inert | `m-ohne-punktpunkt` | `p12-unzulaessig-punktpunkt` (grün: `-punkt`, `-leer`, `-absolut`) |
| … mit `.` | `m-ohne-punkt` | `p12-unzulaessig-punkt` (grün: die übrigen) |
| … mit leerem Segment | `m-ohne-leeres-segment` | `p12-unzulaessig-leer` (grün: die übrigen) |
| … absolut | `m-absolut-relativ`: Prüfung auf `/` entfernt **und** `inert` liest ohne führenden `/` (zwei Änderungen); die Prüfung allein zu entfernen (`m-ohne-absolut`) bleibt grün | `p12-unzulaessig-absolut` (grün unter `m-ohne-absolut`: äquivalent, die Positivliste lehnt `/docs/user/x.md` am Pfadanfang ab) |
| Untracked zählt nach seinem Pfad | `m-untracked-inert`: ein nicht getrackter Pfad ist inert | `p18-untracked-nicht-inert` (grün: `p14-neu-nicht-inert`) |
| Gelöschte und umbenannte Dateien zählen mit beiden Pfaden | `m-nur-arbeitsbaum`: Pfade nur aus der Liste des Arbeitsbaums | `p14-geloescht-index-nicht-inert`, `p14-umbenennung-nach-inert` (grün: `p14-geloescht-nicht-inert`, `p14-umbenennung-von-inert`, `p14-neu-nicht-inert`) |
| … der neue Pfad zählt | `m-nur-stand`: Pfade nur aus der Liste des Stands | `p14-umbenennung-von-inert`, `p14-neu-nicht-inert`, `p18-untracked-nicht-inert` (grün: `p14-geloescht-nicht-inert`, `p14-geloescht-index-nicht-inert`, `p14-umbenennung-nach-inert`) |
| Eine Datei auf der Grenze zählt nicht (Datei `docs/reviews`, `docs/user`) | siehe die Grenze oben | `p12-grenze-datei-reviews`, `p12-grenze-datei-user` |
| Ein geänderter Eintrag der Art `X` ist nie inert | `m-art-x-wie-datei`: die Zeile entfernt | `p12-art-x` |
| Ein Symlink zählt mit seinem Ziel | `m-link-ohne-ziel`: nur der Link-Pfad | alle Ziel- und Art-Fälle von Punkt 15 (acht); grün: `p15-link-inert` |
| Ein Ziel, das ein Verzeichnis der Liste benennt, ist nicht inert | `m-verzeichnis-inert` | `p15-ziel-verzeichnis` (grün: `p15-ziel-ausserhalb-liste`, `p15-link-inert`) |
| Ein Ziel über die Wurzel hinaus ist nicht inert | `m-ueber-wurzel-bleibt-innen`: `..` an der Wurzel bleibt an der Wurzel | `p15-ziel-aussen`, `p15-ziel-ueber-wurzel` |
| Ein absolutes Ziel ist nicht inert | `m-absolut-wie-relativ`: `ABS=1` entfernt | `p15-ziel-absolut` (grün: `p15-ziel-absolut-aussen`) |
| … außerhalb des Repos | `m-absolut-ohne-zweig`: der absolute Zweig entfällt | `p15-ziel-absolut`, `p15-ziel-absolut-aussen` |
| Art-Wechsel in beide Richtungen | `m-art-nur-stand`; `m-art-nur-arbeitsbaum`: nur eine Seite geprüft | `p15-art-datei-zu-link`; `p15-art-link-zu-datei` (je der andere grün) |
| Die Spiegelregel | `m-ohne-spiegel` | `p15-spiegel-liste`, `-verzeichnis`, `-vorfahr`, `-wurzel`, `-absolut` |
| … je Zieltyp | `m-spiegel-ohne-liste`; `-ohne-verzeichnis`; `-ohne-vorfahr`; `-ohne-wurzel`; `-ohne-absolut` | `p15-spiegel-liste`; `-verzeichnis` (und `-absolut`, dessen Ziel ein Verzeichnis ist); `-vorfahr`; `-wurzel`; `-absolut` |
| Ohne grünen Stand der schnellen Kette keine Freigabe | `m-schnell-nicht-geprueft` | `p10-fehlt`, `p10-fremd`, `p10-leer`, `p10-unlesbar` |
| Der Stand gehört zum aktuellen Arbeitsbaum | `m-schnell-nur-vorhanden`: `[ -n "$gespeichert" ]` | `p10-fremd` (grün: `p10-fehlt`, `p10-leer`) |
| Eine leere Stand-Datei ist kein Stand | `m-schnell-leer-gilt` | `p10-leer` (grün: `p10-fehlt`, `p10-fremd`) |
| Eine nicht lesbare Stand-Datei ist kein Stand | `m-schnell-unlesbar-gilt`: `cat … \|\| printf "$hash"` | `p10-unlesbar` (als Nutzer 1000 gelaufen), `p10-fehlt` (grün: `p10-leer`, `p10-fremd`) |
| Ein scheiterndes Hash-Skript sperrt | `m-skript-ueberlesen`: `\|\| true` | `p10-skript-scheitert` (grün: `p10-skript-fehlt`, `p10-skript-leer`) |
| Ein fehlendes Hash-Skript sperrt | `m-skript-fehlt-ersatz`: bei Fehlen den gespeicherten Stand nehmen; `m-skript-ueberlesen` bleibt hier grün, äquivalent: die leere Ausgabe wird abgelehnt | `p10-skript-fehlt` (grün: `p10-skript-leer`, `p10-skript-scheitert`) |
| Eine leere Ausgabe des Hash-Skripts sperrt | `m-hash-leer-gilt`: die Prüfung entfernt | `p10-skript-leer` (grün: `p10-skript-fehlt`, `p10-skript-scheitert`) |
| Der Schleifenschutz steht vor dem Aufruf des Skripts | `m-schutz-nach-hash` | `p10-skript-fehlt-schleifenschutz` |
| Ein leerer Diff gibt ohne passenden Stand der schnellen Kette frei | `m-schnell-bei-identitaet`: auch dort den Stand verlangen | `p16-identisch` (fehlend), `p16-schnell-leer`, `p16-schnell-fremd` |
| Der Hash kommt aus dem Skript, auch im Unterverzeichnis und in leerer Umgebung | `m-22-hash-relativ`: Aufruf über `$OLDPWD`; `m-22-hash-home`: `${HOME:?}` vor dem Aufruf | `p22-freigabe-unterverzeichnis` (grün: `p22-freigabe-leere-umgebung`, `p9-freigabe`); `p22-freigabe-leere-umgebung` (grün: `p22-freigabe-unterverzeichnis`, `p9-freigabe`) |
| Der Konflikt steht vor der Identität (Fall von (a); die Zeile `L_SCHLUSS` der Gegenprobe ist auf den Aufruf `freigabe` nachgezogen) | `m-identitaet-zuerst` | `p17-konflikt-identisch` |

**Grüne Mutanten, eingeordnet:** (1) `m-pfad-vor-stand-fehlt`: äquivalent, ohne Stand bricht die Freigabe beim Lesen ab; die Zusage trägt `m-stand-fehlt-als-unveraendert` mit
`p9-stand-fehlt`. (2) `m-ohne-absolut`: äquivalent, die Positivliste hält jeden absoluten Pfad am Anfang draußen; die Prüfung auf `/` ist zweite Verteidigung, rot erst mit `m-absolut-relativ`. (3) `m-skript-ueberlesen`
bei fehlendem Skript: äquivalent, die Prüfung auf leere Ausgabe fängt es; rot mit `m-skript-fehlt-ersatz` und `m-hash-leer-gilt` je an ihrem Fall. Die Seite des Arbeitsbaums hat für
unzulässige Pfade keinen herstellbaren Fall (`SPEC-052` Punkt 12, akzeptiertes Negativ). Die Tabelle des Plans (§6) nennt für „leere Datei“, „`|| true`“ und „absolut“ Mutationen, die in dieser Fassung
des Hooks äquivalent wären; sie sind durch die oben genannten ersetzt, die nur ihre Zusage brechen.

### Belege der Nacharbeit nach dem Review (b), 2026-10-14

**Code-Commits:** `b6c2fbe` (F-713, `tools/hook/commit-msg-gegenprobe.sh`), `9e727fb` (F-720, F-716, F-714: Hook, `stop-hook-gegenprobe.sh`, Kopf von `record-gates.sh`). Der Hook ist in einer Kopie im
Scratchpad geändert und erst nach grünem Lauf der Gegenprobe (138 Fälle) eingespielt worden; zurück: `git checkout 785ed69 -- .claude/hooks/stop-require-gates.sh tools/harness/ tools/hook/`.
F-717 und F-718 sind Entscheidungen des Architect ohne Code; F-721 und F-722 bleiben ohne Code (Entscheidung des Architect); F-719 löst der `reason` (unten).

**F-713, Befund vorher und nachher.** Kopie des Arbeitsbaums (alle Dateien des Index, eigenes Repo), je Zeile eine Menge aus dem Planungsbaum entfernt; „alt“ ist
`tools/hook/commit-msg-gegenprobe.sh` aus `1b7c390`, „neu“ das Skript aus `b6c2fbe`, jeweils der Exit-Code:

| Planungsbaum der Kopie | alt | neu |
|---|---|---|
| unverändert | 0 | 0 |
| ohne Slice in `open/` (der Fall der Meldung des Reviews) | 1 | 0 |
| ohne Slice in `done/` | 1 | 0 |
| ohne flache Welle | 1 | 0 |
| ohne geschlossene Welle | 1 | 0 |
| ohne Closure-Notiz | 1 | 0 |
| ohne jede Datei unter `docs/plan/planning/` | 1 | 0 |

Das „alt“ ist zugleich die Mutante der Tabelle in §6 (Zeile *Kein Fall von `hook-gegenprobe` bezieht Slice oder Welle aus dem Index*): das Skript liest wieder aus dem Index, Exit 1, das Skript der Lieferung 0.
Beleg von Hand, kein Dauerfall (akzeptiertes Negativ in §6).

**F-713, Zusagen des Trägers (`SPEC-050`) je Fall mit Mutation.** Die Mutante steht in einer Kopie des Trägers (`.githooks/commit-msg` oder `tools/harness/commit-msg-traceability.sh`) in einem Verzeichnis ohne
Planungsbaum; das Skript der Lieferung läuft dort (die Unverändert-Kontrolle: grün). Rot ist der Fall, den die Meldung `ROT — erwartet …` nennt; mehr rote Meldungen stehen dahinter:

| Fall (Zusage) | Mutation | roter Fall (außer ihm) |
|---|---|---|
| ohne jede Kennung abgelehnt | `exec bash …traceability.sh` am Ende durch `exit 0` ersetzt | `Arbeit ohne Kennung`; zugleich die Fälle erfundene Slice-Kennung, `slice-mv`, Suffix, Wortteil, Kommentar, Scissors, erfundene Welle, `-results` |
| erfundene Slice-Kennung abgelehnt | `case "$kennung" in slice-erfunden*) exit 0` vor der Prüfung | `slice-erfunden-gegenprobe` allein |
| Name des Lifecycle-Werkzeugs (`slice-mv`) abgelehnt | `case … slice-mv) exit 0` | `slice-mv` allein |
| Slice in `open/` angenommen | Muster nur `done/` | `slice-gegenprobe-offen`, Kennung im Rumpf |
| Slice in `done/` angenommen | Muster nur `open/` | `slice-gegenprobe-fertig` allein |
| Kennung nur im Rumpf angenommen | nur die erste Zeile lesen (`head -n 1`) | Rumpf-Fall allein |
| angehängtes Suffix abgelehnt | `(-[a-z0-9]+)*` auf `{0,1}` gekappt | `…-zusatz`; zugleich `-results` |
| vorangestellter Wortteil abgelehnt | Wortgrenze `(^|[^A-Za-z0-9_-])` entfernt | `xslice-gegenprobe-offen` allein |
| Kommentarzeile zählt nicht | Filter `grep -v '^[[:space:]]*#'` entfernt | Kommentar-Fall allein |
| Kennung unter der Scissors-Zeile zählt nicht | `sed '/^# -…/,$d'` durch `cat` ersetzt | Scissors-Fall allein |
| Lastenheft-Kennung angenommen | `LH-[A-Z]{2}-[0-9]{2}` aus `patterns=` entfernt | `Arbeit (LH-QA-04)` allein |
| flache Welle angenommen | Muster `docs/plan/planning/$kennung.md` entfernt | `welle-gegenprobe-offen` allein |
| Welle in `done/` angenommen | Muster `*/$kennung.md` der Welle entfernt | `welle-gegenprobe-fertig` allein |
| erfundene Welle abgelehnt | Muster der Welle auf `docs/plan/planning/*` | `welle-erfunden-gegenprobe` allein; ebenso `case … welle-erfunden*) exit 0` |
| Closure-Notiz (`-results`) zählt nicht | Zeile `*-results) continue ;;` gestrichen | `welle-gegenprobe-fertig-results` allein |

Das sind fünfzehn Fälle, je einer rot unter der Mutation seiner Zusage.

**F-720 und F-719, der `reason`.** Der Hook gibt bei fehlendem Stand der schnellen Kette mit inerter Änderung ``Run `make gates-schnell` (or `make gates`).`` aus; alle übrigen Sperren behalten ``Run `make gates`.``
`SPEC-052` Punkt 9 und §6 nennen den Fall „`p9-stand-fehlt`“; **dieser Fall prüft den Stand des vollen Laufs** (er legt nur den schnellen Stand an, Punkte 9 und 19) und sperrt mit ``Run `make gates`.``.
Der Fall des Standes der schnellen Kette bei inerter Änderung ist `p10-fehlt`; die Prüfung von `grund … 'make gates-schnell'` steht dort. Der neue Fall `p9-stand-fehlt-nicht-inert` (voller Stand da, schneller fehlt, Änderung in
`internal/x.go`) verlangt `kein_grund … 'make gates-schnell'`. Zuordnung als Frage an den Architect im Bericht. Mutanten:

| Zusage | Mutation | roter Fall (grün bleibt) |
|---|---|---|
| Fehlt der Stand der schnellen Kette bei inerter Änderung, nennt der `reason` `make gates-schnell` | `m-hinweis-nie`: die Zeile `schnell_gruen` ohne Hinweis | `p10-fehlt` (grün: `p9-stand-fehlt-nicht-inert`) |
| Jeder andere Fall nennt es nie | `m-hinweis-immer`: Standardhinweis von `voll` mit `make gates-schnell` | `p9-stand-fehlt-nicht-inert` (grün: `p10-fehlt`) |
| Die Pfadregel steht vor dem Stand der schnellen Kette (beobachtbar am `reason`, F-719) | `m-stand-vor-pfadregel`: `schnell_gruen` vor `pfadregel` | `p9-stand-fehlt-nicht-inert` (grün: `p10-fehlt`) |

Der bestehende Mutant `m-schnell-nur-vorhanden` (`voll "x"`) färbte beim ersten Lauf auch `p10-fehlt` rot, weil sein Text den Hinweis nicht trug; er trägt ihn jetzt und färbt nur `p10-fremd`.

**F-716, Klasse „Teilform des leeren Segments“.** Merkmal: der Pfad hat ein leeres Segment, und das Segment steht in der Mitte oder am Ende. Neuer Fall `p12-unzulaessig-leer-ende` (Pfad `docs/user/x/` mit gültigem Hash im Stand);
Mutanten: `m-ohne-leeres-segment-ende` (`| */` aus `zulaessig` entfernt) färbt `p12-unzulaessig-leer-ende` allein, `m-ohne-leeres-segment-mitte` (`*//*` entfernt) `p12-unzulaessig-leer` allein; der Mutant, der die ganze Zeile entfernt
(`m-ohne-leeres-segment`), färbt beide.

**F-714.** Der Kopf von `record-gates.sh` nennt die Freigabe über die Positivliste und den Stand der schnellen Kette und die Fälle `p16-identisch`, `p9-freigabe`; `grep -rn "gibt nur frei" tools` bleibt ohne Treffer.

**Größe** (Metrik aus §4, `git diff --numstat 8349a65..9e727fb -- tools harness .claude`): 670 Einfügungen, 60 Löschungen, Summe **730** (Stand vor der Nacharbeit 639; Schätzung 700 bis 740, Marke 900). Anteil der Nacharbeit (`b512b23..9e727fb`): 71/42, Summe **113**
(Schätzung 60 bis 100; darüber, weil Gegenprobe des Hooks 45 statt 25 bis 40: zwei zusätzliche Mutanten für das leere Segment und die dritte Mutante des `reason`).
Dateien: `commit-msg-gegenprobe.sh` 23/28 (netto −5), Hook 8/3, `stop-hook-gegenprobe.sh` 36/9, `record-gates.sh` 4/2.

**Laufzeit.** `make hook-gegenprobe`: vorher 0,17 s (ein Lauf am Stand `1b7c390`), nachher 0,18 bis 0,19 s (drei Läufe); `make stop-hook-gegenprobe`: 2 min 9,4 s mit 138 Fällen (vorher 2 min 3,8 s mit 136 Fällen, im Scratchpad gemessen);
die Laufzeit des Hooks (F-721) ist nicht verändert und nicht zugesagt.

**Läufe der engen Sensoren** (Code-Commit `9e727fb`): `make hook-gegenprobe` grün, `make stop-hook-gegenprobe` grün (138 Fälle), `make gates-schnell-gegenprobe` grün, `make gates-schnell` grün.

**Stand der Gates der Nacharbeit** (Code-Commit `9e727fb`, Belege-Commit davor: `git log -1 --format=%h -- docs/plan/planning/in-progress/slice-harness-gate-kette-stop-hook-freigabe.md` am Stand des Laufs; Baum sauber): `make gates` Exit 0, 7 min 22,9 s;
darin `d-check: 694 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.0 OK`, `commit-msg-gegenprobe: gruen` (hook-gegenprobe), `stop-hook-gegenprobe: gruen — 138 Fälle`, `gates-schnell-gegenprobe: gruen`,
`kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `a-check-negativ: gruen`, `run-integration-tests: gruen`. Danach `bash .claude/hooks/stop-require-gates.sh </dev/null` im Arbeitsbaum: `{"decision":"approve"}`.

**Nacharbeit nach der Verifikation** (V-340 bis V-342 des Reports vom 2026-10-13): Code-Commit `e48f655` (Kopf von `.claude/hooks/stop-require-gates.sh` und
`tools/harness/stop-hook-gegenprobe.sh` sagen den `reason` nur zu, wo ein Fall ihn prüft: fehlender Stand bei inerter Änderung `p10-fehlt`, bei nicht inerter
`p9-stand-fehlt-nicht-inert`; `tools/harness/record-gates.sh` und `harness/mk/stop-hook-gegenprobe.mk` tragen den Satz nicht, `grep` nach „anderen Arbeitsbaum“
und „gates-schnell“ im Kommentar ohne weiteren Treffer); die Änderung ist nur Kommentar, die Gegenprobe lief in einer Kopie im Scratchpad vor dem Commit
(`make stop-hook-gegenprobe` Exit 0, 138 Fälle, 2 min 10 s). Kopf des Plans (V-342) in `1459432`. V-343 ist Sache des Architect (Wortlaut von Punkt 22 (b)),
V-344 und V-345 ohne Folge.

**Stand der Gates nach der Verifikation** (V-341; Code-Commit `e48f655`, Kopf-Commit `1459432`; Baum sauber): `make gates` Exit 0, 7 min 32,0 s; darin
`d-check: 695 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.0 OK`, `commit-msg-gegenprobe: gruen`, `stop-hook-gegenprobe: gruen — 138 Fälle`,
`gates-schnell-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `a-check-negativ: gruen`, `run-integration-tests: gruen`.
Danach `bash .claude/hooks/stop-require-gates.sh </dev/null` im Arbeitsbaum: `{"decision":"approve"}`.

### Closure (Planner, 2026-10-13)

**Gegenstand und Lieferung.** Liefer-Punkte 1 bis 3 sind vom Verifier bestätigt (Report `docs/reviews/2026-10-13-verifikation-slice-harness-gate-kette-stop-hook-freigabe.md`, Commit `67cf485`: 0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO). Review: `docs/reviews/2026-10-13-review-slice-harness-gate-kette-stop-hook-freigabe.md`, Commit `785ed69` (0 HIGH, 1 MEDIUM, 3 LOW, 7 INFO). Entscheidungen des Architect: `23414ff`, `2f49e7e`, `d788954`, `5822f9f`, `1b7c390`, `278140c`.

**Gate-Lauf (Belege des Implementers, nicht des Planners).** `make gates` Exit 0, 7 min 32,0 s am Stand `e48f655`/`1459432` (letzter Stand mit Änderung an `tools`, `harness`, `.claude`, `spec`); danach kamen `0af4744` und diese Closure, beide nur unter `docs/plan/planning/` (inert), und `make gates-schnell` (Exit 0). Der Hook im Arbeitsbaum antwortete am `HEAD` `0af4744` `approve`. **Die Messung am Stand nach der Closure steht aus** (Auftraggeber); der Planner hat `make gates` nicht gefahren. Gefahren hat er `make hook-gegenprobe` (grün) und `make docs-check` und `make kopf-check`.

**F-713 ist geheilt (vom Planner geprüft).** `tools/hook/commit-msg-gegenprobe.sh` (`b6c2fbe`) legt Slices und Wellen seiner Positivfälle in einem eigenen Temp-Repo an und liest den Planungsbaum des Arbeitsbaums nicht. Probe des Planners: das Skript mit seinem Träger und der Traceability-Prüfung allein in einem Verzeichnis ohne `docs/` gefahren (kein Planungsbaum): Exit 0, Meldung `commit-msg-gegenprobe: gruen`; `make hook-gegenprobe` im Arbeitsbaum: grün. Die Verifikation fuhr dieselbe Probe mit Kopien, in denen Slices verschoben waren (Tabelle in *Belege der Nacharbeit*: alt Exit 1, neu Exit 0). Damit liest kein Gate der Ausnahmeliste `docs/plan/planning/`, und [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) Entscheidung 3 gilt für diesen Stand.

**Befunde ohne Folge.** F-717 (`SPEC-052` Punkt 21 nachgezogen), F-718 (Ersetzung bestätigt), F-719 (gelöst durch die `reason`-Fälle), F-720 (behoben), F-721 (kein `timeout`, keine Dauer zugesagt; Register unten), F-722 (`export LC_ALL=C` ohne Zusage und Fall, Entscheidung des Architect), F-723 (Beleg von Punkt 13 vollständig) und V-344 (Nacharbeit nach dem Review (b) hat keinen zweiten Report, die Rollenfolge verlangt ihn nicht), V-345 (die Probe eines Handbuch-Beispiels gegen das gebaute Binary bleibt Pflicht von Implementer und Verifier eines Slice mit Handbuch-Änderung; dieser ändert keines): INFO, kein Folge-Slice. F-714, F-715, F-716, V-340, V-342 sind behoben (`9e727fb`, `2f34a98`, `e48f655`, `1459432`); V-341 ist mit dem Gate-Lauf oben beantwortet.

**Offener Punkt: V-343 (benannte Spec-Lücke).** `SPEC-052` Punkt 22 (b) liest sich einzeln als „leere Umgebung: bei identischer Liste `approve`, bei jedem Unterschied `block`“; erst der folgende Satz („beide Aufrufarten gelten auch für die Freigabe nach Punkt 9“) hebt das auf. Die Gegenprobe prüft die Freigabe in beiden Aufrufarten (`p22-freigabe-unterverzeichnis`, `p22-freigabe-leere-umgebung`, Mutanten `m-22-hash-relativ`, `m-22-hash-home`); der Code ist richtig. Entscheidung des Architect: als offenen Punkt nennen, die Spezifikation in diesem Slice nicht ändern. Der Punkt hat keinen Nehmer: Weder Slice noch Register-Zeile nehmen ihn an (`AGENTS.md` §3.13 verlangt vor einer Zuweisung den Blick in den Nehmer); eine Berichtigung des Wortlauts wäre eine Änderung von `SPEC-052` im Beschluss des Architect, nicht in einem Slice-Plan. Wer `SPEC-052` Punkt 22 das nächste Mal anfasst (etwa der Umzug nach `harness/vertraege/`, [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)), liest diesen Absatz.

**Größe** (Urteil des Planners nach der Marke aus §4): Summe **730** (`git diff --numstat 8349a65..9e727fb -- tools harness .claude`: 670 Einfügungen, 60 Löschungen; vom Planner nachgerechnet), Schätzung des Architect 550 bis 750, Marke 900: **nicht eingetreten**, keine Rückführung. Die Nacharbeit nach `9e727fb` ist Kommentar und Plantext. Reviewer (F-721) und Verifier: in einer Sitzung prüfbar. Die Schätzung in der Metrik von `AGENTS.md` §3.14 traf (730 gegen 550 bis 750); beim Geber `slice-harness-gate-kette-stop-hook-liste` lag die Schätzung ohne diese Metrik bei 440 bis 620 und die Messung bei 1065. Das ist ein Beleg dafür, dass §3.14 wirkt; ein Zähler im Register entsteht daraus nicht (Zeile im `state.md` von `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, ohne Evidence-Datei).

**Risiken aus §6, je ein Ausgang.**

| Risiko (§6) | Ausgang | Begründung bzw. Träger |
|---|---|---|
| Der Hook ist das Restrisiko der Entscheidung | entfallen | Die Gegenprobe macht je Randform einen Fall (138), jede Mutation der Tabellen ist rot gesehen; der Verifier fuhr zwanzig eigene Mutationen, alle rot an den benannten Fällen |
| Die Gegenprobe prüft nicht die Umgebung von Claude Code | weiter offen | `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (2×, Evidence F-721) |
| Ein Pfad der Positivliste wird später von einem langsamen Gate gelesen | entfallen | Ein späterer Leser ändert Skript oder Gate und löst den vollen Lauf aus (akzeptiertes Negativ der ADR); der vorhandene Leser (F-713) ist geheilt |
| Zwei Läufe gleichzeitig | entfallen | `SPEC-052` bindet die Zusagen an Läufe nacheinander; ein halb geschriebener Stand endet bei voll |
| Die Positivliste liest sich als Zusage, die die Gegenprobe nur auf und neben der Grenze prüft | entfallen | Der Kopf des Hooks sagt nur zu, was ein Fall prüft (V-340 behoben); `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` bleibt bei 2× (kein neuer Beleg) |
| Der geänderte Hook entscheidet über seine eigene Änderung | entfallen | Die Änderung am Hook ist nicht inert (Fall `p12-selbst-hook`); akzeptiertes Negativ von `SPEC-052` Punkt 12 |
| Der Hook gibt öfter frei als der volle Lauf für dieselbe Menge | entfallen | Jeder Schritt hat Fall und Mutation; F-713 war der einzige vorhandene Leser und ist geheilt (Verifikation: kein weiterer unter den langsamen Gates) |
| Der Umzug der Verträge nach `harness/vertraege/` berührt diesen Slice nicht | entfallen | Der Slice liegt vor dem Umzug und ändert den Ort nicht. Der Nachzug der Verweise auf `SPEC-052` (Hook, Gegenprobe, Plan) ist Folgepflicht in [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) *Konsequenzen*; der Slice dafür ist in der Roadmap vorgemerkt und noch nicht angelegt, hat also keine Kennung. Das Risiko ist daher nicht *eingetreten* mit Kennung; die Vormerkung der Roadmap ist der Träger, ein Register-Eintrag wäre ein zweiter Mechanismus für dieselbe Pflicht |
| Der Nutzen schrumpft, wenn die Gegenproben schneller werden | entfallen | Der Slice sagt keine Zeitersparnis zu; die Messung steht als Evidence im Register |

**Register** (Commit `300362f`; Zähler = Dateien unter `evidence/`, nachgezählt):

| Eintrag | Beleg | Zähler, Stand |
|---|---|---|
| `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` | Messung dieses Slice | 4×, **verkörpert** (Zielort Hook, schnelle Kette, `SPEC-052`; Anker `seit slice-harness-gate-kette-stop-hook-freigabe`); der Weg eines Doku-Stops sinkt von 7 min 32 s auf 8 s plus 4,5 s, die Laufzeit von `make gates` wuchs (2 min 9 s allein `stop-hook-gegenprobe`) |
| `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | F-716, F-719 | 32×, verkörpert (`AGENTS.md` §3.10), Regel bleibt |
| `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | F-714, F-717, V-340 | 51×, verkörpert (`AGENTS.md` §3.11), Regel bleibt |
| `BEO-REPO/plan-folgt-korrektur-nicht` | F-715, F-718, V-342 | 41×, verkörpert (`AGENTS.md` §3.9), Regel bleibt |
| `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` | Risiko 2, F-721 | 2×, offen |
| `BEO-REPO/beleg-ueber-den-bestand-ohne-probe` (neu) | F-713 | 1×, offen |

F-713 steht bewusst nicht unter `BEO-REPO/spec-randform-erst-im-review-entschieden`: Dort fehlt eine Randform in der Spezifikation; hier stand ein falscher Beleg des Architect in ihr (Urteil des Planners). Kein Eintrag steht mit diesem Slice neu bei 3× ohne Ausgang.

**Lerneintrag — benannte Spec-Lücke und zwei Beobachtungen.** (1) Form *benannte Spec-Lücke*: V-343, oben; sie ist als Lücke benannt und nicht beauftragt. (2) Ein Beleg des Architect („kein langsames Gate liest diese Pfade“) wurde im Review widerlegt (F-713). Ein Beleg über den Bestand ist eine Aussage zu einem Stand und braucht eine Probe, die ihn hält; das Review fuhr sie (Slices verschieben, `hook-gegenprobe`), der Architect hatte sie nicht gefahren. Das ist einmal aufgetreten und steht als `BEO-REPO/beleg-ueber-den-bestand-ohne-probe` (1×, offen) im Register; eine Regel folgt erst ab 3× (Schwelle des Registers, und die Rolle des Architect sieht einen einmaligen Blindfleck als akzeptiertes Negativ vor). Wird der dritte Beleg gezählt, ist der Vorschlag: ein Satz in `AGENTS.md` §3.12, ein Beleg aus dem Bestand nenne die Probe, die ihn wieder prüft. (3) Die Schätzung in der Metrik von `AGENTS.md` §3.14 hat getroffen (730 gegen 550 bis 750): kein Anlass für eine Änderung der Regel. Keine Regel in `AGENTS.md`, `.claude/` oder `harness/` wird mit dieser Closure geändert; Entscheidungen des Nutzers: keine offen.

**Paarungen.** (a) Anker: `grep -rn "seit slice-harness-gate-kette-stop-hook-freigabe" docs/plan/planning/observations` findet `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/state.md`. (b) Folge-Slice: kein Risiko ist *eingetreten*; V-343 und der Umzug haben keinen angelegten Nehmer (oben). (c) Register: die genannten Kennungen existieren als Verzeichnis mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/`
gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand `f1f4415`, 2026-10-12). Treffer:

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (3×, *geplant* mit der Kennung dieses Slice, seit
  `slice-harness-gate-kette-schnell`) — Gegenstand: Mit der Messung dieses Slice (DoD) steht der Eintrag
  bei 4×; sein Ausgang lautet *verkörpert*, Zielort
  Hook und Kette, Herkunfts-Anker `seit slice-harness-gate-kette-stop-hook-freigabe`, sofern der Hook die
  Laufzeit eines Doku-Commits zeigbar senkt, sonst *offen* mit der Begründung.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, offen) — berührt: Symlink und nicht lesbarer
  Eintrag sind Randformen (§6); übergeht der Hook einen Eintrag still, ist es das zweite Auftreten.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — berührt: Die Positivliste liest
  sich als Zusage; die Gegenprobe prüft Pfade auf und neben der Grenze. Tritt es wieder auf, steht der
  Eintrag bei 3×.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — den dritten Beleg legt
  `slice-harness-gate-kette-stop-hook-liste` an; hier berührt: §4 Größe und Rückführung.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen entscheidet der Architect vor dem Code des Gebers, §4 *Start* nennt es.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) und
  `BEO-REPO/slice-waechst-durch-uebernahmen` (verkörpert, `AGENTS.md` §3.13) — Mutationen in §6, kein
  Kommentar sagt mehr zu, als die Gegenprobe prüft, die Sendungen sind nachgezählt.

Dieser Slice führt den vierten Beleg für `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (den dritten lieferte
`slice-harness-gate-kette-schnell`, der Eintrag ist seither *geplant* mit der Kennung dieses Slice) und ist mit
`slice-harness-gate-kette-schnell` selbst die Antwort darauf; sonst erreicht keiner der Einträge mit
diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (Positivliste mit Freigabe und Stand der schnellen
Kette; Symlink-Regel; Fälle der Gegenprobe mit der Gate-Zeile). Messung als Beleg, Gate-Läufe, Review und
Closure-Notiz zählen nicht. Schichten: 1 — Harness-Werkzeug (`.claude/hooks`, `tools/harness`,
`harness/README.md`, `harness/conventions.md`); keine Spezifikation (sie ändert
`slice-harness-gate-kette-schnell`), keine Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap.
**Sendungen** des Gebers `slice-harness-gate-kette-stop-hook-liste` (Liste mit Leser, Gerüst von
`stop-hook-gegenprobe`, Prüfung nicht aufgelöster Index-Einträge vor der Pfad-Entscheidung, Name der
Dateiliste) und der Randformen und Mutationen der Hook-Seite aus `slice-harness-gate-kette-schnell`, soweit
verteilt: nachgezählt mit den Sendungen 3 Liefer-Punkte und 1 Schicht; der Name der Dateiliste ist Teil
der Liefer-Punkte 1 bis 3, kein vierter. Nachgezählt am 2026-10-13 bei der Closure von `slice-harness-gate-kette-stop-hook-liste` (Sendung: der Name
`.harness/state/gates-passed.dateiliste`, die Funktionen und das Gerüst aus §1, die Bindung aus §6): Namen und Bindung ändern die Liefer-Punkte nicht
(der Name ist Teil der Liefer-Punkte 1 bis 3, die Bindung ein Satz in §6), die Zählung bleibt 3 Liefer-Punkte und 1 Schicht. Die Bindung an `SPEC-052` Teil 1 (§1, Punkte 5 bis 7 des Gebers
`slice-harness-gate-kette-schnell`, bei dessen Closure eingetragen) ist kein vierter Liefer-Punkt und keine Schicht: Der Hook
liest den Stand ohnehin (Liefer-Punkt 1); nachgezählt 3 Liefer-Punkte, 1 Schicht. Geschätzt 280 bis 420 Zeilen mit Tests (§4). Nachgezählt am
2026-10-13 bei der Umordnung (vor die Release-Slices): Die Sendungen sind unverändert, die Zählung bleibt 3
Liefer-Punkte und 1 Schicht. Nachgezählt vom Architect am 2026-10-13 mit der erweiterten Positivliste (`docs/user/`,
`README.md`): zwei Mutationen mehr in der Tabelle, vier Pfadklassen mehr in den Fällen der Positivliste, keine neue
Logik (dieselbe Präfix-Regel, dazu ein Gleichheitsvergleich für `README.md`); die Zählung bleibt 3 Liefer-Punkte
und 1 Schicht, die Schätzung 260 bis 390 Zeilen steigt um etwa 20 bis 30 (bis 420), keine Rückgabe an den Planner; vom Architect am 2026-10-13 als 280 bis 420 übernommen.

Nachgezählt vom Architect am 2026-10-14 mit den Zusätzen von `SPEC-052` (Reihenfolge, Hash-Aufruf, Zulässigkeit, Art `X`, Symlink-Einzelheiten,
Selbständerung): 3 Liefer-Punkte, 1 Schicht. Die Zusätze liegen im Hook (Liefer-Punkt 1 und 2) und in der Gegenprobe (Liefer-Punkt 3), kein vierter
Punkt, keine neue Schicht. Die Größe steht in §4 in der Metrik von `AGENTS.md` §3.14 (550 bis 750 Zeilen, Marke 900). Sendungen: Geber
`slice-harness-gate-kette-schnell` und `slice-harness-gate-kette-stop-hook-liste` stehen mit Kennung in §1.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Vermerk Architect, 2026-10-13/14:** Prüfung vor dem Code; Randformen, Größe, Fail-open-Grenzen und Sendungen entschieden wie in §1, §4 und §6;
`SPEC-052` Punkte 9, 10, 12, 15, 22 geschärft. Der Plan ist bereit für den Start (offen: nichts für den Planner außer Verteilung der Kennungen eines
späteren Schnitts, falls die Marke eintritt).
