# Verifikation slice-harness-devcontainer-einrichten-start

**Rolle:** Verifier. **Stand:** Zweig `rev-start`, Commit `1f93ed00`, Basis `aad621a8`. **Datum:** 2026-10-18.
**Gegenstand:** DoD (§2) des Slice-Plans `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten-start.md` gegen [`SPEC-058`](../../harness/vertraege/SPEC-058-container-start.md) Teil II (Punkte 9 bis 12, 13, 14, 15), [ADR-0045](../plan/adr/0045-container-start-als-werkzeug-im-repo.md) und [ADR-0047](../plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md). Nach Auftrag des Nutzers wiederholt dieser Lauf weder die volle Gegenprobe (246 Mutationen) noch `make gates` und fährt je den engsten Sensor.

## Urteil

Die Liefer-Punkte 1 (Ziel) und 2 (Gegenprobe) sind **erfüllt, soweit meine Läufe reichen**: alle gefahrenen Mutanten färben genau die im Plan genannten Fälle, der echte Lauf gegen Docker und `devcontainer` verhält sich wie der Vertrag. Liefer-Punkt 3 (Dokument) ist bis auf zwei kleine Überdehnungen belegt. **Der Plan selbst ist nicht in Ordnung (V-1, V-2):** Der Commit `1f93ed00` hat §5, §6 und §7 des Plans gelöscht, und die behaupteten Belege des Implementers stehen nicht in der Datei. Closure-Pflichten sind erwartungsgemäß offen.

## Eigene Läufe (Host, Daemon des Hosts; `pgr-dev`, `pgr-dev2`, `pgr-dev3` und `/Development/pgr-dev*` unberührt)

| # | Lauf | Ergebnis |
|---|---|---|
| 1 | Gegenprobe ohne Docker-Gruppe | **Schalter existiert nicht** (`GP_KIND` ist der interne Kind-Modus, `GP_NUR` wählt Fälle nur positiv). Ersatz: unveränderte Kopie mit allen Fällen (Attrappen und Gruppe Docker) zusammen mit Lauf 2: grün |
| 2 | 13 Mutanten einzeln über `GP_MUT` und `GP_ENTDECKE=1`, ein Aufruf | Exit 0 `gruen`, 5 min 27 s einschließlich unveränderter Kopie. Rot genau wie in §4 des Plans: `p9-kollision` (`p1-keine-zerstoerung-kollision`, `p13-kollision-vor-port`, `p9-kollision`), `p9-laeuft-immer` (`p10-belegt-gestoppt`, `p9-created`, `p9-exited`), `p9-update` (`p9-kein-update`), `p10-ohne-pruefung` (`p10-belegt`, `p10-belegt-gestoppt`), `p11-vorgabe-warten` (`p11-vorgabe-warten`), `p12-int-entfaellt` (drei `p12-sigINT-*`), `p12-rueckbau-docker` (alle sechs `p12-sig*`), `p15-container-per-praefix`, `p15-volume-per-praefix`, `p15-image-per-praefix` (je `p15-fremde-docker-objekte`), `p9-docker-neustart-ersetzt` (`p9-created`, `p9-docker-stop-start`, `p9-exited`), `p10-docker-ohne-pruefung` (`p10-docker-port-belegt`), `p11-docker-erfolg` (`p11-docker-anderer-schluessel`). Je Punkt 9, 10, 11, 12, 15 (g) mindestens einer. Nach dem Lauf kein Objekt `pgr-cs-*` (Container, Volumes, Images) und kein `/tmp/pgr-cs-*` |
| 3 | `make container-start` echt, Scratch (`PGR_NAME=pgr-cs-ver1`, Port 2310, `PGR_CPUS=2`, `PGR_MEM=1g`, `PGR_DEV` unter dem Scratchpad) | Erststart Exit 0, 45,7 s, letzte Zeile `container-start: bereit pgr-cs-ver1 127.0.0.1:2310`. Zweiter Aufruf bei laufendem Container 0,3 s, `StartedAt` unverändert. Nach `docker stop` (Status `exited`) dritter Aufruf: `docker start` derselbe Container (gleiche ID), 2,06 s, dieselben zwei Volumes. Vierter Aufruf mit `PGR_MEM=2g` bei laufendem Container: `Memory` blieb 1073741824 (Ziel ändert Limits nie). `known_hosts` des Hosts unverändert (Zeitstempel, kein Eintrag für `[127.0.0.1]:2310`) |
| 3a | Messung der Limits | `docker inspect`: `Memory` 1073741824, `NanoCpus` 2000000000, `OOMKilled=false`; Label `devcontainer.local_folder` gleich `$PGR_DEV` wörtlich; Mounts: zwei benannte Volumes |
| 3b | Entfernen nur per Name | `docker rm -fv pgr-cs-ver1`; die zwei benannten Volumes blieben (erwartet), per Name mit `docker volume rm` entfernt, kein anonymes Volume übrig; `docker rmi` per Name: das Image `…-uid` ging, das Eltern-Image `vsc-pgr-cs-ver1-<Hash>:latest` blieb und ging mit einem zweiten `docker rmi` per Name (V-4). Danach nichts mit `pgr-cs` in `docker ps -a`, `volume ls`, `images` |
| 4 | `make docs-check` | Exit 0, 1080 Dateien, 0 Befunde (vor Anlage dieses Berichts) |
| 4 | `make kopf-check` | Exit 0 |
| 4 | `make gates-schnell` | Exit 0, 10,8 s (warmer Cache) |
| 5 | Metrik `git diff --numstat aad621a8..HEAD -- tools harness docs/developer` | 515 Einfügungen, 127 Löschungen, **Summe 642**; `harness/vertraege/` 25 und 4, zusammen 29; ohne sie **613**. Beides wie im Plan. Unter Meldeschwelle 700 und Marke 900 |

Nicht gefahren (Auftrag): volle Gegenprobe (246 Mutationen, 19 min 36 s laut Implementer) und `make gates` (26 min 7 s laut Implementer); beide Zahlen sind **Behauptung ohne Beleg in der Plan-Datei** (V-2).

## DoD-Punkte

- **Liefer-Punkt 1, Ziel Teil II:** erfüllt, gefahren: Zustand (Container, Port), `start`, `warten`, Meldung `bereit`, zweiter Aufruf lässt den laufenden Container stehen (Lauf 3), Neustart nach `stop` derselbe Container, Limits ändert das Ziel nie (3, `p9-update` rot), Kollision und Port belegt Exit 4 (Mutanten). `harness/README.md` nennt das volle Verhalten (Zeilen `container-start` und `container-start-gegenprobe`). Das README nennt für die Gegenprobe „je Punkt 1 bis 14“, die Zeile führt zugleich Punkt 15 und 16; Wortwahl, kein Befund.
- **Liefer-Punkt 2, Gegenprobe:** erfüllt für die gefahrenen Mutanten; `p11-vorgabe-warten` (Vorgabe 61 liefert 62 Aufrufe) rot nur im Fall `p11-vorgabe-warten`; Gruppe Docker entfernt nur Selbstangelegtes (Lauf 2 und `p15-fremde-docker-objekte`, `p15-fremdes-verzeichnis` nicht einzeln gefahren). **Laufzeit der Gegenprobe gemessen und im Bericht genannt (DoD):** nicht in der Plan-Datei (V-2); von mir gemessen nur der Teilumfang in Lauf 2.
- **Liefer-Punkt 3, Dokument:** §2 bis §5 nennen `make container-start`, Abschnitt 8.1 „Zwei Arten zu arbeiten“ liegt vor. Aussagen gegen Messung: siehe unten.
- **`make gates` grün:** Beleg des Implementers (Exit 0, 26 min 7 s) liegt **nicht in der Datei** (V-2); mein Teilnachweis: `gates-schnell` grün. Stop-Hook verlangt den vollen Lauf auf dem Endstand.
- **Review:** ein Report liegt untracked im Worktree (`2026-10-18-review-slice-harness-devcontainer-einrichten-start.md`), von mir nicht gelesen.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** offen, und die Grundlage dafür fehlt (V-1).

## Dokument gegen Messung (§8.1, §10.1, §12 und Umgebung)

- **`12g`** kommt in §10 (Zeile „damalige Vorgabe `12g`“), §10.1 (Spitze 5,42 und 4,24 GiB „bei Limit `12g`“) und §11 („Vorgabe zur Zeit der Messung, jetzt `8g`“) vor und ist jedes Mal als Messung bei 12g gekennzeichnet; der Lauf bei `8g` steht in §13 als nicht gemessen. In Ordnung.
- **Limits** (§10.1, §4): „`make container-start` ändert sie nie“ ist durch Lauf 3 (Limit blieb bei abweichendem `PGR_MEM`) und `p9-kein-update` belegt. `PGR_CPUS=2`, `PGR_MEM=2g` mit `2000000000` und `2147483648` (§10) trägt die Gruppe Docker (`p15-docker-grenzen`); meine Werte bei `1g` bestätigen den Mechanismus.
- **Neustart** (§4, Tabelle): „derselbe Container“ bestätigt (gleiche ID, gleiche Volumes); die 1,7 bis 1,8 s stammen aus der Messung des Vorgängers bis `ssh` und `dockerd`, mein Aufruf von `make container-start` brauchte 2,06 s einschließlich `devcontainer up` und Probe. Kein Widerspruch, andere Strecke.
- **§12, Absturz des Daemons:** Messung des Koordinators, „nicht `docker start`“ ausdrücklich als Folgerung gekennzeichnet; nicht von mir wiederholt.
- **Keine Namen oder Mail des Nutzers, kein Zugangswert:** `grep` über Dokument, Skripte und Vertrag ohne Treffer; Beispiele nutzen Platzhalter (`<Schlüssel>`, `<Token>`).
- **Quantor Zeile 75:** „Jeder Aufruf lässt sich wiederholen.“ steht ohne Einschränkung; der Plan selbst sagt, dass nur `p12-wiederholung` und die Docker-Folge für `geruest`, `start` und das Starten nach `stop`/`rm` es tragen, und der Vertrag Punkt 12 nennt zwei akzeptierte Negative (Abbruch von `u-boot init` hinterlässt ein „fremdes“ Verzeichnis, Abbruch zwischen den zwei Dateien von `ssh-keygen`). Beide fehlen im Dokument und in §13 (V-3).
- **„höchstens `PGR_WARTEN` Sekunden“ (§4):** Das Skript prüft die Zeit nach der Probe; eine Probe dauert bis `ConnectTimeout=5`, die Wartezeit kann die Grenze also um bis zu rund 5 s überschreiten. Der Vertrag und der Fall `p11-nie` sagen „nicht unter“ zu. „Höchstens“ ist durch keinen Test gedeckt (V-3, aus dem Code gelesen, nicht gemessen).
- **§8.1 (a) „dieses Dokument entstand so“:** Prozessaussage, nicht prüfbar; harmlos, aber nicht belegt.
- **§4 „Image löschen“:** nennt ein Image; es entstehen zwei (V-4).

## Plan gegen Code

- §3 nennt `pgr-dev2` und Port 2223 für die Arbeit; der Lauf war `pgr-dev3` und 2224. Der Plan trägt das selbst als Rückgabe an den Planner nach. Offen im Sinne von §3.9 bleibt, §3 nachzuziehen.
- §3 erwartet `harness/mk/container-start*.mk` unverändert; geändert sind drei Hilfetext- und Kommentarzeilen (3 und 1 Zeile). Folgenlos, aber nicht die erwartete Spalte.
- §3 nennt Schritte, Falle, Gegenprobe, Index, Dokument; alle geliefert. Die Reihenfolge der Gründe im Skript (Kollision, anderer Name, Port) steht als Rückgabe 1 beim Architect, Punkt 13 des Vertrags schweigt dazu.
- Vertrag: 29 Zeilen in `harness/vertraege/` aus `arch-limits` und `arch-spec058-punkt15`, im Plan als Rückgabe 3 genannt.

## Befunde

- **V-1 (hoch, Plan): §5, §6 und §7 des Slice-Plans sind gelöscht.** Der Commit `1f93ed00` („§7 Belege …“) entfernt §5 Closure-Trigger, §6 Risiken und offene Punkte (die Randformen von Teil II samt Fundstelle und alle „Weiteren Risiken“ mit ihren Ausgängen), §7 Closure-Notiz sowie die vorab benannten Rückführungs-Bedingungen in §4. Folgen: Die DoD verlangt weiter „je Randform aus §6 ein Fall“ und „Jedes Risiko aus §6 trägt einen Ausgang“; beides hat keine Grundlage mehr. `AGENTS.md` §3.12 verlangt die Randformen in §6, das Baseline-Regelwerk (`modul-05-planning-harness.md`) die Rückführungs-Bedingung vorab in §4 und einen Ausgang je Risiko bei der Closure. Der Planner muss §5 bis §7 aus `1f93ed00~1` wiederherstellen und die Belege ergänzen, nicht ersetzen. Die Datei kennt auch keinen Abschnitt „§7 Belege“, obwohl die Commit-Nachricht ihn nennt.
- **V-2 (hoch, Belege): Die Belege, auf die der Auftrag und die Commit-Nachricht verweisen, stehen nicht in der Plan-Datei.** Nicht auffindbar: Messung des Labels `devcontainer.local_folder` (Plan §3 verlangt sie vor dem Code), 246 Mutationen grün in 19 min 36 s, `make gates` Exit 0 in 26 min 7 s, die gemessene Laufzeit der Gegenprobe (DoD, Liefer-Punkt 2). Das Label habe ich selbst gemessen (3a: gleich `$PGR_DEV` wörtlich); die beiden Läufe sind unbelegt, solange sie nicht in §7 stehen. Das Mutationsprotokoll (Zusage, Mutation, rot) und die vier grünen Mutanten mit Invariante, Ort und Sonde stehen in §4 und sind nach Form und meinen 13 Stichproben in Ordnung; sie gehören aber nach §7.
- **V-3 (niedrig, Dokument): zwei Aussagen weiter als der Test** — „Jeder Aufruf lässt sich wiederholen“ (Zeile 75) ohne die zwei akzeptierten Negative aus Vertrag Punkt 12, und „höchstens `PGR_WARTEN` Sekunden“ (§4) gegen den Code, der um bis zu eine Probe überschreitet. Enger fassen (`AGENTS.md` §3.11) oder in §13 aufnehmen.
- **V-4 (niedrig, Dokument): §4 „Image löschen“ nennt `docker rmi <Image>` mit Namen `vsc-…`.** Gemessen sind es zwei Images (`vsc-<Name>-<Hash>-uid` und `vsc-<Name>-<Hash>:latest`); `docker rmi` des ersten ließ das zweite stehen. Die Gegenprobe löscht beide (`id1` und `id1%-uid`).
- **V-5 (niedrig, Schalter): Gegenprobe ohne Docker-Gruppe gibt es nicht.** Wer die Attrappen allein fahren will, hat keinen dokumentierten Weg; der Kopf nennt nur `GP_MUT` und `GP_ENTDECKE`. Das Ziel steht nicht im DoD; Hinweis an den Planner, falls die Trennung für Läufe im Dev-Container gebraucht wird ([ADR-0047](../plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) verlegt die Gegenprobe ohnehin an den Host).

## Aufgeräumt

Alle Objekte dieses Laufs (Container `pgr-cs-ver1`, zwei Volumes, zwei Images, Scratch-Verzeichnis `pgr-cs-ver1`) sind per Name entfernt; die Gegenprobe hat ihre eigenen entfernt. Kein `prune`, kein `pkill`. Der Zweig `rev-start` ist unverändert; neu ist nur diese Datei.

## Nachverifikation (Stand d20096d9)

**Basis:** nur gegen `5e51b0c8` (`git diff 5e51b0c8..HEAD`). Engste Sensoren; volle Gegenprobe (264 Mutationen, 19 min 38 s) und `make gates` (4 min 37 s) nicht wiederholt, ihre Zahlen stehen jetzt in §7 des Plans.

| # | Prüfung | Ergebnis |
|---|---|---|
| 1 | `git diff 5e51b0c8..HEAD -- docs/plan` | eine Datei, ein Hunk, 84 Einfügungen und 1 Löschung; die Löschung ist die Platzhalterzeile `*(wird bei der Closure geschrieben)*` in §7. §1 bis §6 und §8 unverändert. §7 trägt: Messung des Labels, Zeiten (19 min 38 s, 4 min 37 s), Mutationsprotokoll, vier grüne Mutanten (+`p2-vorgabe-warten`) mit Invariante, Ort, Sonde, sechs Rückgaben entschieden, Behandlung F-1 bis F-11 und V-1 bis V-5 |
| 2 | `p15-docker-rest` mit `p15-volume-rest-entfernen`, `p15-volume-rest-melden` (`GP_NUR`, `GP_MUT`, `GP_ENTDECKE`) | grün, 67 s; beide Mutanten färben genau `p15-docker-rest` |
| 3 | `p9-name-regex` | rot genau `p9-name-punkt-fremd` (1,8 s) |
| 4 | Stichproben | `p9-kollision-nach-name` → `p9-doppelt`; `p9-zustand-paused/restarting/removing/dead` → je der eigene Fall, `p9-zustand-geprueft` → alle vier; `p9-leere-bindung` → `p9-kein-port`; `p9-geruest-geprueft` und `-nur-verzeichnis` → `-leer` und `-ohne-json`, `-nur-ordner` → `-ohne-json`; `p9-ps-label-fehler` → `p9-docker-ps-label`; `p13-fehler-als-signal` → `p13-fehler-vor-signal`, `p9-up-exit1`; `p9-update`, `-gestoppt`, `-keiner` → `p9-kein-update`. Alles wie in §7. (Die Namen der Mutationen weichen von denen der Fälle in der §7-Tabelle ab, die Tabelle nennt Fälle; Wortwahl) |
| 5 | Takt von `SECONDS`: `p11-kein-sleep` (= „nie“), `p11-zu-kurz`, `p11-vorgabe-warten` | ruhig: `p11-nie`, `p11-nie`, `p11-vorgabe-warten`, 69 s. Unter Last (20 Rechenschleifen, Load 27 bis 42, zugleich zwei Läufe der Gruppe Attrappen): dieselben drei roten Fälle, 80 s. Keine Lastabhängigkeit |
| 6 | Dokument | §12 trennt Beobachtung, Messung, Folgerung (Absturz des Daemons: „nicht gemessen“ für Exit-Code und `postStartCommand`-Folgerung). „Kein `docker update`“ in allen drei Zweigen belegt (Zeile 4). Wiederholung nennt die zwei Ausnahmen; Wartezeit nennt die Überschreitung um eine Probe, „nicht gemessen“; Image löschen nennt zwei Images. Kein Name und keine Mail des Nutzers, kein Zugangswert (nur Variablennamen und Platzhalter). Hinweis: Zeile 40 enthält die öffentliche Adresse `github.com/pt9912/u-boot` (Upstream des Werkzeugs, wie die Baseline-Adresse in `AGENTS.md`), kein Befund |
| 7 | Echter Scratch-Lauf mit u-boot **0.11.0** (`pgr-cs-ver2`, Port 2320, 2 CPUs, 1g) | Erststart Exit 0, 51 s, `bereit pgr-cs-ver2 127.0.0.1:2320`; zweiter Aufruf 0 s; nach `docker stop` Aufruf 2 s, gleiche ID; `Memory` 1073741824, `NanoCpus` 2000000000. `u-boot doctor`: 0 Fehler, **2 Warnungen**, 16 ok, wie bei 0.10.0 (`compose.yaml has no services`, Egress mit `nestedRuntime docker`). `u-boot generate`/`init` meldet fünf `Warning (LH-FA-DEV-006/007/008)` zum Profil (kein `origin`, `--privileged`, `NET_ADMIN`, sshd); gegen 0.10.0 nicht verglichen (ich habe die 0.10.0-Meldungen nicht aufgehoben), kein Fehler. Danach Container, beide Volumes, beide Images per Name entfernt, nichts mit `pgr-cs` übrig |
| 8 | `make docs-check` (1108 Dateien, 0 Befunde), `make kopf-check`, `make gates-schnell` (11 s) | je Exit 0 |
| 9 | Metrik `git diff --numstat 5e51b0c8..HEAD -- tools harness docs/developer` | 588 + 130 = **718**; 49 in `harness/vertraege/`, ohne sie 669; wie im Plan (über Meldeschwelle 700, unter Marke 900) |

**Anmerkung zu Zeile 5.** Meine zwei Hintergrundläufe der Gruppe Attrappen gingen an der unveränderten Kopie bei `p12-sigINT-*` rot: Ich startete sie als Hintergrundjob einer Shell ohne Jobsteuerung, die `SIGINT` ignoriert. Das ist die im Kopf der Gegenprobe genannte GRENZE, kein Fehler der Probe und keine Lastwirkung; als Last war sie trotzdem wirksam.

### Stand der Befunde

- **V-1 behoben.** §5 bis §7 und die Rückführungs-Bedingungen in §4 stehen wie auf `5e51b0c8`; die Belege hängen in §7 an.
- **V-2 behoben.** Label-Messung, 19 min 38 s, 4 min 37 s, Laufzeit der Gegenprobe, Mutationsprotokoll und grüne Mutanten stehen in §7. Die beiden Läufe habe ich nicht wiederholt (Auftrag), die Zahlen sind Belege des Implementers.
- **V-3 behoben** (beide Aussagen enger bzw. mit Ausnahmen).
- **V-4 behoben** (zwei Images im Dokument).
- **V-5 behoben** als Hinweis im Kopf der Gegenprobe (`GP_NUR` und `GP_MUT='^$'`, `docker` und `devcontainer` bleiben Voraussetzung); ein eigener Schalter existiert weiter nicht, im DoD nicht verlangt.

Offen bleiben Closure-Notiz, Register und Risiko-Ausgänge (Aufgabe der Closure) sowie die Rückgabe an den Planner, §3 nachzuziehen (`pgr-dev2`/2223 gegen `pgr-dev3`/2224, Basis `aad621a8`), die §7 selbst nennt. Aufgeräumt: alle Objekte und Verzeichnisse dieses Laufs per Name entfernt, kein `prune`, kein `pkill`; `pgr-dev2`, `pgr-dev3` unberührt.
