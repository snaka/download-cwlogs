# download-cwlogs

This tool exports CloudWatch Logs to a CSV file.

## Installation

```bash
go install github.com/snaka/download-cwlogs/cmd/download-cwlogs@latest
```

## Usage

```bash
download-cwlogs -g <log-group-name> -s <log-stream-name> -r <region> -o <output-file>
```

- `-g`, `--log-group`: Log group name (required)
- `-s`, `--log-stream`: Log stream name (required)
- `-r`, `--region`: AWS region (default: `ap-northeast-1`)
- `-o`, `--output-file`: Output file name (default: `output.csv`)

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
