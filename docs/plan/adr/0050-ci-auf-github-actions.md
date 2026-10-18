# ADR-0050: Die CI läuft auf GitHub Actions, gehärtet durch Pins, leere Rechte und Zeitgrenzen, mit einem Gate für die Workflows; sie ergänzt den lokalen Hook

**Status:** Accepted

**Datum:** 2026-10-18

**Autor:** Architect; angenommen vom Nutzer am 2026-10-18.

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [`LH-QA-01`](../../../spec/lastenheft.md#lh-qa-01--determinismus), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md), [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md)

**Schärft:** [`SPEC-063`](../../../harness/vertraege/SPEC-063-workflows-check.md), [`SPEC-064`](../../../harness/vertraege/SPEC-064-ci.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Prüfung dieses Repos läuft heute an einer Stelle: im Klon des Entwicklers, vor dem Ende einer Sitzung, durch den Stop-Hook und `make gates` ([ADR-0038](0038-gate-kette-voll-und-schnell.md)). Der Hook gilt für jeden Lauf von Claude Code in diesem Klon, nicht für `git commit --no-verify`, nicht für einen frischen Klon und nicht für einen Beitrag von außen. Das Repo ist öffentlich; der erste Release steht an (Auftrag des Nutzers vom 2026-10-18: „Erst das Release voranbringen“), und eine Veröffentlichung, deren Gates nur im Klon des Autors liefen, belegt nichts für den, der sie nutzt. Mit [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md) dauert `make gates` etwa vier bis fünf Minuten und fährt keine Gegenprobe; das ist die Größe, die ein Lauf auf einem fremden Runner tragen muss.

Zwei Befunde der Prüfung vor dem Code prägen die Entscheidung. Erstens: `make doc-commits` und `make doc-immutable` melden in diesem Repo heute „0 Befund(e)“ über jeder Range, weil `.d-check.yml` die Blöcke `commits` und `vcs` nicht führt; der Vorlauf fängt nur eine leere Range. Eine CI, die beide ruft, wäre grün, ohne etwas zu prüfen. Zweitens: Eine Workflow-Datei ist Code mit Rechten (Token, Geheimnisse, fremde Aktionen), den bisher kein Gate dieses Repos liest.

## Entscheidung

Wir empfehlen **Option A: die CI läuft auf GitHub Actions; `ci.yml` fährt `make gates` unverändert und prüft die Commit-Range; jede Workflow-Datei besteht ein Gate vor ihrem ersten Commit; die CI ergänzt den lokalen Hook und ersetzt ihn nicht.**

1. **Träger.** GitHub Actions auf `ubuntu-24.04` (fest, nicht `ubuntu-latest`), `linux/amd64`. Der Remote ist GitHub, das Image wird dort veröffentlicht, und ein anderer Träger bräuchte einen Runner, den dieses Repo nicht betreibt.
2. **Was die CI prüft.** `make gates` ohne Änderung (kein zweites Gate-Ziel nur für die CI) und die Commit-Range: jede Nicht-Merge-Message trägt eine Kennung (`make doc-commits`), eine angenommene ADR bleibt unverändert (`make doc-immutable`). Dafür führt `.d-check.yml` die Blöcke `commits` und `vcs`, und eine Gegenprobe zeigt, dass sie wirken. Die Range berechnet ein Skript aus dem Ereignis; jede Form, die nichts prüft (leere Range, flacher Klon, unauflösbare Basis), ist rot.
3. **Härtung.** Aktionen per vollem Commit-SHA mit Tag-Kommentar, `permissions: {}` oben und im Job nur, was er braucht, eine Zeitgrenze je Job, keine Geheimnisse, keine Aktion außer dem Checkout, kein Cache und kein Hochladen. Ein Gate (`make workflows-check`) prüft Syntax und Ausdrücke mit `actionlint` (Image per Digest), die Rechte und Zeitgrenzen, und dass `make docs-check` die Pins prüft (Modul `workflows` von d-check). Es steht in der Gate-Kette und liegt vor der ersten Datei.
4. **Verhältnis zum Hook.** Die CI ergänzt: Sie läuft klonunabhängig und über die Range. Sie schreibt keinen Nachweis, den der Hook liest, und der Hook liest nichts aus ihr.
5. **Was keine Zusage des Repos ist.** Branch-Schutz („required status check“), Geheimnisse, Sichtbarkeit und die Freigabe von Läufen für Erstbeitragende sind Einstellungen auf GitHub und aus dem Klon nicht prüfbar. Der Name des Jobs, den ein Branch-Schutz verlangt, steht im Vertrag.

Die Regeln im Einzelnen (Text der Datei, Formen der Range, Konfiguration, Fälle, Grenzen) stehen in [`SPEC-064`](../../../harness/vertraege/SPEC-064-ci.md) und [`SPEC-063`](../../../harness/vertraege/SPEC-063-workflows-check.md).

Der Nutzer hat am 2026-10-18 die ADR angenommen, `ubuntu-24.04` statt `ubuntu-latest` (abweichend vom ursprünglichen Auftrag) und den Verzicht auf `concurrency` und `workflow_dispatch` bestätigt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — GitHub Actions, gehärtet, mit Gate für die Workflows** | prüft jede Integration klonunabhängig; die Härtung ist mechanisch gehalten; keine Annahme über Geheimnisse | ein Lauf pro Integration kostet vier bis fünf Minuten und hängt an der Verfügbarkeit von GitHub und der Docker-Registries |
| B — keine CI, nur der Hook | nichts zu bauen | `--no-verify`, ein frischer Klon und Beiträge von außen bleiben ungeprüft; ein Release ohne fremde Prüfung |
| C — ein anderer Träger (selbst betriebener Runner, anderer Dienst) | Kontrolle über die Umgebung | Das Repo liegt auf GitHub, der erste Release geht in dessen Registry; ein eigener Runner ist Betrieb, den niemand trägt |
| D — GitHub Actions ohne Gate, Härtung durch Review | weniger Werkzeug | Was nur das Review hält, hält es bis zum ersten Versehen; eine Datei mit Rechten ist die Stelle, an der ein Versehen die Lieferkette trifft |
| E — die CI ersetzt den Hook | ein Ort der Prüfung | Der Hook schützt die Arbeit im Klon vor dem Commit; die CI sieht erst, was geschoben ist |
| F — `ubuntu-latest` und `concurrency` mit Abbruch (Vorbild d-check) | gängige Form | Ein Wechsel des Images ändert still die Gates; ein Abbruch verwirft die Range des abgebrochenen Push |

## Konsequenzen

- Positiv: Jede Integration läuft durch die Gates auf einer Umgebung, die keiner der Beteiligten pflegt; die Range-Prüfungen, die bisher nichts prüften, prüfen; eine Workflow-Datei besteht ein Gate, bevor sie Rechte bekommt.
- Negativ: Ein Rate-Limit anonymer Pulls vom Docker Hub kann einen Lauf rot färben und bleibt rot (keine Anmeldung, keine Geheimnisse); die Laufzeit auf dem Runner ist unbekannt, bis ein Lauf sie zeigt; ein Pull-Request im selben Repo löst zwei Läufe aus.
- Akzeptiertes Negativ: Die CI ist ohne Branch-Schutz ein Hinweis, kein Sperrriegel. Grund: die Einstellung liegt auf GitHub, nicht im Klon; sie ist eine Handlung des Betreibers.
- Akzeptiertes Negativ: Das Gate der Workflows prüft die Form, nicht die Lauffähigkeit und nicht die Gültigkeit der Pins. Grund: nur ein Lauf auf GitHub und das Netz messen das.
- Akzeptiertes Negativ: Die CI prüft die Anwesenheit einer Kennung in der Message, nicht die Existenz der genannten Plandatei. Grund: der Hook prüft am Index, die Range liest Commits; eine Prüfung am Index der Spitze lehnte alte Commits zu Unrecht ab.
- Folgepflicht: [`SPEC-063`](../../../harness/vertraege/SPEC-063-workflows-check.md), [`SPEC-064`](../../../harness/vertraege/SPEC-064-ci.md); `harness/README.md` §Sensors führt Gate und Gegenproben und nennt die CI als Ergänzung des Hooks und den Branch-Schutz als Handlung des Betreibers; die Workflows der Nachtläufe und der Veröffentlichung fallen vom ersten Commit an unter das Gate.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, Docker, `actionlint`, d-check | Syntax, Rechte, Zeitgrenzen und Bindung der Pin-Prüfung jeder Workflow-Datei | `make workflows-check` |
| bash, make, Docker | Gate der Workflows an Attrappe und Image: Fälle und Mutationen je Punkt | `make workflows-check-gegenprobe` |
| bash, git, Docker, d-check | Text von `ci.yml`, Berechnung der Range, Wirkung der Konfiguration der Range-Prüfungen | `make ci-gegenprobe` |

Dass ein Lauf auf GitHub grün endet, prüft kein Gate; der Beleg ist ein Lauf.

## Re-Evaluierungs-Trigger

Wenn ein Lauf auf GitHub rot endet, den ein Gate dieses Repos hätte fangen sollen; wenn der Runner `make gates` nicht in der Zeitgrenze trägt; wenn das Limit des Docker Hub Läufe regelmäßig färbt (dann Anmeldung oder Spiegel erwägen, mit Geheimnis und neuer Entscheidung); wenn der Nutzer einen anderen Träger will.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Proposed (Empfehlung des Architect nach dem Auftrag des Nutzers, die CI vor die Veröffentlichung zu legen) | — |
| 2026-10-18 | Accepted (Nutzer) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
