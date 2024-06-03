# chiSSL Admin CLI

The `chissl` admin CLI allows you to manage users in the chiSSL server. This includes adding, deleting, retrieving, and listing users. The commands are secured using Basic Authentication.

## Known limitations and TODOs
- Add support of using proxy config
- Add support for using custom TLS configs
- Use client side hashed passwords


## Usage

```sh
chissl admin [subcommand] [options]
```

### Subcommands

* `adduser` - Adds a new user
* `deluser` - Deletes an existing user
* `getuser` - Retrieves info about an existing user
* `listusers` - Lists all users

### Global Options

* `--profile` - Path to profile YAML file containing configuration

### Examples

#### Add User

```sh
chissl admin adduser --username johndoe --password secret --admin --addresses ".*"
```

#### Delete User

```sh
chissl admin deluser --username johndoe
```

#### Get User

```sh
chissl admin getuser --username johndoe
```

To get the raw JSON response:

```sh
chissl admin getuser --username johndoe --raw
```

#### List Users

```sh
chissl admin listusers
```

To get the raw JSON response:

```sh
chissl admin listusers --raw
```

## Detailed Subcommands Usage

### `adduser`

Adds a new user to the chiSSL server.

```sh
chissl admin adduser [options]
```

#### Options

* `--username, -u` - Username for the new user
* `--password, -p` - Password for the new user
* `--addresses, -a` - Comma-separated list of regex expressions for allowed addresses
* `--admin` - Flag to grant admin permissions to the user

### `deluser`

Deletes an existing user from the chiSSL server.

```sh
chissl admin deluser [options]
```

#### Options

* `--username, -u` - Username of the user to delete

### `getuser`

Retrieves details of an existing user from the chiSSL server.

```sh
chissl admin getuser [options]
```

#### Options

* `--username, -u` - Username of the user to get
* `--raw` - Flag to output the raw JSON response

### `listusers`

Lists all users in the chiSSL server.

```sh
chissl admin listusers
```

#### Flags

* `--raw` - Flag to output the raw JSON response

