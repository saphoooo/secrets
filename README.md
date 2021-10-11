# Go Secrets Manager

Go Secret Manager is a tool that allows you to use encrypted secrets in Datadog agent:

[Secrets Management](https://docs.datadoghq.com/agent/guide/secrets-management/?tab=linux)

## AWS

> Your EC2 instances need to be in the same region than the secrets you’re looking for, otherwise copy the secrets to the appropriate region.

### Configure IAM

EC2 instances targeted to use this feature first need a role to be able to read from the Secrets Manager. Let’s start with a simple role:

![fig1](https://a.cl.ly/RBuEKxqj)

And its associated policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "secretsmanager:GetSecretValue",
            "Resource": "*"
        }
    ]
}
```

You need to associate this role to the EC2 instance:

![fig2](https://a.cl.ly/E0uK21rd)

### Setup

> The source code is currently in a private repo on Github, but should be available soon. Build pipeline will also let you directly download the binary file for Linux and Windows.

Download the executable and set the required file permissions:

```sh
sudo chmod 700 <path_to_executable>
sudo chown dd-agent <path_to_exectuable>
```

### Store a new secret

Go to AWS Secrets Manager and store a new secret:

![fig3](https://a.cl.ly/jkuEgbq7)

Choose key/value pair:

![fig4](https://a.cl.ly/z8u1AeqO)

The key MUST be **dd-secret** in order to let the tool retrieve it. Enter the secret value and go to the next step.

![fig5](https://a.cl.ly/YEuPGndB)

Give a name to your secret and an optional description, then go to the next step.

![fig6](https://a.cl.ly/YEuPGndB)

Just click on next.

There is nothing to do for the step 4, just click **store** to store your new secret:

![fig7](https://a.cl.ly/GGu4PGkv)

## Try it out!

### On Linux

Verify your setup:

![fig8](https://a.cl.ly/P8u80AKA)

Try to use it manually:

![fig9](https://a.cl.ly/yAu0X5NR)

```bash
sudo su dd-agent -s /bin/bash -c "echo '{\"version\": \"1.0\", \"secrets\": [\"my-secret-env\"]}'|/usr/local/bin/datadog-secrets-aws"
{"my-secret-env":{"value":"test","error":""}}
```

Use it with the agent:

```
sudo vi /etc/datadog-agent/datadog.yaml

## @param env - string - optional
## @env DD_ENV - string - optional
## The environment name where the agent is running. Attached in-app to every
## metric, event, log, trace, and service check emitted by this Agent.
#
env: "ENC[my-secret-env]"

[...]

## @param secret_backend_command - string - optional
## `secret_backend_command` is the path to the script to execute to fetch secrets.
## The executable must have specific rights that differ on Windows and Linux.
##
## For more information see: https://github.com/DataDog/datadog-agent/blob/main/docs/agent/secrets.md
#
secret_backend_command: /usr/local/bin/datadog-secrets-aws
```

Restart the agent

```bash
sudo systemctl restart datadog-agent
```

Check agent status:

![fig10](https://a.cl.ly/RBuE65EO)

### On Windows

> You need to rely more on the CLI than the WebUI, because when the agent doesn’t start, you have to figure out why with the log file and the agent output.

Copy the .exe on the local filesystem.

Change the configuration of the Datadog agent:

```yaml
## @param env - string - optional
## @env DD_ENV - string - optional
## The environment name where the agent is running. Attached in-app to every
## metric, event, log, trace, and service check emitted by this Agent.
#
env: "ENC[my-secret-env]"
[...]
## @param secret_backend_command - string - optional
## `secret_backend_command` is the path to the script to execute to fetch secrets.
## The executable must have specific rights that differ on Windows and Linux.
##
## For more information see: https://github.com/DataDog/datadog-agent/blob/main/docs/agent/secrets.md
#
secret_backend_command: 'C:\Program Files\Datadog\Datadog Agent\bin\datadog-secrets-aws.exe'
```

Before restarting the agent, check that you’ve set the correct permissions on the .exe

Disable the inheritance and add SYSTEM and Administrator in Read and Execute:

![fig11](https://a.cl.ly/P8u8vZ2m)

Restart the agent:

![fig12](https://a.cl.ly/E0uKWGOj)

### Troubleshooting

If the agent doesn’t restart, try to run it manually:

```
"C:\Program Files\Datadog\Datadog Agent\bin\agent.exe start"
```

If you have an error that the agent can’t find the binary, check that you put the path in single quotes in the `datadog.yaml` file.

If the agent complains about permissions:

```
Error: unable to set up global agent configuration: unable to load Datadog 
config file: unable to decrypt secret from datadog.yaml: 
'S-1-5-21-1163014751-2161214003-1688321825-500' user is not allowed 
to execute secretBackendCommand 
'C:\Program Files\Datadog\Datadog Agent\bin\datadog-secrets-aws.exe'
```

And you have to figure out which user is S-1-5-21-1163014751-2161214003-1688321825-500, run the following command and then update your permissions accordingly:

```
wmic useraccount get name,sid
Name                SID
Administrator       S-1-5-21-1163014751-2161214003-1688321825-500
ddagentuser         S-1-5-21-1163014751-2161214003-1688321825-1008
DefaultAccount      S-1-5-21-1163014751-2161214003-1688321825-503
Guest               S-1-5-21-1163014751-2161214003-1688321825-501
WDAGUtilityAccount  S-1-5-21-1163014751-2161214003-1688321825-504
```