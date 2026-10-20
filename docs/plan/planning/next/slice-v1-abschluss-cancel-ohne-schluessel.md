# Slice slice-v1-abschluss-cancel-ohne-schluessel: CancelRequest ohne Schlüssel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)

**Berührte Spec-Stellen:** `LH-FA-05.e` (Spezifikation, Protokollrand, *Länge des `CancelRequest`*) · `LH-FA-13.b` · `ARC-006`

**Verantwortlich:** pt9912 (Zweig impl-cancel)

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein `CancelRequest` der Länge 12, den ein Client ohne Abbruchschlüssel sendet (pgx, weil der Recorder `BackendKeyData` nicht weitergibt), wird wie jeder `CancelRequest` behandelt: Verbindung geschlossen, Warnung `PGR-W3001`, kein Verbindungsfehler (`LH-FA-05.e`), in `record` und `replay`, mit und ohne TLS. Heute lehnt die PGWire-Bibliothek ihn als „cancel request too short“ ab (gemessen mit v5.11.0: sie nimmt die Längen 16 bis 268 an, die Länge 12 und alle darunter und dazwischen nicht), und der Recorder meldet `PGR-E6001` mit Klasse 6 (Review-Finding F-310 aus `docs/reviews/2026-10-05-review-slice-extended-query-record.md`, Sonde P3). Die Längen, ihre Folgen und die Reihenfolge der Gründe stehen in der Spezifikation (`LH-FA-05.e`, *Länge des `CancelRequest`*); der Slice setzt sie um.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Weitergabe von `BackendKeyData` und echtes Abbrechen am Upstream — `LH-FA-05.e` legt fest, dass ein `CancelRequest` nicht weitergeleitet wird; eine Änderung wäre eine neue Anforderung.
- Übriger Protokollrand — ist mit `slice-v1-abschluss-protokollrand` geliefert (Bestand, `done`); dieser Slice ändert nur die Länge 12 des `CancelRequest`.
- Die Längen 8 bis 11, 13 bis 15 und 269 bis zur Höchstlänge bleiben `PGR-E6001`, die Längen 16 bis 268 bleiben bei der Bibliothek — Bestand bleibt bewusst stehen: Kein Client sendet sie (die Spezifikation entscheidet sie jetzt), und eine eigene Auswertung brächte kein Verhalten, das ein Test heute verlangt.
- Das Verbindungsende einer Hauptverbindung, die ein Client bei Kontext-Abbruch mitten in einer Anfrage schließt — bleibt `PGR-E4003` (Exit-Code `4`) nach dem Abschnitt *Verbindungsende* der Spezifikation; es ist der Abbruch der Anfrage, nicht der `CancelRequest`, und eine andere Anforderung wäre ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Ein `CancelRequest` ohne Schlüssel (Länge 12) schließt die Verbindung mit `PGR-W3001` ohne Verbindungsfehler, in `record` und `replay` (Unit-Test mit 12 Byte). Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): pgx mit Kontext-Abbruch über `record` lässt den Lauf mit Exit-Code 0 enden, wenn sonst kein Fehler auftrat (Integrationstest); der Abbruch trifft eine Anfrage, die der Server vor der Frist des Clients beantwortet, sodass die Hauptverbindung regulär endet (Szenario in §3).
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
| `internal/adapters/driving/pgwire` (`startup`) | update | Nur die Länge 12 ist neu: bei Startcode `CancelRequest` und Länge 12 liest der Adapter die Prozess-ID selbst (die Bibliothek lehnt 12 ab), meldet `PGR-W3001` und schließt; fehlen die 4 Bytes, ist es `PGR-E6001`. Länge 8 bis 11, 13 bis 15 und 269 bis zur Höchstlänge bleiben bei der Bibliothek (`PGR-E6001`, gemessen), 16 bis 268 ebenso (`PGR-W3001`) |
| `internal/adapters/driving/pgwire/server_test.go` | update | Negative nach `LH-FA-05`: je Länge 8, 9, 11, 12, 13, 15, 16, 268, 269 ein Fall; Länge 12 in `record` und `replay`, mit TLS-Pflicht, nach `SSLRequest` mit `N` und über TLS nach `S`, mit abgebrochenem Rest |
| `test/integration` | update | pgx mit Kontext-Abbruch, Exit-Code: Der Kontext einer Anfrage mit `pg_sleep` (unter einer Sekunde) endet früher; `BuildContextWatcherHandler` setzt `CancelRequestContextWatcherHandler` mit einer Frist des Clients über der Dauer der Anfrage (der Standard von pgx sendet bei laufender Anfrage keinen `CancelRequest`); die Anfrage endet regulär, die Hauptverbindung sendet `Terminate`, `SIGTERM` an den Recorder, Exit-Code 0, im Log `PGR-W3001` und kein `PGR-E6001` |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |

**Zusagen und Mutationen** (`AGENTS.md` §3.10; der Bericht nennt je Zeile den roten Test): Länge 12 ist `PGR-W3001` (Mutation: Zweig für 12 entfernt, oder 12 bis 16 angenommen) · Länge 11 und 13 sind `PGR-E6001` (Mutation: Untergrenze 11 oder Bereich 12 bis 15) · 16 und 268 sind `PGR-W3001`, 269 `PGR-E6001` (Fälle pinnen die Bibliothek) · gleich in `record` und `replay` (Mutation: Zweig nur im `record`) · mit TLS-Pflicht kein `PGR-E6003` (Mutation: Prüfung der Pflicht vor dem `CancelRequest`) · nach `SSLRequest` mit `N` und über TLS nach `S` · Rest nicht geliefert ist `PGR-E6001` (Mutation: Lesefehler verworfen) · die Verbindung wird geschlossen (Mutation: sie bleibt offen) · Exit-Code 0 (Integrationstest, Mutation: Länge 12 wie bisher abgelehnt, Lauf endet mit 6).

**Größe** (`AGENTS.md` §3.14): Schätzung in der Metrik von §4, mit dem Doppelten für Tests: Produktcode rund 25, Unit-Tests rund 200, Integrationstest rund 100, Handbuch und README rund 30, zusammen **rund 350**. Zwei Liefer-Punkte, zwei Schichten (§8).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-protokollrand` ist `done`, und `welle-erster-release` ist `done` (Start-Trigger von `welle-v1-nachtrag`; Entscheidung des Nutzers vom 2026-10-11).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Behandlung verlangt die Weitergabe von `BackendKeyData`, oder die **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user README.md`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist, liegt **über 700** (Schätzung in §3: rund 350).
- `in-progress` → `open` (blockiert — Carveout?): der Adapter kann die Länge 12 nicht vor der Bibliothek lesen (der Kopf wird heute vor der Bibliothek gelesen, `startup`; die Bedingung gilt, wenn das mit `LH-FA-23.a` nicht zusammen geht).

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

**Randformen** (`AGENTS.md` §3.12), alle entschieden vor dem Code in der Spezifikation, `spec/spezifikation.md` `LH-FA-05.e`, Aufzählungspunkt *Länge des `CancelRequest`* (Nachfolgend: „Punkt“):

- Länge 12 (Prozess-ID, kein Schlüssel) und 16 bis 268 — `PGR-W3001`, Punkt.
- Länge 8 bis 11, 13 bis 15, 269 bis zur Höchstlänge — `PGR-E6001` mit `ErrorResponse` (Exit-Code 6), Punkt; unter 8 und über der Höchstlänge `PGR-W3003`, Tabelle *Protokollrand*.
- Reihenfolge der Gründe: Länge außerhalb von 8 bis zur Höchstlänge, dann Länge des `CancelRequest`, dann Rest nicht geliefert — Punkt.
- Rest der Nachricht nicht geliefert (Abbruch, Ende mitten in der Nachricht) — `PGR-E6001`, Punkt und Zeile *Startnachricht mit vollständigem Kopf*.
- Stelle: erste Nachricht und nach `SSLRequest` (`N`, oder `S` und TLS) oder `GSSENCRequest` — Punkt; mit TLS-Pflicht ohne `PGR-E6003` — Abschnitt *Unverschlüsselte Clients* von `LH-FA-23.a`.
- Modus: `record` und `replay` gleich, keine Session, keine Aufzeichnung, keine Verbindung zum Server, eine zweite Verbindung unberührt — Punkt und *Modus*.
- Prozess-ID und Schlüssel — nicht ausgewertet, Punkt.
- Exit-Code — die Warnung ändert ihn nicht (`LH-FA-13`); eine Hauptverbindung, die der Client bei Kontext-Abbruch mitten in der Anfrage schließt, ist `PGR-E4003` (Exit-Code 4) nach *Verbindungsende*, Punkt.
- Plattform (`LH-QA-03`): die Eigenschaft hängt weder vom Betriebssystem noch vom Dateisystem ab (Bytes auf einer TCP-Verbindung); Unit- und Integrationstest laufen unter Linux im Container, eine Grenze je Plattform gibt es nicht — Punkt.
- Aufruf-Voraussetzung und Aufrufarten des Treibers: nicht zutreffend (keine Ordnung von Werten über Aufrufe, kein Stand, kein Nachweis).

Risiken:

- Ein Client sendet einen `CancelRequest` in einer Länge, die die Spezifikation als `PGR-E6001` führt (zum Beispiel 13 bis 15 oder über 268) — **Ausgang:** offen bis Closure.
- Der Integrationstest trifft die Zeitfolge (Anfrage endet vor der Frist des Clients) nicht stabil — **Ausgang:** offen bis Closure; Gegenmittel: die Frist des Clients um ein Vielfaches über der Dauer der Anfrage.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — der Fall entstand, weil der Test nur den `CancelRequest` mit Schlüssel prüfte.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (PGWire-Adapter, Dokumentation). Nachgezählt beim Eintragen der Randformen (2026-10-20, Architect, `AGENTS.md` §3.13): zwei Liefer-Punkte (`LH-FA-05`, `LH-FA-13`), zwei Schichten; der Bootstrap bleibt unberührt, der Integrationstest zählt als Test des Binaries; die Spezifikation trägt der Architect vor dem Code, sie zählt nicht in die Marke.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
