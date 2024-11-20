# logveil

## Description

Logveil is a simple CLI tool for log anonymizaiton. If you ever had a need to create a sample log data but were afraid to leak sensitive info this tool will help you make your logs anonymous.

Currently LogVeil works only with Logmanager-created log data. It needs to be either LM Export (CSV) or LM Backup (GZIP).

Once log data input is supplied, LogVeil will go over each log line, apply anonymization to it, output raw anonymized log to standard output or a file path provided by user and write anonymisation proof in format `{"original":"original_value","new":"new_value"}` to `proof.json` file.

Note that LogVeil is made to work with data streams so it does not load whole input file to memory - which is crucial when dealing with large LM Backup files.

## Disclaimer

**While LogVeil is designed to anonymize logs effectively, the outcome depends on the configuration. We recommend verifying that the data shared has been fully anonymized and contains no sensitive information!**

## Usage

There are two components needed to make this work:

1. Your input log data in CSV (LM Export) or GZIP (LM Backup) format.
2. Anonymization data - which is data that will be used to replace original log data.

```
Usage of ./logveil:
  -d value
        Path to directory with anonymizing data
  -i value
        Path to input file containing logs to be anonymized
  -o value
        Path to output file (default: Stdout)
  -v
        Enable verbose logging
  -e
        Change input file type to LM export (default: LM Backup)
  -p
        Disable proof wrtier (default: Enabled)
  -h
        Help for logveil
```

**Examples:**

1. Read log data from LM Export file (CSV), output anonymization result to `output.txt` file and write anonymization proof to `proof.json` file.

`./logveil -d example_anon_data/ -e -i lm_export.csv -o output.txt`

2. Read log data from LM Backup file (GZIP), output anonymization result to `output.txt` file and write anonymization proof to `proof.json` file.

`./logveil -d example_anon_data/ -i lm_backup.gz -o output.txt`

3. Read log data from LM Backup file (GZIP), output anonymization result to `output.txt` file and disable writing anonymization proof.

`./logveil -d example_anon_data/ -i lm_backup.gz -o output.txt -p`

4. Read log data from LM Export file (CSV), output anonymization result to standard output (STDOUT) and disable writing anonymization proof.

`./logveil -d example_anon_data/ -e -i lm_export.csv -p`

5. Read log data from LM Export file (CSV), output anonymization result to standard output (STDOUT), disable writing anonymization proof and enable verbose logging.

`./logveil -d example_anon_data/ -e -i lm_export.csv -p -v`

### How it works

**This is only a simplified example and does not match 1:1 with how anonymization is actually implemented**

Consider below log line. It is formatted in a common `key:value` format.

```
{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip":"89.239.31.49", "username":"test.user@test.cz", "organization":"TESTuser.test.com", "mac": "71:e5:41:18:cb:3e"}
```

First, LogVeil will load anonymization data from supplied directory (`-d example_anon_data/`). Each file in that folder should be named according to the values it will be masking. For example, lets assume we have following directory structure:

1. `username.txt`
2. `organization.txt`

Next, LogVeil will go over each log line in supplied input and extract `key:value` pairs from it. When applied to above log line it would look like this:

1. `"@timestamp": "2024-06-05T14:59:27.000+00:00"`
2. `"src_ip":"89.239.31.49"`
3. `"username":"test.user@test.cz"`
4. `"organization":"TESTuser.test.com"`
5. `"mac": "71:e5:41:18:cb:3e"`

Then, LogVeil will try to match extracted pairs to anonymization data it loaded in previous step. Two paris should be matched:

1. `"src_ip":"89.239.31.49"` with `src_ip.txt`
2. `"username":"test.user@test.cz"` with `username.txt`
3. `"organization":"TESTuser.test.com"` with `organization.txt`

And one pair should be matched by regular expression scanning:

1. `"mac": "71:e5:41:18:cb:3e"`

Now LogVeil will grab values (randomly) from files which filenames matched with keys, generate new value for `mac` key and create a replacement map in format `"original_value":"new_value"`:

1. `"89.239.31.49":"10.20.0.53"`
1. `"test.user@test.cz":"ladislav.dosek"`
2. `"TESTuser.test.com":"Apple"`
3. `"71:e5:41:18:cb:3e": "0f:da:68:92:7f:2b"`

Now each element from the above list will be iterated over and compared against log line. Whenever `original_value` is found it will be replaced with `new_value`. Outcome should look like this:

```
{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip":"10.20.0.53", "username":"ladislav.dosek", "organization":"Apple", "mac": "0f:da:68:92:7f:2b"}
```

```
{"original": "27.221.126.209", "new": "10.20.0.53"},
"{"original":"test.user@test.cz","new":"ladislav.dosek"}"
"{"original":"TESTuser.test.com","new":"Apple"}"
{"original": "71:e5:41:18:cb:3e", "new": "0f:da:68:92:7f:2b"},
```

### Anonymization data

Each `key:value` pair which you want to anonymize data must have its equivalent in anonymization data folder.

If anonymization data does not exist for any given `key:value` pair then LogVeil will attempt to use regular expressions to match and replace common values such as: IPv4, IPv6, MAC, Emails and URLs.

For example, if you want to anonymize values in `organization` and `username` keys, you need to have two files of the same name in anonymization folder containing some random data.

### Output

Anonymized data will be outputted to provided file path in txt format.

Alternatively, if you don't provide file path, output will be written to the console.

## Release

Go to: https://github.com/logmanager-oss/logveil/releases to grab latest version of LogVeil. It is available for Windows, Linux and MacOS (x86_64/Arm64).

We are using Goreleaser (https://goreleaser.com) for building LogVeil release file.

If you wish to create your own release do the following:

1. Clone the repository
2. Run `CGO_ENABLED=0 GOOS=<your_target_OS> GOARCH=<your_target_CPU_architecture> go build -o logveil ./cmd/main.go`
