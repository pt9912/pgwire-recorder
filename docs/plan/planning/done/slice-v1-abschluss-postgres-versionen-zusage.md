# Slice slice-v1-abschluss-postgres-versionen-zusage: Zusage der unterstützten PostgreSQL-Versionen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)

**Berührte Spec-Stellen:** `LH-FA-05.b` · `LH-FA-05.e` · `LH-FA-09.a` · `SPEC-029` · `SPEC-038` · `LH-RB-02`

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

**Ziel:** Der Recorder sagt zu, mit welchen PostgreSQL-Hauptversionen er aufzeichnet, und die Zusage reicht genau so weit wie das Gate aus `slice-v1-abschluss-postgres-versionen-matrix`; Handbuch und README nennen den Bereich und, wo sie ein Verhalten des Servers beschreiben, die Version mit einem nachfahrbaren Lauf.

**Übernimmt:** `slice-v1-abschluss-postgres-versionen` — dessen Teil *Zusage* (Entscheidung des Nutzers vom 2026-10-14, Schnitt nach Lieferwert, Matrix zuerst; dort §7 `Gegenstand:`). Der andere Teil, das Gate, liegt in `slice-v1-abschluss-postgres-versionen-matrix`. Dazu die Sendung der Handbuch-Sätze zur Anmeldung (unten).

**Herkunft:** Die Spezifikation sagt in `LH-FA-05.e` und `SPEC-029` nur „jeder Server, der PGWire 3.0 spricht“, und das Lastenheft nennt keine Version. Entscheidungen des Nutzers (2026-10-05): Unterstützt werden die **fünf jüngsten erschienenen Hauptversionen** als feste Liste; die Spezifikation führt sie, ein späterer Slice hebt sie an; das Lastenheft bekommt einen Satz in [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung). Der Slice steht hinter der Matrix, weil eine Zusage nur sagt, was das Gate belegt (`AGENTS.md` §3.11). Die Liste bestätigt die Matrix beim Start gegen die Release-Liste von PostgreSQL; dieser Slice nennt sie so, wie die Matrix sie führt, und bestätigt sie nicht ein zweites Mal.

- **Zusage in der Spezifikation:** `SPEC-029` und `LH-FA-05.e` nennen die Liste der Hauptversionen als feste Liste und was „unterstützt“ heißt, und ersetzen „jeder Server, der PGWire 3.0 spricht“. „Unterstützt“ heißt genau das, was die Matrix fährt: Gegen jede Version der Liste laufen die Abnahmeszenarien 1, 2 und 7, die Anmeldung mit SCRAM-SHA-256 und der Lebendprüfungs-Abgleich grün; die volle Integrationssuite läuft gegen die Referenzversion 17. Abnahmeszenario 18 und Szenario 16 Satz 3 (nur gegen 17, Entscheidung des Nutzers vom 2026-10-12) nennt die Zusage heute nicht: Kein Test prüft sie, ihre Slices liegen in `next/`; jeder der beiden Slices nimmt den Satz mit seinem Test auf (Sendung unten); der Satz „die automatisierten Abnahmeszenarien laufen gegen jede von ihnen grün“ aus dem Ursprungsplan reicht weiter als das Gate und gilt hier nicht (`AGENTS.md` §3.11). `SPEC-038` (Kompatibilitätstests) nennt die Versionsmatrix. [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung) bekommt einen Satz: Der Record-Modus setzt einen PostgreSQL-Server einer unterstützten Hauptversion voraus, und welche das sind, nennt die Spezifikation; das Lastenheft steht auf Draft, kein Versionssprung.
- **Handbuch und README:** Das Benutzerhandbuch nennt den Bereich (Abschnitt *Voraussetzungen*); Handbuch und README beschreiben, was die Matrix liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11).
- **Sendung — Handbuch-Sätze zur Anmeldung** (Entscheidung des Nutzers vom 2026-10-14, Geber `slice-v1-abschluss-anmeldung` und `slice-v1-abschluss-anmeldung-doku`, dort je §6, Risiko zur Version des Servers; `AGENTS.md` §3.13): Die Sätze „geprüft mit PostgreSQL 17“ zur Anmeldung aus `slice-v1-abschluss-anmeldung` und `slice-v1-abschluss-anmeldung-doku` sind Teil dieser Sendung und werden je Verfahren nachgezogen (SCRAM-SHA-256: Bereich 14 bis 18; Klartext und MD5: 17 bleibt), an diesen Stellen: `docs/user/benutzerhandbuch.md` §1 *Voraussetzungen*, §4 der Hinweis in *Eine Anwendung aufzeichnen*, der Abschnitt *Eine Anwendung mit Passwort-Anmeldung aufzeichnen* (Absatz *Voraussetzung*) und §6 *Rollen und Rechte*; dazu der Satz zur Anmeldung in `README.md`, der denselben Wortlaut trägt (Fund des Planners, nicht in der Entscheidung genannt). Ein Satz nennt den Bereich nur für die Verfahren und Versionen, für die §7 des Implementers einen Lauf führt, den ein Verifier nachfahren kann (`.claude/commands/implement-slice.md` Schritt 19, Entscheidung des Nutzers vom 2026-10-14); die Matrix fährt allein den Fall SCRAM-SHA-256 gegen die Versionen außer 17, Klartext und MD5 fährt sie nicht. Reicht für ein Verfahren kein Lauf, nennt der Satz die Versionen mit Lauf und sagt zu den übrigen nichts. **Entschieden (Nutzer, 2026-10-15, Frage B):** Die Sätze werden auf das begrenzt, was das Gate belegt. SCRAM-SHA-256 gilt für PostgreSQL 14 bis 18 (Läufe: Matrix gegen 14, 15, 16 und 18, `make test-integration` gegen die Referenz 17); Klartext-Passwort und MD5 bleiben bei „PostgreSQL 17“ (Lauf: `make test-integration`); kein Lauf von Hand je Version. Das Handbuch und der Satz in `README.md` nennen für Klartext und MD5 keinen Bereich. Die Regel in Schritt 19 gilt seit `slice-v1-abschluss-anmeldung-doku` (Commit `c41c25d`); dieser Slice ändert `implement-slice.md` nicht, er wendet sie an.

**Nahtstelle von `slice-v1-abschluss-postgres-versionen-matrix`** (eingetragen vom Planner bei deren Closure am 2026-10-14, mit Nachzählung in §8; geprüft mit `grep -n "slice-v1-abschluss-postgres-versionen-matrix"` in dieser Datei). Die Matrix ist geschlossen; ihre Namen gelten:

1. *Liste und Referenz.* Datei `harness/mk/postgres-versionen.mk`, Variablen `POSTGRES_REFERENZ` (`17`) und `POSTGRES_VERSIONEN` (Einträge `<Hauptversion>=<Image>`, je per Digest gepinnt). Die Liste ist **14, 15, 16, 17 und 18**, die Referenz 17; bestätigt am 2026-10-14 gegen die Release-Liste von PostgreSQL (19 ist nicht erschienen, Beta). Die Matrix fährt 14, 15, 16 und 18; 17 fährt `make test-integration` mit der vollen Suite. Das Image der Referenz ist `POSTGRES_IMAGE` aus `harness/mk/integration.mk`. PostgreSQL 14 endet am 2026-11-12; die Liste hebt sich nicht selbst an, ihre Pflege ist eine Handlung von Hand.
2. *Teilmenge im Wortlaut*, auf die „unterstützt“ zeigt (Variablen `POSTGRES_TESTS_PHASE1` und `POSTGRES_TESTS_PHASE2`): Phase 1 mit PostgreSQL: `TestE2ERecordSelect1` (Abnahmeszenario 1), `TestE2EVorbereitungOhnePostgres` (Szenario 2, Vorbereitung), `TestE2EVorbereitungExtendedOhnePostgres` und `TestE2ERecordExtendedPgx` (Szenario 7), `TestE2ERecordAnmeldung/scram-sha-256` (der Fall SCRAM-SHA-256 der Anmeldung), `TestE2EReplayLebendpruefungWiePostgres` (Lebendprüfungs-Abgleich); Phase 2 ohne PostgreSQL: `TestE2EOhnePostgresReplay` (Szenario 2), `TestE2EOhnePostgresExtendedReplay` (Szenario 7). Klartext-Passwort und MD5 fährt die Matrix nicht.
3. *Berichtigung an der Versionsgrenze in `LH-FA-09.a`:* keine. Die Grenze 17 für `\v` gilt für 14, 15, 16 und 18 (Lauf gegen jede Version, die Mutation der Grenze färbt die Version zwischen ihr und 17 rot); für 19 ist sie ungeprüft. Die Zusage ändert `LH-FA-09.a` nicht.
4. *Aufruf je Version:* `make test-postgres-versionen POSTGRES_VERSIONEN='<N>=<Image der Version> 17=<Image der Referenz>'`. Der Eintrag der Referenz gehört zur Liste, das Ziel fährt ihn nicht; die Einträge stehen im Fragment. Gelaufen ist der Aufruf für 16 (6 s): `make test-postgres-versionen POSTGRES_VERSIONEN='16=postgres:16-alpine@sha256:721873c34ceb9f8d8fc265984940dc982404c105f19ad51be9fdc5970a6080ea 17=postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24'`; für 14, 15 und 18 ist derselbe Aufruf mit dem Eintrag der Version aus dem Fragment gelaufen (§7). Aus dem Anmeldetest fährt dieser Lauf nur den Fall SCRAM-SHA-256; Klartext und MD5 laufen allein in `make test-integration` gegen 17.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Gate, die Liste mit Tags und Digests, der Lauf gegen eine Version und jede Berichtigung der Versionsgrenze der Lebendprüfung — `slice-v1-abschluss-postgres-versionen-matrix`: Dieser Slice sagt zu, was jener belegt; ändert er Gate oder Produktcode, hat er die Zusage vor ihren Beleg gestellt.
- Versionen vor dem Bereich — die Entscheidung des Nutzers zieht die Grenze bei den fünf jüngsten Hauptversionen; eine ältere Version bleibt ungeprüft und nicht zugesagt, auch wenn sie PGWire 3.0 spricht. Die Zusage ersetzt „jeder Server, der PGWire 3.0 spricht“, sie ergänzt ihn nicht; eine Prüfung der Version des Servers im Recorder liefert dieser Slice nicht (Produkt-Code unberührt).
- Server, die PGWire 3.0 sprechen, aber kein PostgreSQL der Liste sind (Pooler, andere Datenbanken) — nicht zugesagt; keine Stelle sagt zu, dass der Recorder sie ablehnt oder annimmt (kein Test, `AGENTS.md` §3.11); das Handbuch nennt keine Aussage zu ihnen.
- Eine ADR zur Zusage — die Entscheidung des Nutzers steht in der Spezifikation (`AGENTS.md` §3.12), und die ADR zur Matrix liegt bei `slice-v1-abschluss-postgres-versionen-matrix`. Verlangt der Architect eine ADR, hebt sie die Schichten über zwei (§8) und geht an den Planner.
- Handbuch-Sätze über ein Gegenüber außerhalb der Anmeldung (Client-Versionen, `psql` aus libpq 17, `slice-v1-abschluss-tls-client-doku`) — Bestand: Sie nennen ihre Version bereits (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, Gegenbelege); die Sendung betrifft die Version des Servers.
- Windows und macOS sowie andere Treiber — kein Lauf; das Handbuch sagt dazu nichts (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 3×, offen).
- Ein Sensor, der den Bereich gegen die Release-Liste von PostgreSQL hält — bräuchte Netz zur Laufzeit des Gates; die Pflege ist eine Handlung (§6), kein Lauf.
- Kein Produkt-Code, kein Gate, keine Hilfetexte im Binary — Schicht-Abgrenzung: Der Slice ändert Spezifikation, Lastenheft, `docs/user/benutzerhandbuch.md`, `README.md` und diesen Plan.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung): `SPEC-029` und `LH-FA-05.e` nennen die von der Matrix geführte Liste der Hauptversionen als feste Liste und was „unterstützt“ heißt (genau die Teilmenge der Matrix; die volle Suite gegen 17); `SPEC-038` nennt die Versionsmatrix; [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung) trägt den Satz zur unterstützten Hauptversion; die Historie der Spezifikation nennt die Änderung. Jede Aussage der Zusage ist durch einen Lauf der Matrix belegt (§6).
- [x] `docs/user/benutzerhandbuch.md`: nennt den Bereich; die Sätze zur Anmeldung an den vier Stellen aus §1 nennen für SCRAM-SHA-256 PostgreSQL 14 bis 18 und für Klartext und MD5 PostgreSQL 17 (Entscheidung des Nutzers vom 2026-10-15), und §7 führt zu jeder genannten Version und jedem Verfahren einen Lauf, den ein Verifier nachfahren kann; im Ist-Zustand, ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews.
- [x] `README.md`: der Satz zur Anmeldung nennt je Verfahren dieselben Versionen wie das Handbuch (SCRAM-SHA-256: 14 bis 18; Klartext und MD5: 17); im Ist-Zustand.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` (`LH-FA-05.e`, `SPEC-029`, `SPEC-038`, Historie) | update | Versionsbereich als Liste, Bedeutung von „unterstützt“, Versionsmatrix der Kompatibilitätstests |
| `spec/lastenheft.md` ([`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung)) | update | ein Satz zur unterstützten Hauptversion des Servers (Entscheidung des Nutzers, 2026-10-05) und die Zeile *Nachweis* mit der Teilmenge der Kompatibilitätstests je Hauptversion; Draft, kein Versionssprung |
| `docs/user/benutzerhandbuch.md` | update | Bereich in *Voraussetzungen* (eigener Punkt); in der Voraussetzung der Anmeldung der Client der Läufe (pgx 5.11; `psql` aus libpq 17 gegen 17); die vier Stellen zur Anmeldung (SCRAM-SHA-256: PostgreSQL 14 bis 18; Klartext und MD5: „PostgreSQL 17“, Entscheidung des Nutzers vom 2026-10-15); Läufe der Matrix und der Referenz als Beleg in §7 |
| `README.md` | update | Satz zur Anmeldung (Ist-Zustand) |
| `docs/user/abdeckung-*.md` | update (bedingt) | nur wenn Deklarationen sich ändern; dieser Slice schreibt keinen Test |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-postgres-versionen-matrix` liegt in `done/`, die bestätigte Liste und die Namen der Nahtstelle stehen in §1 (eingetragen vom Planner bei der Closure der Matrix, 2026-10-14), und `slice-v1-abschluss-anmeldung-doku` liegt in `done/` (erfüllt). Der Slice steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) als Schritt 32b direkt hinter der Matrix.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): (a) **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- spec docs/user README.md`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 400** geht der Slice zurück. **Meldeschwelle** neben der Marke: Erreicht die Summe **300**, meldet der Implementer den Stand dem Architect, bevor er weiterarbeitet; über Schnitt oder Anhebung der Marke entscheidet der Planner. **Schätzung** (Planner, 2026-10-14, in dieser Metrik; der Architect schätzt vor dem Code neu): kein Produktcode, kein Test, also kein Faktor; Spezifikation und Lastenheft rund 60 Zeilen, Handbuch rund 80 (Bereich, vier Stellen mit Löschungen und Einfügungen), README rund 10; zusammen rund 150, mit der Review-Runde (20 bis 30 %) **180 bis 250**. (b) Der Architect verlangt eine ADR oder eine Änderung außerhalb von Spezifikation, Lastenheft, Handbuch und README: Die Schichten stiegen über zwei, der Planner schneidet vor dem ersten Code-Commit (§8).
- `in-progress` → `open` (blockiert — Carveout?): Das Gate der Matrix ist rot oder entfällt (die Zusage stünde ohne Beleg), oder der Nutzer zieht die Entscheidung zum Bereich zurück.

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

**Randformen** (`AGENTS.md` §3.12) der geänderten Zusage. Der Slice legt keinen Vertrag eines Werkzeugs an, sondern ändert den Ort der Zusage (`SPEC-029`); entschieden werden sie vom Architect vor dem ersten Code-Commit in der Spezifikation (Technik-Stratum), soweit nicht schon entschieden:

- *Lesart des Bereichs:* fünf jüngste erschienene Hauptversionen als feste Liste (Nutzer, 2026-10-05, übernommen); die Liste kommt aus der Matrix.
- *Szenarien nur gegen 17:* Abnahmeszenario 18 und Szenario 16 Satz 3 laufen nur gegen die Referenz 17 (Nutzer, 2026-10-12, übernommen). **Entschieden (Architect, 2026-10-15, Rückgabe V-1 der Verifikation):** Die Zusage nennt sie **nicht**, solange kein Test sie prüft (`AGENTS.md` §3.11); ihre Slices liegen in `next/`. Sendung (`AGENTS.md` §3.13, Geber dieser Slice): `slice-v1-abschluss-verbindungsende-upstream` nimmt Szenario 18, `slice-v1-abschluss-verbindungsende-vergleich` Szenario 16 Satz 3 in `SPEC-029` Punkt 2 auf, jeder im Commit, der seinen Integrationstest liefert; beide Nehmer tragen es in §1 und DoD mit der Kennung dieses Gebers, Kopf und §8 sind nachgezogen (beide nehmen an, kein Ausschluss trifft die Sendung; Zählung je Nehmer 3 Liefer-Punkte, 2 Schichten, unverändert). Kein Zielstand-Satz im Ist-Zustand.
- *Nebenversion und Vorabversion:* die Zusage gilt für die Hauptversion; das Gate fährt je Hauptversion eine gepinnte Nebenversion, Vorabversionen sind nicht zugesagt. Entscheidung des Architect, Wortlaut in `SPEC-029`.
- *Server außerhalb der Liste:* nicht zugesagt und nicht geprüft (§1). „Nicht abgelehnt“ und „der Recorder prüft die Version nicht“ sagt keine Stelle zu, weil kein Test es prüft (`AGENTS.md` §3.11); `LH-FA-05.e` verliert den Satz „jeder PostgreSQL-Server zulässig, der PGWire 3.0 spricht“ ersatzlos, die Behandlung abweichender Protokollversionen (`PGR-E6002`) bleibt. Entscheidung des Architect (2026-10-15).

**Entscheidungen des Architect vor dem ersten Code-Commit (2026-10-15).** Die Stellen sind geschrieben (Commit `9458483`, Nachweis-Zeile von `LH-RB-02` in `1bc7dc7`) nach diesen Vorgaben. Eine Folge-ADR ist nicht nötig: Es wird nichts Angenommenes berührt ([ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) trägt das Gate; die Zusage ist Spezifikation).

- *Wortlaut von `SPEC-029` (Tabelle §6, Spalte Version):* nennt die Hauptversionen 14, 15, 16, 17 und 18 als Liste, mit 17 als Referenz, **nicht** die Lesart „fünf jüngste“. Grund: Erscheint 19, bliebe eine Regel „fünf jüngste“ als Satz stehen, den die Liste dann verletzt; die Liste allein bleibt wahr, bis jemand sie ändert. Die Spalte „Vertrag-Datei“ bleibt `—`.
- *Nennung der Matrix (Referenz-Richtung):* Spezifikation und Lastenheft nennen die Matrix als **Text** („Versionsmatrix der Kompatibilitätstests“) und geben kein Gate, kein `make`-Ziel, keinen Dateinamen und **keinen Link** auf `harness/vertraege/` an; der Vertrag des Gates steht dort, und die Spezifikation zeigt nicht nach `aussen`. Die Teilmenge steht in Worten (Aufzeichnen einer Abfrage, Abspielen ohne Server, erweitertes Protokoll, Anmeldung mit SCRAM-SHA-256, Lebendprüfung), nicht als Testnamen.
- *„Unterstützt“ heißt genau:* (1) Gegen jede Hauptversion der Liste sind die Abnahmeszenarien 1, 2 und 7, die Anmeldung mit SCRAM-SHA-256 und der Lebendprüfungs-Abgleich geprüft. (2) Gegen die Referenz 17 läuft zusätzlich die volle Integrationssuite; Klartext-Passwort und MD5 sind **nur** dort geprüft. Abnahmeszenario 18 und Szenario 16 Satz 3 stehen erst in Punkt 2, wenn ihr Test läuft (Sendung, Randform „Szenarien nur gegen 17“). (3) Je Hauptversion wird eine feste Nebenversion geprüft; zugesagt ist die Hauptversion, Vorabversionen sind nicht zugesagt. Der Satz „die automatisierten Abnahmeszenarien laufen gegen jede von ihnen grün“ entfällt. Die Prüfung „alle übrigen Szenarien“ in (2) ist nicht gelungen; die Texte sagen „volle Integrationssuite“ (§7).
- *`LH-FA-05.e`:* behält Protokollversion 3.0 und `PGR-E6002`; ersetzt den Satz zu „jeder Server“ durch den Verweis auf `SPEC-029` für die zugesagten Hauptversionen.
- *`SPEC-038`:* Absatz *Kompatibilitätstests* ergänzt die Versionsmatrix in einem Satz (Fälle gegen jede Hauptversion aus `SPEC-029`, volle Suite gegen die Referenz), ohne Gate-Namen.
- *`LH-RB-02`:* ein Satz in vorhandener Form, ohne Versionsnummer und ohne Link nach unten: „Zugesagt ist der Record-Modus gegen die Hauptversionen von PostgreSQL, die die Spezifikation nennt.“ Die Zeile *Nachweis* nennt zusätzlich die Kompatibilitätstests je unterstützter Hauptversion (Text); gefasst als „gegen jede unterstützte Hauptversion läuft mindestens die Teilmenge der Kompatibilitätstests, die die Spezifikation nennt“ (Review F-4: „laufen je Hauptversion“ las sich als alle Abnahmeszenarien; Verifikation V-2: „nur“ ist für die Referenz 17 falsch, dort läuft die volle Suite; „mindestens“ deckt beides und stimmt mit `SPEC-029` Punkt 1 und 2 überein; Draft, kein Bump).
- *Version des Lastenhefts:* kein Bump, keine Historie-Zeile im Lastenheft. Status ist Draft; die Bump-Regel (`harness/conventions.md` §Versionierung des Lastenhefts) gilt „ab `Accepted`“. Wäre es Accepted, wäre es **Major**, weil eine bestehende Anforderung geändert wird. Das ist keine Entscheidung des Nutzers, sondern die Anwendung der Regel; der Satz selbst ist Entscheidung des Nutzers vom 2026-10-05. Die Historie der Spezifikation bekommt eine Zeile (`SPEC-029`, `SPEC-038`, `LH-FA-05.e`).
- *Neue Hauptversion (19) und Ende von 14:* Die Liste ist fest und hebt sich nur durch eine Änderung der Spezifikation an. Eine neue Hauptversion ist bis dahin **nicht zugesagt**; das Support-Ende des Herstellers von 14 ändert die Zusage nicht, 14 bleibt, solange die Liste und die Matrix sie führen. Weder Spezifikation noch Handbuch nennen Daten, Ereignisse oder „bis“ (keine Chronik); das Anheben der Liste ändert Spezifikation und Fragment der Matrix im selben Slice. Träger: Risiko „Feste Liste altert“ unten (akzeptiertes Negativ, siehe dort).
- *Handbuch und README:* *Voraussetzungen* nennt den Bereich als Geprüftes („mit PostgreSQL 14 bis 18 geprüft: …“ die Teilmenge in Worten; „mit PostgreSQL 17 zusätzlich: …“ das Übrige), im Ist-Zustand, ohne Szenarionummern, Kennungen, Verweise auf Spezifikation, ADRs, Slices oder Reviews. Die vier Stellen und `README.md` tragen je Verfahren: SCRAM-SHA-256 „PostgreSQL 14 bis 18“, Klartext-Passwort und MD5 „PostgreSQL 17“; für Klartext und MD5 kein Bereich. Sätze zu Clients (`psql`) und zu einem Verhalten des Servers, das eine Version nennt, bleiben unberührt.
- *Zahl der prüfbaren Aussagen und Träger:* acht. (a) Szenario 1, (b) Szenario 2, (c) Szenario 7, (d) SCRAM-SHA-256, (e) Lebendprüfungs-Abgleich je Version 14 bis 18: Lauf der Matrix beziehungsweise `make test-integration` für 17 (Nahtstelle (2) und (4)). (f) volle Suite nur 17, (g) Klartext, MD5 nur 17: `make test-integration`. (h) Hauptversion, eine Nebenversion je Hauptversion: der Pin im Fragment (Nahtstelle (1)). Die Texte in Spezifikation, Lastenheft, Handbuch und README prüft **kein Gate**: `make docs-check` prüft Verweise, nicht die Versionszahlen; `make abdeckung-check` ist unberührt, der Slice schreibt keinen Test und das Handbuch bekommt kein Beispiel als Datei (kein Beispiel mit Ausgabe). Die Deckung Text zu Lauf prüfen Review und Verifier an §7 (Schritt 19). **Akzeptiertes Negativ:** Spezifikationsliste und Fragment der Matrix laufen auseinander, wenn nur eines von Hand geändert wird; ein Sensor wäre ein eigener Vorgang (§1, ohne Netz machbar, aber ohne beobachtbaren Anlass); der Anlass ist das Anheben, das beide im selben Slice ändert. Keine eigene Folgepflicht, solange `BEO-REPO/postgres-versionsliste-altert` offen ist.
- *Aufruf-Voraussetzung und Aufrufarten des Treibers (`AGENTS.md` §3.12):* nicht einschlägig; der Slice liefert kein Werkzeug und keine Ordnungs- oder Eindeutigkeitszusage.
- *Vorgabe der Größe (§3.14):* die Schätzung in §4 bleibt bei 180 bis 250 (kein Code, kein Faktor); Marke und Meldeschwelle gelten unverändert.
- *Sätze des Handbuchs je Verfahren:* **Entschieden (Nutzer, 2026-10-15, Frage B):** Die Sätze nennen nur, was das Gate belegt. SCRAM-SHA-256 gilt für PostgreSQL 14 bis 18 (Matrix 14/15/16/18 plus Referenz 17); Klartext und MD5 laufen allein in `make test-integration` gegen 17 (Nahtstelle (4) in §1) und bleiben bei „PostgreSQL 17“. Kein Lauf von Hand je Version; die frühere Alternative entfällt. §7 führt als Beleg die Aufrufe der Matrix je Version (Nahtstelle (4)) und `make test-integration` für die Referenz.

Risiken:

- **Feste Liste altert.** Die Liste in der Spezifikation hebt sich nicht selbst an; mit jeder neuen Hauptversion fällt sie hinter die Entscheidung „fünf jüngste“ zurück, und kein Gate merkt es (ein Sensor bräuchte Netz, §1). Eine Folge-Adresse gibt es heute nicht: Ein Slice für das Anheben wäre ohne beobachtbaren Anlass geplant. Kandidat für den Ausgang: *weiter offen* → `BEO-REPO/postgres-versionsliste-altert` (angelegt bei der Closure von `slice-v1-abschluss-postgres-versionen-matrix`, 1×; die Closure dieses Slice hängt ihre Datei an), mit dem beobachtbaren Anlass „eine neue Hauptversion von PostgreSQL ist erschienen“; tritt er ein, plant der Planner den Slice, der die Liste anhebt — **Ausgang:** *weiter offen* → `BEO-REPO/postgres-versionsliste-altert` (2×, `evidence/slice-v1-abschluss-postgres-versionen-zusage.md`).
- **Satz im Lastenheft.** Entschieden (Nutzer, 2026-10-05): ein Satz in [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung), der auf die Liste der Spezifikation zeigt, ohne Versionssprung (Draft). Die Formulierung ist entschieden (Architect, 2026-10-15, siehe Randformen): Das Lastenheft nennt die Versionen nicht selbst, sonst müsste jedes Anheben das Lastenheft ändern — **Ausgang:** *entfallen*: Der Satz steht in `LH-RB-02` in der entschiedenen Fassung (Commits `9458483`, `1bc7dc7`, `8ccf3c3`) und nennt keine Versionsnummer; die Verifikation bestätigte ihn, das Risiko einer offenen Formulierung kann nicht mehr eintreten.
- **Zusage reicht weiter als das Gate.** „Unterstützt“ liest sich als Zusage für alles; das Gate fährt die Teilmenge der Matrix (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Jede Aussage der Zusage trägt einen Lauf, sonst wird sie enger gefasst. Eingetreten in V-1 (Szenario 18 und 16 Satz 3 in Punkt 2 ohne Test) und V-2 („nur“ in der Nachweis-Zeile); beide sind enger gefasst (siehe Randformen), die Szenarien gehen als Sendung an `slice-v1-abschluss-verbindungsende-upstream` und `-vergleich` — **Ausgang:** *eingetreten* → Folge-Slices `slice-v1-abschluss-verbindungsende-upstream` und `slice-v1-abschluss-verbindungsende-vergleich`; beide nehmen an (`grep -n "slice-v1-abschluss-postgres-versionen-zusage"` findet die Sendung in §1 und DoD, Nachzählung in §8: je 3 Liefer-Punkte, 2 Schichten).
- **Handbuch-Satz ohne nachfahrbaren Lauf.** Ein Satz nennt eine Version, zu der §7 keinen Lauf führt (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×; Schritt 19). Die Aufrufe je Version kommen aus der Nahtstelle der Matrix (§1, Punkt 4); fehlt einer, nennt der Satz die Version nicht — **Ausgang:** *entfallen*: §7 führt je Version (14, 15, 16, 18 über die Matrix, 17 über `make test-integration`) und je Verfahren einen Aufruf, und die Verifikation fuhr beide nach; kein Satz nennt eine Version ohne Lauf.
- **MD5 gegen die jüngsten Versionen.** MD5 ist ab PostgreSQL 18 abgekündigt; ob der Server der Liste die Anmeldung mit MD5 noch annimmt, zeigt erst der Lauf. Mit der Entscheidung des Nutzers vom 2026-10-15 nennt das Handbuch MD5 nur für 17; **Ausgang:** *entfallen*, weil kein Satz MD5 für eine andere Version zusagt (Handbuch und README nennen MD5 nur mit 17, Lauf `make test-integration`, Fall `md5`; die Verifikation bestätigte es).
- **Namen der Nahtstelle fehlen.** Die Matrix wählte die Namen von Datei, Variablen und Aufruf; sie stehen seit deren Closure in §1 (Punkte 1 bis 4), und `grep -n "slice-v1-abschluss-postgres-versionen-matrix"` und `grep -n "POSTGRES_VERSIONEN"` in dieser Datei finden sie. Bleibt ein Name aus dem Fragment hinter dem Stand der Matrix zurück (ein Rename), bricht der Aufruf in §7 beim Nachfahren — **Ausgang:** *entfallen*: Die Aufrufe in §7 liefen mit den Namen der Nahtstelle, und die Verifikation fuhr sie mit denselben Namen nach (Exit 0).

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

### Belege des Implementers (Stand: Commit `9458483`, `1bc7dc7` plus dieser Plan, 2026-10-15; nachgezogen nach Review-Report vom 2026-10-15)

Alle Läufe nacheinander im Arbeitsbaum, Docker-only; kein `pkill`, keine parallelen Läufe. Der Runner prüft je Version `server_version` gegen die Hauptversion (Mutation unten) und verlangt je Test eine PASS-Zeile.

**Läufe je Version** (Aufruf: `make test-postgres-versionen POSTGRES_VERSIONEN='<N>=<Image> 17=<Image der Referenz>'`, Einträge aus `harness/mk/postgres-versionen.mk`, Nahtstelle (4) in §1):

| Version | Aufruf (Eintrag `<N>=…`) | Ergebnis |
|---|---|---|
| 14 | `14=postgres:14-alpine@sha256:4ea9e5ed…e77a` | exit 0, 6 s, `Version 14 grün`, 9 PASS-Zeilen (Phase 1: sechs Tests mit Fall `scram-sha-256`; Phase 2: zwei) |
| 15 | `15=postgres:15-alpine@sha256:f7d23353…bf81` | exit 0, 6 s, `Version 15 grün`, 9 PASS-Zeilen |
| 16 | `16=postgres:16-alpine@sha256:721873c3…80ea` | exit 0, 6 s, `Version 16 grün`, 9 PASS-Zeilen |
| 18 | `18=postgres:18-alpine@sha256:77f58511…1873` | exit 0, 6 s, `Version 18 grün`, 9 PASS-Zeilen |
| 17 (Referenz) | `make test-integration` | exit 0, 129 s, 127 PASS-Zeilen ohne `FAIL`, darunter `TestE2ERecordAnmeldung/{scram-sha-256,md5,password}` |

Die Digests stehen vollständig im Fragment; gekürzt nur hier.

**Satz · Beleg** (die Texte prüft kein Gate; Deckung Text zu Lauf):

| Satz (Stelle) | Beleg |
|---|---|
| Liste 14, 15, 16, 17, 18, Referenz 17 (`SPEC-029`, Handbuch *Voraussetzungen*) | die fünf Läufe oben; `POSTGRES_REFERENZ`/`POSTGRES_VERSIONEN` im Fragment |
| Gegen jede Hauptversion: Aufzeichnen einer Abfrage (Szenario 1), Abspielen ohne Server (2), erweitertes Protokoll (7), Anmeldung SCRAM-SHA-256, Lebendprüfung-Abgleich | Phase 1 und 2 der Läufe 14, 15, 16, 18 (`POSTGRES_TESTS_PHASE1/2`); für 17 `make test-integration` (volle Suite) |
| Referenz 17: volle Suite; Klartext-Passwort und MD5 nur dort | `make test-integration` (`TestE2ERecordAnmeldung/md5`, `/password`); die Matrix fährt diese nicht (Teilmenge im Fragment); Punkt 14 (5) von `SPEC-053` |
| Je Hauptversion eine feste Nebenversion; Vorabversion nicht zugesagt | Pin per Digest im Fragment; der Runner bricht ab, wenn `server_version` die Hauptversion nicht trägt (Mutation unten) |
| Handbuch/README: SCRAM-SHA-256 mit 14 bis 18 | die Läufe 14, 15, 16, 18 (Fall `scram-sha-256` des Tests `TestE2ERecordAnmeldung`) und `make test-integration` für 17 |
| Handbuch/README: Klartext-Passwort und MD5 mit 17 | `make test-integration`, Fälle `md5` und `password` |
| Handbuch: Anmeldung ohne Passwort mit 17 (Bestand, nicht erweitert) | `make test-integration` (die Tests ohne Rolle laufen gegen den Server mit `trust`, z. B. `TestE2ERecordSelect1`) |
| Handbuch: Client pgx 5.11 | `go.mod` (`github.com/jackc/pgx/v5 v5.11.0`); `TestE2ERecordAnmeldung` nutzt `pgconn` aus pgx |
| Handbuch: `psql` aus libpq 17 gegen PostgreSQL 17 (Bestand von `slice-v1-abschluss-anmeldung-doku`) | neu gefahren mit dem Aufruf unter der Tabelle (*Aufruf der `psql`-Probe*): `psql` 17.11 aus dem Image der Referenz, vier Anmeldungen (je Verfahren und `trust`) liefern ihre Zeile, 4 `sql:`-Einträge, 0 Treffer für Passwörter, `md5…`, `SCRAM` |
| Andere Hauptversionen und Vorabversionen nicht geprüft | verneinend; kein Lauf behauptet sie |

**Aufruf der `psql`-Probe** (nicht im Repo, kein Vertrag; Voraussetzung `make build`; Datei `psqlprobe.sh` mit dem Text unten anlegen; `PGI` ist das Image der Referenz aus `harness/mk/integration.mk` (`POSTGRES_IMAGE`), `S` ein leeres, beschreibbares Verzeichnis; Aufruf `bash psqlprobe.sh "$PGI" "$S"`; Ausgabe: `psql (PostgreSQL) 17.x`, vier Zeilen mit `current_user`, `Treffer Geheimnis: 0`, `4`):

```bash
set -u
PGI=$1; S=$2
cleanup(){ docker rm -f impl-zus-rec impl-zus-pg >/dev/null 2>&1; docker network rm impl-zus-net >/dev/null 2>&1; rm -rf "$S/rec"; }
trap cleanup EXIT
docker network create impl-zus-net >/dev/null
mkdir -p $S/rec; chmod 777 $S/rec
docker run -d --name impl-zus-pg --network impl-zus-net -e POSTGRES_HOST_AUTH_METHOD=trust $PGI >/dev/null
for i in $(seq 60); do docker exec impl-zus-pg pg_isready -q -h 127.0.0.1 -U postgres && break; sleep 0.5; done
docker exec impl-zus-pg psql -U postgres -qc "SET password_encryption='scram-sha-256'; CREATE ROLE ps LOGIN PASSWORD 'pwS'; SET password_encryption='md5'; CREATE ROLE pm LOGIN PASSWORD 'pwM'; CREATE ROLE pk LOGIN PASSWORD 'pwK';"
docker exec impl-zus-pg sh -c "sed -i '1i host all ps all scram-sha-256\nhost all pm all md5\nhost all pk all password' /var/lib/postgresql/data/pg_hba.conf; psql -U postgres -qc 'select pg_reload_conf()' >/dev/null"
docker run -d --name impl-zus-rec --network impl-zus-net -v $S/rec:/rec pgwire-recorder:dev record --upstream impl-zus-pg:5432 --listen 0.0.0.0:15432 --output /rec/r.yaml >/dev/null
sleep 2
docker run --rm $PGI psql --version
for u in "ps pwS" "pm pwM" "pk pwK"; do set -- $u; docker run --rm --network impl-zus-net -e PGPASSWORD=$2 $PGI psql "host=impl-zus-rec port=15432 user=$1 dbname=postgres sslmode=disable" -At -c "select current_user, version()" | cut -c1-60; done
docker run --rm --network impl-zus-net $PGI psql "host=impl-zus-rec port=15432 user=postgres dbname=postgres sslmode=disable" -At -c "select current_user"
docker stop impl-zus-rec >/dev/null; sleep 1
echo "Treffer Geheimnis: $(grep -c -E 'pw[SMK]|md5[0-9a-f]{32}|SCRAM' $S/rec/r.yaml)"; grep -c "sql:" $S/rec/r.yaml
```

**„Alle übrigen Szenarien“ (Vorgabe des Architect):** Nicht nachgewiesen und deshalb nicht geschrieben. `test/integration` zitiert nur die Szenarien 6, 7, 12 bis 16 im Kommentar; Szenarien wie 9 (Container), 11 (Homebrew) und 17 (Prüfbarkeit des Quellcodes) sind Gegenstand anderer Gates, und eine Zuordnung Szenario zu Test habe ich nicht für jedes Szenario geführt. Die Texte sagen „volle Integrationssuite“ (`make test-integration` fährt ohne Auswahlvariablen alle Tests mit dem Tag `integration`, Kopf von `tools/test/run-integration-tests.sh`) und im Handbuch „alle weiteren Integrationstests“, nicht „alle übrigen Szenarien“.

**Mutation (die Zusage „Hauptversion N“ greift):** Eintrag `14=` mit dem Image von 16 (`POSTGRES_VERSIONEN='14=postgres:16-alpine@sha256:721873c3…80ea 17=…'`): exit 2, `Version 14: der Server meldet server_version '16.15'`, `ROT — 14 (Aufbau)`. Ein Lauf, der die Hauptversion nicht prüfte, bliebe grün. Weitere Mutationen für die Zusagen der Texte gibt es nicht: Texte tragen kein Gate (§6, Zahl der prüfbaren Aussagen); die Mutationen der Matrix selbst gehören zu `slice-v1-abschluss-postgres-versionen-matrix`.

**Prüfungen der Texte:** Der Diff fügt der Spezifikation und dem Lastenheft kein `make`-Ziel, keinen Verweis auf `harness/vertraege` und keinen Namen eines Testziels hinzu: `git diff 510e7ab -- spec | grep '^+' | grep -c "harness/vertraege\|make \|test-postgres\|test-integration"` gibt 0. Im Bestand trifft `grep -n "make " spec/*.md` zwei Zeilen in `spec/spezifikation.md` (`make a-check`, `make a-check-negativ`, Zeit-Port); sie sind nicht im Diff und bleiben unverändert (Bestand, Review F-7). Ohne `make `: `grep -n "harness/vertraege\|test-postgres\|test-integration" spec/*.md` keine Treffer. Quantoren-Durchsicht (Schritt 17: „jede“, „alle“, „jeder“ mit Server, Datenbank, Version) in Handbuch, README, Spezifikation: kein weiterer Satz sagt etwas über „jeden Server“ oder „alle Versionen“ (die einzige Fundstelle, `LH-FA-05.e`, ist ersetzt); „geprüft mit PostgreSQL 17“ steht nach der Änderung nirgends mehr ohne Verfahren (`grep -n "geprüft mit PostgreSQL 17"` in Handbuch und README: keine Treffer). Handbuch ohne Verweis auf Spezifikation, ADR, Slices, Reviews; kein neues Beispiel, also keine Datei zu fahren.

**Gates:** `make gates` exit 0, 10 min 25 s, am Stand nach dem Nachzug des Reviews (Commit `1bc7dc7` und der Plan mit Link auf [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) in dieser Zeile); `make docs-check` danach 770 Dateien, 0 Befunde. Ein erster Lauf in diesem Stand war rot, weil dieser Link in §7 fehlte (`id-unlinked`, Review F-1).

**Größe (§4):** `git diff --numstat 510e7ab..1bc7dc7 -- spec docs/user README.md`: 57 (Summe aus Einfügungen und Löschungen: README 3, Handbuch 24, Lastenheft 7, Spezifikation 23); Meldeschwelle 300, Marke 400 nicht berührt.

**Nachzug nach V-1 und V-2 des Verifier-Reports (2026-10-15):** `SPEC-029` Punkt 2 ohne Szenario 18 und Szenario 16 Satz 3, Nachweis-Zeile von `LH-RB-02` mit „mindestens“; Tabellenzeile „Referenz 17“ in der Deckung und Aussage (g) in §6 ohne beide Szenarien. `grep -rnE "Szenario 18|Szenario 16" spec docs/user README.md`: Treffer nur in `spec/lastenheft.md` (Szenario selbst, Zeilen 810, 1087, 1101) und nirgends im Handbuch, in der README oder in `SPEC-029`; keine Versionszusage verbleibt. Die Sätze in §1, §3 und §6, die beide Szenarien als spätere Sendung nennen, sind Plan-Zusage des Architect und unverändert. `make gates`: erster Lauf rot in `postgres-versionen-gegenprobe` (Mutation `p7-int-treiber`: zusätzlich `p7-sigTERM-gruppe` rot; zeitabhängig), `make test-postgres-versionen-gegenprobe` allein danach grün, zweiter voller `make gates` exit 0. Größe: `git diff --numstat 510e7ab..<Commit der Spezifikation> -- spec docs/user README.md` Summe 56 (README 3, Handbuch 24, Lastenheft 7, Spezifikation 22); Marke 400 nicht berührt.

**Nachzug nach dem Review-Report (2026-10-15):** F-1 Link gesetzt; F-2 Trefferstand der Prüfung genau gefasst; F-3 §1, §3, §6 nachgezogen; F-4 Nachweis-Zeile von `LH-RB-02` enger gefasst (Teilmenge); F-5 Aufruf der `psql`-Probe im Wortlaut oben, am Stand `1bc7dc7` neu gefahren (`psql` 17.11, vier Anmeldungen, `Treffer Geheimnis: 0`, 4 `sql:`-Einträge, Container und Netz danach entfernt); F-6 und F-7 ohne Aktion.

**Rückgaben an den Architect:** keine. Keine Randform gefunden, die §6 nicht nennt; die Vorgabe „alle übrigen Szenarien“ ist oben enger gefasst (Entscheidung der Formulierung nach §6: Prüfung vor dem Satz).

### Closure (Planner, 2026-10-15)

**Stand:** HEAD `5804e96`; voller `make gates` grün für diesen Stand (Meldung des Auftraggebers; diese Closure fährt keine Gates). Review-Report `docs/reviews/2026-10-15-review-slice-v1-abschluss-postgres-versionen-zusage.md` (F-1 bis F-7, alle behoben oder ohne Aktion), Verifier-Reports `…-verify-…-zusage.md` (V-1 MEDIUM, V-2 LOW, beide behoben) und `…-r2.md` (PASS, V-1 und V-2 geschlossen, 1 INFO). Kein Gate rot, kein Carveout.

**Kriterien:** (1) DoD-Punkte 1 bis 3 sind geliefert und vom Verifier gegen die Läufe geprüft; (2) `make gates` grün, `make docs-check` und `make kopf-check` ohne Befund.

**Lerneintrag — benannte Spec-Lücke.** Die Zusage in `SPEC-029` Punkt 2 und die Nachweis-Zeile von `LH-RB-02` beschrieben im Entwurf einen Zielstand (Abnahmeszenario 18 und Szenario 16 Satz 3 „nur gegen 17 geprüft“; „nur“ auch für die Referenz 17), den erst die Tests zweier Folge-Slices prüfen. Die Lücke ist benannt und hat Adressen: `slice-v1-abschluss-verbindungsende-upstream` nimmt den Satz zu Szenario 18, `slice-v1-abschluss-verbindungsende-vergleich` den zu Szenario 16 Satz 3 in `SPEC-029` Punkt 2 auf, je im Commit des Tests; bis dahin steht keiner in der Spezifikation. Der Zähler: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Regel `AGENTS.md` §3.11 gilt, Wiederauftreten, kein neuer Wortlaut); `BEO-REPO/plan-folgt-korrektur-nicht` (F-1 bis F-3, Regel §3.9 gilt). Eine neue Regel ist nicht eingetragen; ob ein Spec-Entwurf vor dem Satz je Zusage den Lauf nennen muss, entscheidet der Nutzer.

**Register-Zeilen:** `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (neuer Beleg, verkörpert, Regel bleibt) · `BEO-REPO/plan-folgt-korrektur-nicht` (neuer Beleg, verkörpert) · `BEO-REPO/postgres-versionsliste-altert` (2×, offen) · neu `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig` (1×, offen: die Mutation `p7-int-treiber` färbte im ersten vollen `make gates` einmal zusätzlich `p7-sigTERM-gruppe`) · Kein Auftreten, Zähler 1: `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`.

**Risiko-Ausgänge:** siehe §6, je Risiko einer (1 *weiter offen*, 1 *eingetreten*, 4 *entfallen*).

**Sendung I-1** (Verifier-Report `…-r2.md`, INFO): Der Numerus des Satzes „Die Anmeldung mit Klartext-Passwort und mit MD5 sind nur dort geprüft“ in `SPEC-029` Punkt 2 wird berichtigt, sobald ein Nehmer die Stelle berührt. Adresse: `slice-v1-abschluss-verbindungsende-upstream` und `slice-v1-abschluss-verbindungsende-vergleich` (beide schreiben einen Satz in Punkt 2), in beiden §1 mit der Kennung dieses Gebers und Nachzählung in §8 (je 3 Liefer-Punkte, 2 Schichten, unverändert).

**Paarungen:** (a) Anker: kein Eintrag mit `liegt in`; (b) Folge-Slice: `slice-v1-abschluss-verbindungsende-upstream`, `slice-v1-abschluss-verbindungsende-vergleich` liegen in `next/` und nennen die Kennung dieses Gebers in §1 und DoD (`grep -n` je Nehmer ≥ 1); (c) Register: jede genannte `BEO-REPO/…` existiert als Verzeichnis mit nicht leerem `evidence/`.

**Folge:** Die Liste der Versionen bleibt eine Handlung von Hand (`BEO-REPO/postgres-versionsliste-altert`); nächster Slice laut Roadmap: `slice-v1-abschluss-protokollrand` (Schritt 33).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-14); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` — 1×, offen; die Sendung der Handbuch-Sätze ist die Antwort für den Handbuch-Teil, der Beleg je Version kommt von der Matrix.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` — 3×, offen; die Läufe je Version laufen unter Linux im Container, das Handbuch sagt zu Windows und macOS nichts (§1). Ein weiteres Auftreten bliebe ein viertes Auftreten eines nicht verkörperten Eintrags.
- `BEO-REPO/postgres-versionsliste-altert` — 1×, offen (angelegt von `slice-v1-abschluss-postgres-versionen-matrix`); die Liste in der Spezifikation altert ohne Sensor, dieser Slice führt sie.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` — 1×, offen; dieser Slice ist die Adresse des Handbuch-Teils der Zusage und der Sendung.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 55×, verkörpert (`AGENTS.md` §3.11); Grund der Reihenfolge und des engeren Wortlauts von „unterstützt“.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 29×, verkörpert (`AGENTS.md` §3.12); die Randformen stehen in §6.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` — 10×, verkörpert (`AGENTS.md` §3.13); die Sendung steht mit der Kennung des Gebers in §1.
- `BEO-REPO/plan-folgt-korrektur-nicht` — 45×, verkörpert (`AGENTS.md` §3.9); gilt als Regel, auch für den Kopf.

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure ihrer Welle, nicht hier.

Nachgezählt beim Anlegen, einschließlich der Sendung der Handbuch-Sätze (Planner, 2026-10-14, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Spezifikation und Lastenheft; Handbuch; README); die Sendung gehört zu den Punkten Handbuch und README und fügt keinen hinzu. Schichten: zwei, Spezifikation (Lastenheft und Spezifikation, wie im Ursprungsplan) und Dokumentation (Handbuch und README, `AGENTS.md` §3.13). Eine ADR würde eine dritte Schicht sein, nicht Teil der Spezifikation: Die Source Precedence führt sie auf einem eigenen Rang (4), und der Ursprungsplan zählte „Spezifikation, Entscheidungen, Harness und Dokumentation“ getrennt; der Plan sieht keine vor (§1). Die Regel in `implement-slice.md` Schritt 19 ist Bestand und zählt nicht. Schätzung: 180 bis 250 Zeilen in der Metrik von §4.

Nachgezählt beim Eintragen der Nahtstelle von `slice-v1-abschluss-postgres-versionen-matrix` (Planner, 2026-10-14, `AGENTS.md` §3.13): Liefer-Punkte weiter 3 (Spezifikation und Lastenheft; Handbuch; README): Die vier Punkte der Nahtstelle sind Eingaben dieser Punkte (Liste und Teilmenge für `SPEC-029`, Aufrufe für den Beleg in §7), kein Artefakt, das hinzukommt. Schichten weiter zwei (Spezifikation, Dokumentation), kein Produkt-Code, kein Gate, keine ADR. Schätzung weiter 180 bis 250 Zeilen in der Metrik von §4; die Liste mit fünf Versionen und die Teilmenge in Worten liegen in den rund 60 Zeilen für Spezifikation und Lastenheft; die Meldeschwelle 300 und die Marke 400 bleiben. Der Nehmer liegt nicht über den Grenzen (höchstens drei Punkte, höchstens zwei Schichten); ein Schnitt ist nicht nötig. Ob der Satz zu Klartext und MD5 auf den Lauf von Hand je Version wächst, hängt an der Entscheidung in §6 (Randformen); mit ihr bliebe die Zählung gleich (ein Beleg in §7, kein Artefakt).

Nachgezählt beim Eintragen der Entscheidung des Nutzers vom 2026-10-15, Frage B (Planner, `AGENTS.md` §3.13, §3.9): Liefer-Punkte weiter 3 (Spezifikation und Lastenheft; Handbuch; README); die Begrenzung der Anmeldungssätze ist Inhalt der Punkte Handbuch und README, kein neues Artefakt. Schichten weiter zwei (Spezifikation, Dokumentation), kein Produkt-Code, kein Gate, keine ADR. Der Lauf von Hand je Version entfällt; die Schätzung bleibt bei 180 bis 250 Zeilen (eher am unteren Ende), Meldeschwelle 300 und Marke 400 bleiben. Der Nehmer liegt nicht über den Grenzen; ein Schnitt ist nicht nötig.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
