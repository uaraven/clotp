# CLOTP: Command-line one-time password generator

## Usage

**WARNING**
<span style="font-size:12pt;" ><span style="color:red;">Actually using this utility could result in security breach. For real.</span> It is just plainly not a sound security practice, to have your one-time password stored on the same device where you might be using it for multi-factor authentication purposes. Anyone with access to the computer will be able to steal your one-time passwords and defeat MFA. Granted, you're probably using your MFA-protected sites from your phone copying the codes from the Google Authenticator installed on the same phone. Just be aware of the dangers, you've been warned.</span>


    clotp [options] command [command-options]


Options:

 - None at the moment

Commands:
 
  - `list` - list all stored OTPs
  - `add` - add an OTP
  - `remove` - remove existing OTP
  - `code` - generate OTP code
  - `decode` - decode "otpauth-migrate" URI

#### List

`list` command show all one-time passwords added previously with the `add` command.

#### Add

`add` command adds new one-time password. OTP information is stored in the system's keyring.

    clotp add [options] <url>

where `url` is [otpauth](https://github.com/google/google-authenticator/wiki/Key-Uri-Format) or Google Authenticator [otpauth-migration](https://github.com/google/google-authenticator-android/issues/118) URL.

CLOTP supports HOTP and TOTP  one-time passwords using SHA1, SHA256, and SHA512 hashes with number of code digits from 6 to 10. TOTP can have any time window defined.

Following options are supported:

  - `--name` - custom name that will override Issuer and account information from the `otpauth` URL.

#### Remove

`remove` removes previously added one-time password

    clotp remove <name>

Where `name` is an OTP name.

#### Code

`code` will generate a new code

    clotp code <name>

