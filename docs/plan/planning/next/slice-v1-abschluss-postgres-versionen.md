# Slice slice-v1-abschluss-postgres-versionen: Unterstützte PostgreSQL-Versionen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung), [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)

**Berührte Spec-Stellen:** `LH-FA-05.b` · `LH-FA-05.e` · `LH-FA-09.a` · `SPEC-029` · `SPEC-038`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Recorder sagt zu, mit welchen PostgreSQL-Hauptversionen er aufzeichnet, und das Gate belegt die Zusage für jede Version des Bereichs: Abnahmeszenarien und Lebendprüfungs-Abgleich laufen gegen jede Version grün.

**Herkunft:** Die Spezifikation sagt in `LH-FA-05.e` und `SPEC-029` nur „jeder Server, der PGWire 3.0 spricht“; das Lastenheft nennt keine Version, und die Integrationstests laufen allein gegen das gepinnte `postgres:17-alpine` (`harness/mk/integration.mk`). Seit [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) entscheidet das Replay selbst, was PostgreSQL als leere Anfrage behandelt (`LH-FA-09.a`, Lebendprüfung); das hängt am Scanner des Servers, etwa daran, ob `\v` als Leerraum zählt. Das Risiko dazu steht in §6 von `slice-extended-query-lebendpruefung`; dieser Slice ist seine Adresse. Entscheidungen des Nutzers (2026-10-05): Unterstützt werden die **fünf jüngsten erschienenen Hauptversionen**, vermutlich 15 bis 19; beim Start wird der Bereich gegen die Release-Liste von PostgreSQL bestätigt, und ist 19 nicht erschienen, gelten die fünf jüngsten erschienenen (§6). PostgreSQL 14 fällt heraus. Die Spezifikation führt den Bereich als feste Liste; ein späterer Slice hebt sie an. Das Lastenheft bekommt einen Satz in [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung). Er startet, wenn `slice-v1-abschluss-anmeldung` geliefert ist (§4).

- **Zusage:** `SPEC-029` und `LH-FA-05.e` nennen den Bereich als feste Liste von Hauptversionen und was „unterstützt“ heißt: Der Record-Modus zeichnet gegen einen Server dieser Versionen auf, und die automatisierten Abnahmeszenarien laufen gegen jede von ihnen grün. `SPEC-038` (Kompatibilitätstests) nennt die Versionsmatrix. [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung) bekommt einen Satz: Der Record-Modus setzt einen PostgreSQL-Server einer unterstützten Hauptversion voraus, und welche das sind, nennt die Spezifikation; das Lastenheft steht auf Draft, kein Versionssprung. Das Benutzerhandbuch nennt den Bereich.
- **Versionsmatrix im Gate:** `make test-integration` bleibt unverändert das volle Gate gegen die Referenzversion 17, die im Bereich liegt. Ein neues Ziel `make test-postgres-versionen` hängt an `GATE_CHECKS` und fährt gegen jede **andere** Version des Bereichs, je per Digest gepinnt, eine deklarierte Teilmenge: die Tests, deren Aussage vom Server abhängt — Abnahmeszenarien 1, 2 und 7 (`TestE2ERecordSelect1`, `TestE2EVorbereitung*` mit `TestE2EOhnePostgres*` in der zweiten Phase, `TestE2ERecordExtendedPgx`) und den Lebendprüfungs-Abgleich `TestE2EReplayLebendpruefungWiePostgres`. Das Testimage wird einmal gebaut. Begründung der Variante in §3.
- **Anmeldung je Version:** Die Matrix läuft nicht nur mit `trust`, sondern auch mit Passwort-Anmeldung: Die Teilmenge enthält den Fall SCRAM-SHA-256 des Anmeldetests aus `slice-v1-abschluss-anmeldung` (`LH-FA-05.b`, dort DoD-Punkt 1), und jeder Server des Bereichs außer der Referenz 17 nimmt für diesen Fall eine Anmeldung mit SCRAM-SHA-256 an. Gegen 17 prüft der Anmeldetest selbst in `make test-integration`. Klartext-Passwort und MD5 bleiben in der Matrix außen vor: Ihr Weg durch den Recorder hängt nicht an der Serverversion, und MD5 ist ab PostgreSQL 18 abgekündigt.
- **Lebendprüfung gegen jede Version:** Die Leerraum-Menge der Lebendprüfung hängt nach `slice-extended-query-lebendpruefung` von der Serverversion ab (`server_version` der Session; vor der Zuordnung und ohne lesbare Version ohne `\v`). Die Sonde dort ergab die Grenze 17: PostgreSQL 14, 15 und 16 lehnen `\v` außerhalb eines Kommentars mit einem Syntaxfehler ab, 17 und 18 nicht; die volle Integrationssuite lief von Hand gegen 14, 16, 17 und 18 grün. `TestE2EReplayLebendpruefungWiePostgres` liest `server_version` des Servers und prüft die Texte mit `\v` gegen die Grenze, am Server und am Replay: Unter 17 erwartet er im Replay den aufgezeichneten Syntaxfehler samt Status `E`, ab 17 die leere Antwort außerhalb der Reihe; eine falsche Produktgrenze fällt so nur gegen eine Version auf, die zwischen ihr und 17 liegt. Im Gate läuft er heute nur gegen 17. Dieser Slice prüft die versionsabhängige Menge gegen jede Version der Liste: `TestE2EReplayLebendpruefungWiePostgres` fährt die Tabelle aus `LH-FA-09.a` gegen jeden Server des Bereichs, und für jede Version stimmt die Antwort des Replays mit der des Servers überein. Weicht eine Version von der Menge ab, die die Erkennung für sie annimmt, berichtigt der Slice die Versionsgrenze in `LH-FA-09.a` und im Tabellentest; eine engste gemeinsame Menge führt er nicht ein.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Versionen vor dem Bereich (14 und älter, sofern die Bestätigung beim Start die Liste 15 bis 19 ergibt) — die Entscheidung des Nutzers zieht die Grenze bei den fünf jüngsten Hauptversionen; eine ältere Version bleibt ungeprüft und nicht zugesagt, auch wenn sie PGWire 3.0 spricht. Der Satz „jeder Server, der PGWire 3.0 spricht“ in `SPEC-029` wird dadurch ersetzt, nicht ergänzt.
- Protokollversion 3.2 (PostgreSQL 18) — bleibt Bestand: `LH-FA-05.e` lehnt jede andere Startversion als 3.0 mit `PGR-E6002` ab, und das gilt unverändert auch gegen einen Server, der 3.2 spricht. Der Slice prüft nur, dass ein Client mit 3.0 gegen jeden Server des Bereichs funktioniert; eine Aushandlung (`NegotiateProtocolVersion`) wäre ein anderer Vorgang mit eigener Anforderung.
- Der Protokollrand (Tabelle in `LH-FA-05.e`) und seine Tests — `slice-v1-abschluss-protokollrand`; dieser Slice ändert an der Tabelle nichts und nimmt ihre Tests nicht in die Matrix.
- Die volle Integrationssuite gegen jede Version im Gate — Laufzeit wächst mit der Zahl der Versionen, ohne dass Tests ohne Serverabhängigkeit (Signale, Gegendruck, Schreibfehler) etwas Neues zeigen. Ein Lauf von Hand bleibt möglich, weil `POSTGRES_IMAGE` in `integration.mk` überschreibbar ist.
- Ein Wechsel der Referenzversion des vollen Gates — bleibt 17; ein Tausch im selben Slice ersetzte eine bestehende Prüfung statt sie zu ergänzen (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`).
- Ein Sensor, der den Bereich gegen die Release-Liste von PostgreSQL hält — bräuchte Netz zur Laufzeit des Gates; die Pflege ist eine Handlung (§6), kein Lauf.
- Kein Produkt-Code außer der Erkennung der Lebendprüfung, und die nur bei einer Abweichung — Schicht-Abgrenzung: Record, Replay-Ablauf und Adapter bleiben unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung): `SPEC-029` und `LH-FA-05.e` nennen den beim Start bestätigten Versionsbereich (die fünf jüngsten erschienenen Hauptversionen) als feste Liste und was „unterstützt“ heißt; `SPEC-038` nennt die Versionsmatrix; das Benutzerhandbuch nennt den Bereich; [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung) trägt den Satz zur unterstützten Hauptversion.
- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): `make test-postgres-versionen` hängt an `make gates` und fährt die Abnahmeszenarien 1, 2 und 7 gegen jede Version des Bereichs außer der Referenz 17, je mit per Digest gepinntem Image, dazu je Version den Fall SCRAM-SHA-256 des Anmeldetests aus `slice-v1-abschluss-anmeldung`, sodass jede Version auch mit Passwort-Anmeldung statt nur mit `trust` läuft (Entscheidung des Nutzers vom 2026-10-07); `harness/README.md` §Sensors führt das Ziel mit Bindung; bewusst gebrochen: ein Image ohne Digest, eine fehlende Version oder eine Testauswahl ohne Treffer färbt das Ziel rot.
- [ ] [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay): `TestE2EReplayLebendpruefungWiePostgres` ist gegen jede Version der Liste grün: Für jede Version entspricht die versionsabhängige Leerraum-Menge der Erkennung der Antwort des Servers, auch für `\v` diesseits und jenseits der Versionsgrenze; weicht eine Version ab, sind die Versionsgrenze in `LH-FA-09.a` und der Tabellentest der Erkennung berichtigt.
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
| `spec/spezifikation.md` (`LH-FA-05.e`, `SPEC-029`, `SPEC-038`, Historie) | update | Versionsbereich als Liste, Bedeutung von „unterstützt“, Versionsmatrix der Kompatibilitätstests |
| `spec/lastenheft.md` ([`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)) | update | ein Satz zur unterstützten Hauptversion des Servers (Entscheidung des Nutzers, 2026-10-05); Draft, kein Versionssprung |
| `docs/plan/adr/` (neue ADR, Index) | new, update (Architect) | [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) nennt *ein* gepinntes PostgreSQL-Image; die Matrix und ihre Teilmenge sind eine Ergänzung mit Alternativen (unten). Ob sie eine eigene ADR braucht oder die Ausführung in `harness/mk/` genügt, entscheidet der Architect vor dem ersten Code-Commit; Empfehlung: eigene ADR, die [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) ergänzt |
| `harness/mk/postgres-versionen.mk` (neu) | new | Liste der Versionen mit Tag und Digest, Testauswahl an genau einer Stelle (einschließlich des Falls SCRAM-SHA-256 aus `slice-v1-abschluss-anmeldung`), Ziel `test-postgres-versionen`, `GATE_CHECKS +=` |
| `tools/test/run-integration-tests.sh` | update | Testauswahl je Phase als Parameter; ohne ihn verhält sich der Runner wie heute (beide Phasen, volle Suite); ein Lauf über mehrere Images baut das Testimage einmal |
| `test/integration` (`lebendpruefung_e2e_test.go` und die Tests der Teilmenge) | update (bedingt) | nur wenn ein Test eine Annahme über die Serverversion trägt (Meldungstexte, Parameter wie `server_version`); die Tabelle der Lebendprüfungen bleibt die aus `LH-FA-09.a`, die Erwartung je Zeile folgt der Version des Servers |
| `internal/hexagon/services` (`lebendpruefung.go` und Tabellentest) | update (bedingt) | nur wenn eine Version der Liste von der Versionsgrenze abweicht, die `slice-extended-query-lebendpruefung` per Sonde bestimmt hat: Grenze berichtigt |
| `harness/README.md` §Sensors, `docs/user/benutzerhandbuch.md`, `docs/user/abdeckung-*.md` | update | neues Ziel mit Bindung; unterstützte Versionen; Abdeckung, falls Deklarationen sich ändern |

**Variante der Matrix** — geprüft gegen [`AGENTS.md`](../../../../AGENTS.md) §3.6 (kein Gate wird gelockert) und die Laufzeit von `make gates`:

- *A — volle Suite gegen jede Version in `make gates`:* belegt am meisten, vervielfacht aber die Laufzeit der Integrationstests mit der Zahl der Versionen, für Tests ohne Serverabhängigkeit ohne Gewinn. Verworfen.
- *B — Matrix als eigenes Ziel außerhalb von `make gates`:* billig, aber die Zusage in `SPEC-029` stünde dann ohne Gate; ein Lauf, den niemand ruft, belegt nichts (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Verworfen.
- **C — gewählt, Laufzeit-Zuwachs vom Nutzer akzeptiert (2026-10-05):** volle Suite gegen die Referenz 17 wie heute (`make test-integration`, unverändert), dazu `make test-postgres-versionen` in `make gates` mit der serverabhängigen Teilmenge gegen jede andere Version. Mehrkosten je Version: ein PostgreSQL-Start und die Teilmenge in zwei Phasen, das Testimage einmal. Keine bestehende Prüfung fällt weg; die Teilmenge ist erweiterbar, wenn ein späterer Slice der Welle einen serverabhängigen Abnahmetest liefert.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-anmeldung` liegt in `done/` (Entscheidung des Nutzers vom 2026-10-07: die Matrix fährt den Fall SCRAM-SHA-256 seines Anmeldetests, §1), und `slice-extended-query-lebendpruefung` ist `done`. Der Vorzug vor dem Welle-Trigger (Entscheidung des Nutzers vom 2026-10-05) ist gegenstandslos, seit welle-replay-semantik done ist; Drift-Log der Roadmap. Erster Schritt nach dem Start: den Versionsbereich gegen die Release-Liste von PostgreSQL bestätigen (§6).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Eine Version des Bereichs verlangt mehr als die Leerraum-Menge der Lebendprüfung — etwa einen anderen Handshake, andere Antworten im Record oder Änderungen am Replay-Ablauf. Dann trennt sich die Anpassung an diese Version als eigener Slice ab; die Matrix ohne sie bleibt hier.
- `in-progress` → `open` (blockiert — Carveout?): Ein Image einer Version des Bereichs ist nicht per Digest zu pinnen oder läuft im internen Docker-Netz des Runners nicht; oder der Nutzer zieht die Entscheidung zum Bereich zurück.

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

- **Versionsbereich unbestätigt.** Die Liste 15 bis 19 ist der Stand des Planers. Entschieden (Nutzer, 2026-10-05): die fünf jüngsten erschienenen Hauptversionen als feste Liste; ist 19 beim Start nicht erschienen, gelten die fünf jüngsten erschienenen (dann 14 bis 18). Beim Start wird die Release-Liste von PostgreSQL gelesen, bevor die Spezifikation die Liste nennt — **Ausgang:** offen bis Closure.
- **Feste Liste altert.** Die Liste in der Spezifikation hebt sich nicht selbst an; mit jeder neuen Hauptversion fällt sie hinter die Entscheidung „fünf jüngste“ zurück, und kein Gate merkt es (ein Sensor bräuchte Netz, §1). Eine Folge-Adresse gibt es heute nicht: Ein Slice für das Anheben wäre ohne beobachtbaren Anlass geplant. Kandidat für den Ausgang: *weiter offen* → neue Beobachtung `BEO-REPO/postgres-versionsliste-altert` bei der Closure, mit dem beobachtbaren Anlass „eine neue Hauptversion von PostgreSQL ist erschienen“; tritt er ein, plant der Planer den Slice, der die Liste anhebt — **Ausgang:** offen bis Closure.
- **Satz im Lastenheft.** Entschieden (Nutzer, 2026-10-05): ein Satz in [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung), der auf die Liste der Spezifikation zeigt, ohne Versionssprung (Draft). Offen bleibt die Formulierung: Das Lastenheft nennt die Versionen nicht selbst, sonst müsste jedes Anheben das Lastenheft ändern — **Ausgang:** offen bis Closure.
- **Versionsgrenze für `\v`.** Die Erkennung der Lebendprüfung hängt nach `slice-extended-query-lebendpruefung` von `server_version` ab; die Grenze, ab der `\v` Leerraum ist, hat dort eine Sonde gegen 14 bis 18 auf 17 bestimmt; 19 ist nicht geprüft. Ob sie für jede Version der Liste stimmt und ob die übrigen Zeichen der Tabelle in `LH-FA-09.a` bei allen gleich gelten, prüft erst dieser Slice. Weicht eine Version ab, berichtigt er die Grenze; weicht sie in mehr ab als der Leerraum-Menge, greift die Rückführung aus §4 — **Ausgang:** offen bis Closure.
- **Teilmenge fällt still leer.** Die Teilmenge wählt Tests über ihren Namen; ein umbenannter Test fiele aus der Matrix, und `go test -run` mit einem Muster ohne Treffer endet grün (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`). Der Runner muss je Version und Phase prüfen, dass die erwarteten Tests liefen, nicht nur, dass keiner rot war — **Ausgang:** offen bis Closure.
- **Laufzeit des Gates.** Je Version kommt ein PostgreSQL-Start und die Teilmenge in zwei Phasen hinzu. Der Zuwachs der Variante C ist vom Nutzer akzeptiert (2026-10-05); gemessen wird `make gates` vor und nach dem Slice, die Closure-Notiz nennt beide Zahlen. Liegt der gemessene Zuwachs deutlich über der Schätzung, ist die Antwort eine kleinere Teilmenge oder ein Carveout, keine stille Herausnahme aus `GATE_CHECKS` ([`AGENTS.md`](../../../../AGENTS.md) §3.6) — **Ausgang:** offen bis Closure.
- **Anmeldung je Image.** Der Fall SCRAM-SHA-256 braucht je Server einen Benutzer mit Passwort und eine Regel in `pg_hba.conf`; wie `slice-v1-abschluss-anmeldung` das gegen 17 einrichtet, muss für jedes Image des Bereichs gleich wirken (etwa `password_encryption`, Einstiegsskripte des Images). Weicht ein Image ab, gehört die Einrichtung in die Matrix, nicht in den Anmeldetest — **Ausgang:** offen bis Closure.
- **Zweite Phase verdeckt.** Der Runner endet nach einer roten ersten Phase; eine Version, deren Record scheitert, zeigt nicht, ob Szenario 2 oder 7 gegen ihre Aufzeichnung trüge (`BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar`). Die Matrix erbt das; der Bericht je Version nennt die Phase, in der er endete — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-05); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` — 1×, `offen`; der Slice baut am Integrations-Runner, und eine Pflichtprüfung darf beim Umbau nicht entfallen. Deshalb bleibt `make test-integration` unverändert, und die Matrix kommt daneben (§1).
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` — 1×, `offen`; eine Testauswahl per Namensmuster liest sich als Zusage und endet ohne Treffer grün (§6).
- `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` — 1×, `offen`; die Matrix erbt die zwei Phasen des Runners (§6).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 3×, über der Schwelle; die Randformen dieses Slice (Lesart des Bereichs, feste Liste oder Regel, Satz im Lastenheft) sind vor dem Code entschieden (Nutzer, 2026-10-05, §6).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 5×; die Zusage in `SPEC-029` darf nicht weiter reichen als die Versionen, die das Gate fährt (Variante B verworfen, §3).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — 5×; das neue Ziel braucht eine Gegenprobe mit fehlender Version und Image ohne Digest (§2).
- `BEO-REPO/plan-folgt-korrektur-nicht` — 5×, `verkörpert` ([`AGENTS.md`](../../../../AGENTS.md) §3.9); gilt als Regel, auch für den Kopf und die Liste der Versionen.

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure ihrer Welle, nicht hier.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
