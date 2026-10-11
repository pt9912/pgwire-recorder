# Slice slice-v1-abschluss-zeitangaben-uhr: Uhr-Port und Systemuhr

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen)

**Berührte Spec-Stellen:** `LH-FA-21.a` · `ARC-004` · `ARC-015` · `ARC-001`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-11.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Kern liest die Zeit und wartet auf einen Zeitpunkt nur über einen Driven Port „Uhr“ ([ADR-0015](../../adr/0015-uhr-port.md)): Der Port, eine Fake-Uhr für Tests, die Systemuhr im Composition Root und die Architekturregel, die den Kern an den Port bindet, liegen vor, bevor ein Slice sie benutzt. Beobachtbar ist der Vertrag an Tests (Warten bis zum Zeitpunkt, Abbruch durch den Kontext, Vorrücken der Fake-Uhr) und am Gate `make a-check`; ein Verhalten des Binaries ändert sich nicht. Die drei Testdateien des Kerns, die heute `time` einbinden, wechseln dazu auf eine Testhilfe (§3); ohne das ist die Regel rot.

**Übernommen aus `slice-v1-abschluss-zeitangaben`** (dort §3, Zeile *Kern, Uhr-Port* und DoD-Punkt 2, die Fake-Uhr): der Port und die Uhren. Die Folgepflicht aus [ADR-0015](../../adr/0015-uhr-port.md) (Architekturregeln halten den Kern frei von Zeitaufrufen der Standardbibliothek) ist hier der dritte Liefer-Punkt; sie war bisher kein Slice zugewiesen.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).

**Lieferwert:** Dieser Slice gibt den beiden Slices, die Zeit brauchen (`slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben`), einen gemeinsamen, getesteten Port, statt ihn in einen von beiden zu legen: Wer ihn in einen legte, brächte dort Kern und Bootstrap mit Logik zu den Schichten der Option dazu — drei statt zwei. Er ist ein Schnitt nach der gemeinsamen Voraussetzung, kein Schicht-Schnitt: Er liefert einen Vertrag mit Gegenprobe, den ein Gate hält.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Messen von `offset_ms` im Record und Warten im Einspielen — `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` (dort §1 und DoD, mit der Kennung dieses Slice): Dieser Slice liefert den Port, nicht seine Benutzer.
- Optionen und Konfiguration — dieselben beiden Slices: Der Port hat keine Option.
- Prüfung der Werte einer Aufzeichnung — `slice-v1-abschluss-zeitangaben-pruefen`: ein anderer Vorgang (Daten, nicht Zeit).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten. Dieser Slice ändert kein Verhalten des Binaries; das Handbuch hat nichts zu beschreiben.
- Eine Zeit, die der Kern aus der Aufzeichnung liest oder in sie schreibt (Zeitstempel in der Aufzeichnung) — Re-Evaluierungs-Trigger von [ADR-0015](../../adr/0015-uhr-port.md), nicht dieser Slice.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Driven Port „Uhr“ im Kern mit zwei Leistungen — die monotone Zeit lesen, bis zu einem Zeitpunkt warten, abbrechbar durch den Kontext — und eine Fake-Uhr, die Zeit nur auf Anweisung vorrückt ([ADR-0015](../../adr/0015-uhr-port.md)); Tests des Vertrags: Warten kehrt beim Erreichen des Zeitpunkts zurück, sofort bei einem Zeitpunkt in der Vergangenheit, mit dem Fehler des Kontexts beim Abbruch (vor und während des Wartens), und die Fake-Uhr weckt alle Wartenden, deren Zeitpunkt erreicht ist (Happy, Boundary, Negativ). Tests: `TestUhrWartetBisZumZeitpunkt`, `TestUhrErreichterZeitpunktSofort`, `TestUhrBeendeterKontextVorDemWarten`, `TestUhrAbbruchWaehrendDesWartens`, `TestUhrBeendeterKontextGehtVorBeimWecken`, `TestUhrWecktAlleErreichten`, `TestUhrRueckenNurAufAnweisung`, `TestUhrVorrueckenNegativPanik`. Ohne den Fix aus dem richtigen Grund rot: Die Verifikation baute die Mutanten je Zusage neu (zwölf der Fake-Uhr, alle rot mit Meldung).
- [x] Systemuhr im Composition Root (`internal/bootstrap`), die den Port erfüllt: Zeit monoton, Warten abbrechbar; getestet an der Systemuhr selbst (Untergrenze der Wartezeit, Abbruch), nicht am Binary. Tests: `TestSystemuhrMonoton`, `TestSystemuhrLesungAbgeschnitten`, `TestSystemuhrUrsprungBeiErzeugung`, `TestSystemuhrUrsprungTraegtMonotonenMesswert`, `TestSystemuhrWartenUntergrenze`, `TestSystemuhrErreichterZeitpunktOhneWarten`, `TestSystemuhrVergangenerZeitpunktSofort`, `TestSystemuhrBeendeterKontextVorDemWarten`, `TestSystemuhrAbbruchWaehrendDesWartens`, `TestSystemuhrFernerZeitpunktDrehtNicht`; die Verifikation sah neun Mutanten der Systemuhr rot.
- [x] Architekturregel in `.a-check.yml`, die den Kern an den Port bindet (a-check drückt sie aus, vom Architect am 2026-10-11 an einer Kopie des Baums gesehen: ein `tech`-Eintrag mit `match: regex`, Muster `^time$`, und der Liste der zugelassenen Orte — `internal/adapters`, `internal/bootstrap`, `cmd`, `test` und die Testhilfen unter `internal/`), mit Gegenprobe in `make a-check-negativ`: abgelehnt Domain Model, Ports und Services (je eine Probe, a-check nennt die Datei), zugelassen PGWire-Adapter und Bootstrap, dazu der Import von `runtime` im Domain Model zugelassen (Muster ist kein Teiltext); die drei Dateien der Tests von `internal/hexagon/services`, die `time` einbinden, benutzen stattdessen eine Testhilfe außerhalb des Kerns (`LH-FA-21.a` *Uhr*). Belege: `make a-check` ohne Befund, `make a-check-negativ` mit 23 Fällen grün (Fälle `model-time`, `ports-time`, `services-time` abgelehnt, `pgwire-time`, `bootstrap-time`, `model-runtime` zugelassen), `TestWartenInMillisekunden` und `TestFristLaeuftNachMillisekundenAb` für die Testhilfe; die Verifikation sah zehn Mutanten der Regel rot.
- [x] `make gates` grün (Stand `6ef6d17`, Lauf des Planners auf sauberem Baum vor den Closure-Commits; das Ergebnis nach dem letzten Commit steht im Bericht, ein Lauf belegt nur den Stand, auf dem er lief).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/ports/driven/uhr.go` | neu | Port „Uhr“: `Jetzt() int64` und `WarteBis(ctx, zeitpunkt int64) error`; die Namen wählte der Implementer, die Zusagen stehen in `LH-FA-21.a` *Uhr* |
| `internal/testuhr` (neu, Package außerhalb des Hexagons) | neu | Fake-Uhr (`Neu`, `Jetzt`, `WarteBis`, `Vorruecken`, `Wartende`), Vorrücken nur auf Anweisung um eine nicht negative Zahl (eine negative löst eine Panik aus); weckt alle Wartenden, deren Zeitpunkt erreicht ist; Ort und Zuschnitt entschieden vom Architect (2026-10-11) |
| `internal/testzeit` (neu, Package außerhalb des Hexagons) | neu | `Nach`, `Schlafe`, `NeueFrist` und `Abgelaufen`: Warten, Frist und Zeitgrenzen in ganzen Millisekunden für die Tests des Kerns, die heute `time` einbinden (`services/play_test.go`, `record_extended_test.go`, `record_zwangsende_test.go`); die einzige Stelle der Tests, die `time` einbindet |
| `internal/hexagon/services` (nur die drei Testdateien) | update | `time.After`/`time.Sleep`/`time.Now` durch `internal/testzeit` ersetzen, ohne Verhalten zu ändern; ohne das wird die Regel rot |
| `internal/bootstrap/systemuhr.go` | neu | Systemuhr mit Logik (Warten mit Abbruch), bereitgestellt als `NeueSystemuhr`; wird hier nicht in einen Service verdrahtet, denn kein Service nimmt noch eine Uhr |
| `.a-check.yml`, `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, `harness/README.md` | update | Regel „Kern bindet `time` nicht ein“ nach [ADR-0015](../../adr/0015-uhr-port.md), Fitness Function; Fallzahl und Beschreibung der Gegenprobe nachziehen; die Regel prüft den Import, nicht den Aufruf |
| `internal/bootstrap/export_test.go`, `internal/testuhr/doc.go`, `internal/testzeit/doc.go` | update / neu | Zugriff des Tests auf den Ursprung der Systemuhr (`UrsprungMonoton`); Paketkommentare der beiden Testhilfen |
| Tests | neu | Vertrag des Ports an der Fake-Uhr (`internal/testuhr/testuhr_test.go`) und an der Systemuhr (`internal/bootstrap/systemuhr_test.go`, `systemuhr_cpu_unix_test.go`), die Testhilfe (`internal/testzeit/testzeit_test.go`), Gegenprobe der Regel (`tools/arch/a-check-negativ.sh`) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Voraussetzung: `slice-v1-abschluss-einspielen` ist `done` (Schritt 8 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)); WIP-Limit 1. Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar, oder der Slice braucht eine dritte Schicht (Systemuhr in einem Adapter statt im Bootstrap); die Gegenprobe der Regel ist dann ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): a-check lehnt `time` im Kern trotz der Konfiguration nicht ab oder lehnt die zugelassenen Orte ab (der Architect sah die Regel am 2026-10-11 an einer Kopie arbeiten); dann zuerst die Entscheidung, ob ein Carveout (Modul 7, [ADR-0015](../../adr/0015-uhr-port.md) Folgepflicht offen) oder ein Slice am Werkzeug sie trägt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — vom Architect vor dem Code am 2026-10-11 entschieden, Ort: [`LH-FA-21.a`](../../../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Uhr*. Entscheidungen des Architect (Empfehlung, bis der Nutzer anderes sagt; der Nutzer hat am 2026-10-11 verfügt, ohne Nachfrage weiterzumachen).

- *Einheit und Typ im Port:* ganze Millisekunden als `int64` seit einem Ursprung, den die Uhr wählt; kein Typ der Zeitbibliothek im Port, abgeschnitten, nie kleiner als die vorige Lesung — entschieden ebenda. Die Namen der Leistungen wählt der Implementer, die Zusagen nicht.
- *Zeitpunkt in der Vergangenheit oder jetzt:* Warten kehrt mit `nil` sofort zurück — entschieden ebenda; ein beendeter Kontext geht vor (Fehler des Kontexts auch bei vergangenem Zeitpunkt).
- *Abbruch:* Kontext endet vor oder während des Wartens; der Aufruf kehrt mit dessen Fehler zurück, kein neuer Code — entschieden ([ADR-0015](../../adr/0015-uhr-port.md), `LH-FA-21.a` *Uhr*).
- *Monotone Zeit:* die Systemuhr zählt aus der monotonen Zeit des Prozesses (Ursprung bei ihrer Erzeugung), nie aus der Wanduhr — entschieden ebenda.
- *Mehrere Wartende gleichzeitig:* die Fake-Uhr weckt jeden, dessen Zeitpunkt erreicht ist; Vorrücken um eine negative Zahl ist ein Programmierfehler der Tests — entschieden ebenda.
- *Ferner Zeitpunkt:* ein Zeitpunkt, dessen Abstand zur Lesung die größte Dauer der Zeitbibliothek übersteigt (`math.MaxInt64`), und einer unterhalb jeder Lesung (`math.MinInt64`): Warten kehrt bei `MinInt64` sofort mit `nil` zurück, bei `MaxInt64` wartet es, bis der Kontext endet; die Systemuhr wartet dazu in Stufen dieser Dauer. Entschieden: `LH-FA-21.a` *Uhr* (warten, bis die Lesung den Zeitpunkt erreicht hat; vergangener Zeitpunkt sofort).
- *Gleichzeitig Zeitpunkt erreicht und Kontext beendet:* der Fehler des Kontexts geht vor, in der Systemuhr und in der Fake-Uhr — entschieden ebenda (*ein beendeter Kontext geht vor*).
- *Regel:* Die Regel greift am Import (`tech` mit `^time$`), nicht am Aufruf; sie lehnt damit auch einen Zeittyp ohne Aufruf ab, was der Port nicht braucht (akzeptiertes Negativ, [ADR-0015](../../adr/0015-uhr-port.md) sagt „Zeitaufrufe“: die Regel ist enger gefasst als nötig, nicht weiter). Auch die Tests des Kerns binden `time` nicht ein: a-check lehnt sie für Domain Model, Ports und Services ab, und sie auszunehmen lockerte die übrigen Regeln für diese Dateien (`AGENTS.md` §3.6); darum die Testhilfe `internal/testzeit`. Entschieden ebenda.

**Risiken:**

- Die Regel greift am Import und verbietet mehr, als [ADR-0015](../../adr/0015-uhr-port.md) wörtlich sagt (akzeptiertes Negativ, Grund oben) — **Ausgang:** offen bis Closure (erwartet: entfallen, mit diesem Grund).
- Die Migration der drei Testdateien ändert deren Zeitverhalten (Frist, Schlafzeit) und macht einen Test flackern — **Ausgang:** offen bis Closure; Gegenmittel: dieselben Werte in Millisekunden, kein Test wird umgeschrieben außer an den Aufrufen von `time`.
- Der Port steht vor seinem ersten Benutzer, und der Benutzer braucht eine andere Signatur (`BEO-REPO/folge-slice-adresse-nimmt-nicht-an`, verkörpert in `AGENTS.md` §3.13) — **Ausgang:** offen bis Closure; Gegenmittel: `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` nennen die Signatur in §1, nicht eine eigene.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Geschrieben bei der Closure (Planner, 2026-10-12, Stand `6ef6d17`).

- **Belege zur DoD:** die Testverweise stehen in der DoD; *Belege des Implementers* unten tragen Zusage, Mutation und roten Test. Die Verifikation (`docs/reviews/2026-10-12-verifikation-slice-v1-abschluss-zeitangaben-uhr.md`) bestätigt alle drei Liefer-Punkte: Die Mutanten je Zusage (25 am Code, 10 an der Regel) sind unabhängig neu gebaut und rot aus dem richtigen Grund, die Gegenprobe zählt 23 Fälle mit wahren Beschreibungen, 75 Läufe der neuen Tests bei 60 Rechenschleifen ohne Fehlschlag. Namensabgleich des Planners: jeder Testname der DoD hat genau eine Zeile `^func <Name>(` unter `internal/`. `make gates` lief grün am Stand `6ef6d17` (Planner, auf sauberem Baum vor den Closure-Commits); nach dem letzten Commit steht das Ergebnis im Bericht.
- **Was hat funktioniert:** Der Schnitt nach der gemeinsamen Voraussetzung trug: Port, Fake-Uhr, Systemuhr und Regel kamen als Vertrag mit Gegenprobe, ohne dass ein Benutzer den Zuschnitt beeinflusste; drei Liefer-Punkte, zwei Schichten, in einer Sitzung prüfbar. Der Architect sah die Regel vor dem Code an einer Kopie des Baums arbeiten, und die Gegenprobe (23 Fälle) fing jeden Mutanten der Regel. Die Migration der drei Testdateien änderte keinen Wert (30000, 2000, 100, 50, 1 ms) und flackerte nicht.
- **Was ging anders als geplant:**
  1. Die Aussage über das Verhalten von a-check („Dateien in keiner Schicht prüft a-check nicht, die Einträge sind redundant“) war falsch für die Regel `tech`: Ohne die Einträge meldet a-check 22 `tech-leak` (F-640, F-641); gemessen wurde erst im Review.
  2. Ein Mutant der Grenze der Systemuhr (`<` statt `<=`) blieb grün (F-642), der Test kam in `a866965`.
  3. Das Gegenmittel zu Risiko 3 — die Nehmer nennen die Signatur — stand nirgends als Handlung, bis Review (F-643) und Verifikation (V-186) es fanden; der Planner trug die Namen bei der Closure ein.
  4. §3 nannte drei Dateien des Diffs nicht (F-645 (4)); der ferne Zeitpunkt (`math.MaxInt64`, `math.MinInt64`) stand im Code und erst nach dem Code in §6, entschieden hat ihn die Spezifikation, der Code nichts Neues (F-645 (5), nicht als Beleg zu `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` gezählt).
- **Summary-Zeilen:** Review `docs/reviews/2026-10-11-review-slice-v1-abschluss-zeitangaben-uhr.md`: F-640 (MEDIUM): Aussage über das Gate ohne Gegenprobe, falsch. F-641 (LOW): Fall `bootstrap-time` trägt die beschriebene Regel nicht. F-642 (LOW): Mutant `<` der Systemuhr grün. F-643 (LOW): Nehmer nennen die Namen nicht. F-644 bis F-646 (INFO). Verifikation: „0 HIGH · 0 MEDIUM · 1 LOW · 2 INFO. V-186: F-643 offen. V-187: `TestRunPlayZweitesSignal` fällt unter Last aus, auch vor dem Slice. V-188: Mutant ‚null ist erlaubt‘ rot durch Panik des Testbinarys.“ Ausgänge: F-640 und F-641 durch `29caf1e`; F-642 durch `a866965`; F-643 und V-186 durch `f55efe1`; F-644 durch `a866965` (Fake-Uhr liefert nach dem Wecken den Fehler des beendeten Kontexts); F-645 (4) durch `7127433`; F-646 und V-188 ohne Aktion (Zeitfenster bei Last 75 Mal ohne Flackern, Mutant rot aus dem richtigen Grund); V-187 unten.
- **Entscheidungen des Planners zur Verifikation:**
  - *V-187, `TestRunPlayZweitesSignal` unter Last:* Neuer Eintrag `BEO-REPO/test-faellt-unter-cpu-last-aus` (1×, offen), mit den Zahlen der Verifikation (3 von 15 am Stand `7127433`, 2 von 12 am Stand `aa07a01`). Keine Adresse: `slice-harness-lint-warten-ohne-frist` prüft Empfänge ohne Frist, der Test trägt eine Frist von 30 s; eine Sendung dort hat kein Ziel, das er führt. Die Ursache ist nicht untersucht; ob sie im Test oder in einer Wettlaufbedingung des Play-Service liegt, ist offen. `slice-v1-abschluss-zeitangaben` ändert diese Signalbehandlung, eine Pflicht entsteht dort nicht.
  - *V-186:* erledigt in `f55efe1`, mit Nachzählung beider Nehmer (§8 dort): `-aufzeichnen` 2 Liefer-Punkte und 2 Schichten, `-zeitangaben` 3 und 2; der Bootstrap ruft nur `NeueSystemuhr()` und reicht das Ergebnis (Verdrahtung, `AGENTS.md` §3.13). Kein Schnitt, keine andere Adresse.
- **Steering-Loop-Eintrag:** Geschärfte Regel: Nennt §6 eines Gebers als Gegenmittel, dass Nehmer eine Signatur nennen, deren Namen der Implementer wählt, trägt der Planner die Namen bei der Closure des Gebers in die Nehmer ein, mit Nachzählung im selben Commit; das Risiko hat keinen Ausgang, solange `grep` nach den Namen im Nehmer leer bleibt. Anlass: F-643 und V-186; das Gegenmittel war geplant und hatte keinen Träger. Liegt in `.claude/agents/planner.md` (seit slice-v1-abschluss-zeitangaben-uhr); gezählt unter `BEO-REPO/plan-folgt-korrektur-nicht`.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist wieder aufgetreten (F-643, F-645 (4)); die Regel bleibt. §3.10 (seit welle-extended-query) ist wieder aufgetreten (F-642); die Regel bleibt, der Sensor ist mit `slice-harness-mutation` geplant. §3.11 ist in der Hälfte *Prosa über ein Gate* (seit slice-harness-d-check-v0-85) wieder aufgetreten (F-640, F-641); die Regel bleibt. §3.12 und §3.13 sind nicht wieder aufgetreten: Die Randformen standen in §6 vor dem Code (der ferne Zeitpunkt entscheidet die Spezifikation), die Nehmer nahmen an.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `6ef6d17` (Zähler = Dateien unter `evidence/`); jeder Beleg trägt `evidence/slice-v1-abschluss-zeitangaben-uhr.md`.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 39× → 40× (F-640, F-641), Stand verkörpert, bleibt.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 30× → 31× (F-643, F-645 (4)), Stand verkörpert, mit der Regelzeile in `.claude/agents/planner.md`.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 27× → 28× (F-642), Stand verkörpert, Sensor `slice-harness-mutation` geplant, bleibt.
  - Neu: `BEO-REPO/test-faellt-unter-cpu-last-aus` (1×, V-187), Stand offen.
  - Weiter offen, kein Auftreten: Risiko 3 unten, als Zeile im Stand von `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt 7×).
  - Gesichtet ohne Beleg: `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (bleibt 2×: F-642 gilt einem fehlenden Test an einem neuen Vertrag, nicht dem Umbau im Bestand; als Beleg dort zählte es eine andere Klasse mit), `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (bleibt 2×; der als äquivalent eingestufte Mutant `return nach dem Timer` ist von der Verifikation neu gebaut und bestätigt), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (bleibt 2×; die Rückführungen aus §4 traten nicht ein, zwei Schichten), `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (bleibt 1×; der Port berührt das Replay nicht), `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (bleibt 2×; F-640 liest die Konfiguration von a-check ebenfalls anders, als sie wirkt, und ist als Aussage im Kommentar unter `zusage-im-kommentar-weiter-als-pruefung` gezählt: ein Fund zählt in einer Klasse, sonst stiege der Eintrag mit diesem Slice auf 3× und der Ausgang fiele auf eine Zählweise).
  - Einmalig und nicht eingetragen: F-644 (Gleichzeitiges Wecken und Kontextende in der Fake-Uhr, Testhilfe, behoben), F-646 und V-188 (ohne Aktion).
- **Folge-Slices:** keiner neu. Der nächste Schritt nach §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) ist `slice-v1-abschluss-zeitangaben-aufzeichnen` (Schritt 21); die beiden Nehmer des Ports tragen die Signatur in §1 und DoD.
- **Risiken aus §6:** drei. *Die Regel greift am Import und verbietet mehr als [ADR-0015](../../adr/0015-uhr-port.md) wörtlich sagt*: **entfallen** — die Eigenschaft ist keine offene Frage mehr, sondern entschieden und gezeigt: §6 *Regel* nennt sie als akzeptiertes Negativ, der Kommentar in `.a-check.yml` ebenso; der Kern braucht keinen Zeittyp (der Port rechnet in `int64`), `make a-check` ist ohne Befund, und der Fall `model-runtime` zeigt, dass die Regel nicht mehr als `time` trifft. *Die Migration der drei Testdateien lässt einen Test flackern*: **entfallen** — die Werte sind gleich geblieben (Diff der Verifikation), und die Pakete liefen unter Last (15 Durchgänge bei 40 Schleifen, alle Pakete einschließlich `internal/hexagon/services`, kein Fehlschlag in den migrierten Dateien; der einzige Ausfall ist `TestRunPlayZweitesSignal` in `internal/bootstrap`, den der Slice nicht berührt, V-187). *Der Port steht vor seinem ersten Benutzer, und der Benutzer braucht eine andere Signatur*: **weiter offen** → Register, als Zeile im Stand von `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`; das Gegenmittel ist umgesetzt (`f55efe1`), das Risiko erst mit dem Schluss beider Nehmer beurteilbar. Eingetreten ist keines.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt das Feld `liegt in` mit `.claude/agents/planner.md`; der Absatz steht dort mit dem Anker `seit slice-v1-abschluss-zeitangaben-uhr`, und `BEO-REPO/plan-folgt-korrektur-nicht/state.md` nennt Zielort und Anker. Folge-Slice: Die Nehmer `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liegen in `next/`; `grep -n "slice-v1-abschluss-zeitangaben-uhr"` findet die Kennung des Gebers in §1 beider (Übernahme und *Ausdrücklich NICHT*) und in ihrem Start, der den Slice in `done/` verlangt; ihr §1 schließt die Sendung nicht aus. Register: Die drei Einträge mit Beleg und der neue tragen `evidence/slice-v1-abschluss-zeitangaben-uhr.md`; `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` besteht als Verzeichnis mit nicht leerem `evidence/`; die übrigen genannten bestehen ebenso. Die nächste Welle-Closure prüft die Paarungen erneut.

### Belege des Implementers

**Läufe** (Stand: der Baum, der diesen Plan mit diesem Abschnitt trägt): `make gates` grün; `make test-integration` grün. Mutanten liefen in frischen Kopien des Baums (`cp`/`tar` ohne mtime, Test mit `go test` im Image `pgwire-recorder:test` ohne Netz, Kopie als Bind-Mount), nie im Repo; die Mutanten der Regel liefen mit dem Skript der Gegenprobe (`tools/arch/a-check-negativ.sh`) in Kopien.

**Zusage · Mutation · roter Test** (Wortlaut `LH-FA-21.a` *Uhr*):

| Zusage | Mutation (nur sie) | roter Test |
|---|---|---|
| Fake: rückt nur auf Anweisung vor | `Jetzt` addiert die Realzeit seit Start | `TestUhrRueckenNurAufAnweisung` (`Jetzt nach Wartezeit ohne Vorruecken: 20`) |
| Fake: rückt nur auf Anweisung vor | `WarteBis` setzt die Uhr auf den Zeitpunkt | `TestUhrWartetBisZumZeitpunkt` (kehrt bei Jetzt 4 zurück) |
| Fake: um die genannte Zahl | `jetzt = ms` statt `+=` | `TestUhrRueckenNurAufAnweisung` (3 statt 10) |
| Fake: negative Zahl ist Programmierfehler | Panik entfernt | `TestUhrVorrueckenNegativPanik` |
| Fake: null ist erlaubt (Grenze) | Panik schon bei `ms <= 0` | `TestUhrRueckenNurAufAnweisung` |
| Fake: weckt jeden Wartenden, dessen Zeitpunkt erreicht ist (a: jeden) | weckt nur den ersten | `TestUhrWecktAlleErreichten` (zweiter Wartender mit 5 kehrt nicht zurück) |
| Fake: (b: dessen Zeitpunkt erreicht ist, gleich) | `<` statt `<=` beim Wecken | `TestUhrWartetBisZumZeitpunkt`, `TestUhrWecktAlleErreichten` |
| Fake: (c: nur diejenigen) | weckt alle | `TestUhrWartetBisZumZeitpunkt` (kehrt bei Jetzt 4 zurück), `TestUhrWecktAlleErreichten` (Zeitpunkt 10 bei Jetzt 7) |
| Warten: erreichter Zeitpunkt sofort | `<` statt `<=` beim Prüfen | `TestUhrErreichterZeitpunktSofort` |
| Warten: beendeter Kontext geht vor, auch bei vergangenem Zeitpunkt | Vorab-Prüfung des Kontexts entfernt | `TestUhrBeendeterKontextVorDemWarten` (`WarteBis(-1)` liefert nil) |
| Warten: Fehler des Kontexts, vor dem Warten | Vorab-Prüfung liefert immer `context.Canceled` | `TestUhrBeendeterKontextVorDemWarten` (Frist abgelaufen) |
| Warten: Fehler des Kontexts, während des Wartens | Abbruch liefert `context.Canceled` | `TestUhrAbbruchWaehrendDesWartens` (Fristablauf) |
| Warten: Abbruch kehrt mit dem Fehler zurück (nicht nil) | Abbruch liefert nil | `TestUhrAbbruchWaehrendDesWartens` |
| Warten: ein Abgebrochener wartet nicht mehr | Abbruch trägt den Wartenden nicht aus | `TestUhrAbbruchWaehrendDesWartens` (`Wartende nach Abbruch: 1`) |
| Systemuhr: Ursprung bei ihrer Erzeugung | Ursprung eine Stunde zurück | `TestSystemuhrUrsprungBeiErzeugung`, `TestSystemuhrLesungAbgeschnitten` |
| Systemuhr: Ursprung bei Erzeugung (Gegenform) | Lesung aus der Wanduhr (`UnixNano`) | `TestSystemuhrUrsprungBeiErzeugung` |
| Systemuhr: abgeschnitten | Lesung gerundet statt abgeschnitten | `TestSystemuhrLesungAbgeschnitten` (`Lesung 8 ms über der Obergrenze 7 ms`) |
| Systemuhr: nie kleiner als die vorige | Lesung plus 0 oder 1 nach Zufall | `TestSystemuhrMonoton` (`Lesung 1 nach 2`) |
| Systemuhr: aus der monotonen Zeit | Ursprung ohne Messwert der monotonen Zeit (`Round(0)`) | `TestSystemuhrUrsprungTraegtMonotonenMesswert` |
| Systemuhr: Warten nicht vor dem Zeitpunkt | halbe Wartezeit und kein Nachlesen nach dem Timer | `TestSystemuhrWartenUntergrenze` (`WarteBis(51) kehrte bei Lesung 26 zurück`) |
| Systemuhr: erreichter oder vergangener Zeitpunkt sofort | Sofort-Rückkehr entfernt | `TestSystemuhrWartenUntergrenze`, `TestSystemuhrVergangenerZeitpunktSofort` (hängt) |
| Systemuhr: vergangener Zeitpunkt, auch `math.MinInt64` | Abstand vor dem Vergleich gebildet (Überlauf) | `TestSystemuhrVergangenerZeitpunktSofort` (`WarteBis(-9223372036854775808)` hängt) |
| Systemuhr: beendeter Kontext geht vor | Prüfung des Kontexts entfernt | `TestSystemuhrBeendeterKontextVorDemWarten` |
| Systemuhr: Fehler des Kontexts, während des Wartens | Abbruch liefert `context.Canceled` / nil | `TestSystemuhrAbbruchWaehrendDesWartens` (Fristablauf / Abbruch) |
| Systemuhr: Abbruch beendet das Warten | `ctx.Done` ignoriert | `TestSystemuhrAbbruchWaehrendDesWartens` (hängt) |
| Systemuhr: ferner Zeitpunkt (Abstand über der größten Dauer) wartet bis zum Abbruch, ohne zu drehen | Begrenzung der Stufe entfernt | `TestSystemuhrFernerZeitpunktDrehtNicht` (`350.757ms CPU-Zeit`); mit zusätzlich entferntem Nachlesen auch `TestSystemuhrAbbruchWaehrendDesWartens` (nil statt Fehler) |
| Testhilfe: ganze Millisekunden | Mikrosekunden | `TestWartenInMillisekunden` |
| Testhilfe: Frist läuft nach ihren Millisekunden ab | Frist sofort / nie abgelaufen | `TestFristLaeuftNachMillisekundenAb` (zwei Mutanten) |
| Regel: `time` im Domain Model, in den Ports und in den Services abgelehnt | Eintrag aus `.a-check.yml` entfernt | `make a-check-negativ`: Fälle `model-time`, `ports-time`, `services-time` rot-erwartet, a-check grün |
| Regel: Muster ist kein Teiltext | Muster `time` statt `^time$` | `make a-check-negativ`: Fall `model-runtime` |
| Regel: Muster ist ein regulärer Ausdruck | `match: regex` entfernt | `make a-check-negativ`: die drei Fälle der Ablehnung |
| Regel: Adapter zugelassen | `internal/adapters` aus der Liste entfernt | `make a-check`: 4 `tech-leak` in den Adaptern; `make a-check-negativ` in der Kopie: alle grün-erwarteten Fälle rot, darunter `pgwire-time` |
| Regel: Integrationstests und Testhilfen zugelassen, jeder Eintrag trägt | je Eintrag einzeln entfernt (`test/`, `internal/testuhr`, `internal/testzeit`, `internal/testpki`) | `make a-check` (Kopie): 17 `tech-leak` unter `test/integration`; 1 in `internal/testuhr/testuhr_test.go`; 2 in `internal/testzeit`; 1 in `internal/testpki/testpki.go` |
| Systemuhr: erreichter Zeitpunkt kehrt ohne Warten zurück | `<` statt `<=` bei „Zeitpunkt erreicht“ | `TestSystemuhrErreichterZeitpunktOhneWarten` (400 Aufrufe brauchen 400 ms) |
| Fake: ein beendeter Kontext geht vor, auch wenn `Vorruecken` den Zeitpunkt zugleich erreicht | nach dem Wecken `nil` statt `ctx.Err()` | `TestUhrBeendeterKontextGehtVorBeimWecken` (Kontext, dessen `Err` erst nach dem Eintritt endet, damit der Fall deterministisch ist) |
| Regel: die Tests des Kerns binden `time` nicht ein | `_ "time"` in `record_zwangsende_test.go` | `make a-check` (`app-impurity`, nennt die Datei) |

**Grüne Mutanten:**

- *Nachlesen nach dem Timer entfernt, ohne halbe Wartezeit* (`return nil` nach `timer.C`): grün, äquivalent. Ein Go-Timer läuft nicht früher ab als seine Dauer; die Lesung ist ab der gelesenen Lesung um mindestens `rest` Millisekunden gewachsen und damit abgeschnitten nicht kleiner als der Zeitpunkt. Die Gegenform (mit halber Wartezeit) ist rot, siehe Tabelle.
- *`internal/bootstrap` oder `cmd` in der Liste der zugelassenen Orte* (die Einträge der ersten Fassung): grün, `composition_root` nimmt beide von der Regel aus; beide Einträge sind entfernt. Gemessen ist das Gegenteil für alle übrigen Einträge: Jeder einzeln entfernt wird rot (Tabelle). Dass a-check Dateien ohne Schicht ungeprüft lässt, gilt für Schichten, nicht für `tech`: die Liste ist ein Muster auf dem Pfad, und jeder Eintrag steht wegen Dateien, die `time` einbinden. Fall `bootstrap-time` der Gegenprobe hält deshalb den Zustand fest (Composition Root nicht abgelehnt), nicht einen Eintrag; ohne die Ausnahme in `composition_root` wird der Baum rot.
- *Begrenzung der Stufe entfernt, Nachlesen bleibt*: war zuerst grün (der Zeitpunkt `math.MaxInt64` endet am Abbruch des Kontexts, nur die Rechnerlast unterscheidet sich); über die CPU-Zeit fangbar, `TestSystemuhrFernerZeitpunktDrehtNicht`.

- *`zeitpunkt < jetzt` statt `<=`* in der Systemuhr war grün (Verzögerung höchstens eine Millisekunde: Timer mit Dauer 0, nächste Lesung); jetzt rot, siehe Tabelle.
- *Simultanes Wecken und Kontextende in der Fake-Uhr* lieferte zufällig `nil` oder den Fehler (`select`); jetzt liefert der Aufruf nach dem Wecken den Fehler des Kontexts, wenn er beendet ist; Test deterministisch über einen Kontext, dessen `Err` erst nach dem Eintritt endet.

**Grenze:** Dass die Systemuhr einer Änderung der Wanduhr nicht folgt, prüft kein Test; geprüft ist, dass der Ursprung den Messwert der monotonen Zeit trägt (`TestSystemuhrUrsprungTraegtMonotonenMesswert`), und dass die Lesung nicht aus der Wanduhr stammt (Mutant oben). Ein Test, der die Wanduhr verstellt, ist im Container nicht möglich.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand des Schnitts, 2026-10-11). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) — der Port ist ein neuer Vertrag, jede Zusage eine Mutation.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×, offen) — betrifft Zeitverhalten im Replay, nicht den Port; keine Neuberührung.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) — dieser Schnitt lässt keine Hälfte an der Grenze: drei kleine Liefer-Punkte, zwei Schichten; die Größe misst der Architect vor dem Code.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 3 (die Migration der Tests gehört zum dritten Punkt, ohne sie ist die Regel rot). Schichten: zwei (Kern mit dem Port und seinen Tests, Bootstrap mit der Systemuhr — Logik: Warten mit Abbruch zählt, `AGENTS.md` §3.13); das Gate (`.a-check.yml`) ist ein Liefer-Punkt, keine Produkt-Schicht, wie bei `slice-v1-abschluss-sqlite-format`. Nachgezählt beim Schnitt und vom Architect vor dem Code (2026-10-11), beim Umsetzen unverändert: drei Liefer-Punkte, zwei Schichten (Kern, Bootstrap). Sitzt die Systemuhr in einem Adapter statt im Bootstrap, sind es zwei Schichten mit Adapter und Kern und derselbe Befund.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
