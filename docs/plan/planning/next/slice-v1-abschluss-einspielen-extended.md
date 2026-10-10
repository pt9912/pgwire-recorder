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

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-18.a` · `SPEC-034` · `SPEC-041` · `ARC-002` · `ARC-007`

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
      Extended-Interaktion ist kein Startfehler `PGR-E6001` mehr (Integrationstest). Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung).
- [ ] Innerhalb einer Extended-Interaktion wartet `play` nach `Sync` auf `ReadyForQuery` und
      nach `Flush` auf die Antwort jeder Client-Nachricht der Gruppe, nicht auf die
      aufgezeichneten Server-Nachrichten; ein Abbruch nach `PGR-E4004` sendet keine weitere
      Gruppe und schließt mit `Terminate`; ein Abbruchsignal endet nach dem `ReadyForQuery`
      der laufenden Extended-Interaktion (Test).
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
| `internal/adapters/driven/postgres` | update | die Client-Nachrichten einer Gruppe senden und die Antworten bis `ReadyForQuery` beziehungsweise bis zur Antwort jeder Client-Nachricht der Gruppe lesen |
| `internal/hexagon/services` (Play-Service), `internal/hexagon/ports/driven` | update | Gruppen einer Extended-Interaktion, Warten nach `Sync` und `Flush`, Abbruch, Fortsetzung, erwarteter Fehler und Signal innerhalb der Interaktion; der Startfehler für Extended-Interaktionen entfällt |
| `internal/hexagon/services` (Tests), `internal/adapters/driven/postgres` (Tests) | update | Warten je Gruppenende, Fehler in der ersten und in einer späteren Gruppe, Signal innerhalb der Interaktion, je Zusage eine Mutation |
| `test/integration` | update | Extended-Szenario nach LH-FA-20 und LH-FA-18; die Tests des Kerns zum Zwischenstand *Aufzeichnung mit Extended-Interaktion* ändern |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |

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
- Braucht das Senden einer Gruppe eine neue Operation am Upstream-Port, ist der Port Teil
  des Diffs; er zählt mit dem Play-Service als eine Schicht wie im Kern (§6 dort, *Größe*),
  sonst wären es drei — **Ausgang:** offen bis Closure.

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

Nachgezählt beim Eintragen der Regel *Handbuch und README beschreiben den Ist-Zustand* aus `slice-v1-abschluss-einspielen-laufsteuerung` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch-Teil liegt im ersten Liefer-Punkt, kein neuer Liefer-Punkt; Handbuch und README zählen als Dokumentation, nicht als Schicht. Liefer-Punkte und Schichten bleiben, wie dieser Plan sie zählt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
