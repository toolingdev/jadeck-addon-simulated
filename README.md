# Simulated devices — PTZ controller add-on

A practice switcher, slides and mixer that touch no real gear.

Install: it ships with the controller (Settings → Add-ons); newer releases of this repository arrive as updates there.

Build every platform and the `.ptzaddon` package:

```bash
go run github.com/toolingdev/ptz-controller/sdk/addon/cmd/ptz-addon build
```

Try it on a controller: Settings → Add-ons → Developer → Load a folder (this folder),
or install `dist/simulated.ptzaddon` with Install from file.

Release: push a tag `v1.2.3`; the workflow builds, signs (secret `PTZ_SIGNING_KEY`)
and attaches `simulated.ptzaddon` to the GitHub release.
