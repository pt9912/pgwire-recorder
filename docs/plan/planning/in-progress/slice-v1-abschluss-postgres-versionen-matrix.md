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
| `docs/plan/adr/0042-versionsmatrix-der-integrationstests.md`, Index | new, update (Architect, Accepted durch Entscheidung des Nutzers vom 2026-10-14) | [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) nennt *ein* gepinntes PostgreSQL-Image; die Matrix und ihre Teilmenge sind eine Ergänzung mit Alternativen. Entscheidung des Architect: eigene ADR, weil die Wahl der Variante C (Nutzer, 2026-10-05) mit dem Slice-Plan archiviert würde. Eine ADR zählt als eigene Schicht (§8) |
| `harness/mk/postgres-versionen.mk` (neu) | new | Liste der Versionen mit Tag und Digest, Testauswahl an genau einer Stelle (einschließlich des Falls SCRAM-SHA-256 aus `slice-v1-abschluss-anmeldung`), Ziel `test-postgres-versionen`, Gegenprobe, `GATE_CHECKS +=` |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | update | `test-postgres-versionen` steht in der Ausnahmeliste (`SPEC-053` Punkt 13; die Gegenprobe steht dort nach ihrem Namen); die Gegenprobe der schnellen Kette kennt das Gate in ihren Listen (`LANGSAM`, Namensmuster) und bekommt eine Mutante für den Eintrag |
| `tools/test/run-integration-tests.sh` | update | Testauswahl je Phase als Parameter; ohne ihn verhält sich der Runner wie heute (beide Phasen, volle Suite); ein Lauf über mehrere Images baut das Testimage einmal; je Version und Phase die Prüfung, dass die erwarteten Tests liefen |
| `tools/test/run-postgres-versionen.sh` (neu) | new | Treiber des Ziels: prüft die Eingabe (neun Formen) vor dem ersten Aufruf von docker, baut das Testimage einmal, fährt je Version den Runner nacheinander, fasst rote Versionen mit Stufe zusammen, leitet `SIGINT` und `SIGTERM` weiter |
| `tools/test/postgres-versionen-gegenprobe.sh` (neu) und `tools/harness/gates-schnell-gegenprobe.sh` | new, update | Gegenprobe des Ziels in Kopien des Arbeitsbaums mit einer Attrappe von `docker`: je Punkt des Vertrags Fälle und Mutationen (`AGENTS.md` §3.10); die Gegenprobe der schnellen Kette bekommt `LANGSAM`, das Namensmuster und je eine Mutante für `test-postgres-versionen` und `test-postgres-versionen-gegenprobe` |
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
- *Rückgaben des Implementers, entschieden vom Architect am 2026-10-14 nach dem Code-Commit fc3e987* (`SPEC-053`): Stufe „Aufbau“ als dritte Stufe neben Phase 1 und 2, jeder Fall der Formen aus Punkt 4 und der Einrichtung aus Punkt 5 prüft sie (Punkt 4 und 6); roter Bau des Testimages beendet das Ziel sofort rot ohne Version, Fall prüft, dass außer dem Bau kein Aufruf von `docker` stattfand (Punkt 6); Zusammenführung von stderr in stdout im Runner ist akzeptiertes Negativ, Exit-Code und Tests gleich (Punkt 8 und 14 (8)); fehlender Test in Phase 2 heißt „Phase 2“, ein Fall mit Mutation (Punkt 3); Signal: Zusage gilt für das Skript, Verhalten von `make` ist akzeptiertes Negativ, weil das Rezept es nicht ändern kann und ein Fall auf die Zeit nach dem Ende von `make` zeitabhängig wäre (Punkt 7 und 14 (7)); „ziehen oder starten“ ist ein Schritt, die Zusage ist auf „starten“ enger gefasst, `p4-start` genügt (Punkt 4; Fall früher `p4-pull`).
- *Ausnahmeliste der schnellen Kette:* ja, `test-postgres-versionen` steht in `GATE_SCHNELL_AUSNAHME`, seine Gegenprobe nach ihrem Namen; `SPEC-052` Punkt 2, 4 und 13 sind angepasst, die Positivliste des Hooks bleibt unverändert (Punkt 13).
- *Der Lauf der Referenz:* `make test-integration` fährt mit dem geänderten Runner wie bisher (Punkt 8, `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`).

Risiken:

- **Versionsbereich unbestätigt.** Entschieden (Nutzer, 2026-10-05): die fünf jüngsten erschienenen Hauptversionen als feste Liste; ist 19 nicht erschienen, gelten 14 bis 18. Bestätigt am 2026-10-14 (Architect, vom Nutzer am selben Tag bestätigt): 19 ist nicht erschienen, die Liste ist 14 bis 18, Referenz 17, Matrix 14, 15, 16 und 18. Erscheint 19 vor dem Start, liest der Implementer die Release-Liste von PostgreSQL erneut, bevor die Liste in `harness/mk/postgres-versionen.mk` steht; PostgreSQL 14 endet am 2026-11-12; die Spezifikation nennt sie erst in `slice-v1-abschluss-postgres-versionen-zusage` — **Ausgang:** offen bis Closure.
- **Versionsgrenze für `\v`.** Die Erkennung der Lebendprüfung hängt nach `slice-extended-query-lebendpruefung` von `server_version` ab; die Grenze, ab der `\v` Leerraum ist, hat dort eine Sonde gegen 14 bis 18 auf 17 bestimmt; 19 ist nicht geprüft. Ob sie für jede Version der Liste stimmt und ob die übrigen Zeichen der Tabelle in `LH-FA-09.a` bei allen gleich gelten, prüft erst dieser Slice. Weicht eine Version ab, greift die Rückführung aus §4 (b) oder (c) — **Ausgang:** offen bis Closure.
- **Teilmenge fällt still leer.** Die Teilmenge wählt Tests über ihren Namen; ein umbenannter Test fiele aus der Matrix, und `go test -run` mit einem Muster ohne Treffer endet grün (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`). Der Runner muss je Version und Phase prüfen, dass die erwarteten Tests liefen, nicht nur, dass keiner rot war — **Ausgang:** offen bis Closure.
- **Laufzeit des Gates.** Je Version kommt ein PostgreSQL-Start und die Teilmenge in zwei Phasen hinzu. Der Zuwachs der Variante C ist vom Nutzer akzeptiert (2026-10-05); gemessen wird `make gates` vor und nach dem Slice, die Closure-Notiz nennt beide Zahlen. Liegt der gemessene Zuwachs deutlich über der Schätzung, ist die Antwort eine kleinere Teilmenge oder ein Carveout, keine stille Herausnahme aus `GATE_CHECKS` ([`AGENTS.md`](../../../../AGENTS.md) §3.6) — **Ausgang:** offen bis Closure.
- **Anmeldung je Image.** Der Fall SCRAM-SHA-256 braucht je Server einen Benutzer mit Passwort und eine Regel in `pg_hba.conf`; wie `slice-v1-abschluss-anmeldung` das gegen 17 einrichtet, muss für jedes Image des Bereichs gleich wirken (etwa `password_encryption`, Einstiegsskripte des Images). Weicht ein Image ab, gehört die Einrichtung in die Matrix, nicht in den Anmeldetest — **Ausgang:** offen bis Closure.
- **Zweite Phase verdeckt.** Der Runner endet nach einer roten ersten Phase; eine Version, deren Record scheitert, zeigt nicht, ob Szenario 2 oder 7 gegen ihre Aufzeichnung trüge (`BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar`). Die Matrix erbt das; der Bericht je Version nennt die Phase, in der er endete — **Ausgang:** offen bis Closure.
- **Namen der Nahtstelle.** Die Namen von Datei, Variablen und Aufruf wählt der Implementer; der Nehmer `slice-v1-abschluss-postgres-versionen-zusage` braucht sie für die Spezifikation und den nachfahrbaren Lauf. Gegenmittel: Der Planner trägt die Namen bei der Closure dieses Slice in §1 und §6 des Nehmers ein, mit Nachzählung im selben Commit; solange `grep` nach ihnen im Nehmer leer bleibt, hat das Risiko keinen Ausgang — **Ausgang:** offen bis Closure.
- **Abweichung hebt die Schichten.** Eine Berichtigung der Versionsgrenze berührt Produktcode und `LH-FA-09.a` und damit mehr Schichten als dieser Plan trägt (§8); Entscheidung des Nutzers vom 2026-10-14 war „kleiner Produktcode nur bei Abweichung“, der Plan setzt dafür die Rückführung (§4 (c)). Entscheidung des Nutzers vom 2026-10-14: Das Gate berichtigt die Grenze nicht im Lauf; die Rückführung aus §4 (c) bleibt, es gibt keine Ausnahme von der Zwei-Schichten-Regel. Das Gate wird bei Abweichung rot, die Berichtigung ist ein eigener Slice (die Sonde von `slice-extended-query-lebendpruefung` deckt 14 bis 18) — **Ausgang:** offen bis Closure.

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

### Belege des Implementers

**Stand.** Basis 978341e (Übergang nach `in-progress/`), Code-Commit fc3e987, Nacharbeit zur Entscheidung des Architect (f0a2280) in 8f97a05. Größe in der Metrik von §4 (`git diff --numstat 978341e..8f97a05 -- harness tools test internal spec docs/plan/adr docs/user`): **823** (Einfügungen plus Löschungen; bei fc3e987 767), unter Meldeschwelle 1200 und Marke 1500; davon `tools/test/postgres-versionen-gegenprobe.sh` 502, `harness/vertraege/SPEC-053-postgres-versionen.md` 40, `tools/test/run-integration-tests.sh` 132, `tools/test/run-postgres-versionen.sh` 100, `harness/mk/postgres-versionen.mk` 31, Rest 18 (`harness/README.md`, `gates-schnell.mk`, Gegenprobe der schnellen Kette).

**Bereich bestätigt (2026-10-14, vor dem Code).** `docker pull postgres:19-alpine` endet mit `manifest unknown`; `https://www.postgresql.org/support/versioning/` führt 18.6, 17.11, 16.15, 15.19, 14.24 als unterstützte Hauptversionen, die Startseite nennt „PostgreSQL 19 Beta 4 Released“ vom 2026-09-24. 19 ist nicht erschienen, die Liste ist 14 bis 18. Digests aus `docker image inspect --format '{{index .RepoDigests 0}}'` nach `docker pull -q postgres:<n>-alpine`: 14 `4ea9e5ed…e77a`, 15 `f7d23353…bf81`, 16 `721873c3…80ea`, 17 `b0f9560a…4b24` (gleich `POSTGRES_IMAGE`), 18 `77f58511…1873` (die vollen Digests stehen im Fragment).

**Läufe** (Docker 29.9.0, Images lokal vorhanden, Versionen nacheinander):

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make gates` vor dem Slice (Basis 978341e, frischer Worktree, Build-Cache anderer Kontextpfad) | grün | 589 s |
| `make gates` nach dem Slice (fc3e987, Arbeitsbaum, schreibt den Nachweis) | grün | 561 s |
| `make test-postgres-versionen` (14, 15, 16, 18; Testimage aus dem Cache) | grün | 22 s |
| dasselbe mit einer Version neben der Referenz | grün | 6 s |
| `make test-postgres-versionen-gegenprobe` (unveränderte Kopie, 87 Mutationen; Stand 8f97a05) | grün | 86 s |
| `make gates-schnell-gegenprobe` (mit den beiden neuen Mutanten; Stand 8f97a05) | grün | 21 s |
| `make gates` nach der Nacharbeit zur Entscheidung vom 2026-10-14 (8f97a05, Arbeitsbaum, schreibt den Nachweis) | grün | 553 s |
| `make test-integration` (Referenz 17, volle Suite, geänderter Runner) | grün | 129 s |
| `make docs-check` | 0 Befunde | — |

Die Gesamtzeit von `make gates` ist kein Vergleich der beiden Zeilen: Die Messung vor dem Slice lief im Worktree mit anderem Kontextpfad und kälterem Cache, die danach im Arbeitsbaum; die Differenz liegt unter der Streuung. Der Zuwachs des Slice ist die Summe seiner eigenen Gates: rund 22 s (`test-postgres-versionen`, mit Pull der Images mehr) plus 80 s (Gegenprobe) plus 3 s (Gegenprobe der schnellen Kette), zusammen rund 105 s.

**Zusage · Mutation · roter Test.** Die Gegenprobe prüft die unveränderte Kopie (alle Fälle grün) und fährt jede Mutation in einer Kopie; Ausgang 0 nur, wenn bei jeder Mutation von den Fällen ihrer Gruppe genau die benannten rot sind. Jede Mutation wurde dabei gesehen; „wirkungslos“ meldet die Gegenprobe selbst. Die Fälle heißen nach dem Punkt von `SPEC-053`; die Attrappe von `docker` ist im Kopf der Gegenprobe beschrieben. Spalte 1 ist der Punkt.

| Punkt | Zusage | Mutation (Name in der Gegenprobe) | roter Test |
|---|---|---|---|
| 1 | Die Liste des Fragments ist gültig (vier Versionen neben der Referenz) | Fragment: `14-alpine@sha256:4ea9` → `14-alpine@sha256:4EA9` (`p1-liste`) | `p1-liste-gueltig` |
| 2a | leere Liste ist rot vor docker | Treiber: Prüfzeile `P2a` entfällt (`p2a`) | `p2a-leere-liste` |
| 2b | Referenz fehlt in der Liste ist rot vor docker | Treiber: Prüfzeile `P2b` entfällt (`p2b`) | `p2b-referenz-fehlt` |
| 2c | Hauptversion doppelt ist rot vor docker | Treiber: Prüfzeile `P2c` entfällt (`p2c`) | `p2c-doppelt` |
| 2c | Hauptversion der Referenz doppelt ist rot vor docker | Treiber: Prüfzeile `P2c-ref` entfällt (`p2c-referenz`) | `p2c-referenz-doppelt` |
| 2d | keine andere Version als die Referenz ist rot vor docker | Treiber: Prüfzeile `P2d` entfällt (`p2d`) | `p2d-keine-andere` |
| 2e | Eintrag ohne `=` ist rot vor docker | Treiber: Prüfzeile `P2e1` entfällt (`p2e1`) | `p2e1-ohne-gleich` |
| 2e | leere Hauptversion ist rot vor docker | Treiber: Prüfzeile `P2e2` entfällt (`p2e2`) | `p2e2-leere-hauptversion` |
| 2e | Hauptversion nicht aus Ziffern ist rot vor docker | Treiber: Prüfzeile `P2e3` entfällt (`p2e3`) | `p2e3-nicht-ziffern` |
| 2e | Hauptversion mit führender Null ist rot vor docker | Treiber: Prüfzeile `P2e4` entfällt (`p2e4`) | `p2e4-fuehrende-null` |
| 2f | Image ohne Tag ist rot vor docker | Treiber: Prüfzeile `P2f` entfällt (`p2f`) | `p2f-ohne-tag` |
| 2g | Image ohne Digest ist rot vor docker | Treiber: Prüfzeile `P2g` entfällt (`p2g`) | `p2g-ohne-digest` |
| 2h | Digest mit 63 Ziffern ist rot (64 Ziffern verlangt) | Treiber: `{64}` → `{1,64}` (`p2h-kurz`) | `p2h-digest-kurz` |
| 2h | Digest mit 65 Ziffern ist rot | Treiber: `{64}` → `{64,}` (`p2h-lang`) | `p2h-digest-lang` |
| 2h | Digest mit Großbuchstaben ist rot („kleine“ Hexziffern) | Treiber: `[0-9a-f]{64}` → `[0-9a-fA-F]{64}` (`p2h-gross`) | `p2h-digest-gross` |
| 2h | Digest mit einem Zeichen außerhalb der Hexziffern ist rot | Treiber: `[0-9a-f]{64}` → `[0-9a-z]{64}` (`p2h-kein-hex`) | `p2h-digest-kein-hex` |
| 2h | Digest eines anderen Verfahrens als sha256 ist rot | Treiber: `^sha256:` → `^sha[0-9]+:` (`p2h-verfahren`) | `p2h-digest-verfahren` |
| 2i | Image der Referenz ungleich `POSTGRES_IMAGE` ist rot vor docker | Treiber: Prüfzeile `P2i` entfällt (`p2i`) | `p2i-referenz-image` |
| 2 | jede der neun Formen endet vor dem ersten Aufruf von docker (Eingabefehler lässt den Lauf fortfahren) | Treiber: `if [ "$fehler" -ne 0 ]; then exit 2; fi` → `if [ "$fehler" -ne 0 ]; then :; fi` (`p2-vor-docker`) | alle 17 Fälle `p2*` |
| 3 | ein übersprungener Test (`SKIP`) macht die Version rot | Runner: `if grep -Eq "^ *--- SKIP: $name \(" "$aus"; then # P3-skip` → `if false; then # P3-skip` (`p3-skip`) | `p3-skip` |
| 3 | eine Auswahl ohne Treffer ist rot, nicht grün | Runner: `if ! grep -Eq '^ *--- PASS: ' "\$aus"; then # P3-treffer` → `if false; then # P3-treffer` (`p3-treffer`) | `p3-ohne-treffer` |
| 3 | ein erwarteter Test ohne PASS-Zeile (umbenannt, nicht gewählt) macht die Version rot | Runner: `rot=1 # P3-fehlt` → `rot=0 # P3-fehlt` (`p3-fehlt`) | `p3-skip`, `p3-fehlt-phase2` und alle 8 Fälle `p3-element-*` |
| 3 | ein in Phase 2 fehlender Test heißt „Phase 2“ (Test, Version, Phase in der Meldung) | Runner: `ph="$2"` → `ph=1` (`p3-phase-fest`) | `p3-fehlt-phase2` und die zwei Fälle `p3-element-` der Phase 2 (sie prüfen die Phase der Meldung ebenso) |
| 3 | dasselbe, nur für den Fall `p3-fehlt-phase2` (Version 16) | Runner: `… ph="$2"; if [ "$version" = 16 ]; then ph=1; fi` (`p3-phase-fest-16`) | `p3-fehlt-phase2` |
| 3 | die Auswahl von `-test.run` nennt jeden Test der Teilmenge | Runner: `top+=("${n%%/*}")` → `:` (`p3-auswahl-top`) | `p3-auswahl` und alle 8 Fälle `p3-element-*` |
| 3 | die Auswahl nennt den Fall SCRAM-SHA-256 (Unterfall) der Anmeldung | Runner: `sub+=("${n#*/}")` → `:` (`p3-auswahl-sub`) | `p3-auswahl`, `p3-element-TestE2ERecordAnmeldung/scram-sha-256` |
| 4 | abweichende, unlesbare Hauptversion des Servers ist rot (Vergleich am laufenden Server) | Runner: `if [ "${gemeldet%%[!0-9]*}" != "$PGR_HAUPTVERSION" ]; the…` → `if false; then # P4-version` (`p4-version`) | `p4-version-abweichend`, `p4-version-unlesbar`, `p4-version-ohne-ziffer` |
| 4 | nicht lesbare Hauptversion (leer) ist rot | Runner: `!= "$PGR_HAUPTVERSION" ]; then # P4-version` → `!= "$PGR_HAUPTVERSION" ] && [ -n "$gemeldet" ]; then # P4…` (`p4-unlesbar`) | `p4-version-unlesbar` |
| 4 | Hauptversion ohne führende Ziffer ist rot | Runner: `${gemeldet%%[!0-9]*}` → `${gemeldet//[!0-9]/}` (`p4-ziffer`) | `p4-version-ohne-ziffer` |
| 4 | der Server ist nach 60 Sekunden nicht bereit: 120 Versuche | Runner: `seq 1 120` → `seq 1 60` (`p4-60a`) | `p4-60-sekunden` |
| 4 | der Server ist nach 60 Sekunden nicht bereit: Wartezeit 0,5 Sekunden | Runner: `sleep 0.5` → `sleep 1` (`p4-60b`) | `p4-60-sekunden` |
| 4 | ein Server, der nicht bereit wird, ist rot und nennt die Version | Runner: `exit 1 # E-bereit` → `: # E-bereit` (`p4-bereit`) | `p4-nicht-bereit` |
| 4 | ein Image, das sich nicht starten lässt (der Start zieht es), ist rot | Runner: `"$postgres_image" >/dev/null` → `"$postgres_image" >/dev/null \|\| true` (`p4-start`) | `p4-start` |
| 4 | Stufe „Aufbau“ bei nicht bereitem Server (nach 60 s) | Runner: `exit 1 # E-bereit` → `phase 1; exit 1 # E-bereit` (`p4-aufbau-bereit`) | `p4-nicht-bereit`, `p4-60-sekunden` |
| 4 | Stufe „Aufbau“ im Fall der 60 Sekunden (Version 16) | Runner: `exit 1 # E-bereit` → `[ "$version" = 16 ] && phase 1; exit 1 # E-bereit` (`p4-aufbau-bereit-16`) | `p4-60-sekunden` |
| 4 | Stufe „Aufbau“ bei einem Image, das sich nicht starten lässt | Runner: `… >/dev/null` → `… >/dev/null \|\| { phase 1; exit 1; }` (`p4-aufbau-start`) | `p4-start` |
| 4 | Stufe „Aufbau“ bei abweichender Hauptversion | Runner: nach der Meldung `case "$gemeldet" in [0-9]*) phase 1 ;; esac` (`p4-aufbau-abweichend`) | `p4-version-abweichend` |
| 4 | Stufe „Aufbau“ bei nicht lesbarer (leerer) Hauptversion | Runner: nach der Meldung `case "$gemeldet" in '') phase 1 ;; esac` (`p4-aufbau-unlesbar`) | `p4-version-unlesbar` |
| 4 | Stufe „Aufbau“ bei Hauptversion ohne führende Ziffer | Runner: nach der Meldung `case "$gemeldet" in [!0-9]?*) phase 1 ;; esac` (`p4-aufbau-ohne-ziffer`) | `p4-version-ohne-ziffer` |
| 5 | je Version ein Benutzer mit Passwort | Runner: `CREATE ROLE play_scram LOGIN` → `CREATE ROLE play_scram NOLOGIN` (`p5-benutzer`) | `p5-benutzer` |
| 5 | Methode `scram-sha-256` für ihn in `pg_hba.conf` | Runner: `play_scram all scram-sha-256` → `play_scram all md5` (`p5-methode`) | `p5-methode`, `p5-hba-vor-allgemeiner` |
| 5 | das Passwort wird mit scram-sha-256 hinterlegt, nichts mit md5 | Runner: `SET password_encryption = 'scram-sha-256';` → `SET password_encryption = 'md5';` (`p5-passwort-verschluesselung`) | `p5-methode`, `p5-kein-md5-klartext` |
| 5 | die Zeile steht vor der allgemeinen Zeile | Runner: `{ printf "%s\n" '"$hba"'; cat "$hba"; }` → `{ cat "$hba"; printf "%s\n" '"$hba"'; }` (`p5-hba-nach-allgemeiner`) | `p5-hba-vor-allgemeiner` |
| 5 | Klartext-Passwort und MD5 laufen in der Matrix nicht | Runner: `!= scram ]; then # P5-rollen` → `!= nie ]; then # P5-rollen` (`p5-rollen`) | `p5-kein-md5-klartext` |
| 5 | scheitert die Einrichtung, ist die Version rot | Runner: `SELECT pg_reload_conf();"` → `SELECT pg_reload_conf();" \|\| true` (`p5-einrichtung`) | `p5-einrichtung-scheitert` |
| 6 | die Referenz fährt das Ziel nicht | Treiber: Prüfzeile `P6-referenz` entfällt (`p6-referenz`) | `p6-referenz-nicht` |
| 6 | die Versionen laufen in der Reihenfolge der Liste | Treiber: `for e in "${eintraege[@]}"; do   haupt="${e%%=*}"   [ "$h…` → `for e in $(printf "%s\n" "${eintraege[@]}" \| tac); do   …` (`p6-reihenfolge`) | `p6-reihenfolge`, `p6-zusammenfassung`, `p7-vor-naechster` |
| 6 | nie zwei Server gleichzeitig | Runner: `"$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 \|\| true` → `"$docker" rm -f "$tests" >/dev/null 2>&1 \|\| true` (`p6-nacheinander`) | `p6-nacheinander`, `p7-container` |
| 6 | ist Phase 1 rot, entfällt Phase 2 | Runner: `  exit "$status"` → `  :` (`p6-phase2`) | `p6-meldung-phase1`, `p6-keine-phase2`, `p6-zusammenfassung` |
| 6 | das Ziel fährt alle Versionen auch nach einer roten | Treiber: `rot+=("$haupt ($ort)") # P6-phase` → `rot+=("$haupt ($ort)"); break # P6-phase` (`p6-weiter`) | `p6-alle-versionen`, `p6-keine-phase2`, `p6-zusammenfassung` |
| 6 | die Zusammenfassung nennt die Stufe jeder roten Version | Treiber: `rot+=("$haupt ($ort)")` → `rot+=("$haupt")` (`p6-phase-fehlt`) | `p6-meldung-phase1`, `p6-meldung-phase2`, `p6-meldung-aufbau`, `p6-zusammenfassung` |
| 6 | die Meldung nennt „Phase 1“ für eine Version mit roter Phase 1 | Treiber: `1) ort="Phase 1"` → `1) ort="Phase 2"` (`p6-ort1`) | `p6-meldung-phase1`, `p6-zusammenfassung` |
| 6 | die Meldung nennt „Phase 2“ für eine Version mit roter Phase 2 | Treiber: `2) ort="Phase 2"` → `2) ort="Phase 1"` (`p6-ort2`) | `p6-meldung-phase2`, `p6-zusammenfassung` |
| 6 | die Meldung nennt „Aufbau“ für eine Version, die vor Phase 1 rot wird (Lesart, siehe Rückgabe) | Treiber: `*) ort="Aufbau"` → `*) ort="Phase 1"` (`p6-ort-aufbau`) | `p6-meldung-aufbau` |
| 6 | das Ziel endet mit Exit-Code ungleich 0, wenn eine Version rot war | Treiber: `exit 1 # P6-exit` → `exit 0 # P6-exit` (`p6-exit`) | `p6-exit` |
| 6 | das Testimage wird einmal gebaut, nicht je Version | Treiber: `PGR_BAUEN=0` → `PGR_BAUEN=1` (`p6-bauen`) | `p6-testimage-einmal` |
| 5 | Stufe „Aufbau“ bei scheiternder Einrichtung der Anmeldung | Runner: `SELECT pg_reload_conf();"` → `… \|\| { phase 1; exit 1; }` (`p5-aufbau-einrichtung`) | `p5-einrichtung-scheitert` |
| 6 | ein roter Bau des Testimages beendet das Ziel rot, ohne Server | Treiber: `if ! wait "$laeuft"; then` → `if false; then` (`p6-build-rot`) | `p6-build-rot` |
| 6 | nach rotem Bau kein weiterer Aufruf von `docker` (kein Netz, kein Lauf) | Treiber: nach der Meldung `"$docker" network create --internal x` (`p6-build-rot-aufruf`) | `p6-build-rot` (Meldung und Exit-Code bleiben, nur die Zahl der Aufrufe im Log bricht) |
| 6 | ein roter Bau endet mit Exit-Code ungleich 0 | Treiber: `exit 1` nach der Meldung → `exit 0` (`p6-build-rot-exit`) | `p6-build-rot` |
| 7 | Container einer Version sind entfernt (nach Ende, Fehler, rotem Lauf, Signal; vor der nächsten Version) | Runner: `"$docker" rm -f "$tests" "$pg" >/dev/null 2>&1 \|\| true` → `:` (`p7-container`) | `p7-container`, `p7-ende-rot`, `p7-ende-bereit`, `p7-ende-gruen`, `p7-vor-naechster`, `p7-sigINT-gruppe`, `p7-sigTERM-gruppe`, `p7-sigINT-skript`, `p7-sigTERM-skript` |
| 7 | das Netz einer Version ist entfernt | Runner: `"$docker" network rm "$netz" >/dev/null 2>&1 \|\| true` → `:` (`p7-netz`) | `p7-netz`, `p7-ende-rot`, `p7-ende-bereit`, `p7-ende-pull`, `p7-ende-gruen`, `p7-vor-naechster`, `p7-sigINT-gruppe`, `p7-sigTERM-gruppe`, `p7-sigINT-skript`, `p7-sigTERM-skript` |
| 7 | das Volume einer Version ist entfernt | Runner: `"$docker" volume rm "$daten" >/dev/null 2>&1 \|\| true` → `:` (`p7-volume`) | `p7-volume`, `p7-ende-rot`, `p7-ende-bereit`, `p7-ende-pull`, `p7-ende-gruen`, `p7-vor-naechster`, `p7-sigINT-gruppe`, `p7-sigTERM-gruppe`, `p7-sigINT-skript`, `p7-sigTERM-skript` |
| 7 | Aufräumen nach einem roten Lauf (Phase 1 rot) | Runner: `  exit "$status"` → `  kill -KILL $$` (`p7-ende-rot`) | `p7-ende-rot` |
| 7 | Aufräumen nach einem Fehler des Aufbaus: Server nicht bereit | Runner: `exit 1 # E-bereit` → `kill -KILL $$ # E-bereit` (`p7-ende-bereit`) | `p7-ende-bereit` |
| 7 | Aufräumen nach einem Fehler des Aufbaus: Image nicht zu ziehen | Runner: `"$postgres_image" >/dev/null` → `"$postgres_image" >/dev/null \|\| kill -KILL $$` (`p7-ende-pull`) | `p7-ende-pull` |
| 7 | Aufräumen nach grünem Ende | Runner: `echo "run-integration-tests: gruen"` → `kill -KILL $$` (`p7-ende-gruen`) | `p7-ende-gruen`, `p7-container`, `p7-netz`, `p7-volume`, `p7-vor-naechster` |
| 7 | Aufräumen nach SIGINT | Treiber: `trap 'beenden 130' INT # S-INT` → `trap 'kill -KILL \$\$' INT # S-INT` (`p7-int-treiber`) | `p7-sigINT-gruppe`, `p7-sigINT-skript` |
| 7 | Aufräumen nach SIGTERM (Treiber) | Treiber: `trap 'beenden 143' TERM # S-TERM` → `trap 'kill -KILL \$\$' TERM # S-TERM` (`p7-term-treiber`) | `p7-sigTERM-gruppe`, `p7-sigTERM-skript` |
| 7 | Aufräumen nach SIGTERM (Runner) | Runner: `trap 'trap "" TERM; exit 143' TERM # S-TERM` → `trap 'kill -KILL \$\$' TERM # S-TERM` (`p7-term-runner`) | `p7-sigINT-gruppe`, `p7-sigINT-skript`, `p7-sigTERM-gruppe`, `p7-sigTERM-skript` |
| 8 | Referenz ohne Auswahl: volle Suite in Phase 1 | Runner: `wahl_1=()` → `wahl_1=(-test.run '^TestE2EOhnePostgres')` (`p8-phase1`) | `p8-phase1` |
| 8 | Referenz ohne Auswahl: in Phase 2 `^TestE2EOhnePostgres` | Runner: `wahl_2=(-test.run '^TestE2EOhnePostgres')` → `wahl_2=(-test.run '^TestE2E')` (`p8-phase2`) | `p8-phase2` |
| 8 | Referenz ohne Auswahl legt alle vier Benutzer an | Runner: `!= scram ]; then # P5-rollen` → `!= scram ] && false; then # P5-rollen` (`p8-rollen`) | `p8-rollen` |
| 8 | Referenz ohne Auswahl prüft keine erwarteten Tests | Runner: `if [ "${#tests_1[@]}" -gt 0 ]; then erwartet` → `if true; then erwartet` (`p8-pruefung`) | `p8-ohne-pruefung`, `p8-phase2` |
| 9 | ein Rezept mit `$(MAKE)` fährt unter `-n`, `-q`, `-t`, `MAKEFLAGS=n` docker (Mutation: Rezept mit `$(MAKE)`) | Fragment: `@POSTGRES_VERSIONEN=` → `@$(MAKE) --version >/dev/null; POSTGRES_VERSIONEN=` (`p9-make`) | `p9-n`, `p9-q`, `p9-t`, `p9-makeflags-n`, `p9-n-und-i` |
| 9 | ein Rezept mit `+` fährt unter `-n`, `-q`, `-t`, `MAKEFLAGS=n` docker (Mutation: Rezept mit `+`) | Fragment: `$\t@POSTGRES_VERSIONEN=` → `$\t+@POSTGRES_VERSIONEN=` (`p9-plus`) | `p9-n`, `p9-q`, `p9-t`, `p9-makeflags-n`, `p9-n-und-i` |
| 9 | das Ziel schreibt keinen Stand | Fragment: `bash tools/test/run-postgres-versionen.sh` → `bash tools/test/run-postgres-versionen.sh && touch .harne…` (`p9-stand`) | `p9-kein-stand` |
| 12 | `test-postgres-versionen` hängt an `GATE_CHECKS` | Fragment: `GATE_CHECKS += test-postgres-versionen ` → `GATE_CHECKS += ` (`p12-ziel`) | `p12-gate-checks-ziel` |
| 12 | `test-postgres-versionen-gegenprobe` hängt an `GATE_CHECKS` | Fragment: `GATE_CHECKS += test-postgres-versionen test-postgres-vers…` → `GATE_CHECKS += test-postgres-versionen` (`p12-gegenprobe`) | `p12-gate-checks-gegenprobe` |
| 3 | Test `TestE2ERecordSelect1` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2ERecordSelect1` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2ERecordSelect1` |
| 3 | Test `TestE2EVorbereitungOhnePostgres` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2EVorbereitungOhnePostgres` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2EVorbereitungOhnePostgres` |
| 3 | Test `TestE2EVorbereitungExtendedOhnePostgres` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2EVorbereitungExtendedOhnePostgres` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2EVorbereitungExtendedOhnePostgres` |
| 3 | Test `TestE2ERecordExtendedPgx` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2ERecordExtendedPgx` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2ERecordExtendedPgx` |
| 3 | Test `TestE2ERecordAnmeldung/scram-sha-256` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2ERecordAnmeldung/scram-sha-256` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2ERecordAnmeldung/scram-sha-256` |
| 3 | Test `TestE2EReplayLebendpruefungWiePostgres` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2EReplayLebendpruefungWiePostgres` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2EReplayLebendpruefungWiePostgres` |
| 3 | Test `TestE2EOhnePostgresReplay` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2EOhnePostgresReplay` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2EOhnePostgresReplay` |
| 3 | Test `TestE2EOhnePostgresExtendedReplay` der Teilmenge: ohne PASS-Zeile ist die Version rot, die Auswahl nennt ihn | Fragment: `TestE2EOhnePostgresExtendedReplay` aus der Liste gestrichen (`p3-element-…`) | `p3-element-TestE2EOhnePostgresExtendedReplay` |

Mutation und Fall der Aufzählung in Punkt 7 („Ende“): grün (`p7-ende-gruen`), rot (`p7-ende-rot`), Aufbau (`p7-ende-bereit`, `p7-ende-pull`), `SIGINT` und `SIGTERM` je an die Prozessgruppe und an das Skript (`p7-sigINT-gruppe`, `p7-sigINT-skript`, `p7-sigTERM-gruppe`, `p7-sigTERM-skript`); die Mutationen der Zeilen `p7-*` oben. Wer das Signal nicht aufräumt, färbt je nach Pfad `p7-sigINT-*` und `p7-sigTERM-*` rot (`p7-int-treiber`, `p7-term-treiber`, `p7-term-runner`).

**Gegen echte Server gesehen** (kein Teil der Gegenprobe; `make test-postgres-versionen` und Mutation in der Arbeitskopie, danach zurückgesetzt, `touch` nach jeder Änderung, damit der Build-Kontext sie überträgt):

| Zusage (Punkt) | Mutation | Ergebnis |
|---|---|---|
| 11, Grenze falsch gegen eine Version zwischen ihr und 17: `\v` | `test/integration/lebendpruefung_e2e_test.go`: `vtAbVersion = 17` → `19` | 14, 15, 16 grün; 18 rot in Phase 1: „PostgreSQL 18.6 beantwortet `\v` mit … die leere Anfrage …; `\v` Leerraum laut LH-FA-09.a: false“, Meldung `ROT — 18 (Phase 1)` |
| 11, Grenze zu niedrig: eine Version unter 17 | `vtAbVersion = 17` → `16` | 14, 15 und 18 grün; 16 rot in Phase 1: „PostgreSQL 16.15 beantwortet `\v` mit fehler 42601 …; `\v` Leerraum laut LH-FA-09.a: true“ |
| 11, Wert ohne führende Ziffer macht den Test rot | `hauptversion`: `v := "x" + conn.ParameterStatus("server_version")` | alle vier rot: `server_version "x14.24": strconv.Atoi: parsing "": invalid syntax` (die Mutation fasst nur die Eingabe an; ein Mutant, der den `Fatalf` entfernt, ist nicht gefahren) |
| 11, Basis ohne Mutation | — | 14, 15, 16, 18 grün; `vtAbVersion = 17` stimmt für 14.24, 15.19, 16.15 und 18.6 (keine Berichtigung nötig, keine Rückführung §4 (b) oder (c)) |
| 4, Digest, der zu einem anderen Tag gehört | Liste `17=… 15=postgres:15-alpine@sha256:<Digest von 16>` | `Version 15: der Server meldet server_version '16.15'`, `ROT — 15 (Aufbau)`, keine Reste |
| 4, Image, das sich nicht ziehen lässt | `15=postgres:15-alpine@sha256:000…0` | `manifest unknown`, `ROT — 15 (Aufbau)`, 18 danach grün, keine Reste |
| 3, Test fehlt (umbenannt) | `POSTGRES_TESTS_PHASE1="TestE2ERecordSelect1 TestE2ENichtDa"` | `Version 18, Phase 1: Test TestE2ENichtDa ohne PASS-Zeile`, `ROT — 18 (Phase 1)` |
| 3, Auswahl ohne Treffer | `POSTGRES_TESTS_PHASE1="TestE2ENichtDa"` | `Auswahl ohne Treffer, kein Test bestanden`, `ROT — 18 (Phase 1)`; `go test` selbst endet dabei grün |
| 7, `SIGINT` und `SIGTERM` an die Prozessgruppe von `make test-postgres-versionen` und an `tools/test/run-postgres-versionen.sh`, sobald der Testcontainer läuft | — | Exit 130 und 143; nach dem Ende der Prozesse 0 Container, 0 Netze, 0 Volumes. Über `make`: Bei `SIGTERM` an die Gruppe endet `make` vor dem Skript; die Reste verschwinden Sekunden später (siehe Rückgabe 5) |

Nicht gefahren, weil ohne Mehrwert für den Slice: ein Lauf des vollen `make gates` gegen eine Version mit abweichender Grenze (keine Abweichung gefunden), ein roter Bau des Testimages gegen echtes Docker (nur Attrappe, `p6-build-rot`).

**Namen der Nahtstelle** (für `slice-v1-abschluss-postgres-versionen-zusage`, §1 dort; der Planner trägt sie bei der Closure ein):

1. Datei `harness/mk/postgres-versionen.mk`; Variablen `POSTGRES_REFERENZ` (`17`) und `POSTGRES_VERSIONEN` (Einträge `<Hauptversion>=<Image>`: 14, 15, 16, 17, 18); Referenz-Image gleich `POSTGRES_IMAGE` aus `harness/mk/integration.mk`.
2. Teilmenge im Wortlaut (Variablen `POSTGRES_TESTS_PHASE1`, `POSTGRES_TESTS_PHASE2`): Phase 1 `TestE2ERecordSelect1` (Szenario 1), `TestE2EVorbereitungOhnePostgres` (Szenario 2, Vorbereitung), `TestE2EVorbereitungExtendedOhnePostgres` und `TestE2ERecordExtendedPgx` (Szenario 7), `TestE2ERecordAnmeldung/scram-sha-256` (Fall SCRAM-SHA-256), `TestE2EReplayLebendpruefungWiePostgres`; Phase 2 `TestE2EOhnePostgresReplay` (Szenario 2), `TestE2EOhnePostgresExtendedReplay` (Szenario 7).
3. Berichtigung an der Versionsgrenze in `LH-FA-09.a`: keine; die Grenze 17 gilt für 14, 15, 16 und 18 (Lauf und Mutationen oben).
4. Nachfahrbarer Lauf je Version, hier 16: `make test-postgres-versionen POSTGRES_VERSIONEN='16=postgres:16-alpine@sha256:721873c34ceb9f8d8fc265984940dc982404c105f19ad51be9fdc5970a6080ea 17=postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'` (der Eintrag der Referenz gehört zur Liste, das Ziel fährt ihn nicht; gelaufen, 6 s). Für eine andere Version den Eintrag aus dem Fragment einsetzen. Gegen die Referenz selbst: `make test-integration`.

**Rückgabe an den Architect** (Randform · Frage; `AGENTS.md` §3.12). **Alle sechs sind entschieden** (Architect, 2026-10-14, Commit f0a2280; `SPEC-053` Punkt 3, 4, 6, 7, 8, 14) und in der Gegenprobe umgesetzt (8f97a05): Aufbau als dritte Stufe, jeder Fall der Formen aus Punkt 4 und 5 prüft sie (`p4-nicht-bereit`, `p4-60-sekunden`, `p4-start`, `p4-version-abweichend`, `p4-version-unlesbar`, `p4-version-ohne-ziffer`, `p5-einrichtung-scheitert`, je eine Mutation „Stufe fällt auf Phase 1“ in den Zeilen `p4-aufbau-*` und `p5-aufbau-einrichtung`); `p6-build-rot` prüft Exit-Code, Meldung und dass der Bau der einzige Aufruf war; `p3-fehlt-phase2` für Rückgabe 3; Rückgabe 4 und 5 sind akzeptierte Negative (Punkt 8, 7 und 14), Rückgabe 6 ist auf „starten“ enger gefasst (`p4-start`). Wortlaut der Fragen, wie sie gestellt waren:

1. *Stufe „Aufbau“.* Ein Fehler vor der ersten Testphase (Image nicht zu ziehen oder zu starten, Server nicht bereit, Hauptversion abweichend oder nicht lesbar, Einrichtung der Anmeldung scheitert) steht in der Zusammenfassung als `<Version> (Aufbau)`. `SPEC-053` Punkt 6 nennt nur die Phase 1 und 2. Soll „Aufbau“ so heißen, oder fällt der Fehler unter Phase 1?
2. *Roter Bau des Testimages.* Er beendet das Ziel sofort mit `ROT — das Testimage ist nicht zu bauen` und Exit 1, ohne Version zu starten. Punkt 6 sagt „Das Testimage wird einmal gebaut“, nicht, was ein roter Bau tut.
3. *Rote Phase 2 und Prüfung der erwarteten Tests in Phase 2.* Eine rote Phase 2 gilt als „Phase 2“; fehlt nur dort ein erwarteter Test, ebenso. Punkt 3 nennt die Meldung „Test, Version und Phase“; die Zuordnung der Prüfung zur Phase ist Lesart.
4. *Zusammengeführte Ausgabe des Runners.* Der Runner leitet die Ausgabe von `docker run` des Testlaufs jetzt über `2>&1 | tee` in eine Datei (Voraussetzung der Prüfung); die Ausgabe der Tests auf stderr erscheint damit auch bei `make test-integration` auf stdout. Der Exit-Code und die Tests sind unverändert (Punkt 8 sagt „wie vor diesem Vertrag“).
5. *Signal über `make`.* Punkt 7 sagt „nach `SIGINT` und `SIGTERM`“. Das Skript räumt vor seinem Ende auf (Fälle `p7-sig*`, echt gesehen). Bei `SIGTERM` an die Prozessgruppe von `make` endet `make` vor dem Skript, das danach Sekunden weiter aufräumt; `SIGINT` an die Gruppe endete mit Exit 130 ohne Reste; `SIGINT` nur an `make` ließ `make` das Skript abwarten, und alle vier Versionen liefen durch (Exit 130). Reicht „das Skript räumt auf“, oder soll ein Fall den Zeitpunkt nach dem Ende von `make` prüfen?
6. *Ein Fall für „ziehen oder starten“.* Punkt 4 nennt „nicht ziehen oder starten“; die Attrappe lässt `docker run -d` scheitern und kann beides nicht trennen (`p4-pull`). Genügt das?

**Behauptung zur DoD** (Bestätigung bleibt beim Verifier): Punkt 1 der DoD: Bereich 14 bis 18 bestätigt (oben), Liste in `harness/mk/postgres-versionen.mk`, Ziel hängt an `GATE_CHECKS` (`p12-*`), fährt 1, 2, 7, den Fall SCRAM-SHA-256 und den Lebendprüfungs-Abgleich gegen 14, 15, 16 und 18, je per Digest (`p2g`, `p2h-*`), `harness/README.md` §Sensors führt beide Ziele, die Ausnahmeliste und die Gegenprobe der schnellen Kette kennen sie (Mutanten `m-test-postgres-versionen*`), gebrochen und rot gesehen: Image ohne Digest, fehlende Version (Referenz fehlt, Hauptversion nicht aus Ziffern, nicht ziehbares Image), Auswahl ohne Treffer. Punkt 2: `TestE2EReplayLebendpruefungWiePostgres` grün gegen 14, 15, 16 und 18; die Grenze 17 stimmt, keine Rückführung. Punkt 3: `make gates` grün (561 s, Nachweis geschrieben). Die Nahtstelle (Punkt 9 der DoD) trägt der Planner bei der Closure ein; die Namen stehen oben.

**Hinweise für die Closure (Ausgänge, Urteil beim Planner).** Versionsbereich: Belege oben, offen bleibt die Pflege von Hand, wenn 19 erscheint oder 14 endet (2026-11-12). Versionsgrenze für `\v`: bestätigt für 14, 15, 16, 18 (Lauf, Mutationen), 19 ungeprüft. Teilmenge fällt leer: Runner und Treiber verlangen je Test eine PASS-Zeile, drei Formen rot gesehen (Attrappe und echt). Laufzeit: oben, Zuwachs rund 105 s. Anmeldung je Image: gleiche Einrichtung wirkte auf 14, 15, 16, 18 (SCRAM-Fall grün). Zweite Phase verdeckt: erbt der Treiber, die Meldung nennt die Phase (`p6-meldung-phase1`). Namen der Nahtstelle: oben. Abweichung hebt die Schichten: nicht eingetreten. Beobachtung für das Register: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` trat nicht wieder auf (die Prüfung der PASS-Zeilen schloss die Klasse); `BEO-REPO/spec-randform-erst-im-review-entschieden` berührt die Rückgabe oben.

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
