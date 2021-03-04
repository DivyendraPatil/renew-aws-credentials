# Renew Aws Credentials

There are millions of people who use aws credentials stored in `~/.aws/credentials` file. It seems mundane that we spend 3-5 minutes of our lives clicking around to refresh our credentials every 60 / 90 / 180 days (Depends on different organizations). This is also a security risk due to compromised credentials.

Your aws credentials need be located in `~/.aws/credentials` in the format:
```
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPLE
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

### Note
To use this executable, you need to have generated and placed your aws credentials at least once.

## Installation

Pull the repository. Enter it and run command:

```bash
go build -o binary
```
An executable file called `renew-aws-credentials` will be generated in the `binary` directory .
<br> 
A prebuilt binary is already stored in the same directory.

## Usage

```
./renew-aws-credentials
```

The script creates new credentials and updates them in the default location in your home directory.  <br /> 
If there are two keys, the old key is deleted. 

You can store the binary and set it as a cron on your system to automatically renew credentials after every few weeks.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
