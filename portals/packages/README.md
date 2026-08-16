# Packages

Shared code used across the portals in `apps/` (e.g. auth wiring, the
`window.configs` loader, common UI primitives) lives here as pnpm workspace
packages. Nothing has been extracted yet — this directory is scaffolding for
when that duplication gets pulled out.

Each package should be a normal pnpm workspace member (its own
`package.json`), consumed by an app via a `workspace:*` dependency.