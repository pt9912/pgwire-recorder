# Slice slice-v1-abschluss-verbindungsende-hilfetext-binary: Verbindungsende im Antwortvergleich (Hilfetext und Abnahmeszenario 16, Satz 3)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md), [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-24.a` · `LH-FA-20.a` · `LH-FA-17.a` · `SPEC-034` · `SPEC-029` · `ARC-005`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Ende des Servers im Antwortvergleich ist am Binary belegt und beschrieben: Der Hilfetext von `--allow-recorded-errors` nennt das aufgezeichnete Ende ohne Vergleich und die Wirkungslosigkeit der Option mit Vergleich, soweit ein Test im CLI-Adapter es prüft; Abnahmeszenario 16, Satz 3, ist gegen PostgreSQL automatisiert (eine Aufzeichnung mit Ende, mit Vergleich eingespielt gegen einen Server, der das Ende an derselben Stelle mit demselben Fehlercode herbeiführt: Exit-Code 0 und die nächste Session läuft; bei anderem Fehlercode `PGR-E5004` und Exit-Code 5; die Zeile zum Ende des Laufs trägt `eingespielt`, `verglichen` und `abweichend`). Der Slice ändert keinen Code des Service: Das Verhalten liefert `slice-v1-abschluss-verbindungsende-vergleich`.

**Geschnitten aus `slice-v1-abschluss-verbindungsende-vergleich`** (Planner, 2026-10-16, beim Anspruch von Schritt 39, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Die Prüfung des Architect schätzte den ungeteilten Plan auf 445 bis 970, Mitte rund 710 (71 % der Marke 1000), und empfahl den Schnitt *Service | Hilfetext und E2E*; dies ist die zweite Hälfte (Schritt 39b). Sie setzt Teil 1 (Schritt 39a, `done/`) voraus: Der Hilfetext beschreibt, was das Binary seit Teil 1 zeigt, und die Fälle am Binary prüfen es. Teil 1 wartet nicht auf diesen Slice.

**Übernommen aus `slice-v1-abschluss-verbindungsende-vergleich`** (dort §1, Weitergabe; Sendung des Planners, 2026-10-16, `AGENTS.md` §3.13): der Hilfetext von `--allow-recorded-errors` und `--compare-responses` im Binary, der Test im CLI-Adapter und die Fälle von Abnahmeszenario 16, Satz 3, am Binary. Ursprung: `slice-v1-abschluss-verbindungsende-einspielen` (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12), das den Hilfetext und das Szenario hierher gab.

**Übernommen aus `slice-v1-abschluss-postgres-versionen-zusage`** (dort §6, Randform „Szenarien nur gegen 17“; Sendung des Architect, 2026-10-15, `AGENTS.md` §3.13, über die Adresse `slice-v1-abschluss-verbindungsende-vergleich` und deren Weitergabe): `SPEC-029` Punkt 2 nennt das erwartete Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) heute nicht, weil kein Test es prüft. Im Commit, der den Integrationstest zu Szenario 16 Satz 3 liefert, ergänzt dieser Slice Punkt 2 um den Satz „Das erwartete Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) ist nur dort geprüft“ (dort = Referenz 17); vorher steht er nicht (`AGENTS.md` §3.11). Dazu berichtigt der erste der beiden Nehmer, der Punkt 2 berührt, den Numerus des vorhandenen Satzes „Die Anmeldung mit Klartext-Passwort und mit MD5 sind nur dort geprüft“ (Verifikation `slice-v1-abschluss-postgres-versionen-zusage`, I-1, INFO): „Die Anmeldung mit Klartext-Passwort und die mit MD5 sind nur dort geprüft“; der zweite prüft, dass es geschehen ist. Ein Wort im selben Satz, kein Liefer-Punkt. Der Nehmer `slice-v1-abschluss-verbindungsende-upstream` (Schritt 41) schreibt ebenfalls einen Satz in Punkt 2; wer von beiden Punkt 2 zuerst berührt, berichtigt, der andere prüft es.

**Akzeptiertes Negativ** (Entscheidung des Nutzers vom 2026-10-12, aus `slice-v1-abschluss-verbindungsende-vergleich` §6): Szenario 16 Satz 3 läuft nur gegen PostgreSQL 17; die Matrix aus `slice-v1-abschluss-postgres-versionen-matrix` fährt es nicht mit (dort §1 als Abgrenzung).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Ende im Vergleich im Service, die Zähler, die Fälle der letzten `Flush`-Gruppe, die drei Kommentare in `play.go` und `play_vergleich.go` — `slice-v1-abschluss-verbindungsende-vergleich` (Geber, `done/`; Klasse 4): geliefert; ein Befund der Fälle am Binary ist die Rückführung in §4, nicht eine Änderung des Service hier.
- Der Einleitungssatz von `LH-FA-24.a` — Wortlaut des Architect vom 2026-10-16, geliefert vor dem Code (Klasse 2).
- Die Einstufung ohne Vergleich und die Fälle dazu am Binary — `slice-v1-abschluss-verbindungsende-einspielen`, `-einspielen-extended`, `-einspielen-grenzen` (alle `done/`; Klasse 4).
- Der Satz in `SPEC-029` Punkt 2 zu Abnahmeszenario 18 — `slice-v1-abschluss-verbindungsende-upstream` (Klasse 1): dort liegt der Test; vorher steht der Satz nicht.
- Aufnahme von Szenario 16 Satz 3 in die Versionsmatrix (`make test-postgres-versionen`) — akzeptiertes Negativ oben (Klasse 2).
- Aufzeichnen und Wiedergeben des Endes — `-record`, `-upstream`, `-replay` (Klasse 1).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13; Klasse 1), dort auch die Sätze zu `FATAL`/`PANIC` und `PGR-E5004` aus `slice-v1-abschluss-antwortvergleich-doku`. Die Handbuch-Lücke bis dahin ist akzeptiert, Entscheidung des Nutzers vom 2026-10-12 (Klasse 1).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen): Der Hilfetext von `--allow-recorded-errors` nennt das aufgezeichnete Ende ohne Vergleich und die Wirkungslosigkeit mit Vergleich, der von `--compare-responses` das erwartete Ende, soweit ein Test es prüft; jeder Satz des Hilfetextes zu diesen Punkten ist durch einen Test im CLI-Adapter gedeckt, eine Mutation (Satz entfernt) färbt ihn rot.
- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Abnahmeszenario 16, Satz 3, gegen PostgreSQL am Binary, dazu der Satz in `SPEC-029` Punkt 2, dass es nur gegen die Referenz 17 geprüft ist (Sendung aus `slice-v1-abschluss-postgres-versionen-zusage`, im Commit des Tests), samt Berichtigung I-1: Eine Aufzeichnung mit Ende wird mit Vergleich gegen den Server eingespielt, der das Ende an derselben Stelle mit demselben Fehlercode herbeiführt (`pg_terminate_backend`, `57P01`): Der Lauf geht mit der nächsten Session weiter und endet mit Exit-Code 0; bei anderem Fehlercode endet er mit `PGR-E5004` und Exit-Code 5; die Zeile zum Ende des Laufs trägt `eingespielt`, `verglichen` und `abweichend` (Integrationstest, eine Aufzeichnung als eigene Datei je Fall); `TestE2EPlayVergleichFatal` ist nachgezogen (Beschreibung und, falls nötig, Erwartung).
- [ ] Je Zusage des Hilfetextes und der Fälle am Binary eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
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
| `internal/adapters/driving/cli/cli.go` | update | Hilfetext von `--allow-recorded-errors` (aufgezeichnetes Ende ohne Vergleich, ohne Wirkung mit Vergleich) und `--compare-responses` (erwartetes Ende), soweit ein Test es prüft; kein Satz ohne Test (`AGENTS.md` §3.11) |
| `internal/adapters/driving/cli` (Test) | update | Test des Hilfetextes: je Satz zu den beiden Optionen ein Fall |
| `test/integration` (`play_vergleich_e2e_test.go` oder eigene Datei, `testdata`, eine Datei je Fall) | update | Szenario 16 Satz 3 gegen PostgreSQL am Binary: gleicher SQLSTATE (Session mit Ende, nächste Session läuft, Exit-Code 0), anderer SQLSTATE (`PGR-E5004`, Exit-Code 5), Zeile zum Ende des Laufs; `TestE2EPlayVergleichFatal` (aufgezeichnet `ERROR`, Server `FATAL`) wird nachgezogen; Namen der Fixtures: Architect vor dem Code |
| `spec/spezifikation.md` (`SPEC-029` Punkt 2) | update im Commit des Tests | Satz zum erwarteten Verbindungsende beim Einspielen, nur gegen die Referenz 17 geprüft; Numerus des vorhandenen Satzes zur Anmeldung (I-1); die Spezifikation zählt nicht als Schicht |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-vergleich` (Schritt 39a) liegt in `done/` (Voraussetzung: das Binary zeigt das Ende im Vergleich); dieser Slice ist Schritt 39b der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** `BASISSHA`, Stand des Hauptzweigs vor dem Anspruch, vom Planner beim Anspruch von 39b eingetragen, danach nicht mehr geändert; die Commits danach tragen die Basis, den Move und die Roadmap ein und ändern keinen Pfad der Metrik); **über 400** geht der Slice zurück zur Zerlegung. **Meldeschwelle 300:** Erreicht die Summe 300, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Architect in dieser Metrik (Prüfung vor dem Code, 2026-10-16): **130 bis 245, Mitte 190 und damit 48 % der Marke** (Hilfetext und CLI-Test 15 bis 35, zwei Fälle am Binary mit dem Nachziehen von `TestE2EPlayVergleichFatal` 105 bis 180, Tabellen 10 bis 30; Maße der Geber: E2E-Fall mit Testdaten 65 bis 70,5).
  **Neurechnung des Planners beim Anspruch** (2026-10-16, Maße von `slice-v1-abschluss-verbindungsende-vergleich` §7: Messung 566 bei Schätzung 305 bis 640, Nacharbeit zum Review +81, bei `-einspielen*` +89; E2E-Fall 46 bis 70,5 je Fall): Hilfetext und CLI-Test 15 bis 35; zwei Fälle am Binary mit Testdaten je 45 bis 70, zusammen 90 bis 140; Nachziehen von `TestE2EPlayVergleichFatal` 0 bis 30; Tabellen 10 bis 30; Nacharbeit zum Review 0 bis 90. Summe **125 bis 325, Mitte rund 225 (56 % der Marke 400)**; das obere Ende liegt über der Meldeschwelle 300 und bei 81 % der Marke. **Urteil: kein Schnitt.** Die Mitte liegt unter 70 %, der Slice hat zwei Liefer-Punkte und eine Schicht (§8), und der obere Wert enthält die volle Nacharbeit von 90, die nur eintritt, wenn das Review so viel findet wie bei den zwei letzten Gebern; ein weiterer Schnitt trennte den Hilfetext vom Beleg am Binary, der ihn trägt (Schnitt nach Lieferwert). Die Meldeschwelle 300 fängt das obere Ende ab, bevor die Marke 400 erreicht ist.
- `in-progress` → `next`: Ein Fall am Binary verlangt eine Änderung an Produktcode des Service, des Ports oder des Einspiel-Adapters — das Verhalten gehört `slice-v1-abschluss-verbindungsende-vergleich`; zurück zur Zerlegung (Architect und Planner), nicht still im Code.
- `in-progress` → `open`: Der Nutzer ändert die Regel „gleiche Stelle, gleicher SQLSTATE“ — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code. Ort: `LH-FA-24.a` *Ende der Aufzeichnung*, *Meldung*, `LH-FA-20.a` (Optionstabelle, Tabelle der Fehlerregeln). Die Randformen des Vergleichs stehen in `slice-v1-abschluss-verbindungsende-vergleich` §6 und gelten hier unverändert; der Slice fügt dem Service keinen Vertrag hinzu, er belegt die Zusagen des Gebers am Binary und schreibt den Hilfetext. Der Architect prüft vor dem Code, ob jede Zeile entschieden ist, und nennt den Wortlaut des Hilfetextes und die Namen der Fixtures.

- *Gegenstand und Fixtures:* Das Ende bewirkt der Server (`pg_terminate_backend(pg_backend_pid())`, `57P01`), `play` löst es nicht aus. Jeder E2E-Fall baut seine Aufzeichnung selbst (eine Datei je Fall, kein gemeinsamer Zustand), damit eine Mutation nur den Fall färbt, den sie bricht; die Aufzeichnung trägt eine zweite Session, damit „die nächste läuft“ geprüft ist; die Hilfen des Gebers `-einspielen` in `play_ende_e2e_test.go` werden benutzt. Die Zeilen je Fall mit Testdaten und die Zähler (`eingespielt`, `verglichen`, `abweichend`) meldet der Implementer in §7 (Messung).
- *Stelle der Nachricht (`AGENTS.md` §3.12):* wie im Geber §6 (*Stelle von `connection_closed` im Vergleich*): letzte Nachricht der Aufzeichnung der Interaktion, unmittelbar nach `FATAL` oder `PANIC`; jede andere Stelle lehnt das Format beim Laden ab (`PGR-E3003`) und ist kein Fall.
- *Gleicher SQLSTATE:* Das Ende ist erwartet, keine Meldung, Exit-Code 0 — `LH-FA-24.a` *Ende der Aufzeichnung* (erster Punkt). *Anderer SQLSTATE:* `PGR-E5004` Art Fehler, Exit-Code 5, ohne `--continue-on-error` Abbruch — ebenda.
- *Hilfetext:* Zwei Aussagen, jede mit einem Fall im CLI-Test: ohne Vergleich erlaubt die Option das aufgezeichnete Ende; mit Vergleich bleibt sie ohne Wirkung (Optionstabelle `LH-FA-20.a`). Die Aussage zu `--compare-responses` (das Ende ist erwartet, wenn Stelle und SQLSTATE gleich sind) hat einen Fall. Ein Satz, den kein Test trägt, steht nicht im Hilfetext. Aufruf: `--help`/`play --help` wie im Bestand; Zeilenumbruch und Einrückung folgen dem Bestand des Hilfetextes.
- *Satz in `SPEC-029` Punkt 2:* steht erst im Commit des Tests (`AGENTS.md` §3.11); der Wortlaut ist der in §1; der Numerus der Anmeldung (I-1) steht im selben Satz.
- *Mutationen:* Zusage 1 (gleicher SQLSTATE, Ende erwartet): Mutant `play_vergleich.go` behandelt das Ende nie als erwartet → roter Fall „gleich“. Zusage 2 (anderer SQLSTATE, `PGR-E5004`): Mutant behandelt jedes Ende als erwartet → roter Fall „anders“. Zusage 3 (Hilfetext): Satz entfernt → roter Fall des CLI-Tests. Jeder Mutant färbt nur den Fall seiner Zusage; der Implementer prüft das. Die E2E-Mutanten liegen im Arbeitsbaum, werden vor dem Fahren von `make test-integration` committed (das Image baut aus dem Arbeitsbaum, der Stand der Gates hängt am Hash) und danach per `git checkout` zurückgenommen; sie sind Wegwerf-Änderungen.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): vor dem Anspruch geschnitten; die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Szenario 16 Satz 3 stützt sich auf `pg_terminate_backend` und `57P01`; andere Serverversionen sind nicht geprüft (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 4× bei Beginn). **Akzeptiertes Negativ** (Entscheidung des Nutzers vom 2026-10-12): nur gegen PostgreSQL 17, die Matrix fährt es nicht mit — **Ausgang:** weiter offen → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (Stand dort `verkörpert` als dokumentierte Grenze in `LH-FA-20.a` *Gruppen* und *Interaktion*, Entscheidung des Nutzers vom 2026-10-16; der Eintrag bleibt zählbar: ein Nehmer, der das Risiko als *weiter offen* meldet, zählt als Auftreten, der Zähler steht bei 4 und geht mit der Closure dieses Slice auf 5; der Satz in `SPEC-029` Punkt 2 ist die Grenze für Szenario 16 Satz 3).
- Handgeschriebene Aufzeichnungen bilden nur ab, was die Spezifikation sagt (`BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`); die Fälle belegen das Einspielen gegen PostgreSQL, nicht das Aufzeichnen (`slice-v1-abschluss-verbindungsende-upstream-extended`) — **Ausgang:** offen bis Closure (erwartet: weiter offen → dieses Register).
- Ein Fall am Binary färbt sich rot, weil der Service etwas anderes tut als Teil 1 belegt (zweite Rückführung in §4) — **Ausgang:** offen bis Closure.
- Der Satz in `SPEC-029` Punkt 2 sagt mehr zu als der Test prüft (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`): Er nennt nur Szenario 16 Satz 3 und die Referenz 17 und steht erst im Commit des Tests — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (4×; §6, weiter offen), `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (§6), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×) und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (§6, `SPEC-029`) — der Plan nennt Schätzung vorab (§4) und führt die Zusage im Satz nur dort, wo ein Test sie trägt. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schnitt von `slice-v1-abschluss-verbindungsende-vergleich`): **Liefer-Punkte: 2** (Hilfetext mit Test im CLI-Adapter; Abnahmeszenario 16 Satz 3 am Binary mit dem Satz in `SPEC-029`). **Schichten: eine** (CLI-Adapter, `ARC-005`); `test/integration` zählt als Test, die Spezifikation und die Planung zählen nicht; `docs/user` (die Abdeckungstabellen, `make abdeckung`) konservativ als zweite markiert, die Nutzerentscheidung dazu steht aus (`slice-harness-regeln-nach-v1-abschluss`). Schätzung des Architect in der Metrik von §4: 130 bis 245, Mitte 190 (48 % der Marke 400), das obere Ende unter der Meldeschwelle 300. Die Sendungen dieses Slice (aus `slice-v1-abschluss-verbindungsende-vergleich`, `slice-v1-abschluss-verbindungsende-einspielen` und `slice-v1-abschluss-postgres-versionen-zusage`) stehen in §1 und DoD Punkt 1 und 2; der Satz und die Berichtigung in `SPEC-029` sind ein Satz im Commit des Tests, weder Liefer-Punkt noch Schicht.

Nachgezählt beim Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): unverändert **Liefer-Punkte: 2**, **Schichten: eine** (CLI-Adapter; `test/integration` ist Test; Spezifikation und Planung zählen nicht); `docs/user` (Abdeckungstabellen) bleibt konservativ als zweite markiert, also höchstens zwei. Schätzung des Planners 125 bis 325, Mitte rund 225 (56 % der Marke 400), Urteil in §4: kein Schnitt. Register beim Anspruch gesichtet: `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` steht bei 4× und `verkörpert`; der Slice meldet das Risiko *weiter offen* (§6), das fünfte Auftreten wird mit der Closure gezählt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
