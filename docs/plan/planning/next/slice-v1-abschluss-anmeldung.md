# Slice slice-v1-abschluss-anmeldung: Anmeldung des Clients im Record-Modus vermitteln

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md)

**Berührte Spec-Stellen:** `LH-FA-05.b` · `LH-FA-02.a` · `SPEC-034` · `ARC-006` · `ARC-007` · `architecture.md §4.4`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record` vermittelt die Anmeldung des Clients beim Upstream transparent (Klartext-Passwort, MD5, SCRAM-SHA-256), sodass Anwendungen gegen Server mit Passwort-Anmeldung aufgezeichnet werden können.

**Übernommen aus `slice-harness-blackbox-driven`** (Review F-444, Verifikation V-85; Abgrenzung in `slice-tests-ueberlebende-mutanten`): ein Test, dass die Session aus `Open` den Lesepuffer des Aufbaus trägt, mit Test-Idee und Grenze in §6; er gehört zu DoD-Punkt 1, weil der Umbau der Aufbau-Schleife in `Open` genau diese Stelle berührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Anmeldung beim Einspielen (`play`) — `slice-v1-abschluss-einspielen`.
- Client-Zertifikate und GSSAPI/SSPI — außerhalb des Funktionsumfangs von v1 (`LH-FA-05.b`).
- Anmeldung gegen jede unterstützte PostgreSQL-Version — `slice-v1-abschluss-postgres-versionen` (dort §1 *Anmeldung je Version* und DoD-Punkt 2; Entscheidung des Nutzers vom 2026-10-07): Dessen Matrix fährt den Fall SCRAM-SHA-256 des Anmeldetests aus diesem Slice gegen jede Version außer 17 und startet erst, wenn dieser Slice in `done/` liegt; hier läuft der Anmeldetest gegen die Referenzversion 17. Der Test muss deshalb als Fall SCRAM-SHA-256 einzeln wählbar sein.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Ein Client meldet sich über `record` an einem Server mit Klartext-Passwort, MD5 und SCRAM-SHA-256 an und führt eine Anfrage aus; Passwort und Anmeldenachrichten stehen nicht in der Aufzeichnung (Integrationstest).
- [ ] Eine fehlgeschlagene Anmeldung geht als Fehlerantwort des Servers unverändert an den Client; ein nicht vermitteltes Verfahren endet mit eindeutigem Meldungscode (Test).
- [ ] Der Port zwischen PGWire-Adapter, Record-Service und Upstream-Adapter trägt den Anmeldeaustausch, ohne dass der Core PGWire-Typen kennt (`make a-check`).
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
| `internal/hexagon/ports/driving`, `…/driven`, `internal/hexagon/services` | update | Anmeldeaustausch als Folge fachlicher Nachrichten zwischen Client und Upstream |
| `internal/adapters/driving/pgwire`, `internal/adapters/driven/postgres` | update | Weiterleitung der Anmeldenachrichten, `SetAuthType` für SASL |
| `test/integration` | update | Server mit `password`, `md5` und `scram-sha-256` in `pg_hba.conf` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-walking-skeleton-record` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Der Anmeldeaustausch verlangt einen anderen Zuschnitt des Record-Ports als eine Erweiterung — zurück zur Zerlegung.
- `in-progress` → `open`: SCRAM mit Kanalbindung lässt sich ohne TLS zum Upstream nicht vermitteln — Carveout.

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

- Die Vermittlung muss die Anmeldenachrichten aus der Aufzeichnung heraushalten, sonst stehen Passwort-Hashes darin (LH-RB-01) — **Ausgang:** offen bis Closure.
- Dass die Session aus `Open` den Lesepuffer des Aufbaus trägt, prüft kein Test (`BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft`, seit slice-harness-blackbox-driven); der Umbau der Aufbau-Schleife berührt genau diese Stelle. Test-Idee: Der Server schickt AuthenticationOk, ReadyForQuery und eine NoticeResponse in einem Flush, der Test ruft nach `Open` einmal `Receive` und erwartet die Notice; rot, wenn `Open` eine Session mit neuem Frontend zurückgibt. Grenze: Dass beide Nachrichten in einem Lesevorgang ankommen, ist über TCP-Loopback stabil, aber nicht zugesagt; `net.Pipe` erreicht `Open` nicht, weil `Upstream.Dialer` ein `net.Dialer` ist — **Ausgang:** offen bis Closure.

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
