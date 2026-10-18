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
`harness/vertraege/`, Ort nach `MR-002`; vom Architect am 2026-10-17 vor dem Anspruch geschrieben; Punkte 4, 5 und 7 im Lauf nachgezogen, im Diff dieses
Slice); `SPEC-055` (Punkt 9, Absatz *Volumes*: `test-container` entfernt seine Container mit `rm -fv`; Architect am 2026-10-17, vor dem Anspruch geschrieben, nicht im Diff, und mit `6978e8f6` um den Absatz zu `$rec`, `$rep`, `$img`, `version` und `psql` ergänzt, im Diff dieses Slice; Geber `slice-v1-abschluss-container`, geschlossen); `SPEC-052` (Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest;
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
- Weitere `docker run` in `tools/`, `harness/mk/` und den Wurzel-Fragmenten (`d-check.mk`, `a-check.mk`, beide `--rm`): `harness/mk/build.mk:22`,
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

- [x] Aufräumen: Der Runner entfernt die Container `$pg` und `$tests` mit `docker rm -fv`;
      die Suche aus *Bestand* über `tools/` und `harness/mk/` ist wiederholt und belegt
      (Fundstelle, Flags, `.Config.Volumes` des Images, auch für `a-check` und `d-migrate`).
      Beleg in §7: Zahl der Volumes dieses Laufs vor und nach einem Lauf von `make test-integration`
      (nach Namen, §6, nicht als Zahl auf dem Host) und die Fundstellen.
- [x] Aufräumen in `test-container`: `run-container.sh` entfernt `$pg` an beiden Stellen und alle übrigen
      Container des Laufs mit `rm -fv`; `container-gegenprobe.sh` hat je Ende (grün, rot, `SIGINT`,
      `SIGTERM`) einen Fall `p9-*-volume` an der Attrappe, die ein Volume für `$pg` anlegt und nur mit `-v`
      entfernt, und die Mutationen `rm -f` an je einer der zwei Stellen, selbst rot gesehen (`SPEC-055`
      Punkt 9); der Fall `p7-pg-vor-replay` zieht die neue Zeile mit. Beleg in §7: Namen der Volumes vor
      und nach einem Handlauf von `make test-container`. Dieser Punkt wächst in Liefer-Punkt 1 (Aufräumen)
      hinein, er zählt nicht als vierter.
- [x] Prüfung im Runner: Vor dem Entfernen liest der Runner die Namen der Volumes seiner
      Container (außer `$daten`), nach dem Entfernen prüft er, dass keines mehr existiert;
      sonst nennt die Meldung die Namen, der Lauf endet mit Exit-Code ungleich 0 (auch bei
      grünen Tests) und entfernt die Reste per Namen (§6). Die Funktionen liegen in einer
      Datei, die Runner und Gegenprobe laden.
- [x] Gegenprobe: `make integration-volumes-gegenprobe` hängt an `GATE_CHECKS` und fährt
      mit realem Docker die neun Fälle aus `SPEC-057` Punkt 7 (Grundfall, Mutant ohne `-v` rot mit dem
      Namen des Volumes, Container nur angelegt, kein Container, fremde Volumes bleiben, Rest macht den
      Lauf rot, Vorrang des Codes, `SIGTERM`, `SIGINT`); je Fall die Namen der Probe vor und nach,
      nicht die Zahl auf dem Host; je Zusage die Mutation selbst rot gesehen, im Bericht
      *Zusage · Mutation · roter Fall*. Die Zeile steht in `harness/README.md` §Sensors, der Vertrag in `SPEC-057`; das Ziel steht in
      `GATE_SCHNELL_AUSNAHME` (`harness/mk/gates-schnell.mk`), und `SPEC-052` Punkt 13 trägt den
      Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest (§1, Sendungen).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-container.sh` | update | `rm -fv` an den zwei Stellen für `$pg` (`:198`, im Aufräumen `:37`); eine Zeile je Stelle (`SPEC-055` Punkt 9) |
| `tools/harness/container-gegenprobe.sh` | update | Attrappe legt für `$pg` ein Volume an und entfernt es nur mit `-v`; Fälle `p9-*-volume` je Ende; Mutationen `rm -f`; die Erwartung von `p7-pg-vor-replay` (`^rm -f …-pg$`) nennt die neue Form |
| `tools/test/run-integration-tests.sh` | update | `rm -fv` statt `rm -f` für `$tests` und `$pg`; ruft die Prüfung aus der geladenen Datei nach dem Entfernen; die Namen liest sie unmittelbar vor dem Entfernen (§6); neben der Falle für `TERM` eine für `INT` (Exit-Code 130), weil bash ein Signal sonst verwirft, solange ein Kind (`docker`, Exit-Code 2) läuft (`SPEC-057` Punkt 5) |
| `tools/test/integration-aufraeumen.sh` | neu | Funktionen zum Lesen der Namen, Entfernen und Prüfen; Runner und Gegenprobe laden sie, damit die Gegenprobe den Code des Runners prüft, nicht eine Kopie |
| `tools/test/integration-volumes-gegenprobe.sh` | neu | Fälle aus §6 *Mutationen* gegen reales Docker, Mutant per Ersetzung in einer Kopie der Funktionsdatei (das Skript bricht ab, wenn die Ersetzung nichts oder mehr als eine Stelle trifft); Zahl der Volumes der Probe nach Namen |
| `tools/harness/gates-schnell-gegenprobe.sh` | update | die Liste `LANGSAM` der Kopien nennt `integration-volumes-gegenprobe`; ohne sie ist die Ausnahmeliste der Kopie ein Eintrag ohne Gate (`p4-eintrag`, im Gate-Lauf rot gesehen) |
| `tools/test/postgres-versionen-gegenprobe.sh` | update | Gegenprobe von `SPEC-053`: die Kopie des Runners trägt jetzt `integration-aufraeumen.sh` (Fixture); die Mutanten der Zeile `rm -f "$tests" "$pg"` und die Erwartung des Falls `p7-vor-naechster` (`^rm -fv …`) folgen der Funktionsdatei |
| `harness/mk/integration.mk` | update | Ziel `integration-volumes-gegenprobe`, an `GATE_CHECKS`, hinter `test-integration` |
| `harness/mk/gates-schnell.mk` | update | `integration-volumes-gegenprobe` in `GATE_SCHNELL_AUSNAHME` (`SPEC-052` Punkt 4; sonst ist `gates-schnell-gegenprobe` rot), im Commit des Ziels |
| `harness/vertraege/SPEC-055-container-und-binaries.md` | Punkt 9 (Architect, vor dem Anspruch; Absatz zu `$rec`, `$rep`, `$img`, `version`, `psql` mit `6978e8f6`, im Diff) | Volumes: `rm -fv` für alle Container des Laufs; die übrigen ohne Fall als akzeptiertes Negativ |
| `harness/vertraege/SPEC-052-gate-kette.md` | geschrieben (Architect) | Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest; Punkt 4 gilt nach dem Namen, ohne Text; liegt vor dem Anspruch und ist nicht Teil des Diffs |
| `harness/README.md` | update | Zeile in §Sensors für das Ziel; die Zeile von `make test-integration` nennt die Prüfung; die Zeile von `gates-schnell` nennt die Ausnahmeliste unverändert (kein Gate-Name dort) |
| `harness/vertraege/SPEC-057-integration-aufraeumen.md`, `harness/vertraege/README.md` | geschrieben (Architect), Punkte 5 und 7 vom Implementer nachgezogen | der Vertrag mit den Entscheidungen aus §6, 9 Fälle und Mutationen in Punkt 7; nachgezogen (`AGENTS.md` §3.9), weil der Code abweicht: Falle für `INT` im Runner, Mutant `p5-sigterm` (Falle für `TERM` entfällt) ist äquivalent und durch `kill -KILL` ersetzt, die Hilfsfunktion `integration_ende` trägt Vorrang des Codes (Punkte 4 und 5), weitere Mutanten je Element |

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
[welle-v1-abschluss](welle-v1-abschluss.md)); er läuft vor den Slices, die nach der Welle
am meisten Gate-Läufe brauchen. Das Leck beträgt rund 75 MB je Lauf von `make gates` (172,9 GB
in zehn Tagen), die Platte hat nach der Bereinigung Reserve. Meldet `make test-integration`
vor dem Start wieder „PostgreSQL nicht bereit“ bei voller Platte, zieht der Planner den Slice
vor die laufende Welle (Umplanung im Drift-Log, Entscheidung des Nutzers). Vor dem ersten
Code-Commit hat der Architect die Randformen aus §6 am 2026-10-17 in `SPEC-057` entschieden
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Größe** (`AGENTS.md` §3.14): Metrik ist die Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools/test tools/harness/run-container.sh
tools/harness/container-gegenprobe.sh harness/mk/integration.mk harness/mk/gates-schnell.mk harness/README.md harness/vertraege`, Tests eingeschlossen. Basis: `de908591` (Elternteil des Anspruchs-Commits `5fba0671`, `next` →
`in-progress`). Schätzung (Arbeitswert des Planners aus den
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

**Prüfung des Architect nach der Umsetzung (2026-10-17, Rückgaben des Implementers aus §7):**

1. Zweites `SIGTERM` während des Aufräumens: **nicht zugesagt**, akzeptiertes Negativ mit Grund in `SPEC-057` Punkt 5
   (Kosten eines Falls mit langsamem `docker rm` stehen nicht im Verhältnis; ein Stopper, der zweimal sendet, endet bei
   `SIGKILL` ohnehin ohne Aufräumen). Der Mutant „Falle für `TERM` entfällt“ ist äquivalent und steht so in Punkt 7; zugesagt ist das
   Ende, nicht die Falle.
2. Falle für `INT`: Folge der Zusage von Punkt 5, keine neue Randform; der Nachtrag stimmt mit dem Code (`trap 'trap "" INT; exit 130' INT`)
   und mit `p5-sigint` überein. Der Mutant `p5-int-entfaellt` blieb grün (der Runner steht im Fall im `wait`): Die Zusage „nötig“ entfällt, die Falle ist nicht zugesagt (akzeptiertes Negativ in `SPEC-057` Punkt 5, Nachtrag 2026-10-17).
3. Punkt 7 nennt je Mutation die Menge der roten Fälle statt „nur ihn“ (Menge aus der Tabelle in §7 übernommen); die Gegenprobe verlangt genau diese.
4. `SPEC-057` Punkt 8 nennt die nur auf 29.9.0 gemessene Version jetzt als akzeptiertes Negativ mit Grund.
5. `d-check.mk` und `a-check.mk` tragen `--rm` (geprüft); da sie nicht geändert werden, genügt die Nachführung in §7, §3 bleibt.
6. `SPEC-055` Punkt 9 stimmt mit dem Diff (`rm -fv` für alle Container, `$pg` zweimal gemutet); die übrigen Container ohne Fall sind als akzeptiertes
   Negativ nachgetragen. `SPEC-052` Punkt 13 stimmt (Gegenprobe liest nur `tools/test/` in Kopien).

7. Review F-1/F-2/F-4 (2026-10-17): „Exit 1 bei Rest“ ist an der Funktion `integration_ende` zugesagt, die Verdrahtung im Runner ist akzeptiertes Negativ
   (`SPEC-057` Punkt 4); die Spalte „rote Fälle“ gilt je Gruppe der Mutation (`nur`), Punkt 7; „kein Inhalt des Volumes“ ist Regel an den Schreiber, nicht zugesagt.

**Rückgabe an den Implementer:** `tools/test/integration-volumes-gegenprobe.sh`: Der Mutant `p5-int-entfaellt` bleibt ohne roten Fall (wie `p5-term-entfaellt`, äquivalent);
die Zeile `p3-nach-entfernen` verlangt die roten Fälle `p3-angelegt`, `p4-rest-rot`, `p5-sigint`, `p5-sigterm` (`SPEC-057` Punkt 7). Die Gegenprobe prüft sonst die Mengen der Tabelle in Punkt 7 (`p2-grundfall`+`p2-ohne-v`, `p3-angelegt`+`p4-rest-rot` usw.); weicht ein Mutant ab, melden.

**Risiken:**

- Die Gegenprobe verlängert `make gates` (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, 1×):
  Containerfälle mit realem Docker kosten Sekunden je Fall — **Ausgang:** weiter offen, ins Register
  (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, Beleg `evidence/slice-harness-integration-volumes.md`); gemessen
  42 s für das Ziel, 1038 s für `make gates` am Endstand (§7), ohne Lockerung von Fällen.
- Die Platte läuft vor dem Start dieses Slice wieder voll, weil jeder Lauf von `make gates`
  rund 75 MB hinterlässt, bis der Slice in `done/` ist — **Ausgang:** entfallen: Das Leck ist behoben und
  gemessen (`make test-integration`, `make test-container` und `make gates` je kein neues Volume nach Namen, §7); der
  Trigger in §4 trat nicht ein. Bestehende Volumes bleiben Sache des Nutzers (§1).
- Ein neues `docker run` ohne `--rm` und ohne `-v` außerhalb des Runners legt das Leck wieder an;
  kein Sensor sieht es, die Gegenprobe prüft nur den Runner — **Ausgang:** weiter offen, ins Register
  (`BEO-REPO/docker-run-ohne-rm-ausserhalb-der-laeufer`, neu angelegt); Grenze (1) in `SPEC-057` Punkt 11.

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

### Closure-Notiz

Stand: Planner am 2026-10-17, vor dem `git mv` nach `done/`. Review (3 MEDIUM, 2 LOW, 1 INFO) und Verifikation (bestanden, 2 LOW) liegen unter `docs/reviews/` (`2026-10-17-review-slice-harness-integration-volumes.md`, `2026-10-17-verify-slice-harness-integration-volumes.md`); alle Findings sind bearbeitet.

- **Was hat funktioniert:** Namen statt Zahlen: Die Prüfung liest die Volumes der eigenen Container (`docker inspect`, auch im Zustand `created`) und erfasst so nur, was der Lauf selbst anlegte; die Gegenprobe lädt dieselbe Funktionsdatei wie der Runner und mutiert sie in Kopien. Vor der Behebung hinterließen beide Läufer je ein Volume, danach keines (`make gates` am Endstand: 49 vor, 49 nach). Der Architect hatte die Randformen vor dem ersten Code in `SPEC-057` entschieden.
- **Was ging anders als geplant:** Der Slice brauchte mehrere Runden. Das Review fand zwei Zusagen ohne Mutation (Verdrahtung `exit "$ende"`, Inhalt der Meldung) und zwei Aussagen weiter als die Prüfung (Menge roter Fälle, Gate-Zeile zu äquivalenten Mutanten); der Verifier fand sie nicht, weil er Zusage für Zusage am Mutanten prüfte, nicht die Verdrahtung des Aufrufers. Die Runden dehnte, dass der Implementer Vertragstext selbst nachzog (Punkte 5 und 7 von `SPEC-057`); die Entscheidung gehört dem Architect.
- **Ausgänge der Findings:** F-1 (Exit 1 im Runner) behoben als Zusage an der Funktion `integration_ende` und akzeptiertes Negativ für die Verdrahtung, `SPEC-057` Punkt 4. F-2 (Mengen roter Fälle) behoben: Die Spalte gilt je Gruppe der Mutation, `SPEC-057` Punkt 7. F-3 (Gate-Zeile, Kopfkommentar) behoben in `51205010`. F-4 (Inhalt des Volumes) Regel an den Schreiber der Meldung, nicht zugesagt, `SPEC-057` Punkt 4. F-5 (Plan folgt) behoben in `51205010`. F-6 (INFO) ohne Handlung. Verifikation F-1 (Basis der Marke) und F-2 (offene Punkte in §7) behoben.
- **Steering-Loop-Eintrag:** Neuer Sensor `make integration-volumes-gegenprobe` mit der Prüfung nach Namen im Runner — liegt in `tools/test/integration-aufraeumen.sh`. Auslöser: `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes` (slice-walking-skeleton-record, slice-v1-abschluss-container-doku, slice-harness-integration-volumes — 3×, Ausgang verkörpert). Die Finding-Klasse „Zusage ohne Mutation an der Verdrahtung des Aufrufers“ geht als Wiederauftreten in `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (Regel `AGENTS.md` §3.10 gilt, keine Änderung von diesem Lauf; ob §3.10 die Verdrahtung ausdrücklich nennt, entscheidet der Nutzer).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-harness-integration-volumes.md` ergänzt in `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes/` (Zähler 3×, Stand verkörpert), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (48×), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (73×), `BEO-REPO/plan-folgt-korrektur-nicht/` (61×) und `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` (7×); neu angelegt `BEO-REPO/docker-run-ohne-rm-ausserhalb-der-laeufer/` (1×, offen).
- **Folge-Slices:** keiner. Die Sendung an `slice-harness-abdeckung-gate` (Kopf der Gegenprobe ohne Deklaration) steht in dessen §1 (`grep` findet die Kennung dieses Slice); er liegt in `next/`.
- **Risiken aus §6:** Laufzeit weiter offen (Register), volle Platte entfallen (Leck behoben und gemessen), neues `docker run` ohne `--rm` weiter offen (Register); Begründungen in §6.
- **Marke (`AGENTS.md` §3.14):** Summe aus Einfügungen und Löschungen von `git diff --numstat de908591..51205010` über die Pfade aus §4: **519** (+474, −45) von 900; Meldeschwelle 750 nicht erreicht.
- **`make gates` am Endstand:** grün, rc=0, 17 min 18 s, am Stand `51205010`; seither ändert sich nur `docs/plan/planning/` (inert).
- **Drei Paarungen:** Anker (`liegt in` oben, Zielort liegt vor) · Folge-Slice (keiner zugewiesen; Sendung an `slice-harness-abdeckung-gate` von ihm angenommen) · Register (fünf Belege ergänzt, ein Eintrag neu).

### Belege des Implementers

Stand: Implementer-Lauf am 2026-10-17, Basis `de908591` (Elternteil des Anspruchs `5fba0671`), Code-Commit `5e15247b`.
Namen statt Zahlen: Die Listen der Volumes (`docker volume ls -q | sort`) vor und nach dem Lauf, Vergleich mit
`comm`. Auf dem Host liefen fremde Prozesse: Zwischen zwei Messungen verschwanden fünf Volumes, die nicht aus
diesem Lauf stammen (anonyme Volumes eines `--rm`-Containers anderer Herkunft); sie zählen nicht, geprüft wird
„neu hinzugekommen“ (`comm -13`).

**Aufräumen im Runner (DoD 1), Handlauf `make test-integration`:** vor 54 Volumes, nach 54, neu: 0 (`comm -13` leer);
Lauf grün (`run-integration-tests: gruen`, Exit 0, 2 min 20 s). Gegenstück vor der Behebung: der Runner aus `5fba0671`
(`git show 5fba0671:tools/test/run-integration-tests.sh`, Stand vor dem Code-Commit) hinterließ bei demselben
Lauf genau ein Volume (`a2035760fb25…`, 54 → 55), das der Implementer danach per Namen entfernte.

**Aufräumen in `test-container` (DoD 2), Handlauf `make test-container`:** vor 49, nach 49, neu: 0; Lauf grün (1 min 55 s).
Gegenstück: `run-container.sh` aus `5fba0671` hinterließ ein Volume (`dbd4b69c7873…`, 49 → 50), per Namen entfernt.

**Fundstellen (Suche wiederholt, Stand dieses Laufs):** `docker run`/`docker create` in `tools/`, `harness/mk/` und den
Wurzel-Fragmenten:
`harness/mk/build.mk:22`, `harness/mk/schema.mk:9` und `:12`, `tools/arch/a-check-negativ.sh:83`,
`tools/harness/traeger-fetch.sh:146`, `tools/harness/run-container.sh:169` (`psql`), `:179` (`version`), `d-check.mk` (alle
`run`-Zeilen), `a-check.mk:16` und `:19` tragen `--rm`; ohne `--rm` und ohne `-v` waren `tools/test/run-integration-tests.sh`
`$pg` (`run -d`, jetzt `rm -fv`), `run-container.sh` `$pg` (`:183`), `$rec` (`:189`), `$rep` (`:200`, `create`), `$img` (`:143`, `create`).
Die Plan-Liste nannte `d-check.mk` und `a-check.mk` im Wurzelverzeichnis nicht; beide tragen `--rm`.
`.Config.Volumes` (`docker image inspect`):
`ghcr.io/pt9912/a-check@sha256:e820…` null · `ghcr.io/pt9912/d-check:v0.86.1` null · `golang:1.27-alpine` null ·
`pgwire-recorder:integration` null · Produkt-Image (`--target runtime`, aus dem `Dockerfile` gebaut) null (`grep VOLUME Dockerfile` ohne
Treffer) · `postgres:17-alpine@sha256:b0f9…` `{"/var/lib/postgresql/data":{}}` ·
`ghcr.io/pt9912/d-migrate@sha256:af9d…` `{"/work":{}}`. Letzteres ist die Ausnahme der Annahme „null“: Das Rezept von
`schema-validate` bindet `tools/schema` nach `/work` (Bind-Mount über dem deklarierten `VOLUME`), und ein Handlauf
`make schema-validate` ließ kein neues Volume zurück (Namen vor und nach, neu: 0). `--rm` genügt auch dort.
Bei `$rec`, `$rep`, `$img` (Produkt-Image, `.Config.Volumes` null) legt Docker kein Volume an; `rm -fv` an diesen
Stellen ist daher in der Attrappe nicht unterscheidbar (Attrappe legt nur für `$pg` ein Volume an): kein Mutant dafür, aus diesem Grund
(gemessen: null) und weil `SPEC-055` Punkt 9 „Nur `$pg` trägt ein `VOLUME`“ nennt.

**Zusage · Mutation · roter Fall.** Zu `SPEC-057` (Gegenprobe `integration-volumes-gegenprobe`, `GP_ENTDECKE=1`,
Mutant per Textersatz genau einer Stelle in einer Kopie von `integration-aufraeumen.sh` bzw. `run-integration-tests.sh`):

| Zusage | Mutation | rote Fälle (gesehen) |
|---|---|---|
| P1 Namen aus `$tests` und `$pg`: nur `$pg` | `for c in "$tests" "$pg"` → `"$pg"` (`p1-nur-pg`) | `p3-angelegt` |
| P1 … nur `$tests` | → `"$tests"` (`p1-nur-tests`) | `p3-angelegt` |
| P1 außer `$daten` | Ausnahme entfällt (`p1-mit-daten`) | `p2-grundfall` |
| P1 Mounts vom Typ `volume` (Bind-Mounts zählen nicht) | `.Type`-Prüfung entfällt (`p1-ohne-typ`) | keiner: äquivalent, ein Mount ohne Volume trägt keinen Namen, die Namensliste bleibt gleich |
| P2 `rm -fv` für `$tests` und `$pg` | `rm -f` (`p2-ohne-v`) | `p2-grundfall`, `p2-ohne-v` |
| P2 … je Element | `$tests` entfällt (`p2-ohne-tests`), `$pg` entfällt (`p2-ohne-pg`) | je `p2-grundfall`, `p2-ohne-v` |
| P3 Namen vor dem Entfernen, auch Zustand `created` | erst nach dem Entfernen lesen (`p3-nach-entfernen`) | `p3-angelegt`, `p4-rest-rot`, `p5-sigint`, `p5-sigterm` |
| P3 kein Container: leere Liste ohne Fehler und Meldung | fehlender Container wird Fehler (`p3-leer-fehler`) | `p3-kein-container`, `p5-sigint`, `p5-sigterm` |
| P4 Rest nur per Namen entfernt | entfernt jedes Volume ohne Container, beschränkt auf die Namen der Probe (`p4-alles-entfernen`) | `p4-fremd` |
| P4 Rest wird entfernt | Entfernen entfällt (`p4-nicht-entfernen`) | `p3-angelegt` |
| P4 Meldung nennt die Namen | Meldung ohne Namen (`p4-ohne-namen`) | `p3-angelegt` |
| P4 Rest macht den Lauf rot (Exit 1) | Prüfung endet mit 0 (`p4-pruefen-0`) · Code der Prüfung verworfen (`p4-code-verworfen`) | je `p3-angelegt`, `p4-rest-rot` |
| P4 Code des Laufs ungleich 0 behält Vorrang | Rückgabe immer 1 (`p4-immer-1`) | `p4-vorrang` |
| P5 Falle `EXIT` | `trap aufraeumen EXIT` entfällt (`p5-exit`) | `p5-sigint`, `p5-sigterm` |
| P5 Prüfung läuft in der Falle | Falle ruft nur das Entfernen (`p5-ohne-pruefung`) | `p5-sigint`, `p5-sigterm` |
| P5 Exit-Code des Laufs (143, 130) | `exit 0` am Ende der Falle (`p5-ausgang-0`) | `p5-sigint`, `p5-sigterm` |
| P5 `SIGTERM` | Falle für `TERM` beendet mit `kill -KILL` (`p5-term`) | `p5-sigterm` |
| P5 `SIGINT` | Falle für `INT` beendet mit `kill -KILL` (`p5-int`) | `p5-sigint` |
| P5 `SIGTERM` | Falle für `TERM` entfällt (`p5-term-entfaellt`, Mutant aus der Tabelle von `SPEC-057` Punkt 7) | keiner: äquivalent. bash führt die Falle `EXIT` auch ohne Falle für `TERM` aus und endet mit 143 (gemessen mit bash 5.2.21, 8 Läufe von `p5-sigterm` und ein Handversuch); `SPEC-057` Punkt 7 ist nachgezogen |
| P5 `SIGINT` | Falle für `INT` entfällt (`p5-int-entfaellt`) | keiner: äquivalent im Fall, der Runner steht im `wait` und endet mit Aufräumen und 130 |

Der Grundfall hat „—“ als Mutation und wird von jedem Mutanten rot, der das Entfernen trifft (`p2-ohne-v` und `p2-grundfall` beide);
das ist der Wortlaut „nur sie“ nicht, aber die Zusage („nach dem Aufräumen existiert kein gelesenes Volume“) ist die des Grundfalls.
`p3-angelegt` hängt für `p3-nach-entfernen` an derselben Prüfung wie `p4-rest-rot` (beide brauchen die gelesenen Namen für den Rest).
Zu `SPEC-055` Punkt 9 (`container-gegenprobe`, `GP_MUT=… GP_ENTDECKE=1`):

| Zusage | Mutation | rote Fälle (gesehen) |
|---|---|---|
| Aufräumen am Ende entfernt `$pg` mit Volume | `rm -f` statt `rm -fv` in der Zeile des Aufräumens (`p9-volume-aufraeumen`) | `p9-sigINT-psql-volume`, `p9-sigTERM-psql-volume` |
| Entfernen von `$pg` vor dem Wiedergeben mit Volume | `rm -f` statt `rm -fv` an der Stelle (`p9-volume-replay`) | `p7-pg-vor-replay`, `p9-gruen-volume`, `p9-rot7-volume` |
| je Ende ein Fall `p9-*-volume` | die vorhandenen Mutanten des Aufräumens färben den Fall des Endes mit: Zeile des Aufräumens entfällt (`p9-container`), Falle `EXIT` entfällt (`p9-exit`), Fallen `INT`/`TERM` ohne Aufräumen (`p9-int`, `p9-term`) | `p9-sigINT-psql-volume` · `p9-sigTERM-psql-volume` (je nach Mutant), mit den bisherigen roten Fällen |
| `p7-pg-vor-replay` kennt die neue Zeile | Erwartung `^rm -fv …-pg$`; Mutant `p7-rm-pg` (Zeile entfällt) | `p7-pg-vor-replay` |

Die Fälle `p9-*-volume` der Enden `rot4` bis `rot6` und der `sig*-build`-Enden sind vom Aufbau her trivial grün (der Container `$pg` ist
dort noch nicht angelegt); die Mutanten fangen nur Enden mit lebendem `$pg` (`sig*-psql`) und das Ende nach dem Wiedergeben (`gruen`, `rot7`).

**Befund im Runner (vor dem Code gemessen, nicht in §6):** Die erste Fassung der Fälle `p5-*` schickte `SIGINT` an den Runner (nur den
Prozess); der Lauf ging weiter. Messung: Ein Kind, das nicht am Signal endet (das Go-Programm `docker` endet auf `SIGINT` mit Exit-Code 2),
lässt bash nach `SIGINT` fortsetzen; ohne Falle für `INT` geht ein `SIGINT` an einen Runner im Vordergrundkommando verloren (an den Prozess
allein hing der erste Lauf; an die Prozessgruppe endete von drei Läufen einer mit Exit-Code 2 durch `set -e` am `docker exec`, einer hing). Mit `trap 'trap "" INT; exit 130' INT` im Runner hielten acht Läufe von `p5` (TERM 143, INT 130). Die Zusage von `SPEC-057` Punkt 5 braucht diese
Falle; der Vertrag ist im selben Commit nachgezogen. Das ist eine Folge der Zusage, keine neue Randform.

**Laufzeit des Ziels** (Risiko in §6): `make integration-volumes-gegenprobe` 42 s (Wanduhr), danach keine Container, Netze oder neuen Volumes
des Präfixes `pgr-ivg-` auf dem Host (Namen vor und nach, neu: 0). `make test-postgres-versionen-gegenprobe` mit der angepassten Kopie: grün.

**Größe der Marke** (`AGENTS.md` §3.14): `git diff --numstat de908591..5e15247b -- tools/test tools/harness/run-container.sh
tools/harness/container-gegenprobe.sh harness/mk/integration.mk harness/mk/gates-schnell.mk harness/README.md harness/vertraege`:
Summe aus Einfügungen und Löschungen **457** (+429, −28); Meldeschwelle 750 und Marke 900 nicht erreicht. Schichten: zwei
(Harness-Werkzeug `tools/test`, `tools/harness`, `harness/mk`; Gate-Index und Vertrag `harness/README.md`, `SPEC-057`); Liefer-Punkte drei. Die Zählung hält.
`tools/test/postgres-versionen-gegenprobe.sh` (6+6 Zeilen) gehört zu `tools/test`.

**Beobachtung für das Register:** Der Plan nannte die Fälle von `SPEC-057` Punkt 7 so, als bräche jede Mutation nur den eigenen Fall; die Zusammenhänge
(Grundfall rot bei `rm -f`, `p4-rest-rot` rot bei „Namen nach dem Entfernen“) und ein äquivalenter Mutant (`TERM`-Falle) standen nicht im Vertrag.

**`make gates` am Endstand** (`dc6da5fa`; Läufe im Hintergrund, Log im Scratchpad; Ergebnis aus Log und `rc`):
1. Lauf rot nach 349 s: `gates-schnell-gegenprobe` meldete „Kopie 'ok' ist nicht grün: p4-eintrag“. Ursache: Die Liste `LANGSAM` in
   `tools/harness/gates-schnell-gegenprobe.sh` kannte `integration-volumes-gegenprobe` nicht, die Ausnahmeliste der Kopie nannte ein Gate ohne
   Eintrag in `GATE_CHECKS`. Behoben in `dc6da5fa` (eine Zeile; die Datei steht nicht in der Pfadliste der Marke, +1/−1).
2. Lauf rot: `test-postgres-versionen-gegenprobe`, Mutant `p7-int-treiber` (rot waren `p7-sigINT-gruppe`, `p7-sigINT-skript` und `p7-sigTERM-skript`,
   erwartet die ersten beiden): der bekannte Flake. Das Gate einzeln wiederholt: zweimal grün; ein dritter grüner Einzellauf lag davor (mit der
   Anpassung der Kopie). Fünf Einzelläufe insgesamt, einer rot.
3. Lauf **grün**, `rc=0`, 1018 s (Wanduhr), `integration-volumes-gegenprobe` darin grün. Namen der Volumes vor und nach dem ganzen Lauf: 49 und 49, neu 0,
   verschwunden 0 (vorher: ein Lauf von `make gates` erhöhte die Zahl um mindestens eins, §1). Kein Container mit dem Präfix `pgr-` danach.
Nicht in `make gates` und einzeln gefahren: `make schema-validate` (Volumes), `make test-integration` und `make test-container` als Handläufe (oben).

**Nacharbeit nach `6978e8f6`:** (a) Die Mutanten `p2-ohne-v`, `p3-leer-fehler` und `p4-code-verworfen` laufen jetzt über alle Fälle (`nur` leer), sodass
mehr rote Fälle als die Spalte „rote Fälle“ von `SPEC-057` Punkt 7 die Gegenprobe rot machen; gesehen: `rm -f` → `p2-grundfall`, `p2-ohne-v`; Exit-Code verworfen →
`p3-angelegt`, `p4-rest-rot`; fehlender Container als Fehler → `p3-kein-container`, `p5-sigint`, `p5-sigterm`. (b) **Abweichung, nicht entschieden:** „Namen erst nach dem Entfernen lesen“
färbt über alle Fälle `p3-angelegt`, `p4-rest-rot`, `p5-sigint` und `p5-sigterm` (die Fälle `p5-*` prüfen, dass die Meldung den Namen nennt); die Spalte nennt nur die ersten zwei. Der Mutant
läuft im committeten Stand mit `nur` = `p3 p4`. (c) **Mutant `p5-int-entfaellt` (Zeile `trap 'trap "" INT; exit 130' INT # S-INT` durch `:`) färbt keinen Fall:**
`p5-sigint` bleibt grün, weil der Runner im Fall in seinem `wait` steht (`IVG_NACH_RUN=3`) und bash dort ohne Falle ebenfalls aufräumt (Exit 130). Die Zeile ist nicht eingebaut;
Satz und Zeile in `SPEC-057` Punkt 5 und 7 entfallen oder der Architect nennt einen Fall, der die Falle fängt (der Verlust des Signals trat im Vordergrundkommando auf, `docker exec`).

**Nacharbeit nach `6fb43cd4`:** `p3-nach-entfernen` läuft über alle Fälle und verlangt genau `p3-angelegt`, `p4-rest-rot`, `p5-sigint`, `p5-sigterm` (gesehen). `p5-int-entfaellt`
(`trap 'trap "" INT; exit 130' INT # S-INT` → `:`) steht wie `p5-term-entfaellt` ohne roten Fall in der Gegenprobe; Lauf grün, das Verhalten stimmt mit `SPEC-057` Punkt 7 überein.
Die Absätze (b) und (c) der Nacharbeit nach `6978e8f6` sind damit entschieden. Die Tabelle oben führt die Mutation `p5-int-entfaellt` als äquivalent, kein Fall rot.

**Offene Punkte:** keine. Die Fragen zu `p5-sigterm`, zur Falle für `INT`, zu den roten Fällen je Gruppe und zum Inhalt des Volumes sind in `SPEC-057` Punkte 4, 5 und 7 und in §6 entschieden; keine neue Sendung an einen Nehmer in diesem Lauf.

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
Harness-Dokumentation nicht ausnimmt; die Verträge sind vor dem Anspruch geschrieben; Änderungen daran im Lauf (`SPEC-057` Punkte 5 und 7, `SPEC-055` Punkt 9) stehen im Diff und in der Metrik); kein Produkt-Code, keine Nutzer-Doku.

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
