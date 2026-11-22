# Download and run Pandafish v1.0

Open the [v1.0 release](https://github.com/philyawj/Pandafish/releases/tag/v1.0)
and expand **Assets**. Choose the archive for the operating system and processor
listed in the [README download table](../README.md#download-v10). Download
`SHA256SUMS` as well. Binaries are release attachments, not tracked Git files.
Each archive contains the engine, `LICENSE` and `THIRD_PARTY_NOTICES.md`.

## Choose a processor

- On macOS, **About This Mac** identifies Apple silicon or Intel. Apple silicon
  uses `darwin-arm64`; Intel uses `darwin-amd64`.
- On Linux, `uname -m` normally reports `x86_64` for `linux-amd64` or `aarch64`
  for `linux-arm64`.
- On Windows, **Settings → System → About → System type** identifies an x64 or
  ARM-based processor. Choose `windows-amd64` or `windows-arm64` respectively.

## macOS and Linux

Extract the matching `.tar.gz` archive in Finder, a file manager, or a terminal:

```sh
tar -xzf Pandafish-v1.0-darwin-arm64.tar.gz
cd Pandafish-v1.0-darwin-arm64
./pandafish
```

Substitute the downloaded archive/directory name for another target. On macOS,
the binaries have no Apple Developer distribution signature and are not notarized.
If macOS blocks the verified executable, review the download and use
**System Settings → Privacy & Security →
Open Anyway** for that executable. Building locally is another option.

Verify the archive's SHA-256 against its line in `SHA256SUMS`:

```sh
shasum -a 256 Pandafish-v1.0-darwin-arm64.tar.gz
```

On Linux, `sha256sum ARCHIVE.tar.gz` provides the same digest.

## Windows

Right-click the `.zip` file and select **Extract All**. Open PowerShell in the
extracted `Pandafish-v1.0-windows-amd64` (or ARM64) directory and run:

```powershell
.\pandafish.exe
```

Verify the downloaded archive with:

```powershell
Get-FileHash .\Pandafish-v1.0-windows-amd64.zip -Algorithm SHA256
```

Compare the result with the matching entry in `SHA256SUMS`.

## Use the executable

Pandafish waits for UCI commands; it has no graphical board. Try `uci`, then
`isready`, or configure a chess GUI to use the extracted executable. The
[README example](../README.md#try-the-engine) shows a complete move request.
For Lichess, point the external client at the executable using the
[Lichess setup guide](LICHESS.md).

The release host tests macOS ARM64 at runtime. macOS Intel, Linux and Windows
builds are cross-compiled; native runtime testing on those systems is not claimed.
`build-manifest.json` records the source revision, Go version, targets and
archive identities. Development is paused after the v1.0 milestone.
