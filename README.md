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
  -c value
        Path to input file with custom anonymization mapping
  -v
        Enable verbose logging
  -e
        Change input file type to LM export (default: LM Backup)
  -p
        Disable proof writer (default: Enabled)
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

6. Read log data from LM Export file (CSV), output anonymization result to standard output (STDOUT) and load custom mapping from custom_mapping.txt

`./logveil -d example_anon_data/ -e -i lm_export.csv -c custom_mapping.txt`


## Anonymization functionality

There are three ways LogVeil anonymizes data:

### Custom anonymization mappings

You can provide custom anonymization mappings for LogVeil to use. They will take precedence over any other anonymization functionality.

Custom mappings can be enabled by using flag `-c <file_path>` and must have the following format:

`<original_value>:<new_value>`

Each custom mapping must be separated by new line. For example:

`test_custom_replacement:test_custom_replacement123`\
`replace_this:with_that`\
`test123:test1234`

### Anonymization data

You can also provide sets of fake data to use when anonymizing.

Consider below log line:

```
{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip":"89.239.31.49", "username":"test.user@test.cz", "organization":"TESTuser.test.com", "mac": "71:e5:41:18:cb:3e", "replacement_test":"replace_this"}
```

If you want to anonymize values in `organization` and `username` keys, you need to have two files of the same name in anonymization data folder and enable them by using `-d <path_to_fake_data_folder>` flag.

1. `username.txt`
2. `organization.txt`

Both files should contain appropriate fake data for the values they will be masking.

### Regexp scanning and dynamic fake data generation

LogVeil implements regular expressions to look for common patterns: IP (v4, v6), Emails, MAC and URL. Once such pattern is found it is replaced with fake data generated on the fly.

## Output

Anonymized data will be written to provided file path in txt format. Alternatively, if you don't provide output file path it will be written to the console (stdout).

Additionally LogVeil will write anonymization proof to `proof.json`, to show which values were anonymized. Proof has a following format:

```
{"original":"<original_value>", "new":"<new_value>}
```

## How it works

**This is only a simplified example and does not match 1:1 with how anonymization is actually implemented**

Consider below log line. It is formatted in a common `key:value` format.

```
{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip":"89.239.31.49", "username":"test.user@test.cz", "organization":"TESTuser.test.com", "mac": "71:e5:41:18:cb:3e", "replacement_test":"replace_this"}
```

First, LogVeil will load anonymization data from supplied directory (`-d example_anon_data/`). Each file in that folder should be named according to the values it will be masking. For example, lets assume we have following directory structure:

1. `username.txt`
2. `organization.txt`

Second, if available, LogVeil will load the custom anonymization mapping from user supplied path. For example, assume we have following file `custom_mapping.txt` with below content:

1. `test_custom_replacement:test_custom_replacement123`
2. `replace_this:with_that`
3. `test123:test1234`

Now anonymization process can start. LogVeil will grab log lines from supplied input, one by one, and apply anonymization to it three steps:

1. Replace values based on custom anonymization mapping
2. Replace values based on loaded anonymization data
3. Replace values based on regular expression matching and fake data generation

Final output should look like this:

```
{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip":"10.20.0.53", "username":"ladislav.dosek", "organization":"Apple", "mac": "0f:da:68:92:7f:2b", "replacement_test":"with_that"}
```

And anonymization proof:

```
{"original":"replace_this", "new":"with_that"}
{"original": "27.221.126.209", "new": "10.20.0.53"},
{"original":"test.user@test.cz","new":"ladislav.dosek"},
{"original":"TESTuser.test.com","new":"Apple"},
{"original": "71:e5:41:18:cb:3e", "new": "0f:da:68:92:7f:2b"},
```

## Release

Go to: https://github.com/logmanager-oss/logveil/releases to grab latest version of LogVeil. It is available for Windows, Linux and MacOS (x86_64/Arm64).

We are using Goreleaser (https://goreleaser.com) for building LogVeil release file.

If you wish to create your own release do the following:

1. Clone the repository
2. Run `CGO_ENABLED=0 GOOS=<your_target_OS> GOARCH=<your_target_CPU_architecture> go build -o logveil ./cmd/main.go`
