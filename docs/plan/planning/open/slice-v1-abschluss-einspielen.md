# Slice slice-v1-abschluss-einspielen: Einspielen einer Aufzeichnung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-17.a` · `LH-FA-03.b` · `LH-FA-02.b` · `SPEC-017` · `SPEC-034` · `SPEC-041` · `ARC-002` · `ARC-003` · `ARC-005` · `ARC-007`

**Verantwortlich:** —
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play` führt die Client-Anfragen einer Aufzeichnung (einfach und Extended) gegen einen PostgreSQL-Server aus, authentifiziert sich als Client, verbindet sich auf Wunsch mit TLS und verhält sich bei Serverfehlern und Abbruchsignalen wie spezifiziert.

**Übernommen aus `slice-replay-semantik-mismatch`:** `play` kennt die Option `--fail-on-unconsumed` nicht (`PGR-E2001`, Exit-Code 2) und lässt ihre Umgebungsvariable `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` unbeachtet, auch mit ungültigem Wert (`LH-FA-03.b` §Andere Kommandos, `LH-FA-17.a`); der Test dafür gehört zu den Optionen von `play`.

**Übernommen aus `slice-replay-semantik-meldungscodes`:** `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play`, mit derselben Wertemenge und Strenge wie bei `record` und `replay` (Optionstabelle, `LH-FA-17.a`).

**Bereinigung aus `slice-replay-semantik-fehlerreplay`:** Die Teile von `LH-FA-20.a` Schritt 3, die aufgezeichnete Interaktionen ohne `ReadyForQuery` behandeln, sind seit der Entscheidung in `LH-FA-02.b` §Fehlerantwort vor dem Abbruch nicht mehr erreichbar: Der Recorder schreibt keine solche Interaktion, und der Leser lehnt sie als beschädigt ab (`PGR-E3003`). Dieser Slice bereinigt sie in der Spezifikation; seine DoD sagt für sie nichts zu.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Zeitangaben und zeitgetreues Einspielen — `slice-v1-abschluss-zeitangaben`.
- Vergleich der Serverantworten mit der Aufzeichnung — `slice-v1-abschluss-antwortvergleich`.
- Paralleles Einspielen — Out-of-Scope von LH-FA-20.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): Eine Aufzeichnung mit DDL- und DML-Anweisungen, einfach und Extended, wird gegen eine leere Instanz eingespielt, und die Datenbank enthält danach deren Wirkung (Abnahmeszenario 12); mehrere Sessions laufen über eigene Verbindungen nacheinander; Authentifizierung mit Klartext, MD5 und SCRAM-SHA-256 sowie `--upstream-tls` (auch mit `--upstream-ca` für eine eigene Zertifizierungsstelle) funktionieren, eine fehlgeschlagene Anmeldung oder TLS-Pflicht ohne Option meldet `PGR-E4005` (auch bei abgelaufenem oder ungültigem Serverzertifikat), eine unlesbare CA-Datei beim Start `PGR-E2007`, `--upstream-ca` ohne TLS `PGR-E2001`, `--fail-on-unconsumed` ist bei `play` unbekannt (`PGR-E2001`) und ihre Umgebungsvariable bleibt dort unbeachtet, auch mit ungültigem Wert; `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` wirken bei `play` mit derselben Wertemenge, Strenge und Schwelle wie bei `record` und `replay`, ein ungültiger Wert ist `PGR-E2001`, auch in der Umgebungsvariable neben gültiger Option (Integrationstest).
- [ ] Eine Fehlerantwort des Servers bricht ab (`PGR-E4004`, Exit-Code 4); `--continue-on-error` läuft weiter und endet mit Exit-Code 4; `--allow-recorded-errors` lässt aufgezeichnete Fehler zu; ein Verbindungsfehler bricht immer ab (Test).
- [ ] `SIGINT` und `SIGTERM` beenden nach der laufenden Interaktion, mit `--finish-session-on-interrupt` nach der laufenden Session; der Exit-Code ist 0 ohne vorherigen Fehler, sonst 4; ohne Vergleich ist jedes Verbindungsende nach dem ersten `ReadyForQuery` `PGR-E4003`; die Rangfolge mit Vergleich (Exit-Code 5) prüft `slice-v1-abschluss-antwortvergleich` (Test).
- [ ] `make gates` grün.
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
| `internal/hexagon/services` (Play-Service), `internal/hexagon/ports/driving` | neu | Einspiel-Use-Case; nutzt nur Driven Ports |
| `internal/adapters/driving/cli` | update | Kommando `play`, Optionen, Passwort aus der Umgebung, Konfigurationsdatei; `--fail-on-unconsumed` bei `play` unbekannt, ihre Umgebungsvariable unbeachtet (Test, `LH-FA-03.b`); `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play` wie bei `record` und `replay` (übernommen aus `slice-replay-semantik-meldungscodes`, Test) |
| `internal/adapters/driven/postgres` | update | Authentifizierung und TLS als Client, Nachrichten der Gruppen senden |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-20 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`, und `slice-v1-abschluss-betrieb` ist `done` (Signalbehandlung und Konfigurationsdatei), und [ADR-0019](../../adr/0019-eigene-zertifizierungsstelle-beim-einspielen.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die Authentifizierungsverfahren des Servers sprengen den Slice — zurück zur Zerlegung.
- `in-progress` → `open`: Das Recording enthält Interaktionen, die sich nicht einspielen lassen — Carveout.


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

- Das Einspielen verändert eine Datenbank; ein Fehlgebrauch gegen eine falsche Instanz ist durch das Handbuch nur gewarnt — **Ausgang:** offen bis Closure.
- Eine Serverantwort mit `FATAL` beendet die Verbindung; die Behandlung gemäß Spezifikation ist erst im Test belegbar — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
