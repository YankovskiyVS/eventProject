# Build project
```
make build
```
# Run application
```
make run
```
# Run unit tests
```
make test
```
# Run integration tests
in the _settings.json_ file add 
```
"go.buildFlags": ["-tags=integration"]
```
to the **"gopls:"**
---
```
make test-integration
```
# Clean build artifacts
```
make clean
```