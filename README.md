# bot-pages

`bot` is a simple command-line utility for looking up command usage examples. You can use ChatGPT for free, but it's really verbose. Like, so many words. Ugh. You may as well read the man page at that point. 

**This tool requires an [OpenAI API key](https://platform.openai.com/account/api-keys).**

## Caveat Emptor

This project uses OpenAI's GPT3 model for completion. The commands that it suggests may be wrong or harmful or annoyingly non-existent.

## Build and Install

```bash
git clone https://github.com/jonbrenner/bot-pages && cd bot-pages
make
sudo make install
```

or

```bash
git clone https://github.com/jonbrenner/bot-pages && cd bot-pages
go build -o bot
```

## Configuration

The `bot` command will prompt you to enter your API key on first use and store it in `$HOME/.bot-pages`:

```json
{
  "api-key": "YOUR_API_KEY_HERE"
}
```

## Usage

Specifying a single command to get some general usage examples:
```
❯ bot nc

$ nc -l -p "$PORT"
Listen for incoming connections on port $PORT.

$ nc -z "$HOST" "$PORT"
Check if the host $HOST is listening on port $PORT.

$ nc -l -p "$PORT" | xargs "$COMMAND"
Listen for incoming connections on port $PORT and pass them as arguments to $COMMAND.

[REFERENCES]
- https://linux.die.net/man/1/nc
```

And being a little more specific:
```
❯ bot socat reverse shell using tls

$ socat -v OPENSSL-LISTEN:$PORT,reuseaddr,cert=$CERT,cafile=$CAFILE,verify=0 EXEC:"$COMMAND"
Create a TLS-secured reverse shell using the certificate $CERT and the CA file $CAFILE, listening on port $PORT, and executing the command $COMMAND.

[REFERENCES]
- https://linux.die.net/man/1/socat
```
