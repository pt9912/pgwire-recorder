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

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-11`](../../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-04`](../../../../spec/lastenheft.md#lh-fa-04--transparenz-für-clients), [`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung), [ADR-0007](../../adr/0007-strict-replay.md), [ADR-0012](../../adr/0012-extended-query-gruppen.md)

**Berührte Spec-Stellen:** `LH-FA-18.a` · `LH-FA-09.a` · `LH-FA-10.a` · `LH-FA-11.a` · `LH-FA-13.a` · `LH-FA-13.b` · `SPEC-025` · `SPEC-033` · `SPEC-041` · `SPEC-011` · `ARC-002` · `ARC-003` · `ARC-006`

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

**Liefert:** Der Replay-Use-Case nimmt Extended-Nachrichten an (`ClientMessage` am Port `Replayer`), vergleicht jede in allen Feldern mit der erwarteten am Cursor (Strict Matcher `abweichung`, eine leere Liste gleicht nil) und gibt die Server-Nachrichten einer Gruppe nach deren `Flush` oder `Sync` frei, vorher keine. Der Cursor führt Interaktion, Gruppe und Nachricht (`ARC-002`). Eine Abweichung ist `PGR-E5001`; ihre Diagnose nennt Session, Interaktion, Gruppe, Nachricht, erwarteten und empfangenen Nachrichtentyp und das abweichende Feld, bei SQL beide Texte, bei `bind`, `execute`, `describe` und `close` das SQL der erwarteten und der empfangenen Anweisung, soweit es nach dem Verlauf der Aufzeichnung vor dem Cursor besteht (`LH-FA-10.a` „erwartete Query“): nur vom Server bestätigte `parse`, `bind` und `close` wirken; `close`, eine einfache Anfrage (unbenanntes Statement und Portal) und das Ende der Transaktion (Portale) beenden ein Objekt, sonst „unbekannt“; nie Parameterwerte. Bei einer Abweichung der Art und nach dem Ende der Aufzeichnung nennt die Diagnose die Extended-Nachricht ebenso mit dem SQL ihrer Anweisung. Bestätigungen zählen je Interaktion, damit eine späte Bestätigung in der folgenden Gruppe mitzählt. SQL-Befehle, die Objekte beenden, bildet die Diagnose nicht nach (§6). Eine einfache Anfrage, wo eine Extended-Nachricht erwartet ist, und umgekehrt, ist eine Abweichung (aus `LH-FA-18.a` „jede eingehende Client-Nachricht muss … entsprechen“ abgeleitet). `replaySitzung` leitet Extended-Nachrichten an den Use Case, statt sie mit `PGR-E6001` abzulehnen, und schließt die Verbindung nach einer Abweichung, ohne weiterzulesen; der vorhandene Handshake genügt pgx im Standardmodus. Beim Herunterfahren beantwortet `replaySitzung` eine begonnene Extended-Interaktion bis zu ihrem `Sync` und endet danach, ohne eine neue zu lesen (`LH-FA-13.a`); ob eine Interaktion läuft, sagt der Use Case (`Shutdown` am Port `Replayer`). Das Warten ist hier nicht begrenzt: Die Frist `--shutdown-timeout` ist im Replay nicht verdrahtet (`slice-v1-abschluss-betrieb`). Schließt sie später die Client-Verbindung, endet die Session wie bei jedem Verbindungsende regulär; die Frist muss eine unvollständige Interaktion selbst als `PGR-E4006` merken (§6). Der Replay-Service prüft beim Start jede Interaktion mit `Validate` (`PGR-E3003`), auf deren Form sich der Cursor verlässt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Parametermatching über den exakten Vergleich hinaus — Folge-Anforderung, falls gewünscht.
- Abbildung der Extended-Nachrichten im PGWire-Adapter (`toClientMessage`, `toMessage`) — geliefert von `slice-extended-query-record`.
- Zwei-Richtungs-Ablauf des Record ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)) im Replay — bleibt bewusst stehen: `replaySitzung` liest und antwortet abwechselnd wie ein PostgreSQL-Server; der Gegendruck-Ablauf des Record läuft über `replay` durch (`TestE2EReplayExtendedGegendruck`).
- `--fail-on-unconsumed` (`PGR-E5002`) — `slice-replay-semantik-mismatch`; die Diagnose der Extended-Abweichung oben ist dort als geliefert vermerkt.
- Frist `--shutdown-timeout` und `PGR-E4006` im Replay — `slice-v1-abschluss-betrieb`; dort steht, dass `replay` die Frist durch Schließen der Client-Verbindung durchsetzt und das Ende einer unvollständigen Interaktion dann `PGR-E4006` ist, kein reguläres Ende.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Mit gestoppter PostgreSQL liefert Replay dem Go-Client mit Prepared Statements dasselbe Ergebnis (Abnahmeszenario 7); beim Herunterfahren beantwortet Replay eine begonnene Extended-Interaktion bis zu ihrem `Sync` und endet danach ([`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)); eine Aufzeichnung, deren Interaktion die Form verfehlt, ist beim Start `PGR-E3003`.
- [x] [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage): Eine Extended-Query-Abweichung wird erkannt und nicht beantwortet; eine aufgezeichnete `ErrorResponse` samt Verwerfen bis zum `Sync` wird reproduziert (Abnahmeszenario 6, Extended). Die Diagnose der Abweichung nennt Stelle, Nachrichtentypen, Feld und die SQL-Texte der Anweisungen ohne Parameterwerte (§1), die Klasse ist 5, und die Verbindung endet.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8) — Review und vier Folge-Reviews bis `1ab9db2`; `5a1da90` nur verifiziert, `a6e3fdb` ohne Review und Verifikation, siehe §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — bei Closure geprüft (§7); im Repo **mit** Wellen prüft sie die Closure von `welle-extended-query` erneut.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (`matcher.go`, `replay.go`, Tests `replay_extended_test.go`) | new, update | Strict Matcher für Extended-Nachrichten; Cursor mit Gruppe und Nachricht, `ClientMessage`, `Shutdown`, Abweichung bei falscher Art, SQL der Anweisungen in der Diagnose, Formprüfung beim Start |
| `internal/hexagon/ports/driving/replay.go` | update | `ClientMessage` und `Shutdown` am Replay-Use-Case |
| `internal/adapters/driving/pgwire` (`server.go`, Tests) | update | `replaySitzung` leitet Extended-Nachrichten an den Use Case und liest beim Herunterfahren eine laufende Interaktion zu Ende; Handshake unverändert, er genügt pgx im Standardmodus (`TestE2EOhnePostgresExtendedReplay`) |
| `test/integration` (`extended_replay_e2e_test.go`, `extended_e2e_test.go`) | new, update | Abnahmeszenario 7, Abweichung, Pipeline, Gegendruck und SIGTERM mitten in einer Folge im Replay; `gegendruckAblauf` aus dem Record-Test herausgelöst |
| `tools/test/run-integration-tests.sh` | update | Kopfkommentar: die erste Phase schreibt mehrere Aufzeichnungen |
| `docs/user/abdeckung-*.md`, `README.md` | update | Abdeckungstabellen (`make abdeckung`); Stand der Wiedergabe |
| `docs/plan/planning/open/slice-replay-semantik-mismatch.md`, `docs/plan/planning/welle-replay-semantik.md`, `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` | update | Betrieb: Frist auch für `replay`; Mismatch: Ziel, DoD und Plan nennen nur noch `--fail-on-unconsumed`; „Bereits geliefert“ nennt die Diagnose der Extended-Abweichung |
| `docs/reviews/` (Review, vier Folge-Reviews, Verifikation, Validierung zu `slice-extended-query-replay`) | new | Berichte der Review- und Verifikationsrollen, Belege der Korrekturen |
| Closure: `docs/plan/planning/open/slice-extended-query-lebendpruefung.md` (neu), `docs/plan/planning/welle-extended-query.md`, `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/slice-replay-semantik-mismatch.md`, `docs/plan/planning/welle-replay-semantik.md`, `docs/plan/planning/open/slice-replay-semantik-meldungscodes.md`, `docs/plan/planning/open/slice-v1-abschluss-betrieb.md`, `docs/plan/planning/observations/BEO-REPO/` | new, update | Folge-Slice der Lebendprüfung in der Welle samt Drift-Log; Adressen für F-327, F-328, F-345 und die Validierung (Befunde 2, 3, 5); Register-Belege (§7) |

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

- Wiederholte Ausführung desselben Prepared Statements: Zuordnung über Position. Stand: strict sequential über den Cursor; belegt mit `TestReplayExtendedWiederholung` (Werte a, b, a mit je eigener Antwort, außer der Reihe `PGR-E5001`) und `TestE2EReplayExtendedPgx` (pgx, Werte eins, zwei, eins) — **Ausgang:** entfallen: Die Zuordnung über die Position trägt; neben den beiden Tests bestätigt die Validierung (Frage 1, Sonde S1) dreimal denselben `INSERT` mit verschiedenen Werten in einem breiteren pgx-Ablauf, zweimal wiedergegeben mit byte-gleicher Ausgabe.

- Späte Antworten einer Flush-Gruppe: Server-Nachrichten, die im Record erst nach der ersten Client-Nachricht der folgenden Gruppe eintrafen, stehen in der folgenden Gruppe (`LH-FA-18.a`). Das Replay gibt sie erst nach deren `Flush` oder `Sync` frei; ein Client, der zwischen Beginn und Ende dieser Gruppe auf die Antwort der vorigen wartet, hält im Replay mit seiner eigenen Aufzeichnung an. Die Spezifikation nennt die Folge nicht; zudem sagt [ADR-0012](../../adr/0012-extended-query-gruppen.md), eine Gruppe enthalte die Server-Nachrichten, „die auf sie antworten“ (aus `slice-extended-query-record`, Review F-308). Stand: umgesetzt, wie `LH-FA-18.a` es sagt (Freigabe nur nach `Flush` oder `Sync` der Gruppe, `TestReplayExtendedGruppen`); ein Client, der nicht wartet, läuft durch (`TestE2EReplayExtendedGegendruck`, verzögerte Ausgabe einer Flush-Gruppe), der Fall eines wartenden Clients mit später Antwort ist nicht geprüft — **Ausgang:** weiter offen → `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` im Register. Der Fall setzt eine Flush-Gruppe voraus, nach der der Client nicht wartet, also die in `LH-FA-18.a` §Gruppen benannte Grenze von v1; dass das Replay einen solchen Client anhalten kann, nennt die Spezifikation nicht. pgx im Standardmodus schließt jede Folge mit `Sync` ab (Validierung, Frage 5), deshalb kein Folge-Slice.

- Leere Liste und nil sind im Domain-Modell gleichbedeutend; der YAML-Leser liefert nil, der PGWire-Adapter womöglich leere Listen. Ein Feldvergleich per `reflect.DeepEqual` meldete dann eine Abweichung, die keine ist (aus `slice-extended-query-modell`). Stand: Der Matcher vergleicht Listen nach Länge und Elementen, Werte nach `Null` und Bytes; belegt mit `TestReplayExtendedFelder` (leere Listen und leerer Wert gegen nil gleich, NULL gegen leeren Wert verschieden) — **Ausgang:** entfallen: Der Matcher vergleicht Listen nach Länge und Elementen statt per `reflect.DeepEqual` (Beleg oben); die Mutation, die `Null` ignoriert, ist rot (Review A4). Eine Abweichung allein aus leerer Liste gegen nil kann damit nicht entstehen.

- Herunterfahren im Replay mitten in einer Extended-Interaktion: Endet ctx, bricht `replaySitzung` das Lesen der nächsten Client-Nachricht ab, auch zwischen zwei Nachrichten einer Gruppe; die Session endet dann vor dem `Sync`, und die Interaktion zählt als nicht verbraucht (`PGR-W2001`). `LH-FA-13.a` sagt, laufende Sessions endeten nach Abschluss ihrer laufenden Interaktion, ohne den Replay-Modus für Extended zu nennen; die Randform war damit still im Code entschieden (Review F-321). Stand: im Slice behoben nach Entscheidung des Nutzers: `replaySitzung` beantwortet die laufende Interaktion bis zu ihrem `Sync` und endet danach (`TestReplayHerunterfahren`, `TestReplayExtendedHerunterfahren`, `TestE2EReplayExtendedSigtermMittenInFolge`, alle mit der vorigen Logik rot); die Frist dazu ist nicht verdrahtet (`slice-v1-abschluss-betrieb`) — **Ausgang:** eingetreten → `slice-v1-abschluss-betrieb` (Frist). Die Randform selbst ist im Slice entschieden und behoben (`LH-FA-13.a` gilt ohne Unterscheidung nach Modus; Folge-Review bestätigt, Verifikation DoD 1); den Rest, die fehlende Obergrenze, trägt das Risiko F-330 unten.

- Gegendruck im Replay mit Flush-Gruppe (Review F-329): `TestE2EReplayExtendedGegendruck` prüft nur die Aufzeichnung, in der die verzögerte große Ausgabe in der folgenden Sync-Gruppe steht. Nicht geprüft ist der Fall einer großen sofortigen Ausgabe einer Flush-Gruppe, während der Client eine große nächste Gruppe schreibt; das Replay schreibt dann die Ausgabe, bevor es weiterliest — **Ausgang:** entfallen: Ein Client, der so schreibt, verklemmt gegen PostgreSQL genauso, denn auch PostgreSQL liest und antwortet abwechselnd; das Replay verlangt vom Client nicht mehr als der Server, und pgx liest beim Schreiben großer Daten mit (Validierung, Frage 7).

- Herunterfahren ohne Obergrenze im Replay (Folge-Review F-330): Seit `replaySitzung` eine begonnene Extended-Interaktion bis zu ihrem `Sync` beantwortet, wartet `replay` nach `SIGTERM` unbegrenzt auf einen Client, der mitten in der Interaktion pausiert; `Serve` wartet auf alle Verbindungen. Die Frist `--shutdown-timeout` mit `PGR-E4006` ist nicht verdrahtet — **Ausgang:** eingetreten → `slice-v1-abschluss-betrieb` (§1, §3, §6 dort nennen `replay`).

- Diagnose ohne SQL-Befehle (Folge-Review 3, F-348): Die Nachbildung der Statements und Portale sieht nur Protokollnachrichten; `DEALLOCATE`, `DISCARD ALL`, `CLOSE` und `ROLLBACK TO SAVEPOINT` beenden Objekte, ohne dass sie es bemerkt. Die Diagnose kann dann ein beendetes Objekt als bestehend nennen; die Erkennung der Abweichung ist davon nicht berührt — **Ausgang:** entfallen: vom Bedarf getragen. pgx räumt seinen Statement-Cache über die Protokoll-Nachricht `Close` auf, die die Diagnose nachbildet; SQL-seitiges `deallocate all` sendet pgx nur auf ausdrücklichen Aufruf, `DISCARD ALL` kommt von externen Poolern, die im Testaufbau nicht zwischen Anwendung und Recorder stehen; selbst dann stimmt die Stelle, nur der SQL-Text ist veraltet (Validierung, Frage 6).

- Ordnung von Wächter, Lesefrist und `geweckt` beim Herunterfahren (Folge-Review F-337): belegt mit `TestReplayHerunterfahrenSpaeteFrist`, der das Setzen der Frist durch den Wächter bis nach der Verarbeitung einer Nachricht anhält; ohne `<-geweckt` rot (M1) — **Ausgang:** entfallen: Die Ordnung ist mit einem Test belegt, der für richtigen Code deterministisch ist (drittes Folge-Review, nach F-341) und ohne `<-geweckt` rot wird (M1).

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

- **Was hat funktioniert:** Der Replay-Modus beantwortet eine Extended-Query-Interaktion strict sequential aus der Aufzeichnung: Vergleich jeder Nachricht in allen Feldern, Freigabe einer Gruppe erst nach ihrem `Flush` oder `Sync`, `PGR-E5001` mit Ende der Verbindung, `ErrorResponse` samt Verwerfen bis zum `Sync`, Herunterfahren nach dem `Sync` der laufenden Interaktion, Formprüfung beim Start (`PGR-E3003`). Die Zusagen zum Matching fingen schon die Tests des ersten Commits (Review, A1 bis A6 und A8 rot). Die Verifikation hat alle drei Liefer-Punkte mit eigenen Läufen bestätigt: Abnahmeszenario 7 über getrennte Integrationsphasen, Mutationen, Sonden P1 bis P9. Die Nachverifikation zu `5a1da90` hat DoD-Punkt 2 bewusst gebrochen (Mutante NA: ohne den Fix ist `TestReplayExtendedDiagnoseArt` aus dem richtigen Grund rot). Die Validierung bestätigt Abnahmeszenario 7 mit einem breiteren pgx-Ablauf (gängige Typen, Transaktionen mit Fehler, Batch mit Verwerfen; zwei Wiedergaben byte-gleich), ebenso Abweichung und Fehlerantwort.
- **Was ging anders als geplant:** Fünf Prüfrunden statt einer, und fast alles, was nach der ersten Runde kam, betraf die Diagnose, nicht das Matching. `LH-FA-10.a` verlangt „erwartete Query“ und „tatsächlich empfangene Query“; was das bei `bind`, `execute`, `describe` und `close` heißt, sagt die Spezifikation nicht. Jede Randform entschied der Code, und erst eine Prüfrunde fand sie: Lebensdauer der Objekte (F-339), Antwort des Servers (F-347), SQL der Extended-Seite bei Art-Abweichung (V-25), späte Bestätigung (V-26), einfache Anfrage am Cursor (F-349). Dazu kam die Randform des Herunterfahrens im Replay (F-321), die der Nutzer entschied. Plan und Folge-Slices folgten den Korrekturen in drei Läufen hintereinander nicht vollständig (F-323, F-330, F-331, F-344; Kopf F-325, F-334, F-343, V-28), obwohl `AGENTS.md` §3.9 genau das verlangt. Die Validierung urteilte bei grüner Verifikation in Frage 2 rot: `pgxpool` und `database/sql` senden nach mehr als 1 s Leerlauf die einfache Anfrage `-- ping`, das strenge Replay lehnt sie ab, und ob sie in Aufzeichnung oder Wiedergabe steht, hängt am Zeitverhalten. Der Plan war erfüllt, der Bedarf aus [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) („genügt … die Umstellung von Host und Port“) für die häufigste Einsatzform nicht. Der Nutzer hat entschieden: tolerieren, umgesetzt in `slice-extended-query-lebendpruefung` in dieser Welle. `5a1da90` (Test Q1 zu F-349) hat kein eigenes Review, aber die Nachverifikation. `a6e3fdb` (Test zu V-31, Zählung in §3 zu V-32) hat weder Review noch Verifikation; ihn belegt `make gates`. Dass sein Test die Mutation NG2 jetzt fängt, hat niemand gezeigt; DoD-Punkt 2 hängt nicht daran (Nachverifikation N1, Sonde P9).
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke, offen: [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) und [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) halten unter strengem Replay ([ADR-0007](../../adr/0007-strict-replay.md)) für Connection-Pools und `database/sql` nicht: Deren Lebendprüfung nach Leerlauf macht Aufzeichnung und Wiedergabe vom Zeitverhalten abhängig, und Spezifikation, ADRs und Planung führten den Fall nicht (Validierung, Frage 2). Entschieden ist „tolerieren“; Spezifikation, eine neue ADR gegen [ADR-0007](../../adr/0007-strict-replay.md) (Architect) und die Umsetzung trägt `slice-extended-query-lebendpruefung`. Zwei weitere benannte Spec-Lücken mit Adresse: Die falsche Protokollart im Replay ist `PGR-E5001`, aber `LH-FA-18.a` §Replay nennt den Fall nicht ausdrücklich (F-327) → `slice-extended-query-lebendpruefung`, weil dessen Satz die Ausnahme für Lebendprüfungen gegen genau diesen Fall abgrenzt. `LH-FA-18.a` §Mismatch fasst die Regel zu Parameterwerten bedingt, `SPEC-033` unbedingt (F-328); dazu fehlt der Index des abweichenden Parameters in der Diagnose (Validierung, Befund 3) → `slice-replay-semantik-meldungscodes`, der `--log-level` liefert. Sensor-Kandidat, nicht gebaut: Ein Test der zweiten Integrationsphase lässt sich mit `make test-integration` nicht bewusst brechen, weil der Runner nach der roten ersten Phase endet (V-30); das steht im Register, die Entscheidung liegt bei der Closure von `welle-extended-query`.
- **Beobachtungs-Register (`../observations/`):** Zähler = Dateien unter `evidence/`; je Eintrag Beleg `evidence/slice-extended-query-replay.md`.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` — Zusage der Diagnose oder des Adapters ohne Test, der ihre Bedingung bricht (F-322/A7, F-337/M1, F-340, F-349/M1, V-31/NG2), **5×**.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` — Kommentare und Abdeckungszeilen sagen mehr zu als geprüft (F-322, F-324, F-332, F-333, F-335, F-340, F-342), **5×**.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` — Herunterfahren im Replay (F-321) und die Randformen von „erwartete Query“ in der Extended-Diagnose (F-339, F-347, V-25, V-26, F-349), **3×**, erstmals über der Schwelle.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` (verkörpert in `AGENTS.md` §3.9 seit welle-walking-skeleton) — Folge-Slice nur teilweise nachgezogen (F-323, F-330, F-331, F-344), Kopf nennt berührte Spec-Stelle nicht (F-325, F-334, F-343, V-28), §3 unvollständig (V-27, V-32), **5×**. Retirement-Check: wieder aufgetreten, die Regel bleibt; dass sie allein nicht greift, zeigt der dritte Lauf in Folge.
  - `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients/` — neu: Lebendprüfung von Pools (Validierung, Frage 2) und Flush-Gruppe mit später Antwort (§6), **1×**.
  - `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar/` — neu: V-30, **1×**.

  Über der Schwelle stehen `negativtests-fehlen-bei-neuem-vertrag` (5×), `zusage-im-kommentar-weiter-als-pruefung` (5×), `spec-randform-erst-im-review-entschieden` (3×; nach Modul 5 braucht sie damit einen Folge-Slice) und `plan-folgt-korrektur-nicht` (5×, verkörpert; Frage nach Schärfung oder Sensor, etwa für den Kopf). Ihre Ausgänge vergibt die Closure von `welle-extended-query`; die `state.md` bleibt bis dahin unverändert. Einmalig und nicht eingetragen: F-326 (Invariante des Lesers am Port), F-336 (Struktur-ID in der Commit-Message), F-341 (Test mit Zeitannahme).
- **Folge-Slices:** `slice-extended-query-lebendpruefung` (neu; Lebendprüfungen im Replay, ADR gegen [ADR-0007](../../adr/0007-strict-replay.md), F-327, Handbuch zu Pings und `ShouldPing`; Welle `welle-extended-query`), `slice-v1-abschluss-betrieb` (Frist im Replay, F-330; Handbuch zu `replay` im Container; Schnittfrage F-345 für sein `open → next`), `slice-replay-semantik-mismatch` (Rest `--fail-on-unconsumed`, Titel und Bezug nachgezogen), `slice-replay-semantik-meldungscodes` (F-328, Index des abweichenden Parameters) — alle vier Dateien in `open/`.
- **Risiken aus §6:** acht, jedes mit genau einem Ausgang — zwei eingetreten (→ `slice-v1-abschluss-betrieb`), fünf entfallen, eines weiter offen (→ Register); siehe §6.
- **Drei Paarungen:** Anker — kein `liegt in`, mit diesem Slice ist nichts verkörpert. Folge-Slice — die vier genannten Dateien liegen in `open/`. Register — die sechs genannten Kennungen bestehen als Verzeichnis mit nicht leerem `evidence/`. Die Closure von `welle-extended-query` prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-05-review-slice-extended-query-replay.md` (`b895565`, `09405f0`; F-321 bis F-329), Folge-Reviews `…-folge-review-slice-extended-query-replay.md` (`c8ebf08`; F-330 bis F-338), `…-folge-review-2-…` (`b697462`; F-339 bis F-346), `…-folge-review-3-…` (`2d3e0f7`; F-347, F-348), `…-folge-review-4-…` (`f528c90`, `1ab9db2`; F-349), Verifikation `docs/reviews/2026-10-05-verifikation-slice-extended-query-replay.md` (bis `f528c90`, V-25 bis V-30; Nachverifikation zu `5a1da90`, V-31, V-32; DoD 1 bis 3 bestätigt), Validierung `docs/reviews/2026-10-05-validierung-slice-extended-query-replay.md` (`5a1da90`; Fragen 1, 3 bis 7 tragen, Frage 3 mit Einschränkung, Frage 2 rot). `make gates` grün an `a6e3fdb`; die Closure ändert nur Planungsdokumente.

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
