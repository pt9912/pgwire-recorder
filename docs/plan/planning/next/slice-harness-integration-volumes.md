# Slice slice-harness-integration-volumes: Der Integrations-Runner hinterlässt keine anonymen Volumes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeordnet am 2026-10-12 (Empfehlung des
Planners; der Nutzer hat verfügt, ohne Nachfrage weiterzumachen): nach den Slices von
welle-v1-abschluss und vor `slice-harness-abdeckung-gate` (WIP-Limit 1); Reihenfolge in
§4 *Start* und im Drift-Log der Roadmap.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert nur das Testgeschirr und ein Gate des
Harness, keine Zusage des Produkts. Bindung an Entscheidungen:
[ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Integrationstests
gegen ein gepinntes PostgreSQL-Image; der Runner ist dort der Träger).

**Berührte Spec-Stellen:** `SPEC-057` (neuer Vertrag *Aufräumen der Integrationstests* in
`harness/vertraege/`, Ort nach `MR-002`; vom Architect am 2026-10-17 geschrieben, nicht im Diff dieses
Slice); `SPEC-055` (Punkt 9, Absatz *Volumes*: `test-container` entfernt seine Container mit `rm -fv`; Architect am 2026-10-17, nicht im Diff; Geber `slice-v1-abschluss-container`, geschlossen); `SPEC-052` (Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest;
Architect am 2026-10-17, nicht im Diff; Punkt 4 gilt nach dem Namen ohne Änderung); `SPEC-048` (nur Bindung: Form der
Abdeckungs-Deklaration im Kopf der Gegenprobe, §1; keine Änderung)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Nach einem Lauf von `make test-integration` liegt kein Docker-Volume mehr auf dem
Host, das einer seiner Container angelegt hat: Das Aufräumen entfernt die anonymen Volumes
der eigenen Container mit (`docker rm -fv`), der Runner prüft es nach Namen, und eine
Gegenprobe im Gate zeigt, dass die Prüfung rot wird, wenn jemand die Option wieder
weglässt. Dasselbe gilt seit dem 2026-10-17 für `make test-container`: `tools/harness/run-container.sh`
entfernt seine Container (`$pg` vor dem Wiedergeben und im Aufräumen) mit `rm -fv`, und die Gegenprobe
`test-container-gegenprobe` hat je Ende einen Fall `p9-*-volume` mit der Mutation `rm -f`
(`SPEC-055` Punkt 9). Der Befund dort (669 Volumes, 169,6 GB, `test-container` rot am vollen
Docker-Speicher) ist derselbe Befund, gemessen vom Implementer von `slice-v1-abschluss-container-doku`.
Geber der Sendung: `slice-v1-abschluss-container` (geschlossen, daher keine Adresse), der Befund kam aus
`slice-v1-abschluss-container-doku`; dieser Slice nimmt an (Empfehlung des Architect vom 2026-10-17; der Nutzer hat nicht entschieden).

**Befund** (Register: `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, 1×):
`tools/test/run-integration-tests.sh` ruft im Aufräumen `docker rm -f "$tests" "$pg"` ohne
`-v` auf. Das gepinnte PostgreSQL-Image deklariert `VOLUME /var/lib/postgresql/data`; jeder
Lauf legt dafür ein anonymes Volume an, das niemand entfernt (das benannte Volume `$daten`
wird gesondert mit `volume rm` entfernt). Am 2026-10-12 war `/docker` zu 100 % voll
(448 von 472 GB), darunter 2314 ungenutzte anonyme Volumes mit 172,9 GB (rund 75 MB je
Volume; die Angabe „seit 2026-10-02“ liegt vor der Einführung des Skripts, `b52a8e3`,
2026-10-04, also gehört nicht jedes Volume sicher diesem Repo); `make test-integration` und
damit `make gates` scheiterten mit „PostgreSQL nicht bereit“. Der Nutzer schuf selbst Platz.

**Bestand** (gemessen vom Planner am 2026-10-12, Stand `79258d3`, Docker 29.9.0):

- Im Runner startet nur `run -d --name "$pg"` (`run-integration-tests.sh:42`) ohne `--rm`; die
  zwei `run --rm` (`:73`, `:84`) laufen mit dem Testimage, dessen `.Config.Volumes` leer ist.
  Das Aufräumen (`:33`) ist `rm -f` ohne `-v`.
- `tools/harness/run-container.sh` (Stand 2026-10-17): `run -d --name "$pg"` (PostgreSQL-Image, `:183`) ohne
  `--rm`, entfernt mit `rm -f` in Zeile `:198` (vor dem Wiedergeben, im Lauf) und im Aufräumen (`:37`); `$rec`
  (`:189`) läuft mit dem Produkt-Image, `.Config.Volumes` misst der Implementer; `psql` (`:169`) und `version`
  (`:179`) tragen `--rm`; `docker create` (`:143`, `:200`) mit dem Produkt-Image. `run-postgres-versionen.sh`
  ruft `run-integration-tests.sh` je Version und legt selbst keinen Container an: Es hat kein eigenes Leck, mit
  der Behebung im Runner ist es mitbehoben.
- Weitere `docker run` in `tools/` und `harness/mk/`: `harness/mk/build.mk:22`,
  `harness/mk/schema.mk:9` und `:12`, `tools/arch/a-check-negativ.sh:83`,
  `tools/harness/traeger-fetch.sh:146`, jedes mit `--rm`; `docker create` und `docker compose`
  kommen nicht vor, `docker build` legt keine Volumes an. `docker run --rm` auf das
  PostgreSQL-Image hinterlässt kein Volume (Probe: 50 vor, 50 nach); `.Config.Volumes` ist leer
  für das Go-Image und das Image von d-check, für `a-check` und `d-migrate` misst es der
  Implementer (§3), `--rm` genügt in jedem Fall.
- Probe mit Container und Volume unter eigenem Namen: `docker create` des PostgreSQL-Images
  zeigt das Volume schon im Zustand `created` in `Mounts`; `docker rm -f` lässt es liegen
  (50 → 51 Volumes), `docker rm -fv` nicht. Ein Lauf `make gates` erhöht die Zahl der Volumes
  auf dem Host um eins (49 → 50).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Bestehende Volumes aufräumen (die 2314 vom 2026-10-12 und was seither liegt), `docker volume
  prune`, `docker system prune`, jede hostweite Bereinigung — anderer Vorgang und
  Betrieb des Hosts: Einem anonymen Volume sieht niemand an, ob es diesem Repo gehört, und das
  Werkzeug des Repos entfernt nur, was es im eigenen Lauf angelegt hat (Namen aus
  `docker inspect` seiner Container). Der Nutzer hat den Platz selbst geschaffen.
- Eingriff in Volumes, die dieser Lauf nicht angelegt hat — dieselbe Grenze als Zusage: Die
  Gegenprobe (Liefer-Punkt 3) legt ein fremdes benanntes und ein fremdes anonymes Volume an
  und prüft, dass beide bleiben.
- Die übrigen `docker run` (siehe *Bestand*) — Bestand bleibt bewusst stehen: Sie tragen
  `--rm`, das anonyme Volumes beim Ende des Containers mit entfernt (Probe), und kein Image
  davon ist ein Befund; der Slice belegt es in §7 und ändert sie nicht.
- Eine Prüfung nach Namen in `run-container.sh` (Namen lesen, nachprüfen, Rest entfernen wie `SPEC-057`
  Punkt 4) — Bestand bleibt bewusst einfach: Nur `$pg` trägt ein `VOLUME`, und die Gegenprobe fängt das
  Weglassen von `-v`; die Prüfung wäre ein zweiter Vertrag am Gate (`SPEC-055` Punkt 9, akzeptiertes Negativ).
  Ein gemeinsamer Hilfsaufruf für beide Läufer ist benannt, nicht angelegt: Entstünde ein dritter Läufer mit
  einem Container eines Images mit `VOLUME`, wäre das der Anlass.
- Eine Diagnose der vollen Platte im Runner („PostgreSQL nicht bereit“ nennt die Ursache
  nicht, `docker logs` des PostgreSQL-Containers erscheinen nur nach einem roten Testlauf) —
  anderer Vorgang: Es wäre ein neuer Vertrag am Runner mit eigenen Randformen; er braucht
  einen eigenen Slice, wenn der Nutzer ihn will, und ist nicht angelegt.
- Produkt-Code, Tests unter `test/integration`, das `Dockerfile` und das Produkt-Image —
  Schicht-Abgrenzung: Der Slice ändert nur Skripte unter `tools/test`, Make-Fragmente unter
  `harness/mk` (`integration.mk`, `gates-schnell.mk`) und die Dokumentation und Verträge des
  Gates (`harness/README.md`; die Verträge `SPEC-057` und `SPEC-052` sind vor dem Code geschrieben).

**Sendungen an Dritte** (`AGENTS.md` §3.13; mit der Kennung dieses Slice als Geber):

- *An `harness/mk/gates-schnell.mk` und `SPEC-052`:* `integration-volumes-gegenprobe` endet auf
  `-gegenprobe` und hängt an `GATE_CHECKS`; `SPEC-052` Punkt 4 verlangt dafür einen Eintrag in
  `GATE_SCHNELL_AUSNAHME`, sonst ist `gates-schnell-gegenprobe` rot. Punkt 13 verlangt für jedes
  Gate der Ausnahmeliste auf `-gegenprobe` den Beleg, dass es `docs/user/` und `README.md` des
  Arbeitsbaums nicht liest; die Gegenprobe fährt in Kopien, die sie selbst anlegt (§3), und der
  Beleg steht in `SPEC-052` Punkt 13 (Architect, geschrieben). Der Eintrag und die Zeile stehen
  im selben Commit wie das Ziel.
- *An `slice-harness-abdeckung-gate` (Nehmer der Bindung, nicht dieser Slice):*
  `tools/test/integration-volumes-gegenprobe.sh` liegt unter dem Muster `tools/**/*-gegenprobe.sh`,
  das `tools/test/abdeckung.sh` mit jenem Slice mitliest. Diese Gegenprobe trägt im Kopf keine
  `#`-Zeile in der Form einer Abdeckungs-Deklaration (`SPEC-048`); `slice-harness-abdeckung-gate`
  hält sie mit einem Fall seiner Gegenprobe fest (dort §1 und §6, Risiko *Andere Gegenproben werden
  gelesen*).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Aufräumen: Der Runner entfernt die Container `$pg` und `$tests` mit `docker rm -fv`;
      die Suche aus *Bestand* über `tools/` und `harness/mk/` ist wiederholt und belegt
      (Fundstelle, Flags, `.Config.Volumes` des Images, auch für `a-check` und `d-migrate`).
      Beleg in §7: Zahl der Volumes dieses Laufs vor und nach einem Lauf von `make test-integration`
      (nach Namen, §6, nicht als Zahl auf dem Host) und die Fundstellen.
- [ ] Aufräumen in `test-container`: `run-container.sh` entfernt `$pg` an beiden Stellen und alle übrigen
      Container des Laufs mit `rm -fv`; `container-gegenprobe.sh` hat je Ende (grün, rot, `SIGINT`,
      `SIGTERM`) einen Fall `p9-*-volume` an der Attrappe, die ein Volume für `$pg` anlegt und nur mit `-v`
      entfernt, und die Mutationen `rm -f` an je einer der zwei Stellen, selbst rot gesehen (`SPEC-055`
      Punkt 9); der Fall `p7-pg-vor-replay` zieht die neue Zeile mit. Beleg in §7: Namen der Volumes vor
      und nach einem Handlauf von `make test-container`. Dieser Punkt wächst in Liefer-Punkt 1 (Aufräumen)
      hinein, er zählt nicht als vierter.
- [ ] Prüfung im Runner: Vor dem Entfernen liest der Runner die Namen der Volumes seiner
      Container (außer `$daten`), nach dem Entfernen prüft er, dass keines mehr existiert;
      sonst nennt die Meldung die Namen, der Lauf endet mit Exit-Code ungleich 0 (auch bei
      grünen Tests) und entfernt die Reste per Namen (§6). Die Funktionen liegen in einer
      Datei, die Runner und Gegenprobe laden.
- [ ] Gegenprobe: `make integration-volumes-gegenprobe` hängt an `GATE_CHECKS` und fährt
      mit realem Docker die neun Fälle aus `SPEC-057` Punkt 7 (Grundfall, Mutant ohne `-v` rot mit dem
      Namen des Volumes, Container nur angelegt, kein Container, fremde Volumes bleiben, Rest macht den
      Lauf rot, Vorrang des Codes, `SIGTERM`, `SIGINT`); je Fall die Namen der Probe vor und nach,
      nicht die Zahl auf dem Host; je Zusage die Mutation selbst rot gesehen, im Bericht
      *Zusage · Mutation · roter Fall*. Die Zeile steht in `harness/README.md` §Sensors, der Vertrag in `SPEC-057`; das Ziel steht in
      `GATE_SCHNELL_AUSNAHME` (`harness/mk/gates-schnell.mk`), und `SPEC-052` Punkt 13 trägt den
      Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest (§1, Sendungen).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-container.sh` | update | `rm -fv` an den zwei Stellen für `$pg` (`:198`, im Aufräumen `:37`); eine Zeile je Stelle (`SPEC-055` Punkt 9) |
| `tools/harness/container-gegenprobe.sh` | update | Attrappe legt für `$pg` ein Volume an und entfernt es nur mit `-v`; Fälle `p9-*-volume` je Ende; Mutationen `rm -f`; die Erwartung von `p7-pg-vor-replay` (`^rm -f …-pg$`) nennt die neue Form |
| `tools/test/run-integration-tests.sh` | update | `rm -fv` statt `rm -f` für `$tests` und `$pg`; ruft die Prüfung aus der geladenen Datei nach dem Entfernen; die Namen liest sie unmittelbar vor dem Entfernen (§6) |
| `tools/test/integration-aufraeumen.sh` | neu | Funktionen zum Lesen der Namen, Entfernen und Prüfen; Runner und Gegenprobe laden sie, damit die Gegenprobe den Code des Runners prüft, nicht eine Kopie |
| `tools/test/integration-volumes-gegenprobe.sh` | neu | Fälle aus §6 *Mutationen* gegen reales Docker, Mutant per Ersetzung in einer Kopie der Funktionsdatei (das Skript bricht ab, wenn die Ersetzung nichts oder mehr als eine Stelle trifft); Zahl der Volumes der Probe nach Namen |
| `harness/mk/integration.mk` | update | Ziel `integration-volumes-gegenprobe`, an `GATE_CHECKS`, hinter `test-integration` |
| `harness/mk/gates-schnell.mk` | update | `integration-volumes-gegenprobe` in `GATE_SCHNELL_AUSNAHME` (`SPEC-052` Punkt 4; sonst ist `gates-schnell-gegenprobe` rot), im Commit des Ziels |
| `harness/vertraege/SPEC-052-gate-kette.md` | geschrieben (Architect) | Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest; Punkt 4 gilt nach dem Namen, ohne Text; liegt vor dem Anspruch und ist nicht Teil des Diffs |
| `harness/README.md` | update | Zeile in §Sensors für das Ziel; die Zeile von `make test-integration` nennt die Prüfung; die Zeile von `gates-schnell` nennt die Ausnahmeliste unverändert (kein Gate-Name dort) |
| `harness/vertraege/SPEC-057-integration-aufraeumen.md`, `harness/vertraege/README.md` | geschrieben (Architect) | der Vertrag mit den Entscheidungen aus §6, 9 Fälle und Mutationen in Punkt 7; weicht der Code ab, zieht der Implementer den Vertrag im selben Commit nach (`AGENTS.md` §3.9) |

- Der Implementer misst die Images von `a-check` und `d-migrate` (`docker image inspect`, Feld
  `.Config.Volumes`) für den Beleg in §7; die Messung ändert nichts am Umfang.
- Alle Proben (Container, Volumes, Netze) laufen unter einem eigenen Namenspräfix und werden
  per Name entfernt; kein `prune`, kein Eingriff in fremde Objekte.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` setzt `Verantwortlich:`; `next` → `in-progress`): Die Slices von
welle-v1-abschluss liegen in `done/` (WIP-Limit 1); dieser Slice steht direkt vor
`slice-harness-abdeckung-gate`, danach `slice-harness-coverage`, die Slices von
welle-erster-release und welle-v1-nachtrag und die übrige Harness-Reihe (Drift-Log der
Roadmap vom 2026-10-12; im Release-Lauf steht er hinter Schritt 44 von welle-v1-abschluss §5).
Technisch hängt er an nichts. Begründung des Platzes: Der Fehler
kostet Platte und keinen Test, und die Reihenfolge der Wellen-Slices steht (§5 von
[welle-v1-abschluss](../welle-v1-abschluss.md)); er läuft vor den Slices, die nach der Welle
am meisten Gate-Läufe brauchen. Das Leck beträgt rund 75 MB je Lauf von `make gates` (172,9 GB
in zehn Tagen), die Platte hat nach der Bereinigung Reserve. Meldet `make test-integration`
vor dem Start wieder „PostgreSQL nicht bereit“ bei voller Platte, zieht der Planner den Slice
vor die laufende Welle (Umplanung im Drift-Log, Entscheidung des Nutzers). Vor dem ersten
Code-Commit hat der Architect die Randformen aus §6 am 2026-10-17 in `SPEC-057` entschieden
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Größe** (`AGENTS.md` §3.14): Metrik ist die Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools/test tools/harness/run-container.sh
tools/harness/container-gegenprobe.sh harness/mk/integration.mk harness/mk/gates-schnell.mk harness/README.md harness/vertraege`, Tests eingeschlossen. Basis: der Commit des Anspruchs (`next` →
`in-progress`), beim Anspruch einzutragen. Schätzung (Arbeitswert des Planners aus den
Bestandsmaßen; vom Architect am 2026-10-17 auf 350 bis 600 korrigiert, Mitte 470, Obergrenze 600; mit der Übernahme von `test-container` am selben Tag um 40 bis 70 auf 390 bis 670, Mitte 530, Obergrenze 670 (zwei Zeilen im Gate, Attrappe mit Volume rund 15, vier Fälle und zwei Mutationen rund 30 bis 50 Zeilen Gegenprobe, Gate-Index): Die Marke 900 und die Meldeschwelle 750 bleiben,
die Obergrenze liegt 80 unter der Meldeschwelle:
Die Verträge `SPEC-057` und `SPEC-052` liegen mit dem Anspruch schon im Bestand und zählen nicht;
Produktcode rund 150 Zeilen (Runner, Funktionsdatei, Fragmente, Gate-Index), Gegenprobe mit dem
Doppelten des Produktcodes und neun Fällen rund 200 bis 300): Die Meldeschwelle bleibt bei 750.
**Meldeschwelle 750:** Wer im Lauf darüber liegt oder
liegen wird, meldet es dem Planner, bevor er weiterbaut. **Marke 900:** darüber gilt die
Rückführung unten, ohne Ermessen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe der Metrik liegt über 900,
  oder die Gegenprobe braucht mehr als die
  neun Fälle aus `SPEC-057` Punkt 7 oder eine eigene Laufzeitumgebung und ist nicht in einer
  Review-Sitzung prüfbar; Schnitt dann: Aufräumen und Prüfung im Runner hier, die Gegenprobe
  mit dem Ziel in einem eigenen Slice direkt dahinter.
- `in-progress` → `open` (blockiert — Carveout?): `docker inspect` nennt die Volumes eines
  Containers nicht verlässlich (im Zustand `created`, auf einem Docker, den der Nutzer
  einsetzt; `SPEC-057` Punkt 8 fängt es in der Gegenprobe ab); dann entscheidet der Architect einen anderen Weg zu den Namen (etwa Volumes
  unter eigenem Namen in `docker run`), bevor der Slice weitergeht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit `test-integration` und `integration-volumes-gegenprobe`,
die Zahlen vor und nach einem Lauf in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — der Slice liefert einen neuen Vertrag (die Prüfung im
Runner, die Gegenprobe als Gate). **Entschieden vom Architect am 2026-10-17 in
[`SPEC-057`](../../../../harness/vertraege/SPEC-057-integration-aufraeumen.md)** (Ort nach `MR-002`, nicht
`SPEC-038`: Der Runner ist ein Harness-Werkzeug); die Empfehlungen des Planners sind bis auf die
Fälle unten übernommen. Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie
zurück.

| Randform | Entscheidung | Stelle |
|---|---|---|
| Abbruch mitten im Lauf | `SIGINT` und `SIGTERM` laufen durch die Falle `EXIT`, je ein Fall (Empfehlung war nur `SIGTERM`: §3.10 zählt je Element); `SIGKILL`, Ausfall des Hosts, Docker-Dienst weg: nicht zugesagt | Punkt 5 |
| Parallele Läufe | nur Namen aus `docker inspect` der eigenen zwei Container, nie Zahl oder Filter des Hosts; Fall *fremde Volumes* | Punkt 4, Fall `p4-fremd` |
| Fehlgeschlagener Start | Namen werden vor dem Entfernen in jedem Zustand gelesen; kein Container ergibt eine leere Liste ohne Fehler | Punkt 3 |
| Welche Volumes zählen | Mounts vom Typ `volume` der zwei Container außer `$daten`; Herkunft, nicht Name | Punkt 1 |
| Rest nach dem Entfernen | Namen melden (kein Inhalt), per Namen entfernen, Exit 1; ein Code ungleich 0 des Laufs behält Vorrang (Fall `p4-vorrang`, neu) | Punkt 4 |
| Docker-Versionen | gemessen 29.9.0; eine blinde Version macht `p2-ohne-v` grün und die Gegenprobe rot, sonst nicht zugesagt | Punkt 8 |
| Image mit ungeprüftem `VOLUME` | erfasst, solange es als Container des Runners läuft; `docker run` außerhalb: Grenze | Punkt 1 und 11 |
| Leere Namensliste | kein Fehler, keine Meldung (Fall `p3-kein-container`) | Punkt 3 |
| Ausnahmeliste und Beleg (Sendung aus §1) | `integration-volumes-gegenprobe` steht in `GATE_SCHNELL_AUSNAHME` (der Implementer trägt die Zeile ein); Beleg in `SPEC-052` Punkt 13 geschrieben | Punkt 10, `SPEC-052` Punkt 13 |
| Aufrufarten von `make` (`-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`) | kein Stand, nur Exit-Code; akzeptiertes Negativ mit Grund | Punkt 9 |

**Mutationen** (`AGENTS.md` §3.10): die Tabelle *Fall · Zusage · Mutation* steht in `SPEC-057`
Punkt 7; der Implementer sieht jede Mutation selbst rot und nennt sie im Bericht als
*Zusage · Mutation · roter Fall*. Hinweis zu den Kosten: Die Fälle `p5-*` starten den Runner bis
zum Container `$pg` (wenige Sekunden); alle übrigen kommen mit `docker create` ohne Start aus.

**Risiken:**

- Die Gegenprobe verlängert `make gates` (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, 1×):
  Containerfälle mit realem Docker kosten Sekunden je Fall — **Ausgang:** offen bis Closure; die
  Laufzeit des Ziels steht mit Zahl in §7, ohne Lockerung von Fällen.
- Die Platte läuft vor dem Start dieses Slice wieder voll, weil jeder Lauf von `make gates`
  rund 75 MB hinterlässt, bis der Slice in `done/` ist — **Ausgang:** offen bis Closure; der Trigger
  in §4 zieht den Slice vor, wenn es eintritt.
- Ein neues `docker run` ohne `--rm` und ohne `-v` außerhalb des Runners legt das Leck wieder an;
  kein Sensor sieht es, die Gegenprobe prüft nur den Runner — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`,
gelesen am Stand `79258d3`, 2026-10-12). Treffer:

- `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes` (1×, offen, neu mit diesem Plan) —
  der Gegenstand; die Ursache behebt dieser Slice.
- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (1×, offen) — berührt, weil die Gegenprobe an
  `GATE_CHECKS` hängt (Risiko in §6).
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos
  und liefert einen neuen Vertrag; der Auslöser steht in §4 *Start*.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) und
  `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — neuer
  Vertrag: Mutationen und Randformen in §6.
- `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` (1×, offen) — betrifft denselben Runner;
  der Slice ändert die zwei Phasen nicht.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — der Slice ändert keine
  Konfiguration eines Gates; die Gegenprobe trägt trotzdem Fälle in beide Richtungen.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 3. Schichten nach der Zählung in `AGENTS.md` §3.13: zwei (Harness-Werkzeug mit
`tools/test` und `harness/mk`; Gate-Index und Spezifikation mit der Zeile in `harness/README.md`
und den Verträgen `SPEC-057` und `SPEC-052`, gezählt als eigene Schicht, weil §3.13 die
Harness-Dokumentation nicht ausnimmt; die Verträge sind vor dem Anspruch geschrieben und nicht Teil des Diffs); kein Produkt-Code, keine Nutzer-Doku.

**Nachgezählt am 2026-10-17 mit der Übernahme von `test-container`** (Geber `slice-v1-abschluss-container`
und der Befund aus `slice-v1-abschluss-container-doku`; `SPEC-055` Punkt 9 im Architect-Commit, nicht im
Diff): Liefer-Punkte weiter 3 (die Behebung wächst in Punkt 1 *Aufräumen* und die Fälle in Punkt 3
*Gegenprobe* hinein; sie sind keine eigenen Punkte, erhöhen aber den Umfang, daher die Metrik in §4),
Schichten weiter 2: `tools/harness` ist wie `tools/test` und `harness/mk` die Schicht Harness-Werkzeug, die
Verträge mit dem Gate-Index die zweite. Die Schnittlinie für die Rückführung bleibt: die Gegenprobe von
`test-container` zuerst abgeben, nicht die Behebung.

**Nachgezählt** mit den Sendungen aus §1 (`gates-schnell.mk` und `SPEC-052`; Geber dieser Slice):
Liefer-Punkte weiter 3 (der Eintrag in die Ausnahmeliste und der Beleg in `SPEC-052` wachsen in
Punkt 3, Gegenprobe, hinein; sie sind keine eigenen Punkte, erhöhen aber den Umfang, daher die
Metrik in §4), Schichten weiter 2: beide Dateien liegen in den zwei genannten Schichten
(`harness/mk` bzw. `harness/vertraege`). Die Bindung an `slice-harness-abdeckung-gate` (§1) ändert
den Umfang dieses Slice nicht, sie bindet den Kopf der Gegenprobe.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
