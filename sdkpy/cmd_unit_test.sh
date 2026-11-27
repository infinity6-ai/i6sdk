#!/bin/bash -e

[ -d src/i6sdkpy/sdkapis ] || go run . generate
PYTHONPATH=src python3 -m unittest discover -s tests
