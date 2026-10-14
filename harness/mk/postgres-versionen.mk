# harness/mk/postgres-versionen.mk — Versionsmatrix der Integrationstests
# (ADR-0042; Vertrag harness/vertraege/SPEC-053-postgres-versionen.md)
# · seit slice-v1-abschluss-postgres-versionen-matrix.
# Die Liste der Versionen und die Teilmenge der Tests stehen hier und nur hier.
# `test-postgres-versionen` fährt die Teilmenge gegen jede Version der Liste außer
# der Referenz (die fährt `test-integration` mit der vollen Suite); die Gegenprobe
# fährt es in Kopien des Arbeitsbaums mit einer Attrappe von docker.
# Bestätigt am 2026-10-14 gegen die Release-Liste von PostgreSQL: 14 bis 18. Die Liste
# wird von Hand gepflegt (SPEC-053 Punkt 14).
# GEPRÜFT DURCH tools/test/postgres-versionen-gegenprobe.sh.
POSTGRES_REFERENZ := 17
POSTGRES_VERSIONEN := \
	14=postgres:14-alpine@sha256:4ea9e5ed06591da7ea23eb65465e8d3187fe79f4d5ec3ae976d29a33b013e77a \
	15=postgres:15-alpine@sha256:f7d23353e1b15400d22ebe31189f4d314b87a4c129cc400c8c2d8d4ca127bf81 \
	16=postgres:16-alpine@sha256:721873c34ceb9f8d8fc265984940dc982404c105f19ad51be9fdc5970a6080ea \
	17=postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24 \
	18=postgres:18-alpine@sha256:77f585114c32fbca283dc835b0596f4e52b51b4c6662d7810b2f4084f60a1873
# Abnahmeszenarien 1, 2 und 7, der Fall SCRAM-SHA-256 der Anmeldung, der Abgleich der
# Lebendprüfung mit dem Server; Phase 1 mit PostgreSQL, Phase 2 ohne.
POSTGRES_TESTS_PHASE1 := TestE2ERecordSelect1 TestE2EVorbereitungOhnePostgres TestE2EVorbereitungExtendedOhnePostgres TestE2ERecordExtendedPgx TestE2ERecordAnmeldung/scram-sha-256 TestE2EReplayLebendpruefungWiePostgres
POSTGRES_TESTS_PHASE2 := TestE2EOhnePostgresReplay TestE2EOhnePostgresExtendedReplay

.PHONY: test-postgres-versionen test-postgres-versionen-gegenprobe

test-postgres-versionen: ## Teilmenge der Integrationstests gegen jede Version der Liste außer der Referenz
	@POSTGRES_VERSIONEN='$(POSTGRES_VERSIONEN)' POSTGRES_REFERENZ='$(POSTGRES_REFERENZ)' POSTGRES_IMAGE='$(POSTGRES_IMAGE)' POSTGRES_TESTS_PHASE1='$(POSTGRES_TESTS_PHASE1)' POSTGRES_TESTS_PHASE2='$(POSTGRES_TESTS_PHASE2)' DOCKER='$(DOCKER)' bash tools/test/run-postgres-versionen.sh

test-postgres-versionen-gegenprobe: ## Gegenprobe der Versionsmatrix: Eingabe, Teilmenge, Ablauf, Aufräumen, Aufrufarten, mit Attrappe von docker
	@bash tools/test/postgres-versionen-gegenprobe.sh

GATE_CHECKS += test-postgres-versionen test-postgres-versionen-gegenprobe
