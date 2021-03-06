# aserto-idp - the CLI for managing idp plugins

The aserto-idp CLI is a tool for importing user data from identity providers (idp) to Aserto or another idp.

---

## Installation

`aserto-idp` is available on Linux, macOS, and Windows platforms.

- Binaries for Linux, Windows, and Mac are available as tarballs in the [release](https://github.com/aserto-dev/aserto-idp/releases) page.

- Via Homebrew for macOS or LinuxBrew for Linux

  ```shell
  brew tap aserto-dev/tap && brew install aserto-idp
  ```

- Via a GO install

  ```shell
  # NOTE: The dev version will be in effect!
  go install github.com/aserto-dev/aserto-idp/cmd/aserto-idp@v0.8.5
  ```

---

## The command line

At first, the help message will look like this:

```
Usage: aserto-idp <command>

Aserto Identity Provider CLI

Commands:
  delete          delete user ids from an user-provider idp
  exec            import users from an user-provided idp to another user-provided idp
  get-plugin      download plugin
  list-plugins    list available plugins
  version         version information

Flags:
  -h, --help             Show context-sensitive help.
  -c, --config=STRING    Path to the config file. Any argument provided to the CLI will take precedence.
  -v, --verbosity=INT    Use to increase output verbosity.
```

The specific flags for a specific plugin will appear only after the plugin was downloaded.

## Plugins

The plugins will be downloaded on the system in a directory under the following path: `$HOME/.aserto/idpplugins` .

Currently, the available plugins are:

- aserto
- okta
- json
- auth0

The plugins can be downloaded in 2 ways:

- using the `get-plugin` command
- calling the `exec` or `delete` command with the name of a plugin that is not on the system (this will automatically download the latest version of that plugin).

`get-plugin` examples:

To download the latest version of a plugin:

```
aserto-idp get-plugin aserto
```

or

```
aserto-idp get-plugin aserto:latest
```

To download a specific version of a plugin:

```
aserto-idp get-plugin aserto:1.0.1
```

To see the plugins that are downloaded on the system and their version, the `list-plugins` command can be used.

`list-plugins` examples:

To list plugins that are currently on the system:

```
aserto-idp list-plugins
```

The output will be similar to:

```
	 aserto:0.8.5
	 auth0:0.8.3
	 json:0.8.4
	 okta:0.8.8
```

To list plugins and versions that are available remotely and can be downloaded:

```
aserto-idp list-plugins --remote
```

The output will be similar to:

```
Available versions for 'aserto'
*	 aserto:0.8.5
	 aserto:0.8.4
	 aserto:0.8.3
	 aserto:0.8.2
	 aserto:0.8.1
	 aserto:0.8.0

Available versions for 'json'
*	 json:0.8.4
	 json:0.8.3
	 json:0.8.2
	 json:0.8.1
	 json:0.8.0

Available versions for 'okta'
*	 okta:0.8.8
	 okta:0.8.7
	 okta:0.8.6
	 okta:0.8.5
	 okta:0.8.4
	 okta:0.8.3
	 okta:0.8.2
	 okta:0.8.1
	 okta:0.8.0

Available versions for 'auth0'
*	 auth0:0.8.3
	 auth0:0.8.2
	 auth0:0.8.1
	 auth0:0.8.0
```

where `*` symbolize the version that is currently on the system.

---

## The config

The config has YAML format and its content should contain credentials for the idp you are trying to use.

Eg.:

```
logging:
  log_level: LEVEL
plugins:
  auth0:
    domain: DOMAIN
    client-id: ID
    client-secret: SECRET
    connection-name: AUTH0_CONNECTION_NAME
  json:
    from-file: PATH_TO_FILE
    to-file: PATH_TO_OUTPUT_FILE
  aserto:
    tenant: TENANT
    authorizer: AUTHORIZER
    api-key: API_KEY
    insecure: BOOL
  okta:
    domain: OKTA_DOMAIN
    api-token: TOKEN
```

---

## Logs

Logs are printed to `stdout`. You can increase detail using the verbosity flag (e.g. `-vvv`).

---

## Usage examples

To import user data from an idp to aserto:

```
aserto-idp exec --from json --to aserto -c PATH_TO_CONFIG
```

Note that if json or aserto plugins are not on the system, using this command, they will be automatically downloaded. Also if there is a newer version of either one of the plugins used, the following message will be prompted:

```
A new version '0.8.4' of the plugin 'json' is available
```

To disable updates checking when using `exec` or `delete` :

```
aserto-idp exec --from json --to aserto -c PATH_TO_CONFIG --no-update-check
```

or

```
aserto-idp exec --from json --to aserto -c PATH_TO_CONFIG -n
```

You can delete a user from aserto knowing its id and using the following:

```
aserto-idp delete --from aserto USER_ID
```

---

## Plugin development

If you want to develop your plugin you can check out our example for a dummy plugin [here](https://github.com/aserto-dev/idp-plugin-sdk/tree/main/examples/dummy)

---
