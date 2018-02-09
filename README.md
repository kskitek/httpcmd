# httpcmd

## Usage:

```bash
> go get gitlab.com/kskitek/httpcmd
> httpcmd
  -port string
    http port (default "8080")
  -scriptPath string
    path to the script (default "UNDEFINED")
```

This exposes three endpoints:

* `/start` - runns script asynchronously and returns id of the job
*  `/startAndWait` - runns the script and waits for it to end; returns job description
*  `/status/{id}` - returns job desciption

## Logging

`httpcmd` saves standard output of the scripts into log files under {id}.log.

## Example return job description:

```json
{
    Started: "2018-02-10T00:41:30.3129+01:00",
    Ended: "2018-02-10T00:41:35.335633+01:00",
    IsRunning: false,
    ErrorString: "",
    Err: null
}

{
    Started: "2018-02-10T00:44:17.025755+01:00",
    Ended: "2018-02-10T00:44:22.06281+01:00",
    IsRunning: false,
    ErrorString: "exit status 127",
    Err: {
        Stderr: null
    }
}
```