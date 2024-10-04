# logveil

## Description

Logveil is a simple CLI tool for log anonymizaiton. If you ever had a need to create a sample log data but were afraid to leak sensitive info this tool will help you make your logs anonymous.

## Disclaimer

**While LogVeil is designed to anonymize logs effectively, the outcome depends on the configuration. We recommend verifying that the data shared has been fully anonymized and contains no sensitive information!**

## Usage

There are two components needed to make this work:

1. Your input log data in CSV format.
2. Anonymization data - which is data that will be used to replace original log data.

```
Usage of ./logveil:
  -d value
        Path to directory with anonymizing data
  -i value
        Path to input file containing logs to be anonymized
  -o value
        Path to output file containing anonymized logs
  -v
        Enable verbose logging
  -h
        Help for logveil
```

**Example:**

`./logveil -d example_anon_data/ -i test_log.csv -o output.txt`

### Input log data

Obviously first you need to provide log data to be anonymized. It needs to be in a CSV format. The columns in you CSV file will mark which values you want to anonymize.

As an example consider below log line. It is formatted in a standard `key:value` format. Key names mark the values.

```
{"@timestamp": "2024-06-05T14:59:27.000+00:00", "msg.src_ip":"89.239.31.49", "username":"test.user@test.cz", "organization":"TESTuser.test.com"}
```

As such we can easily parse it into CSV file:

```
@timestamp,msg.src_ip,msg.username,msg.organization,raw
2024-06-05T14:59:27.000+00:00,89.239.31.49,test.user@test.cz,TESTuser.test.com,"{""@timestamp"": ""2024-06-05T14:59:27.000+00:00"", ""msg.src_ip"":""89.239.31.49"", ""username"":""test.user@test.cz"", ""organization"":""TESTuser.test.com""}"
```

Now key names are simply column names in CSV file. `raw` contains original log line. When you run Logveil, column names will be matched against your anonymization data.

You can easily extract log data in such format from your Logmanager. Refer to Logmanager documentation for more info on how to Export data.

### Anonymization data

Each column for which you want to anonymize data must have its equivalent in anonymization data folder.

For example, if you want to anonymize values in `msg.src_ip` and `msg.username` columns, you need to have two files of the same name in anonymization folder.

### Output

Anonymized data will be outputted to provided file path in txt format (unparsed).

Alternatively, if you don't provide file path, output will be written to the console.

## Release

Go to: https://github.com/logmanager-oss/logveil/releases to grab latest version of LogVeil. It is available for Windows, Linux and MacOS (x86_64/Arm64).

We are using Goreleaser (https://goreleaser.com) for building LogVeil release file.

If you wish to create your own release do the following:

1. Clone the repository
2. Run `CGO_ENABLED=0 GOOS=<your_target_OS> GOARCH=<your_target_CPU_architecture> go build -o logveil ./cmd/main.go`
