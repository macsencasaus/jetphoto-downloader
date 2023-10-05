# Jetphotos Downloader
This script downloads the latest image of some plane (by any search query) from [Jetphotos](https://www.jetphotos.com).

For example:
```
go run main.go n804an
```
This will download a nice image of an American Airlines 787-8 in an `img` directory.

You can download more than one image by seperating the registrations with spaces:

```
go run main.go n352ps n173us n218nn
```

Download latest from a photographer:
```
go run main.go "macsen casaus"
```
or airport:
```
go run main.go "Dallas/Fort Worth Int'l Airport"
```