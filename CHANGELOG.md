# Changelog

All notable changes to this project will be documented in this file.

## 1.0.0 - 13th April 2021

### Added

* Create admin account at db initialization
* Signup User
* Verify user email-id and send code again feature
* Ldap Integration for security
* Login User with and without workspace
* Update and view profile
* JWT token based authentication
* Change, forgot and reset password
* Workspace creation, check status and availability feature for user
* Admin can delete workspace, list workspace and delete account of user
* Workspace login, forgot workspace feature
* Invite member in a Workspace
* Member singup and set password for new members in a Workspace
* Mail service is enabled to send mails using smtp server
* Configuration with TOML file and environment variables
* Set configuration with command-line flags
* Session Management using Redis Cache
* Logout API to delete Session
* Refresh Token API to generate new session and delete old Session
* Proper code commenting according to go-vet
* Mail data(template to use, subject, images to be passed) with toml file
* Endpoint configuration in different file
* Multiple Support Agent and admin account can be created, listed and deleted by admin only.
* Health check logic(checking all services and databases health)
* README contains all environment variables to be set
* CHANGELOG.md file contains all changes in this project with the release
* Compute Instance Creation Request and Request send to the support agent
* Static VPN Connection files Management system
* VPN Connection Username and password manage using VPN.TOML files
