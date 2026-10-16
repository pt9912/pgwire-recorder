# Slice slice-v1-abschluss-verbindungsende-doku: Handbuch und README zu: Verbindungsende durch den Server

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-verbindungsende-format` und `slice-v1-abschluss-verbindungsende-format-schreiber` bis `slice-v1-abschluss-verbindungsende-upstream`)

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

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das die Slices `slice-v1-abschluss-verbindungsende-format`, `-format-schreiber`, `-sqlite`, `-sqlite-schreiber`, `-replay`, `-replay-extended`, `-replay-extended-flush`, `-einspielen`, `-einspielen-extended`, `-einspielen-grenzen`, `-vergleich`, `-hilfetext-binary`, `-record`, `-upstream` und `-upstream-extended` liefern, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice der sieben Code-Slices des Verbindungsendes, gebündelt hinter dem letzten, nach Entscheidung des Nutzers vom 2026-10-12 (Vorschlag des Planners; Begründung: Ein Anwender kann ein Ende erst mit `record` erzeugen; bis dahin kennt kein Weg des Binaries eine Aufzeichnung mit Ende außer handgeschriebenen Dateien, die das Handbuch nicht beschreibt. Der Doku-Teil ist ungeteilt, wie bei `slice-v1-abschluss-zeitangaben-doku` und `slice-v1-abschluss-sqlite-format-doku`.) Der Teil, den die Code-Slices in DoD und §3 trügen, liegt hier, weil die Dokumentation eine eigene Schicht ist (`AGENTS.md` §3.13). Zwischen den Code-Slices und diesem Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen. Die Welle schließt erst mit diesem Slice (§3.11).

**Handbuch-Lücke akzeptiert** (Entscheidung des Nutzers vom 2026-10-12, Vorschlag des Planners): Zwischen `slice-v1-abschluss-verbindungsende-vergleich` und diesem Slice stimmt der Satz des Handbuchs zu `FATAL`/`PANIC` (aus `slice-v1-abschluss-antwortvergleich-doku`) für eine handgeschriebene Aufzeichnung mit `connection_closed` nicht mehr; das Handbuch hinkt dem Binary bis zu diesem Slice hinterher, ein Release liegt nicht dazwischen. Die Lücke steht in `slice-v1-abschluss-verbindungsende-einspielen`, `-einspielen-extended` und `-vergleich` §1 als derselbe Vermerk.

**Übernommen aus `slice-v1-abschluss-antwortvergleich-doku`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): die Sätze zu `FATAL` und `PANIC` bei `--compare-responses` (dort „sind `PGR-E5004`“). Sie beschreiben bis `slice-v1-abschluss-verbindungsende-vergleich` den Ist-Zustand; ab dann gilt die Ausnahme eines aufgezeichneten Endes an derselben Stelle mit demselben SQLSTATE, und dieser Slice ersetzt den Satz.

**Übernommen aus `slice-v1-abschluss-sessions-doku`** (dort §1/§7, *Ungenannt geblieben*; Sendung des Planners bei dessen Closure, 2026-10-13, `AGENTS.md` §3.13): das aufgezeichnete Ende des Servers (`FATAL`, `PANIC`) im Zusammenhang mehrerer Verbindungen und Sitzungen. Das Handbuch nennt im Abschnitt *Mehrere Verbindungen aufzeichnen und wiedergeben* dazu keinen Satz (Lücke vom Nutzer akzeptiert, Entscheidung vom 2026-10-12, wie oben); dieser Slice prüft die Sätze jenes Abschnitts und der Sitzungs-Hinweise bei *Aufzeichnen* und *Wiedergeben* gegen das Verbindungsende (Sitzung mit Ende des Servers, Zuordnung der Sitzungen, Exit-Code, Verbindungsfehler) und sagt dazu nur, was die Probe am Binary zeigt. Sie gehört zu DoD-Punkt 1 (Handbuch) und 2 (README, soweit es Sitzungen nennt).

**Übernommen aus den Slices des Verbindungsendes** (Sendung des Planners, 2026-10-12): die Sätze des Handbuchs, die ein Ende des Servers berühren — zu `--allow-recorded-errors` (der Satz, ein Fehler mit dem Schweregrad `FATAL` breche das Einspielen auch dann ab), zur Version der Aufzeichnung in YAML und in SQLite (`version`, `meta`), zu `record` (Verbindung des Servers endet nach einer Fehlerantwort), zu `replay` und zu den Meldungscodes `PGR-E3002`, `PGR-E3003`, `PGR-E4003`, `PGR-E5004` und `PGR-E6001`. Welche Sätze ins Handbuch kommen, bestimmt die Probe am gebauten Binary; was sie nicht zeigt, bleibt draußen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur das Ende des Servers.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` in den Code-Slices.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Das Benutzerhandbuch beschreibt, was die Slices des Verbindungsendes liefern, im Ist-Zustand des gebauten Binaries (Probe am Binary, gegen PostgreSQL und gegen handgeschriebene Aufzeichnungen): wann `record` ein Ende aufzeichnet und wann nicht (Fehlerantwort `FATAL`/`PANIC` in laufender Anfrage; Ende im Leerlauf nur mit nachfolgender Anfrage; Netzwerkabbruch, Ende im Aufbau und `ERROR` mit Abbruch nicht), was der Client sieht, das Format (`version` 1 und 2 in YAML und SQLite, ein Ende macht die Version 2, ältere Programmversionen lehnen sie ab), was `replay` wiedergibt, wie `play` ein aufgezeichnetes Ende behandelt (`--allow-recorded-errors` mit und ohne `--compare-responses`, `PGR-E4003`, `PGR-E5004`, Zähler) und die Zeilen der Meldungscodes, dazu das Ende des Servers im Zusammenhang mehrerer Verbindungen und Sitzungen (aus `slice-v1-abschluss-sessions-doku`, §1); die Sätze, die dieser Stand überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11).
- [ ] `README.md` nennt, was die Slices des Verbindungsendes liefern, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Stand überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu Aufzeichnen, Wiedergeben, Einspielen (Hinweise zu `--allow-recorded-errors`), Format von YAML und SQLite und Meldungscodes |
| `README.md` | update | Funktionsumfang im Ist-Zustand |
| Beispiel-Dateien für die Probe (`test/integration/testdata` oder Probe-Verzeichnis) | neu (nur wenn die Probe sie braucht) | Jedes Beispiel des Handbuchs läuft als Datei gegen das Binary |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-upstream` und `slice-v1-abschluss-verbindungsende-upstream-extended` (Schritt 41b) liegen in `done/` (Voraussetzung: das Binary zeigt das ganze Verhalten; das Handbuch nennt das Ende einer Extended-Interaktion gegen PostgreSQL erst nach dem Beleg in 41b, Bindung aus `-upstream-extended` §1); `slice-v1-abschluss-antwortvergleich-doku` liegt davor. Dieser Slice ist Schritt 42 der Reihenfolge in §5 von welle-v1-abschluss, hinter Schritt 41b. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Probe am Binary zeigt Verhalten, das die Spezifikation anders regelt, oder eines der Code-Slices liefert etwas, das §1 nicht trägt — Befund an den Planner, zurück zur Zerlegung.
- `in-progress` → `open`: Ein Code-Slice des Verbindungsendes wird zurückgeführt und liefert sein Verhalten nicht — Carveout.

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

**Randformen:** keine; der Slice liefert keinen neuen Vertrag (`AGENTS.md` §3.12), er beschreibt das Verhalten der Code-Slices.

**Risiken:**

- Das Handbuch sagt zu, was nur gegen PostgreSQL 17 und handgeschriebene Aufzeichnungen geprüft ist (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×); es nennt keine Version, die die Probe nicht gefahren hat — **Ausgang:** offen bis Closure.
- Bis `slice-v1-abschluss-verbindungsende-vergleich` (Schritt 39a) beschreiben die Sätze von `slice-v1-abschluss-antwortvergleich-doku` zu `FATAL`/`PANIC` den Ist-Zustand; zwischen den Code-Slices und diesem Slice ist der Satz für eine handgeschriebene Aufzeichnung mit Ende nicht mehr wahr; die Lücke ist akzeptiert (Entscheidung des Nutzers vom 2026-10-12, §1; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure (erwartet: eingetreten und mit diesem Slice ersetzt).

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

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×; geliefert bei den bisherigen Doku-Folge-Slices, ohne Auftreten) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) — die Probe prüft die ersetzten Sätze mit der Suche nach den alten Wortlauten. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Handbuch, README). Schichten: eine (Nutzer- und Wartungs-Doku). Schätzung des Planners: 150 bis 300 Zeilen Text.

Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-sessions-doku` (Planner, 2026-10-13, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Handbuch) und 2 (README), kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku).

Nachgezählt beim Eintragen der Kennung `slice-v1-abschluss-verbindungsende-format-schreiber` (Ziel und Berührte Spec-Stellen; Planner, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-format`): Der Teil beschreibt, was das Binary zeigt, und gehört zu DoD-Punkt 1 (Handbuch) und 2 (README); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Ziel; Planner, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-sqlite`): Der Teil beschreibt, was das Binary zeigt, und gehört zu DoD-Punkt 1 (Handbuch) und 2 (README); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku).

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-replay-extended` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-replay`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion im Replay), und gehört zu DoD-Punkt 1 (Handbuch) und 2 (README); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Schritte 35a bis 41 der Voraussetzung (Schritt 42) schließen 37b ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-upstream-extended` (Ziel, Start und Bindung; Planner, 2026-10-16, `AGENTS.md` §3.13, Schritt 41b hinter `-upstream`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion gegen PostgreSQL 17, Probe am Binary), gehört zum selben Satz des Handbuchs wie 37b und 37c und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 41b ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-replay-extended`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion mit `Flush` im Replay), gehört zum selben Satz des Handbuchs wie 37b und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 37c ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-einspielen-extended` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-einspielen`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion beim Einspielen ohne Vergleich), gehört zum selben Satz des Handbuchs wie 38a und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 38b ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-einspielen-grenzen` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt beim Anspruch aus `-einspielen-extended`): Der Teil beschreibt, was das Binary zeigt (ohne `--allow-recorded-errors` und bei aufgezeichnet ohne Ende an der Stelle ist die `FATAL` `PGR-E4003`), gehört zum selben Satz des Handbuchs wie 38a und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 38c ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-hilfetext-binary` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt beim Anspruch aus `-vergleich`, Schritt 39b): Der Teil beschreibt, was das Binary zeigt (Ende im Vergleich am Binary, Hilfetext von `--allow-recorded-errors`), gehört zum selben Satz des Handbuchs wie 39a (`FATAL`/`PANIC` bei `--compare-responses`) und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 39b ein. Die Handbuch-Lücke (§1, §6) beginnt mit 39a, dem Schritt, der das Verhalten ändert; 39b ändert am Verhalten nichts.
