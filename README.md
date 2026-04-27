# driftwatch

> A CLI tool that detects configuration drift between deployed infrastructure and your IaC definitions.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your IaC directory and let it compare against your live infrastructure:

```bash
# Scan for drift using a Terraform state file
driftwatch scan --provider aws --state ./terraform.tfstate --config ./infra/

# Output results as JSON
driftwatch scan --provider aws --state ./terraform.tfstate --format json

# Watch for drift continuously (every 5 minutes)
driftwatch watch --interval 5m --config ./infra/
```

### Example Output

```
[DRIFT DETECTED]
  Resource:  aws_security_group.web
  Field:     ingress.cidr_blocks
  Expected:  ["10.0.0.0/8"]
  Actual:    ["0.0.0.0/0"]

1 drift(s) found across 42 resource(s).
```

---

## Supported Providers

- AWS
- GCP
- Azure *(coming soon)*

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)