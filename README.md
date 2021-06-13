# Passport

Passport is a Command-Line tool, used to manage scripts and provide an easy way to execute long/complex commands - with further support for managing secrets. Using the concept of workspaces, a workspace represents a directory, along with a number of scripts which can be executed in the given directory.

## Example

As an example, I have a Dockerfile which I need to build. This Dockerfile has a long list of, potentially secret, build arguments that need to be passed to it. Passport can be used to store and execute the long command with ease - this also saves having to type the command over and over again!

Let's say my command looks like this:

```
# C:/MyApp
$ docker build -t MyApp \
    --build-arg "MyName=Reece" \
    --build-arg "PORT=80" \
    --build-arg "MySecret=password123" \
    --file dev.Dockerfile \
    .
```

This can be simplified, by using Passport. Firstly, create a secret for `MySecret` build arg.

```
# C:/MyApp
$ passport secrets add --name "MySecret" --value "password123"
```

This will create a secret value, which will be encrypted at rest. Optionally, it can be stored in plain text by adding the `--plant-text` flag at the end.

Next, we need to create a new script in Passport. When creating a script for the first time inside a new directory, a workspace will be created automatically. But to do this, run:

```
# C:/MyApp
$ passport scripts add --name "Build" --command \
    "docker build -t MyApp \
        --build-arg \"MyName=reece\" \
        --build-arg \"PORT=80\" \
        --build-arg \"MySecret=<secret.MySecret>\" \
        --file dev.Dockerfile ."
```

After running this, the script `Build` will be added to our workspace. Note how in the `MySecret` build arg, the value is `<secret.MySecret>` - secrets can be interpolated into commands like so.

Now to run the script:

```
# C:/MyApp
$ passport run Build
```

The above command simpily just executes the command configured in the previous step.
