# BoxMate Export v0

Hacky scripts to get lift rep max data out of BoxMate for your team.

## Usage

1. Clone this repository
2. Add your BoxMate username and password to you environment variables
   - `export BOXMATE_USERNAME="your_username"`
   - `export BOXMATE_PASSWORD="your_password"`
3. Get your auth token
   - `go run ./v0/auth`
   - This will print out your auth token, copy it
   - `export BOXMATE_AUTH_TOKEN="your_auth_token"`
3. Run the download script to get the data
   - `go run ./v0/download`
4. Run the lift script to convert the data to an XLSX file
   - `go run ./v0/lifts`
