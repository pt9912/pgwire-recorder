# Welle welle-v1-nachtrag: v1-Nachtrag: Homebrew und CancelRequest ohne Schlüssel

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M5.

**Verantwortlich:** pt9912. **Datum:** 2026-10-11.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die zwei Anforderungsreste von v1, die hinter dem ersten Release liegen, sind umgesetzt: Die Homebrew-Formel liegt im Tap, und Abnahmeszenario 11 ([`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)) ist auf macOS und Linux nachgewiesen; ein `CancelRequest` ohne Abbruchschlüssel wird wie jeder `CancelRequest` behandelt ([`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)). Damit sind alle Anforderungen (MUSS und SOLL) umgesetzt. Dazu kommt ein Nachtrag an Tests und Testdaten ohne Produktcode: Das Ende einer Extended-Interaktion durch den Server ist gegen PostgreSQL belegt (`slice-v1-abschluss-verbindungsende-upstream-extended`; Entscheidung des Nutzers vom 2026-10-18, vorher wellenlos hinter dem Release; er liefert keine Anforderung und ist keine Bedingung von M5). Die Slices standen in `welle-v1-abschluss` (Homebrew, `CancelRequest`) und `welle-erster-release` (Nachweis); der Nutzer hat sie am 2026-10-11 hinter den ersten Release gezogen, damit M4 ohne sie erreicht wird.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-erster-release](done/welle-erster-release/welle-erster-release.md) done (Entscheidung des Nutzers vom 2026-10-11).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Abnahmeszenario 11 (Homebrew) ist auf macOS und Linux mit dem echten Tap nachgewiesen — das *Mehr* gegenüber den Slice-DoDs.
- `make gates` grün.
- Closure-Notiz in `welle-v1-nachtrag-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-v1-abschluss-verbindungsende-upstream-extended | Das Ende einer Extended-Interaktion gegen PostgreSQL belegen (Fixtures der Gruppen mit Sync und mit Flush) | [`LH-FA-25`](../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus) |
| slice-v1-abschluss-cancel-ohne-schluessel | CancelRequest ohne Schlüssel | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-v1-abschluss-binary-namen | Release-Binaries ohne Version im Dateinamen | [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |
| slice-v1-abschluss-homebrew | Homebrew-Bereitstellung | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |
| slice-v1-abschluss-release-v0-1-1 | Release v0.1.1: Freigabe, Tag (Handlung des Nutzers), Bezug im Handbuch | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung) |
| slice-erster-release-homebrew-nachweis | Homebrew-Nachweis (Abnahmeszenario 11) | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: Welle [welle-erster-release](done/welle-erster-release/welle-erster-release.md). Der letzte Slice der Welle (`slice-erster-release-homebrew-nachweis`) setzt außerdem einen veröffentlichten stabilen Tag mit dem Schema ohne Version voraus (`v0.1.1`; Handlung des Nutzers nach neuem Ja zum Tag; `slice-v1-abschluss-release-v0-1-1` bereitet ihn vor und belegt ihn); der Closure-Trigger der Welle hängt damit an diesem Tag.
- Innerhalb der Welle: `slice-v1-abschluss-verbindungsende-upstream-extended` setzt `slice-v1-abschluss-verbindungsende-upstream`, `-record`, `-replay-extended` und `-replay-extended-flush` aus [welle-v1-abschluss](done/welle-v1-abschluss/welle-v1-abschluss.md) voraus (alle vor dem Release in `done/`) und hängt an keinem Slice dieser Welle; `slice-v1-abschluss-cancel-ohne-schluessel` setzt `slice-v1-abschluss-protokollrand` aus [welle-v1-abschluss](done/welle-v1-abschluss/welle-v1-abschluss.md) voraus (der übrige Protokollrand); `slice-v1-abschluss-binary-namen` setzt `slice-v1-abschluss-container` aus [welle-v1-abschluss](done/welle-v1-abschluss/welle-v1-abschluss.md) voraus (`make binaries`, `make test-container`) und `slice-erster-release-veroeffentlichung` aus [welle-erster-release](done/welle-erster-release/welle-erster-release.md) (`release.sh`); `slice-v1-abschluss-homebrew` setzt `slice-v1-abschluss-container` (Release-Artefakte) und `slice-v1-abschluss-binary-namen` (die Namen, auf die Vorlage und Fälle rechnen) voraus; `slice-v1-abschluss-release-v0-1-1` setzt `slice-v1-abschluss-binary-namen` (die Namen im getaggten Commit) und `slice-v1-abschluss-homebrew` (der Job in `release.yml`) voraus; `slice-erster-release-homebrew-nachweis` setzt `slice-v1-abschluss-homebrew` (das Verfahren der Formel), `slice-v1-abschluss-binary-namen` und `slice-v1-abschluss-release-v0-1-1` mit dem von ihm vorbereiteten, vom Nutzer gesetzten und veröffentlichten stabilen Tag voraus (nicht mehr `v0.1.0` und den Nachtrag von Hand, Entscheidung des Nutzers vom 2026-10-20).

**Reihenfolge** (Entscheidung des Nutzers vom 2026-10-11, Schritt 0 ergänzt am 2026-10-18; Schritt 4 am 2026-10-20 vom Planner eingefügt, die Stellung ist Empfehlung, den Slice wollte der Nutzer; WIP-Limit 1):

0. `slice-v1-abschluss-verbindungsende-upstream-extended` — Fixtures der Extended-Interaktion mit Ende gegen PostgreSQL; er stand seit dem 2026-10-16 vor dieser Welle, der Platz am Anfang bleibt (Empfehlung des Planners; keine Voraussetzung eines anderen Slice der Welle).
1. `slice-v1-abschluss-cancel-ohne-schluessel` — `CancelRequest` ohne Schlüssel in `record` und `replay`.
2. `slice-v1-abschluss-binary-namen` — Binaries als `pgwire-recorder-<os>-<arch>` ohne Version im Namen (Entscheidung des Nutzers vom 2026-10-20, [ADR-0054](../adr/0054-binaries-ohne-version-im-dateinamen.md)); vor der Formel, weil Vorlage und Fälle die Namen brauchen.
3. `slice-v1-abschluss-homebrew` — Verfahren der Formel, belegt mit einem lokalen Probe-Tap (ein Bare-Repository; ein Probe-Tag entfällt, Entscheidung des Nutzers vom 2026-10-20, [ADR-0053](../adr/0053-homebrew-tap-in-der-release-pipeline.md)); setzt 2 voraus.
4. `slice-v1-abschluss-release-v0-1-1` — `VERSION`, Eintrag der Freigabe und Bezug im Handbuch für den stabilen Tag `v0.1.1`, den der Nutzer nach seinem Ja setzt; setzt 2 und 3 voraus.
5. `slice-erster-release-homebrew-nachweis` — Abnahmeszenario 11 mit dem echten Tap (setzt 3 und 4 voraus, also den veröffentlichten Tag mit dem neuen Schema).

Vor dieser Welle liegen die Gate-Kette (wellenlos, Entscheidung des Nutzers vom 2026-10-13), `slice-harness-integration-volumes`, `slice-harness-abdeckung-gate` und `slice-harness-coverage` (M3) und [welle-erster-release](done/welle-erster-release/welle-erster-release.md) (M4). Danach folgen die übrigen Harness-Slices; die Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; wellenlos) liegt vor allen Release-Slices (Entscheidung des Nutzers vom 2026-10-13), beginnend mit `slice-harness-commit-struktur-id`; ihre Reihenfolge steht in deren §4.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Die Aufnahme in das Standard-Repository von Homebrew — Out-of-Scope von LH-FA-19.
- Weitere Plattformen oder Registries — nicht zugesichert.
- Der erste Release und seine Pipeline — `welle-erster-release`; diese Welle ändert `v0.1.0` nicht (er behält die Namen mit Version und bekommt keine Formel).
- Das Setzen des Tags `v0.1.1` (der erste mit dem Schema ohne Version) — Handlung des Nutzers nach neuem Ja; `slice-v1-abschluss-release-v0-1-1` bereitet ihn vor und setzt ihn nicht, der letzte Slice der Welle setzt ihn voraus.
- Übriger Protokollrand — `slice-v1-abschluss-protokollrand` in `welle-v1-abschluss`.
- Produktcode zum Verbindungsende — `slice-v1-abschluss-verbindungsende-upstream-extended` liefert nur Tests und Testdaten; findet er eine Änderung am Produkt, geht sie an einen Folge-Slice, die Welle ändert sie nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-v1-nachtrag-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
