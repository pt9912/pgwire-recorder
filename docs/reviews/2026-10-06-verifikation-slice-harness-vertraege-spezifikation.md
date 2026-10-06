# Verifikation: slice-harness-vertraege-spezifikation — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-vertraege-spezifikation.md` (Kopf, §1 bis §6, §8) gegen den Gesamt-Diff `e491e0c..c667677` ohne `docs/plan/planning/open/slice-harness-abdeckung-gate.md` (den änderte `cbe650e` für einen anderen Slice). Darin: `afed2f0` (Randformen), `be1401c` (§11 *Harness-Werkzeuge* mit [`SPEC-047`](../../spec/spezifikation.md#spec-047--kopf-der-pläne-kopf-check) und [`SPEC-048`](../../spec/spezifikation.md#spec-048--abdeckung-je-anforderung-und-pfad-abdeckung), Historie als §12), `0139f5b` und `f0c87b6` (Skriptköpfe, Sensor-Dateien, Gegenproben), `d0eed74` (Lesarten bestätigt), `d9de9d0` (Review), `0c45064`, `2ba6b31` und `c667677` (Nacharbeit nach F-417 bis F-426, ohne eigenes Review).

**Eingang:**

- die DoD-Liefer-Punkte, §6 *Randformen* und *Lesarten beim Übertragen* des Plans und die Commit-Messages
- der [Review-Report](2026-10-06-review-slice-harness-vertraege-spezifikation.md) (F-417 bis F-427, Stand `f0c87b6`)
- `git show afed2f0:tools/harness/kopf-check.sh` und `git show afed2f0:tools/test/abdeckung.sh` (alte Köpfe), [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 1 bis 9 und §Konsequenzen, [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) §Entscheidung und §Konsequenzen, §6 des archivierten `slice-harness-kopf-sensor` aus `docs/plan/planning/done/welle-replay-semantik/archiv.zip` (Lesarten (1) bis (5), (a), (b))
- `spec/spezifikation.md` §11 und §12, `harness/sensors/kopf-check.md`, `harness/sensors/abdeckung-check.md`, die Vorlage `.harness/baseline/v6.13.0/templates/harness/sensors/gate.template.md`, `harness/README.md` §Sensors, `harness/mk/kopf-check.mk`, `harness/mk/abdeckung.mk`, `.d-check.yml`
- beide Skripte und beide Gegenproben am Stand `c667677`, dazu die Gegenproben am Stand `f0c87b6`
- die Pläne `docs/plan/planning/open/slice-harness-lint.md`, `docs/plan/planning/open/slice-harness-coverage.md`, `docs/plan/planning/open/slice-harness-mutation.md`, `docs/plan/planning/open/slice-harness-abdeckung-gate.md`, Kopf der vier Umstellungs-Slices

Bei meinem Start war der Arbeitsbaum sauber, HEAD `c667677`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

Alle Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg.

**Wie die Proben gebaut wurden.** Jede Probe lief in einer frischen Kopie aus `git archive c667677` unter `verify-vs/` im Scratchpad. Bei einer Mutation wurde das Werkzeug per `perl` an genau einer Stelle geändert; vor dem Lauf wurde geprüft, dass die Ersetzung gegriffen hat. Danach lief die Gegenprobe der Kopie. „Aus dem richtigen Grund“ heißt: Die Meldung `ROT — Fall '<name>'` nennt den Fall, der dem mutierten Vertragspunkt zugeordnet ist. Für Punkt 9 wurde statt des Skripts das Fragment `harness/mk/kopf-check.mk` mutiert. Für die neuen Fälle lief zusätzlich die Gegenprobe vom Stand `f0c87b6` (vor der Nacharbeit) gegen denselben Mutanten. Die Kopien sind weggeräumt.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — Spezifikation. **Bestätigt mit einer Lücke** (V-67).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| §11 *Harness-Werkzeuge* steht nach §10 *Nicht zugesichert in v1*, Historie ist §12 mit einer neuen Zeile | `spec/spezifikation.md:1834`, `:1843`, `:1917`, `:1975`; letzte Zeile der Historie ist die Zeile vom 2026-10-06 zu Harness-Werkzeugen | bestätigt |
| Bezug nach oben einmal im Einleitungsabsatz, ausdrücklich keine Zusage des Produkts | Einleitung mit Link auf [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (4), Satz „Keine Kennung dieses Abschnitts ist eine Zusage des Produkts …“ | bestätigt |
| `SPEC-047` trägt [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 1 bis 9 je einmal, in derselben Nummer | Abgleich Punkt für Punkt: Gegenstand (1), Kennungen (2), Bereich (3), Gleichheit (4), Kopf (5), §1/§2 (6), Formfehler (7), Ausgabe/Ausgang (8), Start ohne Stufung (9). Ohne Text übernommen sind nur Stellen ohne eigenen Vertragsgehalt: „Stubs eingeschlossen“ (gedeckt durch „flach und archiviert“), „was der Bestand meldet, berichtigt der Slice selbst“ (Vorgehen des damaligen Slice) und das Beispiel „ohne Punkt“. Das Beispiel deckt „als ganzes Wort“ ab, und `LH-FA-01a` in `nr2-schreibweise-abschnitt` hält es. Die Konsequenzen der ADR stehen unter *Grenze* | bestätigt |
| Die Lesarten (1) bis (5), (a) und (b) sind eingearbeitet | (1) in Nr. 2 („Ein Punkt beendet das Wort“), (2) in Nr. 1 („flach“, Unterverzeichnisse nicht), (3) in Nr. 5 (Leerzeile mit Leerzeichen und Tabs), (4) in Nr. 3 (Zeilenumbruch als Leerraum), (5) in Nr. 8 (Exit 2 nur ohne Ablage, fehlendes Lifecycle-Verzeichnis kein Abbruch), (a) in Nr. 5 (Absatzanfang), (b) in Nr. 8 (nicht lesbarer Plan). Jede einmal und dem Archiv-Wortlaut treu | bestätigt |
| Jede ZUSAGE-, GRENZE-, Ausgabe- und Ausgangs-Zeile des alten Kopfs von `kopf-check.sh` steht in `SPEC-047` | alle acht ZUSAGE-Zeilen, GRENZE, Ausgabe und Ausgang wiedergefunden. Genauer als vorher: Präfix `kopf-check: `, Befund-Texte wörtlich, Sortierung nach Befund in Bytefolge | bestätigt |
| `SPEC-048` trägt [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) und den alten Kopf von `abdeckung.sh` | Deklarationsform und Folgezeilen (1), Pfade je Anforderungsart und Lastenheft als Überschrift (2), Nachweisart nach Ort und `TestE2E…` (3), vier Tabellen und „vollständig“ (4), Schreiben und `--check` (5), Ausgänge (6). Aus der ADR: Tabellen je Nachweisart und Gesamtsicht, RTM liest nur Vollständiges, Prüfung ohne Schreiben (5), Negativ als *Grenze*. Die Kopplung an `make doc-trace` steht jetzt als Kopplungs-Kommentar im Code (`tools/test/abdeckung.sh:112`–`:113`) | bestätigt |
| Der Text beschreibt das Verhalten der Werkzeuge | Lesend gegen den Code beider Skripte geprüft. Keine Zusage geht über den Code hinaus. Sonden: Ein Bereich über die Grenze zwischen den Feldern `Bezug` und `Berührte Spec-Stellen` wird nicht aufgelöst, weil die Feldmarke dazwischen steht; also gibt es im Kopf kein ungenanntes „still grün“. `docs/user/` fehlt: siehe V-70 | bestätigt |
| Zu jeder Zusage nennt der Skriptkopf einen Fall, der sie hält | Jeder Fall-Name aus den beiden `GEPRÜFT DURCH`-Blöcken existiert in der Gegenprobe (mechanisch geprüft). Jeder Punkt 1 bis 9 und 1 bis 6 hat mindestens einen Fall, der unter einer Mutation aus dem richtigen Grund rot wird (Abschnitt 3). **Ausnahme:** Die *Grenze* von `SPEC-048` sagt Exit 1 „ohne eine Zeile auf stderr“ zu, wenn das Lastenheft keine Überschrift `### LH-…` führt. `lastenheft-ohne-anforderung` prüft nur den Exit (V-67) | Lücke |

### Punkt 2 — Sensor-Dateien. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Beide Dateien folgen der Vorlage | Überschriften von `harness/sensors/kopf-check.md` und `harness/sensors/abdeckung-check.md` gleich der Vorlage: Titelzeile, *Vertrag*, *Grenze — was das Grün nicht abdeckt*, *Ausgabe und Ausgänge*, *Sperren*, *Bindung*. Kein Status- und kein Datumsfeld | bestätigt |
| Vertrag als Link auf die Kennung, Bindung an die ADR | beide Links lösen auf (`make docs-check` 0 Befunde), Bindung [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) bzw. [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) | bestätigt |
| Sie entscheiden nichts, was der Abschnitt nicht sagt | Grenze: `kopf-check.md` die vier Grün-Grenzen von `SPEC-047`. Den Bereich über eine Absatzgrenze führt es zu Recht nicht, denn er macht strenger, nicht grün. Symlink und unlesbares Verzeichnis sind als offen benannt. Ausgänge stimmen mit Nr. 8 bzw. (6) und der *Grenze* überein, Befund-Texte als Verweis auf Nr. 8. *Sperren*: nur die benannte Abbruchzeile bzw. „Keine“ | bestätigt |
| Target-Zellen in `harness/README.md` §Sensors verlinken sie | `harness/README.md:71`, `:73`; die Zelle von `make kopf-check` ist ein Satz mit Vertrag als Link | bestätigt |

### Punkt 3 — Skriptköpfe und Code. **Bestätigt**, ergänzte Fälle mit der Lücke aus V-67.

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Der Kopf trägt den Verweis auf den Abschnitt und je Punkt den Fall | `tools/harness/kopf-check.sh:4`–`:22`, `tools/test/abdeckung.sh:3`–`:23`: Zweck, `Vertrag: spec/spezifikation.md SPEC-04x`, Sensor-Datei als Text, Aufruf, `GEPRÜFT DURCH` mit einer Zeile je Punkt und *Grenze* | bestätigt |
| Kein Vertragspunkt geht verloren | Abgleich in Punkt 1. Jede alte Zusage steht jetzt im Abschnitt, keine verblieb nur im Skript | bestätigt |
| Ausführbarer Code unverändert | `git diff e491e0c..c667677` der beiden Skripte nach `+`/`-`-Zeilen gefiltert, die keine Kommentarzeilen sind: leer. Die einzige Änderung im Rumpf (`abdeckung.sh:112`–`:113`) ist eine Shell-Kommentarzeile außerhalb jedes awk- oder Heredoc-Texts | bestätigt |
| `make kopf-check-gegenprobe` und `make abdeckung-gegenprobe` grün | in der Kopie und im `make gates`-Lauf grün; nicht als root, die Fälle `nr8-unlesbar-*` liefen also | bestätigt |
| Jeder ergänzte Fall ist gegen eine Mutation rot gesehen | Alle in Plan §6 *Stand* genannten Fälle wurden selbst gegen eine Mutation rot gesehen. Gegen dieselben Mutanten bleibt die Gegenprobe von `f0c87b6` grün: bei `check-zwei-veraltet`, `schreiben-fehlend`, `schreiben-unberuehrt`, den Bäumen ohne `Happy`/`Negative`, `verschachtelt`, `ohne-lastenheft`, `ohne-deklaration-zeile`, `lastenheft-ohne-anforderung` und `neg-bereich-*`. Die Fälle sind also wirklich neu. `lastenheft-ohne-anforderung` hält nur den Exit (V-67) | bestätigt bis auf V-67 |

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| `make gates` grün | bestätigt, Abschnitt 5 |
| Review mit Report | liegt vor (F-417 bis F-427); die Nacharbeit hat kein eigenes Review, ihre Wirkung prüft Abschnitt 2 |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen; §7 trägt Platzhalter, §6 *Risiken* vier Ausgänge „— (bei Closure)“. Das ist Sache der Closure und hier nicht zu bestätigen |

---

## 2. Review-Findings

| Finding | Status | Beleg |
|---|---|---|
| F-417 (`vollstaendig()` nur in `Boundary` rot) | **behoben** | Baum `inhalt` führt `LH-FA-04` ohne `Happy` und `LH-FA-05` ohne `Negative`. Mutationen `vollstaendig()` ohne `Negative`, ohne `Happy`, ohne `Boundary`: je rot in `inhalt-gesamt` und `inhalt-vollstaendig` |
| F-418 („je veralteter Tabelle“) | **behoben** | `check-zwei-veraltet`; Mutation `veraltet=1; break` rot |
| F-419 (Rechte 0644 für alle vier, Meldung auf stdout fehlt) | **behoben** | `SPEC-048` (5) sagt jetzt „nur die Tabellen, die abweichen oder fehlen“, Rechte nur für diese, Meldung je geschriebener Tabelle. Fälle `schreiben-*`; Mutationen „alle schreiben“, „ohne `chmod`“, „`chmod` auf alle“, „Meldung auf stderr“ je rot |
| F-420 (`.git/`, `.harness/` nur an der Wurzel) | **behoben** | `SPEC-048` (1) „an der Wurzel des Repos; ein gleichnamiges Verzeichnis tiefer im Baum wird gelesen“ (`c667677`); Plan §6 Lesart (2) nachgezogen; `verschachtelt` rot unter `-name .harness -prune` |
| F-421 (Bereich über Absatzgrenzen) | **behoben** | *Grenze* von `SPEC-047` nennt es, `neg-bereich-absatzgrenze` und `neg-bereich-regel-absatz` rot, wenn eine Leerzeile den Text trennt |
| F-422 (Ausgänge 2 und 1 ohne Zeile ungenannt) | **behoben**, mit Rest V-67 | `SPEC-048` (6) und *Grenze*, Ausgangs-Tabelle in `harness/sensors/abdeckung-check.md:31`–`:35`. Exit 2 hält `ohne-lastenheft`, Exit 1 hält `lastenheft-ohne-anforderung`; „ohne Zeile“ hält kein Fall |
| F-423 (Sensor-Dateien sagen mehr als der Vertrag) | **behoben** | (a) `SPEC-047` (8) zählt die Befund-Texte auf. (b) `abdeckung-check.md:38`–`:40` „Testdatei und Zeile (bei einem `TestE2E…` ohne Deklaration die Zeile des Tests)“ wie (6); `ohne-deklaration-zeile` rot bei Zeile − 1 |
| F-424 (Gegenprobe in `SPEC-047` (9)) | **behoben** | (9) nennt nur `make kopf-check` an der Gate-Kette; die Gegenprobe steht in der *Bindung* der Sensor-Datei als Tatsache und mit Zeiger auf ihre Zeile in §Sensors |
| F-425 (Folge-Slices halb nachgezogen) | **behoben bis auf den Kopf** (V-68) | (a) Coverage: §1, DoD Punkt 1, §3 (ADR-Zeile und neue Zeile `spec/spezifikation.md`), §6 *Die Zahl*; Mutation: Satz über §2, §3, §6 *Schwelle*. (b) Lint, Coverage, Mutation §6: „ein bestehendes Werkzeug schreibt seine Kennung fort, ein neues bekommt die nächste freie“; `grep` nach „Kennungen vergibt“ in `open/` leer. Der Kopf `Berührte Spec-Stellen` von Coverage und Mutation folgt der neuen §3-Zeile nicht |
| F-426 (`TestMain` mit Deklaration) | **behoben** | Plan §6 *Lesarten*, Grund (1): „Eine Deklaration direkt über `func TestMain` wird dennoch angenommen …“ |
| F-427 (Bezug auf LH-QA-07 knapp) | Hinweis, unverändert | kein Auftrag an den Implementer; bleibt Sache von Architect und Validator |

---

## 3. Proben

### Bewusst brechen, `kopf-check.sh` (Gegenprobe `c667677`)

| Punkt | Mutation | Rot in (Auszug) | Grund richtig |
|---|---|---|---|
| 1 | `-maxdepth 2` | `nr1-nicht-geprueft` (Unterordner gelesen) | ja |
| 1 | `in-progress` aus der Liste | `nr1-in-progress`, `nr8-sortiert` | ja |
| 2 | Unterstrich kein Wortzeichen | `nr2-schreibweise-abschnitt` (`LH-FA-01_x`) | ja |
| 2 | `.ab`-Regel gestrichen | `nr2-punkt-endet-wort` | ja |
| 3 | Bereich auch absteigend (Enden vertauscht) | `nr3-absteigend` | ja |
| 3 | Klassen-Gleichheit gestrichen | `nr3-gemischt` | ja |
| 3 | Zeilen in §1/§2 mit `\n` statt Leerzeichen verbunden | `nr3-bereich-zeilenumbruch-abschnitt`, `neg-bereich-*` | ja |
| 4 | Unterkennung trägt Hauptkennung mit ein | `nr4-unter-deckt-haupt` | ja |
| 5 | Leerzeile nur `^$` | `nr5-leerzeile-mit-leerzeichen`, `nr6-leerzeichen-nach-regel-absatz` | ja |
| 5 | Feldmarke überall Absatzanfang | `nr5-welle-vor-bezug`, `-titel-`, `-linie-`, `nr5-bezug-vor-stellen` | ja |
| 6 | „Regeln dieser Sektion“ an beliebiger Stelle | `nr6-regel-mitten-im-absatz` | ja |
| 6 | `^##` statt `^## ` | `nr6-unterueberschrift` | ja |
| 7 | Befund „Abschnitt ## 2. fehlt“ gestrichen | `nr7-ohne-2`, `nr7-leer` | ja |
| 8 | Sortierung des Befunds umgekehrt | `nr8-sortiert`, `nr7-leer`, `nr3-*` | ja |
| 8 | Befunde auf stdout | jeder `rot`-Fall („schreibt auf stdout“) | ja |
| 8 | Abbruchzeile umbenannt | `nr8-ohne-ablage` | ja |
| 8 | fehlendes Lifecycle-Verzeichnis → Exit 2 | `nr8-leere-ablage` | ja |
| 8 | nicht lesbarer Plan still übergangen | `nr8-unlesbar-davor`, `-zuletzt`, `nr8-awk-scheitert` | ja |
| 8 | Exit 2 statt 1 bei Befund | jeder `rot`-Fall | ja |
| 8 | Kennung je Vorkommen statt je Abschnitt einmal | `nr8-sortiert`, `nr2-hervorhebung` | ja |
| 9 | `kopf-check` aus `GATE_CHECKS` (Fragment) | `nr9-gate-checks` | ja |
| 9 | Rezept `… \|\| true` (Fragment) | `nr9-scharf` | ja |
| Grenze | Codeblöcke verfolgt | `neg-codeblock-ueberschrift` | ja |
| Grenze | Leerzeile trennt den Text in §1/§2 | `neg-bereich-absatzgrenze`, `neg-bereich-regel-absatz` | ja |

Äquivalent und darum ohne Aussage: `a < b` → `a != b`, denn die Schleife läuft bei a > b ohnehin nicht. Sie ist durch die Vertauschung oben ersetzt.

### Bewusst brechen, `abdeckung.sh` (Gegenprobe `c667677`)

| Punkt | Mutation | Rot in (Auszug) | Grund richtig |
|---|---|---|---|
| 1 | `-name .harness -prune` (überall) | `verschachtelt` | ja |
| 1 | `.git` nicht ausgenommen | `nicht-gelesen` | ja |
| 1 | Folgezeile ohne Leerzeichen angefügt | `gueltig-maskierung` | ja |
| 1 | eingerückte Deklaration angenommen | `eingerueckt` | ja |
| 1 | zweite Deklaration ohne Fehler | `doppelt` | ja |
| 1 | Dateiende ohne Fehler | `dateiende` | ja |
| 1 | fehlendes „ — “ still | `ohne-trenner` | ja |
| 1 | „ohne direkt folgenden Test“ gestrichen | `leerzeile` | ja |
| 2 | `LH-QA` mit `Happy` angenommen | `qa-happy` | ja |
| 2 | `LH-FA` mit `Messung` angenommen | `fa-messung` | ja |
| 2 | Lastenheft-Überschrift jeder Ebene | `ueberschrift-ebene` | ja |
| 2 | Lastenheft-Abgleich gestrichen | `unbekannt`, `ueberschrift-ebene` | ja |
| 3 | `test/integration/` als Unit | `inhalt-e2e`, `inhalt-unit`, `inhalt-gesamt`, `inhalt-vollstaendig` | ja |
| 3 | `TestE2E…`-Pflicht nur unter `test/integration/` | `ohne-deklaration-unit`, `verschachtelt` | ja |
| 4 | `vollstaendig()` ohne `Negative` / ohne `Happy` / ohne `Boundary` | je `inhalt-gesamt`, `inhalt-vollstaendig` | ja |
| 4 | Sortierung ohne Test (`-k5,5`) | `inhalt-e2e` | ja |
| 4 | `\|` nicht maskiert | `gueltig-maskierung` | ja |
| 4 | Komma ohne Leerzeichen; `n/a` fehlt; `—` fehlt | je `inhalt-gesamt` | ja |
| 4 | Teilweise in die Vollständig-Tabelle | `inhalt-vollstaendig` | ja |
| 5 | jede Tabelle geschrieben | `check-aktuell`, `schreiben-nur-geaendert`, `schreiben-unberuehrt` | ja |
| 5 | `chmod` gestrichen | `gueltig-rechte`, `schreiben-rechte` | ja |
| 5 | `chmod 0644` auf alle am Ende | `schreiben-unberuehrt` | ja |
| 5 | Meldung „geschrieben“ auf stderr | `schreiben-fehlend` | ja |
| 5 | `--check` meldet nur die erste veraltete | `check-zwei-veraltet` | ja |
| 5 | `--check` schreibt | `check-schreibt-nicht` | ja |
| 6 | Exit 3 bei Fehlform | jeder `abgelehnt`-Fall | ja |
| 6 | Meldung ohne Zeilennummer | jeder `abgelehnt`-Fall | ja |
| 6 | Zeile vor dem Test statt des Tests | `ohne-deklaration-zeile` | ja |
| 6 | Exit 2 bei veralteter Tabelle | `check-veraltet`, `check-zwei-veraltet` | ja |
| 6 | Lastenheft fehlt → Exit 1 | `ohne-lastenheft` | ja |
| Grenze | Lastenheft ohne Anforderung → Exit 0 | `lastenheft-ohne-anforderung` | ja |
| Grenze | Lastenheft ohne Anforderung → Exit 1 **mit** Zeile `abdeckung: keine Anforderung im Lastenheft` | — Gegenprobe **grün** | **Lücke, V-67** |

Äquivalent und darum ohne Aussage: `LH-QA`/`LH-RB` immer vollständig. Eine solche Anforderung steht nur mit einer Deklaration in der Tabelle, und ihr einziger zulässiger Pfad ist `Messung`.

### Referenz-Richtung und Matrix (`make docs-check` in Kopien)

| Probe | Eingriff in `spec/spezifikation.md` | Ergebnis |
|---|---|---|
| D0 | keiner | Exit 0, 259 Dateien, 0 Befunde |
| D1 | blanke ADR-Kennung in `SPEC-047` | Exit 2, `id-unlinked` `:1845` |
| D2 | Link auf eine ADR in `SPEC-047` | Exit 2, `matrix-forbidden` spec-straten → adr |
| D3 | Slice-Name blank in `SPEC-047` | Exit 2, `matrix-forbidden` (`slice-`) |
| D4 | Welle-Name blank in `SPEC-047` | Exit 2, `matrix-forbidden` (`welle-`) |
| D5 | Kennung des Adaptions-Blocks blank | Exit 2, `matrix-forbidden` spec-straten → adaptionsblock |
| D6 | Link auf `harness/sensors/kopf-check.md` | Exit 2, `matrix-forbidden` spec-straten → aussen |
| D7 | Slice-Name in der Historie-Zeile | Exit 2, `matrix-forbidden` (`slice-`) `:2012` |
| D8 | Skriptpfad als Code-Span ohne Link | Exit 0, grün (V-71) |

Im Bestand: Ein `grep` über §11 und die Historie-Zeile nach Skriptpfad, `Makefile`, `.sh`, `sensors`, Slice-, Welle-, ADR- und Adaptions-Kennung findet nichts. Der einzige Treffer ist `.harness/` als geprüfter Pfad in (1).

---

## 4. Plan gegen Stand

| Abschnitt | Ergebnis |
|---|---|
| Kopf | folgt: `Bezug` mit [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) und den drei ADRs; `Berührte Spec-Stellen` `spezifikation.md §11` (`SPEC-047`, `SPEC-048`) · `§12` stimmt mit dem Diff; `make kopf-check` Exit 0 |
| §1 | folgt: Ziel nennt Abschnitt, Sensor-Dateien, Skriptköpfe und ergänzte Fälle; Abgrenzung „kein ausführbarer Code der Werkzeuge“ hält (Abschnitt 1, Punkt 3) |
| §3 | folgt: jede Zeile hat ihren Diff. Die Gegenproben-Zeile nennt die geschärfte Erwartung und den Vertrag im Kopf |
| §6 *Randformen* | folgt bis auf einen Satz in *Was im Skriptkopf bleibt*: „der Skriptkopf ordnet ihr die grünen Fälle `neg-*` zu“. Seit `2ba6b31` sind `neg-bereich-*` rote Fälle (V-69) |
| §6 *Stand* und *Lesarten* | folgt; Lesart (2) auf die Wurzel berichtigt (`c667677`), Grund (1) um `TestMain` ergänzt |
| §6 *Risiken* | vier Risiken, Ausgang offen bis Closure. *Übertragung verliert eine Zusage* ist nach Abschnitt 1 nicht eingetreten; Restlücke V-67 |
| §4 | Rückführung `in-progress → next` nicht ausgelöst: Skript und ADR liegen nirgends im Widerspruch, die Lesarten sind bestätigt |
| §8 | Sichtung am Stand `3dfef35`/`e491e0c` vermerkt; keine Schwelle 3× mit diesem Slice allein |
| Folge-Slices | `slice-harness-abdeckung-gate`: Kopf `SPEC-048`, §6 *Ort* ohne neue Kennung, stimmig. `slice-harness-lint`: §6 *Kennungen* nachgezogen. `slice-harness-coverage` und `slice-harness-mutation`: §1/§2, §3, §6 nachgezogen, der Kopf nicht (V-68). Die vier Umstellungs-Slices nennen den Abschnitt beim Namen |

---

## 5. Gate-Lauf

`make gates` auf `c667677` im Repo, selbst gefahren: Exit 0, Laufzeit 62 s. Grün gemeldet: `abdeckung-gegenprobe`, `a-check` (0 Befunde), `a-check-negativ`, `baseline-verify` (v6.13.0, 54 Dateien), Build, `test`, `test-integration` (`run-integration-tests: gruen`), `docs-check` (259 Dateien, 0 Befunde), `commit-msg-gegenprobe`, `kopf-check`, `kopf-check-gegenprobe`, `abdeckung-check`. Der Nachweis-Stempel liegt unter dem ignorierten `.harness/state/`; `git status` danach leer. `make kopf-check` einzeln: Exit 0.

---

## 6. Befunde

| ID | Befund | Ort | Empfehlung |
|---|---|---|---|
| V-67 | **Die Grenze von `SPEC-048` sagt „ohne eine Zeile auf stderr“ zu, und kein Fall hält das.** Führt das Lastenheft keine Überschrift `### LH-…`, endet die Prüfung laut *Grenze* „mit 1 ohne eine Zeile auf stderr“. Die Sensor-Datei braucht genau dieses Merkmal, um die Ursache zu erkennen („Einen Ausgang 1 ohne Zeile gibt es, wenn …“). `lastenheft-ohne-anforderung` schickt stderr nach `/dev/null` und prüft nur den Exit. Ein Mutant, der in diesem Fall eine Zeile ausgibt und mit 1 endet, lässt die Gegenprobe grün. Das ist dieselbe Klasse wie F-417 und F-418: ein Vertragspunkt mit zwei Bedingungen, nur eine rot gehalten. Betroffen sind DoD Punkt 1 („zu jeder Zusage … der Fall“) und Punkt 3 („jeder ergänzte Fall … rot gesehen“) | `spec/spezifikation.md:1972`–`:1973`; `harness/sensors/abdeckung-check.md:34`, `:41`; `tools/test/abdeckung-gegenprobe.sh:107`–`:116`; `tools/test/abdeckung.sh:23` | Im Fall stderr auffangen und auf leer prüfen, dazu die Mutation oben rot sehen. Der Eingriff liegt nur in der Gegenprobe |
| V-68 | **Der Kopf von Coverage und Mutation folgt der neuen §3-Zeile nicht.** Seit `0c45064` führt §3 beider Pläne `spec/spezifikation.md` als „update“ (Kennung des Gates im Abschnitt für Harness-Werkzeuge). Ihr `Berührte Spec-Stellen` sagt weiter `—`, also „berührt keine Spec-Stelle“. `make kopf-check` ist grün, weil §1 und §2 keine Kennung nennen; die Kennung vergibt erst der Slice. Der Kopf dieses Slice zeigt, wie die Form ohne Kennung aussieht (`spezifikation.md §11`, neu) | `docs/plan/planning/open/slice-harness-coverage.md:27`, `:133`; `docs/plan/planning/open/slice-harness-mutation.md:23`, `:126` | Kopf auf `spezifikation.md §11` (neue Kennung des Gates) setzen (`AGENTS.md` §3.9) |
| V-69 | **Plan §6 *Was im Skriptkopf bleibt* nennt `neg-*` die „grünen Fälle“.** Seit `2ba6b31` gehören `neg-bereich-absatzgrenze` und `neg-bereich-regel-absatz` dazu, und beide sind rot. Der Kopf von `kopf-check.sh` ordnet der *Grenze* `neg-*` richtig zu; nur der Plantext ist stehen geblieben | `docs/plan/planning/in-progress/slice-harness-vertraege-spezifikation.md:371`–`:372` | „grünen“ streichen oder „die Fälle `neg-*`, grün bis auf den Bereich über Absatzgrenzen“ |
| V-70 | **Hinweis: zwei Ränder des Bestands, die `SPEC-048` nicht nennt.** (a) Fehlt `docs/user/`, endet `make abdeckung` mit Exit 1 und der Meldung von `cp`, ohne Fehlform; (6) kennt Exit 1 nur für Fehlform oder veraltete Tabelle. (b) „unter `test/integration/`“ heißt im Code „an der Wurzel“ (`case "$datei" in test/integration/*`). Das ist dieselbe Frage, die `c667677` für `.git/` und `.harness/` beantwortet hat. Der Abschnitt schreibt Pfade zwar durchweg von der Wurzel aus, und kein Mutant „auch tiefer im Baum“ wird rot. Beides sagt nichts zu, was das Werkzeug nicht hält; es lässt nur offen | `tools/test/abdeckung.sh:47`, `:182`–`:183`; `spec/spezifikation.md:1941`, `:1945` | Randform-Rückgabe an den Architect; Ort ist `slice-harness-abdeckung-gate`, das `SPEC-048` ohnehin fortschreibt |
| V-71 | **Hinweis: „kein Skriptpfad, auch nicht als Text“ hält kein Sensor.** Plan §6 *Was der Abschnitt nennen darf* verbietet Skriptpfade, Sensor-Dateien und Makefile im Abschnitt auch als Text. Probe D8 (Skriptpfad als Code-Span) bleibt in `make docs-check` grün, nur Links nach `aussen` fängt die Matrix. Der Bestand ist sauber (`grep`). Für die Folge-Slices, die in den Abschnitt schreiben, bleibt das Urteil des Reviews | `spec/spezifikation.md` §11; `.d-check.yml` `matrix` | keine Änderung nötig; im Review der Folge-Slices mitprüfen |

Kleinigkeiten ohne eigene Nummer, alle sachlich richtig, weil die Nummern von [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) und `SPEC-047` gleich sind: Die Schlusszeile der Gegenprobe (`tools/harness/kopf-check-gegenprobe.sh:340`), die Hilfe-Zeile in `harness/mk/kopf-check.mk:14` und die Zeile der Gegenprobe in `harness/README.md:74` sprechen weiter von „je Nummer der ADR“.

---

**Summary:** 3/3 Liefer-Punkte bestätigt, Punkt 1 und 3 mit der Lücke V-67 (Grenze „ohne Zeile auf stderr“ ohne roten Fall) · Werkzeug-Code außerhalb von Kommentaren unverändert · Proben: 24 Mutationen an `kopf-check` (Punkte 1 bis 9 und Grenze, alle rot aus dem richtigen Grund) und 36 an `abdeckung` (Punkte 1 bis 6 und Grenze), davon 35 rot aus dem richtigen Grund und eine grün (V-67), zwei äquivalent und ersetzt bzw. ausgewiesen; 7 von 7 Richtungs-Proben rot, D8 grün wie erwartet (V-71) · Review-Findings F-417 bis F-426 behoben, F-422 mit Rest V-67, F-425 mit Rest V-68 · `make gates` grün auf `c667677` · Befunde: V-67 (DoD), V-68, V-69 (Plan folgt nicht ganz), V-70, V-71 (Hinweise).
