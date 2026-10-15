# Dieselben Fälle stehen je Format in einer eigenen Tabelle

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die 14 Stellen des Endes stehen als `TestSchreibenNimmtFormverletzungAn` (YAML) und als
zweite Tabelle in `TestSQLiteSchreibenNimmtFormverletzungAn`. Beide unterscheiden sich in
der Prüfung der Version (Text des Kopfes gegen `meta.version`), in den Hilfen und im
Schlüssel je Stelle, der die Mutanten schlüsselt; ein gemeinsamer Treiber ändert beide
Test-Dateien und die Mutanten. Eine Grenze, kein Mangel am Vertrag. Kommt ein dritter
Leser oder Schreiber, wächst eine dritte Tabelle.
