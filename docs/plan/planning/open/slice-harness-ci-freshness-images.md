# Slice slice-harness-ci-freshness-images: Aktualität der Image-Pins im Nachtlauf — Tag-Achse und Digest-Achse

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-ci-nachtlaeufe (Entscheidung des Nutzers vom 2026-10-18), dritte Stelle der Reihenfolge, direkt hinter `slice-harness-ci-freshness` ([welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) §5; bis zur Entscheidung vom 2026-10-18 fünfte Stelle von welle-erster-release, die Welle verlässt dieser Slice). Start-Trigger weicht von der Welle ab: `slice-harness-ci-freshness` in `done/` (Skript, Workflow und Vertrag, die dieser Slice erweitert); dazu der Start-Trigger der Welle, `welle-erster-release` in `done/`.

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice lässt die Pins der Werkzeug- und Produkt-Images nachts gegen ihre Registries prüfen. Aktive ADRs: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Images im `Dockerfile` und in `harness/mk/`, Docker-only), [ADR-0042](../../adr/0042-versionsmatrix-der-integrationstests.md) (die Liste der PostgreSQL-Versionen mit Digest). Eine neue ADR schreibt der Architect nur, wenn eine Entscheidung dieses Slice nicht schon die ADR von `slice-harness-ci-freshness` trägt (Randform 1).

**Berührte Spec-Stellen:** `SPEC-066` (Erweiterung um zwei Arten der Frage und die Leser der Pins; Vertrag in `harness/vertraege/`, vom Architect vor dem ersten Code-Commit, gleiche Nummer wie `slice-harness-ci-freshness`) · `SPEC-063` und `SPEC-064` (gelesen) · [`SPEC-053`](../../../../harness/vertraege/SPEC-053-postgres-versionen.md) (gelesen: die Liste der Versionen mit Digest in `harness/mk/postgres-versionen.mk`).

**Verantwortlich:** —
**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Nachtlauf `upstream-drift.yml` (aus `slice-harness-ci-freshness`) prüft zusätzlich die gepinnten Images in `Dockerfile` und `harness/mk/*.mk` (golang, distroless, golangci-lint, a-check, d-migrate, die fünf PostgreSQL-Versionen, das `homebrew/brew`-Image von `harness/mk/homebrew-nachweis.mk`, das `curl`-Image von `tools/harness/traeger-fetch.sh`) auf **zwei Fragen, die verschieden altern**: (1) **Tag-Achse** — gibt es einen neueren Tag (golangci-lint über die Release-Liste, Go über die stabile Version von go.dev); (2) **Digest-Achse** — trägt derselbe Tag heute einen anderen Digest (`docker buildx imagetools inspect`, ohne Bau). Ein Tag-Pin und ein Digest-Pin altern verschieden; ein Pin mit beiden braucht beide Achsen, sonst meldet die Tag-Achse „aktuell“, während derselbe Tag längst neu gebaut ist. Das Skript `tools/harness/pin-freshness.sh` gewinnt die zwei Arten der Frage dazu und liest die Pins aus ihren Dateien; die Ziele und ihre Schritte im Workflow folgen dem Muster des Vorgängers (Exit-Arten, Überspringen bei Netzfehlern, Zählung, `if: always()`).

**Sendung von `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`** (Ausgang `geplant`, eingetragen vom Planner im Commit der Zuweisung, 2026-10-22, `AGENTS.md` §3.13; Entscheidung des Nutzers vom 2026-10-22, „deine empfehlungen“, so vom Koordinator übermittelt, auf den Vorschlag des Planners im Register; der Eintrag steht bei 3×): `harness/mk/homebrew-nachweis.mk` (der Pin `HOMEBREW_IMAGE`, `homebrew/brew` mit Tag und Digest, gesetzt von `slice-erster-release-homebrew-nachweis`) gehört in die Pin-Liste von Liefer-Punkt 2 und wird von der Digest-Achse gelesen wie die übrigen Dateien mit Tag und Digest; die Zeile trägt hinter dem Pin einen Kommentar (Randform 2). Zählung in §8; `grep` nach der Kennung des Eintrags findet die Sendung in §1, DoD und §6. **Nicht angenommen und im Eintrag weiter offen:** die Pins in `tools/harness/release.sh` (`BUILDKIT_IMAGE`, `TRIVY_IMAGE`: nur Digest, Tag im Kommentar, ein Pfad außerhalb von `Dockerfile` und `harness/mk/`), die Werkzeuge des Runners (`skopeo`, `gh`, `buildx`) und der Wortlaut von `gh release download` mit dem echten `gh`; für sie bleibt ein Lauf der Pipeline der einzige Melder.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Skript mit der Frage nach der Release-Liste, die Ziele für Baseline und d-check, der Workflow, die Exit-Arten und der Vertrag — `slice-harness-ci-freshness` (Geber, vor diesem Slice); dieser Slice erweitert sie und ändert ihre Zusagen nicht.
- Das Heben eines Pins (neuer Tag, neuer Digest) — ein Mensch nach dem Befund, ein Slice je Hebung mit dem Gate-Lauf; dieser Slice meldet.
- Pins ohne Tag (`A_CHECK_IMAGE`, `D_MIGRATE_IMAGE`, `TRAEGER_IMAGE` tragen nur `@sha256:…`; ebenso `WORKFLOWS_ACTIONLINT` in `harness/mk/workflows.mk`, mit dem Tag im Kommentar, Stand 2026-10-22, vom Planner nachgetragen): ein Digest ohne Tag hat nichts, wogegen er gefragt werden könnte. Bestand bleibt bewusst stehen, bis der Architect eine Quelle für den Tag wählt (Randform 3); der Vertrag listet sie als nicht geprüft.
- Die Beschaffung von Digests, das Hochladen von Images, Veröffentlichung — nicht Gegenstand der Prüfung; `release.yml` gehört `slice-erster-release-veroeffentlichung`, sein Builder-Image (Treiber `docker-container`, per Digest gepinnt) `slice-erster-release-veroeffentlichung-image` (dort §1, mit der Kennung dieses Slice).
- Die Aktualität der Action-SHA und der Werkzeug-Versionen in Workflows — kein Slice nimmt sie an (`slice-harness-ci-workflows` §6 Frage 3).
- Ob die zwei Stellen eines doppelt gepinnten Images (`Dockerfile` und `harness/mk/build.mk` für golang, `harness/mk/integration.mk` und `harness/mk/postgres-versionen.mk` für PostgreSQL 17) übereinstimmen — ein Konsistenz-Gate ist ein anderer Vorgang; die Prüfung meldet jede Stelle für sich (Randform 4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Tag-Achse):** `make golangci-freshness` und `make go-freshness` (Namen: `SPEC-066`) lesen den Tag aus `Dockerfile`, fragen die Release-Liste von golangci-lint bzw. die stabile Version von go.dev und enden nach dem Vertrag mit aktuell, anders oder übersprungen; die Form des Go-Pins (Tag mit Nebenversion, ohne Patch, `1.27-alpine`) und der Vergleich stehen im Vertrag. Jede Entscheidung je Randform aus §6 steht vor dem ersten Code-Commit in `SPEC-066`.
- [ ] **Liefer-Punkt 2 (Digest-Achse):** Ein Ziel (Name: `SPEC-066`) liest die Pins mit Tag und Digest aus `Dockerfile`, `harness/mk/build.mk`, `harness/mk/integration.mk`, `harness/mk/postgres-versionen.mk`, `harness/mk/homebrew-nachweis.mk` (Sendung von `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`, §1), `d-check.mk`, fragt den heutigen Digest des Tags per `docker buildx imagetools inspect` und meldet je Pin aktuell, abweichend oder übersprungen; die Schritte im Workflow stehen als Schritte von `upstream-drift.yml` mit `if: always()` und bestehen `make workflows-check` und `make docs-check`. Belegt in §7 mit einem grünen Lauf auf GitHub und einem Lauf mit Probe-Pin (ein absichtlich falscher Digest färbt genau seinen Schritt rot).
- [ ] **Liefer-Punkt 3 (Gegenprobe und Index):** `make pin-freshness-images-gegenprobe` (Name: Architect) zeigt mit Attrappen ohne Netz je Zusage von `SPEC-066` für die zwei neuen Arten und die Leser der Pins und je Element ihrer Aufzählungen (die Dateien, die Pin-Formen, die Antworten der Quelle) einen Fall und eine Mutation, die nur ihn färbt (`AGENTS.md` §3.10); `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` je ein Fall; das Ziel steht in `GEGENPROBEN_CONTAINER`. `harness/README.md` §Sensors führt die Zeilen der neuen Ziele und der Gegenprobe und nennt die nicht geprüften Pins.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün; beide Gegenproben des Nachtlaufs (`pin-freshness-gegenprobe` und die neue) grün auf dem Endstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/pin-freshness.sh`, `harness/mk/freshness.mk` | update | zwei Arten der Frage (go.dev, Digest), Leser der Pins aus den Dateien (sechs, mit `harness/mk/homebrew-nachweis.mk`), neue Ziele |
| `.github/workflows/upstream-drift.yml` | update | Schritte der neuen Ziele, je mit `if: always()` |
| Gegenprobe (`harness/mk/…`, `tools/harness/…`; Namen: Architect) | neu | Fälle und Mutationen; Eintrag in `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-066-…md` | update (Architect) | Vertrag; vor dem Code |
| `harness/README.md` §Sensors | update | Zeilen; nicht geprüfte Pins |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci-freshness` liegt in `done/`, dazu der Start-Trigger von [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) (`welle-erster-release` in `done/`), und der Architect hat `SPEC-066` erweitert, jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Schritt 3 der Reihenfolge in §5 der Welle. **Größe gegen die Regel von `AGENTS.md` §3.14** (Entscheidung des Nutzers vom 2026-10-22, wörtlich „Beim Start je Plan entscheiden (Empfohlen)“): Die Mitte der Schätzung liegt bei 108,6 % der Marke 900, über 70 % (Mitte 977,5; der Schnitt in Tag-Achse und Digest-Achse ist vorgemerkt, Frage 1 in §6); es wird nicht rückwirkend geschnitten. Beim Übergang nach `next/`, spätestens vor dem Anspruch, rechnet der Planner die Mitte erneut gegen 70 % der Marke und schneidet nach der Regel oder legt dem Nutzer die datierte Entscheidung „ohne weiteren Schnitt“ vor.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- .github tools harness Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 900; Basis: beim Anspruch (HEAD vor dem Anspruch, hier einzutragen). Schätzung des Planners in dieser Metrik (2026-10-18, ohne Basis): Skript und Fragment +120 bis 180 (Leser für fünf Dateien, zwei Arten), Workflow +20 bis 40, Gegenprobe das Doppelte des Produktcodes (+260 bis 400), `SPEC-066` +80 bis 130, `harness/README.md` 10 bis 20, zusammen 490 bis 770. Schnitt, falls die Marke eintritt: *Tag-Achse* (Liefer-Punkt 1) und *Digest-Achse* (Punkt 2). **Nachgerechnet am 2026-10-22** (Planner, mit der Sendung von `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`; Faktor Gegenprobe zu Produktcode 3 bis 4,6 statt 2, nach den Messwerten in `slice-harness-schaetzung-planner-faktor` §1: 2,9 bei `slice-harness-devcontainer-einrichten`, 3,2 bei `container-start`, 4,6 bei `archiv-bezug`): Skript und Fragment 130 bis 200 (Leser für sechs Dateien), Workflow 20 bis 40, Gegenprobe 390 bis 920, `SPEC-066` 85 bis 140, `harness/README.md` 10 bis 20, **zusammen 635 bis 1320**, das sind 71 bis 147 % der Marke 900 (Mitte 978, 109 %). Die Sendung trägt höchstens 20 Zeilen dazu (ein Leser, ein Fall und eine Mutation). Der Schnitt in Tag-Achse und Digest-Achse ist damit vor dem Anspruch zu entscheiden (Frage 1 in §6), nicht erst bei Eintritt der Marke; die Sendung ginge mit Liefer-Punkt 2 in den Slice der Digest-Achse.
- `in-progress` → `open`: das Register der Images antwortet auf `imagetools inspect` anonym nicht stabil genug (Rate-Limit bei geteilter Adresse, jede Nacht übersprungen), sodass die Achse nichts misst; Entscheidung des Nutzers über Anmeldung oder Streichung (Frage 2).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, ein Lauf von Hand (`workflow_dispatch`) auf GitHub für den Endstand mit allen Achsen grün oder mit benanntem Überspringen und ein Lauf mit Probe-Pin rot als Beleg in §7, `make gates` grün, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ein Digest-Pin des Index und der Digest einer Plattform sind verschieden** (`imagetools inspect` liefert den Digest der Manifestliste; ein Pin kann der einer Plattform sein). Ein Vergleich verschiedener Dinge meldet immer „abweichend“ oder nie. — **Ausgang (vorgesehen):** entfallen, wenn der Vertrag nennt, welcher Digest gemeint ist, die Gegenprobe beide Fälle zeigt und der Probe-Lauf in §7 an einem echten Pin grün endet.
- **Das Floating-Tag `1.27-alpine` trägt die Nebenversion; ein neuer Digest ist dort täglich möglich** (Sicherheits-Patches), die Achse meldet dann fast jede Nacht „abweichend“ und verliert ihr Signal (Dauerrot, vgl. Risiko in `slice-harness-ci-freshness`). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/nachtlauf-ohne-leser` (Beleg bei der Closure), Entscheidung zu Frage 3 vorab.
- **Rate-Limit der Registries für anonyme Abfragen** (Docker Hub für `golang`, `postgres`, `curlimages/curl`; geteilte Adressen der Runner). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/ci-pull-rate-limit` (Beleg bei der Closure, falls dort gezählt; gleiche Beobachtung wie in `slice-harness-ci`).
- **Die Pins ohne Tag bleiben ungeprüft** (`A_CHECK_IMAGE`, `D_MIGRATE_IMAGE`, `TRAEGER_IMAGE`): ein Wächter mit Lücken behauptet mehr, als er hält. — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/pin-ohne-tag-nicht-pruefbar` bei der Closure; die Liste der nicht geprüften Pins steht im Vertrag und in `harness/README.md` (`AGENTS.md` §3.11); dazu gehören die Pins in `tools/harness/release.sh` aus der Sendung in §1.
- **Die Größe liegt an der oberen Hälfte der Spanne** (490 bis 770, Marke 900). — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt; darüber meldet der Implementer (§4 Rückführung *zu groß*).

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem ersten Code-Commit in `SPEC-066` (Technik-Stratum, `harness/vertraege/`). Vorbereitet vom Planner am 2026-10-18 (Empfehlungen, kein Beschluss):

1. **Entscheidung der Vorgänger-ADR gilt.** Wann rot, wann Überspringen, Meldung nur als Rot, fail-open: gelten hier unverändert. Eine neue ADR nur, wenn eine Achse davon abweicht.
2. **Quellen der Pins.** `Dockerfile`: `FROM <image>:<tag>@sha256:…` mit `--platform` und `AS`, drei Stufen; `harness/mk/*.mk`: Variablen mit `?=` und die Liste in `postgres-versionen.mk` (`<Hauptversion>=<image>:<tag>@sha256:…` mit Zeilenfortsetzung); `d-check.mk`: Tag und Digest in zwei Variablen (der Digest sticht den Tag). Fälle: Zeile fehlt, doppelt, Form ohne Tag, ohne Digest, Kommentar mit Muster, Fortsetzungszeile, Leerraum. Mit `harness/mk/homebrew-nachweis.mk` kommt eine Form dazu: `HOMEBREW_IMAGE ?= <image>:<tag>@sha256:<digest> # <Kommentar>` — Variable mit `?=`, Tag und Digest und ein **Kommentar in derselben Zeile**, der den Tag und das Wort „Digest“ trägt; der Fall „Kommentar mit Muster“ gilt für diese Zeile (Sendung von `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`, §1). Empfehlung: Leser pro Datei mit eigenem Namen im Vertrag; eine nicht lesbare Form ist Exit 2.
3. **Pins ohne Tag.** Optionen: (a) ungeprüft und im Vertrag aufgeführt; (b) ein Tag in einem Kommentar oder einer Variable neben dem Pin (`# pin-tag: v…`), den der Leser liest — Änderung der Dateien; (c) Frage an `latest`. Empfehlung: (a), (b) später als eigener Slice; (c) abgelehnt, weil es ein anderes Image meint.
4. **Doppelte Pins.** `golang:1.27-alpine@…` steht in `Dockerfile` und `harness/mk/build.mk`; PostgreSQL 17 in `integration.mk` und `postgres-versionen.mk`. Empfehlung: jede Stelle einzeln geprüft und genannt, kein Abgleich der Stellen (anderer Vorgang, Abgrenzung in §1).
5. **Digest des Index oder der Plattform.** Empfehlung: der Digest der Manifestliste, und der Vertrag sagt, dass ein Pin auf eine Plattform als abweichend gilt; ein Fall je Form.
6. **Go-Version.** Der Pin trägt `1.27-alpine` (Nebenversion). go.dev liefert `go1.27.3`; der Vergleich der Nebenversion (`1.27` gegen `1.27`) meldet erst bei `1.28`. Empfehlung: Nebenversion gegen Nebenversion; Patch-Stände deckt die Digest-Achse.
7. **Vorabversionen und Alias-Tags** (`latest`, `nonroot`): `nonroot` ist ein Tag auf demselben Image und folgt der Digest-Achse; eine Vorabversion als Pin ist ein Fall.
8. **Netz, Rate-Limit, Zeit.** wie `slice-harness-ci-freshness` (Überspringen, Zählung); Abfragen nacheinander mit Obergrenze der Zeit je Abfrage, nicht parallel; zusammen höchstens ein Dutzend Abfragen je Nacht (neun Images, fünf PostgreSQL-Versionen, ein `curl`), der Vertrag nennt die Obergrenze der Gesamtzeit. Mit `harness/mk/homebrew-nachweis.mk` sind es bis zu zwölf Digest-Abfragen (`Dockerfile` 3, `build.mk` 1, `integration.mk` 1, `postgres-versionen.mk` 5, `d-check.mk` 1, `homebrew-nachweis.mk` 1) und zwei der Tag-Achse; der Architect nennt die Obergrenze neu, ein Dutzend trägt die Sendung nicht mehr (Stand 2026-10-22).
9. **Aufrufarten des Treibers** (`AGENTS.md` §3.12): `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`: wie im Vorgänger; die Ziele schreiben nichts.
10. **Docker im Nachtlauf.** `docker buildx imagetools inspect` ist auf dem Runner vorhanden und baut nichts; lokal braucht das Ziel Docker mit `buildx` (Voraussetzung am Host des Aufrufers, im Vertrag genannt); fehlt `buildx`, ist das ein Aufruffehler (Exit 2).
11. **Reihenfolge im Workflow.** Die Schritte der Achsen stehen in einer festen Reihenfolge, jeder mit `if: always()`; ein Fall je Schritt, dass er `if: always()` trägt.

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Schnitt.** Beauftragt war ein Slice `slice-harness-ci-freshness`; ungeteilt lag die Schätzung bei 900 bis 1400 (Marke 900). Empfehlung: zwei Slices wie angelegt (Baseline und d-check zuerst, dann die Images). Alternative: ein Slice, dann Rückführung nach `next` bei Überschreitung.
2. **Anmeldung bei Docker Hub** im Nachtlauf, falls das Rate-Limit die Achse dauerhaft überspringen lässt. Empfehlung: nein (ein Geheimnis mehr, ein Recht mehr); die Achse meldet das Überspringen, und die Beobachtung entscheidet.
3. **Floating-Tag `1.27-alpine`**: tägliche Digest-Änderungen als Rot oder als Hinweis (Exit 0 mit Zeile)? Empfehlung: Rot nur bei Tag-Achse; bei der Digest-Achse einer Floating-Tag-Zeile ein Hinweis mit Zählung, und der Architect legt die Liste dieser Tags im Vertrag fest. Gegenposition: alles rot, das Signal bleibt bei jedem Digest-Wechsel.
4. **Pins ohne Tag** (Randform 3): Empfehlung (a), ungeprüft und genannt.

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
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12) — die zwei Arten der Frage sind neue Zusagen; `BEO-REPO/postgres-versionsliste-altert` — die Digest-Achse ist ihr Sensor für den Digest (die Liste der Versionen selbst bleibt ein Mensch); `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` — nicht berührt. Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt keine Sendung eines Gebers; er erweitert das Skript, den Workflow und den Vertrag von `slice-harness-ci-freshness` (Geber der Erweiterung, dort §1 genannt). Liefer-Punkte 3 (Tag-Achse | Digest-Achse | Gegenprobe samt Index). Schichten 2: Workflows (`upstream-drift.yml`, zwei Schritte) und Harness; Planung, Vertrag und ADR sind keine Schicht; Wartungs-Doku und Handbuch bleiben unberührt. Schätzung 490 bis 770 (§4), unter der Marke.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Welle wechselt; Entscheidung des Nutzers vom 2026-10-18: der erste Release hat Vorrang vor den Nachtläufen): Der Slice wechselt von welle-erster-release in welle-ci-nachtlaeufe und startet nach `welle-erster-release` in `done/`; Inhalt, Sendungen, Liefer-Punkte und Schichten bleiben unverändert (Liefer-Punkte und Schichten wie in der Zählung oben, Marke in §4 unverändert).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der den Ausgang `geplant` von `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder` einträgt): Die Sendung steht in §1, in Liefer-Punkt 2 der DoD, in §3, in den Randformen 2 und 8 und im Risiko zu den Pins ohne Tag als je ein Satz; die Datei `harness/mk/homebrew-nachweis.mk` ist ein Eintrag in der Liste von Liefer-Punkt 2. Liefer-Punkte bleiben 3 (Tag-Achse | Digest-Achse | Gegenprobe samt Index), Schichten bleiben 2 (Workflows | Harness; die Datei liegt unter `harness/`). Schätzung in der Metrik von §4 mit dem Faktor 3 bis 4,6 für Gegenprobe-Code: **635 bis 1320** statt 490 bis 770; das obere Ende liegt über der Marke 900, die Mitte bei 109 %. Der Plan liegt damit über 70 % der Marke; vorgemerkt ist der Schnitt in Tag-Achse und Digest-Achse (§4), die Entscheidung trifft der Planner vor dem Anspruch (Frage 1 in §6), und die Sendung geht dann mit Liefer-Punkt 2 an den Digest-Slice. Die Sendung selbst trägt höchstens 20 Zeilen.
