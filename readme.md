# CLOTP: Command-line one-time password generator

## Usage

**WARNING**
<span style="font-size:12pt;" ><span style="color:red;">Actually using this utility could result in security breach. For real.</span> Storing MFA codes on the same device where you might be using it for multi-factor authentication purposes is a bad idea that goes against the purpose of the MFA. Anyone with access to the computer will be able to steal your one-time passwords and defeat MFA. Granted, you're probably using your MFA-protected sites from your phone copying the codes from the Google Authenticator installed on the same phone. Just be aware of the dangers, you've been warned. At least use strong passphrase for the clotp keyring and do not reuse that passphrase anywhere else.</span>


    clotp command [command-options]

Commands:
 
  - `add` - add an OTP
  - `remove` - remove existing OTP
  - `list` - list all stored OTPs
  - `view` - view OTP code details
  - `code` - generate OTP code
  - `decode` - decode "otpauth-migration" URI
  - `set-counter` - set HOTP counter value
  - `scan` - decode QR code from the image

#### Add

`add` command adds new one-time password. OTP information is stored in the system's keyring.

    clotp add [options] <url>

where `url` is [otpauth](https://github.com/google/google-authenticator/wiki/Key-Uri-Format) or Google Authenticator [otpauth-migration](https://github.com/google/google-authenticator-android/issues/118) URL.

CLOTP supports HOTP and TOTP one-time passwords using SHA1, SHA256, and SHA512 hashes with number of code digits from 6 to 10. TOTP can have any time window defined.

Following options are supported:

  - `--name` - custom name that will override account label from the `otpauth` URL. `--name` option is ignored when used with `otpauth-migration` URIs

#### Remove

`remove` removes previously added one-time password

    clotp remove <name>

Where `name` is an OTP name.

#### List

`list` command show all one-time passwords added previously with the `add` command.

#### View

`view` command shows detailed information about a selected OTP code, including Hash algorithm, number of digits, time step, counter, and (**WARNING!!!**), secret key

    clotp view <name>

Example:

```
$ ./clotp view VPN

        Name: VPN
Account Name: uaraven
      Issuer: SomeVPNProvider
    OTP Type: TOTP
   Hash Type: SHA1
 Code Digits: 6
 Time offset: 0
   Time step: 30
      Secret: XXXXXXX
    Auth URI: otpauth://totp/uaraven@SomeVPNProvider?secret=XXXXXXX
```


#### Code

`code` will generate a new code from either otp id or otp name

    clotp code [options] <name>

Following options are supported:
  
  - `--counter` - set counter for HOTP code. The counter will be stored for future uses. This option is ignored if the name refers to Time-based OTP
  - `-c` or `--copy` - copies generated code to the clipboard

#### Decode

`decode` can be used to decode Google Authenticator "export" URI. `decode` command accepts `otpauth-migration` URI and prints out all `otpauth` urls encoded in that URI.

    clotp decode <url>

For example,

    $ clotp decode otpauth-migration://offline?data=XXXXX

    otpauth://totp/aaa:bbb?secret=YYYYY
    otpauth://hotp/aaa:ccc?secret=ZZZZZ&counter=1023

#### Set-Counter

`set-counter` allows to set the HTOP counter value

    clotp set-counter <name> <counter-value>


#### Scan

`scan` command can be used to decode QR code from the image. It supports `otpauth` and `otpauth-migration` URIs.

    clotp scan [--decode [--parse]] <image>

`<image>` is a path to a PNG or JPG image file.