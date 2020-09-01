# Jeera

Maybe maybe maybe something something JIRA.

## Config

`~/.config/jeera/settings`

```json
{
  "JIRA": {
    "URL": "",
    "EpicCustomFieldID": ""
  },
  "OAuth": {
    "Token": "",
    "Secret": "",
    "PEM": "",
    "ConsumerKey": "",
    "ConsumerSsecret": ""
  }
}
```

## Usage

```console
$ go run cmd/viz/*go <EPIC KEY> | dot -Tsvg > out.svg
...
```
