# Slice slice-v1-abschluss-einspielen-extended: Extended-Interaktionen einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0012](../../adr/0012-extended-query-gruppen.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-18.a` · `SPEC-034` · `SPEC-041` · `ARC-002` · `ARC-004` · `ARC-007`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play` spielt Extended-Interaktionen ein — die Nachrichten jeder Gruppe in der aufgezeichneten Reihenfolge, mit dem Warten nach `Sync` und nach `Flush` nach `LH-FA-20.a` *Gruppen* — und hält Abbruch, Fortsetzung, erwarteten Fehler und Abbruchsignal auch innerhalb einer Extended-Interaktion ein.

**Übernimmt:** `slice-v1-abschluss-einspielen` — dessen Teil *Extended* (Marke [E] in dessen
§6) und die Randformen mit Marke [L·E] nach dem zweiten Schnitt vom 2026-10-09 (Option O3
des Architect, dort §6 Risiko *Größe des Kerns*; Entscheidung des Nutzers; dort §1,
*Abgegeben*); die Teile [L·E] kommen über die Abgrenzung von
`slice-v1-abschluss-einspielen-laufsteuerung` (dort §1), der vor diesem Slice liegt. Im
Einzelnen:

- das Senden der Client-Nachrichten jeder Gruppe einer Extended-Interaktion (`LH-FA-20.a`
  Schritt 3, `SPEC-041`); aus Ziel und DoD-Punkt 1 jenes Slice bis zum Schnitt „einfach und
  Extended“, hier mit einem eigenen Extended-Szenario;
- das Warten innerhalb einer Extended-Interaktion: nach `Sync` auf `ReadyForQuery`, nach
  `Flush` auf die Antwort jeder Client-Nachricht der Gruppe, nicht nach den aufgezeichneten
  Server-Nachrichten (`LH-FA-20.a` *Gruppen*);
- der Abbruch nach `PGR-E4004` ohne weitere Gruppe (`LH-FA-20.a` Schritt 6) und das
  Abbruchsignal in einer Extended-Interaktion bis zu ihrem `ReadyForQuery` (*Abbruchsignal*);
- [L·E] von `slice-v1-abschluss-einspielen-laufsteuerung`: Fortsetzung in einer
  Extended-Interaktion (der Server verwirft bis `Sync`, nach einer `ErrorResponse` nur noch
  das Warten auf `ReadyForQuery`, die übrigen Gruppen ohne Warten) und der erwartete Fehler
  bei `--allow-recorded-errors` mit einer `error_response` in jeder Gruppe;
- die Ablösung des Zwischenstands *Aufzeichnung mit Extended-Interaktion* jenes Slice (§6
  unten);
- die Randformen [E] und [L·E] aus §6 jenes Slice (§6 unten).

Dieser Slice ist auch die Adresse der Abgrenzung *Einspielen selbst und Fehlersemantik der
Serverfehler* von `slice-v1-abschluss-antwortvergleich` (dort §1), soweit sie
Extended-Interaktionen betrifft.

**Aufsetzen.** Der Slice setzt auf dem Kern `slice-v1-abschluss-einspielen` auf (Kommando
`play`, Einspielen einfacher Anfragen, Fehlerregeln, Abbruchsignal) und auf
`slice-v1-abschluss-einspielen-laufsteuerung` (die Optionen `--continue-on-error` und
`--allow-recorded-errors` mit ihrer Wirkung bei einfachen Anfragen). Bis zu diesem Slice ist
eine Extended-Interaktion bei `play` eine Interaktion einer Art, die `play` nicht einspielt:
Startfehler `PGR-E6001`, Exit-Code 6, keine Verbindung (§6 des Kerns, *Zwischenstand*);
diesen Zwischenstand ersetzt dieser Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-einspielen-extended-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Der Vergleich der Antworten einer Extended-Interaktion, auch in `Flush`-Gruppen —
  `slice-v1-abschluss-antwortvergleich` (dort DoD-Punkt 2); er setzt das Einspielen voraus.
- Die Optionen der Laufsteuerung selbst und ihre Wirkung bei einfachen Anfragen —
  `slice-v1-abschluss-einspielen-laufsteuerung`; dieser Slice setzt sie voraus und liefert nur
  die Teile innerhalb einer Extended-Interaktion.
- Anmeldung mit Passwort und TLS zum Server — `slice-v1-abschluss-einspielen-anmeldung` und
  `slice-v1-abschluss-einspielen-tls`; getestet wird gegen einen Server ohne Passwort und
  ohne TLS.
- Das Aufzeichnen und Wiedergeben von Extended-Interaktionen — Bestand aus
  welle-extended-query bleibt stehen; hier wird nur das gelesene Format eingespielt.
- Die Regel *Art der Interaktion* in `LH-FA-20.a` — Bestand bleibt stehen: Im Zielstand hat
  sie keinen Fall, sie gilt für eine künftige Art des Formats (Entscheidung des Nutzers, §6).
- Code im CLI-Adapter, im Bootstrap und im PGWire-Adapter — Schicht-Abgrenzung: Gruppen
  senden und Antworten lesen liegt im Upstream-Adapter, Warten, Fehler und Signal innerhalb
  der Interaktion im Play-Service mit seinem Port; eine neue Option gibt es nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Eine Aufzeichnung mit Extended-Interaktionen — DDL
      und DML mit Parametern, Gruppen mit `Sync` und mit `Flush`, gemischt mit einfachen
      Anfragen, über mehrere Sessions — wird gegen eine leere Instanz eingespielt, die weder
      Passwort noch TLS verlangt, und die Datenbank enthält danach deren Wirkung; eine
      Extended-Interaktion ist kein Startfehler `PGR-E6001` mehr (Integrationstest).
- [ ] Innerhalb einer Extended-Interaktion wartet `play` nach `Sync` auf `ReadyForQuery` und
      nach `Flush` auf die Antwort jeder Client-Nachricht der Gruppe, nicht auf die
      aufgezeichneten Server-Nachrichten; ein Abbruch nach `PGR-E4004` sendet keine weitere
      Gruppe und schließt mit `Terminate`; ein Abbruchsignal endet nach dem `ReadyForQuery`
      der laufenden Extended-Interaktion; eine Gruppe, deren Nachrichten und Antworten größer
      sind als die Puffer der Verbindung, verklemmt nicht (Test, Gegendruck-Szenario gegen die
      Instanz).
- [ ] [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler): Mit `--continue-on-error` wartet `play` nach einer Fehlerantwort in einer
      Extended-Interaktion nur noch auf deren `ReadyForQuery`, sendet die übrigen Gruppen wie
      aufgezeichnet ohne Warten und endet mit Exit-Code 4; mit `--allow-recorded-errors` gilt
      der Fehler als erwartet, wenn die aufgezeichnete Interaktion in irgendeiner Gruppe eine
      `error_response` trägt (Test). Beleg in §7 für Punkt 1 bis 3: je Zusage Zusage ·
      Mutation · roter Test (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgres` | update | die Client-Nachrichten einer Gruppe nebenläufig zum Lesen senden (§6, *Senden und Lesen unabhängig*) |
| `internal/hexagon/services` (Play-Service), `internal/hexagon/ports/driven` | update | neue Operation `Gruppe` am Port; Gruppen einer Extended-Interaktion, Zählen der Antworten je Gruppe, Warten nach `Sync` und `Flush`, Abbruch, Fortsetzung, erwarteter Fehler und Signal innerhalb der Interaktion; der Startfehler für Extended-Interaktionen entfällt |
| `internal/hexagon/services` (Tests), `internal/adapters/driven/postgres` (Tests) | update | Warten je Gruppenende, Fehler in der ersten und in einer späteren Gruppe, Signal innerhalb der Interaktion, je Zusage eine Mutation; neue Dateien `play_extended_test.go` und `einspielen_gruppe_test.go`; `TestPlayStart` ändert sich (Zwischenstand) |
| `internal/bootstrap` (Tests) | update | `TestRunPlayStartfehler` verliert den Fall *Extended-Interaktion* (Zwischenstand), der Fall *nicht ladbar* bleibt |
| `test/integration` | update | neue Datei `play_extended_e2e_test.go` (Extended-Szenario nach LH-FA-20 und LH-FA-18, Fehlerfälle, Gegendruck-Szenario, Szenario *Warten nach Flush* über einen Recorder zwischen `play` und der Instanz, Auftrag nach `F-580`); `TestE2EPlayZwischenstand` in `play_e2e_test.go` verliert den Fall *Extended in Session 2* |
| Kommentare und Abdeckungs-Deklarationen (Auftrag nach `F-581`, `F-582`, `F-583`, `F-585`, `F-586`) | update | `TestEinspielGruppeNichtAbbildbar`: Zitat statt `LH-FA-20.a` *Interaktion* der Satz „Verhalten des Adapters für eine nicht abbildbare Nachricht; `Validate` lässt sie nie durch (§6, akzeptiertes Negativ)“, im Test und (über `make abdeckung`) in der Zeile von `abdeckung-unit.md`; Kommentar vor `TestEinspielGruppeSendefehler`: „der nächsten Gruppe“ statt „der nächsten Operation“ (Port-Kommentar sagt es schon so); Deklaration von `TestE2EPlayExtendedGegendruck` (Kommentar und Zeile in `abdeckung-e2e.md`): Zusage „auf einem Host mit Socket-Puffern unter 16 MB“, das Warten nach `Flush` nennt sie nicht mehr; `TestPlayExtendedErwarteterFehler` variiert die Gruppe des Serverfehlers und die der Aufzeichnung; Kommentar vor `sendeGruppe` und der Struct `einspielSession` |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` aus den Deklarationen der Tests erzeugt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` und
`slice-v1-abschluss-einspielen-laufsteuerung` liegen in `done/`. Schritt 10 der Reihenfolge
in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) (zweiter Schnitt vom 2026-10-09).
Die Randformen aus §6 entschied der Architect am 2026-10-09 vor dem Code in `LH-FA-20.a`; vor
dem ersten Code-Commit prüft er die Liste gegen den gelieferten Kern und die gelieferte
Laufsteuerung, besonders gegen die Form, in der der Upstream-Adapter eine Interaktion sendet
und liest (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar. Schnitt dann: Fortsetzung und erwarteter Fehler innerhalb einer
  Extended-Interaktion (DoD-Punkt 3) als eigener Slice; Senden, Warten, Abbruch und Signal
  sind ohne sie lieferbar, weil ohne die Optionen jeder Fehler abbricht.
- `in-progress` → `open` (blockiert — Carveout?): Eine aufgezeichnete Extended-Interaktion
  trägt nicht genug, um die Client-Nachrichten ihrer Gruppen wiederherzustellen (`SPEC-041`);
  dann zuerst die Entscheidung des Architect über das Format.

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

**Randformen** (`AGENTS.md` §3.12) — die mit Marke [E] und [L·E] aus §6 von
`slice-v1-abschluss-einspielen`, geprüft vom Architect am 2026-10-09 vor dem ersten
Code-Commit jenes Slice und für den zweiten Schnitt markiert; die Teile [L·E] liefert und
testet dieser Slice, weil er in der Reihenfolge nach `slice-v1-abschluss-einspielen-laufsteuerung`
liegt. Neu entschieden heißt: im Commit jener Prüfung in `LH-FA-20.a`, sonst an der
genannten Stelle. Offen ist keine.

- **Warten innerhalb einer Extended-Interaktion** [E] — nach `Sync` auf `ReadyForQuery`, nach
  `Flush` auf die Antwort jeder Client-Nachricht der Gruppe (Tabelle der Antworten), nicht
  nach den aufgezeichneten Server-Nachrichten; nach einer `ErrorResponse` nur noch auf das
  `ReadyForQuery`, übrige Gruppen ohne Warten; neu entschieden in `LH-FA-20.a` *Gruppen*. Die
  Antwort hängt nicht an der Datenlage (Zeilenzahl), darum nicht an der Aufzeichnung.
- **Abbruch nach `PGR-E4004`, keine weitere Gruppe** [E] — keine weitere Antwort gelesen,
  keine weitere Nachricht oder Gruppe, `Terminate`; neu entschieden (Lesen) in `LH-FA-20.a`
  *Interaktion*, sonst bestätigt, Schritt 6.
- **Abbruchsignal in einer Extended-Interaktion** [E] — das Einspielen endet nach ihrem
  `ReadyForQuery`; neu entschieden in `LH-FA-20.a` *Abbruchsignal*.
- **Erwarteter Fehler „in jeder Gruppe“** [L·E] — bei `--allow-recorded-errors` entscheidet
  allein, ob die aufgezeichnete Interaktion irgendwo eine `error_response` trägt, in jeder
  Gruppe; keine Meldung; neu entschieden in `LH-FA-20.a` *Interaktion*.
- **Fortsetzung bei einer Extended-Interaktion** [L·E] — der Server verwirft bis `Sync`, die
  übrigen Gruppen ohne Warten, Rest der Interaktion wie aufgezeichnet, Exit-Code 4 am Ende
  nach einem `PGR-E4004`; bestätigt, `LH-FA-20.a` Schritt 6 und *Gruppen*.
- **Zählen der Antworten einer `Flush`-Gruppe** [E] — je Antwort, nicht je Art: `Describe`
  einer Anweisung zählt zwei, `Flush` keine; `DataRow`, `NoticeResponse`, `ParameterStatus`
  und `NotificationResponse` zählen nicht; eine Antwort anderer Art als erwartet ändert das
  Warten nicht; eine `ErrorResponse` beendet das Zählen der Interaktion. Neu entschieden in
  `LH-FA-20.a` *Gruppen* (Commit der Prüfung vom 2026-10-10).
- **Senden und Lesen unabhängig** [E] — Option A (gewählt): `Gruppe` am Port beginnt das Senden
  der Gruppe und kehrt zurück; ein Fehler des Sendens kommt als `PGR-E4003` aus `Naechste`
  oder der nächsten Gruppe, `Schliesse` beendet den Sender, und ein Fehler danach bleibt
  ohne Folge; Nebenläufigkeit je Verbindung ist Sache des Adapters
  ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md), dort Kontext und Option E), der
  Service zählt und entscheidet. Verworfen: B, `Gruppe` sendet vor dem Lesen — ein
  `pgx`-Batch ist eine `Sync`-Gruppe, und sind Eingabe und Ausgabe größer als die Puffer,
  steht das Senden am Server und `play` wartet ohne eigene Frist; C, ein Puffer für die
  Antworten — der Speicher wächst mit der Ausgabe ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md), Option B). Neu entschieden in
  `LH-FA-20.a` *Gruppen*; mit [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)
  vereinbar, das die Anfragen und Sessions nacheinander festlegt, nicht Senden und Lesen
  innerhalb einer Interaktion. Die Gleichzeitigkeitszusage im Kommentar des Ports ändert sich
  mit (`AGENTS.md` §3.11).
- **Parameterwerte, Formate und Namen** [E] — wie aufgezeichnet gesendet: `Null` ist
  SQL-NULL (Länge −1), ein leerer Wert ist ein leerer Wert, nicht NULL, gleich ob die
  geladenen Bytes `nil` oder leer sind; `param_formats` und `result_formats` unverändert,
  auch leer; benannte und unbenannte Anweisung und Portal wie aufgezeichnet. Entschieden in
  `SPEC-041` und `SPEC-003`. Passt eine Anzahl nicht, entscheidet der Server (`PGR-E4004`).
- **Erwarteter Fehler über alle Gruppen** [L·E] — `mitFehlerantwort` prüft auch die
  Server-Nachrichten der Gruppen, nicht nur `Responses` (siehe oben, *Erwarteter Fehler
  „in jeder Gruppe“*).
- **Bestätigt ohne Änderung** [E], `LH-FA-20.a` *Interaktion*: `FATAL`/`PANIC` in einer Gruppe
  `PGR-E4003` sofort; `CopyInResponse`, `CopyOutResponse`, `CopyBothResponse` oder eine nicht
  lesbare Nachricht in einer Gruppe `PGR-E6001`; `NoticeResponse`, `ParameterStatus`,
  `NotificationResponse` zwischen Antworten gelesen und verworfen; Verbindungsende in einer
  Gruppe `PGR-E4003`.
- **Akzeptierte Negative** (kein Code, kein Test; Grund je Punkt): *leere Gruppe*,
  *Gruppe ohne Client-Nachricht* und *Gruppe mit `Sync` mitten in der Interaktion* — `Validate`
  weist sie beim Laden ab (`PGR-E3003`), `play` sieht sie nie. *`CopyData` des Clients* — das
  Format kennt keine solche Nachricht. *Anweisung über Sessions wiederverwenden* — jede Session
  trägt ihr `Parse` selbst; fehlt es nach einem Fehler, antwortet der Server (`26000`,
  `PGR-E4004`) wie bei jedem Fehler. *Vorbereitete Anweisung ohne `Close`* — `Terminate` und
  das Verbindungsende räumen sie. *Pipelining über das `Sync` hinaus* — die aufgezeichnete
  Folge nach einem `Sync` gehört zur nächsten Interaktion (`LH-FA-18.a`); `play` wartet vor ihr
  auf das `ReadyForQuery` ([ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)), die Wirkung auf dem Server ist dieselbe. *`ReadyForQuery`
  in einer `Flush`-Gruppe* — ein Server sendet es ohne `Sync` nicht, ein Protokollbruch; `play`
  wartet wie auf jede Antwort ohne Frist, das zweite Signal beendet es. *Antwort anderer Art
  als erwartet* — gezählt wird je Antwort. *Serverversion* —
  `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, Referenzversion 17.
  *Nicht abbildbare Client-Nachricht* (Zielart außer `statement` und `portal`, unbekannter
  Typ) — `Gruppe` liefert `PGR-E1000` und sendet nichts von der Gruppe, wie der
  Upstream-Adapter es für diese Nachricht schon tat; `Validate` lässt sie beim Laden nie durch
  (`PGR-E3003`), der Pfad ist im Zielstand unerreichbar, darum steht er nicht in `LH-FA-20.a`
  (entschieden vom Architect am 2026-10-10, auf `F-582`); `TestEinspielGruppeNichtAbbildbar`
  hält das Verhalten des Adapters, nicht eine Regel der Spezifikation. *Unbegrenzte
  Warteschlange der Gruppen im Adapter* — sie hält nur Verweise auf die geladene Aufzeichnung,
  und der Ablauf hält sie klein: `play` liest nach jeder `Flush`-Gruppe, ehe es die nächste
  einreiht, und reiht nur nach einer Fehlerantwort die übrigen Gruppen der einen Interaktion
  ein (`F-583`). *Gegendruck-Test hängt an der Puffergröße des Hosts* — der Mutant „`Gruppe`
  sendet selbst“ verklemmt nur bei Socket-Puffern unter der Summe der Parameter und der
  Ausgabe (16 und 32 MB); ein Host mit größeren Puffern lässt den Mutanten grün. Die Zusage
  des Tests gilt darum *auf einem Host mit Puffern unter 16 MB*; Puffer klein zu halten
  bräuchte Socket-Optionen im Docker-Netz und wöge mehr als der Rest des Falls (`F-585`).
  *Goroutinen-Zählung im Adapter-Test* — `senderZahl` zählt im ganzen Prozess und hängt an
  jedem Test des Pakets, der `Schliesse` nicht erreicht; die Meldung nennt das (`F-585`).
- **Ablösung des Zwischenstands *Aufzeichnung mit Extended-Interaktion*** — bis zu diesem
  Slice ist eine Extended-Interaktion ein Startfehler `PGR-E6001` (§6 von
  `slice-v1-abschluss-einspielen`, *Zwischenstand*, allgemein gefasst in `LH-FA-20.a` *Art
  der Interaktion*). Dieser Slice ersetzt ihn durch das Zielverhalten (Schritt 3 für
  Extended-Interaktionen) und ändert die Tests des Kerns dazu; die Spezifikation ändert er
  dafür nicht. *Akzeptiertes Negativ:* Im Zielstand hat die Regel *Art der Interaktion* keinen
  Fall; sie bleibt als allgemeine Regel stehen (Entscheidung des Nutzers) und gilt für eine
  künftige Art des Formats.

**Risiken:**

- Das Warten nach `Flush` folgt der Tabelle der Antworten je Client-Nachricht; nennt sie für
  eine Nachricht eine Antwort, die der Server nicht sendet, wartet `play` ohne eigene Frist
  (`LH-FA-20.a` *Interaktion*). Belegbar nur gegen einen echten Server, und nur gegen die
  Referenzversion 17 (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, §8) —
  **Ausgang:** offen bis Closure.
- Das Senden einer Gruppe braucht eine neue Operation am Upstream-Port (`Gruppe`, entschieden
  oben); der Port ist Teil des Diffs und zählt mit dem Play-Service als eine Schicht wie im
  Kern (§6 dort, *Größe*), sonst wären es drei — **Ausgang:** offen bis Closure.

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

**Stand und Läufe.** Code- und Test-Stand: Commit `1d045fa` (Code, Tests, Plan §3) und `a9bc23a` (Abdeckungstabelle); auf `a9bc23a` lief `make gates` grün (Exit 0; darin `make test`, `make test-integration`, `make lint`, `make lint-gegenprobe`, `make a-check`, `make a-check-negativ`, `make kopf-check`, `make docs-check`, `make abdeckung-check`). Danach kamen `TestEinspielGruppeSchliesseVerwirftWarteschlange`, die Abdeckungstabelle dazu und dieser Abschnitt; der Lauf von `make gates` auf dem Commit, der sie trägt, steht mit Hash im Bericht (ein Commit kann seinen eigenen Hash nicht nennen). `make abdeckung` hat die Tabellen geschrieben. Größe des Diffs gegen `b23ff2c`: rund 1500 Zeilen eingefügt, rund 100 geändert oder entfernt, davon rund 1150 Zeilen in neuen Testdateien. Integrationslauf mit eigenem Netz und Container, vom Runner entfernt.

**Mutationen.** Jede Mutation lief in einer frischen Kopie (`cp -r` ohne `-p`, gofmt-sauber, eine Änderung je Kopie), die Unit-Mutationen über `go test` im Image der Stufe `deps`, die Integrations-Mutationen über einen Build der Stufe `integration` der Kopie. Rot heißt: der genannte Test schlägt an der genannten Stelle fehl (Ablauf der Aufrufe weicht ab, Meldung oder Bytes weichen ab); keine Gesamtfrist trägt eine Zeile außer den ausdrücklich genannten Hänge-Fällen.

*Warten innerhalb einer Extended-Interaktion (Service, `play.go`)* — jede Bedingung eine Zeile:

| Zusage | Mutation | roter Test |
|---|---|---|
| auf `Describe` einer Anweisung zwei Antworten | `n += 2` → `n++` | `TestPlayExtendedWarten` |
| ebenso gegen den echten Server: nicht mehr als zwei | `n += 2` → `n += 3` | `TestE2EPlayExtended` (hängt, Frist 30 s) |
| auf `Flush` keine Antwort | `Flush` zählt eine | `TestPlayExtendedNurFlush`, `TestPlayExtendedWarten`, `TestPlayExtendedAndereArt` |
| ebenso gegen den echten Server | `Flush` zählt eine | `TestE2EPlayExtended` (hängt, Frist 30 s) |
| auf `Parse` eine Antwort | `Parse` zählt keine | `TestPlayExtendedWarten` |
| auf `Bind` eine Antwort | `Bind` zählt keine | `TestPlayExtendedWarten` |
| auf `Close` eine Antwort | `Close` zählt keine | `TestPlayExtendedWarten` |
| auf `Execute` eine Antwort | `Execute` zählt keine | `TestPlayExtendedWarten`, `TestPlayExtendedAndereArt` |
| auf `Describe` eines Portals eine Antwort | `Describe` zählt keine (Anweisung eingeschlossen) | `TestPlayExtendedWarten` |
| `DataRow` zählt nicht | `DataRow` zählt | `TestPlayExtendedWarten` |
| `NoticeResponse` zählt nicht | `NoticeResponse` zählt | `TestPlayExtendedWarten` |
| `ParameterStatus` zählt nicht | `ParameterStatus` zählt | `TestPlayExtendedWarten` |
| die aufgezeichneten Server-Nachrichten bestimmen das Warten nicht | `n := len(g.Server)` | `TestPlayExtendedWarten`, `TestPlayExtendedAndereArt` |
| eine Antwort anderer Art ändert das Warten nicht | `PortalSuspended` wird nicht gezählt | `TestPlayExtendedAndereArt` |
| nach `Flush` wird gewartet, bevor die nächste Gruppe geht | kein Warten nach `Flush` (`… || true` in der Bedingung vor der Leseschleife) | `TestPlayExtendedWarten` |
| ebenso gegen den echten Server (Recorder zwischen `play` und der Instanz: die Antworten der ersten Gruppe stehen in deren `server:`-Liste) | kein Warten nach `Flush` | `TestE2EPlayExtendedFlushWarten` (rot: die Antworten stehen in der `server:`-Liste der zweiten Gruppe) |
| ebenso | die nächste Gruppe geht vor dem Lesen der Gruppe | `TestPlayExtendedWarten` |
| nach `Sync` wird bis zum `ReadyForQuery` gelesen, ehe die nächste Interaktion beginnt | nach der letzten Gruppe kein Lesen | `TestPlayExtendedWarten`, `TestPlayExtendedAndereArt` |
| `Gruppe` und `Naechste` liefern ihre Fehler weiter | Fehler von `Gruppe` ignoriert; Fehler von `Naechste` ignoriert | `TestPlayExtendedSendenLesenScheitert` (je Mutation) |
| einfache Anfrage und Extended-Interaktion laufen nacheinander in einer Session | `Extended` wird wie eine einfache Anfrage gesendet | `TestPlayExtendedWarten`, `TestPlayExtendedGemischt` |
| der Zwischenstand (`PGR-E6001` für Extended) entfällt | die Prüfung des Startfehlers wird wieder eingebaut | `TestPlayExtendedWarten`, `TestPlayStart` |

*Fehler, Fortsetzung, erwarteter Fehler, Signal (Service)*:

| Zusage | Mutation | roter Test |
|---|---|---|
| ohne Optionen bricht eine Fehlerantwort ab: keine weitere Antwort gelesen | nach dem Fehler wird bis `ReadyForQuery` gelesen | `TestPlayExtendedFehlerBrichtAb` |
| ebenso: keine weitere Gruppe gesendet | nach dem Fehler werden die übrigen Gruppen gesendet | `TestPlayExtendedFehlerBrichtAb` |
| mit `--continue-on-error` zählt `play` nach der Fehlerantwort nicht mehr (übrige Gruppen ohne Warten dazwischen) | der Fehler beendet das Zählen nicht | `TestPlayExtendedFortsetzung` |
| ebenso: danach wird bis `ReadyForQuery` gelesen, ehe die nächste Interaktion beginnt | nach einem Fehler kein Lesen bis `ReadyForQuery` | `TestPlayExtendedFortsetzung` |
| erwarteter Fehler: error_response in der ersten Gruppe der Aufzeichnung | nur die erste Gruppe wird geprüft (`i > 0` übersprungen) — rot an den Fällen mit der Aufzeichnung in der mittleren und der letzten Gruppe | `TestPlayExtendedErwarteterFehler` |
| ebenso: in einer mittleren Gruppe | nur die mittlere Gruppe wird geprüft (`i != 1` übersprungen) | `TestPlayExtendedErwarteterFehler` |
| ebenso: in der letzten Gruppe | nur die letzte Gruppe wird geprüft (`i != len-1` übersprungen) | `TestPlayExtendedErwarteterFehler` |
| ebenso: die erste Gruppe zählt | die erste Gruppe wird nicht geprüft (`i == 0` übersprungen) | `TestPlayExtendedErwarteterFehler` |
| ebenso: die letzte Gruppe zählt | die letzte Gruppe wird nicht geprüft (`i == len-1` übersprungen) | `TestPlayExtendedErwarteterFehler` |
| ebenso: in irgendeiner Gruppe, nicht nur in `Responses` | die Gruppen werden nicht geprüft | `TestPlayExtendedErwarteterFehler` |
| erwartet gilt, gleich in welcher Gruppe der Server den Fehler sendet: in einer Gruppe mit `Flush` nach der ersten | die Fehlerantwort in einer `Flush`-Gruppe ist nur in der ersten Gruppe erwartet (`erwartet && g.Client[0].Type == model.ClientParse`) | `TestPlayExtendedErwarteterFehler` (Fälle mit Server in der mittleren Gruppe) |
| ebenso: in der Gruppe mit `Sync` | die Fehlerantwort beim Lesen bis `ReadyForQuery` ist nie erwartet (`false` statt `erwartet` in `bisBereit`) | `TestPlayExtendedErwarteterFehler` (Fälle mit Server in der letzten Gruppe) |
| ebenso: der Fehler setzt `fehlerGesehen` | `fehlerGesehen = false` statt `true` | `TestPlayExtendedErwarteterFehler` (Aufzeichnung erste, Server erste), `TestPlayExtendedFortsetzung` |
| `FATAL`/`PANIC` ist `PGR-E4003`, auch wenn der Fehler erwartet wäre | die Prüfung von `FATAL` entfällt; sie steht hinter der des erwarteten Fehlers | `TestPlayExtendedFatal` (je Mutation) |
| Abbruchsignal: die Extended-Interaktion läuft bis zu ihrem `ReadyForQuery` | `ctx` wird zwischen den Gruppen geprüft | `TestPlayExtendedSignal` |

*Port-Operation `Gruppe` (Adapter, `einspielen.go`)*:

| Zusage | Mutation | roter Test |
|---|---|---|
| `Gruppe` kehrt zurück, ohne auf das Senden zu warten | `Gruppe` sendet selbst | `TestEinspielGruppeUnabhaengig` (Frist 5 s) |
| ebenso gegen die Instanz mit großer Gruppe (Gegendruck), auf einem Host mit Socket-Puffern unter 16 MB | `Gruppe` sendet selbst | `TestE2EPlayExtendedGegendruck` (hängt, Frist 60 s) |
| das Lesen ist unabhängig vom Senden | `Naechste` hält die Sperre des Sendens | `TestEinspielGruppeUnabhaengig`, `TestEinspielGruppeNachrichten` |
| mehrere Gruppen gehen in der Reihenfolge der Aufrufe | die Warteschlange wird rückwärts abgearbeitet | `TestEinspielGruppeNachrichten`, `TestEinspielGruppeUnabhaengig` |
| ein Parameter mit `Null` ist SQL-NULL | `Null` wird zu leer | `TestEinspielGruppeNachrichten` |
| ein leerer Wert ist leer, nicht NULL (Bytes nil und leer) | leere Bytes werden zu NULL | `TestEinspielGruppeNachrichten` |
| `param_formats` unverändert | `ParameterFormatCodes: nil` | `TestEinspielGruppeNachrichten` |
| `result_formats` unverändert | `ResultFormatCodes: nil` | `TestEinspielGruppeNachrichten` |
| `Execute` mit `max_rows` | `MaxRows` entfällt | `TestEinspielGruppeNachrichten` |
| Zielart von `Describe` und `Close` | `S` und `P` vertauscht | `TestEinspielGruppeNachrichten` |
| Namen von Anweisung und Portal wie aufgezeichnet | Portal und Anweisung in `Bind` vertauscht; `ParameterOIDs` in `Parse` entfällt | `TestEinspielGruppeNachrichten` (je Mutation) |
| Sendefehler ist `PGR-E4003` aus `Naechste`, auch wenn es wartet | `Naechste` fragt den Fehler nach dem Lesefehler nicht ab | `TestEinspielGruppeSendefehler` |
| ebenso: Naechste liest nicht weiter, wenn die Antwort schon bereitliegt | `Naechste` fragt den Fehler vor dem Lesen nicht ab | `TestEinspielGruppeSendefehler` |
| ebenso: die nächste Gruppe liefert den Fehler | `Gruppe` fragt den Fehler nicht ab | `TestEinspielGruppeSendefehler` |
| ein Sendefehler schließt die Verbindung, ein wartendes `Naechste` endet | der Sender schließt nicht | `TestEinspielGruppeSendefehler`, `TestEinspielGruppeSchliesseBeimSenden` (Frist 5 s) |
| `Schliesse` beendet den Sender (Sender läuft nicht weiter) | `Schliesse` schließt `ende` nicht | `TestEinspielGruppeSchliesseBeendetSender` |
| `Schliesse` wartet nicht auf das Senden und schreibt kein `Terminate` dazwischen | `Schliesse` nimmt die Sperre mit `Lock` statt `TryLock` | `TestEinspielGruppeSchliesseBeimSenden`, `TestEinspielSchliesseBeimSenden` (Frist 5 s) |
| ein Fehler nach `Schliesse` bleibt ohne Folge | der Sender merkt auch nach `Schliesse` einen Fehler | `TestEinspielGruppeSchliesseBeimSenden` |
| nach `Schliesse` nimmt die Session keine Gruppe mehr an | die Prüfung entfällt | `TestEinspielGruppeSchliesseBeimSenden` |
| `Schliesse` verwirft die eingereihten Gruppen: der Sender schreibt nach `Schliesse` nichts mehr | die Prüfung von `geschlossen` vor dem Senden entfällt | `TestEinspielGruppeSchliesseVerwirftWarteschlange` |

**Grüne Mutanten.** Keiner, der in den Zeilen oben fehlt; jeder gefahrene Mutant war rot. Nicht gegen den echten Server gezeigt: Eine Fehlerantwort in einer `Flush`-Gruppe (die Fehlerfälle gegen die Instanz enden in der `Sync`-Gruppe, weil `pgx` so sendet); sie ist nur im Unit-Test belegt (`TestPlayExtendedFortsetzung`, `TestPlayExtendedFehlerBrichtAb`); die Mutation `n = 0` → `n--` im Zählen blieb in `TestE2EPlayExtendedFehler` grün, weil dort kein Fehler in einer `Flush`-Gruppe steht, und ist im Unit-Test rot (Zeile *mit `--continue-on-error` zählt `play` nicht mehr*).

**Gegendruck-Szenario (DoD Punkt 2).** `TestE2EPlayExtendedGegendruck`: eine `Sync`-Gruppe mit 32 MB Ausgabe, dann 16 MB Parameter, eine `Flush`-Gruppe mit verzögerter 32-MB-Ausgabe vor einer `Sync`-Gruppe; Frist des Laufs 60 s als Literal (`starteBisFrist`); die Datenbank trägt danach die Länge des Parameters und 7. Der Lauf dauert gegen die Instanz rund 2 s; mit synchron sendendem `Gruppe` hängt er (Zeile oben). Die Zusage gilt auf einem Host mit Socket-Puffern unter 16 MB (`F-585`, §6); das Warten nach `Flush` belegt dieser Test nicht, das tut `TestE2EPlayExtendedFlushWarten`.

**Randformen.** Keine neu entschieden. Bei der Umsetzung angefallen und dem Architect gemeldet; er hat sie am 2026-10-10 als akzeptierte Negative in §6 entschieden (Review `F-582`, `F-583`): (1) Eine Client-Nachricht, die sich nicht abbilden lässt (Zielart außer `statement` und `portal`, unbekannter Typ), liefert `Gruppe` als `PGR-E1000` und sendet nichts von der Gruppe; `Validate` lässt sie beim Laden nie durch; `TestEinspielGruppeNichtAbbildbar` hält das Verhalten des Adapters, nicht eine Regel der Spezifikation. (2) Ein `ReadyForQuery` in einer `Flush`-Gruppe zählt wie jede andere Antwort. (3) Die Warteschlange der Gruppen im Adapter ist unbegrenzt; sie hält nur Verweise auf die geladene Aufzeichnung, und der Ablauf hält sie klein.

**Nacharbeit zum Review (F-580 bis F-586).** *F-580 (MEDIUM):* Der E2E-Beleg war machbar und ist geliefert: `TestE2EPlayExtendedFlushWarten` spielt eine Aufzeichnung (erste Gruppe mit `Flush` und `pg_sleep(1)`, zweite mit `Sync`) über einen `record` zwischen `play` und der Instanz; der Recorder ordnet eine Antwort der zuletzt begonnenen Gruppe zu, also stehen die Antworten der ersten Gruppe nur dann in deren `server:`-Liste, wenn `play` die zweite Gruppe erst nach ihnen sendet. Frist des Laufs 30 s als Literal (`starteBisFrist`, `SPEC-038`). Der Mutant „kein Warten nach `Flush`“ (`… || true` in der Bedingung vor der Leseschleife in `extended`) färbt den Test rot, in einer frischen Kopie ohne `-p`, gofmt-sauber: die Antworten der ersten Gruppe stehen in der `server:`-Liste der zweiten. Kommentar und Zeile von `TestE2EPlayExtendedGegendruck` sagen das Warten nicht mehr zu. *F-581 (LOW):* `TestPlayExtendedErwarteterFehler` variiert beides, die Gruppe der Aufzeichnung und die des Serverfehlers (erste, mittlere, letzte, auch gekreuzt, dazu *keine Gruppe* mit Server in der ersten und der mittleren); Klasse des Merkmals „welche Gruppe wird geprüft“: fünf Mutanten (nur erste, nur mittlere, nur letzte, nicht erste, nicht letzte), alle rot an genau den Fällen, denen die Gruppe fehlt, ohne Panik (die ersten Fassungen mit `Groups[:1]` panikten bei einer einfachen Anfrage und zählten nicht). Die Zeilen stehen in der Tabelle oben. *F-582 und Auftrag §3:* Zitat von `TestEinspielGruppeNichtAbbildbar` ersetzt, Kommentar vor `TestEinspielGruppeSendefehler` sagt „der nächsten Gruppe“, die Abdeckungszeilen sind mit `make abdeckung` nachgezogen. *F-583:* gelesen, kein Auftrag; (3) oben ist als akzeptiertes Negativ in §6. *F-584:* Die Lücke zwischen der Rückkehr von `Flush` im Sender und dem `Unlock`, in der das `TryLock` in `Schliesse` scheitert und ein Abbruch ohne `Terminate` endet, deckt *soweit möglich* (`LH-FA-20.a` Schritt 6); kein Test hält sie, und ich habe keinen Mutanten dafür gefahren. *F-585:* Die Zusage des Gegendruck-Tests gilt auf einem Host mit Socket-Puffern unter 16 MB (Kommentar, Zeile in `abdeckung-e2e.md`, §6); die Goroutinen-Zählung im Adapter-Test bleibt wie sie ist, ihre Meldung nennt die Kopplung. *F-586:* Kommentar vor `sendeGruppe` umgebrochen, der Kommentar von `einspielSession` nennt `senderLaeuft`, `weck` und `ende`.

**Gelaufene Sensoren, die nicht Teil von `make gates` sind.** `make abdeckung` (Tabellen geschrieben und von `make abdeckung-check` im Gate bestätigt); die Mutationen oben. Nicht gelaufen und für diesen Slice nicht berührt: `make a-check-graph`, `make doc-trace` (Werkzeuge, kein Gate).

**Offen für die Closure.** Die zwei Risiken aus §6 tragen noch keinen Ausgang. Zu *Warten nach `Flush` gegen den echten Server*: Gegen PostgreSQL 17 (gepinntes Image) laufen die Pipeline mit `Prepare` + `Flush` und das Gegendruck-Szenario durch; die Tabelle der Antworten stimmte dort für `Parse`, `Describe` einer Anweisung (zwei Antworten) und `Execute`; für `Bind`, `Close` und `Describe` eines Portals in einer `Flush`-Gruppe gibt es keinen Lauf gegen den Server (nur Unit-Tests gegen die Tabelle). Zu *Port ist Teil des Diffs*: Der Diff berührt Play-Service, Port und Upstream-Adapter.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `19f5512` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×) — das Warten nach `Flush`
  hängt am Verhalten des Servers und läuft nur gegen die Referenzversion 17 (§6 Risiko); der
  Eintrag bleibt unter der Schwelle, Beleg erst bei Closure, falls ein Review eine
  Abweichung findet.
- `BEO-REPO/slice-waechst-durch-uebernahmen`, `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel`
  und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` — betreffen den Schnitt des Gebers und
  werden dort gezählt (§8 von `slice-v1-abschluss-einspielen`). Nachgezählt beim Eintragen
  der Sendungen: drei Liefer-Punkte, zwei Schichten (Play-Service mit Port,
  Upstream-Adapter).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12),
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`,
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (§3.11),
  `BEO-REPO/plan-folgt-korrektur-nicht` (§3.9) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (§3.13) — verkörpert; die Randformen in §6
  sind vor dem Code entschieden, je Zusage eine Mutation mit Beleg in §7, die Geber
  `slice-v1-abschluss-einspielen` und `slice-v1-abschluss-einspielen-laufsteuerung` zeigen im
  selben Commit hierher.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-einspielen-extended-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 3. Schichten: zwei Schichten (Kern mit Play-Service und Port, Upstream-Adapter).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
