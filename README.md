# Drone build Surge Protector
Plugin for Drone that kills old builds when new ones are started.

## Why?
Do you have some expensive builds? (For example, large/slow test suites)

Do you have some jerk that pushes a bunch of commits to their PR branch while working
without caring about the results until they are done? (For example, Brian)

Do you have an overworked drone cluster that occasionally can't keep up with the barrage described above?

Try this maybe?

## What?
When this plugin runs, it will look for any currently running builds for the current repo that
matches the current builds properties and kills them.

Specifically:
* If the build is a `pull_request` build, it cancels any other older PR builds for the same PR
* If the build is a `push` build, it cancels any other older push builds for the same branch
* If the build is a `tag` build, it cancels any other older tag builds for the same tag
* If the build is a `deployment` build, it cancels any other older deployment builds for the same deployment target

## How?
```
  kill_previous_builds:
    image: dankirberger/build-surge-protector
    secrets: [drone_token]
    drone_host: https://drone6.target.com
    when:
    # You might only want to run this on certain events, for example, 
    # you might not want to cancel an in progress deployment
      event: [push, deployment, tag, pull_request]
```

Where `drone_token` is a secret containing your token from here: https://drone6.target.com/account/token

Use a token for a NUID account.

This plugin was only tested against Drone 0.8.x