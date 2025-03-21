# Evaluation for IIC3413

## Requirements
To run this app you must have [docker](https://docs.docker.com/engine/install/)
and [make](https://linux.die.net/man/1/make) installed in your machine.

## Commands
- `make build`: build docker image.
- `make run`: runs docker image.

## Usage
For every lab, the submissions and tests should be placed in new
sub-directories inside the `io/submissions` and `io/tests` directories
respectively. The names of this sub-directories (not their path) should be
specified in the `submissions` and `test` fields in `config/config.yaml` file.

In order to prevent submissions from reporting fraudulent data or altering
the results output file a random secret must be set for both the `output_name`
and `verification_code` fields in the same configuration file. A random 16
byte string encoded in base 32 should do the job.

Tests should be written such that their output are of the form:
```
{verification_code} {points_obtained}
```
which means that a verification code must be chosen in advance.

After running the tests using the `make run` command the output can be found
in `io/results/{output_name}.csv`.
