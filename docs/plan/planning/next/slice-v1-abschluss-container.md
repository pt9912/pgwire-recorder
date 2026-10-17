# Slice slice-v1-abschluss-container: Container-Image und Betriebsdokumentation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)

**Berührte Spec-Stellen:** `LH-FA-16.a` · `SPEC-031` · `SPEC-035` · [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) (Vertrag dieses Slice, Architect 2026-10-17) · [`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (gelesen: das Gate hält die zwei Handbuch-Tabellen)

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

**Ziel:** Ein Docker/OCI-Image (`linux/amd64`, `linux/arm64`) und ein reproduzierbarer Binary-Build existieren, und die Betriebsdokumentation nennt Optionen, Exit-Codes, Meldungskatalog und verwendete Umgebungen.

**Zuschnitt nach dem Vorlauf des Architect (2026-10-17, [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md), [ADR-0043](../../adr/0043-binaries-und-container-gate.md) `Accepted`, Entscheidung des Nutzers vom 2026-10-17):** Das Werkzeug `make binaries` baut die sechs Binaries (`linux`, `darwin`, `windows`, je `amd64`, `arm64`) mit `SHA256SUMS`; das Gate `make test-container` baut sie zweimal ohne Cache der Bau-Stufen und vergleicht die Summen, prüft den Kopf jedes Binaries, baut das Image je Linux-Plattform (mit demselben Binary) und fährt auf der Plattform des Docker-Servers den Smoke im Container (`version`, `record` gegen PostgreSQL, `replay` ohne PostgreSQL, `SELECT 1`). Neu sind das Werkzeug, das Gate, der Smoke und die Stufe `dist`; vorhanden sind das Multistage-`Dockerfile` mit `-trimpath`, `-buildvcs=false` und gepinnten Images, `make build` und die Stufe `runtime` (ein Probe-Build von `windows/arm64` und `darwin/amd64` baute, `windows/arm64` zweimal ohne Cache der Stufe `build` mit derselben Summe).

**Übernommen aus `slice-v1-abschluss-herunterfahren` (aus `slice-v1-abschluss-betrieb` hervorgegangen) und `slice-v1-abschluss-homebrew`** (dort je §1, Abgrenzung): das Container-Image (Docker/OCI); das ist das Ziel oben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Veröffentlichungs-Mechanismus der Registries (`ghcr.io`, `docker.io`) — Gegenstand des Release-Verfahrens in `docs/maintainer/releasing.md`; dieser Slice liefert das Image.
- Die Manifestliste für beide Linux-Plattformen in einer Referenz und die Anhänge am GitHub-Release — `slice-erster-release-veroeffentlichung`: Der Standard-Treiber von Docker baut keine Manifestliste, und sie gehört zum Push (`SPEC-055` Punkt 13). Jener Slice ruft `make binaries` und hängt `dist/` mit `SHA256SUMS` an, statt ein zweites Mal zu bauen (Hinweis des Architect; Voraussetzung in dessen §1 eingetragen, 2026-10-17).
- Der Smoke des Images auf der jeweils anderen Plattform (Emulator) und ein Lauf der Binaries für `darwin` und `windows` — kein Emulator und keine Gegenstelle in der Prüfumgebung; belegt sind Bau, Kopf und Summe (`SPEC-055` Punkt 13, Grenze 3 und 4).
- Reproduzierbarkeit des Images — `LH-QA-01` verlangt sie nur für das Binary (`SPEC-055` Punkt 13).
- Gate für Code-Tabelle und Katalog — liefert `slice-harness-meldungskatalog-gate` (wellenlos, Entscheidung des Nutzers vom 2026-10-07) direkt vor diesem Slice (Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)); das Gate liegt in `done/` und prüft seit `slice-harness-meldungskatalog-gate` (`SPEC-054` Punkt 1, 4c, 4d) auch die Tabellen `### Fehlercodes` und `### Warnungen` des Handbuchs: Ihre Codes sind genau die der Code-Tabelle, andere Betriebsdokumente liest es nicht. Schreibt dieser Slice diese Tabellen fort, hält er `make meldungskatalog-check` grün; das Gate selbst hält die Gleichheit, ein eigener Liefer-Punkt dafür entfällt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Das Image je Linux-Plattform ist gebaut und trägt dasselbe Binary wie `make binaries`; das Image startet `replay` ohne PostgreSQL und besteht den Walking-Skeleton-Smoke im Container. Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung). Aus `slice-doku-ist-stand` (dort §6, *Installation*): Der Slice liefert im Handbuch §2 *Container* und *Das Binary aus dem Image* (Image und seine Plattformen, statt „für die Architektur des Rechners, auf dem Sie gebaut haben“) und die Plattform-Zeile in §1 *Voraussetzungen* („Zum Bauen: … Zum Ausführen: Linux auf der Architektur des Rechners“), soweit das gebaute Image sie belegt.
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus): Das Binary ist für alle Zielplattformen aus `SPEC-035` reproduzierbar gebaut (zwei Builds desselben Quellstands, dieselbe Prüfsumme, `make test-container` nach `SPEC-055`), `make binaries` schreibt es mit `SHA256SUMS`; Reproduzierbarkeit verlangt `LH-QA-01` nur für das Binary. Die Betriebsdokumentation ist vorhanden und benennt die Abnahme-Umgebungen.
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
| `Dockerfile` | update | Stufe `dist` (nur die Binaries, `FROM scratch`); die Bau-Stufe bleibt |
| `harness/mk/container.mk`, `tools/harness/binaries.sh`, `tools/test/run-container.sh`, `tools/test/container-gegenprobe.sh` | neu | `binaries`, `test-container`, `test-container-gegenprobe` nach `SPEC-055` |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh`, `harness/vertraege/SPEC-052-gate-kette.md`, `harness/README.md` | update | `test-container` und seine Gegenprobe in die Ausnahmeliste (SPEC-052 Punkt 2 und 4); Zeilen in §Sensors |
| `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md` | update | auf den Stand der umgesetzten Funktionen prüfen und anpassen; Optionen, Exit-Codes, Meldungskatalog stehen im Handbuch, neu sind die Abnahme-Umgebungen (`LH-QA-03`) und §2 *Installation* |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die anderen Slices der Welle sind `done`, `slice-harness-meldungskatalog-gate` liegt in `done/` (§1), und [ADR-0043](../../adr/0043-binaries-und-container-gate.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-17; erfüllt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- Dockerfile harness tools test docs/user docs/maintainer README.md` liegt über 1200; Basis ist `d258c0b7` (HEAD vor dem Anspruch; der Vertrag `SPEC-055` unter `harness/vertraege/` und die ADR liegen davor, der Plan-Commit und der Move berühren die Pfade nicht). Schätzung des Architect (2026-10-17, in dieser Metrik, `AGENTS.md` §3.14): Spanne 780 bis 1080, Mitte 930, Obergrenze 1080 (unter der Marke), Meldeschwelle 1000 bleibt. Aufteilung: Produktcode 190 (`Dockerfile` 12, `container.mk` 35, `binaries.sh` 45, `run-container.sh` 100), Gegenprobe 380 bis 570 (Faktor Zwei bis Drei des Produktcodes, Attrappe von `docker` ohne Docker-Lauf), Ausnahmeliste samt Gegenprobe, Zeilen in `SPEC-052` und `harness/README.md` 30 bis 40 (der Vertrag `SPEC-055` und die [ADR-0043](../../adr/0043-binaries-und-container-gate.md) stehen im Vorlauf vor der Basis und zählen nicht mit), Handbuch, README und `releasing.md` 180 bis 280. Die Marke 1200 bleibt: Ein Schnitt vor dem Code ist nicht geboten; der natürliche Schnitt, falls die Marke eintritt, liegt zwischen den zwei Liefer-Punkten (Image mit Smoke und Handbuch | Binaries mit Gate und Betriebsdokumentation). Wer 1000 überschreitet, meldet es vor dem Handbuch.
- `in-progress` → `next`: die Doku wächst über einen Liefer-Punkt hinaus — zurück zur Zerlegung.
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
- Fehlende, leere oder unzulässige `VERSION` (Punkt 2); Wert in `-ldflags` und Dateiname.
- `DIST` vorhanden und nicht leer, Teilmenge bei einem fehlgeschlagenen Build, Format von `SHA256SUMS`, Endung `.exe` (Punkt 3).
- Zwei Builds: abweichende Summe je Plattform (sechs Fälle), Umgehen des Caches der Stufen `build` und `dist` als Voraussetzung dafür, dass der zweite Build ein Build ist (Punkt 4). Akzeptiertes Negativ: Das Gate vergleicht auf demselben Docker-Server mit denselben Eingangs-Images; ein zweiter Host und der Pin der Images werden nicht geprüft (Punkt 13, Grenze 2), weil Anhebungen bewusste Commits sind ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)) und es keine zweite Umgebung zur Abnahme gibt.
- Zeitstempel und VCS-Stempel: kein Fall; sie gehen nicht in das Binary ein (`-buildvcs=false`, Go schreibt keinen Zeitstempel), und ein Build mit eingeschleustem Zeitstempel wird von Punkt 4 rot (der Implementer zeigt es einmal von Hand und nennt es im Bericht; keine Gegenprobe mit echtem Docker, Punkt 12).
- Netz im Build: nur die Stufe `deps` hat Netz, `build` läuft mit `--network=none` ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)); der Bau mit kaltem Cache braucht einmal Netz, wie `make build` (kein Fall).
- Kopf der Binaries je Format und Architektur, anderes Format (Punkt 5).
- Image: falsche Plattform in den Metadaten, abweichende Summe der Datei im Image, Bau ohne Emulation (Punkt 6).
- Smoke: Plattform des Docker-Servers außerhalb von `amd64` und `arm64`; `version`; Aufzeichnen und Wiedergeben, leere und ungleiche Sicht, Exit-Code 0 nach `SIGTERM`, PostgreSQL vor dem Wiedergeben entfernt, kein `--user`, Volume oder Bind-Mount (Punkt 7); Schreibrecht des Benutzers `nonroot` auf `/tmp` ist damit belegt.
- Reihenfolge der Gründe bei mehreren roten Punkten (Punkt 8).
- Abbruch und Herunterfahren: Aufräumen nach rotem Lauf, `SIGINT`, `SIGTERM` (Punkt 9). Akzeptiertes Negativ: `SIGKILL` und Ausfall des Docker-Dienstes (Punkt 13, Grenze 6), wie bei `SPEC-053`.
- Aufrufarten von `make`: `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung, auch in Kombination, je Ziel `binaries` und `test-container` (Punkt 11); Aufruf-Voraussetzung: nacheinander, `make -j gates` nicht zugesagt (Punkt 13, Grenze 5).
- Ausnahmeliste der schnellen Kette und Lesen des Arbeitsbaums (Punkte 10 und 12).
- Verhalten je Serverversion: Der Smoke läuft gegen die Referenz-PostgreSQL-Version (`POSTGRES_IMAGE`); die Matrix bleibt bei `SPEC-053`, kein Fall.
- Die Manifestliste, die Registries und der Smoke des Images der anderen Plattform liegen nicht in diesem Slice (§1; `slice-erster-release-veroeffentlichung`, `SPEC-055` Punkt 13).

**Risiken:**

- Handbuch und Releasing-Dokument beschreiben Verhalten vor der Software; Abweichungen zeigen sich erst im Bau — **Ausgang:** offen bis Closure.
- Die [ADR-0043](../../adr/0043-binaries-und-container-gate.md) war `Proposed` — **Ausgang:** entfallen; der Nutzer hat sie am 2026-10-17 angenommen (`Accepted`).
- Die Gate-Laufzeit wächst um die Zahl der Builds (zwölf der Stufe `build` ohne Cache, zwei Images, ein Smoke); der Lauf des Gates in `make gates` ist die Probe — **Ausgang:** offen bis Closure; der Bericht nennt die Dauer.
- Die Abdeckungstabellen (`make abdeckung-check`, `SPEC-048`, `[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)`) könnten für `LH-FA-16`, `LH-QA-01` oder `LH-QA-03` einen Beleg des Gates verlangen; der Implementer prüft es mit dem ersten Lauf und deklariert dann nach `[ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md)` im Kopf der Gegenprobe, ohne den Vertrag von `SPEC-048` zu ändern — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-17): Das Register ist seither gefüllt, aber kein Eintrag nennt Container, Image oder reproduzierbaren Build; einschlägig sind `BEO-REPO/verhalten-nur-unter-linux-geprueft` und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` für Aussagen des Handbuchs über Umgebungen, sowie `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` für die Schätzung in §4.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (Dockerfile mit Fragmenten, Dokumentation).

**Nachgezählt** mit der Sendung aus `slice-harness-meldungskatalog-gate` (2026-10-17, `AGENTS.md` §3.13; Gate und Handbuch-Tabellen in §1, Marke und Schätzung in §4): zwei Liefer-Punkte (§2: Image mit Smoke und Handbuch, reproduzierbarer Build mit Betriebsdokumentation), unverändert; Schichten: zwei (Dockerfile mit Fragmenten, Dokumentation), die Prüfung der Reproduzierbarkeit und der Smoke gehören zur ersten; die Bindung ist kein Liefer-Punkt, das Gate hält sie. Die Vorgabe der Marke steht beim Architect.

**Nachgezählt** vom Architect (2026-10-17, `AGENTS.md` §3.13, §3.14): zwei Liefer-Punkte (§2: Image mit Smoke und Handbuch | Binaries, Gate und Betriebsdokumentation), unverändert, kein dritter; Schichten zwei: Container-Bau mit Gate (`Dockerfile`, Fragment, Skripte; Smoke und Gegenprobe gehören dazu) und Dokumentation (Handbuch, README, `releasing.md`). Der Produktcode (`cmd/`, `internal/`) wird nicht berührt; der Bau für `darwin` und `windows` ist mit dem Probe-Build belegt. Eine Sendung an andere Slices geht nicht hinaus: `slice-erster-release-veroeffentlichung` ruft `make binaries`, das ist eine Voraussetzung dort, kein Eintrag hier. `dist/` gehört in eine neue `.gitignore` (außerhalb der Pfade der Marke), damit `make binaries` den Hash des Arbeitsbaums nicht ändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
