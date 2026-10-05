# Slice slice-extended-query-replay: Extended Query im Replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-extended-query.

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-04`](../../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients), [`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung), [ADR-0007](../../adr/0007-strict-replay.md)

**Berührte Spec-Stellen:** `LH-FA-18.a` · `LH-FA-09.a` · `LH-FA-10.a` · `LH-FA-13.a` · `SPEC-033` · `SPEC-041` · `SPEC-011` · `ARC-002` · `ARC-003` · `ARC-006`

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

**Ziel:** Der Replay-Modus beantwortet eine Extended-Query-Interaktion strict sequential aus der Aufzeichnung; Abweichungen und Fehlerantworten werden erkannt beziehungsweise reproduziert.

**Liefert:** Der Replay-Use-Case nimmt Extended-Nachrichten an (`ClientMessage` am Port `Replayer`), vergleicht jede in allen Feldern mit der erwarteten am Cursor (Strict Matcher `abweichung`, eine leere Liste gleicht nil) und gibt die Server-Nachrichten einer Gruppe nach deren `Flush` oder `Sync` frei, vorher keine. Der Cursor führt Interaktion, Gruppe und Nachricht (`ARC-002`). Eine Abweichung ist `PGR-E5001`; ihre Diagnose nennt Session, Interaktion, Gruppe, Nachricht, erwarteten und empfangenen Nachrichtentyp und das abweichende Feld, bei SQL beide Texte, nie Parameterwerte. Eine einfache Anfrage, wo eine Extended-Nachricht erwartet ist, und umgekehrt, ist eine Abweichung (aus `LH-FA-18.a` „jede eingehende Client-Nachricht muss … entsprechen“ abgeleitet). `replaySitzung` leitet Extended-Nachrichten an den Use Case, statt sie mit `PGR-E6001` abzulehnen, und schließt die Verbindung nach einer Abweichung, ohne weiterzulesen; der vorhandene Handshake genügt pgx im Standardmodus. Beim Herunterfahren beantwortet `replaySitzung` eine begonnene Extended-Interaktion bis zu ihrem `Sync` und endet danach, ohne eine neue zu lesen (`LH-FA-13.a`); ob eine Interaktion läuft, sagt der Use Case (`Shutdown` am Port `Replayer`). Die Frist `--shutdown-timeout` ist im Replay nicht verdrahtet; sie schließt später die Client-Verbindung, was das Warten wie jedes Verbindungsende beendet. Der Replay-Service prüft beim Start jede Interaktion mit `Validate` (`PGR-E3003`), auf deren Form sich der Cursor verlässt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Parametermatching über den exakten Vergleich hinaus — Folge-Anforderung, falls gewünscht.
- Abbildung der Extended-Nachrichten im PGWire-Adapter (`toClientMessage`, `toMessage`) — geliefert von `slice-extended-query-record`.
- Zwei-Richtungs-Ablauf des Record ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)) im Replay — bleibt bewusst stehen: `replaySitzung` liest und antwortet abwechselnd wie ein PostgreSQL-Server; der Gegendruck-Ablauf des Record läuft über `replay` durch (`TestE2EReplayExtendedGegendruck`).
- `--fail-on-unconsumed` (`PGR-E5002`) — `slice-replay-semantik-mismatch`; die Diagnose der Extended-Abweichung oben ist dort als geliefert vermerkt.
- Frist `--shutdown-timeout` und `PGR-E4006` im Replay — `slice-v1-abschluss-betrieb`; hier ist nur die Stelle vorbereitet, an der die Frist anschließt (Schließen der Client-Verbindung).


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Mit gestoppter PostgreSQL liefert Replay dem Go-Client mit Prepared Statements dasselbe Ergebnis (Abnahmeszenario 7).
- [ ] [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage): Eine Extended-Query-Abweichung wird erkannt und nicht beantwortet; eine aufgezeichnete `ErrorResponse` samt Verwerfen bis zum `Sync` wird reproduziert (Abnahmeszenario 6, Extended). Die Diagnose der Abweichung nennt Stelle, Nachrichtentypen und Feld ohne Parameterwerte (§1), die Klasse ist 5.
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
| `internal/hexagon/services` (`matcher.go`, `replay.go`) | new, update | Strict Matcher für Extended-Nachrichten; Cursor mit Gruppe und Nachricht, `ClientMessage`, `Shutdown`, Abweichung bei falscher Art, Formprüfung beim Start |
| `internal/hexagon/ports/driving/replay.go` | update | `ClientMessage` und `Shutdown` am Replay-Use-Case |
| `internal/adapters/driving/pgwire` (`server.go`, Tests) | update | `replaySitzung` leitet Extended-Nachrichten an den Use Case und liest beim Herunterfahren eine laufende Interaktion zu Ende; Handshake unverändert, er genügt pgx im Standardmodus (`TestE2EOhnePostgresExtendedReplay`) |
| `test/integration` (`extended_replay_e2e_test.go`, `extended_e2e_test.go`) | new, update | Abnahmeszenario 7, Abweichung, Pipeline, Gegendruck und SIGTERM mitten in einer Folge im Replay; `gegendruckAblauf` aus dem Record-Test herausgelöst |
| `tools/test/run-integration-tests.sh` | update | Kopfkommentar: die erste Phase schreibt mehrere Aufzeichnungen |
| `docs/user/abdeckung-*.md`, `README.md` | update | Abdeckungstabellen (`make abdeckung`); Stand der Wiedergabe |
| `docs/plan/planning/open/slice-replay-semantik-mismatch.md` | update | Ziel, DoD und Plan nennen nur noch `--fail-on-unconsumed`; „Bereits geliefert“ nennt die Diagnose der Extended-Abweichung |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-extended-query-record` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: das Matching verlangt eine Änderung der Strict-Replay-Entscheidung — Entscheidung klären.
- `in-progress` → `open`: Der Client verlangt Verhalten außerhalb von `LH-FA-18` — Carveout.


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

- Wiederholte Ausführung desselben Prepared Statements: Zuordnung über Position. Stand: strict sequential über den Cursor; belegt mit `TestReplayExtendedWiederholung` (Werte a, b, a mit je eigener Antwort, außer der Reihe `PGR-E5001`) und `TestE2EReplayExtendedPgx` (pgx, Werte eins, zwei, eins) — **Ausgang:** offen bis Closure.

- Späte Antworten einer Flush-Gruppe: Server-Nachrichten, die im Record erst nach der ersten Client-Nachricht der folgenden Gruppe eintrafen, stehen in der folgenden Gruppe (`LH-FA-18.a`). Das Replay gibt sie erst nach deren `Flush` oder `Sync` frei; ein Client, der zwischen Beginn und Ende dieser Gruppe auf die Antwort der vorigen wartet, hält im Replay mit seiner eigenen Aufzeichnung an. Die Spezifikation nennt die Folge nicht; zudem sagt [ADR-0012](../../adr/0012-extended-query-gruppen.md), eine Gruppe enthalte die Server-Nachrichten, „die auf sie antworten“ (aus `slice-extended-query-record`, Review F-308). Stand: umgesetzt, wie `LH-FA-18.a` es sagt (Freigabe nur nach `Flush` oder `Sync` der Gruppe, `TestReplayExtendedGruppen`); ein Client, der nicht wartet, läuft durch (`TestE2EReplayExtendedGegendruck`, verzögerte Ausgabe einer Flush-Gruppe), der Fall eines wartenden Clients mit später Antwort ist nicht geprüft — **Ausgang:** offen bis Closure.

- Leere Liste und nil sind im Domain-Modell gleichbedeutend; der YAML-Leser liefert nil, der PGWire-Adapter womöglich leere Listen. Ein Feldvergleich per `reflect.DeepEqual` meldete dann eine Abweichung, die keine ist (aus `slice-extended-query-modell`). Stand: Der Matcher vergleicht Listen nach Länge und Elementen, Werte nach `Null` und Bytes; belegt mit `TestReplayExtendedFelder` (leere Listen und leerer Wert gegen nil gleich, NULL gegen leeren Wert verschieden) — **Ausgang:** offen bis Closure.

- Herunterfahren im Replay mitten in einer Extended-Interaktion: Endet ctx, bricht `replaySitzung` das Lesen der nächsten Client-Nachricht ab, auch zwischen zwei Nachrichten einer Gruppe; die Session endet dann vor dem `Sync`, und die Interaktion zählt als nicht verbraucht (`PGR-W2001`). `LH-FA-13.a` sagt, laufende Sessions endeten nach Abschluss ihrer laufenden Interaktion, ohne den Replay-Modus für Extended zu nennen; die Randform war damit still im Code entschieden (Review F-321). Stand: im Slice behoben nach Entscheidung des Nutzers: `replaySitzung` beantwortet die laufende Interaktion bis zu ihrem `Sync` und endet danach (`TestReplayHerunterfahren`, `TestReplayExtendedHerunterfahren`, `TestE2EReplayExtendedSigtermMittenInFolge`, alle mit der vorigen Logik rot); die Frist dazu ist nicht verdrahtet (`slice-v1-abschluss-betrieb`) — **Ausgang:** offen bis Closure.

- Gegendruck im Replay mit Flush-Gruppe (Review F-329): `TestE2EReplayExtendedGegendruck` prüft nur die Aufzeichnung, in der die verzögerte große Ausgabe in der folgenden Sync-Gruppe steht. Nicht geprüft ist der Fall einer großen sofortigen Ausgabe einer Flush-Gruppe, während der Client eine große nächste Gruppe schreibt; das Replay schreibt dann die Ausgabe, bevor es weiterliest — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand `main`, 2026-10-05). Alle Einträge liegen in der einzigen Sub-Area `REPO`; für diesen Slice einschlägig:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — 4×, `offen`. Der Replay-Pfad für Extended Query ist ein neuer Vertrag am Protokollrand; die Zusagen aus `LH-FA-18.a` und `LH-FA-10` brauchen Tests, die eine Mutation fangen.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — 4×, `offen`. Kommentare am Matcher sagen nur zu, was ein Test prüft.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 2×, `offen`. Trifft dieser Slice sie erneut, erreicht sie 3× und braucht einen Folge-Slice.
- `BEO-REPO/plan-folgt-korrektur-nicht` — `verkörpert` (`AGENTS.md` §3.9); gilt für diesen Slice als Regel.

Die beiden 4×-Einträge stehen über der Schwelle ohne Ausgang; ihr Ausgang wird bei der Closure von welle-extended-query vergeben (Lese-Schritt), nicht in diesem Slice. Die übrigen Einträge (`gate-*`, `kern-fremdimporte-*`, `record-fehlerantwort-*`, `werkzeug-commit-*`) betreffen Gates, Architektur-Regeln oder den Record-Pfad — keine Treffer für diesen Slice.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
