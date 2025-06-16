# Evaluation for IIC3413

## Requirements

To run this app you must have [docker](https://docs.docker.com/engine/install/)
and [make](https://linux.die.net/man/1/make) installed in your machine.

## Commands

- `make build`: build docker image.
- `make run`: runs docker image.

## Usage

For a given lab `lab_n`, the submissions, tests, input databases and expected
outputs must be placed in the `io` directory under the following structure:

```
io
|
+ --- submissions
|     |
|     + --- lab_n
|             submission_1.zip
|             submission_2.zip
+ --- tests
|     |
|     + --- lab_n
|             test_1.cc
|             test_1.cc
+ --- data
      |
      + --- lab_n
            |
            + --- eval_dbs
            |     |
            |     + --- test_1_db
            |     |       catalog.dat
            |     |       table
            |     + --- test_2_db
            |             catalog.dat
            |             table
            + --- outputs
                    test_1.output
                    test_2.output
```

A successful test, should write its output to an `outputs` directory inside the
path where it's executable is called following the naming scheem of the files
in `io/data/{lab_n}/outputs`. If a test encounters an error it may not write an
output only if it exists with a none 0 code. In this case the evaluation
program will write a predetermined output under the name
`{test_file_name}.output`, ensuring that this will be evaluated as wrong.

To run the tests use the following command:

```
make LAB_NAME={lab_n} run
```

Where `{lab_n}` corresponds to the name of the lab used when creating the
aforementioned directories. The output will be written as a csv file in the
`io/results` directory.

## Security

In order to prevent submissions from overwriting tests, reading the expected
outputs or tampering with other submissions they are executed as a separate
user with only the following permissions:

- Executing files in `wkdir/build/Relese/bin`.
- Writing to `wkdir/data` and `wkdir/outputs`.
