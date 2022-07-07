# Simple Distributed File System

Inspired by the [Google File System](https://static.googleusercontent.com/media/research.google.com/en//archive/gfs-sosp2003.pdf).

## Usage

### Client

```bash
sdfs ls # lists files in the system
sdfs create largefile.mp4 # uploads large file into the system
sdfs read largefile.mp4 > cat
sdfs rm largefile.mp4 # deletes large file into the system
```
