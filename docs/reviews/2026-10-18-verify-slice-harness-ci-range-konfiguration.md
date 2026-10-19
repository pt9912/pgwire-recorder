# Verifikation: slice-harness-ci-range-konfiguration

Verifier, 2026-10-18. Stand: Zweig `rev-rk`, Kopf `d8086425` (enthält den Review-Report `6c0fd266`), Basis `368e3d43` (origin/main), Basis der Größenmetrik `eac90f75`. Plan: [slice-harness-ci-range-konfiguration](../plan/planning/done/welle-erster-release/slice-harness-ci-range-konfiguration.md) (§2 DoD, §6 Risiken, §7 Belege und Nacharbeit). Vertrag: [SPEC-064](../../harness/vertraege/SPEC-064-ci.md) Punkte 10, 11, 12, 15, 16, 17, 18. Entscheidungen: [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md), [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md). Hard Rules: [AGENTS.md](../../AGENTS.md) §3.9 bis §3.12, §3.14. Review: [2026-10-18-review-slice-harness-ci-range-konfiguration](2026-10-18-review-slice-harness-ci-range-konfiguration.md) (0 HIGH, 2 MEDIUM, 3 LOW, 5 INFO).

Frage dieser Rolle: Bauen wir es richtig (gegen Plan und DoD)? Nichts repariert; committet ist nur dieser Report. Host: uid 1000, Docker am Host, keine volle Gegenprobe, kein `make gates`.

## Urteil

**Bestanden, mit zwei MEDIUM-Auflagen vor der Closure** (V-1, V-2), kein HIGH. Alle DoD-Punkte sind im normalen Pfad erfüllt und belegt; die Wirkung ist im frischen Klon echt gemessen. V-1 ist ein Fehler im ausgelieferten Skript, den nur die dokumentierte Arbeitshilfe `GP_MUT` auslöst (nicht der volle Lauf); V-2 sind Mutanten gegen Zusagen des Vertrags, die grün bleiben. Beide fallen in die Klasse, die das Review mit F-1 und F-2 schon trug, und beide sind unsichtbar für den vollen Lauf.

## Ergebnis je Punkt

| # | Prüfung | Ergebnis |
|---|---|---|
| 1 | `.d-check.yml` am Kopf: `commits` und `vcs` wörtlich nach Punkt 11 und 12 | erfüllt. `commits.id-patterns` mit genau den vier Mustern, `exempt-pattern: '^(Merge \|Revert )'`; `vcs.paths`, `immutable-when`, `exclude-sections: [Geschichte]`, `status-line`, `head-allow`, `ignore-link-targets: true` Zeichen für Zeichen wie im Vertrag (Diff `368e3d43..HEAD` von `.d-check.yml`: 19 Zeilen, nur diese Blöcke und ihr Kommentar) |
| 2 | `workflows` nicht in `modules:` | erfüllt. Zeile 10: `modules: [links, anchors, ids, matrix, spans]`; `workflows` steht nirgends in der Datei |
| 3 | Wirkung im frischen Klon (`git clone --no-hardlinks --branch rev-rk`, ohne Netz) | erfüllt, siehe Abschnitt Wirkung |
| 4 | Gegenprobe, Mutationen am Host (`GP_MUT`, `GP_ENTDECKE=1`) | 41 Mutationen und die unveränderte Kopie gefahren, alle mit genau der erwarteten roten Menge (die sechs `m-p18-*` nur in einer Kopie ohne den Fehler V-1), Details unten |
| 5 | `make docs-check` am Kopf | `d-check: 1175 Datei(en) geprüft, 0 Befund(e)` |
| 6 | `make gates-schnell-gegenprobe` (DoD Liefer-Punkt 2) | grün, die Namens- und Zuordnungsregel nimmt das neue Ziel an (`GEGENPROBEN_CONTAINER += range-konfiguration-gegenprobe`, nicht in `GATE_CHECKS`) |
| 7 | Größenmetrik am Kopf | siehe Abschnitt Metrik: 433 |
| 8 | Fragment, Skript, `harness/README.md`-Zeile (DoD Liefer-Punkte 2 und 3) | vorhanden, `DOCKER` kommt aus `harness/mk/build.mk`; Zeile mit Vertrag, Bindung und Prosa, die auf die Quelle zeigt |

## Wirkung echt (Punkt 3)

Frischer Klon des Hauptklons (Zweig `rev-rk`, `d8086425`), `RANGE=368e3d43..HEAD` (8 Commits, `history-range-guard` OK):

| Aufruf | Befund |
|---|---|
| `make doc-commits` | `0 Befund(e)` |
| `make doc-immutable` | `0 Befund(e)` |
| Probe: leerer Commit `probe ohne Kennung`, dann `doc-commits` über Basis `d8086425..HEAD` | `commit-untraceable  probe ohne Kennung`, 1 Befund, Exit 2 (über `make`; d-check selbst Exit 1) |
| Probe: Commit mit einer ADR-Kennung (Nummer 0001) im Betreff ändert eine Zeile des Kerns von `0001-hexagonale-architektur.md` (Status `Accepted`), dann `doc-immutable` | `core-drift-vcs  Core einer immutablen Datei hat sich über die Commit-Range geändert`, 1 Befund, Exit 2 |
| Gegenprobe zur Probe: im selben Klon `.d-check.yml` aus `368e3d43` (ohne die Blöcke) | beide Targets `0 Befund(e)` über derselben Range |

Die Targets prüfen damit nicht mehr „über nichts grün“; ohne die Blöcke bleiben sie es. Der Klon ist entfernt.

## Mutationen (Punkt 4)

Am Host mit `GP_ENTDECKE=1` in vier Läufen gesehen; genau die benannte Menge rot, wie §7 sie nennt:

- `m-p11-block-fehlt` (9 Fälle); Muster fehlt: `adr`, `lh`, `mr` (mit `p11-kennung-im-rumpf`), `slice-welle`; Alternative `slice` und `welle` je einzeln; Ausnahme fehlt, `Merge `, `Revert `, ohne Anker.
- Zu F-1 (Form): Weitung `adr`, `lh`, `mr`, `slice-welle` (je der `-knapp`-Fall; `slice-welle` färbt `p11-slice-knapp` und `p11-welle-knapp`); Anker `adr`, `lh`, `mr` (mit `p11-kennung-im-rumpf`), `slice-welle`; Struktur-Kennung ergänzt (`p11-struktur-allein`).
- Zu Punkt 12: `m-p12-block-fehlt`, `paths-leer`, `paths-fehlt`, `immutable-when-jeder-status`, `-anderer-status`, `-fehlt` (alle acht Fälle), `abschnitt-fehlt`, `-falsch`, `status-zeile-fehlt`, `-eng`, `head-allow-fehlt`, `-eng`, `-ohne-accepted`, `link-ziel-fehlt`, `-aus`.
- Zu F-2 (Abbruch): die sechs `m-p18-*` sind in der Kopie des Skripts ohne den Fehler V-1 jeweils genau auf dem eigenen Signal rot (`grün (6 Läufe)`); im ausgelieferten Skript nur ohne `GP_MUT`, das auf sie passt, siehe V-1.
- Die unveränderte Kopie mit allen 27 Fällen (`ok`) und der Lauf `fall-ohne-funktion`: `grün (2 Läufe)`.

Falsche Röte der Fälle: Jede Mutation färbt die Fälle der eigenen Gruppe und keine anderen (die Gegenprobe verlangt die Menge; mehr und weniger machen sie rot), also ist jeder Fall ohne den Fix aus dem richtigen Grund rot (Röte des Falls = Zusage gebrochen), nicht durch Umgebung.

## Plan-vs-Code-Diff (§3.9)

- Der Diff der Implementer-Commits am Plan (`dfeb2d23`, `00e072bf`, `d8086425`) betrifft §3 (eine Zeile: Fälle und Mutationen 27 und 41) und §7; §1, §2, §6 und der Kopf blieben unberührt. Das ist die erlaubte Menge. Die Zeilen von §6 (Randformen 8 und 9) und die Nachschätzung in §8 stammen vom Architect (`313d6459`).
- Dateien aus §3 liegen alle vor und sonst nichts: `.d-check.yml`, Fragment, Skript, `harness/README.md`, `SPEC-064` (Architect). Der Diff gegen die Basis enthält zusätzlich den Review-Report und eine Ergänzung in `slice-harness-ci` (Architect, Punkt 18 als Randform 18 dort); das passt zu AGENTS.md §3.13.
- §7: Zählung stimmt mit dem Code (17 + 8 + 2 Fälle, 20 + 15 + 6 Mutationen = 41, 43 Läufe).
- **Lücke (P-1, LOW, Planner):** Weil nur §7 und §3 geändert werden durften, trägt der Plan den Abbruch (Punkt 18) und die Form der Muster nur in §3, §6 (Randformen 8 und 9) und §8, nicht in §1 und in der DoD: Liefer-Punkt 2 nennt `p11-*` und `p12-*` und keine `p18-*`; der Kopf (Berührte Spec-Stellen) führt Punkt 16 und 17, nicht 18; §1 nennt den Abbruch weder als Ziel noch als Abgrenzung. Das ist nach `AGENTS.md` §3.9 ein Nachzug, den der Planner schuldet (der Plan nennt, was geliefert wird; Punkt 18 ist geliefert, `kopf-check` prüft nur Kennungen aus §1/§2). Kein Code-Befund.

## Metrik (Punkt 5)

`git diff --numstat eac90f75..HEAD -- tools harness .d-check.yml Makefile`: Einfügungen 426, Löschungen 7, **Summe 433** (Skript 386, `.d-check.yml` 19, `SPEC-064` 16, Fragment 11, `harness/README.md` 1). Wie beim Implementer; unter der Meldeschwelle 450 und der Marke 500. Risiko 5 aus §6: entfallen, belegt.

## Befunde

| ID | Schwere | Befund |
|---|---|---|
| V-1 | MEDIUM | **`GP_MUT` auf eine Mutation des Abbruchs löst eine unbeschränkte Rekursion aus.** Der Kopf des Skripts verspricht „`GP_MUT=<Muster>` fährt nur die Mutationen, deren Name passt“ und „`GP_INNEN=1` lässt die Kopie die p18-Fälle auslassen“. Die Fälle `p18-sigint`/`p18-sigterm` starten eine Kopie des Skripts; die Kopie erbt `GP_MUT` aus der Umgebung und `GP_INNEN` schaltet die p18-Fälle nur für die unveränderte Kopie (`ok`) ab, nicht in `mut18`. Passt `GP_MUT` auf `m-p18-*`, fährt die Kopie wieder `mut18`, startet ihrerseits Kopien, und so fort. Beobachtet: `GP_MUT='^m-p18-.*$' make range-konfiguration-gegenprobe` hinterließ über 1000 angehaltene Prozesse mit verschachtelten Pfaden `…/m-p18-…/tmp.…/m-p18-…`, Load-Average 87, falsche Röte (`m-p18-term-verzeichnis-zu-frueh`: erwartet `p18-sigterm`, rot waren beide Fälle) und Läufe über 10 Minuten; der Aufräumer des Skripts fing das nicht (eigene Prozessgruppe der Kopie). Ursache bestätigt: In einer Kopie des Skripts, die die Kopie mit `env -u GP_MUT -u GP_ENTDECKE` startet, sind alle sechs `m-p18-*` mit `GP_PARALLEL=2` in 60 s genau rot auf dem eigenen Signal. Im vollen Lauf (ohne `GP_MUT`) erreicht die Kopie `mut18` erst spät und wird vorher vom Signal beendet; ob das Rennen immer so ausgeht, ist nicht belegt (der Implementer sah es in seinen Läufen nie). Der Auftrag an den Verifier und der Text des Kopfes verweisen auf `GP_MUT` als Arbeitshilfe; ein Muster, das eine der sechs trifft, ist damit gefährlich. Behebung: die Kopie ohne `GP_MUT` und `GP_ENTDECKE` starten (`env -u`), den Kopfsatz berichtigen, und für die Rekursion einen Fall (die Kopie startet keine zweite Kopie). Ausgang vor der Closure: Auflage an den Implementer; da das Skript zugesagtes Verhalten aus Punkt 18 trägt, wird die Wirkung der Hilfe nicht neu entschieden. |
| V-2 | MEDIUM | **Grüne Mutanten gegen Zusagen von Punkt 11 (Zeichenklassen der Muster).** Der Vertrag sagt „`LH-` mit zwei Großbuchstaben“ und „`slice-` oder `welle-` und mindestens ein Zeichen aus `a-z0-9`“. Eine Mutation von `.d-check.yml` je Klasse färbt keinen Fall: `'LH-[A-Z]{2}-[0-9]{2}'` → `'LH-[A-Za-z]{2}-[0-9]{2}'` (kleine Buchstaben treffen), `'(slice\|welle)-[a-z0-9]+(…)'` → `[a-z]+` (Ziffer am Anfang trifft nicht) und → `[0-9a-f]+` (Buchstaben `g` bis `z` am Anfang treffen nicht). Mit `GP_ENTDECKE=1` in einer Kopie des Skripts mit drei zusätzlichen `mut`-Zeilen gesehen: alle `=> nichts`. Die Fälle haben nur große Buchstaben für `LH-FA`, nur Buchstaben in `slice-aa-bb`/`welle-aa`, kein Namensteil mit Ziffer. Das ist die Klasse von F-1 (Form des Musters), die der Implementer mit „Länge der Ziffernfolge, Zahl der Buchstaben, Mindestlänge des Namens“ benennt, ohne die Zeichenklasse; `AGENTS.md` §3.10 zählt je Element einer Aufzählung (`a-z`, `0-9`, Großbuchstaben) eine Mutation. Zwei Wege: ein Fall je Klasse (`LH-fa-01` rot, `slice-1x` grün, `welle-X` rot) mit Mutation, oder der Architect nimmt die Klasse aus der Zusage und nennt sie als Grenze. Der Vertrag selbst listet diese Mutationen nicht, der Implementer hat ihn also erfüllt; die Lücke liegt in der Fallliste. Ausgang vor der Closure: Architect entscheidet, danach Implementer. |
| V-3 | INFO | Zwei weitere Mutanten bleiben grün, beide ohne Zusage im Vertrag: der Schwanz `(-[a-z0-9]+)*` des Slice-Musters entfernt, und das Leerzeichen in `^(Merge\|Revert)` entfernt (`v-schwanz`, `v-leerzeichen`, beide `=> nichts`). Der Schwanz ist für die Anwesenheitsprüfung äquivalent (Invariante: ein Treffer des längeren Musters enthält einen Treffer des kürzeren; Ort: die Suche nach dem Muster in der Zeile, nicht ein Ganzzeichen-Vergleich; Sonde: `v-schwanz`, kein Fall wird rot). Das Leerzeichen ist nicht äquivalent (`Mergex` würde ausgenommen); der Vertrag nennt die Ausnahme als `Merge `/`Revert ` mit Leerzeichen, ein Fall dazu fehlt. Nur zur Kenntnis. |
| V-4 | LOW | Der Kommentar im Skript (`GP_INNEN=1` lässt die Kopie die p18-Fälle auslassen) sagt mehr, als der Code tut (nur für `ok`); gleiche Ursache wie V-1, mit V-1 zu beheben (`AGENTS.md` §3.11). |
| P-1 | LOW | Plan: Kopf, §1 und DoD Liefer-Punkt 2 tragen Punkt 18 (Abbruch) nicht, siehe Plan-vs-Code. Planner. |

Kein HIGH. Kein Befund zu: Metrik, Wortlaut der Blöcke, Eintrag in `GEGENPROBEN_CONTAINER`, Abwesenheit von `workflows`, Wirkung im Klon.

## §3.11 Kommentare

- `.d-check.yml`: die Zusage „Beide Module stehen nicht in `modules:`, die Targets aktivieren sie mit `--enable`“ ist wahr (`d-check.mk` Zeilen 41 und 45); „ohne diese Blöcke melden beide Targets über jeder Range 0 Befund(e)“ habe ich im Klon gemessen; „fehlt ein Block oder ein Eintrag (bei `paths` auch leer), wird `range-konfiguration-gegenprobe` rot“ trägt je Eintrag eine Mutation (`m-p11-block-fehlt`, `m-p12-block-fehlt`, Muster-, `m-p11-ausnahme-fehlt`, `m-p12-paths-fehlt`/`-leer`, `immutable-when-fehlt`, `abschnitt-fehlt`, `status-zeile-fehlt`, `head-allow-fehlt`, `link-ziel-fehlt`), alle gesehen.
- Fragment `harness/mk/range-konfiguration.mk`: Kopf nennt Ziel, Liste und Voraussetzungen; wahr.
- Skript: Kopf bis auf V-1/V-4 getragen; „Kein Fall hängt von Dateirechten ab“ belegt der Implementer mit Läufen als root und unter `umask 077` (nicht von mir wiederholt).
- `harness/README.md`-Zeile: nennt Regel an den Schreiber („wer einen der Blöcke ändert, fährt sie, die Regeln selbst stehen in `.d-check.yml`“) und zeigt auf `SPEC-064`; die Behauptung Exit 130 bzw. 143 ist durch `p18-*` belegt. Keine Nacherzählung der Blöcke.

## Risiken aus §6: Ausgang für die Closure

| # | Risiko | Vorgesehener Ausgang | Stand |
|---|---|---|---|
| 1 | gepinntes d-check liest die Konfiguration anders | entfallen | trägt: jeder Fall der Punkte 11 und 12 zeigt es mit dem echten Image (meine Läufe); nach V-2 mit den drei Klassen ergänzen |
| 2 | Bestand färbt die erste Range rot | entfallen | trägt: `368e3d43..HEAD` mit 0 Befunden gemessen; der Implementer nennt zusätzlich den Bestand (`d027ecc5…`, 2076 Commits, 0 Befunde, von mir nicht wiederholt); am Endstand noch einmal über die Range des Slice messen |
| 3 | Gegenprobe braucht das Image | entfallen | trägt: Ziel in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS` |
| 4 | Löschen der Blöcke färbt die Gegenprobe, nicht die CI | entfallen mit Grund aus Punkt 15 | trägt unter der Bedingung „`make gegenproben` am Ende jeder Welle“ (`AGENTS.md` §6 Schritt 6); sonst weiter offen → Register |
| 5 | Größe unter der Marke | entfallen | belegt: 433 < 500 |

Nicht im Plan, aber für die Closure nötig: (a) V-1 als eingetretenes Risiko der Arbeitshilfe `GP_MUT` mit Behebung oder als Folge-Slice mit Kennung; (b) V-2 als Fall und Mutation oder als Grenze in `SPEC-064`; (c) Register: der Implementer nennt in §7 den Treffer `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (Stand offen, Zähler 2, `observation.md`/`state.md` zeigen „Der Zähler bleibt bei 2“) — mit diesem Slice 3×, nach [Modul 5 §Ziel-Form: Sub-Area-Modus-Begründung](../../.harness/baseline/v6.18.0/regelwerk/modul-05-planning-harness.md) und Modul 6 keine Notiz mehr, sondern eine Lücke mit Folge-Slice; das Register ist noch nicht fortgeschrieben (kommt mit der Closure). Ebenso `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (V-1 zeigt eine weitere Umgebungs- und Aufrufabhängigkeit der Gegenprobe) und die Beobachtung zu CRLF (`zeilenende-crlf-im-gelesenen-wert`, der Implementer hat gemessen, §7).

## Belege des Implementers (Punkt 1 der Aufgabe)

Alle von der DoD verlangten Belege stehen in §7: Befundzeilen vorher/nachher (Tabelle), Range des Slice mit 0 Befunden, Probe-Commits der Gegenprobe, Mutationen mit roten Fällen, Läufe `make range-konfiguration-gegenprobe` (`grün (43 Läufe)`), `make gates-schnell-gegenprobe`, `make gates` Exit 0, `make docs-check`, Metrik 433. Den Beleg zu `make gates` (4 min 20 s) habe ich nicht wiederholt (Auftrag). Die Läufe im Dev-Container sind, soweit ich sie am Host wiederholt habe, bestätigt (43 Läufe entsprechen 41 Mutationen, `ok`, `fall-ohne-funktion`); die Behauptung „drei Läufe der sechs Mutationen des Abbruchs mit derselben roten Menge“ (Sonde 5) habe ich am Host nur in der Kopie ohne V-1 bestätigt.

## Nicht geprüft

Voller Lauf der Gegenprobe (läuft im Container grün laut §7), `make gates`, Läufe als root und unter `umask 077`, `make doc-trace`, `make kopf-check` (läuft in `make gates`).

Aufgeräumt: Der Klon, die Kopie des Skripts und die Temp-Verzeichnisse des Rekursionsfehlers sind entfernt, alle dabei entstandenen Prozesse beendet.
