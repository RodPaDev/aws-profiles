# AWS Profiles

CLI to manage aws profiles by using `default` as the currently selected profiles.

## Why?

Many projects I work with use [Terraform](https://www.terraform.io/) and since many of the projects use different set of credentials, I find myself having to manually switch the profiles in `~/aws/credentials` and `~/aws/config`.

I'm aware I can use named profiles and specify which profile I want to use but I simply don't want to because **I'm stubborn and I want it to work the way I want it to work**.

## Installation


1. Go to release and download the zip matching your os & architecture
2. Unzip to `/usr/local/bin` or if you're on windows `C:\Program Files\aws-profiles`
3. If you're on windows you need to [add the program to path](https://gist.github.com/RodPaDev/e9365fbb6f0c5d8553ceb84ad87110b1)


>I intend on publishing these binaries on `homebrew`, `chocolatey`, & need to figure what the best option if for `linux`.

## Usage


To select a profile run:

```
$ aws-profiles select <my-project-name>
```

To display help and available commands run:

```
$ aws-profiles
```

This is the list of commands and their description:

```
list                   List all profiles
select <profile-name>  Select a profile to select as the default aws profile
current                Display the current Default Profile's name
choose                 Choose a profile from an interactive shell
add                    Add a new aws profile
delete <profile-name>  Delete a aws profile
edit                   Edit an existing aws profile
help [command]         display help for command
```

## How does this work?

`aws-cli` uses ini files to store the credentials and configuration and uses named sections to group key-value pairs that belong to a profile.

So I use the `Default` profile (initally configured by running `aws configure`), as the current/active profile and keep the credentials of named profiles in seperate sections.

Example:

```ini
[default]
aws_access_key_id=foo_id
aws_secret_access_key=foo_secret

[foo]
aws_access_key_id=foo_id
aws_secret_access_key=foo_secret

[foo]
aws_access_key_id=bar_id
aws_secret_access_key=bar_secret
```

As you can see above, the values of the profile foo are used in default as well, this means the foo profile is active and that's how you switch profiles by essentially setting the values of the intended profile to default.
