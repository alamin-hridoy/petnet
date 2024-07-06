How to test the templates
===

* get SMTP password for your gmail
  * to get SMTP password for gmail: if you have 2FA: 
    generate App password for your gmail,
    * https://support.google.com/accounts/answer/185833?hl=en
  * or if you don't - enable "lesser secure apps" and use normal password
    * https://support.google.com/accounts/answer/6010255?hl=en

* install tools
  * cd apix/tools
  * ./install

* go to apix/email
* `go run cli/main.go your-e-mail@gmail.com your-gmail-password where-to-send@foo.com`

* if you change any template - you need to run `go generate` in `email`
  before running `go run cli/main.go ...` again
