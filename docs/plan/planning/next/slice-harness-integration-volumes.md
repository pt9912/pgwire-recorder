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

**Berührte Spec-Stellen:** `SPEC-038` (neuer Absatz *Aufräumen der Integrationstests*, vom
Architect vor dem Code anzulegen); `SPEC-052` (Punkt 4 und 13: Eintrag der neuen Gegenprobe in
die Ausnahmeliste der schnellen Kette und der Beleg, dass sie `docs/user/` und `README.md` nicht
liest; Architect, vor dem Code); `SPEC-048` (nur Bindung: Form der Abdeckungs-Deklaration im Kopf der Gegenprobe, §1; keine Änderung)

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
weglässt.

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
- Eine Diagnose der vollen Platte im Runner („PostgreSQL nicht bereit“ nennt die Ursache
  nicht, `docker logs` des PostgreSQL-Containers erscheinen nur nach einem roten Testlauf) —
  anderer Vorgang: Es wäre ein neuer Vertrag am Runner mit eigenen Randformen; er braucht
  einen eigenen Slice, wenn der Nutzer ihn will, und ist nicht angelegt.
- Produkt-Code, Tests unter `test/integration`, das `Dockerfile` und das Produkt-Image —
  Schicht-Abgrenzung: Der Slice ändert nur Skripte unter `tools/test`, Make-Fragmente unter
  `harness/mk` (`integration.mk`, `gates-schnell.mk`) und die Dokumentation und Verträge des
  Gates (`harness/README.md`, `SPEC-038`, `SPEC-052`).

**Sendungen an Dritte** (`AGENTS.md` §3.13; mit der Kennung dieses Slice als Geber):

- *An `harness/mk/gates-schnell.mk` und `SPEC-052`:* `integration-volumes-gegenprobe` endet auf
  `-gegenprobe` und hängt an `GATE_CHECKS`; `SPEC-052` Punkt 4 verlangt dafür einen Eintrag in
  `GATE_SCHNELL_AUSNAHME`, sonst ist `gates-schnell-gegenprobe` rot. Punkt 13 verlangt für jedes
  Gate der Ausnahmeliste auf `-gegenprobe` den Beleg, dass es `docs/user/` und `README.md` des
  Arbeitsbaums nicht liest; die Gegenprobe fährt in Kopien, die sie selbst anlegt (§3), und der
  Beleg steht in `SPEC-052` Punkt 13 (Architect, vor dem Code). Der Eintrag und die Zeile stehen
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
- [ ] Prüfung im Runner: Vor dem Entfernen liest der Runner die Namen der Volumes seiner
      Container (außer `$daten`), nach dem Entfernen prüft er, dass keines mehr existiert;
      sonst nennt die Meldung die Namen, der Lauf endet mit Exit-Code ungleich 0 (auch bei
      grünen Tests) und entfernt die Reste per Namen (§6). Die Funktionen liegen in einer
      Datei, die Runner und Gegenprobe laden.
- [ ] Gegenprobe: `make integration-volumes-gegenprobe` hängt an `GATE_CHECKS` und fährt
      in Kopien des Arbeitsbaums mit realem Docker die Fälle aus §6 *Mutationen*: Grundfall
      ohne Rest, Mutant ohne `-v` rot mit dem Namen des Volumes, Abbruch durch SIGTERM, Container
      nur angelegt, fremde Volumes bleiben; je Fall die Zahl der Volumes der Probe vor und nach
      (nach Namen). Die Zeile steht in `harness/README.md` §Sensors, der Vertrag in `SPEC-038`; das Ziel steht in
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
| `tools/test/run-integration-tests.sh` | update | `rm -fv` statt `rm -f` für `$tests` und `$pg`; ruft die Prüfung aus der geladenen Datei nach dem Entfernen; die Namen liest sie unmittelbar vor dem Entfernen (§6) |
| `tools/test/integration-aufraeumen.sh` | neu | Funktionen zum Lesen der Namen, Entfernen und Prüfen; Runner und Gegenprobe laden sie, damit die Gegenprobe den Code des Runners prüft, nicht eine Kopie |
| `tools/test/integration-volumes-gegenprobe.sh` | neu | Fälle aus §6 *Mutationen* gegen reales Docker, Mutant per Ersetzung in einer Kopie der Funktionsdatei (das Skript bricht ab, wenn die Ersetzung nichts oder mehr als eine Stelle trifft); Zahl der Volumes der Probe nach Namen |
| `harness/mk/integration.mk` | update | Ziel `integration-volumes-gegenprobe`, an `GATE_CHECKS`, hinter `test-integration` |
| `harness/mk/gates-schnell.mk` | update | `integration-volumes-gegenprobe` in `GATE_SCHNELL_AUSNAHME` (`SPEC-052` Punkt 4; sonst ist `gates-schnell-gegenprobe` rot), im Commit des Ziels |
| `harness/vertraege/SPEC-052-gate-kette.md` | update (Architect) | Punkt 13: Beleg, dass die Gegenprobe `docs/user/` und `README.md` nicht liest, und Punkt 4, soweit die Aufzählung die Gegenprobe nennt; vor dem ersten Code-Commit |
| `harness/README.md` | update | Zeile in §Sensors für das Ziel; die Zeile von `make test-integration` nennt die Prüfung; die Zeile von `gates-schnell` nennt die Ausnahmeliste unverändert (kein Gate-Name dort) |
| `spec/spezifikation.md` (`SPEC-038`) | update (Architect) | Absatz *Aufräumen der Integrationstests* mit den Entscheidungen aus §6, vor dem ersten Code-Commit |

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
Code-Commit entscheidet der Architect die Randformen aus §6 in `SPEC-038`
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Größe** (`AGENTS.md` §3.14): Metrik ist die Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools/test harness/mk/integration.mk
harness/mk/gates-schnell.mk harness/README.md harness/vertraege/SPEC-052-gate-kette.md
spec/spezifikation.md`, Tests eingeschlossen. Basis: der Commit des Anspruchs (`next` →
`in-progress`), beim Anspruch einzutragen. Schätzung (Arbeitswert des Planners aus den
Bestandsmaßen, vom Architect vor dem Code zu bestätigen; die Gegenprobe zählt mit dem Doppelten
des Produktcodes): 450 bis 750, Mitte 600. **Meldeschwelle 750:** Wer im Lauf darüber liegt oder
liegen wird, meldet es dem Planner, bevor er weiterbaut. **Marke 900:** darüber gilt die
Rückführung unten, ohne Ermessen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe der Metrik liegt über 900,
  oder die Gegenprobe braucht mehr als die
  fünf Fälle aus §6 *Mutationen* oder eine eigene Laufzeitumgebung und ist nicht in einer
  Review-Sitzung prüfbar; Schnitt dann: Aufräumen und Prüfung im Runner hier, die Gegenprobe
  mit dem Ziel in einem eigenen Slice direkt dahinter.
- `in-progress` → `open` (blockiert — Carveout?): `docker inspect` nennt die Volumes eines
  Containers nicht verlässlich (im Zustand `created`, auf einem Docker, den der Nutzer
  einsetzt); dann entscheidet der Architect einen anderen Weg zu den Namen (etwa Volumes
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
Runner, die Gegenprobe als Gate). **Zu entscheiden vom Architect vor dem ersten Code-Commit,
Ort: `SPEC-038`, Absatz *Aufräumen der Integrationstests*;** die Empfehlung des Planners steht
je Punkt dabei und gilt, bis der Architect oder der Nutzer anderes sagt (der Nutzer hat verfügt,
ohne Nachfrage weiterzumachen). Eine nicht genannte oder offene Randform entscheidet der
Implementer nicht, er gibt sie zurück.

- *Abbruch mitten im Lauf* — Empfehlung: SIGINT und SIGTERM beenden den Runner über den
  `trap … EXIT`, das Aufräumen läuft, die Prüfung auch; SIGKILL und ein Absturz des Hosts
  räumen nichts auf (akzeptiertes Negativ, kein Gegenmittel im Skript: ein Neustart des
  Runners findet die Container des alten Laufs nicht, weil ihr Name die PID trägt). Die
  Gegenprobe prüft den Fall SIGTERM mit Containern im Zustand `created`.
- *Parallele Läufe* (zwei `make test-integration` auf demselben Host, auch anderer Projekte) —
  Empfehlung: Die Prüfung liest und prüft nur die Namen aus `docker inspect` der Container
  des eigenen Laufs (Präfix `pgr-it-<PID>`); nie die Zahl der Volumes auf dem Host, nie
  `docker volume ls --filter dangling=true`, nie ein Name, den dieser Lauf nicht gelesen hat.
  Die Gegenprobe zählt deshalb je Fall nur die Namen ihrer eigenen Container.
- *Fehlgeschlagener Start* — Empfehlung: Scheitert `docker run -d` nach dem Anlegen des
  Containers (Netz fehlt, Port, Image-Start), steht der Container im Zustand `created` und das
  Volume schon in `Mounts`; das Aufräumen liest die Namen für jeden Zustand unmittelbar vor
  `rm -fv`. Existiert der Container nicht (Image nicht ladbar), gibt es keine Namen, und die
  Prüfung ist ohne Fehler leer. Der Runner endet dann wie heute mit dem Fehler des Starts.
- *Welche Volumes zählen* — Empfehlung: alle Mounts vom Typ `volume` der Container `$pg` und
  `$tests` außer `$daten`; `$daten` entfernt weiter `volume rm` unter seinem Namen. Ob ein
  Volume „anonym“ heißt, entscheidet nicht das Skript, sondern die Herkunft aus dem Container.
- *Rest nach dem Entfernen* — Empfehlung: Die Prüfung nennt die Namen, entfernt sie per Namen
  (nur diese) und setzt den Exit-Code auf 1, sofern der Lauf nicht schon mit einem Code
  ungleich 0 endet; der Code des Testlaufs behält Vorrang. Die Meldung nennt keinen Inhalt des
  Volumes (ein PostgreSQL-Datenverzeichnis kann Zugangsdaten tragen).
- *Docker-Versionen* — Empfehlung: `docker inspect` mit `Mounts` im Zustand `created` ist mit
  Docker 29.9.0 gemessen (2026-10-12); ältere Versionen sind nicht gemessen (akzeptiertes
  Negativ, `DOCKER` bleibt der einzige Parameter wie heute).
- *Image mit ungeprüftem `VOLUME`* — Empfehlung: Ein künftiges Image im Runner, das ein
  `VOLUME` deklariert, erfasst die Prüfung ohne Änderung, solange es als Container des Laufs
  startet; ein neues `docker run` außerhalb des Runners bleibt unberührt (kein Sensor dafür,
  siehe Risiko unten).
- *Leere Namensliste* — Empfehlung: kein Fehler, keine Meldung.
- *Ausnahmeliste und Beleg der schnellen Kette* (Sendung aus §1, `SPEC-052` Punkt 4 und 13) —
  Empfehlung: `integration-volumes-gegenprobe` steht in `GATE_SCHNELL_AUSNAHME`; der Beleg in
  Punkt 13 sagt, dass die Gegenprobe nur Skripte unter `tools/test` und `harness/mk` in ihre
  Kopien legt und `docs/user/`, `README.md` und `docs/plan/planning/` weder kopiert noch liest; ein
  Fall der Gegenprobe der Kette (`gates-schnell-gegenprobe`) hält die Ausnahmeliste.

**Mutationen** (`AGENTS.md` §3.10, für das Gate): je Zusage der Fall der Gegenprobe, vom
Implementer selbst rot gesehen, im Bericht als *Zusage · Mutation · roter Fall*.

| Zusage | Mutation (nur sie) | roter Fall |
|---|---|---|
| Nach einem Lauf bleibt kein Volume der eigenen Container | `rm -f` statt `rm -fv` in der Kopie der Funktionsdatei | Grundfall der Gegenprobe, Meldung nennt den Namen |
| Die Prüfung liest die Namen vor dem Entfernen, auch im Zustand `created` | Namen werden erst nach dem Entfernen gelesen | Fall *Container nur angelegt* |
| Das Aufräumen läuft beim Abbruch | `trap` entfällt in der Kopie des Runners | Fall *SIGTERM* |
| Das Werkzeug berührt kein fremdes Volume | die Prüfung entfernt alle Volumes ohne Container | Fall *fremde Volumes* |
| Rest macht den Lauf rot, auch bei grünen Tests | der Exit-Code der Prüfung wird verworfen | Grundfall des Mutanten ohne `-v` (Exit-Code) |

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
und den Absätzen in `SPEC-038` und `SPEC-052`, gezählt als eigene Schicht, weil §3.13 die
Harness-Dokumentation nicht ausnimmt); kein Produkt-Code, keine Nutzer-Doku.

**Nachgezählt** mit den Sendungen aus §1 (`gates-schnell.mk` und `SPEC-052`; Geber dieser Slice):
Liefer-Punkte weiter 3 (der Eintrag in die Ausnahmeliste und der Beleg in `SPEC-052` wachsen in
Punkt 3, Gegenprobe, hinein; sie sind keine eigenen Punkte, erhöhen aber den Umfang, daher die
Metrik in §4), Schichten weiter 2: beide Dateien liegen in den zwei genannten Schichten
(`harness/mk` bzw. `harness/vertraege`). Die Bindung an `slice-harness-abdeckung-gate` (§1) ändert
den Umfang dieses Slice nicht, sie bindet den Kopf der Gegenprobe.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
