%title: okta-cli: A command line tool for Okta
%author: joel.franusic
%date: 2017-03-24

-> *Team:* <-

* Brian Retterer
* Evan Klein
* Joël Franusic

-------------------------------------------------


-> *Problem:* <-

It takes too long to get started with Okta:

* Create a Developer Edition account
* Log in to Okta and change your password
* Create a sample user
* Create an OIDC app in Okta
    * Name the app
    * Configure a URL
    * Assign the app to "Everyone"
* Configure CORS
* Download a sample app
* Configure the sample app
    * Figure out what the URL should be
    * Find the Client ID and Client Secret for the app

-------------------------------------------------

-> *Solution:* <-

Use the power of your operating system's teletype emulator!

* Create a Developer Edition account
* Log in to Okta and change your password
* Generate an API token
* Download Okta CLI

    $ okta init
    $ okta sample

-------------------------------------------------

-> *Demo:* <-

Already done:
* Created a fresh Okta org
* Changed password for the admin

What I'll be showing:
* That the Okta org is empty
* `$ okta init`
* Creating an API token
* `$ okta sample`
* Run the sample
* Log in with `paul.cook` and `L0rdn1k0n`

-------------------------------------------------

-> *Thank You!* <-
