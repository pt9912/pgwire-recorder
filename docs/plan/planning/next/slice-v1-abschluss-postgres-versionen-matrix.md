# Slice slice-v1-abschluss-postgres-versionen-matrix: Versionsmatrix im Gate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md), [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md)

**Berührte Spec-Stellen:** `LH-FA-05.b` · `LH-FA-05.e` · `LH-FA-09.a` · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) · [`SPEC-053`](../../../../harness/vertraege/SPEC-053-postgres-versionen.md)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Gate belegt für jede Version des Bereichs, dass der Recorder gegen sie aufzeichnet: Die Abnahmeszenarien 1, 2 und 7, der Fall SCRAM-SHA-256 der Anmeldung und der Lebendprüfungs-Abgleich laufen gegen jede Version grün.

**Übernimmt:** `slice-v1-abschluss-postgres-versionen` — dessen Teil *Versionsmatrix* (Entscheidung des Nutzers vom 2026-10-14, Schnitt nach Lieferwert, Matrix zuerst; dort §7 `Gegenstand:`). Der andere Teil, die Zusage in Spezifikation, Lastenheft, Handbuch und README, liegt in `slice-v1-abschluss-postgres-versionen-zusage`.

**Herkunft:** Die Spezifikation sagt in `LH-FA-05.e` nur „jeder Server, der PGWire 3.0 spricht“, und die Integrationstests laufen allein gegen das gepinnte `postgres:17-alpine` (`harness/mk/integration.mk`). Seit [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) entscheidet das Replay selbst, was PostgreSQL als leere Anfrage behandelt (`LH-FA-09.a`, Lebendprüfung); das hängt am Scanner des Servers, etwa daran, ob `\v` als Leerraum zählt. Das Risiko dazu steht in §6 von `slice-extended-query-lebendpruefung`; dieser Slice ist seine Adresse (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`). Entscheidungen des Nutzers (2026-10-05): Unterstützt werden die **fünf jüngsten erschienenen Hauptversionen**; der Bereich wird gegen die Release-Liste von PostgreSQL bestätigt. Bestätigt am 2026-10-14 (Architect, `postgresql.org/support/versioning`, Docker Hub): 19 ist nicht erschienen (Beta 4 vom 2026-09-24, kein Image `19-alpine`), es gelten **14 bis 18** (Referenz 17; die Matrix fährt 14, 15, 16 und 18). PostgreSQL 14 endet am 2026-11-12; die Liste verschiebt sich, sobald 19 erscheint, und das ist eine Pflege der Liste von Hand ([`SPEC-053`](../../../../harness/vertraege/SPEC-053-postgres-versionen.md) Punkt 14). Die Matrix hat zwei Gründe, sie vor die Zusage zu stellen: Eine Zusage sagt nur, was das Gate belegt (`AGENTS.md` §3.11), und die bestätigte Liste, die das Gate fährt, ist die Liste, die die Spezifikation später nennt.

- **Versionsmatrix im Gate:** `make test-integration` bleibt unverändert das volle Gate gegen die Referenzversion 17, die im Bereich liegt. Ein neues Ziel `make test-postgres-versionen` hängt an `GATE_CHECKS` und fährt gegen jede **andere** Version des Bereichs, je per Digest gepinnt, eine deklarierte Teilmenge: die Tests, deren Aussage vom Server abhängt — Abnahmeszenarien 1, 2 und 7 (`TestE2ERecordSelect1`, `TestE2EVorbereitung*` mit `TestE2EOhnePostgres*` in der zweiten Phase, `TestE2ERecordExtendedPgx`) und den Lebendprüfungs-Abgleich `TestE2EReplayLebendpruefungWiePostgres`. Das Testimage wird einmal gebaut. Begründung der Variante in §3. Die Liste der Versionen mit Tag und Digest steht an genau einer Stelle (`harness/mk/postgres-versionen.mk`); sie ist die Quelle, die `slice-v1-abschluss-postgres-versionen-zusage` in der Spezifikation nennt.
- **Anmeldung je Version:** Die Matrix läuft nicht nur mit `trust`, sondern auch mit Passwort-Anmeldung: Die Teilmenge enthält den Fall SCRAM-SHA-256 des Anmeldetests aus `slice-v1-abschluss-anmeldung` (`LH-FA-05.b`, dort DoD-Punkt 1), und jeder Server des Bereichs außer der Referenz 17 nimmt für diesen Fall eine Anmeldung mit SCRAM-SHA-256 an. Gegen 17 prüft der Anmeldetest selbst in `make test-integration`. Klartext-Passwort und MD5 bleiben in der Matrix außen vor: Ihr Weg durch den Recorder hängt nicht an der Serverversion, und MD5 ist ab PostgreSQL 18 abgekündigt.
- **Lebendprüfung gegen jede Version:** Die Leerraum-Menge der Lebendprüfung hängt nach `slice-extended-query-lebendpruefung` von der Serverversion ab (`server_version` der Session; vor der Zuordnung und ohne lesbare Version ohne `\v`). Die Sonde dort ergab die Grenze 17: PostgreSQL 14, 15 und 16 lehnen `\v` außerhalb eines Kommentars mit einem Syntaxfehler ab, 17 und 18 nicht; die volle Integrationssuite lief von Hand gegen 14, 16, 17 und 18 grün. `TestE2EReplayLebendpruefungWiePostgres` liest `server_version` des Servers und prüft die Texte mit `\v` gegen die Grenze, am Server und am Replay: Unter 17 erwartet er im Replay den aufgezeichneten Syntaxfehler samt Status `E`, ab 17 die leere Antwort außerhalb der Reihe; eine falsche Produktgrenze fällt so nur gegen eine Version auf, die zwischen ihr und 17 liegt. Im Gate läuft er heute nur gegen 17. Dieser Slice prüft die versionsabhängige Menge gegen jede Version der Liste: `TestE2EReplayLebendpruefungWiePostgres` fährt die Tabelle aus `LH-FA-09.a` gegen jeden Server des Bereichs, und für jede Version stimmt die Antwort des Replays mit der des Servers überein. Weicht eine Version von der Menge ab, die die Erkennung für sie annimmt, ist das ein Fund mit Rückführung (§4), keine stille Berichtigung im Lauf; eine engste gemeinsame Menge führt der Slice nicht ein.

**Nahtstelle für `slice-v1-abschluss-postgres-versionen-zusage`** (Eintrag in dessen §1 und §6 durch den Planner, vor dem Übergang dieses Slice nach `done/`, mit Nachzählung im selben Commit; Geprüft mit `grep -n "slice-v1-abschluss-postgres-versionen-matrix"` im Nehmer): (1) die bestätigte Liste der Hauptversionen und die Referenzversion, wie sie `harness/mk/postgres-versionen.mk` führt; (2) die Teilmenge der Matrix im Wortlaut (Szenarien 1, 2, 7, Fall SCRAM-SHA-256, Lebendprüfungs-Abgleich), denn „unterstützt“ in der Spezifikation heißt genau das; (3) jede Berichtigung, die die Matrix an der Versionsgrenze in `LH-FA-09.a` brauchte, samt der Version, bei der sie auffiel; (4) der Aufruf, mit dem ein Verifier den Lauf je Version nachfährt (`POSTGRES_IMAGE` und die Testauswahl), weil das Handbuch-Sätze über das Gegenüber mit einem nachfahrbaren Lauf in §7 belegt (`.claude/commands/implement-slice.md` Schritt 19). Die Namen, die der Implementer für Datei, Variablen und Aufruf wählt, trägt der Planner bei der Closure in den Nehmer ein (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Zusage in Spezifikation und Lastenheft (Liste der Versionen, Bedeutung von „unterstützt“, Versionsmatrix der Kompatibilitätstests, Satz zur Einsatzumgebung), Handbuch und README — `slice-v1-abschluss-postgres-versionen-zusage` (dort §1 und DoD; Entscheidung des Nutzers vom 2026-10-14): Eine Zusage sagt nur, was das Gate belegt (`AGENTS.md` §3.11); sie kommt nach der Matrix. Dieser Slice liefert kein Verhalten des Binaries, also keinen Handbuchteil (`AGENTS.md` §3.11).
- Die Handbuch-Sätze „geprüft mit PostgreSQL 17“ zur Anmeldung — `slice-v1-abschluss-postgres-versionen-zusage` (dort §1; `AGENTS.md` §3.13): Sie gehören zum Handbuch, und welche von ihnen auf den Bereich reichen, hängt an dem, was diese Matrix fährt (nur der Fall SCRAM-SHA-256).
- Versionen vor dem Bereich (13 und älter nach der Bestätigung vom 2026-10-14; erscheint 19 vor dem Start, bestätigt der Implementer die Liste neu und der Bereich ist dann 15 bis 19) — die Entscheidung des Nutzers zieht die Grenze bei den fünf jüngsten Hauptversionen; eine ältere Version bleibt ungeprüft und nicht zugesagt, auch wenn sie PGWire 3.0 spricht.
- Protokollversion 3.2 (PostgreSQL 18) — bleibt Bestand: `LH-FA-05.e` lehnt jede andere Startversion als 3.0 mit `PGR-E6002` ab, und das gilt unverändert auch gegen einen Server, der 3.2 spricht. Der Slice prüft nur, dass ein Client mit 3.0 gegen jeden Server des Bereichs funktioniert; eine Aushandlung (`NegotiateProtocolVersion`) wäre ein anderer Vorgang mit eigener Anforderung.
- Der Protokollrand (Tabelle in `LH-FA-05.e`) und seine Tests — `slice-v1-abschluss-protokollrand`; dieser Slice ändert an der Tabelle nichts und nimmt ihre Tests nicht in die Matrix.
- Die volle Integrationssuite gegen jede Version im Gate — Laufzeit wächst mit der Zahl der Versionen, ohne dass Tests ohne Serverabhängigkeit (Signale, Gegendruck, Schreibfehler) etwas Neues zeigen. Ein Lauf von Hand bleibt möglich, weil `POSTGRES_IMAGE` in `integration.mk` überschreibbar ist.
- Abnahmeszenario 18 und Szenario 16 Satz 3 (Verbindungsende durch den Server, `slice-v1-abschluss-verbindungsende-upstream` und `-vergleich`) in der Matrix — Entscheidung des Nutzers vom 2026-10-12: Beide laufen nur gegen PostgreSQL 17, die Referenzversion des vollen Gates. Begründung: Sie stützen sich auf `pg_terminate_backend`, `57P01` und `idle_session_timeout`, deren Verhalten je Hauptversion eigene Prüfung bräuchte; der Zuwachs an Gate-Laufzeit steht dafür nicht im Verhältnis. Die Lücke bleibt benannt (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`).
- Ein Wechsel der Referenzversion des vollen Gates — bleibt 17; ein Tausch im selben Slice ersetzte eine bestehende Prüfung statt sie zu ergänzen (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`).
- Ein Sensor, der den Bereich gegen die Release-Liste von PostgreSQL hält — bräuchte Netz zur Laufzeit des Gates; die Pflege ist eine Handlung (§6), kein Lauf.
- Kein Produkt-Code außer der Erkennung der Lebendprüfung, und die nur nach einer Rückführung (§4) — Schicht-Abgrenzung: Record, Replay-Ablauf und Adapter bleiben unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Der Versionsbereich ist gegen die Release-Liste von PostgreSQL bestätigt (14 bis 18, Stand 2026-10-14, §6) und steht als Liste in `harness/mk/postgres-versionen.mk`; `make test-postgres-versionen` hängt an `make gates` und fährt die Abnahmeszenarien 1, 2 und 7 gegen jede Version des Bereichs außer der Referenz 17, je mit per Digest gepinntem Image, dazu je Version den Fall SCRAM-SHA-256 des Anmeldetests aus `slice-v1-abschluss-anmeldung`, sodass jede Version auch mit Passwort-Anmeldung statt nur mit `trust` läuft (Entscheidung des Nutzers vom 2026-10-07); `harness/README.md` §Sensors führt das Ziel mit Bindung, die Ausnahmeliste der schnellen Kette (`SPEC-052`) nennt es (entschieden, `SPEC-053` Punkt 13) und `gates-schnell-gegenprobe` kennt es, und eine Gegenprobe des Ziels hängt an `GATE_CHECKS`; bewusst gebrochen: ein Image ohne Digest, eine fehlende Version oder eine Testauswahl ohne Treffer färbt das Ziel rot.
- [ ] [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay): `TestE2EReplayLebendpruefungWiePostgres` ist gegen jede Version der Liste grün: Für jede Version entspricht die versionsabhängige Leerraum-Menge der Erkennung der Antwort des Servers, auch für `\v` diesseits und jenseits der Versionsgrenze; weicht eine Version ab, greift die Rückführung aus §4.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
- [ ] Die Nahtstelle aus §1 steht in `slice-v1-abschluss-postgres-versionen-zusage` (§1 und §6), mit den Namen des Implementers und der Nachzählung im selben Commit (`AGENTS.md` §3.13); `grep -n "slice-v1-abschluss-postgres-versionen-matrix"` findet sie dort.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-053-postgres-versionen.md`, `harness/vertraege/README.md`, `harness/vertraege/SPEC-052-gate-kette.md` | new, update (Architect, geschrieben vor dem ersten Code-Commit) | Vertrag des Werkzeugs und seiner Gegenprobe nach [MR-002](../../../../harness/conventions.md#mr-002): Randformen aus §6, Rot-Bedingungen, Grenzen; `SPEC-052` Punkt 2, 4 und 13 nehmen `test-postgres-versionen` in die Ausnahmeliste auf |
| `docs/plan/adr/0042-versionsmatrix-der-integrationstests.md`, Index | new, update (Architect, Proposed) | [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) nennt *ein* gepinntes PostgreSQL-Image; die Matrix und ihre Teilmenge sind eine Ergänzung mit Alternativen. Entscheidung des Architect: eigene ADR, weil die Wahl der Variante C (Nutzer, 2026-10-05) mit dem Slice-Plan archiviert würde. Sie wird erst mit der Annahme durch den Nutzer `Accepted`; eine ADR zählt als eigene Schicht (§8) |
| `harness/mk/postgres-versionen.mk` (neu) | new | Liste der Versionen mit Tag und Digest, Testauswahl an genau einer Stelle (einschließlich des Falls SCRAM-SHA-256 aus `slice-v1-abschluss-anmeldung`), Ziel `test-postgres-versionen`, Gegenprobe, `GATE_CHECKS +=` |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | update | `test-postgres-versionen` steht in der Ausnahmeliste (`SPEC-053` Punkt 13; die Gegenprobe steht dort nach ihrem Namen); die Gegenprobe der schnellen Kette kennt das Gate in ihren Listen (`LANGSAM`, Namensmuster) und bekommt eine Mutante für den Eintrag |
| `tools/test/run-integration-tests.sh` | update | Testauswahl je Phase als Parameter; ohne ihn verhält sich der Runner wie heute (beide Phasen, volle Suite); ein Lauf über mehrere Images baut das Testimage einmal; je Version und Phase die Prüfung, dass die erwarteten Tests liefen |
| `tools/` (Gegenprobe des Ziels, neu) und Mutanten-Fälle | new | je Punkt des Vertrags ein Fall und eine Mutation, die nur ihn bricht (`AGENTS.md` §3.10); Image ohne Digest, fehlende Version, Testauswahl ohne Treffer |
| `test/integration` (`lebendpruefung_e2e_test.go` und die Tests der Teilmenge) | update (bedingt) | nur wenn ein Test eine Annahme über die Serverversion trägt (Meldungstexte, Parameter wie `server_version`); die Tabelle der Lebendprüfungen bleibt die aus `LH-FA-09.a`, die Erwartung je Zeile folgt der Version des Servers |
| `harness/README.md` §Sensors, `docs/user/abdeckung-*.md` | update | neues Ziel mit Bindung; Abdeckung, falls Deklarationen sich ändern |

Nicht in dieser Tabelle, weil nach einer Rückführung (§4) in einem eigenen Slice: `internal/hexagon/services` (`lebendpruefung.go` und Tabellentest) und `LH-FA-09.a`, falls eine Version von der Versionsgrenze abweicht, die `slice-extended-query-lebendpruefung` per Sonde bestimmt hat.

**Variante der Matrix** — geprüft gegen [`AGENTS.md`](../../../../AGENTS.md) §3.6 (kein Gate wird gelockert) und die Laufzeit von `make gates`:

- *A — volle Suite gegen jede Version in `make gates`:* belegt am meisten, vervielfacht aber die Laufzeit der Integrationstests mit der Zahl der Versionen, für Tests ohne Serverabhängigkeit ohne Gewinn. Verworfen.
- *B — Matrix als eigenes Ziel außerhalb von `make gates`:* billig, aber die Zusage in `SPEC-029` stünde dann ohne Gate; ein Lauf, den niemand ruft, belegt nichts (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Verworfen.
- **C — gewählt, Laufzeit-Zuwachs vom Nutzer akzeptiert (2026-10-05):** volle Suite gegen die Referenz 17 wie heute (`make test-integration`, unverändert), dazu `make test-postgres-versionen` in `make gates` mit der serverabhängigen Teilmenge gegen jede andere Version. Mehrkosten je Version: ein PostgreSQL-Start und die Teilmenge in zwei Phasen, das Testimage einmal. Keine bestehende Prüfung fällt weg; die Teilmenge ist erweiterbar, wenn ein späterer Slice der Welle einen serverabhängigen Abnahmetest liefert.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-anmeldung` liegt in `done/` (Entscheidung des Nutzers vom 2026-10-07: die Matrix fährt den Fall SCRAM-SHA-256 seines Anmeldetests, §1), und `slice-extended-query-lebendpruefung` ist `done`; beides ist erfüllt. Er steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) als Schritt 32a, direkt hinter `slice-v1-abschluss-anmeldung-doku`. Erster Schritt nach dem Start: den Versionsbereich gegen die Release-Liste von PostgreSQL bestätigen (§6).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): (a) **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- harness tools test internal spec docs/plan/adr docs/user`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 1500** geht der Slice zurück. **Meldeschwelle** neben der Marke: Erreicht die Summe **1200**, meldet der Implementer den Stand dem Architect, bevor er weiterarbeitet; über Schnitt oder Anhebung der Marke entscheidet der Planner. **Schätzung** (Planner, 2026-10-14: 900 bis 1150; Architect, 2026-10-14, bestätigt als obere Grenze, in dieser Metrik): Vertrag `SPEC-053`, [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) und die Änderung an `SPEC-052` liegen vor der Basis (der Architect schreibt sie vor dem Übergang nach `in-progress/`) und zählen nicht in die Messung. Gemessen werden Werkzeugcode rund 200 Zeilen (Fragment 40, Runner 90, Prüfung der Eingabe 50, Anpassung der Tests und der Listen 20), Gegenprobe nach dem Faktor 2 rund 400 bis 450 (rund 25 Fälle mit Mutation, dazu die Attrappe von `docker`), Verweise (`harness/README.md`) rund 15, zusammen rund 650 bis 700, mit der Review-Runde (20 bis 30 %) **750 bis 950**. Marke und Meldeschwelle bleiben; ein Schnitt ist nicht nötig. (b) Eine Version des Bereichs verlangt mehr als die Leerraum-Menge der Lebendprüfung — etwa einen anderen Handshake, andere Antworten im Record oder Änderungen am Replay-Ablauf. (c) Eine Version weicht von der Versionsgrenze der Erkennung ab: Die Berichtigung berührt `internal/hexagon/services` und `LH-FA-09.a` und hebt die Schichten über zwei (§8). In (b) und (c) trennt sich die Anpassung als eigener Slice ab, den der Planner schneidet und vor diesen stellt; die Matrix ohne sie bleibt hier und wartet.
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

**Randformen** (`AGENTS.md` §3.12) des neuen Vertrags (Werkzeug `test-postgres-versionen` mit Gegenprobe). Entschieden vor dem Code in `harness/vertraege/` (MR-002, Technik-Stratum der Werkzeuge) durch den Architect; bereits entschieden und zu übernehmen:

- *Lesart des Bereichs:* fünf jüngste erschienene Hauptversionen als feste Liste; ist 19 beim Start nicht erschienen, gelten die fünf jüngsten erschienenen (Nutzer, 2026-10-05).
- *Laufzeit:* Der Zuwachs der Variante C ist akzeptiert (Nutzer, 2026-10-05).
- *Anmeldung:* Teilmenge enthält den Fall SCRAM-SHA-256, Klartext und MD5 nicht (Nutzer, 2026-10-07).
- *Szenario 18 und Szenario 16 Satz 3:* nur gegen 17, nicht in der Matrix (Nutzer, 2026-10-12).

Entschieden vom Architect am 2026-10-14, vor dem Code, in [`SPEC-053`](../../../../harness/vertraege/SPEC-053-postgres-versionen.md) (Punkt in Klammern); die Gegenprobe führt je Form einen Fall mit Mutation:

- *Fehlender, leerer, unbekannter Wert:* neun Formen der Eingabe sind rot vor dem ersten Aufruf von Docker (Punkt 2: leere Liste, Referenz fehlt, Hauptversion doppelt oder Alias, keine andere Version, Eintrag ohne `=` oder mit falscher Hauptversion, Image ohne Tag, ohne Digest, Digest ungültig, Referenz-Image ungleich `POSTGRES_IMAGE`); Digest, der zu einem anderen Tag gehört, nicht ziehbares Image und unlesbare Hauptversion sind rot am laufenden Server (Punkt 4); `server_version` ohne führende Ziffer macht den Lebendprüfungs-Test rot (Punkt 11).
- *Abbruch und Herunterfahren:* rote Phase 1 lässt Phase 2 dieser Version entfallen, der Bericht nennt die Phase; das Ziel fährt alle Versionen und fasst die roten zusammen (Punkt 6); Aufräumen nach Ende, Fehler, `SIGINT` und `SIGTERM`, nicht nach `SIGKILL` (Punkt 7 und 14).
- *Verhalten je Serverversion:* SCRAM-SHA-256 je Version, Klartext und MD5 außerhalb der Matrix (Punkt 5); die Lesart von `server_version` und die Grenze der Produktspezifikation (Punkt 11).
- *Aufruf-Voraussetzung:* nacheinander, ein Server zur Zeit; gleichzeitige Läufe und `make -j gates` (beide Docker-Gates gleichzeitig) sind nicht zugesagt und haben keinen Fall, akzeptiertes Negativ mit dem Grund von `SPEC-052` Punkt 8 (Punkt 10 und 14).
- *Aufrufarten des Treibers:* Das Ziel schreibt keinen Stand; ein einzeiliges Rezept ohne `$(MAKE)` und ohne `+` fährt unter `-n`, `-q`, `-t` und `MAKEFLAGS` kein Docker, `-i` überliest den Exit-Code und belegt nichts (Punkt 9); je Flag ein Fall.
- *Ausnahmeliste der schnellen Kette:* ja, `test-postgres-versionen` steht in `GATE_SCHNELL_AUSNAHME`, seine Gegenprobe nach ihrem Namen; `SPEC-052` Punkt 2, 4 und 13 sind angepasst, die Positivliste des Hooks bleibt unverändert (Punkt 13).
- *Der Lauf der Referenz:* `make test-integration` fährt mit dem geänderten Runner wie bisher (Punkt 8, `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`).

Risiken:

- **Versionsbereich unbestätigt.** Entschieden (Nutzer, 2026-10-05): die fünf jüngsten erschienenen Hauptversionen als feste Liste; ist 19 nicht erschienen, gelten 14 bis 18. Bestätigt am 2026-10-14 (Architect): 19 ist nicht erschienen, die Liste ist 14 bis 18. Erscheint 19 vor dem Start, liest der Implementer die Release-Liste von PostgreSQL erneut, bevor die Liste in `harness/mk/postgres-versionen.mk` steht; PostgreSQL 14 endet am 2026-11-12; die Spezifikation nennt sie erst in `slice-v1-abschluss-postgres-versionen-zusage` — **Ausgang:** offen bis Closure.
- **Versionsgrenze für `\v`.** Die Erkennung der Lebendprüfung hängt nach `slice-extended-query-lebendpruefung` von `server_version` ab; die Grenze, ab der `\v` Leerraum ist, hat dort eine Sonde gegen 14 bis 18 auf 17 bestimmt; 19 ist nicht geprüft. Ob sie für jede Version der Liste stimmt und ob die übrigen Zeichen der Tabelle in `LH-FA-09.a` bei allen gleich gelten, prüft erst dieser Slice. Weicht eine Version ab, greift die Rückführung aus §4 (b) oder (c) — **Ausgang:** offen bis Closure.
- **Teilmenge fällt still leer.** Die Teilmenge wählt Tests über ihren Namen; ein umbenannter Test fiele aus der Matrix, und `go test -run` mit einem Muster ohne Treffer endet grün (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`). Der Runner muss je Version und Phase prüfen, dass die erwarteten Tests liefen, nicht nur, dass keiner rot war — **Ausgang:** offen bis Closure.
- **Laufzeit des Gates.** Je Version kommt ein PostgreSQL-Start und die Teilmenge in zwei Phasen hinzu. Der Zuwachs der Variante C ist vom Nutzer akzeptiert (2026-10-05); gemessen wird `make gates` vor und nach dem Slice, die Closure-Notiz nennt beide Zahlen. Liegt der gemessene Zuwachs deutlich über der Schätzung, ist die Antwort eine kleinere Teilmenge oder ein Carveout, keine stille Herausnahme aus `GATE_CHECKS` ([`AGENTS.md`](../../../../AGENTS.md) §3.6) — **Ausgang:** offen bis Closure.
- **Anmeldung je Image.** Der Fall SCRAM-SHA-256 braucht je Server einen Benutzer mit Passwort und eine Regel in `pg_hba.conf`; wie `slice-v1-abschluss-anmeldung` das gegen 17 einrichtet, muss für jedes Image des Bereichs gleich wirken (etwa `password_encryption`, Einstiegsskripte des Images). Weicht ein Image ab, gehört die Einrichtung in die Matrix, nicht in den Anmeldetest — **Ausgang:** offen bis Closure.
- **Zweite Phase verdeckt.** Der Runner endet nach einer roten ersten Phase; eine Version, deren Record scheitert, zeigt nicht, ob Szenario 2 oder 7 gegen ihre Aufzeichnung trüge (`BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar`). Die Matrix erbt das; der Bericht je Version nennt die Phase, in der er endete — **Ausgang:** offen bis Closure.
- **Namen der Nahtstelle.** Die Namen von Datei, Variablen und Aufruf wählt der Implementer; der Nehmer `slice-v1-abschluss-postgres-versionen-zusage` braucht sie für die Spezifikation und den nachfahrbaren Lauf. Gegenmittel: Der Planner trägt die Namen bei der Closure dieses Slice in §1 und §6 des Nehmers ein, mit Nachzählung im selben Commit; solange `grep` nach ihnen im Nehmer leer bleibt, hat das Risiko keinen Ausgang — **Ausgang:** offen bis Closure.
- **Abweichung hebt die Schichten.** Eine Berichtigung der Versionsgrenze berührt Produktcode und `LH-FA-09.a` und damit mehr Schichten als dieser Plan trägt (§8); Entscheidung des Nutzers vom 2026-10-14 war „kleiner Produktcode nur bei Abweichung“, der Plan setzt dafür die Rückführung (§4 (c)). Empfehlung des Architect (keine Entscheidung des Nutzers): die Rückführung bleibt, keine Ausnahme von der Zwei-Schichten-Regel; die Sonde von `slice-extended-query-lebendpruefung` deckt 14 bis 18, und die volle Suite lief von Hand gegen 14, 16, 17 und 18 grün, die Abweichung ist damit unwahrscheinlich. Ob der Nutzer sie trotzdem im Lauf berichtigt sehen will, ist offen (Frage an den Nutzer) — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-14); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` — 1×, offen; dieser Slice ist seine Adresse.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` — 1×, offen; der Slice baut am Integrations-Runner, und eine Pflichtprüfung darf beim Umbau nicht entfallen. Deshalb bleibt `make test-integration` unverändert, und die Matrix kommt daneben (§1).
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` — 2×, offen; eine Testauswahl per Namensmuster liest sich als Zusage und endet ohne Treffer grün (§6). Ein drittes Auftreten machte den Eintrag zur Lücke: Der Runner prüft je Version und Phase, dass die erwarteten Tests liefen.
- `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` — 1×, offen; die Matrix erbt die zwei Phasen des Runners (§6).
- `BEO-REPO/meldeschwelle-ohne-artefakt` (1×, offen) und `BEO-REPO/groessenprobe-ohne-messbaren-stand`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (verkörpert, `AGENTS.md` §3.14) — §4 nennt Marke mit Metrik, Meldeschwelle und Schätzung im Faktor der Tests.
- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` — verkörpert; die Gegenprobe des Ziels und die Ausnahmeliste der schnellen Kette sind deshalb Randform (§6).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 29×, verkörpert (`AGENTS.md` §3.12); die Randformen stehen in §6 vor dem Code.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 55×, verkörpert (`AGENTS.md` §3.11); der Grund für den Schnitt mit der Matrix zuerst.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — 34×, verkörpert (`AGENTS.md` §3.10); die Gegenprobe mit fehlender Version und Image ohne Digest steht in §2.
- `BEO-REPO/plan-folgt-korrektur-nicht` — 45×, verkörpert (`AGENTS.md` §3.9); gilt als Regel, auch für den Kopf und die Liste der Versionen.

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure ihrer Welle, nicht hier.

Nachgezählt beim Anlegen (Planner, 2026-10-14, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Gate mit Teilmenge und Anmeldung; Lebendprüfung gegen jede Version). Schichten: zwei, Harness mit seinen Tests (`harness/`, `tools/`, `test/integration`) und, wenn der Architect eine ADR schreibt, Entscheidungen. Eine ADR zählt als eigene Schicht, nicht mit der Spezifikation: Die Source Precedence führt sie auf einem eigenen Rang (4) neben den Straten 1 bis 3, und der Ursprungsplan zählte „Spezifikation, Entscheidungen, Harness und Dokumentation“ als vier. Der Vertrag des Werkzeugs liegt unter `harness/vertraege/` und zählt zur Schicht Harness (MR-002). Ohne ADR: eine Schicht. Der Abweichungsfall (§4 (c)) kommt nicht hinzu, sondern geht zurück. Schätzung: 900 bis 1150 Zeilen in der Metrik von §4.

Nachgezählt vom Architect (2026-10-14, vor dem Code): Liefer-Punkte weiter 2, Schichten zwei (Harness mit `tools/` und `test/integration`; Entscheidungen mit [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md)). Die Änderung an `SPEC-052` und `SPEC-053` zählt zur Schicht Harness (MR-002). Schätzung in der Metrik von §4: 750 bis 950; Vertrag, ADR und `SPEC-052` liegen vor der Basis. Größe bestätigt, kein Schnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
