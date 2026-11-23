#!/bin/bash -e

PYTHONPATH=src python3 -m unittest discover -s tests
