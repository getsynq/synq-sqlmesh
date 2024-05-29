#!/bin/bash
git describe --tags --abbrev=0 > build/version.txt
date +%FT%T%z > build/time.txt