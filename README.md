# Gostrufig

### Struct driven config management for your Go application!

As a developer, how many times have you been at the mercy of an underlying configuration library to populate configuration data stored in a struct? If you're anything like me, too many to count. That's why there's now Gostrufig!  Using a few simple tags and the Gostrufig library, you can manage all of your configuration items right in the struct! Plus, as an added bonus, there's an included persistence layer that will pull data from Etcd.  Want an example? Here you go:

```Go
type MyConfigInfo struct {
		DecodeDir   string `cfg-def:"/home/user/decoder"`
		Environment string `cfg-ns:"true" cfg-def:"developer"`
		Timer       int
		Type        string
		TestTimeout float64
}
ns := MyConfigInfo{}
gostrufig := GetGoStruFig("appname", "etcd", "http://localhost:2379")
gostrufig.RetrieveConfig(&ns)
```
## Struct Tags And Their Function

Does this wet your appetite? Good. Let me fill you in on those sweet sweet tags!

### cfg-def

This dandy little tag will allow you to set a default value for that field.  Don't worry, Gostrufig handles all the reflection to figure out how to turn that string into the value your struct desires.  ***Set it and forget it!***

### cfg-ns

This guy might be a bit confusing at first but, in the age of microservices and centralized configs, it'll make more sense.  If your app runs in multiple environments but shares a common config, this namespace tag will construct a unique path for that value. For the example above, the configuration would be located at /appname/developer in Etcd. If Environment was set to "production" instead, it would be /appname/production.

## Environmental Awareness

To go along with setting a default value in the struct itself, Gostrufig will look at the shell environment for any configuration items.  If you wish to override the Environment field above in the shell, just go ahead and set the value for **APPNAME_ENVIRONMENT**.  A good example would be:

```bash
export APPNAME_ENVIRONMENT=production
```

Gostrufig will automagically look for the capitalized names combining the application name and the field name. ***It's just that easy!***

## TODO
This is just the start.  Coming soon will be:

* Allowing a program to commit a config back to the persistence layer (you know, if your app has never run against an Etcd server, save your config to Etcd on first run)
* A YAML and/or JSON persistence layer
* Notification when a configuration changes and allow an app to update the config without a restart