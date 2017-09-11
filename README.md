# Local Bookmark (lbm)

## Summary

A simple bookmark manager for terminal.

## Installation



## Usage

```
  Usage:
    Add a bookmark
      $ lbm add [directory]
    Return a bookmark
      $ lbm get [id]
    Go to a bookmark
      $ cd ${lbm get [id]}
    Remove a bookmark
      $ lbm rm [id]
    List bookmarks
      $ lbm ls
```

## How it works

Local Bookmark uses a plain text file named `.lbm_dump` to store
bookmarks in CSV format. Each line in the data file represent an instance
of a bookmark. A bookmark instance consists of two attributes, a ID and a
directory location.

## License

MIT
