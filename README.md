# ptg-weather-demo

A small **weather CLI** — resolve the current location by IP and print the current
temperature from the free, key-less [Open-Meteo](https://open-meteo.com) API.

This repository is a **dogfooding audit fixture** for Compozy's *Parallel Task
Groups* feature. It is decomposed into dependency-independent task groups so each
is executed concurrently in its own git worktree/branch and lands as its own small
pull request. See `AUDIT.md` for the full scenario log and evidence.
