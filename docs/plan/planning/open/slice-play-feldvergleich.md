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

**Berührte Spec-Stellen:** `LH-FA-24.a` (*Meldung*: „Je Interaktion wird die erste Abweichung gemeldet“; *Randformen des Vergleichs*; Tabelle „Fehler und Abweichungen“ mit `PGR-E5004`) · `LH-FA-20.a` (Optionen von `play`; Folgen einer Abweichung, `--continue-on-error`) · `LH-FA-17.a` (Optionstabelle: Umgebungsvariable und Schlüssel der Konfigurationsdatei) · `LH-FA-24` (Beschreibung und Out-of-Scope: „Vergleich der Zeilenwerte“; gelesen, siehe §1 und §6) · `SPEC-027` (Zeile `PGR-E5004` der Fehlerklassen; gelesen, Exit-Code 5 bleibt) · `SPEC-034` (Meldungskatalog; gelesen, kein neuer Code)

**Entscheidung:** Entscheidung des Nutzers vom 2026-10-20 (Memory `play-voller-feldvergleich-per-option`): Der Standard bleibt die erste Abweichung; eine **neue Option** schaltet die vollständige Ausgabe aller abweichenden Felder je Interaktion zu. Der Name (Vorschlag `--diff-all`), die Form der Ausgabe und die Frage des Lastenhefts sind **nicht** entschieden; sie stehen als Fragen in §6 und sind vor dem Start (§4) zu entscheiden. Keine ADR.

**Verantwortlich:** —
**Autor:** Planner. **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Das gebaute Binary meldet bei `play --compare-responses` mit einer neuen Option (Name offen, Vorschlag `--diff-all`) je Interaktion **alle** Abweichungen der verglichenen Felder statt der ersten; ohne die Option bleibt es bei der ersten Abweichung, Exit-Code 5 und die Zähler der Zusammenfassung bleiben. Das ist eine **Änderung der Spezifikation**, kein Fehler im Code: Die Spezifikation sagt heute „Je Interaktion wird die erste Abweichung gemeldet“ (`LH-FA-24.a` *Meldung*) und sieht die Option nicht vor; der Architect legt sie vor dem Start fest (§4, §6). Der Slice liefert die Sammlung der Abweichungen im Vergleich des Services und die Option der Kommandozeile mit Hilfetext, Umgebungsvariable und Schlüssel der Konfigurationsdatei.

**Folge-Slice:** `slice-play-feldvergleich-doku` — liefert Handbuch und README im Ist-Zustand für die Option (`AGENTS.md` §3.11; Doku ist eine eigene Schicht, §8). Er steht in `open/`, nimmt die Sendung in seinem §1 an und folgt direkt hinter diesem Slice in derselben Welle bzw. wellenlos; ein Release mit der Option entsteht erst nach beiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Vergleich der **Zeilenwerte** und der **Zeilenzahl** — Bestand bleibt bewusst stehen (Klasse 2). Grund: `LH-FA-24` schließt den Vergleich der Zeilenwerte ausdrücklich aus (Out-of-Scope) und sagt, Werte und Anzahl der Zeilen würden nicht verglichen; „alle abweichenden Felder“ meint hier die Felder, die der Vergleich heute schon liest (Nachrichtenart, Spalten, Parameter-Typen, Befehl, SQLSTATE, Transaktionsstatus). Wollte der Nutzer Werte vergleichen, wäre das eine Änderung des Lastenhefts (Major, `harness/conventions.md` §Versionierung des Lastenhefts) und ein eigener Vorgang (Klasse 3); siehe Frage 3 in §6.
- Die Standard-Ausgabe (erste Abweichung) — Bestand (Klasse 2). Grund: Entscheidung des Nutzers vom 2026-10-20: Der Standard bleibt.
- Neue Meldungscodes, ein anderer Exit-Code als 5 und neue Attribute der Zeile zum Ende — Bestand (Klasse 2). Grund: Die Abweichung bleibt `PGR-E5004`; `abweichend` zählt weiter Interaktionen, nicht Abweichungen. Hält der Architect einen Zähler der Abweichungen für nötig, ist das eine Randform in §6, kein Zusatz des Implementers.
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
- [ ] **Liefer-Punkt 2 — die Option** ([`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen)): `play` nimmt die Option an, mit Umgebungsvariable `PGWIRE_RECORDER_<NAME>` und Schlüssel der Konfigurationsdatei nach dem Muster von `--compare-responses` (Optionstabelle der Spezifikation); ihre Verwendung ohne `--compare-responses`, bei `record` und `replay`, ihr Wert und der Hilfetext von `play` verhalten sich wie die Spezifikation es festlegt, die Meldung bleibt `PGR-E2001` bzw. `PGR-E2004` ohne neuen Code. Ein Test je Zusage und je Element der Aufzählungen (Kommandozeile, Umgebung, Datei, Wertemenge, Kombination, Kommandos), mit Abdeckungs-Deklaration; `make abdeckung-check` und `make meldungskatalog-check` sind grün. Die Option liegt im Record-/Replay-/Play-Pfad, soweit sie dessen Verhalten ändert: Mutation je Zusage, die das Verhalten bestimmt (Wirkung, Kombination); für den Hilfetext und das Annehmen des Namens genügt ein Test je Zusage (`AGENTS.md` §3.10).
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
| `internal/hexagon/services/play_vergleich.go` | update | neben `erste` eine Funktion, die alle Abweichungen der normalisierten Folge liefert (`unterschied` je Position, Zuordnung und Ende des Servers nach der Spezifikation); `text()` für die Liste und die Obergrenze; `erste` und seine Meldung bleiben |
| `internal/hexagon/services/play.go` | update | Feld der `PlayOptionen` für die Option; der Aufrufstelle des Vergleichs wählt `erste` oder die Liste; Meldung `PGR-E5004` je Interaktion; Zähler `abweichend` unverändert |
| `internal/hexagon/services/play_vergleich*_test.go` | update/neu | Fälle je Zusage und je Element, Mutationen aus DoD 1; Standardfall unverändert |
| `internal/adapters/driving/cli/cli.go` | update | Option in der Tabelle der Optionen von `play` (Wahrheitswert, Standard `false`, Umgebung, Schlüssel), Kombination mit `--compare-responses` in der Reihenfolge der Prüfungen, Hilfetext von `play` |
| `internal/adapters/driving/cli/play_test.go` und Nachbarn | update/neu | Fälle je Zusage und Element aus DoD 2 |
| `test/integration/play_vergleich_e2e_test.go` | update | E2E gegen das gebaute Binary und PostgreSQL: eine Interaktion mit mehreren Abweichungen, mit und ohne Option; mit Abdeckungs-Deklaration |
| `docs/user/abdeckung-*.md` | update (erzeugt) | `make abdeckung` nach den neuen Deklarationen |

**Nicht in dieser Liste, weil vor dem Start:** `spec/spezifikation.md` (`LH-FA-24.a`, `LH-FA-20.a`, Optionstabelle, Hilfe von `play`, Zeile der Historie) und, falls entschieden, `spec/lastenheft.md` — die Festlegung aus §6 steht dort, im Commit des Architect **vor** dem ersten Code-Commit (`AGENTS.md` §3.12); dieser Commit gehört nicht zum Slice, er ist sein Start-Trigger (§4). Ändert der Architect die Spezifikation **während** des Slice (Rückgabe durch den Implementer), zählt sie ab dann als dritte Schicht — dann greift der vorgemerkte Schnitt aus §4.

Nicht in dieser Liste: `internal/bootstrap` (die Option wird als Feld durchgereicht, reine Verdrahtung, zählt nicht, `AGENTS.md` §3.13; liest der Bootstrap Datei oder Umgebung selbst, ist das Logik und der Plan prüft die Grenze neu), Handbuch und `README.md` (Folge-Slice), Driven-Adapter, Gates, Verträge unter `harness/vertraege/`, `internal/hexagon/model` (kein neuer Meldungscode, keine Logik).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): (1) Der Nutzer hat die Fragen aus §6 entschieden (Name der Option, Lastenheft ja oder nein, Form der Ausgabe); eine Entscheidung steht hier erst, wenn der Nutzer sie getroffen hat. (2) Der Architect hat die Randformen aus §6 in `spec/spezifikation.md` (und das Lastenheft, falls entschieden) festgelegt und die Liste geprüft (`AGENTS.md` §3.12). (3) Der Doku-Folge-Slice `slice-play-feldvergleich-doku` liegt in `open/`. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik und Messbefehl: `AGENTS.md` §3.14; Pfade des Plans: `internal/hexagon/services internal/adapters/driving/cli test`); die Basis ist der Commit des Übergangs nach `in-progress/`, sein Hash wird beim Übergang hier eingetragen: `<Hash>`. Dann hält der Implementer an und meldet dem Planner.
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

**Fragen an den Nutzer** (vor dem Start, §4; je eine Empfehlung des Planners — Empfehlung, nicht Entscheidung):

1. *Name der Option.* Empfehlung: `--diff-all` (Vorschlag des Nutzers; kurz, benennt die Wirkung; Umgebungsvariable `PGWIRE_RECORDER_DIFF_ALL`, Schlüssel `diff_all`). Alternative: `--report-all-differences`, länger und eindeutiger. Das Kürzel `-d` wird nicht belegt (kein Kurzoptions-Schema, Lehre aus `slice-version-option`).
2. *Lastenheft ändern?* Empfehlung: **nein.** `LH-FA-24` verlangt, „Abweichungen zu melden“, nennt Interaktion und Art und legt „erste“ nicht fest; die Beschränkung auf die erste Abweichung steht allein in der Spezifikation (`LH-FA-24.a`), die Optionsnamen legt das Technik-Stratum fest. Wertet der Nutzer die Option als neue Anforderung, ist das ein Minor-Bump (neue Anforderung oder neues Akzeptanzkriterium; `harness/conventions.md` §Versionierung des Lastenhefts) vor dem Start. Der Architect prüft das gegen den Wortlaut.
3. *Form der Ausgabe und Umfang.* Empfehlung: **eine** Zeile `PGR-E5004` je Interaktion (die Regel „je Interaktion eine Zeile“ und der Zähler `abweichend` bleiben), die Meldung nennt hinter der ersten Abweichung die weiteren als Liste `Position N, Art …` in aufsteigender Position, bei der Art Fehler mit beiden SQLSTATE, bei `Spaltenbeschreibung` je abweichende Spalte (Nummer, Feld: Name oder Typ-OID, beide Seiten), bei `Befehl` beide Befehle, bei `Transaktionsstatus` beide Zeichen; **nie Zeilenwerte oder Fehlertext** (`LH-FA-24` Out-of-Scope, `LH-FA-24.a` *Meldung*); Obergrenze 20 Positionen je Interaktion, danach „und N weitere“ mit der Gesamtzahl. Ein Vergleich der Zeilenwerte ist **nicht** Teil dieser Empfehlung und wäre eine Änderung des Lastenhefts (Major) — sagt der Nutzer „Feldvergleich“ und meint Werte, ist das ein anderer Slice.

**Randformen** (`AGENTS.md` §3.12; `--compare-responses` liegt im Play-Pfad, also im Geltungsbereich). Je Randform steht der **Vorschlag** des Planners; **entschieden ist keine** — der Architect legt jede in `spec/spezifikation.md` fest und prüft die Liste vor dem Code; ein Commit mit Code, Tests oder Gates trägt sie nicht ein. Geschwister-Optionen: `--compare-responses` (`LH-FA-24.a`, Optionstabelle) und `--allow-recorded-errors` (Kombination ohne Wirkung, ohne dass sie ein Fehler ist).

1. *Option ohne `--compare-responses`.* Vorschlag: `PGR-E2001`, Exit-Code 2 (ungültige Kombination; Muster `--upstream-ca` ohne TLS), geprüft nach der Zusammenführung der Quellen, in der Reihenfolge der Kombinationsprüfungen. Alternative: ohne Wirkung wie `--allow-recorded-errors`; dagegen spricht, dass eine stille Option dem Anwender vorspiegelt, es werde verglichen. Zu klären: Platz in der Reihenfolge der Gründe, auch bei Kommandozeile `--diff-all` und Umgebung `PGWIRE_RECORDER_COMPARE_RESPONSES=false` und umgekehrt.
2. *Quellen und Wertemenge.* Kommandozeile `--diff-all[=true|false]`, Umgebungsvariable, Schlüssel im Abschnitt `play:` der Datei (Muster der Optionstabelle, Zeile von `--compare-responses`); Rangfolge Kommandozeile vor Umgebung vor Datei; `true` und `false` genau, jeder andere Wert und der leere ist `PGR-E2001` (Umgebung, Kommandozeile) bzw. `PGR-E2004` (Datei); bei `record` und `replay` ist Option und Umgebungsvariable eine unbekannte Option bzw. fremde Variable (`PGR-E2001`), der Schlüssel in ihrem Abschnitt `PGR-E2004`.
3. *`--continue-on-error`.* Vorschlag: Mit der Option steht in jeder Interaktion mit Abweichung die vollständige Liste; ohne `--continue-on-error` endet `play` nach der Interaktion mit der Abweichung wie heute, deren Liste vollständig ist; Exit-Code 5 unverändert.
4. *`--allow-recorded-errors`.* Unverändert: mit `--compare-responses` ohne Wirkung; die Option ändert nicht, welche Abweichungen zählen (`LH-FA-24.a` *Ende der Aufzeichnung*).
5. *`--fail-on-unconsumed`.* Entfällt: eine Option von `replay`; bei `play` `PGR-E2001`, vom Bestand gedeckt (§1).
6. *Verglichene und ausgegebene Felder.* Dieselben wie heute: Nachrichtentyp, Spaltenanzahl und je Spalte Name und Typ-OID, Typ-OIDs der Parameter, Befehl ohne Zahlen, SQLSTATE, Transaktionsstatus; nie Zeilenwerte, Zeilenzahl, Hinweise, Einstellungen, Fehlertext. Reihenfolge in der Meldung: aufsteigende Position, innerhalb einer Position die Felder in der Reihenfolge der Tabelle `LH-FA-24.a`.
7. *Zuordnung der Positionen.* Vorschlag: Position für Position in der normalisierten Folge, wie die erste Abweichung heute gefunden wird. Fehlt oder fehlt zusätzlich eine Nachricht in der Mitte, verschiebt sich der Rest, und die Folgepositionen sind Folgen der Verschiebung; der Vorschlag meldet sie trotzdem (kein Abgleich der Folgen, kein Diff-Algorithmus), die Spezifikation sagt es. Haben die Folgen verschiedene Länge, ist jede Position hinter der kürzeren Folge eine Abweichung der Art Nachrichtenart, die erste zählt; ob jede einzeln steht oder eine Position „ab hier fehlen N Nachrichten“ meldet, entscheidet der Architect.
8. *Obergrenze der Ausgabe.* Vorschlag: 20 Positionen je Interaktion, dann „und N weitere“ mit Gesamtzahl; die Grenze ist eine Zahl der Spezifikation, der Fall prüft 20 genau, 21 und 0.
9. *Ende des Servers.* Fehlerantwort `FATAL` oder `PANIC` ist die letzte Position; hinter ihr gibt es nichts zu vergleichen; erwartetes Ende (gleicher SQLSTATE an der aufgezeichneten Stelle) ist keine Abweichung und keine Meldung, auch mit der Option; vorherige Positionen werden vollständig gelistet, die Position des Endes zählt als Art `Fehler`.
10. *Stelle in der Folge.* Je Nachrichtenart die Stellen: `ready_for_query` (letzte Position, Transaktionsstatus), `error_response` (Mitte, `ErrorResponse` in der Extended-Interaktion ändert den Vergleich nicht), `connection_closed` (nur nach `FATAL`/`PANIC`, von `play` angefügt), letzte Gruppe mit `Flush`: Vergleich nur bis zur Antwort auf jede Client-Nachricht; je Stelle gilt die Regel des Bestands, nur dass der Vergleich nach der ersten Abweichung weiterläuft.
11. *Signal.* Erstes Signal: die laufende Interaktion wird mit Vergleich beendet, ihre Liste ist vollständig; zweites Signal unterbricht ohne Vergleich, keine Meldung der Interaktion (Bestand).
12. *Reihenfolge der Zeilen und Zähler.* Die Zeilen `error` stehen wie heute (`LH-FA-20.a`); `verglichen` und `abweichend` zählen Interaktionen und ändern sich nicht (`abweichend` höchstens eine je Interaktion trotz mehrerer Abweichungen).
13. *Hilfetext.* `play --help` nennt die Option und dass sie nur mit `--compare-responses` gilt (Aussage der Spezifikation); die globale Hilfe nennt sie nicht; ein Test je Aussage.
14. *Gleichzeitigkeit und Aufrufe.* Der Vergleich ist je Interaktion, Interaktionen laufen nacheinander (`LH-FA-20.a`); keine Aussage über Ordnung über mehrere Aufrufe. Aufrufarten von `make` entfallen: kein Hook und kein Gate liest das Ergebnis.
15. *Plattform.* Meldung und Exit-Code hängen am Betriebssystem nicht; geprüft nur unter Linux, macOS und Windows sind Grenze an der Zusage (`LH-QA-03`); die Spezifikation hält die Grenze, und der Handbuchsatz nennt nur die geprüfte Plattform.
16. *Serverversion.* Der Vergleich gleicht Versionen nicht aus (`LH-FA-24.a` *Serverversion*); die Option ändert das nicht.

**Risiken:**

- Der Nutzer meint mit „Feldvergleich“ auch die Werte der Zeilen; die Option lieferte dann weniger als erwartet, und das Lastenheft schließt Werte aus — **Ausgang:** offen bis Closure (Frage 3; entscheidet der Nutzer vor dem Start gegen die Empfehlung, ist der Slice neu zu planen).
- Folgeabweichungen nach einer verschobenen Folge machen die Liste lang und irreführend (eine fehlende Nachricht erzeugt viele Positionen) — **Ausgang:** offen bis Closure (Obergrenze und Satz in der Spezifikation; ob die Praxis einen Abgleich der Folgen verlangt, wäre ein Folge-Slice).
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

**Zählung** (Planner, 2026-10-23, `AGENTS.md` §3.13; Marke `AGENTS.md` §3.14): **Liefer-Punkte: 2** (1: der Vergleich sammelt alle Abweichungen; 2: die Option mit Hilfetext, Umgebungsvariable und Schlüssel). **Schichten: zwei** — der Kern des Hexagons mit seinen Tests (`internal/hexagon/services/`: Vergleich und Meldung) und der CLI-Adapter mit seinen Tests (`internal/adapters/driving/cli/`; der E2E-Fall unter `test/integration/` prüft dieselbe Fläche). Die Spezifikation ist eine eigene Schicht, liegt aber vor dem Start und damit außerhalb des Slice (§3, §4). Der Bootstrap reicht ein Feld durch (Verdrahtung, zählt nicht, §3). **Nutzer-Doku (`docs/user/`, `README.md`) wäre die dritte Schicht**: Sie geht in `slice-play-feldvergleich-doku` (`AGENTS.md` §3.11, §3.13). Keine Sendung in diesen Slice eingetragen (kein Geber nennt ihn als Nehmer), die Zählung ist die des Plans selbst.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und Entscheidung gehen dem Code voraus; der Architect legt die Festlegung vor dem ersten Code-Commit fest).
