# secretly

Secretly is a tool which connects to AWS SecretsManager and loads one or more "secrets" as environment variables, which are then passed as environment variables into the specified command.

Each Secret, in AWS terminology, is one or more key-value pairs and can be set within AWS via the web Console, CLI or API.

This can be useful when you need to inject environment variables into a third-party command/binary which you need to run, allowing for centralised management of the secrets for use across containers or instances.  

## Prerequisites

GoLang ^1.11 [Download here](https://golang.org/project/)

## Examples

Assuming you wanted to run a command called `mybinary` (along with a couple of arguments) and you had set up a Secret within SecretsManager called `mySMsecret` and it consisted of the following key-value pairs:

```
FOO: bar
mypass: secret
```

You could then use secretly as follows:

```
secretly run mySMsecret -- mybinary arg1 arg2
```

This is equivalent to running:

```
FOO=bar mypass=secret mybinary arg1 arg2
```

You can pass in multiple Secrets. The Secrets are evaluated in order, and if the same Key is specified in multiple secrets, the _last_ occurance always takes precedence.

```
secretly run mySMsecret myOtherSecret -- mybinary arg1 arg2
```

You can use the tool to check that secrets all have the same keys available. Useful for making sure that the configuration of a new feature, for example, is rolled out across all environments. If not the same, the application exists with a non-zero exit code, and some debug data.

```
secretly compare myProdSecrets myStagingSecrets myQASecrets
```

## Using Env Files

It can be useful to test this with files when developing locally, as opposed to calling out to SecretsManager. When the `--use-files` (`-f`) flag is specified, the secrets that follow are treated as filenames.

```
secretly run --use-files myEnvFile1 myOtherFile -- mybinary arg1 arg2
```

This uses the godotenv package to read the files, so as long as the files are formatted correctly for godotenv to read them (and there are several formatting options) then they will be loaded in.

Using a mix of local files and SecretsManager variables is not currently supported.