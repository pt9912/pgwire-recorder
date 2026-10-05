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

**Berührte Spec-Stellen:** `LH-FA-05.e` · `LH-FA-09.a` · `SPEC-029` · `SPEC-038`

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

**Ziel:** Der Recorder sagt zu, mit welchen PostgreSQL-Hauptversionen er aufzeichnet, und das Gate belegt die Zusage für jede Version des Bereichs: Abnahmeszenarien und Lebendprüfungs-Abgleich laufen gegen jede Version grün.

**Herkunft:** Die Spezifikation sagt in `LH-FA-05.e` und `SPEC-029` nur „jeder Server, der PGWire 3.0 spricht“; das Lastenheft nennt keine Version, und die Integrationstests laufen allein gegen das gepinnte `postgres:17-alpine` (`harness/mk/integration.mk`). Seit [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) entscheidet das Replay selbst, was PostgreSQL als leere Anfrage behandelt (`LH-FA-09.a`, Lebendprüfung); das hängt am Scanner des Servers, etwa daran, ob `\v` als Leerraum zählt. Das Risiko dazu steht in §6 von `slice-extended-query-lebendpruefung`; dieser Slice ist seine Adresse. Entscheidung des Nutzers (2026-10-05): unterstützt werden die PostgreSQL-Hauptversionen der letzten fünf Jahre — nach Stand des Planers 14 (2021) bis 18 (2025); beim Start gegen die Release-Liste von PostgreSQL bestätigt (§6).

- **Zusage:** `SPEC-029` und `LH-FA-05.e` nennen den Bereich als Liste von Hauptversionen (ob die Liste fest ist oder einer Regel folgt, ist offen, §6) und was „unterstützt“ heißt: Der Record-Modus zeichnet gegen einen Server dieser Versionen auf, und die automatisierten Abnahmeszenarien laufen gegen jede von ihnen grün. `SPEC-038` (Kompatibilitätstests) nennt die Versionsmatrix. Ob das Lastenheft einen Satz braucht, entscheidet der Nutzer (§6). Das Benutzerhandbuch nennt den Bereich.
- **Versionsmatrix im Gate:** `make test-integration` bleibt unverändert das volle Gate gegen die Referenzversion 17. Ein neues Ziel `make test-postgres-versionen` hängt an `GATE_CHECKS` und fährt gegen jede **andere** Version des Bereichs, je per Digest gepinnt, eine deklarierte Teilmenge: die Tests, deren Aussage vom Server abhängt — Abnahmeszenarien 1, 2 und 7 (`TestE2ERecordSelect1`, `TestE2EVorbereitung*` mit `TestE2EOhnePostgres*` in der zweiten Phase, `TestE2ERecordExtendedPgx`) und den Lebendprüfungs-Abgleich `TestE2EReplayLebendpruefungWiePostgres`. Das Testimage wird einmal gebaut. Begründung der Variante in §3.
- **Lebendprüfung gegen jede Version:** `TestE2EReplayLebendpruefungWiePostgres` prüft jede Lebendprüfung der Tabelle in `LH-FA-09.a` gegen jeden Server des Bereichs. Weicht eine Version ab (etwa `\v` kein Leerraum), schärft der Slice `LH-FA-09.a` auf die engste gemeinsame Menge, und die reine Funktion der Erkennung im Replay-Service folgt mit ihrem Tabellentest.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Versionen vor dem Bereich (13 und älter) — die Entscheidung des Nutzers zieht die Grenze bei fünf Jahren; eine ältere Version bleibt ungeprüft und nicht zugesagt, auch wenn sie PGWire 3.0 spricht. Der Satz „jeder Server, der PGWire 3.0 spricht“ in `SPEC-029` wird dadurch ersetzt, nicht ergänzt.
- Protokollversion 3.2 (PostgreSQL 18) — bleibt Bestand: `LH-FA-05.e` lehnt jede andere Startversion als 3.0 mit `PGR-E6002` ab, und das gilt unverändert auch gegen einen Server, der 3.2 spricht. Der Slice prüft nur, dass ein Client mit 3.0 gegen jeden Server des Bereichs funktioniert; eine Aushandlung (`NegotiateProtocolVersion`) wäre ein anderer Vorgang mit eigener Anforderung.
- Der Protokollrand (Tabelle in `LH-FA-05.e`) und seine Tests — `slice-v1-abschluss-protokollrand`; dieser Slice ändert an der Tabelle nichts und nimmt ihre Tests nicht in die Matrix.
- Anmeldeverfahren je Version (Passwort, SCRAM gegen ein Standard-Image) — `slice-v1-abschluss-anmeldung`; die Matrix läuft wie die Integrationstests mit `POSTGRES_HOST_AUTH_METHOD=trust`.
- Die volle Integrationssuite gegen jede Version im Gate — Laufzeit wächst mit der Zahl der Versionen, ohne dass Tests ohne Serverabhängigkeit (Signale, Gegendruck, Schreibfehler) etwas Neues zeigen. Ein Lauf von Hand bleibt möglich, weil `POSTGRES_IMAGE` in `integration.mk` überschreibbar ist.
- Ein Wechsel der Referenzversion des vollen Gates — bleibt 17; ein Tausch im selben Slice ersetzte eine bestehende Prüfung statt sie zu ergänzen (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`).
- Ein Sensor, der den Bereich gegen die Release-Liste von PostgreSQL hält — bräuchte Netz zur Laufzeit des Gates; die Pflege ist eine Handlung (§6), kein Lauf.
- Kein Produkt-Code außer der Erkennung der Lebendprüfung, und die nur bei einer Abweichung — Schicht-Abgrenzung: Record, Replay-Ablauf und Adapter bleiben unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): `SPEC-029` und `LH-FA-05.e` nennen den bestätigten Versionsbereich als Liste und was „unterstützt“ heißt; `SPEC-038` nennt die Versionsmatrix; das Benutzerhandbuch nennt den Bereich; das Lastenheft trägt den Satz, falls der Nutzer ihn bestätigt (§6).
- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): `make test-postgres-versionen` hängt an `make gates` und fährt die Abnahmeszenarien 1, 2 und 7 gegen jede Version des Bereichs außer der Referenz 17, je mit per Digest gepinntem Image; `harness/README.md` §Sensors führt das Ziel mit Bindung; bewusst gebrochen: ein Image ohne Digest, eine fehlende Version oder eine Testauswahl ohne Treffer färbt das Ziel rot.
- [ ] [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay): `TestE2EReplayLebendpruefungWiePostgres` ist gegen jede Version des Bereichs grün; weicht eine ab, nennt `LH-FA-09.a` die engste gemeinsame Menge und der Tabellentest der Erkennung folgt ihr.
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
| `spec/lastenheft.md` | update (bedingt) | nur nach Bestätigung des Nutzers (§6); Ort wäre die Einsatzumgebung ([`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)) oder die Beschreibung von [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol); Draft, kein Versionssprung |
| `docs/plan/adr/` (neue ADR, Index) | new, update (Architect) | [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) nennt *ein* gepinntes PostgreSQL-Image; die Matrix und ihre Teilmenge sind eine Ergänzung mit Alternativen (unten). Ob sie eine eigene ADR braucht oder die Ausführung in `harness/mk/` genügt, entscheidet der Architect vor dem ersten Code-Commit; Empfehlung: eigene ADR, die [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) ergänzt |
| `harness/mk/postgres-versionen.mk` (neu) | new | Liste der Versionen mit Tag und Digest, Testauswahl an genau einer Stelle, Ziel `test-postgres-versionen`, `GATE_CHECKS +=` |
| `tools/test/run-integration-tests.sh` | update | Testauswahl je Phase als Parameter; ohne ihn verhält sich der Runner wie heute (beide Phasen, volle Suite); ein Lauf über mehrere Images baut das Testimage einmal |
| `test/integration` (`lebendpruefung_e2e_test.go` und die Tests der Teilmenge) | update (bedingt) | nur wenn ein Test eine Annahme über die Serverversion trägt (Meldungstexte, Parameter wie `server_version`); die Tabelle der Lebendprüfungen bleibt die aus `LH-FA-09.a` |
| `internal/hexagon/services` (`lebendpruefung.go` und Tabellentest) | update (bedingt) | nur bei einer Abweichung einer Version: Leerraum-Menge auf die engste gemeinsame Menge |
| `harness/README.md` §Sensors, `docs/user/benutzerhandbuch.md`, `docs/user/abdeckung-*.md` | update | neues Ziel mit Bindung; unterstützte Versionen; Abdeckung, falls Deklarationen sich ändern |

**Variante der Matrix** — geprüft gegen [`AGENTS.md`](../../../../AGENTS.md) §3.6 (kein Gate wird gelockert) und die Laufzeit von `make gates`:

- *A — volle Suite gegen jede Version in `make gates`:* belegt am meisten, vervielfacht aber die Laufzeit der Integrationstests mit der Zahl der Versionen, für Tests ohne Serverabhängigkeit ohne Gewinn. Verworfen.
- *B — Matrix als eigenes Ziel außerhalb von `make gates`:* billig, aber die Zusage in `SPEC-029` stünde dann ohne Gate; ein Lauf, den niemand ruft, belegt nichts (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Verworfen.
- **C — gewählt:** volle Suite gegen die Referenz 17 wie heute (`make test-integration`, unverändert), dazu `make test-postgres-versionen` in `make gates` mit der serverabhängigen Teilmenge gegen jede andere Version. Mehrkosten je Version: ein PostgreSQL-Start und die Teilmenge in zwei Phasen, das Testimage einmal. Keine bestehende Prüfung fällt weg; die Teilmenge ist erweiterbar, wenn ein späterer Slice der Welle einen serverabhängigen Abnahmetest liefert.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-extended-query-lebendpruefung` ist `done`. Erster Schritt nach dem Start: den Versionsbereich gegen die Release-Liste von PostgreSQL bestätigen (§6).

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

- **Versionsbereich unbestätigt.** Der Bereich 14 bis 18 ist der Stand des Planers, nicht geprüft. Zwei Lesarten von „letzte fünf Jahre“ geben verschiedene Listen: PostgreSQL 14 erschien nach Kenntnis des Planers Ende September 2021, also knapp mehr als fünf Jahre vor dem 2026-10-05; PostgreSQL 19 kann im September 2026 erschienen sein. „Die fünf jüngsten Hauptversionen“ und „erschienen in den letzten fünf Jahren“ ergeben dann 15 bis 19 oder 14 bis 18. Beim Start wird der Bereich gegen die Release-Liste von PostgreSQL bestätigt und die Lesart vom Nutzer entschieden, bevor die Spezifikation ihn nennt. Ebenfalls offen: ob die Zusage eine feste Liste ist, die ein späterer Slice anhebt, oder eine Regel, die mit der Zeit wandert — eine wandernde Regel ändert die Zusage ohne Änderung am Text, und das Gate prüft eine feste Liste — **Ausgang:** offen bis Closure.
- **Satz im Lastenheft.** Das Lastenheft nennt keine Version; ohne Satz trägt allein die Spezifikation die Zusage, und eine Version außerhalb des Bereichs wäre nach dem Lastenheft weiter „ein Server“. Ob [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung) oder [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) einen Satz bekommt, bestätigt der Nutzer vor dem ersten Code-Commit; das Lastenheft steht auf Draft, kein Versionssprung. Empfehlung des Planers: ein Satz in `LH-RB-02`, weil die Serverversion eine Randbedingung der Einsatzumgebung ist und keine Funktion — **Ausgang:** offen bis Closure.
- **Scanner älterer Versionen.** Ob PostgreSQL 14 bis 16 `\v` und die übrigen Zeichen der Tabelle in `LH-FA-09.a` als Leerraum lesen, ist nicht geprüft. Weicht eine Version ab, verengt der Slice die Definition auf die gemeinsame Menge; eine Anfrage, die dann keine Lebendprüfung mehr ist, wird wie jede andere streng verglichen, und das Replay einer gegen 17 aufgezeichneten Lebendprüfung mit `\v` verliert die Toleranz außerhalb der Reihe. Weicht eine Version in mehr ab als der Leerraum-Menge, greift die Rückführung aus §4 — **Ausgang:** offen bis Closure.
- **Teilmenge fällt still leer.** Die Teilmenge wählt Tests über ihren Namen; ein umbenannter Test fiele aus der Matrix, und `go test -run` mit einem Muster ohne Treffer endet grün (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`). Der Runner muss je Version und Phase prüfen, dass die erwarteten Tests liefen, nicht nur, dass keiner rot war — **Ausgang:** offen bis Closure.
- **Laufzeit des Gates.** Je Version kommt ein PostgreSQL-Start und die Teilmenge in zwei Phasen hinzu. Gemessen wird `make gates` vor und nach dem Slice; die Closure-Notiz nennt beide Zahlen. Ist der Zuwachs dem Nutzer zu hoch, ist die Antwort eine kleinere Teilmenge oder ein Carveout, keine stille Herausnahme aus `GATE_CHECKS` ([`AGENTS.md`](../../../../AGENTS.md) §3.6) — **Ausgang:** offen bis Closure.
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
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 3×, über der Schwelle; die Randformen dieses Slice (Lesart des Bereichs, feste Liste oder Regel, Satz im Lastenheft) stehen in §6 und werden vor dem Code entschieden.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 5×; die Zusage in `SPEC-029` darf nicht weiter reichen als die Versionen, die das Gate fährt (Variante B verworfen, §3).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — 5×; das neue Ziel braucht eine Gegenprobe mit fehlender Version und Image ohne Digest (§2).
- `BEO-REPO/plan-folgt-korrektur-nicht` — 5×, `verkörpert` ([`AGENTS.md`](../../../../AGENTS.md) §3.9); gilt als Regel, auch für den Kopf und die Liste der Versionen.

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure ihrer Welle, nicht hier.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
