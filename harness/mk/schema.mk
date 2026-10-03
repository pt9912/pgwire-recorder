# harness/mk/schema.mk — neutrales Schema der SQLite-Aufzeichnung (tools/schema/schema.yaml)
# mit d-migrate. Beide Targets sind Werkzeuge und hängen an keinem GATE_CHECKS-Eintrag:
# sie brauchen das d-migrate-Image, das lokal vorhanden sein muss. Der Pin ist
# ein bewusster Commit (Modul 14); er stimmt mit dem Pin der Schwester-Repos überein.
D_MIGRATE_IMAGE ?= ghcr.io/pt9912/d-migrate@sha256:af9d3eb323a6cfd13eb63f012823788e510f764509918933d2ff98c8e68e4212

.PHONY: schema-validate schema-generate
schema-validate: ## d-migrate: neutrales Schema der SQLite-Aufzeichnung prüfen (netzlos, kein Gate)
	docker run --rm --network none -v "$(CURDIR)/tools/schema":/work:ro $(D_MIGRATE_IMAGE) schema validate --source /work/schema.yaml

schema-generate: ## d-migrate: SQL für SQLite aus dem neutralen Schema auf stdout (netzlos, kein Gate)
	docker run --rm --network none -v "$(CURDIR)/tools/schema":/work:ro $(D_MIGRATE_IMAGE) schema generate --source /work/schema.yaml --target sqlite
