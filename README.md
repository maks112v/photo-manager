# Photo Manager

A simple photo manager that organizes photos into a folder structure based on the date the photo was taken.

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

- Name
- Year
- Month
- PhotoCount

#### Photo File Template Variables

Default `{{.Name}}{{.Ext}}`

- Name
- Ext
- CreatedAt (Date the photo was taken)