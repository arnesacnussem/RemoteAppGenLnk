# Simple RDP Shortcut Creator for [kimmknight/RemoteAppTool](https://github.com/kimmknight/remoteapptool)

[![Build & Release](https://github.com/arnesacnussem/RemoteAppGenLnk/actions/workflows/build.yaml/badge.svg)](https://github.com/arnesacnussem/RemoteAppGenLnk/releases)

I don't know why I can't just bundle every remote apps into one .msi, so I write this simple program and I can create
some shortcut with icon easily.

### Check the help message first
```powershell
ragl -h
```

### Example:

```powershell
ragl -p ./path_to_rdp_and_ico_files -l path_to_lnk_file_will_be -n '{{ .Name }} (remote)'
```
