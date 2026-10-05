# Slice slice-extended-query-lebendpruefung: Lebendprüfungen im Replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-extended-query.

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung), [ADR-0007](../../adr/0007-strict-replay.md)

**Berührte Spec-Stellen:** `LH-FA-09.a` · `LH-FA-18.a` · `LH-FA-06.a` · `SPEC-011` · `ARC-002` · `ARC-006`

**Verantwortlich:** —
**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Eine Anwendung mit `pgxpool` oder `database/sql` läuft im Replay unabhängig von ihren Pausen: Reine Kommentar- oder Leer-Anfragen, wie sie Pools nach Leerlauf als Lebendprüfung senden (`-- ping`), beantwortet das Replay außerhalb der Reihe, ohne den Cursor zu bewegen.

**Herkunft:** Validierung von `slice-extended-query-replay` (`docs/reviews/2026-10-05-validierung-slice-extended-query-replay.md`, Frage 2, Sonden S4 bis S7). Ob die Lebendprüfung in der Aufzeichnung oder in der Wiedergabe steht, hängt an Pausen über 1 s; das strenge Replay meldet dann `PGR-E5001`. Der Nutzer hat entschieden, solche Anfragen zu tolerieren, statt die Zusage „nur Host und Port“ aus [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) zu verengen.

- Eine einfache Anfrage, deren Text nur aus Leerraum und Kommentaren besteht, beantwortet das Replay zwischen zwei Interaktionen mit `EmptyQueryResponse` und `ReadyForQuery`, wie PostgreSQL es tut; der Cursor bleibt stehen.
- Eine aufgezeichnete Lebendprüfung, die in der Wiedergabe ausbleibt, führt nicht zur Abweichung. Ob der Record sie gar nicht erst aufzeichnet oder das Replay sie am Cursor überspringt, entscheidet die ADR; vorhandene Aufzeichnungen mit Lebendprüfungen bleiben lesbar.
- Jede andere einfache Anfrage außer der Reihe bleibt eine Abweichung, auch eine, die mitten in einer Extended-Interaktion kommt (`PGR-E5001`). Diesen Fall nennt `LH-FA-18.a` §Replay heute nicht ausdrücklich (Review F-327 zu `slice-extended-query-replay`); der Satz, der die Ausnahme für Lebendprüfungen schreibt, grenzt ihn mit ab.
- **ADR:** Die Lockerung widerspricht [ADR-0007](../../adr/0007-strict-replay.md) und `LH-FA-09.a` („Kommentare werden nicht ignoriert“). Der Architect schreibt in diesem Slice eine neue ADR, die [ADR-0007](../../adr/0007-strict-replay.md) für Lebendprüfungen einschränkt (ob als Ergänzung oder mit `Supersedes`, entscheidet er), angenommen vor dem ersten Code-Commit. Die Spezifikation (`LH-FA-09.a`, `LH-FA-18.a` §Replay, bei Bedarf `LH-FA-06.a`) folgt ihr im selben Slice.
- **Handbuch:** Der Abschnitt „Mit einem Datenbanktreiber arbeiten“ nennt das Verhalten bei Lebendprüfungen von Pools und `database/sql` (`ShouldPing`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Toleranz für andere Anfragen (Normalisierung, Kommentare in fachlichen Anfragen, Reihenfolge) — [ADR-0007](../../adr/0007-strict-replay.md) bleibt für alles außer Lebendprüfungen bestehen; Out-of-Scope der Welle.
- Gleichzeitig genutzte Verbindungen eines Pools (Validierung, Sonde S8) — `slice-v1-abschluss-sessions` §6 führt sie.
- Frist beim Herunterfahren im Replay — `slice-v1-abschluss-betrieb`.
- Passwort-Anmeldung beim Aufzeichnen gegen ein Standard-Image (Validierung, Sonde S12) — `slice-v1-abschluss-anmeldung`.
- Lebendprüfungen über das Extended Query Protocol — kein beobachteter Treiber sendet sie; pgx pingt mit einer einfachen Anfrage (Validierung, pgx-Quelltext). Ein solcher Fall wäre ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay): Eine reine Kommentar- oder Leer-Anfrage zwischen zwei Interaktionen beantwortet das Replay mit `EmptyQueryResponse` und `ReadyForQuery`, ohne den Cursor zu bewegen; eine aufgezeichnete, die ausbleibt, ist keine Abweichung; jede andere einfache Anfrage außer der Reihe, auch mitten in einer Extended-Interaktion, bleibt `PGR-E5001` (Unit-Tests).
- [ ] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung): `pgxpool` und `database/sql` mit pgx in der Default-Konfiguration laufen mit Pausen über 1 s beim Aufzeichnen, beim Wiedergeben oder bei beiden ohne Abweichung durch (Integrationstest nach den Lagen S4 bis S7 der Validierung).
- [ ] Die neue ADR zu Lebendprüfungen ist angenommen und im ADR-Index; `LH-FA-09.a` und `LH-FA-18.a` §Replay nennen die Ausnahme und die falsche Protokollart; das Benutzerhandbuch nennt das Verhalten bei Lebendprüfungen.
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
| `docs/plan/adr/` (neue ADR, Index) | new, update | Einschränkung von [ADR-0007](../../adr/0007-strict-replay.md) für Lebendprüfungen; Architect |
| `spec/spezifikation.md` (`LH-FA-09.a`, `LH-FA-18.a` §Replay, bei Bedarf `LH-FA-06.a`) | update | Ausnahme für Lebendprüfungen; falsche Protokollart im Replay (F-327) |
| `internal/hexagon/services` (Replay-Service, Matcher; bei Bedarf Record-Service) | update | Lebendprüfung erkennen und außerhalb der Reihe beantworten; aufgezeichnete Lebendprüfung nach der ADR behandeln |
| `internal/adapters/driving/pgwire` | update | nur falls die Antwort außerhalb der Reihe eine Übersetzung braucht, die der Adapter heute nicht hat |
| `test/integration` | new | `pgxpool` und `database/sql` mit und ohne Pausen (Lagen S4 bis S7) |
| `docs/user/benutzerhandbuch.md`, `docs/user/abdeckung-*.md` | update | Verhalten bei Lebendprüfungen; Abdeckungstabellen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-extended-query-replay` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die ADR verlangt eine Änderung am Recording-Format (etwa eine Markierung aufgezeichneter Lebendprüfungen) — das ist eine dritte Schicht; Format und Replay werden getrennt geschnitten.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer verwirft die Einschränkung von [ADR-0007](../../adr/0007-strict-replay.md) im Review der ADR — dann gilt der zweite Weg der Validierung (Zusage verengen, Änderung am Lastenheft), und dieser Slice geht mit Grund zurück.

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

- Transaktionsstatus im `ReadyForQuery` der Antwort außerhalb der Reihe: Er muss dem Stand nach der letzten beantworteten Interaktion entsprechen, auch innerhalb einer Transaktion oder nach einem Fehler in ihr — **Ausgang:** offen bis Closure.
- Eine fachliche Anfrage, die nur aus Kommentar besteht, ist von einer Lebendprüfung nicht zu unterscheiden. PostgreSQL antwortet auf beide gleich (`EmptyQueryResponse`); ob die Ausnahme damit für jede solche Anfrage gelten darf, entscheidet die ADR — **Ausgang:** offen bis Closure.
- Vorhandene Aufzeichnungen enthalten Lebendprüfungen als Interaktionen (Validierung, Sonde S6). Bleibt die Wiedergabe einer solchen Aufzeichnung verträglich, wenn der Record sie künftig nicht mehr aufzeichnet? — **Ausgang:** offen bis Closure.
- Ob [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) im Lastenheft einen Satz braucht oder die Spezifikation genügt, bestätigt der Nutzer vor dem ersten Code-Commit — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand der Closure von `slice-extended-query-replay`, 2026-10-05); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` — 1×, `offen`; der Gegenstand dieses Slice ist sein erster Fund.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 3×, über der Schwelle; die Randformen in §6 gehören vor den Code in die ADR, nicht in eine Prüfrunde.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — je 5×; die Ausnahme ist ein neuer Vertrag am Replay, ihre Grenze (welche Anfrage, welcher Cursor) braucht Tests, die eine Mutation fangen.
- `BEO-REPO/plan-folgt-korrektur-nicht` — 5×, `verkörpert` (`AGENTS.md` §3.9); gilt als Regel, auch für den Kopf.

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure von `welle-extended-query`, nicht hier.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
