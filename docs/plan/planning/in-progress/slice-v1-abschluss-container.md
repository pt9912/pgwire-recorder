# Slice slice-v1-abschluss-container: Container-Image und reproduzierbarer Binary-Build

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus)

**Berührte Spec-Stellen:** `LH-FA-16.a` · `SPEC-031` · `SPEC-035` · [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) (Vertrag dieses Slice, Architect 2026-10-17)

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

**Ziel:** Ein Docker/OCI-Image (`linux/amd64`, `linux/arm64`) und ein reproduzierbarer Binary-Build existieren, belegt durch das Gate `make test-container`. Handbuch, README und Betriebsdokumentation dazu liefert `slice-v1-abschluss-container-doku` (Schnitt vom 2026-10-17, §4).

**Zuschnitt nach dem Vorlauf des Architect (2026-10-17, [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md), [ADR-0043](../../adr/0043-binaries-und-container-gate.md) `Accepted`, Entscheidung des Nutzers vom 2026-10-17):** Das Werkzeug `make binaries` baut die sechs Binaries (`linux`, `darwin`, `windows`, je `amd64`, `arm64`) mit `SHA256SUMS`; das Gate `make test-container` baut sie zweimal ohne Cache der Bau-Stufen und vergleicht die Summen, prüft den Kopf jedes Binaries, baut das Image je Linux-Plattform (mit demselben Binary) und fährt auf der Plattform des Docker-Servers den Smoke im Container (`version`, `record` gegen PostgreSQL, `replay` ohne PostgreSQL, `SELECT 1`). Neu sind das Werkzeug, das Gate, der Smoke und die Stufe `dist`; vorhanden sind das Multistage-`Dockerfile` mit `-trimpath`, `-buildvcs=false` und gepinnten Images, `make build` und die Stufe `runtime` (ein Probe-Build von `windows/arm64` und `darwin/amd64` baute, `windows/arm64` zweimal ohne Cache der Stufe `build` mit derselben Summe).

**Übernommen aus `slice-v1-abschluss-herunterfahren` (aus `slice-v1-abschluss-betrieb` hervorgegangen) und `slice-v1-abschluss-homebrew`** (dort je §1, Abgrenzung): das Container-Image (Docker/OCI); das ist das Ziel oben.

**Abgegeben** an `slice-v1-abschluss-container-doku` (dort §1, *Übernommen aus*, mit der Kennung dieses Slice; Schnitt vom 2026-10-17 beim Erreichen der Meldeschwelle, §4): Handbuch §1 *Voraussetzungen* (Abnahme-Umgebungen, `LH-QA-03`) und §2 *Installation* (Image und seine Plattformen, `make binaries`, `SHA256SUMS`), `README.md` und `docs/maintainer/releasing.md`. Der Nehmer ruft, was dieser Slice liefert (`make binaries`, das Image der Stufe `runtime`), und prüft seine Beispiele daran.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Handbuch, `README.md` und `docs/maintainer/releasing.md` — `slice-v1-abschluss-container-doku`: Handbuch und README beschreiben den Ist-Zustand des gebauten Binaries und des Images (`AGENTS.md` §3.11), und der Code-Teil allein überschritt mit dem Doku-Teil die Marke von §4; der Doku-Teil ist eine zweite Schicht (§3.13).
- Der Veröffentlichungs-Mechanismus der Registries (`ghcr.io`, `docker.io`) — Gegenstand des Release-Verfahrens in `docs/maintainer/releasing.md`; dieser Slice liefert das Image.
- Die Manifestliste für beide Linux-Plattformen in einer Referenz und die Anhänge am GitHub-Release — `slice-erster-release-veroeffentlichung`: Der Standard-Treiber von Docker baut keine Manifestliste, und sie gehört zum Push (`SPEC-055` Punkt 13). Jener Slice ruft `make binaries` und hängt `dist/` mit `SHA256SUMS` an, statt ein zweites Mal zu bauen (Hinweis des Architect; Voraussetzung in dessen §1 eingetragen, 2026-10-17).
- Der Smoke des Images auf der jeweils anderen Plattform (Emulator) und ein Lauf der Binaries für `darwin` und `windows` — kein Emulator und keine Gegenstelle in der Prüfumgebung; belegt sind Bau, Kopf und Summe (`SPEC-055` Punkt 13, Grenze 3 und 4).
- Reproduzierbarkeit des Images — `LH-QA-01` verlangt sie nur für das Binary (`SPEC-055` Punkt 13).
- Gate für Code-Tabelle und Katalog — liefert `slice-harness-meldungskatalog-gate` (wellenlos, Entscheidung des Nutzers vom 2026-10-07, `done/`); dieser Slice berührt das Handbuch nicht, die Tabellen `### Fehlercodes` und `### Warnungen` hält `make meldungskatalog-check` im Nehmer `slice-v1-abschluss-container-doku`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus) (Container-Bau mit Gate): `make binaries` schreibt die Binaries aller Zielplattformen aus `SPEC-035` mit `SHA256SUMS`; `make test-container` baut sie zweimal mit derselben Prüfsumme (Reproduzierbarkeit verlangt `LH-QA-01` nur für das Binary), prüft den Kopf je Binary, baut das Image je Linux-Plattform mit demselben Binary und besteht auf der Plattform des Docker-Servers den Walking-Skeleton-Smoke im Container (`replay` ohne PostgreSQL); die Gegenprobe `make test-container-gegenprobe` fängt je Zusage eine Mutation (`AGENTS.md` §3.10).
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
| `Dockerfile` | update | Stufe `dist` (nur die Binaries, `FROM scratch`); die Bau-Stufe bleibt |
| `.gitignore`, `harness/mk/container.mk`, `tools/harness/binaries.sh`, `tools/harness/run-container.sh`, `tools/harness/container-eingabe.sh`, `tools/harness/container-gegenprobe.sh` | neu | `binaries`, `test-container`, `test-container-gegenprobe` nach `SPEC-055`; `container-eingabe.sh` trägt die Prüfung der Eingabe (Punkt 1, 2) für beide Skripte; `.gitignore` nennt `dist/`, damit `make binaries` den Hash des Arbeitsbaums nicht ändert |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh`, `harness/vertraege/SPEC-052-gate-kette.md`, `harness/README.md` | update | `test-container` und seine Gegenprobe in die Ausnahmeliste (SPEC-052 Punkt 2 und 4); Zeilen in §Sensors |

Handbuch, README und `docs/maintainer/releasing.md` stehen nicht hier, sondern in `slice-v1-abschluss-container-doku` (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die anderen Slices der Welle sind `done`, `slice-harness-meldungskatalog-gate` liegt in `done/` (§1), und [ADR-0043](../../adr/0043-binaries-und-container-gate.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-17; erfüllt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- Dockerfile harness tools test` liegt über 1200; Basis ist `d258c0b7` (HEAD vor dem Anspruch). Die Pfade `docs/user`, `docs/maintainer` und `README.md` gehören seit dem Schnitt nicht mehr dazu. **Durch den Schnitt behandelt (2026-10-17):** Die Messung erreichte beim Code 1147 (Meldeschwelle 1000 überschritten, Marke 1200 nicht), und mit dem Doku-Teil (Schätzung 180 bis 280) hätte die Schätzung (rund 1330 bis 1430) die Marke überschritten; der Schnitt nach Lieferwert (*Container-Bau mit Gate | Handbuch, README, `releasing.md`*, Zuschnitt: Planner, Nutzervorgabe: nur schneiden, wenn die Schätzung die Marke überschreitet) liegt zwischen den zwei Liefer-Punkten, wie der Architect ihn vorgesehen hatte. Der Rest dieses Slice (Review-Befunde, Closure) ist klein; wächst er über 1200, geht er zurück. Schätzung des Architect (in dieser Metrik, `AGENTS.md` §3.14) für den Code: Produktcode 190, Gegenprobe 380 bis 570, Ausnahmeliste, `SPEC-052` und `harness/README.md` 30 bis 40; gemessen 1147.
- `in-progress` → `open`: Ein Binary einer Zielplattform lässt sich ohne Änderung am Produktcode nicht bauen oder zwei Builds liefern verschiedene Summen aus einer Ursache, die das `Dockerfile` nicht beheben kann — Carveout (der Probe-Build von `windows/arm64` und `darwin/amd64` baute, die Summe von `windows/arm64` war bei zwei Builds gleich).

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

**Randformen des neuen Vertrags** (`AGENTS.md` §3.12), vom Architect am 2026-10-17 vor dem ersten Code-Commit entschieden; die Stelle ist [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md), der Punkt in Klammern. Zu jeder Zusage steht eine Mutation (§3.10), vom Implementer selbst rot gesehen; der Bericht nennt je Zusage: Zusage · Mutation · roter Test.

- Fehlende, leere, unbekannte, doppelte oder zu kleine Liste der Zielplattformen und eine Image-Plattform außerhalb der Binaries (Punkt 1, Eingabe: kein `docker` aufgerufen); die Listen sind keine Eingabe des Aufrufers.
- Fehlende, leere oder unzulässige `VERSION` (Punkt 2); Wert in `-ldflags` und Dateiname. Review F-1, Architect 2026-10-17: Die Vorgabe `dev` maskiert ein fehlendes `-X main.version=`; `test-container` setzt `VERSION := 0.0.0-test` und lehnt `dev` vor `docker` ab, der Smoke am echten `Dockerfile` belegt den Wert; die Gegenprobe prüft dazu die Vorgabe, die Ablehnung und die Zeile `-X main.version=$VERSION` im `Dockerfile`.
- `DIST` vorhanden und nicht leer, leerer Wert von `DIST`, Sortierung von `SHA256SUMS` in Byte-Reihenfolge, Teilmenge bei einem fehlgeschlagenen Build, Format von `SHA256SUMS`, Endung `.exe` (Punkt 3).
- Zwei Builds: abweichende Summe je Plattform (sechs Fälle), Umgehen des Caches der Stufen `build` und `dist` als Voraussetzung dafür, dass der zweite Build ein Build ist (Punkt 4). Akzeptiertes Negativ: Das Gate vergleicht auf demselben Docker-Server mit denselben Eingangs-Images; ein zweiter Host und der Pin der Images werden nicht geprüft (Punkt 13, Grenze 2), weil Anhebungen bewusste Commits sind ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)) und es keine zweite Umgebung zur Abnahme gibt.
- Zeitstempel und VCS-Stempel: kein Fall; sie gehen nicht in das Binary ein (`-buildvcs=false`, Go schreibt keinen Zeitstempel), und ein Build mit eingeschleustem Zeitstempel wird von Punkt 4 rot (der Implementer zeigt es einmal von Hand und nennt es im Bericht, geschehen in §7; keine Gegenprobe mit echtem Docker, Punkt 12).
- Netz im Build: nur die Stufe `deps` hat Netz, `build` läuft mit `--network=none` ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)); der Bau mit kaltem Cache braucht einmal Netz, wie `make build` (kein Fall).
- Kopf der Binaries je Format und Architektur, anderes Format (Punkt 5).
- Image: falsche Plattform in den Metadaten, abweichende Summe der Datei im Image, Bau ohne Emulation (Punkt 6).
- Smoke: Plattform des Docker-Servers außerhalb von `amd64` und `arm64`; `version`; Aufzeichnen und Wiedergeben, leere und ungleiche Sicht, Exit-Code 0 nach `SIGTERM`, PostgreSQL vor dem Wiedergeben entfernt, kein `--user`, Volume oder Bind-Mount (Punkt 7); Schreibrecht des Benutzers `nonroot` auf `/tmp` ist damit belegt.
- Reihenfolge der Gründe bei mehreren roten Punkten (Punkt 8).
- Abbruch und Herunterfahren: Aufräumen nach rotem Lauf, `SIGINT`, `SIGTERM`, auch des Stage-Verzeichnisses von `make binaries`; `DIST` bleibt bei einem Fehler leer (Punkte 3 und 9; Review F-7). Akzeptiertes Negativ: `SIGKILL` und Ausfall des Docker-Dienstes (Punkt 13, Grenze 6), wie bei `SPEC-053`.
- Aufrufarten von `make`: `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung, auch in Kombination, je Ziel `binaries` und `test-container` (Punkt 11); Aufruf-Voraussetzung: nacheinander, `make -j gates` nicht zugesagt (Punkt 13, Grenze 5).
- Ausnahmeliste der schnellen Kette und Lesen des Arbeitsbaums (Punkte 10 und 12).
- Verhalten je Serverversion: Der Smoke läuft gegen die Referenz-PostgreSQL-Version (`POSTGRES_IMAGE`); die Matrix bleibt bei `SPEC-053`, kein Fall.
- Die Manifestliste, die Registries und der Smoke des Images der anderen Plattform liegen nicht in diesem Slice (§1; `slice-erster-release-veroeffentlichung`, `SPEC-055` Punkt 13).
- Die Sätze des Handbuchs zu Image, Plattformen und `SHA256SUMS` liegen in `slice-v1-abschluss-container-doku`; sie sagen nur zu, was `make test-container` belegt (`AGENTS.md` §3.11): kein Lauf von `darwin` und `windows`, kein Smoke der anderen Linux-Plattform (§1, Grenze 3 und 4 in `SPEC-055` Punkt 13). Das Gate dieses Slice ist die Quelle (Bindung an den Nehmer, §3.13).

**Risiken:**

- Handbuch und Releasing-Dokument beschreiben Verhalten vor der Software; Abweichungen zeigen sich erst im Bau — **Ausgang:** entfallen: Handbuch und `releasing.md` gehören seit dem Schnitt vom 2026-10-17 nicht mehr zu diesem Slice; das Risiko steht in §6 von `slice-v1-abschluss-container-doku`, dessen Beispiele gegen das gebaute Binary und das Image laufen.
- Die [ADR-0043](../../adr/0043-binaries-und-container-gate.md) war `Proposed` — **Ausgang:** entfallen; der Nutzer hat sie am 2026-10-17 angenommen (`Accepted`).
- Die Gate-Laufzeit wächst um die Zahl der Builds (zwölf der Stufe `build` ohne Cache, zwei Images, ein Smoke); der Lauf des Gates in `make gates` ist die Probe; der Bericht nennt die Dauer — **Ausgang:** *weiter offen* → Register `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (verkörpert, +1, `evidence/slice-v1-abschluss-container.md`): Eingetreten ist ein Zuwachs von rund 5 Minuten (Gate etwa 1 Minute, Gegenprobe 3 m 40 s bis 3 m 45 s, `make gates` am Endstand 927 s); er liegt innerhalb dessen, was die Klasse trägt, kein Carveout und keine Lockerung (`AGENTS.md` §3.6), und die Gegenprobe steht in der Ausnahmeliste der schnellen Kette.
- Die Abdeckungstabellen (`make abdeckung-check`, `SPEC-048`, `[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)`) könnten für `LH-FA-16`, `LH-QA-01` oder `LH-QA-03` einen Beleg des Gates verlangen; der Implementer prüft es mit dem ersten Lauf und deklariert dann nach `[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)` im Kopf der Gegenprobe, ohne den Vertrag von `SPEC-048` zu ändern — **Ausgang:** *entfallen*: `make abdeckung-check` ist am Endstand grün, ohne dass eine Tabelle `LH-FA-16`, `LH-QA-01` oder `LH-QA-03` führt, und verlangt keinen Beleg des Gates ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) erlaubt die Deklaration, verlangt sie nicht); es wurde keine eingetragen, `SPEC-048` blieb unverändert.

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

Stand: Endstand `e365401f` (Code `91a1bd02`, `7e5ecb79`; Nacharbeit `39c07812`, `62a287cb`, `e365401f`), Basis der Messung `d258c0b7`. Docker Server 29.9.0, `linux/amd64`. Die Läufe unten nennen, wo sie vom Endstand abweichen.

**Läufe (jeweils einzeln, nacheinander)**

| Lauf | Ergebnis |
|---|---|
| `make binaries` (leeres `dist/`) | sechs Binaries und `SHA256SUMS`, `sha256sum -c` grün, 29 s |
| `make test-container` (echtes Docker) | grün, 63 s: 12 Builds, Summen gleich, Kopf von 6 Binaries, Image `linux/amd64` und `linux/arm64`, Smoke auf `linux/amd64`; danach keine Container, kein Netz, keine Tags, keine Temp-Verzeichnisse |
| `make test-container-gegenprobe` (über `make`, Endstand `e365401f`) | grün, 3 m 40 s |
| `make gates-schnell` (Endstand) | grün; Verifier 10 s |
| `make abdeckung-check`, `make meldungskatalog-check`, `make kopf-check` (Endstand) | je Exit 0 |
| `make gates` (Endstand `e365401f`) | grün, 927 s (15 m 27 s) im dritten Anlauf; im zweiten trat der Flake `p7-int-treiber` auf (`BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`, nicht dieser Slice); der Nachweis `.harness/state/gates-passed.diffsha` passt zum Arbeitsbaum, die Verifikation hat ihn geprüft |
| Verifier am Endstand | `make test-container` grün in 1 m 08 s, `make binaries` 29 s, Stichprobe der Mutationen aus allen Gruppen wie behauptet rot |
| Zeile davor (Stand `7e5ecb79`, vor der Nacharbeit) | `make test-container` 63 s, Gegenprobe 3 m 45 s, `make gates` 937 s |

**Risiko 3 (Laufzeit):** `test-container` kostet etwa 1 Minute, die Gegenprobe 3 m 40 s bis 3 m 45 s; der volle Lauf dauerte 927 s (Endstand) bzw. 937 s (davor).

**Risiko 4 (Abdeckung):** `make abdeckung-check` ist grün, ohne dass eine Tabelle `LH-FA-16`, `LH-QA-01` oder `LH-QA-03` führt, und verlangt keinen Beleg des Gates. Befund: keine Deklaration eingetragen ([ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) erlaubt sie, verlangt sie nicht).

**Messung (Metrik §4, `git diff --numstat d258c0b7..HEAD -- Dockerfile harness tools test`, Stand nach der Nacharbeit zum Review):** 1191 gesamt (1176 Einfügungen, 15 Löschungen). Marke 1200 nicht erreicht, Meldeschwelle 1000 überschritten (gemeldet).

**Handlauf Zeitstempel (SPEC-055 Punkt 4, §6):** Kopie des Baums (`git ls-files` nach `…/scratchpad/zeit`, außerhalb des Arbeitsbaums), in der Stufe `build` des `Dockerfile` der Kopie `-X main.version=$VERSION` durch `-X main.version=${VERSION}-$(date +%s%N)` ersetzt (`sed -i`), dort `make test-container`. Ausgang, Exit 2 nach 1 m 47 s: `test-container: ROT — die Summen zweier Builds weichen ab: pgwire-recorder_dev_darwin_amd64 pgwire-recorder_dev_darwin_arm64 pgwire-recorder_dev_linux_amd64 pgwire-recorder_dev_linux_arm64 pgwire-recorder_dev_windows_amd64.exe pgwire-recorder_dev_windows_arm64.exe`. Rot bei Punkt 4 mit den Namen aller sechs Binaries, wie erwartet; kein Image gebaut. Kopie entfernt; danach keine Container, Netze, Tags oder Temp-Verzeichnisse des Gates.

**Handlauf `-X main.version` (Review F-1, SPEC-055 Punkt 2):** Kopie des Baums (`git ls-files` nach `…/scratchpad/xver`, außerhalb des Arbeitsbaums), in der Stufe `build` des `Dockerfile` der Kopie ` -X main.version=$VERSION` per `sed -i` entfernt (`-ldflags="-s -w"`), dort `make test-container` (Vorgabe `VERSION=0.0.0-test`). Ausgang nach 1 m 56 s: Punkte 4 bis 6 grün, `test-container: ROT — version gibt 'pgwire-recorder dev' aus statt 'pgwire-recorder 0.0.0-test' und LF` (Schritt 7a, Smoke `version`), Exit 1, `make` 2. Kopie entfernt; keine Container, Netze, Tags oder Temp-Verzeichnisse des Gates zurück.

**Wege der Mutationen:** `GP_MUT='<Muster>' GP_ENTDECKE=1 bash tools/harness/container-gegenprobe.sh` nennt je Mutation die roten Fälle; `mutiere` ersetzt in einer Kopie des Baums (Textersatz, kein `cp -p`), jede Mutation hat eine eigene Kopie. Docker-teure Mutationen gibt es nicht: Die Gegenprobe läuft mit der Attrappe von `docker`; der Lauf des echten `Dockerfile` ist `make test-container`. Das Einschleusen eines Zeitstempels in das echte `Dockerfile` (SPEC-055 Punkt 4, §6) ist als Handlauf oben gefahren, kein Gate wiederholt es. Entscheidung des Architect 2026-10-17: Der Implementer zeigt es einmal von Hand in einer Kopie des Baums außerhalb des Arbeitsbaums (in der Stufe `build` des `Dockerfile` `-X main.version=${VERSION}-$(date +%s%N)` einsetzen, dort `make test-container` fahren; erwartet rot bei Punkt 4 mit den Namen der Binaries) und hält Befehl und Ausgang hier fest; das ist der Handlauf Zeitstempel oben. Die sechs Fälle der Attrappe belegen nur den Vergleich, nicht das echte `Dockerfile` (§3.11); ohne Handlauf bliebe die Aussage „ein eingeschleuster Zeitstempel wird rot“ ungeprüft und dürfte nirgends stehen.

**Mutationstabelle** (Zusage · Mutation · roter Fall; „nur“ heißt: die Mutation färbt in ihrer Gruppe genau diese Fälle; Datei `tools/harness/…`)

Punkt 1 (Plattformen, vor dem ersten `docker`):

| Zusage | Mutation | roter Fall |
|---|---|---|
| leere Liste | Zeile `P1a` entfernt (`container-eingabe.sh`) | `p1a-bin-leer`, `p1a-img-leer` |
| Eintrag ohne `/` | `P1b1` | `p1b1-ohne-slash` |
| leeres Betriebssystem | `P1b2` | `p1b2-os-leer` |
| leere Architektur | `P1b3` | `p1b3-arch-leer` |
| Betriebssystem unbekannt | `P1c1` | `p1c1-os-unbekannt` |
| Architektur unbekannt | `P1c2` | `p1c2-arch-unbekannt` |
| Eintrag doppelt | `P1d` | `p1d-doppelt` |
| Plattform fehlt, je eine der sechs und der zwei | `ERF_BINARIES`/`ERF_IMAGES` ohne diese Plattform | `p1e-bin-<plattform>` (sechs, bei `darwin/arm64` auch `p1e-binaries`), `p1e-img-<plattform>` (zwei) |
| Image-Plattform nicht in den Binaries | `P1f` | `p1f-image-nicht-in-binaries` |
| Eingabefehler endet vor `docker` (Gate) | `P8vor-docker` (`run-container.sh`) | alle Fälle von Punkt 1 und 2, dazu `p2-dev-gate`, `p8-1-vor-5` |
| dasselbe bei `make binaries` | `exit 2` in `binaries.sh` entfernt | `p1e-binaries`, `p2-binaries`, `p3-dist-nicht-leer`, `p3-zweiter-lauf` |

Punkt 2 (`VERSION`):

| Zusage | Mutation | roter Fall |
|---|---|---|
| leer unzulässig | `-z`-Zweig vor dem Muster | `p2-leer`, `p2-binaries` |
| Leerzeichen / `$` / `'` / `"` unzulässig | Zeichen in `rest` aufgenommen (je eins) | `p2-leerzeichen`, `p2-dollar`, `p2-einfaches-anfuehrungszeichen`, `p2-doppeltes-anfuehrungszeichen` |
| führendes `-` unzulässig | `-` in `erst` | `p2-fuehrendes-minus` |
| `.`, `+`, `_`, `-` zulässig (je ein Zeichen) | Zeichen aus `rest` entfernt (vier Mutationen) | `p2-gueltig` |
| `VERSION` geht in den Namen/Build ein | feste Zeichenkette `dev` im Vergleich der Ausgabe | `p2-gueltig` |
| Vorgabe am Ziel `test-container` ist `0.0.0-test` | Zeile `test-container: VERSION := 0.0.0-test` entfernt | `p2-vorgabe-test` |
| `dev` am Gate ist rot vor `docker` | Prüfung in `run-container.sh` entfernt | `p2-dev-gate` |
| das `Dockerfile` setzt `-X main.version=$VERSION` | Text aus dem `Dockerfile` der Kopie entfernt | `p2-dockerfile-x` |
| `-X main.version=$VERSION` im echten `Dockerfile` | Zeile in der Kopie entfernt | `p2-dockerfile-x` |

Punkt 3 (`make binaries`):

| Zusage | Mutation | roter Fall |
|---|---|---|
| `DIST` nicht leer oder Datei: rot vor `docker` | `P3dist` | `p3-dist-nicht-leer`, `p3-dist-datei`, `p3-zweiter-lauf` |
| `DIST=` (leerer Wert): Eingabefehler mit Meldung, vor `docker` | `P3distleer` | `p3-dist-leerer-wert` |
| `DIST` fehlt: wird angelegt | `mkdir -p` entfernt | `p3-dist-fehlt` und die Fälle, die in ein fehlendes `out` schreiben |
| Endung `.exe` unter Windows | `P3exe` | `p3-namen`, `p3-exe`, `p3-sums-je-binary` |
| Name trägt die `VERSION` | Name ohne `VERSION` | `p3-exe`, `p3-namen`, `p3-sums-format`, `p3-sums-je-binary` |
| `VERSION` im Build | `--build-arg` entfernt | `p3-version-im-build` |
| sortiert | `P3sort` | `p3-sums-sortiert` |
| zwei Leerzeichen | `sha256sum --binary` | `p3-sums-format`, `p3-sums-je-binary` |
| Summe SHA-256 | `sha1sum` | `p3-sums-format`, `p3-sums-pruefbar` |
| Name ohne Pfad | Namen mit `./` | `p3-sums-format`, `p3-sums-je-binary` |
| LF | `\r` angehängt | `p3-sums-lf`, `p3-sums-format`, `p3-sums-je-binary` |
| je Binary eine Zeile (sechs) | `grep -v` einer Plattform | `p3-sums-format`, `p3-sums-je-binary` (je Plattform) |
| Build scheitert: rot mit Plattform, keine Teilmenge | `break` statt `exit 1`; Binary direkt nach `DIST`; `SHA256SUMS` aus Teilmenge | `p3-rot-nennt-windows-arm64`, `p3-rot-ohne-sums`, `p3-rot-keine-teilmenge` (je nach Mutation) |
| Werkzeug ohne Cache-Umgehung | Cache-Flags immer | `p3-nur-binaries` |

Punkt 4 (zwei Builds):

| Zusage | Mutation | roter Fall |
|---|---|---|
| abweichende Summe je Plattform (sechs) | Zeile der Plattform aus dem Vergleich genommen | `p4-abweichung-<plattform>` (bei `linux/amd64` auch `p8-4-vor-6`) |
| Namen der abweichenden Binaries | alle Namen statt der abweichenden | sechs `p4-abweichung-*`, `p4-namen` |
| zweiter Lauf umgeht den Cache | `P4zweiter`; `--no-cache-filter build` entfernt; `--no-cache-filter dist` entfernt | je `p4-cache-umgehen`, `p4-gruen` |
| Verzeichnisse außerhalb des Arbeitsbaums | `DIST` relativ im Baum | `p10-arbeitsbaum` |

Punkt 5 (Kopf; der Lauf endet beim ersten roten Kopf, Mutationen sind auf ihre Fälle begrenzt):

| Zusage | Mutation | roter Fall |
|---|---|---|
| ELF-Magic / Maschine | `P5elf-magic` / `P5elf-maschine` | `p5-format-linux-*` / `p5-arch-linux-*` |
| Mach-O-Magic / CPU-Typ | `P5macho-magic` / `P5macho-cpu` | `p5-format-darwin-*` / `p5-arch-darwin-*` |
| PE `MZ` / Maschine | `P5pe-magic` / `P5pe-maschine` | `p5-format-windows-*` / `p5-arch-windows-*` |
| Architektur je Plattform (sechs) | Wert der Plattform vertauscht | `p5-gruen`, `p5-arch-<plattform>` |
| PE: Offset aus Byte 0x3c | fester Offset 0x84; Byte 0x38 statt 0x3c | `p5-gruen`, `p5-arch-windows-*` |

Punkt 6 (Image):

| Zusage | Mutation | roter Fall |
|---|---|---|
| Betriebssystem / Architektur / Summe | `P6os` / `P6arch` / `P6summe` | `p6-os-*` / `p6-arch-*` / `p6-summe-*` |
| dasselbe je Plattform (zwei mal drei) | `|| [ "$p" = <plattform> ]` | `p6-<form>-<plattform>` |
| Bau scheitert: rot mit Plattform | `true` statt `rot` | `p6-bau-rot` |
| Stufe `runtime` ohne `RUN` | `RUN true` im `Dockerfile` | `p6-runtime-ohne-run` |

Punkt 7 (Smoke):

| Zusage | Mutation | roter Fall |
|---|---|---|
| Server-Plattform nur `amd64`/`arm64` | `P7server`; nur `amd64`; nur `arm64`; Tag fest `amd64` | `p7-server-fremd`, `p7-server-arm64`, `p7-gruen` |
| Netz `--internal` | ohne `--internal` | `p7-intern` |
| `version`: Exit 0 / genau die Ausgabe / LF / keine zweite Zeile | `|| true`; `P7a-ausgabe`; Vergleich über `$(…)`; nur erste Zeile | `p7a-exit`; `p7a-ausgabe-abweichend`, `p7a-ohne-lf`, `p7a-zusatzzeile`; `p7a-ohne-lf`; `p7a-zusatzzeile` |
| Sicht A nicht leer / nennt `eins` (zwei leere Sichten sind gleich) | Prüfung entfernt (je) | `p7b-sicht-a-leer`, `p7b-beide-leer`; `p7b-ohne-eins` |
| PostgreSQL bereit | `|| true` | `p7b-pg-nicht-bereit` |
| Aufzeichnung entnommen / nicht leer | `|| true`; `P7aufz` | `p7b-cp-rot`; `p7b-aufzeichnung-leer`; `p7-cp-entnehmen` |
| Exit 0 je Container | Ausnahme für `record` / `replay` | `p7b-exit-rec` / `p7c-exit-rep` |
| PostgreSQL entfernt vor dem Wiedergeben | `P7pg`; `rm -f` des PostgreSQL-Containers entfernt | `p7c-pg-lebt`; `p7-pg-vor-replay` |
| Aufzeichnung per `docker cp` in den Container | `cp` entfernt | `p7-cp-hinein` |
| Sicht B gleich Sicht A | `P7gleich` | `p7c-sicht-ungleich` |
| kein Terminal / `--user` / `-v` / `--mount` | Flag an der `record`-Zeile (je) | `p7-kein-terminal`, `p7-kein-user`, `p7-kein-volume`, `p7-kein-mount` |
| kein `-v` / `-u` als Kurzflag | `run -dv /tmp:/tmp`; `run -d -u0` an der `record`-Zeile | `p7-kein-volume`; `p7-kein-user` (je nur dieser) |

Punkt 8: `P8vor-docker` (oben); Gate fährt nach rotem Punkt 4 weiter (`rot` durch `echo` ersetzt) färbt `p8-4-vor-6`, alle `p4-abweichung-*`, `p4-namen`, `p4-cache-umgehen`.

Punkt 9 (Aufräumen; Enden: grün, rot 4 bis 7, `SIGINT`, `SIGTERM`, `SIGINT` an die Gruppe, `binaries.sh` allein mit `SIGINT`/`SIGTERM`):

| Zusage | Mutation | rote Fälle |
|---|---|---|
| Container / Netz / Tags / Temp-Verzeichnis entfernt | `P9container` / `P9netz` / `P9tags` / `P9werk` | `p9-<ende>-container` / `-netz` / `-tags` / `-werk` für die Enden, an denen die Ressource besteht |
| nach `SIGINT` / `SIGTERM` | Fang auf `kill -KILL`; EXIT-Fang entfernt | `p9-sigINT-*`, `p9-sigTERM-*` |
| Signal an den laufenden Aufruf weitergegeben | `kill -TERM "$laeuft"` entfernt (Gate, `binaries.sh`) | `p9-<ende>-signal` |
| Stage von `binaries.sh` entfernt | EXIT-, INT-, TERM-Fang in `binaries.sh` | `p9-<ende>-stage`, `p9-binINT-*`, `p9-binTERM-*` |
| Namen tragen die Kennung des Laufs | `pgr-ct-x`; `pgr-ct-1` | `p9-kennung-ziffern`, `p9-kennung-eindeutig` |

Punkt 11 (Aufrufarten): `$(MAKE)` im Rezept oder `+` je Ziel färbt `-n`, `-q`, `-t`, `MAKEFLAGS=n/q/t` und die Kombinationen `-n -i`, `-q -i`, `-t -i` (je Ziel neun Fälle); ein Rezept, das unter `-i` nichts fährt, färbt `p11-<ziel>-i` und `p11-<ziel>-makeflags-i`; `$(MAKE)` oder `+` im Rezept der Gegenprobe färbt `p11-gegenprobe-rezept`.

Punkt 12: Anschluss an `GATE_CHECKS` je Ziel entfernt färbt `p12-gate-checks-ziel`, `p12-gate-checks-gegenprobe`; die Ausnahmeliste der schnellen Kette färben die Mutanten in `gates-schnell-gegenprobe.sh` (`m-test-container`, `m-test-container-gegenprobe`).

**Äquivalente und Grenzen**

- `DIST` als reguläre Datei: Die Mutation von `P3dist` allein färbt `p3-dist-datei` nicht rot, wenn `exit 2` fehlt, weil `mkdir -p "$dist"` (`binaries.sh`, vor dem ersten Aufruf von `docker`) auf einer Datei scheitert. Invariante: `mkdir -p` auf eine Datei endet mit Exit-Code ungleich 0, bevor `docker` ruft. Mit der Mutation `%P3dist` wird der Fall rot (die Meldung fehlt).
- `DIST=` (leerer Wert): vom Architect am 2026-10-17 entschieden, `SPEC-055` Punkt 3: ungültige Eingabe mit Meldung vor dem ersten `docker` (Marke `P3distleer` in `binaries.sh`, Fall `p3-dist-leerer-wert`, Mutation `%P3distleer`); umgesetzt, Mutation `%P3distleer` färbt genau `p3-dist-leerer-wert` (`p3-dist-nicht-leer`, `p3-dist-datei`, `p3-zweiter-lauf` bleiben grün). Wird das `exit 2` entfernt (Mutation `p3-vor-docker`), bleibt der Fall grün, weil `mkdir -p ""` vor dem ersten `docker` scheitert; die Meldung bleibt.
- `LC_ALL=C sort` zu `sort` in `binaries.sh` (Review F-8): äquivalent. Invariante: Die sechs Namen unterscheiden sich nur in Betriebssystem und Architektur bei gleichem Präfix (`darwin_amd64` bis `windows_arm64.exe`); die Reihenfolge ist in jeder Locale gleich. Ort: Namensschema `pgwire-recorder_<VERSION>_<os>_<arch>` in `binaries.sh`. `P3sort` entfernt das ganze `sort` und ist nicht äquivalent.
- Akzeptiertes Negativ (Review F-10): `bereit` in `run-container.sh` liest `docker logs … | grep -q` unter `pipefail`; endet `grep -q` vor dem letzten Schreiben von `docker logs`, kann die Pipeline mit 141 als nicht gefunden gelten. Grund: Die Ausgabe ist eine Meldung, der Fall trat in keinem Lauf auf; es gibt keinen Fall. Flackert das Gate, ist das ein Befund.
- `LC_ALL=C` für die Sortierung von `SHA256SUMS`: `SPEC-055` Punkt 3 nennt seit dem 2026-10-17 die Byte-Reihenfolge (C-Locale); der Code bleibt, Fall `p3-sums-sortiert` (`LC_ALL=C sort -c`), Mutation `P3sort`.
- Wartezeit auf PostgreSQL und auf die Startmeldung (je 60 s) sind Wahl des Codes, nicht des Vertrags; kein Fall prüft die Dauer.

**Beobachtungen**

- Der Lauf endet bei Punkt 5 beim ersten roten Kopf; die Mutationen verdecken die späteren Plattformen, darum sind sie begrenzt (`nur`).
- `binaries.sh` läuft unter dem Gate als Hintergrundjob und ignoriert `SIGINT`; sein INT-Fang gilt nur im Aufruf allein.
- Die Weitergabe von `SIGTERM` ist durch das folgende `wait` nur über den Vermerk der Attrappe (`term`) beobachtbar, nicht über die Zeit.
- Die Gegenprobe hat ein Zehntel der Fälle in `lauft`-Gruppen, die einen echten Gate-Lauf der Attrappe fahren; der Umfang `nur` hält die Mutationsläufe kurz.
- Zeile `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: Schätzung für die Gegenprobe 380 bis 570, gemessen 764 Zeilen für die Gegenprobe am Endstand (Faktor 2,1 gegen die Schätzung; Messung in der Closure des Planners unten).

### Closure des Planners (2026-10-17)

**Gate-Nachweis.** Voller `make gates` am Code-Stand `e365401f` grün (927 s, siehe Tabelle oben); seither liegen nur Review, Verifikation und diese Closure (Positivliste des Stop-Hooks). `make kopf-check` und `make docs-check` laufen nach der Closure nacheinander.

**Gegenstand:** geliefert: `make binaries`, Stufe `dist`, Gate `make test-container` mit Smoke, Gegenprobe `make test-container-gegenprobe`, `SPEC-055` (Architect), Ausnahmeliste und Zeilen in `harness/README.md`.

**Review und Verifikation.** Review (`docs/reviews/2026-10-17-review-slice-v1-abschluss-container.md`): 0 HIGH, 1 MEDIUM, 6 LOW, 5 INFO. Verifikation (`docs/reviews/2026-10-17-verify-slice-v1-abschluss-container.md`): konform, V-1 LOW, V-2 INFO. Vorgabe des Nutzers: LOW und INFO ohne Zusage ohne Test sind „akzeptiert, mit Ausgang in §7“. Ausgang je Befund:

| Befund | Ausgang |
|---|---|
| F-1 (MEDIUM) `-X main.version=$VERSION` im echten `Dockerfile` ohne roten Fall (Vorgabe `dev` maskiert) | behoben `39c07812`: `VERSION := 0.0.0-test` am Gate, `dev` vor `docker` abgelehnt, Fälle `p2-vorgabe-test`, `p2-dev-gate`, `p2-dockerfile-x`; Handlauf `-X main.version` oben und Brechprobe des Verifiers rot; Register `negativtests-fehlen-bei-neuem-vertrag` |
| F-2 (LOW) `p7-kein-volume`/`-user` nur eine Schreibweise | behoben `39c07812` (Kurzflags `-dv`, `-u0`), vom Verifier gebrochen: je nur der eigene Fall rot; Register `negativtests-fehlen-bei-neuem-vertrag` (Aufzählung je Element, `AGENTS.md` §3.10) |
| F-3 (LOW) Kopf nennt Exit 2/1, kein Fall | behoben `39c07812`: Köpfe enger gefasst; Register `zusage-im-kommentar-weiter-als-pruefung` |
| F-4 (LOW) Kopf der Gegenprobe: Zeitstempel „belegt“ | behoben `39c07812`: Kopf nennt den Handlauf als einmalig; Register `zusage-im-kommentar-weiter-als-pruefung` |
| F-5 (LOW) `harness/README.md` „schreibt nichts in den Arbeitsbaum“ weiter als `p10-arbeitsbaum` | erledigt: der Satz in `harness/README.md` §Sensors ist enger gefasst („legt nichts im Arbeitsbaum an“, Punkt 10); Register `zusage-im-kommentar-weiter-als-pruefung` |
| F-6 (LOW) Aufschlüsselung der Messung in §7 veraltet | behoben `62a287cb` (Aufschlüsselung entfernt, Messung 1191 nachgerechnet); Register `plan-folgt-korrektur-nicht` |
| F-7 (LOW) Stage-Verzeichnis von `binaries.sh` im Vertrag nicht genannt | Vertrag erledigt (`SPEC-055` Punkt 3 und 9 nennen das Stage-Verzeichnis); Register `spec-randform-erst-im-review-entschieden` (der Implementer hatte es im Code entschieden und nicht zurückgegeben; `randform-im-code-entschieden-dann-zurueckgegeben` trifft nicht zu: keine Rückgabe) |
| F-8 (INFO) `LC_ALL=C sort` äquivalent | akzeptiert: Invariante und Ort in §7 (Äquivalente und Grenzen), der Verifier prüfte unter vier Locales; zählt nicht zu `gruener-mutant-faelschlich-aequivalent` |
| F-9 (INFO) Bestätigung der Äquivalenz `DIST` als Datei | akzeptiert: kein Befund |
| F-10 (INFO) `grep -q` unter `pipefail` in `bereit` | akzeptiert: Negativ mit Grund in §7; flackert das Gate, ist das ein Befund, kein Folge-Slice |
| F-11 (INFO) Wartezeiten, Mutationen ohne Tabellenzeile | akzeptiert: Dauer ist Wahl des Codes, keine Zusage; die Mutationen sind rot (Review gefahren) |
| F-12 (INFO) Größe, Aufbau der Gegenprobe | akzeptiert: Messung 1191 von 1200 nachgerechnet (Verifikation); kein Eingriff |
| V-1 (LOW) §7 ohne Läufe am Endstand | behoben mit dieser Closure: Stand `e365401f`, Läufe und Dauer in der Tabelle oben |
| V-2 (INFO) `SPEC-055` Punkt 2 „eine Variable der Umgebung gilt nicht“ ohne Fall | **Hinweis an den Architect:** bei der nächsten Berührung von `SPEC-055` den Satz enger fassen oder streichen oder einen Fall nachziehen (`AGENTS.md` §3.10, §3.11); die Spezifikation ändert der Planner nicht, und eine Änderung unter `harness/` verlangt einen vollen `make gates`. Ausgang: akzeptiert; die Zusage hat keinen Fall und steht hier benannt (Register `zusage-im-kommentar-weiter-als-pruefung`, derselbe Eintrag wie F-3 bis F-5, ein Evidence je Vorgang) |

**Messung der Marke** (§4, Metrik `AGENTS.md` §3.14): 1191 (1176 Einfügungen, 15 Löschungen) gegen Marke 1200; Meldeschwelle 1000 überschritten und gemeldet; der Schnitt in `-container-doku` stand im Drift-Log der Roadmap (2026-10-17). Die Rückführung `in-progress` → `next` trat nicht ein. Schätzung des Architect 780 bis 1080 (Produktcode 190, Gegenprobe 380 bis 570, Ausnahmeliste und Vertrag 30 bis 40): Obergrenze um 111 überschritten; Gegenprobe 764 statt 380 bis 570 (Faktor 2,1 gegen den angesetzten Doppelten des Produktcodes). Nach der Zählregel vom 2026-10-16 ist das ein Auftreten von `schaetzung-in-zeilen-unterschaetzt-tests`.

**Risiken aus §6 — Ausgang je Risiko** (die Zeilen stehen in §6 am Risiko): (1) Handbuch beschreibt vor der Software *entfallen* (Schnitt, steht im Nehmer); (2) [ADR-0043](../../adr/0043-binaries-und-container-gate.md) `Proposed` *entfallen* (`Accepted`); (3) Gate-Laufzeit *weiter offen* → Register `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (+1; abweichend vom Vorschlag der Verifikation, *entfallen*: Der Zuwachs von rund 5 Minuten ist eingetreten und gehört zur Klasse, die der Eintrag zählt); (4) Abdeckung *entfallen* (`make abdeckung-check` verlangt keinen Beleg, [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)). Ob die Ausgänge tragen: (1) und (2) am Nehmer und an der ADR (`Accepted`); (3) an der Evidence-Datei; (4) am Lauf von `make abdeckung-check`, den die Verifikation wiederholte.

**Register** (Zähler = Dateien unter `evidence/`; alle `verkörpert` oder `geplant` und bleiben): `negativtests-fehlen-bei-neuem-vertrag` 45 → 46 (F-1, F-2); `zusage-im-kommentar-weiter-als-pruefung` 69 → 70 (F-3, F-4, F-5, V-2); `spec-randform-erst-im-review-entschieden` 38 → 39 (F-7); `plan-folgt-korrektur-nicht` 58 → 59 (F-6); `schaetzung-in-zeilen-unterschaetzt-tests` 9 → 10 (1191 gegen 780 bis 1080); `gate-laufzeit-waechst-mit-gegenprobe` 5 → 6 (Risiko 3); `gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig` 3 → 4 (Flake `p7-int-treiber` im zweiten Gate-Lauf, Ausgang `geplant` auf `slice-harness-gegenprobe-matrix-p7-flanke` bleibt). Vermerk ohne Zuwachs: `schnitt-laesst-haelfte-an-der-grenze` und `groesse-erst-in-der-pruefung-des-architect-sichtbar` (der Schnitt kam beim Erreichen der Meldeschwelle vor der Marke, die Größe war vor der Prüfung des Architect bekannt); `gruener-mutant-faelschlich-aequivalent` (F-8 bestätigt äquivalent). Kein Eintrag erreicht neu 3× ohne Verkörperung und ohne `geplant`. Die Verifikation-Klasse „Zusage am echten Artefakt, die eine Attrappe nur mittelbar sieht“ passt in `negativtests-fehlen-bei-neuem-vertrag` (Zusage ohne fangende Mutation), kein eigener Eintrag.

**Lerneintrag — neuer Sensor und eine Beobachtung.** Neuer Sensor: `make binaries`, `make test-container` und `make test-container-gegenprobe` ([ADR-0043](../../adr/0043-binaries-und-container-gate.md), Herkunfts-Anker `seit slice-v1-abschluss-container` in `harness/README.md` §Sensors). Beobachtung, keine geänderte Regel: Ein Vorgabewert des Codes (`VERSION=dev`) maskierte eine Zusage am echten `Dockerfile`, die die Attrappe nur mittelbar sah; erst der Wert, der nicht die Vorgabe ist (`0.0.0-test`), plus Smoke am echten Artefakt fängt den Bruch. Wer eine Zusage am echten Artefakt gibt, wählt am Gate einen Wert, der von der Vorgabe des Codes abweicht (`AGENTS.md` §3.10 trägt die Klasse, im Register kein eigener Eintrag). Spec-Lücke benannt: `SPEC-055` Punkt 2, Satz „eine Variable der Umgebung gilt nicht“ ohne Fall (V-2, beim Architect).

**Paarungen.** (a) Anker: Der Lerneintrag nennt `seit slice-v1-abschluss-container` für den Sensor; die Zeilen in `harness/README.md` tragen ihn. (b) Folge-Slice: Kein Risiko hat den Ausgang *eingetreten*; die Sendung an `slice-v1-abschluss-container-doku` (§1 *Abgegeben*) nimmt der Nehmer an (`grep` nach dieser Kennung findet sie dort in §1, §6 und §8). (c) Register: Alle genannten `BEO-REPO/…` bestehen als Verzeichnis mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-17): Das Register ist seither gefüllt, aber kein Eintrag nennt Container, Image oder reproduzierbaren Build; einschlägig sind `BEO-REPO/verhalten-nur-unter-linux-geprueft` und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` für Aussagen des Handbuchs über Umgebungen, sowie `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` für die Schätzung in §4.

**Nachgezählt beim Schnitt** (2026-10-17, Planner, `AGENTS.md` §3.13; die Sendung geht an `slice-v1-abschluss-container-doku`, dort ebenfalls nachgezählt): ein Liefer-Punkt (§2: Container-Bau mit Gate: `make binaries`, `make test-container`, Image, Smoke, Gegenprobe), vorher zwei; Schichten: eine (Container-Bau und Harness: `Dockerfile`, Fragmente, Skripte, Vertrag, Gate-Index). Der Produktcode (`cmd/`, `internal/`) wird nicht berührt. Die frühere Zählung (zwei Schichten, die Dokumentation als zweite, 2026-10-10 bis 2026-10-17) gilt für den Nehmer. `dist/` steht in `.gitignore` (außerhalb der Pfade der Marke), damit `make binaries` den Hash des Arbeitsbaums nicht ändert.

Die Prüfung des Architect vom 2026-10-17 (zwei Liefer-Punkte, zwei Schichten) und die Zählung mit der Sendung aus `slice-harness-meldungskatalog-gate` sind im Schnitt aufgegangen; eine Sendung an andere Slices geht außer an `slice-v1-abschluss-container-doku` nicht hinaus: `slice-erster-release-veroeffentlichung` ruft `make binaries`, das ist eine Voraussetzung dort (jetzt auch von `-container-doku`), kein Eintrag hier.

**Stand bei der Closure** (2026-10-17): Die Zähler der einschlägigen Einträge stehen in §7 (*Register*); `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` ist der Eintrag, den dieser Slice über die Spanne des Architect führte.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Sub-Area `REPO` ist GF).
