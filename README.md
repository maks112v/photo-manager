# Photo Manager

Automatically organize photos from an SD card or folder into albums based on the date the photo was taken. When photos haven't been taken in 24 hours, it will create a new album. It will traverse any subdirectories and organize all photos into the destination directory.

## Example

```
/source
  /fuji-101
    - DSCF0001.jpg (2024-01-01 12:00:00)
    - DSCF0002.jpg (2024-01-01 12:00:01)
    - DSCF0003.jpg (2024-02-03 12:00:00)
  /fuji-102
    - DSCF0004.jpg (2024-02-03 12:00:01)
    - DSCF0005.jpg (2024-02-03 12:00:02)
    - DSCF0006.jpg (2024-02-03 12:00:03)

/destination
  /2024-01 Album 1
    - DSCF0001.jpg
    - DSCF0002.jpg
  /2024-02 Album 2
    - DSCF0003.jpg
    - DSCF0004.jpg
    - DSCF0005.jpg
    - DSCF0006.jpg
```

## Install

```bash
brew tap maks112v/tap
brew install photomanager
```

## Usage

```bash
photomanager init
photomanager organize
```

## Reference

#### Album Path Template Variables

Default `{{.Year}}-{{.Month}} {{.Name}}`

- Name (`Album #`)
- Year (Year of the first photo in the album)
- Month (Month of the first photo in the album)
- PhotoCount

#### Photo File Template Variables

Default `{{.Name}}{{.Ext}}`

- Name (Original file name without extension)
- Ext (File extension with dot e.g. `.jpg`)
- CreatedAt (Date the photo was taken)