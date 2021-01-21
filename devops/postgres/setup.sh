#!/bin/bash
set -e

psql postgres -d db -f /fixtures/db.sql