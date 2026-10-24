# Slice slice-go-toolchain-1-27-3: Go-Builder auf 1.27.3 heben, damit der Scan des Release-Images den Befund der Standardbibliothek nicht mehr meldet

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-24). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist, und keine Anforderung einer Welle hängt an ihm.

**Bezug:** Kein Anspruch des Produkts (`LH-*`): Der Slice hebt den Pin eines Build-Werkzeugs. Auslöser: Der Trockenlauf des Release `v0.2.0` (GitHub-Lauf `40071373625`, Commit `d720e44f`) ist im Schritt `Scan` rot (`CVE-2026-111943`, Paket `stdlib`, eingebaut `v1.27.2`, behoben in `1.26.10` und `1.27.3`, `HIGH`, für `linux/amd64` und `linux/arm64`). Aktive ADRs: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Builder im `Dockerfile`, Pins werden bewusst angehoben), [ADR-0052](../../adr/0052-scan-des-images-vor-dem-push.md) (Schwelle und Behebung statt Ausnahme).

**Berührte Spec-Stellen:** [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 16 (der Scan; gelesen, nicht geändert) · [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) (Bau und Summen der Binaries; gelesen, nicht geändert).

**Verantwortlich:** pt9912 (Zweig impl-go127)
**Autor:** Planner. **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Pin des Go-Builders steht an beiden Stellen (`Dockerfile`, `harness/mk/build.mk` `GO_IMAGE`) auf `golang:1.27.3-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c` (Index-Digest, am 2026-10-24 mit `docker buildx imagetools inspect` gelesen), sodass die Binaries und das Image mit Go 1.27.3 gebaut werden und der Scan des Release-Pfads den Befund `CVE-2026-111943` nicht mehr meldet.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Tag, `VERSION`, Release-Commit, Lauf des Tags, Eintrag `freigaben/0.2.0.md` — `slice-release-v0-2-0` (in `in-progress/`); dort bleiben Tag und Ja des Nutzers. Dieser Slice setzt keinen Tag, startet keinen Lauf und pusht nicht (andere Rolle, andere Handlung).
- Eine Ausnahme für den Befund oder eine Senkung der Schwelle des Scans — `SPEC-067` Punkt 16 und [ADR-0052](../../adr/0052-scan-des-images-vor-dem-push.md) lassen keine zu; ein Befund mit Behebung wird behoben (`AGENTS.md` §3.6).
- Eine neue ADR — das Heben eines Pins lockert kein Gate; [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) nennt als Folgepflicht, dass Pins bewusst angehoben werden, und die Kopfzeile des `Dockerfile` sagt dasselbe („eine Anhebung ist ein bewusster Commit“).
- `go.mod` (`go 1.27`) — die `go`-Zeile ist die Untergrenze; Toolchain liefert das Image (`GOTOOLCHAIN=local`), eine Patch-Version in `go.mod` änderte die Zusage an Nutzer des Moduls ohne Grund. Bestand bleibt bewusst stehen.
- Die Pins von `golangci-lint`, `gcr.io/distroless/static`, PostgreSQL, a-check, d-migrate, Trivy, BuildKit — der Scan meldet für sie nichts; ein Heben ohne Befund ist ein anderer Vorgang (Nachtlauf: `slice-harness-ci-freshness-images`, dort „Das Heben eines Pins … ein Slice je Hebung; dieser Slice meldet“, also keine Sendung an ihn).
- Der Nachtlauf gegen alternde Pins — `slice-harness-ci-freshness-images` (Klasse 1, Adresse). Ein Sensor, der beide Pin-Stellen abgleicht, ist **nicht übernommen**: `slice-harness-ci-freshness-images` schließt den Abgleich der Stellen in §1 aus (jede Stelle einzeln geprüft), kein anderer Slice nimmt ihn an (`AGENTS.md` §3.13; Bestand bleibt, Klasse 2: der Abgleich ist hier eine Handlung, `grep` in der DoD, kein Gate; ein Sensor wäre ein neuer Slice des Planners).
- Das Ausführen der `darwin`- und `windows`-Binaries und von `linux/arm64` — Grenze wie bei jedem Release (§6 Randform Plattform); der Slice ändert nur die Toolchain, nicht die Prüftiefe.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkte:

- [x] **(1) Pin an beiden Stellen.** `Dockerfile` (Stufe `deps`) und `harness/mk/build.mk` (`GO_IMAGE`) tragen `golang:1.27.3-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c`. Beleg: `grep -rn 'golang:1\.27' Dockerfile harness/mk` nennt genau diese zwei Zeilen mit demselben Text, und `grep -rn 'f92b6ef800e4' --exclude-dir=.git --exclude-dir=baseline .` findet im Baum keinen Pin mehr (Erwähnungen in Plänen und Belegen unter `docs/` sind Chronik und bleiben). Der Digest ist am Tag der Lieferung neu gegen die Registry gelesen: `docker buildx imagetools inspect golang:1.27.3-alpine` nennt den Index-Digest oben (Ausgabe im Bericht).
- [ ] **(2) Der Scan-Befund ist weg.** Beleg ist der Trockenlauf des Release auf dem Stand dieses Slice (`workflow_dispatch`, ohne Tag; Auslöser und Push sind Handlung der Koordination, der Implementer pusht nicht): Lauf-Kennung, Hash, Zeile `release: scan: 0 Befund(e) über der Schwelle, <m> ohne Behebung` im Log, und kein Befund `CVE-2026-111943` für `linux/amd64` und `linux/arm64`; Lauf-Kennung, Hash und Zeile stehen in §7. Lokal läuft kein Scanner (`SPEC-067` Punkt 16: Trivy nur in der Pipeline); die Gegenprobe der Pipeline belegt den Scan nicht (`harness/README.md` §Sensors, `release-gegenprobe`), nur der Lauf. Ist der Lauf rot aus einem anderen Grund, gilt §4 (`in-progress` → `open`), nicht ein Haken.

Gates und Closure:

- [x] `make gates` grün am Endstand (darin `make test-container`: zwei Builds der sechs Binaries mit gleichen Summen, [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md)); Hash, Datum, Exit, Dauer im Bericht.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `evidence/slice-go-toolchain-1-27-3.md` in `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder/` (alternder Go-Pin, vierte Datei); **kein Zähler wird gesetzt**, er folgt aus den Dateien.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` (Zeile 9, Stufe `deps`; alle Stufen `test`, `build`, `integration` leiten von ihr ab) | update | einzige `golang`-Stelle im `Dockerfile` (gelesen 2026-10-24); `golangci-lint` und `distroless` bleiben |
| `harness/mk/build.mk` (Zeile 15, `GO_IMAGE`) | update | zweite Stelle, nur für `make go-mod-tidy`; derselbe Text wie im `Dockerfile` |

Geprüft, nichts zu ändern: `go.mod` (`go 1.27`, keine `toolchain`-Zeile), `tools/` (kein Pin des Builders; `grep` nach `f92b6ef800e4` und `golang:1.` ohne `docs/`), `docs/user/` und `docs/maintainer/` (nennen keine Go-Version).

Ansatz:

- Tag und Digest stehen zusammen (`golang:1.27.3-alpine@sha256:…`), damit der Pin die gemessene Patch-Version nennt; vorher stand `1.27-alpine`, das auf 1.27.2 zeigte.
- Der Implementer liest den Digest selbst neu (DoD 1), statt ihn aus diesem Plan zu kopieren.
- Ein Test entsteht nicht: Der Pin ist Text, keine Zusage an einer Wire-, Datei- oder Formatgrenze und kein Vertrag eines Gates (`AGENTS.md` §3.10); der Beleg ist der Gate-Lauf und der Scan.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Nutzer hat der Koordination den Auftrag gegeben (Trockenlauf rot, Befund beheben, kein Tag). Der Slice blockiert `slice-release-v0-2-0` (Scan ist dort DoD 6); er läuft vor dem Tag.

**Marke 1500** (Metrik: `AGENTS.md` §3.14, dort der Befehl; Pfade des Plans: `Dockerfile harness/mk`; Basis: der Commit des Übergangs nach `in-progress/`, Hash beim Übergang eintragen: `7cfad5ab`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: zu groß, wenn die Summe die Marke 1500 erreicht; erwartet sind zwei Zeilen je Stelle, ein Diff darüber hinaus heißt, dass ein anderer Gegenstand mitkam.
- `in-progress` → `open`: (i) `make gates` bleibt rot aus einem Grund, der nicht der Pin ist (Gate gegen Go 1.27.3: `go vet`, Lint-Profil, Tests; Befund wird dem Architect oder Planner zurückgegeben, kein Slice-Wachstum). (ii) Der Trockenlauf auf dem neuen Stand meldet weiter einen Befund mit Behebung (anderer CVE, etwa im Basis-Image): Entscheidung des Nutzers, ein weiterer Pin ist ein weiterer Slice; keine Ausnahme (`SPEC-067` Punkt 16). (iii) Die Koordination startet den Trockenlauf nicht: ein Slice, der wartet, hält das WIP-Limit.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. Beobachtbar: (a) `make gates` grün am Endstand, (b) der Trockenlauf des Release auf dem Stand dieses Slice endet im Schritt `Scan` mit `release: scan: 0 Befund(e) über der Schwelle`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keine Wire-, Datei- oder Formatgrenze und keinen Vertrag; die Liste nennt, was sich am Bau ändert.

| Randform | Entscheidung | Stelle |
|---|---|---|
| (a) Plattform | Der Bau im Gate (`make test-container`) erzeugt die sechs Binaries `linux`, `darwin`, `windows`, je `amd64` und `arm64`, und das Image je Linux-Plattform, alle mit derselben Toolchain; geprüft ist Kopf je Format und Summe, ausgeführt wird der Smoke auf dem Image der Plattform des Docker-Servers. Der Scan prüft nur das Image (`linux/amd64`, `linux/arm64`); die Standardbibliothek der `darwin`- und `windows`-Binaries scannt niemand, sie entstehen aber aus demselben Pin. `darwin`, `windows` und `linux/arm64` sind gebaut, nicht ausgeführt: Grenze wie in `slice-release-v0-2-0` (h) | `SPEC-055`; `SPEC-067` Punkt 15; `LH-QA-03` |
| (b) Summen | Die Summen der Binaries ändern sich gegenüber den Binaries mit 1.27.2 (andere Toolchain); zugesagt ist nur Gleichheit zweier Builds mit demselben Pin, nicht über Toolchains hinweg | `SPEC-055` (Gleichheit der zwei Builds) |
| (c) Digest | Index-Digest der Manifestliste, nicht der einer Plattform; die Stufen `deps`ff. laufen mit `--platform=$BUILDPLATFORM` und lösen die Plattform aus dem Index | Bestand, Form wie bisher |
| (d) Zwei Stellen | Es gibt keinen Sensor, der die zwei Pin-Stellen abgleicht; DoD 1 nennt den `grep` | kein Sensor, kein Slice nimmt den Abgleich an (`slice-harness-ci-freshness-images` §1 schließt ihn aus; Bestand bleibt) |
| (e) Aufrufarten, Abbruch | entfällt: kein Werkzeug, kein Ergebnis, das ein Hook oder Gate später liest; der Pin ist Text | `AGENTS.md` §3.12 |
| (f) Lint-Stufe | `golangci-lint` bringt seine eigene Go-Version mit und baut nichts, was ausgeliefert wird; sie läuft mit `GOTOOLCHAIN=local` gegen die `go`-Zeile `1.27`; bleibt unverändert, ein Scheitern wäre Rückgabe (§4 (i)) | `SPEC-049` Punkt 2 |

**Risiken:**

- Der Trockenlauf meldet nach der Hebung einen weiteren Befund über der Schwelle (anderer CVE der Standardbibliothek, des Basis-Images `distroless` oder von `golang` selbst) — **Ausgang:** offen bis zur Closure.
- Go 1.27.3 färbt ein Gate rot (`go vet`, Lint, Test) — **Ausgang:** offen bis zur Closure; erwartet gering (Patch-Stand), Folge wäre §4 (i).
- Der Pin wird nur an einer der beiden Stellen gehoben, weil kein Sensor sie abgleicht — **Ausgang:** offen bis zur Closure; Gegenmittel ist DoD 1 (`grep`).
- Der Pin altert wieder, und erst der Scan des nächsten Release meldet es (viertes Auftreten von `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`, davon zweites mit dem Go-Pin) — **Ausgang:** weiter offen: → `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder` (dort `geplant`: `slice-harness-ci-freshness-images`, `open/`); die Datei `evidence/slice-go-toolchain-1-27-3.md` ist DoD.
- Der Trockenlauf läuft auf einem Stand, der nicht der Stand des Tags ist (`slice-release-v0-2-0` Randform (c): `R` soll die Spitze von `main` sein) — **Ausgang:** eingetreten durch den Schnitt, wenn dieser Slice nach `R` landet; die Anpassung steht in `slice-release-v0-2-0` (Planner, nicht dieser Slice).

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

(noch nicht geschlossen)

**Belege des Implementers** (2026-10-24):

- Digest neu gelesen: `docker buildx imagetools inspect golang:1.27.3-alpine` nennt `Digest: sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c`; gleich dem Plan.
- Pin-Commit `143e9994`. `grep -rn 'golang:1\.27' Dockerfile harness/mk` nennt genau `Dockerfile:9` und `harness/mk/build.mk:15` mit demselben Text `golang:1.27.3-alpine@sha256:b8c9…505e9c`; `grep -rn 'f92b6ef800e4' --exclude-dir=.git --exclude-dir=baseline --exclude-dir=docs .` ohne Treffer (Exit 1); unter `docs/` nur dieser Plan (Chronik).
- `make gates` am Stand `143e9994`: Start 2026-10-24T06:51:45Z, Ende 07:00:33Z, Dauer 528 s, Exit 0 (darin `test-container`, `test-integration`, `test-postgres-versionen`, `lint`, `workflows-check`).
- Mutation: keine (Plan §6, Pin ist Text, `AGENTS.md` §3.10). DoD 2 (Scan) belegt die Koordination über den Trockenlauf; nicht abgehakt.


## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `REPO` (gesamtes Repo, `harness/conventions.md`); die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (2026-10-24, gemergter Stand). Treffer: `release-pipeline-pin-und-werkzeug-altern-ohne-melder` (drei Dateien unter `evidence/`, Zähler 3, Ausgang `geplant`: `slice-harness-ci-freshness-images`; dessen zweites Auftreten war derselbe Go-Pin, `slice-erster-release-freigabe`; mit der Closure dieses Slice kommt die vierte Datei dazu, der Eintrag ist keine Notiz mehr, sondern trägt seine Adresse und wird dort weiter geführt).

**Zählung** (Planner, 2026-10-24, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 2** (1: Pin an beiden Stellen; 2: Beleg des Scans im Trockenlauf). **Schichten: eine** (Build-Werkzeug: `Dockerfile` und Fragment `harness/mk/build.mk`); Doku (`docs/user/`, `docs/maintainer/`, `README.md`) wird nicht berührt, Planung zählt nicht, Produkt-Code und Tests bleiben. **Sendungen an andere Slices:** keine (die Hebung ist die Handlung, die `slice-harness-ci-freshness-images` §1 ausdrücklich einem eigenen Slice zuweist).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
