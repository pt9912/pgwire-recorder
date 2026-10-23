# Slice slice-play-feldvergleich: `play --compare-responses` meldet auf Wunsch alle Abweichungen einer Interaktion

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** — (wellenlos). `welle-v1-nachtrag` geprüft und **nicht** gewählt: Ihr Ziel sind die zwei Anforderungsreste von v1 (Homebrew, `CancelRequest` ohne Schlüssel, Meilenstein M5); `LH-FA-24` ist umgesetzt, der Slice liefert keine Anforderung dieser Welle, und ihn aufzunehmen dehnte ihren Closure-Trigger (Out-of-Scope-Disziplin, §6 ihres Plans). Eine Welle für die Arbeit nach v1 gibt es noch nicht; entsteht eine, kann sie ihn einsammeln. Die Roadmap führt wellenlose Arbeit nicht.

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)

**Berührte Spec-Stellen:** `LH-FA-24.a` (*Meldung*: „Je Interaktion wird die erste Abweichung gemeldet“; *Randformen des Vergleichs*; neu *Alle Abweichungen* mit allen Randformen der Option, festgelegt vor dem Start; Tabelle „Fehler und Abweichungen“ mit `PGR-E5004`) · `LH-FA-20.a` (Optionen von `play`; Folgen einer Abweichung, `--continue-on-error`) · `LH-FA-17.a` (Optionstabelle: Umgebungsvariable und Schlüssel der Konfigurationsdatei) · `LH-FA-24` (Beschreibung und Out-of-Scope: „Vergleich der Zeilenwerte“; gelesen, siehe §1 und §6) · `SPEC-027` (Zeile `PGR-E5004` der Fehlerklassen; gelesen, Exit-Code 5 bleibt) · `SPEC-034` (Meldungskatalog; gelesen, kein neuer Code)

**Entscheidung:** Entscheidung des Nutzers vom 2026-10-20 (Memory `play-voller-feldvergleich-per-option`): Der Standard bleibt die erste Abweichung; eine **neue Option** schaltet die vollständige Ausgabe aller abweichenden Felder je Interaktion zu. Name, Lastenheft und Form der Ausgabe: Entscheidung des Nutzers vom 2026-10-23 (Name `--diff-all`, Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL`, Schlüssel `diff_all`, kein Kurzname; Lastenheft unverändert; eine Zeile `PGR-E5004` je Interaktion, beide Seiten je Art, nie Zeilenwerte oder Fehlertext, Obergrenze 20 Positionen). Festgelegt in `spec/spezifikation.md` `LH-FA-24.a` *Alle Abweichungen* (Architect, 2026-10-23, vor dem Start). Keine ADR.

**Verantwortlich:** Implementer (Zweig impl-diff-all)
**Autor:** Planner. **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Das gebaute Binary meldet bei `play --compare-responses` mit der Option `--diff-all` (Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL`, Schlüssel `diff_all`, kein Kurzname) je Interaktion **alle** Abweichungen der verglichenen Felder in einer Zeile `PGR-E5004` statt der ersten; ohne die Option bleibt es bei der ersten Abweichung, Wort für Wort die heutige Meldung, Exit-Code 5 und die Zähler der Zusammenfassung bleiben. Das ist eine **Änderung der Spezifikation**, kein Fehler im Code: Der Architect hat sie vor dem Start in `spec/spezifikation.md` festgelegt (`LH-FA-24.a` *Alle Abweichungen*, Optionstabellen, Hilfe von `play`, Fehlertabelle, Änderungsprotokoll; Entscheidung des Nutzers vom 2026-10-23). Der Slice liefert die Sammlung der Abweichungen im Vergleich des Services und die Option der Kommandozeile mit Hilfetext, Umgebungsvariable und Schlüssel der Konfigurationsdatei.

**Folge-Slice:** `slice-play-feldvergleich-doku` — liefert Handbuch und README im Ist-Zustand für die Option (`AGENTS.md` §3.11; Doku ist eine eigene Schicht, §8). Er steht in `open/`, nimmt die Sendung in seinem §1 an und folgt direkt hinter diesem Slice in derselben Welle bzw. wellenlos; ein Release mit der Option entsteht erst nach beiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Vergleich der **Zeilenwerte** und der **Zeilenzahl** — Bestand bleibt bewusst stehen (Klasse 2). Grund: `LH-FA-24` schließt den Vergleich der Zeilenwerte ausdrücklich aus (Out-of-Scope) und sagt, Werte und Anzahl der Zeilen würden nicht verglichen; „alle abweichenden Felder“ meint hier die Felder, die der Vergleich heute schon liest (Nachrichtenart, Spalten, Parameter-Typen, Befehl, SQLSTATE, Transaktionsstatus). Der Nutzer hat am 2026-10-23 die Form ohne Zeilenwerte und ohne Fehlertext bestätigt; Werte zu vergleichen wäre eine Änderung des Lastenhefts (Major, `harness/conventions.md` §Versionierung des Lastenhefts) und ein eigener Vorgang (Klasse 3).
- Die Standard-Ausgabe (erste Abweichung) — Bestand (Klasse 2). Grund: Entscheidung des Nutzers vom 2026-10-20: Der Standard bleibt.
- Neue Meldungscodes, ein anderer Exit-Code als 5 und neue Attribute der Zeile zum Ende — Bestand (Klasse 2). Grund: Die Abweichung bleibt `PGR-E5004`; `abweichend` zählt weiter Interaktionen, nicht Abweichungen. Einen Zähler der Abweichungen gibt es nicht (entschieden, *Zähler* in `LH-FA-24.a` *Alle Abweichungen*); die Gesamtzahl steht nur im Kopf der Meldung.
- Toleranzregeln für einzelne Felder, Ausnahmelisten, das Verändern der Aufzeichnung anhand des Vergleichs — anderer Vorgang (Klasse 3). Grund: Out-of-Scope von `LH-FA-24`.
- Handbuch und `README.md` — Folge-Slice `slice-play-feldvergleich-doku` (Klasse 1, mit Kennung). Grund: Doku ist die dritte Schicht (§8); ohne sie nähme der Slice drei Schichten auf (`AGENTS.md` §3.13).
- `--fail-on-unconsumed` — Bestand (Klasse 2). Grund: Die Option gehört zu `replay`; bei `play` ist sie eine unbekannte Option (`PGR-E2001`), es gibt kein Zusammenspiel zu regeln (§6).
- Produkt-Code außerhalb von `internal/hexagon/services` und dem CLI-Adapter, insbesondere die Driven-Adapter (Aufzeichnung, PGWire-Client) und das Aufzeichnungsformat — Schicht-Abgrenzung (Klasse 4). Grund: Der Vergleich liest die Folge, die er schon liest; kein Port, keine Datei und kein Wire-Verhalten ändert sich.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 — der Vergleich sammelt alle Abweichungen** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), `LH-FA-24.a`): Mit der Option meldet `play` je Interaktion alle Abweichungen der Folge in der Form, die die Spezifikation festlegt (Zuordnung der Positionen, Inhalt der Meldung je Art, Obergrenze, Ende des Servers, Signal, Zusammenspiel mit `--continue-on-error`); ohne die Option bleibt es bei der ersten Abweichung, Wort für Wort die heutige Meldung. Je Zusage der Spezifikationsstelle ein Test, bei einer Aufzählung ein Fall je Element, und je Zusage eine **Mutation** (`AGENTS.md` §3.10: der Vergleich liegt im Play-Pfad), vom Implementer selbst rot gesehen; der Bericht nennt je Zusage: Zusage · Mutation · roter Test. Der Standardfall ist gegen die Option gebrochen und rot, wenn der Code die Option ignoriert oder immer sammelt.
- [ ] **Liefer-Punkt 2 — die Option** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)): `play` nimmt `--diff-all` an, mit Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL` und Schlüssel der Konfigurationsdatei nach dem Muster von `--compare-responses` (Optionstabelle der Spezifikation); ihre Verwendung ohne `--compare-responses`, bei `record` und `replay`, ihr Wert und der Hilfetext von `play` verhalten sich wie die Spezifikation es festlegt, die Meldung bleibt `PGR-E2001` bzw. `PGR-E2004` ohne neuen Code. Ein Test je Zusage und je Element der Aufzählungen (Kommandozeile, Umgebung, Datei, Wertemenge, Kombination, Kommandos), mit Abdeckungs-Deklaration; `make abdeckung-check` und `make meldungskatalog-check` sind grün. Die Option liegt im Record-/Replay-/Play-Pfad, soweit sie dessen Verhalten ändert: Mutation je Zusage, die das Verhalten bestimmt (Wirkung, Kombination); für den Hilfetext und das Annehmen des Namens genügt ein Test je Zusage (`AGENTS.md` §3.10).
- [ ] `make gates` am Endstand: grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services/play_vergleich.go` | update | neben `erste` eine Funktion, die alle Abweichungen der normalisierten Folge liefert (`unterschied` je Position mit den Einzelheiten beider Seiten, Zuordnung, Folgen verschiedener Länge einzeln und Ende des Servers nach `LH-FA-24.a` *Alle Abweichungen*); Text der Liste mit Gesamtzahl, Obergrenze 20 Positionen und 20 Spalten je Position; `erste` und seine Meldung bleiben |
| `internal/hexagon/services/play.go` | update | Feld der `PlayOptionen` für die Option; der Aufrufstelle des Vergleichs wählt `erste` oder die Liste; Meldung `PGR-E5004` je Interaktion; Zähler `abweichend` unverändert |
| `internal/hexagon/services/play_vergleich*_test.go` | update/neu | Fälle je Zusage und je Element, Mutationen aus DoD 1; Standardfall unverändert |
| `internal/adapters/driving/cli/cli.go` | update | Option in der Tabelle der Optionen von `play` (Wahrheitswert, Standard `false`, Umgebung, Schlüssel), Kombination mit `--compare-responses` in der Reihenfolge der Prüfungen, Hilfetext von `play` |
| `internal/adapters/driving/cli/play_test.go` und Nachbarn | update/neu | Fälle je Zusage und Element aus DoD 2 |
| `test/integration/play_vergleich_e2e_test.go` | update | E2E gegen das gebaute Binary und PostgreSQL: eine Interaktion mit mehreren Abweichungen, mit und ohne Option; mit Abdeckungs-Deklaration |
| `docs/user/abdeckung-*.md` | update (erzeugt) | `make abdeckung` nach den neuen Deklarationen |

**Nicht in dieser Liste, weil vor dem Start:** `spec/spezifikation.md` (`LH-FA-24.a` *Alle Abweichungen*, Optionstabellen, Hilfe von `play`, Fehlertabelle, Änderungsprotokoll); die Festlegung aus §6 steht dort, im Commit des Architect **vor** dem ersten Code-Commit (`AGENTS.md` §3.12); das Lastenheft ändert sich nicht (Entscheidung des Nutzers vom 2026-10-23). Dieser Commit gehört nicht zum Slice. Ändert der Architect die Spezifikation **während** des Slice (Rückgabe durch den Implementer), zählt sie ab dann als dritte Schicht — dann greift der vorgemerkte Schnitt aus §4.

Nicht in dieser Liste: `internal/bootstrap` (die Option wird als Feld durchgereicht, reine Verdrahtung, zählt nicht, `AGENTS.md` §3.13; liest der Bootstrap Datei oder Umgebung selbst, ist das Logik und der Plan prüft die Grenze neu), Handbuch und `README.md` (Folge-Slice), Driven-Adapter, Gates, Verträge unter `harness/vertraege/`, `internal/hexagon/model` (kein neuer Meldungscode, keine Logik).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): (1) Die Fragen aus §6 sind entschieden (Entscheidung des Nutzers vom 2026-10-23). (2) Der Architect hat die Randformen aus §6 in `spec/spezifikation.md` festgelegt und die Liste geprüft (`AGENTS.md` §3.12; erledigt am 2026-10-23, §6). (3) Der Doku-Folge-Slice `slice-play-feldvergleich-doku` liegt in `open/`. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik und Messbefehl: `AGENTS.md` §3.14; Pfade des Plans: `internal/hexagon/services internal/adapters/driving/cli test`); die Basis ist der Commit des Übergangs nach `in-progress/`, sein Hash: `bf68b876` (Zweig `impl-diff-all`). Dann hält der Implementer an und meldet dem Planner.
- `in-progress` → `next`: **vorgemerkter Schnitt.** Ändert der Architect die Spezifikation im Slice (§3) oder verlangt die Festlegung Logik im Bootstrap, zählt die Zählung neu; liegt der Slice dann über zwei Schichten, geht die Option der Kommandozeile (Liefer-Punkt 2) in einen eigenen Slice direkt dahinter, dieser Slice behält den Vergleich.
- `in-progress` → `open`: Der Architect gibt eine Entscheidung zurück, die der Nutzer treffen muss und §6 nicht nennt — der Slice wartet nicht in `in-progress/`, er geht zurück und hält das WIP-Limit frei.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben; der Bericht des Implementers nennt je Zusage Mutation und roten Test.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Entscheidungen des Nutzers vom 2026-10-23** (zu den drei Empfehlungen des Planners „ja“, so vom aufrufenden Lauf übermittelt):

1. *Name der Option:* `--diff-all`, Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL`, Schlüssel `diff_all`, kein Kurzname.
2. *Lastenheft:* wird nicht geändert, kein Versions-Bump; Optionsname und Form der Ausgabe gehören ins Technik-Stratum. Der Architect hat den Wortlaut geprüft: `LH-FA-24` legt „die erste“ Abweichung nicht fest, die Meldung nennt weiter Interaktion und Art, die Einzelheiten (Namen, Typ-OID, SQLSTATE, Befehl ohne Zahlen, Transaktionsstatus) sind Felder, die der Vergleich schon liest, keine Zeilenwerte und keine Zeilenzahl; kein Widerspruch zum Out-of-Scope.
3. *Form der Ausgabe:* eine Zeile `PGR-E5004` je Interaktion, Abweichungen in aufsteigender Position, je Art beide Seiten, nie Zeilenwerte oder Fehlertext, höchstens 20 Positionen, danach „und N weitere“ mit Gesamtzahl; der Standard bleibt die erste Abweichung.

**Randformen — entschieden vom Architect, 2026-10-23, in `spec/spezifikation.md` `LH-FA-24.a` *Alle Abweichungen*** (Stelle; Punkte der Liste des Planners; Geschwister-Optionen `--compare-responses` und `--allow-recorded-errors`; Abweichung vom Vorschlag genannt):

1. *Option ohne `--compare-responses`:* `PGR-E2001`, Exit-Code 2, Kombination nach der Zusammenführung, vor `--upstream` und vor jeder anderen Kombination von `play`; maßgebend sind die wirksamen Werte (`--diff-all` `true`, `--compare-responses` `false`), auch bei verschiedenen Quellen und bei der Datei als Quelle (kein `PGR-E2004`); `--diff-all=false` ist auch ohne `--compare-responses` gültig (Spiegel von `--allow-plaintext=false`). Stelle: *Ohne `--compare-responses`*. Zusatz zum Vorschlag: die Fälle `false` und die gemischten Quellen sind entschieden.
2. *Quellen und Wertemenge:* wie `--compare-responses`; `record`/`replay` kennen die Option nicht (Kommandozeile und Umgebung `PGR-E2001`, Schlüssel `PGR-E2004`). Stelle: *Quellen und Werte*, *Andere Kommandos*; Optionstabelle in `LH-FA-17.a` (Zeile hinter `--compare-responses`, damit steht die Variable an dieser Stelle der Prüfung der Umgebung).
3. *`--continue-on-error`:* wie Vorschlag. Stelle: *Zusammenspiel*.
4. *`--allow-recorded-errors`:* unverändert ohne Wirkung. Stelle: *Zusammenspiel*.
5. *`--fail-on-unconsumed`:* entfällt, bei `play` unbekannt. Stelle: *Zusammenspiel*.
6. *Verglichene und ausgegebene Felder:* dieselben wie heute; je Art die beiden Seiten; Reihenfolge innerhalb einer Position: erst die Zahl, dann je Nummer aufsteigend, Name vor Typ-OID. Stelle: *Einzelheiten je Art, beide Seiten*. **Abweichung vom Vorschlag, neue Randform:** der Befehl erscheint (ohne Zahlen) und bei Spalten und Parametern gilt eine Obergrenze von 20 je Position mit „und N weitere“; bei verschiedener Zahl der Spalten oder Parameter steht die Zahl beider Seiten, Spalten ohne Gegenstück stehen nicht einzeln.
7. *Zuordnung der Positionen:* Position für Position, ohne Abgleich der Folgen; Folgepositionen einer Verschiebung werden gemeldet. **Offene Frage entschieden:** Folgen verschiedener Länge **jede Position einzeln** (keine Sammelmeldung), weil ein Format, eine Zählung und die Obergrenze die Länge der Meldung begrenzen. Stelle: *Zuordnung der Positionen*.
8. *Obergrenze:* 20 Positionen, danach „und N weitere“, Gesamtzahl im Kopf der Meldung; Fälle 20 genau, 21 und 0 (0: keine Zeile). Stelle: *Obergrenze*, *Form der Meldung*.
9. *Ende des Servers:* wie Vorschlag; Abweichungen der aufgezeichneten Folge hinter der Position des Endes erscheinen nicht; ein erwartetes Ende ist keine Position, Abweichungen davor bleiben. Stelle: *Ende des Servers*.
10. *Stelle in der Folge:* je Nachrichtenart die Stelle des Bestands, nur dass der Vergleich nach der ersten Abweichung weiterläuft; verschiedene Typen an einer Position (`ready_for_query` gegen `connection_closed`) sind eine Position der Art Nachrichtenart. Stelle: *Stelle in der Folge*.
11. *Signal:* wie Vorschlag. Stelle: *Signal*.
12. *Zeilen `error` und Zähler:* unverändert; `abweichend` höchstens eine je Interaktion. Stelle: *Zähler*.
13. *Hilfetext:* `play --help` trägt drei Aussagen (alle statt der ersten, höchstens 20 Positionen, nur mit `--compare-responses`); die globale Hilfe nennt die Option nicht; ein Test je Aussage. Stelle: Hilfe von `play` in `LH-FA-20.a`.
14. *Gleichzeitigkeit und Aufrufe:* keine Ordnung über mehrere Aufrufe; `make`-Aufrufarten entfallen (kein Hook, kein Gate liest das Ergebnis). Stelle: *Serverversion, Aufruf, Plattform*.
15. *Plattform:* **Abweichung vom Vorschlag, akzeptiertes Negativ:** Meldung und Exit-Code hängen weder von Betriebssystem noch Dateisystem ab (reine Zeichenfolge auf `stderr`), es gibt keine Eigenschaft, für die eine Plattformgrenze zuzusagen wäre; die Spezifikation sagt das an der Zusage. Stelle: *Serverversion, Aufruf, Plattform*.
16. *Serverversion:* der Vergleich gleicht Versionen nicht aus, die Option ändert das nicht. Stelle: *Serverversion, Aufruf, Plattform*.
17. *Form der Meldung (neu, aus der Entscheidung 3):* genau eine Zeile je Interaktion mit Abweichung, mit Session, Sequenznummer und Gesamtzahl, auch bei nur einer Abweichung die Listenform mit Einzelheiten; ohne Abweichung keine Zeile. Der Wortlaut der Zeile ist nicht festgelegt, der Inhalt je Art ist es; die Tests prüfen den Inhalt. Stelle: *Form der Meldung*.


**Risiken:**

- Der Nutzer meint mit „Feldvergleich“ auch die Werte der Zeilen; die Option lieferte dann weniger als erwartet, und das Lastenheft schließt Werte aus — **Ausgang:** offen bis Closure (der Nutzer hat am 2026-10-23 die Form ohne Zeilenwerte bestätigt, Entscheidung 3).
- Folgeabweichungen nach einer verschobenen Folge machen die Liste lang und irreführend (eine fehlende Nachricht erzeugt viele Positionen) — **Ausgang:** offen bis Closure (Obergrenze 20, Folgen verschiedener Länge einzeln und der Satz zur Verschiebung stehen in der Spezifikation; ob die Praxis einen Abgleich der Folgen verlangt, wäre ein Folge-Slice).
- Der Slice erreicht die Marke 1500 (Vergleich, Obergrenze und Mutationen je Zusage und je Element der Aufzählungen sind viele Tests) — **Ausgang:** offen bis Closure (Rückführung §4, vorgemerkter Schnitt Option/Vergleich).
- Ein Handbuch oder README, das die Option vor dem Folge-Slice nennt, oder ein Release dazwischen — **Ausgang:** offen bis Closure (`slice-play-feldvergleich-doku` direkt dahinter; Release erst nach beiden, `AGENTS.md` §3.11).

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

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (gemergter Stand, 2026-10-23). Treffer, die diesen Slice betreffen: `negativtests-fehlen-bei-neuem-vertrag` (verkörpert in `AGENTS.md` §3.10; DoD 1 und 2 verlangen je Zusage eine Mutation), `spec-randform-erst-im-review-entschieden` (verkörpert in §3.12; §6 nennt die Randformen, der Architect legt sie vor dem Code fest), `randform-wellenlos-ohne-architect-vor-code` (der Slice ist wellenlos; §4 macht die Festlegung des Architect zur Start-Bedingung), `zusage-im-kommentar-weiter-als-pruefung` (verkörpert in §3.11; Hilfetext und Handbuch sagen nur zu, was ein Test prüft), `plan-folgt-korrektur-nicht` (verkörpert in §3.9; der Plan folgt der Festlegung des Architect) und `gruener-mutant-faelschlich-aequivalent` (§3.10; ein äquivalenter Mutant braucht Invariante, Ort und Sonde). Keiner dieser Einträge erreicht mit diesem Slice eine Schwelle neu (die verkörperten zählen erst wieder, wenn die Regel trotz ihrer Verkörperung versagt).

**Zählung** (Planner, 2026-10-23, `AGENTS.md` §3.13; Marke `AGENTS.md` §3.14): **Liefer-Punkte: 2** (1: der Vergleich sammelt alle Abweichungen; 2: die Option mit Hilfetext, Umgebungsvariable und Schlüssel). **Schichten: zwei** — der Kern des Hexagons mit seinen Tests (`internal/hexagon/services/`: Vergleich und Meldung) und der CLI-Adapter mit seinen Tests (`internal/adapters/driving/cli/`; der E2E-Fall unter `test/integration/` prüft dieselbe Fläche). Die Spezifikation ist eine eigene Schicht, liegt aber vor dem Start und damit außerhalb des Slice (§3, §4). Der Bootstrap reicht ein Feld durch (Verdrahtung, zählt nicht, §3). **Nutzer-Doku (`docs/user/`, `README.md`) wäre die dritte Schicht**: Sie geht in `slice-play-feldvergleich-doku` (`AGENTS.md` §3.11, §3.13). Keine Sendung in diesen Slice eingetragen (kein Geber nennt ihn als Nehmer), die Zählung ist die des Plans selbst. **Nachgezählt vom Architect (2026-10-23, nach der Festlegung der Randformen):** Liefer-Punkte 2, Schichten 2 (Kern und CLI-Adapter), unverändert; die Spezifikation liegt vor dem Start, der Bootstrap bleibt Verdrahtung, solange die Prüfung der Kombination `--diff-all` ohne `--compare-responses` im CLI-Adapter liegt (sie ist Logik des Adapters, keine des Bootstraps). Gewachsen sind die Fälle je Zusage (17 Randformen, Aufzählung „Einzelheiten je Art“ mit fünf Arten und zwei Obergrenzen); die Größe misst die Marke 1500 (§4), nicht die Zahl der Punkte.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und Entscheidung gehen dem Code voraus; der Architect legt die Festlegung vor dem ersten Code-Commit fest).
